package generic

import (
	"fmt"
	"regexp"
	"strings"

	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func ConvertResourceAttributesToDataSourceSingular(sourceRsAttributes map[string]rsschema.Attribute, isRequiredFunc func(string) bool) map[string]dsschema.Attribute {

	targetDsAttributes := make(map[string]dsschema.Attribute)
	for attrName, rsAttribute := range sourceRsAttributes {
		targetDsAttributes[attrName] = convertResourceAttr2DataSourceAttr(rsAttribute, isRequiredFunc(attrName))
	}
	return targetDsAttributes
}

func ConvertResourceAttributesToDataSourcePlural(sourceRsAttributes map[string]rsschema.Attribute, nestedListName string,
	isRootAndRequiredAttrFunc func(string) bool, isNestedListAttrFunc func(string, rsschema.Attribute) (bool, dsschema.Attribute)) map[string]dsschema.Attribute {

	targetRootDsAttributes := make(map[string]dsschema.Attribute)
	targetNestedListDsAttributes := make(map[string]dsschema.Attribute)
	for attrName, rsAttribute := range sourceRsAttributes {
		if isRootAndRequiredAttrFunc(attrName) {
			targetRootDsAttributes[attrName] = convertResourceAttr2DataSourceAttr(rsAttribute, true)
		}
		if ok, dsAttribute := isNestedListAttrFunc(attrName, rsAttribute); ok {
			if dsAttribute == nil {
				dsAttribute = convertResourceAttr2DataSourceAttr(rsAttribute, false)
			}
			targetNestedListDsAttributes[attrName] = dsAttribute
		}
	}

	targetRootDsAttributes[nestedListName] = dsschema.ListNestedAttribute{
		Computed: true,
		NestedObject: dsschema.NestedAttributeObject{
			Attributes: targetNestedListDsAttributes,
		},
	}

	return targetRootDsAttributes
}

func convertResourceAttr2DataSourceAttr(rsAttribute rsschema.Attribute, required bool) dsschema.Attribute {

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
			dsAttributes[nKey] = convertResourceAttr2DataSourceAttr(nValue, false)
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
		for nKey, nValue := range typed.Attributes {
			dsAttributes[nKey] = convertResourceAttr2DataSourceAttr(nValue, false)
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
				MarkdownDescription: descCleanupRegex.ReplaceAllLiteralString(typed.MarkdownDescription, ""),
				DeprecationMessage:  typed.DeprecationMessage,
			},
			DerivedType: typed.DerivedType,
		}
	case rsschema.SetNestedAttribute:
		dsAttributes := make(map[string]dsschema.Attribute)
		for nKey, nValue := range typed.NestedObject.Attributes {
			dsAttributes[nKey] = convertResourceAttr2DataSourceAttr(nValue, false)
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
			dsAttributes[nKey] = convertResourceAttr2DataSourceAttr(nValue, false)
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

func GetSubcategorySuffixFromMarkdownDescription(markdownDescription string) string {
	return " ||| " + strings.Split(markdownDescription, " ||| ")[1]
}
