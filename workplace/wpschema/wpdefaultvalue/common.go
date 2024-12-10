// Process: If the config value is null (i.e. no value provided in config) then return the hard-coded default value
// provided by the schema instead.
//
// Why we continue to use our own implementation even though the Framework supports default values natively by now: The
// Framework default value algorithm strictly follows the config and therefore does not work recursively on object
// attributes when the object itself had been set as a default value. Example: If the default for some_attr is an empty
// map/object, then the framework would not visit some_attr.some_prop and therefore it would be impossible to also set a
// default value for a child property.
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
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
// defaultValuePlanModify
//

func defaultValuePlanModify[T attr.Value](ctx context.Context, diags *diag.Diagnostics, schema rsschema.Schema, path path.Path, configValue T, defaultValueRaw any, responsePlanValue *T) {

	if configValue.IsNull() {

		valueFw := getFwValueFromAny(ctx, diags, &schema, path, defaultValueRaw)
		if diags.HasError() {
			return
		}

		valueFwTyped, ok := valueFw.(T)
		if !ok {
			diags.AddError("Error setting default value", fmt.Sprintf("defaultValuePlanModify: valueFw is not %T but %T", *new(T), valueFw))
		}

		*responsePlanValue = valueFwTyped
	}
}

//
// getFwValueFromAny
//

func getFwValueFromAny(ctx context.Context, diags *diag.Diagnostics, schema *rsschema.Schema, path path.Path, valueAny any) attr.Value {

	typeFw, diags2 := schema.TypeAtPath(ctx, path)
	diags.Append(diags2...)
	if diags.HasError() {
		return nil
	}

	valueTf, err := getFwValueFromAnyInner(ctx, schema, path, valueAny)
	if err != nil {
		diags.AddError("getFwValueFromAny", fmt.Sprintf("Unable to get valueTf from raw for %s: %s", valueAny, err.Error()))
		return nil
	}

	valueFw, err := typeFw.ValueFromTerraform(ctx, valueTf)
	if err != nil {
		diags.AddError("getFwValueFromAny", fmt.Sprintf("Unable to get valueFw from valueTf for %s: %s", valueAny, err.Error()))
		return nil
	}

	return valueFw
}

func getFwValueFromAnyInner(ctx context.Context, schema *rsschema.Schema, path path.Path, valueAny any) (tftypes.Value, error) {
	// As github.com/hashicorp/terraform-plugin-framework/internal/reflect can create objects only from structs but not
	// from map[string]any, we use a simplified form of ToFromGraphTranslator.TerraformFromRaw here

	var err error

	typeFw, diags := schema.TypeAtPath(ctx, path)
	if diags.HasError() {
		errDiag := diags.Errors()[0]
		return tftypes.Value{}, errors.New(errDiag.Summary() + ": " + errDiag.Detail())
	}
	typeTf := typeFw.TerraformType(ctx)

	switch valueTyped := valueAny.(type) {

	// primitive types
	case nil, bool, uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64, string:
		return tftypes.NewValue(typeTf, valueTyped), nil

	// slice
	case []any:

		// first case statement will not match for nil slices - for whatever reason
		if valueTyped == nil {
			return tftypes.NewValue(typeTf, nil), nil
		}

		// No need to worry about a specific value here to step into a set
		var setDummyKeyFw attr.Value
		if typeFwSet, ok := typeFw.(types.SetType); ok {
			// seems to be the most simple way to create a nil value (framework types do not seem to have a common NewValue function)
			setDummyKeyFw, err = typeFwSet.ElementType().ValueFromTerraform(ctx, tftypes.NewValue(typeFwSet.ElementType().TerraformType(ctx), nil))
			if err != nil {
				return tftypes.Value{}, err
			}
		}

		var valuesTf []tftypes.Value
		for i, v := range valueTyped {

			if _, ok := typeFw.(types.SetType); ok {
				path = path.AtSetValue(setDummyKeyFw)
			} else {
				path = path.AtListIndex(i)
			}

			valueTf, err := getFwValueFromAnyInner(ctx, schema, path, v)
			if err != nil {
				return tftypes.Value{}, err
			}
			valuesTf = append(valuesTf, valueTf)

			path = path.ParentPath()
		}

		return tftypes.NewValue(typeTf, valuesTf), nil

	// map
	case map[string]any:

		// first case statement will not match for nil maps - for whatever reason
		if valueTyped == nil {
			return tftypes.NewValue(typeTf, nil), nil
		}

		vals := make(map[string]tftypes.Value)
		for k, v := range valueTyped {

			if _, ok := typeFw.(types.ObjectType); ok {
				path = path.AtName(k)
			} else {
				path = path.AtMapKey(k)
			}

			val, err := getFwValueFromAnyInner(ctx, schema, path, v)
			if err != nil {
				return tftypes.Value{}, err
			}
			vals[k] = val

			path = path.ParentPath()
		}

		if typeTfObject, ok := typeTf.(tftypes.Object); ok {

			// Set any missing attributes to Null or Unknown (for computed values) as it is done in the framework as well:
			// https://github.com/hashicorp/terraform-plugin-framework/blob/v1.13.0/internal/fwserver/server_planresourcechange.go#L400

			for k, at := range typeTfObject.AttributeTypes {

				if _, ok := vals[k]; !ok {

					attribute, diags := schema.AttributeAtPath(ctx, path.AtName(k))
					if diags.HasError() {
						errDiag := diags.Errors()[0]
						return tftypes.Value{}, errors.New(errDiag.Summary() + ": " + errDiag.Detail())
					}

					var nullOrUnknown any
					if attribute.IsComputed() {
						nullOrUnknown = tftypes.UnknownValue
					} else {
						nullOrUnknown = nil
					}
					vals[k] = tftypes.NewValue(at, nullOrUnknown)
				}
			}
		}

		return tftypes.NewValue(typeTf, vals), nil

	default:
		return tftypes.Value{}, fmt.Errorf("unsupported raw type at %s: %T", path, valueTyped)
	}
}
