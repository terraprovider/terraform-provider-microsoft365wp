package services

import (
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

	DeviceConfigurationPluralDataSource = generic.CreateGenericDataSourcePluralFromSingular(
		&DeviceConfigurationSingularDataSource, "")
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
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: `DateTime the object was created.`,
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: `Admin provided description of the Device Configuration.`,
		},
		"device_management_applicability_rule_device_mode": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // deviceManagementApplicabilityRuleDeviceMode
				"device_mode": schema.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf("standardConfiguration", "sModeConfiguration"),
					},
					MarkdownDescription: `Applicability rule for device mode.`,
				},
				"name": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: `Name for object.`,
				},
				"rule_type": schema.StringAttribute{
					Required:            true,
					Validators:          []validator.String{stringvalidator.OneOf("include", "exclude")},
					MarkdownDescription: `Applicability Rule type.`,
				},
			},
			MarkdownDescription: `The device mode applicability rule for this Policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceManagementApplicabilityRuleDeviceMode?view=graph-rest-beta`,
		},
		"device_management_applicability_rule_os_edition": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // deviceManagementApplicabilityRuleOsEdition
				"name": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: `Name for object.`,
				},
				"os_edition_types": schema.SetAttribute{
					ElementType: types.StringType,
					Required:    true,
					Validators: []validator.Set{
						setvalidator.ValueStringsAre(
							stringvalidator.OneOf("windows10Enterprise", "windows10EnterpriseN", "windows10Education", "windows10EducationN", "windows10MobileEnterprise", "windows10HolographicEnterprise", "windows10Professional", "windows10ProfessionalN", "windows10ProfessionalEducation", "windows10ProfessionalEducationN", "windows10ProfessionalWorkstation", "windows10ProfessionalWorkstationN", "notConfigured", "windows10Home", "windows10HomeChina", "windows10HomeN", "windows10HomeSingleLanguage", "windows10Mobile", "windows10IoTCore", "windows10IoTCoreCommercial"),
						),
					},
					MarkdownDescription: `Applicability rule OS edition type.`,
				},
				"rule_type": schema.StringAttribute{
					Required:            true,
					Validators:          []validator.String{stringvalidator.OneOf("include", "exclude")},
					MarkdownDescription: `Applicability Rule type.`,
				},
			},
			MarkdownDescription: `The OS edition applicability for this Policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceManagementApplicabilityRuleOsEdition?view=graph-rest-beta`,
		},
		"device_management_applicability_rule_os_version": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // deviceManagementApplicabilityRuleOsVersion
				"max_os_version": schema.StringAttribute{
					Optional:            true,
					Description:         `maxOSVersion`, // custom MS Graph attribute name
					MarkdownDescription: `Max OS version for Applicability Rule.`,
				},
				"min_os_version": schema.StringAttribute{
					Optional:            true,
					Description:         `minOSVersion`, // custom MS Graph attribute name
					MarkdownDescription: `Min OS version for Applicability Rule.`,
				},
				"name": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: `Name for object.`,
				},
				"rule_type": schema.StringAttribute{
					Required:            true,
					Validators:          []validator.String{stringvalidator.OneOf("include", "exclude")},
					MarkdownDescription: `Applicability Rule type.`,
				},
			},
			MarkdownDescription: `The OS version applicability rule for this Policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceManagementApplicabilityRuleOsVersion?view=graph-rest-beta`,
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: `Admin provided name of the device configuration.`,
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: `DateTime the object was last modified.`,
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: `List of Scope Tags for this Entity instance.`,
		},
		"supports_scope_tags": schema.BoolAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
		},
		"version": schema.Int64Attribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
			MarkdownDescription: `Version of the device configuration.`,
		},
		"assignments": deviceAndAppManagementAssignment,
		"android_device_owner_general_device": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.androidDeviceOwnerGeneralDeviceConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // androidDeviceOwnerGeneralDeviceConfiguration
					"accounts_block_modification": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not adding or removing accounts is disabled.`,
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
									MarkdownDescription: `Information about the app like Name, AppStoreUrl, Publisher and AppId / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-appListItem?view=graph-rest-beta`,
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
									MarkdownDescription: `List of scopes an app has been assigned.`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Specifies the list of managed apps with app details and its associated delegated scope(s). This collection can contain a maximum of 500 elements. / Represents one item in the list of managed apps with app details and its associated delegated scope(s). / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androidDeviceOwnerDelegatedScopeAppSetting?view=graph-rest-beta`,
					},
					"apps_allow_install_from_unknown_sources": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not the user is allowed to enable to unknown sources setting.`,
					},
					"apps_auto_update_policy": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "userChoice", "never", "wiFiOnly", "always"),
						},
						MarkdownDescription: `Indicates the value of the app auto update policy.`,
					},
					"apps_default_permission_policy": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "prompt", "autoGrant", "autoDeny"),
						},
						MarkdownDescription: `Indicates the permission policy for requests for runtime permissions if one is not defined for the app specifically.`,
					},
					"apps_recommend_skipping_first_use_hints": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether or not to recommend all apps skip any first-time-use hints they may have added.`,
					},
					"azure_ad_shared_device_data_clear_apps": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationAppListItemAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `A list of managed apps that will have their data cleared during a global sign-out in AAD shared device mode. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-appListItem?view=graph-rest-beta`,
					},
					"bluetooth_block_configuration": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to block a user from configuring bluetooth.`,
					},
					"bluetooth_block_contact_sharing": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to block a user from sharing contacts via bluetooth.`,
					},
					"camera_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to disable the use of the camera.`,
					},
					"cellular_block_wifi_tethering": schema.BoolAttribute{
						Optional:            true,
						Description:         `cellularBlockWiFiTethering`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to block Wi-Fi tethering.`,
					},
					"certificate_credential_configuration_disabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to block users from any certificate credential configuration.`,
					},
					"cross_profile_policies_allow_copy_paste": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not text copied from one profile (personal or work) can be pasted in the other.`,
					},
					"cross_profile_policies_allow_data_sharing": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "crossProfileDataSharingBlocked", "dataSharingFromWorkToPersonalBlocked", "crossProfileDataSharingAllowed", "unkownFutureValue"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: `Indicates whether data from one profile (personal or work) can be shared with apps in the other profile.`,
					},
					"cross_profile_policies_show_work_contacts_in_personal_profile": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not contacts stored in work profile are shown in personal profile contact searches/incoming calls.`,
					},
					"data_roaming_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to block a user from data roaming.`,
					},
					"date_time_configuration_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to block the user from manually changing the date or time on the device`,
					},
					"detailed_help_text": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // androidDeviceOwnerUserFacingMessage
							"default_message": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: `The default message displayed if the user's locale doesn't match with any of the localized messages`,
							},
							"localized_messages": schema.SetNestedAttribute{
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: deviceConfigurationKeyValuePairAttributes,
								},
								PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
								Computed:            true,
								MarkdownDescription: `The list of <locale, message> pairs. This collection can contain a maximum of 500 elements. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyValuePair?view=graph-rest-beta`,
							},
						},
						MarkdownDescription: `Represents the customized detailed help text provided to users when they attempt to modify managed settings on their device. / Represents a user-facing message with locale information as well as a default message to be used if the user's locale doesn't match with any of the localized messages / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androidDeviceOwnerUserFacingMessage?view=graph-rest-beta`,
					},
					"device_location_mode": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "disabled", "unknownFutureValue"),
						},
						MarkdownDescription: `Indicates the location setting configuration for fully managed devices (COBO) and corporate owned devices with a work profile (COPE)`,
					},
					"device_owner_lock_screen_message": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // androidDeviceOwnerUserFacingMessage
							"default_message": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: `The default message displayed if the user's locale doesn't match with any of the localized messages`,
							},
							"localized_messages": schema.SetNestedAttribute{
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: deviceConfigurationKeyValuePairAttributes,
								},
								PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
								Computed:            true,
								MarkdownDescription: `The list of <locale, message> pairs. This collection can contain a maximum of 500 elements. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyValuePair?view=graph-rest-beta`,
							},
						},
						MarkdownDescription: `Represents the customized lock screen message provided to users when they attempt to modify managed settings on their device. / Represents a user-facing message with locale information as well as a default message to be used if the user's locale doesn't match with any of the localized messages / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androidDeviceOwnerUserFacingMessage?view=graph-rest-beta`,
					},
					"enrollment_profile": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "dedicatedDevice", "fullyManaged"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: `Indicates which enrollment profile you want to configure.`,
					},
					"factory_reset_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not the factory reset option in settings is disabled.`,
					},
					"factory_reset_device_administrator_emails": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `List of Google account emails that will be required to authenticate after a device is factory reset before it can be set up.`,
					},
					"global_proxy": schema.SingleNestedAttribute{
						Optional:   true,
						Attributes: map[string]schema.Attribute{ // androidDeviceOwnerGlobalProxy
						},
						MarkdownDescription: `Proxy is set up directly with host, port and excluded hosts. / Android Device Owner Global Proxy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androidDeviceOwnerGlobalProxy?view=graph-rest-beta`,
					},
					"google_accounts_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not google accounts will be blocked.`,
					},
					"kiosk_customization_device_settings_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether a user can access the device's Settings app while in Kiosk Mode.`,
					},
					"kiosk_customization_power_button_actions_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether the power menu is shown when a user long presses the Power button of a device in Kiosk Mode.`,
					},
					"kiosk_customization_status_bar": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "notificationsAndSystemInfoEnabled", "systemInfoOnly"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: `Indicates whether system info and notifications are disabled in Kiosk Mode.`,
					},
					"kiosk_customization_system_error_warnings": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether system error dialogs for crashed or unresponsive apps are shown in Kiosk Mode.`,
					},
					"kiosk_customization_system_navigation": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "navigationEnabled", "homeButtonOnly"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: `Indicates which navigation features are enabled in Kiosk Mode.`,
					},
					"kiosk_mode_app_order_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether or not to enable app ordering in Kiosk Mode.`,
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
									MarkdownDescription: `Item to be arranged / Represents an item on the Android Device Owner Managed Home Screen (application, weblink or folder / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androidDeviceOwnerKioskModeHomeScreenItem?view=graph-rest-beta`,
								},
								"position": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: `Position of the item on the grid. Valid values 0 to 9999999`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `The ordering of items on Kiosk Mode Managed Home Screen. This collection can contain a maximum of 500 elements. / An item in the list of app positions that sets the order of items on the Managed Home Screen / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androidDeviceOwnerKioskModeAppPositionItem?view=graph-rest-beta`,
					},
					"kiosk_mode_apps": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationAppListItemAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `A list of managed apps that will be shown when the device is in Kiosk Mode. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-appListItem?view=graph-rest-beta`,
					},
					"kiosk_mode_apps_in_folder_ordered_by_name": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether or not to alphabetize applications within a folder in Kiosk Mode.`,
					},
					"kiosk_mode_bluetooth_configuration_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether or not to allow a user to configure Bluetooth settings in Kiosk Mode.`,
					},
					"kiosk_mode_debug_menu_easy_access_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether or not to allow a user to easy access to the debug menu in Kiosk Mode.`,
					},
					"kiosk_mode_exit_code": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Exit code to allow a user to escape from Kiosk Mode when the device is in Kiosk Mode.`,
					},
					"kiosk_mode_flashlight_configuration_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether or not to allow a user to use the flashlight in Kiosk Mode.`,
					},
					"kiosk_mode_folder_icon": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "darkSquare", "darkCircle", "lightSquare", "lightCircle"),
						},
						MarkdownDescription: `Folder icon configuration for managed home screen in Kiosk Mode.`,
					},
					"kiosk_mode_grid_height": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of rows for Managed Home Screen grid with app ordering enabled in Kiosk Mode. Valid values 1 to 9999999`,
					},
					"kiosk_mode_grid_width": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of columns for Managed Home Screen grid with app ordering enabled in Kiosk Mode. Valid values 1 to 9999999`,
					},
					"kiosk_mode_icon_size": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "smallest", "small", "regular", "large", "largest"),
						},
						MarkdownDescription: `Icon size configuration for managed home screen in Kiosk Mode.`,
					},
					"kiosk_mode_lock_home_screen": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether or not to lock home screen to the end user in Kiosk Mode.`,
					},
					"kiosk_mode_managed_folders": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // androidDeviceOwnerKioskModeManagedFolder
								"folder_identifier": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: `Unique identifier for the folder`,
								},
								"folder_name": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: `Display name for the folder`,
								},
								"items": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{ // androidDeviceOwnerKioskModeFolderItem
										},
									},
									PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: `Items to be added to managed folder. This collection can contain a maximum of 500 elements. / Represents an item that can be added to Android Device Owner folder (application or weblink) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androidDeviceOwnerKioskModeFolderItem?view=graph-rest-beta`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `A list of managed folders for a device in Kiosk Mode. This collection can contain a maximum of 500 elements. / A folder containing pages of apps and weblinks on the Managed Home Screen / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androidDeviceOwnerKioskModeManagedFolder?view=graph-rest-beta`,
					},
					"kiosk_mode_managed_home_screen_auto_signout": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether or not to automatically sign-out of MHS and Shared device mode applications after inactive for Managed Home Screen.`,
					},
					"kiosk_mode_managed_home_screen_inactive_sign_out_delay_in_seconds": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of seconds to give user notice before automatically signing them out for Managed Home Screen. Valid values 0 to 9999999`,
					},
					"kiosk_mode_managed_home_screen_inactive_sign_out_notice_in_seconds": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of seconds device is inactive before automatically signing user out for Managed Home Screen. Valid values 0 to 9999999`,
					},
					"kiosk_mode_managed_home_screen_pin_complexity": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "simple", "complex"),
						},
						MarkdownDescription: `Complexity of PIN for sign-in session for Managed Home Screen.`,
					},
					"kiosk_mode_managed_home_screen_pin_required": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether or not require user to set a PIN for sign-in session for Managed Home Screen.`,
					},
					"kiosk_mode_managed_home_screen_pin_required_to_resume": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether or not required user to enter session PIN if screensaver has appeared for Managed Home Screen.`,
					},
					"kiosk_mode_managed_home_screen_sign_in_background": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Custom URL background for sign-in screen for Managed Home Screen.`,
					},
					"kiosk_mode_managed_home_screen_sign_in_branding_logo": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Custom URL branding logo for sign-in screen and session pin page for Managed Home Screen.`,
					},
					"kiosk_mode_managed_home_screen_sign_in_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether or not show sign-in screen for Managed Home Screen.`,
					},
					"kiosk_mode_managed_settings_entry_disabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether or not to display the Managed Settings entry point on the managed home screen in Kiosk Mode.`,
					},
					"kiosk_mode_media_volume_configuration_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether or not to allow a user to change the media volume in Kiosk Mode.`,
					},
					"kiosk_mode_screen_orientation": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "portrait", "landscape", "autoRotate"),
						},
						MarkdownDescription: `Screen orientation configuration for managed home screen in Kiosk Mode.`,
					},
					"kiosk_mode_screen_saver_configuration_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether or not to enable screen saver mode or not in Kiosk Mode.`,
					},
					"kiosk_mode_screen_saver_detect_media_disabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether or not the device screen should show the screen saver if audio/video is playing in Kiosk Mode.`,
					},
					"kiosk_mode_screen_saver_display_time_in_seconds": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `The number of seconds that the device will display the screen saver for in Kiosk Mode. Valid values 0 to 9999999`,
					},
					"kiosk_mode_screen_saver_image_url": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `URL for an image that will be the device's screen saver in Kiosk Mode.`,
					},
					"kiosk_mode_screen_saver_start_delay_in_seconds": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `The number of seconds the device needs to be inactive for before the screen saver is shown in Kiosk Mode. Valid values 1 to 9999999`,
					},
					"kiosk_mode_show_app_notification_badge": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether or not to display application notification badges in Kiosk Mode.`,
					},
					"kiosk_mode_show_device_info": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether or not to allow a user to access basic device information.`,
					},
					"kiosk_mode_use_managed_home_screen_app": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "singleAppMode", "multiAppMode"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: `Whether or not to use single app kiosk mode or multi-app kiosk mode.`,
					},
					"kiosk_mode_virtual_home_button_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether or not to display a virtual home button when the device is in Kiosk Mode.`,
					},
					"kiosk_mode_virtual_home_button_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "swipeUp", "floating"),
						},
						MarkdownDescription: `Indicates whether the virtual home button is a swipe up home button or a floating home button.`,
					},
					"kiosk_mode_wallpaper_url": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `URL to a publicly accessible image to use for the wallpaper when the device is in Kiosk Mode.`,
					},
					"kiosk_mode_wifi_allowed_ssids": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `The restricted set of WIFI SSIDs available for the user to configure in Kiosk Mode. This collection can contain a maximum of 500 elements.`,
					},
					"kiosk_mode_wifi_configuration_enabled": schema.BoolAttribute{
						Optional:            true,
						Description:         `kioskModeWiFiConfigurationEnabled`, // custom MS Graph attribute name
						MarkdownDescription: `Whether or not to allow a user to configure Wi-Fi settings in Kiosk Mode.`,
					},
					"locate_device_lost_mode_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not LocateDevice for devices with lost mode (COBO, COPE) is enabled.`,
					},
					"locate_device_userless_disabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not LocateDevice for userless (COSU) devices is disabled.`,
					},
					"microphone_force_mute": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to block unmuting the microphone on the device.`,
					},
					"microsoft_launcher_configuration_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to you want configure Microsoft Launcher.`,
					},
					"microsoft_launcher_custom_wallpaper_allow_user_modification": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not the user can modify the wallpaper to personalize their device.`,
					},
					"microsoft_launcher_custom_wallpaper_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to configure the wallpaper on the targeted devices.`,
					},
					"microsoft_launcher_custom_wallpaper_image_url": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates the URL for the image file to use as the wallpaper on the targeted devices.`,
					},
					"microsoft_launcher_dock_presence_allow_user_modification": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not the user can modify the device dock configuration on the device.`,
					},
					"microsoft_launcher_dock_presence_configuration": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "show", "hide", "disabled"),
						},
						MarkdownDescription: `Indicates whether or not you want to configure the device dock.`,
					},
					"microsoft_launcher_feed_allow_user_modification": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not the user can modify the launcher feed on the device.`,
					},
					"microsoft_launcher_feed_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not you want to enable the launcher feed on the device.`,
					},
					"microsoft_launcher_search_bar_placement_configuration": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "top", "bottom", "hide"),
						},
						MarkdownDescription: `Indicates the search bar placement configuration on the device.`,
					},
					"network_escape_hatch_allowed": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not the device will allow connecting to a temporary network connection at boot time.`,
					},
					"nfc_block_outgoing_beam": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to block NFC outgoing beam.`,
					},
					"password_block_keyguard": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not the keyguard is disabled.`,
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
						MarkdownDescription: `List of device keyguard features to block. This collection can contain a maximum of 11 elements.`,
					},
					"password_expiration_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the amount of time that a password can be set for before it expires and a new password will be required. Valid values 1 to 365`,
					},
					"password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the minimum length of the password required on the device. Valid values 4 to 16`,
					},
					"password_minimum_letter_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the minimum number of letter characters required for device password. Valid values 1 to 16`,
					},
					"password_minimum_lower_case_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the minimum number of lower case characters required for device password. Valid values 1 to 16`,
					},
					"password_minimum_non_letter_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the minimum number of non-letter characters required for device password. Valid values 1 to 16`,
					},
					"password_minimum_numeric_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the minimum number of numeric characters required for device password. Valid values 1 to 16`,
					},
					"password_minimum_symbol_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the minimum number of symbol characters required for device password. Valid values 1 to 16`,
					},
					"password_minimum_upper_case_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the minimum number of upper case letter characters required for device password. Valid values 1 to 16`,
					},
					"password_minutes_of_inactivity_before_screen_timeout": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Minutes of inactivity before the screen times out.`,
					},
					"password_previous_password_count_to_block": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the length of password history, where the user will not be able to enter a new password that is the same as any password in the history. Valid values 0 to 24`,
					},
					"password_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "required", "numeric", "numericComplex", "alphabetic", "alphanumeric", "alphanumericWithSymbols", "lowSecurityBiometric", "customPassword"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: `Indicates the minimum password quality required on the device.`,
					},
					"password_require_unlock": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "daily", "unkownFutureValue"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: `Indicates the timeout period after which a device must be unlocked using a form of strong authentication.`,
					},
					"password_sign_in_failure_count_before_factory_reset": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the number of times a user can enter an incorrect password before the device is wiped. Valid values 4 to 11`,
					},
					"personal_profile_apps_allow_install_from_unknown_sources": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether the user can install apps from unknown sources on the personal profile.`,
					},
					"personal_profile_camera_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether to disable the use of the camera on the personal profile.`,
					},
					"personal_profile_personal_applications": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationAppListItemAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Policy applied to applications in the personal profile. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-appListItem?view=graph-rest-beta`,
					},
					"personal_profile_play_store_mode": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "blockedApps", "allowedApps"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: `Used together with PersonalProfilePersonalApplications to control how apps in the personal profile are allowed or blocked.`,
					},
					"personal_profile_screen_capture_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether to disable the capability to take screenshots on the personal profile.`,
					},
					"play_store_mode": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "allowList", "blockList"),
						},
						MarkdownDescription: `Indicates the Play Store mode of the device.`,
					},
					"screen_capture_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to disable the capability to take screenshots.`,
					},
					"security_common_criteria_mode_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Represents the security common criteria mode enabled provided to users when they attempt to modify managed settings on their device.`,
					},
					"security_developer_settings_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not the user is allowed to access developer settings like developer options and safe boot on the device.`,
					},
					"security_require_verify_apps": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not verify apps is required.`,
					},
					"share_device_location_disabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not location sharing is disabled for fully managed devices (COBO), and corporate owned devices with a work profile (COPE)`,
					},
					"short_help_text": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // androidDeviceOwnerUserFacingMessage
							"default_message": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: `The default message displayed if the user's locale doesn't match with any of the localized messages`,
							},
							"localized_messages": schema.SetNestedAttribute{
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: deviceConfigurationKeyValuePairAttributes,
								},
								PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
								Computed:            true,
								MarkdownDescription: `The list of <locale, message> pairs. This collection can contain a maximum of 500 elements. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyValuePair?view=graph-rest-beta`,
							},
						},
						MarkdownDescription: `Represents the customized short help text provided to users when they attempt to modify managed settings on their device. / Represents a user-facing message with locale information as well as a default message to be used if the user's locale doesn't match with any of the localized messages / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androidDeviceOwnerUserFacingMessage?view=graph-rest-beta`,
					},
					"status_bar_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or the status bar is disabled, including notifications, quick settings and other screen overlays.`,
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
						MarkdownDescription: `List of modes in which the device's display will stay powered-on. This collection can contain a maximum of 4 elements.`,
					},
					"storage_allow_usb": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to allow USB mass storage.`,
					},
					"storage_block_external_media": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to block external media.`,
					},
					"storage_block_usb_file_transfer": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to block USB file transfer.`,
					},
					"system_update_freeze_periods": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // androidDeviceOwnerSystemUpdateFreezePeriod
								"end_day": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: `The day of the end date of the freeze period. Valid values 1 to 31`,
								},
								"end_month": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: `The month of the end date of the freeze period. Valid values 1 to 12`,
								},
								"start_day": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: `The day of the start date of the freeze period. Valid values 1 to 31`,
								},
								"start_month": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: `The month of the start date of the freeze period. Valid values 1 to 12`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Indicates the annually repeating time periods during which system updates are postponed. This collection can contain a maximum of 500 elements. / Represents one item in the list of freeze periods for Android Device Owner system updates / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androidDeviceOwnerSystemUpdateFreezePeriod?view=graph-rest-beta`,
					},
					"system_update_install_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "postpone", "windowed", "automatic"),
						},
						MarkdownDescription: `The type of system update configuration.`,
					},
					"system_update_window_end_minutes_after_midnight": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the number of minutes after midnight that the system update window ends. Valid values 0 to 1440`,
					},
					"system_update_window_start_minutes_after_midnight": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the number of minutes after midnight that the system update window starts. Valid values 0 to 1440`,
					},
					"system_windows_blocked": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether or not to block Android system prompt windows, like toasts, phone activities, and system alerts.`,
					},
					"users_block_add": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not adding users and profiles is disabled.`,
					},
					"users_block_remove": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to disable removing other users from the device.`,
					},
					"volume_block_adjustment": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not adjusting the master volume is disabled.`,
					},
					"vpn_always_on_lockdown_mode": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `If an always on VPN package name is specified, whether or not to lock network traffic when that VPN is disconnected.`,
					},
					"vpn_always_on_package_identifier": schema.StringAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
						Computed:            true,
						MarkdownDescription: `Android app package name for app that will handle an always-on VPN connection.`,
					},
					"wifi_block_edit_configurations": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to block the user from editing the wifi connection settings.`,
					},
					"wifi_block_edit_policy_defined_configurations": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to block the user from editing just the networks defined by the policy.`,
					},
					"work_profile_password_expiration_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the number of days that a work profile password can be set before it expires and a new password will be required. Valid values 1 to 365`,
					},
					"work_profile_password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the minimum length of the work profile password. Valid values 4 to 16`,
					},
					"work_profile_password_minimum_letter_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the minimum number of letter characters required for the work profile password. Valid values 1 to 16`,
					},
					"work_profile_password_minimum_lower_case_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the minimum number of lower-case characters required for the work profile password. Valid values 1 to 16`,
					},
					"work_profile_password_minimum_non_letter_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the minimum number of non-letter characters required for the work profile password. Valid values 1 to 16`,
					},
					"work_profile_password_minimum_numeric_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the minimum number of numeric characters required for the work profile password. Valid values 1 to 16`,
					},
					"work_profile_password_minimum_symbol_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the minimum number of symbol characters required for the work profile password. Valid values 1 to 16`,
					},
					"work_profile_password_minimum_upper_case_characters": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the minimum number of upper-case letter characters required for the work profile password. Valid values 1 to 16`,
					},
					"work_profile_password_previous_password_count_to_block": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the length of the work profile password history, where the user will not be able to enter a new password that is the same as any password in the history. Valid values 0 to 24`,
					},
					"work_profile_password_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "required", "numeric", "numericComplex", "alphabetic", "alphanumeric", "alphanumericWithSymbols", "lowSecurityBiometric", "customPassword"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: `Indicates the minimum password quality required on the work profile password.`,
					},
					"work_profile_password_require_unlock": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "daily", "unkownFutureValue"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: `Indicates the timeout period after which a work profile must be unlocked using a form of strong authentication.`,
					},
					"work_profile_password_sign_in_failure_count_before_factory_reset": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Indicates the number of times a user can enter an incorrect work profile password before the device is wiped. Valid values 4 to 11`,
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: `This topic provides descriptions of the declared methods, properties and relationships exposed by the androidDeviceOwnerGeneralDeviceConfiguration resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androidDeviceOwnerGeneralDeviceConfiguration?view=graph-rest-beta`,
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
						MarkdownDescription: `Payload. (UTF8 encoded byte array)`,
					},
					"file_name": schema.StringAttribute{
						Required:            true,
						Description:         `payloadFileName`, // custom MS Graph attribute name
						MarkdownDescription: `Payload file name (*.mobileconfig | *.xml).`,
					},
					"name": schema.StringAttribute{
						Required:            true,
						Description:         `payloadName`, // custom MS Graph attribute name
						MarkdownDescription: `Name that is displayed to the user.`,
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: `This topic provides descriptions of the declared methods, properties and relationships exposed by the iosCustomConfiguration resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosCustomConfiguration?view=graph-rest-beta`,
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
						MarkdownDescription: `An array of AirPrint printers that should always be shown. This collection can contain a maximum of 500 elements. / Represents an AirPrint destination. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-airPrintDestination?view=graph-rest-beta`,
					},
					"asset_tag_template": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Asset tag information for the device, displayed on the login window and lock screen.`,
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
											MarkdownDescription: `Additional URLs allowed for access`,
										},
										"blocked_urls": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: `Additional URLs blocked for access`,
										},
									},
									Validators: []validator.Object{
										deviceConfigurationIosWebContentFilterBaseValidator,
									},
									MarkdownDescription: `Represents an iOS Web Content Filter setting type, which enables iOS automatic filter feature and allows for additional URL access control. When constructed with no property values, the iOS device will enable the automatic filter regardless. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosWebContentFilterAutoFilter?view=graph-rest-beta`,
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
											MarkdownDescription: `URL bookmarks which will be installed into built-in browser and user is only allowed to access websites through bookmarks. This collection can contain a maximum of 500 elements. / iOS URL bookmark / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosBookmark?view=graph-rest-beta`,
										},
										"website_list": schema.SetNestedAttribute{
											Optional: true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: deviceConfigurationIosBookmarkAttributes,
											},
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: `URL bookmarks which will be installed into built-in browser and user is only allowed to access websites through bookmarks. This collection can contain a maximum of 500 elements. / iOS URL bookmark / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosBookmark?view=graph-rest-beta`,
										},
									},
									Validators: []validator.Object{
										deviceConfigurationIosWebContentFilterBaseValidator,
									},
									MarkdownDescription: `Represents an iOS Web Content Filter setting type, which installs URL bookmarks into iOS built-in browser. An example scenario is in the classroom where teachers would like the students to navigate websites through browser bookmarks configured on their iOS devices, and no access to other sites. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosWebContentFilterSpecificWebsitesAccess?view=graph-rest-beta`,
								},
							},
						},
						MarkdownDescription: `Gets or sets iOS Web Content Filter settings, supervised mode only / Represents an iOS Web Content Filter setting base type. An empty and abstract base. Caller should use one of derived types for configurations. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosWebContentFilterBase?view=graph-rest-beta`,
					},
					"home_screen_dock_icons": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // iosHomeScreenItem
								"display_name": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `Name of the app`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `A list of app and folders to appear on the Home Screen Dock. This collection can contain a maximum of 500 elements. / Represents an item on the iOS Home Screen / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosHomeScreenItem?view=graph-rest-beta`,
					},
					"home_screen_grid_height": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Gets or sets the number of rows to render when configuring iOS home screen layout settings. If this value is configured, homeScreenGridWidth must be configured as well.`,
					},
					"home_screen_grid_width": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Gets or sets the number of columns to render when configuring iOS home screen layout settings. If this value is configured, homeScreenGridHeight must be configured as well.`,
					},
					"home_screen_pages": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // iosHomeScreenPage
								"display_name": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `Name of the page`,
								},
								"icons": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{ // iosHomeScreenItem
											"display_name": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: `Name of the app`,
											},
										},
									},
									PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: `A list of apps, folders, and web clips to appear on a page. This collection can contain a maximum of 500 elements. / Represents an item on the iOS Home Screen / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosHomeScreenItem?view=graph-rest-beta`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `A list of pages on the Home Screen. This collection can contain a maximum of 500 elements. / A page containing apps, folders, and web clips on the Home Screen. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosHomeScreenPage?view=graph-rest-beta`,
					},
					"ios_single_sign_on_extension": schema.SingleNestedAttribute{
						Optional:   true,
						Attributes: map[string]schema.Attribute{ // iosSingleSignOnExtension
						},
						MarkdownDescription: `Gets or sets a single sign-on extension profile. / An abstract base class for all iOS-specific single sign-on extension types. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosSingleSignOnExtension?view=graph-rest-beta`,
					},
					"lock_screen_footnote": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `A footnote displayed on the login window and lock screen. Available in iOS 9.3.1 and later.`,
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
									MarkdownDescription: `Indicates the type of alert for notifications for this app.`,
								},
								"app_name": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `Application name to be associated with the bundleID.`,
								},
								"badges_enabled": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: `Indicates whether badges are allowed for this app.`,
								},
								"bundle_id": schema.StringAttribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
									Computed:            true,
									Description:         `bundleID`, // custom MS Graph attribute name
									MarkdownDescription: `Bundle id of app to which to apply these notification settings.`,
								},
								"enabled": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: `Indicates whether notifications are allowed for this app.`,
								},
								"preview_visibility": schema.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.OneOf("notConfigured", "alwaysShow", "hideWhenLocked", "neverShow"),
									},
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
									Computed:            true,
									MarkdownDescription: `Overrides the notification preview policy set by the user on an iOS device.`,
								},
								"publisher": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `Publisher to be associated with the bundleID.`,
								},
								"show_in_notification_center": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: `Indicates whether notifications can be shown in notification center.`,
								},
								"show_on_lock_screen": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: `Indicates whether notifications can be shown on the lock screen.`,
								},
								"sounds_enabled": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: `Indicates whether sounds are allowed for this app.`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Notification settings for each bundle id. Applicable to devices in supervised mode only (iOS 9.3 and later). This collection can contain a maximum of 500 elements. / An item describing notification setting. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosNotificationSettings?view=graph-rest-beta`,
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
								MarkdownDescription: `List of app identifiers that are allowed to use this login. If this field is omitted, the login applies to all applications on the device. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-appListItem?view=graph-rest-beta`,
							},
							"allowed_urls": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
								Computed:            true,
								MarkdownDescription: `List of HTTP URLs that must be matched in order to use this login. With iOS 9.0 or later, a wildcard characters may be used.`,
							},
							"display_name": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: `The display name of login settings shown on the receiving device.`,
							},
							"kerberos_principal_name": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: `A Kerberos principal name. If not provided, the user is prompted for one during profile installation.`,
							},
							"kerberos_realm": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: `A Kerberos realm name. Case sensitive.`,
							},
						},
						MarkdownDescription: `The Kerberos login settings that enable apps on receiving devices to authenticate smoothly. / iOS Kerberos authentication settings for single sign-on / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosSingleSignOnSettings?view=graph-rest-beta`,
					},
					"wallpaper_display_location": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "lockScreen", "homeScreen", "lockAndHomeScreens"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: `A wallpaper display location specifier.`,
					},
					"wallpaper_image": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mimeContent
							"type": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: `Indicates the content mime type.`,
							},
							"value": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: `The byte array that contains the actual content.`,
							},
						},
						MarkdownDescription: `A wallpaper image must be in either PNG or JPEG format. It requires a supervised device with iOS 8 or later version. / Contains properties for a generic mime content. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mimeContent?view=graph-rest-beta`,
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: `iOS Device Features Configuration Profile. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosDeviceFeaturesConfiguration?view=graph-rest-beta`,
			},
		},
		"ios_eas_email_profile": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.iosEasEmailProfileConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // iosEasEmailProfileConfiguration
					"custom_domain_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Custom domain name value used while generating an email profile before installing on the device.`,
					},
					"user_domain_name_source": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("fullDomainName", "netBiosDomainName"),
						},
						MarkdownDescription: `UserDomainname attribute that is picked from AAD and injected into this profile before installing on the device.`,
					},
					"username_aad_source": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userPrincipalName", "primarySmtpAddress", "samAccountName"),
						},
						Description:         `usernameAADSource`, // custom MS Graph attribute name
						MarkdownDescription: `Name of the AAD field, that will be used to retrieve UserName for email profile.`,
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
						MarkdownDescription: `Username attribute that is picked from AAD and injected into this profile before installing on the device.`,
					},
					"account_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: `Account name.`,
					},
					"authentication_method": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("usernameAndPassword", "certificate", "derivedCredential"),
						},
						MarkdownDescription: `Authentication method for this Email profile.`,
					},
					"block_moving_messages_to_other_email_accounts": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to block moving messages to other email accounts.`,
					},
					"block_sending_email_from_third_party_apps": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to block sending email from third party apps.`,
					},
					"block_syncing_recently_used_email_addresses": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether or not to block syncing recently used email addresses, for instance - when composing new email.`,
					},
					"duration_of_email_to_sync": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userDefined", "oneDay", "threeDays", "oneWeek", "twoWeeks", "oneMonth", "unlimited"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("userDefined")},
						Computed:            true,
						MarkdownDescription: `Duration of time email should be synced back to. `,
					},
					"eas_services": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							wpvalidator.FlagEnumValues("none", "calendars", "contacts", "email", "notes", "reminders"),
						},
						MarkdownDescription: `Exchange data to sync.`,
					},
					"eas_services_user_override_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Allow users to change sync settings.`,
					},
					"email_address_source": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userPrincipalName", "primarySmtpAddress"),
						},
						MarkdownDescription: `Email attribute that is picked from AAD and injected into this profile before installing on the device.`,
					},
					"encryption_certificate_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("none", "certificate", "derivedCredential"),
						},
						MarkdownDescription: `Encryption Certificate type for this Email profile.`,
					},
					"host_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: `Exchange location that (URL) that the native mail app connects to.`,
					},
					"per_app_vpn_profile_id": schema.StringAttribute{
						Optional:            true,
						Description:         `perAppVPNProfileId`, // custom MS Graph attribute name
						MarkdownDescription: `Profile ID of the Per-App VPN policy to be used to access emails from the native Mail client`,
					},
					"require_ssl": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to use SSL.`,
					},
					"signing_certificate_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("none", "certificate", "derivedCredential"),
						},
						MarkdownDescription: `Signing Certificate type for this Email profile.`,
					},
					"use_oauth": schema.BoolAttribute{
						Optional:            true,
						Description:         `useOAuth`, // custom MS Graph attribute name
						MarkdownDescription: `Specifies whether the connection should use OAuth for authentication.`,
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: `By providing configurations in this profile you can instruct the native email client on iOS devices to communicate with an Exchange server and get email, contacts, calendar, reminders, and notes. Furthermore, you can also specify how much email to sync and how often the device should sync. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosEasEmailProfileConfiguration?view=graph-rest-beta`,
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
						MarkdownDescription: `Indicates whether or not to allow account modification when the device is in supervised mode.`,
					},
					"activation_lock_allow_when_supervised": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow activation lock when the device is in the supervised mode.`,
					},
					"airdrop_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `airDropBlocked`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to allow AirDrop when the device is in supervised mode.`,
					},
					"airdrop_force_unmanaged_drop_target": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `airDropForceUnmanagedDropTarget`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to cause AirDrop to be considered an unmanaged drop target (iOS 9.0 and later).`,
					},
					"airplay_force_pairing_password_for_outgoing_requests": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `airPlayForcePairingPasswordForOutgoingRequests`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to enforce all devices receiving AirPlay requests from this device to use a pairing password.`,
					},
					"airprint_block_credentials_storage": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `airPrintBlockCredentialsStorage`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not keychain storage of username and password for Airprint is blocked (iOS 11.0 and later).`,
					},
					"airprint_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `airPrintBlocked`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not AirPrint is blocked (iOS 11.0 and later).`,
					},
					"airprint_blocki_beacon_discovery": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `airPrintBlockiBeaconDiscovery`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not iBeacon discovery of AirPrint printers is blocked. This prevents spurious AirPrint Bluetooth beacons from phishing for network traffic (iOS 11.0 and later).`,
					},
					"airprint_force_trusted_tls": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `airPrintForceTrustedTLS`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates if trusted certificates are required for TLS printing communication (iOS 11.0 and later).`,
					},
					"app_clips_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Prevents a user from adding any App Clips and removes any existing App Clips on the device.`,
					},
					"apple_news_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from using News when the device is in supervised mode (iOS 9.0 and later).`,
					},
					"apple_personalized_ads_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Limits Apple personalized advertising when true. Available in iOS 14 and later.`,
					},
					"apple_watch_block_pairing": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow Apple Watch pairing when the device is in supervised mode (iOS 9.0 and later).`,
					},
					"apple_watch_force_wrist_detection": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to force a paired Apple Watch to use Wrist Detection (iOS 8.2 and later).`,
					},
					"app_removal_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates if the removal of apps is allowed.`,
					},
					"apps_single_app_mode_list": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationAppListItemAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Gets or sets the list of iOS apps allowed to autonomously enter Single App Mode. Supervised only. iOS 7.0 and later. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-appListItem?view=graph-rest-beta`,
					},
					"app_store_block_automatic_downloads": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the automatic downloading of apps purchased on other devices when the device is in supervised mode (iOS 9.0 and later).`,
					},
					"app_store_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from using the App Store. Requires a supervised device for iOS 13 and later.`,
					},
					"app_store_block_in_app_purchases": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from making in app purchases.`,
					},
					"app_store_block_ui_app_installation": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `appStoreBlockUIAppInstallation`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to block the App Store app, not restricting installation through Host apps. Applies to supervised mode only (iOS 9.0 and later).`,
					},
					"app_store_require_password": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to require a password when using the app store.`,
					},
					"apps_visibility_list": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationAppListItemAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `List of apps in the visibility list (either visible/launchable apps list or hidden/unlaunchable apps list, controlled by AppsVisibilityListType) (iOS 9.3 and later). This collection can contain a maximum of 10000 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-appListItem?view=graph-rest-beta`,
					},
					"apps_visibility_list_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("none", "appsInListCompliant", "appsNotInListCompliant"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
						Computed:            true,
						MarkdownDescription: `Type of list that is in the AppsVisibilityList.`,
					},
					"auto_fill_force_authentication": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to force user authentication before autofilling passwords and credit card information in Safari and other apps on supervised devices.`,
					},
					"auto_unlock_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Blocks users from unlocking their device with Apple Watch. Available for devices running iOS and iPadOS versions 14.5 and later.`,
					},
					"block_system_app_removal": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not the removal of system apps from the device is blocked on a supervised device (iOS 11.0 and later).`,
					},
					"bluetooth_block_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow modification of Bluetooth settings when the device is in supervised mode (iOS 10.0 and later).`,
					},
					"camera_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from accessing the camera of the device. Requires a supervised device for iOS 13 and later.`,
					},
					"cellular_block_data_roaming": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block data roaming.`,
					},
					"cellular_block_global_background_fetch_while_roaming": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block global background fetch while roaming.`,
					},
					"cellular_block_per_app_data_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow changes to cellular app data usage settings when the device is in supervised mode.`,
					},
					"cellular_block_personal_hotspot": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block Personal Hotspot.`,
					},
					"cellular_block_personal_hotspot_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from modifying the personal hotspot setting (iOS 12.2 or later).`,
					},
					"cellular_block_plan_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow users to change the settings of the cellular plan on a supervised device.`,
					},
					"cellular_block_voice_roaming": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block voice roaming.`,
					},
					"certificates_block_untrusted_tls_certificates": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block untrusted TLS certificates.`,
					},
					"classroom_app_block_remote_screen_observation": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow remote screen observation by Classroom app when the device is in supervised mode (iOS 9.3 and later).`,
					},
					"classroom_app_force_unprompted_screen_observation": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to automatically give permission to the teacher of a managed course on the Classroom app to view a student's screen without prompting when the device is in supervised mode.`,
					},
					"classroom_force_automatically_join_classes": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to automatically give permission to the teacher's requests, without prompting the student, when the device is in supervised mode.`,
					},
					"classroom_force_request_permission_to_leave_classes": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether a student enrolled in an unmanaged course via Classroom will request permission from the teacher when attempting to leave the course (iOS 11.3 and later).`,
					},
					"classroom_force_unprompted_app_and_device_lock": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow the teacher to lock apps or the device without prompting the student. Supervised only.`,
					},
					"compliant_app_list_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("none", "appsInListCompliant", "appsNotInListCompliant"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
						Computed:            true,
						MarkdownDescription: `List that is in the AppComplianceList.`,
					},
					"compliant_apps_list": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationAppListItemAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `List of apps in the compliance (either allow list or block list, controlled by CompliantAppListType). This collection can contain a maximum of 10000 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-appListItem?view=graph-rest-beta`,
					},
					"configuration_profile_block_changes": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from installing configuration profiles and certificates interactively when the device is in supervised mode.`,
					},
					"contacts_allow_managed_to_unmanaged_write": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not managed apps can write contacts to unmanaged contacts accounts (iOS 12.0 and later).`,
					},
					"contacts_allow_unmanaged_to_managed_read": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not unmanaged apps can read from managed contacts accounts (iOS 12.0 or later).`,
					},
					"continuous_path_keyboard_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the continuous path keyboard when the device is supervised (iOS 13 or later).`,
					},
					"date_and_time_force_set_automatically": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not the Date and Time "Set Automatically" feature is enabled and cannot be turned off by the user (iOS 12.0 and later).`,
					},
					"definition_lookup_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block definition lookup when the device is in supervised mode (iOS 8.1.3 and later ).`,
					},
					"device_block_enable_restrictions": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow the user to enables restrictions in the device settings when the device is in supervised mode.`,
					},
					"device_block_erase_content_and_settings": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow the use of the 'Erase all content and settings' option on the device when the device is in supervised mode.`,
					},
					"device_block_name_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow device name modification when the device is in supervised mode (iOS 9.0 and later).`,
					},
					"diagnostic_data_block_submission": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block diagnostic data submission.`,
					},
					"diagnostic_data_block_submission_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow diagnostics submission settings modification when the device is in supervised mode (iOS 9.3.2 and later).`,
					},
					"documents_block_managed_documents_in_unmanaged_apps": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from viewing managed documents in unmanaged apps.`,
					},
					"documents_block_unmanaged_documents_in_managed_apps": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from viewing unmanaged documents in managed apps.`,
					},
					"email_in_domain_suffixes": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `An email address lacking a suffix that matches any of these strings will be considered out-of-domain.`,
					},
					"enterprise_app_block_trust": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from trusting an enterprise app.`,
					},
					"enterprise_book_block_backup": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not Enterprise book back up is blocked.`,
					},
					"enterprise_book_block_metadata_sync": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not Enterprise book notes and highlights sync is blocked.`,
					},
					"esim_block_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow the addition or removal of cellular plans on the eSIM of a supervised device.`,
					},
					"facetime_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `faceTimeBlocked`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to block the user from using FaceTime. Requires a supervised device for iOS 13 and later.`,
					},
					"files_network_drive_access_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates if devices can access files or other resources on a network server using the Server Message Block (SMB) protocol. Available for devices running iOS and iPadOS, versions 13.0 and later.`,
					},
					"files_usb_drive_access_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates if sevices with access can connect to and open files on a USB drive. Available for devices running iOS and iPadOS, versions 13.0 and later.`,
					},
					"find_my_device_in_find_my_app_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block Find My Device when the device is supervised (iOS 13 or later).`,
					},
					"find_my_friends_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block changes to Find My Friends when the device is in supervised mode.`,
					},
					"find_my_friends_in_find_my_app_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block Find My Friends when the device is supervised (iOS 13 or later).`,
					},
					"game_center_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from using Game Center when the device is in supervised mode.`,
					},
					"gaming_block_game_center_friends": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from having friends in Game Center. Requires a supervised device for iOS 13 and later.`,
					},
					"gaming_block_multiplayer": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from using multiplayer gaming. Requires a supervised device for iOS 13 and later.`,
					},
					"host_pairing_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `indicates whether or not to allow host pairing to control the devices an iOS device can pair with when the iOS device is in supervised mode.`,
					},
					"ibooks_store_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iBooksStoreBlocked`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to block the user from using the iBooks Store when the device is in supervised mode.`,
					},
					"ibooks_store_block_erotica": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iBooksStoreBlockErotica`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to block the user from downloading media from the iBookstore that has been tagged as erotica.`,
					},
					"icloud_block_activity_continuation": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iCloudBlockActivityContinuation`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to block the user from continuing work they started on iOS device to another iOS or macOS device.`,
					},
					"icloud_block_backup": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iCloudBlockBackup`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to block iCloud backup. Requires a supervised device for iOS 13 and later.`,
					},
					"icloud_block_document_sync": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iCloudBlockDocumentSync`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to block iCloud document sync. Requires a supervised device for iOS 13 and later.`,
					},
					"icloud_block_managed_apps_sync": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iCloudBlockManagedAppsSync`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to block Managed Apps Cloud Sync.`,
					},
					"icloud_block_photo_library": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iCloudBlockPhotoLibrary`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to block iCloud Photo Library.`,
					},
					"icloud_block_photo_stream_sync": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iCloudBlockPhotoStreamSync`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to block iCloud Photo Stream Sync.`,
					},
					"icloud_block_shared_photo_stream": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iCloudBlockSharedPhotoStream`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to block Shared Photo Stream.`,
					},
					"icloud_private_relay_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iCloudPrivateRelayBlocked`, // custom MS Graph attribute name
						MarkdownDescription: `iCloud private relay is an iCloud+ service that prevents networks and servers from monitoring a person's activity across the internet. By blocking iCloud private relay, Apple will not encrypt the traffic leaving the device. Available for devices running iOS 15 and later.`,
					},
					"icloud_require_encrypted_backup": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iCloudRequireEncryptedBackup`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to require backups to iCloud be encrypted.`,
					},
					"itunes_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iTunesBlocked`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to block the iTunes app. Requires a supervised device for iOS 13 and later.`,
					},
					"itunes_block_explicit_content": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iTunesBlockExplicitContent`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to block the user from accessing explicit content in iTunes and the App Store. Requires a supervised device for iOS 13 and later.`,
					},
					"itunes_block_music_service": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iTunesBlockMusicService`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to block Music service and revert Music app to classic mode when the device is in supervised mode (iOS 9.3 and later and macOS 10.12 and later).`,
					},
					"itunes_block_radio": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `iTunesBlockRadio`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to block the user from using iTunes Radio when the device is in supervised mode (iOS 9.3 and later).`,
					},
					"keyboard_block_auto_correct": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block keyboard auto-correction when the device is in supervised mode (iOS 8.1.3 and later).`,
					},
					"keyboard_block_dictation": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from using dictation input when the device is in supervised mode.`,
					},
					"keyboard_block_predictive": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block predictive keyboards when device is in supervised mode (iOS 8.1.3 and later).`,
					},
					"keyboard_block_shortcuts": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block keyboard shortcuts when the device is in supervised mode (iOS 9.0 and later).`,
					},
					"keyboard_block_spell_check": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block keyboard spell-checking when the device is in supervised mode (iOS 8.1.3 and later).`,
					},
					"keychain_block_cloud_sync": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not iCloud keychain synchronization is blocked. Requires a supervised device for iOS 13 and later.`,
					},
					"kiosk_mode_allow_assistive_speak": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow assistive speak while in kiosk mode.`,
					},
					"kiosk_mode_allow_assistive_touch_settings": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow access to the Assistive Touch Settings while in kiosk mode.`,
					},
					"kiosk_mode_allow_color_inversion_settings": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow access to the Color Inversion Settings while in kiosk mode.`,
					},
					"kiosk_mode_allow_voice_control_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow the user to toggle voice control in kiosk mode.`,
					},
					"kiosk_mode_allow_voice_over_settings": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow access to the voice over settings while in kiosk mode.`,
					},
					"kiosk_mode_allow_zoom_settings": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow access to the zoom settings while in kiosk mode.`,
					},
					"kiosk_mode_app_store_url": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `URL in the app store to the app to use for kiosk mode. Use if KioskModeManagedAppId is not known.`,
					},
					"kiosk_mode_app_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "appStoreApp", "managedApp", "builtInApp"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: `Type of app to run in kiosk mode.`,
					},
					"kiosk_mode_block_auto_lock": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block device auto lock while in kiosk mode.`,
					},
					"kiosk_mode_block_ringer_switch": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block use of the ringer switch while in kiosk mode.`,
					},
					"kiosk_mode_block_screen_rotation": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block screen rotation while in kiosk mode.`,
					},
					"kiosk_mode_block_sleep_button": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block use of the sleep button while in kiosk mode.`,
					},
					"kiosk_mode_block_touchscreen": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block use of the touchscreen while in kiosk mode.`,
					},
					"kiosk_mode_block_volume_buttons": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the volume buttons while in Kiosk Mode.`,
					},
					"kiosk_mode_built_in_app_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `ID for built-in apps to use for kiosk mode. Used when KioskModeManagedAppId and KioskModeAppStoreUrl are not set.`,
					},
					"kiosk_mode_enable_voice_control": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to enable voice control in kiosk mode.`,
					},
					"kiosk_mode_managed_app_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Managed app id of the app to use for kiosk mode. If KioskModeManagedAppId is specified then KioskModeAppStoreUrl will be ignored.`,
					},
					"kiosk_mode_require_assistive_touch": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to require assistive touch while in kiosk mode.`,
					},
					"kiosk_mode_require_color_inversion": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to require color inversion while in kiosk mode.`,
					},
					"kiosk_mode_require_mono_audio": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to require mono audio while in kiosk mode.`,
					},
					"kiosk_mode_require_voice_over": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to require voice over while in kiosk mode.`,
					},
					"kiosk_mode_require_zoom": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to require zoom while in kiosk mode.`,
					},
					"lock_screen_block_control_center": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from using control center on the lock screen.`,
					},
					"lock_screen_block_notification_view": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from using the notification view on the lock screen.`,
					},
					"lock_screen_block_passbook": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from using passbook when the device is locked.`,
					},
					"lock_screen_block_today_view": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from using the Today View on the lock screen.`,
					},
					"managed_pasteboard_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Open-in management controls how people share data between unmanaged and managed apps. Setting this to true enforces copy/paste restrictions based on how you configured <b>Block viewing corporate documents in unmanaged apps </b> and <b> Block viewing non-corporate documents in corporate apps.</b>`,
					},
					"media_content_rating_apps": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("allAllowed", "allBlocked", "agesAbove4", "agesAbove9", "agesAbove12", "agesAbove17"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("allAllowed")},
						Computed:            true,
						MarkdownDescription: `Media content rating settings for Apps`,
					},
					"media_content_rating_australia": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mediaContentRatingAustralia
							"movie_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "parentalGuidance", "mature", "agesAbove15", "agesAbove18"),
								},
								MarkdownDescription: `Movies rating selected for Australia`,
							},
							"tv_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "preschoolers", "children", "general", "parentalGuidance", "mature", "agesAbove15", "agesAbove15AdultViolence"),
								},
								MarkdownDescription: `TV rating selected for Australia`,
							},
						},
						MarkdownDescription: `Media content rating settings for Australia / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediaContentRatingAustralia?view=graph-rest-beta`,
					},
					"media_content_rating_canada": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mediaContentRatingCanada
							"movie_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "parentalGuidance", "agesAbove14", "agesAbove18", "restricted"),
								},
								MarkdownDescription: `Movies rating selected for Canada`,
							},
							"tv_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "children", "childrenAbove8", "general", "parentalGuidance", "agesAbove14", "agesAbove18"),
								},
								MarkdownDescription: `TV rating selected for Canada`,
							},
						},
						MarkdownDescription: `Media content rating settings for Canada / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediaContentRatingCanada?view=graph-rest-beta`,
					},
					"media_content_rating_france": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mediaContentRatingFrance
							"movie_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "agesAbove10", "agesAbove12", "agesAbove16", "agesAbove18"),
								},
								MarkdownDescription: `Movies rating selected for France`,
							},
							"tv_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "agesAbove10", "agesAbove12", "agesAbove16", "agesAbove18"),
								},
								MarkdownDescription: `TV rating selected for France`,
							},
						},
						MarkdownDescription: `Media content rating settings for France / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediaContentRatingFrance?view=graph-rest-beta`,
					},
					"media_content_rating_germany": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mediaContentRatingGermany
							"movie_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "agesAbove6", "agesAbove12", "agesAbove16", "adults"),
								},
								MarkdownDescription: `Movies rating selected for Germany`,
							},
							"tv_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "agesAbove6", "agesAbove12", "agesAbove16", "adults"),
								},
								MarkdownDescription: `TV rating selected for Germany`,
							},
						},
						MarkdownDescription: `Media content rating settings for Germany / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediaContentRatingGermany?view=graph-rest-beta`,
					},
					"media_content_rating_ireland": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mediaContentRatingIreland
							"movie_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "parentalGuidance", "agesAbove12", "agesAbove15", "agesAbove16", "adults"),
								},
								MarkdownDescription: `Movies rating selected for Ireland`,
							},
							"tv_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "children", "youngAdults", "parentalSupervision", "mature"),
								},
								MarkdownDescription: `TV rating selected for Ireland`,
							},
						},
						MarkdownDescription: `Media content rating settings for Ireland / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediaContentRatingIreland?view=graph-rest-beta`,
					},
					"media_content_rating_japan": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mediaContentRatingJapan
							"movie_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "parentalGuidance", "agesAbove15", "agesAbove18"),
								},
								MarkdownDescription: `Movies rating selected for Japan`,
							},
							"tv_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "explicitAllowed"),
								},
								MarkdownDescription: `TV rating selected for Japan`,
							},
						},
						MarkdownDescription: `Media content rating settings for Japan / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediaContentRatingJapan?view=graph-rest-beta`,
					},
					"media_content_rating_new_zealand": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mediaContentRatingNewZealand
							"movie_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "parentalGuidance", "mature", "agesAbove13", "agesAbove15", "agesAbove16", "agesAbove18", "restricted", "agesAbove16Restricted"),
								},
								MarkdownDescription: `Movies rating selected for New Zealand`,
							},
							"tv_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "parentalGuidance", "adults"),
								},
								MarkdownDescription: `TV rating selected for New Zealand`,
							},
						},
						MarkdownDescription: `Media content rating settings for New Zealand / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediaContentRatingNewZealand?view=graph-rest-beta`,
					},
					"media_content_rating_united_kingdom": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mediaContentRatingUnitedKingdom
							"movie_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "universalChildren", "parentalGuidance", "agesAbove12Video", "agesAbove12Cinema", "agesAbove15", "adults"),
								},
								MarkdownDescription: `Movies rating selected for United Kingdom`,
							},
							"tv_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "caution"),
								},
								MarkdownDescription: `TV rating selected for United Kingdom`,
							},
						},
						MarkdownDescription: `Media content rating settings for United Kingdom / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediaContentRatingUnitedKingdom?view=graph-rest-beta`,
					},
					"media_content_rating_united_states": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // mediaContentRatingUnitedStates
							"movie_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "general", "parentalGuidance", "parentalGuidance13", "restricted", "adults"),
								},
								MarkdownDescription: `Movies rating selected for United States`,
							},
							"tv_rating": schema.StringAttribute{
								Required: true,
								Validators: []validator.String{
									stringvalidator.OneOf("allAllowed", "allBlocked", "childrenAll", "childrenAbove7", "general", "parentalGuidance", "childrenAbove14", "adults"),
								},
								MarkdownDescription: `TV rating selected for United States`,
							},
						},
						MarkdownDescription: `Media content rating settings for United States / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediaContentRatingUnitedStates?view=graph-rest-beta`,
					},
					"messages_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from using the Messages app on the supervised device.`,
					},
					"network_usage_rules": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // iosNetworkUsageRule
								"cellular_data_blocked": schema.BoolAttribute{
									Required:            true,
									MarkdownDescription: `If set to true, corresponding managed apps will not be allowed to use cellular data at any time.`,
								},
								"cellular_data_block_when_roaming": schema.BoolAttribute{
									Required:            true,
									MarkdownDescription: `If set to true, corresponding managed apps will not be allowed to use cellular data when roaming.`,
								},
								"managed_apps": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: deviceConfigurationAppListItemAttributes,
									},
									MarkdownDescription: `Information about the managed apps that this rule is going to apply to. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-appListItem?view=graph-rest-beta`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `List of managed apps and the network rules that applies to them. This collection can contain a maximum of 1000 elements. / Network Usage Rules allow enterprises to specify how managed apps use networks, such as cellular data networks. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosNetworkUsageRule?view=graph-rest-beta`,
					},
					"nfc_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Disable NFC to prevent devices from pairing with other NFC-enabled devices. Available for iOS/iPadOS devices running 14.2 and later.`,
					},
					"notifications_block_settings_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow notifications settings modification (iOS 9.3 and later).`,
					},
					"on_device_only_dictation_forced": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Disables connections to Siri servers so that users can’t use Siri to dictate text. Available for devices running iOS and iPadOS versions 14.5 and later.`,
					},
					"on_device_only_translation_forced": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `When set to TRUE, the setting disables connections to Siri servers so that users can’t use Siri to translate text. When set to FALSE, the setting allows connections to to Siri servers to users can use Siri to translate text. Available for devices running iOS and iPadOS versions 15.0 and later.`,
					},
					"passcode_block_fingerprint_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Block modification of registered Touch ID fingerprints when in supervised mode.`,
					},
					"passcode_block_fingerprint_unlock": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block fingerprint unlock.`,
					},
					"passcode_block_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow passcode modification on the supervised device (iOS 9.0 and later).`,
					},
					"passcode_block_simple": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block simple passcodes.`,
					},
					"passcode_expiration_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of days before the passcode expires. Valid values 1 to 65535`,
					},
					"passcode_minimum_character_set_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of character sets a passcode must contain. Valid values 0 to 4`,
					},
					"passcode_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Minimum length of passcode. Valid values 4 to 14`,
					},
					"passcode_minutes_of_inactivity_before_lock": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Minutes of inactivity before a passcode is required.`,
					},
					"passcode_minutes_of_inactivity_before_screen_timeout": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Minutes of inactivity before the screen times out.`,
					},
					"passcode_previous_passcode_block_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of previous passcodes to block. Valid values 1 to 24`,
					},
					"passcode_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to require a passcode.`,
					},
					"passcode_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "alphanumeric", "numeric"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: `Type of passcode that is required.`,
					},
					"passcode_sign_in_failure_count_before_wipe": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of sign in failures allowed before wiping the device. Valid values 2 to 11`,
					},
					"password_block_airdrop_sharing": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `passwordBlockAirDropSharing`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to block sharing passwords with the AirDrop passwords feature iOS 12.0 and later).`,
					},
					"password_block_auto_fill": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates if the AutoFill passwords feature is allowed (iOS 12.0 and later).`,
					},
					"password_block_proximity_requests": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block requesting passwords from nearby devices (iOS 12.0 and later).`,
					},
					"pki_block_ota_updates": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `pkiBlockOTAUpdates`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not over-the-air PKI updates are blocked. Setting this restriction to false does not disable CRL and OCSP checks (iOS 7.0 and later).`,
					},
					"podcasts_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from using podcasts on the supervised device (iOS 8.0 and later).`,
					},
					"privacy_force_limit_ad_tracking": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates if ad tracking is limited.(iOS 7.0 and later).`,
					},
					"proximity_block_setup_to_new_device": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to enable the prompt to setup nearby devices with a supervised device.`,
					},
					"safari_block_autofill": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from using Auto fill in Safari. Requires a supervised device for iOS 13 and later.`,
					},
					"safari_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from using Safari. Requires a supervised device for iOS 13 and later.`,
					},
					"safari_block_javascript": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `safariBlockJavaScript`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to block JavaScript in Safari.`,
					},
					"safari_block_popups": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block popups in Safari.`,
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
						MarkdownDescription: `Cookie settings for Safari.`,
					},
					"safari_managed_domains": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `URLs matching the patterns listed here will be considered managed.`,
					},
					"safari_password_auto_fill_domains": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Users can save passwords in Safari only from URLs matching the patterns listed here. Applies to devices in supervised mode (iOS 9.3 and later).`,
					},
					"safari_require_fraud_warning": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to require fraud warning in Safari.`,
					},
					"screen_capture_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from taking Screenshots.`,
					},
					"shared_device_block_temporary_sessions": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block temporary sessions on Shared iPads (iOS 13.4 or later).`,
					},
					"siri_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from using Siri.`,
					},
					"siri_blocked_when_locked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block the user from using Siri when locked.`,
					},
					"siri_block_user_generated_content": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block Siri from querying user-generated content when used on a supervised device.`,
					},
					"siri_require_profanity_filter": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to prevent Siri from dictating, or speaking profane language on supervised device.`,
					},
					"software_updates_enforced_delay_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Sets how many days a software update will be delyed for a supervised device. Valid values 0 to 90`,
					},
					"software_updates_force_delayed": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to delay user visibility of software updates when the device is in supervised mode.`,
					},
					"spotlight_block_internet_results": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block Spotlight search from returning internet results on supervised device.`,
					},
					"unpaired_external_boot_to_recovery_allowed": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Allow users to boot devices into recovery mode with unpaired devices. Available for devices running iOS and iPadOS versions 14.5 and later.`,
					},
					"usb_restricted_mode_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates if connecting to USB accessories while the device is locked is allowed (iOS 11.4.1 and later).`,
					},
					"voice_dialing_blocked": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block voice dialing.`,
					},
					"vpn_block_creation": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not the creation of VPN configurations is blocked (iOS 11.0 and later).`,
					},
					"wallpaper_block_modification": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to allow wallpaper modification on supervised device (iOS 9.0 and later) .`,
					},
					"wifi_connect_only_to_configured_networks": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `wiFiConnectOnlyToConfiguredNetworks`, // custom MS Graph attribute name
						MarkdownDescription: `Indicates whether or not to force the device to use only Wi-Fi networks from configuration profiles when the device is in supervised mode. Available for devices running iOS and iPadOS versions 14.4 and earlier. Devices running 14.5+ should use the setting, “WiFiConnectToAllowedNetworksOnlyForced.`,
					},
					"wifi_connect_to_allowed_networks_only_forced": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `wiFiConnectToAllowedNetworksOnlyForced`, // custom MS Graph attribute name
						MarkdownDescription: `Require devices to use Wi-Fi networks set up via configuration profiles. Available for devices running iOS and iPadOS versions 14.5 and later.`,
					},
					"wifi_power_on_forced": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not Wi-Fi remains on, even when device is in airplane mode. Available for devices running iOS and iPadOS, versions 13.0 and later.`,
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: `This topic provides descriptions of the declared methods, properties and relationships exposed by the iosGeneralDeviceConfiguration resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosGeneralDeviceConfiguration?view=graph-rest-beta`,
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
						MarkdownDescription: `Active Hours End (active hours mean the time window when updates install should not happen)`,
					},
					"active_hours_start": schema.StringAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.String{
							wpdefaultvalue.StringDefaultValue("00:00:00.0000000"),
						},
						Computed:            true,
						MarkdownDescription: `Active Hours Start (active hours mean the time window when updates install should not happen)`,
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
									MarkdownDescription: `End day of the time window`,
								},
								"end_time": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: `End time of the time window`,
								},
								"start_day": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf("sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"),
									},
									MarkdownDescription: `Start day of the time window`,
								},
								"start_time": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: `Start time of the time window`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `If update schedule type is set to use time window scheduling, custom time windows when updates will be scheduled. This collection can contain a maximum of 20 elements. / Custom update time window / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-customUpdateTimeWindow?view=graph-rest-beta`,
					},
					"desired_os_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `If left unspecified, devices will update to the latest version of the OS.`,
					},
					"enforced_software_update_delay_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Days before software updates are visible to iOS devices ranging from 0 to 90 inclusive`,
					},
					"is_enabled": schema.BoolAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
						MarkdownDescription: `Is setting enabled in UI`,
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
						MarkdownDescription: `Days in week for which active hours are configured. This collection can contain a maximum of 7 elements.`,
					},
					"update_schedule_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("updateOutsideOfActiveHours", "alwaysUpdate", "updateDuringTimeWindows", "updateOutsideOfTimeWindows"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("alwaysUpdate")},
						Computed:            true,
						MarkdownDescription: `Update schedule type`,
					},
					"utc_time_offset_in_minutes": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `UTC Time Offset indicated in minutes`,
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: `IOS Update Configuration, allows you to configure time window within week to install iOS updates / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosUpdateConfiguration?view=graph-rest-beta`,
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
						MarkdownDescription: `Associated Domains`,
					},
					"authentication_method": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("certificate", "usernameAndPassword", "sharedSecret", "derivedCredential", "azureAD"),
						},
						MarkdownDescription: `Authentication method for this VPN connection.`,
					},
					"connection_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: `Connection name displayed to the user.`,
					},
					"connection_type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("ciscoAnyConnect", "pulseSecure", "f5EdgeClient", "dellSonicWallMobileConnect", "checkPointCapsuleVpn", "customVpn", "ciscoIPSec", "citrix", "ciscoAnyConnectV2", "paloAltoGlobalProtect", "zscalerPrivateAccess", "f5Access2018", "citrixSso", "paloAltoGlobalProtectV2", "ikEv2", "alwaysOn", "microsoftTunnel", "netMotionMobility", "microsoftProtect"),
						},
						MarkdownDescription: `Connection type.`,
					},
					"custom_data": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // keyValue
								"key": schema.StringAttribute{
									Optional: true,
								},
								"value": schema.StringAttribute{
									Optional: true,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Custom data when connection type is set to Custom VPN. Use this field to enable functionality not supported by Intune, but available in your VPN solution. Contact your VPN vendor to learn how to add these key/value pairs. This collection can contain a maximum of 25 elements. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyValue?view=graph-rest-beta`,
					},
					"disable_on_demand_user_override": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Toggle to prevent user from disabling automatic VPN in the Settings app`,
					},
					"disconnect_on_idle": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Whether to disconnect after on-demand connection idles`,
					},
					"disconnect_on_idle_timer_in_seconds": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `The length of time in seconds to wait before disconnecting an on-demand connection. Valid values 0 to 65535`,
					},
					"enable_per_app": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Setting this to true creates Per-App VPN payload which can later be associated with Apps that can trigger this VPN conneciton on the end user's iOS device.`,
					},
					"enable_split_tunneling": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Send all network traffic through VPN.`,
					},
					"excluded_domains": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Domains that are accessed through the public internet instead of through VPN, even when per-app VPN is activated`,
					},
					"identifier": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Identifier provided by VPN vendor when connection type is set to Custom VPN. For example: Cisco AnyConnect uses an identifier of the form com.cisco.anyconnect.applevpn.plugin`,
					},
					"login_group_or_domain": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Login group or domain when connection type is set to Dell SonicWALL Mobile Connection.`,
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
									MarkdownDescription: `Action.`,
								},
								"dns_search_domains": schema.SetAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: `DNS Search Domains.`,
								},
								"dns_server_address_match": schema.SetAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: `DNS Search Server Address.`,
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
									MarkdownDescription: `Domain Action (Only applicable when Action is evaluate connection).`,
								},
								"domains": schema.SetAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: `Domains (Only applicable when Action is evaluate connection).`,
								},
								"interface_type_match": schema.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.OneOf("notConfigured", "ethernet", "wiFi", "cellular"),
									},
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
									Computed:            true,
									MarkdownDescription: `Network interface to trigger VPN.`,
								},
								"probe_required_url": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `Probe Required Url (Only applicable when Action is evaluate connection and DomainAction is connect if needed).`,
								},
								"probe_url": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `A URL to probe. If this URL is successfully fetched (returning a 200 HTTP status code) without redirection, this rule matches.`,
								},
								"ssids": schema.SetAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: `Network Service Set Identifiers (SSIDs).`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `On-Demand Rules. This collection can contain a maximum of 500 elements. / VPN On-Demand Rule definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-vpnOnDemandRule?view=graph-rest-beta`,
					},
					"opt_in_to_device_id_sharing": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Opt-In to sharing the device's Id to third-party vpn clients for use during network access control validation.`,
					},
					"provider_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "appProxy", "packetTunnel"),
						},
						MarkdownDescription: `Provider type for per-app VPN.`,
					},
					"proxy_server": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // vpnProxyServer
							"address": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: `Address.`,
							},
							"automatic_configuration_script_url": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: `Proxy's automatic configuration script url.`,
							},
							"port": schema.Int64Attribute{
								Optional:            true,
								MarkdownDescription: `Port. Valid values 0 to 65535`,
							},
						},
						MarkdownDescription: `Proxy Server. / VPN Proxy Server. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-vpnProxyServer?view=graph-rest-beta`,
					},
					"realm": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Realm when connection type is set to Pulse Secure.`,
					},
					"role": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Role when connection type is set to Pulse Secure.`,
					},
					"safari_domains": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Safari domains when this VPN per App setting is enabled. In addition to the apps associated with this VPN, Safari domains specified here will also be able to trigger this VPN connection.`,
					},
					"server": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{ // vpnServer
							"address": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: `Address (IP address, FQDN or URL)`,
							},
							"description": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: `Description.`,
							},
							"is_default_server": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
								Computed:            true,
								MarkdownDescription: `Default server.`,
							},
						},
						MarkdownDescription: `VPN Server on the network. Make sure end users can access this network location. / VPN Server definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-vpnServer?view=graph-rest-beta`,
					},
					"cloud_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Zscaler only. Zscaler cloud which the user is assigned to.`,
					},
					"exclude_list": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Zscaler only. List of network addresses which are not sent through the Zscaler cloud.`,
					},
					"microsoft_tunnel_site_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Microsoft Tunnel site ID.`,
					},
					"strict_enforcement": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Zscaler only. Blocks network traffic until the user signs into Zscaler app. "True" means traffic is blocked.`,
					},
					"targeted_mobile_apps": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationAppListItemAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Targeted mobile apps. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-appListItem?view=graph-rest-beta`,
					},
					"user_domain": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Zscaler only. Enter a static domain to pre-populate the login field with in the Zscaler app. If this is left empty, the user's Azure Active Directory domain will be used instead.`,
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: `By providing the configurations in this profile you can instruct the iOS device to connect to desired VPN endpoint. By specifying the authentication method and security types expected by VPN endpoint you can make the VPN connection seamless for end user. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosVpnConfiguration?view=graph-rest-beta`,
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
							stringvalidator.OneOf("deviceChannel", "userChannel"),
						},
						MarkdownDescription: `Indicates the channel used to deploy the configuration profile. Available choices are DeviceChannel, UserChannel.`,
					},
					"payload": schema.StringAttribute{
						Required:            true,
						Validators:          []validator.String{wpvalidator.TranslateValueEncodeBase64()},
						MarkdownDescription: `Payload. (UTF8 encoded byte array)`,
					},
					"file_name": schema.StringAttribute{
						Required:            true,
						Description:         `payloadFileName`, // custom MS Graph attribute name
						MarkdownDescription: `Payload file name (*.mobileconfig | *.xml).`,
					},
					"name": schema.StringAttribute{
						Required:            true,
						Description:         `payloadName`, // custom MS Graph attribute name
						MarkdownDescription: `Name that is displayed to the user.`,
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: `This topic provides descriptions of the declared methods, properties and relationships exposed by the macOSCustomConfiguration resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macOSCustomConfiguration?view=graph-rest-beta`,
			},
		},
		"macos_custom_app": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.macOSCustomAppConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // macOSCustomAppConfiguration
					"bundle_id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: `Bundle id for targeting.`,
					},
					"configuration_xml": schema.StringAttribute{
						Required:            true,
						Validators:          []validator.String{wpvalidator.TranslateValueEncodeBase64()},
						MarkdownDescription: `Configuration xml. (UTF8 encoded byte array)`,
					},
					"file_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: `Configuration file name (*.plist | *.xml).`,
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: `This topic provides descriptions of the declared methods, properties and relationships exposed by the macOSCustomAppConfiguration resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macOSCustomAppConfiguration?view=graph-rest-beta`,
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
						MarkdownDescription: `An array of AirPrint printers that should always be shown. This collection can contain a maximum of 500 elements. / Represents an AirPrint destination. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-airPrintDestination?view=graph-rest-beta`,
					},
					"admin_show_host_info": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Whether to show admin host information on the login window.`,
					},
					"app_associated_domains": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // macOSAssociatedDomainsItem
								"application_identifier": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: `The application identifier of the app to associate domains with.`,
								},
								"direct_downloads_enabled": schema.BoolAttribute{
									Required:            true,
									MarkdownDescription: `Determines whether data should be downloaded directly or via a CDN.`,
								},
								"domains": schema.SetAttribute{
									ElementType:         types.StringType,
									Required:            true,
									MarkdownDescription: `The list of domains to associate.`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Gets or sets a list that maps apps to their associated domains. Application identifiers must be unique. This collection can contain a maximum of 500 elements. / A mapping of application identifiers to associated domains. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macOSAssociatedDomainsItem?view=graph-rest-beta`,
					},
					"authorized_users_list_hidden": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Whether to show the name and password dialog or a list of users on the login window.`,
					},
					"authorized_users_list_hide_admin_users": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Whether to hide admin users in the authorized users list on the login window.`,
					},
					"authorized_users_list_hide_local_users": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Whether to show only network and system users in the authorized users list on the login window.`,
					},
					"authorized_users_list_hide_mobile_accounts": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Whether to hide mobile users in the authorized users list on the login window.`,
					},
					"authorized_users_list_include_network_users": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Whether to show network users in the authorized users list on the login window.`,
					},
					"authorized_users_list_show_other_managed_users": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Whether to show other users in the authorized users list on the login window.`,
					},
					"auto_launch_items": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // macOSLaunchItem
								"hide": schema.BoolAttribute{
									Required:            true,
									MarkdownDescription: `Whether or not to hide the item from the Users and Groups List.`,
								},
								"path": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: `Path to the launch item.`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `List of applications, files, folders, and other items to launch when the user logs in. This collection can contain a maximum of 500 elements. / Represents an app in the list of macOS launch items / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macOSLaunchItem?view=graph-rest-beta`,
					},
					"console_access_disabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Whether the Other user will disregard use of the '>console> special user name.`,
					},
					"content_caching_block_deletion": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Prevents content caches from purging content to free up disk space for other apps.`,
					},
					"content_caching_client_listen_ranges": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationIpRangeAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `A list of custom IP ranges content caches will use to listen for clients. This collection can contain a maximum of 500 elements. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-ipRange?view=graph-rest-beta`,
					},
					"content_caching_client_policy": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "clientsInLocalNetwork", "clientsWithSamePublicIpAddress", "clientsInCustomLocalNetworks", "clientsInCustomLocalNetworksWithFallback"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: `Determines the method in which content caching servers will listen for clients.`,
					},
					"content_caching_data_path": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The path to the directory used to store cached content. The value must be (or end with) /Library/Application Support/Apple/AssetCache/Data`,
					},
					"content_caching_disable_connection_sharing": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Disables internet connection sharing.`,
					},
					"content_caching_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Enables content caching and prevents it from being disabled by the user.`,
					},
					"content_caching_force_connection_sharing": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Forces internet connection sharing. contentCachingDisableConnectionSharing overrides this setting.`,
					},
					"content_caching_keep_awake": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Prevent the device from sleeping if content caching is enabled.`,
					},
					"content_caching_log_client_identities": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Enables logging of IP addresses and ports of clients that request cached content.`,
					},
					"content_caching_max_size_bytes": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `The maximum number of bytes of disk space that will be used for the content cache. A value of 0 (default) indicates unlimited disk space. `,
					},
					"content_caching_parents": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `A list of IP addresses representing parent content caches.`,
					},
					"content_caching_parent_selection_policy": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "roundRobin", "firstAvailable", "urlPathHash", "random", "stickyAvailable"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: `Determines the method in which content caching servers will select parents if multiple are present.`,
					},
					"content_caching_peer_filter_ranges": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationIpRangeAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `A list of custom IP ranges content caches will use to query for content from peers caches. This collection can contain a maximum of 500 elements. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-ipRange?view=graph-rest-beta`,
					},
					"content_caching_peer_listen_ranges": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationIpRangeAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `A list of custom IP ranges content caches will use to listen for peer caches. This collection can contain a maximum of 500 elements. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-ipRange?view=graph-rest-beta`,
					},
					"content_caching_peer_policy": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "peersInLocalNetwork", "peersWithSamePublicIpAddress", "peersInCustomLocalNetworks"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: `Determines the method in which content caches peer with other caches.`,
					},
					"content_caching_port": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Sets the port used for content caching. If the value is 0, a random available port will be selected. Valid values 0 to 65535`,
					},
					"content_caching_public_ranges": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: deviceConfigurationIpRangeAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `A list of custom IP ranges that Apple's content caching service should use to match clients to content caches. This collection can contain a maximum of 500 elements. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-ipRange?view=graph-rest-beta`,
					},
					"content_caching_show_alerts": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Display content caching alerts as system notifications.`,
					},
					"content_caching_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "userContentOnly", "sharedContentOnly"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: `Determines what type of content is allowed to be cached by Apple's content caching service.`,
					},
					"login_window_text": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							wpvalidator.TranslateGraphNullWithTfDefault(tftypes.String),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
						Computed:            true,
						MarkdownDescription: `Custom text to be displayed on the login window.`,
					},
					"logout_disabled_while_logged_in": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `logOutDisabledWhileLoggedIn`, // custom MS Graph attribute name
						MarkdownDescription: `Whether the Log Out menu item on the login window will be disabled while the user is logged in.`,
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
											MarkdownDescription: `An optional list of additional bundle IDs allowed to use the AAD extension for single sign-on.`,
										},
										"configurations": schema.SetNestedAttribute{
											Optional: true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: deviceConfigurationKeyTypedValuePairAttributes,
											},
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: `Gets or sets a list of typed key-value pairs used to configure Credential-type profiles. This collection can contain a maximum of 500 elements. / A key-value pair with a string key and a typed value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyTypedValuePair?view=graph-rest-beta`,
										},
										"enable_shared_device_mode": schema.BoolAttribute{
											Optional:            true,
											PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
											Computed:            true,
											MarkdownDescription: `Enables or disables shared device mode.`,
										},
									},
									Validators: []validator.Object{
										deviceConfigurationMacOSSingleSignOnExtensionValidator,
									},
									MarkdownDescription: `Represents an Azure AD-type Single Sign-On extension profile for macOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macOSAzureAdSingleSignOnExtension?view=graph-rest-beta`,
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
											MarkdownDescription: `Gets or sets a list of typed key-value pairs used to configure Credential-type profiles. This collection can contain a maximum of 500 elements. / A key-value pair with a string key and a typed value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyTypedValuePair?view=graph-rest-beta`,
										},
										"domains": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: `Gets or sets a list of hosts or domain names for which the app extension performs SSO.`,
										},
										"extension_identifier": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: `Gets or sets the bundle ID of the app extension that performs SSO for the specified URLs.`,
										},
										"realm": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: `Gets or sets the case-sensitive realm name for this profile.`,
										},
										"team_identifier": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: `Gets or sets the team ID of the app extension that performs SSO for the specified URLs.`,
										},
									},
									Validators: []validator.Object{
										deviceConfigurationMacOSSingleSignOnExtensionValidator,
									},
									MarkdownDescription: `Represents a Credential-type Single Sign-On extension profile for macOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macOSCredentialSingleSignOnExtension?view=graph-rest-beta`,
								},
							},
							"kerberos": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.macOSKerberosSingleSignOnExtension",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // macOSKerberosSingleSignOnExtension
										"active_directory_site_code": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: `Gets or sets the Active Directory site.`,
										},
										"block_active_directory_site_auto_discovery": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: `Enables or disables whether the Kerberos extension can automatically determine its site name.`,
										},
										"block_automatic_login": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: `Enables or disables Keychain usage.`,
										},
										"cache_name": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: `Gets or sets the Generic Security Services name of the Kerberos cache to use for this profile.`,
										},
										"credential_bundle_id_access_control_list": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											MarkdownDescription: `Gets or sets a list of app Bundle IDs allowed to access the Kerberos Ticket Granting Ticket.`,
										},
										"credentials_cache_monitored": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: `When set to True, the credential is requested on the next matching Kerberos challenge or network state change. When the credential is expired or missing, a new credential is created. Available for devices running macOS versions 12 and later.`,
										},
										"domain_realms": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											MarkdownDescription: `Gets or sets a list of realms for custom domain-realm mapping. Realms are case sensitive.`,
										},
										"domains": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											MarkdownDescription: `Gets or sets a list of hosts or domain names for which the app extension performs SSO.`,
										},
										"is_default_realm": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: `When true, this profile's realm will be selected as the default. Necessary if multiple Kerberos-type profiles are configured.`,
										},
										"kerberos_apps_in_bundle_id_acl_included": schema.BoolAttribute{
											Required:            true,
											Description:         `kerberosAppsInBundleIdACLIncluded`, // custom MS Graph attribute name
											MarkdownDescription: `When set to True, the Kerberos extension allows any apps entered with the app bundle ID, managed apps, and standard Kerberos utilities, such as TicketViewer and klist, to access and use the credential. Available for devices running macOS versions 12 and later.`,
										},
										"managed_apps_in_bundle_id_acl_included": schema.BoolAttribute{
											Required:            true,
											Description:         `managedAppsInBundleIdACLIncluded`, // custom MS Graph attribute name
											MarkdownDescription: `When set to True, the Kerberos extension allows managed apps, and any apps entered with the app bundle ID to access the credential. When set to False, the Kerberos extension allows all apps to access the credential. Available for devices running iOS and iPadOS versions 14 and later.`,
										},
										"mode_credential_used": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: `Select how other processes use the Kerberos Extension credential.`,
										},
										"password_block_modification": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: `Enables or disables password changes.`,
										},
										"password_change_url": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: `Gets or sets the URL that the user will be sent to when they initiate a password change.`,
										},
										"password_enable_local_sync": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: `Enables or disables password syncing. This won't affect users logged in with a mobile account on macOS.`,
										},
										"password_expiration_days": schema.Int64Attribute{
											Optional:            true,
											MarkdownDescription: `Overrides the default password expiration in days. For most domains, this value is calculated automatically.`,
										},
										"password_expiration_notification_days": schema.Int64Attribute{
											Optional:            true,
											MarkdownDescription: `Gets or sets the number of days until the user is notified that their password will expire (default is 15).`,
										},
										"password_minimum_age_days": schema.Int64Attribute{
											Optional:            true,
											MarkdownDescription: `Gets or sets the minimum number of days until a user can change their password again.`,
										},
										"password_minimum_length": schema.Int64Attribute{
											Optional:            true,
											MarkdownDescription: `Gets or sets the minimum length of a password.`,
										},
										"password_previous_password_block_count": schema.Int64Attribute{
											Optional:            true,
											MarkdownDescription: `Gets or sets the number of previous passwords to block.`,
										},
										"password_require_active_directory_complexity": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: `Enables or disables whether passwords must meet Active Directory's complexity requirements.`,
										},
										"password_requirements_description": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: `Gets or sets a description of the password complexity requirements.`,
										},
										"preferred_kdcs": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											Description:         `preferredKDCs`, // custom MS Graph attribute name
											MarkdownDescription: `Add creates an ordered list of preferred Key Distribution Centers (KDCs) to use for Kerberos traffic. This list is used when the servers are not discoverable using DNS. When the servers are discoverable, the list is used for both connectivity checks, and used first for Kerberos traffic. If the servers don’t respond, then the device uses DNS discovery. Delete removes an existing list, and devices use DNS discovery. Available for devices running macOS versions 12 and later.`,
										},
										"realm": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: `Gets or sets the case-sensitive realm name for this profile.`,
										},
										"require_user_presence": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: `Gets or sets whether to require authentication via Touch ID, Face ID, or a passcode to access the keychain entry.`,
										},
										"sign_in_help_text": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: `Text displayed to the user at the Kerberos sign in window. Available for devices running iOS and iPadOS versions 14 and later.`,
										},
										"tls_for_ldap_required": schema.BoolAttribute{
											Required:            true,
											Description:         `tlsForLDAPRequired`, // custom MS Graph attribute name
											MarkdownDescription: `When set to True, LDAP connections are required to use Transport Layer Security (TLS). Available for devices running macOS versions 11 and later.`,
										},
										"username_label_custom": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: `This label replaces the user name shown in the Kerberos extension. You can enter a name to match the name of your company or organization. Available for devices running macOS versions 11 and later.`,
										},
										"user_principal_name": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: `Gets or sets the principle user name to use for this profile. The realm name does not need to be included.`,
										},
										"user_setup_delayed": schema.BoolAttribute{
											Required:            true,
											MarkdownDescription: `When set to True, the user isn’t prompted to set up the Kerberos extension until the extension is enabled by the admin, or a Kerberos challenge is received. Available for devices running macOS versions 11 and later.`,
										},
									},
									Validators: []validator.Object{
										deviceConfigurationMacOSSingleSignOnExtensionValidator,
									},
									MarkdownDescription: `Represents a Kerberos-type Single Sign-On extension profile for macOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macOSKerberosSingleSignOnExtension?view=graph-rest-beta`,
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
											MarkdownDescription: `Gets or sets a list of typed key-value pairs used to configure Credential-type profiles. This collection can contain a maximum of 500 elements. / A key-value pair with a string key and a typed value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyTypedValuePair?view=graph-rest-beta`,
										},
										"extension_identifier": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: `Gets or sets the bundle ID of the app extension that performs SSO for the specified URLs.`,
										},
										"team_identifier": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: `Gets or sets the team ID of the app extension that performs SSO for the specified URLs.`,
										},
										"url_prefixes": schema.SetAttribute{
											ElementType:         types.StringType,
											Required:            true,
											MarkdownDescription: `One or more URL prefixes of identity providers on whose behalf the app extension performs single sign-on. URLs must begin with http:// or https://. All URL prefixes must be unique for all profiles.`,
										},
									},
									Validators: []validator.Object{
										deviceConfigurationMacOSSingleSignOnExtensionValidator,
									},
									MarkdownDescription: `Represents a Redirect-type Single Sign-On extension profile for macOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macOSRedirectSingleSignOnExtension?view=graph-rest-beta`,
								},
							},
						},
						Description:         `macOSSingleSignOnExtension`, // custom MS Graph attribute name
						MarkdownDescription: `Gets or sets a single sign-on extension profile. / An abstract base class for all macOS-specific single sign-on extension types. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macOSSingleSignOnExtension?view=graph-rest-beta`,
					},
					"poweroff_disabled_while_logged_in": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `powerOffDisabledWhileLoggedIn`, // custom MS Graph attribute name
						MarkdownDescription: `Whether the Power Off menu item on the login window will be disabled while the user is logged in.`,
					},
					"restart_disabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Whether to hide the Restart button item on the login window.`,
					},
					"restart_disabled_while_logged_in": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Whether the Restart menu item on the login window will be disabled while the user is logged in.`,
					},
					"screen_lock_disable_immediate": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Whether to disable the immediate screen lock functions.`,
					},
					"shutdown_disabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `shutDownDisabled`, // custom MS Graph attribute name
						MarkdownDescription: `Whether to hide the Shut Down button item on the login window.`,
					},
					"shutdown_disabled_while_logged_in": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `shutDownDisabledWhileLoggedIn`, // custom MS Graph attribute name
						MarkdownDescription: `Whether the Shut Down menu item on the login window will be disabled while the user is logged in.`,
					},
					"sleep_disabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Whether to hide the Sleep menu item on the login window.`,
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: `MacOS device features configuration profile. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macOSDeviceFeaturesConfiguration?view=graph-rest-beta`,
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
						MarkdownDescription: `All kernel extensions validly signed by the team identifiers in this list will be allowed to load.`,
					},
					"kernel_extension_overrides_allowed": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `If set to true, users can approve additional kernel extensions not explicitly allowed by configurations profiles.`,
					},
					"kernel_extensions_allowed": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // macOSKernelExtension
								"bundle_id": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: `Bundle ID of the kernel extension.`,
								},
								"team_identifier": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `The team identifier that was used to sign the kernel extension.`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `A list of kernel extensions that will be allowed to load. . This collection can contain a maximum of 500 elements. / Represents a specific macOS kernel extension. A macOS kernel extension can be described by its team identifier plus its bundle identifier. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macOSKernelExtension?view=graph-rest-beta`,
					},
					"system_extensions_allowed": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // macOSSystemExtension
								"bundle_id": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: `Gets or sets the bundle identifier of the system extension.`,
								},
								"team_identifier": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `Gets or sets the team identifier that was used to sign the system extension.`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Gets or sets a list of allowed macOS system extensions. This collection can contain a maximum of 500 elements. / Represents a specific macOS system extension. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macOSSystemExtension?view=graph-rest-beta`,
					},
					"system_extensions_allowed_team_identifiers": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Gets or sets a list of allowed team identifiers. Any system extension signed with any of the specified team identifiers will be approved.`,
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
									MarkdownDescription: `Gets or sets the allowed macOS system extension types.`,
								},
								"team_identifier": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: `Gets or sets the team identifier used to sign the system extension.`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Gets or sets a list of allowed macOS system extension types. This collection can contain a maximum of 500 elements. / Represents a mapping between team identifiers for macOS system extensions and system extension types. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macOSSystemExtensionTypeMapping?view=graph-rest-beta`,
					},
					"system_extensions_block_override": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Gets or sets whether to allow the user to approve additional system extensions not explicitly allowed by configuration profiles.`,
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: `MacOS extensions configuration profile. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macOSExtensionsConfiguration?view=graph-rest-beta`,
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
						MarkdownDescription: `Update behavior for all other updates.`,
					},
					"config_data_update_behavior": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "default", "downloadOnly", "installASAP", "notifyOnly", "installLater"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: `Update behavior for configuration data file updates.`,
					},
					"critical_update_behavior": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "default", "downloadOnly", "installASAP", "notifyOnly", "installLater"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: `Update behavior for critical updates.`,
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
									MarkdownDescription: `End day of the time window`,
								},
								"end_time": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: `End time of the time window`,
								},
								"start_day": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf("sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"),
									},
									MarkdownDescription: `Start day of the time window`,
								},
								"start_time": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: `Start time of the time window`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Custom Time windows when updates will be allowed or blocked. This collection can contain a maximum of 20 elements. / Custom update time window / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-customUpdateTimeWindow?view=graph-rest-beta`,
					},
					"firmware_update_behavior": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "default", "downloadOnly", "installASAP", "notifyOnly", "installLater"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: `Update behavior for firmware updates.`,
					},
					"max_user_deferrals_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `The maximum number of times the system allows the user to postpone an update before it’s installed. Supported values: 0 - 366. Valid values 0 to 365`,
					},
					"priority": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("low", "high", "unknownFutureValue"),
						},
						MarkdownDescription: `The scheduling priority for downloading and preparing the requested update. Default: Low. Possible values: Null, Low, High`,
					},
					"update_schedule_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("alwaysUpdate", "updateDuringTimeWindows", "updateOutsideOfTimeWindows"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("alwaysUpdate")},
						Computed:            true,
						MarkdownDescription: `Update schedule type`,
					},
					"update_time_window_utc_offset_in_minutes": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Minutes indicating UTC offset for each update time window`,
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: `MacOS Software Update Configuration / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macOSSoftwareUpdateConfiguration?view=graph-rest-beta`,
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
						MarkdownDescription: `Enables device health monitoring on the device`,
					},
					"custom_scope": schema.StringAttribute{
						Optional:            true,
						Description:         `configDeviceHealthMonitoringCustomScope`, // custom MS Graph attribute name
						MarkdownDescription: `Specifies custom set of events collected from the device where health monitoring is enabled`,
					},
					"scope": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							wpvalidator.FlagEnumValues("undefined", "healthMonitoring", "bootPerformance", "windowsUpdates", "privilegeManagement"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("undefined")},
						Computed:            true,
						Description:         `configDeviceHealthMonitoringScope`, // custom MS Graph attribute name
						MarkdownDescription: `Specifies set of events collected from the device where health monitoring is enabled`,
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: `Windows device health monitoring configuration / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowsHealthMonitoringConfiguration?view=graph-rest-beta`,
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
						MarkdownDescription: `When TRUE, allows eligible Windows 10 devices to upgrade to Windows 11. When FALSE, implies the device stays on the existing operating system. Returned by default. Query parameters are not supported.`,
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
						MarkdownDescription: `The Automatic Update Mode. Possible values are: UserDefined, NotifyDownload, AutoInstallAtMaintenanceTime, AutoInstallAndRebootAtMaintenanceTime, AutoInstallAndRebootAtScheduledTime, AutoInstallAndRebootWithoutEndUserControl, WindowsDefault. UserDefined is the default value, no intent. Returned by default. Query parameters are not supported.`,
					},
					"auto_restart_notification_dismissal": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "automatic", "user", "unknownFutureValue"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: `Specify the method by which the auto-restart required notification is dismissed. Possible values are: NotConfigured, Automatic, User. Returned by default. Query parameters are not supported.`,
					},
					"business_ready_updates_only": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userDefined", "all", "businessReadyOnly", "windowsInsiderBuildFast", "windowsInsiderBuildSlow", "windowsInsiderBuildRelease"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("userDefined")},
						Computed:            true,
						MarkdownDescription: `Determines which branch devices will receive their updates from. Possible values are: UserDefined, All, BusinessReadyOnly, WindowsInsiderBuildFast, WindowsInsiderBuildSlow, WindowsInsiderBuildRelease. Returned by default. Query parameters are not supported.`,
					},
					"deadline_for_feature_updates_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of days before feature updates are installed automatically with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.`,
					},
					"deadline_for_quality_updates_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of days before quality updates are installed automatically with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.`,
					},
					"deadline_grace_period_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of days after deadline until restarts occur automatically with valid range from 0 to 7 days. Returned by default. Query parameters are not supported.`,
					},
					"delivery_optimization_mode": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userDefined", "httpOnly", "httpWithPeeringNat", "httpWithPeeringPrivateGroup", "httpWithInternetPeering", "simpleDownload", "bypassMode"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("userDefined")},
						Computed:            true,
						MarkdownDescription: `The Delivery Optimization Mode. Possible values are: UserDefined, HttpOnly, HttpWithPeeringNat, HttpWithPeeringPrivateGroup, HttpWithInternetPeering, SimpleDownload, BypassMode. UserDefined allows the user to set. Returned by default. Query parameters are not supported.`,
					},
					"drivers_excluded": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `When TRUE, excludes Windows update Drivers. When FALSE, does not exclude Windows update Drivers. Returned by default. Query parameters are not supported.`,
					},
					"engaged_restart_deadline_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Deadline in days before automatically scheduling and executing a pending restart outside of active hours, with valid range from 2 to 30 days. Returned by default. Query parameters are not supported.`,
					},
					"engaged_restart_snooze_schedule_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of days a user can snooze Engaged Restart reminder notifications with valid range from 1 to 3 days. Returned by default. Query parameters are not supported.`,
					},
					"engaged_restart_transition_schedule_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of days before transitioning from Auto Restarts scheduled outside of active hours to Engaged Restart, which requires the user to schedule, with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.`,
					},
					"feature_updates_deferral_period_in_days": schema.Int64Attribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
						Computed:            true,
						MarkdownDescription: `Defer Feature Updates by these many days with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.`,
					},
					"feature_updates_paused": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `When TRUE, assigned devices are paused from receiving feature updates for up to 35 days from the time you pause the ring. When FALSE, does not pause Feature Updates. Returned by default. Query parameters are not supported.s`,
					},
					"feature_updates_rollback_window_in_days": schema.Int64Attribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(10)},
						Computed:            true,
						MarkdownDescription: `The number of days after a Feature Update for which a rollback is valid with valid range from 2 to 60 days. Returned by default. Query parameters are not supported.`,
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
											MarkdownDescription: `Active Hours End`,
										},
										"start": schema.StringAttribute{
											Optional: true,
											PlanModifiers: []planmodifier.String{
												wpdefaultvalue.StringDefaultValue("08:00:00.0000000"),
											},
											Computed:            true,
											Description:         `activeHoursStart`, // custom MS Graph attribute name
											MarkdownDescription: `Active Hours Start`,
										},
									},
									Validators: []validator.Object{
										deviceConfigurationWindowsUpdateInstallScheduleTypeValidator,
									},
									MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowsUpdateActiveHoursInstall?view=graph-rest-beta`,
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
											MarkdownDescription: `Scheduled Install Day in week`,
										},
										"time": schema.StringAttribute{
											Optional: true,
											PlanModifiers: []planmodifier.String{
												wpdefaultvalue.StringDefaultValue("03:00:00.0000000"),
											},
											Computed:            true,
											Description:         `scheduledInstallTime`, // custom MS Graph attribute name
											MarkdownDescription: `Scheduled Install Time during day`,
										},
									},
									Validators: []validator.Object{
										deviceConfigurationWindowsUpdateInstallScheduleTypeValidator,
									},
									MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowsUpdateScheduledInstall?view=graph-rest-beta`,
								},
							},
						},
						MarkdownDescription: `The Installation Schedule. Possible values are: ActiveHoursStart, ActiveHoursEnd, ScheduledInstallDay, ScheduledInstallTime. Returned by default. Query parameters are not supported. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowsUpdateInstallScheduleType?view=graph-rest-beta`,
					},
					"microsoft_update_service_allowed": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
						Computed:            true,
						MarkdownDescription: `When TRUE, allows Microsoft Update Service. When FALSE, does not allow Microsoft Update Service. Returned by default. Query parameters are not supported.`,
					},
					"postpone_reboot_until_after_deadline": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `When TRUE the device should wait until deadline for rebooting outside of active hours. When FALSE the device should not wait until deadline for rebooting outside of active hours. Returned by default. Query parameters are not supported.`,
					},
					"prerelease_features": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("userDefined", "settingsOnly", "settingsAndExperimentations", "notAllowed"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("userDefined")},
						Computed:            true,
						MarkdownDescription: `The Pre-Release Features. Possible values are: UserDefined, SettingsOnly, SettingsAndExperimentations, NotAllowed. UserDefined is the default value, no intent. Returned by default. Query parameters are not supported.`,
					},
					"quality_updates_deferral_period_in_days": schema.Int64Attribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
						Computed:            true,
						MarkdownDescription: `Defer Quality Updates by these many days with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.`,
					},
					"quality_updates_paused": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `When TRUE, assigned devices are paused from receiving quality updates for up to 35 days from the time you pause the ring. When FALSE, does not pause Quality Updates. Returned by default. Query parameters are not supported.`,
					},
					"schedule_imminent_restart_warning_in_minutes": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Specify the period for auto-restart imminent warning notifications. Supported values: 15, 30 or 60 (minutes). Returned by default. Query parameters are not supported.`,
					},
					"schedule_restart_warning_in_hours": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Specify the period for auto-restart warning reminder notifications. Supported values: 2, 4, 8, 12 or 24 (hours). Returned by default. Query parameters are not supported.`,
					},
					"skip_checks_before_restart": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `When TRUE, skips all checks before restart: Battery level = 40%, User presence, Display Needed, Presentation mode, Full screen mode, phone call state, game mode etc. When FALSE, does not skip all checks before restart. Returned by default. Query parameters are not supported.`,
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
						MarkdownDescription: `Specifies what Windows Update notifications users see. Possible values are: NotConfigured, DefaultNotifications, RestartWarningsOnly, DisableAllNotifications. Returned by default. Query parameters are not supported.`,
					},
					"update_weeks": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							wpvalidator.FlagEnumValues("userDefined", "firstWeek", "secondWeek", "thirdWeek", "fourthWeek", "everyWeek", "unknownFutureValue"),
						},
						MarkdownDescription: `Schedule the update installation on the weeks of the month. Possible values are: UserDefined, FirstWeek, SecondWeek, ThirdWeek, FourthWeek, EveryWeek. Returned by default. Query parameters are not supported.`,
					},
					"user_pause_access": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("enabled")},
						Computed:            true,
						MarkdownDescription: `Specifies whether to enable end user’s access to pause software updates. Possible values are: NotConfigured, Enabled, Disabled. Returned by default. Query parameters are not supported.`,
					},
					"user_windows_update_scan_access": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "enabled", "disabled"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("enabled")},
						Computed:            true,
						MarkdownDescription: `Specifies whether to disable user’s access to scan Windows Update. Possible values are: NotConfigured, Enabled, Disabled. Returned by default. Query parameters are not supported.`,
					},
				},
				Validators:          []validator.Object{deviceConfigurationDeviceConfigurationValidator},
				MarkdownDescription: `Windows Update for business configuration, allows you to specify how and when Windows as a Service updates your Windows 10/11 devices with feature and quality updates. Supports ODATA clauses that DeviceConfiguration entity supports: $filter by types of DeviceConfiguration, $top, $select only DeviceConfiguration base properties, $orderby only DeviceConfiguration base properties, and $skip. The query parameter '$search' is not supported. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowsUpdateForBusinessConfiguration?view=graph-rest-beta`,
			},
		},
	},
	MarkdownDescription: `Device Configuration. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceConfiguration?view=graph-rest-beta`,
}

var deviceConfigurationDeviceConfigurationValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("android_device_owner_general_device"),
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
)

var deviceConfigurationAppListItemAttributes = map[string]schema.Attribute{ // appListItem
	"app_id": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: `The application or bundle identifier of the application`,
	},
	"app_store_url": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: `The Store URL of the application`,
	},
	"name": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: `The application name`,
	},
	"publisher": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: `The publisher of the application`,
	},
}

var deviceConfigurationKeyValuePairAttributes = map[string]schema.Attribute{ // keyValuePair
	"name": schema.StringAttribute{
		Required: true,
	},
	"value": schema.StringAttribute{
		Optional: true,
	},
}

var deviceConfigurationAirPrintDestinationAttributes = map[string]schema.Attribute{ // airPrintDestination
	"force_tls": schema.BoolAttribute{
		Required:            true,
		MarkdownDescription: `If true AirPrint connections are secured by Transport Layer Security (TLS). Default is false. Available in iOS 11.0 and later.`,
	},
	"ip_address": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: `The IP Address of the AirPrint destination.`,
	},
	"port": schema.Int64Attribute{
		Optional:            true,
		MarkdownDescription: `The listening port of the AirPrint destination. If this key is not specified AirPrint will use the default port. Available in iOS 11.0 and later.`,
	},
	"resource_path": schema.StringAttribute{
		Required: true,
	},
}

var deviceConfigurationIosWebContentFilterBaseValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("auto_filter"),
	path.MatchRelative().AtParent().AtName("specific_websites_access"),
)

var deviceConfigurationIosBookmarkAttributes = map[string]schema.Attribute{ // iosBookmark
	"bookmark_folder": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: `The folder into which the bookmark should be added in Safari`,
	},
	"display_name": schema.StringAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
		Computed:            true,
		MarkdownDescription: `The display name of the bookmark`,
	},
	"url": schema.StringAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
		Computed:            true,
		MarkdownDescription: `URL allowed to access`,
	},
}

var deviceConfigurationIpRangeAttributes = map[string]schema.Attribute{ // ipRange
	"v4": generic.OdataDerivedTypeNestedAttributeRs{
		DerivedType: "#microsoft.graph.iPv4Range",
		SingleNestedAttribute: schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // iPv4Range
				"lower_address": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: `Lower address.`,
				},
				"upper_address": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: `Upper address.`,
				},
			},
			Validators:          []validator.Object{deviceConfigurationIpRangeValidator},
			MarkdownDescription: `IPv4 Range definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-iPv4Range?view=graph-rest-beta`,
		},
	},
	"v4_cidr": generic.OdataDerivedTypeNestedAttributeRs{
		DerivedType: "#microsoft.graph.iPv4CidrRange",
		SingleNestedAttribute: schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // iPv4CidrRange
				"cidr_address": schema.StringAttribute{
					Required: true,
				},
			},
			Validators:          []validator.Object{deviceConfigurationIpRangeValidator},
			MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/iPv4CidrRange?view=graph-rest-beta`,
		},
	},
	"v6": generic.OdataDerivedTypeNestedAttributeRs{
		DerivedType: "#microsoft.graph.iPv6Range",
		SingleNestedAttribute: schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // iPv6Range
				"lower_address": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: `Lower address.`,
				},
				"upper_address": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: `Upper address.`,
				},
			},
			Validators:          []validator.Object{deviceConfigurationIpRangeValidator},
			MarkdownDescription: `IPv6 Range definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-iPv6Range?view=graph-rest-beta`,
		},
	},
	"v6_cidr": generic.OdataDerivedTypeNestedAttributeRs{
		DerivedType: "#microsoft.graph.iPv6CidrRange",
		SingleNestedAttribute: schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // iPv6CidrRange
				"cidr_address": schema.StringAttribute{
					Required: true,
				},
			},
			Validators:          []validator.Object{deviceConfigurationIpRangeValidator},
			MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/iPv6CidrRange?view=graph-rest-beta`,
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

var deviceConfigurationKeyTypedValuePairAttributes = map[string]schema.Attribute{ // keyTypedValuePair
	"key": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: `The string key of the key-value pair.`,
	},
	"boolean": generic.OdataDerivedTypeNestedAttributeRs{
		DerivedType: "#microsoft.graph.keyBooleanValuePair",
		SingleNestedAttribute: schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // keyBooleanValuePair
				"value": schema.BoolAttribute{
					Required:            true,
					MarkdownDescription: `The Boolean value of the key-value pair.`,
				},
			},
			Validators:          []validator.Object{deviceConfigurationKeyTypedValuePairValidator},
			MarkdownDescription: `A key-value pair with a string key and a Boolean value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyBooleanValuePair?view=graph-rest-beta`,
		},
	},
	"integer": generic.OdataDerivedTypeNestedAttributeRs{
		DerivedType: "#microsoft.graph.keyIntegerValuePair",
		SingleNestedAttribute: schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // keyIntegerValuePair
				"value": schema.Int64Attribute{
					Required:            true,
					MarkdownDescription: `The integer value of the key-value pair.`,
				},
			},
			Validators:          []validator.Object{deviceConfigurationKeyTypedValuePairValidator},
			MarkdownDescription: `A key-value pair with a string key and an integer value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyIntegerValuePair?view=graph-rest-beta`,
		},
	},
	"real": generic.OdataDerivedTypeNestedAttributeRs{
		DerivedType: "#microsoft.graph.keyRealValuePair",
		SingleNestedAttribute: schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // keyRealValuePair
				"value": schema.NumberAttribute{
					Required:            true,
					MarkdownDescription: `The real (floating-point) value of the key-value pair.`,
				},
			},
			Validators:          []validator.Object{deviceConfigurationKeyTypedValuePairValidator},
			MarkdownDescription: `A key-value pair with a string key and a real (floating-point) value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyRealValuePair?view=graph-rest-beta`,
		},
	},
	"string": generic.OdataDerivedTypeNestedAttributeRs{
		DerivedType: "#microsoft.graph.keyStringValuePair",
		SingleNestedAttribute: schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // keyStringValuePair
				"value": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: `The string value of the key-value pair.`,
				},
			},
			Validators:          []validator.Object{deviceConfigurationKeyTypedValuePairValidator},
			MarkdownDescription: `A key-value pair with a string key and a string value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyStringValuePair?view=graph-rest-beta`,
		},
	},
}

var deviceConfigurationKeyTypedValuePairValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("boolean"),
	path.MatchRelative().AtParent().AtName("integer"),
	path.MatchRelative().AtParent().AtName("real"),
	path.MatchRelative().AtParent().AtName("string"),
)

var deviceConfigurationWindowsUpdateInstallScheduleTypeValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("active_hours"),
	path.MatchRelative().AtParent().AtName("scheduled"),
)
