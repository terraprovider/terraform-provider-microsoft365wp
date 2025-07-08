package services

import (
	"strings"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpjsontypes"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	DeviceManagementConfigurationPolicyJsonResource = generic.GenericResource{
		TypeNameSuffix: "device_management_configuration_policy_json",
		SpecificSchema: deviceManagementConfigurationPolicyJsonResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/deviceManagement/configurationPolicies",
			ReadOptions: generic.ReadOptions{
				ODataExpand: "settings",
				ExtraRequests: []generic.ReadExtraRequest{
					{
						Attribute: "assignments",
					},
				},
			},
			WriteOptions: generic.WriteOptions{
				UsePutForUpdate: true,
				SubActions: []generic.WriteSubAction{
					&generic.WriteSubActionAllInOne{
						WriteSubActionBase: generic.WriteSubActionBase{
							Attributes: []string{"assignments"},
							UriSuffix:  "assign",
						},
					},
				},
			},
		},
	}

	DeviceManagementConfigurationPolicyJsonSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&DeviceManagementConfigurationPolicyJsonResource)

	DeviceManagementConfigurationPolicyJsonPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&DeviceManagementConfigurationPolicyJsonResource, "device_management_configuration_policies_json")
)

var deviceManagementConfigurationPolicyJsonResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // deviceManagementConfigurationPolicy
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Key of the policy document. Automatically generated.",
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
			Computed:            true,
			PlanModifiers:       []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
			MarkdownDescription: "Policy assignment status. This property is read-only.",
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
				wpvalidator.FlagEnumValues("none", "android", "iOS", "macOS", "windows10X", "windows10", "linux", "unknownFutureValue", "androidEnterprise", "aosp", "visionOS", "tvOS"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("windows10")},
			Computed:            true,
			MarkdownDescription: "Platforms for this policy. / Supported platform types; possible values are: `none` (Default. No platform type specified.), `android` (Settings for Android platform.), `iOS` (Settings for iOS platform.), `macOS` (Settings for MacOS platform.), `windows10X` (Windows 10 X.), `windows10` (Settings for Windows 10 platform.), `linux` (Settings for Linux platform.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.), `androidEnterprise` (Settings for Corporate Owned Android Enterprise devices.), `aosp` (Settings for Android Open Source Project platform.), `visionOS` (Settings for visionOS platform.), `tvOS` (Settings for tvOS platform.). The _provider_ default value is `\"windows10\"`.",
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: "List of Scope Tags for this Entity instance. The _provider_ default value is `[\"0\"]`.",
		},
		"setting_count": schema.Int64Attribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
			MarkdownDescription: "Number of settings",
		},
		"technologies": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				wpvalidator.FlagEnumValues("none", "mdm", "windows10XManagement", "configManager", "intuneManagementExtension", "thirdParty", "documentGateway", "appleRemoteManagement", "microsoftSense", "exchangeOnline", "mobileApplicationManagement", "linuxMdm", "enrollment", "endpointPrivilegeManagement", "unknownFutureValue", "windowsOsRecovery", "android"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("mdm")},
			Computed:            true,
			MarkdownDescription: "Technologies for this policy. / Describes which technology this setting can be deployed with; possible values are: `none` (Default. Indicates the setting cannot be deployed through any channel.), `mdm` (Indicates the settings that can be deployed through the Mobile Device Management (MDM) channel.), `windows10XManagement` (Indicates the settings that can be deployed through the Windows10XManagement channel.), `configManager` (Indicates the settings that can be deployed through the ConfigManager channel.), `intuneManagementExtension`, `thirdParty`, `documentGateway`, `appleRemoteManagement` (Indicates the settings that can be deployed through the AppleRemoteManagement channel.), `microsoftSense` (Indicates the settings that can be deployed through the SENSE agent channel.), `exchangeOnline` (Indicates the settings that can be deployed through the Exchange Online agent channel.), `mobileApplicationManagement` (Indicates the settings that can be deployed through the Mobile Application Management (MAM) channel.), `linuxMdm` (Indicates the settings that can be deployed through the Linux Mobile Device Management (MDM) channel.), `enrollment` (Indicates the settings that can be deployed through device enrollment.), `endpointPrivilegeManagement` (Indicates the settings that can be deployed through the Endpoint privilege management channel.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.), `windowsOsRecovery` (Indicates the settings that can be applied during a Windows operating system recovery. Such as accept license, base language, create recovery partition, format disk, keyboard settings.), `android` (Indicates the settings that can be deployed through the Android channel.). The _provider_ default value is `\"mdm\"`.",
		},
		"template_reference": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // deviceManagementConfigurationPolicyTemplateReference
				"template_display_name": schema.StringAttribute{
					Computed:            true,
					PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
					MarkdownDescription: "Template Display Name of the referenced template. This property is read-only.",
				},
				"template_display_version": schema.StringAttribute{
					Computed:            true,
					PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
					MarkdownDescription: "Template Display Version of the referenced Template. This property is read-only.",
				},
				"template_family": schema.StringAttribute{
					Computed: true,
					Validators: []validator.String{
						stringvalidator.OneOf("none", "endpointSecurityAntivirus", "endpointSecurityDiskEncryption", "endpointSecurityFirewall", "endpointSecurityEndpointDetectionAndResponse", "endpointSecurityAttackSurfaceReduction", "endpointSecurityAccountProtection", "endpointSecurityApplicationControl", "endpointSecurityEndpointPrivilegeManagement", "enrollmentConfiguration", "appQuietTime", "baseline", "unknownFutureValue", "deviceConfigurationScripts", "deviceConfigurationPolicies", "windowsOsRecoveryPolicies", "companyPortal"),
					},
					PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
					MarkdownDescription: "Template Family of the referenced Template. This property is read-only. / Describes the TemplateFamily for the Template entity; possible values are: `none` (Default for Template Family when Policy is not linked to a Template), `endpointSecurityAntivirus` (Template Family for EndpointSecurityAntivirus that manages the discrete group of antivirus settings for managed devices), `endpointSecurityDiskEncryption` (Template Family for EndpointSecurityDiskEncryption that provides settings that are relevant for a devices built-in encryption  method, like FileVault or BitLocker), `endpointSecurityFirewall` (Template Family for EndpointSecurityFirewall that helps configure a devices built-in firewall for device that run macOS and Windows 10), `endpointSecurityEndpointDetectionAndResponse` (Template Family for EndpointSecurityEndpointDetectionAndResponse that facilitates management of the EDR settings and onboard devices to Microsoft Defender for Endpoint), `endpointSecurityAttackSurfaceReduction` (Template Family for EndpointSecurityAttackSurfaceReduction that help reduce your attack surfaces, by minimizing the places where your organization is vulnerable to cyberthreats and attacks), `endpointSecurityAccountProtection` (Template Family for EndpointSecurityAccountProtection that facilitates protecting the identity and accounts of users), `endpointSecurityApplicationControl` (Template Family for ApplicationControl that helps mitigate security threats by restricting the applications that users can run and the code that runs in the System Core (kernel)), `endpointSecurityEndpointPrivilegeManagement` (Template Family for EPM Elevation Rules), `enrollmentConfiguration` (Template Family for EnrollmentConfiguration), `appQuietTime` (Template Family for QuietTimeIndicates Template Family for all the Apps QuietTime policies and templates), `baseline` (Template Family for Baseline), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.), `deviceConfigurationScripts` (Template Family for device configuration scripts), `deviceConfigurationPolicies` (Template Family for device configuration policies), `windowsOsRecoveryPolicies` (Template Family for windowsOsRecovery that can be applied during a Windows operating system recovery), `companyPortal` (Template Family for Company Portal settings)",
				},
				"template_id": schema.StringAttribute{
					Optional: true,
					PlanModifiers: []planmodifier.String{
						wpdefaultvalue.StringDefaultValue(""),
						stringplanmodifier.RequiresReplace(),
					},
					Computed:            true,
					MarkdownDescription: "Template id. The _provider_ default value is `\"\"`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Template reference information / Policy template reference information / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicytemplatereference?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
		"assignments": deviceAndAppManagementAssignment,
		"settings": schema.SetNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // deviceManagementConfigurationSetting
					"instance_json": schema.StringAttribute{
						Required:            true,
						CustomType:          wpjsontypes.NormalizedType{WpObjectFilterFunc: deviceManagementConfigurationPolicyJsonSettingsInstanceFilter},
						Description:         `settingInstance`, // custom MS Graph attribute name
						MarkdownDescription: "Setting Instance / Setting instance within policy / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationsettinginstance?view=graph-rest-beta",
					},
				},
			},
			Validators:          []validator.Set{setvalidator.SizeBetween(1, 5000)},
			MarkdownDescription: "Policy settings / Setting instance within policy / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationsetting?view=graph-rest-beta",
		},
	},
	MarkdownDescription: "Device Management Configuration Policy / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfigv2-devicemanagementconfigurationpolicy?view=graph-rest-beta ||| MS Graph: Device configuration",
}

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
