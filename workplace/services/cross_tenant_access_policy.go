package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	CrossTenantAccessPolicyResource = generic.GenericResource{
		TypeNameSuffix: "cross_tenant_access_policy",
		SpecificSchema: crossTenantAccessPolicyResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:     "/policies/crossTenantAccessPolicy",
			IsSingleton: true,
		},
	}

	CrossTenantAccessPolicySingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&CrossTenantAccessPolicyResource)
)

var crossTenantAccessPolicyResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // crossTenantAccessPolicy
		"id": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"display_name": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The display name of the cross-tenant access policy.",
		},
		"allowed_cloud_endpoints": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Used to specify which Microsoft clouds an organization would like to collaborate with. By default, this value is empty. Supported values for this field are: `microsoftonline.com`, `microsoftonline.us`, and `partner.microsoftonline.cn`. The _provider_ default value is `[]`.",
		},
	},
	MarkdownDescription: "Represents the base policy in the directory for cross-tenant access settings. / https://learn.microsoft.com/en-us/graph/api/resources/crosstenantaccesspolicy?view=graph-rest-beta ||| MS Graph: Cross-tenant access",
}
