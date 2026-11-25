package generic

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

// Timing of deleteEntityRaw's verification (variables to allow for shorter values in tests)
var (
	deleteVerifyPollInterval   = 5 * time.Second
	deleteVerifyResendInterval = 20 * time.Second
	deleteVerifyNotFoundReads  = 3
)

// deleteEntityRaw deletes the entity of a resource. Without WriteOptions.DeleteRetryNotFoundFor,
// WriteOptions.DeleteVerifyFor and WriteOptions.WriteRetryFunc it is the same as DeleteRaw.
//
// Some entities (e.g. Conditional Access policies) are not consistent in MS Graph for a while after they have been
// created: a DELETE might return 404 (or a 403 "CheckAccess returned 'NotAllowed'") although the entity exists, or
// return 204 without the entity actually being deleted. Accepting such a 404 as "already gone" leaves an orphaned
// entity (see issue #241).
//   - DeleteRetryNotFoundFor: retry such 404/403 responses for this long, before accepting a 404 as "already gone".
//   - DeleteVerifyFor: after a successful DELETE, read the entity until MS Graph reports it as missing several times
//     in a row, sending the DELETE again if it still exists after a while. Fails if it still exists after this long.
func (aps *AccessParams) deleteEntityRaw(ctx context.Context, diags *diag.Diagnostics,
	baseUri string, id string, idAttributer GetAttributer) {

	if aps.WriteOptions.DeleteRetryNotFoundFor <= 0 && aps.WriteOptions.DeleteVerifyFor <= 0 && aps.writeRetryFunc() == nil {
		aps.DeleteRaw(ctx, diags, baseUri, id, idAttributer)
		return
	}

	uri := aps.GetUriWithIdForUD(ctx, diags, baseUri, id, idAttributer)
	if diags.HasError() {
		return
	}

	notFound, err := aps.deleteRetryingInconsistency(ctx, uri)
	if err != nil {
		diags.AddError("Error deleting from MS Graph", fmt.Sprintf("Original Error: %s", err.Error()))
		return
	}

	existed := false
	if aps.WriteOptions.DeleteVerifyFor > 0 {
		// also verify after a 404: MS Graph might return 404 for a DELETE of an entity that has just been created for
		// longer than DeleteRetryNotFoundFor, while the entity still exists
		existed = aps.verifyDeleted(ctx, diags, uri)
		if diags.HasError() {
			return
		}
	}
	if notFound && !existed {
		diags.AddWarning("Unable to delete from MS Graph", fmt.Sprintf("Item %q does not exist (anymore).", uri.Entity))
	}
}

// deleteRetryingInconsistency sends a DELETE for uri, retrying (for WriteOptions.DeleteRetryNotFoundFor) responses that
// MS Graph returns for entities that have just been created, and (for WriteOptions.WriteRetryFor) the responses
// WriteOptions.WriteRetryFunc accepts. It returns whether the entity was (still) not found.
func (aps *AccessParams) deleteRetryingInconsistency(ctx context.Context, uri msgraph.Uri) (bool, error) {
	var notFoundDeadline time.Time
	if d := aps.WriteOptions.DeleteRetryNotFoundFor; d > 0 {
		notFoundDeadline = time.Now().Add(d)
	}
	writeRetryFunc := aps.writeRetryFunc()
	consistencyFailureFunc := func(resp *http.Response, o *odata.OData) bool {
		return (time.Now().Before(notFoundDeadline) && isDeleteInconsistencyResponse(resp, o)) ||
			(writeRetryFunc != nil && writeRetryFunc(resp, o))
	}
	_, status, _, err := aps.graphClient.Delete(ctx, msgraph.DeleteHttpRequestInput{
		Uri:                    uri,
		ValidStatusCodes:       []int{http.StatusOK, http.StatusNoContent},
		ConsistencyFailureFunc: consistencyFailureFunc,
	})
	if err != nil && status == http.StatusNotFound {
		return true, nil
	}
	return false, err
}

// isDeleteInconsistencyResponse reports whether resp is a 404, or a 403 that MS Graph returns for a short while for
// entities that have just been created (although the caller is allowed to delete them).
func isDeleteInconsistencyResponse(resp *http.Response, o *odata.OData) bool {
	if resp == nil {
		return false
	}
	switch resp.StatusCode {
	case http.StatusNotFound:
		return true
	case http.StatusForbidden:
		return o != nil && o.Error != nil && strings.Contains(o.Error.String(), "CheckAccess returned 'NotAllowed'")
	}
	return false
}

// verifyDeleted reads uri until MS Graph returns 404 for deleteVerifyNotFoundReads reads in a row. If the entity still
// exists deleteVerifyResendInterval after the last DELETE, the DELETE is sent again. Adds an error if the entity
// still exists after WriteOptions.DeleteVerifyFor. Returns whether MS Graph returned the entity at all.
func (aps *AccessParams) verifyDeleted(ctx context.Context, diags *diag.Diagnostics, uri msgraph.Uri) (existed bool) {
	deadline := time.Now().Add(aps.WriteOptions.DeleteVerifyFor)
	lastDelete := time.Now()
	notFoundReads := 0
	for {
		_, status, _, err := aps.graphClient.Get(ctx, msgraph.GetHttpRequestInput{
			Uri:              uri,
			ValidStatusCodes: []int{http.StatusOK},
		})
		switch {
		case status == http.StatusNotFound:
			notFoundReads++
			if notFoundReads >= deleteVerifyNotFoundReads {
				return
			}
		case err == nil:
			existed = true
			notFoundReads = 0
			if time.Since(lastDelete) >= deleteVerifyResendInterval {
				tflog.Info(ctx, fmt.Sprintf("Item %q still exists after it has been deleted, sending DELETE again", uri.Entity))
				if _, err := aps.deleteRetryingInconsistency(ctx, uri); err != nil {
					diags.AddError("Error deleting from MS Graph", fmt.Sprintf("Original Error: %s", err.Error()))
					return
				}
				lastDelete = time.Now()
			}
		default:
			diags.AddError("Error verifying deletion in MS Graph", fmt.Sprintf("Original Error: %s", err.Error()))
			return
		}

		if time.Now().After(deadline) {
			diags.AddError("Item still exists in MS Graph after deleting it",
				fmt.Sprintf("Item %q has been deleted, but MS Graph still returned it %s later. Please check and delete it manually if needed.",
					uri.Entity, aps.WriteOptions.DeleteVerifyFor))
			return
		}

		select {
		case <-ctx.Done():
			diags.AddError("Error verifying deletion in MS Graph", ctx.Err().Error())
			return
		case <-time.After(deleteVerifyPollInterval):
		}
	}
}
