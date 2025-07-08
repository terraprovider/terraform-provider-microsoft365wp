package generic

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func (aps *AccessParams) CreateRaw(ctx context.Context, diags *diag.Diagnostics,
	baseUri string, parentIdAttributer GetAttributer, rawVal map[string]any) (id string, rawResult map[string]any) {
	return aps.CreateRaw2(ctx, diags, baseUri, parentIdAttributer, rawVal, aps.WriteOptions.UsePutForCreate, aps.WriteOptions.UsePutForCreate)
}

func (aps *AccessParams) CreateRaw2(ctx context.Context, diags *diag.Diagnostics,
	baseUri string, parentIdAttributer GetAttributer, rawVal map[string]any,
	allowEmptyResponse bool, usePut bool) (id string, rawResult map[string]any) {

	uri := aps.GetBaseUri(ctx, diags, baseUri, parentIdAttributer)
	if diags.HasError() {
		return
	}

	rawResult = CreateRaw(ctx, diags, aps.graphClient, uri, rawVal, nil, allowEmptyResponse, usePut)
	if diags.HasError() || (allowEmptyResponse && len(rawResult) == 0) { // if response is completely empty, also just return
		return
	}

	// Currently we only read the new item id from the response body (as navigation property values will not be returned from POST anyway)
	id, ok := rawResult[aps.EntityId.AttrNameGraph].(string)
	if !ok {
		diags.AddError("Unable to parse MS Graph result", fmt.Sprintf("property `%s` not found or not of type string", aps.EntityId.AttrNameGraph))
		return
	}

	return
}

func CreateRaw(ctx context.Context, diags *diag.Diagnostics, graphClient *msgraph.Client,
	uri msgraph.Uri, rawVal map[string]any, validStatusCodesExtra []int, allowEmptyResponse bool, usePut bool) (rawResult map[string]any) {

	jsonVal := ConvertOdataRawToJson(ctx, diags, rawVal, "Apply")
	if diags.HasError() {
		return
	}

	validStatusCodes := []int{http.StatusOK, http.StatusCreated}
	if allowEmptyResponse {
		validStatusCodes = append(validStatusCodes, http.StatusNoContent)
	}
	validStatusCodes = append(validStatusCodes, validStatusCodesExtra...)

	var graphResp *http.Response
	var odata *odata.OData
	var err error
	if !usePut {
		graphResp, _, odata, err = graphClient.Post(ctx, msgraph.PostHttpRequestInput{Uri: uri, Body: jsonVal, ValidStatusCodes: validStatusCodes})
	} else {
		graphResp, _, odata, err = graphClient.Put(ctx, msgraph.PutHttpRequestInput{Uri: uri, Body: jsonVal, ValidStatusCodes: validStatusCodes})
	}
	if err != nil {
		diags.AddError("Unable to create with MS Graph", err.Error())
		if odata != nil && odata.Error != nil && odata.Error.Code != nil && *odata.Error.Code == "BadRequest" {
			diags.AddError("HTTP Request Body", string(jsonVal))
		}
		return
	}

	defer graphResp.Body.Close()
	respBody, err := io.ReadAll(graphResp.Body)
	if err != nil {
		diags.AddError("io.ReadAll()", err.Error())
		return
	}

	// if response is completely empty, just return
	if allowEmptyResponse && len(respBody) == 0 {
		return
	}

	rawResult = ConvertOdataJsonToRaw(ctx, diags, respBody)
	if diags.HasError() {
		return
	}

	return
}
