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

// GitHub #40: tenants that never set allowUserConsentForRiskyApps return null, which must be read as the default
// value false (otherwise every plan shows a diff), while false must still be written as false.
func TestAuthorizationPolicyRiskyAppsNullReadsAsFalse(t *testing.T) {
	ctx := context.Background()
	resp := resource.SchemaResponse{}
	AuthorizationPolicyResource.Schema(ctx, resource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema: %v", resp.Diagnostics)
	}
	s := resp.Schema

	for _, tc := range []struct {
		name  string
		graph any
		want  bool
	}{
		{"null", nil, false},
		{"false", false, false},
		{"true", true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			diags := diag.Diagnostics{}
			raw := map[string]any{"id": "authorizationPolicy", "allowUserConsentForRiskyApps": tc.graph}
			val := generic.ConvertOdataRawToTerraform(ctx, &diags, s, raw, "", nil, "authorizationPolicy", false, nil)
			if diags.HasError() {
				t.Fatalf("diags: %v", diags)
			}
			var attrs map[string]tftypes.Value
			if err := val.As(&attrs); err != nil {
				t.Fatalf("As object: %v", err)
			}
			var got bool
			if err := attrs["allow_user_consent_for_risky_apps"].As(&got); err != nil {
				t.Fatalf("As bool (value %v): %v", attrs["allow_user_consent_for_risky_apps"], err)
			}
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}

			// and back: the value from Terraform is sent as is (no translation to null)
			back := generic.ConvertTerraformToOdataRaw(ctx, &diags, s, false, val, false, tfsdk.Config{}, nil, "Plan")
			if diags.HasError() {
				t.Fatalf("diags: %v", diags)
			}
			if sent, ok := back["allowUserConsentForRiskyApps"].(bool); !ok || sent != tc.want {
				t.Errorf("sent %#v, want %v", back["allowUserConsentForRiskyApps"], tc.want)
			}
		})
	}
}
