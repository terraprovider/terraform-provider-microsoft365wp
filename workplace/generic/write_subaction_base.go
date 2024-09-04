package generic

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type WriteSubAction interface {
	Initialize()
	CheckRunAction(wsaOperation OperationType) bool
	ExecutePre(ctx context.Context, diags *diag.Diagnostics, r *GenericResource, wsaOperation OperationType,
		rawVal map[string]any, subActionsData map[string]any, id string, idAttributer GetAttributer)
	ExecutePost(ctx context.Context, diags *diag.Diagnostics, r *GenericResource, wsaOperation OperationType,
		id string, idAttributer GetAttributer, subActionsData map[string]any, reqState *tfsdk.State, reqPlan *tfsdk.Plan)
}

type WriteSubActionBase struct {
	Attributes        []string
	AttributesMap     map[string]string
	UriSuffix         string
	UpdateOnly        bool
	EmptyBeforeDelete bool
}

func (a *WriteSubActionBase) Initialize() {
	if a.AttributesMap == nil {
		a.AttributesMap = make(map[string]string)
	}
	if a.Attributes != nil {
		for _, at := range a.Attributes {
			a.AttributesMap[at] = at
		}
		a.Attributes = nil
	}
}

func (a *WriteSubActionBase) CheckRunAction(wsaOperation OperationType) bool {
	return (wsaOperation == OperationCreate && !a.UpdateOnly) ||
		wsaOperation == OperationUpdate ||
		(wsaOperation == OperationDelete && a.EmptyBeforeDelete)
}

func (a *WriteSubActionBase) ExecutePre(ctx context.Context, diags *diag.Diagnostics, r *GenericResource, wsaOperation OperationType,
	rawVal map[string]any, subActionsData map[string]any, id string, idAttributer GetAttributer,
	executeImpl func(ctx context.Context, diags *diag.Diagnostics, r *GenericResource, id string, idAttributer GetAttributer, subActionsData map[string]any)) {
	switch wsaOperation {
	case OperationCreate, OperationUpdate:
		// just extract values to use them in ExecutePost later
		for k := range a.AttributesMap {
			// only save values that actually exist to allow manipulation in CreateModifyFunc, UpdateModifyFunc
			if v, ok := rawVal[k]; ok {
				subActionsData[k] = v
				delete(rawVal, k)
			}
		}
	case OperationDelete:
		if a.EmptyBeforeDelete {
			if executeImpl == nil {
				panic("executeImpl == nil but EmptyBeforeDelete is set")
			}
			// fill with emtpy arrays and execute request now (i.e. before actually deleting the resource)
			for k := range a.AttributesMap {
				subActionsData[k] = []any{}
			}
			executeImpl(ctx, diags, r, id, idAttributer, subActionsData)
		}
	}
}
