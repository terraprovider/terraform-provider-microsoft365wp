package generic

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

var _ WriteSubAction = &WriteSubActionPostAndPatch{}

type WriteSubActionPostAndPatch struct {
	AttributesPatch        []string
	AttributesPatchExclude []string
	AttributesKeep         []string
}

func (a *WriteSubActionPostAndPatch) Initialize() {
}

func (a *WriteSubActionPostAndPatch) CheckRunAction(wsaOperation OperationType) bool {
	// run on create only
	return wsaOperation == OperationCreate
}

func (a *WriteSubActionPostAndPatch) ExecutePre(ctx context.Context, diags *diag.Diagnostics, wsaReq *WriteSubActionRequest) {
	// copy all attributes and only remove the ones that should not be included in the first POST request
	attributesForPatch := make(map[string]any)
	for k, v := range wsaReq.RawVal {
		if len(a.AttributesPatch) > 0 && !slices.Contains(a.AttributesPatch, k) {
			continue
		}
		if !slices.Contains(a.AttributesPatchExclude, k) {
			attributesForPatch[k] = v
		}
		if !slices.Contains(a.AttributesKeep, k) {
			delete(wsaReq.RawVal, k)
		}
	}
	wsaReq.SubActionsData["WriteSubActionPostAndPatch"] = attributesForPatch
}

func (a *WriteSubActionPostAndPatch) ExecutePost(ctx context.Context, diags *diag.Diagnostics, wsaReq *WriteSubActionRequest) {

	attributesForPatch := wsaReq.SubActionsData["WriteSubActionPostAndPatch"].(map[string]any)
	if len(attributesForPatch) == 0 {
		// nothing to do for us here
		return
	}

	wsaReq.GenRes.AccessParams.UpdateRaw(ctx, diags, "", wsaReq.Id, wsaReq.IdAttributer, attributesForPatch)
}
