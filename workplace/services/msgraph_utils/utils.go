package msgraphutils

import (
	"context"
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func SingletonSyntheticIdGraphToTerraformMiddleware(_ context.Context, _ *diag.Diagnostics, params *generic.GraphToTerraformMiddlewareParams) generic.GraphToTerraformMiddlewareReturns {
	// MS Graph does indeed sometimes not return anything for id, so we add a synthetic one
	if params.RawVal["id"] == nil {
		params.RawVal["id"] = "singleton_synthetic_id"
	}
	return nil
}
