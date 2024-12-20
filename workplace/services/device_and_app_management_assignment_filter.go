package services

import (
	"context"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	DeviceAndAppManagementAssignmentFilterResource = generic.GenericResource{
		TypeNameSuffix: "device_and_app_management_assignment_filter",
		SpecificSchema: deviceAndAppManagementAssignmentFilterResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:                    "/deviceManagement/assignmentFilters",
			ReadOptions:                deviceAndAppManagementAssignmentFilterReadOptions,
			TerraformToGraphMiddleware: deviceAndAppManagementAssignmentFilterTerraformToGraphMiddleware,
		},
	}

	DeviceAndAppManagementAssignmentFilterSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&DeviceAndAppManagementAssignmentFilterResource)

	DeviceAndAppManagementAssignmentFilterPluralDataSource = generic.CreateGenericDataSourcePluralFromSingular(
		&DeviceAndAppManagementAssignmentFilterSingularDataSource, "")
)

var deviceAndAppManagementAssignmentFilterReadOptions = generic.ReadOptions{
	PluralNoFilterSupport: true,
}

func deviceAndAppManagementAssignmentFilterTerraformToGraphMiddleware(ctx context.Context, params generic.TerraformToGraphMiddlewareParams) generic.TerraformToGraphMiddlewareReturns {
	if params.IsUpdate {
		// platform cannot be updated and may not even be written again after creation
		delete(params.RawVal, "platform")
	}
	return nil
}

var deviceAndAppManagementAssignmentFilterResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // deviceAndAppManagementAssignmentFilter
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Key of the Assignment Filter.",
		},
		"assignment_filter_management_type": schema.StringAttribute{
			Computed: true,
			Validators: []validator.String{
				stringvalidator.OneOf("devices", "apps", "unknownFutureValue"),
			},
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Indicates filter is applied to either 'devices' or 'apps' management type. Possible values are devices, apps. Default filter will be applied to 'devices'. / Supported filter management types whether its devices or apps; possible values are: `devices` (Indicates when filter is supported based on device properties. This is the default value when management type resolution fails.), `apps` (Indicates when filter is supported based on app properties.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.)",
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The creation time of the assignment filter. The value cannot be modified and is automatically populated during new assignment filter process. The timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 would look like this: '2014-01-01T00:00:00Z'.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "Optional description of the Assignment Filter. The _provider_ default value is `\"\"`.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The name of the Assignment Filter.",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Last modified time of the Assignment Filter. The timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 would look like this: '2014-01-01T00:00:00Z'",
		},
		"platform": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf("android", "androidForWork", "iOS", "macOS", "windowsPhone81", "windows81AndLater", "windows10AndLater", "androidWorkProfile", "unknown", "androidAOSP", "androidMobileApplicationManagement", "iOSMobileApplicationManagement", "unknownFutureValue", "windowsMobileApplicationManagement"),
			},
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			MarkdownDescription: "Indicates filter is applied to which flatform. Possible values are android,androidForWork,iOS,macOS,windowsPhone81,windows81AndLater,windows10AndLater,androidWorkProfile, unknown, androidAOSP, androidMobileApplicationManagement, iOSMobileApplicationManagement, windowsMobileApplicationManagement. Default filter will be applied to 'unknown'. / Supported platform types; possible values are: `android` (Indicates device platform type is android.), `androidForWork` (Indicates device platform type is android for work.), `iOS` (Indicates device platform type is iOS.), `macOS` (Indicates device platform type is macOS.), `windowsPhone81` (Indicates device platform type is WindowsPhone 8.1.), `windows81AndLater` (Indicates device platform type is Windows 8.1 and later.), `windows10AndLater` (Indicates device platform type is Windows 10 and later.), `androidWorkProfile` (Indicates device platform type is Android Work Profile.), `unknown` (This is the default value when device platform type resolution fails), `androidAOSP` (Indicates device platform type is Android AOSP.), `androidMobileApplicationManagement` (Indicates device platform type associated with Mobile Application Management (MAM) for android devices.), `iOSMobileApplicationManagement` (Indicates device platform type associated with Mobile Application Management (MAM) for iOS devices.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.), `windowsMobileApplicationManagement` (Indicates device platform type associated with Mobile Application Management (MAM) for Windows devices.)",
		},
		"role_scope_tags": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: "Indicates role scope tags assigned for the assignment filter. The _provider_ default value is `[\"0\"]`.",
		},
		"rule": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Rule definition of the assignment filter.",
		},
	},
	MarkdownDescription: "A class containing the properties used for Assignment Filter. / https://learn.microsoft.com/en-us/graph/api/resources/intune-policyset-deviceandappmanagementassignmentfilter?view=graph-rest-beta ||| MS Graph: Assignments",
}
