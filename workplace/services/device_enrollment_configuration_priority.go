package services

import (
	"context"
	"math"

	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func init() {
	actions := DeviceEnrollmentConfigurationResource.AccessParams.WriteOptions.SubActions
	// Assignment failure can stop execution before the priority action sees the copied plan.
	actions[0] = &enrollmentPriorityCheckpoint{WriteSubAction: actions[0]}
	for i, action := range actions {
		if original, ok := action.(*generic.WriteSubActionAllInOne); ok && original.UriSuffix == "setPriority" {
			actions[i] = &enrollmentPriorityAction{
				WriteSubActionAllInOne: generic.WriteSubActionAllInOne{WriteSubActionBase: original.WriteSubActionBase},
			}
		}
	}
}

type enrollmentPriorityCheckpoint struct{ generic.WriteSubAction }

func (a *enrollmentPriorityCheckpoint) ExecutePost(ctx context.Context, diags *diag.Diagnostics, req *generic.WriteSubActionRequest) {
	if _, requested := req.SubActionsData["priority"]; requested {
		prior := types.Int64Null()
		if req.ReqState != nil {
			diags.Append(req.ReqState.GetAttribute(ctx, path.Root("priority"), &prior)...)
		}
		diags.Append(req.RespState.SetAttribute(ctx, path.Root("priority"), prior)...)
	}
	if !diags.HasError() {
		a.WriteSubAction.ExecutePost(ctx, diags, req)
	}
}

type enrollmentPriorityAction struct{ generic.WriteSubActionAllInOne }

func (a *enrollmentPriorityAction) ExecutePost(ctx context.Context, diags *diag.Diagnostics, req *generic.WriteSubActionRequest) {
	if _, requested := req.SubActionsData["priority"]; !requested {
		return
	}
	var desired types.Int64
	diags.Append(req.ReqPlan.GetAttribute(ctx, path.Root("priority"), &desired)...)
	if diags.HasError() {
		return
	}
	observed, ok := enrollmentReadPriority(ctx, diags, req)
	if !ok || observed.Equal(desired) {
		return
	}
	a.WriteSubActionAllInOne.ExecutePost(ctx, diags, req)
	// An error can follow a partial change, so read back even if the action failed.
	observed, ok = enrollmentReadPriority(ctx, diags, req)
	if ok && !observed.Equal(desired) {
		diags.AddError("Enrollment priority did not converge", "Graph did not return the requested priority; state contains the observed value.")
	}
}

func enrollmentReadPriority(ctx context.Context, diags *diag.Diagnostics, req *generic.WriteSubActionRequest) (types.Int64, bool) {
	var readDiags diag.Diagnostics
	defer func() { diags.Append(readDiags...) }()
	uri := req.GenRes.AccessParams.GetUriWithIdForUD(ctx, &readDiags, "", req.Id, req.IdAttributer)
	if readDiags.HasError() {
		return types.Int64Null(), false
	}
	observed := req.GenRes.AccessParams.ReadRaw(ctx, &readDiags, uri.Entity, false)
	if readDiags.HasError() {
		return types.Int64Null(), false
	}
	value, ok := observed["priority"].(float64)
	if observed["id"] != req.Id || !ok || value < 0 || value > math.MaxInt32 || value != math.Trunc(value) {
		readDiags.AddError("Unable to verify enrollment priority", "Graph did not return the requested profile with a valid priority; retaining the last known value.")
		return types.Int64Null(), false
	}
	priority := types.Int64Value(int64(value))
	readDiags.Append(req.RespState.SetAttribute(ctx, path.Root("priority"), priority)...)
	return priority, !readDiags.HasError()
}
