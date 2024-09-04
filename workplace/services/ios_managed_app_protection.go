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
	IosManagedAppProtectionResource = generic.GenericResource{
		TypeNameSuffix: "ios_managed_app_protection",
		SpecificSchema: iosManagedAppProtectionResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:         "/deviceAppManagement/iosManagedAppProtections",
			ReadOptions:     iosManagedAppProtectionReadOptions,
			WriteSubActions: iosManagedAppProtectionWriteSubActions,
		},
	}

	IosManagedAppProtectionSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&IosManagedAppProtectionResource)

	IosManagedAppProtectionPluralDataSource = generic.CreateGenericDataSourcePluralFromSingular(
		&IosManagedAppProtectionSingularDataSource, "")
)

var iosManagedAppProtectionReadOptions = generic.ReadOptions{
	ODataExpand: "apps,assignments",
}

var iosManagedAppProtectionWriteSubActions = []generic.WriteSubAction{
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

var iosManagedAppProtectionResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // iosManagedAppProtection
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
			Computed: true,
		},
		"allowed_data_storage_locations": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Validators: []validator.Set{
				setvalidator.ValueStringsAre(
					stringvalidator.OneOf("oneDriveForBusiness", "sharePoint", "box", "localStorage", "photoLibrary"),
				),
			},
			PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:      true,
		},
		"allowed_inbound_data_transfer_sources": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("allApps", "managedApps", "none"),
			},
			PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("allApps")},
			Computed:      true,
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
			Computed: true,
		},
		"allowed_outbound_data_transfer_destinations": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("allApps", "managedApps", "none"),
			},
			PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("allApps")},
			Computed:      true,
		},
		"app_action_if_device_compliance_required": schema.StringAttribute{
			Optional:      true,
			Validators:    []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
			PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("block")},
			Computed:      true,
		},
		"app_action_if_maximum_pin_retries_exceeded": schema.StringAttribute{
			Optional:      true,
			Validators:    []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
			PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("block")},
			Computed:      true,
		},
		"app_action_if_unable_to_authenticate_user": schema.StringAttribute{
			Optional:   true,
			Validators: []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
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
			PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("allApps")},
			Computed:      true,
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
			PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
			Computed:      true,
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
			PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
			Computed:      true,
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
		},
		"mobile_threat_defense_remediation_action": schema.StringAttribute{
			Optional:      true,
			Validators:    []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
			PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("block")},
			Computed:      true,
		},
		"notification_restriction": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("allow", "blockOrganizationalData", "block"),
			},
			PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("allow")},
			Computed:      true,
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
			PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("numeric")},
			Computed:      true,
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
			PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("anyApp")},
			Computed:      true,
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
			MarkdownDescription: `Public Apps selection: group or individual`,
		},
		"is_assigned": schema.BoolAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
			MarkdownDescription: `Indicates if the policy is deployed to any inclusion groups or not.`,
		},
		"targeted_app_management_levels": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				wpvalidator.FlagEnumValues("unspecified", "unmanaged", "mdm", "androidEnterprise", "androidEnterpriseDedicatedDevicesWithAzureAdSharedMode", "androidOpenSourceProjectUserAssociated", "androidOpenSourceProjectUserless", "unknownFutureValue"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unspecified")},
			Computed:            true,
			MarkdownDescription: `The intended app management levels for this policy`,
		},
		"assignments": deviceAndAppManagementAssignment,
		"allowed_ios_device_models": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: `Semicolon seperated list of device models allowed, as a string, for the managed app to work.`,
		},
		"allow_widget_content_sync": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: `Indicates  if content sync for widgets is allowed for iOS on App Protection Policies`,
		},
		"app_action_if_account_is_clocked_out": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("block")},
			Computed:            true,
			MarkdownDescription: `Defines a managed app behavior, either block or warn, if the user is clocked out (non-working time).`,
		},
		"app_action_if_ios_device_model_not_allowed": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("block", "wipe", "warn")},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("block")},
			Computed:            true,
			MarkdownDescription: `Defines a managed app behavior, either block or wipe, if the specified device model is not allowed.`,
		},
		"app_data_encryption_type": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("useDeviceSettings", "afterDeviceRestart", "whenDeviceLockedExceptOpenFiles", "whenDeviceLocked"),
			},
			PlanModifiers: []planmodifier.String{
				wpdefaultvalue.StringDefaultValue("whenDeviceLocked"),
			},
			Computed:            true,
			MarkdownDescription: `Type of encryption which should be used for data in a managed app.`,
		},
		"custom_browser_protocol": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: `A custom browser protocol to open weblink on iOS. When this property is configured, ManagedBrowserToOpenLinksRequired should be true.`,
		},
		"custom_dialer_app_protocol": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: `Protocol of a custom dialer app to click-to-open a phone number on iOS, for example, skype:.`,
		},
		"deployed_app_count": schema.Int64Attribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
			MarkdownDescription: `Count of apps to which the current policy is deployed.`,
		},
		"disable_protection_of_managed_outbound_open_in_data": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: `Disable protection of data transferred to other apps through IOS OpenIn option. This setting is only allowed to be True when AllowedOutboundDataTransferDestinations is set to ManagedApps.`,
		},
		"exempted_app_protocols": schema.SetNestedAttribute{
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
			MarkdownDescription: `Apps in this list will be exempt from the policy and will be able to receive data from managed apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyValuePair?view=graph-rest-beta`,
		},
		"exempted_universal_links": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: `A list of custom urls that are allowed to invocate an unmanaged app`,
		},
		"face_id_blocked": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: `Indicates whether use of the FaceID is allowed in place of a pin if PinRequired is set to True.`,
		},
		"filter_open_in_to_only_managed_apps": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: `Defines if open-in operation is supported from the managed app to the filesharing locations selected. This setting only applies when AllowedOutboundDataTransferDestinations is set to ManagedApps and DisableProtectionOfManagedOutboundOpenInData is set to False.`,
		},
		"managed_universal_links": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: `A list of custom urls that are allowed to invocate a managed app`,
		},
		"messaging_redirect_app_url_scheme": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: `When a specific app redirection is enforced by protectedMessagingRedirectAppType in an App Protection Policy, this value defines the app url redirect schemes which are allowed to be used.`,
		},
		"minimum_required_sdk_version": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: `Versions less than the specified version will block the managed app from accessing company data.`,
		},
		"minimum_warning_sdk_version": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: `Versions less than the specified version will result in warning message on the managed app from accessing company data.`,
		},
		"minimum_wipe_sdk_version": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: `Versions less than the specified version will block the managed app from accessing company data.`,
		},
		"protect_inbound_data_from_unknown_sources": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: `Protect incoming data from unknown source. This setting is only allowed to be True when AllowedInboundDataTransferSources is set to AllApps.`,
		},
		"third_party_keyboards_blocked": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: `Defines if third party keyboards are allowed while accessing a managed app`,
		},
		"apps": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // managedMobileApp
					"app_id": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{ // mobileAppIdentifier
							"ios": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.iosMobileAppIdentifier",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // iosMobileAppIdentifier
										"bundle_id": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: `The identifier for an app, as specified in the app store.`,
										},
									},
									Validators: []validator.Object{
										iosManagedAppProtectionMobileAppIdentifierValidator,
									},
									MarkdownDescription: `The identifier for an iOS app.`,
								},
							},
						},
						Description:         `mobileAppIdentifier`, // custom MS Graph attribute name
						MarkdownDescription: `The identifier for an app with it's operating system type. / The identifier for a mobile app.`,
					},
				},
			},
			PlanModifiers: []planmodifier.Set{
				wpdefaultvalue.SetDefaultValueEmpty(),
				&wpplanmodifier.IgnoreOnOtherAttributeValuePlanModifier{OtherAttributePath: path.Root("app_group_type"), ValuesRespect: []attr.Value{types.StringValue("selectedPublicApps")}},
			},
			Computed:            true,
			MarkdownDescription: `List of apps to which the policy is deployed. / The identifier for the deployment an app.`,
		},
	},
	MarkdownDescription: `Policy used to configure detailed management settings targeted to specific security groups and for a specified set of apps on an iOS device / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-iosManagedAppProtection?view=graph-rest-beta`,
}

var iosManagedAppProtectionMobileAppIdentifierValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("ios"),
)
