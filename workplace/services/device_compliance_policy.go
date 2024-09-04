package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	DeviceCompliancePolicyResource = generic.GenericResource{
		TypeNameSuffix: "device_compliance_policy",
		SpecificSchema: deviceCompliancePolicyResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:         "/deviceManagement/deviceCompliancePolicies",
			ReadOptions:     deviceCompliancePolicyReadOptions,
			WriteSubActions: deviceCompliancePolicyWriteSubActions,
		},
	}

	DeviceCompliancePolicySingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&DeviceCompliancePolicyResource)

	DeviceCompliancePolicyPluralDataSource = generic.CreateGenericDataSourcePluralFromSingular(
		&DeviceCompliancePolicySingularDataSource, "")
)

var deviceCompliancePolicyReadOptions = generic.ReadOptions{
	ODataExpand: "scheduledActionsForRule($expand=scheduledActionConfigurations),assignments",
}

var deviceCompliancePolicyWriteSubActions = []generic.WriteSubAction{
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			AttributesMap: map[string]string{"scheduledActionsForRule": "deviceComplianceScheduledActionForRules"},
			UriSuffix:     "scheduleActionsForRules",
			UpdateOnly:    true,
		},
	},
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			Attributes: []string{"assignments"},
			UriSuffix:  "assign",
		},
	},
}

var deviceCompliancePolicyResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // deviceCompliancePolicy
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
		"version": schema.Int64Attribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
			MarkdownDescription: `Version of the device configuration.`,
		},
		"assignments": deviceAndAppManagementAssignment,
		"scheduled_actions_for_rule": schema.SetNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // deviceComplianceScheduledActionForRule
					"scheduled_action_configurations": schema.SetNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // deviceComplianceActionItem
								"action_type": schema.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.OneOf("noAction", "notification", "block", "retire", "wipe", "removeResourceAccessProfiles", "pushNotification", "remoteLock"),
									},
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("block")},
									Computed:            true,
									MarkdownDescription: `What action to take`,
								},
								"grace_period_hours": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: `Number of hours to wait till the action will be enforced. Valid values 0 to 8760`,
								},
								"notification_message_cc_list": schema.SetAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
									Computed:            true,
									Description:         `notificationMessageCCList`, // custom MS Graph attribute name
									MarkdownDescription: `A list of group IDs to speicify who to CC this notification message to.`,
								},
								"notification_template_id": schema.StringAttribute{
									Optional: true,
									PlanModifiers: []planmodifier.String{
										wpdefaultvalue.StringDefaultValue("00000000-0000-0000-0000-000000000000"),
									},
									Computed:            true,
									MarkdownDescription: `What notification Message template to use`,
								},
							},
						},
						Validators:          []validator.Set{setvalidator.SizeAtLeast(1)},
						MarkdownDescription: `The list of scheduled action configurations for this compliance policy. Compliance policy must have one and only one block scheduled action. / Scheduled Action Configuration / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceComplianceActionItem?view=graph-rest-beta`,
					},
				},
			},
			Validators:          []validator.Set{setvalidator.SizeBetween(1, 1)},
			MarkdownDescription: `The list of scheduled action per rule for this compliance policy. This is a required property when creating any individual per-platform compliance policies. / Scheduled Action for Rule / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceComplianceScheduledActionForRule?view=graph-rest-beta`,
		},
		"android": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.androidCompliancePolicy",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // androidCompliancePolicy
					"advanced_threat_protection_required_security_level": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unavailable")},
						Computed:            true,
						MarkdownDescription: `MDATP Require Mobile Threat Protection minimum risk level to report noncompliance.`,
					},
					"condition_statement_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Condition statement id.`,
					},
					"device_threat_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require that devices have enabled device threat protection.`,
					},
					"device_threat_protection_required_security_level": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unavailable")},
						Computed:            true,
						MarkdownDescription: `Require Mobile Threat Protection minimum risk level to report noncompliance.`,
					},
					"min_android_security_patch_level": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Minimum Android security patch level.`,
					},
					"os_maximum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Maximum Android version.`,
					},
					"os_minimum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Minimum Android version.`,
					},
					"password_expiration_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of days before the password expires. Valid values 1 to 365`,
					},
					"password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Minimum password length. Valid values 4 to 16`,
					},
					"password_minutes_of_inactivity_before_lock": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Minutes of inactivity before a password is required.`,
					},
					"password_previous_password_block_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of previous passwords to block. Valid values 1 to 24`,
					},
					"password_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require a password to unlock device.`,
					},
					"password_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "alphabetic", "alphanumeric", "alphanumericWithSymbols", "lowSecurityBiometric", "numeric", "numericComplex", "any"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: `Type of characters in password`,
					},
					"password_sign_in_failure_count_before_factory_reset": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of sign-in failures allowed before factory reset. Valid values 1 to 16`,
					},
					"required_password_complexity": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("none", "low", "medium", "high"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
						Computed:            true,
						MarkdownDescription: `Indicates the required password complexity on Android. One of: NONE, LOW, MEDIUM, HIGH. This is a new API targeted to Android 11+.`,
					},
					"restricted_apps": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // appListItem
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
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Require the device to not have the specified apps installed. This collection can contain a maximum of 100 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-appListItem?view=graph-rest-beta`,
					},
					"security_block_device_administrator_managed_devices": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Block device administrator managed devices.`,
					},
					"security_block_jailbroken_devices": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Devices must not be jailbroken or rooted.`,
					},
					"security_disable_usb_debugging": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Disable USB debugging on Android devices.`,
					},
					"security_prevent_install_apps_from_unknown_sources": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require that devices disallow installation of apps from unknown sources.`,
					},
					"security_require_company_portal_app_integrity": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require the device to pass the Company Portal client app runtime integrity check.`,
					},
					"security_require_google_play_services": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require Google Play Services to be installed and enabled on the device.`,
					},
					"security_require_safety_net_attestation_basic_integrity": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require the device to pass the SafetyNet basic integrity check.`,
					},
					"security_require_safety_net_attestation_certified_device": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require the device to pass the SafetyNet certified device check.`,
					},
					"security_require_up_to_date_security_providers": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require the device to have up to date security providers. The device will require Google Play Services to be enabled and up to date.`,
					},
					"security_require_verify_apps": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require the Android Verify apps feature is turned on.`,
					},
					"storage_require_encryption": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require encryption on Android devices.`,
					},
				},
				Validators: []validator.Object{
					deviceCompliancePolicyDeviceCompliancePolicyValidator,
				},
				MarkdownDescription: `This class contains compliance settings for Android. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androidCompliancePolicy?view=graph-rest-beta`,
			},
		},
		"android_device_owner": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.androidDeviceOwnerCompliancePolicy",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // androidDeviceOwnerCompliancePolicy
					"advanced_threat_protection_required_security_level": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unavailable")},
						Computed:            true,
						MarkdownDescription: `MDATP Require Mobile Threat Protection minimum risk level to report noncompliance.`,
					},
					"device_threat_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require that devices have enabled device threat protection.`,
					},
					"device_threat_protection_required_security_level": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unavailable")},
						Computed:            true,
						MarkdownDescription: `Require Mobile Threat Protection minimum risk level to report noncompliance.`,
					},
					"min_android_security_patch_level": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Minimum Android security patch level.`,
					},
					"os_maximum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Maximum Android version.`,
					},
					"os_minimum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Minimum Android version.`,
					},
					"password_expiration_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of days before the password expires. Valid values 1 to 365`,
					},
					"password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Minimum password length. Valid values 4 to 16`,
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
					"password_minutes_of_inactivity_before_lock": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Minutes of inactivity before a password is required.`,
					},
					"password_previous_password_count_to_block": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of previous passwords to block. Valid values 1 to 24`,
					},
					"password_required": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Require a password to unlock device.`,
					},
					"password_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "required", "numeric", "numericComplex", "alphabetic", "alphanumeric", "alphanumericWithSymbols", "lowSecurityBiometric", "customPassword"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("numeric")},
						Computed:            true,
						MarkdownDescription: `Type of characters in password`,
					},
					"require_no_pending_system_updates": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Require device to have no pending Android system updates.`,
					},
					"security_required_android_safety_net_evaluation_type": schema.StringAttribute{
						Optional:            true,
						Validators:          []validator.String{stringvalidator.OneOf("basic", "hardwareBacked")},
						MarkdownDescription: `Require a specific Play Integrity evaluation type for compliance.`,
					},
					"security_require_intune_app_integrity": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `If setting is set to true, checks that the Intune app installed on fully managed, dedicated, or corporate-owned work profile Android Enterprise enrolled devices, is the one provided by Microsoft from the Managed Google Playstore. If the check fails, the device will be reported as non-compliant.`,
					},
					"security_require_safety_net_attestation_basic_integrity": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require the device to pass the Play Integrity basic integrity check.`,
					},
					"security_require_safety_net_attestation_certified_device": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require the device to pass the Play Integrity device integrity check.`,
					},
					"storage_require_encryption": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Require encryption on Android devices.`,
					},
				},
				Validators: []validator.Object{
					deviceCompliancePolicyDeviceCompliancePolicyValidator,
				},
				MarkdownDescription: `This topic provides descriptions of the declared methods, properties and relationships exposed by the AndroidDeviceOwnerCompliancePolicy resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androidDeviceOwnerCompliancePolicy?view=graph-rest-beta`,
			},
		},
		"android_work_profile": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.androidWorkProfileCompliancePolicy",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // androidWorkProfileCompliancePolicy
					"advanced_threat_protection_required_security_level": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unavailable")},
						Computed:            true,
						MarkdownDescription: `MDATP Require Mobile Threat Protection minimum risk level to report noncompliance.`,
					},
					"device_threat_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require that devices have enabled device threat protection.`,
					},
					"device_threat_protection_required_security_level": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unavailable")},
						Computed:            true,
						MarkdownDescription: `Require Mobile Threat Protection minimum risk level to report noncompliance.`,
					},
					"min_android_security_patch_level": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Minimum Android security patch level.`,
					},
					"os_maximum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Maximum Android version.`,
					},
					"os_minimum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Minimum Android version.`,
					},
					"password_expiration_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of days before the password expires. Valid values 1 to 365`,
					},
					"password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Minimum password length. Valid values 4 to 16`,
					},
					"password_minutes_of_inactivity_before_lock": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Minutes of inactivity before a password is required.`,
					},
					"password_previous_password_block_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of previous passwords to block. Valid values 1 to 24`,
					},
					"password_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require a password to unlock device.`,
					},
					"password_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "alphabetic", "alphanumeric", "alphanumericWithSymbols", "lowSecurityBiometric", "numeric", "numericComplex", "any"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: `Type of characters in password`,
					},
					"password_sign_in_failure_count_before_factory_reset": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of sign-in failures allowed before factory reset. Valid values 1 to 16`,
					},
					"required_password_complexity": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("none", "low", "medium", "high"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
						Computed:            true,
						MarkdownDescription: `Indicates the required device password complexity on Android. One of: NONE, LOW, MEDIUM, HIGH. This is a new API targeted to Android API 12+.`,
					},
					"security_block_jailbroken_devices": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Devices must not be jailbroken or rooted.`,
					},
					"security_disable_usb_debugging": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Disable USB debugging on Android devices.`,
					},
					"security_prevent_install_apps_from_unknown_sources": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require that devices disallow installation of apps from unknown sources.`,
					},
					"security_require_company_portal_app_integrity": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require the device to pass the Company Portal client app runtime integrity check.`,
					},
					"security_required_android_safety_net_evaluation_type": schema.StringAttribute{
						Optional:            true,
						Validators:          []validator.String{stringvalidator.OneOf("basic", "hardwareBacked")},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("basic")},
						Computed:            true,
						MarkdownDescription: `Require a specific SafetyNet evaluation type for compliance.`,
					},
					"security_require_google_play_services": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require Google Play Services to be installed and enabled on the device.`,
					},
					"security_require_safety_net_attestation_basic_integrity": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require the device to pass the Play Integrity basic integrity check.`,
					},
					"security_require_safety_net_attestation_certified_device": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require the device to pass the Play Integrity device integrity check.`,
					},
					"security_require_up_to_date_security_providers": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require the device to have up to date security providers. The device will require Google Play Services to be enabled and up to date.`,
					},
					"security_require_verify_apps": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require the Android Verify apps feature is turned on.`,
					},
					"storage_require_encryption": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require encryption on Android devices.`,
					},
					"work_profile_inactive_before_screen_lock_in_minutes": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Minutes of inactivity before the screen times out.`,
					},
					"work_profile_password_expiration_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of days before the work profile password expires. Valid values 1 to 365`,
					},
					"work_profile_password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Minimum length of work profile password. Valid values 4 to 16`,
					},
					"work_profile_password_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "lowSecurityBiometric", "required", "atLeastNumeric", "numericComplex", "atLeastAlphabetic", "atLeastAlphanumeric", "alphanumericWithSymbols"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: `Type of work profile password that is required.`,
					},
					"work_profile_previous_password_block_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of previous work profile passwords to block. Valid values 0 to 24`,
					},
					"work_profile_required_password_complexity": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("none", "low", "medium", "high"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
						Computed:            true,
						MarkdownDescription: `Indicates the required work profile password complexity on Android. One of: NONE, LOW, MEDIUM, HIGH. This is a new API targeted to Android 12+.`,
					},
					"work_profile_require_password": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Password is required or not for work profile`,
					},
				},
				Validators: []validator.Object{
					deviceCompliancePolicyDeviceCompliancePolicyValidator,
				},
				MarkdownDescription: `This class contains compliance settings for Android Work Profile. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androidWorkProfileCompliancePolicy?view=graph-rest-beta`,
			},
		},
		"aosp_device_owner": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.aospDeviceOwnerCompliancePolicy",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // aospDeviceOwnerCompliancePolicy
					"min_android_security_patch_level": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Minimum Android security patch level.`,
					},
					"os_maximum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Maximum Android version.`,
					},
					"os_minimum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Minimum Android version.`,
					},
					"password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Minimum password length. Valid values 4 to 16`,
					},
					"password_minutes_of_inactivity_before_lock": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Minutes of inactivity before a password is required. Valid values 1 to 8640`,
					},
					"password_required": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Require a password to unlock device.`,
					},
					"password_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "required", "numeric", "numericComplex", "alphabetic", "alphanumeric", "alphanumericWithSymbols", "lowSecurityBiometric", "customPassword"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: `Type of characters in password`,
					},
					"security_block_jailbroken_devices": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Devices must not be jailbroken or rooted.`,
					},
					"storage_require_encryption": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Require encryption on Android devices.`,
					},
				},
				Validators: []validator.Object{
					deviceCompliancePolicyDeviceCompliancePolicyValidator,
				},
				MarkdownDescription: `This topic provides descriptions of the declared methods, properties and relationships exposed by the AndroidDeviceOwnerAOSPCompliancePolicy resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-aospDeviceOwnerCompliancePolicy?view=graph-rest-beta`,
			},
		},
		"ios": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.iosCompliancePolicy",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // iosCompliancePolicy
					"advanced_threat_protection_required_security_level": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unavailable")},
						Computed:            true,
						MarkdownDescription: `MDATP Require Mobile Threat Protection minimum risk level to report noncompliance.`,
					},
					"device_threat_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require that devices have enabled device threat protection .`,
					},
					"device_threat_protection_required_security_level": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unavailable")},
						Computed:            true,
						MarkdownDescription: `Require Mobile Threat Protection minimum risk level to report noncompliance.`,
					},
					"managed_email_profile_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to require a managed email profile.`,
					},
					"os_maximum_build_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Maximum IOS build version.`,
					},
					"os_maximum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Maximum IOS version.`,
					},
					"os_minimum_build_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Minimum IOS build version.`,
					},
					"os_minimum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Minimum IOS version.`,
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
						MarkdownDescription: `The number of character sets required in the password.`,
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
						MarkdownDescription: `The required passcode type.`,
					},
					"restricted_apps": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // appListItem
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
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Require the device to not have the specified apps installed. This collection can contain a maximum of 100 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-appListItem?view=graph-rest-beta`,
					},
					"security_block_jailbroken_devices": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Devices must not be jailbroken or rooted.`,
					},
				},
				Validators: []validator.Object{
					deviceCompliancePolicyDeviceCompliancePolicyValidator,
				},
				MarkdownDescription: `This class contains compliance settings for IOS. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosCompliancePolicy?view=graph-rest-beta`,
			},
		},
		"macos": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.macOSCompliancePolicy",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // macOSCompliancePolicy
					"advanced_threat_protection_required_security_level": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unavailable")},
						Computed:            true,
						MarkdownDescription: `MDATP Require Mobile Threat Protection minimum risk level to report noncompliance.`,
					},
					"device_threat_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require that devices have enabled device threat protection.`,
					},
					"device_threat_protection_required_security_level": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unavailable")},
						Computed:            true,
						MarkdownDescription: `Require Mobile Threat Protection minimum risk level to report noncompliance.`,
					},
					"firewall_block_all_incoming": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Corresponds to the “Block all incoming connections” option.`,
					},
					"firewall_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Whether the firewall should be enabled or not.`,
					},
					"firewall_enable_stealth_mode": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Corresponds to “Enable stealth mode.”`,
					},
					"gatekeeper_allowed_app_source": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "macAppStore", "macAppStoreAndIdentifiedDevelopers", "anywhere"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: `System and Privacy setting that determines which download locations apps can be run from on a macOS device.`,
					},
					"os_maximum_build_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Maximum MacOS build version.`,
					},
					"os_maximum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Maximum MacOS version.`,
					},
					"os_minimum_build_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Minimum MacOS build version.`,
					},
					"os_minimum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Minimum MacOS version.`,
					},
					"password_block_simple": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block simple passwords.`,
					},
					"password_expiration_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of days before the password expires. Valid values 1 to 65535`,
					},
					"password_minimum_character_set_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `The number of character sets required in the password.`,
					},
					"password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Minimum length of password. Valid values 4 to 14`,
					},
					"password_minutes_of_inactivity_before_lock": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Minutes of inactivity before a password is required.`,
					},
					"password_previous_password_block_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Number of previous passwords to block. Valid values 1 to 24`,
					},
					"password_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Whether or not to require a password.`,
					},
					"password_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "alphanumeric", "numeric"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: `The required password type.`,
					},
					"storage_require_encryption": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require encryption on Mac OS devices.`,
					},
					"system_integrity_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require that devices have enabled system integrity protection.`,
					},
				},
				Validators: []validator.Object{
					deviceCompliancePolicyDeviceCompliancePolicyValidator,
				},
				MarkdownDescription: `This class contains compliance settings for Mac OS. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macOSCompliancePolicy?view=graph-rest-beta`,
			},
		},
		"windows10": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.windows10CompliancePolicy",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // windows10CompliancePolicy
					"active_firewall_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require active firewall on Windows devices.`,
					},
					"anti_spyware_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require any AntiSpyware solution registered with Windows Decurity Center to be on and monitoring (e.g. Symantec, Windows Defender).`,
					},
					"antivirus_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require any Antivirus solution registered with Windows Decurity Center to be on and monitoring (e.g. Symantec, Windows Defender).`,
					},
					"bitlocker_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `bitLockerEnabled`, // custom MS Graph attribute name
						MarkdownDescription: `Require devices to be reported healthy by Windows Device Health Attestation - bit locker is enabled`,
					},
					"code_integrity_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require devices to be reported as healthy by Windows Device Health Attestation.`,
					},
					"configuration_manager_compliance_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require to consider SCCM Compliance state into consideration for Intune Compliance State.`,
					},
					"defender_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require Windows Defender Antimalware on Windows devices.`,
					},
					"defender_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Require Windows Defender Antimalware minimum version on Windows devices.`,
					},
					"device_compliance_policy_script": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // deviceCompliancePolicyScript
							"device_compliance_script_id": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: `Device compliance script Id.`,
							},
							"rules_content_base64": schema.StringAttribute{
								Optional:            true,
								Description:         `rulesContent`, // custom MS Graph attribute name
								MarkdownDescription: `Json of the rules.`,
							},
						},
						MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceCompliancePolicyScript?view=graph-rest-beta`,
					},
					"device_threat_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require that devices have enabled device threat protection.`,
					},
					"device_threat_protection_required_security_level": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unavailable")},
						Computed:            true,
						MarkdownDescription: `Require Device Threat Protection minimum risk level to report noncompliance.`,
					},
					"early_launch_anti_malware_driver_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require devices to be reported as healthy by Windows Device Health Attestation - early launch antimalware driver is enabled.`,
					},
					"firmware_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `When TRUE, indicates that Firmware protection is required to be reported as healthy by Microsoft Azure Attestion. When FALSE, indicates that Firmware protection is not required to be reported as healthy. Devices that support either Dynamic Root of Trust for Measurement (DRTM) or Firmware Attack Surface Reduction (FASR) will report compliant for this setting. Default value is FALSE.`,
					},
					"kernel_dma_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `When TRUE, indicates that Kernel Direct Memory Access (DMA) protection is required to be reported as healthy by Microsoft Azure Attestion. When FALSE, indicates that Kernel DMA Protection is not required to be reported as healthy. Default value is FALSE.`,
					},
					"memory_integrity_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `When TRUE, indicates that Memory Integrity as known as Hypervisor-protected Code Integrity (HVCI) or Hypervisor Enforced Code Integrity protection is required to be reported as healthy by Microsoft Azure Attestion. When FALSE, indicates that Memory Integrity Protection is not required to be reported as healthy. Default value is FALSE.`,
					},
					"mobile_os_maximum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Maximum Windows Phone version.`,
					},
					"mobile_os_minimum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Minimum Windows Phone version.`,
					},
					"os_maximum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Maximum Windows 10 version.`,
					},
					"os_minimum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Minimum Windows 10 version.`,
					},
					"password_block_simple": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether or not to block simple password.`,
					},
					"password_expiration_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `The password expiration in days.`,
					},
					"password_minimum_character_set_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `The number of character sets required in the password.`,
					},
					"password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `The minimum password length.`,
					},
					"password_minutes_of_inactivity_before_lock": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `Minutes of inactivity before a password is required.`,
					},
					"password_previous_password_block_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `The number of previous passwords to prevent re-use of.`,
					},
					"password_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require a password to unlock Windows device.`,
					},
					"password_required_to_unlock_from_idle": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require a password to unlock an idle device.`,
					},
					"password_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "alphanumeric", "numeric"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: `The required password type.`,
					},
					"require_healthy_device_report": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require devices to be reported as healthy by Windows Device Health Attestation.`,
					},
					"rtp_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require Windows Defender Antimalware Real-Time Protection on Windows devices.`,
					},
					"secure_boot_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require devices to be reported as healthy by Windows Device Health Attestation - secure boot is enabled.`,
					},
					"signature_out_of_date": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require Windows Defender Antimalware Signature to be up to date on Windows devices.`,
					},
					"storage_require_encryption": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require encryption on windows devices.`,
					},
					"tpm_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Require Trusted Platform Module(TPM) to be present.`,
					},
					"valid_operating_system_build_ranges": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // operatingSystemVersionRange
								"description": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: `The description of this range (e.g. Valid 1702 builds)`,
								},
								"highest_version": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: `The highest inclusive version that this range contains.`,
								},
								"lowest_version": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: `The lowest inclusive version that this range contains.`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `The valid operating system build ranges on Windows devices. This collection can contain a maximum of 10000 elements. / Operating System version range. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-operatingSystemVersionRange?view=graph-rest-beta`,
					},
					"virtualization_based_security_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `When TRUE, indicates that Virtualization-based Security is required to be reported as healthy by Microsoft Azure Attestion. When FALSE, indicates that Virtualization-based Security is not required to be reported as healthy. Default value is FALSE.`,
					},
				},
				Validators: []validator.Object{
					deviceCompliancePolicyDeviceCompliancePolicyValidator,
				},
				MarkdownDescription: `This class contains compliance settings for Windows 10. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windows10CompliancePolicy?view=graph-rest-beta`,
			},
		},
	},
	MarkdownDescription: `This is the base class for Compliance policy. Compliance policies are platform specific and individual per-platform compliance policies inherit from here.  / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceCompliancePolicy?view=graph-rest-beta`,
}

var deviceCompliancePolicyDeviceCompliancePolicyValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("android"),
	path.MatchRelative().AtParent().AtName("android_device_owner"),
	path.MatchRelative().AtParent().AtName("android_work_profile"),
	path.MatchRelative().AtParent().AtName("aosp_device_owner"),
	path.MatchRelative().AtParent().AtName("ios"),
	path.MatchRelative().AtParent().AtName("macos"),
	path.MatchRelative().AtParent().AtName("windows10"),
)
