package services

import (
	"context"
	"testing"

	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Omitted `assignments` must not be sent (the plan would carry the refreshed value from state and `assign` would
// overwrite changes made elsewhere); an explicitly configured value, even `[]`, must be sent.
func TestRoleScopeTagTerraformToGraphAssignments(t *testing.T) {
	ctx := context.Background()
	resp := resource.SchemaResponse{}
	RoleScopeTagResource.Schema(ctx, resource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema: %v", resp.Diagnostics)
	}
	s := resp.Schema
	objType := s.Type().TerraformType(ctx).(tftypes.Object)
	assignmentsType := objType.AttributeTypes["assignments"]

	cases := []struct {
		name        string
		assignments tftypes.Value
		wantSent    bool
	}{
		{"omitted", tftypes.NewValue(assignmentsType, nil), false},
		{"empty", tftypes.NewValue(assignmentsType, []tftypes.Value{}), true},
		{"unknown", tftypes.NewValue(assignmentsType, tftypes.UnknownValue), true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			config := tfsdk.Config{Schema: s, Raw: tftypes.NewValue(objType, map[string]tftypes.Value{
				"id":           tftypes.NewValue(tftypes.String, nil),
				"description":  tftypes.NewValue(tftypes.String, nil),
				"display_name": tftypes.NewValue(tftypes.String, "tag"),
				"is_built_in":  tftypes.NewValue(tftypes.Bool, nil),
				"assignments":  c.assignments,
			})}
			params := generic.TerraformToGraphMiddlewareParams{
				Config:   config,
				RawVal:   map[string]any{"displayName": "tag", "assignments": []any{}},
				IsUpdate: true,
			}
			diags := diag.Diagnostics{}
			if err := roleScopeTagTerraformToGraphMiddleware(ctx, &diags, &params); err != nil {
				t.Fatalf("middleware: %v", err)
			}
			if diags.HasError() {
				t.Fatalf("diags: %v", diags)
			}
			if _, sent := params.RawVal["assignments"]; sent != c.wantSent {
				t.Errorf("assignments sent = %v, want %v", sent, c.wantSent)
			}
		})
	}
}
