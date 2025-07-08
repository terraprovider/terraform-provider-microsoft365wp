package services

import (
	"context"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	IntuneBrandingProfileResource = generic.GenericResource{
		TypeNameSuffix: "intune_branding_profile",
		SpecificSchema: intuneBrandingProfileResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/deviceManagement/intuneBrandingProfiles",
			ReadOptions: generic.ReadOptions{
				ODataExpand: "assignments",
				ExtraRequests: []generic.ReadExtraRequest{
					{
						Attribute: "landingPageCustomizedImage",
					},
					{
						Attribute: "lightBackgroundLogo",
					},
					{
						Attribute: "themeColorLogo",
					},
				},
				DataSource: generic.DataSourceOptions{
					ExtraFilterAttributes: []string{"display_name", "is_default_profile", "profile_name"},
					Plural: generic.PluralOptions{
						ExtraAttributes: []string{"is_default_profile", "profile_name"},
					},
				},
			},
			WriteOptions: generic.WriteOptions{
				SubActions: []generic.WriteSubAction{
					&generic.WriteSubActionPostAndPatch{
						AttributesPatchExclude: []string{"assignments"},
						AttributesKeep:         []string{"profileName", "profileDescription", "roleScopeTagIds", "assignments"},
					},
					&generic.WriteSubActionAllInOne{
						WriteSubActionBase: generic.WriteSubActionBase{
							Attributes: []string{"assignments"},
							UriSuffix:  "assign",
						},
					},
				},
			},
			TerraformToGraphMiddleware: intuneBrandingProfileTerraformToGraphMiddleware,
			GraphToTerraformMiddleware: intuneBrandingProfileGraphToTerraformMiddleware,
			CreateModifyFunc:           intuneBrandingProfileCreateModifyFunc,
			DeleteModifyFunc:           intuneBrandingProfileDeleteModifyFunc,
		},
	}

	IntuneBrandingProfileSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&IntuneBrandingProfileResource)

	IntuneBrandingProfilePluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&IntuneBrandingProfileResource, "")
)

var kIbpDefaultProfileName = "###TfWorkplaceDefault###"

func intuneBrandingProfileTerraformToGraphMiddleware(ctx context.Context, diags *diag.Diagnostics, params *generic.TerraformToGraphMiddlewareParams) generic.TerraformToGraphMiddlewareReturns {
	// We cannot easily use isDefaultProfile here since it's read-only in TF. Therefore we use profile_name here.
	profileName, profileNameOk := params.RawVal["profileName"].(string)
	if profileNameOk && profileName == kIbpDefaultProfileName {
		// remove unchangeable attributes
		delete(params.RawVal, "profileName")
		delete(params.RawVal, "profileDescription")
		delete(params.RawVal, "assignments")
	}
	return nil
}

func intuneBrandingProfileGraphToTerraformMiddleware(ctx context.Context, diags *diag.Diagnostics, params *generic.GraphToTerraformMiddlewareParams) generic.GraphToTerraformMiddlewareReturns {
	// When reading from MS Graph, profileName cannot be our marker to assert that it's the default profile, so we assert using isDefaultProfile
	isDefaultProfile, isDefaultProfileOk := params.RawVal["isDefaultProfile"].(bool)
	if isDefaultProfileOk && isDefaultProfile {
		// set unchangeable attributes to convention values
		params.RawVal["profileName"] = kIbpDefaultProfileName
		params.RawVal["profileDescription"] = "" // schema default value
		params.RawVal["assignments"] = []any{}   // schema default value
	}
	return nil
}

func intuneBrandingProfileCreateModifyFunc(ctx context.Context, diags *diag.Diagnostics, params *generic.CreateModifyFuncParams) {
	var profileName string
	params.Req.Plan.GetAttribute(ctx, path.Root("profile_name"), &profileName)
	if profileName == kIbpDefaultProfileName {
		// Apparently the default profile is targeted, so try to determine the id of the existing entity for it to get
		// updated instead of a new entity to get created.
		odataFilter := "isDefaultProfile eq true"
		defaultProfileId := params.R.AccessParams.ReadId(ctx, diags, "", nil, true, odataFilter, nil)
		if diags.HasError() {
			return
		}
		params.UpdateExisting = true
		params.Id = defaultProfileId
	}
}

func intuneBrandingProfileDeleteModifyFunc(ctx context.Context, diags *diag.Diagnostics, params *generic.DeleteModifyFuncParams) {
	// If it's the default profile we must not delete it (it would probably fail anyway) but just skip deletion instead.
	var profileName string
	params.Req.State.GetAttribute(ctx, path.Root("profile_name"), &profileName)
	if profileName == kIbpDefaultProfileName {
		params.SkipDelete = true
	}
}

var intuneBrandingProfileResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // intuneBrandingProfile
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Profile Key",
		},
		"company_portal_blocked_actions": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // companyPortalBlockedAction
					"action": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("unknown", "remove", "reset"),
						},
						MarkdownDescription: "Device Action. / Action on a device that can be executed in the Company Portal; possible values are: `unknown` (Unknown device action), `remove` (Remove device from Company Portal), `reset` (Reset device enrolled in Company Portal)",
					},
					"owner_type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("unknown", "company", "personal"),
						},
						MarkdownDescription: "Device ownership type. / Owner type of device; possible values are: `unknown` (Unknown.), `company` (Owned by company.), `personal` (Owned by person.)",
					},
					"platform": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("android", "androidForWork", "iOS", "macOS", "windowsPhone81", "windows81AndLater", "windows10AndLater", "androidWorkProfile", "unknown", "androidAOSP", "androidMobileApplicationManagement", "iOSMobileApplicationManagement", "unknownFutureValue", "windowsMobileApplicationManagement"),
						},
						MarkdownDescription: "Device OS/Platform. / Supported platform types; possible values are: `android` (Indicates device platform type is android.), `androidForWork` (Indicates device platform type is android for work.), `iOS` (Indicates device platform type is iOS.), `macOS` (Indicates device platform type is macOS.), `windowsPhone81` (Indicates device platform type is WindowsPhone 8.1.), `windows81AndLater` (Indicates device platform type is Windows 8.1 and later.), `windows10AndLater` (Indicates device platform type is Windows 10 and later.), `androidWorkProfile` (Indicates device platform type is Android Work Profile.), `unknown` (This is the default value when device platform type resolution fails), `androidAOSP` (Indicates device platform type is Android AOSP.), `androidMobileApplicationManagement` (Indicates device platform type associated with Mobile Application Management (MAM) for android devices.), `iOSMobileApplicationManagement` (Indicates device platform type associated with Mobile Application Management (MAM) for iOS devices.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.), `windowsMobileApplicationManagement` (Indicates device platform type associated with Mobile Application Management (MAM) for Windows devices.)",
					},
				},
			},
			Default: setdefault.StaticValue(types.SetValueMust(
				types.ObjectType{AttrTypes: map[string]attr.Type{"action": types.StringType, "owner_type": types.StringType, "platform": types.StringType}},
				[]attr.Value{types.ObjectValueMust(
					map[string]attr.Type{"action": types.StringType, "owner_type": types.StringType, "platform": types.StringType},
					map[string]attr.Value{"action": types.StringValue("remove"), "owner_type": types.StringValue("company"), "platform": types.StringValue("windows10AndLater")},
				)},
			)),
			Computed:            true,
			MarkdownDescription: "Collection of blocked actions on the company portal as per platform and device ownership types. / Blocked actions on the company portal as per platform and device ownership types / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-companyportalblockedaction?view=graph-rest-beta  \nProvider Note: The _provider_ default value is `[{\"action\":\"remove,\"owner_type\":\"company\",\"platform\":\"windows10AndLater\"}]`.",
		},
		"contact_it_email_address": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			Description:         `contactITEmailAddress`, // custom MS Graph attribute name
			MarkdownDescription: "E-mail address of the person/organization responsible for IT support. The _provider_ default value is `\"\"`.",
		},
		"contact_it_name": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			Description:         `contactITName`, // custom MS Graph attribute name
			MarkdownDescription: "Name of the person/organization responsible for IT support. The _provider_ default value is `\"\"`.",
		},
		"contact_it_notes": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			Description:         `contactITNotes`, // custom MS Graph attribute name
			MarkdownDescription: "Text comments regarding the person/organization responsible for IT support. The _provider_ default value is `\"\"`.",
		},
		"contact_it_phone_number": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			Description:         `contactITPhoneNumber`, // custom MS Graph attribute name
			MarkdownDescription: "Phone number of the person/organization responsible for IT support. The _provider_ default value is `\"\"`.",
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Time when the BrandingProfile was created",
		},
		"custom_can_see_privacy_message": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "Text comments regarding what the admin has access to on the device. The _provider_ default value is `\"\"`.",
		},
		"custom_cant_see_privacy_message": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "Text comments regarding what the admin doesn't have access to on the device. The _provider_ default value is `\"\"`.",
		},
		"disable_device_category_selection": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Boolean that indicates if Device Category Selection will be shown in Company Portal. The _provider_ default value is `false`.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Company/organization name that is displayed to end users",
		},
		"enrollment_availability": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("availableWithPrompts", "availableWithoutPrompts", "unavailable"),
			},
			PlanModifiers: []planmodifier.String{
				wpdefaultvalue.StringDefaultValue("availableWithPrompts"),
			},
			Computed:            true,
			MarkdownDescription: "Customized device enrollment flow displayed to the end user . / Options available for enrollment flow customization; possible values are: `availableWithPrompts` (Device enrollment flow is shown to the end user with guided enrollment prompts), `availableWithoutPrompts` (Device enrollment flow is available to the end user without guided enrollment prompts), `unavailable` (Device enrollment flow is unavailable to the enduser). The _provider_ default value is `\"availableWithPrompts\"`.",
		},
		"is_default_profile": schema.BoolAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
			MarkdownDescription: "Boolean that represents whether the profile is used as default or not",
		},
		"is_factory_reset_disabled": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Boolean that represents whether the adminsistrator has disabled the 'Factory Reset' action on corporate owned devices. The _provider_ default value is `false`.",
		},
		"is_remove_device_disabled": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Boolean that represents whether the adminsistrator has disabled the 'Remove Device' action on corporate owned devices. The _provider_ default value is `false`.",
		},
		"landing_page_customized_image": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // mimeContent
				"type": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Indicates the content mime type.",
				},
				"value_base64": schema.StringAttribute{
					Required:            true,
					Description:         `value`, // custom MS Graph attribute name
					MarkdownDescription: "The byte array that contains the actual content.",
				},
			},
			MarkdownDescription: "Customized image displayed in Company Portal apps landing page / Contains properties for a generic mime content. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mimecontent?view=graph-rest-beta",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Time when the BrandingProfile was last modified",
		},
		"light_background_logo": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // mimeContent
				"type": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Indicates the content mime type.",
				},
				"value_base64": schema.StringAttribute{
					Required:            true,
					Description:         `value`, // custom MS Graph attribute name
					MarkdownDescription: "The byte array that contains the actual content.",
				},
			},
			MarkdownDescription: "Logo image displayed in Company Portal apps which have a light background behind the logo / Contains properties for a generic mime content. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mimecontent?view=graph-rest-beta",
		},
		"online_support_site_name": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "Display name of the company/organization’s IT helpdesk site. The _provider_ default value is `\"\"`.",
		},
		"online_support_site_url": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "URL to the company/organization’s IT helpdesk site. The _provider_ default value is `\"\"`.",
		},
		"privacy_url": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "URL to the company/organization’s privacy policy",
		},
		"profile_description": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "Description of the profile. The _provider_ default value is `\"\"`.",
		},
		"profile_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Name of the profile",
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: "List of scope tags assigned to the branding profile. The _provider_ default value is `[\"0\"]`.",
		},
		"send_device_ownership_change_push_notification": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:            true,
			MarkdownDescription: "Boolean that indicates if a push notification is sent to users when their device ownership type changes from personal to corporate. The _provider_ default value is `true`.",
		},
		"show_azure_ad_enterprise_apps": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			Description:         `showAzureADEnterpriseApps`, // custom MS Graph attribute name
			MarkdownDescription: "Boolean that indicates if AzureAD Enterprise Apps will be shown in Company Portal. The _provider_ default value is `false`.",
		},
		"show_configuration_manager_apps": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:            true,
			MarkdownDescription: "Boolean that indicates if Configuration Manager Apps will be shown in Company Portal. The _provider_ default value is `true`.",
		},
		"show_display_name_next_to_logo": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:            true,
			MarkdownDescription: "Boolean that represents whether the administrator-supplied display name will be shown next to the logo image or not. The _provider_ default value is `true`.",
		},
		"show_logo": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Boolean that represents whether the administrator-supplied logo images are shown or not. The _provider_ default value is `false`.",
		},
		"show_office_web_apps": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Boolean that indicates if Office WebApps will be shown in Company Portal. The _provider_ default value is `false`.",
		},
		"theme_color": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // rgbColor
				"b": schema.Int64Attribute{
					Required:            true,
					MarkdownDescription: "Blue value",
				},
				"g": schema.Int64Attribute{
					Required:            true,
					MarkdownDescription: "Green value",
				},
				"r": schema.Int64Attribute{
					Required:            true,
					MarkdownDescription: "Red value",
				},
			},
			Default: objectdefault.StaticValue(types.ObjectValueMust(
				map[string]attr.Type{"r": types.Int64Type, "g": types.Int64Type, "b": types.Int64Type},
				map[string]attr.Value{"r": types.Int64Value(0), "g": types.Int64Value(114), "b": types.Int64Value(198)},
			)),
			Computed:            true,
			MarkdownDescription: "Primary theme color used in the Company Portal applications and web portal / Color in RGB. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-rgbcolor?view=graph-rest-beta  \nProvider Note: The _provider_ default value is `{\"r\":0,\"g\":114,\"b\":198}`.",
		},
		"theme_color_logo": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // mimeContent
				"type": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Indicates the content mime type.",
				},
				"value_base64": schema.StringAttribute{
					Required:            true,
					Description:         `value`, // custom MS Graph attribute name
					MarkdownDescription: "The byte array that contains the actual content.",
				},
			},
			MarkdownDescription: "Logo image displayed in Company Portal apps which have a theme color background behind the logo / Contains properties for a generic mime content. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mimecontent?view=graph-rest-beta",
		},
		"assignments": deviceAndAppManagementAssignment,
	},
	MarkdownDescription: "This entity contains data which is used in customizing the tenant level appearance of the Company Portal applications as well as the end user web portal. / https://learn.microsoft.com/en-us/graph/api/resources/intune-wip-intunebrandingprofile?view=graph-rest-beta\n\nProvider Note: There have been **nonspecific errors reported when trying to _assign_ with application permissions** (even though the MS docs state that it should be possible)! Empty `assignments` or using delegated permissions both seemed to work in our tests, though. ||| MS Graph: Mobile app management (MAM)",
}
