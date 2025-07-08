package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var (
	ExternalIdentitiesPolicyResource = generic.GenericResource{
		TypeNameSuffix: "external_identities_policy",
		SpecificSchema: externalIdentitiesPolicyResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:     "/policies/externalIdentitiesPolicy",
			IsSingleton: true,
		},
	}

	ExternalIdentitiesPolicySingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&ExternalIdentitiesPolicyResource)
)

var externalIdentitiesPolicyResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // externalIdentitiesPolicy
		"id": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"description": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"display_name": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The policy name.",
		},
		"allow_external_identities_to_leave": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:            true,
			MarkdownDescription: "Defines whether external users can leave the guest tenant. If set to `false`, self-service controls are disabled, and the admin of the guest tenant must manually remove the external user from the guest tenant. When the external user leaves the tenant, their data in the guest tenant is first soft-deleted then permanently deleted in 30 days. The _provider_ default value is `true`.",
		},
	},
	MarkdownDescription: "Represents the tenant-wide policy that controls whether external users can leave the guest Microsoft Entra tenant via self-service controls. When permitted by the administrator, external users can leave the guest Microsoft Entra tenant through the **organizations** menu of the [My Account](https://myaccount.microsoft.com/) portal. / https://learn.microsoft.com/en-us/graph/api/resources/externalidentitiespolicy?view=graph-rest-beta ||| MS Graph: Identity and sign-in",
}
