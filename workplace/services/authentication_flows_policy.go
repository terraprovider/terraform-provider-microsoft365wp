package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var (
	AuthenticationFlowsPolicyResource = generic.GenericResource{
		TypeNameSuffix: "authentication_flows_policy",
		SpecificSchema: authenticationFlowsPolicyResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:     "/policies/authenticationFlowsPolicy",
			IsSingleton: true,
		},
	}

	AuthenticationFlowsPolicySingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&AuthenticationFlowsPolicyResource)
)

var authenticationFlowsPolicyResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // authenticationFlowsPolicy
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Inherited property. The ID of the authentication flows policy. Optional. Read-only.",
		},
		"description": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Inherited property. A description of the policy. This property isn't a key. Optional. Read-only.",
		},
		"display_name": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Inherited property. The human-readable name of the policy. This property isn't a key. Optional. Read-only.",
		},
		"self_service_sign_up": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // selfServiceSignUpAuthenticationFlowConfiguration
				"is_enabled": schema.BoolAttribute{
					Optional:            true,
					Default:             booldefault.StaticBool(false),
					Computed:            true,
					MarkdownDescription: "Indicates whether self-service sign-up flow is enabled or disabled. The default value is `false`. This property isn't a key. Required. <br/> The _provider_ default value is `false`.",
				},
			},
			Default: wpdefaultvalue.ObjectDefaultValue(map[string]any{
				"is_enabled": false,
			}),
			Computed:            true,
			MarkdownDescription: "Contains [selfServiceSignUpAuthenticationFlowConfiguration](https://learn.microsoft.com/en-us/graph/api/resources/selfservicesignupauthenticationflowconfiguration?view=graph-rest-beta) settings that convey whether self-service sign-up is enabled or disabled. This property isn't a key. Optional. Read-only. / Represents the configurations related to self-service sign-up. Also see [Microsoft docs for selfServiceSignUpAuthenticationFlowConfiguration](https://learn.microsoft.com/en-us/graph/api/resources/selfservicesignupauthenticationflowconfiguration?view=graph-rest-beta). <br/> The _provider_ default value is `{\"is_enabled\":false}`. <br> ",
		},
	},
	MarkdownDescription: "Represents the [policy configuration of self-service sign-up experience](https://learn.microsoft.com/en-us/graph/api/resources/selfservicesignupauthenticationflowconfiguration?view=graph-rest-beta) at a tenant level that lets external users request to sign up for approval. It contains information about the ID, display name, and description, and indicates whether self-service sign up is enabled for the policy. <br/> Also see [Microsoft docs for authenticationFlowsPolicy](https://learn.microsoft.com/en-us/graph/api/resources/authenticationflowspolicy?view=graph-rest-beta). ||| MS Graph: Policies",
}
