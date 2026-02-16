package wpdefaultvalue

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
)

//
// Init
//

type DefaultValueWithInit interface {
	Init(context.Context, *diag.Diagnostics, *rsschema.Schema, path.Path)
}

// fwschema.AttributeWithObjectDefaultValue etc. are internal, so we have to duplicate them here
type WithObjectDefaultValue interface {
	ObjectDefaultValue() defaults.Object
}
type WithSetDefaultValue interface {
	SetDefaultValue() defaults.Set
}
type WithListDefaultValue interface {
	ListDefaultValue() defaults.List
}

func Init(ctx context.Context, diags *diag.Diagnostics, schema *rsschema.Schema) {
	for name, attribute := range schema.GetAttributes() {
		initAttribute(ctx, diags, schema, path.Root(name), attribute)
	}
}

func initAttribute(ctx context.Context, diags *diag.Diagnostics, schema *rsschema.Schema, path path.Path, attribute any) {

	if _, ok := attribute.(rsschema.Attribute); !ok {
		panic(fmt.Sprintf("initAttribute: attribute is not rsschema.Attribute but %T", attribute))
	}

	// As per definition only complex or collection types can have a default of DefaultValueWithInit, we can ignore simple types completely for now
	var defaultValueWithInit DefaultValueWithInit
	switch typedAttribute := attribute.(type) {
	case WithObjectDefaultValue:
		defaultValueWithInit, _ = typedAttribute.ObjectDefaultValue().(DefaultValueWithInit)
	case WithSetDefaultValue:
		defaultValueWithInit, _ = typedAttribute.SetDefaultValue().(DefaultValueWithInit)
	case WithListDefaultValue:
		defaultValueWithInit, _ = typedAttribute.ListDefaultValue().(DefaultValueWithInit)
	}
	if defaultValueWithInit != nil {
		// There would be no real benefit in delaying the actual type conversion until a default value actually gets used
		// since the TF framework will anyway request all default values while validating all resources and data sources.
		defaultValueWithInit.Init(ctx, diags, schema, path)
	}

	// Now traverse nested attributes
	if nestedAttribute, ok := attribute.(rsschema.NestedAttribute); ok {
		nestedAttributeObject := nestedAttribute.GetNestedObject()
		for name, attribute := range nestedAttributeObject.GetAttributes() {
			initAttribute(ctx, diags, schema, path.AtName(name), attribute)
		}
	}
}

//
// defaultValueBase
//

type defaultValueBase struct {
	defaultValRaw any
}

func (d defaultValueBase) Description(_ context.Context) string {
	return fmt.Sprintf("value defaults to %+v", d.defaultValRaw)
}

func (d defaultValueBase) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("value defaults to `%+v`", d.defaultValRaw)
}
