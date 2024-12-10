package wpdefaultvalue

import (
	"context"

	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

//
// ObjectDefaultValue / objectDefaultValueAttributePlanModifier
//

type objectDefaultValueAttributePlanModifier struct {
	defaultValueAttributePlanModifier
	defaultValueRaw map[string]any
}

// ATTENTION: Plan modifiers on the nested object's attributes will overwrite any attribute set here!
// For background see note at top of commong.go.
func ObjectDefaultValue(defaultValueRaw map[string]any) planmodifier.Object {
	return objectDefaultValueAttributePlanModifier{defaultValueRaw: defaultValueRaw}
}

func ObjectDefaultValueEmpty() planmodifier.Object {
	return ObjectDefaultValue(map[string]any{})
}

func (attributePlanModifier objectDefaultValueAttributePlanModifier) PlanModifyObject(ctx context.Context, request planmodifier.ObjectRequest, response *planmodifier.ObjectResponse) {
	defaultValuePlanModify(ctx, &response.Diagnostics, request.Config.Schema.(rsschema.Schema), request.Path, request.ConfigValue, attributePlanModifier.defaultValueRaw, &response.PlanValue)
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
	defaultValuePlanModify(ctx, &response.Diagnostics, request.Config.Schema.(rsschema.Schema), request.Path, request.ConfigValue, attributePlanModifier.defaultValueRaw, &response.PlanValue)
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
	defaultValuePlanModify(ctx, &response.Diagnostics, request.Config.Schema.(rsschema.Schema), request.Path, request.ConfigValue, attributePlanModifier.defaultValueRaw, &response.PlanValue)
}
