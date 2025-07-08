package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	AuthenticationCombinationConfigurationResource = generic.GenericResource{
		TypeNameSuffix: "authentication_combination_configuration",
		SpecificSchema: authenticationCombinationConfigurationResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/policies/authenticationStrengthPolicies",
			ParentEntities: generic.ParentEntities{
				{
					ParentIdField: path.Root("authentication_strength_policy_id"),
					UriSuffix:     "combinationConfigurations",
				},
			},
		},
	}

	AuthenticationCombinationConfigurationSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&AuthenticationCombinationConfigurationResource)

	AuthenticationCombinationConfigurationPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&AuthenticationCombinationConfigurationResource, "")
)

var authenticationCombinationConfigurationResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // authenticationCombinationConfiguration
		"authentication_strength_policy_id": schema.StringAttribute{
			Required:            true,
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			MarkdownDescription: "Provider Note: ID of the authentication strength policy that this resource belongs to. Required.",
		},
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "A unique system-generated identifier.",
		},
		"applies_to_combinations": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Validators: []validator.Set{
				setvalidator.ValueStringsAre(
					wpvalidator.FlagEnumValues("password", "voice", "hardwareOath", "softwareOath", "sms", "fido2", "windowsHelloForBusiness", "microsoftAuthenticatorPush", "deviceBasedPush", "temporaryAccessPassOneTime", "temporaryAccessPassMultiUse", "email", "x509CertificateSingleFactor", "x509CertificateMultiFactor", "federatedSingleFactor", "federatedMultiFactor", "unknownFutureValue", "qrCodePin"),
				),
			},
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Which authentication method combinations this configuration applies to. Must be an **allowedCombinations** object defined for the [authenticationStrengthPolicy](../resources/authenticationstrengthpolicy.md). For **fido2combinationConfigurations** use `\"fido2\"`, for **x509certificatecombinationconfiguration** use `\"x509CertificateSingleFactor\"` or `\"x509CertificateMultiFactor\"`. / Possible values are: `password`, `voice`, `hardwareOath`, `softwareOath`, `sms`, `fido2`, `windowsHelloForBusiness`, `microsoftAuthenticatorPush`, `deviceBasedPush`, `temporaryAccessPassOneTime`, `temporaryAccessPassMultiUse`, `email`, `x509CertificateSingleFactor`, `x509CertificateMultiFactor`, `federatedSingleFactor`, `federatedMultiFactor`, `unknownFutureValue`, `qrCodePin`. The _provider_ default value is `[]`.",
		},
		"fido2": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.fido2CombinationConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // fido2CombinationConfiguration
					"allowed_aaguids": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						Description:         `allowedAAGUIDs`, // custom MS Graph attribute name
						MarkdownDescription: "A list of AAGUIDs allowed to be used as part of the specified authentication method combinations. The _provider_ default value is `[]`.",
					},
				},
				Validators: []validator.Object{
					authenticationCombinationConfigurationAuthenticationCombinationConfigurationValidator,
				},
				MarkdownDescription: "Configuration to require specific FIDO2 key types in an authentication strength. An administrator may use this entity to specify which Authenticator Attestations GUIDs (AAGUIDs) are allowed, as part of certain authentication method combinations, in an [authentication strength](authenticationstrengthpolicy.md).\n\nInherits and derived from [authenticationCombinationConfiguration](../resources/authenticationcombinationconfiguration.md). / https://learn.microsoft.com/en-us/graph/api/resources/fido2combinationconfiguration?view=graph-rest-beta",
			},
		},
		"x509_certificate": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.x509CertificateCombinationConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // x509CertificateCombinationConfiguration
					"allowed_issuer_skis": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "A list of allowed subject key identifier values. The _provider_ default value is `[]`.",
					},
					"allowed_policy_oids": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						Description:         `allowedPolicyOIDs`, // custom MS Graph attribute name
						MarkdownDescription: "A list of allowed policy OIDs. The _provider_ default value is `[]`.",
					},
				},
				Validators: []validator.Object{
					authenticationCombinationConfigurationAuthenticationCombinationConfigurationValidator,
				},
				MarkdownDescription: "Configuration to require specific certificate properties. You can use this entity to specify the certificate issuer or policy OID that are allowed, as part of certificate-based authentication, in an [authentication strength policy](authenticationstrengthpolicy.md). / https://learn.microsoft.com/en-us/graph/api/resources/x509certificatecombinationconfiguration?view=graph-rest-beta",
			},
		},
	},
	MarkdownDescription: "Sets restrictions on specific types, modes, or versions of an authentication method that is tied to specific auth method combinations used in an [authentication strength](authenticationstrengths-overview.md). The following resources inherit from this abstract and define the various types of combination configurations: / https://learn.microsoft.com/en-us/graph/api/resources/authenticationcombinationconfiguration?view=graph-rest-beta\n\nProvider Note: To import this resource, an ID consisting of `authentication_strength_policy_id` and `id` being joined by a forward slash (`/`) must be used. ||| MS Graph: Authentication",
}

var authenticationCombinationConfigurationAuthenticationCombinationConfigurationValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("fido2"),
	path.MatchRelative().AtParent().AtName("x509_certificate"),
)
