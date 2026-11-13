package mobileapputils

import (
	"context"
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// An officeSuiteApp is either configured with the settings of the Intune designer or with an Office Customization Tool
// XML file (officeConfigurationXml). The mode gets fixed when the app is created (MS Graph rejects any PATCH that tries
// to switch it). For apps configured with an XML file, MS Graph
// - rejects any designer property in a PATCH (400 "This app was created with an Office Customization tool XML
//   configuration file. It can only be patched with another Office Customization tool XML configuration file."),
//   regardless of its value,
// - silently ignores designer properties in a POST,
// - returns placeholder values for the designer properties on read ("none", "notConfigured", [], false, null).
// See GitHub issue #26.

const officeSuiteOdataType = "#microsoft.graph.officeSuiteApp"

var officeSuiteDesignerProperties = []string{
	"autoAcceptEula",
	"excludedApps",
	"installProgressDisplayLevel",
	"localesToInstall",
	"officePlatformArchitecture",
	"officeSuiteAppDefaultFileFormat",
	"productIds",
	"shouldUninstallOlderVersionsOfOffice",
	"targetVersion",
	"updateChannel",
	"updateVersion",
	"useSharedComputerActivation",
}

func officeSuiteHasXmlConfiguration(rawVal map[string]any) bool {
	if rawVal["@odata.type"] != officeSuiteOdataType {
		return false
	}
	xml, _ := rawVal["officeConfigurationXml"].(string)
	return xml != ""
}

// OfficeSuiteRemoveDesignerProperties removes all designer properties from the write request of an officeSuiteApp
// that is configured with an XML file (MS Graph would reject them on PATCH and ignore them on POST).
func OfficeSuiteRemoveDesignerProperties(rawVal map[string]any) {
	if !officeSuiteHasXmlConfiguration(rawVal) {
		return
	}
	for _, p := range officeSuiteDesignerProperties {
		delete(rawVal, p)
	}
}

// OfficeSuiteKeepDesignerStateRerc replaces the placeholder values that MS Graph returns for the designer properties of
// an officeSuiteApp configured with an XML file with the values from prior state. Since these properties never get sent
// to MS Graph for such apps, this avoids a perpetual diff against the placeholders. Without prior state (e.g. right after
// import) the values from MS Graph are kept.
func OfficeSuiteKeepDesignerStateRerc(ctx context.Context, diags *diag.Diagnostics, params generic.ReadExtraRequestCustomParams) {

	if !officeSuiteHasXmlConfiguration(params.RawVal) {
		return
	}
	if params.ReqState == nil || params.ReqState.Raw.IsNull() || !params.ReqState.Raw.IsKnown() {
		return
	}

	translator := generic.NewToFromGraphTranslator(params.ReqState.Schema, false, nil)
	stateRaw, err := translator.TerraformAsRaw(ctx, params.ReqState.Raw)
	if err != nil {
		diags.AddWarning("Unable to read prior state of windows_office_suite", err.Error())
		return
	}
	if stateRaw["@odata.type"] != officeSuiteOdataType {
		return
	}

	for _, p := range officeSuiteDesignerProperties {
		if v, ok := stateRaw[p]; ok {
			params.RawVal[p] = v
		}
	}
	tflog.Debug(ctx, "OfficeSuiteKeepDesignerStateRerc: kept designer properties from prior state for XML configured app")
}

// OfficeSuiteXmlModeChangeRequiresReplace forces replacement if office_configuration_xml changes between set and
// unset, since MS Graph does not allow to switch an existing app between designer and XML configuration. Changing the
// XML content itself works in place.
func OfficeSuiteXmlModeChangeRequiresReplace() planmodifier.String {
	return stringplanmodifier.RequiresReplaceIf(
		func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
			if req.PlanValue.IsUnknown() {
				return // cannot tell yet, MS Graph will reject the update if the mode would change
			}
			stateHasXml := !req.StateValue.IsNull() && !req.StateValue.IsUnknown() && req.StateValue.ValueString() != ""
			planHasXml := !req.PlanValue.IsNull() && req.PlanValue.ValueString() != ""
			resp.RequiresReplace = stateHasXml != planHasXml
		},
		"Switching between designer and XML configuration requires replacement.",
		"Switching between designer and XML configuration requires replacement.",
	)
}

// OfficeSuiteDesignerOnlyRequiresReplace forces replacement on changes of an attribute (e.g. office_platform_architecture)
// only for apps configured with the designer. For apps configured with an XML file the attribute does not get sent to
// MS Graph at all (the XML defines it instead).
func OfficeSuiteDesignerOnlyRequiresReplace() planmodifier.String {
	return stringplanmodifier.RequiresReplaceIf(
		func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
			var xml types.String
			resp.Diagnostics.Append(req.Config.GetAttribute(ctx, req.Path.ParentPath().AtName("office_configuration_xml"), &xml)...)
			if resp.Diagnostics.HasError() || xml.IsUnknown() {
				return
			}
			resp.RequiresReplace = xml.IsNull() || xml.ValueString() == ""
		},
		"Changing this value requires replacement unless office_configuration_xml is used.",
		"Changing this value requires replacement unless `office_configuration_xml` is used.",
	)
}
