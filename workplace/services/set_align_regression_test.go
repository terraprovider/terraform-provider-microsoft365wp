package services

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// planTestProvider is a minimal provider that serves the given resources, to run plans through the framework.
type planTestProvider struct {
	resources []func() resource.Resource
}

func (planTestProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "microsoft365wp"
}
func (planTestProvider) Schema(context.Context, provider.SchemaRequest, *provider.SchemaResponse) {}
func (planTestProvider) Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse) {
}
func (p planTestProvider) Resources(context.Context) []func() resource.Resource     { return p.resources }
func (planTestProvider) DataSources(context.Context) []func() datasource.DataSource { return nil }

// jsonToTf converts a JSON value to a tftypes.Value of the given type. Missing object attributes are null. The
// elements of sets keep the order of the JSON array (which is the order sent to the provider).
func jsonToTf(t *testing.T, typ tftypes.Type, v any) tftypes.Value {
	t.Helper()
	if v == nil {
		return tftypes.NewValue(typ, nil)
	}
	switch {
	case typ.Is(tftypes.Object{}):
		ot := typ.(tftypes.Object)
		m := v.(map[string]any)
		for k := range m {
			if _, ok := ot.AttributeTypes[k]; !ok {
				t.Fatalf("jsonToTf: unknown attribute %q", k)
			}
		}
		attrs := make(map[string]tftypes.Value, len(ot.AttributeTypes))
		for name, at := range ot.AttributeTypes {
			attrs[name] = jsonToTf(t, at, m[name])
		}
		return tftypes.NewValue(typ, attrs)
	case typ.Is(tftypes.Set{}), typ.Is(tftypes.List{}):
		var et tftypes.Type
		if st, ok := typ.(tftypes.Set); ok {
			et = st.ElementType
		} else {
			et = typ.(tftypes.List).ElementType
		}
		var els []tftypes.Value
		for _, e := range v.([]any) {
			els = append(els, jsonToTf(t, et, e))
		}
		return tftypes.NewValue(typ, els)
	case typ.Is(tftypes.Map{}):
		et := typ.(tftypes.Map).ElementType
		els := make(map[string]tftypes.Value)
		for k, e := range v.(map[string]any) {
			els[k] = jsonToTf(t, et, e)
		}
		return tftypes.NewValue(typ, els)
	case typ.Is(tftypes.Number):
		return tftypes.NewValue(typ, big.NewFloat(v.(float64)))
	default:
		return tftypes.NewValue(typ, v)
	}
}

// TestSetAlignToConfigOrderRegression runs plans through the framework with unchanged configs that omit a defaulted
// attribute in one set element. The proposed set is ordered like Terraform/OpenTofu order it (cty orders by value, null
// sorts last), i.e. differently from the config. Before the fix (issue #246), the default of the element without the
// attribute was applied to the other element (verified with OpenTofu 1.13.9 and v0.21.6).
func TestSetAlignToConfigOrderRegression(t *testing.T) {
	ctx := context.Background()
	server := providerserver.NewProtocol6(planTestProvider{resources: []func() resource.Resource{
		func() resource.Resource { return &MobileAppResource },
		func() resource.Resource { return &CloudPcProvisioningPolicyResource },
	}})()
	schemaResp, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil || len(schemaResp.Diagnostics) > 0 {
		t.Fatalf("GetProviderSchema: %v %v", err, schemaResp.Diagnostics)
	}

	target := func(group string) string {
		return `{"group": {"group_id": "` + group + `"}, "filter_type": "none"}`
	}
	for _, tc := range []struct {
		name      string
		typeName  string
		attribute string
		// prior state as read from MS Graph, config, and the proposed new state (prior state with the set elements in
		// the order Terraform/OpenTofu send them)
		prior, config, proposed string
	}{
		{
			name: "mobile_app assignments: intent uninstall and omitted", typeName: "microsoft365wp_mobile_app", attribute: "assignments",
			prior: `{"id": "a1", "display_name": "x", "web_link": {"app_url": "https://contoso.test"}, "assignments": [
				{"intent": "uninstall", "target": ` + target("g1") + `}, {"intent": "required", "target": ` + target("g2") + `}]}`,
			config: `{"display_name": "x", "web_link": {"app_url": "https://contoso.test"}, "assignments": [
				{"intent": "uninstall", "target": {"group": {"group_id": "g1"}}}, {"target": {"group": {"group_id": "g2"}}}]}`,
			proposed: `{"id": "a1", "display_name": "x", "web_link": {"app_url": "https://contoso.test"}, "assignments": [
				{"intent": "required", "target": ` + target("g2") + `}, {"intent": "uninstall", "target": ` + target("g1") + `}]}`,
		},
		{
			name: "cloud_pc_provisioning_policy domain_join_configurations: hybridAzureADJoin and omitted", typeName: "microsoft365wp_cloud_pc_provisioning_policy",
			attribute: "domain_join_configurations",
			prior: `{"id": "p1", "display_name": "x", "image_id": "img", "image_display_name": "img", "image_type": "gallery", "domain_join_configurations": [
				{"domain_join_type": "hybridAzureADJoin", "type": "hybridAzureADJoin", "on_premises_connection_id": "c1"},
				{"domain_join_type": "azureADJoin", "type": "azureADJoin", "region_group": "default", "region_name": "automatic"}]}`,
			config: `{"display_name": "x", "image_id": "img", "image_display_name": "img", "image_type": "gallery", "domain_join_configurations": [
				{"domain_join_type": "hybridAzureADJoin", "type": "hybridAzureADJoin", "on_premises_connection_id": "c1"},
				{"region_group": "default", "region_name": "automatic"}]}`,
			proposed: `{"id": "p1", "display_name": "x", "image_id": "img", "image_display_name": "img", "image_type": "gallery", "domain_join_configurations": [
				{"domain_join_type": "azureADJoin", "type": "azureADJoin", "region_group": "default", "region_name": "automatic"},
				{"domain_join_type": "hybridAzureADJoin", "type": "hybridAzureADJoin", "on_premises_connection_id": "c1"}]}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			objType := schemaResp.ResourceSchemas[tc.typeName].ValueType()
			parse := func(s string) tftypes.Value {
				var v any
				if err := json.Unmarshal([]byte(s), &v); err != nil {
					t.Fatal(err)
				}
				return jsonToTf(t, objType, v)
			}
			dyn := func(v tftypes.Value) *tfprotov6.DynamicValue {
				d, err := tfprotov6.NewDynamicValue(objType, v)
				if err != nil {
					t.Fatal(err)
				}
				return &d
			}
			prior := parse(tc.prior)

			resp, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
				TypeName:         tc.typeName,
				PriorState:       dyn(prior),
				ProposedNewState: dyn(parse(tc.proposed)),
				Config:           dyn(parse(tc.config)),
			})
			if err != nil {
				t.Fatal(err)
			}
			for _, d := range resp.Diagnostics {
				if d.Severity == tfprotov6.DiagnosticSeverityError {
					t.Fatalf("diagnostic: %s: %s", d.Summary, d.Detail)
				}
			}
			planned, err := resp.PlannedState.Unmarshal(objType)
			if err != nil {
				t.Fatal(err)
			}
			got, _, _ := tftypes.WalkAttributePath(planned, tftypes.NewAttributePath().WithAttributeName(tc.attribute))
			want, _, _ := tftypes.WalkAttributePath(prior, tftypes.NewAttributePath().WithAttributeName(tc.attribute))
			if !sameSetElements(got.(tftypes.Value), want.(tftypes.Value)) {
				t.Errorf("unchanged config must not change %s:\ngot  %s\nwant %s", tc.attribute, got, want)
			}
		})
	}
}
