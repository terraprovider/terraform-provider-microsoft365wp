package generic

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (aps *AccessParams) ReadSingleCompleteTf(ctx context.Context, diags *diag.Diagnostics, schema tftypes.AttributePathStepper,
	readOptions ReadOptions, id string, idAttributer GetAttributer, tolerateNotFound bool) tftypes.Value {
	return aps.ReadSingleCompleteTf2(ctx, diags, schema, readOptions, id, idAttributer, "", tolerateNotFound, nil, nil)
}

func (aps *AccessParams) ReadSingleCompleteTf2(ctx context.Context, diags *diag.Diagnostics, schema tftypes.AttributePathStepper,
	readOptions ReadOptions, id string, idAttributer GetAttributer, odataFilter string, tolerateNotFound bool,
	reqState *tfsdk.State, respPrivate PrivateDataGetSetter) tftypes.Value {
	return aps.ReadSingleCompleteTf3(ctx, diags, schema, readOptions, id, idAttributer, odataFilter, "", 0, tolerateNotFound, reqState, respPrivate)
}

func (aps *AccessParams) ReadSingleCompleteTf3(ctx context.Context, diags *diag.Diagnostics, schema tftypes.AttributePathStepper,
	readOptions ReadOptions, id string, idAttributer GetAttributer, odataFilter string, odataOrderby string, odataTop int,
	tolerateNotFound bool, reqState *tfsdk.State, respPrivate PrivateDataGetSetter) tftypes.Value {

	id = aps.GetId(ctx, diags, id, idAttributer) // we need the id anyway for GraphToTerraformMiddleware
	if diags.HasError() {
		return tftypes.Value{}
	}

	if odataFilter == "" {
		odataFilter = readOptions.ODataFilter
	}
	uri, odataFilter2 := aps.GetUriWithIdForR(ctx, diags, "", id, idAttributer, readOptions.SingleItemUseODataFilter, odataFilter)
	if diags.HasError() {
		return tftypes.Value{}
	}

	rawVal := aps.ReadRaw3(ctx, diags, uri, readOptions.ODataExpand, odataFilter2, readOptions.ODataSelect, odataOrderby, odataTop, tolerateNotFound)
	if rawVal == nil || diags.HasError() {
		return tftypes.Value{}
	}

	if len(readOptions.ExtraRequests) > 0 || len(readOptions.ExtraRequestsCustom) > 0 {

		// also see ConvertOdataRawToTerraform (but we need to do this up here for the extra requests to work)
		if items, ok := rawVal["value"].([]any); ok && len(rawVal) <= 3 {
			if len(items) == 1 {
				if singleValue, ok := items[0].(map[string]any); ok {
					rawVal = singleValue
				}
			}
		}

		// need actual id for next requests in case the main item had been found using odataFilter etc.
		if id == "" && odataFilter != "" {
			var ok bool
			if id, ok = rawVal[aps.EntityId.AttrNameGraph].(string); !ok {
				diags.AddError("Error reading from MS Graph", "Unable to determine entity id for extra requests")
				return tftypes.Value{}
			}
			uri.Entity = fmt.Sprintf("%s/%s", uri.Entity, id)
		}

		for _, r := range readOptions.ExtraRequests {
			uriSuffix := r.UriSuffix
			if uriSuffix == "" {
				uriSuffix = r.Attribute
			}
			leafUri := fmt.Sprintf("%s/%s", uri.Entity, uriSuffix)
			rawLeafVal := aps.ReadRaw(ctx, diags, leafUri, tolerateNotFound)
			if diags.HasError() {
				return tftypes.Value{}
			}

			if rawLeafVal != nil {
				if valueItems, ok := rawLeafVal["value"].([]any); ok {
					rawVal[r.Attribute] = valueItems
				} else {
					delete(rawLeafVal, "@odata.context") // ensure that we focus on the real content attributes
					if len(rawLeafVal) > 0 {
						rawVal[r.Attribute] = rawLeafVal
					} else {
						rawVal[r.Attribute] = nil
					}
				}
			}
		}

		for _, c := range readOptions.ExtraRequestsCustom {
			params := ReadExtraRequestCustomParams{
				Client:           aps.graphClient,
				Uri:              uri,
				TolerateNotFound: tolerateNotFound,
				ReqState:         reqState,
				RespPrivate:      respPrivate,
				RawVal:           rawVal,
			}
			c(ctx, diags, params)
			if diags.HasError() {
				return tftypes.Value{}
			}
		}
	}

	tfVal := ConvertOdataRawToTerraform(ctx, diags, schema, rawVal, "", aps.GraphToTerraformMiddleware, id, aps.GraphToTerraformMiddlewareTargetSetRunOnRawVal)
	if diags.HasError() {
		return tftypes.Value{}
	}

	// if query was by OData filter, we might only notice down here
	if (tfVal == tftypes.Value{}) && !tolerateNotFound {
		diags.AddError(fmt.Sprintf("No single entity was found in MS Graph for query %q and OData filter %q", uri.Entity, odataFilter2), "MS Graph returned none or multiple entities.")
		return tftypes.Value{}
	}

	return tfVal
}

func (aps *AccessParams) ReadRaw(ctx context.Context, diags *diag.Diagnostics, uriEntity string, tolerateNotFound bool) map[string]any {
	return aps.ReadRaw2(ctx, diags, msgraph.Uri{Entity: uriEntity}, "", "", nil, tolerateNotFound)
}

func (aps *AccessParams) ReadRaw2(ctx context.Context, diags *diag.Diagnostics, uri msgraph.Uri,
	odataExpand string, odataFilter string, odataSelect []string, tolerateNotFound bool) map[string]any {
	return aps.ReadRaw3(ctx, diags, uri, odataExpand, odataFilter, odataSelect, "", 0, tolerateNotFound)
}
func (aps *AccessParams) ReadRaw3(ctx context.Context, diags *diag.Diagnostics, uri msgraph.Uri,
	odataExpand string, odataFilter string, odataSelect []string, odataOrderby string, odataTop int, tolerateNotFound bool) map[string]any {
	odataQuery := odata.Query{
		Filter: odataFilter,
		Select: odataSelect,
		Top:    odataTop,
	}
	if odataExpand != "" {
		odataQuery.Expand = odata.Expand{Relationship: odataExpand}
	}
	if odataOrderby != "" {
		// do not use odata.Query.OrderBy here as it would be a lot more work to parse
		if uri.Params == nil {
			uri.Params = url.Values{}
		}
		uri.Params.Set("$orderby", odataOrderby)
	}
	return ReadRaw2(ctx, diags, aps.graphClient, uri, &odataQuery, aps.ReadOptions.ValidStatusCodesExtra, tolerateNotFound)
}

func (aps *AccessParams) ReadId(ctx context.Context, diags *diag.Diagnostics,
	baseUri string, parentIdAttributer GetAttributer, useOdataSelect bool, odataFilter string, idAssertionRegex *regexp.Regexp) (idResult string) {

	baseUriUri := aps.GetBaseUri(ctx, diags, baseUri, parentIdAttributer)
	if diags.HasError() {
		return
	}
	errorSummary := fmt.Sprintf("Error determining MS Graph id for %q", baseUriUri.Entity)

	odataSelect := []string{}
	if useOdataSelect {
		odataSelect = append(odataSelect, aps.EntityId.AttrNameGraph)
	}
	idResultRaw := aps.ReadRaw2(ctx, diags, baseUriUri, "", odataFilter, odataSelect, false)
	if diags.HasError() {
		return
	}

	if aps.GraphToTerraformMiddleware != nil {
		params := GraphToTerraformMiddlewareParams{RawVal: idResultRaw}
		if err := aps.GraphToTerraformMiddleware(ctx, diags, &params); err != nil {
			diags.AddError(errorSummary, fmt.Sprintf("GraphToTerraformMiddleware returned an error: %s", err.Error()))
		}
		if diags.HasError() {
			return
		}
	}

	if id, ok := idResultRaw[aps.EntityId.AttrNameGraph].(string); ok {
		idResult = id
	} else if valueSlice, ok := idResultRaw["value"].([]any); ok && len(valueSlice) == 1 {
		if id, ok := valueSlice[0].(map[string]any)[aps.EntityId.AttrNameGraph].(string); ok {
			idResult = id
		}
	}
	if idResult == "" {
		diags.AddError(errorSummary, "Azure query result did not have the expected format")
		return
	}

	if idAssertionRegex != nil {
		if !idAssertionRegex.MatchString(idResult) {
			diags.AddError(errorSummary, fmt.Sprintf("id returned from Azure query does not match assertion regex (id is %q, idAssertionRegex is %v)", idResult, idAssertionRegex))
			return
		}
	}

	return
}

func (aps *AccessParams) CheckEntityExistence(ctx context.Context, diags *diag.Diagnostics,
	id string, idAttributer GetAttributer) (entityExists bool) {

	uri, odataFilter2 := aps.GetUriWithIdForR(ctx, diags, "", id, idAttributer, aps.ReadOptions.SingleItemUseODataFilter, aps.ReadOptions.ODataFilter)
	if diags.HasError() {
		return
	}

	entityExists = aps.ReadRaw2(ctx, diags, uri, "", odataFilter2, nil, true) != nil
	return
}

func (aps *AccessParams) PopulateStateParentIdsFromRequest(ctx context.Context, diags *diag.Diagnostics, dst *tfsdk.State, src GetAttributer) {
	for _, pe := range aps.ParentEntities {
		err := CopyValueAtPath(ctx, dst, src, pe.ParentIdField)
		if err != nil {
			diags.AddError(fmt.Sprintf("CopyValueAtPath(), ParentIdField: %s", pe.ParentIdField), err.Error())
			return
		}
	}
}

func ReadRaw(ctx context.Context, diags *diag.Diagnostics, graphClient *msgraph.Client, uriEntity string, tolerateNotFound bool) map[string]any {
	return ReadRaw2(ctx, diags, graphClient, msgraph.Uri{Entity: uriEntity}, nil, nil, tolerateNotFound)
}

func ReadRaw2(ctx context.Context, diags *diag.Diagnostics, graphClient *msgraph.Client, uri msgraph.Uri, odataQuery *odata.Query,
	validStatusCodesExtra []int, tolerateNotFound bool) map[string]any {

	graphRequest := msgraph.GetHttpRequestInput{
		Uri:              uri,
		ValidStatusCodes: append([]int{http.StatusOK}, validStatusCodesExtra...),
	}
	if odataQuery != nil {
		graphRequest.OData = *odataQuery
	}

	graphResp, status, odata, err := graphClient.Get(ctx, graphRequest)
	if err != nil {
		if status == http.StatusNotFound || (odata != nil && odata.Error != nil && *odata.Error.Code == "ResourceNotFound") {
			if tolerateNotFound {
				tflog.Info(ctx, fmt.Sprintf("No entity found in MS Graph for query %q and OData filter %q", uri.Entity, odataQuery.Filter))
			} else {
				diags.AddError(fmt.Sprintf("No entity found in MS Graph for query %q and OData filter %q", uri.Entity, odataQuery.Filter),
					fmt.Sprintf("MS Graph returned a resource not found error. Original Error: %s", err.Error()))
			}
		} else {
			diags.AddError(fmt.Sprintf("Error reading from MS Graph for query %q and OData filter %q", uri.Entity, odataQuery.Filter),
				fmt.Sprintf("Original Error: %s", err.Error()))
		}
		return nil
	}

	defer graphResp.Body.Close()
	respBody, err := io.ReadAll(graphResp.Body)
	if err != nil {
		diags.AddError("io.ReadAll()", err.Error())
		return nil
	}

	rawVal := ConvertOdataJsonToRaw(ctx, diags, respBody)
	if diags.HasError() {
		return nil
	}

	return rawVal
}
