package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	UnifiedRoleManagementPolicyResource = generic.GenericResource{
		TypeNameSuffix: "unified_role_management_policy",
		SpecificSchema: unifiedRoleManagementPolicyResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/policies/roleManagementPolicies",
			ReadOptions: generic.ReadOptions{
				ODataExpand: "rules,effectiveRules",
				DataSource: generic.DataSourceOptions{
					Plural: generic.PluralOptions{
						ExtraAttributes: []string{"is_organization_default", "scope_id", "scope_type"},
					},
				},
			},
			WriteOptions: generic.WriteOptions{
				UpdateInsteadOfCreate: true,
				SkipDelete:            true,
			},
		},
	}

	UnifiedRoleManagementPolicySingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&UnifiedRoleManagementPolicyResource)

	UnifiedRoleManagementPolicyPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&UnifiedRoleManagementPolicyResource, "")
)

var unifiedRoleManagementPolicyResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // unifiedRoleManagementPolicy
		"id": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Unique identifier for the policy.",
		},
		"description": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Description for the policy.",
		},
		"display_name": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Display name for the policy.",
		},
		"is_organization_default": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "This can only be set to `true` for a single tenant-wide policy which will apply to all scopes and roles. Set the scopeId to `/` and scopeType to `Directory`. Supports `$filter` (`eq`, `ne`).",
		},
		"last_modified_by": schema.SingleNestedAttribute{
			Computed: true,
			Attributes: map[string]schema.Attribute{ // identity
				"display_name": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The display name of the identity. <br/><br/>For drive items, the display name might not always be available or up to date. For example, if a user changes their display name the API might show the new value in a future response, but the items associated with the user don't show up as changed when using [delta](../api/driveitem-delta.md).",
				},
				"id": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Unique identifier for the identity or actor. For example, in the access reviews decisions API, this property might record the **id** of the principal, that is, the group, user, or application that's subject to review.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpplanmodifier.ObjectUseStateForUnknown()},
			MarkdownDescription: "The identity who last modified the role setting. / Represents an identity of an _actor_. For example, an actor can be a user, device, or application. Multiple Microsoft Graph APIs share this resource and the data they return varies depending on the API.\n\nBase type of [userIdentity](useridentity.md). / https://learn.microsoft.com/en-us/graph/api/resources/identity?view=graph-rest-beta",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The time when the role setting was last modified.",
		},
		"scope_id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The identifier of the scope where the policy is created. Can be `/` for the tenant or a group ID. Required.",
		},
		"scope_type": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The type of the scope where the policy is created. One of `Directory`, `DirectoryRole`, `Group`. Required.",
		},
		"effective_rules": schema.SetNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // unifiedRoleManagementPolicyRule
					"id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Identifier for the rule. Read-only.",
					},
					"target": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{ // unifiedRoleManagementPolicyRuleTarget
							"caller": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "The type of caller that's the target of the policy rule. Allowed values are: `None`, `Admin`, `EndUser`.",
							},
							"enforced_settings": schema.SetAttribute{
								ElementType:         types.StringType,
								Computed:            true,
								MarkdownDescription: "The list of role settings that are enforced and cannot be overridden by child scopes. Use `All` for all settings.",
							},
							"inheritable_settings": schema.SetAttribute{
								ElementType:         types.StringType,
								Computed:            true,
								MarkdownDescription: "The list of role settings that can be inherited by child scopes. Use `All` for all settings.",
							},
							"level": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "The role assignment type that's the target of policy rule. Allowed values are: `Eligibility`, `Assignment`.",
							},
							"operations": schema.SetAttribute{
								ElementType:         types.StringType,
								Computed:            true,
								MarkdownDescription: "The role management operations that are the target of the policy rule. Allowed values are: `All`, `Activate`, `Deactivate`, `Assign`, `Update`, `Remove`, `Extend`, `Renew`.",
							},
						},
						MarkdownDescription: "**Not implemented.** Defines details of scope that's targeted by role management policy rule. The details can include the principal type, the role assignment type, and actions affecting a role. Supports `$filter` (`eq`, `ne`). / Defines details of the scope that's targeted by role management policy rule. The details can include the principal type, the role assignment type, and actions affecting a role. / https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyruletarget?view=graph-rest-beta",
					},
					"approval": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.unifiedRoleManagementPolicyApprovalRule",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{ // unifiedRoleManagementPolicyApprovalRule
								"setting": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{ // approvalSettings
										"approval_mode": schema.StringAttribute{
											Computed:            true,
											MarkdownDescription: "One of `SingleStage`, `Serial`, `Parallel`, `NoApproval` (default). `NoApproval` is used when `isApprovalRequired` is `false`.",
										},
										"approval_stages": schema.SetNestedAttribute{
											Computed: true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{ // approvalStage
													"approval_stage_time_out_in_days": schema.Int64Attribute{
														Computed:            true,
														MarkdownDescription: "The number of days that a request can be pending a response before it is automatically denied.",
													},
													"escalation_approvers": schema.SetNestedAttribute{
														Computed: true,
														NestedObject: schema.NestedAttributeObject{
															Attributes: map[string]schema.Attribute{ // userSet
																"is_backup": schema.BoolAttribute{
																	Computed:            true,
																	MarkdownDescription: "For a user in an approval stage, this property indicates whether the user is a backup fallback approver.",
																},
																"connected_organization_members": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.connectedOrganizationMembers",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Computed: true,
																		Attributes: map[string]schema.Attribute{ // connectedOrganizationMembers
																			"id": schema.StringAttribute{
																				Computed:            true,
																				MarkdownDescription: "The ID of the connected organization in entitlement management.",
																			},
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the request settings of an [access package assignment policy](accesspackageassignmentpolicy.md). The `@odata.type` value `#microsoft.graph.connectedOrganizationMembers` indicates that this type identifies a collection of users who are associated with a [connected organization](connectedorganization.md) and are allowed to request an access package. / https://learn.microsoft.com/en-us/graph/api/resources/connectedorganizationmembers?view=graph-rest-beta",
																	},
																},
																"external_sponsors": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.externalSponsors",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Computed:   true,
																		Attributes: map[string]schema.Attribute{ // externalSponsors
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the approval stage of an [access package assignment policy](accesspackageassignmentpolicy.md).\nIt's a subtype of [userSet](userset.md), in which the `@odata.type` value `#microsoft.graph.externalSponsors` indicates that a requesting user's connected organization external sponsors are to be the approver. This approver is only applicable to requests from users who are part of a connected organization.  When creating an access package assignment policy approval stage with externalSponsors, also include another approver, such as a single user or group member, in case the connected organization doesn't have an external sponsor. / https://learn.microsoft.com/en-us/graph/api/resources/externalsponsors?view=graph-rest-beta",
																	},
																},
																"group_members": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.groupMembers",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Computed: true,
																		Attributes: map[string]schema.Attribute{ // groupMembers
																			"id": schema.StringAttribute{
																				Computed:            true,
																				MarkdownDescription: "The ID of the group in Microsoft Entra ID.",
																			},
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the request, approval, and assignment review settings of an [access package assignment policy](accesspackageassignmentpolicy.md).\nThe `@odata.type` value `#microsoft.graph.groupMembers` indicates that this type identifies a collection of users in the tenant who are allowed as requestor, approver, or reviewer, who are the members of a specific group. / https://learn.microsoft.com/en-us/graph/api/resources/groupmembers?view=graph-rest-beta",
																	},
																},
																"internal_sponsors": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.internalSponsors",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Computed:   true,
																		Attributes: map[string]schema.Attribute{ // internalSponsors
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the approval stage of an [access package assignment policy](accesspackageassignmentpolicy.md).\nIt's a subtype of [userSet](userset.md), in which the `@odata.type` value `#microsoft.graph.internalSponsors` indicates that a requesting user's connected organization internal sponsors are to be the approver. This approver is only applicable to requests from users who are part of a connected organization.  When creating an access package assignment policy approval stage with internalSponsors, also include another approver, such as a single user or group member, in case the connected organization doesn't have an internal sponsor. / https://learn.microsoft.com/en-us/graph/api/resources/internalsponsors?view=graph-rest-beta",
																	},
																},
																"requestor_manager": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.requestorManager",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Computed: true,
																		Attributes: map[string]schema.Attribute{ // requestorManager
																			"manager_level": schema.Int64Attribute{
																				Computed:            true,
																				MarkdownDescription: "The hierarchical level of the manager with respect to the requestor. For example, the direct manager of a requestor would have a managerLevel of 1, while the manager of the requestor's manager would have a managerLevel of 2. Default value for managerLevel is 1. Possible values for this property range from 1 to 2.",
																			},
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the approval settings of an [access package assignment policy](accesspackageassignmentpolicy.md).\nIt's a subtype of [userSet](userset.md), in which the `@odata.type` value `#microsoft.graph.requestorManager` indicates that a requesting user's manager is to be the approver. Include another approver When creating an access package assignment policy approval stage with requestorManager, in case the requesting user doesn't have a manager. Including another approver, such as a single user or group member, covers the case where the requesting user doesn't have a manager. / https://learn.microsoft.com/en-us/graph/api/resources/requestormanager?view=graph-rest-beta",
																	},
																},
																"single_user": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.singleUser",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Computed: true,
																		Attributes: map[string]schema.Attribute{ // singleUser
																			"id": schema.StringAttribute{
																				Computed:            true,
																				MarkdownDescription: "The ID of the user in Microsoft Entra ID.",
																			},
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the request, approval, and assignment review settings of an [access package assignment policy](accesspackageassignmentpolicy.md). The  `@odata.type` value `#microsoft.graph.singleUser` indicates that this userSet identifies a specific user in the tenant who will be allowed as a requestor, approver, or reviewer. / https://learn.microsoft.com/en-us/graph/api/resources/singleuser?view=graph-rest-beta",
																	},
																},
																"target_user_sponsors": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.targetUserSponsors",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Computed:   true,
																		Attributes: map[string]schema.Attribute{ // targetUserSponsors
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the approval settings of an [access package assignment policy](accesspackageassignmentpolicy.md).\nIt's a subtype of [userSet](userset.md), in which the `@odata.type` value `#microsoft.graph.targetUserSponsors` indicates that a requesting user's sponsors are the approvers. When creating an access package assignment policy approval stage with **targetUserSponsors**, also include another approver, such as a single user or group member, in case the requesting user doesn't have sponsors. / https://learn.microsoft.com/en-us/graph/api/resources/targetusersponsors?view=graph-rest-beta",
																	},
																},
															},
														},
														MarkdownDescription: "The users who are asked to approve requests if escalation is enabled and the primary approvers don't respond before the escalation time. This property can be a collection of [singleUser](singleuser.md), [groupMembers](groupmembers.md), [requestorManager](requestormanager.md), [internalSponsors](internalsponsors.md), and [externalSponsors](externalsponsors.md). When you create or update a [policy](accesspackageassignmentpolicy.md), if there are no escalation approvers, or escalation approvers aren't required for the stage, assign an empty collection to this property. / Used in the request, approval, and assignment review settings of an [access package assignment policy](accesspackageassignmentpolicy.md). It's an abstract base type inherited by the following resource types: / https://learn.microsoft.com/en-us/graph/api/resources/userset?view=graph-rest-beta",
													},
													"escalation_time_in_minutes": schema.Int64Attribute{
														Computed:            true,
														MarkdownDescription: "If escalation is required, the time a request can be pending a response from a primary approver.",
													},
													"is_approver_justification_required": schema.BoolAttribute{
														Computed:            true,
														MarkdownDescription: "Indicates whether the approver is required to provide a justification for approving a request.",
													},
													"is_escalation_enabled": schema.BoolAttribute{
														Computed:            true,
														MarkdownDescription: "If true, then one or more escalation approvers are configured in this approval stage.",
													},
													"primary_approvers": schema.SetNestedAttribute{
														Computed: true,
														NestedObject: schema.NestedAttributeObject{
															Attributes: map[string]schema.Attribute{ // userSet
																"is_backup": schema.BoolAttribute{
																	Computed:            true,
																	MarkdownDescription: "For a user in an approval stage, this property indicates whether the user is a backup fallback approver.",
																},
																"connected_organization_members": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.connectedOrganizationMembers",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Computed: true,
																		Attributes: map[string]schema.Attribute{ // connectedOrganizationMembers
																			"id": schema.StringAttribute{
																				Computed:            true,
																				MarkdownDescription: "The ID of the connected organization in entitlement management.",
																			},
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the request settings of an [access package assignment policy](accesspackageassignmentpolicy.md). The `@odata.type` value `#microsoft.graph.connectedOrganizationMembers` indicates that this type identifies a collection of users who are associated with a [connected organization](connectedorganization.md) and are allowed to request an access package. / https://learn.microsoft.com/en-us/graph/api/resources/connectedorganizationmembers?view=graph-rest-beta",
																	},
																},
																"external_sponsors": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.externalSponsors",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Computed:   true,
																		Attributes: map[string]schema.Attribute{ // externalSponsors
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the approval stage of an [access package assignment policy](accesspackageassignmentpolicy.md).\nIt's a subtype of [userSet](userset.md), in which the `@odata.type` value `#microsoft.graph.externalSponsors` indicates that a requesting user's connected organization external sponsors are to be the approver. This approver is only applicable to requests from users who are part of a connected organization.  When creating an access package assignment policy approval stage with externalSponsors, also include another approver, such as a single user or group member, in case the connected organization doesn't have an external sponsor. / https://learn.microsoft.com/en-us/graph/api/resources/externalsponsors?view=graph-rest-beta",
																	},
																},
																"group_members": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.groupMembers",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Computed: true,
																		Attributes: map[string]schema.Attribute{ // groupMembers
																			"id": schema.StringAttribute{
																				Computed:            true,
																				MarkdownDescription: "The ID of the group in Microsoft Entra ID.",
																			},
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the request, approval, and assignment review settings of an [access package assignment policy](accesspackageassignmentpolicy.md).\nThe `@odata.type` value `#microsoft.graph.groupMembers` indicates that this type identifies a collection of users in the tenant who are allowed as requestor, approver, or reviewer, who are the members of a specific group. / https://learn.microsoft.com/en-us/graph/api/resources/groupmembers?view=graph-rest-beta",
																	},
																},
																"internal_sponsors": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.internalSponsors",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Computed:   true,
																		Attributes: map[string]schema.Attribute{ // internalSponsors
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the approval stage of an [access package assignment policy](accesspackageassignmentpolicy.md).\nIt's a subtype of [userSet](userset.md), in which the `@odata.type` value `#microsoft.graph.internalSponsors` indicates that a requesting user's connected organization internal sponsors are to be the approver. This approver is only applicable to requests from users who are part of a connected organization.  When creating an access package assignment policy approval stage with internalSponsors, also include another approver, such as a single user or group member, in case the connected organization doesn't have an internal sponsor. / https://learn.microsoft.com/en-us/graph/api/resources/internalsponsors?view=graph-rest-beta",
																	},
																},
																"requestor_manager": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.requestorManager",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Computed: true,
																		Attributes: map[string]schema.Attribute{ // requestorManager
																			"manager_level": schema.Int64Attribute{
																				Computed:            true,
																				MarkdownDescription: "The hierarchical level of the manager with respect to the requestor. For example, the direct manager of a requestor would have a managerLevel of 1, while the manager of the requestor's manager would have a managerLevel of 2. Default value for managerLevel is 1. Possible values for this property range from 1 to 2.",
																			},
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the approval settings of an [access package assignment policy](accesspackageassignmentpolicy.md).\nIt's a subtype of [userSet](userset.md), in which the `@odata.type` value `#microsoft.graph.requestorManager` indicates that a requesting user's manager is to be the approver. Include another approver When creating an access package assignment policy approval stage with requestorManager, in case the requesting user doesn't have a manager. Including another approver, such as a single user or group member, covers the case where the requesting user doesn't have a manager. / https://learn.microsoft.com/en-us/graph/api/resources/requestormanager?view=graph-rest-beta",
																	},
																},
																"single_user": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.singleUser",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Computed: true,
																		Attributes: map[string]schema.Attribute{ // singleUser
																			"id": schema.StringAttribute{
																				Computed:            true,
																				MarkdownDescription: "The ID of the user in Microsoft Entra ID.",
																			},
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the request, approval, and assignment review settings of an [access package assignment policy](accesspackageassignmentpolicy.md). The  `@odata.type` value `#microsoft.graph.singleUser` indicates that this userSet identifies a specific user in the tenant who will be allowed as a requestor, approver, or reviewer. / https://learn.microsoft.com/en-us/graph/api/resources/singleuser?view=graph-rest-beta",
																	},
																},
																"target_user_sponsors": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.targetUserSponsors",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Computed:   true,
																		Attributes: map[string]schema.Attribute{ // targetUserSponsors
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the approval settings of an [access package assignment policy](accesspackageassignmentpolicy.md).\nIt's a subtype of [userSet](userset.md), in which the `@odata.type` value `#microsoft.graph.targetUserSponsors` indicates that a requesting user's sponsors are the approvers. When creating an access package assignment policy approval stage with **targetUserSponsors**, also include another approver, such as a single user or group member, in case the requesting user doesn't have sponsors. / https://learn.microsoft.com/en-us/graph/api/resources/targetusersponsors?view=graph-rest-beta",
																	},
																},
															},
														},
														MarkdownDescription: "The users who are asked to approve requests. A collection of [singleUser](singleuser.md), [groupMembers](groupmembers.md), [requestorManager](requestormanager.md), [internalSponsors](internalsponsors.md), [externalSponsors](externalsponsors.md), and [targetUserSponsors](targetusersponsors.md). When creating or updating a [policy](accesspackageassignmentpolicy.md), include at least one **userSet** in this collection. / Used in the request, approval, and assignment review settings of an [access package assignment policy](accesspackageassignmentpolicy.md). It's an abstract base type inherited by the following resource types: / https://learn.microsoft.com/en-us/graph/api/resources/userset?view=graph-rest-beta",
													},
												},
											},
											MarkdownDescription: "If approval is required, the one or two elements of this collection define each of the stages of approval. An empty array if no approval is required. / In entitlement management, used for the **approvalStages** property of approval settings in the **requestApprovalSettings** property of an [access package assignment policy](accesspackageassignmentpolicy.md). Specifies the primary, fallback, and escalation approvers of each stage.\n\nIn PIM, defines the settings of the approval stages in a [unifiedRoleManagementPolicyApprovalRule](unifiedrolemanagementpolicyapprovalrule.md) object. Specifies the primary and escalation approvers of each stage and whether approvals and escalations are required. / https://learn.microsoft.com/en-us/graph/api/resources/approvalstage?view=graph-rest-beta",
										},
										"is_approval_required": schema.BoolAttribute{
											Computed:            true,
											MarkdownDescription: "Indicates whether approval is required for requests in this policy.",
										},
										"is_approval_required_for_extension": schema.BoolAttribute{
											Computed:            true,
											MarkdownDescription: "Indicates whether approval is required for a user to extend their assignment.",
										},
										"is_requestor_justification_required": schema.BoolAttribute{
											Computed:            true,
											MarkdownDescription: "Indicates whether the requestor is required to supply a justification in their request.",
										},
									},
									MarkdownDescription: "The settings for approval of the role assignment. / The settings for approval as defined in a role management policy rule. / https://learn.microsoft.com/en-us/graph/api/resources/approvalsettings?view=graph-rest-beta",
								},
							},
							Validators: []validator.Object{
								unifiedRoleManagementPolicyUnifiedRoleManagementPolicyRuleValidator,
							},
							MarkdownDescription: "A type derived from the [unifiedRoleManagementPolicyRule](../resources/unifiedrolemanagementpolicyrule.md) resource type that defines rules for approving a role assignment. / https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyapprovalrule?view=graph-rest-beta",
						},
					},
					"authentication_context": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.unifiedRoleManagementPolicyAuthenticationContextRule",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{ // unifiedRoleManagementPolicyAuthenticationContextRule
								"claim_value": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The value of the authentication context claim.",
								},
								"is_enabled": schema.BoolAttribute{
									Computed:            true,
									MarkdownDescription: "Whether this rule is enabled.",
								},
							},
							Validators: []validator.Object{
								unifiedRoleManagementPolicyUnifiedRoleManagementPolicyRuleValidator,
							},
							MarkdownDescription: "A type derived from the [unifiedRoleManagementPolicyRule](../resources/unifiedrolemanagementpolicyrule.md) resource type that defines the authentication context rule for the conditional access policy associated with a role management policy. / https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyauthenticationcontextrule?view=graph-rest-beta",
						},
					},
					"enablement": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.unifiedRoleManagementPolicyEnablementRule",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{ // unifiedRoleManagementPolicyEnablementRule
								"enabled_rules": schema.SetAttribute{
									ElementType:         types.StringType,
									Computed:            true,
									MarkdownDescription: "The collection of rules that are enabled for this policy rule. For example, `MultiFactorAuthentication`, `Ticketing`, and `Justification`.",
								},
							},
							Validators: []validator.Object{
								unifiedRoleManagementPolicyUnifiedRoleManagementPolicyRuleValidator,
							},
							MarkdownDescription: "A type derived from the [unifiedRoleManagementPolicyRule](../resources/unifiedrolemanagementpolicyrule.md) resource type that defines the rules to enable the assignment, for example, enable MFA, justification on assignments or ticketing information. / https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyenablementrule?view=graph-rest-beta",
						},
					},
					"expiration": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.unifiedRoleManagementPolicyExpirationRule",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{ // unifiedRoleManagementPolicyExpirationRule
								"is_expiration_required": schema.BoolAttribute{
									Computed:            true,
									MarkdownDescription: "Indicates whether expiration is required or if it's a permanently active assignment or eligibility.",
								},
								"maximum_duration": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The maximum duration allowed for eligibility or assignment that isn't permanent. Required when **isExpirationRequired** is `true`.",
								},
							},
							Validators: []validator.Object{
								unifiedRoleManagementPolicyUnifiedRoleManagementPolicyRuleValidator,
							},
							MarkdownDescription: "A type derived from the [unifiedRoleManagementPolicyRule](../resources/unifiedrolemanagementpolicyrule.md) resource type that defines the maximum duration a role can be assigned to a principal (either through direct assignment or through activation of eligibility / https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyexpirationrule?view=graph-rest-beta",
						},
					},
					"notification": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.unifiedRoleManagementPolicyNotificationRule",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{ // unifiedRoleManagementPolicyNotificationRule
								"is_default_recipients_enabled": schema.BoolAttribute{
									Computed:            true,
									MarkdownDescription: "Indicates whether a default recipient will receive the notification email.",
								},
								"notification_level": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The level of notification. The possible values are `None`, `Critical`, `All`.",
								},
								"notification_recipients": schema.SetAttribute{
									ElementType:         types.StringType,
									Computed:            true,
									MarkdownDescription: "The list of recipients of the email notifications.",
								},
								"notification_type": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The type of notification. Only `Email` is supported.",
								},
								"recipient_type": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The type of recipient of the notification. The possible values are `Requestor`, `Approver`, `Admin`.",
								},
							},
							Validators: []validator.Object{
								unifiedRoleManagementPolicyUnifiedRoleManagementPolicyRuleValidator,
							},
							MarkdownDescription: "A type derived from the [unifiedRoleManagementPolicyRule](../resources/unifiedrolemanagementpolicyrule.md) resource type that defines the email notification rules for role assignments, activations, and approvals. / https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicynotificationrule?view=graph-rest-beta",
						},
					},
				},
			},
			PlanModifiers:       []planmodifier.Set{wpplanmodifier.SetUseStateForUnknown()},
			MarkdownDescription: "The list of effective rules like approval rules and expiration rules evaluated based on inherited referenced rules. For example, if there is a tenant-wide policy to enforce enabling an approval rule, the effective rule will be to enable approval even if the policy has a rule to disable approval. Supports `$expand`. / An abstract type that defines the rules associated with role management policies. This abstract type is inherited by the following resources that define the various types of rules and their settings associated with role management policies. / https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyrule?view=graph-rest-beta",
		},
		"rules": schema.SetNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // unifiedRoleManagementPolicyRule
					"id": schema.StringAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
						Computed:            true,
						MarkdownDescription: "Identifier for the rule. Read-only. The _provider_ default value is `\"\"`.",
					},
					"target": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // unifiedRoleManagementPolicyRuleTarget
							"caller": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "The type of caller that's the target of the policy rule. Allowed values are: `None`, `Admin`, `EndUser`.",
							},
							"enforced_settings": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
								Computed:            true,
								MarkdownDescription: "The list of role settings that are enforced and cannot be overridden by child scopes. Use `All` for all settings. The _provider_ default value is `[]`.",
							},
							"inheritable_settings": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
								Computed:            true,
								MarkdownDescription: "The list of role settings that can be inherited by child scopes. Use `All` for all settings. The _provider_ default value is `[]`.",
							},
							"level": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "The role assignment type that's the target of policy rule. Allowed values are: `Eligibility`, `Assignment`.",
							},
							"operations": schema.SetAttribute{
								ElementType:         types.StringType,
								Optional:            true,
								PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
								Computed:            true,
								MarkdownDescription: "The role management operations that are the target of the policy rule. Allowed values are: `All`, `Activate`, `Deactivate`, `Assign`, `Update`, `Remove`, `Extend`, `Renew`. The _provider_ default value is `[]`.",
							},
						},
						MarkdownDescription: "**Not implemented.** Defines details of scope that's targeted by role management policy rule. The details can include the principal type, the role assignment type, and actions affecting a role. Supports `$filter` (`eq`, `ne`). / Defines details of the scope that's targeted by role management policy rule. The details can include the principal type, the role assignment type, and actions affecting a role. / https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyruletarget?view=graph-rest-beta",
					},
					"approval": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.unifiedRoleManagementPolicyApprovalRule",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // unifiedRoleManagementPolicyApprovalRule
								"setting": schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // approvalSettings
										"approval_mode": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "One of `SingleStage`, `Serial`, `Parallel`, `NoApproval` (default). `NoApproval` is used when `isApprovalRequired` is `false`.",
										},
										"approval_stages": schema.SetNestedAttribute{
											Optional: true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{ // approvalStage
													"approval_stage_time_out_in_days": schema.Int64Attribute{
														Optional:            true,
														MarkdownDescription: "The number of days that a request can be pending a response before it is automatically denied.",
													},
													"escalation_approvers": schema.SetNestedAttribute{
														Optional: true,
														NestedObject: schema.NestedAttributeObject{
															Attributes: map[string]schema.Attribute{ // userSet
																"is_backup": schema.BoolAttribute{
																	Optional:            true,
																	PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
																	Computed:            true,
																	MarkdownDescription: "For a user in an approval stage, this property indicates whether the user is a backup fallback approver. The _provider_ default value is `false`.",
																},
																"connected_organization_members": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.connectedOrganizationMembers",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Optional: true,
																		Attributes: map[string]schema.Attribute{ // connectedOrganizationMembers
																			"id": schema.StringAttribute{
																				Required:            true,
																				MarkdownDescription: "The ID of the connected organization in entitlement management.",
																			},
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the request settings of an [access package assignment policy](accesspackageassignmentpolicy.md). The `@odata.type` value `#microsoft.graph.connectedOrganizationMembers` indicates that this type identifies a collection of users who are associated with a [connected organization](connectedorganization.md) and are allowed to request an access package. / https://learn.microsoft.com/en-us/graph/api/resources/connectedorganizationmembers?view=graph-rest-beta",
																	},
																},
																"external_sponsors": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.externalSponsors",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Optional:   true,
																		Attributes: map[string]schema.Attribute{ // externalSponsors
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the approval stage of an [access package assignment policy](accesspackageassignmentpolicy.md).\nIt's a subtype of [userSet](userset.md), in which the `@odata.type` value `#microsoft.graph.externalSponsors` indicates that a requesting user's connected organization external sponsors are to be the approver. This approver is only applicable to requests from users who are part of a connected organization.  When creating an access package assignment policy approval stage with externalSponsors, also include another approver, such as a single user or group member, in case the connected organization doesn't have an external sponsor. / https://learn.microsoft.com/en-us/graph/api/resources/externalsponsors?view=graph-rest-beta",
																	},
																},
																"group_members": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.groupMembers",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Optional: true,
																		Attributes: map[string]schema.Attribute{ // groupMembers
																			"id": schema.StringAttribute{
																				Required:            true,
																				MarkdownDescription: "The ID of the group in Microsoft Entra ID.",
																			},
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the request, approval, and assignment review settings of an [access package assignment policy](accesspackageassignmentpolicy.md).\nThe `@odata.type` value `#microsoft.graph.groupMembers` indicates that this type identifies a collection of users in the tenant who are allowed as requestor, approver, or reviewer, who are the members of a specific group. / https://learn.microsoft.com/en-us/graph/api/resources/groupmembers?view=graph-rest-beta",
																	},
																},
																"internal_sponsors": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.internalSponsors",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Optional:   true,
																		Attributes: map[string]schema.Attribute{ // internalSponsors
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the approval stage of an [access package assignment policy](accesspackageassignmentpolicy.md).\nIt's a subtype of [userSet](userset.md), in which the `@odata.type` value `#microsoft.graph.internalSponsors` indicates that a requesting user's connected organization internal sponsors are to be the approver. This approver is only applicable to requests from users who are part of a connected organization.  When creating an access package assignment policy approval stage with internalSponsors, also include another approver, such as a single user or group member, in case the connected organization doesn't have an internal sponsor. / https://learn.microsoft.com/en-us/graph/api/resources/internalsponsors?view=graph-rest-beta",
																	},
																},
																"requestor_manager": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.requestorManager",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Optional: true,
																		Attributes: map[string]schema.Attribute{ // requestorManager
																			"manager_level": schema.Int64Attribute{
																				Optional:            true,
																				MarkdownDescription: "The hierarchical level of the manager with respect to the requestor. For example, the direct manager of a requestor would have a managerLevel of 1, while the manager of the requestor's manager would have a managerLevel of 2. Default value for managerLevel is 1. Possible values for this property range from 1 to 2.",
																			},
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the approval settings of an [access package assignment policy](accesspackageassignmentpolicy.md).\nIt's a subtype of [userSet](userset.md), in which the `@odata.type` value `#microsoft.graph.requestorManager` indicates that a requesting user's manager is to be the approver. Include another approver When creating an access package assignment policy approval stage with requestorManager, in case the requesting user doesn't have a manager. Including another approver, such as a single user or group member, covers the case where the requesting user doesn't have a manager. / https://learn.microsoft.com/en-us/graph/api/resources/requestormanager?view=graph-rest-beta",
																	},
																},
																"single_user": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.singleUser",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Optional: true,
																		Attributes: map[string]schema.Attribute{ // singleUser
																			"id": schema.StringAttribute{
																				Required:            true,
																				MarkdownDescription: "The ID of the user in Microsoft Entra ID.",
																			},
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the request, approval, and assignment review settings of an [access package assignment policy](accesspackageassignmentpolicy.md). The  `@odata.type` value `#microsoft.graph.singleUser` indicates that this userSet identifies a specific user in the tenant who will be allowed as a requestor, approver, or reviewer. / https://learn.microsoft.com/en-us/graph/api/resources/singleuser?view=graph-rest-beta",
																	},
																},
																"target_user_sponsors": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.targetUserSponsors",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Optional:   true,
																		Attributes: map[string]schema.Attribute{ // targetUserSponsors
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the approval settings of an [access package assignment policy](accesspackageassignmentpolicy.md).\nIt's a subtype of [userSet](userset.md), in which the `@odata.type` value `#microsoft.graph.targetUserSponsors` indicates that a requesting user's sponsors are the approvers. When creating an access package assignment policy approval stage with **targetUserSponsors**, also include another approver, such as a single user or group member, in case the requesting user doesn't have sponsors. / https://learn.microsoft.com/en-us/graph/api/resources/targetusersponsors?view=graph-rest-beta",
																	},
																},
															},
														},
														PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
														Computed:            true,
														MarkdownDescription: "The users who are asked to approve requests if escalation is enabled and the primary approvers don't respond before the escalation time. This property can be a collection of [singleUser](singleuser.md), [groupMembers](groupmembers.md), [requestorManager](requestormanager.md), [internalSponsors](internalsponsors.md), and [externalSponsors](externalsponsors.md). When you create or update a [policy](accesspackageassignmentpolicy.md), if there are no escalation approvers, or escalation approvers aren't required for the stage, assign an empty collection to this property. / Used in the request, approval, and assignment review settings of an [access package assignment policy](accesspackageassignmentpolicy.md). It's an abstract base type inherited by the following resource types: / https://learn.microsoft.com/en-us/graph/api/resources/userset?view=graph-rest-beta. The _provider_ default value is `[]`.",
													},
													"escalation_time_in_minutes": schema.Int64Attribute{
														Optional:            true,
														MarkdownDescription: "If escalation is required, the time a request can be pending a response from a primary approver.",
													},
													"is_approver_justification_required": schema.BoolAttribute{
														Optional:            true,
														MarkdownDescription: "Indicates whether the approver is required to provide a justification for approving a request.",
													},
													"is_escalation_enabled": schema.BoolAttribute{
														Optional:            true,
														MarkdownDescription: "If true, then one or more escalation approvers are configured in this approval stage.",
													},
													"primary_approvers": schema.SetNestedAttribute{
														Optional: true,
														NestedObject: schema.NestedAttributeObject{
															Attributes: map[string]schema.Attribute{ // userSet
																"is_backup": schema.BoolAttribute{
																	Optional:            true,
																	PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
																	Computed:            true,
																	MarkdownDescription: "For a user in an approval stage, this property indicates whether the user is a backup fallback approver. The _provider_ default value is `false`.",
																},
																"connected_organization_members": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.connectedOrganizationMembers",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Optional: true,
																		Attributes: map[string]schema.Attribute{ // connectedOrganizationMembers
																			"id": schema.StringAttribute{
																				Required:            true,
																				MarkdownDescription: "The ID of the connected organization in entitlement management.",
																			},
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the request settings of an [access package assignment policy](accesspackageassignmentpolicy.md). The `@odata.type` value `#microsoft.graph.connectedOrganizationMembers` indicates that this type identifies a collection of users who are associated with a [connected organization](connectedorganization.md) and are allowed to request an access package. / https://learn.microsoft.com/en-us/graph/api/resources/connectedorganizationmembers?view=graph-rest-beta",
																	},
																},
																"external_sponsors": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.externalSponsors",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Optional:   true,
																		Attributes: map[string]schema.Attribute{ // externalSponsors
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the approval stage of an [access package assignment policy](accesspackageassignmentpolicy.md).\nIt's a subtype of [userSet](userset.md), in which the `@odata.type` value `#microsoft.graph.externalSponsors` indicates that a requesting user's connected organization external sponsors are to be the approver. This approver is only applicable to requests from users who are part of a connected organization.  When creating an access package assignment policy approval stage with externalSponsors, also include another approver, such as a single user or group member, in case the connected organization doesn't have an external sponsor. / https://learn.microsoft.com/en-us/graph/api/resources/externalsponsors?view=graph-rest-beta",
																	},
																},
																"group_members": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.groupMembers",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Optional: true,
																		Attributes: map[string]schema.Attribute{ // groupMembers
																			"id": schema.StringAttribute{
																				Required:            true,
																				MarkdownDescription: "The ID of the group in Microsoft Entra ID.",
																			},
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the request, approval, and assignment review settings of an [access package assignment policy](accesspackageassignmentpolicy.md).\nThe `@odata.type` value `#microsoft.graph.groupMembers` indicates that this type identifies a collection of users in the tenant who are allowed as requestor, approver, or reviewer, who are the members of a specific group. / https://learn.microsoft.com/en-us/graph/api/resources/groupmembers?view=graph-rest-beta",
																	},
																},
																"internal_sponsors": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.internalSponsors",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Optional:   true,
																		Attributes: map[string]schema.Attribute{ // internalSponsors
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the approval stage of an [access package assignment policy](accesspackageassignmentpolicy.md).\nIt's a subtype of [userSet](userset.md), in which the `@odata.type` value `#microsoft.graph.internalSponsors` indicates that a requesting user's connected organization internal sponsors are to be the approver. This approver is only applicable to requests from users who are part of a connected organization.  When creating an access package assignment policy approval stage with internalSponsors, also include another approver, such as a single user or group member, in case the connected organization doesn't have an internal sponsor. / https://learn.microsoft.com/en-us/graph/api/resources/internalsponsors?view=graph-rest-beta",
																	},
																},
																"requestor_manager": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.requestorManager",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Optional: true,
																		Attributes: map[string]schema.Attribute{ // requestorManager
																			"manager_level": schema.Int64Attribute{
																				Optional:            true,
																				MarkdownDescription: "The hierarchical level of the manager with respect to the requestor. For example, the direct manager of a requestor would have a managerLevel of 1, while the manager of the requestor's manager would have a managerLevel of 2. Default value for managerLevel is 1. Possible values for this property range from 1 to 2.",
																			},
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the approval settings of an [access package assignment policy](accesspackageassignmentpolicy.md).\nIt's a subtype of [userSet](userset.md), in which the `@odata.type` value `#microsoft.graph.requestorManager` indicates that a requesting user's manager is to be the approver. Include another approver When creating an access package assignment policy approval stage with requestorManager, in case the requesting user doesn't have a manager. Including another approver, such as a single user or group member, covers the case where the requesting user doesn't have a manager. / https://learn.microsoft.com/en-us/graph/api/resources/requestormanager?view=graph-rest-beta",
																	},
																},
																"single_user": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.singleUser",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Optional: true,
																		Attributes: map[string]schema.Attribute{ // singleUser
																			"id": schema.StringAttribute{
																				Required:            true,
																				MarkdownDescription: "The ID of the user in Microsoft Entra ID.",
																			},
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the request, approval, and assignment review settings of an [access package assignment policy](accesspackageassignmentpolicy.md). The  `@odata.type` value `#microsoft.graph.singleUser` indicates that this userSet identifies a specific user in the tenant who will be allowed as a requestor, approver, or reviewer. / https://learn.microsoft.com/en-us/graph/api/resources/singleuser?view=graph-rest-beta",
																	},
																},
																"target_user_sponsors": generic.OdataDerivedTypeNestedAttributeRs{
																	DerivedType: "#microsoft.graph.targetUserSponsors",
																	SingleNestedAttribute: schema.SingleNestedAttribute{
																		Optional:   true,
																		Attributes: map[string]schema.Attribute{ // targetUserSponsors
																		},
																		Validators:          []validator.Object{unifiedRoleManagementPolicyUserSetValidator},
																		MarkdownDescription: "Used in the approval settings of an [access package assignment policy](accesspackageassignmentpolicy.md).\nIt's a subtype of [userSet](userset.md), in which the `@odata.type` value `#microsoft.graph.targetUserSponsors` indicates that a requesting user's sponsors are the approvers. When creating an access package assignment policy approval stage with **targetUserSponsors**, also include another approver, such as a single user or group member, in case the requesting user doesn't have sponsors. / https://learn.microsoft.com/en-us/graph/api/resources/targetusersponsors?view=graph-rest-beta",
																	},
																},
															},
														},
														PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
														Computed:            true,
														MarkdownDescription: "The users who are asked to approve requests. A collection of [singleUser](singleuser.md), [groupMembers](groupmembers.md), [requestorManager](requestormanager.md), [internalSponsors](internalsponsors.md), [externalSponsors](externalsponsors.md), and [targetUserSponsors](targetusersponsors.md). When creating or updating a [policy](accesspackageassignmentpolicy.md), include at least one **userSet** in this collection. / Used in the request, approval, and assignment review settings of an [access package assignment policy](accesspackageassignmentpolicy.md). It's an abstract base type inherited by the following resource types: / https://learn.microsoft.com/en-us/graph/api/resources/userset?view=graph-rest-beta. The _provider_ default value is `[]`.",
													},
												},
											},
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: "If approval is required, the one or two elements of this collection define each of the stages of approval. An empty array if no approval is required. / In entitlement management, used for the **approvalStages** property of approval settings in the **requestApprovalSettings** property of an [access package assignment policy](accesspackageassignmentpolicy.md). Specifies the primary, fallback, and escalation approvers of each stage.\n\nIn PIM, defines the settings of the approval stages in a [unifiedRoleManagementPolicyApprovalRule](unifiedrolemanagementpolicyapprovalrule.md) object. Specifies the primary and escalation approvers of each stage and whether approvals and escalations are required. / https://learn.microsoft.com/en-us/graph/api/resources/approvalstage?view=graph-rest-beta. The _provider_ default value is `[]`.",
										},
										"is_approval_required": schema.BoolAttribute{
											Optional:            true,
											MarkdownDescription: "Indicates whether approval is required for requests in this policy.",
										},
										"is_approval_required_for_extension": schema.BoolAttribute{
											Optional:            true,
											MarkdownDescription: "Indicates whether approval is required for a user to extend their assignment.",
										},
										"is_requestor_justification_required": schema.BoolAttribute{
											Optional:            true,
											MarkdownDescription: "Indicates whether the requestor is required to supply a justification in their request.",
										},
									},
									MarkdownDescription: "The settings for approval of the role assignment. / The settings for approval as defined in a role management policy rule. / https://learn.microsoft.com/en-us/graph/api/resources/approvalsettings?view=graph-rest-beta",
								},
							},
							Validators: []validator.Object{
								unifiedRoleManagementPolicyUnifiedRoleManagementPolicyRuleValidator,
							},
							MarkdownDescription: "A type derived from the [unifiedRoleManagementPolicyRule](../resources/unifiedrolemanagementpolicyrule.md) resource type that defines rules for approving a role assignment. / https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyapprovalrule?view=graph-rest-beta",
						},
					},
					"authentication_context": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.unifiedRoleManagementPolicyAuthenticationContextRule",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // unifiedRoleManagementPolicyAuthenticationContextRule
								"claim_value": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The value of the authentication context claim.",
								},
								"is_enabled": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: "Whether this rule is enabled.",
								},
							},
							Validators: []validator.Object{
								unifiedRoleManagementPolicyUnifiedRoleManagementPolicyRuleValidator,
							},
							MarkdownDescription: "A type derived from the [unifiedRoleManagementPolicyRule](../resources/unifiedrolemanagementpolicyrule.md) resource type that defines the authentication context rule for the conditional access policy associated with a role management policy. / https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyauthenticationcontextrule?view=graph-rest-beta",
						},
					},
					"enablement": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.unifiedRoleManagementPolicyEnablementRule",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // unifiedRoleManagementPolicyEnablementRule
								"enabled_rules": schema.SetAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: "The collection of rules that are enabled for this policy rule. For example, `MultiFactorAuthentication`, `Ticketing`, and `Justification`. The _provider_ default value is `[]`.",
								},
							},
							Validators: []validator.Object{
								unifiedRoleManagementPolicyUnifiedRoleManagementPolicyRuleValidator,
							},
							MarkdownDescription: "A type derived from the [unifiedRoleManagementPolicyRule](../resources/unifiedrolemanagementpolicyrule.md) resource type that defines the rules to enable the assignment, for example, enable MFA, justification on assignments or ticketing information. / https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyenablementrule?view=graph-rest-beta",
						},
					},
					"expiration": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.unifiedRoleManagementPolicyExpirationRule",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // unifiedRoleManagementPolicyExpirationRule
								"is_expiration_required": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: "Indicates whether expiration is required or if it's a permanently active assignment or eligibility.",
								},
								"maximum_duration": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The maximum duration allowed for eligibility or assignment that isn't permanent. Required when **isExpirationRequired** is `true`.",
								},
							},
							Validators: []validator.Object{
								unifiedRoleManagementPolicyUnifiedRoleManagementPolicyRuleValidator,
							},
							MarkdownDescription: "A type derived from the [unifiedRoleManagementPolicyRule](../resources/unifiedrolemanagementpolicyrule.md) resource type that defines the maximum duration a role can be assigned to a principal (either through direct assignment or through activation of eligibility / https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyexpirationrule?view=graph-rest-beta",
						},
					},
					"notification": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.unifiedRoleManagementPolicyNotificationRule",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // unifiedRoleManagementPolicyNotificationRule
								"is_default_recipients_enabled": schema.BoolAttribute{
									Optional:            true,
									MarkdownDescription: "Indicates whether a default recipient will receive the notification email.",
								},
								"notification_level": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The level of notification. The possible values are `None`, `Critical`, `All`.",
								},
								"notification_recipients": schema.SetAttribute{
									ElementType:         types.StringType,
									Optional:            true,
									PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: "The list of recipients of the email notifications. The _provider_ default value is `[]`.",
								},
								"notification_type": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The type of notification. Only `Email` is supported.",
								},
								"recipient_type": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The type of recipient of the notification. The possible values are `Requestor`, `Approver`, `Admin`.",
								},
							},
							Validators: []validator.Object{
								unifiedRoleManagementPolicyUnifiedRoleManagementPolicyRuleValidator,
							},
							MarkdownDescription: "A type derived from the [unifiedRoleManagementPolicyRule](../resources/unifiedrolemanagementpolicyrule.md) resource type that defines the email notification rules for role assignments, activations, and approvals. / https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicynotificationrule?view=graph-rest-beta",
						},
					},
				},
			},
			MarkdownDescription: "The collection of rules like approval rules and expiration rules. Supports `$expand`. / An abstract type that defines the rules associated with role management policies. This abstract type is inherited by the following resources that define the various types of rules and their settings associated with role management policies. / https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyrule?view=graph-rest-beta",
		},
	},
	MarkdownDescription: "Specifies the various policies associated with scopes and roles. For policies that apply to Azure RBAC, use the [Azure REST PIM API for role management policies](/rest/api/authorization/role-management-policies). / https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicy?view=graph-rest-beta\n\nProvider Note: MS Graph entities of type `unifiedRoleManagementPolicy` cannot be created (or deleted) but only existing ones\ncan be updated. Therefore this resource will always update an existing entity (determined by the `id` attribute)\nand automatically import it into TF state upon creation.  \nAlso, entities will never be deleted from MS Graph but only get removed from TF state instead. ||| MS Graph: Role management",
}

var unifiedRoleManagementPolicyUnifiedRoleManagementPolicyRuleValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("approval"),
	path.MatchRelative().AtParent().AtName("authentication_context"),
	path.MatchRelative().AtParent().AtName("enablement"),
	path.MatchRelative().AtParent().AtName("expiration"),
	path.MatchRelative().AtParent().AtName("notification"),
)

var unifiedRoleManagementPolicyUserSetValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("connected_organization_members"),
	path.MatchRelative().AtParent().AtName("external_sponsors"),
	path.MatchRelative().AtParent().AtName("group_members"),
	path.MatchRelative().AtParent().AtName("internal_sponsors"),
	path.MatchRelative().AtParent().AtName("requestor_manager"),
	path.MatchRelative().AtParent().AtName("single_user"),
	path.MatchRelative().AtParent().AtName("target_user_sponsors"),
)
