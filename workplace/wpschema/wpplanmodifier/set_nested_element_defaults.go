package wpplanmodifier

import (
	"context"
	"fmt"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SetNestedElementDefaults returns a plan modifier for a SetNestedAttribute that plans each element as its config
// element with the null attributes listed in defaults set to their default values (these nested attributes must be
// Optional and Computed, without defaults or plan modifiers of their own).
//
// Defaults (or default plan modifiers) on the nested attributes themselves do not work reliably for elements of sets:
//   - The framework runs the plan modifiers of nested attributes with the config element at the same INDEX as the
//     planned element, but config and plan sets are not ordered the same way. A planned element can therefore get the
//     defaults of another element's null config values (verified with OpenTofu 1.13.9: an unchanged configuration
//     planned the configured values of a site back to the defaults).
//   - The framework's `Default` looks up the config value by the planned element's VALUE. Terraform/OpenTofu fills
//     the null Optional+Computed attributes of the proposed element from the prior state, so the value is not found in
//     the config, the config value counts as null, and the default overwrites configured values as well.
//
// Planning the elements from the config avoids both, and also resets an attribute to its default when it gets removed
// from the config (instead of keeping its prior value as Optional+Computed attributes usually do).
func SetNestedElementDefaults(defaults map[string]attr.Value) planmodifier.Set {
	return setNestedElementDefaults{defaults: defaults}
}

var _ planmodifier.Set = (*setNestedElementDefaults)(nil)

type setNestedElementDefaults struct {
	defaults map[string]attr.Value
}

func (pm setNestedElementDefaults) Description(_ context.Context) string {
	return "Plans each set element from its config element with default values for null attributes."
}

func (pm setNestedElementDefaults) MarkdownDescription(ctx context.Context) string {
	return pm.Description(ctx)
}

func (pm setNestedElementDefaults) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	// null config: the attribute's own default (if any) has already been applied; unknown config: nothing to do
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elemType, ok := req.ConfigValue.ElementType(ctx).(types.ObjectType)
	if !ok {
		resp.Diagnostics.AddAttributeError(req.Path, "SetNestedElementDefaults: invalid element type",
			fmt.Sprintf("The set elements must be objects, but are %s. This is always an error in the provider.", req.ConfigValue.ElementType(ctx)))
		return
	}

	elements := make([]attr.Value, 0, len(req.ConfigValue.Elements()))
	for _, element := range req.ConfigValue.Elements() {
		obj, ok := element.(types.Object)
		if !ok || obj.IsNull() || obj.IsUnknown() {
			elements = append(elements, element)
			continue
		}
		attrs := maps.Clone(obj.Attributes())
		for name, defaultValue := range pm.defaults {
			if v, ok := attrs[name]; ok && v.IsNull() {
				attrs[name] = defaultValue
			}
		}
		planned, diags := types.ObjectValue(elemType.AttrTypes, attrs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		elements = append(elements, planned)
	}

	planned, diags := types.SetValue(elemType, elements)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.PlanValue = planned
}
