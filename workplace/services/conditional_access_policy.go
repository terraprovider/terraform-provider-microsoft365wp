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
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
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
									MarkdownDescription: "; possible values are: `include`, `exclude`",
								},
								"rule": schema.StringAttribute{
									Required: true,
								},
							},
							MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessFilter?view=graph-rest-beta",
						},
						"exclude_applications": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
						},
						"include_applications": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
						},
						"include_authentication_context_class_references": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
						},
						"include_user_actions": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
						},
					},
					MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessApplications?view=graph-rest-beta",
				},
				"authentication_flows": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // conditionalAccessAuthenticationFlows
						"transfer_methods": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								wpvalidator.FlagEnumValues("none", "deviceCodeFlow", "authenticationTransfer", "unknownFutureValue"),
							},
							MarkdownDescription: "; possible values are: `none`, `deviceCodeFlow`, `authenticationTransfer`, `unknownFutureValue`",
						},
					},
					MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessAuthenticationFlows?view=graph-rest-beta",
				},
				"client_applications": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // conditionalAccessClientApplications
						"exclude_service_principals": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
						},
						"include_service_principals": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
						},
						"service_principal_filter": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // conditionalAccessFilter
								"mode": schema.StringAttribute{
									Required:            true,
									Validators:          []validator.String{stringvalidator.OneOf("include", "exclude")},
									MarkdownDescription: "; possible values are: `include`, `exclude`",
								},
								"rule": schema.StringAttribute{
									Required: true,
								},
							},
							MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessFilter?view=graph-rest-beta",
						},
					},
					MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessClientApplications?view=graph-rest-beta",
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
					MarkdownDescription: "; possible values are: `all`, `browser`, `mobileAppsAndDesktopClients`, `exchangeActiveSync`, `easSupported`, `other`, `unknownFutureValue`",
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
									MarkdownDescription: "; possible values are: `include`, `exclude`",
								},
								"rule": schema.StringAttribute{
									Required: true,
								},
							},
							MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessFilter?view=graph-rest-beta",
						},
						"exclude_devices": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
						},
						"exclude_device_states": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
						},
						"include_devices": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
						},
						"include_device_states": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
						},
					},
					MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessDevices?view=graph-rest-beta",
				},
				"device_states": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // conditionalAccessDeviceStates
						"exclude_states": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
						},
						"include_states": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
						},
					},
					MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessDeviceStates?view=graph-rest-beta",
				},
				"insider_risk_levels": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						wpvalidator.FlagEnumValues("minor", "moderate", "elevated", "unknownFutureValue"),
					},
					MarkdownDescription: "; possible values are: `minor`, `moderate`, `elevated`, `unknownFutureValue`",
				},
				"locations": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // conditionalAccessLocations
						"exclude_locations": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
						},
						"include_locations": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
						},
					},
					MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessLocations?view=graph-rest-beta",
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
							MarkdownDescription: "; possible values are: `android`, `iOS`, `windows`, `windowsPhone`, `macOS`, `all`, `unknownFutureValue`, `linux`",
						},
						"include_platforms": schema.SetAttribute{
							ElementType: types.StringType,
							Required:    true,
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.OneOf("android", "iOS", "windows", "windowsPhone", "macOS", "all", "unknownFutureValue", "linux"),
								),
							},
							MarkdownDescription: "; possible values are: `android`, `iOS`, `windows`, `windowsPhone`, `macOS`, `all`, `unknownFutureValue`, `linux`",
						},
					},
					MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessPlatforms?view=graph-rest-beta",
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
					MarkdownDescription: "; possible values are: `low`, `medium`, `high`, `hidden`, `none`, `unknownFutureValue`",
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
					MarkdownDescription: "; possible values are: `low`, `medium`, `high`, `hidden`, `none`, `unknownFutureValue`",
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
					MarkdownDescription: "; possible values are: `low`, `medium`, `high`, `hidden`, `none`, `unknownFutureValue`",
				},
				"users": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // conditionalAccessUsers
						"exclude_groups": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
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
											MarkdownDescription: "; possible values are: `all`, `enumerated`, `unknownFutureValue`",
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
												MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessAllExternalTenants?view=graph-rest-beta",
											},
										},
										"enumerated": generic.OdataDerivedTypeNestedAttributeRs{
											DerivedType: "#microsoft.graph.conditionalAccessEnumeratedExternalTenants",
											SingleNestedAttribute: schema.SingleNestedAttribute{
												Optional: true,
												Attributes: map[string]schema.Attribute{ // conditionalAccessEnumeratedExternalTenants
													"members": schema.SetAttribute{
														ElementType: types.StringType,
														Required:    true,
														Validators:  []validator.Set{setvalidator.SizeAtLeast(1)},
													},
												},
												Validators: []validator.Object{
													conditionalAccessPolicyConditionalAccessExternalTenantsValidator,
												},
												MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessEnumeratedExternalTenants?view=graph-rest-beta",
											},
										},
									},
									MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessExternalTenants?view=graph-rest-beta",
								},
								"guest_or_external_user_types": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										wpvalidator.FlagEnumValues("none", "internalGuest", "b2bCollaborationGuest", "b2bCollaborationMember", "b2bDirectConnectUser", "otherExternalUser", "serviceProvider", "unknownFutureValue"),
									},
									MarkdownDescription: "; possible values are: `none`, `internalGuest`, `b2bCollaborationGuest`, `b2bCollaborationMember`, `b2bDirectConnectUser`, `otherExternalUser`, `serviceProvider`, `unknownFutureValue`",
								},
							},
							MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessGuestsOrExternalUsers?view=graph-rest-beta",
						},
						"exclude_roles": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
						},
						"exclude_users": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
						},
						"include_groups": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
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
											MarkdownDescription: "; possible values are: `all`, `enumerated`, `unknownFutureValue`",
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
												MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessAllExternalTenants?view=graph-rest-beta",
											},
										},
										"enumerated": generic.OdataDerivedTypeNestedAttributeRs{
											DerivedType: "#microsoft.graph.conditionalAccessEnumeratedExternalTenants",
											SingleNestedAttribute: schema.SingleNestedAttribute{
												Optional: true,
												Attributes: map[string]schema.Attribute{ // conditionalAccessEnumeratedExternalTenants
													"members": schema.SetAttribute{
														ElementType: types.StringType,
														Required:    true,
														Validators:  []validator.Set{setvalidator.SizeAtLeast(1)},
													},
												},
												Validators: []validator.Object{
													conditionalAccessPolicyConditionalAccessExternalTenantsValidator,
												},
												MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessEnumeratedExternalTenants?view=graph-rest-beta",
											},
										},
									},
									MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessExternalTenants?view=graph-rest-beta",
								},
								"guest_or_external_user_types": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										wpvalidator.FlagEnumValues("none", "internalGuest", "b2bCollaborationGuest", "b2bCollaborationMember", "b2bDirectConnectUser", "otherExternalUser", "serviceProvider", "unknownFutureValue"),
									},
									MarkdownDescription: "; possible values are: `none`, `internalGuest`, `b2bCollaborationGuest`, `b2bCollaborationMember`, `b2bDirectConnectUser`, `otherExternalUser`, `serviceProvider`, `unknownFutureValue`",
								},
							},
							MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessGuestsOrExternalUsers?view=graph-rest-beta",
						},
						"include_roles": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
						},
						"include_users": schema.SetAttribute{
							ElementType:   types.StringType,
							Optional:      true,
							PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:      true,
						},
					},
					MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessUsers?view=graph-rest-beta",
				},
			},
			MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessConditionSet?view=graph-rest-beta",
		},
		"created_date_time": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"display_name": schema.StringAttribute{
			Required: true,
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
					MarkdownDescription: "; possible values are: `block`, `mfa`, `compliantDevice`, `domainJoinedDevice`, `approvedApplication`, `compliantApplication`, `passwordChange`, `unknownFutureValue`",
				},
				"custom_authentication_factors": schema.SetAttribute{
					ElementType:   types.StringType,
					Optional:      true,
					PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
					Computed:      true,
				},
				"operator": schema.StringAttribute{
					Required:   true,
					Validators: []validator.String{stringvalidator.OneOf("AND", "OR")},
				},
				"terms_of_use": schema.SetAttribute{
					ElementType:   types.StringType,
					Optional:      true,
					PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
					Computed:      true,
				},
				"authentication_strength": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // authenticationStrengthPolicy
						"id": schema.StringAttribute{
							Required: true,
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
							MarkdownDescription: "; possible values are: `password`, `voice`, `hardwareOath`, `softwareOath`, `sms`, `fido2`, `windowsHelloForBusiness`, `microsoftAuthenticatorPush`, `deviceBasedPush`, `temporaryAccessPassOneTime`, `temporaryAccessPassMultiUse`, `email`, `x509CertificateSingleFactor`, `x509CertificateMultiFactor`, `federatedSingleFactor`, `federatedMultiFactor`, `unknownFutureValue`",
						},
						"created_date_time": schema.StringAttribute{
							Computed:      true,
							PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						},
						"description": schema.StringAttribute{
							Computed:      true,
							PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						},
						"display_name": schema.StringAttribute{
							Computed:      true,
							PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						},
						"modified_date_time": schema.StringAttribute{
							Computed:      true,
							PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						},
						"policy_type": schema.StringAttribute{
							Computed: true,
							Validators: []validator.String{
								stringvalidator.OneOf("builtIn", "custom", "unknownFutureValue"),
							},
							PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
							MarkdownDescription: "; possible values are: `builtIn`, `custom`, `unknownFutureValue`",
						},
						"requirements_satisfied": schema.StringAttribute{
							Computed: true,
							Validators: []validator.String{
								wpvalidator.FlagEnumValues("none", "mfa", "unknownFutureValue"),
							},
							PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
							MarkdownDescription: "; possible values are: `none`, `mfa`, `unknownFutureValue`",
						},
					},
					MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/authenticationStrengthPolicy?view=graph-rest-beta",
				},
			},
			MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessGrantControls?view=graph-rest-beta",
		},
		"session_controls": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // conditionalAccessSessionControls
				"application_enforced_restrictions": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // applicationEnforcedRestrictionsSessionControl
						"is_enabled": schema.BoolAttribute{
							Optional: true,
						},
					},
					MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/applicationEnforcedRestrictionsSessionControl?view=graph-rest-beta",
				},
				"cloud_app_security": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // cloudAppSecuritySessionControl
						"is_enabled": schema.BoolAttribute{
							Optional: true,
						},
						"cloud_app_security_type": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("mcasConfigured", "monitorOnly", "blockDownloads", "unknownFutureValue"),
							},
							MarkdownDescription: "; possible values are: `mcasConfigured`, `monitorOnly`, `blockDownloads`, `unknownFutureValue`",
						},
					},
					MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/cloudAppSecuritySessionControl?view=graph-rest-beta",
				},
				"continuous_access_evaluation": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // continuousAccessEvaluationSessionControl
						"mode": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("strictEnforcement", "disabled", "unknownFutureValue", "strictLocation"),
							},
							MarkdownDescription: "; possible values are: `strictEnforcement`, `disabled`, `unknownFutureValue`, `strictLocation`",
						},
					},
					MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/continuousAccessEvaluationSessionControl?view=graph-rest-beta",
				},
				"disable_resilience_defaults": schema.BoolAttribute{
					Optional: true,
					Validators: []validator.Bool{
						wpvalidator.TranslateGraphNullWithTfDefault(tftypes.Bool),
					},
					PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
					Computed:      true,
				},
				"persistent_browser": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // persistentBrowserSessionControl
						"is_enabled": schema.BoolAttribute{
							Optional: true,
						},
						"mode": schema.StringAttribute{
							Optional:            true,
							Validators:          []validator.String{stringvalidator.OneOf("always", "never")},
							MarkdownDescription: "; possible values are: `always`, `never`",
						},
					},
					MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/persistentBrowserSessionControl?view=graph-rest-beta",
				},
				"secure_sign_in_session": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // secureSignInSessionControl
						"is_enabled": schema.BoolAttribute{
							Optional: true,
						},
					},
					MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/secureSignInSessionControl?view=graph-rest-beta",
				},
				"sign_in_frequency": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // signInFrequencySessionControl
						"is_enabled": schema.BoolAttribute{
							Optional: true,
						},
						"authentication_type": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("primaryAndSecondaryAuthentication", "secondaryAuthentication", "unknownFutureValue"),
							},
							MarkdownDescription: "; possible values are: `primaryAndSecondaryAuthentication`, `secondaryAuthentication`, `unknownFutureValue`",
						},
						"frequency_interval": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("timeBased", "everyTime", "unknownFutureValue"),
							},
							MarkdownDescription: "; possible values are: `timeBased`, `everyTime`, `unknownFutureValue`",
						},
						"type": schema.StringAttribute{
							Optional:            true,
							Validators:          []validator.String{stringvalidator.OneOf("days", "hours")},
							MarkdownDescription: "; possible values are: `days`, `hours`",
						},
						"value": schema.Int64Attribute{
							Optional: true,
						},
					},
					MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/signInFrequencySessionControl?view=graph-rest-beta",
				},
			},
			MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessSessionControls?view=graph-rest-beta",
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
			MarkdownDescription: "; possible values are: `enabled`, `disabled`, `enabledForReportingButNotEnforced`",
		},
	},
	MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/conditionalAccessPolicy?view=graph-rest-beta",
}

var conditionalAccessPolicyConditionalAccessExternalTenantsValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("all"),
	path.MatchRelative().AtParent().AtName("enumerated"),
)
