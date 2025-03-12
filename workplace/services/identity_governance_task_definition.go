package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	identityGovernanceTaskDefinitionResource = generic.GenericResource{
		TypeNameSuffix: "identity_governance_task_definition",
		SpecificSchema: identityGovernanceTaskDefinitionResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:     "/identityGovernance/lifecycleWorkflows/taskDefinitions",
			ReadOptions: identityGovernanceTaskDefinitionReadOptions,
		},
	}

	IdentityGovernanceTaskDefinitionSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&identityGovernanceTaskDefinitionResource)

	IdentityGovernanceTaskDefinitionPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&identityGovernanceTaskDefinitionResource, "")
)

var identityGovernanceTaskDefinitionReadOptions = generic.ReadOptions{
	DataSource: generic.DataSourceOptions{
		ExtraFilterAttributes: []string{"display_name"},
		Plural: generic.PluralOptions{
			ExtraAttributes: []string{"category"},
		},
	},
}

var identityGovernanceTaskDefinitionResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // identityGovernance.taskDefinition
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the taskDefinition.",
		},
		"category": schema.StringAttribute{
			Validators: []validator.String{
				wpvalidator.FlagEnumValues("joiner", "leaver", "unknownFutureValue", "mover"),
			},
			MarkdownDescription: "The category of the HR function that the tasks created using this definition can be used with. The This is a multi-valued enumeration whose allowed combinations are `joiner`, `joiner,leaver`, or `leaver`.<br><br>Supports `$filter`(`eq`, `ne`, `has`) and `$orderby`. / Possible values are: `joiner`, `leaver`, `unknownFutureValue`, `mover`",
		},
		"continue_on_error": schema.BoolAttribute{
			MarkdownDescription: "Defines if the workflow will continue if the task has an error.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The description of the taskDefinition.",
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: "The display name of the taskDefinition.<br><br>Supports `$filter`(`eq`, `ne`) and `$orderby`.",
		},
		"parameters": schema.SetNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // identityGovernance.parameter
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the parameter.",
					},
					"values": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						MarkdownDescription: "The values of the parameter.",
					},
					"value_type": schema.StringAttribute{
						Validators: []validator.String{
							stringvalidator.OneOf("enum", "string", "int", "bool", "unknownFutureValue"),
						},
						MarkdownDescription: "The value type of the parameter. The / Possible values are: `enum`, `string`, `int`, `bool`, `unknownFutureValue`",
					},
				},
			},
			MarkdownDescription: "The parameters that must be supplied when creating a workflow task object.<br><br>Supports `$filter`(`any`). / Represents the allowed arguments that are defined in a [taskDefinition](../resources/identitygovernance-taskdefinition.md) for built-in lifecycle workflow tasks. / https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-parameter?view=graph-rest-beta",
		},
		"version": schema.Int64Attribute{
			MarkdownDescription: "The version number of the taskDefinition. New records are pushed when we add support for new parameters.<br><br>Supports `$filter`(`ge`, `gt`, `le`, `lt`, `eq`, `ne`) and `$orderby`.",
		},
	},
	MarkdownDescription: "Represents the built-in tasks that you can use to construct tasks for lifecycle workflows. Each task has a unique template identifier. For a full list of available built-in tasks, see [Configure the arguments for built-in Lifecycle Workflow tasks](/graph/identitygovernance-lifecycleworkflows-task-arguments). / https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-taskdefinition?view=graph-rest-beta ||| MS Graph: Lifecycle workflows",
}
