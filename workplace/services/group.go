package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	groupResource = generic.GenericResource{
		TypeNameSuffix: "group",
		SpecificSchema: groupResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/groups",
			ReadOptions: generic.ReadOptions{
				ODataSelect: []string{"displayName", "assignedLicenses"},
				DataSource: generic.DataSourceOptions{
					Plural: generic.PluralOptions{
						ExtraAttributes: []string{"assigned_licenses"},
					},
				},
			},
		},
	}

	GroupSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&groupResource)

	GroupPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&groupResource, "")
)

var groupResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // group
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the group. <br><br>Returned by default. Key. Not nullable. Read-only. <br><br>Supports `$filter` (`eq`, `ne`, `not`, `in`).",
		},
		"assigned_licenses": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // assignedLicense
					"disabled_plans": schema.SetAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "A collection of the unique identifiers for plans that have been disabled. IDs are available in **servicePlans** > **servicePlanId** in the tenant's [subscribedSkus](../resources/subscribedsku.md) or **serviceStatus** > **servicePlanId** in the tenant's [companySubscription](../resources/subscribedsku.md).",
					},
					"sku_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The unique identifier for the SKU. Corresponds to the **skuId** from [subscribedSkus](../resources/subscribedsku.md) or [companySubscription](../resources/companysubscription.md).",
					},
				},
			},
			MarkdownDescription: "The licenses that are assigned to the group. <br><br>Returned only on `$select`. Supports `$filter` (`eq`). Read-only. / Represents a license assigned to a user or group. The **assignedLicenses** property of the [user](user.md) or [group](group.md) entitity is a collection of **assignedLicense** objects. / https://learn.microsoft.com/en-us/graph/api/resources/assignedlicense?view=graph-rest-beta",
		},
		"display_name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The display name for the group. Required. Maximum length is 256 characters. <br><br>Returned by default. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values), `$search`, and `$orderby`.",
		},
	},
	MarkdownDescription: "Represents a Microsoft Entra group, which can be a Microsoft 365 group, a team in Microsoft Teams, or a security group. This resource is an open type that allows other properties to be passed in.\n\nFor performance reasons, the [create](../api/group-post-groups.md), [get](../api/group-get.md), and [list](../api/group-list.md) operations return only a subset of more commonly used properties by default. These _default_ properties are noted in the [Properties](#properties) section. To get any of the properties not returned by default, specify them in a `$select` OData query option.\n\nThis resource supports: / https://learn.microsoft.com/en-us/graph/api/resources/group?view=graph-rest-beta\n\nProvider Note: This data source is only provided as a companion to `azuread_group` to allow for OData filtering. It is not planned to add more attributes or even provide a resource for this entity. ||| MS Graph: Entra ID",
}
