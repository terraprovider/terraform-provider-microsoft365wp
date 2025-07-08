package wpplanmodifier

import (
	"context"
	"fmt"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpschemautil"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

type IsImmutable struct {
	// We only convert OnlyIfCurrentValueIs to attr.Value inside PlanModify to keep provider loading duration low (since
	// schemas of _all_ resources get created at provider loading time no matter whether they are actually being used by
	// the practitioner or not).
	OnlyIfCurrentValueIs any
}

func (pm IsImmutable) Description(ctx context.Context) string {
	return "attribute value cannot be changed"
}
func (pm IsImmutable) MarkdownDescription(ctx context.Context) string {
	return pm.Description(ctx)
}

var _ planmodifier.Bool = (*IsImmutable)(nil)
var _ planmodifier.String = (*IsImmutable)(nil)
var _ planmodifier.Int64 = (*IsImmutable)(nil)
var _ planmodifier.Object = (*IsImmutable)(nil)
var _ planmodifier.Set = (*IsImmutable)(nil)

func (pm IsImmutable) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, res *planmodifier.BoolResponse) {
	// do nothing on create or on destroy
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	pm.planModifyImpl(ctx, &res.Diagnostics, req.Config.Schema.(rsschema.Schema), req.Path, req.StateValue, req.PlanValue)
}

func (pm IsImmutable) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, res *planmodifier.StringResponse) {
	pm.planModifyImpl(ctx, &res.Diagnostics, req.Config.Schema.(rsschema.Schema), req.Path, req.StateValue, req.PlanValue)
}

func (pm IsImmutable) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, res *planmodifier.Int64Response) {
	pm.planModifyImpl(ctx, &res.Diagnostics, req.Config.Schema.(rsschema.Schema), req.Path, req.StateValue, req.PlanValue)
}

func (pm IsImmutable) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, res *planmodifier.ObjectResponse) {
	pm.planModifyImpl(ctx, &res.Diagnostics, req.Config.Schema.(rsschema.Schema), req.Path, req.StateValue, req.PlanValue)
}

func (pm IsImmutable) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, res *planmodifier.SetResponse) {
	pm.planModifyImpl(ctx, &res.Diagnostics, req.Config.Schema.(rsschema.Schema), req.Path, req.StateValue, req.PlanValue)
}

func (pm IsImmutable) planModifyImpl(ctx context.Context, diags *diag.Diagnostics, schema rsschema.Schema, path path.Path, stateValue attr.Value, planValue attr.Value) {
	// do nothing on create or on destroy
	if stateValue.IsNull() || planValue.IsNull() {
		return
	}

	if pm.OnlyIfCurrentValueIs != nil {
		onlyIfCurrentValueIsFw := wpschemautil.GetFwValueFromAny(ctx, diags, &schema, path, pm.OnlyIfCurrentValueIs)
		if diags.HasError() {
			return
		}
		// do nothing if the current value is different from pm.OnlyIfCurrentValueIs
		if !stateValue.Equal(onlyIfCurrentValueIsFw) {
			return
		}
	}

	if !planValue.Equal(stateValue) {
		// 2025/06: In case of `tofu destroy` the plan modifiers unfortunately get called twice for the same entity during
		// the planning phase already. Once with a plan based on the current config as it exists on disk at the time
		// the command gets executed (not expected) and again thereafter with a null plan (this is expected).
		// Therefore we can unfortunately only emit a warning here since otherwise a practitioner would not even be able
		// to do a `tofu destroy` in case an immutable attribute value had changed in the config (and RequiresReplace
		// would not work either).
		// Note that if a resource has actually been removed from the config, the plan modifiers will only get called
		// once with a null plan (as expected).
		diags.AddAttributeWarning(
			path,
			"Apply Might Fail or Be Incomplete as Attribute Value Cannot Be Changed",
			fmt.Sprintf("Attribute '%s' state (current): %s, plan (new): %s", path, stateValue, planValue),
		)
	}
}
