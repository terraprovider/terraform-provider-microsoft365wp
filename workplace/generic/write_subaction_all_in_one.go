package generic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

var _ WriteSubAction = &WriteSubActionAllInOne{}

type WriteSubActionAllInOne struct {
	WriteSubActionBase

	GraphErrorsAsWarnings bool
	UsePatch              bool
	// Read the entity first and skip the request if MS Graph already has the values (compared by the source attribute
	// names of AttributesMap, i.e. for actions like setPriority whose parameters are attributes of the entity). Avoids
	// redundant actions that MS Graph rejects (e.g. setPriority with the current priority returns 400).
	SkipIfUnchangedInGraph bool
}

func (a *WriteSubActionAllInOne) Initialize() {
	a.WriteSubActionBase.Initialize()
}

func (a *WriteSubActionAllInOne) CheckRunAction(wsaOperation OperationType) bool {
	return a.WriteSubActionBase.CheckRunAction(wsaOperation)
}

func (a *WriteSubActionAllInOne) ExecutePre(ctx context.Context, diags *diag.Diagnostics, wsaReq *WriteSubActionRequest) {
	a.WriteSubActionBase.ExecutePre(ctx, diags, wsaReq, a.executeRequest)
}

func (a *WriteSubActionAllInOne) ExecutePost(ctx context.Context, diags *diag.Diagnostics, wsaReq *WriteSubActionRequest) {
	// for deletes this is a no-op (for now) since we already emptied before delete
	if wsaReq.Operation != OperationDelete {
		a.executeRequest(ctx, diags, wsaReq)
	}
}

func (a *WriteSubActionAllInOne) executeRequest(ctx context.Context, diags *diag.Diagnostics, wsaReq *WriteSubActionRequest) {

	rawBody := make(map[string]any)
	for source_key, target_key := range a.AttributesMap {
		if v, ok := wsaReq.SubActionsData[source_key]; ok {
			rawBody[target_key] = v
		}
	}
	if len(rawBody) == 0 {
		// nothing to do for us here
		return
	}

	jsonBody, err := json.Marshal(rawBody)
	if err != nil {
		diags.AddError(
			"Creation Of OData Sub Action JSON Unsuccessful",
			fmt.Sprintf("Unable to create OData JSON from Terraform %s. Original Error: %s", "Plan", err.Error()),
		)
		return
	}

	parentUri := wsaReq.GenRes.AccessParams.GetUriWithIdForUD(ctx, diags, "", wsaReq.Id, wsaReq.IdAttributer)
	if diags.HasError() {
		return
	}

	if a.SkipIfUnchangedInGraph && wsaReq.Operation != OperationDelete && a.unchangedInGraph(ctx, diags, wsaReq, parentUri) {
		return
	}

	validStatusCodes := []int{http.StatusOK, http.StatusCreated, http.StatusNoContent}
	uri := msgraph.Uri{Entity: parentUri.Entity}
	if a.UriSuffix != "" {
		uri.Entity += "/" + a.UriSuffix
	}
	if a.UsePatch {
		_, _, _, err = wsaReq.GenRes.AccessParams.graphClient.Patch(ctx, msgraph.PatchHttpRequestInput{Body: jsonBody, ValidStatusCodes: validStatusCodes, Uri: uri})
	} else {
		_, _, _, err = wsaReq.GenRes.AccessParams.graphClient.Post(ctx, msgraph.PostHttpRequestInput{Body: jsonBody, ValidStatusCodes: validStatusCodes, Uri: uri})
	}
	if err != nil {
		summary := "Error creating/updating with MS Graph in WriteSubActionAllInOne"
		detail := err.Error()
		if a.GraphErrorsAsWarnings {
			diags.AddWarning(summary, detail)
		} else {
			diags.AddError(summary, detail)
		}
		return
	}
}

// unchangedInGraph reports whether the entity in MS Graph already has all values of the sub-action (see
// SkipIfUnchangedInGraph). A read error is reported as a warning and counts as changed, so that the request still gets
// sent. The entity is read without $select, as not all endpoints support it.
func (a *WriteSubActionAllInOne) unchangedInGraph(ctx context.Context, diags *diag.Diagnostics, wsaReq *WriteSubActionRequest, entityUri msgraph.Uri) bool {
	var readDiags diag.Diagnostics
	current := wsaReq.GenRes.AccessParams.ReadRaw2(ctx, &readDiags, entityUri, "", "", nil, false)
	if readDiags.HasError() {
		for _, d := range readDiags {
			diags.AddWarning("Unable to read current values in WriteSubActionAllInOne, sending the request anyway", d.Summary()+": "+d.Detail())
		}
		return false
	}
	for sourceKey := range a.AttributesMap {
		planned, ok := wsaReq.SubActionsData[sourceKey]
		if !ok {
			continue
		}
		if !jsonEqual(planned, current[sourceKey]) {
			return false
		}
	}
	return true
}

// jsonEqual compares two raw values by their JSON representation (e.g. numbers from plan and MS Graph).
func jsonEqual(a, b any) bool {
	aj, errA := json.Marshal(a)
	bj, errB := json.Marshal(b)
	return errA == nil && errB == nil && string(aj) == string(bj)
}
