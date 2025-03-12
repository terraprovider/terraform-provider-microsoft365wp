package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	CrossTenantAccessPolicyConfigurationDefaultResource = generic.GenericResource{
		TypeNameSuffix: "cross_tenant_access_policy_configuration_default",
		SpecificSchema: crossTenantAccessPolicyConfigurationDefaultResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:     "/policies/crossTenantAccessPolicy/default",
			IsSingleton: true,
		},
	}

	CrossTenantAccessPolicyConfigurationDefaultSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&CrossTenantAccessPolicyConfigurationDefaultResource)
)

var crossTenantAccessPolicyConfigurationDefaultResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyConfigurationDefault
		"id": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"automatic_user_consent_settings": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // inboundOutboundPolicyConfiguration
				"inbound_allowed": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
					Computed:            true,
					MarkdownDescription: "Defines whether external users coming inbound are allowed. The _provider_ default value is `false`.",
				},
				"outbound_allowed": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
					Computed:            true,
					MarkdownDescription: "Defines whether internal users are allowed to go outbound. The _provider_ default value is `false`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Determines the default configuration for automatic user consent settings. The **inboundAllowed** and **outboundAllowed** properties are always `false` and can't be updated in the default configuration. Read-only. / Defines the inbound and outbound rulesets for particular configurations within cross-tenant access settings. / https://learn.microsoft.com/en-us/graph/api/resources/inboundoutboundpolicyconfiguration?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
		"b2b_collaboration_inbound": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyB2BSetting
				"applications": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTargetConfiguration
						"access_type": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("allowed", "blocked", "unknownFutureValue"),
							},
							PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("allowed")},
							Computed:            true,
							MarkdownDescription: "Defines whether access is allowed or blocked. The / Possible values are: `allowed`, `blocked`, `unknownFutureValue`. The _provider_ default value is `\"allowed\"`.",
						},
						"targets": schema.SetNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTarget
									"target": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "Defines the target for cross-tenant access policy settings and can have one of the following values: <li> The unique identifier of the user, group, or application <li> `AllUsers` <li> `AllApplications` - Refers to any [Microsoft cloud application](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#microsoft-cloud-applications). <li> `Office365` - Includes the applications mentioned as part of the [Office 365](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#office-365) suite.",
									},
									"target_type": schema.StringAttribute{
										Optional: true,
										Validators: []validator.String{
											stringvalidator.OneOf("user", "group", "application", "unknownFutureValue"),
										},
										PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("application")},
										Computed:            true,
										MarkdownDescription: "The type of resource that you want to target. The / Possible values are: `user`, `group`, `application`, `unknownFutureValue`. The _provider_ default value is `\"application\"`.",
									},
								},
							},
							PlanModifiers: []planmodifier.Set{
								wpdefaultvalue.SetDefaultValue([]any{map[string]any{
									"target": "AllApplications",
								}}),
							},
							Computed:            true,
							MarkdownDescription: "Specifies whether to target users, groups, or applications with this rule. / Defines how to target your cross-tenant access policy settings. Settings can be targeted to specific users, groups, or applications. You can also use keywords to target specific groups or applications. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytarget?view=graph-rest-beta. The _provider_ default value is `[{\"target\":\"AllApplications\"}]`.",
						},
					},
					PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "The list of applications targeted with your cross-tenant access policy. / Defines the target of a cross-tenant access policy setting configuration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytargetconfiguration?view=graph-rest-beta. The _provider_ default value is `{}`.",
				},
				"users_and_groups": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTargetConfiguration
						"access_type": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("allowed", "blocked", "unknownFutureValue"),
							},
							PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("allowed")},
							Computed:            true,
							MarkdownDescription: "Defines whether access is allowed or blocked. The / Possible values are: `allowed`, `blocked`, `unknownFutureValue`. The _provider_ default value is `\"allowed\"`.",
						},
						"targets": schema.SetNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTarget
									"target": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "Defines the target for cross-tenant access policy settings and can have one of the following values: <li> The unique identifier of the user, group, or application <li> `AllUsers` <li> `AllApplications` - Refers to any [Microsoft cloud application](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#microsoft-cloud-applications). <li> `Office365` - Includes the applications mentioned as part of the [Office 365](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#office-365) suite.",
									},
									"target_type": schema.StringAttribute{
										Optional: true,
										Validators: []validator.String{
											stringvalidator.OneOf("user", "group", "application", "unknownFutureValue"),
										},
										PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("user")},
										Computed:            true,
										MarkdownDescription: "The type of resource that you want to target. The / Possible values are: `user`, `group`, `application`, `unknownFutureValue`. The _provider_ default value is `\"user\"`.",
									},
								},
							},
							PlanModifiers: []planmodifier.Set{
								wpdefaultvalue.SetDefaultValue([]any{map[string]any{
									"target": "AllUsers",
								}}),
							},
							Computed:            true,
							MarkdownDescription: "Specifies whether to target users, groups, or applications with this rule. / Defines how to target your cross-tenant access policy settings. Settings can be targeted to specific users, groups, or applications. You can also use keywords to target specific groups or applications. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytarget?view=graph-rest-beta. The _provider_ default value is `[{\"target\":\"AllUsers\"}]`.",
						},
					},
					PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "The list of users and groups targeted with your cross-tenant access policy. / Defines the target of a cross-tenant access policy setting configuration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytargetconfiguration?view=graph-rest-beta. The _provider_ default value is `{}`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			Description:         `b2bCollaborationInbound`, // custom MS Graph attribute name
			MarkdownDescription: "Defines your default configuration for users from other organizations accessing your resources via Microsoft Entra B2B collaboration. / Defines the inbound and outbound rulesets for Microsoft Entra B2B collaboration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyb2bsetting?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
		"b2b_collaboration_outbound": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyB2BSetting
				"applications": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTargetConfiguration
						"access_type": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("allowed", "blocked", "unknownFutureValue"),
							},
							PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("allowed")},
							Computed:            true,
							MarkdownDescription: "Defines whether access is allowed or blocked. The / Possible values are: `allowed`, `blocked`, `unknownFutureValue`. The _provider_ default value is `\"allowed\"`.",
						},
						"targets": schema.SetNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTarget
									"target": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "Defines the target for cross-tenant access policy settings and can have one of the following values: <li> The unique identifier of the user, group, or application <li> `AllUsers` <li> `AllApplications` - Refers to any [Microsoft cloud application](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#microsoft-cloud-applications). <li> `Office365` - Includes the applications mentioned as part of the [Office 365](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#office-365) suite.",
									},
									"target_type": schema.StringAttribute{
										Optional: true,
										Validators: []validator.String{
											stringvalidator.OneOf("user", "group", "application", "unknownFutureValue"),
										},
										PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("application")},
										Computed:            true,
										MarkdownDescription: "The type of resource that you want to target. The / Possible values are: `user`, `group`, `application`, `unknownFutureValue`. The _provider_ default value is `\"application\"`.",
									},
								},
							},
							PlanModifiers: []planmodifier.Set{
								wpdefaultvalue.SetDefaultValue([]any{map[string]any{
									"target": "AllApplications",
								}}),
							},
							Computed:            true,
							MarkdownDescription: "Specifies whether to target users, groups, or applications with this rule. / Defines how to target your cross-tenant access policy settings. Settings can be targeted to specific users, groups, or applications. You can also use keywords to target specific groups or applications. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytarget?view=graph-rest-beta. The _provider_ default value is `[{\"target\":\"AllApplications\"}]`.",
						},
					},
					PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "The list of applications targeted with your cross-tenant access policy. / Defines the target of a cross-tenant access policy setting configuration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytargetconfiguration?view=graph-rest-beta. The _provider_ default value is `{}`.",
				},
				"users_and_groups": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTargetConfiguration
						"access_type": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("allowed", "blocked", "unknownFutureValue"),
							},
							PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("allowed")},
							Computed:            true,
							MarkdownDescription: "Defines whether access is allowed or blocked. The / Possible values are: `allowed`, `blocked`, `unknownFutureValue`. The _provider_ default value is `\"allowed\"`.",
						},
						"targets": schema.SetNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTarget
									"target": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "Defines the target for cross-tenant access policy settings and can have one of the following values: <li> The unique identifier of the user, group, or application <li> `AllUsers` <li> `AllApplications` - Refers to any [Microsoft cloud application](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#microsoft-cloud-applications). <li> `Office365` - Includes the applications mentioned as part of the [Office 365](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#office-365) suite.",
									},
									"target_type": schema.StringAttribute{
										Optional: true,
										Validators: []validator.String{
											stringvalidator.OneOf("user", "group", "application", "unknownFutureValue"),
										},
										PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("user")},
										Computed:            true,
										MarkdownDescription: "The type of resource that you want to target. The / Possible values are: `user`, `group`, `application`, `unknownFutureValue`. The _provider_ default value is `\"user\"`.",
									},
								},
							},
							PlanModifiers: []planmodifier.Set{
								wpdefaultvalue.SetDefaultValue([]any{map[string]any{
									"target": "AllUsers",
								}}),
							},
							Computed:            true,
							MarkdownDescription: "Specifies whether to target users, groups, or applications with this rule. / Defines how to target your cross-tenant access policy settings. Settings can be targeted to specific users, groups, or applications. You can also use keywords to target specific groups or applications. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytarget?view=graph-rest-beta. The _provider_ default value is `[{\"target\":\"AllUsers\"}]`.",
						},
					},
					PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "The list of users and groups targeted with your cross-tenant access policy. / Defines the target of a cross-tenant access policy setting configuration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytargetconfiguration?view=graph-rest-beta. The _provider_ default value is `{}`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			Description:         `b2bCollaborationOutbound`, // custom MS Graph attribute name
			MarkdownDescription: "Defines your default configuration for users in your organization going outbound to access resources in another organization via Microsoft Entra B2B collaboration. / Defines the inbound and outbound rulesets for Microsoft Entra B2B collaboration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyb2bsetting?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
		"b2b_direct_connect_inbound": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyB2BSetting
				"applications": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTargetConfiguration
						"access_type": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("allowed", "blocked", "unknownFutureValue"),
							},
							PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("blocked")},
							Computed:            true,
							MarkdownDescription: "Defines whether access is allowed or blocked. The / Possible values are: `allowed`, `blocked`, `unknownFutureValue`. The _provider_ default value is `\"blocked\"`.",
						},
						"targets": schema.SetNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTarget
									"target": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "Defines the target for cross-tenant access policy settings and can have one of the following values: <li> The unique identifier of the user, group, or application <li> `AllUsers` <li> `AllApplications` - Refers to any [Microsoft cloud application](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#microsoft-cloud-applications). <li> `Office365` - Includes the applications mentioned as part of the [Office 365](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#office-365) suite.",
									},
									"target_type": schema.StringAttribute{
										Optional: true,
										Validators: []validator.String{
											stringvalidator.OneOf("user", "group", "application", "unknownFutureValue"),
										},
										PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("application")},
										Computed:            true,
										MarkdownDescription: "The type of resource that you want to target. The / Possible values are: `user`, `group`, `application`, `unknownFutureValue`. The _provider_ default value is `\"application\"`.",
									},
								},
							},
							PlanModifiers: []planmodifier.Set{
								wpdefaultvalue.SetDefaultValue([]any{map[string]any{
									"target": "AllApplications",
								}}),
							},
							Computed:            true,
							MarkdownDescription: "Specifies whether to target users, groups, or applications with this rule. / Defines how to target your cross-tenant access policy settings. Settings can be targeted to specific users, groups, or applications. You can also use keywords to target specific groups or applications. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytarget?view=graph-rest-beta. The _provider_ default value is `[{\"target\":\"AllApplications\"}]`.",
						},
					},
					PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "The list of applications targeted with your cross-tenant access policy. / Defines the target of a cross-tenant access policy setting configuration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytargetconfiguration?view=graph-rest-beta. The _provider_ default value is `{}`.",
				},
				"users_and_groups": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTargetConfiguration
						"access_type": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("allowed", "blocked", "unknownFutureValue"),
							},
							PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("blocked")},
							Computed:            true,
							MarkdownDescription: "Defines whether access is allowed or blocked. The / Possible values are: `allowed`, `blocked`, `unknownFutureValue`. The _provider_ default value is `\"blocked\"`.",
						},
						"targets": schema.SetNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTarget
									"target": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "Defines the target for cross-tenant access policy settings and can have one of the following values: <li> The unique identifier of the user, group, or application <li> `AllUsers` <li> `AllApplications` - Refers to any [Microsoft cloud application](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#microsoft-cloud-applications). <li> `Office365` - Includes the applications mentioned as part of the [Office 365](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#office-365) suite.",
									},
									"target_type": schema.StringAttribute{
										Optional: true,
										Validators: []validator.String{
											stringvalidator.OneOf("user", "group", "application", "unknownFutureValue"),
										},
										PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("user")},
										Computed:            true,
										MarkdownDescription: "The type of resource that you want to target. The / Possible values are: `user`, `group`, `application`, `unknownFutureValue`. The _provider_ default value is `\"user\"`.",
									},
								},
							},
							PlanModifiers: []planmodifier.Set{
								wpdefaultvalue.SetDefaultValue([]any{map[string]any{
									"target": "AllUsers",
								}}),
							},
							Computed:            true,
							MarkdownDescription: "Specifies whether to target users, groups, or applications with this rule. / Defines how to target your cross-tenant access policy settings. Settings can be targeted to specific users, groups, or applications. You can also use keywords to target specific groups or applications. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytarget?view=graph-rest-beta. The _provider_ default value is `[{\"target\":\"AllUsers\"}]`.",
						},
					},
					PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "The list of users and groups targeted with your cross-tenant access policy. / Defines the target of a cross-tenant access policy setting configuration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytargetconfiguration?view=graph-rest-beta. The _provider_ default value is `{}`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			Description:         `b2bDirectConnectInbound`, // custom MS Graph attribute name
			MarkdownDescription: "Defines your default configuration for users from other organizations accessing your resources via Microsoft Entra B2B direct connect. / Defines the inbound and outbound rulesets for Microsoft Entra B2B collaboration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyb2bsetting?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
		"b2b_direct_connect_outbound": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyB2BSetting
				"applications": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTargetConfiguration
						"access_type": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("allowed", "blocked", "unknownFutureValue"),
							},
							PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("blocked")},
							Computed:            true,
							MarkdownDescription: "Defines whether access is allowed or blocked. The / Possible values are: `allowed`, `blocked`, `unknownFutureValue`. The _provider_ default value is `\"blocked\"`.",
						},
						"targets": schema.SetNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTarget
									"target": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "Defines the target for cross-tenant access policy settings and can have one of the following values: <li> The unique identifier of the user, group, or application <li> `AllUsers` <li> `AllApplications` - Refers to any [Microsoft cloud application](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#microsoft-cloud-applications). <li> `Office365` - Includes the applications mentioned as part of the [Office 365](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#office-365) suite.",
									},
									"target_type": schema.StringAttribute{
										Optional: true,
										Validators: []validator.String{
											stringvalidator.OneOf("user", "group", "application", "unknownFutureValue"),
										},
										PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("application")},
										Computed:            true,
										MarkdownDescription: "The type of resource that you want to target. The / Possible values are: `user`, `group`, `application`, `unknownFutureValue`. The _provider_ default value is `\"application\"`.",
									},
								},
							},
							PlanModifiers: []planmodifier.Set{
								wpdefaultvalue.SetDefaultValue([]any{map[string]any{
									"target": "AllApplications",
								}}),
							},
							Computed:            true,
							MarkdownDescription: "Specifies whether to target users, groups, or applications with this rule. / Defines how to target your cross-tenant access policy settings. Settings can be targeted to specific users, groups, or applications. You can also use keywords to target specific groups or applications. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytarget?view=graph-rest-beta. The _provider_ default value is `[{\"target\":\"AllApplications\"}]`.",
						},
					},
					PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "The list of applications targeted with your cross-tenant access policy. / Defines the target of a cross-tenant access policy setting configuration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytargetconfiguration?view=graph-rest-beta. The _provider_ default value is `{}`.",
				},
				"users_and_groups": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTargetConfiguration
						"access_type": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("allowed", "blocked", "unknownFutureValue"),
							},
							PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("blocked")},
							Computed:            true,
							MarkdownDescription: "Defines whether access is allowed or blocked. The / Possible values are: `allowed`, `blocked`, `unknownFutureValue`. The _provider_ default value is `\"blocked\"`.",
						},
						"targets": schema.SetNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTarget
									"target": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "Defines the target for cross-tenant access policy settings and can have one of the following values: <li> The unique identifier of the user, group, or application <li> `AllUsers` <li> `AllApplications` - Refers to any [Microsoft cloud application](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#microsoft-cloud-applications). <li> `Office365` - Includes the applications mentioned as part of the [Office 365](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#office-365) suite.",
									},
									"target_type": schema.StringAttribute{
										Optional: true,
										Validators: []validator.String{
											stringvalidator.OneOf("user", "group", "application", "unknownFutureValue"),
										},
										PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("user")},
										Computed:            true,
										MarkdownDescription: "The type of resource that you want to target. The / Possible values are: `user`, `group`, `application`, `unknownFutureValue`. The _provider_ default value is `\"user\"`.",
									},
								},
							},
							PlanModifiers: []planmodifier.Set{
								wpdefaultvalue.SetDefaultValue([]any{map[string]any{
									"target": "AllUsers",
								}}),
							},
							Computed:            true,
							MarkdownDescription: "Specifies whether to target users, groups, or applications with this rule. / Defines how to target your cross-tenant access policy settings. Settings can be targeted to specific users, groups, or applications. You can also use keywords to target specific groups or applications. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytarget?view=graph-rest-beta. The _provider_ default value is `[{\"target\":\"AllUsers\"}]`.",
						},
					},
					PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "The list of users and groups targeted with your cross-tenant access policy. / Defines the target of a cross-tenant access policy setting configuration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytargetconfiguration?view=graph-rest-beta. The _provider_ default value is `{}`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			Description:         `b2bDirectConnectOutbound`, // custom MS Graph attribute name
			MarkdownDescription: "Defines your default configuration for users in your organization going outbound to access resources in another organization via Microsoft Entra B2B direct connect. / Defines the inbound and outbound rulesets for Microsoft Entra B2B collaboration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyb2bsetting?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
		"inbound_trust": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyInboundTrust
				"is_compliant_device_accepted": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
					Computed:            true,
					MarkdownDescription: "Specifies whether compliant devices from external Microsoft Entra organizations are trusted. The _provider_ default value is `false`.",
				},
				"is_hybrid_azure_ad_joined_device_accepted": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
					Computed:            true,
					Description:         `isHybridAzureADJoinedDeviceAccepted`, // custom MS Graph attribute name
					MarkdownDescription: "Specifies whether Microsoft Entra hybrid joined devices from external Microsoft Entra organizations are trusted. The _provider_ default value is `false`.",
				},
				"is_mfa_accepted": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
					Computed:            true,
					MarkdownDescription: "Specifies whether MFA from external Microsoft Entra organizations is trusted. The _provider_ default value is `false`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Determines the default configuration for trusting other Conditional Access claims from external Microsoft Entra organizations. / Defines the Conditional Access claims you want to accept from other Microsoft Entra organizations via your cross-tenant access policy configuration. These can be configured in your default configuration, partner-specific configuration, or both. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyinboundtrust?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
		"invitation_redemption_identity_provider_configuration": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // defaultInvitationRedemptionIdentityProviderConfiguration
				"fallback_identity_provider": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf("azureActiveDirectory", "externalFederation", "socialIdentityProviders", "emailOneTimePasscode", "microsoftAccount", "defaultConfiguredIdp", "unknownFutureValue"),
					},
					PlanModifiers: []planmodifier.String{
						wpdefaultvalue.StringDefaultValue("defaultConfiguredIdp"),
					},
					Computed:            true,
					MarkdownDescription: "The fallback identity provider to be used in case no primary identity provider can be used for guest invitation redemption.or `microsoftAccount`. / Possible values are: `azureActiveDirectory`, `externalFederation`, `socialIdentityProviders`, `emailOneTimePasscode`, `microsoftAccount`, `defaultConfiguredIdp`, `unknownFutureValue`. The _provider_ default value is `\"defaultConfiguredIdp\"`.",
				},
				"primary_identity_provider_precedence_order": schema.SetAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Validators: []validator.Set{
						setvalidator.ValueStringsAre(
							stringvalidator.OneOf("azureActiveDirectory", "externalFederation", "socialIdentityProviders", "emailOneTimePasscode", "microsoftAccount", "defaultConfiguredIdp", "unknownFutureValue"),
						),
					},
					PlanModifiers: []planmodifier.Set{
						wpdefaultvalue.SetDefaultValue([]any{"azureActiveDirectory", "externalFederation", "socialIdentityProviders"}),
					},
					Computed:            true,
					MarkdownDescription: "Collection of identity providers in priority order of preference to be used for guest invitation redemption.or `socialIdentityProviders`. / Possible values are: `azureActiveDirectory`, `externalFederation`, `socialIdentityProviders`, `emailOneTimePasscode`, `microsoftAccount`, `defaultConfiguredIdp`, `unknownFutureValue`. The _provider_ default value is `[\"azureActiveDirectory\",\"externalFederation\",\"socialIdentityProviders\"]`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Defines the priority order based on which an identity provider is selected during invitation redemption for a guest user. / Defines the invitation redemption provider configuration to set redemption flow settings for Microsoft Entra ID B2B collaboration. / https://learn.microsoft.com/en-us/graph/api/resources/defaultinvitationredemptionidentityproviderconfiguration?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
		"is_service_default": schema.BoolAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
			MarkdownDescription: "If `true`, the default configuration is set to the system default configuration. If `false`, the default settings are customized.",
		},
		"tenant_restrictions": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTenantRestrictions
				"applications": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTargetConfiguration
						"access_type": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("allowed", "blocked", "unknownFutureValue"),
							},
							PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("blocked")},
							Computed:            true,
							MarkdownDescription: "Defines whether access is allowed or blocked. The / Possible values are: `allowed`, `blocked`, `unknownFutureValue`. The _provider_ default value is `\"blocked\"`.",
						},
						"targets": schema.SetNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTarget
									"target": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "Defines the target for cross-tenant access policy settings and can have one of the following values: <li> The unique identifier of the user, group, or application <li> `AllUsers` <li> `AllApplications` - Refers to any [Microsoft cloud application](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#microsoft-cloud-applications). <li> `Office365` - Includes the applications mentioned as part of the [Office 365](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#office-365) suite.",
									},
									"target_type": schema.StringAttribute{
										Optional: true,
										Validators: []validator.String{
											stringvalidator.OneOf("user", "group", "application", "unknownFutureValue"),
										},
										PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("application")},
										Computed:            true,
										MarkdownDescription: "The type of resource that you want to target. The / Possible values are: `user`, `group`, `application`, `unknownFutureValue`. The _provider_ default value is `\"application\"`.",
									},
								},
							},
							PlanModifiers: []planmodifier.Set{
								wpdefaultvalue.SetDefaultValue([]any{map[string]any{
									"target": "AllApplications",
								}}),
							},
							Computed:            true,
							MarkdownDescription: "Specifies whether to target users, groups, or applications with this rule. / Defines how to target your cross-tenant access policy settings. Settings can be targeted to specific users, groups, or applications. You can also use keywords to target specific groups or applications. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytarget?view=graph-rest-beta. The _provider_ default value is `[{\"target\":\"AllApplications\"}]`.",
						},
					},
					PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "The list of applications targeted with your cross-tenant access policy. / Defines the target of a cross-tenant access policy setting configuration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytargetconfiguration?view=graph-rest-beta. The _provider_ default value is `{}`.",
				},
				"users_and_groups": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTargetConfiguration
						"access_type": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("allowed", "blocked", "unknownFutureValue"),
							},
							PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("blocked")},
							Computed:            true,
							MarkdownDescription: "Defines whether access is allowed or blocked. The / Possible values are: `allowed`, `blocked`, `unknownFutureValue`. The _provider_ default value is `\"blocked\"`.",
						},
						"targets": schema.SetNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTarget
									"target": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "Defines the target for cross-tenant access policy settings and can have one of the following values: <li> The unique identifier of the user, group, or application <li> `AllUsers` <li> `AllApplications` - Refers to any [Microsoft cloud application](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#microsoft-cloud-applications). <li> `Office365` - Includes the applications mentioned as part of the [Office 365](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#office-365) suite.",
									},
									"target_type": schema.StringAttribute{
										Optional: true,
										Validators: []validator.String{
											stringvalidator.OneOf("user", "group", "application", "unknownFutureValue"),
										},
										PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("user")},
										Computed:            true,
										MarkdownDescription: "The type of resource that you want to target. The / Possible values are: `user`, `group`, `application`, `unknownFutureValue`. The _provider_ default value is `\"user\"`.",
									},
								},
							},
							PlanModifiers: []planmodifier.Set{
								wpdefaultvalue.SetDefaultValue([]any{map[string]any{
									"target": "AllUsers",
								}}),
							},
							Computed:            true,
							MarkdownDescription: "Specifies whether to target users, groups, or applications with this rule. / Defines how to target your cross-tenant access policy settings. Settings can be targeted to specific users, groups, or applications. You can also use keywords to target specific groups or applications. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytarget?view=graph-rest-beta. The _provider_ default value is `[{\"target\":\"AllUsers\"}]`.",
						},
					},
					PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "The list of users and groups targeted with your cross-tenant access policy. / Defines the target of a cross-tenant access policy setting configuration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytargetconfiguration?view=graph-rest-beta. The _provider_ default value is `{}`.",
				},
				"devices": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // devicesFilter
						"mode": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("allowed", "blocked", "unknownFutureValue"),
							},
							MarkdownDescription: "Determines whether devices that satisfy the rule should be allowed or blocked. The / Possible values are: `allowed`, `blocked`, `unknownFutureValue`",
						},
						"rule": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "Defines the rule to filter the devices. For example, `device.deviceAttribute2 -eq 'PrivilegedAccessWorkstation'`.",
						},
					},
					MarkdownDescription: "Defines the rule for filtering devices and whether devices satisfying the rule should be allowed or blocked. This property isn't supported on the server side yet. / Defines a rule to filter the devices and whether devices that satisfy the rule should be allowed or blocked. / https://learn.microsoft.com/en-us/graph/api/resources/devicesfilter?view=graph-rest-beta",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Defines the default tenant restrictions configuration for users in your organization who access an external organization on your network or devices. / Defines how to target your tenant restrictions settings. Tenant restrictions give you control over the external organizations that your users can access from your network or devices when they use external identities. Settings can be targeted to specific users, groups, or applications. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytenantrestrictions?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
	},
	MarkdownDescription: "Represents the default configuration for cross-tenant access and tenant restrictions. Cross-tenant access settings include inbound and outbound settings of Microsoft Entra B2B collaboration and B2B direct connect. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyconfigurationdefault?view=graph-rest-beta ||| MS Graph: Cross-tenant access",
}
