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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	WindowsDriverUpdateProfileResource = generic.GenericResource{
		TypeNameSuffix: "windows_driver_update_profile",
		SpecificSchema: windowsDriverUpdateProfileResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/deviceManagement/windowsDriverUpdateProfiles",
			ReadOptions: generic.ReadOptions{
				ODataExpand: "assignments,driverInventories",
				DataSource: generic.DataSourceOptions{
					NoFilterSupport: true,
				},
			},
			WriteOptions: generic.WriteOptions{
				SubActions: []generic.WriteSubAction{
					&generic.WriteSubActionAllInOne{
						WriteSubActionBase: generic.WriteSubActionBase{
							Attributes: []string{"assignments"},
							UriSuffix:  "assign",
						},
					},
				},
			},
			TerraformToGraphMiddleware: windowsDriverUpdateProfileTerraformToGraphMiddleware,
		},
	}

	WindowsDriverUpdateProfileSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&WindowsDriverUpdateProfileResource)

	WindowsDriverUpdateProfilePluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&WindowsDriverUpdateProfileResource, "")
)

func windowsDriverUpdateProfileTerraformToGraphMiddleware(ctx context.Context, diags *diag.Diagnostics, params *generic.TerraformToGraphMiddlewareParams) generic.TerraformToGraphMiddlewareReturns {
	if params.IsUpdate {
		// remove unchangeable attribute(s) from PATCH
		delete(params.RawVal, "approvalType")
	}
	return nil
}

var windowsDriverUpdateProfileResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // windowsDriverUpdateProfile
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The Intune policy id.",
		},
		"approval_type": schema.StringAttribute{
			Required:            true,
			Validators:          []validator.String{stringvalidator.OneOf("manual", "automatic")},
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			MarkdownDescription: "Driver update profile approval type. For example, manual or automatic approval. / An enum type to represent approval type of a driver update profile; possible values are: `manual` (This indicates a driver and firmware profile needs to be approved manually.), `automatic` (This indicates a driver and firmware profile is approved automatically.)",
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The date time that the profile was created.",
		},
		"deployment_deferral_in_days": schema.Int64Attribute{
			Optional:            true,
			MarkdownDescription: "Deployment deferral settings in days, only applicable when ApprovalType is set to automatic approval.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "The description of the profile which is specified by the user. The _provider_ default value is `\"\"`.",
		},
		"device_reporting": schema.Int64Attribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
			MarkdownDescription: "Number of devices reporting for this profile",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The display name for the profile.",
		},
		"inventory_sync_status": schema.SingleNestedAttribute{
			Computed: true,
			Attributes: map[string]schema.Attribute{ // windowsDriverUpdateProfileInventorySyncStatus
				"driver_inventory_sync_state": schema.StringAttribute{
					Computed: true,
					Validators: []validator.String{
						stringvalidator.OneOf("pending", "success", "failure"),
					},
					MarkdownDescription: "The state of the latest sync. / Windows DnF update inventory sync state; possible values are: `pending` (Pending sync.), `success` (Successful sync.), `failure` (Failed sync.)",
				},
				"last_successful_sync_date_time": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The last successful sync date and time in UTC.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpplanmodifier.ObjectUseStateForUnknown()},
			MarkdownDescription: "Driver inventory sync status for this profile. / A complex type to store the status of a driver and firmware profile inventory sync. The status includes the last successful sync date time and the state of the last sync. / https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-windowsdriverupdateprofileinventorysyncstatus?view=graph-rest-beta",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The date time that the profile was last modified.",
		},
		"new_updates": schema.Int64Attribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
			MarkdownDescription: "Number of new driver updates available for this profile.",
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: "List of Scope Tags for this Driver Update entity. The _provider_ default value is `[\"0\"]`.",
		},
		"assignments": deviceAndAppManagementAssignment,
		"driver_inventories": schema.SetNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // windowsDriverUpdateInventory
					"id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The id of the driver.",
					},
					"applicable_device_count": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "The number of devices for which this driver is applicable.",
					},
					"approval_status": schema.StringAttribute{
						Computed: true,
						Validators: []validator.String{
							stringvalidator.OneOf("needsReview", "declined", "approved", "suspended"),
						},
						MarkdownDescription: "The approval status for this driver. / An enum type to represent approval status of a driver; possible values are: `needsReview` (This indicates a driver needs IT admin's review.), `declined` (This indicates IT admin has declined a driver.), `approved` (This indicates IT admin has approved a driver.), `suspended` (This indicates IT admin has suspended a driver.)",
					},
					"category": schema.StringAttribute{
						Computed: true,
						Validators: []validator.String{
							stringvalidator.OneOf("recommended", "previouslyApproved", "other"),
						},
						MarkdownDescription: "The category for this driver. / An enum type to represent which category a driver belongs to; possible values are: `recommended` (This indicates a driver is recommended by Microsoft.), `previouslyApproved` (This indicates a driver was recommended by Microsoft and IT admin has taken some approval action on it.), `other` (This indicates a driver is never recommended by Microsoft.)",
					},
					"deploy_date_time": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The date time when a driver should be deployed if approvalStatus is approved.",
					},
					"driver_class": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The class of the driver.",
					},
					"manufacturer": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The manufacturer of the driver.",
					},
					"name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The name of the driver.",
					},
					"release_date_time": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The release date time of the driver.",
					},
					"version": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The version of the driver.",
					},
				},
			},
			PlanModifiers:       []planmodifier.Set{wpplanmodifier.SetUseStateForUnknown()},
			MarkdownDescription: "Driver inventories for this profile. / A new entity to represent driver inventories. / https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-windowsdriverupdateinventory?view=graph-rest-beta",
		},
	},
	MarkdownDescription: "Windows Driver Update Profile / https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-windowsdriverupdateprofile?view=graph-rest-beta ||| MS Graph: Software updates",
}
