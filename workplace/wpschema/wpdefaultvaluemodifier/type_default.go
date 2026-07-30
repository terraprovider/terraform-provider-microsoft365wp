package wpdefaultvaluemodifier

import (
	"context"

	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

type typeDefaultValueAttributePlanModifier struct {
	cachedDefaultValueAttributePlanModifier
}

func TypeDefaultValue() *typeDefaultValueAttributePlanModifier {
	return &typeDefaultValueAttributePlanModifier{}
}

func (attributePlanModifier *typeDefaultValueAttributePlanModifier) PlanModifyBool(ctx context.Context, request planmodifier.BoolRequest, response *planmodifier.BoolResponse) {
	defaultValuePlanModify(ctx, &response.Diagnostics, request.Config.Schema.(rsschema.Schema), request.Path, request.ConfigValue, false, &attributePlanModifier.cachedDefaultValueAttributePlanModifier, &response.PlanValue)
}

func (attributePlanModifier *typeDefaultValueAttributePlanModifier) PlanModifyFloat64(ctx context.Context, request planmodifier.Float64Request, response *planmodifier.Float64Response) {
	defaultValuePlanModify(ctx, &response.Diagnostics, request.Config.Schema.(rsschema.Schema), request.Path, request.ConfigValue, 0, &attributePlanModifier.cachedDefaultValueAttributePlanModifier, &response.PlanValue)
}

func (attributePlanModifier *typeDefaultValueAttributePlanModifier) PlanModifyInt64(ctx context.Context, request planmodifier.Int64Request, response *planmodifier.Int64Response) {
	defaultValuePlanModify(ctx, &response.Diagnostics, request.Config.Schema.(rsschema.Schema), request.Path, request.ConfigValue, 0, &attributePlanModifier.cachedDefaultValueAttributePlanModifier, &response.PlanValue)
}

func (attributePlanModifier *typeDefaultValueAttributePlanModifier) PlanModifyString(ctx context.Context, request planmodifier.StringRequest, response *planmodifier.StringResponse) {
	defaultValuePlanModify(ctx, &response.Diagnostics, request.Config.Schema.(rsschema.Schema), request.Path, request.ConfigValue, "", &attributePlanModifier.cachedDefaultValueAttributePlanModifier, &response.PlanValue)
}

func (attributePlanModifier *typeDefaultValueAttributePlanModifier) PlanModifyObject(ctx context.Context, request planmodifier.ObjectRequest, response *planmodifier.ObjectResponse) {
	defaultValuePlanModify(ctx, &response.Diagnostics, request.Config.Schema.(rsschema.Schema), request.Path, request.ConfigValue, map[string]any{}, &attributePlanModifier.cachedDefaultValueAttributePlanModifier, &response.PlanValue)
}

func (attributePlanModifier *typeDefaultValueAttributePlanModifier) PlanModifyList(ctx context.Context, request planmodifier.ListRequest, response *planmodifier.ListResponse) {
	defaultValuePlanModify(ctx, &response.Diagnostics, request.Config.Schema.(rsschema.Schema), request.Path, request.ConfigValue, []any{}, &attributePlanModifier.cachedDefaultValueAttributePlanModifier, &response.PlanValue)
}

func (attributePlanModifier *typeDefaultValueAttributePlanModifier) PlanModifySet(ctx context.Context, request planmodifier.SetRequest, response *planmodifier.SetResponse) {
	defaultValuePlanModify(ctx, &response.Diagnostics, request.Config.Schema.(rsschema.Schema), request.Path, request.ConfigValue, []any{}, &attributePlanModifier.cachedDefaultValueAttributePlanModifier, &response.PlanValue)
}
