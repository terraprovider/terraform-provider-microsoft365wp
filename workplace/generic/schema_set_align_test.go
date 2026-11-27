package generic

import (
	"testing"

	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func countAlign(modifiers []planmodifier.Set) int {
	n := 0
	for _, m := range modifiers {
		if wpplanmodifier.IsSetAlignToConfigOrder(m) {
			n++
		}
	}
	return n
}

func TestWithSetAlignToConfigOrder(t *testing.T) {
	inner := rsschema.SetNestedAttribute{
		Optional:     true,
		NestedObject: rsschema.NestedAttributeObject{Attributes: map[string]rsschema.Attribute{"x": rsschema.StringAttribute{Optional: true}}},
	}
	// shared like the assignments schema
	shared := map[string]rsschema.Attribute{"inner": inner}
	attrs := map[string]rsschema.Attribute{
		"outer": rsschema.SetNestedAttribute{
			Optional:      true,
			NestedObject:  rsschema.NestedAttributeObject{Attributes: shared},
			PlanModifiers: []planmodifier.Set{wpplanmodifier.SetSortByStringKey("x")},
		},
		"single":  rsschema.SingleNestedAttribute{Optional: true, Attributes: shared},
		"derived": OdataDerivedTypeNestedAttributeRs{SingleNestedAttribute: rsschema.SingleNestedAttribute{Optional: true, Attributes: shared}},
		"list":    rsschema.ListNestedAttribute{Optional: true, NestedObject: rsschema.NestedAttributeObject{Attributes: shared}},
		"plain":   rsschema.StringAttribute{Optional: true},
	}

	result := withSetAlignToConfigOrder(attrs)
	// applying it twice must not add the modifier twice
	result = withSetAlignToConfigOrder(result)

	outer := result["outer"].(rsschema.SetNestedAttribute)
	if countAlign(outer.PlanModifiers) != 1 || len(outer.PlanModifiers) != 2 || !wpplanmodifier.IsSetAlignToConfigOrder(outer.PlanModifiers[1]) {
		t.Errorf("outer: want the existing modifier followed by exactly one SetAlignToConfigOrder, got %v", outer.PlanModifiers)
	}
	for _, nested := range []map[string]rsschema.Attribute{
		outer.NestedObject.Attributes,
		result["single"].(rsschema.SingleNestedAttribute).Attributes,
		result["derived"].(OdataDerivedTypeNestedAttributeRs).Attributes,
		result["list"].(rsschema.ListNestedAttribute).NestedObject.Attributes,
	} {
		if n := countAlign(nested["inner"].(rsschema.SetNestedAttribute).PlanModifiers); n != 1 {
			t.Errorf("inner: want exactly one SetAlignToConfigOrder, got %d", n)
		}
	}

	// the original (shared) attributes must be unchanged
	if len(attrs["outer"].(rsschema.SetNestedAttribute).PlanModifiers) != 1 || len(shared["inner"].(rsschema.SetNestedAttribute).PlanModifiers) != 0 {
		t.Error("the original attributes have been modified")
	}
}
