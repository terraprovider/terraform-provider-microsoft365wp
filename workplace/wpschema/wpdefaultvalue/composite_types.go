package wpdefaultvalue

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

//
// ObjectDefaultValue / objectDefaultValueAttributePlanModifier
//

type objectDefaultValueAttributePlanModifier struct {
	defaultValueAttributePlanModifier
	defaultValueRaw map[string]any
}

func ObjectDefaultValue(defaultValueRaw map[string]any) planmodifier.Object {
	return objectDefaultValueAttributePlanModifier{defaultValueRaw: defaultValueRaw}
}

func ObjectDefaultValueEmpty() planmodifier.Object {
	return ObjectDefaultValue(map[string]any{})
}

func (attributePlanModifier objectDefaultValueAttributePlanModifier) PlanModifyObject(ctx context.Context, request planmodifier.ObjectRequest, response *planmodifier.ObjectResponse) {
	defaultValuePlanModify(ctx, request.ConfigValue, &response.PlanValue, &response.Diagnostics, attributePlanModifier.defaultValueRaw)
}

//
// ListDefaultValue / listDefaultValueAttributePlanModifier
//

type listDefaultValueAttributePlanModifier struct {
	defaultValueAttributePlanModifier
	defaultValueRaw []any
}

func ListDefaultValue(defaultValueRaw []any) planmodifier.List {
	return listDefaultValueAttributePlanModifier{defaultValueRaw: defaultValueRaw}
}

func ListDefaultValueEmpty() planmodifier.List {
	return ListDefaultValue([]any{})
}

func (attributePlanModifier listDefaultValueAttributePlanModifier) PlanModifyList(ctx context.Context, request planmodifier.ListRequest, response *planmodifier.ListResponse) {
	defaultValuePlanModify(ctx, request.ConfigValue, &response.PlanValue, &response.Diagnostics, attributePlanModifier.defaultValueRaw)
}

//
// SetDefaultValue / setDefaultValueAttributePlanModifier
//

type setDefaultValueAttributePlanModifier struct {
	defaultValueAttributePlanModifier
	defaultValueRaw []any
}

func SetDefaultValue(defaultValueRaw []any) planmodifier.Set {
	return setDefaultValueAttributePlanModifier{defaultValueRaw: defaultValueRaw}
}

func SetDefaultValueEmpty() planmodifier.Set {
	return SetDefaultValue([]any{})
}

func (attributePlanModifier setDefaultValueAttributePlanModifier) PlanModifySet(ctx context.Context, request planmodifier.SetRequest, response *planmodifier.SetResponse) {
	defaultValuePlanModify(ctx, request.ConfigValue, &response.PlanValue, &response.Diagnostics, attributePlanModifier.defaultValueRaw)
}
