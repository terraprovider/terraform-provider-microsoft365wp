package services

import (
	"context"
	"strings"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	DeviceConfigurationResource = generic.GenericResource{
		TypeNameSuffix: "device_configuration",
		SpecificSchema: deviceConfigurationResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:         "/deviceManagement/deviceConfigurations",
			ReadOptions:     deviceConfigurationReadOptions,
			WriteSubActions: deviceConfigurationWriteSubActions,
		},
	}

	DeviceConfigurationSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&DeviceConfigurationResource)

	DeviceConfigurationPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&DeviceConfigurationResource, "")
)

var deviceConfigurationReadOptions = generic.ReadOptions{
	ODataFilter: "not(isof('microsoft.graph.windows10CustomConfiguration'))",
	ExtraRequests: []generic.ReadExtraRequest{
		{
			ParentAttribute: "assignments",
			UriSuffix:       "assignments",
		},
	},
}

var deviceConfigurationWriteSubActions = []generic.WriteSubAction{
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			Attributes: []string{"assignments"},
			UriSuffix:  "assign",
		},
	},
}

var deviceConfigurationResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // deviceConfiguration
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Key of the entity.",
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "DateTime the object was created.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Admin provided description of the Device Configuration.",
		},
		"device_management_applicability_rule_device_mode": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // deviceManagementApplicabilityRuleDeviceMode
				"device_mode": schema.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf("standardConfiguration", "sModeConfiguration"),
					},
					MarkdownDescription: "Applicability rule for device mode. / Windows 10 Device Mode type; possible values are: `standardConfiguration` (Standard Configuration), `sModeConfiguration` (S Mode Configuration)",
				},
				"name": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Name for object.",
				},
				"rule_type": schema.StringAttribute{
					Required:            true,
					Validators:          []validator.String{stringvalidator.OneOf("include", "exclude")},
					MarkdownDescription: "Applicability Rule type. / Supported Applicability rule types for Device Configuration; possible values are: `include` (Include), `exclude` (Exclude)",
				},
			},
			MarkdownDescription: "The device mode applicability rule for this Policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-devicemanagementapplicabilityruledevicemode?view=graph-rest-beta",
		},
		"device_management_applicability_rule_os_edition": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // deviceManagementApplicabilityRuleOsEdition
				"name": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Name for object.",
				},
				"os_edition_types": schema.SetAttribute{
					ElementType: types.StringType,
					Required:    true,
					Validators: []validator.Set{
						setvalidator.ValueStringsAre(
							stringvalidator.OneOf("windows10Enterprise", "windows10EnterpriseN", "windows10Education", "windows10EducationN", "windows10MobileEnterprise", "windows10HolographicEnterprise", "windows10Professional", "windows10ProfessionalN", "windows10ProfessionalEducation", "windows10ProfessionalEducationN", "windows10ProfessionalWorkstation", "windows10ProfessionalWorkstationN", "notConfigured", "windows10Home", "windows10HomeChina", "windows10HomeN", "windows10HomeSingleLanguage", "windows10Mobile", "windows10IoTCore", "windows10IoTCoreCommercial"),
						),
					},
					MarkdownDescription: "Applicability rule OS edition type. / Windows 10 Edition type; possible values are: `windows10Enterprise` (Windows 10 Enterprise), `windows10EnterpriseN` (Windows 10 EnterpriseN), `windows10Education` (Windows 10 Education), `windows10EducationN` (Windows 10 EducationN), `windows10MobileEnterprise` (Windows 10 Mobile Enterprise), `windows10HolographicEnterprise` (Windows 10 Holographic Enterprise), `windows10Professional` (Windows 10 Professional), `windows10ProfessionalN` (Windows 10 ProfessionalN), `windows10ProfessionalEducation` (Windows 10 Professional Education), `windows10ProfessionalEducationN` (Windows 10 Professional EducationN), `windows10ProfessionalWorkstation` (Windows 10 Professional for Workstations), `windows10ProfessionalWorkstationN` (Windows 10 Professional for Workstations N), `notConfigured` (NotConfigured), `windows10Home` (Windows 10 Home), `windows10HomeChina` (Windows 10 Home China), `windows10HomeN` (Windows 10 Home N), `windows10HomeSingleLanguage` (Windows 10 Home Single Language), `windows10Mobile` (Windows 10 Mobile), `windows10IoTCore` (Windows 10 IoT Core), `windows10IoTCoreCommercial` (Windows 10 IoT Core Commercial)",
				},
				"rule_type": schema.StringAttribute{
					Required:            true,
					Validators:          []validator.String{stringvalidator.OneOf("include", "exclude")},
					MarkdownDescription: "Applicability Rule type. / Supported Applicability rule types for Device Configuration; possible values are: `include` (Include), `exclude` (Exclude)",
				},
			},
			MarkdownDescription: "The OS edition applicability for this Policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-devicemanagementapplicabilityruleosedition?view=graph-rest-beta",
		},
		"device_management_applicability_rule_os_version": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // deviceManagementApplicabilityRuleOsVersion
				"max_os_version": schema.StringAttribute{
					Optional:            true,
					Description:         `maxOSVersion`, // custom MS Graph attribute name
					MarkdownDescription: "Max OS version for Applicability Rule.",
				},
				"min_os_version": schema.StringAttribute{
					Optional:            true,
					Description:         `minOSVersion`, // custom MS Graph attribute name
					MarkdownDescription: "Min OS version for Applicability Rule.",
				},
				"name": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "Name for object.",
				},
				"rule_type": schema.StringAttribute{
					Required:            true,
					Validators:          []validator.String{stringvalidator.OneOf("include", "exclude")},
					MarkdownDescription: "Applicability Rule type. / Supported Applicability rule types for Device Configuration; possible values are: `include` (Include), `exclude` (Exclude)",
				},
			},
			MarkdownDescription: "The OS version applicability rule for this Policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-devicemanagementapplicabilityruleosversion?view=graph-rest-beta",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Admin provided name of the device configuration.",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "DateTime the object was last modified.",
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: "List of Scope Tags for this Entity instance. The _provider_ default value is `[\"0\"]`.",
		},
		"supports_scope_tags": schema.BoolAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
			MarkdownDescription: "Indicates whether or not the underlying Device Configuration supports the assignment of scope tags. Assigning to the ScopeTags property is not allowed when this value is false and entities will not be visible to scoped users. This occurs for Legacy policies created in Silverlight and can be resolved by deleting and recreating the policy in the Azure Portal. This property is read-only.",
		},
		"version": schema.Int64Attribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
			MarkdownDescription: "Version of the device configuration.",
		},
		"assignments": deviceAndAppManagementAssignment,
		"android_device_owner_general_device": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.androidDeviceOwnerGeneralDeviceConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // androidDeviceOwnerGeneralDeviceConfiguration
					"accounts_block_modification": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not adding or removing accounts is disabled.",
					},
					"android_device_owner_delegated_scope_app_settings": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // androidDeviceOwnerDelegatedScopeAppSetting
								"app_detail": schema.SingleNestedAttribute{
									Optional:            true,
									Attributes:          deviceConfigurationAppListItemAttributes,
									PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: "Information about the app like Name, AppStoreUrl, Publisher and AppId / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `{}`.",
								},
								"app_scopes": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Validators: []validator.Set{
										setvalidator.ValueStringsAre(
											stringvalidator.OneOf("unspecified", "certificateInstall", "captureNetworkActivityLog", "captureSecurityLog", "unknownFutureValue"),
										),
									},
									PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: "List of scopes an app has been assigned. / An enum representing possible values for delegated app scope; possible values are: `unspecified` (Unspecified; this value defaults to DELEGATED_SCOPE_UNSPECIFIED.), `certificateInstall` (Specifies that the admin has given app permission to install and manage certificates on device.), `captureNetworkActivityLog` (Specifies that the admin has given app permission to capture network activity logs on device. More info on Network activity logs: https://developer.android.com/work/dpc/logging), `captureSecurityLog` (Specified that the admin has given permission to capture security logs on device. More info on Security logs: https://developer.android.com/work/dpc/security#log_enterprise_device_activity), `unknownFutureValue` (Unknown future value (reserved, not used right now)). The _provider_ default value is `[]`.",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Specifies the list of managed apps with app details and its associated delegated scope(s). This collection can contain a maximum of 500 elements. / Represents one item in the list of managed apps with app details and its associated delegated scope(s). / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceownerdelegatedscopeappsetting?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"apps_allow_install_from_unknown_sources": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not the user is allowed to enable to unknown sources setting.",
					},
					"apps_auto_update_policy": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "userChoice", "never", "wiFiOnly", "always"),
						},
						MarkdownDescription: "Indicates the value of the app auto update policy. / Android Device Owner possible values for states of the device's app auto update policy; possible values are: `notConfigured` (Not configured; this value is ignored.), `userChoice` (The user can control auto-updates.), `never` (Apps are never auto-updated.), `wiFiOnly` (Apps are auto-updated over Wi-Fi only.), `always` (Apps are auto-updated at any time. Data charges may apply.)",
					},
					"apps_default_permission_policy": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "prompt", "autoGrant", "autoDeny"),
						},
						MarkdownDescription: "Indicates the permission policy for requests for runtime permissions if one is not defined for the app specifically. / Android Device Owner default app permission policy type; possible values are: `deviceDefault` (Device default value, no intent.), `prompt` (Prompt.), `autoGrant` (Auto grant.), `autoDeny` (Auto deny.)",
					},
					"apps_recommend_skipping_first_use_hints": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not to recommend all apps skip any first-time-use hints they may have added.",
					},
					"azure_ad_shared_device_data_clear_apps": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationAppListItemAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "A list of managed apps that will have their data cleared during a global sign-out in AAD shared device mode. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"bluetooth_block_configuration": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block a user from configuring bluetooth.",
					},
					"bluetooth_block_contact_sharing": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block a user from sharing contacts via bluetooth.",
					},
					"camera_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to disable the use of the camera.",
					},
					"cellular_block_wifi_tethering": schema.BoolAttribute{
						Optional:            true,
						Description:         `cellularBlockWiFiTethering`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to block Wi-Fi tethering.",
					},
					"certificate_credential_configuration_disabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block users from any certificate credential configuration.",
					},
					"cross_profile_policies_allow_copy_paste": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not text copied from one profile (personal or work) can be pasted in the other.",
					},
					"cross_profile_policies_allow_data_sharing": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "crossProfileDataSharingBlocked", "dataSharingFromWorkToPersonalBlocked", "crossProfileDataSharingAllowed", "unkownFutureValue"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Indicates whether data from one profile (personal or work) can be shared with apps in the other profile. / An enum representing possible values for cross profile data sharing; possible values are: `notConfigured` (Not configured; this value defaults to CROSS_PROFILE_DATA_SHARING_UNSPECIFIED.), `crossProfileDataSharingBlocked` (Data cannot be shared from both the personal profile to work profile and the work profile to the personal profile.), `dataSharingFromWorkToPersonalBlocked` (Prevents users from sharing data from the work profile to apps in the personal profile. Personal data can be shared with work apps.), `crossProfileDataSharingAllowed` (Data from either profile can be shared with the other profile.), `unkownFutureValue` (Unknown future value (reserved, not used right now)). The _provider_ default value is `\"notConfigured\"`.",
					},
					"cross_profile_policies_show_work_contacts_in_personal_profile": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not contacts stored in work profile are shown in personal profile contact searches/incoming calls.",
					},
					"data_roaming_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block a user from data roaming.",
					},
					"date_time_configuration_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block the user from manually changing the date or time on the device",
					},
					"detailed_help_text": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // androidDeviceOwnerUserFacingMessage
							"default_message": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "The default message displayed if the user's locale doesn't match with any of the localized messages",
							},
							"localized_messages": schema.SetNestedAttribute{
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: deviceConfigurationKeyValuePairAttributes,
								},
								PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
								Computed:            true,
								MarkdownDescription: "The list of <locale, message> pairs. This collection can contain a maximum of 500 elements. / Key value pair for storing custom settings / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`.",
							},
						},
						MarkdownDescription: "Represents the customized detailed help text provided to users when they attempt to modify managed settings on their device. / Represents a user-facing message with locale information as well as a default message to be used if the user's locale doesn't match with any of the localized messages / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceowneruserfacingmessage?view=graph-rest-beta",
					},
					"device_location_mode": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "disabled", "unknownFutureValue"),
						},
						MarkdownDescription: "Indicates the location setting configuration for fully managed devices (COBO) and corporate owned devices with a work profile (COPE). / Android Device Owner Location Mode Type; possible values are: `notConfigured` (No restrictions on the location setting and no specific behavior is set or enforced. This is the default), `disabled` (Location setting is disabled on the device), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use)",
					},
					"device_owner_lock_screen_message": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // androidDeviceOwnerUserFacingMessage
							"default_message": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "The default message displayed if the user's locale doesn't match with any of the localized messages",
							},
							"localized_messages": schema.SetNestedAttribute{
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: deviceConfigurationKeyValuePairAttributes,
								},
								PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
								Computed:            true,
								MarkdownDescription: "The list of <locale, message> pairs. This collection can contain a maximum of 500 elements. / Key value pair for storing custom settings / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`.",
							},
						},
						MarkdownDescription: "Represents the customized lock screen message provided to users when they attempt to modify managed settings on their device. / Represents a user-facing message with locale information as well as a default message to be used if the user's locale doesn't match with any of the localized messages / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceowneruserfacingmessage?view=graph-rest-beta",
					},
					"enrollment_profile": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "dedicatedDevice", "fullyManaged"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Indicates which enrollment profile you want to configure. / Android Device Owner Enrollment Profile types; possible values are: `notConfigured` (Not configured; this value is ignored.), `dedicatedDevice` (Dedicated device.), `fullyManaged` (Fully managed.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"factory_reset_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not the factory reset option in settings is disabled.",
					},
					"factory_reset_device_administrator_emails": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "List of Google account emails that will be required to authenticate after a device is factory reset before it can be set up. The _provider_ default value is `[]`.",
					},
					"global_proxy": schema.SingleNestedAttribute{
						Optional:   true,
						Attributes: map[string]schema.Attribute{ // androidDeviceOwnerGlobalProxy
						},
						MarkdownDescription: "Proxy is set up directly with host, port and excluded hosts. / Android Device Owner Global Proxy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceownerglobalproxy?view=graph-rest-beta",
					},
					"google_accounts_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not google accounts will be blocked.",
					},
					"kiosk_customization_device_settings_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether a user can access the device's Settings app while in Kiosk Mode.",
					},
					"kiosk_customization_power_button_actions_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether the power menu is shown when a user long presses the Power button of a device in Kiosk Mode.",
					},
					"kiosk_customization_status_bar": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "notificationsAndSystemInfoEnabled", "systemInfoOnly"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Indicates whether system info and notifications are disabled in Kiosk Mode. / An enum representing possible values for kiosk customization system navigation; possible values are: `notConfigured` (Not configured; this value defaults to STATUS_BAR_UNSPECIFIED.), `notificationsAndSystemInfoEnabled` (System info and notifications are shown on the status bar in kiosk mode.), `systemInfoOnly` (Only system info is shown on the status bar in kiosk mode.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"kiosk_customization_system_error_warnings": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether system error dialogs for crashed or unresponsive apps are shown in Kiosk Mode.",
					},
					"kiosk_customization_system_navigation": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "navigationEnabled", "homeButtonOnly"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Indicates which navigation features are enabled in Kiosk Mode. / An enum representing possible values for kiosk customization system navigation; possible values are: `notConfigured` (Not configured; this value defaults to NAVIGATION_DISABLED.), `navigationEnabled` (Home and overview buttons are enabled.), `homeButtonOnly` (Only the home button is enabled.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"kiosk_mode_app_order_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not to enable app ordering in Kiosk Mode.",
					},
					"kiosk_mode_app_positions": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // androidDeviceOwnerKioskModeAppPositionItem
								"item": schema.SingleNestedAttribute{
									Optional:   true,
									Attributes: map[string]schema.Attribute{ // androidDeviceOwnerKioskModeHomeScreenItem
									},
									PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: "Item to be arranged / Represents an item on the Android Device Owner Managed Home Screen (application, weblink or folder / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceownerkioskmodehomescreenitem?view=graph-rest-beta. The _provider_ default value is `{}`.",
								},
								"position": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: "Position of the item on the grid. Valid values 0 to 9999999. The _provider_ default value is `0`.",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "The ordering of items on Kiosk Mode Managed Home Screen. This collection can contain a maximum of 500 elements. / An item in the list of app positions that sets the order of items on the Managed Home Screen / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceownerkioskmodeapppositionitem?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"kiosk_mode_apps": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationAppListItemAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "A list of managed apps that will be shown when the device is in Kiosk Mode. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"kiosk_mode_apps_in_folder_ordered_by_name": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not to alphabetize applications within a folder in Kiosk Mode.",
					},
					"kiosk_mode_bluetooth_configuration_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not to allow a user to configure Bluetooth settings in Kiosk Mode.",
					},
					"kiosk_mode_debug_menu_easy_access_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not to allow a user to easy access to the debug menu in Kiosk Mode.",
					},
					"kiosk_mode_exit_code": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Exit code to allow a user to escape from Kiosk Mode when the device is in Kiosk Mode.",
					},
					"kiosk_mode_flashlight_configuration_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not to allow a user to use the flashlight in Kiosk Mode.",
					},
					"kiosk_mode_folder_icon": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "darkSquare", "darkCircle", "lightSquare", "lightCircle"),
						},
						MarkdownDescription: "Folder icon configuration for managed home screen in Kiosk Mode. / Android Device Owner Kiosk Mode folder icon type; possible values are: `notConfigured` (Not configured; this value is ignored.), `darkSquare` (Folder icon appears as dark square.), `darkCircle` (Folder icon appears as dark circle.), `lightSquare` (Folder icon appears as light square.), `lightCircle` (Folder icon appears as light circle  .)",
					},
					"kiosk_mode_grid_height": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of rows for Managed Home Screen grid with app ordering enabled in Kiosk Mode. Valid values 1 to 9999999",
					},
					"kiosk_mode_grid_width": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of columns for Managed Home Screen grid with app ordering enabled in Kiosk Mode. Valid values 1 to 9999999",
					},
					"kiosk_mode_icon_size": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "smallest", "small", "regular", "large", "largest"),
						},
						MarkdownDescription: "Icon size configuration for managed home screen in Kiosk Mode. / Android Device Owner Kiosk Mode managed home screen icon size; possible values are: `notConfigured` (Not configured; this value is ignored.), `smallest` (Smallest icon size.), `small` (Small icon size.), `regular` (Regular icon size.), `large` (Large icon size.), `largest` (Largest icon size.)",
					},
					"kiosk_mode_lock_home_screen": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not to lock home screen to the end user in Kiosk Mode.",
					},
					"kiosk_mode_managed_folders": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // androidDeviceOwnerKioskModeManagedFolder
								"folder_identifier": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "Unique identifier for the folder",
								},
								"folder_name": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "Display name for the folder",
								},
								"items": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{ // androidDeviceOwnerKioskModeFolderItem
										},
									},
									PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: "Items to be added to managed folder. This collection can contain a maximum of 500 elements. / Represents an item that can be added to Android Device Owner folder (application or weblink) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceownerkioskmodefolderitem?view=graph-rest-beta. The _provider_ default value is `[]`.",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "A list of managed folders for a device in Kiosk Mode. This collection can contain a maximum of 500 elements. / A folder containing pages of apps and weblinks on the Managed Home Screen / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceownerkioskmodemanagedfolder?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"kiosk_mode_managed_home_screen_auto_signout": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not to automatically sign-out of MHS and Shared device mode applications after inactive for Managed Home Screen.",
					},
					"kiosk_mode_managed_home_screen_inactive_sign_out_delay_in_seconds": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of seconds to give user notice before automatically signing them out for Managed Home Screen. Valid values 0 to 9999999",
					},
					"kiosk_mode_managed_home_screen_inactive_sign_out_notice_in_seconds": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of seconds device is inactive before automatically signing user out for Managed Home Screen. Valid values 0 to 9999999",
					},
					"kiosk_mode_managed_home_screen_pin_complexity": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "simple", "complex"),
						},
						MarkdownDescription: "Complexity of PIN for sign-in session for Managed Home Screen. / Complexity of PIN for Managed Home Screen sign-in session; possible values are: `notConfigured` (Not configured.), `simple` (Numeric values only.), `complex` (Alphanumerical value.)",
					},
					"kiosk_mode_managed_home_screen_pin_required": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not require user to set a PIN for sign-in session for Managed Home Screen.",
					},
					"kiosk_mode_managed_home_screen_pin_required_to_resume": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not required user to enter session PIN if screensaver has appeared for Managed Home Screen.",
					},
					"kiosk_mode_managed_home_screen_sign_in_background": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Custom URL background for sign-in screen for Managed Home Screen.",
					},
					"kiosk_mode_managed_home_screen_sign_in_branding_logo": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Custom URL branding logo for sign-in screen and session pin page for Managed Home Screen.",
					},
					"kiosk_mode_managed_home_screen_sign_in_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not show sign-in screen for Managed Home Screen.",
					},
					"kiosk_mode_managed_settings_entry_disabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not to display the Managed Settings entry point on the managed home screen in Kiosk Mode.",
					},
					"kiosk_mode_media_volume_configuration_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not to allow a user to change the media volume in Kiosk Mode.",
					},
					"kiosk_mode_screen_orientation": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "portrait", "landscape", "autoRotate"),
						},
						MarkdownDescription: "Screen orientation configuration for managed home screen in Kiosk Mode. / Android Device Owner Kiosk Mode managed home screen orientation; possible values are: `notConfigured` (Not configured; this value is ignored.), `portrait` (Portrait orientation.), `landscape` (Landscape orientation.), `autoRotate` (Auto rotate between portrait and landscape orientations.)",
					},
					"kiosk_mode_screen_saver_configuration_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not to enable screen saver mode or not in Kiosk Mode.",
					},
					"kiosk_mode_screen_saver_detect_media_disabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not the device screen should show the screen saver if audio/video is playing in Kiosk Mode.",
					},
					"kiosk_mode_screen_saver_display_time_in_seconds": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The number of seconds that the device will display the screen saver for in Kiosk Mode. Valid values 0 to 9999999",
					},
					"kiosk_mode_screen_saver_image_url": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "URL for an image that will be the device's screen saver in Kiosk Mode.",
					},
					"kiosk_mode_screen_saver_start_delay_in_seconds": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The number of seconds the device needs to be inactive for before the screen saver is shown in Kiosk Mode. Valid values 1 to 9999999",
					},
					"kiosk_mode_show_app_notification_badge": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not to display application notification badges in Kiosk Mode.",
					},
					"kiosk_mode_show_device_info": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not to allow a user to access basic device information.",
					},
					"kiosk_mode_use_managed_home_screen_app": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "singleAppMode", "multiAppMode"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Whether or not to use single app kiosk mode or multi-app kiosk mode. / Possible values of Android Kiosk Mode; possible values are: `notConfigured` (Not configured), `singleAppMode` (Run in single-app mode), `multiAppMode` (Run in multi-app mode). The _provider_ default value is `\"notConfigured\"`.",
					},
					"kiosk_mode_virtual_home_button_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not to display a virtual home button when the device is in Kiosk Mode.",
					},
					"kiosk_mode_virtual_home_button_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "swipeUp", "floating"),
						},
						MarkdownDescription: "Indicates whether the virtual home button is a swipe up home button or a floating home button. / Android Device Owner Kiosk Mode managed home screen virtual home button type; possible values are: `notConfigured` (Not configured; this value is ignored.), `swipeUp` (Swipe-up for home button.), `floating` (Floating home button.)",
					},
					"kiosk_mode_wallpaper_url": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "URL to a publicly accessible image to use for the wallpaper when the device is in Kiosk Mode.",
					},
					"kiosk_mode_wifi_allowed_ssids": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "The restricted set of WIFI SSIDs available for the user to configure in Kiosk Mode. This collection can contain a maximum of 500 elements. The _provider_ default value is `[]`.",
					},
					"kiosk_mode_wifi_configuration_enabled": schema.BoolAttribute{
						Optional:            true,
						Description:         `kioskModeWiFiConfigurationEnabled`, // custom MS Graph attribute name
						MarkdownDescription: "Whether or not to allow a user to configure Wi-Fi settings in Kiosk Mode.",
					},
					"locate_device_lost_mode_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not LocateDevice for devices with lost mode (COBO, COPE) is enabled.",
					},
					"locate_device_userless_disabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not LocateDevice for userless (COSU) devices is disabled.",
					},
					"microphone_force_mute": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block unmuting the microphone on the device.",
					},
					"microsoft_launcher_configuration_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to you want configure Microsoft Launcher.",
					},
					"microsoft_launcher_custom_wallpaper_allow_user_modification": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not the user can modify the wallpaper to personalize their device.",
					},
					"microsoft_launcher_custom_wallpaper_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to configure the wallpaper on the targeted devices.",
					},
					"microsoft_launcher_custom_wallpaper_image_url": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates the URL for the image file to use as the wallpaper on the targeted devices.",
					},
					"microsoft_launcher_dock_presence_allow_user_modification": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not the user can modify the device dock configuration on the device.",
					},
					"microsoft_launcher_dock_presence_configuration": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "show", "hide", "disabled"),
						},
						MarkdownDescription: "Indicates whether or not you want to configure the device dock. / Microsoft Launcher Dock Presence selection; possible values are: `notConfigured` (Not configured; this value is ignored.), `show` (Indicates the device's dock will be displayed on the device.), `hide` (Indicates the device's dock will be hidden on the device, but the user can access the dock by dragging the handler on the bottom of the screen.), `disabled` (Indicates the device's dock will be disabled on the device.)",
					},
					"microsoft_launcher_feed_allow_user_modification": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not the user can modify the launcher feed on the device.",
					},
					"microsoft_launcher_feed_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not you want to enable the launcher feed on the device.",
					},
					"microsoft_launcher_search_bar_placement_configuration": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "top", "bottom", "hide"),
						},
						MarkdownDescription: "Indicates the search bar placement configuration on the device. / Microsoft Launcher Search Bar Placement selection; possible values are: `notConfigured` (Not configured; this value is ignored.), `top` (Indicates that the search bar will be displayed on the top of the device.), `bottom` (Indicates that the search bar will be displayed on the bottom of the device.), `hide` (Indicates that the search bar will be hidden on the device.)",
					},
					"network_escape_hatch_allowed": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not the device will allow connecting to a temporary network connection at boot time.",
					},
					"nfc_block_outgoing_beam": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block NFC outgoing beam.",
					},
					"password_block_keyguard": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not the keyguard is disabled.",
					},
					"password_block_keyguard_features": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(
								stringvalidator.OneOf("notConfigured", "camera", "notifications", "unredactedNotifications", "trustAgents", "fingerprint", "remoteInput", "allFeatures", "face", "iris", "biometrics"),
							),
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "List of device keyguard features to block. This collection can contain a maximum of 11 elements. / Android keyguard feature; possible values are: `notConfigured` (Not configured; this value is ignored.), `camera` (Camera usage when on secure keyguard screens.), `notifications` (Showing notifications when on secure keyguard screens.), `unredactedNotifications` (Showing unredacted notifications when on secure keyguard screens.), `trustAgents` (Trust agent state when on secure keyguard screens.), `fingerprint` (Fingerprint sensor usage when on secure keyguard screens.), `remoteInput` (Notification text entry when on secure keyguard screens.), `allFeatures` (All keyguard features when on secure keyguard screens.), `face` (Face authentication on secure keyguard screens.), `iris` (Iris authentication on secure keyguard screens.), `biometrics` (All biometric authentication on secure keyguard screens.). The _provider_ default value is `[]`.",
					},
					"password_expiration_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the amount of time that a password can be set for before it expires and a new password will be required. Valid values 1 to 365",
					},
					"password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the minimum length of the password required on the device. Valid values 4 to 16",
					},
					"password_minimum_letter_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the minimum number of letter characters required for device password. Valid values 1 to 16",
					},
					"password_minimum_lower_case_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the minimum number of lower case characters required for device password. Valid values 1 to 16",
					},
					"password_minimum_non_letter_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the minimum number of non-letter characters required for device password. Valid values 1 to 16",
					},
					"password_minimum_numeric_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the minimum number of numeric characters required for device password. Valid values 1 to 16",
					},
					"password_minimum_symbol_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the minimum number of symbol characters required for device password. Valid values 1 to 16",
					},
					"password_minimum_upper_case_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the minimum number of upper case letter characters required for device password. Valid values 1 to 16",
					},
					"password_minutes_of_inactivity_before_screen_timeout": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before the screen times out.",
					},
					"password_previous_password_count_to_block": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the length of password history, where the user will not be able to enter a new password that is the same as any password in the history. Valid values 0 to 24",
					},
					"password_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "required", "numeric", "numericComplex", "alphabetic", "alphanumeric", "alphanumericWithSymbols", "lowSecurityBiometric", "customPassword"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: "Indicates the minimum password quality required on the device. / Android Device Owner policy required password type; possible values are: `deviceDefault` (Device default value, no intent.), `required` (There must be a password set, but there are no restrictions on type.), `numeric` (At least numeric.), `numericComplex` (At least numeric with no repeating or ordered sequences.), `alphabetic` (At least alphabetic password.), `alphanumeric` (At least alphanumeric password), `alphanumericWithSymbols` (At least alphanumeric with symbols.), `lowSecurityBiometric` (Low security biometrics based password required.), `customPassword` (Custom password set by the admin.). The _provider_ default value is `\"deviceDefault\"`.",
					},
					"password_require_unlock": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "daily", "unkownFutureValue"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: "Indicates the timeout period after which a device must be unlocked using a form of strong authentication. / An enum representing possible values for required password unlock; possible values are: `deviceDefault` (Timeout period before strong authentication is required is set to the device's default.), `daily` (Timeout period before strong authentication is required is set to 24 hours.), `unkownFutureValue` (Unknown future value (reserved, not used right now)). The _provider_ default value is `\"deviceDefault\"`.",
					},
					"password_sign_in_failure_count_before_factory_reset": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the number of times a user can enter an incorrect password before the device is wiped. Valid values 4 to 11",
					},
					"personal_profile_apps_allow_install_from_unknown_sources": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether the user can install apps from unknown sources on the personal profile.",
					},
					"personal_profile_camera_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether to disable the use of the camera on the personal profile.",
					},
					"personal_profile_personal_applications": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationAppListItemAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Policy applied to applications in the personal profile. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"personal_profile_play_store_mode": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "blockedApps", "allowedApps"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Used together with PersonalProfilePersonalApplications to control how apps in the personal profile are allowed or blocked. / Used together with personalApplications to control how apps in the personal profile are allowed or blocked; possible values are: `notConfigured` (Not configured.), `blockedApps` (Blocked Apps.), `allowedApps` (Allowed Apps.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"personal_profile_screen_capture_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether to disable the capability to take screenshots on the personal profile.",
					},
					"play_store_mode": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "allowList", "blockList"),
						},
						MarkdownDescription: "Indicates the Play Store mode of the device. / Android Device Owner Play Store mode type; possible values are: `notConfigured` (Not Configured), `allowList` (Only apps that are in the policy are available and any app not in the policy will be automatically uninstalled from the device.), `blockList` (All apps are available and any app that should not be on the device should be explicitly marked as 'BLOCKED' in the applications policy.)",
					},
					"screen_capture_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to disable the capability to take screenshots.",
					},
					"security_common_criteria_mode_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Represents the security common criteria mode enabled provided to users when they attempt to modify managed settings on their device.",
					},
					"security_developer_settings_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not the user is allowed to access developer settings like developer options and safe boot on the device.",
					},
					"security_require_verify_apps": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not verify apps is required. The _provider_ default value is `true`.",
					},
					"share_device_location_disabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not location sharing is disabled for fully managed devices (COBO), and corporate owned devices with a work profile (COPE)",
					},
					"short_help_text": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // androidDeviceOwnerUserFacingMessage
							"default_message": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "The default message displayed if the user's locale doesn't match with any of the localized messages",
							},
							"localized_messages": schema.SetNestedAttribute{
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: deviceConfigurationKeyValuePairAttributes,
								},
								PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
								Computed:            true,
								MarkdownDescription: "The list of <locale, message> pairs. This collection can contain a maximum of 500 elements. / Key value pair for storing custom settings / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`.",
							},
						},
						MarkdownDescription: "Represents the customized short help text provided to users when they attempt to modify managed settings on their device. / Represents a user-facing message with locale information as well as a default message to be used if the user's locale doesn't match with any of the localized messages / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceowneruserfacingmessage?view=graph-rest-beta",
					},
					"status_bar_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or the status bar is disabled, including notifications, quick settings and other screen overlays.",
					},
					"stay_on_modes": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(
								stringvalidator.OneOf("notConfigured", "ac", "usb", "wireless"),
							),
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "List of modes in which the device's display will stay powered-on. This collection can contain a maximum of 4 elements. / Android Device Owner possible values for states of the device's plugged-in power modes; possible values are: `notConfigured` (Not configured; this value is ignored.), `ac` (Power source is an AC charger.), `usb` (Power source is a USB port.), `wireless` (Power source is wireless.). The _provider_ default value is `[]`.",
					},
					"storage_allow_usb": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to allow USB mass storage.",
					},
					"storage_block_external_media": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block external media.",
					},
					"storage_block_usb_file_transfer": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block USB file transfer.",
					},
					"system_update_freeze_periods": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // androidDeviceOwnerSystemUpdateFreezePeriod
								"end_day": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: "The day of the end date of the freeze period. Valid values 1 to 31. The _provider_ default value is `0`.",
								},
								"end_month": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: "The month of the end date of the freeze period. Valid values 1 to 12. The _provider_ default value is `0`.",
								},
								"start_day": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: "The day of the start date of the freeze period. Valid values 1 to 31. The _provider_ default value is `0`.",
								},
								"start_month": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: "The month of the start date of the freeze period. Valid values 1 to 12. The _provider_ default value is `0`.",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Indicates the annually repeating time periods during which system updates are postponed. This collection can contain a maximum of 500 elements. / Represents one item in the list of freeze periods for Android Device Owner system updates / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceownersystemupdatefreezeperiod?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"system_update_install_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "postpone", "windowed", "automatic"),
						},
						MarkdownDescription: "The type of system update configuration. / System Update Types for Android Device Owner; possible values are: `deviceDefault` (Device default behavior, which typically prompts the user to accept system updates.), `postpone` (Postpone automatic install of updates up to 30 days.), `windowed` (Install automatically inside a daily maintenance window.), `automatic` (Automatically install updates as soon as possible.)",
					},
					"system_update_window_end_minutes_after_midnight": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the number of minutes after midnight that the system update window ends. Valid values 0 to 1440",
					},
					"system_update_window_start_minutes_after_midnight": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the number of minutes after midnight that the system update window starts. Valid values 0 to 1440",
					},
					"system_windows_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not to block Android system prompt windows, like toasts, phone activities, and system alerts.",
					},
					"users_block_add": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not adding users and profiles is disabled.",
					},
					"users_block_remove": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to disable removing other users from the device.",
					},
					"volume_block_adjustment": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not adjusting the master volume is disabled.",
					},
					"vpn_always_on_lockdown_mode": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "If an always on VPN package name is specified, whether or not to lock network traffic when that VPN is disconnected. The _provider_ default value is `false`.",
					},
					"vpn_always_on_package_identifier": schema.StringAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
						Computed:            true,
						MarkdownDescription: "Android app package name for app that will handle an always-on VPN connection. The _provider_ default value is `\"\"`.",
					},
					"wifi_block_edit_configurations": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block the user from editing the wifi connection settings.",
					},
					"wifi_block_edit_policy_defined_configurations": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block the user from editing just the networks defined by the policy.",
					},
					"work_profile_password_expiration_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the number of days that a work profile password can be set before it expires and a new password will be required. Valid values 1 to 365",
					},
					"work_profile_password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the minimum length of the work profile password. Valid values 4 to 16",
					},
					"work_profile_password_minimum_letter_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the minimum number of letter characters required for the work profile password. Valid values 1 to 16",
					},
					"work_profile_password_minimum_lower_case_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the minimum number of lower-case characters required for the work profile password. Valid values 1 to 16",
					},
					"work_profile_password_minimum_non_letter_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the minimum number of non-letter characters required for the work profile password. Valid values 1 to 16",
					},
					"work_profile_password_minimum_numeric_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the minimum number of numeric characters required for the work profile password. Valid values 1 to 16",
					},
					"work_profile_password_minimum_symbol_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the minimum number of symbol characters required for the work profile password. Valid values 1 to 16",
					},
					"work_profile_password_minimum_upper_case_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the minimum number of upper-case letter characters required for the work profile password. Valid values 1 to 16",
					},
					"work_profile_password_previous_password_count_to_block": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the length of the work profile password history, where the user will not be able to enter a new password that is the same as any password in the history. Valid values 0 to 24",
					},
					"work_profile_password_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "required", "numeric", "numericComplex", "alphabetic", "alphanumeric", "alphanumericWithSymbols", "lowSecurityBiometric", "customPassword"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: "Indicates the minimum password quality required on the work profile password. / Android Device Owner policy required password type; possible values are: `deviceDefault` (Device default value, no intent.), `required` (There must be a password set, but there are no restrictions on type.), `numeric` (At least numeric.), `numericComplex` (At least numeric with no repeating or ordered sequences.), `alphabetic` (At least alphabetic password.), `alphanumeric` (At least alphanumeric password), `alphanumericWithSymbols` (At least alphanumeric with symbols.), `lowSecurityBiometric` (Low security biometrics based password required.), `customPassword` (Custom password set by the admin.). The _provider_ default value is `\"deviceDefault\"`.",
					},
					"work_profile_password_require_unlock": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "daily", "unkownFutureValue"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: "Indicates the timeout period after which a work profile must be unlocked using a form of strong authentication. / An enum representing possible values for required password unlock; possible values are: `deviceDefault` (Timeout period before strong authentication is required is set to the device's default.), `daily` (Timeout period before strong authentication is required is set to 24 hours.), `unkownFutureValue` (Unknown future value (reserved, not used right now)). The _provider_ default value is `\"deviceDefault\"`.",
					},
					"work_profile_password_sign_in_failure_count_before_factory_reset": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Indicates the number of times a user can enter an incorrect work profile password before the device is wiped. Valid values 4 to 11",
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: "This topic provides descriptions of the declared methods, properties and relationships exposed by the androidDeviceOwnerGeneralDeviceConfiguration resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceownergeneraldeviceconfiguration?view=graph-rest-beta",
			},
		},
		"android_work_profile_general_device": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.androidWorkProfileGeneralDeviceConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // androidWorkProfileGeneralDeviceConfiguration
					"allowed_google_account_domains": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Determine domains allow-list for accounts that can be added to work profile. The _provider_ default value is `[]`.",
					},
					"block_unified_password_for_work_profile": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Prevent using unified password for unlocking device and work profile. The _provider_ default value is `false`.",
					},
					"password_block_face_unlock": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block face unlock. The _provider_ default value is `false`.",
					},
					"password_block_fingerprint_unlock": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block fingerprint unlock. The _provider_ default value is `false`.",
					},
					"password_block_iris_unlock": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block iris unlock. The _provider_ default value is `false`.",
					},
					"password_block_trust_agents": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block Smart Lock and other trust agents. The _provider_ default value is `false`.",
					},
					"password_expiration_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days before the password expires. Valid values 1 to 365",
					},
					"password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minimum length of passwords. Valid values 4 to 16",
					},
					"password_minutes_of_inactivity_before_screen_timeout": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before the screen times out.",
					},
					"password_previous_password_block_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of previous passwords to block. Valid values 0 to 24",
					},
					"password_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "lowSecurityBiometric", "required", "atLeastNumeric", "numericComplex", "atLeastAlphabetic", "atLeastAlphanumeric", "alphanumericWithSymbols"),
						},
						PlanModifiers: []planmodifier.String{
							wpdefaultvalue.StringDefaultValue("atLeastNumeric"),
						},
						Computed:            true,
						MarkdownDescription: "Type of password that is required. / Android Work Profile required password type; possible values are: `deviceDefault` (Device default value, no intent.), `lowSecurityBiometric` (Low security biometrics based password required.), `required` (Required.), `atLeastNumeric` (At least numeric password required.), `numericComplex` (Numeric complex password required.), `atLeastAlphabetic` (At least alphabetic password required.), `atLeastAlphanumeric` (At least alphanumeric password required.), `alphanumericWithSymbols` (At least alphanumeric with symbols password required.). The _provider_ default value is `\"atLeastNumeric\"`.",
					},
					"password_sign_in_failure_count_before_factory_reset": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of sign in failures allowed before factory reset. Valid values 1 to 16",
					},
					"required_password_complexity": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("none", "low", "medium", "high"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
						Computed:            true,
						MarkdownDescription: "Indicates the required device password complexity on Android. One of: NONE, LOW, MEDIUM, HIGH. This is a new API targeted to Android 12+. / The password complexity types that can be set on Android. One of: NONE, LOW, MEDIUM, HIGH. This is an API targeted to Android 11+; possible values are: `none` (Device default value, no password.), `low` (The required password complexity on the device is of type low as defined by the Android documentation.), `medium` (The required password complexity on the device is of type medium as defined by the Android documentation.), `high` (The required password complexity on the device is of type high as defined by the Android documentation.). The _provider_ default value is `\"none\"`.",
					},
					"security_require_verify_apps": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require the Android Verify apps feature is turned on. The _provider_ default value is `false`.",
					},
					"vpn_always_on_package_identifier": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Enable lockdown mode for always-on VPN.",
					},
					"vpn_enable_always_on_lockdown_mode": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Enable lockdown mode for always-on VPN. The _provider_ default value is `false`.",
					},
					"work_profile_account_use": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("allowAllExceptGoogleAccounts", "blockAll", "allowAll", "unknownFutureValue"),
						},
						PlanModifiers: []planmodifier.String{
							wpdefaultvalue.StringDefaultValue("allowAllExceptGoogleAccounts"),
						},
						Computed:            true,
						MarkdownDescription: "Control user's ability to add accounts in work profile including Google accounts. / An enum representing possible values for account use in work profile; possible values are: `allowAllExceptGoogleAccounts` (Allow additon of all accounts except Google accounts in Android Work Profile.), `blockAll` (Block any account from being added in Android Work Profile.), `allowAll` (Allow addition of all accounts (including Google accounts) in Android Work Profile.), `unknownFutureValue` (Unknown future value for evolvable enum patterns.). The _provider_ default value is `\"allowAllExceptGoogleAccounts\"`.",
					},
					"work_profile_allow_app_installs_from_unknown_sources": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether to allow installation of apps from unknown sources. The _provider_ default value is `false`.",
					},
					"work_profile_allow_widgets": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Allow widgets from work profile apps. The _provider_ default value is `false`.",
					},
					"work_profile_block_adding_accounts": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Block users from adding/removing accounts in work profile. The _provider_ default value is `false`.",
					},
					"work_profile_block_camera": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Block work profile camera. The _provider_ default value is `false`.",
					},
					"work_profile_block_cross_profile_caller_id": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Block display work profile caller ID in personal profile. The _provider_ default value is `false`.",
					},
					"work_profile_block_cross_profile_contacts_search": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Block work profile contacts availability in personal profile. The _provider_ default value is `false`.",
					},
					"work_profile_block_cross_profile_copy_paste": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
						Computed:            true,
						MarkdownDescription: "Boolean that indicates if the setting disallow cross profile copy/paste is enabled. The _provider_ default value is `true`.",
					},
					"work_profile_block_notifications_while_device_locked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block notifications while device locked. The _provider_ default value is `false`.",
					},
					"work_profile_block_personal_app_installs_from_unknown_sources": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Prevent app installations from unknown sources in the personal profile. The _provider_ default value is `false`.",
					},
					"work_profile_block_screen_capture": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Block screen capture in work profile. The _provider_ default value is `false`.",
					},
					"work_profile_bluetooth_enable_contact_sharing": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Allow bluetooth devices to access enterprise contacts. The _provider_ default value is `false`.",
					},
					"work_profile_data_sharing_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "preventAny", "allowPersonalToWork", "noRestrictions"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: "Type of data sharing that is allowed. / Android Work Profile cross profile data sharing type; possible values are: `deviceDefault` (Device default value, no intent.), `preventAny` (Prevent any sharing.), `allowPersonalToWork` (Allow data sharing request from personal profile to work profile.), `noRestrictions` (No restrictions on sharing.). The _provider_ default value is `\"deviceDefault\"`.",
					},
					"work_profile_default_app_permission_policy": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "prompt", "autoGrant", "autoDeny"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: "Type of password that is required. / Android Work Profile default app permission policy type; possible values are: `deviceDefault` (Device default value, no intent.), `prompt` (Prompt.), `autoGrant` (Auto grant.), `autoDeny` (Auto deny.). The _provider_ default value is `\"deviceDefault\"`.",
					},
					"work_profile_password_block_face_unlock": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block face unlock for work profile. The _provider_ default value is `false`.",
					},
					"work_profile_password_block_fingerprint_unlock": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block fingerprint unlock for work profile. The _provider_ default value is `false`.",
					},
					"work_profile_password_block_iris_unlock": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block iris unlock for work profile. The _provider_ default value is `false`.",
					},
					"work_profile_password_block_trust_agents": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block Smart Lock and other trust agents for work profile. The _provider_ default value is `false`.",
					},
					"work_profile_password_expiration_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days before the work profile password expires. Valid values 1 to 365",
					},
					"work_profile_password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minimum length of work profile password. Valid values 4 to 16",
					},
					"work_profile_password_min_letter_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minimum # of letter characters required in work profile password. Valid values 1 to 10",
					},
					"work_profile_password_min_lower_case_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minimum # of lower-case characters required in work profile password. Valid values 1 to 10",
					},
					"work_profile_password_min_non_letter_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minimum # of non-letter characters required in work profile password. Valid values 1 to 10",
					},
					"work_profile_password_min_numeric_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minimum # of numeric characters required in work profile password. Valid values 1 to 10",
					},
					"work_profile_password_min_symbol_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minimum # of symbols required in work profile password. Valid values 1 to 10",
					},
					"work_profile_password_min_upper_case_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minimum # of upper-case characters required in work profile password. Valid values 1 to 10",
					},
					"work_profile_password_minutes_of_inactivity_before_screen_timeout": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before the screen times out.",
					},
					"work_profile_password_previous_password_block_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of previous work profile passwords to block. Valid values 0 to 24",
					},
					"work_profile_password_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "lowSecurityBiometric", "required", "atLeastNumeric", "numericComplex", "atLeastAlphabetic", "atLeastAlphanumeric", "alphanumericWithSymbols"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: "Type of work profile password that is required. / Android Work Profile required password type; possible values are: `deviceDefault` (Device default value, no intent.), `lowSecurityBiometric` (Low security biometrics based password required.), `required` (Required.), `atLeastNumeric` (At least numeric password required.), `numericComplex` (Numeric complex password required.), `atLeastAlphabetic` (At least alphabetic password required.), `atLeastAlphanumeric` (At least alphanumeric password required.), `alphanumericWithSymbols` (At least alphanumeric with symbols password required.). The _provider_ default value is `\"deviceDefault\"`.",
					},
					"work_profile_password_sign_in_failure_count_before_factory_reset": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of sign in failures allowed before work profile is removed and all corporate data deleted. Valid values 1 to 16",
					},
					"work_profile_required_password_complexity": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("none", "low", "medium", "high"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
						Computed:            true,
						MarkdownDescription: "Indicates the required work profile password complexity on Android. One of: NONE, LOW, MEDIUM, HIGH. This is a new API targeted to Android 12+. / The password complexity types that can be set on Android. One of: NONE, LOW, MEDIUM, HIGH. This is an API targeted to Android 11+; possible values are: `none` (Device default value, no password.), `low` (The required password complexity on the device is of type low as defined by the Android documentation.), `medium` (The required password complexity on the device is of type medium as defined by the Android documentation.), `high` (The required password complexity on the device is of type high as defined by the Android documentation.). The _provider_ default value is `\"none\"`.",
					},
					"work_profile_require_password": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Password is required or not for work profile. The _provider_ default value is `false`.",
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: "Android Work Profile general device configuration. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androidworkprofilegeneraldeviceconfiguration?view=graph-rest-beta",
			},
		},
		"edition_upgrade": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.editionUpgradeConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // editionUpgradeConfiguration
					"license": schema.StringAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
						Computed:            true,
						MarkdownDescription: "Edition Upgrade License File Content. The _provider_ default value is `\"\"`.",
					},
					"license_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("productKey", "licenseFile", "notConfigured"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Edition Upgrade License Type. / Edition Upgrade License type; possible values are: `productKey` (Product Key Type), `licenseFile` (License File Type), `notConfigured` (NotConfigured). The _provider_ default value is `\"notConfigured\"`.",
					},
					"product_key": schema.StringAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.String{
							wpdefaultvalue.StringDefaultValue(""),
							deviceConfigurationEditionUpgradeProductKeyPlanModifier{},
						},
						Computed:            true,
						MarkdownDescription: "Edition Upgrade Product Key. The _provider_ default value is `\"\"`.  \nProvider Note: MS Graph returns this attribute with most of the digits replaced by `#`. Therefore only the remaining visible digits will be compared to ensure that state matches plan, i.e. one of the visible digits must change for the entity to get updated in MS Graph.",
					},
					"target_edition": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("windows10Enterprise", "windows10EnterpriseN", "windows10Education", "windows10EducationN", "windows10MobileEnterprise", "windows10HolographicEnterprise", "windows10Professional", "windows10ProfessionalN", "windows10ProfessionalEducation", "windows10ProfessionalEducationN", "windows10ProfessionalWorkstation", "windows10ProfessionalWorkstationN", "notConfigured", "windows10Home", "windows10HomeChina", "windows10HomeN", "windows10HomeSingleLanguage", "windows10Mobile", "windows10IoTCore", "windows10IoTCoreCommercial"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Edition Upgrade Target Edition. / Windows 10 Edition type; possible values are: `windows10Enterprise` (Windows 10 Enterprise), `windows10EnterpriseN` (Windows 10 EnterpriseN), `windows10Education` (Windows 10 Education), `windows10EducationN` (Windows 10 EducationN), `windows10MobileEnterprise` (Windows 10 Mobile Enterprise), `windows10HolographicEnterprise` (Windows 10 Holographic Enterprise), `windows10Professional` (Windows 10 Professional), `windows10ProfessionalN` (Windows 10 ProfessionalN), `windows10ProfessionalEducation` (Windows 10 Professional Education), `windows10ProfessionalEducationN` (Windows 10 Professional EducationN), `windows10ProfessionalWorkstation` (Windows 10 Professional for Workstations), `windows10ProfessionalWorkstationN` (Windows 10 Professional for Workstations N), `notConfigured` (NotConfigured), `windows10Home` (Windows 10 Home), `windows10HomeChina` (Windows 10 Home China), `windows10HomeN` (Windows 10 Home N), `windows10HomeSingleLanguage` (Windows 10 Home Single Language), `windows10Mobile` (Windows 10 Mobile), `windows10IoTCore` (Windows 10 IoT Core), `windows10IoTCoreCommercial` (Windows 10 IoT Core Commercial). The _provider_ default value is `\"notConfigured\"`.",
					},
					"windows_s_mode": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("noRestriction", "block", "unlock"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("noRestriction")},
						Computed:            true,
						MarkdownDescription: "S mode configuration. / The possible options to configure S mode unlock; possible values are: `noRestriction` (This option will remove all restrictions to unlock S mode - default), `block` (This option will block the user to unlock the device from S mode), `unlock` (This option will unlock the device from S mode). The _provider_ default value is `\"noRestriction\"`.",
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: "Windows 10 Edition Upgrade configuration. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-editionupgradeconfiguration?view=graph-rest-beta  \nProvider Note: There have been reports that creating this type with an Entra Id application would fail even though it did have the permission `DeviceManagementConfiguration.ReadWrite.All`. Since this would be a limitation in MS Graph there is nothing that the provider could do about it.",
			},
		},
		"ios_custom": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.iosCustomConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // iosCustomConfiguration
					"payload": schema.StringAttribute{
						Required:            true,
						Validators:          []validator.String{wpvalidator.TranslateValueEncodeBase64()},
						MarkdownDescription: "Payload. (UTF8 encoded byte array)",
					},
					"file_name": schema.StringAttribute{
						Required:            true,
						Description:         `payloadFileName`, // custom MS Graph attribute name
						MarkdownDescription: "Payload file name (*.mobileconfig",
					},
					"name": schema.StringAttribute{
						Required:            true,
						Description:         `payloadName`, // custom MS Graph attribute name
						MarkdownDescription: "Name that is displayed to the user.",
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: "This topic provides descriptions of the declared methods, properties and relationships exposed by the iosCustomConfiguration resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioscustomconfiguration?view=graph-rest-beta",
			},
		},
		"ios_device_features": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.iosDeviceFeaturesConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // iosDeviceFeaturesConfiguration
					"airprint_destinations": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationAirPrintDestinationAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						Description:         `airPrintDestinations`, // custom MS Graph attribute name
						MarkdownDescription: "An array of AirPrint printers that should always be shown. This collection can contain a maximum of 500 elements. / Represents an AirPrint destination. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-airprintdestination?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"asset_tag_template": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Asset tag information for the device, displayed on the login window and lock screen.",
					},
					"content_filter_settings": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // iosWebContentFilterBase
							"auto_filter": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.iosWebContentFilterAutoFilter",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // iosWebContentFilterAutoFilter
										"allowed_urls": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: "Additional URLs allowed for access. The _provider_ default value is `[]`.",
										},
										"blocked_urls": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: "Additional URLs blocked for access. The _provider_ default value is `[]`.",
										},
									},
									Validators: []validator.Object{
										deviceConfigurationIosWebContentFilterBaseValidator,
									},
									MarkdownDescription: "Represents an iOS Web Content Filter setting type, which enables iOS automatic filter feature and allows for additional URL access control. When constructed with no property values, the iOS device will enable the automatic filter regardless. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioswebcontentfilterautofilter?view=graph-rest-beta",
								},
							},
							"specific_websites_access": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.iosWebContentFilterSpecificWebsitesAccess",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // iosWebContentFilterSpecificWebsitesAccess
										"specific_websites_only": schema.SetNestedAttribute{
											Optional: true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: deviceConfigurationIosBookmarkAttributes,
											},
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: "URL bookmarks which will be installed into built-in browser and user is only allowed to access websites through bookmarks. This collection can contain a maximum of 500 elements. / iOS URL bookmark / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosbookmark?view=graph-rest-beta. The _provider_ default value is `[]`.",
										},
										"website_list": schema.SetNestedAttribute{
											Optional: true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: deviceConfigurationIosBookmarkAttributes,
											},
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: "URL bookmarks which will be installed into built-in browser and user is only allowed to access websites through bookmarks. This collection can contain a maximum of 500 elements. / iOS URL bookmark / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosbookmark?view=graph-rest-beta. The _provider_ default value is `[]`.",
										},
									},
									Validators: []validator.Object{
										deviceConfigurationIosWebContentFilterBaseValidator,
									},
									MarkdownDescription: "Represents an iOS Web Content Filter setting type, which installs URL bookmarks into iOS built-in browser. An example scenario is in the classroom where teachers would like the students to navigate websites through browser bookmarks configured on their iOS devices, and no access to other sites. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioswebcontentfilterspecificwebsitesaccess?view=graph-rest-beta",
								},
							},
						},
						MarkdownDescription: "Gets or sets iOS Web Content Filter settings, supervised mode only / Represents an iOS Web Content Filter setting base type. An empty and abstract base. Caller should use one of derived types for configurations. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioswebcontentfilterbase?view=graph-rest-beta",
					},
					"home_screen_dock_icons": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // iosHomeScreenItem
								"display_name": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Name of the app",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "A list of app and folders to appear on the Home Screen Dock. This collection can contain a maximum of 500 elements. / Represents an item on the iOS Home Screen / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioshomescreenitem?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"home_screen_grid_height": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Gets or sets the number of rows to render when configuring iOS home screen layout settings. If this value is configured, homeScreenGridWidth must be configured as well.",
					},
					"home_screen_grid_width": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Gets or sets the number of columns to render when configuring iOS home screen layout settings. If this value is configured, homeScreenGridHeight must be configured as well.",
					},
					"home_screen_pages": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // iosHomeScreenPage
								"display_name": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Name of the page",
								},
								"icons": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{ // iosHomeScreenItem
											"display_name": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "Name of the app",
											},
										},
									},
									PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: "A list of apps, folders, and web clips to appear on a page. This collection can contain a maximum of 500 elements. / Represents an item on the iOS Home Screen / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioshomescreenitem?view=graph-rest-beta. The _provider_ default value is `[]`.",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "A list of pages on the Home Screen. This collection can contain a maximum of 500 elements. / A page containing apps, folders, and web clips on the Home Screen. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioshomescreenpage?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"ios_single_sign_on_extension": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // iosSingleSignOnExtension
							"azure_ad": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.iosAzureAdSingleSignOnExtension",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // iosAzureAdSingleSignOnExtension
										"bundle_id_access_control_list": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: "An optional list of additional bundle IDs allowed to use the AAD extension for single sign-on. The _provider_ default value is `[]`.",
										},
										"configurations": schema.SetNestedAttribute{
											Optional: true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: deviceConfigurationKeyTypedValuePairAttributes,
											},
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: "Gets or sets a list of typed key-value pairs used to configure Credential-type profiles. This collection can contain a maximum of 500 elements. / A key-value pair with a string key and a typed value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keytypedvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`.",
										},
										"enable_shared_device_mode": schema.BoolAttribute{
											Optional:            true,
											PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
											Computed:            true,
											MarkdownDescription: "Enables or disables shared device mode. The _provider_ default value is `false`.",
										},
									},
									Validators: []validator.Object{
										deviceConfigurationIosSingleSignOnExtensionValidator,
									},
									MarkdownDescription: "Represents an Azure AD-type Single Sign-On extension profile for iOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosazureadsinglesignonextension?view=graph-rest-beta",
								},
							},
							"credential": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.iosCredentialSingleSignOnExtension",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // iosCredentialSingleSignOnExtension
										"configurations": schema.SetNestedAttribute{
											Optional: true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: deviceConfigurationKeyTypedValuePairAttributes,
											},
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: "Gets or sets a list of typed key-value pairs used to configure Credential-type profiles. This collection can contain a maximum of 500 elements. / A key-value pair with a string key and a typed value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keytypedvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`.",
										},
										"domains": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: "Gets or sets a list of hosts or domain names for which the app extension performs SSO. The _provider_ default value is `[]`.",
										},
										"extension_identifier": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "Gets or sets the bundle ID of the app extension that performs SSO for the specified URLs.",
										},
										"realm": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "Gets or sets the case-sensitive realm name for this profile.",
										},
										"team_identifier": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "Gets or sets the team ID of the app extension that performs SSO for the specified URLs.",
										},
									},
									Validators: []validator.Object{
										deviceConfigurationIosSingleSignOnExtensionValidator,
									},
									MarkdownDescription: "Represents a Credential-type Single Sign-On extension profile for iOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioscredentialsinglesignonextension?view=graph-rest-beta",
								},
							},
							"kerberos": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.iosKerberosSingleSignOnExtension",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // iosKerberosSingleSignOnExtension
										"active_directory_site_code": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "Gets or sets the Active Directory site.",
										},
										"block_active_directory_site_auto_discovery": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: "Enables or disables whether the Kerberos extension can automatically determine its site name.",
										},
										"block_automatic_login": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: "Enables or disables Keychain usage.",
										},
										"cache_name": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "Gets or sets the Generic Security Services name of the Kerberos cache to use for this profile.",
										},
										"credential_bundle_id_access_control_list": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											MarkdownDescription: "Gets or sets a list of app Bundle IDs allowed to access the Kerberos Ticket Granting Ticket.",
										},
										"domain_realms": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											MarkdownDescription: "Gets or sets a list of realms for custom domain-realm mapping. Realms are case sensitive.",
										},
										"domains": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											MarkdownDescription: "Gets or sets a list of hosts or domain names for which the app extension performs SSO.",
										},
										"is_default_realm": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: "When true, this profile's realm will be selected as the default. Necessary if multiple Kerberos-type profiles are configured.",
										},
										"managed_apps_in_bundle_id_acl_included": schema.BoolAttribute{
											Required:            true,
											Description:         `managedAppsInBundleIdACLIncluded`, // custom MS Graph attribute name
											MarkdownDescription: "When set to True, the Kerberos extension allows managed apps, and any apps entered with the app bundle ID to access the credential. When set to False, the Kerberos extension allows all apps to access the credential. Available for devices running iOS and iPadOS versions 14 and later.",
										},
										"password_block_modification": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: "Enables or disables password changes.",
										},
										"password_change_url": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "Gets or sets the URL that the user will be sent to when they initiate a password change.",
										},
										"password_enable_local_sync": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: "Enables or disables password syncing. This won't affect users logged in with a mobile account on macOS.",
										},
										"password_expiration_days": schema.Int64Attribute{
											Optional:            true,
											MarkdownDescription: "Overrides the default password expiration in days. For most domains, this value is calculated automatically.",
										},
										"password_expiration_notification_days": schema.Int64Attribute{
											Optional:            true,
											MarkdownDescription: "Gets or sets the number of days until the user is notified that their password will expire (default is 15).",
										},
										"password_minimum_age_days": schema.Int64Attribute{
											Optional:            true,
											MarkdownDescription: "Gets or sets the minimum number of days until a user can change their password again.",
										},
										"password_minimum_length": schema.Int64Attribute{
											Optional:            true,
											MarkdownDescription: "Gets or sets the minimum length of a password.",
										},
										"password_previous_password_block_count": schema.Int64Attribute{
											Optional:            true,
											MarkdownDescription: "Gets or sets the number of previous passwords to block.",
										},
										"password_require_active_directory_complexity": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: "Enables or disables whether passwords must meet Active Directory's complexity requirements.",
										},
										"password_requirements_description": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "Gets or sets a description of the password complexity requirements.",
										},
										"realm": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "Gets or sets the case-sensitive realm name for this profile.",
										},
										"require_user_presence": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: "Gets or sets whether to require authentication via Touch ID, Face ID, or a passcode to access the keychain entry.",
										},
										"sign_in_help_text": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "Text displayed to the user at the Kerberos sign in window. Available for devices running iOS and iPadOS versions 14 and later.",
										},
										"user_principal_name": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "Gets or sets the principle user name to use for this profile. The realm name does not need to be included.",
										},
									},
									Validators: []validator.Object{
										deviceConfigurationIosSingleSignOnExtensionValidator,
									},
									MarkdownDescription: "Represents a Kerberos-type Single Sign-On extension profile for iOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioskerberossinglesignonextension?view=graph-rest-beta",
								},
							},
							"redirect": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.iosRedirectSingleSignOnExtension",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // iosRedirectSingleSignOnExtension
										"configurations": schema.SetNestedAttribute{
											Optional: true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: deviceConfigurationKeyTypedValuePairAttributes,
											},
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: "Gets or sets a list of typed key-value pairs used to configure Credential-type profiles. This collection can contain a maximum of 500 elements. / A key-value pair with a string key and a typed value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keytypedvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`.",
										},
										"extension_identifier": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "Gets or sets the bundle ID of the app extension that performs SSO for the specified URLs.",
										},
										"team_identifier": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "Gets or sets the team ID of the app extension that performs SSO for the specified URLs.",
										},
										"url_prefixes": schema.SetAttribute{
											ElementType:         types.StringType,
											Required:            true,
											MarkdownDescription: "One or more URL prefixes of identity providers on whose behalf the app extension performs single sign-on. URLs must begin with http:// or https://. All URL prefixes must be unique for all profiles.",
										},
									},
									Validators: []validator.Object{
										deviceConfigurationIosSingleSignOnExtensionValidator,
									},
									MarkdownDescription: "Represents a Redirect-type Single Sign-On extension profile for iOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosredirectsinglesignonextension?view=graph-rest-beta",
								},
							},
						},
						MarkdownDescription: "Gets or sets a single sign-on extension profile. / An abstract base class for all iOS-specific single sign-on extension types. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iossinglesignonextension?view=graph-rest-beta",
					},
					"lock_screen_footnote": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "A footnote displayed on the login window and lock screen. Available in iOS 9.3.1 and later.",
					},
					"notification_settings": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // iosNotificationSettings
								"alert_type": schema.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.OneOf("deviceDefault", "banner", "modal", "none"),
									},
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
									Computed:            true,
									MarkdownDescription: "Indicates the type of alert for notifications for this app. / Notification Settings Alert Type; possible values are: `deviceDefault` (Device default value, no intent.), `banner` (Banner.), `modal` (Modal.), `none` (None.). The _provider_ default value is `\"deviceDefault\"`.",
								},
								"app_name": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Application name to be associated with the bundleID.",
								},
								"badges_enabled": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: "Indicates whether badges are allowed for this app.",
								},
								"bundle_id": schema.StringAttribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
									Computed:            true,
									Description:         `bundleID`, // custom MS Graph attribute name
									MarkdownDescription: "Bundle id of app to which to apply these notification settings. The _provider_ default value is `\"\"`.",
								},
								"enabled": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: "Indicates whether notifications are allowed for this app.",
								},
								"preview_visibility": schema.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.OneOf("notConfigured", "alwaysShow", "hideWhenLocked", "neverShow"),
									},
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
									Computed:            true,
									MarkdownDescription: "Overrides the notification preview policy set by the user on an iOS device. / Determines when notification previews are visible on an iOS device. Previews can include things like text (from Messages and Mail) and invitation details (from Calendar). When configured, it will override the user's defined preview settings; possible values are: `notConfigured` (Notification preview settings will not be overwritten.), `alwaysShow` (Always show notification previews.), `hideWhenLocked` (Only show notification previews when the device is unlocked.), `neverShow` (Never show notification previews.). The _provider_ default value is `\"notConfigured\"`.",
								},
								"publisher": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Publisher to be associated with the bundleID.",
								},
								"show_in_notification_center": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: "Indicates whether notifications can be shown in notification center.",
								},
								"show_on_lock_screen": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: "Indicates whether notifications can be shown on the lock screen.",
								},
								"sounds_enabled": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: "Indicates whether sounds are allowed for this app.",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Notification settings for each bundle id. Applicable to devices in supervised mode only (iOS 9.3 and later). This collection can contain a maximum of 500 elements. / An item describing notification setting. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosnotificationsettings?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"single_sign_on_settings": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // iosSingleSignOnSettings
							"allowed_apps_list": schema.SetNestedAttribute{
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: deviceConfigurationAppListItemAttributes,
								},
								PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
								Computed:            true,
								MarkdownDescription: "List of app identifiers that are allowed to use this login. If this field is omitted, the login applies to all applications on the device. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `[]`.",
							},
							"allowed_urls": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
								Computed:            true,
								MarkdownDescription: "List of HTTP URLs that must be matched in order to use this login. With iOS 9.0 or later, a wildcard characters may be used. The _provider_ default value is `[]`.",
							},
							"display_name": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "The display name of login settings shown on the receiving device.",
							},
							"kerberos_principal_name": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "A Kerberos principal name. If not provided, the user is prompted for one during profile installation.",
							},
							"kerberos_realm": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "A Kerberos realm name. Case sensitive.",
							},
						},
						MarkdownDescription: "The Kerberos login settings that enable apps on receiving devices to authenticate smoothly. / iOS Kerberos authentication settings for single sign-on / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iossinglesignonsettings?view=graph-rest-beta",
					},
					"wallpaper_display_location": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "lockScreen", "homeScreen", "lockAndHomeScreens"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "A wallpaper display location specifier. / An enum type for wallpaper display location specifier; possible values are: `notConfigured` (No location specified for wallpaper display.), `lockScreen` (A configured wallpaper image is displayed on Lock screen.), `homeScreen` (A configured wallpaper image is displayed on Home (icon list) screen.), `lockAndHomeScreens` (A configured wallpaper image is displayed on Lock screen and Home screen.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"wallpaper_image": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mimeContent
							"type": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Indicates the content mime type.",
							},
							"value": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "The byte array that contains the actual content.",
							},
						},
						MarkdownDescription: "A wallpaper image must be in either PNG or JPEG format. It requires a supervised device with iOS 8 or later version. / Contains properties for a generic mime content. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mimecontent?view=graph-rest-beta",
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: "iOS Device Features Configuration Profile. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosdevicefeaturesconfiguration?view=graph-rest-beta",
			},
		},
		"ios_eas_email_profile": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.iosEasEmailProfileConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // iosEasEmailProfileConfiguration
					"custom_domain_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Custom domain name value used while generating an email profile before installing on the device.",
					},
					"user_domain_name_source": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("fullDomainName", "netBiosDomainName"),
						},
						MarkdownDescription: "UserDomainname attribute that is picked from AAD and injected into this profile before installing on the device. / Domainname source; possible values are: `fullDomainName` (Full domain name.), `netBiosDomainName` (net bios domain name.)",
					},
					"username_aad_source": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userPrincipalName", "primarySmtpAddress", "samAccountName"),
						},
						Description:         `usernameAADSource`, // custom MS Graph attribute name
						MarkdownDescription: "Name of the AAD field, that will be used to retrieve UserName for email profile. / Username source; possible values are: `userPrincipalName` (User principal name.), `primarySmtpAddress` (Primary SMTP address.), `samAccountName` (The user sam account name.)",
					},
					"username_source": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userPrincipalName", "primarySmtpAddress"),
						},
						PlanModifiers: []planmodifier.String{
							wpdefaultvalue.StringDefaultValue("userPrincipalName"),
						},
						Computed:            true,
						MarkdownDescription: "Username attribute that is picked from AAD and injected into this profile before installing on the device. / Possible values for username source or email source; possible values are: `userPrincipalName` (User principal name.), `primarySmtpAddress` (Primary SMTP address.). The _provider_ default value is `\"userPrincipalName\"`.",
					},
					"account_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Account name.",
					},
					"authentication_method": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("usernameAndPassword", "certificate", "derivedCredential"),
						},
						MarkdownDescription: "Authentication method for this Email profile. / Exchange Active Sync authentication method; possible values are: `usernameAndPassword` (Authenticate with a username and password.), `certificate` (Authenticate with a certificate.), `derivedCredential` (Authenticate with derived credential.)",
					},
					"block_moving_messages_to_other_email_accounts": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block moving messages to other email accounts.",
					},
					"block_sending_email_from_third_party_apps": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block sending email from third party apps.",
					},
					"block_syncing_recently_used_email_addresses": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates whether or not to block syncing recently used email addresses, for instance - when composing new email.",
					},
					"duration_of_email_to_sync": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userDefined", "oneDay", "threeDays", "oneWeek", "twoWeeks", "oneMonth", "unlimited"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("userDefined")},
						Computed:            true,
						MarkdownDescription: "Duration of time email should be synced back to. . / Possible values for email sync duration; possible values are: `userDefined` (User Defined, default value, no intent.), `oneDay` (Sync one day of email.), `threeDays` (Sync three days of email.), `oneWeek` (Sync one week of email.), `twoWeeks` (Sync two weeks of email.), `oneMonth` (Sync one month of email.), `unlimited` (Sync an unlimited duration of email.). The _provider_ default value is `\"userDefined\"`.",
					},
					"eas_services": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							wpvalidator.FlagEnumValues("none", "calendars", "contacts", "email", "notes", "reminders"),
						},
						MarkdownDescription: "Exchange data to sync. / Exchange Active Sync services; possible values are: `none`, `calendars` (Enables synchronization of calendars.), `contacts` (Enables synchronization of contacts.), `email` (Enables synchronization of email.), `notes` (Enables synchronization of notes.), `reminders` (Enables synchronization of reminders.)",
					},
					"eas_services_user_override_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Allow users to change sync settings.",
					},
					"email_address_source": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userPrincipalName", "primarySmtpAddress"),
						},
						MarkdownDescription: "Email attribute that is picked from AAD and injected into this profile before installing on the device. / Possible values for username source or email source; possible values are: `userPrincipalName` (User principal name.), `primarySmtpAddress` (Primary SMTP address.)",
					},
					"encryption_certificate_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("none", "certificate", "derivedCredential"),
						},
						MarkdownDescription: "Encryption Certificate type for this Email profile. / Supported certificate sources for email signing and encryption; possible values are: `none` (Do not use a certificate as a source.), `certificate` (Use an certificate for certificate source.), `derivedCredential` (Use a derived credential for certificate source.)",
					},
					"host_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Exchange location that (URL) that the native mail app connects to.",
					},
					"per_app_vpn_profile_id": schema.StringAttribute{
						Optional:            true,
						Description:         `perAppVPNProfileId`, // custom MS Graph attribute name
						MarkdownDescription: "Profile ID of the Per-App VPN policy to be used to access emails from the native Mail client",
					},
					"require_ssl": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to use SSL. The _provider_ default value is `true`.",
					},
					"signing_certificate_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("none", "certificate", "derivedCredential"),
						},
						MarkdownDescription: "Signing Certificate type for this Email profile. / Supported certificate sources for email signing and encryption; possible values are: `none` (Do not use a certificate as a source.), `certificate` (Use an certificate for certificate source.), `derivedCredential` (Use a derived credential for certificate source.)",
					},
					"use_oauth": schema.BoolAttribute{
						Optional:            true,
						Description:         `useOAuth`, // custom MS Graph attribute name
						MarkdownDescription: "Specifies whether the connection should use OAuth for authentication.",
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: "By providing configurations in this profile you can instruct the native email client on iOS devices to communicate with an Exchange server and get email, contacts, calendar, reminders, and notes. Furthermore, you can also specify how much email to sync and how often the device should sync. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioseasemailprofileconfiguration?view=graph-rest-beta",
			},
		},
		"ios_general_device": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.iosGeneralDeviceConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // iosGeneralDeviceConfiguration
					"account_block_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow account modification when the device is in supervised mode. The _provider_ default value is `false`.",
					},
					"activation_lock_allow_when_supervised": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow activation lock when the device is in the supervised mode. The _provider_ default value is `false`.",
					},
					"airdrop_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `airDropBlocked`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to allow AirDrop when the device is in supervised mode. The _provider_ default value is `false`.",
					},
					"airdrop_force_unmanaged_drop_target": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `airDropForceUnmanagedDropTarget`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to cause AirDrop to be considered an unmanaged drop target (iOS 9.0 and later). The _provider_ default value is `false`.",
					},
					"airplay_force_pairing_password_for_outgoing_requests": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `airPlayForcePairingPasswordForOutgoingRequests`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to enforce all devices receiving AirPlay requests from this device to use a pairing password. The _provider_ default value is `false`.",
					},
					"airprint_block_credentials_storage": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `airPrintBlockCredentialsStorage`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not keychain storage of username and password for Airprint is blocked (iOS 11.0 and later). The _provider_ default value is `false`.",
					},
					"airprint_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `airPrintBlocked`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not AirPrint is blocked (iOS 11.0 and later). The _provider_ default value is `false`.",
					},
					"airprint_blocki_beacon_discovery": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `airPrintBlockiBeaconDiscovery`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not iBeacon discovery of AirPrint printers is blocked. This prevents spurious AirPrint Bluetooth beacons from phishing for network traffic (iOS 11.0 and later). The _provider_ default value is `false`.",
					},
					"airprint_force_trusted_tls": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `airPrintForceTrustedTLS`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates if trusted certificates are required for TLS printing communication (iOS 11.0 and later). The _provider_ default value is `false`.",
					},
					"app_clips_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Prevents a user from adding any App Clips and removes any existing App Clips on the device. The _provider_ default value is `false`.",
					},
					"apple_news_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using News when the device is in supervised mode (iOS 9.0 and later). The _provider_ default value is `false`.",
					},
					"apple_personalized_ads_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Limits Apple personalized advertising when true. Available in iOS 14 and later. The _provider_ default value is `false`.",
					},
					"apple_watch_block_pairing": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow Apple Watch pairing when the device is in supervised mode (iOS 9.0 and later). The _provider_ default value is `false`.",
					},
					"apple_watch_force_wrist_detection": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to force a paired Apple Watch to use Wrist Detection (iOS 8.2 and later). The _provider_ default value is `false`.",
					},
					"app_removal_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates if the removal of apps is allowed. The _provider_ default value is `false`.",
					},
					"apps_single_app_mode_list": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationAppListItemAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Gets or sets the list of iOS apps allowed to autonomously enter Single App Mode. Supervised only. iOS 7.0 and later. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"app_store_block_automatic_downloads": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the automatic downloading of apps purchased on other devices when the device is in supervised mode (iOS 9.0 and later). The _provider_ default value is `false`.",
					},
					"app_store_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using the App Store. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.",
					},
					"app_store_block_in_app_purchases": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from making in app purchases. The _provider_ default value is `false`.",
					},
					"app_store_block_ui_app_installation": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `appStoreBlockUIAppInstallation`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to block the App Store app, not restricting installation through Host apps. Applies to supervised mode only (iOS 9.0 and later). The _provider_ default value is `false`.",
					},
					"app_store_require_password": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to require a password when using the app store. The _provider_ default value is `false`.",
					},
					"apps_visibility_list": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationAppListItemAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "List of apps in the visibility list (either visible/launchable apps list or hidden/unlaunchable apps list, controlled by AppsVisibilityListType) (iOS 9.3 and later). This collection can contain a maximum of 10000 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"apps_visibility_list_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("none", "appsInListCompliant", "appsNotInListCompliant"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
						Computed:            true,
						MarkdownDescription: "Type of list that is in the AppsVisibilityList. / Possible values of the compliance app list; possible values are: `none` (Default value, no intent.), `appsInListCompliant` (The list represents the apps that will be considered compliant (only apps on the list are compliant).), `appsNotInListCompliant` (The list represents the apps that will be considered non compliant (all apps are compliant except apps on the list).). The _provider_ default value is `\"none\"`.",
					},
					"auto_fill_force_authentication": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to force user authentication before autofilling passwords and credit card information in Safari and other apps on supervised devices. The _provider_ default value is `false`.",
					},
					"auto_unlock_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Blocks users from unlocking their device with Apple Watch. Available for devices running iOS and iPadOS versions 14.5 and later. The _provider_ default value is `false`.",
					},
					"block_system_app_removal": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not the removal of system apps from the device is blocked on a supervised device (iOS 11.0 and later). The _provider_ default value is `false`.",
					},
					"bluetooth_block_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow modification of Bluetooth settings when the device is in supervised mode (iOS 10.0 and later). The _provider_ default value is `false`.",
					},
					"camera_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from accessing the camera of the device. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.",
					},
					"cellular_block_data_roaming": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block data roaming. The _provider_ default value is `false`.",
					},
					"cellular_block_global_background_fetch_while_roaming": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block global background fetch while roaming. The _provider_ default value is `false`.",
					},
					"cellular_block_per_app_data_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow changes to cellular app data usage settings when the device is in supervised mode. The _provider_ default value is `false`.",
					},
					"cellular_block_personal_hotspot": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block Personal Hotspot. The _provider_ default value is `false`.",
					},
					"cellular_block_personal_hotspot_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from modifying the personal hotspot setting (iOS 12.2 or later). The _provider_ default value is `false`.",
					},
					"cellular_block_plan_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow users to change the settings of the cellular plan on a supervised device. The _provider_ default value is `false`.",
					},
					"cellular_block_voice_roaming": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block voice roaming. The _provider_ default value is `false`.",
					},
					"certificates_block_untrusted_tls_certificates": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block untrusted TLS certificates. The _provider_ default value is `false`.",
					},
					"classroom_app_block_remote_screen_observation": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow remote screen observation by Classroom app when the device is in supervised mode (iOS 9.3 and later). The _provider_ default value is `false`.",
					},
					"classroom_app_force_unprompted_screen_observation": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to automatically give permission to the teacher of a managed course on the Classroom app to view a student's screen without prompting when the device is in supervised mode. The _provider_ default value is `false`.",
					},
					"classroom_force_automatically_join_classes": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to automatically give permission to the teacher's requests, without prompting the student, when the device is in supervised mode. The _provider_ default value is `false`.",
					},
					"classroom_force_request_permission_to_leave_classes": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether a student enrolled in an unmanaged course via Classroom will request permission from the teacher when attempting to leave the course (iOS 11.3 and later). The _provider_ default value is `false`.",
					},
					"classroom_force_unprompted_app_and_device_lock": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow the teacher to lock apps or the device without prompting the student. Supervised only. The _provider_ default value is `false`.",
					},
					"compliant_app_list_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("none", "appsInListCompliant", "appsNotInListCompliant"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
						Computed:            true,
						MarkdownDescription: "List that is in the AppComplianceList. / Possible values of the compliance app list; possible values are: `none` (Default value, no intent.), `appsInListCompliant` (The list represents the apps that will be considered compliant (only apps on the list are compliant).), `appsNotInListCompliant` (The list represents the apps that will be considered non compliant (all apps are compliant except apps on the list).). The _provider_ default value is `\"none\"`.",
					},
					"compliant_apps_list": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationAppListItemAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "List of apps in the compliance (either allow list or block list, controlled by CompliantAppListType). This collection can contain a maximum of 10000 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"configuration_profile_block_changes": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from installing configuration profiles and certificates interactively when the device is in supervised mode. The _provider_ default value is `false`.",
					},
					"contacts_allow_managed_to_unmanaged_write": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not managed apps can write contacts to unmanaged contacts accounts (iOS 12.0 and later). The _provider_ default value is `false`.",
					},
					"contacts_allow_unmanaged_to_managed_read": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not unmanaged apps can read from managed contacts accounts (iOS 12.0 or later). The _provider_ default value is `false`.",
					},
					"continuous_path_keyboard_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the continuous path keyboard when the device is supervised (iOS 13 or later). The _provider_ default value is `false`.",
					},
					"date_and_time_force_set_automatically": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not the Date and Time \"Set Automatically\" feature is enabled and cannot be turned off by the user (iOS 12.0 and later). The _provider_ default value is `false`.",
					},
					"definition_lookup_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block definition lookup when the device is in supervised mode (iOS 8.1.3 and later ). The _provider_ default value is `false`.",
					},
					"device_block_enable_restrictions": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow the user to enables restrictions in the device settings when the device is in supervised mode. The _provider_ default value is `false`.",
					},
					"device_block_erase_content_and_settings": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow the use of the 'Erase all content and settings' option on the device when the device is in supervised mode. The _provider_ default value is `false`.",
					},
					"device_block_name_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow device name modification when the device is in supervised mode (iOS 9.0 and later). The _provider_ default value is `false`.",
					},
					"diagnostic_data_block_submission": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block diagnostic data submission. The _provider_ default value is `false`.",
					},
					"diagnostic_data_block_submission_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow diagnostics submission settings modification when the device is in supervised mode (iOS 9.3.2 and later). The _provider_ default value is `false`.",
					},
					"documents_block_managed_documents_in_unmanaged_apps": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from viewing managed documents in unmanaged apps. The _provider_ default value is `false`.",
					},
					"documents_block_unmanaged_documents_in_managed_apps": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from viewing unmanaged documents in managed apps. The _provider_ default value is `false`.",
					},
					"email_in_domain_suffixes": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "An email address lacking a suffix that matches any of these strings will be considered out-of-domain. The _provider_ default value is `[]`.",
					},
					"enterprise_app_block_trust": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from trusting an enterprise app. The _provider_ default value is `false`.",
					},
					"enterprise_book_block_backup": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not Enterprise book back up is blocked. The _provider_ default value is `false`.",
					},
					"enterprise_book_block_metadata_sync": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not Enterprise book notes and highlights sync is blocked. The _provider_ default value is `false`.",
					},
					"esim_block_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow the addition or removal of cellular plans on the eSIM of a supervised device. The _provider_ default value is `false`.",
					},
					"facetime_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `faceTimeBlocked`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to block the user from using FaceTime. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.",
					},
					"files_network_drive_access_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates if devices can access files or other resources on a network server using the Server Message Block (SMB) protocol. Available for devices running iOS and iPadOS, versions 13.0 and later. The _provider_ default value is `false`.",
					},
					"files_usb_drive_access_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates if sevices with access can connect to and open files on a USB drive. Available for devices running iOS and iPadOS, versions 13.0 and later. The _provider_ default value is `false`.",
					},
					"find_my_device_in_find_my_app_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block Find My Device when the device is supervised (iOS 13 or later). The _provider_ default value is `false`.",
					},
					"find_my_friends_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block changes to Find My Friends when the device is in supervised mode. The _provider_ default value is `false`.",
					},
					"find_my_friends_in_find_my_app_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block Find My Friends when the device is supervised (iOS 13 or later). The _provider_ default value is `false`.",
					},
					"game_center_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using Game Center when the device is in supervised mode. The _provider_ default value is `false`.",
					},
					"gaming_block_game_center_friends": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from having friends in Game Center. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.",
					},
					"gaming_block_multiplayer": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using multiplayer gaming. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.",
					},
					"host_pairing_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "indicates whether or not to allow host pairing to control the devices an iOS device can pair with when the iOS device is in supervised mode. The _provider_ default value is `false`.",
					},
					"ibooks_store_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iBooksStoreBlocked`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to block the user from using the iBooks Store when the device is in supervised mode. The _provider_ default value is `false`.",
					},
					"ibooks_store_block_erotica": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iBooksStoreBlockErotica`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to block the user from downloading media from the iBookstore that has been tagged as erotica. The _provider_ default value is `false`.",
					},
					"icloud_block_activity_continuation": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iCloudBlockActivityContinuation`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to block the user from continuing work they started on iOS device to another iOS or macOS device. The _provider_ default value is `false`.",
					},
					"icloud_block_backup": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iCloudBlockBackup`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to block iCloud backup. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.",
					},
					"icloud_block_document_sync": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iCloudBlockDocumentSync`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to block iCloud document sync. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.",
					},
					"icloud_block_managed_apps_sync": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iCloudBlockManagedAppsSync`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to block Managed Apps Cloud Sync. The _provider_ default value is `false`.",
					},
					"icloud_block_photo_library": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iCloudBlockPhotoLibrary`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to block iCloud Photo Library. The _provider_ default value is `false`.",
					},
					"icloud_block_photo_stream_sync": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iCloudBlockPhotoStreamSync`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to block iCloud Photo Stream Sync. The _provider_ default value is `false`.",
					},
					"icloud_block_shared_photo_stream": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iCloudBlockSharedPhotoStream`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to block Shared Photo Stream. The _provider_ default value is `false`.",
					},
					"icloud_private_relay_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iCloudPrivateRelayBlocked`, // custom MS Graph attribute name
						MarkdownDescription: "iCloud private relay is an iCloud+ service that prevents networks and servers from monitoring a person's activity across the internet. By blocking iCloud private relay, Apple will not encrypt the traffic leaving the device. Available for devices running iOS 15 and later. The _provider_ default value is `false`.",
					},
					"icloud_require_encrypted_backup": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iCloudRequireEncryptedBackup`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to require backups to iCloud be encrypted. The _provider_ default value is `false`.",
					},
					"itunes_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iTunesBlocked`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to block the iTunes app. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.",
					},
					"itunes_block_explicit_content": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iTunesBlockExplicitContent`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to block the user from accessing explicit content in iTunes and the App Store. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.",
					},
					"itunes_block_music_service": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iTunesBlockMusicService`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to block Music service and revert Music app to classic mode when the device is in supervised mode (iOS 9.3 and later and macOS 10.12 and later). The _provider_ default value is `false`.",
					},
					"itunes_block_radio": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iTunesBlockRadio`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to block the user from using iTunes Radio when the device is in supervised mode (iOS 9.3 and later). The _provider_ default value is `false`.",
					},
					"keyboard_block_auto_correct": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block keyboard auto-correction when the device is in supervised mode (iOS 8.1.3 and later). The _provider_ default value is `false`.",
					},
					"keyboard_block_dictation": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using dictation input when the device is in supervised mode. The _provider_ default value is `false`.",
					},
					"keyboard_block_predictive": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block predictive keyboards when device is in supervised mode (iOS 8.1.3 and later). The _provider_ default value is `false`.",
					},
					"keyboard_block_shortcuts": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block keyboard shortcuts when the device is in supervised mode (iOS 9.0 and later). The _provider_ default value is `false`.",
					},
					"keyboard_block_spell_check": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block keyboard spell-checking when the device is in supervised mode (iOS 8.1.3 and later). The _provider_ default value is `false`.",
					},
					"keychain_block_cloud_sync": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not iCloud keychain synchronization is blocked. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.",
					},
					"kiosk_mode_allow_assistive_speak": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow assistive speak while in kiosk mode. The _provider_ default value is `false`.",
					},
					"kiosk_mode_allow_assistive_touch_settings": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow access to the Assistive Touch Settings while in kiosk mode. The _provider_ default value is `false`.",
					},
					"kiosk_mode_allow_color_inversion_settings": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow access to the Color Inversion Settings while in kiosk mode. The _provider_ default value is `false`.",
					},
					"kiosk_mode_allow_voice_control_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow the user to toggle voice control in kiosk mode. The _provider_ default value is `false`.",
					},
					"kiosk_mode_allow_voice_over_settings": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow access to the voice over settings while in kiosk mode. The _provider_ default value is `false`.",
					},
					"kiosk_mode_allow_zoom_settings": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow access to the zoom settings while in kiosk mode. The _provider_ default value is `false`.",
					},
					"kiosk_mode_app_store_url": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "URL in the app store to the app to use for kiosk mode. Use if KioskModeManagedAppId is not known.",
					},
					"kiosk_mode_app_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "appStoreApp", "managedApp", "builtInApp"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Type of app to run in kiosk mode. / App source options for iOS kiosk mode; possible values are: `notConfigured` (Device default value, no intent.), `appStoreApp` (The app to be run comes from the app store.), `managedApp` (The app to be run is built into the device.), `builtInApp` (The app to be run is a managed app.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"kiosk_mode_block_auto_lock": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block device auto lock while in kiosk mode. The _provider_ default value is `false`.",
					},
					"kiosk_mode_block_ringer_switch": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block use of the ringer switch while in kiosk mode. The _provider_ default value is `false`.",
					},
					"kiosk_mode_block_screen_rotation": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block screen rotation while in kiosk mode. The _provider_ default value is `false`.",
					},
					"kiosk_mode_block_sleep_button": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block use of the sleep button while in kiosk mode. The _provider_ default value is `false`.",
					},
					"kiosk_mode_block_touchscreen": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block use of the touchscreen while in kiosk mode. The _provider_ default value is `false`.",
					},
					"kiosk_mode_block_volume_buttons": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the volume buttons while in Kiosk Mode. The _provider_ default value is `false`.",
					},
					"kiosk_mode_built_in_app_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "ID for built-in apps to use for kiosk mode. Used when KioskModeManagedAppId and KioskModeAppStoreUrl are not set.",
					},
					"kiosk_mode_enable_voice_control": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to enable voice control in kiosk mode. The _provider_ default value is `false`.",
					},
					"kiosk_mode_managed_app_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Managed app id of the app to use for kiosk mode. If KioskModeManagedAppId is specified then KioskModeAppStoreUrl will be ignored.",
					},
					"kiosk_mode_require_assistive_touch": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to require assistive touch while in kiosk mode. The _provider_ default value is `false`.",
					},
					"kiosk_mode_require_color_inversion": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to require color inversion while in kiosk mode. The _provider_ default value is `false`.",
					},
					"kiosk_mode_require_mono_audio": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to require mono audio while in kiosk mode. The _provider_ default value is `false`.",
					},
					"kiosk_mode_require_voice_over": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to require voice over while in kiosk mode. The _provider_ default value is `false`.",
					},
					"kiosk_mode_require_zoom": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to require zoom while in kiosk mode. The _provider_ default value is `false`.",
					},
					"lock_screen_block_control_center": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using control center on the lock screen. The _provider_ default value is `false`.",
					},
					"lock_screen_block_notification_view": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using the notification view on the lock screen. The _provider_ default value is `false`.",
					},
					"lock_screen_block_passbook": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using passbook when the device is locked. The _provider_ default value is `false`.",
					},
					"lock_screen_block_today_view": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using the Today View on the lock screen. The _provider_ default value is `false`.",
					},
					"managed_pasteboard_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Open-in management controls how people share data between unmanaged and managed apps. Setting this to true enforces copy/paste restrictions based on how you configured <b>Block viewing corporate documents in unmanaged apps </b> and <b> Block viewing non-corporate documents in corporate apps.</b>. The _provider_ default value is `false`.",
					},
					"media_content_rating_apps": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("allAllowed", "allBlocked", "agesAbove4", "agesAbove9", "agesAbove12", "agesAbove17"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("allAllowed")},
						Computed:            true,
						MarkdownDescription: "Media content rating settings for Apps. / Apps rating as in media content; possible values are: `allAllowed` (Default value, allow all apps content), `allBlocked` (Do not allow any apps content), `agesAbove4` (4+, age 4 and above), `agesAbove9` (9+, age 9 and above), `agesAbove12` (12+, age 12 and above), `agesAbove17` (17+, age 17 and above). The _provider_ default value is `\"allAllowed\"`.",
					},
					"media_content_rating_australia": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mediaContentRatingAustralia
							"movie_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "parentalGuidance", "mature", "agesAbove15", "agesAbove18"),
								},
								MarkdownDescription: "Movies rating selected for Australia. / Movies rating labels in Australia; possible values are: `allAllowed` (Default value, allow all movies content), `allBlocked` (Do not allow any movies content), `general` (The G classification is suitable for everyone), `parentalGuidance` (The PG recommends viewers under 15 with guidance from parents or guardians), `mature` (The M classification is not recommended for viewers under 15), `agesAbove15` (The MA15+ classification is not suitable for viewers under 15), `agesAbove18` (The R18+ classification is not suitable for viewers under 18)",
							},
							"tv_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "preschoolers", "children", "general", "parentalGuidance", "mature", "agesAbove15", "agesAbove15AdultViolence"),
								},
								MarkdownDescription: "TV rating selected for Australia. / TV content rating labels in Australia; possible values are: `allAllowed` (Default value, allow all TV shows content), `allBlocked` (Do not allow any TV shows content), `preschoolers` (The P classification is intended for preschoolers), `children` (The C classification is intended for children under 14), `general` (The G classification is suitable for all ages), `parentalGuidance` (The PG classification is recommended for young viewers), `mature` (The M classification is recommended for viewers over 15), `agesAbove15` (The MA15+ classification is not suitable for viewers under 15), `agesAbove15AdultViolence` (The AV15+ classification is not suitable for viewers under 15, adult violence-specific)",
							},
						},
						MarkdownDescription: "Media content rating settings for Australia / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediacontentratingaustralia?view=graph-rest-beta",
					},
					"media_content_rating_canada": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mediaContentRatingCanada
							"movie_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "parentalGuidance", "agesAbove14", "agesAbove18", "restricted"),
								},
								MarkdownDescription: "Movies rating selected for Canada. / Movies rating labels in Canada; possible values are: `allAllowed` (Default value, allow all movies content), `allBlocked` (Do not allow any movies content), `general` (The G classification is suitable for all ages), `parentalGuidance` (The PG classification advises parental guidance), `agesAbove14` (The 14A classification is suitable for viewers above 14 or older), `agesAbove18` (The 18A classification is suitable for viewers above 18 or older), `restricted` (The R classification is restricted to 18 years and older)",
							},
							"tv_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "children", "childrenAbove8", "general", "parentalGuidance", "agesAbove14", "agesAbove18"),
								},
								MarkdownDescription: "TV rating selected for Canada. / TV content rating labels in Canada; possible values are: `allAllowed` (Default value, allow all TV shows content), `allBlocked` (Do not allow any TV shows content), `children` (The C classification is suitable for children ages of 2 to 7 years), `childrenAbove8` (The C8 classification is suitable for children ages 8+), `general` (The G classification is suitable for general audience), `parentalGuidance` (PG, Parental Guidance), `agesAbove14` (The 14+ classification is intended for viewers ages 14 and older), `agesAbove18` (The 18+ classification is intended for viewers ages 18 and older)",
							},
						},
						MarkdownDescription: "Media content rating settings for Canada / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediacontentratingcanada?view=graph-rest-beta",
					},
					"media_content_rating_france": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mediaContentRatingFrance
							"movie_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "agesAbove10", "agesAbove12", "agesAbove16", "agesAbove18"),
								},
								MarkdownDescription: "Movies rating selected for France. / Movies rating labels in France; possible values are: `allAllowed` (Default value, allow all movies content), `allBlocked` (Do not allow any movies content), `agesAbove10` (The 10 classification prohibits the screening of the film to minors under 10), `agesAbove12` (The 12 classification prohibits the screening of the film to minors under 12), `agesAbove16` (The 16 classification prohibits the screening of the film to minors under 16), `agesAbove18` (The 18 classification prohibits the screening to minors under 18)",
							},
							"tv_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "agesAbove10", "agesAbove12", "agesAbove16", "agesAbove18"),
								},
								MarkdownDescription: "TV rating selected for France. / TV content rating labels in France; possible values are: `allAllowed` (Default value, allow all TV shows content), `allBlocked` (Do not allow any TV shows content), `agesAbove10` (The -10 classification is not recommended for children under 10), `agesAbove12` (The -12 classification is not recommended for children under 12), `agesAbove16` (The -16 classification is not recommended for children under 16), `agesAbove18` (The -18 classification is not recommended for persons under 18)",
							},
						},
						MarkdownDescription: "Media content rating settings for France / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediacontentratingfrance?view=graph-rest-beta",
					},
					"media_content_rating_germany": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mediaContentRatingGermany
							"movie_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "agesAbove6", "agesAbove12", "agesAbove16", "adults"),
								},
								MarkdownDescription: "Movies rating selected for Germany. / Movies rating labels in Germany; possible values are: `allAllowed` (Default value, allow all movies content), `allBlocked` (Do not allow any movies content), `general` (Ab 0 Jahren, no age restrictions), `agesAbove6` (Ab 6 Jahren, ages 6 and older), `agesAbove12` (Ab 12 Jahren, ages 12 and older), `agesAbove16` (Ab 16 Jahren, ages 16 and older), `adults` (Ab 18 Jahren, adults only)",
							},
							"tv_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "agesAbove6", "agesAbove12", "agesAbove16", "adults"),
								},
								MarkdownDescription: "TV rating selected for Germany. / TV content rating labels in Germany; possible values are: `allAllowed` (Default value, allow all TV shows content), `allBlocked` (Do not allow any TV shows content), `general` (Ab 0 Jahren, no age restrictions), `agesAbove6` (Ab 6 Jahren, ages 6 and older), `agesAbove12` (Ab 12 Jahren, ages 12 and older), `agesAbove16` (Ab 16 Jahren, ages 16 and older), `adults` (Ab 18 Jahren, adults only)",
							},
						},
						MarkdownDescription: "Media content rating settings for Germany / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediacontentratinggermany?view=graph-rest-beta",
					},
					"media_content_rating_ireland": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mediaContentRatingIreland
							"movie_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "parentalGuidance", "agesAbove12", "agesAbove15", "agesAbove16", "adults"),
								},
								MarkdownDescription: "Movies rating selected for Ireland. / Movies rating labels in Ireland; possible values are: `allAllowed` (Default value, allow all movies content), `allBlocked` (Do not allow any movies content), `general` (Suitable for children of school going age), `parentalGuidance` (The PG classification advises parental guidance), `agesAbove12` (The 12A classification is suitable for viewers of 12 or older), `agesAbove15` (The 15A classification is suitable for viewers of 15 or older), `agesAbove16` (The 16 classification is suitable for viewers of 16 or older), `adults` (The 18 classification, suitable only for adults)",
							},
							"tv_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "children", "youngAdults", "parentalSupervision", "mature"),
								},
								MarkdownDescription: "TV rating selected for Ireland. / TV content rating labels in Ireland; possible values are: `allAllowed` (Default value, allow all TV shows content), `allBlocked` (Do not allow any TV shows content), `general` (The GA classification is suitable for all audiences), `children` (The CH classification is suitable for children), `youngAdults` (The YA classification is suitable for teenage audience), `parentalSupervision` (The PS classification invites parents and guardians to consider restriction children’s access), `mature` (The MA classification is suitable for adults)",
							},
						},
						MarkdownDescription: "Media content rating settings for Ireland / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediacontentratingireland?view=graph-rest-beta",
					},
					"media_content_rating_japan": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mediaContentRatingJapan
							"movie_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "parentalGuidance", "agesAbove15", "agesAbove18"),
								},
								MarkdownDescription: "Movies rating selected for Japan. / Movies rating labels in Japan; possible values are: `allAllowed` (Default value, allow all movies content), `allBlocked` (Do not allow any movies content), `general` (Suitable for all ages), `parentalGuidance` (The PG-12 classification requests parental guidance for young people under 12), `agesAbove15` (The R15+ classification is suitable for viewers of 15 or older), `agesAbove18` (The R18+ classification is suitable for viewers of 18 or older)",
							},
							"tv_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "explicitAllowed"),
								},
								MarkdownDescription: "TV rating selected for Japan. / TV content rating labels in Japan; possible values are: `allAllowed` (Default value, allow all TV shows content), `allBlocked` (Do not allow any TV shows content), `explicitAllowed` (All TV content is explicitly allowed)",
							},
						},
						MarkdownDescription: "Media content rating settings for Japan / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediacontentratingjapan?view=graph-rest-beta",
					},
					"media_content_rating_new_zealand": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mediaContentRatingNewZealand
							"movie_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "parentalGuidance", "mature", "agesAbove13", "agesAbove15", "agesAbove16", "agesAbove18", "restricted", "agesAbove16Restricted"),
								},
								MarkdownDescription: "Movies rating selected for New Zealand. / Movies rating labels in New Zealand; possible values are: `allAllowed` (Default value, allow all movies content), `allBlocked` (Do not allow any movies content), `general` (Suitable for general audience), `parentalGuidance` (The PG classification recommends parental guidance), `mature` (The M classification is suitable for mature audience), `agesAbove13` (The R13 classification is restricted to persons 13 years and over), `agesAbove15` (The R15 classification is restricted to persons 15 years and over), `agesAbove16` (The R16 classification is restricted to persons 16 years and over), `agesAbove18` (The R18 classification is restricted to persons 18 years and over), `restricted` (The R classification is restricted to a certain audience), `agesAbove16Restricted` (The RP16 classification requires viewers under 16 accompanied by a parent or an adult)",
							},
							"tv_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "parentalGuidance", "adults"),
								},
								MarkdownDescription: "TV rating selected for New Zealand. / TV content rating labels in New Zealand; possible values are: `allAllowed` (Default value, allow all TV shows content), `allBlocked` (Do not allow any TV shows content), `general` (The G classification excludes materials likely to harm children under 14), `parentalGuidance` (The PGR classification encourages parents and guardians to supervise younger viewers), `adults` (The AO classification is not suitable for children)",
							},
						},
						MarkdownDescription: "Media content rating settings for New Zealand / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediacontentratingnewzealand?view=graph-rest-beta",
					},
					"media_content_rating_united_kingdom": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mediaContentRatingUnitedKingdom
							"movie_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "universalChildren", "parentalGuidance", "agesAbove12Video", "agesAbove12Cinema", "agesAbove15", "adults"),
								},
								MarkdownDescription: "Movies rating selected for United Kingdom. / Movies rating labels in United Kingdom; possible values are: `allAllowed` (Default value, allow all movies content), `allBlocked` (Do not allow any movies content), `general` (The U classification is suitable for all ages), `universalChildren` (The UC classification is suitable for pre-school children, an old rating label), `parentalGuidance` (The PG classification is suitable for mature), `agesAbove12Video` (12, video release suitable for 12 years and over), `agesAbove12Cinema` (12A, cinema release suitable for 12 years and over), `agesAbove15` (15, suitable only for 15 years and older), `adults` (Suitable only for adults)",
							},
							"tv_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "caution"),
								},
								MarkdownDescription: "TV rating selected for United Kingdom. / TV content rating labels in United Kingdom; possible values are: `allAllowed` (Default value, allow all TV shows content), `allBlocked` (Do not allow any TV shows content), `caution` (Allowing TV contents with a warning message)",
							},
						},
						MarkdownDescription: "Media content rating settings for United Kingdom / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediacontentratingunitedkingdom?view=graph-rest-beta",
					},
					"media_content_rating_united_states": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mediaContentRatingUnitedStates
							"movie_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "parentalGuidance", "parentalGuidance13", "restricted", "adults"),
								},
								MarkdownDescription: "Movies rating selected for United States. / Movies rating labels in United States; possible values are: `allAllowed` (Default value, allow all movies content), `allBlocked` (Do not allow any movies content), `general` (G, all ages admitted), `parentalGuidance` (PG, some material may not be suitable for children), `parentalGuidance13` (PG13, some material may be inappropriate for children under 13), `restricted` (R, viewers under 17 require accompanying parent or adult guardian), `adults` (NC17, adults only)",
							},
							"tv_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "childrenAll", "childrenAbove7", "general", "parentalGuidance", "childrenAbove14", "adults"),
								},
								MarkdownDescription: "TV rating selected for United States. / TV content rating labels in United States; possible values are: `allAllowed` (Default value, allow all TV shows content), `allBlocked` (Do not allow any TV shows content), `childrenAll` (TV-Y, all children), `childrenAbove7` (TV-Y7, children age 7 and above), `general` (TV-G, suitable for all ages), `parentalGuidance` (TV-PG, parental guidance), `childrenAbove14` (TV-14, children age 14 and above), `adults` (TV-MA, adults only)",
							},
						},
						MarkdownDescription: "Media content rating settings for United States / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediacontentratingunitedstates?view=graph-rest-beta",
					},
					"messages_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using the Messages app on the supervised device. The _provider_ default value is `false`.",
					},
					"network_usage_rules": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // iosNetworkUsageRule
								"cellular_data_blocked": schema.BoolAttribute{
									Required:            true,
									MarkdownDescription: "If set to true, corresponding managed apps will not be allowed to use cellular data at any time.",
								},
								"cellular_data_block_when_roaming": schema.BoolAttribute{
									Required:            true,
									MarkdownDescription: "If set to true, corresponding managed apps will not be allowed to use cellular data when roaming.",
								},
								"managed_apps": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: deviceConfigurationAppListItemAttributes,
									},
									MarkdownDescription: "Information about the managed apps that this rule is going to apply to. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "List of managed apps and the network rules that applies to them. This collection can contain a maximum of 1000 elements. / Network Usage Rules allow enterprises to specify how managed apps use networks, such as cellular data networks. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosnetworkusagerule?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"nfc_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Disable NFC to prevent devices from pairing with other NFC-enabled devices. Available for iOS/iPadOS devices running 14.2 and later. The _provider_ default value is `false`.",
					},
					"notifications_block_settings_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow notifications settings modification (iOS 9.3 and later). The _provider_ default value is `false`.",
					},
					"on_device_only_dictation_forced": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Disables connections to Siri servers so that users can’t use Siri to dictate text. Available for devices running iOS and iPadOS versions 14.5 and later. The _provider_ default value is `false`.",
					},
					"on_device_only_translation_forced": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "When set to TRUE, the setting disables connections to Siri servers so that users can’t use Siri to translate text. When set to FALSE, the setting allows connections to to Siri servers to users can use Siri to translate text. Available for devices running iOS and iPadOS versions 15.0 and later. The _provider_ default value is `false`.",
					},
					"passcode_block_fingerprint_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Block modification of registered Touch ID fingerprints when in supervised mode. The _provider_ default value is `false`.",
					},
					"passcode_block_fingerprint_unlock": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block fingerprint unlock. The _provider_ default value is `false`.",
					},
					"passcode_block_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow passcode modification on the supervised device (iOS 9.0 and later). The _provider_ default value is `false`.",
					},
					"passcode_block_simple": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block simple passcodes. The _provider_ default value is `false`.",
					},
					"passcode_expiration_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days before the passcode expires. Valid values 1 to 65535",
					},
					"passcode_minimum_character_set_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of character sets a passcode must contain. Valid values 0 to 4",
					},
					"passcode_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minimum length of passcode. Valid values 4 to 14",
					},
					"passcode_minutes_of_inactivity_before_lock": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before a passcode is required.",
					},
					"passcode_minutes_of_inactivity_before_screen_timeout": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before the screen times out.",
					},
					"passcode_previous_passcode_block_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of previous passcodes to block. Valid values 1 to 24",
					},
					"passcode_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to require a passcode. The _provider_ default value is `false`.",
					},
					"passcode_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "alphanumeric", "numeric"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: "Type of passcode that is required. / Possible values of required passwords; possible values are: `deviceDefault` (Device default value, no intent.), `alphanumeric` (Alphanumeric password required.), `numeric` (Numeric password required.). The _provider_ default value is `\"deviceDefault\"`.",
					},
					"passcode_sign_in_failure_count_before_wipe": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of sign in failures allowed before wiping the device. Valid values 2 to 11",
					},
					"password_block_airdrop_sharing": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `passwordBlockAirDropSharing`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to block sharing passwords with the AirDrop passwords feature iOS 12.0 and later). The _provider_ default value is `false`.",
					},
					"password_block_auto_fill": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates if the AutoFill passwords feature is allowed (iOS 12.0 and later). The _provider_ default value is `false`.",
					},
					"password_block_proximity_requests": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block requesting passwords from nearby devices (iOS 12.0 and later). The _provider_ default value is `false`.",
					},
					"pki_block_ota_updates": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `pkiBlockOTAUpdates`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not over-the-air PKI updates are blocked. Setting this restriction to false does not disable CRL and OCSP checks (iOS 7.0 and later). The _provider_ default value is `false`.",
					},
					"podcasts_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using podcasts on the supervised device (iOS 8.0 and later). The _provider_ default value is `false`.",
					},
					"privacy_force_limit_ad_tracking": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates if ad tracking is limited.(iOS 7.0 and later). The _provider_ default value is `false`.",
					},
					"proximity_block_setup_to_new_device": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to enable the prompt to setup nearby devices with a supervised device. The _provider_ default value is `false`.",
					},
					"safari_block_autofill": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using Auto fill in Safari. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.",
					},
					"safari_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using Safari. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.",
					},
					"safari_block_javascript": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `safariBlockJavaScript`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to block JavaScript in Safari. The _provider_ default value is `false`.",
					},
					"safari_block_popups": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block popups in Safari. The _provider_ default value is `false`.",
					},
					"safari_cookie_settings": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("browserDefault", "blockAlways", "allowCurrentWebSite", "allowFromWebsitesVisited", "allowAlways"),
						},
						PlanModifiers: []planmodifier.String{
							wpdefaultvalue.StringDefaultValue("browserDefault"),
						},
						Computed:            true,
						MarkdownDescription: "Cookie settings for Safari. / Web Browser Cookie Settings; possible values are: `browserDefault` (Browser default value, no intent.), `blockAlways` (Always block cookies.), `allowCurrentWebSite` (Allow cookies from current Web site.), `allowFromWebsitesVisited` (Allow Cookies from websites visited.), `allowAlways` (Always allow cookies.). The _provider_ default value is `\"browserDefault\"`.",
					},
					"safari_managed_domains": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "URLs matching the patterns listed here will be considered managed. The _provider_ default value is `[]`.",
					},
					"safari_password_auto_fill_domains": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Users can save passwords in Safari only from URLs matching the patterns listed here. Applies to devices in supervised mode (iOS 9.3 and later). The _provider_ default value is `[]`.",
					},
					"safari_require_fraud_warning": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to require fraud warning in Safari. The _provider_ default value is `false`.",
					},
					"screen_capture_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from taking Screenshots. The _provider_ default value is `false`.",
					},
					"shared_device_block_temporary_sessions": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block temporary sessions on Shared iPads (iOS 13.4 or later). The _provider_ default value is `false`.",
					},
					"siri_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using Siri. The _provider_ default value is `false`.",
					},
					"siri_blocked_when_locked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using Siri when locked. The _provider_ default value is `false`.",
					},
					"siri_block_user_generated_content": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block Siri from querying user-generated content when used on a supervised device. The _provider_ default value is `false`.",
					},
					"siri_require_profanity_filter": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to prevent Siri from dictating, or speaking profane language on supervised device. The _provider_ default value is `false`.",
					},
					"software_updates_enforced_delay_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Sets how many days a software update will be delyed for a supervised device. Valid values 0 to 90",
					},
					"software_updates_force_delayed": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to delay user visibility of software updates when the device is in supervised mode. The _provider_ default value is `false`.",
					},
					"spotlight_block_internet_results": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block Spotlight search from returning internet results on supervised device. The _provider_ default value is `false`.",
					},
					"unpaired_external_boot_to_recovery_allowed": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Allow users to boot devices into recovery mode with unpaired devices. Available for devices running iOS and iPadOS versions 14.5 and later. The _provider_ default value is `false`.",
					},
					"usb_restricted_mode_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates if connecting to USB accessories while the device is locked is allowed (iOS 11.4.1 and later). The _provider_ default value is `false`.",
					},
					"voice_dialing_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block voice dialing. The _provider_ default value is `false`.",
					},
					"vpn_block_creation": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not the creation of VPN configurations is blocked (iOS 11.0 and later). The _provider_ default value is `false`.",
					},
					"wallpaper_block_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow wallpaper modification on supervised device (iOS 9.0 and later) . The _provider_ default value is `false`.",
					},
					"wifi_connect_only_to_configured_networks": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `wiFiConnectOnlyToConfiguredNetworks`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to force the device to use only Wi-Fi networks from configuration profiles when the device is in supervised mode. Available for devices running iOS and iPadOS versions 14.4 and earlier. Devices running 14.5+ should use the setting, “WiFiConnectToAllowedNetworksOnlyForced. The _provider_ default value is `false`.",
					},
					"wifi_connect_to_allowed_networks_only_forced": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `wiFiConnectToAllowedNetworksOnlyForced`, // custom MS Graph attribute name
						MarkdownDescription: "Require devices to use Wi-Fi networks set up via configuration profiles. Available for devices running iOS and iPadOS versions 14.5 and later. The _provider_ default value is `false`.",
					},
					"wifi_power_on_forced": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not Wi-Fi remains on, even when device is in airplane mode. Available for devices running iOS and iPadOS, versions 13.0 and later. The _provider_ default value is `false`.",
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: "This topic provides descriptions of the declared methods, properties and relationships exposed by the iosGeneralDeviceConfiguration resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosgeneraldeviceconfiguration?view=graph-rest-beta",
			},
		},
		"ios_update": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.iosUpdateConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // iosUpdateConfiguration
					"active_hours_end": schema.StringAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.String{
							wpdefaultvalue.StringDefaultValue("00:00:00.0000000"),
						},
						Computed:            true,
						MarkdownDescription: "Active Hours End (active hours mean the time window when updates install should not happen). The _provider_ default value is `\"00:00:00.0000000\"`.",
					},
					"active_hours_start": schema.StringAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.String{
							wpdefaultvalue.StringDefaultValue("00:00:00.0000000"),
						},
						Computed:            true,
						MarkdownDescription: "Active Hours Start (active hours mean the time window when updates install should not happen). The _provider_ default value is `\"00:00:00.0000000\"`.",
					},
					"custom_update_time_windows": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // customUpdateTimeWindow
								"end_day": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf("sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"),
									},
									MarkdownDescription: "End day of the time window. / Possible values for a weekday; possible values are: `sunday` (Sunday.), `monday` (Monday.), `tuesday` (Tuesday.), `wednesday` (Wednesday.), `thursday` (Thursday.), `friday` (Friday.), `saturday` (Saturday.)",
								},
								"end_time": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "End time of the time window",
								},
								"start_day": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf("sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"),
									},
									MarkdownDescription: "Start day of the time window. / Possible values for a weekday; possible values are: `sunday` (Sunday.), `monday` (Monday.), `tuesday` (Tuesday.), `wednesday` (Wednesday.), `thursday` (Thursday.), `friday` (Friday.), `saturday` (Saturday.)",
								},
								"start_time": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "Start time of the time window",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "If update schedule type is set to use time window scheduling, custom time windows when updates will be scheduled. This collection can contain a maximum of 20 elements. / Custom update time window / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-customupdatetimewindow?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"desired_os_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "If left unspecified, devices will update to the latest version of the OS.",
					},
					"enforced_software_update_delay_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Days before software updates are visible to iOS devices ranging from 0 to 90 inclusive",
					},
					"is_enabled": schema.BoolAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
						MarkdownDescription: "Is setting enabled in UI",
					},
					"scheduled_install_days": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(
								stringvalidator.OneOf("sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"),
							),
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Days in week for which active hours are configured. This collection can contain a maximum of 7 elements. / Possible values for a weekday; possible values are: `sunday` (Sunday.), `monday` (Monday.), `tuesday` (Tuesday.), `wednesday` (Wednesday.), `thursday` (Thursday.), `friday` (Friday.), `saturday` (Saturday.). The _provider_ default value is `[]`.",
					},
					"update_schedule_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("updateOutsideOfActiveHours", "alwaysUpdate", "updateDuringTimeWindows", "updateOutsideOfTimeWindows"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("alwaysUpdate")},
						Computed:            true,
						MarkdownDescription: "Update schedule type. / Update schedule type for iOS software updates; possible values are: `updateOutsideOfActiveHours` (Update outside of active hours.), `alwaysUpdate` (Always update.), `updateDuringTimeWindows` (Update during time windows.), `updateOutsideOfTimeWindows` (Update outside of time windows.). The _provider_ default value is `\"alwaysUpdate\"`.",
					},
					"utc_time_offset_in_minutes": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "UTC Time Offset indicated in minutes",
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: "IOS Update Configuration, allows you to configure time window within week to install iOS updates / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosupdateconfiguration?view=graph-rest-beta",
			},
		},
		"ios_vpn": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.iosVpnConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // iosVpnConfiguration
					"associated_domains": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Associated Domains. The _provider_ default value is `[]`.",
					},
					"authentication_method": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("certificate", "usernameAndPassword", "sharedSecret", "derivedCredential", "azureAD"),
						},
						MarkdownDescription: "Authentication method for this VPN connection. / VPN Authentication Method; possible values are: `certificate` (Authenticate with a certificate.), `usernameAndPassword` (Use username and password for authentication.), `sharedSecret` (Use Shared Secret for Authentication.  Only valid for iOS IKEv2.), `derivedCredential` (Use Derived Credential for Authentication.), `azureAD` (Use Azure AD for authentication.)",
					},
					"connection_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Connection name displayed to the user.",
					},
					"connection_type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("ciscoAnyConnect", "pulseSecure", "f5EdgeClient", "dellSonicWallMobileConnect", "checkPointCapsuleVpn", "customVpn", "ciscoIPSec", "citrix", "ciscoAnyConnectV2", "paloAltoGlobalProtect", "zscalerPrivateAccess", "f5Access2018", "citrixSso", "paloAltoGlobalProtectV2", "ikEv2", "alwaysOn", "microsoftTunnel", "netMotionMobility", "microsoftProtect"),
						},
						MarkdownDescription: "Connection type. / Apple VPN connection type; possible values are: `ciscoAnyConnect` (Cisco AnyConnect.), `pulseSecure` (Pulse Secure.), `f5EdgeClient` (F5 Edge Client.), `dellSonicWallMobileConnect` (Dell SonicWALL Mobile Connection.), `checkPointCapsuleVpn` (Check Point Capsule VPN.), `customVpn` (Custom VPN.), `ciscoIPSec` (Cisco (IPSec).), `citrix` (Citrix.), `ciscoAnyConnectV2` (Cisco AnyConnect V2.), `paloAltoGlobalProtect` (Palo Alto Networks GlobalProtect.), `zscalerPrivateAccess` (Zscaler Private Access.), `f5Access2018` (F5 Access 2018.), `citrixSso` (Citrix Sso.), `paloAltoGlobalProtectV2` (Palo Alto Networks GlobalProtect V2.), `ikEv2` (IKEv2.), `alwaysOn` (AlwaysOn.), `microsoftTunnel` (Microsoft Tunnel.), `netMotionMobility` (NetMotion Mobility.), `microsoftProtect` (Microsoft Protect.)",
					},
					"custom_data": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // keyValue
								"key": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Key.",
								},
								"value": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Value.",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Custom data when connection type is set to Custom VPN. Use this field to enable functionality not supported by Intune, but available in your VPN solution. Contact your VPN vendor to learn how to add these key/value pairs. This collection can contain a maximum of 25 elements. / Key Value definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyvalue?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"disable_on_demand_user_override": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Toggle to prevent user from disabling automatic VPN in the Settings app",
					},
					"disconnect_on_idle": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether to disconnect after on-demand connection idles",
					},
					"disconnect_on_idle_timer_in_seconds": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The length of time in seconds to wait before disconnecting an on-demand connection. Valid values 0 to 65535",
					},
					"enable_per_app": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Setting this to true creates Per-App VPN payload which can later be associated with Apps that can trigger this VPN conneciton on the end user's iOS device.",
					},
					"enable_split_tunneling": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Send all network traffic through VPN. The _provider_ default value is `false`.",
					},
					"excluded_domains": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Domains that are accessed through the public internet instead of through VPN, even when per-app VPN is activated. The _provider_ default value is `[]`.",
					},
					"identifier": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Identifier provided by VPN vendor when connection type is set to Custom VPN. For example: Cisco AnyConnect uses an identifier of the form com.cisco.anyconnect.applevpn.plugin",
					},
					"login_group_or_domain": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Login group or domain when connection type is set to Dell SonicWALL Mobile Connection.",
					},
					"on_demand_rules": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // vpnOnDemandRule
								"action": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf("connect", "evaluateConnection", "ignore", "disconnect"),
									},
									MarkdownDescription: "Action. / VPN On-Demand Rule Connection Action; possible values are: `connect` (Connect.), `evaluateConnection` (Evaluate Connection.), `ignore` (Ignore.), `disconnect` (Disconnect.)",
								},
								"dns_search_domains": schema.SetAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: "DNS Search Domains. The _provider_ default value is `[]`.",
								},
								"dns_server_address_match": schema.SetAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: "DNS Search Server Address. The _provider_ default value is `[]`.",
								},
								"domain_action": schema.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.OneOf("connectIfNeeded", "neverConnect"),
									},
									PlanModifiers: []planmodifier.String{
										wpdefaultvalue.StringDefaultValue("connectIfNeeded"),
									},
									Computed:            true,
									MarkdownDescription: "Domain Action (Only applicable when Action is evaluate connection). / VPN On-Demand Rule Connection Domain Action; possible values are: `connectIfNeeded` (Connect if needed.), `neverConnect` (Never connect.). The _provider_ default value is `\"connectIfNeeded\"`.",
								},
								"domains": schema.SetAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: "Domains (Only applicable when Action is evaluate connection). The _provider_ default value is `[]`.",
								},
								"interface_type_match": schema.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.OneOf("notConfigured", "ethernet", "wiFi", "cellular"),
									},
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
									Computed:            true,
									MarkdownDescription: "Network interface to trigger VPN. / VPN On-Demand Rule Connection network interface type; possible values are: `notConfigured` (NotConfigured), `ethernet` (Ethernet.), `wiFi` (WiFi.), `cellular` (Cellular.). The _provider_ default value is `\"notConfigured\"`.",
								},
								"probe_required_url": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Probe Required Url (Only applicable when Action is evaluate connection and DomainAction is connect if needed).",
								},
								"probe_url": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "A URL to probe. If this URL is successfully fetched (returning a 200 HTTP status code) without redirection, this rule matches.",
								},
								"ssids": schema.SetAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: "Network Service Set Identifiers (SSIDs). The _provider_ default value is `[]`.",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "On-Demand Rules. This collection can contain a maximum of 500 elements. / VPN On-Demand Rule definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-vpnondemandrule?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"opt_in_to_device_id_sharing": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Opt-In to sharing the device's Id to third-party vpn clients for use during network access control validation.",
					},
					"provider_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "appProxy", "packetTunnel"),
						},
						MarkdownDescription: "Provider type for per-app VPN. / Provider type for per-app VPN; possible values are: `notConfigured` (Tunnel traffic is not explicitly configured.), `appProxy` (Tunnel traffic at the application layer.), `packetTunnel` (Tunnel traffic at the IP layer.)",
					},
					"proxy_server": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // vpnProxyServer
							"address": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Address.",
							},
							"automatic_configuration_script_url": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Proxy's automatic configuration script url.",
							},
							"port": schema.Int64Attribute{
								Optional:            true,
								MarkdownDescription: "Port. Valid values 0 to 65535",
							},
						},
						MarkdownDescription: "Proxy Server. / VPN Proxy Server. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-vpnproxyserver?view=graph-rest-beta",
					},
					"realm": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Realm when connection type is set to Pulse Secure.",
					},
					"role": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Role when connection type is set to Pulse Secure.",
					},
					"safari_domains": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Safari domains when this VPN per App setting is enabled. In addition to the apps associated with this VPN, Safari domains specified here will also be able to trigger this VPN connection. The _provider_ default value is `[]`.",
					},
					"server": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{ // vpnServer
							"address": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Address (IP address, FQDN or URL)",
							},
							"description": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Description.",
							},
							"is_default_server": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
								Computed:            true,
								MarkdownDescription: "Default server. The _provider_ default value is `true`.",
							},
						},
						MarkdownDescription: "VPN Server on the network. Make sure end users can access this network location. / VPN Server definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-vpnserver?view=graph-rest-beta",
					},
					"cloud_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Zscaler only. Zscaler cloud which the user is assigned to.",
					},
					"exclude_list": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Zscaler only. List of network addresses which are not sent through the Zscaler cloud. The _provider_ default value is `[]`.",
					},
					"microsoft_tunnel_site_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Microsoft Tunnel site ID.",
					},
					"strict_enforcement": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Zscaler only. Blocks network traffic until the user signs into Zscaler app. \"True\" means traffic is blocked.",
					},
					"targeted_mobile_apps": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationAppListItemAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Targeted mobile apps. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"user_domain": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Zscaler only. Enter a static domain to pre-populate the login field with in the Zscaler app. If this is left empty, the user's Azure Active Directory domain will be used instead.",
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: "By providing the configurations in this profile you can instruct the iOS device to connect to desired VPN endpoint. By specifying the authentication method and security types expected by VPN endpoint you can make the VPN connection seamless for end user. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosvpnconfiguration?view=graph-rest-beta",
			},
		},
		"macos_custom": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.macOSCustomConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // macOSCustomConfiguration
					"deployment_channel": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceChannel", "userChannel", "unknownFutureValue"),
						},
						MarkdownDescription: "Indicates the channel used to deploy the configuration profile. Available choices are DeviceChannel, UserChannel. / Indicates the channel used to deploy the configuration profile. Available choices are DeviceChannel, UserChannel; possible values are: `deviceChannel` (Send payload down over Device Channel.), `userChannel` (Send payload down over User Channel.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.)",
					},
					"payload": schema.StringAttribute{
						Required:            true,
						Validators:          []validator.String{wpvalidator.TranslateValueEncodeBase64()},
						MarkdownDescription: "Payload. (UTF8 encoded byte array)",
					},
					"file_name": schema.StringAttribute{
						Required:            true,
						Description:         `payloadFileName`, // custom MS Graph attribute name
						MarkdownDescription: "Payload file name (*.mobileconfig",
					},
					"name": schema.StringAttribute{
						Required:            true,
						Description:         `payloadName`, // custom MS Graph attribute name
						MarkdownDescription: "Name that is displayed to the user.",
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: "This topic provides descriptions of the declared methods, properties and relationships exposed by the macOSCustomConfiguration resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macoscustomconfiguration?view=graph-rest-beta",
			},
		},
		"macos_custom_app": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.macOSCustomAppConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // macOSCustomAppConfiguration
					"bundle_id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Bundle id for targeting.",
					},
					"configuration_xml": schema.StringAttribute{
						Required:            true,
						Validators:          []validator.String{wpvalidator.TranslateValueEncodeBase64()},
						MarkdownDescription: "Configuration xml. (UTF8 encoded byte array)",
					},
					"file_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Configuration file name (*.plist",
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: "This topic provides descriptions of the declared methods, properties and relationships exposed by the macOSCustomAppConfiguration resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macoscustomappconfiguration?view=graph-rest-beta",
			},
		},
		"macos_device_features": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.macOSDeviceFeaturesConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // macOSDeviceFeaturesConfiguration
					"airprint_destinations": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationAirPrintDestinationAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						Description:         `airPrintDestinations`, // custom MS Graph attribute name
						MarkdownDescription: "An array of AirPrint printers that should always be shown. This collection can contain a maximum of 500 elements. / Represents an AirPrint destination. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-airprintdestination?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"admin_show_host_info": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether to show admin host information on the login window. The _provider_ default value is `false`.",
					},
					"app_associated_domains": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // macOSAssociatedDomainsItem
								"application_identifier": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The application identifier of the app to associate domains with.",
								},
								"direct_downloads_enabled": schema.BoolAttribute{
									Required:            true,
									MarkdownDescription: "Determines whether data should be downloaded directly or via a CDN.",
								},
								"domains": schema.SetAttribute{
									ElementType:         types.StringType,
									Required:            true,
									MarkdownDescription: "The list of domains to associate.",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Gets or sets a list that maps apps to their associated domains. Application identifiers must be unique. This collection can contain a maximum of 500 elements. / A mapping of application identifiers to associated domains. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macosassociateddomainsitem?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"authorized_users_list_hidden": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether to show the name and password dialog or a list of users on the login window. The _provider_ default value is `false`.",
					},
					"authorized_users_list_hide_admin_users": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether to hide admin users in the authorized users list on the login window. The _provider_ default value is `false`.",
					},
					"authorized_users_list_hide_local_users": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether to show only network and system users in the authorized users list on the login window. The _provider_ default value is `false`.",
					},
					"authorized_users_list_hide_mobile_accounts": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether to hide mobile users in the authorized users list on the login window. The _provider_ default value is `false`.",
					},
					"authorized_users_list_include_network_users": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether to show network users in the authorized users list on the login window. The _provider_ default value is `false`.",
					},
					"authorized_users_list_show_other_managed_users": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether to show other users in the authorized users list on the login window. The _provider_ default value is `false`.",
					},
					"auto_launch_items": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // macOSLaunchItem
								"hide": schema.BoolAttribute{
									Required:            true,
									MarkdownDescription: "Whether or not to hide the item from the Users and Groups List.",
								},
								"path": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "Path to the launch item.",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "List of applications, files, folders, and other items to launch when the user logs in. This collection can contain a maximum of 500 elements. / Represents an app in the list of macOS launch items / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macoslaunchitem?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"console_access_disabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether the Other user will disregard use of the `console` special user name. The _provider_ default value is `false`.",
					},
					"content_caching_block_deletion": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Prevents content caches from purging content to free up disk space for other apps. The _provider_ default value is `false`.",
					},
					"content_caching_client_listen_ranges": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationIpRangeAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "A list of custom IP ranges content caches will use to listen for clients. This collection can contain a maximum of 500 elements. / IP range base class for representing IPV4, IPV6 address ranges / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-iprange?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"content_caching_client_policy": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "clientsInLocalNetwork", "clientsWithSamePublicIpAddress", "clientsInCustomLocalNetworks", "clientsInCustomLocalNetworksWithFallback"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Determines the method in which content caching servers will listen for clients. / Determines which clients a content cache will serve; possible values are: `notConfigured` (Defaults to clients in local network.), `clientsInLocalNetwork` (Content caches will provide content to devices only in their immediate local network.), `clientsWithSamePublicIpAddress` (Content caches will provide content to devices that share the same public IP address.), `clientsInCustomLocalNetworks` (Content caches will provide content to devices in contentCachingClientListenRanges.), `clientsInCustomLocalNetworksWithFallback` (Content caches will provide content to devices in contentCachingClientListenRanges, contentCachingPeerListenRanges, and contentCachingParents.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"content_caching_data_path": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The path to the directory used to store cached content. The value must be (or end with) /Library/Application Support/Apple/AssetCache/Data",
					},
					"content_caching_disable_connection_sharing": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Disables internet connection sharing. The _provider_ default value is `false`.",
					},
					"content_caching_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Enables content caching and prevents it from being disabled by the user. The _provider_ default value is `false`.",
					},
					"content_caching_force_connection_sharing": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Forces internet connection sharing. contentCachingDisableConnectionSharing overrides this setting. The _provider_ default value is `false`.",
					},
					"content_caching_keep_awake": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Prevent the device from sleeping if content caching is enabled. The _provider_ default value is `false`.",
					},
					"content_caching_log_client_identities": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Enables logging of IP addresses and ports of clients that request cached content. The _provider_ default value is `false`.",
					},
					"content_caching_max_size_bytes": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The maximum number of bytes of disk space that will be used for the content cache. A value of 0 (default) indicates unlimited disk space.",
					},
					"content_caching_parents": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "A list of IP addresses representing parent content caches. The _provider_ default value is `[]`.",
					},
					"content_caching_parent_selection_policy": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "roundRobin", "firstAvailable", "urlPathHash", "random", "stickyAvailable"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Determines the method in which content caching servers will select parents if multiple are present. / Determines how content caches select a parent cache; possible values are: `notConfigured` (Defaults to round-robin strategy.), `roundRobin` (Rotate through the parents in order. Use this policy for load balancing.), `firstAvailable` (Always use the first available parent in the Parents list. Use this policy to designate permanent primary, secondary, and subsequent parents.), `urlPathHash` (Hash the path part of the requested URL so that the same parent is always used for the same URL. This is useful for maximizing the size of the combined caches of the parents.), `random` (Choose a parent at random. Use this policy for load balancing.), `stickyAvailable` (Use the first available parent that is available in the Parents list until it becomes unavailable, then advance to the next one. Use this policy for designating floating primary, secondary, and subsequent parents.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"content_caching_peer_filter_ranges": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationIpRangeAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "A list of custom IP ranges content caches will use to query for content from peers caches. This collection can contain a maximum of 500 elements. / IP range base class for representing IPV4, IPV6 address ranges / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-iprange?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"content_caching_peer_listen_ranges": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationIpRangeAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "A list of custom IP ranges content caches will use to listen for peer caches. This collection can contain a maximum of 500 elements. / IP range base class for representing IPV4, IPV6 address ranges / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-iprange?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"content_caching_peer_policy": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "peersInLocalNetwork", "peersWithSamePublicIpAddress", "peersInCustomLocalNetworks"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Determines the method in which content caches peer with other caches. / Determines which content caches other content caches will peer with; possible values are: `notConfigured` (Defaults to peers in local network.), `peersInLocalNetwork` (Content caches will only peer with caches in their immediate local network.), `peersWithSamePublicIpAddress` (Content caches will only peer with caches that share the same public IP address.), `peersInCustomLocalNetworks` (Content caches will use contentCachingPeerFilterRanges and contentCachingPeerListenRanges to determine which caches to peer with.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"content_caching_port": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Sets the port used for content caching. If the value is 0, a random available port will be selected. Valid values 0 to 65535",
					},
					"content_caching_public_ranges": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationIpRangeAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "A list of custom IP ranges that Apple's content caching service should use to match clients to content caches. This collection can contain a maximum of 500 elements. / IP range base class for representing IPV4, IPV6 address ranges / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-iprange?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"content_caching_show_alerts": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Display content caching alerts as system notifications. The _provider_ default value is `false`.",
					},
					"content_caching_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "userContentOnly", "sharedContentOnly"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Determines what type of content is allowed to be cached by Apple's content caching service. / Indicates the type of content allowed to be cached by Apple's content caching service; possible values are: `notConfigured` (Default. Both user iCloud data and non-iCloud data will be cached.), `userContentOnly` (Allow Apple's content caching service to cache user iCloud data.), `sharedContentOnly` (Allow Apple's content caching service to cache non-iCloud data (e.g. app and software updates).). The _provider_ default value is `\"notConfigured\"`.",
					},
					"login_window_text": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							wpvalidator.TranslateGraphNullWithTfDefault(tftypes.String),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
						Computed:            true,
						MarkdownDescription: "Custom text to be displayed on the login window. The _provider_ default value is `\"\"`.",
					},
					"logout_disabled_while_logged_in": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `logOutDisabledWhileLoggedIn`, // custom MS Graph attribute name
						MarkdownDescription: "Whether the Log Out menu item on the login window will be disabled while the user is logged in. The _provider_ default value is `false`.",
					},
					"macos_single_sign_on_extension": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // macOSSingleSignOnExtension
							"azure_ad": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.macOSAzureAdSingleSignOnExtension",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // macOSAzureAdSingleSignOnExtension
										"bundle_id_access_control_list": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: "An optional list of additional bundle IDs allowed to use the AAD extension for single sign-on. The _provider_ default value is `[]`.",
										},
										"configurations": schema.SetNestedAttribute{
											Optional: true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: deviceConfigurationKeyTypedValuePairAttributes,
											},
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: "Gets or sets a list of typed key-value pairs used to configure Credential-type profiles. This collection can contain a maximum of 500 elements. / A key-value pair with a string key and a typed value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keytypedvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`.",
										},
										"enable_shared_device_mode": schema.BoolAttribute{
											Optional:            true,
											PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
											Computed:            true,
											MarkdownDescription: "Enables or disables shared device mode. The _provider_ default value is `false`.",
										},
									},
									Validators: []validator.Object{
										deviceConfigurationMacOSSingleSignOnExtensionValidator,
									},
									MarkdownDescription: "Represents an Azure AD-type Single Sign-On extension profile for macOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macosazureadsinglesignonextension?view=graph-rest-beta",
								},
							},
							"credential": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.macOSCredentialSingleSignOnExtension",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // macOSCredentialSingleSignOnExtension
										"configurations": schema.SetNestedAttribute{
											Optional: true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: deviceConfigurationKeyTypedValuePairAttributes,
											},
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: "Gets or sets a list of typed key-value pairs used to configure Credential-type profiles. This collection can contain a maximum of 500 elements. / A key-value pair with a string key and a typed value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keytypedvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`.",
										},
										"domains": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: "Gets or sets a list of hosts or domain names for which the app extension performs SSO. The _provider_ default value is `[]`.",
										},
										"extension_identifier": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "Gets or sets the bundle ID of the app extension that performs SSO for the specified URLs.",
										},
										"realm": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "Gets or sets the case-sensitive realm name for this profile.",
										},
										"team_identifier": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "Gets or sets the team ID of the app extension that performs SSO for the specified URLs.",
										},
									},
									Validators: []validator.Object{
										deviceConfigurationMacOSSingleSignOnExtensionValidator,
									},
									MarkdownDescription: "Represents a Credential-type Single Sign-On extension profile for macOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macoscredentialsinglesignonextension?view=graph-rest-beta",
								},
							},
							"kerberos": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.macOSKerberosSingleSignOnExtension",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // macOSKerberosSingleSignOnExtension
										"active_directory_site_code": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "Gets or sets the Active Directory site.",
										},
										"block_active_directory_site_auto_discovery": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: "Enables or disables whether the Kerberos extension can automatically determine its site name.",
										},
										"block_automatic_login": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: "Enables or disables Keychain usage.",
										},
										"cache_name": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "Gets or sets the Generic Security Services name of the Kerberos cache to use for this profile.",
										},
										"credential_bundle_id_access_control_list": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											MarkdownDescription: "Gets or sets a list of app Bundle IDs allowed to access the Kerberos Ticket Granting Ticket.",
										},
										"credentials_cache_monitored": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: "When set to True, the credential is requested on the next matching Kerberos challenge or network state change. When the credential is expired or missing, a new credential is created. Available for devices running macOS versions 12 and later.",
										},
										"domain_realms": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											MarkdownDescription: "Gets or sets a list of realms for custom domain-realm mapping. Realms are case sensitive.",
										},
										"domains": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											MarkdownDescription: "Gets or sets a list of hosts or domain names for which the app extension performs SSO.",
										},
										"is_default_realm": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: "When true, this profile's realm will be selected as the default. Necessary if multiple Kerberos-type profiles are configured.",
										},
										"kerberos_apps_in_bundle_id_acl_included": schema.BoolAttribute{
											Required:            true,
											Description:         `kerberosAppsInBundleIdACLIncluded`, // custom MS Graph attribute name
											MarkdownDescription: "When set to True, the Kerberos extension allows any apps entered with the app bundle ID, managed apps, and standard Kerberos utilities, such as TicketViewer and klist, to access and use the credential. Available for devices running macOS versions 12 and later.",
										},
										"managed_apps_in_bundle_id_acl_included": schema.BoolAttribute{
											Required:            true,
											Description:         `managedAppsInBundleIdACLIncluded`, // custom MS Graph attribute name
											MarkdownDescription: "When set to True, the Kerberos extension allows managed apps, and any apps entered with the app bundle ID to access the credential. When set to False, the Kerberos extension allows all apps to access the credential. Available for devices running iOS and iPadOS versions 14 and later.",
										},
										"mode_credential_used": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "Select how other processes use the Kerberos Extension credential.",
										},
										"password_block_modification": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: "Enables or disables password changes.",
										},
										"password_change_url": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "Gets or sets the URL that the user will be sent to when they initiate a password change.",
										},
										"password_enable_local_sync": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: "Enables or disables password syncing. This won't affect users logged in with a mobile account on macOS.",
										},
										"password_expiration_days": schema.Int64Attribute{
											Optional:            true,
											MarkdownDescription: "Overrides the default password expiration in days. For most domains, this value is calculated automatically.",
										},
										"password_expiration_notification_days": schema.Int64Attribute{
											Optional:            true,
											MarkdownDescription: "Gets or sets the number of days until the user is notified that their password will expire (default is 15).",
										},
										"password_minimum_age_days": schema.Int64Attribute{
											Optional:            true,
											MarkdownDescription: "Gets or sets the minimum number of days until a user can change their password again.",
										},
										"password_minimum_length": schema.Int64Attribute{
											Optional:            true,
											MarkdownDescription: "Gets or sets the minimum length of a password.",
										},
										"password_previous_password_block_count": schema.Int64Attribute{
											Optional:            true,
											MarkdownDescription: "Gets or sets the number of previous passwords to block.",
										},
										"password_require_active_directory_complexity": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: "Enables or disables whether passwords must meet Active Directory's complexity requirements.",
										},
										"password_requirements_description": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "Gets or sets a description of the password complexity requirements.",
										},
										"preferred_kdcs": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											Description:         `preferredKDCs`, // custom MS Graph attribute name
											MarkdownDescription: "Add creates an ordered list of preferred Key Distribution Centers (KDCs) to use for Kerberos traffic. This list is used when the servers are not discoverable using DNS. When the servers are discoverable, the list is used for both connectivity checks, and used first for Kerberos traffic. If the servers don’t respond, then the device uses DNS discovery. Delete removes an existing list, and devices use DNS discovery. Available for devices running macOS versions 12 and later.",
										},
										"realm": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "Gets or sets the case-sensitive realm name for this profile.",
										},
										"require_user_presence": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: "Gets or sets whether to require authentication via Touch ID, Face ID, or a passcode to access the keychain entry.",
										},
										"sign_in_help_text": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "Text displayed to the user at the Kerberos sign in window. Available for devices running iOS and iPadOS versions 14 and later.",
										},
										"tls_for_ldap_required": schema.BoolAttribute{
											Required:            true,
											Description:         `tlsForLDAPRequired`, // custom MS Graph attribute name
											MarkdownDescription: "When set to True, LDAP connections are required to use Transport Layer Security (TLS). Available for devices running macOS versions 11 and later.",
										},
										"username_label_custom": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "This label replaces the user name shown in the Kerberos extension. You can enter a name to match the name of your company or organization. Available for devices running macOS versions 11 and later.",
										},
										"user_principal_name": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "Gets or sets the principle user name to use for this profile. The realm name does not need to be included.",
										},
										"user_setup_delayed": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: "When set to True, the user isn’t prompted to set up the Kerberos extension until the extension is enabled by the admin, or a Kerberos challenge is received. Available for devices running macOS versions 11 and later.",
										},
									},
									Validators: []validator.Object{
										deviceConfigurationMacOSSingleSignOnExtensionValidator,
									},
									MarkdownDescription: "Represents a Kerberos-type Single Sign-On extension profile for macOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macoskerberossinglesignonextension?view=graph-rest-beta",
								},
							},
							"redirect": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.macOSRedirectSingleSignOnExtension",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // macOSRedirectSingleSignOnExtension
										"configurations": schema.SetNestedAttribute{
											Optional: true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: deviceConfigurationKeyTypedValuePairAttributes,
											},
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: "Gets or sets a list of typed key-value pairs used to configure Credential-type profiles. This collection can contain a maximum of 500 elements. / A key-value pair with a string key and a typed value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keytypedvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`.",
										},
										"extension_identifier": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "Gets or sets the bundle ID of the app extension that performs SSO for the specified URLs.",
										},
										"team_identifier": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "Gets or sets the team ID of the app extension that performs SSO for the specified URLs.",
										},
										"url_prefixes": schema.SetAttribute{
											ElementType:         types.StringType,
											Required:            true,
											MarkdownDescription: "One or more URL prefixes of identity providers on whose behalf the app extension performs single sign-on. URLs must begin with http:// or https://. All URL prefixes must be unique for all profiles.",
										},
									},
									Validators: []validator.Object{
										deviceConfigurationMacOSSingleSignOnExtensionValidator,
									},
									MarkdownDescription: "Represents a Redirect-type Single Sign-On extension profile for macOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macosredirectsinglesignonextension?view=graph-rest-beta",
								},
							},
						},
						Description:         `macOSSingleSignOnExtension`, // custom MS Graph attribute name
						MarkdownDescription: "Gets or sets a single sign-on extension profile. / An abstract base class for all macOS-specific single sign-on extension types. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macossinglesignonextension?view=graph-rest-beta",
					},
					"poweroff_disabled_while_logged_in": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `powerOffDisabledWhileLoggedIn`, // custom MS Graph attribute name
						MarkdownDescription: "Whether the Power Off menu item on the login window will be disabled while the user is logged in. The _provider_ default value is `false`.",
					},
					"restart_disabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether to hide the Restart button item on the login window. The _provider_ default value is `false`.",
					},
					"restart_disabled_while_logged_in": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether the Restart menu item on the login window will be disabled while the user is logged in. The _provider_ default value is `false`.",
					},
					"screen_lock_disable_immediate": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether to disable the immediate screen lock functions. The _provider_ default value is `false`.",
					},
					"shutdown_disabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `shutDownDisabled`, // custom MS Graph attribute name
						MarkdownDescription: "Whether to hide the Shut Down button item on the login window. The _provider_ default value is `false`.",
					},
					"shutdown_disabled_while_logged_in": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `shutDownDisabledWhileLoggedIn`, // custom MS Graph attribute name
						MarkdownDescription: "Whether the Shut Down menu item on the login window will be disabled while the user is logged in. The _provider_ default value is `false`.",
					},
					"sleep_disabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether to hide the Sleep menu item on the login window. The _provider_ default value is `false`.",
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: "MacOS device features configuration profile. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macosdevicefeaturesconfiguration?view=graph-rest-beta",
			},
		},
		"macos_extensions": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.macOSExtensionsConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // macOSExtensionsConfiguration
					"kernel_extension_allowed_team_identifiers": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "All kernel extensions validly signed by the team identifiers in this list will be allowed to load. The _provider_ default value is `[]`.",
					},
					"kernel_extension_overrides_allowed": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "If set to true, users can approve additional kernel extensions not explicitly allowed by configurations profiles. The _provider_ default value is `false`.",
					},
					"kernel_extensions_allowed": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // macOSKernelExtension
								"bundle_id": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "Bundle ID of the kernel extension.",
								},
								"team_identifier": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The team identifier that was used to sign the kernel extension.",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "A list of kernel extensions that will be allowed to load. . This collection can contain a maximum of 500 elements. / Represents a specific macOS kernel extension. A macOS kernel extension can be described by its team identifier plus its bundle identifier. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macoskernelextension?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"system_extensions_allowed": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // macOSSystemExtension
								"bundle_id": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "Gets or sets the bundle identifier of the system extension.",
								},
								"team_identifier": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Gets or sets the team identifier that was used to sign the system extension.",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Gets or sets a list of allowed macOS system extensions. This collection can contain a maximum of 500 elements. / Represents a specific macOS system extension. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macossystemextension?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"system_extensions_allowed_team_identifiers": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Gets or sets a list of allowed team identifiers. Any system extension signed with any of the specified team identifiers will be approved. The _provider_ default value is `[]`.",
					},
					"system_extensions_allowed_types": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // macOSSystemExtensionTypeMapping
								"allowed_types": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										wpvalidator.FlagEnumValues("driverExtensionsAllowed", "networkExtensionsAllowed", "endpointSecurityExtensionsAllowed"),
									},
									MarkdownDescription: "Gets or sets the allowed macOS system extension types. / Flag enum representing the allowed macOS system extension types; possible values are: `driverExtensionsAllowed` (Enables driver extensions.), `networkExtensionsAllowed` (Enables network extensions.), `endpointSecurityExtensionsAllowed` (Enables endpoint security extensions.)",
								},
								"team_identifier": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "Gets or sets the team identifier used to sign the system extension.",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Gets or sets a list of allowed macOS system extension types. This collection can contain a maximum of 500 elements. / Represents a mapping between team identifiers for macOS system extensions and system extension types. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macossystemextensiontypemapping?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"system_extensions_block_override": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Gets or sets whether to allow the user to approve additional system extensions not explicitly allowed by configuration profiles. The _provider_ default value is `false`.",
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: "MacOS extensions configuration profile. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macosextensionsconfiguration?view=graph-rest-beta",
			},
		},
		"macos_software_update": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.macOSSoftwareUpdateConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // macOSSoftwareUpdateConfiguration
					"all_other_update_behavior": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "default", "downloadOnly", "installASAP", "notifyOnly", "installLater"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Update behavior for all other updates. / Update behavior options for macOS software updates; possible values are: `notConfigured` (Not configured.), `default` (Download and/or install the software update, depending on the current device state.), `downloadOnly` (Download the software update without installing it.), `installASAP` (Install an already downloaded software update.), `notifyOnly` (Download the software update and notify the user via the App Store.), `installLater` (Download the software update and install it at a later time.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"config_data_update_behavior": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "default", "downloadOnly", "installASAP", "notifyOnly", "installLater"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Update behavior for configuration data file updates. / Update behavior options for macOS software updates; possible values are: `notConfigured` (Not configured.), `default` (Download and/or install the software update, depending on the current device state.), `downloadOnly` (Download the software update without installing it.), `installASAP` (Install an already downloaded software update.), `notifyOnly` (Download the software update and notify the user via the App Store.), `installLater` (Download the software update and install it at a later time.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"critical_update_behavior": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "default", "downloadOnly", "installASAP", "notifyOnly", "installLater"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Update behavior for critical updates. / Update behavior options for macOS software updates; possible values are: `notConfigured` (Not configured.), `default` (Download and/or install the software update, depending on the current device state.), `downloadOnly` (Download the software update without installing it.), `installASAP` (Install an already downloaded software update.), `notifyOnly` (Download the software update and notify the user via the App Store.), `installLater` (Download the software update and install it at a later time.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"custom_update_time_windows": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // customUpdateTimeWindow
								"end_day": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf("sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"),
									},
									MarkdownDescription: "End day of the time window. / Possible values for a weekday; possible values are: `sunday` (Sunday.), `monday` (Monday.), `tuesday` (Tuesday.), `wednesday` (Wednesday.), `thursday` (Thursday.), `friday` (Friday.), `saturday` (Saturday.)",
								},
								"end_time": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "End time of the time window",
								},
								"start_day": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf("sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"),
									},
									MarkdownDescription: "Start day of the time window. / Possible values for a weekday; possible values are: `sunday` (Sunday.), `monday` (Monday.), `tuesday` (Tuesday.), `wednesday` (Wednesday.), `thursday` (Thursday.), `friday` (Friday.), `saturday` (Saturday.)",
								},
								"start_time": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "Start time of the time window",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Custom Time windows when updates will be allowed or blocked. This collection can contain a maximum of 20 elements. / Custom update time window / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-customupdatetimewindow?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"firmware_update_behavior": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "default", "downloadOnly", "installASAP", "notifyOnly", "installLater"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Update behavior for firmware updates. / Update behavior options for macOS software updates; possible values are: `notConfigured` (Not configured.), `default` (Download and/or install the software update, depending on the current device state.), `downloadOnly` (Download the software update without installing it.), `installASAP` (Install an already downloaded software update.), `notifyOnly` (Download the software update and notify the user via the App Store.), `installLater` (Download the software update and install it at a later time.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"max_user_deferrals_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The maximum number of times the system allows the user to postpone an update before it’s installed. Supported values: 0 - 366. Valid values 0 to 365",
					},
					"priority": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("low", "high", "unknownFutureValue"),
						},
						MarkdownDescription: "The scheduling priority for downloading and preparing the requested update. Default: Low. Possible values: Null, Low, High. / The scheduling priority options for downloading and preparing the requested mac OS update; possible values are: `low` (Indicates low scheduling priority for downloading and preparing the requested update), `high` (Indicates high scheduling priority for downloading and preparing the requested update), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.)",
					},
					"update_schedule_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("alwaysUpdate", "updateDuringTimeWindows", "updateOutsideOfTimeWindows"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("alwaysUpdate")},
						Computed:            true,
						MarkdownDescription: "Update schedule type. / Update schedule type for macOS software updates; possible values are: `alwaysUpdate` (Always update.), `updateDuringTimeWindows` (Update during time windows.), `updateOutsideOfTimeWindows` (Update outside of time windows.). The _provider_ default value is `\"alwaysUpdate\"`.",
					},
					"update_time_window_utc_offset_in_minutes": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes indicating UTC offset for each update time window",
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: "MacOS Software Update Configuration / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macossoftwareupdateconfiguration?view=graph-rest-beta",
			},
		},
		"windows_health_monitoring": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.windowsHealthMonitoringConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // windowsHealthMonitoringConfiguration
					"allow": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						Description:         `allowDeviceHealthMonitoring`, // custom MS Graph attribute name
						MarkdownDescription: "Enables device health monitoring on the device. / Possible values of a property; possible values are: `notConfigured` (Device default value, no intent.), `enabled` (Enables the setting on the device.), `disabled` (Disables the setting on the device.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"custom_scope": schema.StringAttribute{
						Optional:            true,
						Description:         `configDeviceHealthMonitoringCustomScope`, // custom MS Graph attribute name
						MarkdownDescription: "Specifies custom set of events collected from the device where health monitoring is enabled",
					},
					"scope": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							wpvalidator.FlagEnumValues("undefined", "healthMonitoring", "bootPerformance", "windowsUpdates", "privilegeManagement"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("undefined")},
						Computed:            true,
						Description:         `configDeviceHealthMonitoringScope`, // custom MS Graph attribute name
						MarkdownDescription: "Specifies set of events collected from the device where health monitoring is enabled. / Device health monitoring scope; possible values are: `undefined` (Undefined), `healthMonitoring` (Basic events for windows device health monitoring), `bootPerformance` (Boot performance events), `windowsUpdates` (Windows updates events), `privilegeManagement` (PrivilegeManagement). The _provider_ default value is `\"undefined\"`.",
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: "Windows device health monitoring configuration / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowshealthmonitoringconfiguration?view=graph-rest-beta",
			},
		},
		"windows_update_for_business": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.windowsUpdateForBusinessConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // windowsUpdateForBusinessConfiguration
					"allow_windows11_upgrade": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `allowWindows11Upgrade`, // custom MS Graph attribute name
						MarkdownDescription: "When TRUE, allows eligible Windows 10 devices to upgrade to Windows 11. When FALSE, implies the device stays on the existing operating system. Returned by default. Query parameters are not supported. The _provider_ default value is `false`.",
					},
					"automatic_update_mode": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userDefined", "notifyDownload", "autoInstallAtMaintenanceTime", "autoInstallAndRebootAtMaintenanceTime", "autoInstallAndRebootAtScheduledTime", "autoInstallAndRebootWithoutEndUserControl", "windowsDefault"),
						},
						PlanModifiers: []planmodifier.String{
							wpdefaultvalue.StringDefaultValue("autoInstallAtMaintenanceTime"),
						},
						Computed:            true,
						MarkdownDescription: "The Automatic Update Mode. Possible values are: UserDefined, NotifyDownload, AutoInstallAtMaintenanceTime, AutoInstallAndRebootAtMaintenanceTime, AutoInstallAndRebootAtScheduledTime, AutoInstallAndRebootWithoutEndUserControl, WindowsDefault. UserDefined is the default value, no intent. Returned by default. Query parameters are not supported. / Possible values for automatic update mode; possible values are: `userDefined` (User Defined, default value, no intent.), `notifyDownload` (Notify on download.), `autoInstallAtMaintenanceTime` (Auto-install at maintenance time.), `autoInstallAndRebootAtMaintenanceTime` (Auto-install and reboot at maintenance time.), `autoInstallAndRebootAtScheduledTime` (Auto-install and reboot at scheduled time.), `autoInstallAndRebootWithoutEndUserControl` (Auto-install and restart without end-user control), `windowsDefault` (Reset to Windows default value.). The _provider_ default value is `\"autoInstallAtMaintenanceTime\"`.",
					},
					"auto_restart_notification_dismissal": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "automatic", "user", "unknownFutureValue"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Specify the method by which the auto-restart required notification is dismissed. Possible values are: NotConfigured, Automatic, User. Returned by default. Query parameters are not supported. / Auto restart required notification dismissal method; possible values are: `notConfigured` (Not configured), `automatic` (Auto dismissal Indicates that the notification is automatically dismissed without user intervention), `user` (User dismissal. Allows the user to dismiss the notification), `unknownFutureValue` (Evolvable enum member). The _provider_ default value is `\"notConfigured\"`.",
					},
					"business_ready_updates_only": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userDefined", "all", "businessReadyOnly", "windowsInsiderBuildFast", "windowsInsiderBuildSlow", "windowsInsiderBuildRelease"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("userDefined")},
						Computed:            true,
						MarkdownDescription: "Determines which branch devices will receive their updates from. Possible values are: UserDefined, All, BusinessReadyOnly, WindowsInsiderBuildFast, WindowsInsiderBuildSlow, WindowsInsiderBuildRelease. Returned by default. Query parameters are not supported. / Which branch devices will receive their updates from; possible values are: `userDefined` (Allow the user to set.), `all` (Semi-annual Channel (Targeted). Device gets all applicable feature updates from Semi-annual Channel (Targeted).), `businessReadyOnly` (Semi-annual Channel. Device gets feature updates from Semi-annual Channel.), `windowsInsiderBuildFast` (Windows Insider build - Fast), `windowsInsiderBuildSlow` (Windows Insider build - Slow), `windowsInsiderBuildRelease` (Release Windows Insider build). The _provider_ default value is `\"userDefined\"`.",
					},
					"deadline_for_feature_updates_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days before feature updates are installed automatically with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.",
					},
					"deadline_for_quality_updates_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days before quality updates are installed automatically with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.",
					},
					"deadline_grace_period_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days after deadline until restarts occur automatically with valid range from 0 to 7 days. Returned by default. Query parameters are not supported.",
					},
					"delivery_optimization_mode": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userDefined", "httpOnly", "httpWithPeeringNat", "httpWithPeeringPrivateGroup", "httpWithInternetPeering", "simpleDownload", "bypassMode"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("userDefined")},
						Computed:            true,
						MarkdownDescription: "The Delivery Optimization Mode. Possible values are: UserDefined, HttpOnly, HttpWithPeeringNat, HttpWithPeeringPrivateGroup, HttpWithInternetPeering, SimpleDownload, BypassMode. UserDefined allows the user to set. Returned by default. Query parameters are not supported. / Delivery optimization mode for peer distribution; possible values are: `userDefined` (Allow the user to set.), `httpOnly` (HTTP only, no peering), `httpWithPeeringNat` (OS default – Http blended with peering behind the same network address translator), `httpWithPeeringPrivateGroup` (HTTP blended with peering across a private group), `httpWithInternetPeering` (HTTP blended with Internet peering), `simpleDownload` (Simple download mode with no peering), `bypassMode` (Bypass mode. Do not use Delivery Optimization and use BITS instead). The _provider_ default value is `\"userDefined\"`.",
					},
					"drivers_excluded": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "When TRUE, excludes Windows update Drivers. When FALSE, does not exclude Windows update Drivers. Returned by default. Query parameters are not supported. The _provider_ default value is `false`.",
					},
					"engaged_restart_deadline_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Deadline in days before automatically scheduling and executing a pending restart outside of active hours, with valid range from 2 to 30 days. Returned by default. Query parameters are not supported.",
					},
					"engaged_restart_snooze_schedule_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days a user can snooze Engaged Restart reminder notifications with valid range from 1 to 3 days. Returned by default. Query parameters are not supported.",
					},
					"engaged_restart_transition_schedule_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days before transitioning from Auto Restarts scheduled outside of active hours to Engaged Restart, which requires the user to schedule, with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.",
					},
					"feature_updates_deferral_period_in_days": schema.Int64Attribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
						Computed:            true,
						MarkdownDescription: "Defer Feature Updates by these many days with valid range from 0 to 30 days. Returned by default. Query parameters are not supported. The _provider_ default value is `0`.",
					},
					"feature_updates_paused": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "When TRUE, assigned devices are paused from receiving feature updates for up to 35 days from the time you pause the ring. When FALSE, does not pause Feature Updates. Returned by default. Query parameters are not supported.s. The _provider_ default value is `false`.",
					},
					"feature_updates_rollback_window_in_days": schema.Int64Attribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(10)},
						Computed:            true,
						MarkdownDescription: "The number of days after a Feature Update for which a rollback is valid with valid range from 2 to 60 days. Returned by default. Query parameters are not supported. The _provider_ default value is `10`.",
					},
					"installation_schedule": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // windowsUpdateInstallScheduleType
							"active_hours": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.windowsUpdateActiveHoursInstall",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // windowsUpdateActiveHoursInstall
										"end": schema.StringAttribute{
											Optional: true,
											PlanModifiers: []planmodifier.String{
												wpdefaultvalue.StringDefaultValue("17:00:00.0000000"),
											},
											Computed:            true,
											Description:         `activeHoursEnd`, // custom MS Graph attribute name
											MarkdownDescription: "Active Hours End. The _provider_ default value is `\"17:00:00.0000000\"`.",
										},
										"start": schema.StringAttribute{
											Optional: true,
											PlanModifiers: []planmodifier.String{
												wpdefaultvalue.StringDefaultValue("08:00:00.0000000"),
											},
											Computed:            true,
											Description:         `activeHoursStart`, // custom MS Graph attribute name
											MarkdownDescription: "Active Hours Start. The _provider_ default value is `\"08:00:00.0000000\"`.",
										},
									},
									Validators: []validator.Object{
										deviceConfigurationWindowsUpdateInstallScheduleTypeValidator,
									},
									MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowsupdateactivehoursinstall?view=graph-rest-beta",
								},
							},
							"scheduled": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.windowsUpdateScheduledInstall",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // windowsUpdateScheduledInstall
										"day": schema.StringAttribute{
											Optional: true,
											Validators: []validator.String{
												stringvalidator.OneOf("userDefined", "everyday", "sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "noScheduledScan"),
											},
											PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("everyday")},
											Computed:            true,
											Description:         `scheduledInstallDay`, // custom MS Graph attribute name
											MarkdownDescription: "Scheduled Install Day in week. / Possible values for a weekly schedule; possible values are: `userDefined` (User Defined, default value, no intent.), `everyday` (Everyday.), `sunday` (Sunday.), `monday` (Monday.), `tuesday` (Tuesday.), `wednesday` (Wednesday.), `thursday` (Thursday.), `friday` (Friday.), `saturday` (Saturday.), `noScheduledScan` (No Scheduled Scan). The _provider_ default value is `\"everyday\"`.",
										},
										"time": schema.StringAttribute{
											Optional: true,
											PlanModifiers: []planmodifier.String{
												wpdefaultvalue.StringDefaultValue("03:00:00.0000000"),
											},
											Computed:            true,
											Description:         `scheduledInstallTime`, // custom MS Graph attribute name
											MarkdownDescription: "Scheduled Install Time during day. The _provider_ default value is `\"03:00:00.0000000\"`.",
										},
									},
									Validators: []validator.Object{
										deviceConfigurationWindowsUpdateInstallScheduleTypeValidator,
									},
									MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowsupdatescheduledinstall?view=graph-rest-beta",
								},
							},
						},
						MarkdownDescription: "The Installation Schedule. Possible values are: ActiveHoursStart, ActiveHoursEnd, ScheduledInstallDay, ScheduledInstallTime. Returned by default. Query parameters are not supported. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowsupdateinstallscheduletype?view=graph-rest-beta",
					},
					"microsoft_update_service_allowed": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
						Computed:            true,
						MarkdownDescription: "When TRUE, allows Microsoft Update Service. When FALSE, does not allow Microsoft Update Service. Returned by default. Query parameters are not supported. The _provider_ default value is `true`.",
					},
					"postpone_reboot_until_after_deadline": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "When TRUE the device should wait until deadline for rebooting outside of active hours. When FALSE the device should not wait until deadline for rebooting outside of active hours. Returned by default. Query parameters are not supported.",
					},
					"prerelease_features": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userDefined", "settingsOnly", "settingsAndExperimentations", "notAllowed"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("userDefined")},
						Computed:            true,
						MarkdownDescription: "The Pre-Release Features. Possible values are: UserDefined, SettingsOnly, SettingsAndExperimentations, NotAllowed. UserDefined is the default value, no intent. Returned by default. Query parameters are not supported. / Possible values for pre-release features; possible values are: `userDefined` (User Defined, default value, no intent.), `settingsOnly` (Settings only pre-release features.), `settingsAndExperimentations` (Settings and experimentations pre-release features.), `notAllowed` (Pre-release features not allowed.). The _provider_ default value is `\"userDefined\"`.",
					},
					"quality_updates_deferral_period_in_days": schema.Int64Attribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
						Computed:            true,
						MarkdownDescription: "Defer Quality Updates by these many days with valid range from 0 to 30 days. Returned by default. Query parameters are not supported. The _provider_ default value is `0`.",
					},
					"quality_updates_paused": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "When TRUE, assigned devices are paused from receiving quality updates for up to 35 days from the time you pause the ring. When FALSE, does not pause Quality Updates. Returned by default. Query parameters are not supported. The _provider_ default value is `false`.",
					},
					"schedule_imminent_restart_warning_in_minutes": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Specify the period for auto-restart imminent warning notifications. Supported values: 15, 30 or 60 (minutes). Returned by default. Query parameters are not supported.",
					},
					"schedule_restart_warning_in_hours": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Specify the period for auto-restart warning reminder notifications. Supported values: 2, 4, 8, 12 or 24 (hours). Returned by default. Query parameters are not supported.",
					},
					"skip_checks_before_restart": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "When TRUE, skips all checks before restart: Battery level = 40%, User presence, Display Needed, Presentation mode, Full screen mode, phone call state, game mode etc. When FALSE, does not skip all checks before restart. Returned by default. Query parameters are not supported. The _provider_ default value is `false`.",
					},
					"update_notification_level": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "defaultNotifications", "restartWarningsOnly", "disableAllNotifications", "unknownFutureValue"),
						},
						PlanModifiers: []planmodifier.String{
							wpdefaultvalue.StringDefaultValue("defaultNotifications"),
						},
						Computed:            true,
						MarkdownDescription: "Specifies what Windows Update notifications users see. Possible values are: NotConfigured, DefaultNotifications, RestartWarningsOnly, DisableAllNotifications. Returned by default. Query parameters are not supported. / Windows Update Notification Display Options; possible values are: `notConfigured` (Not configured), `defaultNotifications` (Use the default Windows Update notifications.), `restartWarningsOnly` (Turn off all notifications, excluding restart warnings.), `disableAllNotifications` (Turn off all notifications, including restart warnings.), `unknownFutureValue` (Evolvable enum member). The _provider_ default value is `\"defaultNotifications\"`.",
					},
					"update_weeks": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							wpvalidator.FlagEnumValues("userDefined", "firstWeek", "secondWeek", "thirdWeek", "fourthWeek", "everyWeek", "unknownFutureValue"),
						},
						MarkdownDescription: "Schedule the update installation on the weeks of the month. Possible values are: UserDefined, FirstWeek, SecondWeek, ThirdWeek, FourthWeek, EveryWeek. Returned by default. Query parameters are not supported. / Scheduled the update installation on the weeks of the month; possible values are: `userDefined` (Allow the user to set.), `firstWeek` (Scheduled the update installation on the first week of the month), `secondWeek` (Scheduled the update installation on the second week of the month), `thirdWeek` (Scheduled the update installation on the third week of the month), `fourthWeek` (Scheduled the update installation on the fourth week of the month), `everyWeek` (Scheduled the update installation on every week of the month), `unknownFutureValue` (Evolvable enum member)",
					},
					"user_pause_access": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("enabled")},
						Computed:            true,
						MarkdownDescription: "Specifies whether to enable end user’s access to pause software updates. Possible values are: NotConfigured, Enabled, Disabled. Returned by default. Query parameters are not supported. / Possible values of a property; possible values are: `notConfigured` (Device default value, no intent.), `enabled` (Enables the setting on the device.), `disabled` (Disables the setting on the device.). The _provider_ default value is `\"enabled\"`.",
					},
					"user_windows_update_scan_access": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("enabled")},
						Computed:            true,
						MarkdownDescription: "Specifies whether to disable user’s access to scan Windows Update. Possible values are: NotConfigured, Enabled, Disabled. Returned by default. Query parameters are not supported. / Possible values of a property; possible values are: `notConfigured` (Device default value, no intent.), `enabled` (Enables the setting on the device.), `disabled` (Disables the setting on the device.). The _provider_ default value is `\"enabled\"`.",
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: "Windows Update for business configuration, allows you to specify how and when Windows as a Service updates your Windows 10/11 devices with feature and quality updates. Supports ODATA clauses that DeviceConfiguration entity supports: $filter by types of DeviceConfiguration, $top, $select only DeviceConfiguration base properties, $orderby only DeviceConfiguration base properties, and $skip. The query parameter '$search' is not supported. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowsupdateforbusinessconfiguration?view=graph-rest-beta",
			},
		},
		"windows10_general": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.windows10GeneralConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // windows10GeneralConfiguration
					"accounts_block_adding_non_microsoft_account_email": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Block the user from adding email accounts to the device that are not associated with a Microsoft account. The _provider_ default value is `false`.",
					},
					"activate_apps_with_voice": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Specifies if Windows apps can be activated by voice. / Possible values of a property; possible values are: `notConfigured` (Device default value, no intent.), `enabled` (Enables the setting on the device.), `disabled` (Disables the setting on the device.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"anti_theft_mode_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from selecting an AntiTheft mode preference (Windows 10 Mobile only). The _provider_ default value is `false`.",
					},
					"app_management_msi_allow_user_control_over_install": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `appManagementMSIAllowUserControlOverInstall`, // custom MS Graph attribute name
						MarkdownDescription: "This policy setting permits users to change installation options that typically are available only to system administrators. The _provider_ default value is `false`.",
					},
					"app_management_msi_always_install_with_elevated_privileges": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `appManagementMSIAlwaysInstallWithElevatedPrivileges`, // custom MS Graph attribute name
						MarkdownDescription: "This policy setting directs Windows Installer to use elevated permissions when it installs any program on the system. The _provider_ default value is `false`.",
					},
					"app_management_package_family_names_to_launch_after_log_on": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "List of semi-colon delimited Package Family Names of Windows apps. Listed Windows apps are to be launched after logon.​ The _provider_ default value is `[]`.",
					},
					"apps_allow_trusted_apps_sideloading": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "blocked", "allowed"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Indicates whether apps from AppX packages signed with a trusted certificate can be side loaded. / State Management Setting; possible values are: `notConfigured` (Not configured.), `blocked` (Blocked.), `allowed` (Allowed.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"apps_block_windows_store_originated_apps": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to disable the launch of all apps from Windows Store that came pre-installed or were downloaded. The _provider_ default value is `false`.",
					},
					"authentication_allow_secondary_device": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Allows secondary authentication devices to work with Windows. The _provider_ default value is `false`.",
					},
					"authentication_preferred_azure_ad_tenant_domain_name": schema.StringAttribute{
						Optional:            true,
						Description:         `authenticationPreferredAzureADTenantDomainName`, // custom MS Graph attribute name
						MarkdownDescription: "Specifies the preferred domain among available domains in the Azure AD tenant.",
					},
					"authentication_web_sign_in": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not Web Credential Provider will be enabled. / Possible values of a property; possible values are: `notConfigured` (Device default value, no intent.), `enabled` (Enables the setting on the device.), `disabled` (Disables the setting on the device.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"bluetooth_allowed_services": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Specify a list of allowed Bluetooth services and profiles in hex formatted strings. The _provider_ default value is `[]`.",
					},
					"bluetooth_block_advertising": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether or not to Block the user from using bluetooth advertising. The _provider_ default value is `false`.",
					},
					"bluetooth_block_discoverable_mode": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether or not to Block the user from using bluetooth discoverable mode. The _provider_ default value is `false`.",
					},
					"bluetooth_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether or not to Block the user from using bluetooth. The _provider_ default value is `false`.",
					},
					"bluetooth_block_pre_pairing": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether or not to block specific bundled Bluetooth peripherals to automatically pair with the host device. The _provider_ default value is `false`.",
					},
					"bluetooth_block_prompted_proximal_connections": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether or not to block the users from using Swift Pair and other proximity based scenarios. The _provider_ default value is `false`.",
					},
					"camera_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether or not to Block the user from accessing the camera of the device. The _provider_ default value is `false`.",
					},
					"cellular_block_data_when_roaming": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether or not to Block the user from using data over cellular while roaming. The _provider_ default value is `false`.",
					},
					"cellular_block_vpn": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether or not to Block the user from using VPN over cellular. The _provider_ default value is `false`.",
					},
					"cellular_block_vpn_when_roaming": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether or not to Block the user from using VPN when roaming over cellular. The _provider_ default value is `false`.",
					},
					"cellular_data": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("blocked", "required", "allowed", "notConfigured"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("allowed")},
						Computed:            true,
						MarkdownDescription: "Whether or not to allow the cellular data channel on the device. If not configured, the cellular data channel is allowed and the user can turn it off. / Possible values of the ConfigurationUsage list; possible values are: `blocked` (Disallowed.), `required` (Required.), `allowed` (Optional.), `notConfigured` (Not Configured.). The _provider_ default value is `\"allowed\"`.",
					},
					"certificates_block_manual_root_certificate_installation": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether or not to Block the user from doing manual root certificate installation. The _provider_ default value is `false`.",
					},
					"configure_time_zone": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Specifies the time zone to be applied to the device. This is the standard Windows name for the target time zone.",
					},
					"connected_devices_service_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether or not to block Connected Devices Service which enables discovery and connection to other devices, remote messaging, remote app sessions and other cross-device experiences. The _provider_ default value is `false`.",
					},
					"copy_paste_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether or not to Block the user from using copy paste. The _provider_ default value is `false`.",
					},
					"cortana_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether or not to Block the user from using Cortana. The _provider_ default value is `false`.",
					},
					"cryptography_allow_fips_algorithm_policy": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Specify whether to allow or disallow the Federal Information Processing Standard (FIPS) policy. The _provider_ default value is `false`.",
					},
					"data_protection_block_direct_memory_access": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "This policy setting allows you to block direct memory access (DMA) for all hot pluggable PCI downstream ports until a user logs into Windows.",
					},
					"defender_block_end_user_access": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether or not to block end user access to Defender. The _provider_ default value is `false`.",
					},
					"defender_block_on_access_protection": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Allows or disallows Windows Defender On Access Protection functionality. The _provider_ default value is `false`.",
					},
					"defender_cloud_block_level": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "high", "highPlus", "zeroTolerance"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Specifies the level of cloud-delivered protection. / Possible values of Cloud Block Level; possible values are: `notConfigured` (Default value, uses the default Windows Defender Antivirus blocking level and provides strong detection without increasing the risk of detecting legitimate files), `high` (High applies a strong level of detection.), `highPlus` (High + uses the High level and applies addition protection measures), `zeroTolerance` (Zero tolerance blocks all unknown executables). The _provider_ default value is `\"notConfigured\"`.",
					},
					"defender_cloud_extended_timeout": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Timeout extension for file scanning by the cloud. Valid values 0 to 50",
					},
					"defender_cloud_extended_timeout_in_seconds": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Timeout extension for file scanning by the cloud. Valid values 0 to 50",
					},
					"defender_days_before_deleting_quarantined_malware": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days before deleting quarantined malware. Valid values 0 to 90",
					},
					"defender_detected_malware_actions": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // defenderDetectedMalwareActions
							"high_severity": schema.StringAttribute{
								Optional: true,
								Validators: []validator.String{
									stringvalidator.OneOf("deviceDefault", "clean", "quarantine", "remove", "allow", "userDefined", "block"),
								},
								PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
								Computed:            true,
								MarkdownDescription: "Indicates a Defender action to take for high severity Malware threat detected. / Defender’s default action to take on detected Malware threats; possible values are: `deviceDefault` (Apply action based on the update definition.), `clean` (Clean the detected threat.), `quarantine` (Quarantine the detected threat.), `remove` (Remove the detected threat.), `allow` (Allow the detected threat.), `userDefined` (Allow the user to determine the action to take with the detected threat.), `block` (Block the detected threat.). The _provider_ default value is `\"deviceDefault\"`.",
							},
							"low_severity": schema.StringAttribute{
								Optional: true,
								Validators: []validator.String{
									stringvalidator.OneOf("deviceDefault", "clean", "quarantine", "remove", "allow", "userDefined", "block"),
								},
								PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
								Computed:            true,
								MarkdownDescription: "Indicates a Defender action to take for low severity Malware threat detected. / Defender’s default action to take on detected Malware threats; possible values are: `deviceDefault` (Apply action based on the update definition.), `clean` (Clean the detected threat.), `quarantine` (Quarantine the detected threat.), `remove` (Remove the detected threat.), `allow` (Allow the detected threat.), `userDefined` (Allow the user to determine the action to take with the detected threat.), `block` (Block the detected threat.). The _provider_ default value is `\"deviceDefault\"`.",
							},
							"moderate_severity": schema.StringAttribute{
								Optional: true,
								Validators: []validator.String{
									stringvalidator.OneOf("deviceDefault", "clean", "quarantine", "remove", "allow", "userDefined", "block"),
								},
								PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
								Computed:            true,
								MarkdownDescription: "Indicates a Defender action to take for moderate severity Malware threat detected. / Defender’s default action to take on detected Malware threats; possible values are: `deviceDefault` (Apply action based on the update definition.), `clean` (Clean the detected threat.), `quarantine` (Quarantine the detected threat.), `remove` (Remove the detected threat.), `allow` (Allow the detected threat.), `userDefined` (Allow the user to determine the action to take with the detected threat.), `block` (Block the detected threat.). The _provider_ default value is `\"deviceDefault\"`.",
							},
							"severe_severity": schema.StringAttribute{
								Optional: true,
								Validators: []validator.String{
									stringvalidator.OneOf("deviceDefault", "clean", "quarantine", "remove", "allow", "userDefined", "block"),
								},
								PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
								Computed:            true,
								MarkdownDescription: "Indicates a Defender action to take for severe severity Malware threat detected. / Defender’s default action to take on detected Malware threats; possible values are: `deviceDefault` (Apply action based on the update definition.), `clean` (Clean the detected threat.), `quarantine` (Quarantine the detected threat.), `remove` (Remove the detected threat.), `allow` (Allow the detected threat.), `userDefined` (Allow the user to determine the action to take with the detected threat.), `block` (Block the detected threat.). The _provider_ default value is `\"deviceDefault\"`.",
							},
						},
						MarkdownDescription: "Gets or sets Defender’s actions to take on detected Malware per threat level. / Specify Defender’s actions to take on detected Malware per threat level. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-defenderdetectedmalwareactions?view=graph-rest-beta",
					},
					"defender_disable_catchup_full_scan": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "When blocked, catch-up scans for scheduled full scans will be turned off. The _provider_ default value is `false`.",
					},
					"defender_disable_catchup_quick_scan": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "When blocked, catch-up scans for scheduled quick scans will be turned off. The _provider_ default value is `false`.",
					},
					"defender_file_extensions_to_exclude": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "File extensions to exclude from scans and real time protection. The _provider_ default value is `[]`.",
					},
					"defender_files_and_folders_to_exclude": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Files and folder to exclude from scans and real time protection. The _provider_ default value is `[]`.",
					},
					"defender_monitor_file_activity": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userDefined", "disable", "monitorAllFiles", "monitorIncomingFilesOnly", "monitorOutgoingFilesOnly"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("userDefined")},
						Computed:            true,
						MarkdownDescription: "Value for monitoring file activity. / Possible values for monitoring file activity; possible values are: `userDefined` (User Defined, default value, no intent.), `disable` (Disable monitoring file activity.), `monitorAllFiles` (Monitor all files.), `monitorIncomingFilesOnly` (Monitor incoming files only.), `monitorOutgoingFilesOnly` (Monitor outgoing files only.). The _provider_ default value is `\"userDefined\"`.",
					},
					"defender_potentially_unwanted_app_action": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "block", "audit"),
						},
						MarkdownDescription: "Gets or sets Defender’s action to take on Potentially Unwanted Application (PUA), which includes software with behaviors of ad-injection, software bundling, persistent solicitation for payment or subscription, etc. Defender alerts user when PUA is being downloaded or attempts to install itself. Added in Windows 10 for desktop. / Defender’s action to take on detected Potentially Unwanted Application (PUA); possible values are: `deviceDefault` (PUA Protection is off. Defender will not protect against potentially unwanted applications.), `block` (PUA Protection is on. Detected items are blocked. They will show in history along with other threats.), `audit` (Audit mode. Defender will detect potentially unwanted applications, but take no actions. You can review information about applications Defender would have taken action against by searching for events created by Defender in the Event Viewer.)",
					},
					"defender_potentially_unwanted_app_action_setting": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userDefined", "enable", "auditMode", "warn", "notConfigured"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("userDefined")},
						Computed:            true,
						MarkdownDescription: "Gets or sets Defender’s action to take on Potentially Unwanted Application (PUA), which includes software with behaviors of ad-injection, software bundling, persistent solicitation for payment or subscription, etc. Defender alerts user when PUA is being downloaded or attempts to install itself. Added in Windows 10 for desktop. / Possible values of Defender PUA Protection; possible values are: `userDefined` (Device default value, no intent.), `enable` (Block functionality.), `auditMode` (Allow functionality but generate logs.), `warn` (Warning message to end user with ability to bypass block from attack surface reduction rule.), `notConfigured` (Not configured.). The _provider_ default value is `\"userDefined\"`.",
					},
					"defender_processes_to_exclude": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Processes to exclude from scans and real time protection. The _provider_ default value is `[]`.",
					},
					"defender_prompt_for_sample_submission": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userDefined", "alwaysPrompt", "promptBeforeSendingPersonalData", "neverSendData", "sendAllDataWithoutPrompting"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("userDefined")},
						Computed:            true,
						MarkdownDescription: "The configuration for how to prompt user for sample submission. / Possible values for prompting user for samples submission; possible values are: `userDefined` (User Defined, default value, no intent.), `alwaysPrompt` (Always prompt.), `promptBeforeSendingPersonalData` (Send safe samples automatically.), `neverSendData` (Never send data.), `sendAllDataWithoutPrompting` (Send all data without prompting.). The _provider_ default value is `\"userDefined\"`.",
					},
					"defender_require_behavior_monitoring": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to require behavior monitoring. The _provider_ default value is `false`.",
					},
					"defender_require_cloud_protection": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to require cloud protection. The _provider_ default value is `false`.",
					},
					"defender_require_network_inspection_system": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to require network inspection system. The _provider_ default value is `false`.",
					},
					"defender_require_real_time_monitoring": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to require real time monitoring. The _provider_ default value is `false`.",
					},
					"defender_scan_archive_files": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to scan archive files. The _provider_ default value is `false`.",
					},
					"defender_scan_downloads": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to scan downloads. The _provider_ default value is `false`.",
					},
					"defender_scan_incoming_mail": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to scan incoming mail messages. The _provider_ default value is `false`.",
					},
					"defender_scan_mapped_network_drives_during_full_scan": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to scan mapped network drives during full scan. The _provider_ default value is `false`.",
					},
					"defender_scan_max_cpu": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Max CPU usage percentage during scan. Valid values 0 to 100",
					},
					"defender_scan_network_files": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to scan files opened from a network folder. The _provider_ default value is `false`.",
					},
					"defender_scan_removable_drives_during_full_scan": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to scan removable drives during full scan. The _provider_ default value is `false`.",
					},
					"defender_scan_scripts_loaded_in_internet_explorer": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to scan scripts loaded in Internet Explorer browser. The _provider_ default value is `false`.",
					},
					"defender_scan_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userDefined", "disabled", "quick", "full"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("userDefined")},
						Computed:            true,
						MarkdownDescription: "The defender system scan type. / Possible values for system scan type; possible values are: `userDefined` (User Defined, default value, no intent.), `disabled` (System scan disabled.), `quick` (Quick system scan.), `full` (Full system scan.). The _provider_ default value is `\"userDefined\"`.",
					},
					"defender_scheduled_quick_scan_time": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The time to perform a daily quick scan.",
					},
					"defender_scheduled_scan_time": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The defender time for the system scan.",
					},
					"defender_schedule_scan_enable_low_cpu_priority": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "When enabled, low CPU priority will be used during scheduled scans. The _provider_ default value is `false`.",
					},
					"defender_signature_update_interval_in_hours": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The signature update interval in hours. Specify 0 not to check. Valid values 0 to 24",
					},
					"defender_submit_samples_consent_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("sendSafeSamplesAutomatically", "alwaysPrompt", "neverSend", "sendAllSamplesAutomatically"),
						},
						MarkdownDescription: "Checks for the user consent level in Windows Defender to send data. / Possible values for DefenderSubmitSamplesConsentType; possible values are: `sendSafeSamplesAutomatically` (Send safe samples automatically), `alwaysPrompt` (Always prompt), `neverSend` (Never send), `sendAllSamplesAutomatically` (Send all samples automatically)",
					},
					"defender_system_scan_schedule": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userDefined", "everyday", "sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "noScheduledScan"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("userDefined")},
						Computed:            true,
						MarkdownDescription: "Defender day of the week for the system scan. / Possible values for a weekly schedule; possible values are: `userDefined` (User Defined, default value, no intent.), `everyday` (Everyday.), `sunday` (Sunday.), `monday` (Monday.), `tuesday` (Tuesday.), `wednesday` (Wednesday.), `thursday` (Thursday.), `friday` (Friday.), `saturday` (Saturday.), `noScheduledScan` (No Scheduled Scan). The _provider_ default value is `\"userDefined\"`.",
					},
					"developer_unlock_setting": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "blocked", "allowed"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow developer unlock. / State Management Setting; possible values are: `notConfigured` (Not configured.), `blocked` (Blocked.), `allowed` (Allowed.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"device_management_block_factory_reset_on_mobile": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Block the user from resetting their phone. The _provider_ default value is `false`.",
					},
					"device_management_block_manual_unenroll": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Block the user from doing manual un-enrollment from device management. The _provider_ default value is `false`.",
					},
					"diagnostics_data_submission_mode": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userDefined", "none", "basic", "enhanced", "full"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("userDefined")},
						Computed:            true,
						MarkdownDescription: "Gets or sets a value allowing the device to send diagnostic and usage telemetry data, such as Watson. / Allow the device to send diagnostic and usage telemetry data, such as Watson; possible values are: `userDefined` (Allow the user to set.), `none` (No telemetry data is sent from OS components. Note: This value is only applicable to enterprise and server devices. Using this setting on other devices is equivalent to setting the value of 1.), `basic` (Sends basic telemetry data.), `enhanced` (Sends enhanced telemetry data including usage and insights data.), `full` (Sends full telemetry data including diagnostic data, such as system state.). The _provider_ default value is `\"userDefined\"`.",
					},
					"display_app_list_with_gdi_dpi_scaling_turned_off": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						Description:         `displayAppListWithGdiDPIScalingTurnedOff`, // custom MS Graph attribute name
						MarkdownDescription: "List of legacy applications that have GDI DPI Scaling turned off. The _provider_ default value is `[]`.",
					},
					"display_app_list_with_gdi_dpi_scaling_turned_on": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						Description:         `displayAppListWithGdiDPIScalingTurnedOn`, // custom MS Graph attribute name
						MarkdownDescription: "List of legacy applications that have GDI DPI Scaling turned on. The _provider_ default value is `[]`.",
					},
					"edge_allow_start_pages_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Allow users to change Start pages on Edge. Use the EdgeHomepageUrls to specify the Start pages that the user would see by default when they open Edge. The _provider_ default value is `false`.",
					},
					"edge_block_access_to_about_flags": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to prevent access to about flags on Edge browser. The _provider_ default value is `false`.",
					},
					"edge_block_address_bar_dropdown": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Block the address bar dropdown functionality in Microsoft Edge. Disable this settings to minimize network connections from Microsoft Edge to Microsoft services. The _provider_ default value is `false`.",
					},
					"edge_block_autofill": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block auto fill. The _provider_ default value is `false`.",
					},
					"edge_block_compatibility_list": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Block Microsoft compatibility list in Microsoft Edge. This list from Microsoft helps Edge properly display sites with known compatibility issues. The _provider_ default value is `false`.",
					},
					"edge_block_developer_tools": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block developer tools in the Edge browser. The _provider_ default value is `false`.",
					},
					"edge_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Block the user from using the Edge browser. The _provider_ default value is `false`.",
					},
					"edge_block_edit_favorites": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Block the user from making changes to Favorites. The _provider_ default value is `false`.",
					},
					"edge_block_extensions": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block extensions in the Edge browser. The _provider_ default value is `false`.",
					},
					"edge_block_full_screen_mode": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Allow or prevent Edge from entering the full screen mode. The _provider_ default value is `false`.",
					},
					"edge_block_in_private_browsing": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block InPrivate browsing on corporate networks, in the Edge browser. The _provider_ default value is `false`.",
					},
					"edge_block_java_script": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Block the user from using JavaScript. The _provider_ default value is `false`.",
					},
					"edge_block_live_tile_data_collection": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Block the collection of information by Microsoft for live tile creation when users pin a site to Start from Microsoft Edge. The _provider_ default value is `false`.",
					},
					"edge_block_password_manager": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Block password manager. The _provider_ default value is `false`.",
					},
					"edge_block_popups": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block popups. The _provider_ default value is `false`.",
					},
					"edge_block_prelaunch": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Decide whether Microsoft Edge is prelaunched at Windows startup. The _provider_ default value is `false`.",
					},
					"edge_block_printing": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Configure Edge to allow or block printing. The _provider_ default value is `false`.",
					},
					"edge_block_saving_history": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Configure Edge to allow browsing history to be saved or to never save browsing history. The _provider_ default value is `false`.",
					},
					"edge_block_search_engine_customization": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from adding new search engine or changing the default search engine. The _provider_ default value is `false`.",
					},
					"edge_block_search_suggestions": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from using the search suggestions in the address bar. The _provider_ default value is `false`.",
					},
					"edge_block_sending_do_not_track_header": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Block the user from sending the do not track header. The _provider_ default value is `false`.",
					},
					"edge_block_sending_intranet_traffic_to_internet_explorer": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to switch the intranet traffic from Edge to Internet Explorer. Note: the name of this property is misleading; the property is obsolete, use EdgeSendIntranetTrafficToInternetExplorer instead. The _provider_ default value is `false`.",
					},
					"edge_block_sideloading_extensions": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether the user can sideload extensions. The _provider_ default value is `false`.",
					},
					"edge_block_tab_preloading": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Configure whether Edge preloads the new tab page at Windows startup. The _provider_ default value is `false`.",
					},
					"edge_block_web_content_on_new_tab_page": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Configure to load a blank page in Edge instead of the default New tab page and prevent users from changing it. The _provider_ default value is `false`.",
					},
					"edge_clear_browsing_data_on_exit": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Clear browsing data on exiting Microsoft Edge. The _provider_ default value is `false`.",
					},
					"edge_cookie_policy": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userDefined", "allow", "blockThirdParty", "blockAll"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("userDefined")},
						Computed:            true,
						MarkdownDescription: "Indicates which cookies to block in the Edge browser. / Possible values to specify which cookies are allowed in Microsoft Edge; possible values are: `userDefined` (Allow the user to set.), `allow` (Allow.), `blockThirdParty` (Block only third party cookies.), `blockAll` (Block all cookies.). The _provider_ default value is `\"userDefined\"`.",
					},
					"edge_disable_first_run_page": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Block the Microsoft web page that opens on the first use of Microsoft Edge. This policy allows enterprises, like those enrolled in zero emissions configurations, to block this page. The _provider_ default value is `false`.",
					},
					"edge_enterprise_mode_site_list_location": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Indicates the enterprise mode site list location. Could be a local file, local network or http location.",
					},
					"edge_favorites_bar_visibility": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "hide", "show"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Get or set a value that specifies whether to set the favorites bar to always be visible or hidden on any page. / Generic visibility state; possible values are: `notConfigured` (Not configured.), `hide` (Hide.), `show` (Show.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"edge_favorites_list_location": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The location of the favorites list to provision. Could be a local file, local network or http location.",
					},
					"edge_first_run_url": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The first run URL for when Edge browser is opened for the first time.",
					},
					"edge_home_button_configuration": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // edgeHomeButtonConfiguration
							"hidden": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.edgeHomeButtonHidden",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional:   true,
									Attributes: map[string]schema.Attribute{ // edgeHomeButtonHidden
									},
									Validators: []validator.Object{
										deviceConfigurationEdgeHomeButtonConfigurationValidator,
									},
									MarkdownDescription: "Hide the home button. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-edgehomebuttonhidden?view=graph-rest-beta",
								},
							},
							"loads_start_page": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.edgeHomeButtonLoadsStartPage",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional:   true,
									Attributes: map[string]schema.Attribute{ // edgeHomeButtonLoadsStartPage
									},
									Validators: []validator.Object{
										deviceConfigurationEdgeHomeButtonConfigurationValidator,
									},
									MarkdownDescription: "Show the home button; clicking the home button loads the Start page - this is also the default value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-edgehomebuttonloadsstartpage?view=graph-rest-beta",
								},
							},
							"opens_custom_url": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.edgeHomeButtonOpensCustomURL",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // edgeHomeButtonOpensCustomURL
										"home_button_custom_url": schema.StringAttribute{
											Optional:            true,
											Description:         `homeButtonCustomURL`, // custom MS Graph attribute name
											MarkdownDescription: "The specific URL to load.",
										},
									},
									Validators: []validator.Object{
										deviceConfigurationEdgeHomeButtonConfigurationValidator,
									},
									MarkdownDescription: "Show the home button; clicking the home button loads a specific URL. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-edgehomebuttonopenscustomurl?view=graph-rest-beta",
								},
							},
							"opens_new_tab": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.edgeHomeButtonOpensNewTab",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional:   true,
									Attributes: map[string]schema.Attribute{ // edgeHomeButtonOpensNewTab
									},
									Validators: []validator.Object{
										deviceConfigurationEdgeHomeButtonConfigurationValidator,
									},
									MarkdownDescription: "Show the home button; clicking the home button loads the New tab page. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-edgehomebuttonopensnewtab?view=graph-rest-beta",
								},
							},
						},
						MarkdownDescription: "Causes the Home button to either hide, load the default Start page, load a New tab page, or a custom URL / The home button configuration base class used to identify the available options / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-edgehomebuttonconfiguration?view=graph-rest-beta",
					},
					"edge_home_button_configuration_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Enable the Home button configuration. The _provider_ default value is `false`.",
					},
					"edge_homepage_urls": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "The list of URLs for homepages shodwn on MDM-enrolled devices on Edge browser. The _provider_ default value is `[]`.",
					},
					"edge_kiosk_mode_restriction": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "digitalSignage", "normalMode", "publicBrowsingSingleApp", "publicBrowsingMultiApp"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Controls how the Microsoft Edge settings are restricted based on the configure kiosk mode. / Specify how the Microsoft Edge settings are restricted based on kiosk mode; possible values are: `notConfigured` (Not configured (unrestricted).), `digitalSignage` (Interactive/Digital signage in single-app mode.), `normalMode` (Normal mode (full version of Microsoft Edge).), `publicBrowsingSingleApp` (Public browsing in single-app mode.), `publicBrowsingMultiApp` (Public browsing (inPrivate) in multi-app mode.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"edge_kiosk_reset_after_idle_time_in_minutes": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Specifies the time in minutes from the last user activity before Microsoft Edge kiosk resets.  Valid values are 0-1440. The default is 5. 0 indicates no reset. Valid values 0 to 1440",
					},
					"edge_new_tab_page_url": schema.StringAttribute{
						Optional:            true,
						Description:         `edgeNewTabPageURL`, // custom MS Graph attribute name
						MarkdownDescription: "Specify the page opened when new tabs are created.",
					},
					"edge_opens_with": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "startPage", "newTabPage", "previousPages", "specificPages"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Specify what kind of pages are open at start. / Possible values for the EdgeOpensWith setting; possible values are: `notConfigured` (Not configured.), `startPage` (StartPage.), `newTabPage` (NewTabPage.), `previousPages` (PreviousPages.), `specificPages` (SpecificPages.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"edge_prevent_certificate_error_override": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Allow or prevent users from overriding certificate errors. The _provider_ default value is `false`.",
					},
					"edge_required_extension_package_family_names": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Specify the list of package family names of browser extensions that are required and cannot be turned off by the user. The _provider_ default value is `[]`.",
					},
					"edge_require_smart_screen": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Require the user to use the smart screen filter. The _provider_ default value is `false`.",
					},
					"edge_search_engine": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // edgeSearchEngineBase
							"custom": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.edgeSearchEngineCustom",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // edgeSearchEngineCustom
										"edge_search_engine_open_search_xml_url": schema.StringAttribute{
											Optional:            true,
											PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
											Computed:            true,
											MarkdownDescription: "Points to a https link containing the OpenSearch xml file that contains, at minimum, the short name and the URL to the search Engine. The _provider_ default value is `\"\"`.",
										},
									},
									Validators:          []validator.Object{deviceConfigurationEdgeSearchEngineBaseValidator},
									MarkdownDescription: "Allows IT admins to set a custom default search engine for MDM-Controlled devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-edgesearchenginecustom?view=graph-rest-beta",
								},
							},
							"predefined": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.edgeSearchEngine",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // edgeSearchEngine
										"edge_search_engine_type": schema.StringAttribute{
											Optional:            true,
											Validators:          []validator.String{stringvalidator.OneOf("default", "bing")},
											PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("default")},
											Computed:            true,
											MarkdownDescription: "Allows IT admins to set a predefined default search engine for MDM-Controlled devices. / Allows IT admind to set a predefined default search engine for MDM-Controlled devices; possible values are: `default` (Uses factory settings of Edge to assign the default search engine as per the user market), `bing` (Sets Bing as the default search engine). The _provider_ default value is `\"default\"`.",
										},
									},
									Validators:          []validator.Object{deviceConfigurationEdgeSearchEngineBaseValidator},
									MarkdownDescription: "Allows IT admins to set a predefined default search engine for MDM-Controlled devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-edgesearchengine?view=graph-rest-beta",
								},
							},
						},
						MarkdownDescription: "Allows IT admins to set a default search engine for MDM-Controlled devices. Users can override this and change their default search engine provided the AllowSearchEngineCustomization policy is not set. / Allows IT admins to set a default search engine for MDM-Controlled devices. Users can override this and change their default search engine provided the AllowSearchEngineCustomization policy is not set. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-edgesearchenginebase?view=graph-rest-beta",
					},
					"edge_send_intranet_traffic_to_internet_explorer": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to switch the intranet traffic from Edge to Internet Explorer. The _provider_ default value is `false`.",
					},
					"edge_show_message_when_opening_internet_explorer_sites": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "disabled", "enabled", "keepGoing"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Controls the message displayed by Edge before switching to Internet Explorer. / What message will be displayed by Edge before switching to Internet Explorer; possible values are: `notConfigured` (Not configured.), `disabled` (Disabled.), `enabled` (Enabled.), `keepGoing` (KeepGoing.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"edge_sync_favorites_with_internet_explorer": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Enable favorites sync between Internet Explorer and Microsoft Edge. Additions, deletions, modifications and order changes to favorites are shared between browsers. The _provider_ default value is `false`.",
					},
					"edge_telemetry_for_microsoft365_analytics": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "intranet", "internet", "intranetAndInternet"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						Description:         `edgeTelemetryForMicrosoft365Analytics`, // custom MS Graph attribute name
						MarkdownDescription: "Specifies what type of telemetry data (none, intranet, internet, both) is sent to Microsoft 365 Analytics. / Type of browsing data sent to Microsoft 365 analytics; possible values are: `notConfigured` (Default – No telemetry data collected or sent), `intranet` (Allow sending intranet history only: Only send browsing history data for intranet sites), `internet` (Allow sending internet history only: Only send browsing history data for internet sites), `intranetAndInternet` (Allow sending both intranet and internet history: Send browsing history data for intranet and internet sites). The _provider_ default value is `\"notConfigured\"`.",
					},
					"enable_automatic_redeployment": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Allow users with administrative rights to delete all user data and settings using CTRL + Win + R at the device lock screen so that the device can be automatically re-configured and re-enrolled into management. The _provider_ default value is `false`.",
					},
					"energy_saver_on_battery_threshold_percentage": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "This setting allows you to specify battery charge level at which Energy Saver is turned on. While on battery, Energy Saver is automatically turned on at (and below) the specified battery charge level. Valid input range (0-100). Valid values 0 to 100",
					},
					"energy_saver_plugged_in_threshold_percentage": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "This setting allows you to specify battery charge level at which Energy Saver is turned on. While plugged in, Energy Saver is automatically turned on at (and below) the specified battery charge level. Valid input range (0-100). Valid values 0 to 100",
					},
					"enterprise_cloud_print_discovery_end_point": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Endpoint for discovering cloud printers.",
					},
					"enterprise_cloud_print_discovery_max_limit": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Maximum number of printers that should be queried from a discovery endpoint. This is a mobile only setting. Valid values 1 to 65535",
					},
					"enterprise_cloud_print_mopria_discovery_resource_identifier": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "OAuth resource URI for printer discovery service as configured in Azure portal.",
					},
					"enterprise_cloud_print_oauth_authority": schema.StringAttribute{
						Optional:            true,
						Description:         `enterpriseCloudPrintOAuthAuthority`, // custom MS Graph attribute name
						MarkdownDescription: "Authentication endpoint for acquiring OAuth tokens.",
					},
					"enterprise_cloud_print_oauth_client_identifier": schema.StringAttribute{
						Optional:            true,
						Description:         `enterpriseCloudPrintOAuthClientIdentifier`, // custom MS Graph attribute name
						MarkdownDescription: "GUID of a client application authorized to retrieve OAuth tokens from the OAuth Authority.",
					},
					"enterprise_cloud_print_resource_identifier": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "OAuth resource URI for print service as configured in the Azure portal.",
					},
					"experience_block_device_discovery": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to enable device discovery UX. The _provider_ default value is `false`.",
					},
					"experience_block_error_dialog_when_no_sim": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `experienceBlockErrorDialogWhenNoSIM`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to allow the error dialog from displaying if no SIM card is detected. The _provider_ default value is `false`.",
					},
					"experience_block_task_switcher": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to enable task switching on the device. The _provider_ default value is `false`.",
					},
					"experience_do_not_sync_browser_settings": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "blockedWithUserOverride", "blocked"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Allow or prevent the syncing of Microsoft Edge Browser settings. Option for IT admins to prevent syncing across devices, but allow user override. / Allow(Not Configured) or prevent(Block) the syncing of Microsoft Edge Browser settings. Option to prevent syncing across devices, but allow user override; possible values are: `notConfigured` (Default – Allow syncing of browser settings across devices.), `blockedWithUserOverride` (Prevent syncing of browser settings across user devices, allow user override of setting.), `blocked` (Absolutely prevent syncing of browser settings across user devices.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"find_my_files": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Controls if the user can configure search to Find My Files mode, which searches files in secondary hard drives and also outside of the user profile. Find My Files does not allow users to search files or locations to which they do not have access. / Possible values of a property; possible values are: `notConfigured` (Device default value, no intent.), `enabled` (Enables the setting on the device.), `disabled` (Disables the setting on the device.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"game_dvr_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block DVR and broadcasting. The _provider_ default value is `false`.",
					},
					"ink_workspace_access": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Controls the user access to the ink workspace, from the desktop and from above the lock screen. / Values for the InkWorkspaceAccess setting; possible values are: `notConfigured` (Not configured.), `enabled` (Enabled.), `disabled` (Disabled.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"ink_workspace_access_state": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "blocked", "allowed"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Controls the user access to the ink workspace, from the desktop and from above the lock screen. / State Management Setting; possible values are: `notConfigured` (Not configured.), `blocked` (Blocked.), `allowed` (Allowed.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"ink_workspace_block_suggested_apps": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Specify whether to show recommended app suggestions in the ink workspace. The _provider_ default value is `false`.",
					},
					"internet_sharing_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Block the user from using internet sharing. The _provider_ default value is `false`.",
					},
					"location_services_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Block the user from location services. The _provider_ default value is `false`.",
					},
					"lock_screen_activate_apps_with_voice": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "This policy setting specifies whether Windows apps can be activated by voice while the system is locked. / Possible values of a property; possible values are: `notConfigured` (Device default value, no intent.), `enabled` (Enables the setting on the device.), `disabled` (Disables the setting on the device.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"lock_screen_allow_timeout_configuration": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Specify whether to show a user-configurable setting to control the screen timeout while on the lock screen of Windows 10 Mobile devices. If this policy is set to Allow, the value set by lockScreenTimeoutInSeconds is ignored. The _provider_ default value is `false`.",
					},
					"lock_screen_block_action_center_notifications": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block action center notifications over lock screen. The _provider_ default value is `false`.",
					},
					"lock_screen_block_cortana": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not the user can interact with Cortana using speech while the system is locked. The _provider_ default value is `false`.",
					},
					"lock_screen_block_toast_notifications": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether to allow toast notifications above the device lock screen. The _provider_ default value is `false`.",
					},
					"lock_screen_timeout_in_seconds": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Set the duration (in seconds) from the screen locking to the screen turning off for Windows 10 Mobile devices. Supported values are 11-1800. Valid values 11 to 1800",
					},
					"logon_block_fast_user_switching": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Disables the ability to quickly switch between users that are logged on simultaneously without logging off. The _provider_ default value is `false`.",
					},
					"messaging_block_mms": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `messagingBlockMMS`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to block the MMS send/receive functionality on the device. The _provider_ default value is `false`.",
					},
					"messaging_block_rich_communication_services": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the RCS send/receive functionality on the device. The _provider_ default value is `false`.",
					},
					"messaging_block_sync": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block text message back up and restore and Messaging Everywhere. The _provider_ default value is `false`.",
					},
					"microsoft_account_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Block a Microsoft account. The _provider_ default value is `false`.",
					},
					"microsoft_account_block_settings_sync": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Block Microsoft account settings sync. The _provider_ default value is `false`.",
					},
					"microsoft_account_sign_in_assistant_settings": schema.StringAttribute{
						Optional:            true,
						Validators:          []validator.String{stringvalidator.OneOf("notConfigured", "disabled")},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Controls the Microsoft Account Sign-In Assistant (wlidsvc) NT service. / Values for the SignInAssistantSettings; possible values are: `notConfigured` (Not configured - wlidsvc Start will be set to SERVICE_DEMAND_START.), `disabled` (Disabled - wlidsvc Start will be set to SERVICE_DISABLED.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"network_proxy_apply_settings_device_wide": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "If set, proxy settings will be applied to all processes and accounts in the device. Otherwise, it will be applied to the user account that’s enrolled into MDM. The _provider_ default value is `false`.",
					},
					"network_proxy_automatic_configuration_url": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Address to the proxy auto-config (PAC) script you want to use.",
					},
					"network_proxy_disable_auto_detect": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Disable automatic detection of settings. If enabled, the system will try to find the path to a proxy auto-config (PAC) script. The _provider_ default value is `false`.",
					},
					"network_proxy_server": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // windows10NetworkProxyServer
							"address": schema.StringAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
								Computed:            true,
								MarkdownDescription: "Address to the proxy server. Specify an address in the format <server>\\[“:”<port>\\]. The _provider_ default value is `\"\"`.",
							},
							"exceptions": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
								Computed:            true,
								MarkdownDescription: "Addresses that should not use the proxy server. The system will not use the proxy server for addresses beginning with what is specified in this node. The _provider_ default value is `[]`.",
							},
							"use_for_local_addresses": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: "Specifies whether the proxy server should be used for local (intranet) addresses. The _provider_ default value is `false`.",
							},
						},
						MarkdownDescription: "Specifies manual proxy server settings. / Network Proxy Server Policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windows10networkproxyserver?view=graph-rest-beta",
					},
					"nfc_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Block the user from using near field communication. The _provider_ default value is `false`.",
					},
					"one_drive_disable_file_sync": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Gets or sets a value allowing IT admins to prevent apps and features from working with files on OneDrive. The _provider_ default value is `false`.",
					},
					"password_block_simple": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Specify whether PINs or passwords such as \"1111\" or \"1234\" are allowed. For Windows 10 desktops, it also controls the use of picture passwords. The _provider_ default value is `false`.",
					},
					"password_expiration_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The password expiration in days. Valid values 0 to 730",
					},
					"password_minimum_age_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "This security setting determines the period of time (in days) that a password must be used before the user can change it. Valid values 0 to 998",
					},
					"password_minimum_character_set_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The number of character sets required in the password.",
					},
					"password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The minimum password length. Valid values 4 to 16",
					},
					"password_minutes_of_inactivity_before_screen_timeout": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The minutes of inactivity before the screen times out.",
					},
					"password_previous_password_block_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The number of previous passwords to prevent reuse of. Valid values 0 to 50",
					},
					"password_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to require the user to have a password. The _provider_ default value is `false`.",
					},
					"password_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "alphanumeric", "numeric"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: "The required password type. / Possible values of required passwords; possible values are: `deviceDefault` (Device default value, no intent.), `alphanumeric` (Alphanumeric password required.), `numeric` (Numeric password required.). The _provider_ default value is `\"deviceDefault\"`.",
					},
					"password_require_when_resume_from_idle_state": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to require a password upon resuming from an idle state. The _provider_ default value is `false`.",
					},
					"password_sign_in_failure_count_before_factory_reset": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The number of sign in failures before factory reset. Valid values 0 to 999",
					},
					"personalization_desktop_image_url": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "A http or https Url to a jpg, jpeg or png image that needs to be downloaded and used as the Desktop Image or a file Url to a local image on the file system that needs to used as the Desktop Image.",
					},
					"personalization_lock_screen_image_url": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "A http or https Url to a jpg, jpeg or png image that neeeds to be downloaded and used as the Lock Screen Image or a file Url to a local image on the file system that needs to be used as the Lock Screen Image.",
					},
					"power_button_action_on_battery": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "noAction", "sleep", "hibernate", "shutdown"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "This setting specifies the action that Windows takes when a user presses the Power button while on battery. / Power action types; possible values are: `notConfigured` (Not configured), `noAction` (No action), `sleep` (Put device in sleep state), `hibernate` (Put device in hibernate state), `shutdown` (Shutdown device). The _provider_ default value is `\"notConfigured\"`.",
					},
					"power_button_action_plugged_in": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "noAction", "sleep", "hibernate", "shutdown"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "This setting specifies the action that Windows takes when a user presses the Power button while plugged in. / Power action types; possible values are: `notConfigured` (Not configured), `noAction` (No action), `sleep` (Put device in sleep state), `hibernate` (Put device in hibernate state), `shutdown` (Shutdown device). The _provider_ default value is `\"notConfigured\"`.",
					},
					"power_hybrid_sleep_on_battery": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "This setting allows you to turn off hybrid sleep while on battery. If you set this setting to disable, a hiberfile is not generated when the system transitions to sleep (Stand By). If you set this setting to enable or do not configure this policy setting, users control this setting. / Possible values of a property; possible values are: `notConfigured` (Device default value, no intent.), `enabled` (Enables the setting on the device.), `disabled` (Disables the setting on the device.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"power_hybrid_sleep_plugged_in": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "This setting allows you to turn off hybrid sleep while plugged in. If you set this setting to disable, a hiberfile is not generated when the system transitions to sleep (Stand By). If you set this setting to enable or do not configure this policy setting, users control this setting. / Possible values of a property; possible values are: `notConfigured` (Device default value, no intent.), `enabled` (Enables the setting on the device.), `disabled` (Disables the setting on the device.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"power_lid_close_action_on_battery": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "noAction", "sleep", "hibernate", "shutdown"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "This setting specifies the action that Windows takes when a user closes the lid on a mobile PC while on battery. / Power action types; possible values are: `notConfigured` (Not configured), `noAction` (No action), `sleep` (Put device in sleep state), `hibernate` (Put device in hibernate state), `shutdown` (Shutdown device). The _provider_ default value is `\"notConfigured\"`.",
					},
					"power_lid_close_action_plugged_in": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "noAction", "sleep", "hibernate", "shutdown"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "This setting specifies the action that Windows takes when a user closes the lid on a mobile PC while plugged in. / Power action types; possible values are: `notConfigured` (Not configured), `noAction` (No action), `sleep` (Put device in sleep state), `hibernate` (Put device in hibernate state), `shutdown` (Shutdown device). The _provider_ default value is `\"notConfigured\"`.",
					},
					"power_sleep_button_action_on_battery": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "noAction", "sleep", "hibernate", "shutdown"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "This setting specifies the action that Windows takes when a user presses the Sleep button while on battery. / Power action types; possible values are: `notConfigured` (Not configured), `noAction` (No action), `sleep` (Put device in sleep state), `hibernate` (Put device in hibernate state), `shutdown` (Shutdown device). The _provider_ default value is `\"notConfigured\"`.",
					},
					"power_sleep_button_action_plugged_in": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "noAction", "sleep", "hibernate", "shutdown"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "This setting specifies the action that Windows takes when a user presses the Sleep button while plugged in. / Power action types; possible values are: `notConfigured` (Not configured), `noAction` (No action), `sleep` (Put device in sleep state), `hibernate` (Put device in hibernate state), `shutdown` (Shutdown device). The _provider_ default value is `\"notConfigured\"`.",
					},
					"printer_block_addition": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Prevent user installation of additional printers from printers settings. The _provider_ default value is `false`.",
					},
					"printer_default_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Name (network host name) of an installed printer.",
					},
					"printer_names": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Automatically provision printers based on their names (network host names). The _provider_ default value is `[]`.",
					},
					"privacy_advertising_id": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "blocked", "allowed"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Enables or disables the use of advertising ID. Added in Windows 10, version 1607. / State Management Setting; possible values are: `notConfigured` (Not configured.), `blocked` (Blocked.), `allowed` (Allowed.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"privacy_auto_accept_pairing_and_consent_prompts": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow the automatic acceptance of the pairing and privacy user consent dialog when launching apps. The _provider_ default value is `false`.",
					},
					"privacy_block_activity_feed": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Blocks the usage of cloud based speech services for Cortana, Dictation, or Store applications. The _provider_ default value is `false`.",
					},
					"privacy_block_input_personalization": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the usage of cloud based speech services for Cortana, Dictation, or Store applications. The _provider_ default value is `false`.",
					},
					"privacy_block_publish_user_activities": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Blocks the shared experiences/discovery of recently used resources in task switcher etc. The _provider_ default value is `false`.",
					},
					"privacy_disable_launch_experience": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "This policy prevents the privacy experience from launching during user logon for new and upgraded users.​ The _provider_ default value is `false`.",
					},
					"reset_protection_mode_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Block the user from reset protection mode. The _provider_ default value is `false`.",
					},
					"safe_search_filter": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userDefined", "strict", "moderate"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("userDefined")},
						Computed:            true,
						MarkdownDescription: "Specifies what filter level of safe search is required. / Specifies what level of safe search (filtering adult content) is required; possible values are: `userDefined` (User Defined, default value, no intent.), `strict` (Strict, highest filtering against adult content.), `moderate` (Moderate filtering against adult content (valid search results will not be filtered).). The _provider_ default value is `\"userDefined\"`.",
					},
					"screen_capture_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Block the user from taking Screenshots. The _provider_ default value is `false`.",
					},
					"search_block_diacritics": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Specifies if search can use diacritics. The _provider_ default value is `false`.",
					},
					"search_block_web_results": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the web search. The _provider_ default value is `false`.",
					},
					"search_disable_auto_language_detection": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Specifies whether to use automatic language detection when indexing content and properties. The _provider_ default value is `false`.",
					},
					"search_disable_indexer_backoff": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to disable the search indexer backoff feature. The _provider_ default value is `false`.",
					},
					"search_disable_indexing_encrypted_items": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block indexing of WIP-protected items to prevent them from appearing in search results for Cortana or Explorer. The _provider_ default value is `false`.",
					},
					"search_disable_indexing_removable_drive": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow users to add locations on removable drives to libraries and to be indexed. The _provider_ default value is `false`.",
					},
					"search_disable_location": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Specifies if search can use location information. The _provider_ default value is `false`.",
					},
					"search_disable_use_location": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Specifies if search can use location information. The _provider_ default value is `false`.",
					},
					"search_enable_automatic_index_size_manangement": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Specifies minimum amount of hard drive space on the same drive as the index location before indexing stops. The _provider_ default value is `false`.",
					},
					"search_enable_remote_queries": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block remote queries of this computer’s index. The _provider_ default value is `false`.",
					},
					"security_block_azure_ad_joined_devices_auto_encryption": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `securityBlockAzureADJoinedDevicesAutoEncryption`, // custom MS Graph attribute name
						MarkdownDescription: "Specify whether to allow automatic device encryption during OOBE when the device is Azure AD joined (desktop only). The _provider_ default value is `false`.",
					},
					"settings_block_accounts_page": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block access to Accounts in Settings app. The _provider_ default value is `false`.",
					},
					"settings_block_add_provisioning_package": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from installing provisioning packages. The _provider_ default value is `false`.",
					},
					"settings_block_apps_page": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block access to Apps in Settings app. The _provider_ default value is `false`.",
					},
					"settings_block_change_language": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from changing the language settings. The _provider_ default value is `false`.",
					},
					"settings_block_change_power_sleep": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from changing power and sleep settings. The _provider_ default value is `false`.",
					},
					"settings_block_change_region": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from changing the region settings. The _provider_ default value is `false`.",
					},
					"settings_block_change_system_time": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from changing date and time settings. The _provider_ default value is `false`.",
					},
					"settings_block_devices_page": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block access to Devices in Settings app. The _provider_ default value is `false`.",
					},
					"settings_block_ease_of_access_page": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block access to Ease of Access in Settings app. The _provider_ default value is `false`.",
					},
					"settings_block_edit_device_name": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from editing the device name. The _provider_ default value is `false`.",
					},
					"settings_block_gaming_page": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block access to Gaming in Settings app. The _provider_ default value is `false`.",
					},
					"settings_block_network_internet_page": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block access to Network & Internet in Settings app. The _provider_ default value is `false`.",
					},
					"settings_block_personalization_page": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block access to Personalization in Settings app. The _provider_ default value is `false`.",
					},
					"settings_block_privacy_page": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block access to Privacy in Settings app. The _provider_ default value is `false`.",
					},
					"settings_block_remove_provisioning_package": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the runtime configuration agent from removing provisioning packages. The _provider_ default value is `false`.",
					},
					"settings_block_settings_app": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block access to Settings app. The _provider_ default value is `false`.",
					},
					"settings_block_system_page": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block access to System in Settings app. The _provider_ default value is `false`.",
					},
					"settings_block_time_language_page": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block access to Time & Language in Settings app. The _provider_ default value is `false`.",
					},
					"settings_block_update_security_page": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block access to Update & Security in Settings app. The _provider_ default value is `false`.",
					},
					"shared_user_app_data_allowed": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block multiple users of the same app to share data. The _provider_ default value is `false`.",
					},
					"smartscreen_app_install_control": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "anywhere", "storeOnly", "recommendations", "preferStore"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						Description:         `smartScreenAppInstallControl`, // custom MS Graph attribute name
						MarkdownDescription: "Added in Windows 10, version 1703. Allows IT Admins to control whether users are allowed to install apps from places other than the Store. / App Install control Setting; possible values are: `notConfigured` (Not configured), `anywhere` (Turn off app recommendations), `storeOnly` (Allow apps from Store only), `recommendations` (Show me app recommendations), `preferStore` (Warn me before installing apps from outside the Store). The _provider_ default value is `\"notConfigured\"`.",
					},
					"smartscreen_block_prompt_override": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `smartScreenBlockPromptOverride`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not users can override SmartScreen Filter warnings about potentially malicious websites. The _provider_ default value is `false`.",
					},
					"smartscreen_block_prompt_override_for_files": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `smartScreenBlockPromptOverrideForFiles`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not users can override the SmartScreen Filter warnings about downloading unverified files. The _provider_ default value is `false`.",
					},
					"start_block_unpinning_apps_from_taskbar": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block the user from unpinning apps from taskbar. The _provider_ default value is `false`.",
					},
					"start_menu_app_list_visibility": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							wpvalidator.FlagEnumValues("userDefined", "collapse", "remove", "disableSettingsApp"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("userDefined")},
						Computed:            true,
						MarkdownDescription: "Setting the value of this collapses the app list, removes the app list entirely, or disables the corresponding toggle in the Settings app. / Type of start menu app list visibility; possible values are: `userDefined` (User defined. Default value.), `collapse` (Collapse the app list on the start menu.), `remove` (Removes the app list entirely from the start menu.), `disableSettingsApp` (Disables the corresponding toggle (Collapse or Remove) in the Settings app.). The _provider_ default value is `\"userDefined\"`.",
					},
					"start_menu_hide_change_account_settings": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Enabling this policy hides the change account setting from appearing in the user tile in the start menu. The _provider_ default value is `false`.",
					},
					"start_menu_hide_frequently_used_apps": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Enabling this policy hides the most used apps from appearing on the start menu and disables the corresponding toggle in the Settings app. The _provider_ default value is `false`.",
					},
					"start_menu_hide_hibernate": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Enabling this policy hides hibernate from appearing in the power button in the start menu. The _provider_ default value is `false`.",
					},
					"start_menu_hide_lock": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Enabling this policy hides lock from appearing in the user tile in the start menu. The _provider_ default value is `false`.",
					},
					"start_menu_hide_power_button": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Enabling this policy hides the power button from appearing in the start menu. The _provider_ default value is `false`.",
					},
					"start_menu_hide_recent_jump_lists": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Enabling this policy hides recent jump lists from appearing on the start menu/taskbar and disables the corresponding toggle in the Settings app. The _provider_ default value is `false`.",
					},
					"start_menu_hide_recently_added_apps": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Enabling this policy hides recently added apps from appearing on the start menu and disables the corresponding toggle in the Settings app. The _provider_ default value is `false`.",
					},
					"start_menu_hide_restart_options": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Enabling this policy hides “Restart/Update and Restart” from appearing in the power button in the start menu. The _provider_ default value is `false`.",
					},
					"start_menu_hide_shutdown": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `startMenuHideShutDown`, // custom MS Graph attribute name
						MarkdownDescription: "Enabling this policy hides shut down/update and shut down from appearing in the power button in the start menu. The _provider_ default value is `false`.",
					},
					"start_menu_hide_sign_out": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Enabling this policy hides sign out from appearing in the user tile in the start menu. The _provider_ default value is `false`.",
					},
					"start_menu_hide_sleep": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Enabling this policy hides sleep from appearing in the power button in the start menu. The _provider_ default value is `false`.",
					},
					"start_menu_hide_switch_account": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Enabling this policy hides switch account from appearing in the user tile in the start menu. The _provider_ default value is `false`.",
					},
					"start_menu_hide_user_tile": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Enabling this policy hides the user tile from appearing in the start menu. The _provider_ default value is `false`.",
					},
					"start_menu_layout_edge_assets_xml": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "This policy setting allows you to import Edge assets to be used with startMenuLayoutXml policy. Start layout can contain secondary tile from Edge app which looks for Edge local asset file. Edge local asset would not exist and cause Edge secondary tile to appear empty in this case. This policy only gets applied when startMenuLayoutXml policy is modified. The value should be a UTF-8 Base64 encoded byte array.",
					},
					"start_menu_layout_xml": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Allows admins to override the default Start menu layout and prevents the user from changing it. The layout is modified by specifying an XML file based on a layout modification schema. XML needs to be in a UTF8 encoded byte array format.",
					},
					"start_menu_mode": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userDefined", "fullScreen", "nonFullScreen"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("userDefined")},
						Computed:            true,
						MarkdownDescription: "Allows admins to decide how the Start menu is displayed. / Type of display modes for the start menu; possible values are: `userDefined` (User defined. Default value.), `fullScreen` (Full screen.), `nonFullScreen` (Non-full screen.). The _provider_ default value is `\"userDefined\"`.",
					},
					"start_menu_pinned_folder_documents": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "hide", "show"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Enforces the visibility (Show/Hide) of the Documents folder shortcut on the Start menu. / Generic visibility state; possible values are: `notConfigured` (Not configured.), `hide` (Hide.), `show` (Show.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"start_menu_pinned_folder_downloads": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "hide", "show"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Enforces the visibility (Show/Hide) of the Downloads folder shortcut on the Start menu. / Generic visibility state; possible values are: `notConfigured` (Not configured.), `hide` (Hide.), `show` (Show.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"start_menu_pinned_folder_file_explorer": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "hide", "show"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Enforces the visibility (Show/Hide) of the FileExplorer shortcut on the Start menu. / Generic visibility state; possible values are: `notConfigured` (Not configured.), `hide` (Hide.), `show` (Show.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"start_menu_pinned_folder_home_group": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "hide", "show"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Enforces the visibility (Show/Hide) of the HomeGroup folder shortcut on the Start menu. / Generic visibility state; possible values are: `notConfigured` (Not configured.), `hide` (Hide.), `show` (Show.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"start_menu_pinned_folder_music": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "hide", "show"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Enforces the visibility (Show/Hide) of the Music folder shortcut on the Start menu. / Generic visibility state; possible values are: `notConfigured` (Not configured.), `hide` (Hide.), `show` (Show.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"start_menu_pinned_folder_network": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "hide", "show"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Enforces the visibility (Show/Hide) of the Network folder shortcut on the Start menu. / Generic visibility state; possible values are: `notConfigured` (Not configured.), `hide` (Hide.), `show` (Show.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"start_menu_pinned_folder_personal_folder": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "hide", "show"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Enforces the visibility (Show/Hide) of the PersonalFolder shortcut on the Start menu. / Generic visibility state; possible values are: `notConfigured` (Not configured.), `hide` (Hide.), `show` (Show.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"start_menu_pinned_folder_pictures": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "hide", "show"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Enforces the visibility (Show/Hide) of the Pictures folder shortcut on the Start menu. / Generic visibility state; possible values are: `notConfigured` (Not configured.), `hide` (Hide.), `show` (Show.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"start_menu_pinned_folder_settings": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "hide", "show"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Enforces the visibility (Show/Hide) of the Settings folder shortcut on the Start menu. / Generic visibility state; possible values are: `notConfigured` (Not configured.), `hide` (Hide.), `show` (Show.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"start_menu_pinned_folder_videos": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "hide", "show"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Enforces the visibility (Show/Hide) of the Videos folder shortcut on the Start menu. / Generic visibility state; possible values are: `notConfigured` (Not configured.), `hide` (Hide.), `show` (Show.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"storage_block_removable_storage": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Block the user from using removable storage. The _provider_ default value is `false`.",
					},
					"storage_require_mobile_device_encryption": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicating whether or not to require encryption on a mobile device. The _provider_ default value is `false`.",
					},
					"storage_restrict_app_data_to_system_volume": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether application data is restricted to the system drive. The _provider_ default value is `false`.",
					},
					"storage_restrict_app_install_to_system_volume": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether the installation of applications is restricted to the system drive. The _provider_ default value is `false`.",
					},
					"system_telemetry_proxy_server": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Gets or sets the fully qualified domain name (FQDN) or IP address of a proxy server to forward Connected User Experiences and Telemetry requests.",
					},
					"task_manager_block_end_task": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Specify whether non-administrators can use Task Manager to end tasks. The _provider_ default value is `false`.",
					},
					"tenant_lockdown_require_network_during_out_of_box_experience": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether the device is required to connect to the network. The _provider_ default value is `false`.",
					},
					"uninstall_built_in_apps": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to uninstall a fixed list of built-in Windows apps. The _provider_ default value is `false`.",
					},
					"usb_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Block the user from USB connection. The _provider_ default value is `false`.",
					},
					"voice_recording_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Block the user from voice recording. The _provider_ default value is `false`.",
					},
					"webrtc_block_localhost_ip_address": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `webRtcBlockLocalhostIpAddress`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not user's localhost IP address is displayed while making phone calls using the WebRTC. The _provider_ default value is `false`.",
					},
					"wifi_block_automatic_connect_hotspots": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `wiFiBlockAutomaticConnectHotspots`, // custom MS Graph attribute name
						MarkdownDescription: "Indicating whether or not to block automatically connecting to Wi-Fi hotspots. Has no impact if Wi-Fi is blocked. The _provider_ default value is `false`.",
					},
					"wifi_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `wiFiBlocked`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to Block the user from using Wi-Fi. The _provider_ default value is `false`.",
					},
					"wifi_block_manual_configuration": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `wiFiBlockManualConfiguration`, // custom MS Graph attribute name
						MarkdownDescription: "Indicates whether or not to Block the user from using Wi-Fi manual configuration. The _provider_ default value is `false`.",
					},
					"wifi_scan_interval": schema.Int64Attribute{
						Optional:            true,
						Description:         `wiFiScanInterval`, // custom MS Graph attribute name
						MarkdownDescription: "Specify how often devices scan for Wi-Fi networks. Supported values are 1-500, where 100 = default, and 500 = low frequency. Valid values 1 to 500",
					},
					"windows10_apps_force_update_schedule": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // windows10AppsForceUpdateSchedule
							"recurrence": schema.StringAttribute{
								Optional: true,
								Validators: []validator.String{
									stringvalidator.OneOf("none", "daily", "weekly", "monthly"),
								},
								PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
								Computed:            true,
								MarkdownDescription: "Recurrence schedule. / Possible values for App update on Windows10 recurrence; possible values are: `none` (Default value, specifies a single occurence.), `daily` (Daily.), `weekly` (Weekly.), `monthly` (Monthly.). The _provider_ default value is `\"none\"`.",
							},
							"run_immediately_if_after_start_date_time": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: "If true, runs the task immediately if StartDateTime is in the past, else, runs at the next recurrence. The _provider_ default value is `false`.",
							},
							"start_date_time": schema.StringAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
								Computed:            true,
								MarkdownDescription: "The start time for the force restart. The _provider_ default value is `\"\"`.",
							},
						},
						Description:         `windows10AppsForceUpdateSchedule`, // custom MS Graph attribute name
						MarkdownDescription: "Windows 10 force update schedule for Apps. / Windows 10 force update schedule for Apps / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windows10appsforceupdateschedule?view=graph-rest-beta",
					},
					"windows_spotlight_block_consumer_specific_features": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Allows IT admins to block experiences that are typically for consumers only, such as Start suggestions, Membership notifications, Post-OOBE app install and redirect tiles. The _provider_ default value is `false`.",
					},
					"windows_spotlight_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Allows IT admins to turn off all Windows Spotlight features. The _provider_ default value is `false`.",
					},
					"windows_spotlight_block_on_action_center": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Block suggestions from Microsoft that show after each OS clean install, upgrade or in an on-going basis to introduce users to what is new or changed. The _provider_ default value is `false`.",
					},
					"windows_spotlight_block_tailored_experiences": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Block personalized content in Windows spotlight based on user’s device usage. The _provider_ default value is `false`.",
					},
					"windows_spotlight_block_third_party_notifications": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Block third party content delivered via Windows Spotlight. The _provider_ default value is `false`.",
					},
					"windows_spotlight_block_welcome_experience": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Block Windows Spotlight Windows welcome experience. The _provider_ default value is `false`.",
					},
					"windows_spotlight_block_windows_tips": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Allows IT admins to turn off the popup of Windows Tips. The _provider_ default value is `false`.",
					},
					"windows_spotlight_configure_on_lock_screen": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "disabled", "enabled"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "Specifies the type of Spotlight. / Allows IT admind to set a predefined default search engine for MDM-Controlled devices; possible values are: `notConfigured` (Spotlight on lock screen is not configured), `disabled` (Disable Windows Spotlight on lock screen), `enabled` (Enable Windows Spotlight on lock screen). The _provider_ default value is `\"notConfigured\"`.",
					},
					"windows_store_block_auto_update": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block automatic update of apps from Windows Store. The _provider_ default value is `false`.",
					},
					"windows_store_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to Block the user from using the Windows store. The _provider_ default value is `false`.",
					},
					"windows_store_enable_private_store_only": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to enable Private Store Only. The _provider_ default value is `false`.",
					},
					"wireless_display_block_projection_to_this_device": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow other devices from discovering this PC for projection. The _provider_ default value is `false`.",
					},
					"wireless_display_block_user_input_from_receiver": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to allow user input from wireless display receiver. The _provider_ default value is `false`.",
					},
					"wireless_display_require_pin_for_pairing": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to require a PIN for new devices to initiate pairing. The _provider_ default value is `false`.",
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: "This topic provides descriptions of the declared methods, properties and relationships exposed by the windows10GeneralConfiguration resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windows10generalconfiguration?view=graph-rest-beta",
			},
		},
	},
	MarkdownDescription: "Device Configuration. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceconfiguration?view=graph-rest-beta ||| MS Graph: Device configuration",
}

var deviceConfigurationDeviceConfigurationValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("android_device_owner_general_device"),
	path.MatchRelative().AtParent().AtName("android_work_profile_general_device"),
	path.MatchRelative().AtParent().AtName("edition_upgrade"),
	path.MatchRelative().AtParent().AtName("ios_custom"),
	path.MatchRelative().AtParent().AtName("ios_device_features"),
	path.MatchRelative().AtParent().AtName("ios_eas_email_profile"),
	path.MatchRelative().AtParent().AtName("ios_general_device"),
	path.MatchRelative().AtParent().AtName("ios_update"),
	path.MatchRelative().AtParent().AtName("ios_vpn"),
	path.MatchRelative().AtParent().AtName("macos_custom"),
	path.MatchRelative().AtParent().AtName("macos_custom_app"),
	path.MatchRelative().AtParent().AtName("macos_device_features"),
	path.MatchRelative().AtParent().AtName("macos_extensions"),
	path.MatchRelative().AtParent().AtName("macos_software_update"),
	path.MatchRelative().AtParent().AtName("windows_health_monitoring"),
	path.MatchRelative().AtParent().AtName("windows_update_for_business"),
	path.MatchRelative().AtParent().AtName("windows10_general"),
)

var deviceConfigurationAppListItemAttributes = map[string]schema.Attribute{ // appListItem
	"app_id": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The application or bundle identifier of the application",
	},
	"app_store_url": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The Store URL of the application",
	},
	"name": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The application name",
	},
	"publisher": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The publisher of the application",
	},
}

var deviceConfigurationKeyValuePairAttributes = map[string]schema.Attribute{ // keyValuePair
	"name": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "Name for this key-value pair",
	},
	"value": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "Value for this key-value pair",
	},
}

var deviceConfigurationAirPrintDestinationAttributes = map[string]schema.Attribute{ // airPrintDestination
	"force_tls": schema.BoolAttribute{
		Required:            true,
		MarkdownDescription: "If true AirPrint connections are secured by Transport Layer Security (TLS). Default is false. Available in iOS 11.0 and later.",
	},
	"ip_address": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The IP Address of the AirPrint destination.",
	},
	"port": schema.Int64Attribute{
		Optional:            true,
		MarkdownDescription: "The listening port of the AirPrint destination. If this key is not specified AirPrint will use the default port. Available in iOS 11.0 and later.",
	},
	"resource_path": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The Resource Path associated with the printer. This corresponds to the rp parameter of the _ipps.tcp Bonjour record. For example: printers/Canon_MG5300_series, printers/Xerox_Phaser_7600, ipp/print, Epson_IPP_Printer.",
	},
}

var deviceConfigurationIosWebContentFilterBaseValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("auto_filter"),
	path.MatchRelative().AtParent().AtName("specific_websites_access"),
)

var deviceConfigurationIosBookmarkAttributes = map[string]schema.Attribute{ // iosBookmark
	"bookmark_folder": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The folder into which the bookmark should be added in Safari",
	},
	"display_name": schema.StringAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
		Computed:            true,
		MarkdownDescription: "The display name of the bookmark. The _provider_ default value is `\"\"`.",
	},
	"url": schema.StringAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
		Computed:            true,
		MarkdownDescription: "URL allowed to access. The _provider_ default value is `\"\"`.",
	},
}

var deviceConfigurationIosSingleSignOnExtensionValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("azure_ad"),
	path.MatchRelative().AtParent().AtName("credential"),
	path.MatchRelative().AtParent().AtName("kerberos"),
	path.MatchRelative().AtParent().AtName("redirect"),
)

var deviceConfigurationKeyTypedValuePairAttributes = map[string]schema.Attribute{ // keyTypedValuePair
	"key": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The string key of the key-value pair.",
	},
	"boolean": generic.OdataDerivedTypeNestedAttributeRs{
		DerivedType: "#microsoft.graph.keyBooleanValuePair",
		SingleNestedAttribute: schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // keyBooleanValuePair
				"value": schema.BoolAttribute{
					Required:            true,
					MarkdownDescription: "The Boolean value of the key-value pair.",
				},
			},
			Validators:          []validator.Object{deviceConfigurationKeyTypedValuePairValidator},
			MarkdownDescription: "A key-value pair with a string key and a Boolean value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keybooleanvaluepair?view=graph-rest-beta",
		},
	},
	"integer": generic.OdataDerivedTypeNestedAttributeRs{
		DerivedType: "#microsoft.graph.keyIntegerValuePair",
		SingleNestedAttribute: schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // keyIntegerValuePair
				"value": schema.Int64Attribute{
					Required:            true,
					MarkdownDescription: "The integer value of the key-value pair.",
				},
			},
			Validators:          []validator.Object{deviceConfigurationKeyTypedValuePairValidator},
			MarkdownDescription: "A key-value pair with a string key and an integer value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyintegervaluepair?view=graph-rest-beta",
		},
	},
	"real": generic.OdataDerivedTypeNestedAttributeRs{
		DerivedType: "#microsoft.graph.keyRealValuePair",
		SingleNestedAttribute: schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // keyRealValuePair
				"value": schema.NumberAttribute{
					Required:            true,
					MarkdownDescription: "The real (floating-point) value of the key-value pair.",
				},
			},
			Validators:          []validator.Object{deviceConfigurationKeyTypedValuePairValidator},
			MarkdownDescription: "A key-value pair with a string key and a real (floating-point) value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyrealvaluepair?view=graph-rest-beta",
		},
	},
	"string": generic.OdataDerivedTypeNestedAttributeRs{
		DerivedType: "#microsoft.graph.keyStringValuePair",
		SingleNestedAttribute: schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // keyStringValuePair
				"value": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "The string value of the key-value pair.",
				},
			},
			Validators:          []validator.Object{deviceConfigurationKeyTypedValuePairValidator},
			MarkdownDescription: "A key-value pair with a string key and a string value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keystringvaluepair?view=graph-rest-beta",
		},
	},
}

var deviceConfigurationKeyTypedValuePairValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("boolean"),
	path.MatchRelative().AtParent().AtName("integer"),
	path.MatchRelative().AtParent().AtName("real"),
	path.MatchRelative().AtParent().AtName("string"),
)

var deviceConfigurationIpRangeAttributes = map[string]schema.Attribute{ // ipRange
	"v4": generic.OdataDerivedTypeNestedAttributeRs{
		DerivedType: "#microsoft.graph.iPv4Range",
		SingleNestedAttribute: schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // iPv4Range
				"lower_address": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Lower address.",
				},
				"upper_address": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Upper address.",
				},
			},
			Validators:          []validator.Object{deviceConfigurationIpRangeValidator},
			MarkdownDescription: "IPv4 Range definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-ipv4range?view=graph-rest-beta",
		},
	},
	"v4_cidr": generic.OdataDerivedTypeNestedAttributeRs{
		DerivedType: "#microsoft.graph.iPv4CidrRange",
		SingleNestedAttribute: schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // iPv4CidrRange
				"cidr_address": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "IPv4 address in CIDR notation. Not nullable.",
				},
			},
			Validators:          []validator.Object{deviceConfigurationIpRangeValidator},
			MarkdownDescription: "Represents an IPv4 range using the Classless Inter-Domain Routing (CIDR) notation. / https://learn.microsoft.com/en-us/graph/api/resources/ipv4cidrrange?view=graph-rest-beta",
		},
	},
	"v6": generic.OdataDerivedTypeNestedAttributeRs{
		DerivedType: "#microsoft.graph.iPv6Range",
		SingleNestedAttribute: schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // iPv6Range
				"lower_address": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Lower address.",
				},
				"upper_address": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Upper address.",
				},
			},
			Validators:          []validator.Object{deviceConfigurationIpRangeValidator},
			MarkdownDescription: "IPv6 Range definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-ipv6range?view=graph-rest-beta",
		},
	},
	"v6_cidr": generic.OdataDerivedTypeNestedAttributeRs{
		DerivedType: "#microsoft.graph.iPv6CidrRange",
		SingleNestedAttribute: schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // iPv6CidrRange
				"cidr_address": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "IPv6 address in CIDR notation. Not nullable.",
				},
			},
			Validators:          []validator.Object{deviceConfigurationIpRangeValidator},
			MarkdownDescription: "Represents an IPv6 range using the Classless Inter-Domain Routing (CIDR) notation. / https://learn.microsoft.com/en-us/graph/api/resources/ipv6cidrrange?view=graph-rest-beta",
		},
	},
}

var deviceConfigurationIpRangeValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("v4"),
	path.MatchRelative().AtParent().AtName("v4_cidr"),
	path.MatchRelative().AtParent().AtName("v6"),
	path.MatchRelative().AtParent().AtName("v6_cidr"),
)

var deviceConfigurationMacOSSingleSignOnExtensionValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("azure_ad"),
	path.MatchRelative().AtParent().AtName("credential"),
	path.MatchRelative().AtParent().AtName("kerberos"),
	path.MatchRelative().AtParent().AtName("redirect"),
)

var deviceConfigurationWindowsUpdateInstallScheduleTypeValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("active_hours"),
	path.MatchRelative().AtParent().AtName("scheduled"),
)

var deviceConfigurationEdgeHomeButtonConfigurationValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("hidden"),
	path.MatchRelative().AtParent().AtName("loads_start_page"),
	path.MatchRelative().AtParent().AtName("opens_custom_url"),
	path.MatchRelative().AtParent().AtName("opens_new_tab"),
)

var deviceConfigurationEdgeSearchEngineBaseValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("custom"),
	path.MatchRelative().AtParent().AtName("predefined"),
)

type deviceConfigurationEditionUpgradeProductKeyPlanModifier struct{ generic.EmptyDescriber }

var _ planmodifier.String = deviceConfigurationEditionUpgradeProductKeyPlanModifier{}

func (apm deviceConfigurationEditionUpgradeProductKeyPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, res *planmodifier.StringResponse) {
	// do nothing on create or on destroy
	if req.StateValue.IsNull() || req.PlanValue.IsNull() {
		return
	}

	stateValue := req.StateValue.ValueString()
	planValue := req.PlanValue.ValueString()
	stateValueShort := strings.TrimLeft(stateValue, "#-")
	planValueShort := planValue[len(planValue)-len(stateValueShort):]
	if planValueShort == stateValueShort {
		res.PlanValue = req.StateValue
	}
}
