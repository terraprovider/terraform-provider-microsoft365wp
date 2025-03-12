package services

import (
	"context"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	WindowsFeatureUpdateProfileResource = generic.GenericResource{
		TypeNameSuffix: "windows_feature_update_profile",
		SpecificSchema: windowsFeatureUpdateProfileResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:                    "/deviceManagement/windowsFeatureUpdateProfiles",
			ReadOptions:                windowsFeatureUpdateProfileReadOptions,
			WriteSubActions:            windowsFeatureUpdateProfileWriteSubActions,
			TerraformToGraphMiddleware: windowsFeatureUpdateProfileTerraformToGraphMiddleware,
		},
	}

	WindowsFeatureUpdateProfileSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&WindowsFeatureUpdateProfileResource)

	WindowsFeatureUpdateProfilePluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&WindowsFeatureUpdateProfileResource, "")
)

var windowsFeatureUpdateProfileReadOptions = generic.ReadOptions{
	ODataExpand: "assignments",
	DataSource: generic.DataSourceOptions{
		NoFilterSupport: true,
	},
}

var windowsFeatureUpdateProfileWriteSubActions = []generic.WriteSubAction{
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			Attributes: []string{"assignments"},
			UriSuffix:  "assign",
		},
	},
}

func windowsFeatureUpdateProfileTerraformToGraphMiddleware(ctx context.Context, diags *diag.Diagnostics, params *generic.TerraformToGraphMiddlewareParams) generic.TerraformToGraphMiddlewareReturns {
	if params.IsUpdate {
		// installLatestWindows10OnWindows11IneligibleDevice cannot be updated and may not even be written again after creation
		delete(params.RawVal, "installLatestWindows10OnWindows11IneligibleDevice")
	}
	return nil
}

var windowsFeatureUpdateProfileResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // windowsFeatureUpdateProfile
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The Identifier of the entity.",
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The date time that the profile was created.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The description of the profile which is specified by the user.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The display name of the profile.",
		},
		"end_of_support_date": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The last supported date for a feature update",
		},
		"feature_update_version": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The feature update version that will be deployed to the devices targeted by this profile. The version could be any supported version for example 1709, 1803 or 1809 and so on.",
		},
		"install_feature_updates_optional": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "If true, the Windows 11 update will become optional. The _provider_ default value is `false`.",
		},
		"install_latest_windows10_on_windows11_ineligible_device": schema.BoolAttribute{
			Optional: true,
			PlanModifiers: []planmodifier.Bool{
				wpdefaultvalue.BoolDefaultValue(false),
				boolplanmodifier.RequiresReplace(),
			},
			Computed:            true,
			Description:         `installLatestWindows10OnWindows11IneligibleDevice`, // custom MS Graph attribute name
			MarkdownDescription: "If true, the latest Microsoft Windows 10 update will be installed on devices ineligible for Microsoft Windows 11. The _provider_ default value is `false`.",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The date time that the profile was last modified.",
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: "List of Scope Tags for this Feature Update entity. The _provider_ default value is `[\"0\"]`.",
		},
		"rollout_settings": schema.SingleNestedAttribute{
			Required: true,
			Attributes: map[string]schema.Attribute{ // windowsUpdateRolloutSettings
				"offer_end_date_time_in_utc": schema.StringAttribute{
					Optional:            true,
					Description:         `offerEndDateTimeInUTC`, // custom MS Graph attribute name
					MarkdownDescription: "The feature update's ending  of release date and time to be set, update, and displayed for a feature Update profile for example: 2020-06-09T10:00:00Z.",
				},
				"offer_interval_in_days": schema.Int64Attribute{
					Optional:            true,
					MarkdownDescription: "The number of day(s) between each set of offers to be set, updated, and displayed for a feature update profile, for example: if OfferStartDateTimeInUTC is 2020-06-09T10:00:00Z, and OfferIntervalInDays is 1, then the next two sets of offers will be made consecutively on 2020-06-10T10:00:00Z (next day at the same specified time) and 2020-06-11T10:00:00Z (next next day at the same specified time) with 1 day in between each set of offers.",
				},
				"offer_start_date_time_in_utc": schema.StringAttribute{
					Optional:            true,
					Description:         `offerStartDateTimeInUTC`, // custom MS Graph attribute name
					MarkdownDescription: "The feature update's starting date and time to be set, update, and displayed for a feature Update profile for example: 2020-06-09T10:00:00Z.",
				},
			},
			MarkdownDescription: "The windows update rollout settings, including offer start date time, offer end date time, and days between each set of offers. / A complex type to store the windows update rollout settings including offer start date time, offer end date time, and days between each set of offers. / https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-windowsupdaterolloutsettings?view=graph-rest-beta",
		},
		"assignments": deviceAndAppManagementAssignment,
	},
	MarkdownDescription: "Windows Feature Update Profile / https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-windowsfeatureupdateprofile?view=graph-rest-beta ||| MS Graph: Software updates",
}
