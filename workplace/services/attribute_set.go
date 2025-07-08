package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var (
	AttributeSetResource = generic.GenericResource{
		TypeNameSuffix: "attribute_set",
		SpecificSchema: attributeSetResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/directory/attributeSets",
			ReadOptions: generic.ReadOptions{
				DataSource: generic.DataSourceOptions{
					NoFilterSupport: true,
					Plural: generic.PluralOptions{
						ExtraAttributes: []string{"max_attributes_per_set"},
					},
				},
			},
			WriteOptions: generic.WriteOptions{
				UpdateIfExistsOnCreate: true,
				SkipDelete:             true,
			},
		},
	}

	AttributeSetSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&AttributeSetResource)

	AttributeSetPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&AttributeSetResource, "")
)

var attributeSetResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // attributeSet
		"id": schema.StringAttribute{
			Required:            true,
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			MarkdownDescription: "Identifier for the attribute set that is unique within a tenant. Can be up to 32 characters long and include Unicode characters. Cannot contain spaces or special characters. Cannot be changed later. Case insensitive.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Description of the attribute set. Can be up to 128 characters long and include Unicode characters. Can be changed later.",
		},
		"max_attributes_per_set": schema.Int64Attribute{
			Optional:            true,
			MarkdownDescription: "Maximum number of custom security attributes that can be defined in this attribute set. Default value is `null`. If not specified, the administrator can add up to the maximum of 500 active attributes per tenant. Can be changed later.",
		},
	},
	MarkdownDescription: "Represents a group of related custom security attribute definitions.\n\nYou can define up to 500 **attributeSet** objects in a tenant. The **attributeSet** object can't be renamed or deleted. / https://learn.microsoft.com/en-us/graph/api/resources/attributeset?view=graph-rest-beta ||| MS Graph: Custom security attributes",
}
