package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var (
	IdentitySecurityDefaultsEnforcementPolicyResource = generic.GenericResource{
		TypeNameSuffix: "identity_security_defaults_enforcement_policy",
		SpecificSchema: identitySecurityDefaultsEnforcementPolicyResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:     "/policies/identitySecurityDefaultsEnforcementPolicy",
			IsSingleton: true,
		},
	}

	IdentitySecurityDefaultsEnforcementPolicySingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&IdentitySecurityDefaultsEnforcementPolicyResource)
)

var identitySecurityDefaultsEnforcementPolicyResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // identitySecurityDefaultsEnforcementPolicy
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Identifier for this policy. Read-only.",
		},
		"description": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Description for this policy. Read-only.",
		},
		"display_name": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Display name for this policy. Read-only.",
		},
		"is_enabled": schema.BoolAttribute{
			Required:            true,
			MarkdownDescription: "If set to `true`, Microsoft Entra security defaults are enabled for the tenant.",
		},
	},
	MarkdownDescription: "Represents the Microsoft Entra ID [security defaults](/azure/active-directory/fundamentals/concept-fundamentals-security-defaults) policy. Security defaults contain preconfigured security settings that protect against common attacks. / https://learn.microsoft.com/en-us/graph/api/resources/identitysecuritydefaultsenforcementpolicy?view=graph-rest-beta ||| MS Graph: Policies",
}
