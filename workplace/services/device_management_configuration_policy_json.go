package services

// INFO#1
// ==========
// There seems to be a bug in the TF plugin framework so that Computed: true does not seem to work for ListAttribute. TF
// still complains if anything gets written to `settings_json` that differs from the literal config.
// Therefore we currently:
// - Save just the literal config values to the state on create or update (and do not refresh from MS Graph after the
//   create/update).
// - On read/refresh we strip down the JSON returned from MS Graph by removing all non-required attributes and then only
//   save the result to the state. This alone will work fine for Terraform as long as the config matches exactly the
//   stripped down response from MS Graph.
// - On plan, we additionally use a plan modifier that takes both the config/plan JSON settings as well as the state
//   JSON settings, strips them both down again (as above), sorts the list and then compares the results. If there is no
//   relevant difference left, we replace the plan settings with the state settings and TF will happily recognize that
//   there is no change required (TF will accept config or state for plan, but would panic for any other plan due to the
//   bug mentioned above).
//   This plan modifier helps whenever the config does not match exactly the stripped down response from MS Graph (e.g.
//   contains additional fields from a JSON dump etc.). The reason we also have to strip down the state again is that
//   directly after creation the state will contain exactly what was in the config (e.g. also additional fields) and
//   therefore needs to be stripped down for comparison as well.
// - Small caveat: If the config only contains the attributes that are required then the TF plan UI will properly show
//   only the individual settings items to be updated that actually have changed. But if the config contains attributes
//   that usually get filtered out then also these items will be shown as changed to user even if nothing relevant has
//   changed on them.

// INFO#2
// ==========
// `settings_json`` vs. `settings.setting_instance`: I could not get a schema closer to the actual REST call (based on
// schema.SetNestedAttribute with an inner setting_instance) to work as its PlanModifiers would not be enough for it
// to be recoqnized as "eqal". Therefore we stick with the simpler interface for now (which seems ok as settingInstance
// really seems to only exist due to syntactic requirements in MS Graph)

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	DeviceManagementConfigurationPolicyJsonResource = generic.GenericResource{
		TypeNameSuffix: "device_management_configuration_policy_json",
		SpecificSchema: deviceManagementConfigurationPolicyJsonResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:         "/deviceManagement/configurationPolicies",
			ReadOptions:     deviceManagementConfigurationPolicyJsonReadOptions,
			WriteSubActions: deviceManagementConfigurationPolicyJsonWriteSubActions,
			UsePutForUpdate: true,
			// See INFO#1, INFO#2
			TerraformToGraphMiddleware: deviceManagementConfigurationPolicyJsonTerraformToGraphMiddleware,
			GraphToTerraformMiddleware: deviceManagementConfigurationPolicyJsonGraphToTerraformMiddleware,
		},
	}

	DeviceManagementConfigurationPolicyJsonSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&DeviceManagementConfigurationPolicyJsonResource)

	DeviceManagementConfigurationPolicyJsonPluralDataSource = generic.CreateGenericDataSourcePluralFromSingular(
		&DeviceManagementConfigurationPolicyJsonSingularDataSource, "device_management_configuration_policies_json")
)

var deviceManagementConfigurationPolicyJsonReadOptions = generic.ReadOptions{
	ODataExpand: "settings",
	ExtraRequests: []generic.ReadExtraRequest{
		{
			ParentAttribute: "assignments",
			UriSuffix:       "assignments",
		},
	},
}

var deviceManagementConfigurationPolicyJsonWriteSubActions = []generic.WriteSubAction{
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			Attributes: []string{"assignments"},
			UriSuffix:  "assign",
		},
	},
}

var deviceManagementConfigurationPolicyJsonResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // deviceManagementConfigurationPolicy
		"id": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Policy creation date and time",
		},
		"creation_source": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Policy creation source",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Policy description",
		},
		"is_assigned": schema.BoolAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Policy last modification date and time",
		},
		"name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Policy name",
		},
		"platforms": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				wpvalidator.FlagEnumValues("none", "android", "iOS", "macOS", "windows10X", "windows10", "linux", "unknownFutureValue"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("windows10")},
			Computed:            true,
			MarkdownDescription: "Platforms for this policy / Supported platform types; possible values are: `none` (Indicates the settings contained in this configuration don't have a platform set.), `android` (Indicates that the settings contained in associated configuration applies to the Android operating system. ), `iOS` (Indicates that the settings contained in associated configuration applies to the iOS operating system.), `macOS` (Indicates that the settings contained in associated configuration applies to the MacOS operating system.), `windows10X` (Indicates that the settings contained in associated configuration applies to the Windows 10X operating system.), `windows10` (Indicates that the settings contained in associated configuration applies to the Windows 10 operating system.), `linux` (Indicates that the settings contained in associated configuration applies to the Linux operating system.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.)",
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: "List of Scope Tags for this Entity instance.",
		},
		"setting_count": schema.Int64Attribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
			MarkdownDescription: "Number of settings",
		},
		"technologies": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				wpvalidator.FlagEnumValues("none", "mdm", "windows10XManagement", "configManager", "appleRemoteManagement", "microsoftSense", "exchangeOnline", "mobileApplicationManagement", "linuxMdm", "enrollment", "endpointPrivilegeManagement", "unknownFutureValue", "windowsOsRecovery"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("mdm")},
			Computed:            true,
			MarkdownDescription: "Technologies for this policy / Describes which technology this setting can be deployed with; possible values are: `none` (Default. Setting cannot be deployed through any channel.), `mdm` (Setting can be deployed through the MDM channel.), `windows10XManagement` (Setting can be deployed through the Windows10XManagement channel), `configManager` (Setting can be deployed through the ConfigManager channel.), `appleRemoteManagement` (Setting can be deployed through the AppleRemoteManagement channel.), `microsoftSense` (Setting can be deployed through the SENSE agent channel.), `exchangeOnline` (Setting can be deployed through the Exchange Online agent channel.), `mobileApplicationManagement` (Setting can be deployed through the Mobile Application Management (MAM) channel), `linuxMdm` (Setting can be deployed through the Linux Mdm channel.), `enrollment` (Setting can be deployed through device enrollment.), `endpointPrivilegeManagement` (Setting can be deployed using the Endpoint privilege management channel), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.), `windowsOsRecovery` (Setting can be deployed using the Operating System Recovery channel)",
		},
		"template_reference": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // deviceManagementConfigurationPolicyTemplateReference
				"template_display_name": schema.StringAttribute{
					Computed:      true,
					PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
				},
				"template_display_version": schema.StringAttribute{
					Computed:      true,
					PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
				},
				"template_family": schema.StringAttribute{
					Computed: true,
					Validators: []validator.String{
						stringvalidator.OneOf("none", "endpointSecurityAntivirus", "endpointSecurityDiskEncryption", "endpointSecurityFirewall", "endpointSecurityEndpointDetectionAndResponse", "endpointSecurityAttackSurfaceReduction", "endpointSecurityAccountProtection", "endpointSecurityApplicationControl", "endpointSecurityEndpointPrivilegeManagement", "enrollmentConfiguration", "appQuietTime", "baseline", "unknownFutureValue", "deviceConfigurationScripts", "deviceConfigurationPolicies", "windowsOsRecoveryPolicies", "companyPortal"),
					},
					PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
					MarkdownDescription: "Describes the TemplateFamily for the Template entity; possible values are: `none` (Default for Template Family when Policy is not linked to a Template), `endpointSecurityAntivirus` (Template Family for EndpointSecurityAntivirus that manages the discrete group of antivirus settings for managed devices), `endpointSecurityDiskEncryption` (Template Family for EndpointSecurityDiskEncryption that provides settings that are relevant for a devices built-in encryption  method, like FileVault or BitLocker), `endpointSecurityFirewall` (Template Family for EndpointSecurityFirewall that helps configure a devices built-in firewall for device that run macOS and Windows 10), `endpointSecurityEndpointDetectionAndResponse` (Template Family for EndpointSecurityEndpointDetectionAndResponse that facilitates management of the EDR settings and onboard devices to Microsoft Defender for Endpoint), `endpointSecurityAttackSurfaceReduction` (Template Family for EndpointSecurityAttackSurfaceReduction that help reduce your attack surfaces, by minimizing the places where your organization is vulnerable to cyberthreats and attacks), `endpointSecurityAccountProtection` (Template Family for EndpointSecurityAccountProtection that facilitates protecting the identity and accounts of users), `endpointSecurityApplicationControl` (Template Family for ApplicationControl that helps mitigate security threats by restricting the applications that users can run and the code that runs in the System Core (kernel)), `endpointSecurityEndpointPrivilegeManagement` (Template Family for EPM Elevation Rules), `enrollmentConfiguration` (Template Family for EnrollmentConfiguration), `appQuietTime` (Template Family for QuietTimeIndicates Template Family for all the Apps QuietTime policies and templates), `baseline` (Template Family for Baseline), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.), `deviceConfigurationScripts` (Template Family for device configuration scripts), `deviceConfigurationPolicies` (Template Family for device configuration policies), `windowsOsRecoveryPolicies` (Template Family for windowsOsRecovery that can be applied during a Windows operating system recovery), `companyPortal` (Template Family for Company Portal settings)",
				},
				"template_id": schema.StringAttribute{
					Optional: true,
					PlanModifiers: []planmodifier.String{
						wpdefaultvalue.StringDefaultValue(""),
						stringplanmodifier.RequiresReplace(),
					},
					Computed:            true,
					MarkdownDescription: "Template id",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Template reference information / Policy template reference information / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationPolicyTemplateReference?view=graph-rest-beta",
		},
		"assignments": deviceAndAppManagementAssignment,
		"settings_json": schema.SetAttribute{
			// See INFO#2
			Required:    true,
			ElementType: wpschema.JSONStringType,
			Validators: []validator.Set{
				setvalidator.SizeBetween(1, 5000),
			},
			PlanModifiers:       []planmodifier.Set{&deviceManagementConfigurationPolicyJsonSettingsInstancePlanModifier{}},
			Description:         `settings`, // custom MS Graph attribute name
			MarkdownDescription: `Policy settings`,
		},
	},
	MarkdownDescription: "Device Management Configuration Policy / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-deviceManagementConfigurationPolicy?view=graph-rest-beta",
}

//
// See INFO#1
//

// Filters all properties that do not seem to be required or returned from MS Graph.
// This helps to keep settings_json consistent and comparable throughout roundtrips to MS Graph.
func deviceManagementConfigurationPolicyJsonSettingsInstanceFilter(path string, value any) (pass bool, err error) {
	sliceValue, valueIsSlice := value.([]any)
	pass = true && // formatting only
		!(strings.HasSuffix(path, "/settingInstanceTemplateReference") && value == nil) &&
		!(strings.HasSuffix(path, "/settingInstanceTemplateReference/@odata.type")) &&
		!(strings.HasSuffix(path, "/settingValueTemplateReference") && value == nil) &&
		!(strings.HasSuffix(path, "/settingValueTemplateReference/@odata.type")) &&
		!(strings.HasSuffix(path, "/choiceSettingValue/@odata.type")) &&
		!(strings.HasSuffix(path, "/@odata.type") && value == "#microsoft.graph.deviceManagementConfigurationGroupSettingValue") &&
		!(strings.HasSuffix(path, "/children@odata.type")) &&
		!(strings.HasSuffix(path, "/simpleSettingCollectionValue@odata.type")) &&
		!(strings.HasSuffix(path, "/groupSettingCollectionValue@odata.type")) &&
		!(strings.HasSuffix(path, "/children") && (value == nil || (valueIsSlice && len(sliceValue) == 0)))
	return
}

//
// See INFO#1, INFO#2
//

func deviceManagementConfigurationPolicyJsonTerraformToGraphMiddleware(params generic.TerraformToGraphMiddlewareParams) generic.TerraformToGraphMiddlewareReturns {

	type deviceManagementConfigurationPolicySettingModelJson struct {
		SettingInstance json.RawMessage `json:"settingInstance"`
	}

	settings := params.RawVal["settings"].([]any)
	for i := range settings {
		settings[i] = deviceManagementConfigurationPolicySettingModelJson{
			SettingInstance: json.RawMessage(settings[i].(string)),
		}
	}

	return nil
}

func deviceManagementConfigurationPolicyJsonGraphToTerraformMiddleware(params generic.GraphToTerraformMiddlewareParams) generic.GraphToTerraformMiddlewareReturns {

	settingsRaw := params.RawVal["settings"]
	if settingsRaw == nil {
		return nil // settings might have not been returned from Graph
	}
	settings, ok := settingsRaw.([]any)
	if !ok {
		return fmt.Errorf("unsupported type for settings: %T", settingsRaw)
	}

	for i := range settings {
		settingInstance := settings[i].(map[string]any)["settingInstance"].(map[string]any)

		// See INFO#1. Ugly hack relying on go json doing the sorting automatically internally
		newSettingInstance, err := wpschema.Traverse(settingInstance, deviceManagementConfigurationPolicyJsonSettingsInstanceFilter)
		if err != nil {
			return err
		}
		newSettingInstanceJson, err := json.Marshal(newSettingInstance)
		if err != nil {
			return err
		}

		settings[i] = string(newSettingInstanceJson)
	}

	return nil
}

//
// See INFO#1
//

type deviceManagementConfigurationPolicyJsonSettingsInstancePlanModifier struct{ generic.EmptyDescriber }

var _ planmodifier.Set = (*deviceManagementConfigurationPolicyJsonSettingsInstancePlanModifier)(nil)

func (apm *deviceManagementConfigurationPolicyJsonSettingsInstancePlanModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, res *planmodifier.SetResponse) {
	// do nothing on create or on destroy
	if req.StateValue.IsNull() || req.PlanValue.IsNull() {
		return
	}

	planElements := req.PlanValue.Elements()
	stateElements := req.StateValue.Elements()
	// do nothing (i.e. keep plan) if sizes differ
	if len(planElements) != len(stateElements) {
		return
	}

	planCleanStrings := apm.cleanAndSortElements(&res.Diagnostics, planElements, "planElements")
	stateCleanStrings := apm.cleanAndSortElements(&res.Diagnostics, stateElements, "stateElements")

	// updating only the really different elements in the plan would unfortunately make TF panic (see INFO#1 above)
	equalsEssentially := true
	for i := range planCleanStrings {
		if planCleanStrings[i] != stateCleanStrings[i] {
			tflog.Debug(ctx, fmt.Sprintf("planCleanStrings[%d]:  %s", i, planCleanStrings[i]))
			tflog.Debug(ctx, fmt.Sprintf("stateCleanStrings[%d]: %s", i, stateCleanStrings[i]))
			equalsEssentially = false
			break
		}
	}
	if equalsEssentially {
		res.PlanValue = req.StateValue
	}
}

func (apm *deviceManagementConfigurationPolicyJsonSettingsInstancePlanModifier) cleanAndSortElements(diags *diag.Diagnostics, elements []attr.Value, source string) []string {

	result := make([]string, len(elements))
	for i, element := range elements {
		elementClean, err := wpschema.TraverseJson([]byte(element.(wpschema.JSONString).ValueJSONString()), deviceManagementConfigurationPolicyJsonSettingsInstanceFilter)
		if err != nil {
			diags.AddError(fmt.Sprintf("TraverseJson(%s[%d])", source, i), err.Error())
			return nil
		}
		result = append(result, string(elementClean))
	}

	// Even though we use schema.SetAttribute we still need to sort our elements for comparison to make sure that both
	// lists are in the same order when WE compare them (above).
	sort.Strings(result)

	return result
}
