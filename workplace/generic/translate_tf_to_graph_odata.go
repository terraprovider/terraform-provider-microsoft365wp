package generic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func ConvertTerraformToOdataJson(ctx context.Context, diags *diag.Diagnostics, schema tftypes.AttributePathStepper, includeNullObjects bool, val tftypes.Value,
	isUpdate bool, middlewareFunc func(TerraformToGraphMiddlewareParams) TerraformToGraphMiddlewareReturns, valSource string) []byte {

	rawVal := ConvertTerraformToOdataRaw(ctx, diags, schema, includeNullObjects, val, isUpdate, middlewareFunc, valSource)
	if diags.HasError() {
		return nil
	}

	jsonVal := ConvertOdataRawToJson(ctx, diags, rawVal, valSource)
	if diags.HasError() {
		return nil
	}

	return jsonVal
}

func ConvertTerraformToOdataRaw(ctx context.Context, diags *diag.Diagnostics, schema tftypes.AttributePathStepper, includeNullObjects bool, val tftypes.Value,
	isUpdate bool, middlewareFunc func(TerraformToGraphMiddlewareParams) TerraformToGraphMiddlewareReturns, valSource string) map[string]any {

	translator := NewToFromGraphTranslator(schema, includeNullObjects)
	rawVal, err := translator.TerraformAsRaw(ctx, val)

	if err != nil {
		diags.AddError(
			"Creation Of OData Raw Unsuccessful",
			fmt.Sprintf("Unable to create a raw value from Terraform %s. Original Error: %s", valSource, err.Error()),
		)
		return nil
	}

	if middlewareFunc != nil {
		params := TerraformToGraphMiddlewareParams{
			Ctx:      ctx,
			RawVal:   rawVal,
			IsUpdate: isUpdate,
		}
		if err := middlewareFunc(params); err != nil {
			diags.AddError(
				"Creation Of OData Raw Unsuccessful",
				fmt.Sprintf("TerraformToGraphMiddleware returned an error: %s", err.Error()),
			)
			return nil
		}
	}

	return rawVal
}

func ConvertOdataRawToJson(ctx context.Context, diags *diag.Diagnostics, raw map[string]any, valSource string) []byte {

	jsonVal, err := json.Marshal(raw)
	if err != nil {
		diags.AddError(
			"Creation Of OData JSON Unsuccessful",
			fmt.Sprintf("Unable to create OData JSON from Terraform %s. Original Error: %s", valSource, err.Error()),
		)
		return nil
	}

	return jsonVal
}
