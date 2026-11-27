package wpplanmodifier

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSetAlignToConfigOrder(t *testing.T) {
	ctx := context.Background()
	attrTypes := map[string]attr.Type{"intent": types.StringType, "group": types.StringType, "tags": types.SetType{ElemType: types.StringType}}
	elemType := types.ObjectType{AttrTypes: attrTypes}
	obj := func(intent attr.Value, group string, tags attr.Value) attr.Value {
		return types.ObjectValueMust(attrTypes, map[string]attr.Value{"intent": intent, "group": types.StringValue(group), "tags": tags})
	}
	noTags := types.SetNull(types.StringType)
	emptyTags := types.SetValueMust(types.StringType, []attr.Value{})
	oneTag := types.SetValueMust(types.StringType, []attr.Value{types.StringValue("t")})

	groups := func(s types.Set) []string {
		var result []string
		for _, e := range s.Elements() {
			result = append(result, e.(types.Object).Attributes()["group"].(types.String).ValueString())
		}
		return result
	}

	for _, tc := range []struct {
		name           string
		config, plan   []attr.Value
		wantGroupOrder []string
	}{
		{
			name:           "plan filled from prior and ordered differently",
			config:         []attr.Value{obj(types.StringValue("uninstall"), "g1", noTags), obj(types.StringNull(), "g2", noTags)},
			plan:           []attr.Value{obj(types.StringValue("required"), "g2", emptyTags), obj(types.StringValue("uninstall"), "g1", emptyTags)},
			wantGroupOrder: []string{"g1", "g2"},
		},
		{
			name:           "already aligned",
			config:         []attr.Value{obj(types.StringValue("a"), "g1", noTags), obj(types.StringValue("b"), "g2", noTags)},
			plan:           []attr.Value{obj(types.StringValue("a"), "g1", emptyTags), obj(types.StringValue("b"), "g2", emptyTags)},
			wantGroupOrder: []string{"g1", "g2"},
		},
		{
			// the config element with null values must not take the element of the more specific one
			name:           "exact match first",
			config:         []attr.Value{obj(types.StringNull(), "g1", noTags), obj(types.StringValue("x"), "g1", oneTag)},
			plan:           []attr.Value{obj(types.StringValue("x"), "g1", oneTag), obj(types.StringValue("required"), "g1", emptyTags)},
			wantGroupOrder: []string{"g1", "g1"},
		},
		{
			name:           "nested set values must match",
			config:         []attr.Value{obj(types.StringNull(), "g1", oneTag), obj(types.StringNull(), "g2", emptyTags)},
			plan:           []attr.Value{obj(types.StringValue("r"), "g2", emptyTags), obj(types.StringValue("r"), "g1", oneTag)},
			wantGroupOrder: []string{"g1", "g2"},
		},
		{
			name:           "unknown config values match anything",
			config:         []attr.Value{obj(types.StringUnknown(), "g1", noTags), obj(types.StringNull(), "g2", noTags)},
			plan:           []attr.Value{obj(types.StringValue("r"), "g2", emptyTags), obj(types.StringUnknown(), "g1", emptyTags)},
			wantGroupOrder: []string{"g1", "g2"},
		},
		{
			name:           "plan element without match is kept at the end",
			config:         []attr.Value{obj(types.StringNull(), "g2", noTags)},
			plan:           []attr.Value{obj(types.StringValue("r"), "g9", emptyTags), obj(types.StringValue("r"), "g2", emptyTags)},
			wantGroupOrder: []string{"g2", "g9"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := types.SetValueMust(elemType, tc.config)
			plan := types.SetValueMust(elemType, tc.plan)
			resp := &planmodifier.SetResponse{PlanValue: plan}
			SetAlignToConfigOrder().PlanModifySet(ctx, planmodifier.SetRequest{ConfigValue: config, PlanValue: plan}, resp)
			if resp.Diagnostics.HasError() {
				t.Fatal(resp.Diagnostics)
			}
			if !resp.PlanValue.Equal(plan) {
				t.Errorf("the values must not change: got %s, want %s", resp.PlanValue, plan)
			}
			if got := groups(resp.PlanValue); !reflect.DeepEqual(got, tc.wantGroupOrder) {
				t.Errorf("order: got %v, want %v", got, tc.wantGroupOrder)
			}
			if tc.name == "exact match first" {
				first := resp.PlanValue.Elements()[0].(types.Object).Attributes()["intent"].(types.String).ValueString()
				if first != "required" {
					t.Errorf("first planned element should be the one of the null config element, got intent %q", first)
				}
			}
		})
	}

	// null or unknown config or plan: unchanged
	plan := types.SetValueMust(elemType, []attr.Value{obj(types.StringValue("r"), "g2", emptyTags), obj(types.StringValue("r"), "g1", emptyTags)})
	for _, config := range []types.Set{types.SetNull(elemType), types.SetUnknown(elemType)} {
		resp := &planmodifier.SetResponse{PlanValue: plan}
		SetAlignToConfigOrder().PlanModifySet(ctx, planmodifier.SetRequest{ConfigValue: config, PlanValue: plan}, resp)
		if got := groups(resp.PlanValue); !reflect.DeepEqual(got, []string{"g2", "g1"}) {
			t.Errorf("config %s: order changed to %v", config, got)
		}
	}
}
