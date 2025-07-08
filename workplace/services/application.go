package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var (
	applicationResource = generic.GenericResource{
		TypeNameSuffix: "application",
		SpecificSchema: applicationResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/applications",
			ReadOptions: generic.ReadOptions{
				DataSource: generic.DataSourceOptions{
					ExtraFilterAttributes: []string{"app_id"},
					Plural: generic.PluralOptions{
						ExtraAttributes: []string{"app_id"},
					},
				},
			},
		},
	}

	ApplicationSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&applicationResource)

	ApplicationPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&applicationResource, "")
)

var applicationResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // application
		"id": schema.StringAttribute{
			MarkdownDescription: "Unique identifier for the application object. This property is referred to as **Object ID** in the Microsoft Entra admin center. Key. Not nullable. Read-only. Supports `$filter` (`eq`, `ne`, `not`, `in`).",
		},
		"app_id": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The unique identifier for the application that is assigned by Microsoft Entra ID. Not nullable. Read-only. Alternate key. Supports `$filter` (`eq`).",
		},
		"display_name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The display name for the application. Maximum length is 256 characters. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values), `$search`, and `$orderby`.",
		},
	},
	MarkdownDescription: "Represents an application. Any application that outsources authentication to Microsoft Entra ID must be registered in the Microsoft identity platform. Application registration involves telling Microsoft Entra ID about your application, including the URL where it's located, the URL to send replies after authentication, the URI to identify your application, and more.\n\nThis resource is an open type that allows other properties to be passed in.\n\nThis resource supports: / https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta\n\nProvider Note: This data source is only provided as a companion to `azuread_application` to allow for OData filtering. It is not planned to add more attributes or even provide a resource for this entity. ||| MS Graph: Entra ID",
}
