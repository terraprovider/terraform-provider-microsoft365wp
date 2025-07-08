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
	DeviceComplianceScriptResource = generic.GenericResource{
		TypeNameSuffix: "device_compliance_script",
		SpecificSchema: deviceComplianceScriptResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/deviceManagement/deviceComplianceScripts",
			ReadOptions: generic.ReadOptions{
				ODataExpand: "assignments",
				DataSource: generic.DataSourceOptions{
					ExtraFilterAttributes: []string{"display_name", "publisher", "run_as_account"},
					Plural: generic.PluralOptions{
						ExtraAttributes: []string{"display_name", "publisher", "run_as_account"},
					},
				},
			},
			WriteOptions: generic.WriteOptions{
				SubActions: []generic.WriteSubAction{
					&generic.WriteSubActionAllInOne{
						WriteSubActionBase: generic.WriteSubActionBase{
							AttributesMap: map[string]string{"assignments": "deviceHealthScriptAssignments"},
							UriSuffix:     "assign",
						},
					},
				},
			},
		},
	}

	DeviceComplianceScriptSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&DeviceComplianceScriptResource)

	DeviceComplianceScriptPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&DeviceComplianceScriptResource, "")
)

var deviceComplianceScriptResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // deviceComplianceScript
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Unique Identifier for the device compliance script",
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The timestamp of when the device compliance script was created. This property is read-only.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "Description of the device compliance script. The _provider_ default value is `\"\"`.",
		},
		"detection_script_content": schema.StringAttribute{
			Required:            true,
			Validators:          []validator.String{wpvalidator.TranslateValueEncodeBase64()},
			MarkdownDescription: "The entire content of the detection powershell script",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Name of the device compliance script",
		},
		"enforce_signature_check": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Indicate whether the script signature needs be checked. The _provider_ default value is `false`.",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The timestamp of when the device compliance script was modified. This property is read-only.",
		},
		"publisher": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "Name of the device compliance script publisher. The _provider_ default value is `\"\"`.",
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: "List of Scope Tag IDs for the device compliance script. The _provider_ default value is `[\"0\"]`.",
		},
		"run_as_32_bit": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			Description:         `runAs32Bit`, // custom MS Graph attribute name
			MarkdownDescription: "Indicate whether PowerShell script(s) should run as 32-bit. The _provider_ default value is `false`.",
		},
		"run_as_account": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("system", "user")},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("system")},
			Computed:            true,
			MarkdownDescription: "Indicates the type of execution context. / Indicates the type of execution context the app runs in; possible values are: `system` (System context), `user` (User context). The _provider_ default value is `\"system\"`.",
		},
		"version": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Version of the device compliance script",
		},
		"assignments": deviceAndAppManagementAssignment,
	},
	MarkdownDescription: "Intune will provide customer the ability to run their Powershell Compliance scripts (detection) on the enrolled windows 10 Azure Active Directory joined devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-devicecompliancescript?view=graph-rest-beta ||| MS Graph: Device management",
}
