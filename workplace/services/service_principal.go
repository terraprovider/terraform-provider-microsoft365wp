package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	servicePrincipalResource = generic.GenericResource{
		TypeNameSuffix: "service_principal",
		SpecificSchema: servicePrincipalResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/servicePrincipals",
			ReadOptions: generic.ReadOptions{
				DataSource: generic.DataSourceOptions{
					ExtraFilterAttributes: []string{"account_enabled", "app_id", "service_principal_type"},
					Plural: generic.PluralOptions{
						ExtraAttributes: []string{"account_enabled", "app_id", "service_principal_type"},
					},
				},
			},
		},
	}

	ServicePrincipalSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&servicePrincipalResource)

	ServicePrincipalPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&servicePrincipalResource, "")
)

var servicePrincipalResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // servicePrincipal
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the service principal. Key. Not nullable. Read-only. Supports `$filter` (`eq`, `ne`, `not`, `in`).",
		},
		"account_enabled": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "`true` if the service principal account is enabled; otherwise, `false`. If set to `false`, then no users are able to sign in to this app, even if they're assigned to it. Supports `$filter` (`eq`, `ne`, `not`, `in`).",
		},
		"app_id": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The unique identifier for the associated application (its **appId** property). Alternate key. Supports `$filter` (`eq`, `ne`, `not`, `in`, `startsWith`).",
		},
		"display_name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The display name for the service principal. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values), `$search`, and `$orderby`.",
		},
		"service_principal_type": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Identifies if the service principal represents an application or a managed identity. This is set by Microsoft Entra ID internally. For a service principal that represents an [application](./application.md) this is set as __Application__. For a service principal that represents a [managed identity](/azure/active-directory/managed-identities-azure-resources/overview) this is set as __ManagedIdentity__. The __SocialIdp__ type is for internal use.",
		},
	},
	MarkdownDescription: "Represents an instance of an application in a directory.\n\nThis resource supports using [delta query](/graph/delta-query-overview) to track incremental additions, deletions, and updates, by providing a [delta](../api/serviceprincipal-delta.md) function. This resource is an open type that allows other properties to be passed in. / https://learn.microsoft.com/en-us/graph/api/resources/serviceprincipal?view=graph-rest-beta\n\nProvider Note: This data source is only provided as a companion to `azuread_service_principal` to allow for OData filtering. It is not planned to add more attributes or even provide a resource for this entity. ||| MS Graph: Entra ID",
}
