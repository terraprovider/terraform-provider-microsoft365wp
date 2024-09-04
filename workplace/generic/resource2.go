// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package generic

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// populateUnknownValues populates and unknown values in State with values from the current resource description.
// Lifted from https://github.com/hashicorp/terraform-provider-awscc/blob/v0.50.0/internal/generic/resource.go
func (r *GenericResource) populateUnknownValues(ctx context.Context, state *tfsdk.State, tfSource tftypes.Value) diag.Diagnostics {
	var diags diag.Diagnostics

	unknowns, err := UnknownValuePaths(ctx, state.Raw)
	if err != nil {
		diags.AddError("Creation Of Terraform State Unsuccessful", fmt.Sprintf("Unable to set Terraform State Unknown values from REST API. Original Error: %s", err.Error()))
		return diags
	}

	if len(unknowns) == 0 {
		// no paths to unknown values, so nothing more to do
		return nil
	}

	if tfSource.IsNull() {
		readOptions := ReadOptions{
			ODataExpand:              r.AccessParams.ReadOptions.ODataExpand,
			ODataFilter:              r.AccessParams.ReadOptions.ODataFilter,
			SingleItemUseODataFilter: r.AccessParams.ReadOptions.SingleItemUseODataFilter,
			// leave out ExtraRequests, ExtraRequestsCustom
		}
		tfSource = r.AccessParams.ReadSingleRaw(ctx, &diags, state.Schema, readOptions, "", state, true)
		if diags.HasError() {
			return diags
		}
	}

	err = SetUnknownValuesFromResourceModel(ctx, state, unknowns, tfSource)
	if err != nil {
		diags.AddError("Creation Of Terraform State Unsuccessful", fmt.Sprintf("Unable to set Terraform State Unknown values from REST API. Original Error: %s", err.Error()))
		return diags
	}

	return nil
}
