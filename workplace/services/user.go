package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	userResource = generic.GenericResource{
		TypeNameSuffix: "user",
		SpecificSchema: userResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/users",
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

	UserSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&userResource)

	UserPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&userResource, "")
)

var userResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // user
		"id": schema.StringAttribute{
			MarkdownDescription: "The user identifier.",
		},
		"assigned_licenses": schema.SetNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // assignedLicense
					"disabled_plans": schema.SetAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "A collection of the unique identifiers for plans that have been disabled. IDs are available in **servicePlans** > **servicePlanId** in the tenant's [subscribedSkus](https://learn.microsoft.com/en-us/graph/api/resources/subscribedsku?view=graph-rest-beta) or **serviceStatus** > **servicePlanId** in the tenant's [companySubscription](https://learn.microsoft.com/en-us/graph/api/resources/subscribedsku?view=graph-rest-beta).",
					},
					"sku_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The unique identifier for the SKU. Corresponds to the **skuId** from [subscribedSkus](https://learn.microsoft.com/en-us/graph/api/resources/subscribedsku?view=graph-rest-beta) or [companySubscription](https://learn.microsoft.com/en-us/graph/api/resources/companysubscription?view=graph-rest-beta).",
					},
				},
			},
			MarkdownDescription: "Represents a license assigned to a user or group. The **assignedLicenses** property of the [user](user.md) or [group](group.md) entitity is a collection of **assignedLicense** objects. Also see [Microsoft docs for assignedLicense](https://learn.microsoft.com/en-us/graph/api/resources/assignedlicense?view=graph-rest-beta). <br> ",
		},
		"display_name": schema.StringAttribute{
			Optional: true,
		},
	},
	MarkdownDescription: "Represents an Azure Active Directory user object. <br/> Also see [Microsoft docs for user](https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-user?view=graph-rest-beta).\n\n_Provider_ Note: This data source is only provided as a companion to `azuread_user` to allow for OData filtering. It is not planned to add more attributes or even provide a resource for this entity. ||| MS Graph: Entra ID",
}
