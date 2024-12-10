package wpplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// The framework UseStateForUnknown modifiers are unable to deal with null values in state, so we duplicate them here.

type useStateForUnknownModifierBase struct{}

func (pm useStateForUnknownModifierBase) Description(ctx context.Context) string {
	return "Replacement for framework UseStateForUnknown"
}
func (pm useStateForUnknownModifierBase) MarkdownDescription(ctx context.Context) string {
	return pm.Description(ctx)
}

//
// Bool
//

func BoolUseStateForUnknown() planmodifier.Bool {
	return boolUseStateForUnknown{}
}

var _ planmodifier.Bool = (*boolUseStateForUnknown)(nil)

type boolUseStateForUnknown struct {
	useStateForUnknownModifierBase
}

func (pm boolUseStateForUnknown) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, res *planmodifier.BoolResponse) {
	// do nothing on create or on destroy
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.PlanValue.IsUnknown() && !req.ConfigValue.IsUnknown() {
		res.PlanValue = req.StateValue
	}
}

//
// String
//

func StringUseStateForUnknown() planmodifier.String {
	return stringUseStateForUnknown{}
}

var _ planmodifier.String = (*stringUseStateForUnknown)(nil)

type stringUseStateForUnknown struct {
	useStateForUnknownModifierBase
}

func (pm stringUseStateForUnknown) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, res *planmodifier.StringResponse) {
	// do nothing on create or on destroy
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.PlanValue.IsUnknown() && !req.ConfigValue.IsUnknown() {
		res.PlanValue = req.StateValue
	}
}

//
// Int64
//

func Int64UseStateForUnknown() planmodifier.Int64 {
	return int64UseStateForUnknown{}
}

var _ planmodifier.Int64 = (*int64UseStateForUnknown)(nil)

type int64UseStateForUnknown struct {
	useStateForUnknownModifierBase
}

func (pm int64UseStateForUnknown) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, res *planmodifier.Int64Response) {
	// do nothing on create or on destroy
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.PlanValue.IsUnknown() && !req.ConfigValue.IsUnknown() {
		res.PlanValue = req.StateValue
	}
}

//
// Object
//

func ObjectUseStateForUnknown() planmodifier.Object {
	return objectUseStateForUnknown{}
}

var _ planmodifier.Object = (*objectUseStateForUnknown)(nil)

type objectUseStateForUnknown struct {
	useStateForUnknownModifierBase
}

func (pm objectUseStateForUnknown) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, res *planmodifier.ObjectResponse) {
	// do nothing on create or on destroy
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.PlanValue.IsUnknown() && !req.ConfigValue.IsUnknown() {
		res.PlanValue = req.StateValue
	}
}

//
// Set
//

func SetUseStateForUnknown() planmodifier.Set {
	return setUseStateForUnknown{}
}

var _ planmodifier.Set = (*setUseStateForUnknown)(nil)

type setUseStateForUnknown struct {
	useStateForUnknownModifierBase
}

func (pm setUseStateForUnknown) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, res *planmodifier.SetResponse) {
	// do nothing on create or on destroy
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.PlanValue.IsUnknown() && !req.ConfigValue.IsUnknown() {
		res.PlanValue = req.StateValue
	}
}
