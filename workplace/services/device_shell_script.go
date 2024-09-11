package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	DeviceShellScriptResource = generic.GenericResource{
		TypeNameSuffix: "device_shell_script",
		SpecificSchema: deviceShellScriptResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:         "/deviceManagement/deviceShellScripts",
			ReadOptions:     deviceShellScriptReadOptions,
			WriteSubActions: deviceShellScriptWriteSubActions,
		},
	}

	DeviceShellScriptSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&DeviceShellScriptResource)

	DeviceShellScriptPluralDataSource = generic.CreateGenericDataSourcePluralFromSingular(
		&DeviceShellScriptSingularDataSource, "")
)

var deviceShellScriptReadOptions = generic.ReadOptions{
	ODataExpand: "assignments",
}

var deviceShellScriptWriteSubActions = []generic.WriteSubAction{
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			AttributesMap: map[string]string{"assignments": "deviceManagementScriptAssignments"},
			UriSuffix:     "assign",
		},
	},
}

var deviceShellScriptResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // deviceShellScript
		"id": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"block_execution_notifications": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Does not notify the user a script is being executed",
		},
		"created_date_time": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"description": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "Optional description for the device management script.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Name of the device management script.",
		},
		"execution_frequency": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("PT0S")},
			Computed:            true,
			MarkdownDescription: "The interval for script to run. If not defined the script will run once",
		},
		"file_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Script file name.",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"retry_count": schema.Int64Attribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
			Computed:            true,
			MarkdownDescription: "Number of times for the script to be retried if it fails",
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: "List of Scope Tag IDs for this PowerShellScript instance.",
		},
		"run_as_account": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("system", "user")},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("user")},
			Computed:            true,
			MarkdownDescription: "Indicates the type of execution context. / Indicates the type of execution context the app runs in; possible values are: `system` (System context), `user` (User context)",
		},
		"script_content": schema.StringAttribute{
			Required:            true,
			Validators:          []validator.String{wpvalidator.TranslateValueEncodeBase64()},
			MarkdownDescription: "The script content.",
		},
		"assignments": deviceAndAppManagementAssignment,
	},
	MarkdownDescription: "Intune will provide customer the ability to run their Shell scripts on the enrolled Mac OS devices. The script can be run once or periodically.",
}
