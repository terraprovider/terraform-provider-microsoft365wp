package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	AndroidManagedAppProtectionResource = generic.GenericResource{
		TypeNameSuffix: "android_managed_app_protection",
		SpecificSchema: androidManagedAppProtectionResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:         "/deviceAppManagement/androidManagedAppProtections",
			ReadOptions:     androidManagedAppProtectionReadOptions,
			WriteSubActions: androidManagedAppProtectionWriteSubActions,
		},
	}

	AndroidManagedAppProtectionSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&AndroidManagedAppProtectionResource)

	AndroidManagedAppProtectionPluralDataSource = generic.CreateGenericDataSourcePluralFromSingular(
		&AndroidManagedAppProtectionSingularDataSource, "")
)

var androidManagedAppProtectionReadOptions = generic.ReadOptions{
	ODataExpand: "apps,assignments",
}

var androidManagedAppProtectionWriteSubActions = []generic.WriteSubAction{
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			Attributes: []string{"appGroupType", "apps"},
			UriSuffix:  "targetApps",
			UpdateOnly: true,
		},
	},
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			Attributes: []string{"assignments"},
			UriSuffix:  "assign",
			UpdateOnly: true,
		},
	},
}

var androidManagedAppProtectionResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // androidManagedAppProtection
		"id": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"created_date_time": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"description": schema.StringAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:      true,
		},
		"display_name": schema.StringAttribute{
			Required: true,
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:   types.StringType,
			Optional:      true,
			PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:      true,
		},
		"version": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"allowed_data_ingestion_locations": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Validators: []validator.Set{
				setvalidator.ValueStringsAre(
					stringvalidator.OneOf("oneDriveForBusiness", "sharePoint", "camera", "photoLibrary"),
				),
			},
			PlanModifiers: []planmodifier.Set{
				wpdefaultvalue.SetDefaultValue([]any{"oneDriveForBusiness", "sharePoint", "camera", "photoLibrary"}),
			},
			Computed:            true,
			MarkdownDescription: "Locations which can be used to bring data into organization documents; possible values are: `oneDriveForBusiness` (OneDrive for business), `sharePoint` (SharePoint Online), `camera` (The device's camera), `photoLibrary` (The device's photo library)",
		},
		"allowed_data_storage_locations": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Validators: []validator.Set{
				setvalidator.ValueStringsAre(
					stringvalidator.OneOf("oneDriveForBusiness", "sharePoint", "box", "localStorage", "photoLibrary"),
				),
			},
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Storage locations where managed apps can potentially store their data; possible values are: `oneDriveForBusiness` (OneDrive for business), `sharePoint` (SharePoint), `box` (Box), `localStorage` (Local storage on the device), `photoLibrary` (The device's photo library)",
		},
		"allowed_inbound_data_transfer_sources": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("allApps", "managedApps", "none"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("allApps")},
			Computed:            true,
			MarkdownDescription: "Data can be transferred from/to these classes of apps; possible values are: `allApps` (All apps.), `managedApps` (Managed apps.), `none` (No apps.)",
		},
		"allowed_outbound_clipboard_sharing_exception_length": schema.Int64Attribute{
			Optional:      true,
			PlanModifiers: []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
			Computed:      true,
		},
		"allowed_outbound_clipboard_sharing_level": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("allApps", "managedAppsWithPasteIn", "managedApps", "blocked"),
			},
			PlanModifiers: []planmodifier.String{
				wpdefaultvalue.StringDefaultValue("managedAppsWithPasteIn"),
			},
			Computed:            true,
			MarkdownDescription: "Represents the level to which the device's clipboard may be shared between apps; possible values are: `allApps` (Sharing is allowed between all apps, managed or not), `managedAppsWithPasteIn` (Sharing is allowed between all managed apps with paste in enabled), `managedApps` (Sharing is allowed between all managed apps), `blocked` (Sharing between apps is disabled)",
		},
		"allowed_outbound_data_transfer_destinations": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("allApps", "managedApps", "none"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("allApps")},
			Computed:            true,
			MarkdownDescription: "Data can be transferred from/to these classes of apps; possible values are: `allApps` (All apps.), `managedApps` (Managed apps.), `none` (No apps.)",
		},
		"app_action_if_device_compliance_required": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("block")},
			Computed:            true,
			MarkdownDescription: "An admin initiated action to be applied on a managed app; possible values are: `block` (app and the corresponding company data to be blocked), `wipe` (app and the corresponding company data to be wiped), `warn` (app and the corresponding user to be warned)",
		},
		"app_action_if_maximum_pin_retries_exceeded": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("block")},
			Computed:            true,
			MarkdownDescription: "An admin initiated action to be applied on a managed app; possible values are: `block` (app and the corresponding company data to be blocked), `wipe` (app and the corresponding company data to be wiped), `warn` (app and the corresponding user to be warned)",
		},
		"app_action_if_unable_to_authenticate_user": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
			MarkdownDescription: "An admin initiated action to be applied on a managed app; possible values are: `block` (app and the corresponding company data to be blocked), `wipe` (app and the corresponding company data to be wiped), `warn` (app and the corresponding user to be warned)",
		},
		"block_data_ingestion_into_organization_documents": schema.BoolAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:      true,
		},
		"contact_sync_blocked": schema.BoolAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:      true,
		},
		"data_backup_blocked": schema.BoolAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:      true,
		},
		"device_compliance_required": schema.BoolAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:      true,
		},
		"dialer_restriction_level": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("allApps", "managedApps", "customApp", "blocked"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("allApps")},
			Computed:            true,
			MarkdownDescription: "The classes of apps that are allowed to click-to-open a phone number, for making phone calls or sending text messages; possible values are: `allApps` (Sharing is allowed to all apps.), `managedApps` (Sharing is allowed to all managed apps.), `customApp` (Sharing is allowed to a custom app.), `blocked` (Sharing between apps is blocked.)",
		},
		"disable_app_pin_if_device_pin_is_set": schema.BoolAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:      true,
		},
		"fingerprint_blocked": schema.BoolAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:      true,
		},
		"grace_period_to_block_apps_during_off_clock_hours": schema.StringAttribute{
			Optional: true,
		},
		"managed_browser": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				wpvalidator.FlagEnumValues("notConfigured", "microsoftEdge"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
			Computed:            true,
			MarkdownDescription: "Type of managed browser; possible values are: `notConfigured` (Not configured), `microsoftEdge` (Microsoft Edge)",
		},
		"managed_browser_to_open_links_required": schema.BoolAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:      true,
		},
		"maximum_allowed_device_threat_level": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("notConfigured", "secured", "low", "medium", "high"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
			Computed:            true,
			MarkdownDescription: "The maxium threat level allowed for an app to be compliant; possible values are: `notConfigured` (Value not configured), `secured` (Device needs to have no threat), `low` (Device needs to have a low threat.), `medium` (Device needs to have not more than medium threat.), `high` (Device needs to have not more than high threat)",
		},
		"maximum_pin_retries": schema.Int64Attribute{
			Optional:      true,
			PlanModifiers: []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(5)},
			Computed:      true,
		},
		"maximum_required_os_version": schema.StringAttribute{
			Optional: true,
		},
		"maximum_warning_os_version": schema.StringAttribute{
			Optional: true,
		},
		"maximum_wipe_os_version": schema.StringAttribute{
			Optional: true,
		},
		"minimum_pin_length": schema.Int64Attribute{
			Optional:      true,
			PlanModifiers: []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(4)},
			Computed:      true,
		},
		"minimum_required_app_version": schema.StringAttribute{
			Optional: true,
		},
		"minimum_required_os_version": schema.StringAttribute{
			Optional: true,
		},
		"minimum_warning_app_version": schema.StringAttribute{
			Optional: true,
		},
		"minimum_warning_os_version": schema.StringAttribute{
			Optional: true,
		},
		"minimum_wipe_app_version": schema.StringAttribute{
			Optional: true,
		},
		"minimum_wipe_os_version": schema.StringAttribute{
			Optional: true,
		},
		"mobile_threat_defense_partner_priority": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("defenderOverThirdPartyPartner", "thirdPartyPartnerOverDefender", "unknownFutureValue"),
			},
			MarkdownDescription: "Determines the conflict resolution strategy, when more than one Mobile Threat Defense provider is enabled; possible values are: `defenderOverThirdPartyPartner` (Indicates use of Microsoft Defender Endpoint over 3rd party MTD connectors), `thirdPartyPartnerOverDefender` (Indicates use of a 3rd party MTD connector over Microsoft Defender Endpoint), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.)",
		},
		"mobile_threat_defense_remediation_action": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("block")},
			Computed:            true,
			MarkdownDescription: "An admin initiated action to be applied on a managed app; possible values are: `block` (app and the corresponding company data to be blocked), `wipe` (app and the corresponding company data to be wiped), `warn` (app and the corresponding user to be warned)",
		},
		"notification_restriction": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("allow", "blockOrganizationalData", "block"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("allow")},
			Computed:            true,
			MarkdownDescription: "Restrict managed app notification; possible values are: `allow` (Share all notifications.), `blockOrganizationalData` (Do not share Orgnizational data in notifications.), `block` (Do not share notifications.)",
		},
		"organizational_credentials_required": schema.BoolAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:      true,
		},
		"period_before_pin_reset": schema.StringAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("PT0S")},
			Computed:      true,
		},
		"period_offline_before_access_check": schema.StringAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("PT12H")},
			Computed:      true,
		},
		"period_offline_before_wipe_is_enforced": schema.StringAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("P90D")},
			Computed:      true,
		},
		"period_online_before_access_check": schema.StringAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("PT30M")},
			Computed:      true,
		},
		"pin_character_set": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("numeric", "alphanumericAndSymbol"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("numeric")},
			Computed:            true,
			MarkdownDescription: "Character set which is to be used for a user's app PIN; possible values are: `numeric` (Numeric characters), `alphanumericAndSymbol` (Alphanumeric and symbolic characters)",
		},
		"pin_required": schema.BoolAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:      true,
		},
		"pin_required_instead_of_biometric_timeout": schema.StringAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("PT30M")},
			Computed:      true,
		},
		"previous_pin_block_count": schema.Int64Attribute{
			Optional:      true,
			PlanModifiers: []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
			Computed:      true,
		},
		"print_blocked": schema.BoolAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:      true,
		},
		"protected_messaging_redirect_app_type": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("anyApp", "anyManagedApp", "specificApps", "blocked"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("anyApp")},
			Computed:            true,
			MarkdownDescription: "Defines how app messaging redirection is protected by an App Protection Policy. Default is anyApp; possible values are: `anyApp` (App protection policy will allow messaging redirection to any app.), `anyManagedApp` (App protection policy will allow messaging redirection to any managed application.), `specificApps` (App protection policy will allow messaging redirection only to specified applications in related App protection policy settings. See related settings 'messagingRedirectAppDisplayName', 'messagingRedirectAppPackageId' and 'messagingRedirectAppUrlScheme'.), `blocked` (App protection policy will block messaging redirection to any app.)",
		},
		"save_as_blocked": schema.BoolAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:      true,
		},
		"simple_pin_blocked": schema.BoolAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:      true,
		},
		"app_group_type": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf("selectedPublicApps", "allCoreMicrosoftApps", "allMicrosoftApps", "allApps"),
			},
			MarkdownDescription: "Public Apps selection: group or individual / Indicates a collection of apps to target which can be one of several pre-defined lists of apps or a manually selected list of apps; possible values are: `selectedPublicApps` (Target the collection of apps manually selected by the admin.), `allCoreMicrosoftApps` (Target the core set of Microsoft apps (Office, Edge, etc).), `allMicrosoftApps` (Target all apps with Microsoft as publisher.), `allApps` (Target all apps with an available assignment.)",
		},
		"is_assigned": schema.BoolAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
			MarkdownDescription: "Indicates if the policy is deployed to any inclusion groups or not.",
		},
		"targeted_app_management_levels": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				wpvalidator.FlagEnumValues("unspecified", "unmanaged", "mdm", "androidEnterprise", "androidEnterpriseDedicatedDevicesWithAzureAdSharedMode", "androidOpenSourceProjectUserAssociated", "androidOpenSourceProjectUserless", "unknownFutureValue"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unspecified")},
			Computed:            true,
			MarkdownDescription: "The intended app management levels for this policy / Management levels for apps; possible values are: `unspecified` (Unspecified), `unmanaged` (Unmanaged), `mdm` (MDM), `androidEnterprise` (Android Enterprise), `androidEnterpriseDedicatedDevicesWithAzureAdSharedMode` (Android Enterprise dedicated devices with Azure AD Shared mode), `androidOpenSourceProjectUserAssociated` (Android Open Source Project (AOSP) devices), `androidOpenSourceProjectUserless` (Android Open Source Project (AOSP) userless devices), `unknownFutureValue` (Place holder for evolvable enum)",
		},
		"assignments": deviceAndAppManagementAssignment,
		"allowed_android_device_manufacturers": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Semicolon seperated list of device manufacturers allowed, as a string, for the managed app to work.",
		},
		"allowed_android_device_models": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "List of device models allowed, as a string, for the managed app to work.",
		},
		"app_action_if_account_is_clocked_out": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
			MarkdownDescription: "Defines a managed app behavior, either block or warn, if the user is clocked out (non-working time). / An admin initiated action to be applied on a managed app; possible values are: `block` (app and the corresponding company data to be blocked), `wipe` (app and the corresponding company data to be wiped), `warn` (app and the corresponding user to be warned)",
		},
		"app_action_if_android_device_manufacturer_not_allowed": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("block")},
			Computed:            true,
			MarkdownDescription: "Defines a managed app behavior, either block or wipe, if the specified device manufacturer is not allowed. / An admin initiated action to be applied on a managed app; possible values are: `block` (app and the corresponding company data to be blocked), `wipe` (app and the corresponding company data to be wiped), `warn` (app and the corresponding user to be warned)",
		},
		"app_action_if_android_device_model_not_allowed": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("block")},
			Computed:            true,
			MarkdownDescription: "Defines a managed app behavior, either block or wipe, if the specified device model is not allowed. / An admin initiated action to be applied on a managed app; possible values are: `block` (app and the corresponding company data to be blocked), `wipe` (app and the corresponding company data to be wiped), `warn` (app and the corresponding user to be warned)",
		},
		"app_action_if_android_safety_net_apps_verification_failed": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("block")},
			Computed:            true,
			MarkdownDescription: "Defines a managed app behavior, either warn or block, if the specified Android App Verification requirement fails. / An admin initiated action to be applied on a managed app; possible values are: `block` (app and the corresponding company data to be blocked), `wipe` (app and the corresponding company data to be wiped), `warn` (app and the corresponding user to be warned)",
		},
		"app_action_if_android_safety_net_device_attestation_failed": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("block")},
			Computed:            true,
			MarkdownDescription: "Defines a managed app behavior, either warn or block, if the specified Android SafetyNet Attestation requirement fails. / An admin initiated action to be applied on a managed app; possible values are: `block` (app and the corresponding company data to be blocked), `wipe` (app and the corresponding company data to be wiped), `warn` (app and the corresponding user to be warned)",
		},
		"app_action_if_device_lock_not_set": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("block")},
			Computed:            true,
			MarkdownDescription: "Defines a managed app behavior, either warn, block or wipe, if the screen lock is required on android device but is not set. / An admin initiated action to be applied on a managed app; possible values are: `block` (app and the corresponding company data to be blocked), `wipe` (app and the corresponding company data to be wiped), `warn` (app and the corresponding user to be warned)",
		},
		"app_action_if_device_passcode_complexity_less_than_high": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
			MarkdownDescription: "If the device does not have a passcode of high complexity or higher, trigger the stored action. / An admin initiated action to be applied on a managed app; possible values are: `block` (app and the corresponding company data to be blocked), `wipe` (app and the corresponding company data to be wiped), `warn` (app and the corresponding user to be warned)",
		},
		"app_action_if_device_passcode_complexity_less_than_low": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
			MarkdownDescription: "If the device does not have a passcode of low complexity or higher, trigger the stored action. / An admin initiated action to be applied on a managed app; possible values are: `block` (app and the corresponding company data to be blocked), `wipe` (app and the corresponding company data to be wiped), `warn` (app and the corresponding user to be warned)",
		},
		"app_action_if_device_passcode_complexity_less_than_medium": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
			MarkdownDescription: "If the device does not have a passcode of medium complexity or higher, trigger the stored action. / An admin initiated action to be applied on a managed app; possible values are: `block` (app and the corresponding company data to be blocked), `wipe` (app and the corresponding company data to be wiped), `warn` (app and the corresponding user to be warned)",
		},
		"app_action_if_samsung_knox_attestation_required": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
			MarkdownDescription: "Defines the behavior of a managed app when Samsung Knox Attestation is required. Possible values are null, warn, block & wipe. If the admin does not set this action, the default is null, which indicates this setting is not configured. / An admin initiated action to be applied on a managed app; possible values are: `block` (app and the corresponding company data to be blocked), `wipe` (app and the corresponding company data to be wiped), `warn` (app and the corresponding user to be warned)",
		},
		"approved_keyboards": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // keyValuePair
					"name": schema.StringAttribute{
						Required: true,
					},
					"value": schema.StringAttribute{
						Optional: true,
					},
				},
			},
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "If Keyboard Restriction is enabled, only keyboards in this approved list will be allowed. A key should be Android package id for a keyboard and value should be a friendly name / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyValuePair?view=graph-rest-beta",
		},
		"biometric_authentication_blocked": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Indicates whether use of the biometric authentication is allowed in place of a pin if PinRequired is set to True.",
		},
		"block_after_company_portal_update_deferral_in_days": schema.Int64Attribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
			Computed:            true,
			MarkdownDescription: "Maximum number of days Company Portal update can be deferred on the device or app access will be blocked.",
		},
		"connect_to_vpn_on_launch": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Whether the app should connect to the configured VPN on launch.",
		},
		"custom_browser_display_name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Friendly name of the preferred custom browser to open weblink on Android. When this property is configured, ManagedBrowserToOpenLinksRequired should be true.",
		},
		"custom_browser_package_id": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Unique identifier of the preferred custom browser to open weblink on Android. When this property is configured, ManagedBrowserToOpenLinksRequired should be true.",
		},
		"custom_dialer_app_display_name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Friendly name of a custom dialer app to click-to-open a phone number on Android.",
		},
		"custom_dialer_app_package_id": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "PackageId of a custom dialer app to click-to-open a phone number on Android.",
		},
		"deployed_app_count": schema.Int64Attribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
			MarkdownDescription: "Count of apps to which the current policy is deployed.",
		},
		"device_lock_required": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Defines if any kind of lock must be required on android device",
		},
		"disable_app_encryption_if_device_encryption_is_enabled": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "When this setting is enabled, app level encryption is disabled if device level encryption is enabled",
		},
		"encrypt_app_data": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:            true,
			MarkdownDescription: "Indicates whether application data for managed apps should be encrypted",
		},
		"exempted_app_packages": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // keyValuePair
					"name": schema.StringAttribute{
						Required: true,
					},
					"value": schema.StringAttribute{
						Optional: true,
					},
				},
			},
			PlanModifiers: []planmodifier.Set{
				wpdefaultvalue.SetDefaultValueEmpty(),
				&wpplanmodifier.IgnoreOnOtherAttributeValuePlanModifier{OtherAttributePath: path.Root("allowed_outbound_data_transfer_destinations"), ValuesRespect: []attr.Value{types.StringValue("managedApps")}},
			},
			Computed:            true,
			MarkdownDescription: "App packages in this list will be exempt from the policy and will be able to receive data from managed apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyValuePair?view=graph-rest-beta",
		},
		"fingerprint_and_biometric_enabled": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:            true,
			MarkdownDescription: "If null, this setting will be ignored. If false both fingerprints and biometrics will not be enabled. If true, both fingerprints and biometrics will be enabled.",
		},
		"keyboards_restricted": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Indicates if keyboard restriction is enabled. If enabled list of approved keyboards must be provided as well.",
		},
		"messaging_redirect_app_display_name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "When a specific app redirection is enforced by protectedMessagingRedirectAppType in an App Protection Policy, this value defines the app name which is allowed to be used.",
		},
		"messaging_redirect_app_package_id": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "When a specific app redirection is enforced by protectedMessagingRedirectAppType in an App Protection Policy, this value defines the app package id which is allowed to be used.",
		},
		"minimum_required_company_portal_version": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Minimum version of the Company portal that must be installed on the device or app access will be blocked",
		},
		"minimum_required_patch_version": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("0000-00-00")},
			Computed:            true,
			MarkdownDescription: "Define the oldest required Android security patch level a user can have to gain secure access to the app.",
		},
		"minimum_warning_company_portal_version": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Minimum version of the Company portal that must be installed on the device or the user will receive a warning",
		},
		"minimum_warning_patch_version": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("0000-00-00")},
			Computed:            true,
			MarkdownDescription: "Define the oldest recommended Android security patch level a user can have for secure access to the app.",
		},
		"minimum_wipe_company_portal_version": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Minimum version of the Company portal that must be installed on the device or the company data on the app will be wiped",
		},
		"minimum_wipe_patch_version": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("0000-00-00")},
			Computed:            true,
			MarkdownDescription: "Android security patch level  less than or equal to the specified value will wipe the managed app and the associated company data.",
		},
		"require_class3_biometrics": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			Description:         `requireClass3Biometrics`, // custom MS Graph attribute name
			MarkdownDescription: "Require user to apply Class 3 Biometrics on their Android device.",
		},
		"required_android_safety_net_apps_verification_type": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("none", "enabled")},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
			Computed:            true,
			MarkdownDescription: "Defines the Android SafetyNet Apps Verification requirement for a managed app to work. / An admin enforced Android SafetyNet Device Attestation requirement on a managed app; possible values are: `none` (no requirement set), `enabled` (require that Android device has SafetyNet Apps Verification enabled)",
		},
		"required_android_safety_net_device_attestation_type": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("none", "basicIntegrity", "basicIntegrityAndDeviceCertification"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
			Computed:            true,
			MarkdownDescription: "Defines the Android SafetyNet Device Attestation requirement for a managed app to work. / An admin enforced Android SafetyNet Device Attestation requirement on a managed app; possible values are: `none` (no requirement set), `basicIntegrity` (require that Android device passes SafetyNet Basic Integrity validation), `basicIntegrityAndDeviceCertification` (require that Android device passes SafetyNet Basic Integrity and Device Certification validations)",
		},
		"required_android_safety_net_evaluation_type": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("basic", "hardwareBacked")},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("basic")},
			Computed:            true,
			MarkdownDescription: "Defines the Android SafetyNet evaluation type requirement for a managed app to work. / An admin enforced Android SafetyNet evaluation type requirement on a managed app; possible values are: `basic` (Require basic evaluation), `hardwareBacked` (Require hardware backed evaluation)",
		},
		"require_pin_after_biometric_change": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "A PIN prompt will override biometric prompts if class 3 biometrics are updated on the device.",
		},
		"screen_capture_blocked": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Indicates whether a managed user can take screen captures of managed apps",
		},
		"warn_after_company_portal_update_deferral_in_days": schema.Int64Attribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
			Computed:            true,
			MarkdownDescription: "Maximum number of days Company Portal update can be deferred on the device or the user will receive the warning",
		},
		"wipe_after_company_portal_update_deferral_in_days": schema.Int64Attribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
			Computed:            true,
			MarkdownDescription: "Maximum number of days Company Portal update can be deferred on the device or the company data on the app will be wiped",
		},
		"apps": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // managedMobileApp
					"app_id": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{ // mobileAppIdentifier
							"android": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.androidMobileAppIdentifier",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // androidMobileAppIdentifier
										"package_id": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "The identifier for an app, as specified in the play store.",
										},
									},
									Validators: []validator.Object{
										androidManagedAppProtectionMobileAppIdentifierValidator,
									},
									MarkdownDescription: "The identifier for an Android app.",
								},
							},
						},
						Description:         `mobileAppIdentifier`, // custom MS Graph attribute name
						MarkdownDescription: "The identifier for an app with it's operating system type. / The identifier for a mobile app.",
					},
				},
			},
			PlanModifiers: []planmodifier.Set{
				&wpplanmodifier.IgnoreOnOtherAttributeValuePlanModifier{OtherAttributePath: path.Root("app_group_type"), ValuesRespect: []attr.Value{types.StringValue("selectedPublicApps")}},
			},
			Computed:            true,
			MarkdownDescription: "List of apps to which the policy is deployed. / The identifier for the deployment an app.",
		},
	},
	MarkdownDescription: "Policy used to configure detailed management settings targeted to specific security groups and for a specified set of apps on an Android device / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-androidManagedAppProtection?view=graph-rest-beta",
}

var androidManagedAppProtectionMobileAppIdentifierValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("android"),
)
