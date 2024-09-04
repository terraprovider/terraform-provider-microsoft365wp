// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package generic

import (
	"context"
	"errors"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	tfdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// Lifted from https://github.com/hashicorp/terraform-provider-awscc/blob/v0.50.0/internal/tfresource/error.go
func DiagsError(diags tfdiag.Diagnostics) error {
	var errs *multierror.Error

	for _, diag := range diags {
		if diag == nil {
			continue
		}

		if diag.Severity() == tfdiag.SeverityError {
			errs = multierror.Append(errs, errors.New(diag.Detail()))
		}
	}

	return errs.ErrorOrNil()
}

// CopyValueAtPath copies the value at a specified path from source State to destination State.
// Lifted from https://github.com/hashicorp/terraform-provider-awscc/blob/v0.50.0/internal/generic/state.go
func CopyValueAtPath(ctx context.Context, dst, src *tfsdk.State, p path.Path) error {
	var val attr.Value
	diags := src.GetAttribute(ctx, p, &val)

	if diags.HasError() {
		return DiagsError(diags)
	}

	diags = dst.SetAttribute(ctx, p, val)

	if diags.HasError() {
		return DiagsError(diags)
	}

	return nil
}
