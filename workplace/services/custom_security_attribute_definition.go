package services

import (
	"context"
	"fmt"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	CustomSecurityAttributeDefinitionResource = generic.GenericResource{
		TypeNameSuffix: "custom_security_attribute_definition",
		SpecificSchema: customSecurityAttributeDefinitionResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/directory/customSecurityAttributeDefinitions",
			EntityId: generic.EntityIdOptions{
				FallbackIdGetterFunc: customSecurityAttributeDefinitionFallbackIdGetterFunc,
			},
			ReadOptions: generic.ReadOptions{
				ODataExpand: "allowedValues",
				DataSource: generic.DataSourceOptions{
					ExtraFilterAttributes: []string{"attribute_set", "status", "type"},
					Plural: generic.PluralOptions{
						ExtraAttributes: []string{"attribute_set", "status", "type"},
					},
				},
			},
			WriteOptions: generic.WriteOptions{
				UpdateIfExistsOnCreate: true,
				SkipDelete:             true,
				SubActions: []generic.WriteSubAction{
					&generic.WriteSubActionIndividual{
						ComparisonKeyAttribute: "id",
						SetNestedPath:          tftypes.NewAttributePath().WithAttributeName("allowed_values"),
						WriteSubActionBase: generic.WriteSubActionBase{
							Attributes: []string{"allowedValues"},
							UriSuffix:  "allowedValues",
							UpdateOnly: true,
						},
					},
				},
			},
			TerraformToGraphMiddleware: customSecurityAttributeDefinitionTerraformToGraphMiddleware,
		},
	}

	CustomSecurityAttributeDefinitionSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&CustomSecurityAttributeDefinitionResource)

	CustomSecurityAttributeDefinitionPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&CustomSecurityAttributeDefinitionResource, "")
)

func customSecurityAttributeDefinitionFallbackIdGetterFunc(ctx context.Context, diags *diag.Diagnostics, idAttributer generic.GetAttributer) string {
	var attributeSet, name types.String // might be null
	diags.Append(idAttributer.GetAttribute(ctx, path.Root("attribute_set"), &attributeSet)...)
	diags.Append(idAttributer.GetAttribute(ctx, path.Root("name"), &name)...)
	return fmt.Sprintf("%s_%s", attributeSet.ValueString(), name.ValueString()) // nil will become "" in ValueString()
}

func customSecurityAttributeDefinitionTerraformToGraphMiddleware(ctx context.Context, diags *diag.Diagnostics, params *generic.TerraformToGraphMiddlewareParams) generic.TerraformToGraphMiddlewareReturns {
	if params.IsUpdate {
		// some attributes cannot be updated and may not even be written again after creation
		delete(params.RawVal, "attributeSet")
		delete(params.RawVal, "name")
		delete(params.RawVal, "isCollection")
		delete(params.RawVal, "isSearchable")
		delete(params.RawVal, "type")
		// usePreDefinedValuesOnly cannot be updated from false to true and may not even be sent in a PATCH requests if it's true
		if params.RawVal["usePreDefinedValuesOnly"] == true {
			delete(params.RawVal, "usePreDefinedValuesOnly")
		}
	}
	return nil
}

var customSecurityAttributeDefinitionResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // customSecurityAttributeDefinition
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Identifier of the custom security attribute, which is a combination of the attribute set name and the custom security attribute name separated by an underscore (`attributeSet`_`name`). The **id** property is auto generated and cannot be set. Case insensitive.",
		},
		"attribute_set": schema.StringAttribute{
			Required:            true,
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			MarkdownDescription: "Name of the attribute set. Case insensitive.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Description of the custom security attribute. Can be up to 128 characters long and include Unicode characters. Can be changed later.",
		},
		"is_collection": schema.BoolAttribute{
			Optional: true,
			PlanModifiers: []planmodifier.Bool{
				wpdefaultvalue.BoolDefaultValue(false),
				wpplanmodifier.IsImmutable{},
			},
			Computed:            true,
			MarkdownDescription: "Indicates whether multiple values can be assigned to the custom security attribute. Cannot be changed later. If **type** is set to `Boolean`, **isCollection** cannot be set to `true`. The _provider_ default value is `false`.",
		},
		"is_searchable": schema.BoolAttribute{
			Optional: true,
			PlanModifiers: []planmodifier.Bool{
				wpdefaultvalue.BoolDefaultValue(true),
				wpplanmodifier.IsImmutable{},
			},
			Computed:            true,
			MarkdownDescription: "Indicates whether custom security attribute values are indexed for searching on objects that are assigned attribute values. Cannot be changed later. The _provider_ default value is `true`.",
		},
		"name": schema.StringAttribute{
			Required:            true,
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			MarkdownDescription: "Name of the custom security attribute. Must be unique within an attribute set. Can be up to 32 characters long and include Unicode characters. Cannot contain spaces or special characters. Cannot be changed later. Case insensitive.",
		},
		"status": schema.StringAttribute{
			Optional:            true,
			Validators:          []validator.String{stringvalidator.OneOf("Available", "Deprecated")},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("Available")},
			Computed:            true,
			MarkdownDescription: "Specifies whether the custom security attribute is active or deactivated. Acceptable values are: `Available` and `Deprecated`. Can be changed later. The _provider_ default value is `\"Available\"`.",
		},
		"type": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf("Boolean", "Integer", "String"),
			},
			PlanModifiers:       []planmodifier.String{wpplanmodifier.IsImmutable{}},
			MarkdownDescription: "Data type for the custom security attribute values. Supported types are: `Boolean`, `Integer`, and `String`. Cannot be changed later.",
		},
		"use_pre_defined_values_only": schema.BoolAttribute{
			Optional: true,
			PlanModifiers: []planmodifier.Bool{
				wpdefaultvalue.BoolDefaultValue(false),
				wpplanmodifier.IsImmutable{OnlyIfCurrentValueIs: false},
			},
			Computed:            true,
			MarkdownDescription: "Indicates whether only predefined values can be assigned to the custom security attribute. If set to `false`, free-form values are allowed. Can later be changed from `true` to `false`, but cannot be changed from `false` to `true`. If **type** is set to `Boolean`, **usePreDefinedValuesOnly** cannot be set to `true`. The _provider_ default value is `false`.",
		},
		"allowed_values": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // allowedValue
					"id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Identifier for the predefined value. Can be up to 64 characters long and include Unicode characters. Can include spaces, but some special characters aren't allowed. Can't be changed later. Case sensitive.",
					},
					"is_active": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
						Computed:            true,
						MarkdownDescription: "Indicates whether the predefined value is active or deactivated. If set to `false`, this predefined value can't be assigned to any more supported directory objects. The _provider_ default value is `true`.",
					},
				},
			},
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Values that are predefined for this custom security attribute. This navigation property is not returned by default and must be specified in an `$expand` query. For example, `/directory/customSecurityAttributeDefinitions?$expand=allowedValues`. / Represents a predefined value that is allowed for a custom security attribute definition.\n\nYou can define up to 100 **allowedValue** objects per [customSecurityAttributeDefinition](customsecurityattributedefinition.md). The **allowedValue** object can't be renamed or deleted, but it can be deactivated by using the [Update allowedValue](../api/../api/allowedvalue-update.md) operation. This object is defined as a navigation property on the [customSecurityAttributeDefinition](customsecurityattributedefinition.md) resource and its value is returned only on `$expand`. / https://learn.microsoft.com/en-us/graph/api/resources/allowedvalue?view=graph-rest-beta. The _provider_ default value is `[]`.",
		},
	},
	MarkdownDescription: "Represents the schema of a custom security attribute (key-value pair). For example, the custom security attribute name, description, data type, and allowed values.\n\nYou can define up to 500 active objects in a tenant. The **customSecurityAttributeDefinition** object can't be renamed or deleted, but it can be deactivated by using the [Update customSecurityAttributeDefinition](../api/customsecurityattributedefinition-update.md) operation. Must be part of an [attributeSet](../resources/attributeset.md). / https://learn.microsoft.com/en-us/graph/api/resources/customsecurityattributedefinition?view=graph-rest-beta ||| MS Graph: Custom security attributes",
}
