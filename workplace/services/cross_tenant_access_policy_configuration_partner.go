package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	CrossTenantAccessPolicyConfigurationPartnerResource = generic.GenericResource{
		TypeNameSuffix: "cross_tenant_access_policy_configuration_partner",
		SpecificSchema: crossTenantAccessPolicyConfigurationPartnerResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/policies/crossTenantAccessPolicy/partners",
			EntityId: generic.EntityIdOptions{
				AttrNameGraph: "tenantId",
			},
			ReadOptions: generic.ReadOptions{
				ODataExpand: "identitySynchronization",
				DataSource: generic.DataSourceOptions{
					Plural: generic.PluralOptions{
						ExtraAttributes: []string{"is_in_multi_tenant_organization", "is_service_provider"},
					},
				},
			},
		},
	}

	CrossTenantAccessPolicyConfigurationPartnerSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&CrossTenantAccessPolicyConfigurationPartnerResource)

	CrossTenantAccessPolicyConfigurationPartnerPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&CrossTenantAccessPolicyConfigurationPartnerResource, "")
)

var crossTenantAccessPolicyConfigurationPartnerResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyConfigurationPartner
		"automatic_user_consent_settings": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // inboundOutboundPolicyConfiguration
				"inbound_allowed": schema.BoolAttribute{
					Optional:            true,
					MarkdownDescription: "Defines whether external users coming inbound are allowed.",
				},
				"outbound_allowed": schema.BoolAttribute{
					Optional:            true,
					MarkdownDescription: "Defines whether internal users are allowed to go outbound.",
				},
			},
			MarkdownDescription: "Determines the partner-specific configuration for automatic user consent settings. Unless configured, the **inboundAllowed** and **outboundAllowed** properties are `null` and inherit from the default settings, which is always `false`. / Defines the inbound and outbound rulesets for particular configurations within cross-tenant access settings. / https://learn.microsoft.com/en-us/graph/api/resources/inboundoutboundpolicyconfiguration?view=graph-rest-beta",
		},
		"b2b_collaboration_inbound": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          crossTenantAccessPolicyConfigurationPartnerCrossTenantAccessPolicyB2BSettingAttributes,
			Description:         `b2bCollaborationInbound`, // custom MS Graph attribute name
			MarkdownDescription: "Defines your partner-specific configuration for users from other organizations accessing your resources via Microsoft Entra B2B collaboration. / Defines the inbound and outbound rulesets for Microsoft Entra B2B collaboration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyb2bsetting?view=graph-rest-beta",
		},
		"b2b_collaboration_outbound": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          crossTenantAccessPolicyConfigurationPartnerCrossTenantAccessPolicyB2BSettingAttributes,
			Description:         `b2bCollaborationOutbound`, // custom MS Graph attribute name
			MarkdownDescription: "Defines your partner-specific configuration for users in your organization going outbound to access resources in another organization via Microsoft Entra B2B collaboration. / Defines the inbound and outbound rulesets for Microsoft Entra B2B collaboration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyb2bsetting?view=graph-rest-beta",
		},
		"b2b_direct_connect_inbound": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          crossTenantAccessPolicyConfigurationPartnerCrossTenantAccessPolicyB2BSettingAttributes,
			Description:         `b2bDirectConnectInbound`, // custom MS Graph attribute name
			MarkdownDescription: "Defines your partner-specific configuration for users from other organizations accessing your resources via Azure B2B direct connect. / Defines the inbound and outbound rulesets for Microsoft Entra B2B collaboration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyb2bsetting?view=graph-rest-beta",
		},
		"b2b_direct_connect_outbound": schema.SingleNestedAttribute{
			Optional:            true,
			Attributes:          crossTenantAccessPolicyConfigurationPartnerCrossTenantAccessPolicyB2BSettingAttributes,
			Description:         `b2bDirectConnectOutbound`, // custom MS Graph attribute name
			MarkdownDescription: "Defines your partner-specific configuration for users in your organization going outbound to access resources in another organization via Microsoft Entra B2B direct connect. / Defines the inbound and outbound rulesets for Microsoft Entra B2B collaboration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyb2bsetting?view=graph-rest-beta",
		},
		"inbound_trust": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyInboundTrust
				"is_compliant_device_accepted": schema.BoolAttribute{
					Optional:            true,
					MarkdownDescription: "Specifies whether compliant devices from external Microsoft Entra organizations are trusted.",
				},
				"is_hybrid_azure_ad_joined_device_accepted": schema.BoolAttribute{
					Optional:            true,
					Description:         `isHybridAzureADJoinedDeviceAccepted`, // custom MS Graph attribute name
					MarkdownDescription: "Specifies whether Microsoft Entra hybrid joined devices from external Microsoft Entra organizations are trusted.",
				},
				"is_mfa_accepted": schema.BoolAttribute{
					Optional:            true,
					MarkdownDescription: "Specifies whether MFA from external Microsoft Entra organizations is trusted.",
				},
			},
			MarkdownDescription: "Determines the partner-specific configuration for trusting other Conditional Access claims from external Microsoft Entra organizations. / Defines the Conditional Access claims you want to accept from other Microsoft Entra organizations via your cross-tenant access policy configuration. These can be configured in your default configuration, partner-specific configuration, or both. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyinboundtrust?view=graph-rest-beta",
		},
		"is_in_multi_tenant_organization": schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Identifies whether a tenant is a member of a multitenant organization.",
		},
		"is_service_provider": schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Identifies whether the partner-specific configuration is a Cloud Service Provider for your organization.",
		},
		"tenant_id": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "The tenant identifier for the partner Microsoft Entra organization. Read-only. Key. The _provider_ default value is `\"\"`.",
		},
		"tenant_restrictions": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicyTenantRestrictions
				"applications": schema.SingleNestedAttribute{
					Optional:            true,
					Attributes:          crossTenantAccessPolicyConfigurationPartnerCrossTenantAccessPolicyTargetConfigurationAttributes,
					MarkdownDescription: "The list of applications targeted with your cross-tenant access policy. / Defines the target of a cross-tenant access policy setting configuration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytargetconfiguration?view=graph-rest-beta",
				},
				"users_and_groups": schema.SingleNestedAttribute{
					Optional:            true,
					Attributes:          crossTenantAccessPolicyConfigurationPartnerCrossTenantAccessPolicyTargetConfigurationAttributes,
					MarkdownDescription: "The list of users and groups targeted with your cross-tenant access policy. / Defines the target of a cross-tenant access policy setting configuration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytargetconfiguration?view=graph-rest-beta",
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
			MarkdownDescription: "Defines the partner-specific tenant restrictions configuration for users in your organization who access a partner organization using partner supplied identities on your network or devices. / Defines how to target your tenant restrictions settings. Tenant restrictions give you control over the external organizations that your users can access from your network or devices when they use external identities. Settings can be targeted to specific users, groups, or applications. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytenantrestrictions?view=graph-rest-beta",
		},
		"identity_synchronization": schema.SingleNestedAttribute{
			Computed: true,
			Attributes: map[string]schema.Attribute{ // crossTenantIdentitySyncPolicyPartner
				"display_name": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Display name for the cross-tenant user synchronization policy. Use the name of the partner Microsoft Entra tenant to easily identify the policy. Optional.",
				},
				"tenant_id": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "Tenant identifier for the partner Microsoft Entra organization. Read-only.",
				},
				"user_sync_inbound": schema.SingleNestedAttribute{
					Computed: true,
					Attributes: map[string]schema.Attribute{ // crossTenantUserSyncInbound
						"is_sync_allowed": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Defines whether user objects should be synchronized from the partner tenant. `false` causes any current user synchronization from the source tenant to the target tenant to stop. This property has no impact on existing users who have already been synchronized.",
						},
					},
					MarkdownDescription: "Defines whether users can be synchronized from the partner tenant. Key. / Defines whether users can be synchronized from the partner tenant. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantusersyncinbound?view=graph-rest-beta",
				},
			},
			MarkdownDescription: "Defines the cross-tenant policy for the synchronization of users from a partner tenant. Use this user synchronization policy to streamline collaboration between users in a multitenant organization by automating the creation, update, and deletion of users from one tenant to another. / Defines the cross-tenant policy for synchronization of users from a partner tenant. Use this user synchronization policy to streamline collaboration between users in a multi-tenant organization by automating the creation, update, and deletion of users from one tenant to another. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantidentitysyncpolicypartner?view=graph-rest-beta",
		},
	},
	MarkdownDescription: "Represents the partner-specific configuration for cross-tenant access and tenant restrictions. Cross-tenant access settings include inbound and outbound settings of Microsoft Entra B2B collaboration and B2B direct connect.\n\nFor any partner-specific property that is `null`, these settings inherit the behavior configured in your [default cross-tenant access settings](../resources/crosstenantaccesspolicyconfigurationdefault.md). / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicyconfigurationpartner?view=graph-rest-beta ||| MS Graph: Cross-tenant access",
}

var crossTenantAccessPolicyConfigurationPartnerCrossTenantAccessPolicyB2BSettingAttributes = map[string]schema.Attribute{ // crossTenantAccessPolicyB2BSetting
	"applications": schema.SingleNestedAttribute{
		Optional:            true,
		Attributes:          crossTenantAccessPolicyConfigurationPartnerCrossTenantAccessPolicyTargetConfigurationAttributes,
		MarkdownDescription: "The list of applications targeted with your cross-tenant access policy. / Defines the target of a cross-tenant access policy setting configuration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytargetconfiguration?view=graph-rest-beta",
	},
	"users_and_groups": schema.SingleNestedAttribute{
		Optional:            true,
		Attributes:          crossTenantAccessPolicyConfigurationPartnerCrossTenantAccessPolicyTargetConfigurationAttributes,
		MarkdownDescription: "The list of users and groups targeted with your cross-tenant access policy. / Defines the target of a cross-tenant access policy setting configuration. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytargetconfiguration?view=graph-rest-beta",
	},
}

var crossTenantAccessPolicyConfigurationPartnerCrossTenantAccessPolicyTargetConfigurationAttributes = map[string]schema.Attribute{ // crossTenantAccessPolicyTargetConfiguration
	"access_type": schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			stringvalidator.OneOf("allowed", "blocked", "unknownFutureValue"),
		},
		MarkdownDescription: "Defines whether access is allowed or blocked. The / Possible values are: `allowed`, `blocked`, `unknownFutureValue`",
	},
	"targets": schema.SetNestedAttribute{
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: crossTenantAccessPolicyConfigurationPartnerCrossTenantAccessPolicyTargetAttributes,
		},
		PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
		Computed:            true,
		MarkdownDescription: "Specifies whether to target users, groups, or applications with this rule. / Defines how to target your cross-tenant access policy settings. Settings can be targeted to specific users, groups, or applications. You can also use keywords to target specific groups or applications. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicytarget?view=graph-rest-beta. The _provider_ default value is `[]`.",
	},
}

var crossTenantAccessPolicyConfigurationPartnerCrossTenantAccessPolicyTargetAttributes = map[string]schema.Attribute{ // crossTenantAccessPolicyTarget
	"target": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "Defines the target for cross-tenant access policy settings and can have one of the following values: <li> The unique identifier of the user, group, or application <li> `AllUsers` <li> `AllApplications` - Refers to any [Microsoft cloud application](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#microsoft-cloud-applications). <li> `Office365` - Includes the applications mentioned as part of the [Office 365](/azure/active-directory/conditional-access/concept-conditional-access-cloud-apps#office-365) suite.",
	},
	"target_type": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("user", "group", "application", "unknownFutureValue"),
		},
		MarkdownDescription: "The type of resource that you want to target. The / Possible values are: `user`, `group`, `application`, `unknownFutureValue`",
	},
}
