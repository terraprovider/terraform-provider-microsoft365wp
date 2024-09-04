package generic

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/manicminer/hamilton/msgraph"
)

func (aps *AccessParams) ReadSingleRaw(ctx context.Context, diags *diag.Diagnostics, schema tftypes.AttributePathStepper,
	readOptions ReadOptions, id string, idAttributer GetAttributer, tolerateNotFound bool) tftypes.Value {

	uri, odataFilter := aps.GetUriWithIdForR(ctx, diags, "", id, idAttributer, readOptions.SingleItemUseODataFilter, readOptions.ODataFilter)
	if diags.HasError() {
		return tftypes.Value{}
	}

	rawVal := aps.ReadRaw(ctx, diags, uri, readOptions.ODataExpand, odataFilter, nil, tolerateNotFound)
	if rawVal == nil {
		return tftypes.Value{}
	}

	for _, r := range readOptions.ExtraRequests {
		leafUri := msgraph.Uri{
			Entity: fmt.Sprintf("%s/%s", uri.Entity, r.UriSuffix),
		}
		rawLeafVal := aps.ReadRaw(ctx, diags, leafUri, r.ODataExpand, odataFilter, nil, tolerateNotFound)
		if diags.HasError() {
			return tftypes.Value{}
		}

		if rawLeafVal != nil {
			rawVal[r.ParentAttribute] = rawLeafVal["value"]
		}
	}

	for _, c := range readOptions.ExtraRequestsCustom {
		params := ReadExtraRequestCustomFuncParams{
			Ctx:              ctx,
			Diags:            diags,
			Client:           aps.graphClient,
			Uri:              uri,
			TolerateNotFound: tolerateNotFound,
			RawVal:           rawVal,
		}
		c(params)
		if diags.HasError() {
			return tftypes.Value{}
		}
	}

	tfVal := ConvertOdataRawToTerraform(ctx, diags, schema, rawVal, "", aps.GraphToTerraformMiddleware)
	if diags.HasError() {
		return tftypes.Value{}
	}

	return tfVal
}

func (aps *AccessParams) ReadRaw(ctx context.Context, diags *diag.Diagnostics,
	uri msgraph.Uri, odataExpand string, odataFilter string, odataSelect []string, tolerateNotFound bool) map[string]any {
	return ReadRaw(ctx, diags, aps.graphClient, uri, odataExpand, odataFilter, odataSelect, tolerateNotFound)
}

func ReadRaw(ctx context.Context, diags *diag.Diagnostics, graphClient *msgraph.Client,
	uri msgraph.Uri, odataExpand string, odataFilter string, odataSelect []string, tolerateNotFound bool) map[string]any {

	odataQuery := odata.Query{
		Filter: odataFilter,
		Select: odataSelect,
	}
	if odataExpand != "" {
		odataQuery.Expand = odata.Expand{Relationship: odataExpand}
	}
	graphResp, status, odata, err := graphClient.Get(ctx, msgraph.GetHttpRequestInput{
		OData:            odataQuery,
		ValidStatusCodes: []int{http.StatusOK},
		Uri:              uri,
	})
	if err != nil {
		if status == http.StatusNotFound || (odata != nil && odata.Error != nil && *odata.Error.Code == "ResourceNotFound") {
			if tolerateNotFound {
				diags.AddWarning(fmt.Sprintf("Item %q not found in MS Graph", uri.Entity),
					"Automatically removing from Terraform State instead of returning the error.")
			} else {
				diags.AddError(fmt.Sprintf("Item %q not found in MS Graph", uri.Entity),
					fmt.Sprintf("After attempting to read, MS Graph returned a resource not found error for the id provided. Original Error: %s", err.Error()))
			}
		} else {
			diags.AddError("Error reading from MS Graph",
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
