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
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The identifier of the policy.",
		},
		"description": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "A description of the policy.",
		},
		"display_name": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The name of the policy.",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The date and time of the last update to the policy.",
		},
		"policy_migration_state": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("preMigration", "migrationInProgress", "migrationComplete", "unknownFutureValue"),
			},
			PlanModifiers: []planmodifier.String{
				wpdefaultvalue.StringDefaultValue("migrationInProgress"),
			},
			Computed:            true,
			MarkdownDescription: "The state of migration of the authentication methods policy from the legacy multifactor authentication and self-service password reset (SSPR) policies. The possible values are: <br/><li>`premigration` - means the authentication methods policy is used for authentication only, legacy policies are respected. <li>`migrationInProgress` - means the authentication methods policy is used for both authentication and SSPR, legacy policies are respected. <li>`migrationComplete` - means the authentication methods policy is used for authentication and SSPR, legacy policies are ignored. <li>`unknownFutureValue` - Evolvable enumeration sentinel value. Don't use. / Possible values are: `preMigration`, `migrationInProgress`, `migrationComplete`, `unknownFutureValue`. The _provider_ default value is `\"migrationInProgress\"`.",
		},
		"policy_version": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The version of the policy in use.",
		},
		"registration_enforcement": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // registrationEnforcement
				"authentication_methods_registration_campaign": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // authenticationMethodsRegistrationCampaign
						"enforce_registration_after_allowed_snoozes": schema.BoolAttribute{
							Optional:            true,
							PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
							Computed:            true,
							MarkdownDescription: "Specifies whether a user is required to perform registration after snoozing 3 times. If `true`, the user is required to register after 3 snoozes. If `false`, the user can snooze indefinitely. The default value is `true`. The _provider_ default value is `true`.",
						},
						"exclude_targets": schema.SetNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: authenticationMethodsPolicyExcludeTargetAttributes,
							},
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "Users and groups of users that are excluded from being prompted to set up the authentication method. / Represents the users or groups of users that are excluded from a policy. / https://learn.microsoft.com/en-us/graph/api/resources/excludetarget?view=graph-rest-beta. The _provider_ default value is `[]`.",
						},
						"include_targets": schema.SetNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{ // authenticationMethodsRegistrationCampaignIncludeTarget
									"id": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "The object identifier of a Microsoft Entra user or group.",
									},
									"targeted_authentication_method": schema.StringAttribute{
										Optional: true,
										PlanModifiers: []planmodifier.String{
											wpdefaultvalue.StringDefaultValue("microsoftAuthenticator"),
										},
										Computed:            true,
										MarkdownDescription: "The authentication method that the user is prompted to register. The value must be `microsoftAuthenticator`. The _provider_ default value is `\"microsoftAuthenticator\"`.",
									},
									"target_type": schema.StringAttribute{
										Required: true,
										Validators: []validator.String{
											stringvalidator.OneOf("user", "group", "unknownFutureValue"),
										},
										MarkdownDescription: "The type of the authentication method target. / Possible values are: `user`, `group`, `unknownFutureValue`",
									},
								},
							},
							PlanModifiers:       []planmodifier.Set{authenticationMethodsPolicyTargetsDefaultAllUsers},
							Computed:            true,
							MarkdownDescription: "Users and groups of users that are prompted to set up the authentication method. / Represents the users and groups that are targeted for authentication method registration campaigns. Only users and groups that are enabled by the policy to set up the authentication method are targeted. / https://learn.microsoft.com/en-us/graph/api/resources/authenticationmethodsregistrationcampaignincludetarget?view=graph-rest-beta. The _provider_ default value is `authenticationMethodsPolicyTargetsDefaultAllUsers`.",
						},
						"snooze_duration_in_days": schema.Int64Attribute{
							Optional:            true,
							PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(1)},
							Computed:            true,
							MarkdownDescription: "Specifies the number of days that the user sees a prompt again if they select \"Not now\" and snoozes the prompt. Minimum 0 days. Maximum: 14 days. If the value is `0` – The user is prompted during every MFA attempt. The _provider_ default value is `1`.",
						},
						"state": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("default", "enabled", "disabled", "unknownFutureValue"),
							},
							PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("default")},
							Computed:            true,
							MarkdownDescription: "Enable or disable the feature. The `default` value is used when the configuration hasn't been explicitly set and uses the default behavior of Microsoft Entra ID for the setting. The default value is `disabled`. / Possible values are: `default`, `enabled`, `disabled`, `unknownFutureValue`. The _provider_ default value is `\"default\"`.",
						},
					},
					PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "Run campaigns to remind users to set up targeted authentication methods. / Represents the settings used to run campaigns to push users to set up targeted authentication methods. Users are prompted to set up the authentication method after they successfully complete a MFA challenge. Only available for the Microsoft Authenticator app for MFA. / https://learn.microsoft.com/en-us/graph/api/resources/authenticationmethodsregistrationcampaign?view=graph-rest-beta. The _provider_ default value is `{}`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Enforce registration at sign-in time. This property can be used to remind users to set up targeted authentication methods. / Enforce registration at sign-in time. This can currently only be used to remind users to set up targeted authentication methods (Microsoft Authenticator) using the 'authenticationMethodsRegistrationCampaign`. / https://learn.microsoft.com/en-us/graph/api/resources/registrationenforcement?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
		"report_suspicious_activity_settings": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // reportSuspiciousActivitySettings
				"include_target": schema.SingleNestedAttribute{
					Optional:            true,
					Attributes:          authenticationMethodsPolicyIncludeTargetAttributes,
					PlanModifiers:       []planmodifier.Object{authenticationMethodsPolicyTargetDefaultAllUsers},
					Computed:            true,
					MarkdownDescription: "Group IDs in scope for report suspicious activity. / Defines the users and groups that are included in a set of changes. / https://learn.microsoft.com/en-us/graph/api/resources/includetarget?view=graph-rest-beta. The _provider_ default value is `authenticationMethodsPolicyTargetDefaultAllUsers`.",
				},
				"state": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf("default", "enabled", "disabled", "unknownFutureValue"),
					},
					PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("default")},
					Computed:            true,
					MarkdownDescription: "Specifies the state of the reportSuspiciousActivitySettings object. The Setting to `default` results in a disabled state. / Possible values are: `default`, `enabled`, `disabled`, `unknownFutureValue`. The _provider_ default value is `\"default\"`.",
				},
				"voice_reporting_code": schema.Int64Attribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
					Computed:            true,
					MarkdownDescription: "Specifies the number the user enters on their phone to report the MFA prompt as suspicious. The _provider_ default value is `0`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Enable users to report unexpected voice call or phone app notification multi-factor authentication prompts as suspicious. / Defines the report suspicious activity settings for the tenant, whether it's enabled and which group of users is enabled for use. Report suspicious activity enables users to report a suspicious voice or phone app notification multifactor authentication prompt as suspicious. These users have their user risk set to `high`, and a [risk detection](riskdetection.md) **riskEventType** of `userReportedSuspiciousActivity` is emitted. / https://learn.microsoft.com/en-us/graph/api/resources/reportsuspiciousactivitysettings?view=graph-rest-beta. The _provider_ default value is `{}`.",
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
					MarkdownDescription: "Users and groups excluded from the preferred authentication method experience of the system. / Represents the users or groups of users that are excluded from a policy. / https://learn.microsoft.com/en-us/graph/api/resources/excludetarget?view=graph-rest-beta. The _provider_ default value is `[]`.",
				},
				"include_targets": schema.SetNestedAttribute{
					Optional: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: authenticationMethodsPolicyIncludeTargetAttributes,
					},
					PlanModifiers:       []planmodifier.Set{authenticationMethodsPolicyTargetsDefaultAllUsers},
					Computed:            true,
					MarkdownDescription: "Users and groups included in the preferred authentication method experience of the system. / Defines the users and groups that are included in a set of changes. / https://learn.microsoft.com/en-us/graph/api/resources/includetarget?view=graph-rest-beta. The _provider_ default value is `authenticationMethodsPolicyTargetsDefaultAllUsers`.",
				},
				"state": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf("default", "enabled", "disabled", "unknownFutureValue"),
					},
					PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("default")},
					Computed:            true,
					MarkdownDescription: "Indicates whether the feature is enabled or disabled. The `default` value is used when the configuration hasn't been explicitly set, and uses the default behavior of Microsoft Entra ID for the setting. The default value is `disabled`. / Possible values are: `default`, `enabled`, `disabled`, `unknownFutureValue`. The _provider_ default value is `\"default\"`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Prompt users with their most-preferred credential for multifactor authentication. / Dynamically detects and prompts users with their preferred multifactor authentication method from the registered methods. / https://learn.microsoft.com/en-us/graph/api/resources/systemcredentialpreferences?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
		"authentication_method_configurations": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // authenticationMethodConfiguration
					"id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The policy name.",
					},
					"exclude_targets": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: authenticationMethodsPolicyExcludeTargetAttributes,
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "Groups of users that are excluded from a policy. / Represents the users or groups of users that are excluded from a policy. / https://learn.microsoft.com/en-us/graph/api/resources/excludetarget?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"state": schema.StringAttribute{
						Required:            true,
						Validators:          []validator.String{stringvalidator.OneOf("enabled", "disabled")},
						MarkdownDescription: "The state of the policy. / Possible values are: `enabled`, `disabled`",
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
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("default")},
									Computed:            true,
									MarkdownDescription: "Determines whether email OTP is usable by external users for authentication. Tenants in the `default` state who didn't use the *beta* API automatically have email OTP enabled beginning in October 2021. / Possible values are: `default`, `enabled`, `disabled`, `unknownFutureValue`. The _provider_ default value is `\"default\"`.",
								},
								"include_targets": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: authenticationMethodsPolicyAuthenticationMethodTargetAttributes,
									},
									PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: "A collection of groups that are enabled to use the authentication method. / A collection of groups that are enabled to use an authentication method as part of an authentication method policy in Microsoft Entra ID.\n\nThe following types are derived from this resource type: / https://learn.microsoft.com/en-us/graph/api/resources/authenticationmethodtarget?view=graph-rest-beta. The _provider_ default value is `[]`.",
								},
							},
							Validators: []validator.Object{
								authenticationMethodsPolicyAuthenticationMethodConfigurationValidator,
							},
							MarkdownDescription: "Represents this tenant's email one-time passcode (OTP) authentication methods policy. Authentication methods policies define configuration settings and users or groups who are enabled to use the authentication method. The tenant's cloud-native users may use email OTP for self-service password reset. External users can use email OTP for authentication during invitation redemption and self-service sign-up for specific apps in user flows. / https://learn.microsoft.com/en-us/graph/api/resources/emailauthenticationmethodconfiguration?view=graph-rest-beta",
						},
					},
					"fido2": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.fido2AuthenticationMethodConfiguration",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // fido2AuthenticationMethodConfiguration
								"is_attestation_enforced": schema.BoolAttribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
									Computed:            true,
									MarkdownDescription: "Determines whether attestation must be enforced for FIDO2 security key registration. The _provider_ default value is `false`.",
								},
								"is_self_service_registration_allowed": schema.BoolAttribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
									Computed:            true,
									MarkdownDescription: "Determines if users can register new FIDO2 security keys. The _provider_ default value is `true`.",
								},
								"key_restrictions": schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // fido2KeyRestrictions
										"aa_guids": schema.SetAttribute{
											ElementType:         types.StringType,
											Optional:            true,
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: "A collection of Authenticator Attestation GUIDs. AADGUIDs define key types and manufacturers. The _provider_ default value is `[]`.",
										},
										"enforcement_type": schema.StringAttribute{
											Optional: true,
											Validators: []validator.String{
												stringvalidator.OneOf("allow", "block", "unknownFutureValue"),
											},
											PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("block")},
											Computed:            true,
											MarkdownDescription: "Enforcement type. / Possible values are: `allow`, `block`, `unknownFutureValue`. The _provider_ default value is `\"block\"`.",
										},
										"is_enforced": schema.BoolAttribute{
											Optional:            true,
											PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
											Computed:            true,
											MarkdownDescription: "Determines if the configured key enforcement is enabled. The _provider_ default value is `false`.",
										},
									},
									PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: "Controls whether key restrictions are enforced on FIDO2 security keys, either allowing or disallowing certain key types as defined by Authenticator Attestation GUID (AAGUID), an identifier that indicates the type (e.g. make and model) of the authenticator. / Represents the key restrictions that are enforced as part of the [FIDO2 security keys authentication methods policy](../resources/fido2authenticationmethodconfiguration.md). / https://learn.microsoft.com/en-us/graph/api/resources/fido2keyrestrictions?view=graph-rest-beta. The _provider_ default value is `{}`.",
								},
								"include_targets": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: authenticationMethodsPolicyPasskeyAuthenticationMethodTargetAttributes,
									},
									PlanModifiers:       []planmodifier.Set{authenticationMethodsPolicyTargetsDefaultAllUsers},
									Computed:            true,
									MarkdownDescription: "A collection of groups that are enabled to use the authentication method. / A collection of groups that are enabled to use a passkey (FIDO2) authentication method as part of a passkey (FIDO2) authentication method policy in Microsoft Entra ID. / https://learn.microsoft.com/en-us/graph/api/resources/passkeyauthenticationmethodtarget?view=graph-rest-beta. The _provider_ default value is `authenticationMethodsPolicyTargetsDefaultAllUsers`.",
								},
							},
							Validators: []validator.Object{
								authenticationMethodsPolicyAuthenticationMethodConfigurationValidator,
							},
							MarkdownDescription: "Represents a FIDO2 authentication methods policy. Authentication methods policies define configuration settings and users or groups who are enabled to use the authentication method. / https://learn.microsoft.com/en-us/graph/api/resources/fido2authenticationmethodconfiguration?view=graph-rest-beta",
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
									MarkdownDescription: "A collection of groups that are enabled to use the authentication method. Expanded by default. / A collection of groups that are enabled to use an authentication method as part of an authentication method policy in Microsoft Entra ID.\n\nThe following types are derived from this resource type: / https://learn.microsoft.com/en-us/graph/api/resources/authenticationmethodtarget?view=graph-rest-beta. The _provider_ default value is `authenticationMethodsPolicyTargetsDefaultAllUsers`.",
								},
							},
							Validators: []validator.Object{
								authenticationMethodsPolicyAuthenticationMethodConfigurationValidator,
							},
							MarkdownDescription: "Represents a Hardware OATH authentication method policy. Authentication method policies define configuration settings and users or groups that are enabled to use the authentication method. / https://learn.microsoft.com/en-us/graph/api/resources/hardwareoathauthenticationmethodconfiguration?view=graph-rest-beta",
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
											MarkdownDescription: "Determines whether users are able to approve push notifications on other Microsoft applications such as Outlook Mobile. / Defines the features that are allowed for different authentication methods. For each authentication method, defines the users who are enabled to use or excluded from using the feature. / https://learn.microsoft.com/en-us/graph/api/resources/authenticationmethodfeatureconfiguration?view=graph-rest-beta. The _provider_ default value is `{}`.",
										},
										"display_app_information_required_state": schema.SingleNestedAttribute{
											Optional:            true,
											Attributes:          authenticationMethodsPolicyAuthenticationMethodFeatureConfigurationAttributes,
											PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: "Determines whether the user's Authenticator app shows them the client app they're signing into. / Defines the features that are allowed for different authentication methods. For each authentication method, defines the users who are enabled to use or excluded from using the feature. / https://learn.microsoft.com/en-us/graph/api/resources/authenticationmethodfeatureconfiguration?view=graph-rest-beta. The _provider_ default value is `{}`.",
										},
										"display_location_information_required_state": schema.SingleNestedAttribute{
											Optional:            true,
											Attributes:          authenticationMethodsPolicyAuthenticationMethodFeatureConfigurationAttributes,
											PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: "Determines whether the user's Authenticator app shows them the geographic location of where the authentication request originated from. / Defines the features that are allowed for different authentication methods. For each authentication method, defines the users who are enabled to use or excluded from using the feature. / https://learn.microsoft.com/en-us/graph/api/resources/authenticationmethodfeatureconfiguration?view=graph-rest-beta. The _provider_ default value is `{}`.",
										},
									},
									PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: "A collection of Microsoft Authenticator settings such as number matching and location context, and whether they are enabled for all users or specific users only. / Represents Microsoft Authenticator settings such as number matching and location context, and whether they're enabled for all users or specific users only. / https://learn.microsoft.com/en-us/graph/api/resources/microsoftauthenticatorfeaturesettings?view=graph-rest-beta. The _provider_ default value is `{}`.",
								},
								"is_software_oath_enabled": schema.BoolAttribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
									Computed:            true,
									MarkdownDescription: "`true` if users can use the OTP code generated by the Microsoft Authenticator app, `false` otherwise. The _provider_ default value is `false`.",
								},
								"include_targets": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{ // microsoftAuthenticatorAuthenticationMethodTarget
											"id": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "Object identifier of a Microsoft Entra user or group.",
											},
											"is_registration_required": schema.BoolAttribute{
												Optional:            true,
												PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
												Computed:            true,
												MarkdownDescription: "Determines whether the user is enforced to register the authentication method. **Not supported**. The _provider_ default value is `false`.",
											},
											"target_type": schema.StringAttribute{
												Required: true,
												Validators: []validator.String{
													stringvalidator.OneOf("user", "group", "unknownFutureValue"),
												},
												MarkdownDescription: "and `unknownFutureValue`. From December 2022, targeting individual users using `user` is no longer recommended. Existing targets remain but we recommend moving the individual users to a targeted group. / Possible values are: `user`, `group`, `unknownFutureValue`",
											},
											"authentication_mode": schema.StringAttribute{
												Optional: true,
												Validators: []validator.String{
													stringvalidator.OneOf("deviceBasedPush", "push", "any"),
												},
												PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("any")},
												Computed:            true,
												MarkdownDescription: "Determines which types of notifications can be used for sign-in. The(passwordless only), `push`, and `any`. / Possible values are: `deviceBasedPush`, `push`, `any`. The _provider_ default value is `\"any\"`.",
											},
										},
									},
									PlanModifiers:       []planmodifier.Set{authenticationMethodsPolicyTargetsDefaultAllUsers},
									Computed:            true,
									MarkdownDescription: "A collection of groups that are enabled to use the authentication method. Expanded by default. / A collection of groups enabled to use [Microsoft Authenticator authentication methods policy](../resources/microsoftAuthenticatorAuthenticationMethodConfiguration.md) in Microsoft Entra ID. / https://learn.microsoft.com/en-us/graph/api/resources/microsoftauthenticatorauthenticationmethodtarget?view=graph-rest-beta. The _provider_ default value is `authenticationMethodsPolicyTargetsDefaultAllUsers`.",
								},
							},
							Validators: []validator.Object{
								authenticationMethodsPolicyAuthenticationMethodConfigurationValidator,
							},
							MarkdownDescription: "Represents a Microsoft Authenticator authentication methods policy. Authentication methods policies define configuration settings and users or groups that are enabled to use the authentication method. / https://learn.microsoft.com/en-us/graph/api/resources/microsoftauthenticatorauthenticationmethodconfiguration?view=graph-rest-beta",
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
												Required:            true,
												MarkdownDescription: "Object ID of a Microsoft Entra user or group.",
											},
											"is_registration_required": schema.BoolAttribute{
												Optional:            true,
												PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
												Computed:            true,
												MarkdownDescription: "Determines whether the user is enforced to register the authentication method. **Not supported**. The _provider_ default value is `false`.",
											},
											"target_type": schema.StringAttribute{
												Required: true,
												Validators: []validator.String{
													stringvalidator.OneOf("user", "group", "unknownFutureValue"),
												},
												MarkdownDescription: "and `unknownFutureValue`. From December 2022, targeting individual users using `user` is no longer recommended. Existing targets remain but we recommend moving the individual users to a targeted group. / Possible values are: `user`, `group`, `unknownFutureValue`",
											},
											"is_usable_for_sign_in": schema.BoolAttribute{
												Optional:            true,
												PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
												Computed:            true,
												MarkdownDescription: "Determines if users can use this authentication method to sign in to Microsoft Entra ID. `true` if users can use this method for primary authentication, otherwise `false`. The _provider_ default value is `true`.",
											},
										},
									},
									PlanModifiers:       []planmodifier.Set{authenticationMethodsPolicyTargetsDefaultAllUsers},
									Computed:            true,
									MarkdownDescription: "A collection of groups that are enabled to use the authentication method. / A collection of groups enabled to use [text message authentication methods policy](../resources/smsAuthenticationMethodConfiguration.md) in Microsoft Entra ID. / https://learn.microsoft.com/en-us/graph/api/resources/smsauthenticationmethodtarget?view=graph-rest-beta. The _provider_ default value is `authenticationMethodsPolicyTargetsDefaultAllUsers`.",
								},
							},
							Validators: []validator.Object{
								authenticationMethodsPolicyAuthenticationMethodConfigurationValidator,
							},
							MarkdownDescription: "Represents a Text Message authentication methods policy. Authentication methods policies define configuration settings and users or groups that are enabled to use the authentication method. / https://learn.microsoft.com/en-us/graph/api/resources/smsauthenticationmethodconfiguration?view=graph-rest-beta",
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
									MarkdownDescription: "A collection of groups that are enabled to use the authentication method. Expanded by default. / A collection of groups that are enabled to use an authentication method as part of an authentication method policy in Microsoft Entra ID.\n\nThe following types are derived from this resource type: / https://learn.microsoft.com/en-us/graph/api/resources/authenticationmethodtarget?view=graph-rest-beta. The _provider_ default value is `authenticationMethodsPolicyTargetsDefaultAllUsers`.",
								},
							},
							Validators: []validator.Object{
								authenticationMethodsPolicyAuthenticationMethodConfigurationValidator,
							},
							MarkdownDescription: "Represents the authentication policy for a third-party software OATH authentication method. Authentication methods policies define configuration settings and users or groups that are enabled to use the authentication method. / https://learn.microsoft.com/en-us/graph/api/resources/softwareoathauthenticationmethodconfiguration?view=graph-rest-beta",
						},
					},
					"temporary_access_pass": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.temporaryAccessPassAuthenticationMethodConfiguration",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // temporaryAccessPassAuthenticationMethodConfiguration
								"default_length": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(8)},
									Computed:            true,
									MarkdownDescription: "Default length in characters of a Temporary Access Pass object. Must be between 8 and 48 characters. The _provider_ default value is `8`.",
								},
								"default_lifetime_in_minutes": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(60)},
									Computed:            true,
									MarkdownDescription: "Default lifetime in minutes for a Temporary Access Pass. Value can be any integer between the **minimumLifetimeInMinutes** and **maximumLifetimeInMinutes**. The _provider_ default value is `60`.",
								},
								"is_usable_once": schema.BoolAttribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
									Computed:            true,
									MarkdownDescription: "If `true`, all the passes in the tenant will be restricted to one-time use. If `false`, passes in the tenant can be created to be either one-time use or reusable. The _provider_ default value is `false`.",
								},
								"maximum_lifetime_in_minutes": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(480)},
									Computed:            true,
									MarkdownDescription: "Maximum lifetime in minutes for any Temporary Access Pass created in the tenant. Value can be between 10 and 43200 minutes (equivalent to 30 days). The _provider_ default value is `480`.",
								},
								"minimum_lifetime_in_minutes": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(60)},
									Computed:            true,
									MarkdownDescription: "Minimum lifetime in minutes for any Temporary Access Pass created in the tenant. Value can be between 10 and 43200 minutes (equivalent to 30 days). The _provider_ default value is `60`.",
								},
								"include_targets": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: authenticationMethodsPolicyAuthenticationMethodTargetAttributes,
									},
									PlanModifiers:       []planmodifier.Set{authenticationMethodsPolicyTargetsDefaultAllUsers},
									Computed:            true,
									MarkdownDescription: "A collection of groups that are enabled to use the authentication method. / A collection of groups that are enabled to use an authentication method as part of an authentication method policy in Microsoft Entra ID.\n\nThe following types are derived from this resource type: / https://learn.microsoft.com/en-us/graph/api/resources/authenticationmethodtarget?view=graph-rest-beta. The _provider_ default value is `authenticationMethodsPolicyTargetsDefaultAllUsers`.",
								},
							},
							Validators: []validator.Object{
								authenticationMethodsPolicyAuthenticationMethodConfigurationValidator,
							},
							MarkdownDescription: "Represents a Temporary Access Pass authentication methods policy that defines the configuration settings and users or groups who are enabled to use the [Temporary Access Pass authentication method](temporaryaccesspassauthenticationmethod.md). / https://learn.microsoft.com/en-us/graph/api/resources/temporaryaccesspassauthenticationmethodconfiguration?view=graph-rest-beta",
						},
					},
					"voice": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.voiceAuthenticationMethodConfiguration",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // voiceAuthenticationMethodConfiguration
								"is_office_phone_allowed": schema.BoolAttribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
									Computed:            true,
									MarkdownDescription: "`true` if users can register office phones, otherwise, `false`. The _provider_ default value is `false`.",
								},
								"include_targets": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{ // voiceAuthenticationMethodTarget
											"id": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "Object ID of a Microsoft Entra group.",
											},
											"is_registration_required": schema.BoolAttribute{
												Optional:            true,
												PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
												Computed:            true,
												MarkdownDescription: "Determines whether the user is enforced to register the authentication method. **Not supported**. The _provider_ default value is `false`.",
											},
											"target_type": schema.StringAttribute{
												Required: true,
												Validators: []validator.String{
													stringvalidator.OneOf("user", "group", "unknownFutureValue"),
												},
												MarkdownDescription: "and `unknownFutureValue`. From December 2022, targeting individual users using `user` is no longer recommended. Existing targets remain but we recommend moving the individual users to a targeted group. / Possible values are: `user`, `group`, `unknownFutureValue`",
											},
										},
									},
									PlanModifiers:       []planmodifier.Set{authenticationMethodsPolicyTargetsDefaultAllUsers},
									Computed:            true,
									MarkdownDescription: "A collection of groups that are enabled to use the authentication method. Expanded by default. / A collection of groups enabled to use voice call authentication via the [voice call authentication methods policy](../resources/voiceAuthenticationMethodConfiguration.md) in Microsoft Entra ID. / https://learn.microsoft.com/en-us/graph/api/resources/voiceauthenticationmethodtarget?view=graph-rest-beta. The _provider_ default value is `authenticationMethodsPolicyTargetsDefaultAllUsers`.",
								},
							},
							Validators: []validator.Object{
								authenticationMethodsPolicyAuthenticationMethodConfigurationValidator,
							},
							MarkdownDescription: "Represents a voice call authentication methods policy. Authentication methods policies define configuration settings and users or groups that are enabled to use the authentication method. / https://learn.microsoft.com/en-us/graph/api/resources/voiceauthenticationmethodconfiguration?view=graph-rest-beta",
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
														Optional:            true,
														MarkdownDescription: "The identifier of the X.509 certificate. Required.",
													},
													"issuer_subject_identifier": schema.StringAttribute{
														Optional:            true,
														MarkdownDescription: "The identifier of the certificate issuer.",
													},
													"policy_oid_identifier": schema.StringAttribute{
														Optional:            true,
														MarkdownDescription: "The identifier of the X.509 certificate policyOID.",
													},
													"x509_certificate_authentication_mode": schema.StringAttribute{
														Optional: true,
														Validators: []validator.String{
															stringvalidator.OneOf("x509CertificateSingleFactor", "x509CertificateMultiFactor", "unknownFutureValue"),
														},
														Description:         `x509CertificateAuthenticationMode`, // custom MS Graph attribute name
														MarkdownDescription: "The type of strong authentication mode. The Required. / Possible values are: `x509CertificateSingleFactor`, `x509CertificateMultiFactor`, `unknownFutureValue`",
													},
													"x509_certificate_required_affinity_level": schema.StringAttribute{
														Optional: true,
														Validators: []validator.String{
															stringvalidator.OneOf("low", "high", "unknownFutureValue"),
														},
														Description:         `x509CertificateRequiredAffinityLevel`, // custom MS Graph attribute name
														MarkdownDescription: "The / Possible values are: `low`, `high`, `unknownFutureValue`",
													},
													"x509_certificate_rule_type": schema.StringAttribute{
														Optional: true,
														Validators: []validator.String{
															stringvalidator.OneOf("issuerSubject", "policyOID", "unknownFutureValue", "issuerSubjectAndPolicyOID"),
														},
														Description:         `x509CertificateRuleType`, // custom MS Graph attribute name
														MarkdownDescription: "The type of the X.509 certificate mode configuration rule. The Note that you must use the `Prefer: include-unknown-enum-members` request header to get the following values from this [evolvable enum](/graph/best-practices-concept#handling-future-members-in-evolvable-enumerations): `issuerSubjectAndPolicyOID`. Required. / Possible values are: `issuerSubject`, `policyOID`, `unknownFutureValue`, `issuerSubjectAndPolicyOID`",
													},
												},
											},
											PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
											Computed:            true,
											MarkdownDescription: "Rules are configured in addition to the authentication mode to bind a specific **x509CertificateRuleType** to an **x509CertificateAuthenticationMode**. For example, bind the `policyOID` with identifier `1.32.132.343` to `x509CertificateMultiFactor` authentication mode. / Defines the strong authentication configuration rules for the X.509 certificate. Rules are configured in addition to the authentication mode. / https://learn.microsoft.com/en-us/graph/api/resources/x509certificaterule?view=graph-rest-beta. The _provider_ default value is `[]`.",
										},
										"x509_certificate_authentication_default_mode": schema.StringAttribute{
											Optional: true,
											Validators: []validator.String{
												stringvalidator.OneOf("x509CertificateSingleFactor", "x509CertificateMultiFactor", "unknownFutureValue"),
											},
											PlanModifiers: []planmodifier.String{
												wpdefaultvalue.StringDefaultValue("x509CertificateSingleFactor"),
											},
											Computed:            true,
											Description:         `x509CertificateAuthenticationDefaultMode`, // custom MS Graph attribute name
											MarkdownDescription: "The type of strong authentication mode. The / Possible values are: `x509CertificateSingleFactor`, `x509CertificateMultiFactor`, `unknownFutureValue`. The _provider_ default value is `\"x509CertificateSingleFactor\"`.",
										},
										"x509_certificate_default_required_affinity_level": schema.StringAttribute{
											Optional: true,
											Validators: []validator.String{
												stringvalidator.OneOf("low", "high", "unknownFutureValue"),
											},
											PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("low")},
											Computed:            true,
											Description:         `x509CertificateDefaultRequiredAffinityLevel`, // custom MS Graph attribute name
											MarkdownDescription: "Determines the default value for the tenant affinity binding level. The / Possible values are: `low`, `high`, `unknownFutureValue`. The _provider_ default value is `\"low\"`.",
										},
									},
									PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: "Defines strong authentication configurations. This configuration includes the default authentication mode and the different rules for strong authentication bindings. / Defines the strong authentication configurations for the X.509 certificate. This configuration includes the default authentication mode and the different rules of strong authentication bindings. / https://learn.microsoft.com/en-us/graph/api/resources/x509certificateauthenticationmodeconfiguration?view=graph-rest-beta. The _provider_ default value is `{}`.",
								},
								"certificate_user_bindings": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{ // x509CertificateUserBinding
											"priority": schema.Int64Attribute{
												Required:            true,
												MarkdownDescription: "The priority of the binding. Microsoft Entra ID uses the binding with the highest priority. This value must be a non-negative integer and unique in the collection of objects in the **certificateUserBindings** property of an **x509CertificateAuthenticationMethodConfiguration** object. Required",
											},
											"trust_affinity_level": schema.StringAttribute{
												Optional: true,
												Validators: []validator.String{
													stringvalidator.OneOf("low", "high", "unknownFutureValue"),
												},
												MarkdownDescription: "The affinity level of the username binding rule. The / Possible values are: `low`, `high`, `unknownFutureValue`",
											},
											"user_property": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "Defines the Microsoft Entra user property of the user object to use for the binding. The Required.",
											},
											"x509_certificate_field": schema.StringAttribute{
												Optional:            true,
												Description:         `x509CertificateField`, // custom MS Graph attribute name
												MarkdownDescription: "The field on the X.509 certificate to use for the binding. The",
											},
										},
									},
									PlanModifiers: []planmodifier.Set{
										authenticationMethodsPolicyX509CertificateUserBindingsDefault,
									},
									Computed:            true,
									MarkdownDescription: "Defines fields in the X.509 certificate that map to attributes of the Microsoft Entra user object in order to bind the certificate to the user. The **priority** of the object determines the order in which the binding is carried out. The first binding that matches will be used and the rest ignored. / Defines the fields in the X.509 certificate that map to attributes of the Microsoft Entra user object in order to bind the certificate to the user account. / https://learn.microsoft.com/en-us/graph/api/resources/x509certificateuserbinding?view=graph-rest-beta. The _provider_ default value is `authenticationMethodsPolicyX509CertificateUserBindingsDefault`.",
								},
								"issuer_hints_configuration": schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // x509CertificateIssuerHintsConfiguration
										"state": schema.StringAttribute{
											Optional: true,
											Validators: []validator.String{
												stringvalidator.OneOf("disabled", "enabled", "unknownFutureValue"),
											},
											PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("disabled")},
											Computed:            true,
											MarkdownDescription: "The / Possible values are: `disabled`, `enabled`, `unknownFutureValue`. The _provider_ default value is `\"disabled\"`.",
										},
									},
									PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
									Computed:            true,
									MarkdownDescription: "Determines whether issuer(CA) hints are sent back to the client side to filter the certificates shown in certificate picker. / Determines whether issuer(CA) hints are sent back to the client side to filter the certificates shown in certificate picker. / https://learn.microsoft.com/en-us/graph/api/resources/x509certificateissuerhintsconfiguration?view=graph-rest-beta. The _provider_ default value is `{}`.",
								},
								"include_targets": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: authenticationMethodsPolicyAuthenticationMethodTargetAttributes,
									},
									PlanModifiers:       []planmodifier.Set{authenticationMethodsPolicyTargetsDefaultAllUsers},
									Computed:            true,
									MarkdownDescription: "A collection of groups that are enabled to use the authentication method. / A collection of groups that are enabled to use an authentication method as part of an authentication method policy in Microsoft Entra ID.\n\nThe following types are derived from this resource type: / https://learn.microsoft.com/en-us/graph/api/resources/authenticationmethodtarget?view=graph-rest-beta. The _provider_ default value is `authenticationMethodsPolicyTargetsDefaultAllUsers`.",
								},
							},
							Validators: []validator.Object{
								authenticationMethodsPolicyAuthenticationMethodConfigurationValidator,
							},
							MarkdownDescription: "Represents the details of the Microsoft Entra native Certificate-Based Authentication (CBA) in the tenant, including whether the authentication method is enabled or disabled and the users and groups who can register and use it. / https://learn.microsoft.com/en-us/graph/api/resources/x509certificateauthenticationmethodconfiguration?view=graph-rest-beta",
						},
					},
				},
			},
			MarkdownDescription: "Represents the settings for each authentication method. Automatically expanded on `GET /policies/authenticationMethodsPolicy`. / An abstract type that represents the settings for each authentication method. It has the configuration of whether a specific authentication method is enabled or disabled for the tenant and which users and groups can register and use that method.\n\nThe following authentication methods are derived from the **authenticationMethodConfiguration** resource type: / https://learn.microsoft.com/en-us/graph/api/resources/authenticationmethodconfiguration?view=graph-rest-beta",
		},
	},
	MarkdownDescription: "Defines authentication methods and the users that are allowed to use them to sign in and perform multi-factor authentication (MFA) in Microsoft Entra ID. / https://learn.microsoft.com/en-us/graph/api/resources/authenticationmethodspolicy?view=graph-rest-beta ||| MS Graph: Authentication",
}

var authenticationMethodsPolicyExcludeTargetAttributes = map[string]schema.Attribute{ // excludeTarget
	"id": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The object identifier of a Microsoft Entra group.",
	},
	"target_type": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("user", "group", "unknownFutureValue"),
		},
		MarkdownDescription: "The type of the authentication method target.and `unknownFutureValue`. / Possible values are: `user`, `group`, `unknownFutureValue`",
	},
}

var authenticationMethodsPolicyIncludeTargetAttributes = map[string]schema.Attribute{ // includeTarget
	"id": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The ID of the entity targeted.",
	},
	"target_type": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("user", "group", "unknownFutureValue"),
		},
		MarkdownDescription: "The kind of entity targeted. / Possible values are: `user`, `group`, `unknownFutureValue`",
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
		Required:            true,
		MarkdownDescription: "Object identifier of a Microsoft Entra user or group.",
	},
	"is_registration_required": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		MarkdownDescription: "Determines if the user is enforced to register the authentication method. The _provider_ default value is `false`.",
	},
	"target_type": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("user", "group", "unknownFutureValue"),
		},
		MarkdownDescription: "and `unknownFutureValue`. From December 2022, targeting individual users using `user` is no longer recommended. Existing targets remain but we recommend moving the individual users to a targeted group. / Possible values are: `user`, `group`, `unknownFutureValue`",
	},
}

var authenticationMethodsPolicyPasskeyAuthenticationMethodTargetAttributes = map[string]schema.Attribute{ // passkeyAuthenticationMethodTarget
	"id": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "Object identifier of a Microsoft Entra user or group.",
	},
	"is_registration_required": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		MarkdownDescription: "Indicates whether the user is required to register the authentication method. The _provider_ default value is `false`.",
	},
	"target_type": schema.StringAttribute{
		Required: true,
		Validators: []validator.String{
			stringvalidator.OneOf("user", "group", "unknownFutureValue"),
		},
		MarkdownDescription: "The authentication method type. Theand `unknownFutureValue`. Effective December 2022, the `user` target value is no longer recommended. We recommend moving individual users to a targeted group. / Possible values are: `user`, `group`, `unknownFutureValue`",
	},
}

var authenticationMethodsPolicyAuthenticationMethodFeatureConfigurationAttributes = map[string]schema.Attribute{ // authenticationMethodFeatureConfiguration
	"exclude_target": schema.SingleNestedAttribute{
		Optional:            true,
		Attributes:          authenticationMethodsPolicyFeatureTargetAttributes,
		PlanModifiers:       []planmodifier.Object{authenticationMethodsPolicyTargetDefaultNoGroup},
		Computed:            true,
		MarkdownDescription: "A single entity that's excluded from using this feature. / Defines a single group, Microsoft Entra role, or administrative unit that is included or excluded in the settings specified in the [authenticationMethodFeatureConfiguration](authenticationmethodfeatureconfiguration.md) object. / https://learn.microsoft.com/en-us/graph/api/resources/featuretarget?view=graph-rest-beta. The _provider_ default value is `authenticationMethodsPolicyTargetDefaultNoGroup`.",
	},
	"include_target": schema.SingleNestedAttribute{
		Optional:            true,
		Attributes:          authenticationMethodsPolicyFeatureTargetAttributes,
		PlanModifiers:       []planmodifier.Object{authenticationMethodsPolicyTargetDefaultAllUsers},
		Computed:            true,
		MarkdownDescription: "A single entity that's allowed to use this feature. / Defines a single group, Microsoft Entra role, or administrative unit that is included or excluded in the settings specified in the [authenticationMethodFeatureConfiguration](authenticationmethodfeatureconfiguration.md) object. / https://learn.microsoft.com/en-us/graph/api/resources/featuretarget?view=graph-rest-beta. The _provider_ default value is `authenticationMethodsPolicyTargetDefaultAllUsers`.",
	},
	"state": schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			stringvalidator.OneOf("default", "enabled", "disabled", "unknownFutureValue"),
		},
		PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("default")},
		Computed:            true,
		MarkdownDescription: "Enable or disable the feature. The `default` value is used when the configuration hasn't been explicitly set and uses the default behavior of Microsoft Entra ID for the setting. The default value is `disabled`. / Possible values are: `default`, `enabled`, `disabled`, `unknownFutureValue`. The _provider_ default value is `\"default\"`.",
	},
}

var authenticationMethodsPolicyFeatureTargetAttributes = map[string]schema.Attribute{ // featureTarget
	"id": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The ID of the entity that's targeted in the include or exclude rule or `all_users` to target all users.",
	},
	"target_type": schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			stringvalidator.OneOf("group", "administrativeUnit", "role", "unknownFutureValue"),
		},
		MarkdownDescription: "The kind of entity that's targeted. The / Possible values are: `group`, `administrativeUnit`, `role`, `unknownFutureValue`",
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
