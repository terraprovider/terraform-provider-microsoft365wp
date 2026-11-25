package generic

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type PatchAndVerifyOptions struct {
	PollInterval     time.Duration // between reads
	ResendAfter      time.Duration // send the PATCH again if the written values have not shown up after this long
	ConsecutiveReads int           // reads may be served by different replicas, so require several consecutive matching reads
	Timeout          time.Duration // give up (with an error) after this long
}

var DefaultPatchAndVerifyOptions = PatchAndVerifyOptions{
	PollInterval:     10 * time.Second,
	ResendAfter:      45 * time.Second,
	ConsecutiveReads: 2,
	Timeout:          4 * time.Minute,
}

// PatchAndVerifyRaw PATCHes body to uri and then reads readUri until matches returns true for
// opts.ConsecutiveReads reads in a row.
//
// It is meant for entities where MS Graph acknowledges a write (204) that never takes effect (observed for the
// authenticationContextId of unifiedRbacResourceAction, see issue #239): if the written values have not shown up
// opts.ResendAfter after the last PATCH, the PATCH is sent again. A 404 for the PATCH is not taken as final either (MS
// Graph returned it for clearing a value whose write had been lost), as long as readUri can still be read. Adds an
// error if the values still have not shown up after opts.Timeout, since the write cannot be considered successful
// then (unlike WaitForRawValues, which only adds a warning for stale reads of a successful write).
func PatchAndVerifyRaw(ctx context.Context, diags *diag.Diagnostics, client *msgraph.Client, uri msgraph.Uri,
	body map[string]any, readUri msgraph.Uri, matches func(got map[string]any) bool, opts PatchAndVerifyOptions) {

	jsonBody := ConvertOdataRawToJson(ctx, diags, body, "Plan")
	if diags.HasError() {
		return
	}

	// patch returns false if the PATCH failed for good (error added to diags)
	patch := func() bool {
		_, status, _, err := client.Patch(ctx, msgraph.PatchHttpRequestInput{
			Uri:              uri,
			Body:             jsonBody,
			ValidStatusCodes: []int{http.StatusOK, http.StatusNoContent},
		})
		if err != nil {
			if status == http.StatusNotFound {
				tflog.Info(ctx, fmt.Sprintf("PATCH of %q returned 404, checking the values by reading: %s", uri.Entity, err))
				return true
			}
			diags.AddError("Unable to update with MS Graph", err.Error())
			return false
		}
		return true
	}

	if !patch() {
		return
	}
	start := time.Now()
	lastPatch := start
	matchingReads := 0
	for {
		readDiags := diag.Diagnostics{}
		got := ReadRaw2(ctx, &readDiags, client, readUri, nil, nil, true)
		switch {
		case readDiags.HasError():
			diags.Append(readDiags...)
			return
		case got == nil:
			diags.AddError("Unable to update with MS Graph", fmt.Sprintf("%q does not exist.", readUri.Entity))
			return
		case matches(got):
			matchingReads++
			if matchingReads >= opts.ConsecutiveReads {
				return
			}
		default:
			matchingReads = 0
			if time.Since(lastPatch) >= opts.ResendAfter {
				tflog.Info(ctx, fmt.Sprintf("MS Graph still does not return the values written to %q %s after the last write, writing them again",
					uri.Entity, time.Since(lastPatch).Round(time.Second)))
				if !patch() {
					return
				}
				lastPatch = time.Now()
			}
		}

		if time.Since(start) >= opts.Timeout {
			diags.AddError("MS Graph did not apply the update",
				fmt.Sprintf("%q has been written (and written again) successfully, but MS Graph still returned other values after %s.",
					uri.Entity, opts.Timeout))
			return
		}
		select {
		case <-ctx.Done():
			diags.AddError("Unable to update with MS Graph", ctx.Err().Error())
			return
		case <-time.After(opts.PollInterval):
		}
	}
}
