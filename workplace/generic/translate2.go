// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package generic

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// based on https://github.com/hashicorp/terraform-plugin-framework/blob/v1.1.1/internal/fwschema/schema.go
func (t ToFromGraphTranslator) SchemaTypeAtTerraformPath(ctx context.Context, p *tftypes.AttributePath) (attr.Type, error) {
	rawType, remaining, err := tftypes.WalkAttributePath(t.SchemaRoot, p)

	if err != nil {
		return nil, fmt.Errorf("%v still remains in the path: %w", remaining, err)
	}

	switch typ := rawType.(type) {
	case attr.Type:
		return typ, nil
	case rsschema.Attribute: // interface, will also match dsschema.Attribute
		return typ.GetType(), nil
	case rsschema.Block: // interface, will also match dsschema.Block
		return typ.Type(), nil
	case rsschema.NestedAttributeObject:
		return typ.Type(), nil
	case rsschema.NestedBlockObject:
		return typ.Type(), nil
	case rsschema.Schema:
		return typ.Type(), nil
	case dsschema.NestedAttributeObject:
		return typ.Type(), nil
	case dsschema.NestedBlockObject:
		return typ.Type(), nil
	case dsschema.Schema:
		return typ.Type(), nil
	// does not seem to be required for our purposes :-)
	// case rsschema.UnderlyingAttributes:
	// 	return typ.Type(), nil
	// case dsschema.UnderlyingAttributes:
	// 	return typ.Type(), nil
	default:
		return nil, fmt.Errorf("got unexpected type %T", rawType)
	}
}

// based on https://github.com/hashicorp/terraform-plugin-framework/blob/v1.1.1/internal/fwschema/schema.go
func (t ToFromGraphTranslator) SchemaAttributeAtTerraformPath(ctx context.Context, p *tftypes.AttributePath) (rsschema.Attribute, error) {

	// ensure for WalkAttributePath to get the attribute itself, not some contained NestedObject
	switch p.LastStep().(type) {
	case tftypes.ElementKeyInt, tftypes.ElementKeyString, tftypes.ElementKeyValue:
		p = p.WithoutLastStep()
	}
	rawType, remaining, err := tftypes.WalkAttributePath(t.SchemaRoot, p)
	if err != nil {
		return nil, fmt.Errorf("%v still remains in the path: %w", remaining, err)
	}

	switch typ := rawType.(type) {
	case rsschema.Attribute:
		return typ, nil
	default:
		return nil, fmt.Errorf("got unexpected type %T for path %s", rawType, p)
	}
}
