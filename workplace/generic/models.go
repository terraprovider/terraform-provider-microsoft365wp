package generic

import (
	"fmt"

	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"golang.org/x/exp/slices"
)

func ConvertResourceSchemaToDataSourceSingular(resourceSchema *rsschema.Schema, accessParams *AccessParams) dsschema.Schema {

	dsSchemaAttributes := make(map[string]dsschema.Attribute)

	hasParentItem := !accessParams.HasParentItem.ParentIdField.Equal(path.Path{})
	parentIdFieldString := accessParams.HasParentItem.ParentIdField.String()

	for aKey, aValue := range resourceSchema.Attributes {
		required := (aKey == "id" && !accessParams.IsSingleton) || (hasParentItem && aKey == parentIdFieldString)
		if dsAttribute := convertResource2DataSourceAttribute(aValue, required); dsAttribute != nil {
			dsSchemaAttributes[aKey] = dsAttribute
		}
	}

	return dsschema.Schema{
		Attributes: dsSchemaAttributes,
	}
}

func convertResource2DataSourceAttribute(rsAttribute rsschema.Attribute, required bool) dsschema.Attribute {

	switch typed := rsAttribute.(type) {
	case rsschema.StringAttribute:
		typed.Required = required
		typed.Optional = false
		typed.Computed = !required
		return typed
	case rsschema.BoolAttribute:
		typed.Required = required
		typed.Optional = false
		typed.Computed = !required
		return typed
	case rsschema.Int64Attribute:
		typed.Required = required
		typed.Optional = false
		typed.Computed = !required
		return typed
	case rsschema.NumberAttribute:
		typed.Required = required
		typed.Optional = false
		typed.Computed = !required
		return typed
	case rsschema.SetAttribute:
		typed.Required = required
		typed.Optional = false
		typed.Computed = !required
		return typed
	case rsschema.SingleNestedAttribute:
		dsAttributes := make(map[string]dsschema.Attribute)
		for nKey, nValue := range typed.Attributes {
			dsAttributes[nKey] = convertResource2DataSourceAttribute(nValue, required)
		}
		return dsschema.SingleNestedAttribute{
			Attributes:          dsAttributes,
			CustomType:          typed.CustomType,
			Required:            required,
			Optional:            false,
			Computed:            !required,
			Sensitive:           typed.Sensitive,
			Description:         typed.Description,
			MarkdownDescription: typed.MarkdownDescription,
			DeprecationMessage:  typed.DeprecationMessage,
			Validators:          typed.Validators,
		}
	case OdataDerivedTypeNestedAttributeRs:
		dsAttributes := make(map[string]dsschema.Attribute)
		for nKey, nValue := range typed.Attributes {
			dsAttributes[nKey] = convertResource2DataSourceAttribute(nValue, required)
		}
		return OdataDerivedTypeNestedAttributeDs{
			SingleNestedAttribute: dsschema.SingleNestedAttribute{
				Attributes:          dsAttributes,
				CustomType:          typed.CustomType,
				Required:            required,
				Optional:            false,
				Computed:            !required,
				Sensitive:           typed.Sensitive,
				Description:         typed.Description,
				MarkdownDescription: typed.MarkdownDescription,
				DeprecationMessage:  typed.DeprecationMessage,
			},
			DerivedType: typed.DerivedType,
		}
	case rsschema.SetNestedAttribute:
		dsAttributes := make(map[string]dsschema.Attribute)
		for nKey, nValue := range typed.NestedObject.Attributes {
			dsAttributes[nKey] = convertResource2DataSourceAttribute(nValue, required)
		}
		nestedObject := dsschema.NestedAttributeObject{
			Attributes: dsAttributes,
			CustomType: typed.NestedObject.CustomType,
			Validators: typed.NestedObject.Validators,
		}
		return dsschema.SetNestedAttribute{
			NestedObject:        nestedObject,
			CustomType:          typed.CustomType,
			Required:            required,
			Optional:            false,
			Computed:            !required,
			Sensitive:           typed.Sensitive,
			Description:         typed.Description,
			MarkdownDescription: typed.MarkdownDescription,
			DeprecationMessage:  typed.DeprecationMessage,
		}
	case rsschema.ListNestedAttribute:
		dsAttributes := make(map[string]dsschema.Attribute)
		for nKey, nValue := range typed.NestedObject.Attributes {
			dsAttributes[nKey] = convertResource2DataSourceAttribute(nValue, required)
		}
		nestedObject := dsschema.NestedAttributeObject{
			Attributes: dsAttributes,
			CustomType: typed.NestedObject.CustomType,
			Validators: typed.NestedObject.Validators,
		}
		return dsschema.ListNestedAttribute{
			NestedObject:        nestedObject,
			CustomType:          typed.CustomType,
			Required:            required,
			Optional:            false,
			Computed:            !required,
			Sensitive:           typed.Sensitive,
			Description:         typed.Description,
			MarkdownDescription: typed.MarkdownDescription,
			DeprecationMessage:  typed.DeprecationMessage,
		}
	}

	panic(fmt.Sprintf("Don't know how to convert attribute of type %T", rsAttribute))
}

func ConvertDataSourceSingularSchemaToPlural(singularSchema *dsschema.Schema, nestedSetName string, parentItem *HasParentItem) dsschema.Schema {
	hasParentItem := !parentItem.ParentIdField.Equal(path.Path{})
	parentIdFieldString := parentItem.ParentIdField.String()

	nestedSetAttributeFilter := []string{
		"id",
		"created_date_time",
		"display_name",
		"last_modified_date_time",
		"name",
		"role_scope_tag_ids",
		"role_scope_tags",
		"version",
	}

	rootAttributes := make(map[string]dsschema.Attribute)
	leafAttributes := make(map[string]dsschema.Attribute)
	for k, v := range singularSchema.Attributes {
		if hasParentItem && k == parentIdFieldString {
			rootAttributes[k] = v
		} else if slices.Contains(nestedSetAttributeFilter, k) {
			leafAttributes[k] = v
		} else if typed, ok := v.(OdataDerivedTypeNestedAttributeDs); ok {
			// Also include virtual attributes that signal a derived OData type. They will not have any child attributes or content
			// at all, but it might still be helpful to be able to check for a derived type by comparing with null in Terraform.
			leafAttributes[k] = OdataDerivedTypeNestedAttributeDs{
				SingleNestedAttribute: dsschema.SingleNestedAttribute{
					CustomType:  typed.CustomType,
					Required:    false,
					Optional:    false,
					Computed:    true,
					Sensitive:   typed.Sensitive,
					Description: typed.Description,
					MarkdownDescription: fmt.Sprintf("Please note that this nested object does not have any attributes "+
						"but only exists to be able to test if the parent object is of derived OData type `%s` "+
						"(using e.g. `if x.%s != null`).", typed.DerivedType, k),
					DeprecationMessage: typed.DeprecationMessage,
				},
				DerivedType: typed.DerivedType,
			}
		}
	}

	rootAttributes[nestedSetName] = dsschema.SetNestedAttribute{
		Computed: true,
		NestedObject: dsschema.NestedAttributeObject{
			Attributes: leafAttributes,
		},
	}

	return dsschema.Schema{
		Attributes: rootAttributes,
	}
}
