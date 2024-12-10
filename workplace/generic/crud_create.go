package generic

import (
	"context"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/manicminer/hamilton/msgraph"
)

func (aps *AccessParams) CreateRaw(ctx context.Context, diags *diag.Diagnostics,
	baseUri string, parentIdAttributer GetAttributer, rawVal map[string]any,
	allowEmptyResponse bool) (id string, rawResult map[string]any) {

	uri := aps.GetBaseUri(ctx, diags, baseUri, parentIdAttributer)
	if diags.HasError() {
		return
	}

	jsonVal := ConvertOdataRawToJson(ctx, diags, rawVal, "Plan")
	if diags.HasError() {
		return
	}

	validStatusCodes := []int{http.StatusOK, http.StatusCreated}
	if allowEmptyResponse {
		validStatusCodes = append(validStatusCodes, http.StatusNoContent)
	}

	graphResp, _, odata, err := aps.graphClient.Post(ctx, msgraph.PostHttpRequestInput{
		Uri:              uri,
		Body:             jsonVal,
		ValidStatusCodes: validStatusCodes,
	})
	if err != nil {
		diags.AddError("Unable to create with MS Graph", err.Error())
		if odata != nil && odata.Error != nil && odata.Error.Code != nil && *odata.Error.Code == "BadRequest" {
			diags.AddError("HTTP Request Body", string(jsonVal))
		}
		return
	}

	// Currently we only read the new item id from the response body (as navigation property values will not be returned from POST anyway)
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

	id, ok := rawResult["id"].(string)
	if !ok {
		diags.AddError("Unable to parse MS Graph result", "property `id` not found or not of type string")
		return
	}

	return
}
