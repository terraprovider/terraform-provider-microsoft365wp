package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// settingsCatalogConfigJson returns the config of a settings catalog policy with n settings, each a choice setting
// with a simple collection child (like the Chrome extension settings in test/tfAllInOneResources).
func settingsCatalogConfigJson(n int) string {
	var settings []string
	for i := 0; i < n; i++ {
		settings = append(settings, fmt.Sprintf(`{"instance": {"definition_id": "def_%d", "choice": {"value": {"value": "def_%d_1",
			"children": [{"definition_id": "def_%d_desc", "simple_collection": {"values": [
				{"string": {"value": "a"}}, {"string": {"value": "b"}}, {"string": {"value": "c"}}]}}]}}}}`, i, i, i))
	}
	return `{"name": "bench", "settings": [` + strings.Join(settings, ",") + `]}`
}

func benchmarkSettingsCatalogPlan(b *testing.B, n int) {
	ctx := context.Background()
	server := providerserver.NewProtocol6(planTestProvider{resources: []func() resource.Resource{
		func() resource.Resource { return &DeviceManagementConfigurationPolicyResource },
	}})()
	schemaResp, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil || len(schemaResp.Diagnostics) > 0 {
		b.Fatalf("GetProviderSchema: %v %v", err, schemaResp.Diagnostics)
	}
	objType := schemaResp.ResourceSchemas["microsoft365wp_device_management_configuration_policy"].ValueType()

	var v any
	if err := json.Unmarshal([]byte(settingsCatalogConfigJson(n)), &v); err != nil {
		b.Fatal(err)
	}
	t := &testing.T{}
	config := jsonToTf(t, objType, v)
	if t.Failed() {
		b.Fatal("jsonToTf failed")
	}
	configDyn, err := tfprotov6.NewDynamicValue(objType, config)
	if err != nil {
		b.Fatal(err)
	}
	prior, _ := tfprotov6.NewDynamicValue(objType, tftypes.NewValue(objType, nil))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
			TypeName:         "microsoft365wp_device_management_configuration_policy",
			PriorState:       &prior,
			ProposedNewState: &configDyn,
			Config:           &configDyn,
		})
		if err != nil {
			b.Fatal(err)
		}
		for _, d := range resp.Diagnostics {
			if d.Severity == tfprotov6.DiagnosticSeverityError {
				b.Fatalf("%s: %s", d.Summary, d.Detail)
			}
		}
	}
}

func BenchmarkSettingsCatalogPlan10(b *testing.B) { benchmarkSettingsCatalogPlan(b, 10) }
func BenchmarkSettingsCatalogPlan90(b *testing.B) { benchmarkSettingsCatalogPlan(b, 90) }
