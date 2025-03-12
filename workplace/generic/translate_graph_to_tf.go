// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//
// Lifted from https://github.com/hashicorp/terraform-provider-awscc/blob/v0.48.0/internal/generic/translate.go
// Changes:
//   - Return empty list or set for empty array instead of null
//   - Push OData (sub-)type's attributes to a subobject
//   - Support nil value
//   - Automatic attribute name conversion (snake case vs. camel case)
//   - and probably some more...
//

package generic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// TerraformFromRaw returns the Terraform Value for the specified JSON Properties (raw map[string]any).
func (t *ToFromGraphTranslator) TerraformFromRaw(ctx context.Context, resourceModel map[string]any) (tftypes.Value, error) {
	return t.terraformValueFromRaw(ctx, nil, resourceModel)
}

// TerraformFromJsonString returns the Terraform Value for the specified JSON Properties (string).
func (t *ToFromGraphTranslator) TerraformFromJsonString(ctx context.Context, resourceModel string) (tftypes.Value, error) {
	var v any

	if err := json.Unmarshal([]byte(resourceModel), &v); err != nil {
		return tftypes.Value{}, err
	}

	if v, ok := v.(map[string]any); ok {
		return t.TerraformFromRaw(ctx, v)
	}

	return tftypes.Value{}, fmt.Errorf("unexpected raw type: %T", v)
}

func (t *ToFromGraphTranslator) terraformValueFromRaw(ctx context.Context, path *tftypes.AttributePath, v any) (tftypes.Value, error) {

	// first check if we have any translated value for this attribute - if so, we're finished already
	tfValue, ok, err := t.GraphToTerraformTranslateValue(path, v)
	if err != nil || ok {
		return tfValue, err
	}

	attrType, err := t.SchemaTypeAtTerraformPath(path)
	if err != nil {
		return tftypes.Value{}, fmt.Errorf("getting attribute type at %s: %w", path, err)
	}

	typ := attrType.TerraformType(ctx)

	// if target is string but source is complex then assume that the target string should receive JSON
	if typ.Is(tftypes.String) {
		switch v := v.(type) {
		case []any, map[string]any:
			val, err := json.Marshal(v)
			if err != nil {
				return tftypes.Value{}, err
			}

			return tftypes.NewValue(typ, string(val)), nil
		}
	}

	switch v := v.(type) {
	//
	// Primitive types.
	//
	case nil:
		// translation for existing but null values gets handled further up, so we just want a Terraform default here
		return t.GraphToTerraformDefaultValue(path, false, typ)

	case bool:
		return tftypes.NewValue(tftypes.Bool, v), nil

	case float64:
		return tftypes.NewValue(tftypes.Number, v), nil

	case string:
		return tftypes.NewValue(tftypes.String, v), nil

	//
	// Complex types.
	//
	case []any:

		var vals []tftypes.Value
		for idx, v := range v {
			if typ.Is(tftypes.Set{}) {
				// No need to worry about a specific value here.
				path = path.WithElementKeyValue(tftypes.NewValue(typ.(tftypes.Set).ElementType, nil))
			} else {
				path = path.WithElementKeyInt(idx)
			}
			val, err := t.terraformValueFromRaw(ctx, path, v)
			if err != nil {
				return tftypes.Value{}, err
			}
			vals = append(vals, val)
			path = path.WithoutLastStep()
		}
		return tftypes.NewValue(typ, vals), nil

	case map[string]any:

		vals := make(map[string]tftypes.Value)
		if typ.Is(tftypes.Object{}) {
			// TF target is "real" object

			parentSchemaAttribute, err := t.SchemaAttributeAtTerraformPath(path)
			if err != nil {
				return tftypes.Value{}, err
			}

			if odataType, ok := v["@odata.type"].(string); ok {
				leafAttrs, ok := t.NestedObjectAttributesByOdataType(parentSchemaAttribute, odataType)
				if ok {
					// move attributes belonging to specific (sub-)type to corresponding new leaf object
					leaf := make(map[string]any)
					for _, k := range leafAttrs {
						if v2, ok := v[k]; ok {
							leaf[k] = v2
							delete(v, k)
						}
					}
					v[odataType] = leaf
				}
			}

			for key, v := range v {
				tfAttributeName, ok := t.TerraformAttributeNameFromGraphName(parentSchemaAttribute, key)
				if !ok {
					tflog.Info(ctx, "not found in Terraform schema", map[string]any{"key": key, "path": path})
					continue
				}
				path = path.WithAttributeName(tfAttributeName)
				val, err := t.terraformValueFromRaw(ctx, path, v)
				if err != nil {
					return tftypes.Value{}, err
				}
				vals[tfAttributeName] = val
				path = path.WithoutLastStep()
			}

			// Set any missing attributes to Null.
			for k, at := range typ.(tftypes.Object).AttributeTypes {
				if _, ok := vals[k]; !ok {
					defaultValue, err := t.GraphToTerraformDefaultValue(path.WithAttributeName(k), true, at)
					if err != nil {
						return tftypes.Value{}, err
					}
					vals[k] = defaultValue
				}
			}

		} else {
			// TF target is map[string]

			for key, v := range v {
				path = path.WithElementKeyString(key)
				val, err := t.terraformValueFromRaw(ctx, path, v)
				if err != nil {
					return tftypes.Value{}, err
				}
				vals[key] = val
				path = path.WithoutLastStep()
			}

		}
		return tftypes.NewValue(typ, vals), nil

	default:
		return tftypes.Value{}, fmt.Errorf("unsupported raw type at %s: %T", path, v)
	}
}
