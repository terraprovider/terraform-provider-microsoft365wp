package generic

import (
	"fmt"
	"regexp"

	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func convertResourceAttr2DataSourceAttr(rsAttribute rsschema.Attribute, required bool, odataDerivedTypeIsPlaceholderOnly bool, odataDerivedTypeDescNameTf string) dsschema.Attribute {

	optional := false
	computed := !required

	// remove any "read-only." or "(the )(provider )default value is..." statements from the description
	descCleanupRegex := regexp.MustCompile(`(?i)\s*read-only\.|\s*(the )?(.?provider.? )?default value is \S+`)

	switch typed := rsAttribute.(type) {
	case rsschema.StringAttribute:
		typed.Required = required
		typed.Optional = optional
		typed.Computed = computed
		typed.MarkdownDescription = descCleanupRegex.ReplaceAllLiteralString(typed.MarkdownDescription, "")
		return typed
	case rsschema.BoolAttribute:
		typed.Required = required
		typed.Optional = optional
		typed.Computed = computed
		typed.MarkdownDescription = descCleanupRegex.ReplaceAllLiteralString(typed.MarkdownDescription, "")
		return typed
	case rsschema.Int64Attribute:
		typed.Required = required
		typed.Optional = optional
		typed.Computed = computed
		typed.MarkdownDescription = descCleanupRegex.ReplaceAllLiteralString(typed.MarkdownDescription, "")
		return typed
	case rsschema.NumberAttribute:
		typed.Required = required
		typed.Optional = optional
		typed.Computed = computed
		typed.MarkdownDescription = descCleanupRegex.ReplaceAllLiteralString(typed.MarkdownDescription, "")
		return typed
	case rsschema.SetAttribute:
		typed.Required = required
		typed.Optional = optional
		typed.Computed = computed
		typed.MarkdownDescription = descCleanupRegex.ReplaceAllLiteralString(typed.MarkdownDescription, "")
		return typed
	case rsschema.SingleNestedAttribute:
		dsAttributes := make(map[string]dsschema.Attribute)
		for nKey, nValue := range typed.Attributes {
			dsAttributes[nKey] = convertResourceAttr2DataSourceAttr(nValue, false, odataDerivedTypeIsPlaceholderOnly, nKey)
		}
		return dsschema.SingleNestedAttribute{
			Attributes:          dsAttributes,
			CustomType:          typed.CustomType,
			Required:            required,
			Optional:            optional,
			Computed:            computed,
			Sensitive:           typed.Sensitive,
			Description:         typed.Description,
			MarkdownDescription: descCleanupRegex.ReplaceAllLiteralString(typed.MarkdownDescription, ""),
			DeprecationMessage:  typed.DeprecationMessage,
			Validators:          typed.Validators,
		}
	case OdataDerivedTypeNestedAttributeRs:
		dsAttributes := make(map[string]dsschema.Attribute)
		mdDescription := ""
		if !odataDerivedTypeIsPlaceholderOnly {
			for nKey, nValue := range typed.Attributes {
				dsAttributes[nKey] = convertResourceAttr2DataSourceAttr(nValue, false, odataDerivedTypeIsPlaceholderOnly, nKey)
			}
			mdDescription = descCleanupRegex.ReplaceAllLiteralString(typed.MarkdownDescription, "")
		} else {
			mdDescription = fmt.Sprintf("Please note that this nested object does not have any attributes "+
				"but only exists to be able to test if the parent object is of derived OData type `%s` "+
				"(using e.g. `if x.%s != null`).", typed.DerivedType, odataDerivedTypeDescNameTf)
		}
		return OdataDerivedTypeNestedAttributeDs{
			SingleNestedAttribute: dsschema.SingleNestedAttribute{
				Attributes:          dsAttributes,
				CustomType:          typed.CustomType,
				Required:            required,
				Optional:            optional,
				Computed:            computed,
				Sensitive:           typed.Sensitive,
				Description:         typed.Description,
				MarkdownDescription: mdDescription,
				DeprecationMessage:  typed.DeprecationMessage,
			},
			DerivedType: typed.DerivedType,
		}
	case rsschema.SetNestedAttribute:
		dsAttributes := make(map[string]dsschema.Attribute)
		for nKey, nValue := range typed.NestedObject.Attributes {
			dsAttributes[nKey] = convertResourceAttr2DataSourceAttr(nValue, false, odataDerivedTypeIsPlaceholderOnly, nKey)
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
			Optional:            optional,
			Computed:            computed,
			Sensitive:           typed.Sensitive,
			Description:         typed.Description,
			MarkdownDescription: descCleanupRegex.ReplaceAllLiteralString(typed.MarkdownDescription, ""),
			DeprecationMessage:  typed.DeprecationMessage,
		}
	case rsschema.ListNestedAttribute:
		dsAttributes := make(map[string]dsschema.Attribute)
		for nKey, nValue := range typed.NestedObject.Attributes {
			dsAttributes[nKey] = convertResourceAttr2DataSourceAttr(nValue, false, odataDerivedTypeIsPlaceholderOnly, nKey)
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
			Optional:            optional,
			Computed:            computed,
			Sensitive:           typed.Sensitive,
			Description:         typed.Description,
			MarkdownDescription: descCleanupRegex.ReplaceAllLiteralString(typed.MarkdownDescription, ""),
			DeprecationMessage:  typed.DeprecationMessage,
		}
	}

	panic(fmt.Sprintf("Don't know how to convert attribute of type %T", rsAttribute))
}
