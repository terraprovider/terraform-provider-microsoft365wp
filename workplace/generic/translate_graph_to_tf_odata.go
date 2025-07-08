package generic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func ConvertOdataJsonToRaw(ctx context.Context, diags *diag.Diagnostics, jsonVal []byte) map[string]any {

	var rawVal map[string]any
	if err := json.Unmarshal(jsonVal, &rawVal); err != nil {
		diags.AddError(
			"Creation Of Terraform State Unsuccessful",
			fmt.Sprintf("Unable to create a raw value from the JSON returned by the REST API. Original Error: %s", err.Error()),
		)
		return nil
	}

	return rawVal
}

func ConvertOdataRawToTerraform(ctx context.Context, diags *diag.Diagnostics, schema tftypes.AttributePathStepper,
	rawVal map[string]any, valueToTargetSetName string, middlewareFunc GraphToTerraformMiddlewareFunc,
	middlewareExpectedId string, middlewareFuncTargetSetRunOnRawVal bool) tftypes.Value {

	runMiddleware := func(item map[string]any, isTargetSetOnRawVal bool) bool {
		params := GraphToTerraformMiddlewareParams{
			ExpectedId:          middlewareExpectedId,
			RawVal:              item,
			IsTargetSetOnRawVal: isTargetSetOnRawVal,
		}
		if err := middlewareFunc(ctx, diags, &params); err != nil {
			diags.AddError(
				"Creation Of Terraform State Unsuccessful",
				fmt.Sprintf("GraphToTerraformMiddleware returned an error: %s", err.Error()),
			)
		}
		if diags.HasError() {
			return false
		}

		return true
	}

	if valueToTargetSetName == "" {
		// In rare cases we might actually receive a value collection with zero or just one item (e.g. if there does not
		// exist a direct route an item and instead an OData $filter is required to target it)
		// Added "len(rawVal) <= 3" as an additional sanity check to not trigger on any other entity that might have
		// a "value" attribute containing an array
		if items, ok := rawVal["value"].([]any); ok && len(rawVal) <= 3 {
			singleValueIsValid := false
			if len(items) == 1 {
				if singleValue, ok := items[0].(map[string]any); ok {
					rawVal = singleValue
					singleValueIsValid = true
				}
			}
			if !singleValueIsValid {
				return tftypes.Value{} // gets translated to "not found" upstream
			}
		}
		if middlewareFunc != nil {
			if ok := runMiddleware(rawVal, false); !ok {
				return tftypes.Value{}
			}
			if len(rawVal) == 0 {
				// a completely empty map signals "not found"
				return tftypes.Value{}
			}
		}
	} else {
		if middlewareFunc != nil && middlewareFuncTargetSetRunOnRawVal {
			if ok := runMiddleware(rawVal, true); !ok {
				return tftypes.Value{}
			}
		}
		items, ok := rawVal["value"].([]any)
		if !ok {
			diags.AddError("Creation Of Terraform State Unsuccessful", "JSON property 'value' not found in OData result or not an array.")
			return tftypes.Value{}
		}
		if middlewareFunc != nil && !middlewareFuncTargetSetRunOnRawVal {
			for _, x := range items {
				if ok := runMiddleware(x.(map[string]any), true); !ok {
					return tftypes.Value{}
				}
			}
		}
		// "rename" key from "value" to valueTotargetSetName
		delete(rawVal, "value")
		rawVal[valueToTargetSetName] = items
	}

	translator := NewToFromGraphTranslator(schema, false)
	tfVal, err := translator.TerraformFromRaw(ctx, rawVal)
	if err != nil {
		diags.AddError(
			"Creation Of Terraform State Unsuccessful",
			fmt.Sprintf("Unable to create a Terraform State value from OData Raw value. Original Error: %s", err.Error()),
		)
		return tftypes.Value{}
	}

	return tfVal
}

func ConvertOdataJsonToTerraform(ctx context.Context, diags *diag.Diagnostics, schema tftypes.AttributePathStepper,
	jsonVal []byte, valueToTargetSetName string, middlewareFunc GraphToTerraformMiddlewareFunc,
	middlewareExpectedId string, middlewareFuncTargetSetRunOnRawVal bool) tftypes.Value {

	rawVal := ConvertOdataJsonToRaw(ctx, diags, jsonVal)
	if diags.HasError() {
		return tftypes.Value{}
	}

	tfVal := ConvertOdataRawToTerraform(ctx, diags, schema, rawVal, valueToTargetSetName, middlewareFunc, middlewareExpectedId, middlewareFuncTargetSetRunOnRawVal)
	if diags.HasError() {
		return tftypes.Value{}
	}

	return tfVal
}
