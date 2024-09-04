package generic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/manicminer/hamilton/msgraph"
)

var _ WriteSubAction = &WriteSubActionAllInOne{}

type WriteSubActionAllInOne struct {
	WriteSubActionBase

	GraphErrorsAsWarnings bool
}

func (a *WriteSubActionAllInOne) Initialize() {
	a.WriteSubActionBase.Initialize()
}

func (a *WriteSubActionAllInOne) CheckRunAction(wsaOperation OperationType) bool {
	return a.WriteSubActionBase.CheckRunAction(wsaOperation)
}

func (a *WriteSubActionAllInOne) ExecutePre(ctx context.Context, diags *diag.Diagnostics, r *GenericResource, wsaOperation OperationType,
	rawVal map[string]any, subActionsData map[string]any, id string, idAttributer GetAttributer) {
	a.WriteSubActionBase.ExecutePre(ctx, diags, r, wsaOperation, rawVal, subActionsData, id, idAttributer, a.executeRequest)
}

func (a *WriteSubActionAllInOne) ExecutePost(ctx context.Context, diags *diag.Diagnostics, r *GenericResource, wsaOperation OperationType,
	id string, idAttributer GetAttributer, subActionsData map[string]any, _ *tfsdk.State, _ *tfsdk.Plan) {
	// for deletes this is a no-op (for now) since we already emptied before delete
	if wsaOperation != OperationDelete {
		a.executeRequest(ctx, diags, r, id, idAttributer, subActionsData)
	}
}

func (a *WriteSubActionAllInOne) executeRequest(ctx context.Context, diags *diag.Diagnostics, r *GenericResource,
	id string, idAttributer GetAttributer, subActionsData map[string]any) {

	rawBody := make(map[string]any)
	for source_key, target_key := range a.AttributesMap {
		if v, ok := subActionsData[source_key]; ok {
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

	parentUri := r.AccessParams.GetUriWithIdForUD(ctx, diags, "", id, idAttributer)
	if diags.HasError() {
		return
	}

	_, _, _, err = r.AccessParams.graphClient.Post(ctx, msgraph.PostHttpRequestInput{
		Body:             []byte(jsonBody),
		ValidStatusCodes: []int{http.StatusOK, http.StatusCreated, http.StatusNoContent},
		Uri: msgraph.Uri{
			Entity: fmt.Sprintf("%s/%s", parentUri.Entity, a.UriSuffix),
		},
	})
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
