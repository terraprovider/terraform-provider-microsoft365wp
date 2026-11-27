package generic

import (
	"slices"

	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// withSetAlignToConfigOrder returns a copy of attrs in which every SetNestedAttribute (at any depth) additionally has
// the plan modifier wpplanmodifier.SetAlignToConfigOrder (as its last set-level plan modifier), so that the plan
// modifiers of nested attributes get the matching config element (see issue #246).
//
// The nested attribute maps get copied instead of modified, as they can be shared between resources (e.g. the
// assignments schema).
func withSetAlignToConfigOrder(attrs map[string]rsschema.Attribute) map[string]rsschema.Attribute {
	if attrs == nil {
		return nil
	}
	result := make(map[string]rsschema.Attribute, len(attrs))
	for name, attribute := range attrs {
		result[name] = withSetAlignToConfigOrderAttribute(attribute)
	}
	return result
}

func withSetAlignToConfigOrderAttribute(attribute rsschema.Attribute) rsschema.Attribute {
	switch typed := attribute.(type) {
	case rsschema.SetNestedAttribute:
		typed.NestedObject.Attributes = withSetAlignToConfigOrder(typed.NestedObject.Attributes)
		if !slices.ContainsFunc(typed.PlanModifiers, wpplanmodifier.IsSetAlignToConfigOrder) {
			typed.PlanModifiers = append(slices.Clone(typed.PlanModifiers), wpplanmodifier.SetAlignToConfigOrder())
		}
		return typed
	case rsschema.ListNestedAttribute:
		typed.NestedObject.Attributes = withSetAlignToConfigOrder(typed.NestedObject.Attributes)
		return typed
	case rsschema.MapNestedAttribute:
		typed.NestedObject.Attributes = withSetAlignToConfigOrder(typed.NestedObject.Attributes)
		return typed
	case rsschema.SingleNestedAttribute:
		typed.Attributes = withSetAlignToConfigOrder(typed.Attributes)
		return typed
	case OdataDerivedTypeNestedAttributeRs:
		typed.Attributes = withSetAlignToConfigOrder(typed.Attributes)
		return typed
	default:
		return attribute
	}
}
