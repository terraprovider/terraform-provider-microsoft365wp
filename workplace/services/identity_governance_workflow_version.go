package services

import (
	"context"
	"strconv"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	identityGovernanceWorkflowVersionResource = generic.GenericResource{
		TypeNameSuffix: "identity_governance_workflow_version",
		SpecificSchema: identityGovernanceWorkflowVersionResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/identityGovernance/lifecycleWorkflows/workflows",
			ParentEntities: generic.ParentEntities{
				{
					ParentIdField: path.Root("workflow_id"),
					UriSuffix:     "versions",
				},
			},
			EntityId: generic.EntityIdOptions{
				AttrNameGraph: "versionNumber",
			},
			ReadOptions: generic.ReadOptions{
				ODataExpand: "tasks",
				DataSource: generic.DataSourceOptions{
					Plural: generic.PluralOptions{
						ExtraAttributes: []string{"category"},
					},
				},
			},
			GraphToTerraformMiddleware: identityGovernanceWorkflowVersionGraphToTerraformMiddleware,
		},
	}

	IdentityGovernanceWorkflowVersionSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&identityGovernanceWorkflowVersionResource)

	IdentityGovernanceWorkflowVersionPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&identityGovernanceWorkflowVersionResource, "")
)

func identityGovernanceWorkflowVersionGraphToTerraformMiddleware(ctx context.Context, diags *diag.Diagnostics, params *generic.GraphToTerraformMiddlewareParams) generic.GraphToTerraformMiddlewareReturns {
	// versionNumber is number, but we need a string
	if versionNumber, ok := params.RawVal["versionNumber"].(float64); ok {
		params.RawVal["versionNumber"] = strconv.FormatFloat(versionNumber, 'f', 0, 64)
	}
	return nil
}

var identityGovernanceWorkflowVersionResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // identityGovernance.workflowVersion
		"workflow_id": schema.StringAttribute{
			Required:      true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
		},
		"category": schema.StringAttribute{
			Validators: []validator.String{
				stringvalidator.OneOf("joiner", "leaver", "unknownFutureValue", "mover"),
			},
			MarkdownDescription: "The category of the HR function supported by the workflows created using this template. A workflow can only belong to one category. <br/> Supports `$filter`(`eq`,`ne`) and `$orderby` <br/> _Provider_ allowed values are: `joiner`, `leaver`, `unknownFutureValue`, `mover`.",
		},
		"created_date_time": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The date time when the `workflow` was versioned. <br/> Supports `$filter`(`lt`, `le`, `gt`, `ge`, `eq`, `ne`) and `$orderby`.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The description of the `workflowversion`.",
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: "The display name of the `workflowversion`. <br/> Supports `$filter`(`eq`, `ne`) and `orderby`.",
		},
		"execution_conditions": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // identityGovernance.workflowExecutionConditions
				"on_demand_execution_only": generic.OdataDerivedTypeNestedAttributeRs{
					DerivedType: "#microsoft.graph.identityGovernance.onDemandExecutionOnly",
					SingleNestedAttribute: schema.SingleNestedAttribute{
						Optional:   true,
						Attributes: map[string]schema.Attribute{ // identityGovernance.onDemandExecutionOnly
						},
						Validators: []validator.Object{
							identityGovernanceWorkflowVersionIdentityGovernanceWorkflowExecutionConditionsValidator,
						},
						MarkdownDescription: "Represents the execution condition of a [lifecycle workflow](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-workflow?view=graph-rest-beta) running on-demand only instead of by schedule. Also see [Microsoft docs for identityGovernance.onDemandExecutionOnly](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-ondemandexecutiononly?view=graph-rest-beta). <br> ",
					},
				},
				"trigger_and_scope_based_conditions": generic.OdataDerivedTypeNestedAttributeRs{
					DerivedType: "#microsoft.graph.identityGovernance.triggerAndScopeBasedConditions",
					SingleNestedAttribute: schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // identityGovernance.triggerAndScopeBasedConditions
							"scope": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{ // subjectSet
									"group": generic.OdataDerivedTypeNestedAttributeRs{
										DerivedType: "#microsoft.graph.identityGovernance.groupBasedSubjectSet",
										SingleNestedAttribute: schema.SingleNestedAttribute{
											Optional: true,
											Attributes: map[string]schema.Attribute{ // identityGovernance.groupBasedSubjectSet
												"groups": schema.SetNestedAttribute{
													Optional: true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{ // group
															"id": schema.StringAttribute{
																Required:            true,
																MarkdownDescription: "The unique identifier for the group. <br/> Returned by default. Key. Not nullable. Read-only. <br/> Supports `$filter` (`eq`, `ne`, `not`, `in`).",
															},
														},
													},
													MarkdownDescription: "The specific group a user is interacting with in a [membershipChangeTrigger](identitygovernance-membershipchangetrigger.md) workflow. <br/> Represents a Microsoft Entra group, which can be a Microsoft 365 group, a team in Microsoft Teams, or a security group. <br/> For performance reasons, the [create](https://learn.microsoft.com/en-us/graph/api/group-post-groups?view=graph-rest-beta), [get](https://learn.microsoft.com/en-us/graph/api/group-get?view=graph-rest-beta), and [list](https://learn.microsoft.com/en-us/graph/api/group-list?view=graph-rest-beta) operations return only a subset of more commonly used properties by default. These _default_ properties are noted in the [Properties](#properties) section. To get any of the properties not returned by default, specify them in a `$select` OData query option. Also see [Microsoft docs for group](https://learn.microsoft.com/en-us/graph/api/resources/group?view=graph-rest-beta). <br> ",
												},
											},
											Validators: []validator.Object{
												identityGovernanceWorkflowVersionSubjectSetValidator,
											},
											MarkdownDescription: "Defines the group that is the scope of a lifecycle workflow [membershipChangeTrigger](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-membershipchangetrigger?view=graph-rest-beta) configuration. Also see [Microsoft docs for identityGovernance.groupBasedSubjectSet](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-groupbasedsubjectset?view=graph-rest-beta). <br> ",
										},
									},
									"rule": generic.OdataDerivedTypeNestedAttributeRs{
										DerivedType: "#microsoft.graph.identityGovernance.ruleBasedSubjectSet",
										SingleNestedAttribute: schema.SingleNestedAttribute{
											Optional: true,
											Attributes: map[string]schema.Attribute{ // identityGovernance.ruleBasedSubjectSet
												"rule": schema.StringAttribute{
													MarkdownDescription: "The rule for the subject set. Lifecycle Workflows supports a rich set of [user properties](user.md#properties) for configuring the rules using `$filter` query expressions. For more information, see [supported user and query parameters](#supported-user-properties-and-query-parameters).",
												},
											},
											Validators: []validator.Object{
												identityGovernanceWorkflowVersionSubjectSetValidator,
											},
											MarkdownDescription: "Specifies the rules to define the subjects that are the scope of a lifecycle workflow [triggerAndScopeBasedConditions](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-triggerandscopebasedconditions?view=graph-rest-beta) configuration. Also see [Microsoft docs for identityGovernance.ruleBasedSubjectSet](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-rulebasedsubjectset?view=graph-rest-beta). <br> ",
										},
									},
								},
								MarkdownDescription: "Defines who the workflow runs for. <br/> A shared object that is used in entitlement management access package assignment policies, role management policies, and lifecycle workflows. <br/> This object is an abstract base type from which the following resources are derived:. Also see [Microsoft docs for subjectSet](https://learn.microsoft.com/en-us/graph/api/resources/subjectset?view=graph-rest-beta). <br> ",
							},
							"trigger": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{ // identityGovernance.workflowExecutionTrigger
									"attribute_change": generic.OdataDerivedTypeNestedAttributeRs{
										DerivedType: "#microsoft.graph.identityGovernance.attributeChangeTrigger",
										SingleNestedAttribute: schema.SingleNestedAttribute{
											Optional: true,
											Attributes: map[string]schema.Attribute{ // identityGovernance.attributeChangeTrigger
												"trigger_attributes": schema.SetNestedAttribute{
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{ // identityGovernance.triggerAttribute
															"name": schema.StringAttribute{
																MarkdownDescription: "The name of the trigger attribute that is changed to trigger an [attributeChangeTrigger](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-attributechangetrigger?view=graph-rest-beta) workflow.",
															},
														},
													},
													MarkdownDescription: "The trigger attribute being changed that triggers the workflowexecutiontrigger of a workflow.) / Defines the trigger attribute, which is changed to activate a workflow using an [attributeChangeTrigger](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-attributechangetrigger?view=graph-rest-beta). Also see [Microsoft docs for identityGovernance.triggerAttribute](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-triggerattribute?view=graph-rest-beta). <br> ",
												},
											},
											Validators: []validator.Object{
												identityGovernanceWorkflowVersionIdentityGovernanceWorkflowExecutionTriggerValidator,
											},
											MarkdownDescription: "Represents changes in user attributes that trigger the execution of workload conditions for a user. Also see [Microsoft docs for identityGovernance.attributeChangeTrigger](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-attributechangetrigger?view=graph-rest-beta). <br> ",
										},
									},
									"membership_change": generic.OdataDerivedTypeNestedAttributeRs{
										DerivedType: "#microsoft.graph.identityGovernance.membershipChangeTrigger",
										SingleNestedAttribute: schema.SingleNestedAttribute{
											Optional: true,
											Attributes: map[string]schema.Attribute{ // identityGovernance.membershipChangeTrigger
												"change_type": schema.StringAttribute{
													Validators: []validator.String{
														stringvalidator.OneOf("add", "remove", "unknownFutureValue"),
													},
													MarkdownDescription: "Defines what change that happens to the workflow group to trigger the [workflowExecutionTrigger](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-workflowexecutiontrigger?view=graph-rest-beta). <br/> _Provider_ allowed values are: `add`, `remove`, `unknownFutureValue`.",
												},
											},
											Validators: []validator.Object{
												identityGovernanceWorkflowVersionIdentityGovernanceWorkflowExecutionTriggerValidator,
											},
											MarkdownDescription: "Represents the change in group membership that triggers the execution conditions of a workflow for a user. Also see [Microsoft docs for identityGovernance.membershipChangeTrigger](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-membershipchangetrigger?view=graph-rest-beta). <br> ",
										},
									},
									"time_based_attribute": generic.OdataDerivedTypeNestedAttributeRs{
										DerivedType: "#microsoft.graph.identityGovernance.timeBasedAttributeTrigger",
										SingleNestedAttribute: schema.SingleNestedAttribute{
											Optional: true,
											Attributes: map[string]schema.Attribute{ // identityGovernance.timeBasedAttributeTrigger
												"offset_in_days": schema.Int64Attribute{
													MarkdownDescription: "How many days before or after the time-based attribute specified the workflow should trigger. For example, if the attribute is `employeeHireDate` and offsetInDays is -1, then the workflow should trigger one day before the employee hire date. The value can range between -180 and 180 days.",
												},
												"time_based_attribute": schema.StringAttribute{
													Validators: []validator.String{
														stringvalidator.OneOf("employeeHireDate", "employeeLeaveDateTime", "unknownFutureValue", "createdDateTime"),
													},
													MarkdownDescription: "Determines which time-based identity property to reference. <br/> _Provider_ allowed values are: `employeeHireDate`, `employeeLeaveDateTime`, `unknownFutureValue`, `createdDateTime`.",
												},
											},
											Validators: []validator.Object{
												identityGovernanceWorkflowVersionIdentityGovernanceWorkflowExecutionTriggerValidator,
											},
											MarkdownDescription: "Trigger based on a time-based attribute for initiating the execution of a [lifecycle workflow](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-workflow?view=graph-rest-beta). The combination of scope and trigger conditions determines when a workflow is executed and on which identities. Also see [Microsoft docs for identityGovernance.timeBasedAttributeTrigger](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-timebasedattributetrigger?view=graph-rest-beta). <br> ",
										},
									},
									"user_inactivity": generic.OdataDerivedTypeNestedAttributeRs{
										DerivedType: "#microsoft.graph.identityGovernance.userInactivityTrigger",
										SingleNestedAttribute: schema.SingleNestedAttribute{
											Optional: true,
											Attributes: map[string]schema.Attribute{ // identityGovernance.userInactivityTrigger
												"inactivity_period_in_days": schema.Int64Attribute{
													MarkdownDescription: "The number of days a user must be inactive before triggering workflow execution.",
												},
											},
											Validators: []validator.Object{
												identityGovernanceWorkflowVersionIdentityGovernanceWorkflowExecutionTriggerValidator,
											},
											MarkdownDescription: "Represents a trigger based on user inactivity that initiates workflow execution for a user. Also see [Microsoft docs for identityGovernance.userInactivityTrigger](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-userinactivitytrigger?view=graph-rest-beta). <br> ",
										},
									},
								},
								MarkdownDescription: "What triggers a workflow to run. / The workflowExecutionTrigger type represents the workflow execution trigger when the [workflow runs on schedule](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-triggerandscopebasedconditions?view=graph-rest-beta). Inherited by the following derived types:. Also see [Microsoft docs for identityGovernance.workflowExecutionTrigger](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-workflowexecutiontrigger?view=graph-rest-beta). <br> ",
							},
						},
						Validators: []validator.Object{
							identityGovernanceWorkflowVersionIdentityGovernanceWorkflowExecutionConditionsValidator,
						},
						MarkdownDescription: "Represents a lifecycle workflow running by schedule, who it runs for, and what triggers the workflow to run. Also see [Microsoft docs for identityGovernance.triggerAndScopeBasedConditions](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-triggerandscopebasedconditions?view=graph-rest-beta). <br> ",
					},
				},
			},
			MarkdownDescription: "Conditions describing when to execute the workflow and the criteria to identify in-scope subject set. / The workflowExecutionConditions type notes the workflow execution conditions in [workflowTemplate](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-workflowtemplate?view=graph-rest-beta) and [workflowBase](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-workflowbase?view=graph-rest-beta) objects. Execution conditions define when a workflow runs and rules that identify the users that are the target of the workflow. Also see [Microsoft docs for identityGovernance.workflowExecutionConditions](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-workflowexecutionconditions?view=graph-rest-beta). <br> ",
		},
		"is_enabled": schema.BoolAttribute{
			MarkdownDescription: "Whether the workflow is enabled or disabled. If this setting is `true`, the workflow can be run on demand or on schedule when **isSchedulingEnabled** is `true`. <br/> Supports `$filter`(`eq`, `ne`) and `orderBy`.",
		},
		"is_scheduling_enabled": schema.BoolAttribute{
			MarkdownDescription: "If `true`, the Lifecycle Workflow engine executes the workflow based on the schedule defined by tenant settings. Cannot be `true` for a disabled workflow (where **isEnabled** is `false`). <br/> Supports `$filter`(`eq`, `ne`) and `orderBy`.",
		},
		"last_modified_date_time": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The date time when the `workflow` was last modified. <br/> Supports `$filter`(`lt`, `le`, `gt`, `ge`, `eq`, `ne`) and `$orderby`.",
		},
		"administration_scope_targets": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // directoryObject
					"id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The unique identifier for the object. For example, 12345678-9abc-def0-1234-56789abcde. The value of the **id** property is often but not exclusively in the form of a GUID; treat it as an opaque identifier and do not rely on it being a GUID. Key. Not nullable. Read-only.",
					},
				},
			},
			MarkdownDescription: "The [administrative units](https://learn.microsoft.com/en-us/graph/api/resources/administrativeunit?view=graph-rest-beta) in the scope of the workflow. Optional. <br/> Supports `$expand`. <br/> Represents a Microsoft Entra object. The **directoryObject** type is the base type for the following directory entity types generally referred to as directory objects:. Also see [Microsoft docs for directoryObject](https://learn.microsoft.com/en-us/graph/api/resources/directoryobject?view=graph-rest-beta). <br> ",
		},
		"created_by": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // user
				"id": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "The user identifier.",
				},
			},
			MarkdownDescription: "The user who created the workflow. <br/> Supports `$filter`(`eq`, `ne`) and `$expand`. <br/> Represents an Azure Active Directory user object. Also see [Microsoft docs for user](https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-user?view=graph-rest-beta). <br> ",
		},
		"last_modified_by": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // user
				"id": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "The user identifier.",
				},
			},
			MarkdownDescription: "The user who last modified the workflow. <br/> Supports `$filter`(`eq`, `ne`) and `$expand`. <br/> Represents an Azure Active Directory user object. Also see [Microsoft docs for user](https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-user?view=graph-rest-beta). <br> ",
		},
		"tasks": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // identityGovernance.task
					"id": schema.StringAttribute{
						MarkdownDescription: "Unique identifier for the task. By default, this value will not change if a task is moved from one list to another.",
					},
					"arguments": schema.SetNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // keyValuePair
								"name": schema.StringAttribute{
									MarkdownDescription: "Name for this key-value pair",
								},
								"value": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "Value for this key-value pair",
								},
							},
						},
						MarkdownDescription: "Key value pair for storing custom settings. Also see [Microsoft docs for keyValuePair](https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyvaluepair?view=graph-rest-beta). <br> ",
					},
					"category": schema.StringAttribute{
						Validators: []validator.String{
							wpvalidator.FlagEnumValues("joiner", "leaver", "unknownFutureValue", "mover"),
						},
						MarkdownDescription: "_Provider_ allowed values are: `joiner`, `leaver`, `unknownFutureValue`, `mover`.",
					},
					"continue_on_error": schema.BoolAttribute{},
					"description": schema.StringAttribute{
						Optional: true,
					},
					"display_name": schema.StringAttribute{
						MarkdownDescription: "The name of the task.",
					},
					"execution_sequence": schema.Int64Attribute{},
					"is_enabled":         schema.BoolAttribute{},
					"task_definition_id": schema.StringAttribute{},
				},
			},
			MarkdownDescription: "The tasks in the workflow. / Represents a task, such as a piece of work or personal item, that can be tracked and completed. A **task** is always contained in a [base task list](basetasklist.md). Also see [Microsoft docs for identityGovernance.task](https://learn.microsoft.com/en-us/graph/api/resources/task?view=graph-rest-beta). <br> ",
		},
		"version_number": schema.StringAttribute{
			MarkdownDescription: "The version of the workflow.",
		},
	},
	MarkdownDescription: "Represents a version of a [lifecycle workflow](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-workflowversion?view=graph-rest-beta). Workflow versions are subsequent versions of workflows you can create when you need to change the workflow configuration other than its basic properties. You can view older versions of the workflow and associated reports will note which workflow version had been run. <br/> Also see [Microsoft docs for identityGovernance.workflowVersion](https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-workflowversion?view=graph-rest-beta). ||| MS Graph: Lifecycle workflows",
}

var identityGovernanceWorkflowVersionIdentityGovernanceWorkflowExecutionConditionsValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("on_demand_execution_only"),
	path.MatchRelative().AtParent().AtName("trigger_and_scope_based_conditions"),
)

var identityGovernanceWorkflowVersionSubjectSetValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("group"),
	path.MatchRelative().AtParent().AtName("rule"),
)

var identityGovernanceWorkflowVersionIdentityGovernanceWorkflowExecutionTriggerValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("attribute_change"),
	path.MatchRelative().AtParent().AtName("membership_change"),
	path.MatchRelative().AtParent().AtName("time_based_attribute"),
	path.MatchRelative().AtParent().AtName("user_inactivity"),
)
