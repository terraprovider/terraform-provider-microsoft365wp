package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	ConditionalAccessPolicyResource = generic.GenericResource{
		TypeNameSuffix:           "conditional_access_policy",
		SpecificSchema:           conditionalAccessPolicyResourceSchema,
		SpecificConfigValidators: conditionalAccessPolicyResourceSchemaValidators,
		AccessParams: generic.AccessParams{
			BaseUri: "/identity/conditionalAccess/policies",
		},
	}

	ConditionalAccessPolicySingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&ConditionalAccessPolicyResource)

	ConditionalAccessPolicyPluralDataSource = generic.CreateGenericDataSourcePluralFromSingular(
		&ConditionalAccessPolicySingularDataSource, "")
)

var conditionalAccessPolicyResourceSchemaValidators = []resource.ConfigValidator{
	resourcevalidator.AtLeastOneOf(
		path.MatchRoot("grant_controls"),
		path.MatchRoot("session_controls"),
	),
}

var conditionalAccessPolicyResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // conditionalAccessPolicy
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Specifies the identifier of a conditionalAccessPolicy object. Read-only.",
		},
		"conditions": schema.SingleNestedAttribute{
			Required: true,
			Attributes: map[string]schema.Attribute{ // conditionalAccessConditionSet
				"applications": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // conditionalAccessApplications
						"application_filter": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // conditionalAccessFilter
								"mode": schema.StringAttribute{
									Required:            true,
									Validators:          []validator.String{stringvalidator.OneOf("include", "exclude")},
									MarkdownDescription: "Mode to use for the filter. Possible values are `include` or `exclude`. / Possible values are: `include`, `exclude`",
								},
								"rule": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "Rule syntax is similar to that used for membership rules for groups in Microsoft Entra ID. For details, see [rules with multiple expressions](/azure/active-directory/enterprise-users/groups-dynamic-membership#rules-with-multiple-expressions)",
								},
							},
							MarkdownDescription: "Filter that defines the dynamic-application-syntax rule to include/exclude cloud applications. A filter can use custom security attributes to include/exclude applications. / Represents filter in the policy scope. / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessfilter?view=graph-rest-beta",
						},
						"exclude_applications": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "Can be one of the following: <li> The list of client IDs (**appId**) explicitly excluded from the policy.<li> `Office365` - For the list of apps included in `Office365`, see [Apps included in Conditional Access Office 365 app suite](/entra/identity/conditional-access/reference-office-365-application-contents) <li> `MicrosoftAdminPortals` - For more information, see [Conditional Access Target resources: Microsoft Admin Portals](/entra/identity/conditional-access/concept-conditional-access-cloud-apps#microsoft-admin-portals). The _provider_ default value is `[]`.",
						},
						"include_applications": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "Can be one of the following: <li> The list of client IDs (**appId**) the policy applies to, unless explicitly excluded (in **excludeApplications**) <li> `All` <li> `Office365` - For the list of apps included in `Office365`, see [Apps included in Conditional Access Office 365 app suite](/entra/identity/conditional-access/reference-office-365-application-contents) <li> `MicrosoftAdminPortals` - For more information, see [Conditional Access Target resources: Microsoft Admin Portals](/entra/identity/conditional-access/concept-conditional-access-cloud-apps#microsoft-admin-portals). The _provider_ default value is `[]`.",
						},
						"include_authentication_context_class_references": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "Authentication context class references include. Supported values are `c1` through `c25`. The _provider_ default value is `[]`.",
						},
						"include_user_actions": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "User actions to include. Supported values are `urn:user:registersecurityinfo` and `urn:user:registerdevice`. The _provider_ default value is `[]`.",
						},
					},
					MarkdownDescription: "Applications and user actions included in and excluded from the policy. Required. / Represents the applications and user actions included in and excluded from the conditional access policy. / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessapplications?view=graph-rest-beta",
				},
				"authentication_flows": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // conditionalAccessAuthenticationFlows
						"transfer_methods": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								wpvalidator.FlagEnumValues("none", "deviceCodeFlow", "authenticationTransfer", "unknownFutureValue"),
							},
							MarkdownDescription: "Represents the transfer methods in scope for the policy. The / Possible values are: `none`, `deviceCodeFlow`, `authenticationTransfer`, `unknownFutureValue`",
						},
					},
					MarkdownDescription: "Authentication flows included in the policy scope. For more information, see [Conditional Access: Authentication flows](/entra/identity/conditional-access/concept-authentication-flows). / Represents the authentication flows in scope for the policy. / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessauthenticationflows?view=graph-rest-beta",
				},
				"client_applications": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // conditionalAccessClientApplications
						"exclude_service_principals": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "Service principal IDs excluded from the policy scope. The _provider_ default value is `[]`.",
						},
						"include_service_principals": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "Service principal IDs included in the policy scope, or `ServicePrincipalsInMyTenant`. The _provider_ default value is `[]`.",
						},
						"service_principal_filter": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // conditionalAccessFilter
								"mode": schema.StringAttribute{
									Required:            true,
									Validators:          []validator.String{stringvalidator.OneOf("include", "exclude")},
									MarkdownDescription: "Mode to use for the filter. Possible values are `include` or `exclude`. / Possible values are: `include`, `exclude`",
								},
								"rule": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "Rule syntax is similar to that used for membership rules for groups in Microsoft Entra ID. For details, see [rules with multiple expressions](/azure/active-directory/enterprise-users/groups-dynamic-membership#rules-with-multiple-expressions)",
								},
							},
							MarkdownDescription: "Filter that defines the dynamic-servicePrincipal-syntax rule to include/exclude service principals. A filter can use custom security attributes to include/exclude service principals. / Represents filter in the policy scope. / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessfilter?view=graph-rest-beta",
						},
					},
					MarkdownDescription: "Client applications (service principals and workload identities) included in and excluded from the policy. Either **users** or **clientApplications** is required. / Represents client applications (service principals and workload identities) included in and excluded from the policy scope. / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessclientapplications?view=graph-rest-beta",
				},
				"client_app_types": schema.SetAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Validators: []validator.Set{
						setvalidator.ValueStringsAre(
							stringvalidator.OneOf("all", "browser", "mobileAppsAndDesktopClients", "exchangeActiveSync", "easSupported", "other", "unknownFutureValue"),
						),
					},
					PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"all"})},
					Computed:            true,
					MarkdownDescription: "Client application types included in the policy. Required. <br/><br/> The `easUnsupported` enumeration member is deprecated in favor of `exchangeActiveSync`, which includes EAS supported and unsupported platforms. / Possible values are: `all`, `browser`, `mobileAppsAndDesktopClients`, `exchangeActiveSync`, `easSupported`, `other`, `unknownFutureValue`. The _provider_ default value is `[\"all\"]`.",
				},
				"devices": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // conditionalAccessDevices
						"device_filter": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // conditionalAccessFilter
								"mode": schema.StringAttribute{
									Required:            true,
									Validators:          []validator.String{stringvalidator.OneOf("include", "exclude")},
									MarkdownDescription: "Mode to use for the filter. Possible values are `include` or `exclude`. / Possible values are: `include`, `exclude`",
								},
								"rule": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "Rule syntax is similar to that used for membership rules for groups in Microsoft Entra ID. For details, see [rules with multiple expressions](/azure/active-directory/enterprise-users/groups-dynamic-membership#rules-with-multiple-expressions)",
								},
							},
							MarkdownDescription: "Filter that defines the dynamic-device-syntax rule to include/exclude devices. A filter can use device properties (such as extension attributes) to include/exclude them. Cannot be set if **includeDevices** or **excludeDevices** is set. / Represents filter in the policy scope. / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessfilter?view=graph-rest-beta",
						},
						"exclude_devices": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "States excluded from the scope of the policy. Possible values: `Compliant`, `DomainJoined`. Cannot be set if **deviceFIlter** is set. The _provider_ default value is `[]`.",
						},
						"exclude_device_states": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: ". The _provider_ default value is `[]`.",
						},
						"include_devices": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "States in the scope of the policy. `All` is the only allowed value. Cannot be set if **deviceFilter** is set. The _provider_ default value is `[]`.",
						},
						"include_device_states": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: ". The _provider_ default value is `[]`.",
						},
					},
					MarkdownDescription: "Devices in the policy. / Represents devices in the policy scope. / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessdevices?view=graph-rest-beta",
				},
				"device_states": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // conditionalAccessDeviceStates
						"exclude_states": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "States excluded from the scope of the policy. Possible values: `Compliant`, `DomainJoined`. The _provider_ default value is `[]`.",
						},
						"include_states": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "States in the scope of the policy. `All` is the only allowed value. The _provider_ default value is `[]`.",
						},
					},
					MarkdownDescription: "Device states in the policy. To be deprecated and removed. Use the **devices** property instead. / Represents device states in the policy scope. / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessdevicestates?view=graph-rest-beta",
				},
				"insider_risk_levels": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						wpvalidator.FlagEnumValues("minor", "moderate", "elevated", "unknownFutureValue"),
					},
					MarkdownDescription: "Insider risk levels included in the policy. The / Possible values are: `minor`, `moderate`, `elevated`, `unknownFutureValue`",
				},
				"locations": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // conditionalAccessLocations
						"exclude_locations": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "Location IDs excluded from scope of policy. The _provider_ default value is `[]`.",
						},
						"include_locations": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "Location IDs in scope of policy unless explicitly excluded, `All`, or `AllTrusted`. The _provider_ default value is `[]`.",
						},
					},
					MarkdownDescription: "Locations included in and excluded from the policy. / Represents locations included in and excluded from the policy scope. / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccesslocations?view=graph-rest-beta",
				},
				"platforms": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // conditionalAccessPlatforms
						"exclude_platforms": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.OneOf("android", "iOS", "windows", "windowsPhone", "macOS", "all", "unknownFutureValue", "linux"),
								),
							},
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "Possible values are: `android`, `iOS`, `windows`, `windowsPhone`, `macOS`, `all`, `unknownFutureValue`, `linux`. The _provider_ default value is `[]`.",
						},
						"include_platforms": schema.SetAttribute{
							ElementType: types.StringType,
							Required:    true,
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.OneOf("android", "iOS", "windows", "windowsPhone", "macOS", "all", "unknownFutureValue", "linux"),
								),
							},
							MarkdownDescription: "linux``. / Possible values are: `android`, `iOS`, `windows`, `windowsPhone`, `macOS`, `all`, `unknownFutureValue`, `linux`",
						},
					},
					MarkdownDescription: "Platforms included in and excluded from the policy. / Platforms included in and excluded from the policy scope. / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessplatforms?view=graph-rest-beta",
				},
				"service_principal_risk_levels": schema.SetAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Validators: []validator.Set{
						setvalidator.ValueStringsAre(
							stringvalidator.OneOf("low", "medium", "high", "hidden", "none", "unknownFutureValue"),
						),
						wpvalidator.TranslateGraphMissingWithTfDefault(tftypes.Set{ElementType: tftypes.String}),
					},
					PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "Service principal risk levels included in the policy. / Possible values are: `low`, `medium`, `high`, `hidden`, `none`, `unknownFutureValue`. The _provider_ default value is `[]`.",
				},
				"sign_in_risk_levels": schema.SetAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Validators: []validator.Set{
						setvalidator.ValueStringsAre(
							stringvalidator.OneOf("low", "medium", "high", "hidden", "none", "unknownFutureValue"),
						),
					},
					PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "Sign-in risk levels included in the policy. Required. / Possible values are: `low`, `medium`, `high`, `hidden`, `none`, `unknownFutureValue`. The _provider_ default value is `[]`.",
				},
				"user_risk_levels": schema.SetAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Validators: []validator.Set{
						setvalidator.ValueStringsAre(
							stringvalidator.OneOf("low", "medium", "high", "hidden", "none", "unknownFutureValue"),
						),
					},
					PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "User risk levels included in the policy. Required. / Possible values are: `low`, `medium`, `high`, `hidden`, `none`, `unknownFutureValue`. The _provider_ default value is `[]`.",
				},
				"users": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // conditionalAccessUsers
						"exclude_groups": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "Group IDs excluded from scope of policy. The _provider_ default value is `[]`.",
						},
						"exclude_guests_or_external_users": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // conditionalAccessGuestsOrExternalUsers
								"external_tenants": schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // conditionalAccessExternalTenants
										"membership_kind": schema.StringAttribute{
											Optional: true,
											Validators: []validator.String{
												stringvalidator.OneOf("all", "enumerated", "unknownFutureValue"),
											},
											MarkdownDescription: "The membership kind. The `enumerated` member references an [conditionalAccessEnumeratedExternalTenants](conditionalaccessenumeratedexternaltenants.md) object. / Possible values are: `all`, `enumerated`, `unknownFutureValue`",
										},
										"all": generic.OdataDerivedTypeNestedAttributeRs{
											DerivedType: "#microsoft.graph.conditionalAccessAllExternalTenants",
											SingleNestedAttribute: schema.SingleNestedAttribute{
												Optional:   true,
												Attributes: map[string]schema.Attribute{ // conditionalAccessAllExternalTenants
												},
												Validators: []validator.Object{
													conditionalAccessPolicyConditionalAccessExternalTenantsValidator,
												},
												MarkdownDescription: "Represents all external tenants in a policy scope. / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessallexternaltenants?view=graph-rest-beta",
											},
										},
										"enumerated": generic.OdataDerivedTypeNestedAttributeRs{
											DerivedType: "#microsoft.graph.conditionalAccessEnumeratedExternalTenants",
											SingleNestedAttribute: schema.SingleNestedAttribute{
												Optional: true,
												Attributes: map[string]schema.Attribute{ // conditionalAccessEnumeratedExternalTenants
													"members": schema.SetAttribute{
														ElementType:         types.StringType,
														Required:            true,
														Validators:          []validator.Set{setvalidator.SizeAtLeast(1)},
														MarkdownDescription: "A collection of tenant IDs that define the scope of a policy targeting conditional access for guests and external users.",
													},
												},
												Validators: []validator.Object{
													conditionalAccessPolicyConditionalAccessExternalTenantsValidator,
												},
												MarkdownDescription: "Represents a list of external tenants in a policy scope. / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessenumeratedexternaltenants?view=graph-rest-beta",
											},
										},
									},
									MarkdownDescription: "The tenant IDs of the selected types of external users. Either all B2B tenant or a collection of tenant IDs. External tenants can be specified only when the property **guestOrExternalUserTypes** isn't `null` or an empty String. / An abstract type that represents external tenants in a policy scope.\n\nBase type of [conditionalAccessAllExternalTenants](../resources/conditionalaccessallexternaltenants.md) and [conditionalAccessEnumeratedExternalTenants](conditionalaccessenumeratedexternaltenants.md). / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessexternaltenants?view=graph-rest-beta",
								},
								"guest_or_external_user_types": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										wpvalidator.FlagEnumValues("none", "internalGuest", "b2bCollaborationGuest", "b2bCollaborationMember", "b2bDirectConnectUser", "otherExternalUser", "serviceProvider", "unknownFutureValue"),
									},
									MarkdownDescription: "Indicates internal guests or external user types, and is a multi-valued property. / Possible values are: `none`, `internalGuest`, `b2bCollaborationGuest`, `b2bCollaborationMember`, `b2bDirectConnectUser`, `otherExternalUser`, `serviceProvider`, `unknownFutureValue`",
								},
							},
							MarkdownDescription: "Internal guests or external users excluded from the policy scope. Optionally populated. / Represents internal guests and external users in a policy scope. / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessguestsorexternalusers?view=graph-rest-beta",
						},
						"exclude_roles": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "Role IDs excluded from scope of policy. The _provider_ default value is `[]`.",
						},
						"exclude_users": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "User IDs excluded from scope of policy and/or `GuestsOrExternalUsers`. The _provider_ default value is `[]`.",
						},
						"include_groups": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "Group IDs in scope of policy unless explicitly excluded. The _provider_ default value is `[]`.",
						},
						"include_guests_or_external_users": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // conditionalAccessGuestsOrExternalUsers
								"external_tenants": schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // conditionalAccessExternalTenants
										"membership_kind": schema.StringAttribute{
											Optional: true,
											Validators: []validator.String{
												stringvalidator.OneOf("all", "enumerated", "unknownFutureValue"),
											},
											MarkdownDescription: "The membership kind. The `enumerated` member references an [conditionalAccessEnumeratedExternalTenants](conditionalaccessenumeratedexternaltenants.md) object. / Possible values are: `all`, `enumerated`, `unknownFutureValue`",
										},
										"all": generic.OdataDerivedTypeNestedAttributeRs{
											DerivedType: "#microsoft.graph.conditionalAccessAllExternalTenants",
											SingleNestedAttribute: schema.SingleNestedAttribute{
												Optional:   true,
												Attributes: map[string]schema.Attribute{ // conditionalAccessAllExternalTenants
												},
												Validators: []validator.Object{
													conditionalAccessPolicyConditionalAccessExternalTenantsValidator,
												},
												MarkdownDescription: "Represents all external tenants in a policy scope. / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessallexternaltenants?view=graph-rest-beta",
											},
										},
										"enumerated": generic.OdataDerivedTypeNestedAttributeRs{
											DerivedType: "#microsoft.graph.conditionalAccessEnumeratedExternalTenants",
											SingleNestedAttribute: schema.SingleNestedAttribute{
												Optional: true,
												Attributes: map[string]schema.Attribute{ // conditionalAccessEnumeratedExternalTenants
													"members": schema.SetAttribute{
														ElementType:         types.StringType,
														Required:            true,
														Validators:          []validator.Set{setvalidator.SizeAtLeast(1)},
														MarkdownDescription: "A collection of tenant IDs that define the scope of a policy targeting conditional access for guests and external users.",
													},
												},
												Validators: []validator.Object{
													conditionalAccessPolicyConditionalAccessExternalTenantsValidator,
												},
												MarkdownDescription: "Represents a list of external tenants in a policy scope. / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessenumeratedexternaltenants?view=graph-rest-beta",
											},
										},
									},
									MarkdownDescription: "The tenant IDs of the selected types of external users. Either all B2B tenant or a collection of tenant IDs. External tenants can be specified only when the property **guestOrExternalUserTypes** isn't `null` or an empty String. / An abstract type that represents external tenants in a policy scope.\n\nBase type of [conditionalAccessAllExternalTenants](../resources/conditionalaccessallexternaltenants.md) and [conditionalAccessEnumeratedExternalTenants](conditionalaccessenumeratedexternaltenants.md). / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessexternaltenants?view=graph-rest-beta",
								},
								"guest_or_external_user_types": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										wpvalidator.FlagEnumValues("none", "internalGuest", "b2bCollaborationGuest", "b2bCollaborationMember", "b2bDirectConnectUser", "otherExternalUser", "serviceProvider", "unknownFutureValue"),
									},
									MarkdownDescription: "Indicates internal guests or external user types, and is a multi-valued property. / Possible values are: `none`, `internalGuest`, `b2bCollaborationGuest`, `b2bCollaborationMember`, `b2bDirectConnectUser`, `otherExternalUser`, `serviceProvider`, `unknownFutureValue`",
								},
							},
							MarkdownDescription: "Internal guests or external users included in the policy scope. Optionally populated. / Represents internal guests and external users in a policy scope. / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessguestsorexternalusers?view=graph-rest-beta",
						},
						"include_roles": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "Role IDs in scope of policy unless explicitly excluded. The _provider_ default value is `[]`.",
						},
						"include_users": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "User IDs in scope of policy unless explicitly excluded, `None`, `All`, or `GuestsOrExternalUsers`. The _provider_ default value is `[]`.",
						},
					},
					MarkdownDescription: "Users, groups, and roles included in and excluded from the policy. Either **users** or **clientApplications** is required. / Represents users, groups, and roles included in and excluded from the policy scope. / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessusers?view=graph-rest-beta",
				},
			},
			MarkdownDescription: "Specifies the rules that must be met for the policy to apply. Required. / Represents the type of conditions that govern when the policy applies. / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessconditionset?view=graph-rest-beta",
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Readonly.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Specifies a display name for the conditionalAccessPolicy object.",
		},
		"grant_controls": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // conditionalAccessGrantControls
				"built_in_controls": schema.SetAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Validators: []validator.Set{
						setvalidator.ValueStringsAre(
							stringvalidator.OneOf("block", "mfa", "compliantDevice", "domainJoinedDevice", "approvedApplication", "compliantApplication", "passwordChange", "unknownFutureValue"),
						),
					},
					PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "List of values of built-in controls required by the policy. Possible values: `block`, `mfa`, `compliantDevice`, `domainJoinedDevice`, `approvedApplication`, `compliantApplication`, `passwordChange`, `unknownFutureValue`. / Possible values are: `block`, `mfa`, `compliantDevice`, `domainJoinedDevice`, `approvedApplication`, `compliantApplication`, `passwordChange`, `unknownFutureValue`. The _provider_ default value is `[]`.",
				},
				"custom_authentication_factors": schema.SetAttribute{
					ElementType:         types.StringType,
					Optional:            true,
					PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "List of custom controls IDs required by the policy. To learn more about custom control, see [Custom controls (preview)](/azure/active-directory/conditional-access/controls#custom-controls-preview). The _provider_ default value is `[]`.",
				},
				"operator": schema.StringAttribute{
					Required:            true,
					Validators:          []validator.String{stringvalidator.OneOf("AND", "OR")},
					MarkdownDescription: "Defines the relationship of the grant controls. Possible values: `AND`, `OR`.",
				},
				"terms_of_use": schema.SetAttribute{
					ElementType:         types.StringType,
					Optional:            true,
					PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "List of [terms of use](agreement.md) IDs required by the policy. The _provider_ default value is `[]`.",
				},
				"authentication_strength": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // authenticationStrengthPolicy
						"id": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "The system-generated identifier for this mode.",
						},
						"allowed_combinations": schema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									wpvalidator.FlagEnumValues("password", "voice", "hardwareOath", "softwareOath", "sms", "fido2", "windowsHelloForBusiness", "microsoftAuthenticatorPush", "deviceBasedPush", "temporaryAccessPassOneTime", "temporaryAccessPassMultiUse", "email", "x509CertificateSingleFactor", "x509CertificateMultiFactor", "federatedSingleFactor", "federatedMultiFactor", "unknownFutureValue"),
								),
							},
							PlanModifiers:       []planmodifier.Set{wpplanmodifier.SetUseStateForUnknown()},
							MarkdownDescription: "A collection of authentication method modes that are required be used to satify this authentication strength. / Possible values are: `password`, `voice`, `hardwareOath`, `softwareOath`, `sms`, `fido2`, `windowsHelloForBusiness`, `microsoftAuthenticatorPush`, `deviceBasedPush`, `temporaryAccessPassOneTime`, `temporaryAccessPassMultiUse`, `email`, `x509CertificateSingleFactor`, `x509CertificateMultiFactor`, `federatedSingleFactor`, `federatedMultiFactor`, `unknownFutureValue`",
						},
						"created_date_time": schema.StringAttribute{
							Computed:            true,
							PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
							MarkdownDescription: "The datetime when this policy was created.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
							MarkdownDescription: "The human-readable description of this policy.",
						},
						"display_name": schema.StringAttribute{
							Computed:            true,
							PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
							MarkdownDescription: "The human-readable display name of this policy. <br><br>Supports `$filter` (`eq`, `ne`, `not` , and `in`).",
						},
						"modified_date_time": schema.StringAttribute{
							Computed:            true,
							PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
							MarkdownDescription: "The datetime when this policy was last modified.",
						},
						"policy_type": schema.StringAttribute{
							Computed: true,
							Validators: []validator.String{
								stringvalidator.OneOf("builtIn", "custom", "unknownFutureValue"),
							},
							PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
							MarkdownDescription: "A descriptor of whether this policy is built into Microsoft Entra Conditional Access or created by an admin for the tenant. The <br><br>Supports `$filter` (`eq`, `ne`, `not` , and `in`). / Possible values are: `builtIn`, `custom`, `unknownFutureValue`",
						},
						"requirements_satisfied": schema.StringAttribute{
							Computed: true,
							Validators: []validator.String{
								wpvalidator.FlagEnumValues("none", "mfa", "unknownFutureValue"),
							},
							PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
							MarkdownDescription: "A descriptor of whether this authentication strength grants the MFA claim upon successful satisfaction. The / Possible values are: `none`, `mfa`, `unknownFutureValue`",
						},
					},
					MarkdownDescription: "The authentication strength required by the conditional access policy. Optional. / A collection of settings that define specific combinations of authentication methods and metadata. The authentication strength policy, when applied to a given scenario using Microsoft Entra Conditional Access, defines which authentication methods must be used to authenticate in that scenario. An authentication strength may be built-in or custom (defined by the tenant) and may or may not fulfill the requirements to grant an MFA claim. / https://learn.microsoft.com/en-us/graph/api/resources/authenticationstrengthpolicy?view=graph-rest-beta",
				},
			},
			MarkdownDescription: "Specifies the grant controls that must be fulfilled to pass the policy. / Represents grant controls that must be fulfilled to pass the policy. / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccessgrantcontrols?view=graph-rest-beta",
		},
		"session_controls": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // conditionalAccessSessionControls
				"application_enforced_restrictions": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // applicationEnforcedRestrictionsSessionControl
						"is_enabled": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Specifies whether the session control is enabled or not.",
						},
					},
					MarkdownDescription: "Session control to enforce application restrictions. Only Exchange Online and Sharepoint Online support this session control. / Session control to enforce application restrictions. Inherit from [Conditional Access Session Control](conditionalaccesssessioncontrol.md). / https://learn.microsoft.com/en-us/graph/api/resources/applicationenforcedrestrictionssessioncontrol?view=graph-rest-beta",
				},
				"cloud_app_security": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // cloudAppSecuritySessionControl
						"is_enabled": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Specifies whether the session control is enabled.",
						},
						"cloud_app_security_type": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("mcasConfigured", "monitorOnly", "blockDownloads", "unknownFutureValue"),
							},
							MarkdownDescription: "To learn more about these values, [Deploy Conditional Access App Control for featured apps](/cloud-app-security/proxy-deployment-aad#step-1--configure-your-idp-to-work-with-cloud-app-security). / Possible values are: `mcasConfigured`, `monitorOnly`, `blockDownloads`, `unknownFutureValue`",
						},
					},
					MarkdownDescription: "Session control to apply cloud app security. / Session control used to enforce cloud app security checks. Inehrits from [Conditional Access Session Control](conditionalaccesssessioncontrol.md). / https://learn.microsoft.com/en-us/graph/api/resources/cloudappsecuritysessioncontrol?view=graph-rest-beta",
				},
				"continuous_access_evaluation": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // continuousAccessEvaluationSessionControl
						"mode": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("strictEnforcement", "disabled", "unknownFutureValue", "strictLocation"),
							},
							MarkdownDescription: "Specifies continuous access evaluation settings. The Note that you must use the `Prefer: include-unknown-enum-members` request header to get the following value(s) in this [evolvable enum](/graph/best-practices-concept#handling-future-members-in-evolvable-enumerations): `strictLocation`. / Possible values are: `strictEnforcement`, `disabled`, `unknownFutureValue`, `strictLocation`",
						},
					},
					MarkdownDescription: "Session control for continuous access evaluation settings. / Session control to control continuous access evaluation settings. / https://learn.microsoft.com/en-us/graph/api/resources/continuousaccessevaluationsessioncontrol?view=graph-rest-beta",
				},
				"disable_resilience_defaults": schema.BoolAttribute{
					Optional: true,
					Validators: []validator.Bool{
						wpvalidator.TranslateGraphNullWithTfDefault(tftypes.Bool),
					},
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
					Computed:            true,
					MarkdownDescription: "Session control that determines whether it's acceptable for Microsoft Entra ID to extend existing sessions based on information collected prior to an outage or not. The _provider_ default value is `false`.",
				},
				"persistent_browser": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // persistentBrowserSessionControl
						"is_enabled": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Specifies whether the session control is enabled.",
						},
						"mode": schema.StringAttribute{
							Optional:            true,
							Validators:          []validator.String{stringvalidator.OneOf("always", "never")},
							MarkdownDescription: "Possible values are: `always`, `never`",
						},
					},
					MarkdownDescription: "Session control to define whether to persist cookies or not. All apps should be selected for this session control to work correctly. / Session control to define whether to persist cookies or not. / https://learn.microsoft.com/en-us/graph/api/resources/persistentbrowsersessioncontrol?view=graph-rest-beta",
				},
				"secure_sign_in_session": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // secureSignInSessionControl
						"is_enabled": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Specifies whether the session control is enabled.",
						},
					},
					MarkdownDescription: "Session control to require sign in sessions to be bound to a device. / Session control to require sign in sessions to be bound to a device. / https://learn.microsoft.com/en-us/graph/api/resources/securesigninsessioncontrol?view=graph-rest-beta",
				},
				"sign_in_frequency": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // signInFrequencySessionControl
						"is_enabled": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Specifies whether the session control is enabled.",
						},
						"authentication_type": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("primaryAndSecondaryAuthentication", "secondaryAuthentication", "unknownFutureValue"),
							},
							MarkdownDescription: "The possible values are `primaryAndSecondaryAuthentication`, `secondaryAuthentication`, `unknownFutureValue`. This property isn't required when using **frequencyInterval** with the value of `timeBased`. / Possible values are: `primaryAndSecondaryAuthentication`, `secondaryAuthentication`, `unknownFutureValue`",
						},
						"frequency_interval": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("timeBased", "everyTime", "unknownFutureValue"),
							},
							MarkdownDescription: "The possible values are `timeBased`, `everyTime`, `unknownFutureValue`. Sign-in frequency of `everyTime` is available for risky users, risky sign-ins, Intune device enrollment, any application, authentication context, and user actions. For more information, see [Require reauthentication every time](https://aka.ms/RequireReauthentication). / Possible values are: `timeBased`, `everyTime`, `unknownFutureValue`",
						},
						"type": schema.StringAttribute{
							Optional:            true,
							Validators:          []validator.String{stringvalidator.OneOf("days", "hours")},
							MarkdownDescription: "or `null` if **frequencyInterval** is `everyTime` . / Possible values are: `days`, `hours`",
						},
						"value": schema.Int64Attribute{
							Optional:            true,
							MarkdownDescription: "The number of `days` or `hours`.",
						},
					},
					MarkdownDescription: "Session control to enforce signin frequency. / Session control to enforce sign-in frequency. / https://learn.microsoft.com/en-us/graph/api/resources/signinfrequencysessioncontrol?view=graph-rest-beta",
				},
			},
			MarkdownDescription: "Specifies the session controls that are enforced after sign-in. / Represents session controls that are enforced after sign-in.\nAll the session controls inherit from [conditionalAccessSessionControl](conditionalaccesssessioncontrol.md). / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccesssessioncontrols?view=graph-rest-beta",
		},
		"state": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("enabled", "disabled", "enabledForReportingButNotEnforced"),
			},
			PlanModifiers: []planmodifier.String{
				wpdefaultvalue.StringDefaultValue("enabledForReportingButNotEnforced"),
			},
			Computed:            true,
			MarkdownDescription: "Specifies the state of the conditionalAccessPolicy object. Required. / Possible values are: `enabled`, `disabled`, `enabledForReportingButNotEnforced`. The _provider_ default value is `\"enabledForReportingButNotEnforced\"`.",
		},
	},
	MarkdownDescription: "Represents a Microsoft Entra Conditional Access policy. Conditional access policies are custom rules that define an access scenario. For more information, see the [Conditional access documentation](/azure/active-directory/conditional-access/). / https://learn.microsoft.com/en-us/graph/api/resources/conditionalaccesspolicy?view=graph-rest-beta ||| MS Graph: Conditional access",
}

var conditionalAccessPolicyConditionalAccessExternalTenantsValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("all"),
	path.MatchRelative().AtParent().AtName("enumerated"),
)
