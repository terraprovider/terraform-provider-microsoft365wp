// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//
// Lifted from https://github.com/hashicorp/terraform-plugin-framework-jsontypes/tree/v0.2.0/jsontypes
// Changes:
//   - Added WpObjectFilterFunc
//   - Added WpIgnoreArrayOrder
//

package wpjsontypes

import (
	"context"
	"fmt"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpobjectfilter"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable = (*NormalizedType)(nil)
)

// NormalizedType is an attribute type that represents a valid JSON string (RFC 7159). Semantic equality logic is defined for NormalizedType
// such that inconsequential differences between JSON strings are ignored (whitespace, property order, etc). If you need strict, byte-for-byte,
// string equality, consider using ExactType.
type NormalizedType struct {
	basetypes.StringType
	WpObjectFilterFunc wpobjectfilter.FilterFunc
	// WpIgnoreArrayOrder makes semantic equality ignore the element order of all (nested) arrays.
	// Only use it for JSON where no array order carries meaning.
	WpIgnoreArrayOrder bool
}

// String returns a human readable string of the type name.
func (t NormalizedType) String() string {
	return "wpjsontypes.NormalizedType"
}

// ValueType returns the Value type.
func (t NormalizedType) ValueType(ctx context.Context) attr.Value {
	return Normalized{
		WpObjectFilterFunc: t.WpObjectFilterFunc,
		WpIgnoreArrayOrder: t.WpIgnoreArrayOrder,
	}
}

// Equal returns true if the given type is equivalent.
func (t NormalizedType) Equal(o attr.Type) bool {
	other, ok := o.(NormalizedType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

// ValueFromString returns a StringValuable type given a StringValue.
func (t NormalizedType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return Normalized{
		StringValue:        in,
		WpObjectFilterFunc: t.WpObjectFilterFunc,
		WpIgnoreArrayOrder: t.WpIgnoreArrayOrder,
	}, nil
}

// ValueFromTerraform returns a Value given a tftypes.Value.  This is meant to convert the tftypes.Value into a more convenient Go type
// for the provider to consume the data with.
func (t NormalizedType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}
