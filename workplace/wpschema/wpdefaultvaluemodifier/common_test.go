package wpdefaultvaluemodifier

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// A single plan modifier instance may be shared by attributes of different types (e.g. the include_targets attributes of
// authentication_methods_policy), so the cached default value of one attribute must not be returned for another one.
func TestSetDefaultValueSharedByAttributesOfDifferentTypes(t *testing.T) {
	ctx := context.Background()

	sharedDefault := SetDefaultValue([]any{map[string]any{"id": "all_users", "target_type": "group"}})

	schema := rsschema.Schema{
		Attributes: map[string]rsschema.Attribute{
			"plain_targets": rsschema.SetNestedAttribute{
				Optional: true,
				Computed: true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: map[string]rsschema.Attribute{
						"id":          rsschema.StringAttribute{Required: true},
						"target_type": rsschema.StringAttribute{Required: true},
					},
				},
				PlanModifiers: []planmodifier.Set{sharedDefault},
			},
			"extended_targets": rsschema.SetNestedAttribute{
				Optional: true,
				Computed: true,
				NestedObject: rsschema.NestedAttributeObject{
					Attributes: map[string]rsschema.Attribute{
						"id":          rsschema.StringAttribute{Required: true},
						"target_type": rsschema.StringAttribute{Required: true},
						"extra":       rsschema.StringAttribute{Optional: true},
					},
				},
				PlanModifiers: []planmodifier.Set{sharedDefault},
			},
		},
	}

	for _, attributeName := range []string{"plain_targets", "extended_targets", "plain_targets"} {
		attribute := schema.Attributes[attributeName]
		wantType := attribute.GetType()

		request := planmodifier.SetRequest{
			Path:        path.Root(attributeName),
			Config:      tfsdk.Config{Schema: schema},
			ConfigValue: types.SetNull(wantType.(attr.TypeWithElementType).ElementType()),
		}
		response := planmodifier.SetResponse{PlanValue: request.ConfigValue}
		sharedDefault.PlanModifySet(ctx, request, &response)

		if response.Diagnostics.HasError() {
			t.Fatalf("%s: unexpected diagnostics: %v", attributeName, response.Diagnostics)
		}
		if gotType := response.PlanValue.Type(ctx); !gotType.Equal(wantType) {
			t.Errorf("%s: default value has type %s, want %s", attributeName, gotType, wantType)
		}
		if len(response.PlanValue.Elements()) != 1 {
			t.Errorf("%s: default value has %d elements, want 1", attributeName, len(response.PlanValue.Elements()))
		}
	}
}

func TestGetSchemaPathKeyIgnoresElementKeys(t *testing.T) {
	first := path.Root("a").AtSetValue(types.StringValue("x")).AtName("b")
	second := path.Root("a").AtSetValue(types.StringValue("y")).AtName("b")
	other := path.Root("a").AtName("b")

	if getSchemaPathKey(first) != getSchemaPathKey(second) {
		t.Errorf("keys of different set elements differ: %q vs %q", getSchemaPathKey(first), getSchemaPathKey(second))
	}
	if getSchemaPathKey(first) == getSchemaPathKey(other) {
		t.Errorf("keys of different schema paths are equal: %q", getSchemaPathKey(first))
	}
}
