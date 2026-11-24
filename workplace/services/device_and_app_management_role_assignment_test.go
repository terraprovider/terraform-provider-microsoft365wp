package services

import (
	"context"
	"testing"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// The roleDefinition@odata.bind URL must use the configured endpoint, not graph.microsoft.com (national clouds, #238).
func TestDeviceAndAppManagementRoleAssignmentBindUsesEndpoint(t *testing.T) {
	client := msgraph.NewClient(msgraph.VersionBeta)
	client.Endpoint = "https://graph.microsoft.us"

	params := generic.TerraformToGraphMiddlewareParams{
		RawVal:      map[string]any{"displayName": "x", "roleDefinitionId": "rd-1"},
		GraphClient: &client,
	}
	diags := diag.Diagnostics{}
	if err := deviceAndAppManagementRoleAssignmentTerraformToGraphMiddleware(context.Background(), &diags, &params); err != nil {
		t.Fatalf("middleware: %v", err)
	}

	want := "https://graph.microsoft.us/beta/deviceManagement/roleDefinitions('rd-1')"
	if got := params.RawVal["roleDefinition@odata.bind"]; got != want {
		t.Errorf("roleDefinition@odata.bind = %v, want %s", got, want)
	}
	if _, ok := params.RawVal["roleDefinitionId"]; ok {
		t.Errorf("roleDefinitionId must not be sent to MS Graph")
	}
}
