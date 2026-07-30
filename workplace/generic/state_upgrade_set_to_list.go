package generic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// SetToListStateUpgrader is a reusable StateUpgrader function that transparently
// migrates prior state from Set-typed to List-typed attributes.
//
// It works because Terraform persists state as JSON, which has no Set/List
// distinction (both are JSON arrays). By unmarshalling the prior RawState JSON
// and re-interpreting it through the current schema (where affected attributes
// are now List instead of Set), the values are naturally produced with the
// correct List type. The actual data shape is unchanged.
func SetToListStateUpgrader(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	if req.RawState == nil {
		resp.Diagnostics.AddError("Upgrade State Error", "RawState is nil; cannot upgrade state")
		return
	}

	rawJSON, err := req.RawState.Unmarshal(resp.State.Schema.Type().TerraformType(ctx))
	if err != nil {
		// If unmarshalling with the new schema type fails (because Set→List type mismatch),
		// fall back to unmarshalling as plain JSON and re-encoding through the translator.
		tflog.Debug(ctx, fmt.Sprintf("SetToListStateUpgrader: Unmarshal with new schema type failed (%s), falling back to JSON re-parse", err))

		if req.RawState.JSON == nil {
			resp.Diagnostics.AddError("Upgrade State Error",
				"Prior state has no JSON data; cannot upgrade state from flatmap format")
			return
		}

		var rawVal map[string]any
		if jsonErr := json.Unmarshal(req.RawState.JSON, &rawVal); jsonErr != nil {
			resp.Diagnostics.AddError("Upgrade State Error",
				fmt.Sprintf("Unable to unmarshal prior state JSON: %s", jsonErr))
			return
		}

		translator := NewToFromGraphTranslator(resp.State.Schema, false, req.State)
		tfVal, translateErr := translator.TerraformFromRaw(ctx, rawVal)
		if translateErr != nil {
			resp.Diagnostics.AddError("Upgrade State Error",
				fmt.Sprintf("Unable to translate prior state to new schema: %s", translateErr))
			return
		}

		resp.State = tfsdk.State{
			Schema: resp.State.Schema,
			Raw:    tfVal,
		}
		return
	}

	// Unmarshal succeeded directly with the new schema type — just set the state.
	resp.State = tfsdk.State{
		Schema: resp.State.Schema,
		Raw:    rawJSON,
	}
}
