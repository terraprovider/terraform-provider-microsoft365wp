package generic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/manicminer/hamilton/msgraph"
)

var _ WriteSubAction = &WriteSubActionAllInOne{}

type WriteSubActionAllInOne struct {
	WriteSubActionBase

	GraphErrorsAsWarnings bool
	UsePatch              bool
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
