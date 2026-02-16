package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvaluemodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var (
	CrossTenantIdentitySyncPolicyPartnerResource = generic.GenericResource{
		TypeNameSuffix: "cross_tenant_identity_sync_policy_partner",
		SpecificSchema: crossTenantIdentitySyncPolicyPartnerResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/policies/crossTenantAccessPolicy/partners",
			ParentEntities: generic.ParentEntities{
				{
					ParentIdField: path.Root("tenant_id"),
					UriSuffix:     "identitySynchronization",
				},
			},
			UriNoId: true,
			EntityId: generic.EntityIdOptions{
				AttrNameGraph: "tenantId",
			},
			ReadOptions: generic.ReadOptions{
				DataSource: generic.DataSourceOptions{
					NoFilterSupport: true,
					Singular: generic.SingularOptions{
						NoIdRequired: true,
					},
					Plural: generic.PluralOptions{
						NoDataSource: true,
					},
				},
			},
			WriteOptions: generic.WriteOptions{
				UpdateIfExistsOnCreate: true,
				UsePutForCreate:        true,
			},
		},
	}

	CrossTenantIdentitySyncPolicyPartnerSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&CrossTenantIdentitySyncPolicyPartnerResource)
)

var crossTenantIdentitySyncPolicyPartnerResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // crossTenantIdentitySyncPolicyPartner
		"tenant_id": schema.StringAttribute{
			Required:      true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
		},
		"deleted_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Shows the last date and time the policy was deleted.",
		},
		"display_name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Display name for the cross-tenant user synchronization policy. Use the name of the partner Microsoft Entra tenant to easily identify the policy. Optional.",
		},
		"group_sync_inbound": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // crossTenantGroupSyncInbound
				"is_sync_allowed": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvaluemodifier.BoolDefaultValue(false)},
					Computed:            true,
					MarkdownDescription: "Defines whether group objects should be synchronized from the partner tenant. `false` stops any current group synchronization from the source tenant to the target tenant. This property has no impact on existing groups that were synchronized. <br/> The _provider_ default value is `false`.",
				},
			},
			MarkdownDescription: "Defines whether groups can be synchronized from a partner tenant. Key. / Defines whether groups can be synchronized from a partner tenant, as defined in the **groupSyncInbound** property of [crossTenantIdentitySyncPolicyPartner](https://learn.microsoft.com/en-us/graph/api/resources/crossTenantIdentitySyncPolicyPartner?view=graph-rest-beta) object. Also see [Microsoft docs for crossTenantGroupSyncInbound](https://learn.microsoft.com/en-us/graph/api/resources/crosstenantgroupsyncinbound?view=graph-rest-beta). <br> ",
		},
		"user_sync_inbound": schema.SingleNestedAttribute{
			Required: true,
			Attributes: map[string]schema.Attribute{ // crossTenantUserSyncInbound
				"is_sync_allowed": schema.BoolAttribute{
					Required:            true,
					MarkdownDescription: "Defines whether user objects should be synchronized from the partner tenant. `false` causes any current user synchronization from the source tenant to the target tenant to stop. This property has no impact on existing users who have already been synchronized.",
				},
			},
			MarkdownDescription: "Defines whether users can be synchronized from the partner tenant. Key. / Defines whether users can be synchronized from the partner tenant. Also see [Microsoft docs for crossTenantUserSyncInbound](https://learn.microsoft.com/en-us/graph/api/resources/crosstenantusersyncinbound?view=graph-rest-beta). <br> ",
		},
	},
	MarkdownDescription: "Defines the cross-tenant policy for synchronization of users from a partner tenant. Use this user synchronization policy to streamline collaboration between users in a multi-tenant organization by automating the creation, update, and deletion of users from one tenant to another. <br/> Also see [Microsoft docs for crossTenantIdentitySyncPolicyPartner](https://learn.microsoft.com/en-us/graph/api/resources/crosstenantidentitysyncpolicypartner?view=graph-rest-beta). ||| MS Graph: Cross-tenant access",
}
