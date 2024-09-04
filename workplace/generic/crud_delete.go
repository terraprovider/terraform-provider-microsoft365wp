package generic

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/manicminer/hamilton/msgraph"
)

func (aps *AccessParams) DeleteRaw(ctx context.Context, diags *diag.Diagnostics,
	baseUri string, id string, idAttributer GetAttributer) {

	uri := aps.GetUriWithIdForUD(ctx, diags, baseUri, id, idAttributer)
	if diags.HasError() {
		return
	}

	_, _, _, err := aps.graphClient.Delete(ctx, msgraph.DeleteHttpRequestInput{
		Uri:              uri,
		ValidStatusCodes: []int{http.StatusOK, http.StatusNoContent},
	})
	if err != nil {
		diags.AddError("Unable to delete from MS Graph", err.Error())
		return
	}
}
