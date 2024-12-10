package services

import (
	"context"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	DeviceCustomAttributeShellScriptResource = generic.GenericResource{
		TypeNameSuffix: "device_custom_attribute_shell_script",
		SpecificSchema: deviceCustomAttributeShellScriptResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:                    "/deviceManagement/deviceCustomAttributeShellScripts",
			ReadOptions:                deviceCustomAttributeShellScriptReadOptions,
			WriteSubActions:            deviceCustomAttributeShellScriptWriteSubActions,
			TerraformToGraphMiddleware: deviceCustomAttributeShellScriptTerraformToGraphMiddleware,
		},
	}

	DeviceCustomAttributeShellScriptSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&DeviceCustomAttributeShellScriptResource)

	DeviceCustomAttributeShellScriptPluralDataSource = generic.CreateGenericDataSourcePluralFromSingular(
		&DeviceCustomAttributeShellScriptSingularDataSource, "")
)

var deviceCustomAttributeShellScriptReadOptions = generic.ReadOptions{
	ODataExpand: "assignments",
}

var deviceCustomAttributeShellScriptWriteSubActions = []generic.WriteSubAction{
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			AttributesMap: map[string]string{"assignments": "deviceManagementScriptAssignments"},
			UriSuffix:     "assign",
		},
	},
}

func deviceCustomAttributeShellScriptTerraformToGraphMiddleware(ctx context.Context, params generic.TerraformToGraphMiddlewareParams) generic.TerraformToGraphMiddlewareReturns {
	if params.IsUpdate {
		// displayName, customAttributeType cannot be updated and may not even be written again after creation
		delete(params.RawVal, "displayName")
		delete(params.RawVal, "customAttributeType")
	}
	return nil
}

var deviceCustomAttributeShellScriptResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // deviceCustomAttributeShellScript
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Unique Identifier for the custom attribute entity.",
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The date and time the device management script was created. This property is read-only.",
		},
		"custom_attribute_type": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf("integer", "string", "dateTime"),
			},
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			MarkdownDescription: "The expected type of the custom attribute's value. / Represents the expected type for a macOS custom attribute script value; possible values are: `integer` (Indicates the value for a custom attribute script is an integer.), `string` (Indicates the value for a custom attribute script is a string.), `dateTime` (Indicates the value for a custom attribute script is a date conforming to ISO 8601.)",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "Optional description for the device management script. The _provider_ default value is `\"\"`.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			MarkdownDescription: "Name of the device management script.",
		},
		"file_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Script file name.",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The date and time the device management script was last modified. This property is read-only.",
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: "List of Scope Tag IDs for this PowerShellScript instance. The _provider_ default value is `[\"0\"]`.",
		},
		"run_as_account": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("system", "user")},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("system")},
			Computed:            true,
			MarkdownDescription: "Indicates the type of execution context. / Indicates the type of execution context the app runs in; possible values are: `system` (System context), `user` (User context). The _provider_ default value is `\"system\"`.",
		},
		"script_content": schema.StringAttribute{
			Required:            true,
			Validators:          []validator.String{wpvalidator.TranslateValueEncodeBase64()},
			MarkdownDescription: "The script content.",
		},
		"assignments": deviceAndAppManagementAssignment,
	},
	MarkdownDescription: "Represents a custom attribute script for macOS. / https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-devicecustomattributeshellscript?view=graph-rest-beta ||| MS Graph: Device management",
}
