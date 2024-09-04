// Process: If the planned value is null (i.e. no value provided in config) then return the hard-coded default value
// provided by the schema instead.
//
// 2023-05-03 Felix Storm: Previously, we had also checked if there was no current plan value yet (to honor potential
// previous planModifier's result). But since TF seems to always copy the current state to the proposed plan if the
// resource already exists (before applying config values) this rendered the whole idea pretty much useless (as all
// config values not provided will have already been replaced by the corresponding values from the state before any plan
// modifier gets called). I am pretty sure that I had tested out that removing a value from the config will force it to
// become the default value again but cannot really understand now how this could have worked at all before...
//

package wpdefaultvalue

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

//
// defaultValueAttributePlanModifier
//

type defaultValueAttributePlanModifier struct{}

func (attributePlanModifier defaultValueAttributePlanModifier) Description(_ context.Context) string {
	return "If the value of the attribute is missing, then the value is semantically the same as if the value was present with the default value hard-coded in the provider."
}

func (attributePlanModifier defaultValueAttributePlanModifier) MarkdownDescription(ctx context.Context) string {
	return attributePlanModifier.Description(ctx)
}

//
// GetAttrValueFromRawWithType
//

func GetAttrValueFromRawWithType(ctx context.Context, diags *diag.Diagnostics, attrType attr.Type, rawValue any) attr.Value {
	tfTypeRoot := attrType.TerraformType(ctx)
	tfValue, err := getAttrValueFromRawWithType(ctx, &tfTypeRoot, nil, rawValue)
	if err != nil {
		diags.AddError("GetAttrValueFromRawWithType", fmt.Sprintf("Unable to get tfValue from raw for %s: %s", rawValue, err.Error()))
		return nil
	}

	attrValue, err := attrType.ValueFromTerraform(ctx, tfValue)
	if err != nil {
		diags.AddError("GetAttrValueFromRawWithType", fmt.Sprintf("Unable to get attrValue from tfValue for %s: %s", rawValue, err.Error()))
		return nil
	}

	return attrValue
}

func getAttrValueFromRawWithType(ctx context.Context, root *tftypes.Type, path *tftypes.AttributePath, v any) (tftypes.Value, error) {
	// As github.com/hashicorp/terraform-plugin-framework/internal/reflect can create objects only from structs but not
	// from map[string]any, we use a simplified form of ToFromGraphTranslator.TerraformFromRaw here

	rawType, remaining, err := tftypes.WalkAttributePath(*root, path)
	if err != nil {
		return tftypes.Value{}, fmt.Errorf("%v still remains in the path: %w", remaining, err)
	}
	typ, ok := rawType.(tftypes.Type)
	if !ok {
		return tftypes.Value{}, fmt.Errorf("rawType is not tftypes.Type but %T at %s", rawType, path)
	}

	switch v := v.(type) {

	// primitive types
	case nil, bool, uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64, string:
		return tftypes.NewValue(typ, v), nil

	// slice
	case []any:
		var vals []tftypes.Value
		for idx, v := range v {
			if typ.Is(tftypes.Set{}) {
				// No need to worry about a specific value here.
				path = path.WithElementKeyValue(tftypes.NewValue(typ.(tftypes.Set).ElementType, nil))
			} else {
				path = path.WithElementKeyInt(idx)
			}
			val, err := getAttrValueFromRawWithType(ctx, root, path, v)
			if err != nil {
				return tftypes.Value{}, err
			}
			vals = append(vals, val)
			path = path.WithoutLastStep()
		}
		return tftypes.NewValue(typ, vals), nil

	// map
	case map[string]any:

		vals := make(map[string]tftypes.Value)
		for k, v := range v {
			if typ.Is(tftypes.Object{}) {
				path = path.WithAttributeName(k)
			} else {
				path = path.WithElementKeyString(k)
			}
			val, err := getAttrValueFromRawWithType(ctx, root, path, v)
			if err != nil {
				return tftypes.Value{}, err
			}
			vals[k] = val
			path = path.WithoutLastStep()
		}

		if typ.Is(tftypes.Object{}) {
			// Set any missing attributes to unknown.
			for k, at := range typ.(tftypes.Object).AttributeTypes {
				if _, ok := vals[k]; !ok {
					vals[k] = tftypes.NewValue(at, tftypes.UnknownValue)
				}
			}
		}

		return tftypes.NewValue(typ, vals), nil

	default:
		return tftypes.Value{}, fmt.Errorf("unsupported raw type at %s: %T", path, v)
	}
}

//
// defaultValuePlanModify
//

func defaultValuePlanModify[T attr.Value](ctx context.Context, configValue T, responsePlanValue *T, diags *diag.Diagnostics, defaultValueRaw any) {
	if configValue.IsNull() {
		attrType := configValue.Type(ctx)
		attrValue := GetAttrValueFromRawWithType(ctx, diags, attrType, defaultValueRaw)
		if diags.HasError() {
			return
		}

		typedValue, ok := attrValue.(T)
		if !ok {
			diags.AddError("objectDefaultValueAttributePlanModifier.PlanModifyObject",
				fmt.Sprintf("attrValue is not %T but %T for %s", *new(T), attrValue, attrType))
		}

		*responsePlanValue = typedValue
	}
}
