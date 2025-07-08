package services

import (
	"context"
	"fmt"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	IdentityGovernanceWorkflowResource = generic.GenericResource{
		TypeNameSuffix: "identity_governance_workflow",
		SpecificSchema: identityGovernanceWorkflowResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/identityGovernance/lifecycleWorkflows/workflows",
			ReadOptions: generic.ReadOptions{
				DataSource: generic.DataSourceOptions{
					Plural: generic.PluralOptions{
						ExtraAttributes: []string{"category", "is_enabled", "is_scheduling_enabled"},
					},
				},
			},
			UpdateReplaceFunc: identityGovernanceWorkflowUpdateReplaceFunc,
		},
	}

	IdentityGovernanceWorkflowSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&IdentityGovernanceWorkflowResource)

	IdentityGovernanceWorkflowPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&IdentityGovernanceWorkflowResource, "")
)

func identityGovernanceWorkflowUpdateReplaceFunc(ctx context.Context, diags *diag.Diagnostics, params *generic.UpdateReplaceFuncParams) {
	baseUri := fmt.Sprintf("%s/%s/createNewVersion", params.R.AccessParams.BaseUri, params.Id)
	rawVal := map[string]any{"workflow": params.RawVal}
	_, params.RawResult = params.R.AccessParams.CreateRaw2(ctx, diags, baseUri, params.IdAttributer, rawVal, false, false)
}

var identityGovernanceWorkflowResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // identityGovernance.workflow
		"category": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf("joiner", "leaver", "unknownFutureValue", "mover"),
			},
			MarkdownDescription: "The category of the HR function supported by the workflows created using this template. A workflow can only belong to one category. The Required.<br><br>Supports `$filter`(`eq`,`ne`) and `$orderby` / Possible values are: `joiner`, `leaver`, `unknownFutureValue`, `mover`",
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "When the `workflow` was created.<br><br>Supports `$filter`(`lt`, `le`, `gt`, `ge`, `eq`, `ne`) and `$orderby`.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "The description of the `workflow`. Optional. The _provider_ default value is `\"\"`.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The display name of the `workflow`. Required.<br><br>Supports `$filter`(`eq`, `ne`) and `orderby`.",
		},
		"execution_conditions": schema.SingleNestedAttribute{
			Required: true,
			Attributes: map[string]schema.Attribute{ // identityGovernance.workflowExecutionConditions
				"on_demand_execution_only": generic.OdataDerivedTypeNestedAttributeRs{
					DerivedType: "#microsoft.graph.identityGovernance.onDemandExecutionOnly",
					SingleNestedAttribute: schema.SingleNestedAttribute{
						Optional:   true,
						Attributes: map[string]schema.Attribute{ // identityGovernance.onDemandExecutionOnly
						},
						Validators: []validator.Object{
							identityGovernanceWorkflowIdentityGovernanceWorkflowExecutionConditionsValidator,
						},
						MarkdownDescription: "Represents the execution condition of a [lifecycle workflow](../resources/identitygovernance-workflow.md) running on-demand only instead of by schedule. / https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-ondemandexecutiononly?view=graph-rest-beta",
					},
				},
				"trigger_and_scope_based_conditions": generic.OdataDerivedTypeNestedAttributeRs{
					DerivedType: "#microsoft.graph.identityGovernance.triggerAndScopeBasedConditions",
					SingleNestedAttribute: schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // identityGovernance.triggerAndScopeBasedConditions
							"scope": schema.SingleNestedAttribute{
								Required: true,
								Attributes: map[string]schema.Attribute{ // subjectSet
									"group": generic.OdataDerivedTypeNestedAttributeRs{
										DerivedType: "#microsoft.graph.identityGovernance.groupBasedSubjectSet",
										SingleNestedAttribute: schema.SingleNestedAttribute{
											Optional: true,
											Attributes: map[string]schema.Attribute{ // identityGovernance.groupBasedSubjectSet
												"groups": schema.SetNestedAttribute{
													Required: true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{ // group
															"id": schema.StringAttribute{
																Required:            true,
																MarkdownDescription: "The unique identifier for the group. <br><br>Returned by default. Key. Not nullable. Read-only. <br><br>Supports `$filter` (`eq`, `ne`, `not`, `in`).",
															},
														},
													},
													MarkdownDescription: "The specific group a user is interacting with in a [membershipChangeTrigger](identitygovernance-membershipchangetrigger.md) workflow. / Represents a Microsoft Entra group, which can be a Microsoft 365 group, a team in Microsoft Teams, or a security group. This resource is an open type that allows other properties to be passed in.\n\nFor performance reasons, the [create](../api/group-post-groups.md), [get](../api/group-get.md), and [list](../api/group-list.md) operations return only a subset of more commonly used properties by default. These _default_ properties are noted in the [Properties](#properties) section. To get any of the properties not returned by default, specify them in a `$select` OData query option.\n\nThis resource supports: / https://learn.microsoft.com/en-us/graph/api/resources/group?view=graph-rest-beta",
												},
											},
											Validators:          []validator.Object{identityGovernanceWorkflowSubjectSetValidator},
											MarkdownDescription: "Defines the group that is the scope of a lifecycle workflow [membershipChangeTrigger](../resources/identitygovernance-membershipchangetrigger.md) configuration. / https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-groupbasedsubjectset?view=graph-rest-beta",
										},
									},
									"rule": generic.OdataDerivedTypeNestedAttributeRs{
										DerivedType: "#microsoft.graph.identityGovernance.ruleBasedSubjectSet",
										SingleNestedAttribute: schema.SingleNestedAttribute{
											Optional: true,
											Attributes: map[string]schema.Attribute{ // identityGovernance.ruleBasedSubjectSet
												"rule": schema.StringAttribute{
													Required:            true,
													MarkdownDescription: "The rule for the subject set. Lifecycle Workflows supports a rich set of [user properties](user.md#properties) for configuring the rules using `$filter` query expressions. For more information, see [supported user and query parameters](#supported-user-properties-and-query-parameters).",
												},
											},
											Validators:          []validator.Object{identityGovernanceWorkflowSubjectSetValidator},
											MarkdownDescription: "Specifies the rules to define the subjects that are the scope of a lifecycle workflow [triggerAndScopeBasedConditions](../resources/identitygovernance-triggerandscopebasedconditions.md) configuration. / https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-rulebasedsubjectset?view=graph-rest-beta",
										},
									},
								},
								MarkdownDescription: "Defines who the workflow runs for. / A shared object that is used in entitlement management access package assignment policies, role management policies, and lifecycle workflows.\n\nThis object is an abstract base type from which the following resources are derived: / https://learn.microsoft.com/en-us/graph/api/resources/subjectset?view=graph-rest-beta",
							},
							"trigger": schema.SingleNestedAttribute{
								Required: true,
								Attributes: map[string]schema.Attribute{ // identityGovernance.workflowExecutionTrigger
									"attribute_change": generic.OdataDerivedTypeNestedAttributeRs{
										DerivedType: "#microsoft.graph.identityGovernance.attributeChangeTrigger",
										SingleNestedAttribute: schema.SingleNestedAttribute{
											Optional: true,
											Attributes: map[string]schema.Attribute{ // identityGovernance.attributeChangeTrigger
												"trigger_attributes": schema.SetNestedAttribute{
													Required: true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{ // identityGovernance.triggerAttribute
															"name": schema.StringAttribute{
																Required:            true,
																MarkdownDescription: "The name of the trigger attribute that is changed to trigger an [attributeChangeTrigger](../resources/identitygovernance-attributechangetrigger.md) workflow.",
															},
														},
													},
													MarkdownDescription: "The trigger attribute being changed that triggers the workflowexecutiontrigger of a workflow.) / Defines the trigger attribute, which is changed to activate a workflow using an [attributeChangeTrigger](../resources/identitygovernance-attributechangetrigger.md). / https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-triggerattribute?view=graph-rest-beta",
												},
											},
											Validators: []validator.Object{
												identityGovernanceWorkflowIdentityGovernanceWorkflowExecutionTriggerValidator,
											},
											MarkdownDescription: "Represents changes in user attributes that trigger the execution of workload conditions for a user. / https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-attributechangetrigger?view=graph-rest-beta",
										},
									},
									"membership_change": generic.OdataDerivedTypeNestedAttributeRs{
										DerivedType: "#microsoft.graph.identityGovernance.membershipChangeTrigger",
										SingleNestedAttribute: schema.SingleNestedAttribute{
											Optional: true,
											Attributes: map[string]schema.Attribute{ // identityGovernance.membershipChangeTrigger
												"change_type": schema.StringAttribute{
													Required: true,
													Validators: []validator.String{
														stringvalidator.OneOf("add", "remove", "unknownFutureValue"),
													},
													MarkdownDescription: "Defines what change that happens to the workflow group to trigger the [workflowExecutionTrigger](../resources/identitygovernance-workflowexecutiontrigger.md). / Possible values are: `add`, `remove`, `unknownFutureValue`",
												},
											},
											Validators: []validator.Object{
												identityGovernanceWorkflowIdentityGovernanceWorkflowExecutionTriggerValidator,
											},
											MarkdownDescription: "Represents the change in group membership that triggers the execution conditions of a workflow for a user. / https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-membershipchangetrigger?view=graph-rest-beta",
										},
									},
									"time_based_attribute": generic.OdataDerivedTypeNestedAttributeRs{
										DerivedType: "#microsoft.graph.identityGovernance.timeBasedAttributeTrigger",
										SingleNestedAttribute: schema.SingleNestedAttribute{
											Optional: true,
											Attributes: map[string]schema.Attribute{ // identityGovernance.timeBasedAttributeTrigger
												"offset_in_days": schema.Int64Attribute{
													Required:            true,
													MarkdownDescription: "How many days before or after the time-based attribute specified the workflow should trigger. For example, if the attribute is `employeeHireDate` and offsetInDays is -1, then the workflow should trigger one day before the employee hire date. The value can range between -180 and 180 days.",
												},
												"time_based_attribute": schema.StringAttribute{
													Required: true,
													Validators: []validator.String{
														stringvalidator.OneOf("employeeHireDate", "employeeLeaveDateTime", "unknownFutureValue", "createdDateTime"),
													},
													MarkdownDescription: "Determines which time-based identity property to reference. The / Possible values are: `employeeHireDate`, `employeeLeaveDateTime`, `unknownFutureValue`, `createdDateTime`",
												},
											},
											Validators: []validator.Object{
												identityGovernanceWorkflowIdentityGovernanceWorkflowExecutionTriggerValidator,
											},
											MarkdownDescription: "Trigger based on a time-based attribute for initiating the execution of a [lifecycle workflow](../resources/identitygovernance-workflow.md). The combination of scope and trigger conditions determines when a workflow is executed and on which identities. / https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-timebasedattributetrigger?view=graph-rest-beta",
										},
									},
								},
								MarkdownDescription: "What triggers a workflow to run. / The workflowExecutionTrigger type represents the workflow execution trigger when the [workflow runs on schedule](../resources/identitygovernance-triggerandscopebasedconditions.md). Inherited by the following derived types: / https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-workflowexecutiontrigger?view=graph-rest-beta",
							},
						},
						Validators: []validator.Object{
							identityGovernanceWorkflowIdentityGovernanceWorkflowExecutionConditionsValidator,
						},
						MarkdownDescription: "Represents a lifecycle workflow running by schedule, who it runs for, and what triggers the workflow to run. / https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-triggerandscopebasedconditions?view=graph-rest-beta",
					},
				},
			},
			MarkdownDescription: "Conditions describing when to execute the workflow and the criteria to identify in-scope subject set. Required. / The workflowExecutionConditions type notes the workflow execution conditions in [workflowTemplate](../resources/identitygovernance-workflowtemplate.md) and [workflowBase](../resources/identitygovernance-workflowbase.md) objects. Execution conditions define when a workflow runs and rules that identify the users that are the target of the workflow. The following types are derived from this abstract type: / https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-workflowexecutionconditions?view=graph-rest-beta",
		},
		"is_enabled": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Whether the workflow is enabled or disabled. If this setting is `true`, the workflow can be run on demand or on schedule when **isSchedulingEnabled** is `true`. Optional. Defaults to `true`.<br><br>Supports `$filter`(`eq`, `ne`) and `orderBy`. The _provider_ default value is `false`.",
		},
		"is_scheduling_enabled": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "If `true`, the Lifecycle Workflow engine executes the workflow based on the schedule defined by [tenant settings](identitygovernance-lifecyclemanagementsettings.md). Cannot be `true` for a disabled workflow (where **isEnabled** is `false`). Optional. Defaults to `false`.<br><br>Supports `$filter`(`eq`, `ne`) and `orderBy`. The _provider_ default value is `false`.",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The date time when the `workflow` was last modified.<br><br>Supports `$filter`(`lt`, `le`, `gt`, `ge`, `eq`, `ne`) and `$orderby`.",
		},
		"created_by": schema.SingleNestedAttribute{
			Computed: true,
			Attributes: map[string]schema.Attribute{ // user
				"id": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The user identifier.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpplanmodifier.ObjectUseStateForUnknown()},
			MarkdownDescription: "The unique identifier of the Microsoft Entra user that created the [workflow](../resources/identitygovernance-workflow.md) object.<br><br>Supports `$filter`(`eq`, `ne`) and `$expand`. / Represents an Azure Active Directory user object. / https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-user?view=graph-rest-beta",
		},
		"last_modified_by": schema.SingleNestedAttribute{
			Computed: true,
			Attributes: map[string]schema.Attribute{ // user
				"id": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The user identifier.",
				},
			},
			MarkdownDescription: "The user who last modified the [workflow](../resources/identitygovernance-workflow.md) object.<br><br>Supports `$filter`(`eq`, `ne`) and `$expand`. / Represents an Azure Active Directory user object. / https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-user?view=graph-rest-beta",
		},
		"tasks": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // identityGovernance.task
					"id": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "Unique identifier for the task. By default, this value will not change if a task is moved from one list to another.",
					},
					"arguments": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // keyValuePair
								"name": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "Name for this key-value pair",
								},
								"value": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Value for this key-value pair",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Key value pair for storing custom settings / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"category": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							wpvalidator.FlagEnumValues("joiner", "leaver", "unknownFutureValue", "mover"),
						},
						MarkdownDescription: "Possible values are: `joiner`, `leaver`, `unknownFutureValue`, `mover`",
					},
					"continue_on_error": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: ". The _provider_ default value is `false`.",
					},
					"description": schema.StringAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
						Computed:            true,
						MarkdownDescription: ". The _provider_ default value is `\"\"`.",
					},
					"display_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The name of the task.",
					},
					"execution_sequence": schema.Int64Attribute{
						Computed:      true,
						PlanModifiers: []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
					},
					"is_enabled": schema.BoolAttribute{
						Required: true,
					},
					"task_definition_id": schema.StringAttribute{
						Required: true,
					},
				},
			},
			PlanModifiers:       []planmodifier.List{wpdefaultvalue.ListDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Represents the configured tasks to execute and their execution sequence within a [workflow](../resources/identitygovernance-workflow.md) object. Required. / Represents a task, such as a piece of work or personal item, that can be tracked and completed. A **task** is always contained in a [base task list](basetasklist.md).\n\nThis resource supports the following: / https://learn.microsoft.com/en-us/graph/api/resources/task?view=graph-rest-beta. The _provider_ default value is `[]`.",
		},
		"deleted_date_time": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "When the workflow was deleted.<br><br>Supports `$filter`(`lt`, `le`, `gt`, `ge`, `eq`, `ne`) and `$orderby`.",
		},
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Identifier used for individually addressing a specific workflow.<br><br>Supports `$filter`(`eq`, `ne`) and `$orderby`.",
		},
		"next_schedule_run_date_time": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The date time when the `workflow` is expected to run next based on the schedule interval, if there are any users matching the execution conditions. <br><br>Supports `$filter`(`lt`,`gt`) and `$orderby`.",
		},
		"version": schema.Int64Attribute{
			Computed:            true,
			MarkdownDescription: "The current version number of the workflow. Value is 1 when the workflow is first created.<br><br>Supports `$filter`(`lt`, `le`, `gt`, `ge`, `eq`, `ne`) and `$orderby`.",
		},
	},
	MarkdownDescription: "Represents workflows created using Lifecycle Workflows. Workflows, when triggered by execution conditions, automate parts of the lifecycle management process using tasks. These tasks can either be built-in tasks, or a custom task can be called using the custom task extension which integrate with Azure Logic Apps.\n\nYou can create up to 100 workflows in a tenant. / https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-workflow?view=graph-rest-beta ||| MS Graph: Lifecycle workflows",
}

var identityGovernanceWorkflowIdentityGovernanceWorkflowExecutionConditionsValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("on_demand_execution_only"),
	path.MatchRelative().AtParent().AtName("trigger_and_scope_based_conditions"),
)

var identityGovernanceWorkflowSubjectSetValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("group"),
	path.MatchRelative().AtParent().AtName("rule"),
)

var identityGovernanceWorkflowIdentityGovernanceWorkflowExecutionTriggerValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("attribute_change"),
	path.MatchRelative().AtParent().AtName("membership_change"),
	path.MatchRelative().AtParent().AtName("time_based_attribute"),
)
