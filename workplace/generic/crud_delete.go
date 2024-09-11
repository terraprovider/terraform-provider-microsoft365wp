package generic

import (
	"context"
	"fmt"
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

	_, status, _, err := aps.graphClient.Delete(ctx, msgraph.DeleteHttpRequestInput{
		Uri:              uri,
		ValidStatusCodes: []int{http.StatusOK, http.StatusNoContent},
	})
	if err != nil {
		if status == http.StatusNotFound {
			diags.AddWarning("Unable to delete from MS Graph", fmt.Sprintf("Item %q does not exist (anymore).", uri.Entity))
		} else {
			diags.AddError("Error deleting from MS Graph", fmt.Sprintf("Original Error: %s", err.Error()))
		}
		return
	}
}
