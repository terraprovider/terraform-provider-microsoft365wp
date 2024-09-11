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
	rawVal map[string]any, valueToTargetSetName string, middlewareFunc func(GraphToTerraformMiddlewareParams) GraphToTerraformMiddlewareReturns) tftypes.Value {

	runMiddleware := func(item map[string]any) bool {
		params := GraphToTerraformMiddlewareParams{
			Ctx:    ctx,
			RawVal: item,
		}
		if err := middlewareFunc(params); err != nil {
			diags.AddError(
				"Creation Of Terraform State Unsuccessful",
				fmt.Sprintf("GraphToTerraformMiddleware returned an error: %s", err.Error()),
			)
			return false
		}

		return true
	}

	if valueToTargetSetName == "" {
		// In rare cases we might actually receive a value collection with zero or just one item (e.g. if there does not
		// exist a direct route an item and instead an OData $filter is required to target it)
		// Added "len(rawVal) <= 5" as an additional sanity check to not trigger on any other entity that might have
		// a "value" attribute containing an array
		if items, ok := rawVal["value"].([]any); ok && len(rawVal) <= 5 {
			switch len(items) {
			case 0:
				return tftypes.Value{} // gets translated to "not found" upstream
			case 1:
				if singleValue, ok := items[0].(map[string]any); ok {
					rawVal = singleValue
				}
			}
		}
		if middlewareFunc != nil {
			if ok := runMiddleware(rawVal); !ok {
				return tftypes.Value{}
			}
		}
	} else {
		items := rawVal["value"].([]any)
		if middlewareFunc != nil {
			for _, x := range items {
				if ok := runMiddleware(x.(map[string]any)); !ok {
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

func ConvertOdataJsonToTerraform(ctx context.Context, diags *diag.Diagnostics, schema tftypes.AttributePathStepper, jsonVal []byte,
	valueToTargetSetName string, middlewareFunc func(GraphToTerraformMiddlewareParams) GraphToTerraformMiddlewareReturns) tftypes.Value {

	rawVal := ConvertOdataJsonToRaw(ctx, diags, jsonVal)
	if diags.HasError() {
		return tftypes.Value{}
	}

	tfVal := ConvertOdataRawToTerraform(ctx, diags, schema, rawVal, valueToTargetSetName, middlewareFunc)
	if diags.HasError() {
		return tftypes.Value{}
	}

	return tfVal
}
