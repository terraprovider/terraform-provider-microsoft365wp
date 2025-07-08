package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	AuthenticationStrengthPolicyResource = generic.GenericResource{
		TypeNameSuffix: "authentication_strength_policy",
		SpecificSchema: authenticationStrengthPolicyResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/policies/authenticationStrengthPolicies",
			WriteOptions: generic.WriteOptions{
				SubActions: []generic.WriteSubAction{
					&generic.WriteSubActionAllInOne{
						WriteSubActionBase: generic.WriteSubActionBase{
							Attributes: []string{"allowedCombinations"},
							UriSuffix:  "updateAllowedCombinations",
							UpdateOnly: true,
						},
					},
				},
			},
		},
	}

	AuthenticationStrengthPolicySingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&AuthenticationStrengthPolicyResource)

	AuthenticationStrengthPolicyPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&AuthenticationStrengthPolicyResource, "")
)

var authenticationStrengthPolicyResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // authenticationStrengthPolicy
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The system-generated identifier for this mode.",
		},
		"allowed_combinations": schema.SetAttribute{
			ElementType: types.StringType,
			Required:    true,
			Validators: []validator.Set{
				setvalidator.ValueStringsAre(
					wpvalidator.FlagEnumValues("password", "voice", "hardwareOath", "softwareOath", "sms", "fido2", "windowsHelloForBusiness", "microsoftAuthenticatorPush", "deviceBasedPush", "temporaryAccessPassOneTime", "temporaryAccessPassMultiUse", "email", "x509CertificateSingleFactor", "x509CertificateMultiFactor", "federatedSingleFactor", "federatedMultiFactor", "unknownFutureValue", "qrCodePin"),
				),
			},
			MarkdownDescription: "A collection of authentication method modes that are required be used to satify this authentication strength. / Possible values are: `password`, `voice`, `hardwareOath`, `softwareOath`, `sms`, `fido2`, `windowsHelloForBusiness`, `microsoftAuthenticatorPush`, `deviceBasedPush`, `temporaryAccessPassOneTime`, `temporaryAccessPassMultiUse`, `email`, `x509CertificateSingleFactor`, `x509CertificateMultiFactor`, `federatedSingleFactor`, `federatedMultiFactor`, `unknownFutureValue`, `qrCodePin`",
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The datetime when this policy was created.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "The human-readable description of this policy. The _provider_ default value is `\"\"`.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
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
		"combination_configurations": schema.SetNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // authenticationCombinationConfiguration
					"id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "A unique system-generated identifier.",
					},
				},
			},
			PlanModifiers:       []planmodifier.Set{wpplanmodifier.SetUseStateForUnknown()},
			MarkdownDescription: "Settings that may be used to require specific types or instances of an authentication method to be used when authenticating with a specified combination of authentication methods. / Sets restrictions on specific types, modes, or versions of an authentication method that is tied to specific auth method combinations used in an [authentication strength](authenticationstrengths-overview.md). The following resources inherit from this abstract and define the various types of combination configurations: / https://learn.microsoft.com/en-us/graph/api/resources/authenticationcombinationconfiguration?view=graph-rest-beta",
		},
	},
	MarkdownDescription: "A collection of settings that define specific combinations of authentication methods and metadata. The authentication strength policy, when applied to a given scenario using Microsoft Entra Conditional Access, defines which authentication methods must be used to authenticate in that scenario. An authentication strength may be built-in or custom (defined by the tenant) and may or may not fulfill the requirements to grant an MFA claim. / https://learn.microsoft.com/en-us/graph/api/resources/authenticationstrengthpolicy?view=graph-rest-beta\n\nProvider Note: Use the separate resource `microsoft365wp_authentication_combination_configuration` to create combination configurations for an authentication strength policy. ||| MS Graph: Authentication",
}
