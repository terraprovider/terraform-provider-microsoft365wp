package generic

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Several unknown values inside the same set element must not corrupt the set (see GitHub issue #27 follow-up:
// region_group and geographic_location_type of cloud_pc_provisioning_policy.domain_join_configurations).
func TestSetUnknownValuesFromResourceModelSetElement(t *testing.T) {
	ctx := context.Background()
	s := schema.Schema{Attributes: map[string]schema.Attribute{
		"id":   schema.StringAttribute{Computed: true},
		"kind": schema.StringAttribute{Optional: true, Computed: true},
		"djcs": schema.SetNestedAttribute{Optional: true, NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
			"name":  schema.StringAttribute{Optional: true},
			"group": schema.StringAttribute{Optional: true, Computed: true},
			"geo":   schema.StringAttribute{Computed: true},
		}}},
	}}
	elemType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"name": tftypes.String, "group": tftypes.String, "geo": tftypes.String}}
	rootType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String, "kind": tftypes.String, "djcs": tftypes.Set{ElementType: elemType}}}
	elem := func(name, group, geo any) tftypes.Value {
		return tftypes.NewValue(elemType, map[string]tftypes.Value{
			"name": tftypes.NewValue(tftypes.String, name), "group": tftypes.NewValue(tftypes.String, group), "geo": tftypes.NewValue(tftypes.String, geo)})
	}
	root := func(id, kind any, elems ...tftypes.Value) tftypes.Value {
		return tftypes.NewValue(rootType, map[string]tftypes.Value{
			"id": tftypes.NewValue(tftypes.String, id), "kind": tftypes.NewValue(tftypes.String, kind),
			"djcs": tftypes.NewValue(tftypes.Set{ElementType: elemType}, elems)})
	}

	unknown := tftypes.UnknownValue
	state := tfsdk.State{Schema: s, Raw: root("p", unknown, elem("northeurope", unknown, unknown))}
	source := root("p", "dedicated", elem("northeurope", "ireland", "europe"))

	unknowns, err := UnknownValuePaths(ctx, state.Raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(unknowns) != 3 {
		t.Fatalf("expected 3 unknown paths, got %d", len(unknowns))
	}
	if err := SetUnknownValuesFromResourceModel(ctx, &state, unknowns, source); err != nil {
		t.Fatal(err)
	}

	expected := root("p", "dedicated", elem("northeurope", nil, nil))
	if !state.Raw.Equal(expected) {
		t.Errorf("expected %s, got %s", expected, state.Raw)
	}
	if remaining, _ := UnknownValuePaths(ctx, state.Raw); len(remaining) != 0 {
		t.Errorf("unknown values left: %v", remaining)
	}
}
