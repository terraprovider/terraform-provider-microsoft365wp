package generic

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/manicminer/hamilton/msgraph"
)

func (aps *AccessParams) UpdateRaw(ctx context.Context, diags *diag.Diagnostics,
	baseUri string, id string, idAttributer GetAttributer, rawVal map[string]any, usePut bool) {

	uri := aps.GetUriWithIdForUD(ctx, diags, baseUri, id, idAttributer)
	if diags.HasError() {
		return
	}

	jsonVal := ConvertOdataRawToJson(ctx, diags, rawVal, "Plan")
	if diags.HasError() {
		return
	}

	var err error
	// No need to read body when updating as id will stay the same
	validStatusCodes := []int{http.StatusOK, http.StatusNoContent}
	if !usePut {
		_, _, _, err = aps.graphClient.Patch(ctx, msgraph.PatchHttpRequestInput{Uri: uri, Body: jsonVal, ValidStatusCodes: validStatusCodes})
	} else {
		_, _, _, err = aps.graphClient.Put(ctx, msgraph.PutHttpRequestInput{Uri: uri, Body: jsonVal, ValidStatusCodes: validStatusCodes})
	}
	if err != nil {
		diags.AddError("Unable to update with MS Graph", err.Error())
		return
	}
}
