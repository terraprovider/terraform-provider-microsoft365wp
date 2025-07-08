package generic

import (
	"fmt"
	"terraform-provider-microsoft365wp/workplace/external/strcase"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"

	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type ToFromGraphTranslator struct {
	// use dummy SingleNestedAttribute for SchemaRoot to always return the same type
	SchemaRoot               rsschema.NestedAttribute
	IsDataSource             bool
	IncludeNullObjectsInJson bool
}

func NewToFromGraphTranslator(schema tftypes.AttributePathStepper, includeNullObjectsInJson bool) ToFromGraphTranslator {

	result := ToFromGraphTranslator{
		IncludeNullObjectsInJson: includeNullObjectsInJson,
	}

	// create dummy SingleNestedAttribute here to always return the same type
	// cannot really de-dup the code since GetAttributes() returns internal fwschema type :-(
	switch typ := schema.(type) {
	case rsschema.Schema:
		nestedAttributes := make(map[string]rsschema.Attribute, 0)
		for k, v := range typ.GetAttributes() {
			nestedAttributes[k] = v
		}
		result.SchemaRoot = rsschema.SingleNestedAttribute{
			Attributes:  nestedAttributes,
			Description: typ.GetDescription(),
		}
	case dsschema.Schema:
		nestedAttributes := make(map[string]dsschema.Attribute, 0)
		for k, v := range typ.GetAttributes() {
			nestedAttributes[k] = v
		}
		result.SchemaRoot = dsschema.SingleNestedAttribute{
			Attributes:  nestedAttributes,
			Description: typ.GetDescription(),
		}
		result.IsDataSource = true
	case rsschema.NestedAttribute:
		// support serializing of individual nested objects, e.g. elements of sets or lists
		// ensure SingleNestedAttribute as we mostly receive SetNestedAttribute or ListNestedAttribute
		nestedAttributes := make(map[string]rsschema.Attribute, 0)
		for k, v := range typ.GetNestedObject().GetAttributes() {
			nestedAttributes[k] = v
		}
		result.SchemaRoot = rsschema.SingleNestedAttribute{
			Attributes:  nestedAttributes,
			Description: typ.GetDescription(),
		}
	default:
		panic(fmt.Errorf("got unexpected type %T", schema))
	}

	return result
}

func (t *ToFromGraphTranslator) GraphAttributeNameFromTerraformName(parentAttribute rsschema.Attribute, terraformAttributeName string) (string, bool, bool) {

	// cannot really extract this since GetAttributes() returns internal fwschema type :-(
	parentNested, ok := parentAttribute.(rsschema.NestedAttribute)
	if !ok {
		panic(fmt.Errorf("parentAttribute %s is not NestedAttribute (terraformAttributeName: %s)", parentAttribute, terraformAttributeName))
	}

	if attribute, ok := parentNested.GetNestedObject().GetAttributes()[terraformAttributeName]; ok {
		graphAttributeName := t.GraphAttributeNameFromTerraformNameImpl(terraformAttributeName, attribute)
		attributeIsWritable := attribute.IsRequired() || attribute.IsOptional()
		return graphAttributeName, attributeIsWritable, true
	}

	return "", false, false
}

func (*ToFromGraphTranslator) GraphAttributeNameFromTerraformNameImpl(terraformAttributeName string, attribute GetDescriptioner) string {

	if description := attribute.GetDescription(); description != "" {
		// custom mapping
		return description
	}
	return strcase.ToLowerCamel(terraformAttributeName)
}

func (*ToFromGraphTranslator) TerraformAttributeNameFromGraphName(parentAttribute rsschema.Attribute, graphAttributeName string) (string, bool) {

	// cannot really extract this since GetAttributes() returns internal fwschema type :-(
	parentNested, ok := parentAttribute.(rsschema.NestedAttribute)
	if !ok {
		panic(fmt.Errorf("parentAttribute %s is not NestedAttribute (graphAttributeName: %s)", parentAttribute, graphAttributeName))
	}

	for k, v := range parentNested.GetNestedObject().GetAttributes() {
		description := v.GetDescription()
		if description != "" && description == graphAttributeName {
			// custom mapping
			return k, true
		}
		if derivedTypeNested, ok := v.(DerivedTypeNestedAttribute); ok {
			if derivedTypeNested.GetDerivedType() == graphAttributeName {
				// OData type name -> nested attribute name
				return k, true
			}
		}
	}

	terraformAttributeName := strcase.ToSnake(graphAttributeName)
	_, ok = parentNested.GetNestedObject().GetAttributes()[terraformAttributeName]
	if ok {
		// ok, attribute exists
		return terraformAttributeName, true
	}

	return "", false
}

func (*ToFromGraphTranslator) OdataTypeByTerraformAttributeName(parentAttribute rsschema.Attribute, terraformAttributeName string) (string, bool) {

	// cannot really extract this since GetAttributes() returns internal fwschema type :-(
	parentNested, ok := parentAttribute.(rsschema.NestedAttribute)
	if !ok {
		panic(fmt.Errorf("parentAttribute %s is not NestedAttribute (terraformAttributeName: %s)", parentAttribute, terraformAttributeName))
	}

	if attribute, ok := parentNested.GetNestedObject().GetAttributes()[terraformAttributeName]; ok {
		if derivedTypeNested, ok := attribute.(DerivedTypeNestedAttribute); ok {
			return derivedTypeNested.GetDerivedType(), true
		}
	}

	return "", false
}

func (t *ToFromGraphTranslator) NestedObjectAttributesByOdataType(parentAttribute rsschema.Attribute, odataTypeName string) ([]string, bool) {

	// cannot really extract this since GetAttributes() returns internal fwschema type :-(
	parentNested, ok := parentAttribute.(rsschema.NestedAttribute)
	if !ok {
		panic(fmt.Errorf("parentAttribute %s is not NestedAttribute (odataTypeName: %s)", parentAttribute, odataTypeName))
	}

	for _, v := range parentNested.GetNestedObject().GetAttributes() {
		if derivedTypeNested, ok := v.(DerivedTypeNestedAttribute); ok {
			if derivedTypeNested.GetDerivedType() == odataTypeName {
				nestedAttributeNames := make([]string, 0)
				for k := range derivedTypeNested.GetNestedObject().GetAttributes() {
					graphAttributeName, _, ok := t.GraphAttributeNameFromTerraformName(derivedTypeNested, k)
					if ok {
						nestedAttributeNames = append(nestedAttributeNames, graphAttributeName)
					}
				}
				return nestedAttributeNames, true
			}
		}
	}

	return nil, false
}

func (t *ToFromGraphTranslator) getAllValidatorsForPath(p *tftypes.AttributePath) ([]interface{}, error) {
	attr, err := t.SchemaAttributeAtTerraformPath(p)
	if err != nil {
		return nil, err
	}

	var validators []interface{}
	switch attr := attr.(type) {
	case rsschema.BoolAttribute:
		for _, v := range attr.BoolValidators() {
			validators = append(validators, v)
		}
	case dsschema.BoolAttribute:
		for _, v := range attr.BoolValidators() {
			validators = append(validators, v)
		}
	case rsschema.Float64Attribute:
		for _, v := range attr.Float64Validators() {
			validators = append(validators, v)
		}
	case dsschema.Float64Attribute:
		for _, v := range attr.Float64Validators() {
			validators = append(validators, v)
		}
	case rsschema.StringAttribute:
		for _, v := range attr.StringValidators() {
			validators = append(validators, v)
		}
	case dsschema.StringAttribute:
		for _, v := range attr.StringValidators() {
			validators = append(validators, v)
		}
	case rsschema.SetAttribute:
		for _, v := range attr.SetValidators() {
			validators = append(validators, v)
		}
	case dsschema.SetAttribute:
		for _, v := range attr.SetValidators() {
			validators = append(validators, v)
		}
	}

	return validators, nil
}

//
// We use dummy validators of type wpvalidator.AttributeValueTranslator as internal flags to signal that an attribute
// value being sent to or received from MS Graph might need to be translated to some other value for Terraform.
// See wpvalidator.AttributeValueTranslator for more details.
//

func (t *ToFromGraphTranslator) GraphToTerraformDefaultValue(p *tftypes.AttributePath, attributeIsMissing bool, typ tftypes.Type) (tftypes.Value, error) {

	if attributeIsMissing {
		tfValue, ok, err := t.GraphToTerraformTranslateValue(p, wpvalidator.AttributeMising)
		if err != nil || ok {
			return tfValue, err
		}
	}

	return tftypes.NewValue(typ, nil), nil
}

func (t *ToFromGraphTranslator) GraphToTerraformTranslateValue(p *tftypes.AttributePath, graphValue any) (tftypes.Value, bool, error) {

	validators, err := t.getAllValidatorsForPath(p)
	if err != nil {
		return tftypes.Value{}, false, err
	}

	for _, v := range validators {
		if v, ok := v.(wpvalidator.AttributeValueTranslator); ok {
			tfValue, ok, err := v.TranslateGraphToTerraform(graphValue)
			if err != nil {
				return tftypes.Value{}, false, err
			}
			if ok {
				return tfValue, true, nil
			}
		}
	}

	return tftypes.Value{}, false, nil
}

func (t *ToFromGraphTranslator) TerraformToGraphTranslateValue(p *tftypes.AttributePath, tfValue tftypes.Value) (any, bool, error) {

	validators, err := t.getAllValidatorsForPath(p)
	if err != nil {
		return nil, false, err
	}

	for _, v := range validators {
		if v, ok := v.(wpvalidator.AttributeValueTranslator); ok {
			graphValue, ok, err := v.TranslateTerraformToGraph(tfValue)
			if err != nil {
				return nil, false, err
			}
			if ok {
				return graphValue, true, nil
			}
		}
	}

	return nil, false, nil
}
