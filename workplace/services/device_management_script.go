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
	DeviceManagementScriptResource = generic.GenericResource{
		TypeNameSuffix: "device_management_script",
		SpecificSchema: deviceManagementScriptResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/deviceManagement/deviceManagementScripts",
			ReadOptions: generic.ReadOptions{
				ODataExpand: "assignments",
				DataSource: generic.DataSourceOptions{
					NoIdFilterSupport: true,
				},
			},
			WriteOptions: generic.WriteOptions{
				SubActions: []generic.WriteSubAction{
					&generic.WriteSubActionAllInOne{
						WriteSubActionBase: generic.WriteSubActionBase{
							AttributesMap: map[string]string{"assignments": "deviceManagementScriptAssignments"},
							UriSuffix:     "assign",
						},
					},
				},
			},
		},
	}

	DeviceManagementScriptSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&DeviceManagementScriptResource)

	DeviceManagementScriptPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&DeviceManagementScriptResource, "")
)

var deviceManagementScriptResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // deviceManagementScript
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Unique Identifier for the device management script.",
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The date and time the device management script was created. This property is read-only.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "Optional description for the device management script. The _provider_ default value is `\"\"`.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Name of the device management script.",
		},
		"enforce_signature_check": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Indicate whether the script signature needs be checked. The _provider_ default value is `false`.",
		},
		"file_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Script file name.",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The date and time the device management script was last modified. This property is read-only.",
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: "List of Scope Tag IDs for this PowerShellScript instance. The _provider_ default value is `[\"0\"]`.",
		},
		"run_as_32_bit": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			Description:         `runAs32Bit`, // custom MS Graph attribute name
			MarkdownDescription: "A value indicating whether the PowerShell script should run as 32-bit. The _provider_ default value is `false`.",
		},
		"run_as_account": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("system", "user")},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("user")},
			Computed:            true,
			MarkdownDescription: "Indicates the type of execution context. / Indicates the type of execution context the app runs in; possible values are: `system` (System context), `user` (User context). The _provider_ default value is `\"user\"`.",
		},
		"script_content": schema.StringAttribute{
			Required:            true,
			Validators:          []validator.String{wpvalidator.TranslateValueEncodeBase64()},
			MarkdownDescription: "The script content.",
		},
		"assignments": deviceAndAppManagementAssignment,
	},
	MarkdownDescription: "Intune will provide customer the ability to run their Powershell scripts on the enrolled windows 10 Azure Active Directory joined devices. The script can be run once or periodically. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-devicemanagementscript?view=graph-rest-beta ||| MS Graph: Device management",
}
