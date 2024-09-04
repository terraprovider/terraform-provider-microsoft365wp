package wpdefaultvalue

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

type typeDefaultValueAttributePlanModifier struct {
	defaultValueAttributePlanModifier
}

func TypeDefaultValue() typeDefaultValueAttributePlanModifier {
	return typeDefaultValueAttributePlanModifier{}
}

func (attributePlanModifier typeDefaultValueAttributePlanModifier) PlanModifyBool(ctx context.Context, request planmodifier.BoolRequest, response *planmodifier.BoolResponse) {
	defaultValuePlanModify(ctx, request.ConfigValue, &response.PlanValue, &response.Diagnostics, false)
}

func (attributePlanModifier typeDefaultValueAttributePlanModifier) PlanModifyFloat64(ctx context.Context, request planmodifier.Float64Request, response *planmodifier.Float64Response) {
	defaultValuePlanModify(ctx, request.ConfigValue, &response.PlanValue, &response.Diagnostics, 0)
}

func (attributePlanModifier typeDefaultValueAttributePlanModifier) PlanModifyInt64(ctx context.Context, request planmodifier.Int64Request, response *planmodifier.Int64Response) {
	defaultValuePlanModify(ctx, request.ConfigValue, &response.PlanValue, &response.Diagnostics, 0)
}

func (attributePlanModifier typeDefaultValueAttributePlanModifier) PlanModifyString(ctx context.Context, request planmodifier.StringRequest, response *planmodifier.StringResponse) {
	defaultValuePlanModify(ctx, request.ConfigValue, &response.PlanValue, &response.Diagnostics, "")
}

func (attributePlanModifier typeDefaultValueAttributePlanModifier) PlanModifyObject(ctx context.Context, request planmodifier.ObjectRequest, response *planmodifier.ObjectResponse) {
	defaultValuePlanModify(ctx, request.ConfigValue, &response.PlanValue, &response.Diagnostics, map[string]any{})
}

func (attributePlanModifier typeDefaultValueAttributePlanModifier) PlanModifyList(ctx context.Context, request planmodifier.ListRequest, response *planmodifier.ListResponse) {
	defaultValuePlanModify(ctx, request.ConfigValue, &response.PlanValue, &response.Diagnostics, []any{})
}

func (attributePlanModifier typeDefaultValueAttributePlanModifier) PlanModifySet(ctx context.Context, request planmodifier.SetRequest, response *planmodifier.SetResponse) {
	defaultValuePlanModify(ctx, request.ConfigValue, &response.PlanValue, &response.Diagnostics, []any{})
}
