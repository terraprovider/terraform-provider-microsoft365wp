package services

import (
	"context"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	SharepointSettingsResource = generic.GenericResource{
		TypeNameSuffix: "sharepoint_settings",
		SpecificSchema: sharepointSettingsResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:                    "/admin/sharepoint/settings",
			IsSingleton:                true,
			GraphToTerraformMiddleware: sharepointSettingsGraphToTerraformMiddleware,
		},
	}

	SharepointSettingsSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&SharepointSettingsResource)
)

func sharepointSettingsGraphToTerraformMiddleware(ctx context.Context, diags *diag.Diagnostics, params *generic.GraphToTerraformMiddlewareParams) generic.GraphToTerraformMiddlewareReturns {
	// MS Graph does indeed not return anything for id, so add something helpful
	if params.RawVal["id"] == nil {
		params.RawVal["id"] = "singleton"
	}
	return nil
}

var sharepointSettingsResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // sharepointSettings
		"id": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"allowed_domain_guids_for_sync_app": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Collection of trusted domain GUIDs for the OneDrive sync app. The _provider_ default value is `[]`.",
		},
		"available_managed_paths_for_site_creation": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			PlanModifiers: []planmodifier.Set{
				wpdefaultvalue.SetDefaultValue([]any{"/sites/", "/teams/"}),
			},
			Computed:            true,
			MarkdownDescription: "Collection of managed paths available for site creation. Read-only. The _provider_ default value is `[\"/sites/\",\"/teams/\"]`.",
		},
		"deleted_user_personal_site_retention_period_in_days": schema.Int64Attribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(30)},
			Computed:            true,
			MarkdownDescription: "The number of days for preserving a deleted user's OneDrive. The _provider_ default value is `30`.",
		},
		"excluded_file_extensions_for_sync_app": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Collection of file extensions not uploaded by the OneDrive sync app. The _provider_ default value is `[]`.",
		},
		"idle_session_sign_out": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // idleSessionSignOut
				"is_enabled": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
					Computed:            true,
					MarkdownDescription: "Indicates whether the idle session sign-out policy is enabled. The _provider_ default value is `false`.",
				},
				"sign_out_after_in_seconds": schema.Int64Attribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
					Computed:            true,
					MarkdownDescription: "Number of seconds of inactivity after which a user is signed out. The _provider_ default value is `0`.",
				},
				"warn_after_in_seconds": schema.Int64Attribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
					Computed:            true,
					MarkdownDescription: "Number of seconds of inactivity after which a user is notified that they'll be signed out. The _provider_ default value is `0`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Specifies the idle session sign-out policies for the tenant. / Represents the idle session sign-out policy settings for SharePoint. / https://learn.microsoft.com/en-us/graph/api/resources/idlesessionsignout?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
		"image_tagging_option": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("disabled", "basic", "enhanced", "unknownFutureValue"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("basic")},
			Computed:            true,
			MarkdownDescription: "Specifies the image tagging option for the tenant. / Possible values are: `disabled`, `basic`, `enhanced`, `unknownFutureValue`. The _provider_ default value is `\"basic\"`.",
		},
		"is_commenting_on_site_pages_enabled": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:            true,
			MarkdownDescription: "Indicates whether comments are allowed on modern site pages in SharePoint. The _provider_ default value is `true`.",
		},
		"is_file_activity_notification_enabled": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:            true,
			MarkdownDescription: "Indicates whether push notifications are enabled for OneDrive events. The _provider_ default value is `true`.",
		},
		"is_legacy_auth_protocols_enabled": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:            true,
			MarkdownDescription: "Indicates whether legacy authentication protocols are enabled for the tenant. The _provider_ default value is `true`.",
		},
		"is_loop_enabled": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Indicates whetherif Fluid Framework is allowed on SharePoint sites. The _provider_ default value is `false`.",
		},
		"is_mac_sync_app_enabled": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:            true,
			MarkdownDescription: "Indicates whether files can be synced using the OneDrive sync app for Mac. The _provider_ default value is `true`.",
		},
		"is_require_accepting_user_to_match_invited_user_enabled": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "Indicates whether guests must sign in using the same account to which sharing invitations are sent.",
		},
		"is_resharing_by_external_users_enabled": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:            true,
			MarkdownDescription: "Indicates whether guests are allowed to reshare files, folders, and sites they don't own. The _provider_ default value is `true`.",
		},
		"is_sharepoint_mobile_notification_enabled": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:            true,
			Description:         `isSharePointMobileNotificationEnabled`, // custom MS Graph attribute name
			MarkdownDescription: "Indicates whether mobile push notifications are enabled for SharePoint. The _provider_ default value is `true`.",
		},
		"is_sharepoint_newsfeed_enabled": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			Description:         `isSharePointNewsfeedEnabled`, // custom MS Graph attribute name
			MarkdownDescription: "Indicates whether the newsfeed is allowed on the modern site pages in SharePoint. The _provider_ default value is `false`.",
		},
		"is_site_creation_enabled": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:            true,
			MarkdownDescription: "Indicates whether users are allowed to create sites. The _provider_ default value is `true`.",
		},
		"is_site_creation_ui_enabled": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:            true,
			Description:         `isSiteCreationUIEnabled`, // custom MS Graph attribute name
			MarkdownDescription: "Indicates whether the UI commands for creating sites are shown. The _provider_ default value is `true`.",
		},
		"is_site_pages_creation_enabled": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:            true,
			MarkdownDescription: "Indicates whether creating new modern pages is allowed on SharePoint sites. The _provider_ default value is `true`.",
		},
		"is_sites_storage_limit_automatic": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:            true,
			MarkdownDescription: "Indicates whether site storage space is automatically managed or if specific storage limits are set per site. The _provider_ default value is `true`.",
		},
		"is_sync_button_hidden_on_personal_site": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Indicates whether the sync button in OneDrive is hidden. The _provider_ default value is `false`.",
		},
		"is_unmanaged_sync_app_for_tenant_restricted": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Indicates whether users are allowed to sync files only on PCs joined to specific domains. The _provider_ default value is `false`.",
		},
		"personal_site_default_storage_limit_in_mb": schema.Int64Attribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(1048576)},
			Computed:            true,
			Description:         `personalSiteDefaultStorageLimitInMB`, // custom MS Graph attribute name
			MarkdownDescription: "The default OneDrive storage limit for all new and existing users who are assigned a qualifying license. Measured in megabytes (MB). The _provider_ default value is `1048576`.",
		},
		"sharing_allowed_domain_list": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Collection of email domains that are allowed for sharing outside the organization. The _provider_ default value is `[]`.",
		},
		"sharing_blocked_domain_list": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Collection of email domains that are blocked for sharing outside the organization. The _provider_ default value is `[]`.",
		},
		"sharing_capability": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("disabled", "externalUserSharingOnly", "externalUserAndGuestSharing", "existingExternalUserSharingOnly", "unknownFutureValue"),
			},
			PlanModifiers: []planmodifier.String{
				wpdefaultvalue.StringDefaultValue("externalUserAndGuestSharing"),
			},
			Computed:            true,
			MarkdownDescription: "Sharing capability for the tenant. / Possible values are: `disabled`, `externalUserSharingOnly`, `externalUserAndGuestSharing`, `existingExternalUserSharingOnly`, `unknownFutureValue`. The _provider_ default value is `\"externalUserAndGuestSharing\"`.",
		},
		"sharing_domain_restriction_mode": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("none", "allowList", "blockList", "unknownFutureValue"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
			Computed:            true,
			MarkdownDescription: "Specifies the external sharing mode for domains. / Possible values are: `none`, `allowList`, `blockList`, `unknownFutureValue`. The _provider_ default value is `\"none\"`.",
		},
		"site_creation_default_managed_path": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("/sites/")},
			Computed:            true,
			MarkdownDescription: "The value of the team site managed path. This is the path under which new team sites will be created. The _provider_ default value is `\"/sites/\"`.",
		},
		"site_creation_default_storage_limit_in_mb": schema.Int64Attribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(26214400)},
			Computed:            true,
			Description:         `siteCreationDefaultStorageLimitInMB`, // custom MS Graph attribute name
			MarkdownDescription: "The default storage quota for a new site upon creation. Measured in megabytes (MB). The _provider_ default value is `26214400`.",
		},
		"tenant_default_timezone": schema.StringAttribute{
			Optional: true,
			PlanModifiers: []planmodifier.String{
				wpdefaultvalue.StringDefaultValue("(UTC) Coordinated Universal Time"),
			},
			Computed:            true,
			MarkdownDescription: "The default timezone of a tenant for newly created sites. For a list of possible values, see [SPRegionalSettings.TimeZones property](/sharepoint/dev/schema/regional-settings-schema). The _provider_ default value is `\"(UTC) Coordinated Universal Time\"`.",
		},
	},
	MarkdownDescription: "Represents the tenant-level settings for SharePoint and OneDrive. / https://learn.microsoft.com/en-us/graph/api/resources/sharepointsettings?view=graph-rest-beta ||| MS Graph: Sites and lists",
}
