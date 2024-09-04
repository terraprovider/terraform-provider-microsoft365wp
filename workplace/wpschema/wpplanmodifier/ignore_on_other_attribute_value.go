package wpplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var _ planmodifier.Set = (*IgnoreOnOtherAttributeValuePlanModifier)(nil)

type IgnoreOnOtherAttributeValuePlanModifier struct {
	OtherAttributePath path.Path
	ValuesRespect      []attr.Value
	ValuesIgnore       []attr.Value
}

func (pm *IgnoreOnOtherAttributeValuePlanModifier) Description(ctx context.Context) string {
	return "IgnoreOnOtherAttributeValuePlanModifier"
}
func (pm *IgnoreOnOtherAttributeValuePlanModifier) MarkdownDescription(ctx context.Context) string {
	return pm.Description(ctx)
}

func (pm *IgnoreOnOtherAttributeValuePlanModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, res *planmodifier.SetResponse) {
	// do nothing on create or on destroy
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var otherAttributeValue attr.Value
	diags := req.Plan.GetAttribute(ctx, pm.OtherAttributePath, &otherAttributeValue)
	if diags.HasError() {
		res.Diagnostics.Append(diags...)
		return
	}

	if pm.ValuesIgnore != nil {
		for _, v := range pm.ValuesIgnore {
			if otherAttributeValue == v {
				// ignore
				res.PlanValue = req.StateValue
				return
			}
		}
	}

	if pm.ValuesRespect != nil {
		for _, v := range pm.ValuesRespect {
			if otherAttributeValue == v {
				// do nothing
				return
			}
		}
		// no match, ignore
		res.PlanValue = req.StateValue
	}
}
