package generic

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type WriteSubAction interface {
	Initialize()
	CheckRunAction(wsaOperation OperationType) bool
	ExecutePre(context.Context, *diag.Diagnostics, *WriteSubActionRequest)
	ExecutePost(context.Context, *diag.Diagnostics, *WriteSubActionRequest)
}

type WriteSubActionBase struct {
	Attributes        []string
	AttributesMap     map[string]string
	UriSuffix         string
	UpdateOnly        bool
	EmptyBeforeDelete bool
}

type WriteSubActionRequest struct {
	GenRes         *GenericResource
	Operation      OperationType
	Id             string
	IdAttributer   GetAttributer
	RawVal         map[string]any
	SubActionsData map[string]any
	ReqState       *tfsdk.State
	ReqPlan        *tfsdk.Plan
	RespState      *tfsdk.State
	RespPrivate    PrivateDataGetSetter
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

func (a *WriteSubActionBase) ExecutePre(ctx context.Context, diags *diag.Diagnostics, wsaReq *WriteSubActionRequest,
	executeImpl func(ctx context.Context, diags *diag.Diagnostics, wsaRequest *WriteSubActionRequest)) {
	switch wsaReq.Operation {
	case OperationCreate, OperationUpdate:
		// just extract values to use them in ExecutePost later
		for k := range a.AttributesMap {
			// only save values that actually exist to allow manipulation in CreateModifyFunc, UpdateModifyFunc
			if v, ok := wsaReq.RawVal[k]; ok {
				wsaReq.SubActionsData[k] = v
				delete(wsaReq.RawVal, k)
			}
		}
	case OperationDelete:
		if a.EmptyBeforeDelete {
			if executeImpl == nil {
				panic("executeImpl == nil but EmptyBeforeDelete is set")
			}
			// fill with emtpy arrays and execute request now (i.e. before actually deleting the resource)
			for k := range a.AttributesMap {
				wsaReq.SubActionsData[k] = []any{}
			}
			executeImpl(ctx, diags, wsaReq)
		}
	}
}
