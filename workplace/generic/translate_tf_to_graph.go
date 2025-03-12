// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//
// Lifted from https://github.com/hashicorp/terraform-provider-awscc/blob/v0.48.0/internal/generic/translate.go
// Changes:
//   - Return empty array for empty list or set instead of nil
//   - Do not try to expand JSONStringType
//   - Pull OData (sub-)type's subobject's attributes into root
//   - Automatic attribute name conversion (snake case vs. camel case)
//   - and probably some more...
//

package generic

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpjsontypes"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// TerraformAsRaw returns the raw map[string]any representing JSON from a Terraform Value.
func (t *ToFromGraphTranslator) TerraformAsRaw(ctx context.Context, val tftypes.Value) (map[string]any, error) {
	v, err := t.rawFromTerraformValue(ctx, nil, val)

	if err != nil {
		return nil, err
	}

	if v == nil {
		return make(map[string]any), nil
	}

	if v, ok := v.(map[string]any); ok {
		return v, nil
	}

	return nil, fmt.Errorf("unexpected raw type: %T", v)
}

// TerraformAsJsonString returns the string representing JSON from a Terraform Value.
func (t *ToFromGraphTranslator) TerraformAsJsonString(ctx context.Context, val tftypes.Value) (string, error) {
	v, err := t.TerraformAsRaw(ctx, val)

	if err != nil {
		return "", err
	}

	desiredState, err := json.Marshal(v)

	if err != nil {
		return "", err
	}

	return string(desiredState), nil
}

// rawFromTerraformValue returns the raw value (suitable for JSON marshaling) of the specified Terraform value.
// Terraform attribute names are mapped to JSON property names.
func (t *ToFromGraphTranslator) rawFromTerraformValue(ctx context.Context, path *tftypes.AttributePath, val tftypes.Value) (any, error) {

	// first check if we have any translated value for this attribute - if so, we're finished already
	graphValue, ok, err := t.TerraformToGraphTranslateValue(path, val)
	if err != nil || ok {
		return graphValue, err
	}

	if val.IsNull() || !val.IsKnown() {
		return nil, nil
	}

	attrType, err := t.SchemaTypeAtTerraformPath(path)
	if err != nil {
		return tftypes.Value{}, fmt.Errorf("getting attribute type at %s: %w", path, err)
	}

	typ := val.Type()
	switch {
	//
	// Primitive types.
	//
	case typ.Is(tftypes.Bool):
		var b bool
		if err := val.As(&b); err != nil {
			return nil, err
		}
		return b, nil

	case typ.Is(tftypes.Number):
		n := big.NewFloat(0)
		if err := val.As(&n); err != nil {
			return nil, err
		}
		f, _ := n.Float64()
		return f, nil

	case typ.Is(tftypes.String):
		var s string
		if err := val.As(&s); err != nil {
			return nil, err
		}
		if t := new(wpjsontypes.NormalizedType); t.Equal(attrType) {
			var v any
			err := json.Unmarshal([]byte(s), &v)
			return v, err
		}
		return s, nil

	//
	// Complex types.
	//
	case typ.Is(tftypes.List{}), typ.Is(tftypes.Set{}), typ.Is(tftypes.Tuple{}):
		var vals []tftypes.Value
		if err := val.As(&vals); err != nil {
			return nil, err
		}
		vs := make([]any, 0)
		for idx, val := range vals {
			if typ.Is(tftypes.Set{}) {
				// No need to worry about a specific value here.
				path = path.WithElementKeyValue(tftypes.NewValue(typ.(tftypes.Set).ElementType, nil))
			} else {
				path = path.WithElementKeyInt(idx)
			}
			v, err := t.rawFromTerraformValue(ctx, path, val)
			if err != nil {
				return nil, err
			}
			path = path.WithoutLastStep()
			if v == nil {
				continue
			}
			vs = append(vs, v)
		}
		return vs, nil

	case typ.Is(tftypes.Map{}), typ.Is(tftypes.Object{}):
		var vals map[string]tftypes.Value
		if err := val.As(&vals); err != nil {
			return nil, err
		}

		parentSchemaAttribute, err := t.SchemaAttributeAtTerraformPath(path)
		if err != nil {
			return tftypes.Value{}, err
		}

		vs := make(map[string]any)
		for name, val := range vals {
			if typ.Is(tftypes.Object{}) {
				path = path.WithAttributeName(name)
			} else {
				path = path.WithElementKeyString(name)
			}
			v, err := t.rawFromTerraformValue(ctx, path, val)
			if err != nil {
				return nil, err
			}
			path = path.WithoutLastStep()
			if v == nil && !t.IncludeNullObjectsInJson {
				continue
			}
			odataType, ok := t.OdataTypeByTerraformAttributeName(parentSchemaAttribute, name)
			if ok {
				if v != nil {
					leaf, ok := v.(map[string]any)
					if !ok {
						return nil, fmt.Errorf("OData sub-object for type %s does not satisfy 'map[string]any'", odataType)
					}
					for k, v2 := range leaf {
						vs[k] = v2
					}
					vs["@odata.type"] = odataType
				}
			} else if typ.Is(tftypes.Object{}) {
				propertyName, attributeIsWritable, ok := t.GraphAttributeNameFromTerraformName(parentSchemaAttribute, name)
				if !ok {
					return nil, fmt.Errorf("attribute name mapping not found: %s", name)
				}
				if attributeIsWritable { // some attributes are ok but still not supposed to be included in JSON (e.g. read-only)
					vs[propertyName] = v
				}
			} else {
				vs[name] = v // map, no translation
			}
		}
		return vs, nil
	}

	return nil, fmt.Errorf("unsupported value type: %s", typ)
}
