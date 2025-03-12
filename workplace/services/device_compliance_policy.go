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

	DeviceCompliancePolicyPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&DeviceCompliancePolicyResource, "")
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
		"version": schema.Int64Attribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
			MarkdownDescription: "Version of the device configuration.",
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
									MarkdownDescription: "What action to take. / Scheduled Action Type Enum; possible values are: `noAction` (No Action), `notification` (Send Notification), `block` (Block the device in AAD), `retire` (Retire the device), `wipe` (Wipe the device), `removeResourceAccessProfiles` (Remove Resource Access Profiles from the device), `pushNotification` (Send push notification to device), `remoteLock` (Remotely lock the device). The _provider_ default value is `\"block\"`.",
								},
								"grace_period_hours": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: "Number of hours to wait till the action will be enforced. Valid values 0 to 8760. The _provider_ default value is `0`.",
								},
								"notification_message_cc_list": schema.SetAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
									Computed:            true,
									Description:         `notificationMessageCCList`, // custom MS Graph attribute name
									MarkdownDescription: "A list of group IDs to speicify who to CC this notification message to. The _provider_ default value is `[]`.",
								},
								"notification_template_id": schema.StringAttribute{
									Optional: true,
									PlanModifiers: []planmodifier.String{
										wpdefaultvalue.StringDefaultValue("00000000-0000-0000-0000-000000000000"),
									},
									Computed:            true,
									MarkdownDescription: "What notification Message template to use. The _provider_ default value is `\"00000000-0000-0000-0000-000000000000\"`.",
								},
							},
						},
						Validators:          []validator.Set{setvalidator.SizeAtLeast(1)},
						MarkdownDescription: "The list of scheduled action configurations for this compliance policy. Compliance policy must have one and only one block scheduled action. / Scheduled Action Configuration / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-devicecomplianceactionitem?view=graph-rest-beta",
					},
				},
			},
			Validators:          []validator.Set{setvalidator.SizeBetween(1, 1)},
			MarkdownDescription: "The list of scheduled action per rule for this compliance policy. This is a required property when creating any individual per-platform compliance policies. / Scheduled Action for Rule / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-devicecompliancescheduledactionforrule?view=graph-rest-beta",
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
						MarkdownDescription: "MDATP Require Mobile Threat Protection minimum risk level to report noncompliance. / Device threat protection levels for the Device Threat Protection API; possible values are: `unavailable` (Default Value. Do not use.), `secured` (Device Threat Level requirement: Secured. This is the most secure level, and represents that no threats were found on the device.), `low` (Device Threat Protection level requirement: Low. Low represents a severity of threat that poses minimal risk to the device or device data.), `medium` (Device Threat Protection level requirement: Medium. Medium represents a severity of threat that poses moderate risk to the device or device data.), `high` (Device Threat Protection level requirement: High. High represents a severity of threat that poses severe risk to the device or device data.), `notSet` (Device Threat Protection level requirement: Not Set. Not set represents that there is no requirement for the device to meet a Threat Protection level.). The _provider_ default value is `\"unavailable\"`.",
					},
					"condition_statement_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Condition statement id.",
					},
					"device_threat_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require that devices have enabled device threat protection. The _provider_ default value is `false`.",
					},
					"device_threat_protection_required_security_level": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unavailable")},
						Computed:            true,
						MarkdownDescription: "Require Mobile Threat Protection minimum risk level to report noncompliance. / Device threat protection levels for the Device Threat Protection API; possible values are: `unavailable` (Default Value. Do not use.), `secured` (Device Threat Level requirement: Secured. This is the most secure level, and represents that no threats were found on the device.), `low` (Device Threat Protection level requirement: Low. Low represents a severity of threat that poses minimal risk to the device or device data.), `medium` (Device Threat Protection level requirement: Medium. Medium represents a severity of threat that poses moderate risk to the device or device data.), `high` (Device Threat Protection level requirement: High. High represents a severity of threat that poses severe risk to the device or device data.), `notSet` (Device Threat Protection level requirement: Not Set. Not set represents that there is no requirement for the device to meet a Threat Protection level.). The _provider_ default value is `\"unavailable\"`.",
					},
					"min_android_security_patch_level": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum Android security patch level.",
					},
					"os_maximum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Maximum Android version.",
					},
					"os_minimum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum Android version.",
					},
					"password_expiration_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days before the password expires. Valid values 1 to 365",
					},
					"password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minimum password length. Valid values 4 to 16",
					},
					"password_minutes_of_inactivity_before_lock": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before a password is required.",
					},
					"password_previous_password_block_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of previous passwords to block. Valid values 1 to 24",
					},
					"password_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require a password to unlock device. The _provider_ default value is `false`.",
					},
					"password_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "alphabetic", "alphanumeric", "alphanumericWithSymbols", "lowSecurityBiometric", "numeric", "numericComplex", "any"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: "Type of characters in password. / Android required password type; possible values are: `deviceDefault` (Device default value, no intent.), `alphabetic` (Alphabetic password required.), `alphanumeric` (Alphanumeric password required.), `alphanumericWithSymbols` (Alphanumeric with symbols password required.), `lowSecurityBiometric` (Low security biometrics based password required.), `numeric` (Numeric password required.), `numericComplex` (Numeric complex password required.), `any` (A password or pattern is required, and any is acceptable.). The _provider_ default value is `\"deviceDefault\"`.",
					},
					"password_sign_in_failure_count_before_factory_reset": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of sign-in failures allowed before factory reset. Valid values 1 to 16",
					},
					"required_password_complexity": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("none", "low", "medium", "high"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
						Computed:            true,
						MarkdownDescription: "Indicates the required password complexity on Android. One of: NONE, LOW, MEDIUM, HIGH. This is a new API targeted to Android 11+. / The password complexity types that can be set on Android. One of: NONE, LOW, MEDIUM, HIGH. This is an API targeted to Android 11+; possible values are: `none` (Device default value, no password.), `low` (The required password complexity on the device is of type low as defined by the Android documentation.), `medium` (The required password complexity on the device is of type medium as defined by the Android documentation.), `high` (The required password complexity on the device is of type high as defined by the Android documentation.). The _provider_ default value is `\"none\"`.",
					},
					"restricted_apps": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // appListItem
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
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Require the device to not have the specified apps installed. This collection can contain a maximum of 100 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"security_block_device_administrator_managed_devices": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Block device administrator managed devices. The _provider_ default value is `false`.",
					},
					"security_block_jailbroken_devices": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Devices must not be jailbroken or rooted. The _provider_ default value is `false`.",
					},
					"security_disable_usb_debugging": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Disable USB debugging on Android devices. The _provider_ default value is `false`.",
					},
					"security_prevent_install_apps_from_unknown_sources": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require that devices disallow installation of apps from unknown sources. The _provider_ default value is `false`.",
					},
					"security_require_company_portal_app_integrity": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require the device to pass the Company Portal client app runtime integrity check. The _provider_ default value is `false`.",
					},
					"security_require_google_play_services": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require Google Play Services to be installed and enabled on the device. The _provider_ default value is `false`.",
					},
					"security_require_safety_net_attestation_basic_integrity": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require the device to pass the SafetyNet basic integrity check. The _provider_ default value is `false`.",
					},
					"security_require_safety_net_attestation_certified_device": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require the device to pass the SafetyNet certified device check. The _provider_ default value is `false`.",
					},
					"security_require_up_to_date_security_providers": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require the device to have up to date security providers. The device will require Google Play Services to be enabled and up to date. The _provider_ default value is `false`.",
					},
					"security_require_verify_apps": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require the Android Verify apps feature is turned on. The _provider_ default value is `false`.",
					},
					"storage_require_encryption": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require encryption on Android devices. The _provider_ default value is `false`.",
					},
				},
				Validators: []validator.Object{
					deviceCompliancePolicyDeviceCompliancePolicyValidator,
				},
				MarkdownDescription: "This class contains compliance settings for Android. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androidcompliancepolicy?view=graph-rest-beta",
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
						MarkdownDescription: "MDATP Require Mobile Threat Protection minimum risk level to report noncompliance. / Device threat protection levels for the Device Threat Protection API; possible values are: `unavailable` (Default Value. Do not use.), `secured` (Device Threat Level requirement: Secured. This is the most secure level, and represents that no threats were found on the device.), `low` (Device Threat Protection level requirement: Low. Low represents a severity of threat that poses minimal risk to the device or device data.), `medium` (Device Threat Protection level requirement: Medium. Medium represents a severity of threat that poses moderate risk to the device or device data.), `high` (Device Threat Protection level requirement: High. High represents a severity of threat that poses severe risk to the device or device data.), `notSet` (Device Threat Protection level requirement: Not Set. Not set represents that there is no requirement for the device to meet a Threat Protection level.). The _provider_ default value is `\"unavailable\"`.",
					},
					"device_threat_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require that devices have enabled device threat protection. The _provider_ default value is `false`.",
					},
					"device_threat_protection_required_security_level": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unavailable")},
						Computed:            true,
						MarkdownDescription: "Require Mobile Threat Protection minimum risk level to report noncompliance. / Device threat protection levels for the Device Threat Protection API; possible values are: `unavailable` (Default Value. Do not use.), `secured` (Device Threat Level requirement: Secured. This is the most secure level, and represents that no threats were found on the device.), `low` (Device Threat Protection level requirement: Low. Low represents a severity of threat that poses minimal risk to the device or device data.), `medium` (Device Threat Protection level requirement: Medium. Medium represents a severity of threat that poses moderate risk to the device or device data.), `high` (Device Threat Protection level requirement: High. High represents a severity of threat that poses severe risk to the device or device data.), `notSet` (Device Threat Protection level requirement: Not Set. Not set represents that there is no requirement for the device to meet a Threat Protection level.). The _provider_ default value is `\"unavailable\"`.",
					},
					"min_android_security_patch_level": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum Android security patch level.",
					},
					"os_maximum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Maximum Android version.",
					},
					"os_minimum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum Android version.",
					},
					"password_expiration_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days before the password expires. Valid values 1 to 365",
					},
					"password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minimum password length. Valid values 4 to 16",
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
					"password_minutes_of_inactivity_before_lock": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before a password is required.",
					},
					"password_previous_password_count_to_block": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of previous passwords to block. Valid values 1 to 24",
					},
					"password_required": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require a password to unlock device.",
					},
					"password_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "required", "numeric", "numericComplex", "alphabetic", "alphanumeric", "alphanumericWithSymbols", "lowSecurityBiometric", "customPassword"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("numeric")},
						Computed:            true,
						MarkdownDescription: "Type of characters in password. / Android Device Owner policy required password type; possible values are: `deviceDefault` (Device default value, no intent.), `required` (There must be a password set, but there are no restrictions on type.), `numeric` (At least numeric.), `numericComplex` (At least numeric with no repeating or ordered sequences.), `alphabetic` (At least alphabetic password.), `alphanumeric` (At least alphanumeric password), `alphanumericWithSymbols` (At least alphanumeric with symbols.), `lowSecurityBiometric` (Low security biometrics based password required.), `customPassword` (Custom password set by the admin.). The _provider_ default value is `\"numeric\"`.",
					},
					"require_no_pending_system_updates": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require device to have no pending Android system updates.",
					},
					"security_required_android_safety_net_evaluation_type": schema.StringAttribute{
						Optional:            true,
						Validators:          []validator.String{stringvalidator.OneOf("basic", "hardwareBacked")},
						MarkdownDescription: "Require a specific Play Integrity evaluation type for compliance. / An enum representing the Android Play Integrity API evaluation types; possible values are: `basic` (Default value. Typical measurements and reference data were used.), `hardwareBacked` (Strong Integrity checks (such as a hardware-backed proof of boot integrity) were used.)",
					},
					"security_require_intune_app_integrity": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "If setting is set to true, checks that the Intune app installed on fully managed, dedicated, or corporate-owned work profile Android Enterprise enrolled devices, is the one provided by Microsoft from the Managed Google Playstore. If the check fails, the device will be reported as non-compliant.",
					},
					"security_require_safety_net_attestation_basic_integrity": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require the device to pass the Play Integrity basic integrity check. The _provider_ default value is `false`.",
					},
					"security_require_safety_net_attestation_certified_device": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require the device to pass the Play Integrity device integrity check. The _provider_ default value is `false`.",
					},
					"storage_require_encryption": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require encryption on Android devices.",
					},
				},
				Validators: []validator.Object{
					deviceCompliancePolicyDeviceCompliancePolicyValidator,
				},
				MarkdownDescription: "This topic provides descriptions of the declared methods, properties and relationships exposed by the AndroidDeviceOwnerCompliancePolicy resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceownercompliancepolicy?view=graph-rest-beta",
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
						MarkdownDescription: "MDATP Require Mobile Threat Protection minimum risk level to report noncompliance. / Device threat protection levels for the Device Threat Protection API; possible values are: `unavailable` (Default Value. Do not use.), `secured` (Device Threat Level requirement: Secured. This is the most secure level, and represents that no threats were found on the device.), `low` (Device Threat Protection level requirement: Low. Low represents a severity of threat that poses minimal risk to the device or device data.), `medium` (Device Threat Protection level requirement: Medium. Medium represents a severity of threat that poses moderate risk to the device or device data.), `high` (Device Threat Protection level requirement: High. High represents a severity of threat that poses severe risk to the device or device data.), `notSet` (Device Threat Protection level requirement: Not Set. Not set represents that there is no requirement for the device to meet a Threat Protection level.). The _provider_ default value is `\"unavailable\"`.",
					},
					"device_threat_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require that devices have enabled device threat protection. The _provider_ default value is `false`.",
					},
					"device_threat_protection_required_security_level": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unavailable")},
						Computed:            true,
						MarkdownDescription: "Require Mobile Threat Protection minimum risk level to report noncompliance. / Device threat protection levels for the Device Threat Protection API; possible values are: `unavailable` (Default Value. Do not use.), `secured` (Device Threat Level requirement: Secured. This is the most secure level, and represents that no threats were found on the device.), `low` (Device Threat Protection level requirement: Low. Low represents a severity of threat that poses minimal risk to the device or device data.), `medium` (Device Threat Protection level requirement: Medium. Medium represents a severity of threat that poses moderate risk to the device or device data.), `high` (Device Threat Protection level requirement: High. High represents a severity of threat that poses severe risk to the device or device data.), `notSet` (Device Threat Protection level requirement: Not Set. Not set represents that there is no requirement for the device to meet a Threat Protection level.). The _provider_ default value is `\"unavailable\"`.",
					},
					"min_android_security_patch_level": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum Android security patch level.",
					},
					"os_maximum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Maximum Android version.",
					},
					"os_minimum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum Android version.",
					},
					"password_expiration_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days before the password expires. Valid values 1 to 365",
					},
					"password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minimum password length. Valid values 4 to 16",
					},
					"password_minutes_of_inactivity_before_lock": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before a password is required.",
					},
					"password_previous_password_block_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of previous passwords to block. Valid values 1 to 24",
					},
					"password_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require a password to unlock device. The _provider_ default value is `false`.",
					},
					"password_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "alphabetic", "alphanumeric", "alphanumericWithSymbols", "lowSecurityBiometric", "numeric", "numericComplex", "any"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: "Type of characters in password. / Android required password type; possible values are: `deviceDefault` (Device default value, no intent.), `alphabetic` (Alphabetic password required.), `alphanumeric` (Alphanumeric password required.), `alphanumericWithSymbols` (Alphanumeric with symbols password required.), `lowSecurityBiometric` (Low security biometrics based password required.), `numeric` (Numeric password required.), `numericComplex` (Numeric complex password required.), `any` (A password or pattern is required, and any is acceptable.). The _provider_ default value is `\"deviceDefault\"`.",
					},
					"password_sign_in_failure_count_before_factory_reset": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of sign-in failures allowed before factory reset. Valid values 1 to 16",
					},
					"required_password_complexity": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("none", "low", "medium", "high"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
						Computed:            true,
						MarkdownDescription: "Indicates the required device password complexity on Android. One of: NONE, LOW, MEDIUM, HIGH. This is a new API targeted to Android API 12+. / The password complexity types that can be set on Android. One of: NONE, LOW, MEDIUM, HIGH. This is an API targeted to Android 11+; possible values are: `none` (Device default value, no password.), `low` (The required password complexity on the device is of type low as defined by the Android documentation.), `medium` (The required password complexity on the device is of type medium as defined by the Android documentation.), `high` (The required password complexity on the device is of type high as defined by the Android documentation.). The _provider_ default value is `\"none\"`.",
					},
					"security_block_jailbroken_devices": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Devices must not be jailbroken or rooted. The _provider_ default value is `false`.",
					},
					"security_disable_usb_debugging": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Disable USB debugging on Android devices. The _provider_ default value is `false`.",
					},
					"security_prevent_install_apps_from_unknown_sources": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require that devices disallow installation of apps from unknown sources. The _provider_ default value is `false`.",
					},
					"security_require_company_portal_app_integrity": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require the device to pass the Company Portal client app runtime integrity check. The _provider_ default value is `false`.",
					},
					"security_required_android_safety_net_evaluation_type": schema.StringAttribute{
						Optional:            true,
						Validators:          []validator.String{stringvalidator.OneOf("basic", "hardwareBacked")},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("basic")},
						Computed:            true,
						MarkdownDescription: "Require a specific SafetyNet evaluation type for compliance. / An enum representing the Android Play Integrity API evaluation types; possible values are: `basic` (Default value. Typical measurements and reference data were used.), `hardwareBacked` (Strong Integrity checks (such as a hardware-backed proof of boot integrity) were used.). The _provider_ default value is `\"basic\"`.",
					},
					"security_require_google_play_services": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require Google Play Services to be installed and enabled on the device. The _provider_ default value is `false`.",
					},
					"security_require_safety_net_attestation_basic_integrity": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require the device to pass the Play Integrity basic integrity check. The _provider_ default value is `false`.",
					},
					"security_require_safety_net_attestation_certified_device": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require the device to pass the Play Integrity device integrity check. The _provider_ default value is `false`.",
					},
					"security_require_up_to_date_security_providers": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require the device to have up to date security providers. The device will require Google Play Services to be enabled and up to date. The _provider_ default value is `false`.",
					},
					"security_require_verify_apps": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require the Android Verify apps feature is turned on. The _provider_ default value is `false`.",
					},
					"storage_require_encryption": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require encryption on Android devices. The _provider_ default value is `false`.",
					},
					"work_profile_inactive_before_screen_lock_in_minutes": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before the screen times out.",
					},
					"work_profile_password_expiration_in_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days before the work profile password expires. Valid values 1 to 365",
					},
					"work_profile_password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minimum length of work profile password. Valid values 4 to 16",
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
					"work_profile_previous_password_block_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of previous work profile passwords to block. Valid values 0 to 24",
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
				Validators: []validator.Object{
					deviceCompliancePolicyDeviceCompliancePolicyValidator,
				},
				MarkdownDescription: "This class contains compliance settings for Android Work Profile. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androidworkprofilecompliancepolicy?view=graph-rest-beta",
			},
		},
		"aosp_device_owner": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.aospDeviceOwnerCompliancePolicy",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // aospDeviceOwnerCompliancePolicy
					"min_android_security_patch_level": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum Android security patch level.",
					},
					"os_maximum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Maximum Android version.",
					},
					"os_minimum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum Android version.",
					},
					"password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minimum password length. Valid values 4 to 16",
					},
					"password_minutes_of_inactivity_before_lock": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before a password is required. Valid values 1 to 8640",
					},
					"password_required": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require a password to unlock device.",
					},
					"password_required_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("deviceDefault", "required", "numeric", "numericComplex", "alphabetic", "alphanumeric", "alphanumericWithSymbols", "lowSecurityBiometric", "customPassword"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("deviceDefault")},
						Computed:            true,
						MarkdownDescription: "Type of characters in password. / Android Device Owner policy required password type; possible values are: `deviceDefault` (Device default value, no intent.), `required` (There must be a password set, but there are no restrictions on type.), `numeric` (At least numeric.), `numericComplex` (At least numeric with no repeating or ordered sequences.), `alphabetic` (At least alphabetic password.), `alphanumeric` (At least alphanumeric password), `alphanumericWithSymbols` (At least alphanumeric with symbols.), `lowSecurityBiometric` (Low security biometrics based password required.), `customPassword` (Custom password set by the admin.). The _provider_ default value is `\"deviceDefault\"`.",
					},
					"security_block_jailbroken_devices": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Devices must not be jailbroken or rooted.",
					},
					"storage_require_encryption": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Require encryption on Android devices.",
					},
				},
				Validators: []validator.Object{
					deviceCompliancePolicyDeviceCompliancePolicyValidator,
				},
				MarkdownDescription: "This topic provides descriptions of the declared methods, properties and relationships exposed by the AndroidDeviceOwnerAOSPCompliancePolicy resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-aospdeviceownercompliancepolicy?view=graph-rest-beta",
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
						MarkdownDescription: "MDATP Require Mobile Threat Protection minimum risk level to report noncompliance. / Device threat protection levels for the Device Threat Protection API; possible values are: `unavailable` (Default Value. Do not use.), `secured` (Device Threat Level requirement: Secured. This is the most secure level, and represents that no threats were found on the device.), `low` (Device Threat Protection level requirement: Low. Low represents a severity of threat that poses minimal risk to the device or device data.), `medium` (Device Threat Protection level requirement: Medium. Medium represents a severity of threat that poses moderate risk to the device or device data.), `high` (Device Threat Protection level requirement: High. High represents a severity of threat that poses severe risk to the device or device data.), `notSet` (Device Threat Protection level requirement: Not Set. Not set represents that there is no requirement for the device to meet a Threat Protection level.). The _provider_ default value is `\"unavailable\"`.",
					},
					"device_threat_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require that devices have enabled device threat protection . The _provider_ default value is `false`.",
					},
					"device_threat_protection_required_security_level": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unavailable")},
						Computed:            true,
						MarkdownDescription: "Require Mobile Threat Protection minimum risk level to report noncompliance. / Device threat protection levels for the Device Threat Protection API; possible values are: `unavailable` (Default Value. Do not use.), `secured` (Device Threat Level requirement: Secured. This is the most secure level, and represents that no threats were found on the device.), `low` (Device Threat Protection level requirement: Low. Low represents a severity of threat that poses minimal risk to the device or device data.), `medium` (Device Threat Protection level requirement: Medium. Medium represents a severity of threat that poses moderate risk to the device or device data.), `high` (Device Threat Protection level requirement: High. High represents a severity of threat that poses severe risk to the device or device data.), `notSet` (Device Threat Protection level requirement: Not Set. Not set represents that there is no requirement for the device to meet a Threat Protection level.). The _provider_ default value is `\"unavailable\"`.",
					},
					"managed_email_profile_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to require a managed email profile. The _provider_ default value is `false`.",
					},
					"os_maximum_build_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Maximum IOS build version.",
					},
					"os_maximum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Maximum IOS version.",
					},
					"os_minimum_build_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum IOS build version.",
					},
					"os_minimum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum IOS version.",
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
						MarkdownDescription: "The number of character sets required in the password.",
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
						MarkdownDescription: "The required passcode type. / Possible values of required passwords; possible values are: `deviceDefault` (Device default value, no intent.), `alphanumeric` (Alphanumeric password required.), `numeric` (Numeric password required.). The _provider_ default value is `\"deviceDefault\"`.",
					},
					"restricted_apps": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // appListItem
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
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Require the device to not have the specified apps installed. This collection can contain a maximum of 100 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"security_block_jailbroken_devices": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Devices must not be jailbroken or rooted. The _provider_ default value is `false`.",
					},
				},
				Validators: []validator.Object{
					deviceCompliancePolicyDeviceCompliancePolicyValidator,
				},
				MarkdownDescription: "This class contains compliance settings for IOS. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioscompliancepolicy?view=graph-rest-beta",
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
						MarkdownDescription: "MDATP Require Mobile Threat Protection minimum risk level to report noncompliance. / Device threat protection levels for the Device Threat Protection API; possible values are: `unavailable` (Default Value. Do not use.), `secured` (Device Threat Level requirement: Secured. This is the most secure level, and represents that no threats were found on the device.), `low` (Device Threat Protection level requirement: Low. Low represents a severity of threat that poses minimal risk to the device or device data.), `medium` (Device Threat Protection level requirement: Medium. Medium represents a severity of threat that poses moderate risk to the device or device data.), `high` (Device Threat Protection level requirement: High. High represents a severity of threat that poses severe risk to the device or device data.), `notSet` (Device Threat Protection level requirement: Not Set. Not set represents that there is no requirement for the device to meet a Threat Protection level.). The _provider_ default value is `\"unavailable\"`.",
					},
					"device_threat_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require that devices have enabled device threat protection. The _provider_ default value is `false`.",
					},
					"device_threat_protection_required_security_level": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unavailable")},
						Computed:            true,
						MarkdownDescription: "Require Mobile Threat Protection minimum risk level to report noncompliance. / Device threat protection levels for the Device Threat Protection API; possible values are: `unavailable` (Default Value. Do not use.), `secured` (Device Threat Level requirement: Secured. This is the most secure level, and represents that no threats were found on the device.), `low` (Device Threat Protection level requirement: Low. Low represents a severity of threat that poses minimal risk to the device or device data.), `medium` (Device Threat Protection level requirement: Medium. Medium represents a severity of threat that poses moderate risk to the device or device data.), `high` (Device Threat Protection level requirement: High. High represents a severity of threat that poses severe risk to the device or device data.), `notSet` (Device Threat Protection level requirement: Not Set. Not set represents that there is no requirement for the device to meet a Threat Protection level.). The _provider_ default value is `\"unavailable\"`.",
					},
					"firewall_block_all_incoming": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Corresponds to the “Block all incoming connections” option. The _provider_ default value is `false`.",
					},
					"firewall_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether the firewall should be enabled or not. The _provider_ default value is `false`.",
					},
					"firewall_enable_stealth_mode": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Corresponds to “Enable stealth mode.”. The _provider_ default value is `false`.",
					},
					"gatekeeper_allowed_app_source": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "macAppStore", "macAppStoreAndIdentifiedDevelopers", "anywhere"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
						Computed:            true,
						MarkdownDescription: "System and Privacy setting that determines which download locations apps can be run from on a macOS device. / App source options for macOS Gatekeeper; possible values are: `notConfigured` (Device default value, no intent.), `macAppStore` (Only apps from the Mac AppStore can be run.), `macAppStoreAndIdentifiedDevelopers` (Only apps from the Mac AppStore and identified developers can be run.), `anywhere` (Apps from anywhere can be run.). The _provider_ default value is `\"notConfigured\"`.",
					},
					"os_maximum_build_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Maximum MacOS build version.",
					},
					"os_maximum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Maximum MacOS version.",
					},
					"os_minimum_build_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum MacOS build version.",
					},
					"os_minimum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum MacOS version.",
					},
					"password_block_simple": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block simple passwords. The _provider_ default value is `false`.",
					},
					"password_expiration_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of days before the password expires. Valid values 1 to 65535",
					},
					"password_minimum_character_set_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The number of character sets required in the password.",
					},
					"password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minimum length of password. Valid values 4 to 14",
					},
					"password_minutes_of_inactivity_before_lock": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before a password is required.",
					},
					"password_previous_password_block_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Number of previous passwords to block. Valid values 1 to 24",
					},
					"password_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether or not to require a password. The _provider_ default value is `false`.",
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
					"storage_require_encryption": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require encryption on Mac OS devices. The _provider_ default value is `false`.",
					},
					"system_integrity_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require that devices have enabled system integrity protection. The _provider_ default value is `false`.",
					},
				},
				Validators: []validator.Object{
					deviceCompliancePolicyDeviceCompliancePolicyValidator,
				},
				MarkdownDescription: "This class contains compliance settings for Mac OS. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macoscompliancepolicy?view=graph-rest-beta",
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
						MarkdownDescription: "Require active firewall on Windows devices. The _provider_ default value is `false`.",
					},
					"anti_spyware_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require any AntiSpyware solution registered with Windows Decurity Center to be on and monitoring (e.g. Symantec, Windows Defender). The _provider_ default value is `false`.",
					},
					"antivirus_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require any Antivirus solution registered with Windows Decurity Center to be on and monitoring (e.g. Symantec, Windows Defender). The _provider_ default value is `false`.",
					},
					"bitlocker_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						Description:         `bitLockerEnabled`, // custom MS Graph attribute name
						MarkdownDescription: "Require devices to be reported healthy by Windows Device Health Attestation - bit locker is enabled. The _provider_ default value is `false`.",
					},
					"code_integrity_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require devices to be reported as healthy by Windows Device Health Attestation. The _provider_ default value is `false`.",
					},
					"configuration_manager_compliance_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require to consider SCCM Compliance state into consideration for Intune Compliance State. The _provider_ default value is `false`.",
					},
					"defender_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require Windows Defender Antimalware on Windows devices. The _provider_ default value is `false`.",
					},
					"defender_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Require Windows Defender Antimalware minimum version on Windows devices.",
					},
					"device_compliance_policy_script": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // deviceCompliancePolicyScript
							"device_compliance_script_id": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Device compliance script Id.",
							},
							"rules_content": schema.StringAttribute{
								Optional:            true,
								Validators:          []validator.String{wpvalidator.TranslateValueEncodeBase64()},
								MarkdownDescription: "Json of the rules.",
							},
						},
						MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-devicecompliancepolicyscript?view=graph-rest-beta",
					},
					"device_threat_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require that devices have enabled device threat protection. The _provider_ default value is `false`.",
					},
					"device_threat_protection_required_security_level": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("unavailable", "secured", "low", "medium", "high", "notSet"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unavailable")},
						Computed:            true,
						MarkdownDescription: "Require Device Threat Protection minimum risk level to report noncompliance. / Device threat protection levels for the Device Threat Protection API; possible values are: `unavailable` (Default Value. Do not use.), `secured` (Device Threat Level requirement: Secured. This is the most secure level, and represents that no threats were found on the device.), `low` (Device Threat Protection level requirement: Low. Low represents a severity of threat that poses minimal risk to the device or device data.), `medium` (Device Threat Protection level requirement: Medium. Medium represents a severity of threat that poses moderate risk to the device or device data.), `high` (Device Threat Protection level requirement: High. High represents a severity of threat that poses severe risk to the device or device data.), `notSet` (Device Threat Protection level requirement: Not Set. Not set represents that there is no requirement for the device to meet a Threat Protection level.). The _provider_ default value is `\"unavailable\"`.",
					},
					"early_launch_anti_malware_driver_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require devices to be reported as healthy by Windows Device Health Attestation - early launch antimalware driver is enabled. The _provider_ default value is `false`.",
					},
					"firmware_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "When TRUE, indicates that Firmware protection is required to be reported as healthy by Microsoft Azure Attestion. When FALSE, indicates that Firmware protection is not required to be reported as healthy. Devices that support either Dynamic Root of Trust for Measurement (DRTM) or Firmware Attack Surface Reduction (FASR) will report compliant for this setting. Default value is FALSE. The _provider_ default value is `false`.",
					},
					"kernel_dma_protection_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "When TRUE, indicates that Kernel Direct Memory Access (DMA) protection is required to be reported as healthy by Microsoft Azure Attestion. When FALSE, indicates that Kernel DMA Protection is not required to be reported as healthy. Default value is FALSE. The _provider_ default value is `false`.",
					},
					"memory_integrity_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "When TRUE, indicates that Memory Integrity as known as Hypervisor-protected Code Integrity (HVCI) or Hypervisor Enforced Code Integrity protection is required to be reported as healthy by Microsoft Azure Attestion. When FALSE, indicates that Memory Integrity Protection is not required to be reported as healthy. Default value is FALSE. The _provider_ default value is `false`.",
					},
					"mobile_os_maximum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Maximum Windows Phone version.",
					},
					"mobile_os_minimum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum Windows Phone version.",
					},
					"os_maximum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Maximum Windows 10 version.",
					},
					"os_minimum_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Minimum Windows 10 version.",
					},
					"password_block_simple": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Indicates whether or not to block simple password. The _provider_ default value is `false`.",
					},
					"password_expiration_days": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The password expiration in days.",
					},
					"password_minimum_character_set_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The number of character sets required in the password.",
					},
					"password_minimum_length": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The minimum password length.",
					},
					"password_minutes_of_inactivity_before_lock": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "Minutes of inactivity before a password is required.",
					},
					"password_previous_password_block_count": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The number of previous passwords to prevent re-use of.",
					},
					"password_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require a password to unlock Windows device. The _provider_ default value is `false`.",
					},
					"password_required_to_unlock_from_idle": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require a password to unlock an idle device. The _provider_ default value is `false`.",
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
					"require_healthy_device_report": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require devices to be reported as healthy by Windows Device Health Attestation. The _provider_ default value is `false`.",
					},
					"rtp_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require Windows Defender Antimalware Real-Time Protection on Windows devices. The _provider_ default value is `false`.",
					},
					"secure_boot_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require devices to be reported as healthy by Windows Device Health Attestation - secure boot is enabled. The _provider_ default value is `false`.",
					},
					"signature_out_of_date": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require Windows Defender Antimalware Signature to be up to date on Windows devices. The _provider_ default value is `false`.",
					},
					"storage_require_encryption": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require encryption on windows devices. The _provider_ default value is `false`.",
					},
					"tpm_required": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Require Trusted Platform Module(TPM) to be present. The _provider_ default value is `false`.",
					},
					"valid_operating_system_build_ranges": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // operatingSystemVersionRange
								"description": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The description of this range (e.g. Valid 1702 builds)",
								},
								"highest_version": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The highest inclusive version that this range contains.",
								},
								"lowest_version": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The lowest inclusive version that this range contains.",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "The valid operating system build ranges on Windows devices. This collection can contain a maximum of 10000 elements. / Operating System version range. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-operatingsystemversionrange?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"virtualization_based_security_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "When TRUE, indicates that Virtualization-based Security is required to be reported as healthy by Microsoft Azure Attestion. When FALSE, indicates that Virtualization-based Security is not required to be reported as healthy. Default value is FALSE. The _provider_ default value is `false`.",
					},
					"wsl_distributions": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // wslDistributionConfiguration
								"distribution": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Linux distribution like Debian, Fedora, Ubuntu etc.",
								},
								"maximum_o_s_version": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Maximum supported operating system version of the linux version.",
								},
								"minimum_o_s_version": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Minimum supported operating system version of the linux version.",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-wsldistributionconfiguration?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
				},
				Validators: []validator.Object{
					deviceCompliancePolicyDeviceCompliancePolicyValidator,
				},
				MarkdownDescription: "This class contains compliance settings for Windows 10. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windows10compliancepolicy?view=graph-rest-beta",
			},
		},
	},
	MarkdownDescription: "This is the base class for Compliance policy. Compliance policies are platform specific and individual per-platform compliance policies inherit from here. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-devicecompliancepolicy?view=graph-rest-beta ||| MS Graph: Device configuration",
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
