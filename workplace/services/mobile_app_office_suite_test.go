package services

import (
	"context"
	"reflect"
	"testing"

	"terraform-provider-microsoft365wp/workplace/generic"
	mobileapputils "terraform-provider-microsoft365wp/workplace/services/mobile_app_utils"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// Values as returned by MS Graph for the designer properties of an officeSuiteApp configured with an XML file.
func officeSuiteXmlPlaceholders() map[string]any {
	return map[string]any{
		"@odata.type":                          "#microsoft.graph.officeSuiteApp",
		"id":                                   "app",
		"displayName":                          "Office",
		"officeConfigurationXml":               "PENvbmZpZ3VyYXRpb24+",
		"autoAcceptEula":                       false,
		"excludedApps":                         nil,
		"installProgressDisplayLevel":          "none",
		"localesToInstall":                     []any{},
		"officePlatformArchitecture":           "none",
		"officeSuiteAppDefaultFileFormat":      "notConfigured",
		"productIds":                           []any{},
		"shouldUninstallOlderVersionsOfOffice": false,
		"targetVersion":                        nil,
		"updateChannel":                        "none",
		"updateVersion":                        nil,
		"useSharedComputerActivation":          false,
	}
}

func mobileAppTestState(t *testing.T, ctx context.Context, graphJson string) *tfsdk.State {
	t.Helper()
	schemaResp := resource.SchemaResponse{}
	(&MobileAppResource).Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatal(schemaResp.Diagnostics)
	}
	translator := generic.NewToFromGraphTranslator(schemaResp.Schema, false, nil)
	raw, err := translator.TerraformFromJsonString(ctx, graphJson)
	if err != nil {
		t.Fatal(err)
	}
	return &tfsdk.State{Schema: schemaResp.Schema, Raw: raw}
}

func TestOfficeSuiteKeepDesignerStateFromPriorState(t *testing.T) {
	ctx := context.Background()
	state := mobileAppTestState(t, ctx, `{
		"@odata.type": "#microsoft.graph.officeSuiteApp", "id": "app", "displayName": "Office",
		"officeConfigurationXml": "PENvbmZpZ3VyYXRpb24+",
		"autoAcceptEula": true, "installProgressDisplayLevel": "none", "localesToInstall": ["en-us"],
		"officePlatformArchitecture": "x64", "officeSuiteAppDefaultFileFormat": "officeOpenXMLFormat",
		"productIds": ["o365ProPlusRetail"], "shouldUninstallOlderVersionsOfOffice": true, "targetVersion": "",
		"updateChannel": "current", "updateVersion": "", "useSharedComputerActivation": false,
		"excludedApps": {"access": false, "bing": false, "excel": false, "groove": false, "infoPath": true,
			"lync": false, "oneDrive": false, "oneNote": false, "outlook": false, "powerPoint": false,
			"publisher": false, "sharePointDesigner": true, "teams": false, "visio": false, "word": false}
	}`)

	rawVal := officeSuiteXmlPlaceholders()
	diags := diag.Diagnostics{}
	mobileapputils.OfficeSuiteKeepDesignerStateRerc(ctx, &diags, generic.ReadExtraRequestCustomParams{ReqState: state, RawVal: rawVal})
	if diags.HasError() || diags.WarningsCount() > 0 {
		t.Fatal(diags)
	}

	expected := map[string]any{
		"autoAcceptEula":                       true,
		"officePlatformArchitecture":           "x64",
		"officeSuiteAppDefaultFileFormat":      "officeOpenXMLFormat",
		"updateChannel":                        "current",
		"targetVersion":                        "",
		"updateVersion":                        "",
		"shouldUninstallOlderVersionsOfOffice": true,
		"productIds":                           []any{"o365ProPlusRetail"},
		"localesToInstall":                     []any{"en-us"},
	}
	for k, v := range expected {
		if !reflect.DeepEqual(rawVal[k], v) {
			t.Errorf("%s: expected %#v, got %#v", k, v, rawVal[k])
		}
	}
	excludedApps, ok := rawVal["excludedApps"].(map[string]any)
	if !ok || excludedApps["lync"] != false || excludedApps["infoPath"] != true || excludedApps["sharePointDesigner"] != true {
		t.Errorf("excludedApps not taken from prior state: %#v", rawVal["excludedApps"])
	}
	if rawVal["officeConfigurationXml"] != "PENvbmZpZ3VyYXRpb24+" || rawVal["displayName"] != "Office" {
		t.Errorf("non-designer properties must stay unchanged: %#v", rawVal)
	}
}

func TestOfficeSuiteKeepDesignerStateWithoutXmlOrState(t *testing.T) {
	ctx := context.Background()

	// designer configured app: values from MS Graph must be kept
	designer := officeSuiteXmlPlaceholders()
	designer["officeConfigurationXml"] = nil
	designer["updateChannel"] = "monthlyEnterprise"
	state := mobileAppTestState(t, ctx, `{"@odata.type": "#microsoft.graph.officeSuiteApp", "id": "app", "updateChannel": "current"}`)
	diags := diag.Diagnostics{}
	mobileapputils.OfficeSuiteKeepDesignerStateRerc(ctx, &diags, generic.ReadExtraRequestCustomParams{ReqState: state, RawVal: designer})
	if designer["updateChannel"] != "monthlyEnterprise" {
		t.Errorf("designer app: expected value from MS Graph, got %#v", designer["updateChannel"])
	}

	// no prior state (data source) and state without office suite (import): values from MS Graph must be kept
	for name, reqState := range map[string]*tfsdk.State{
		"nil state":   nil,
		"import only": mobileAppTestState(t, ctx, `{"id": "app"}`),
	} {
		rawVal := officeSuiteXmlPlaceholders()
		mobileapputils.OfficeSuiteKeepDesignerStateRerc(ctx, &diags, generic.ReadExtraRequestCustomParams{ReqState: reqState, RawVal: rawVal})
		if !reflect.DeepEqual(rawVal, officeSuiteXmlPlaceholders()) {
			t.Errorf("%s: expected unchanged values, got %#v", name, rawVal)
		}
	}
	if diags.HasError() {
		t.Fatal(diags)
	}
}

func TestOfficeSuiteRemoveDesignerProperties(t *testing.T) {
	xmlApp := officeSuiteXmlPlaceholders()
	mobileapputils.OfficeSuiteRemoveDesignerProperties(xmlApp)
	expected := map[string]any{
		"@odata.type":            "#microsoft.graph.officeSuiteApp",
		"id":                     "app",
		"displayName":            "Office",
		"officeConfigurationXml": "PENvbmZpZ3VyYXRpb24+",
	}
	if !reflect.DeepEqual(xmlApp, expected) {
		t.Errorf("XML app: expected only non-designer properties, got %#v", xmlApp)
	}

	designerApp := officeSuiteXmlPlaceholders()
	delete(designerApp, "officeConfigurationXml")
	mobileapputils.OfficeSuiteRemoveDesignerProperties(designerApp)
	if _, ok := designerApp["updateChannel"]; !ok {
		t.Error("designer app: designer properties must be kept")
	}

	otherApp := map[string]any{"@odata.type": "#microsoft.graph.winGetApp", "officeConfigurationXml": "x", "updateChannel": "x"}
	mobileapputils.OfficeSuiteRemoveDesignerProperties(otherApp)
	if _, ok := otherApp["updateChannel"]; !ok {
		t.Error("other app types must stay unchanged")
	}
}
