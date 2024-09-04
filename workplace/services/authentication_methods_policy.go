package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	AuthenticationMethodsPolicyResource = generic.GenericResource{
		TypeNameSuffix: "authentication_methods_policy",
		SpecificSchema: authenticationMethodsPolicyResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:     "/policies/authenticationmethodspolicy",
			IsSingleton: true,
		},
	}

	AuthenticationMethodsPolicySingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&AuthenticationMethodsPolicyResource)
)

var authenticationMethodsPolicyResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // authenticationMethodsPolicy
		"id": schema.StringAttribute{
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
		"last_modified_date_time": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"policy_migration_state": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("preMigration", "migrationInProgress", "migrationComplete", "unknownFutureValue"),
			},
			PlanModifiers: []planmodifier.String{
				wpdefaultvalue.StringDefaultValue("migrationInProgress"),
			},
			Computed: true,
		},
		"policy_version": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"registration_enforcement": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // registrationEnforcement
				"authentication_methods_registration_campaign": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // authenticationMethodsRegistrationCampaign
						"enforce_registration_after_allowed_snoozes": schema.BoolAttribute{
							Optional:      true,
							PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
							Computed:      true,
						},
						"exclude_targets": schema.SetNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: authenticationMethodsPolicyExcludeTargetAttributes,
							},
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/excludeTarget?view=graph-rest-beta`,
						},
						"include_targets": schema.SetNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{ // authenticationMethodsRegistrationCampaignIncludeTarget
									"id": schema.StringAttribute{
										Required: true,
									},
									"targeted_authentication_method": schema.StringAttribute{
										Optional: true,
										PlanModifiers: []planmodifier.String{
											wpdefaultvalue.StringDefaultValue("microsoftAuthenticator"),
										},
										Computed: true,
									},
									"target_type": schema.StringAttribute{
										Required: true,
										Validators: []validator.String{
											stringvalidator.OneOf("user", "group", "unknownFutureValue"),
										},
									},
								},
							},
							PlanModifiers:       []planmodifier.Set{authenticationMethodsPolicyTargetsDefaultAllUsers},
							Computed:            true,
							MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/authenticationMethodsRegistrationCampaignIncludeTarget?view=graph-rest-beta`,
						},
						"snooze_duration_in_days": schema.Int64Attribute{
							Optional:      true,
							PlanModifiers: []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(1)},
							Computed:      true,
						},
						"state": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("default", "enabled", "disabled", "unknownFutureValue"),
							},
							PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("default")},
							Computed:      true,
						},
					},
					PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/authenticationMethodsRegistrationCampaign?view=graph-rest-beta`,
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/registrationEnforcement?view=graph-rest-beta`,
		},
		"report_suspicious_activity_settings": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // reportSuspiciousActivitySettings
				"include_target": schema.SingleNestedAttribute{
					Optional:            true,
					Attributes:          authenticationMethodsPolicyIncludeTargetAttributes,
					PlanModifiers:       []planmodifier.Object{authenticationMethodsPolicyTargetDefaultAllUsers},
					Computed:            true,
					MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/includeTarget?view=graph-rest-beta`,
				},
				"state": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf("default", "enabled", "disabled", "unknownFutureValue"),
					},
					PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("default")},
					Computed:      true,
				},
				"voice_reporting_code": schema.Int64Attribute{
					Optional:      true,
					PlanModifiers: []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
					Computed:      true,
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/reportSuspiciousActivitySettings?view=graph-rest-beta`,
		},
		"system_credential_preferences": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // systemCredentialPreferences
				"exclude_targets": schema.SetNestedAttribute{
					Optional: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: authenticationMethodsPolicyExcludeTargetAttributes,
					},
					PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/excludeTarget?view=graph-rest-beta`,
				},
				"include_targets": schema.SetNestedAttribute{
					Optional: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: authenticationMethodsPolicyIncludeTargetAttributes,
					},
					PlanModifiers:       []planmodifier.Set{authenticationMethodsPolicyTargetsDefaultAllUsers},
					Computed:            true,
					MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/includeTarget?view=graph-rest-beta`,
				},
				"state": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf("default", "enabled", "disabled", "unknownFutureValue"),
					},
					PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("default")},
					Computed:      true,
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/systemCredentialPreferences?view=graph-rest-beta`,
		},
		"authentication_method_configurations": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // authenticationMethodConfiguration
					"id": schema.StringAttribute{
						Required: true,
					},
					"exclude_targets": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: authenticationMethodsPolicyExcludeTargetAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/excludeTarget?view=graph-rest-beta`,
					},
					"state": schema.StringAttribute{
						Required:   true,
						Validators: []validator.String{stringvalidator.OneOf("enabled", "disabled")},
					},
					"email": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.emailAuthenticationMethodConfiguration",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // emailAuthenticationMethodConfiguration
								"allow_external_id_to_use_email_otp": schema.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.OneOf("default", "enabled", "disabled", "unknownFutureValue"),
									},
									PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("default")},
									Computed:      true,
								},
								"include_targets": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: authenticationMethodsPolicyAuthenticationMethodTargetAttributes,
									},
									PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/authenticationMethodTarget?view=graph-rest-beta`,
								},
							},
							Validators: []validator.Object{
								authenticationMethodsPolicyAuthenticationMethodConfigurationValidator,
							},
							MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/emailAuthenticationMethodConfiguration?view=graph-rest-beta`,
						},
					},
					"fido2": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.fido2AuthenticationMethodConfiguration",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // fido2AuthenticationMethodConfiguration
								"is_attestation_enforced": schema.BoolAttribute{
									Optional:      true,
									PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
									Computed:      true,
								},
								"is_self_service_registration_allowed": schema.BoolAttribute{
									Optional:      true,
									PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
									Computed:      true,
								},
								"key_restrictions": schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // fido2KeyRestrictions
										"aa_guids": schema.SetAttribute{
											ElementType:   types.StringType,
											Optional:      true,
											PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:      true,
										},
										"enforcement_type": schema.StringAttribute{
											Optional: true,
											Validators: []validator.String{
												stringvalidator.OneOf("allow", "block", "unknownFutureValue"),
											},
											PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("block")},
											Computed:      true,
										},
										"is_enforced": schema.BoolAttribute{
											Optional:      true,
											PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
											Computed:      true,
										},
									},
									PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/fido2KeyRestrictions?view=graph-rest-beta`,
								},
								"include_targets": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: authenticationMethodsPolicyPasskeyAuthenticationMethodTargetAttributes,
									},
									PlanModifiers: []planmodifier.Set{authenticationMethodsPolicyTargetsDefaultAllUsers},
									Computed:      true,
								},
							},
							Validators: []validator.Object{
								authenticationMethodsPolicyAuthenticationMethodConfigurationValidator,
							},
							MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/fido2AuthenticationMethodConfiguration?view=graph-rest-beta`,
						},
					},
					"hardware_oath": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.hardwareOathAuthenticationMethodConfiguration",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // hardwareOathAuthenticationMethodConfiguration
								"include_targets": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: authenticationMethodsPolicyAuthenticationMethodTargetAttributes,
									},
									PlanModifiers:       []planmodifier.Set{authenticationMethodsPolicyTargetsDefaultAllUsers},
									Computed:            true,
									MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/authenticationMethodTarget?view=graph-rest-beta`,
								},
							},
							Validators: []validator.Object{
								authenticationMethodsPolicyAuthenticationMethodConfigurationValidator,
							},
							MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/hardwareOathAuthenticationMethodConfiguration?view=graph-rest-beta`,
						},
					},
					"microsoft_authenticator": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.microsoftAuthenticatorAuthenticationMethodConfiguration",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // microsoftAuthenticatorAuthenticationMethodConfiguration
								"feature_settings": schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // microsoftAuthenticatorFeatureSettings
										"companion_app_allowed_state": schema.SingleNestedAttribute{
											Optional:            true,
											Attributes:          authenticationMethodsPolicyAuthenticationMethodFeatureConfigurationAttributes,
											PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/authenticationMethodFeatureConfiguration?view=graph-rest-beta`,
										},
										"display_app_information_required_state": schema.SingleNestedAttribute{
											Optional:            true,
											Attributes:          authenticationMethodsPolicyAuthenticationMethodFeatureConfigurationAttributes,
											PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/authenticationMethodFeatureConfiguration?view=graph-rest-beta`,
										},
										"display_location_information_required_state": schema.SingleNestedAttribute{
											Optional:            true,
											Attributes:          authenticationMethodsPolicyAuthenticationMethodFeatureConfigurationAttributes,
											PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/authenticationMethodFeatureConfiguration?view=graph-rest-beta`,
										},
									},
									PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/microsoftAuthenticatorFeatureSettings?view=graph-rest-beta`,
								},
								"is_software_oath_enabled": schema.BoolAttribute{
									Optional:      true,
									PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
									Computed:      true,
								},
								"include_targets": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{ // microsoftAuthenticatorAuthenticationMethodTarget
											"id": schema.StringAttribute{
												Required: true,
											},
											"is_registration_required": schema.BoolAttribute{
												Optional:      true,
												PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
												Computed:      true,
											},
											"target_type": schema.StringAttribute{
												Required: true,
												Validators: []validator.String{
													stringvalidator.OneOf("user", "group", "unknownFutureValue"),
												},
											},
											"authentication_mode": schema.StringAttribute{
												Optional: true,
												Validators: []validator.String{
													stringvalidator.OneOf("deviceBasedPush", "push", "any"),
												},
												PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("any")},
												Computed:      true,
											},
										},
									},
									PlanModifiers:       []planmodifier.Set{authenticationMethodsPolicyTargetsDefaultAllUsers},
									Computed:            true,
									MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/microsoftAuthenticatorAuthenticationMethodTarget?view=graph-rest-beta`,
								},
							},
							Validators: []validator.Object{
								authenticationMethodsPolicyAuthenticationMethodConfigurationValidator,
							},
							MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/microsoftAuthenticatorAuthenticationMethodConfiguration?view=graph-rest-beta`,
						},
					},
					"sms": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.smsAuthenticationMethodConfiguration",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // smsAuthenticationMethodConfiguration
								"include_targets": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{ // smsAuthenticationMethodTarget
											"id": schema.StringAttribute{
												Required: true,
											},
											"is_registration_required": schema.BoolAttribute{
												Optional:      true,
												PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
												Computed:      true,
											},
											"target_type": schema.StringAttribute{
												Required: true,
												Validators: []validator.String{
													stringvalidator.OneOf("user", "group", "unknownFutureValue"),
												},
											},
											"is_usable_for_sign_in": schema.BoolAttribute{
												Optional:      true,
												PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
												Computed:      true,
											},
										},
									},
									PlanModifiers:       []planmodifier.Set{authenticationMethodsPolicyTargetsDefaultAllUsers},
									Computed:            true,
									MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/smsAuthenticationMethodTarget?view=graph-rest-beta`,
								},
							},
							Validators: []validator.Object{
								authenticationMethodsPolicyAuthenticationMethodConfigurationValidator,
							},
							MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/smsAuthenticationMethodConfiguration?view=graph-rest-beta`,
						},
					},
					"software_oath": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.softwareOathAuthenticationMethodConfiguration",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // softwareOathAuthenticationMethodConfiguration
								"include_targets": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: authenticationMethodsPolicyAuthenticationMethodTargetAttributes,
									},
									PlanModifiers:       []planmodifier.Set{authenticationMethodsPolicyTargetsDefaultAllUsers},
									Computed:            true,
									MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/authenticationMethodTarget?view=graph-rest-beta`,
								},
							},
							Validators: []validator.Object{
								authenticationMethodsPolicyAuthenticationMethodConfigurationValidator,
							},
							MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/softwareOathAuthenticationMethodConfiguration?view=graph-rest-beta`,
						},
					},
					"temporary_access_pass": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.temporaryAccessPassAuthenticationMethodConfiguration",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // temporaryAccessPassAuthenticationMethodConfiguration
								"default_length": schema.Int64Attribute{
									Optional:      true,
									PlanModifiers: []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(8)},
									Computed:      true,
								},
								"default_lifetime_in_minutes": schema.Int64Attribute{
									Optional:      true,
									PlanModifiers: []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(60)},
									Computed:      true,
								},
								"is_usable_once": schema.BoolAttribute{
									Optional:      true,
									PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
									Computed:      true,
								},
								"maximum_lifetime_in_minutes": schema.Int64Attribute{
									Optional:      true,
									PlanModifiers: []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(480)},
									Computed:      true,
								},
								"minimum_lifetime_in_minutes": schema.Int64Attribute{
									Optional:      true,
									PlanModifiers: []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(60)},
									Computed:      true,
								},
								"include_targets": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: authenticationMethodsPolicyAuthenticationMethodTargetAttributes,
									},
									PlanModifiers:       []planmodifier.Set{authenticationMethodsPolicyTargetsDefaultAllUsers},
									Computed:            true,
									MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/authenticationMethodTarget?view=graph-rest-beta`,
								},
							},
							Validators: []validator.Object{
								authenticationMethodsPolicyAuthenticationMethodConfigurationValidator,
							},
							MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/temporaryAccessPassAuthenticationMethodConfiguration?view=graph-rest-beta`,
						},
					},
					"voice": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.voiceAuthenticationMethodConfiguration",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // voiceAuthenticationMethodConfiguration
								"is_office_phone_allowed": schema.BoolAttribute{
									Optional:      true,
									PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
									Computed:      true,
								},
								"include_targets": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{ // voiceAuthenticationMethodTarget
											"id": schema.StringAttribute{
												Required: true,
											},
											"is_registration_required": schema.BoolAttribute{
												Optional:      true,
												PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
												Computed:      true,
											},
											"target_type": schema.StringAttribute{
												Required: true,
												Validators: []validator.String{
													stringvalidator.OneOf("user", "group", "unknownFutureValue"),
												},
											},
										},
									},
									PlanModifiers:       []planmodifier.Set{authenticationMethodsPolicyTargetsDefaultAllUsers},
									Computed:            true,
									MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/voiceAuthenticationMethodTarget?view=graph-rest-beta`,
								},
							},
							Validators: []validator.Object{
								authenticationMethodsPolicyAuthenticationMethodConfigurationValidator,
							},
							MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/voiceAuthenticationMethodConfiguration?view=graph-rest-beta`,
						},
					},
					"x509_certificate": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.x509CertificateAuthenticationMethodConfiguration",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // x509CertificateAuthenticationMethodConfiguration
								"authentication_mode_configuration": schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // x509CertificateAuthenticationModeConfiguration
										"rules": schema.SetNestedAttribute{
											Optional: true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{ // x509CertificateRule
													"identifier": schema.StringAttribute{
														Optional: true,
													},
													"issuer_subject_identifier": schema.StringAttribute{
														Optional: true,
													},
													"policy_oid_identifier": schema.StringAttribute{
														Optional: true,
													},
													"x509_certificate_authentication_mode": schema.StringAttribute{
														Optional: true,
														Validators: []validator.String{
															stringvalidator.OneOf("x509CertificateSingleFactor", "x509CertificateMultiFactor", "unknownFutureValue"),
														},
														Description: `x509CertificateAuthenticationMode`, // custom MS Graph attribute name
													},
													"x509_certificate_required_affinity_level": schema.StringAttribute{
														Optional: true,
														Validators: []validator.String{
															stringvalidator.OneOf("low", "high", "unknownFutureValue"),
														},
														Description: `x509CertificateRequiredAffinityLevel`, // custom MS Graph attribute name
													},
													"x509_certificate_rule_type": schema.StringAttribute{
														Optional: true,
														Validators: []validator.String{
															stringvalidator.OneOf("issuerSubject", "policyOID", "unknownFutureValue", "issuerSubjectAndPolicyOID"),
														},
														Description: `x509CertificateRuleType`, // custom MS Graph attribute name
													},
												},
											},
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/x509CertificateRule?view=graph-rest-beta`,
										},
										"x509_certificate_authentication_default_mode": schema.StringAttribute{
											Optional: true,
											Validators: []validator.String{
												stringvalidator.OneOf("x509CertificateSingleFactor", "x509CertificateMultiFactor", "unknownFutureValue"),
											},
											PlanModifiers: []planmodifier.String{
												wpdefaultvalue.StringDefaultValue("x509CertificateSingleFactor"),
											},
											Computed:    true,
											Description: `x509CertificateAuthenticationDefaultMode`, // custom MS Graph attribute name
										},
										"x509_certificate_default_required_affinity_level": schema.StringAttribute{
											Optional: true,
											Validators: []validator.String{
												stringvalidator.OneOf("low", "high", "unknownFutureValue"),
											},
											PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("low")},
											Computed:      true,
											Description:   `x509CertificateDefaultRequiredAffinityLevel`, // custom MS Graph attribute name
										},
									},
									PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/x509CertificateAuthenticationModeConfiguration?view=graph-rest-beta`,
								},
								"certificate_user_bindings": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{ // x509CertificateUserBinding
											"priority": schema.Int64Attribute{
												Required: true,
											},
											"trust_affinity_level": schema.StringAttribute{
												Optional: true,
												Validators: []validator.String{
													stringvalidator.OneOf("low", "high", "unknownFutureValue"),
												},
											},
											"user_property": schema.StringAttribute{
												Optional: true,
											},
											"x509_certificate_field": schema.StringAttribute{
												Optional:    true,
												Description: `x509CertificateField`, // custom MS Graph attribute name
											},
										},
									},
									PlanModifiers: []planmodifier.Set{
										authenticationMethodsPolicyX509CertificateUserBindingsDefault,
									},
									Computed:            true,
									MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/x509CertificateUserBinding?view=graph-rest-beta`,
								},
								"issuer_hints_configuration": schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // x509CertificateIssuerHintsConfiguration
										"state": schema.StringAttribute{
											Optional: true,
											Validators: []validator.String{
												stringvalidator.OneOf("disabled", "enabled", "unknownFutureValue"),
											},
											PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("disabled")},
											Computed:      true,
										},
									},
									PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/x509CertificateIssuerHintsConfiguration?view=graph-rest-beta`,
								},
								"include_targets": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: authenticationMethodsPolicyAuthenticationMethodTargetAttributes,
									},
									PlanModifiers:       []planmodifier.Set{authenticationMethodsPolicyTargetsDefaultAllUsers},
									Computed:            true,
									MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/authenticationMethodTarget?view=graph-rest-beta`,
								},
							},
							Validators: []validator.Object{
								authenticationMethodsPolicyAuthenticationMethodConfigurationValidator,
							},
							MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/x509CertificateAuthenticationMethodConfiguration?view=graph-rest-beta`,
						},
					},
				},
			},
			MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/authenticationMethodConfiguration?view=graph-rest-beta`,
		},
	},
	MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/authenticationMethodsPolicy?view=graph-rest-beta`,
}

var authenticationMethodsPolicyExcludeTargetAttributes = map[string]schema.Attribute{ // excludeTarget
	"id": schema.StringAttribute{
		Required: true,
	},
	"target_type": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("user", "group", "unknownFutureValue"),
		},
	},
}

var authenticationMethodsPolicyIncludeTargetAttributes = map[string]schema.Attribute{ // includeTarget
	"id": schema.StringAttribute{
		Required: true,
	},
	"target_type": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("user", "group", "unknownFutureValue"),
		},
	},
}

var authenticationMethodsPolicyAuthenticationMethodConfigurationValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("email"),
	path.MatchRelative().AtParent().AtName("fido2"),
	path.MatchRelative().AtParent().AtName("hardware_oath"),
	path.MatchRelative().AtParent().AtName("microsoft_authenticator"),
	path.MatchRelative().AtParent().AtName("sms"),
	path.MatchRelative().AtParent().AtName("software_oath"),
	path.MatchRelative().AtParent().AtName("temporary_access_pass"),
	path.MatchRelative().AtParent().AtName("voice"),
	path.MatchRelative().AtParent().AtName("x509_certificate"),
)

var authenticationMethodsPolicyAuthenticationMethodTargetAttributes = map[string]schema.Attribute{ // authenticationMethodTarget
	"id": schema.StringAttribute{
		Required: true,
	},
	"is_registration_required": schema.BoolAttribute{
		Optional:      true,
		PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:      true,
	},
	"target_type": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("user", "group", "unknownFutureValue"),
		},
	},
}

var authenticationMethodsPolicyPasskeyAuthenticationMethodTargetAttributes = map[string]schema.Attribute{ // passkeyAuthenticationMethodTarget
	"id": schema.StringAttribute{
		Required: true,
	},
	"is_registration_required": schema.BoolAttribute{
		Optional:      true,
		PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:      true,
	},
	"target_type": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("user", "group", "unknownFutureValue"),
		},
	},
}

var authenticationMethodsPolicyAuthenticationMethodFeatureConfigurationAttributes = map[string]schema.Attribute{ // authenticationMethodFeatureConfiguration
	"exclude_target": schema.SingleNestedAttribute{
		Optional:            true,
		Attributes:          authenticationMethodsPolicyFeatureTargetAttributes,
		PlanModifiers:       []planmodifier.Object{authenticationMethodsPolicyTargetDefaultNoGroup},
		Computed:            true,
		MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/featureTarget?view=graph-rest-beta`,
	},
	"include_target": schema.SingleNestedAttribute{
		Optional:            true,
		Attributes:          authenticationMethodsPolicyFeatureTargetAttributes,
		PlanModifiers:       []planmodifier.Object{authenticationMethodsPolicyTargetDefaultAllUsers},
		Computed:            true,
		MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/featureTarget?view=graph-rest-beta`,
	},
	"state": schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			stringvalidator.OneOf("default", "enabled", "disabled", "unknownFutureValue"),
		},
		PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("default")},
		Computed:      true,
	},
}

var authenticationMethodsPolicyFeatureTargetAttributes = map[string]schema.Attribute{ // featureTarget
	"id": schema.StringAttribute{
		Optional: true,
	},
	"target_type": schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			stringvalidator.OneOf("group", "administrativeUnit", "role", "unknownFutureValue"),
		},
	},
}

var authenticationMethodsPolicyNoGroupRaw = map[string]any{
	"target_type": "group",
	"id":          "00000000-0000-0000-0000-000000000000",
}
var authenticationMethodsPolicyTargetDefaultNoGroup = wpdefaultvalue.ObjectDefaultValue(authenticationMethodsPolicyNoGroupRaw)

var authenticationMethodsPolicyAllUsersRaw = map[string]any{
	"target_type": "group",
	"id":          "all_users",
}
var authenticationMethodsPolicyTargetDefaultAllUsers = wpdefaultvalue.ObjectDefaultValue(authenticationMethodsPolicyAllUsersRaw)
var authenticationMethodsPolicyTargetsDefaultAllUsers = wpdefaultvalue.SetDefaultValue([]any{authenticationMethodsPolicyAllUsersRaw})

var authenticationMethodsPolicyX509CertificateUserBindingsDefault = wpdefaultvalue.SetDefaultValue([]any{
	map[string]any{
		"priority":               1,
		"trust_affinity_level":   "low",
		"user_property":          "userPrincipalName",
		"x509_certificate_field": "PrincipalName",
	},
	map[string]any{
		"priority":               2,
		"trust_affinity_level":   "low",
		"user_property":          "userPrincipalName",
		"x509_certificate_field": "RFC822Name",
	},
	map[string]any{
		"priority":               3,
		"trust_affinity_level":   "high",
		"user_property":          "certificateUserIds",
		"x509_certificate_field": "SubjectKeyIdentifier",
	},
})
