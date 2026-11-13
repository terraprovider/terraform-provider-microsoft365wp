package services

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Shape as returned by MS Graph: every region is listed twice, once with its legacy region group (and no
// geographicLocationType) and once with its current region group.
const cloudPcSupportedRegionsJson = `{"value":[
	{"id":"1","displayName":"northeurope","regionGroup":"europeUnion","geographicLocationType":null},
	{"id":"2","displayName":"northeurope","regionGroup":"ireland","geographicLocationType":"europe"},
	{"id":"3","displayName":"westeurope","regionGroup":"europeUnion","geographicLocationType":null},
	{"id":"4","displayName":"westeurope","regionGroup":"netherlands","geographicLocationType":"europe"},
	{"id":"5","displayName":"germanywestcentral","regionGroup":"germany","geographicLocationType":"europe"},
	{"id":"6","displayName":"eastus","regionGroup":"usEast","geographicLocationType":"usEast"},
	{"id":"7","displayName":"eastus2","regionGroup":"usEast","geographicLocationType":"usEast"}
]}`

func TestCloudPcProvisioningPolicyAddGeographicLocationTypes(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		calls++
		if req.Method != http.MethodGet || req.URL.Path != "/beta/deviceManagement/virtualEndpoint/supportedRegions" {
			t.Errorf("unexpected request: %s %s", req.Method, req.URL)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, cloudPcSupportedRegionsJson)
	}))
	defer server.Close()
	client := msgraph.NewClient(msgraph.VersionBeta)
	client.Endpoint = server.URL
	client.DisableRetries = true

	djc := func(regionGroup, regionName, geo any) map[string]any {
		m := map[string]any{"domainJoinType": "azureADJoin"}
		if regionGroup != nil {
			m["regionGroup"] = regionGroup
		}
		if regionName != nil {
			m["regionName"] = regionName
		}
		if geo != nil {
			m["geographicLocationType"] = geo
		}
		return m
	}
	rawVal := map[string]any{"domainJoinConfigurations": []any{
		djc("ireland", "northeurope", nil),       // reported case: current-only region group
		djc("usEast", "eastus", "europe"),        // stale value from state must be replaced
		djc("germany", "automatic", nil),         // no specific region
		djc("usEast", nil, nil),                  // no region name at all
		djc("europeUnion", "northeurope", nil),   // legacy region group: leave it to MS Graph
		djc(nil, nil, nil),                       // e.g. hybrid join with onPremisesConnectionId only
		djc("unknownGroup", "somewhere", "asia"), // unknown: leave as is
	}}

	diags := diag.Diagnostics{}
	cloudPcProvisioningPolicyAddGeographicLocationTypes(context.Background(), &diags, &client, rawVal)
	if diags.HasError() {
		t.Fatal(diags)
	}
	if calls != 1 {
		t.Fatalf("expected exactly one request, got %d", calls)
	}

	expected := []any{"europe", "usEast", "europe", "usEast", nil, nil, "asia"}
	for i, djc := range rawVal["domainJoinConfigurations"].([]any) {
		if got := djc.(map[string]any)["geographicLocationType"]; got != expected[i] {
			t.Errorf("entry %d: expected %v, got %v", i, expected[i], got)
		}
	}
}

func TestCloudPcProvisioningPolicyAddGeographicLocationTypesNoRequest(t *testing.T) {
	client := msgraph.NewClient(msgraph.VersionBeta)
	client.Endpoint = "http://127.0.0.1:1" // must not be called
	client.DisableRetries = true
	for name, rawVal := range map[string]map[string]any{
		"no domain join configurations": {"displayName": "x"},
		"no region group":               {"domainJoinConfigurations": []any{map[string]any{"onPremisesConnectionId": "id"}}},
	} {
		diags := diag.Diagnostics{}
		cloudPcProvisioningPolicyAddGeographicLocationTypes(context.Background(), &diags, &client, rawVal)
		if diags.HasError() {
			t.Errorf("%s: %v", name, diags)
		}
	}
}

func TestCloudPcDomainJoinRegionModifier(t *testing.T) {
	ctx := context.Background()
	schemaResp := resource.SchemaResponse{}
	(&CloudPcProvisioningPolicyResource).Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatal(schemaResp.Diagnostics)
	}
	translator := generic.NewToFromGraphTranslator(schemaResp.Schema, false, nil)
	fromJson := func(djc string) tftypes.Value {
		raw, err := translator.TerraformFromJsonString(ctx, `{"id":"p","displayName":"x","domainJoinConfigurations":[`+djc+`]}`)
		if err != nil {
			t.Fatal(err)
		}
		return raw
	}

	stateRaw := fromJson(`{"domainJoinType":"azureADJoin","regionName":"northeurope","regionGroup":"ireland","geographicLocationType":"europe"}`)
	djcSetPath := path.Root("domain_join_configurations")

	// mimics the framework: element attributes are addressed by the plan element value, state is paired by index
	elemAttr := func(data interface {
		GetAttribute(context.Context, path.Path, interface{}) diag.Diagnostics
	}, attr string) (types.Object, types.String) {
		var set types.Set
		if diags := data.GetAttribute(ctx, djcSetPath, &set); diags.HasError() {
			t.Fatal(diags)
		}
		elem := set.Elements()[0].(types.Object)
		value, _ := elem.Attributes()[attr].(types.String)
		return elem, value
	}

	run := func(m planmodifier.String, attr string, configDjc string, stateVal tftypes.Value) types.String {
		config := tfsdk.Config{Schema: schemaResp.Schema, Raw: fromJson(configDjc)}
		plan := tfsdk.Plan{Schema: schemaResp.Schema, Raw: config.Raw}
		planElem, configValue := elemAttr(plan, attr)
		state := tfsdk.State{Schema: schemaResp.Schema, Raw: stateRaw}
		var stateValue types.String
		if !stateVal.IsKnown() {
			state = tfsdk.State{Schema: schemaResp.Schema, Raw: tftypes.NewValue(stateRaw.Type(), nil)}
			stateValue = types.StringNull()
		} else {
			_, stateValue = elemAttr(state, attr)
		}
		planValue := configValue
		if configValue.IsNull() {
			planValue = types.StringUnknown()
		}
		req := planmodifier.StringRequest{Path: djcSetPath.AtSetValue(planElem).AtName(attr), Config: config, ConfigValue: configValue,
			State: state, StateValue: stateValue, Plan: plan, PlanValue: planValue}
		resp := planmodifier.StringResponse{PlanValue: planValue}
		m.PlanModifyString(ctx, req, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatal(resp.Diagnostics)
		}
		return resp.PlanValue
	}

	known := tftypes.NewValue(tftypes.String, "")
	created := tftypes.NewValue(tftypes.String, tftypes.UnknownValue)

	regionGroup := cloudPcDomainJoinRegionUseStateUnlessChanged("region_name")
	regionName := cloudPcDomainJoinRegionUseStateUnlessChanged("region_group")
	geo := cloudPcDomainJoinGeographicLocationTypeUnknownOnRegionChange()
	for _, tc := range []struct {
		name     string
		m        planmodifier.String
		attr     string
		config   string
		stateVal tftypes.Value
		expected types.String
	}{
		{"group: only name configured, unchanged", regionGroup, "region_group", `{"regionName":"northeurope"}`, known, types.StringValue("ireland")},
		{"group: only name configured, changed", regionGroup, "region_group", `{"regionName":"westeurope"}`, known, types.StringUnknown()},
		{"group: configured", regionGroup, "region_group", `{"regionName":"westeurope","regionGroup":"netherlands"}`, known, types.StringValue("netherlands")},
		{"group: on create", regionGroup, "region_group", `{"regionName":"northeurope"}`, created, types.StringUnknown()},
		{"name: only group configured, unchanged", regionName, "region_name", `{"regionGroup":"ireland"}`, known, types.StringValue("northeurope")},
		{"name: only group configured, changed", regionName, "region_name", `{"regionGroup":"germany"}`, known, types.StringUnknown()},
		{"geo: both unchanged", geo, "geographic_location_type", `{"regionName":"northeurope","regionGroup":"ireland"}`, known, types.StringValue("europe")},
		{"geo: only name configured, unchanged", geo, "geographic_location_type", `{"regionName":"northeurope"}`, known, types.StringValue("europe")},
		{"geo: name changed", geo, "geographic_location_type", `{"regionName":"westeurope"}`, known, types.StringUnknown()},
		{"geo: group changed", geo, "geographic_location_type", `{"regionName":"northeurope","regionGroup":"germany"}`, known, types.StringUnknown()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := run(tc.m, tc.attr, tc.config, tc.stateVal); !got.Equal(tc.expected) {
				t.Errorf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}
