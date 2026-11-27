package services

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// TestDeviceEnrollmentConfigurationEspEmptyAppsWsa checks that an empty selectedMobileAppIds is written separately
// only together with allowNonBlockingAppInstallation = true (MS Graph rejects that combination in one PATCH).
func TestDeviceEnrollmentConfigurationEspEmptyAppsWsa(t *testing.T) {
	for _, tc := range []struct {
		name          string
		rawVal        map[string]any
		wantSeparate  bool
		wantInRawBody bool
	}{
		{"non-blocking with empty apps", map[string]any{"@odata.type": kEspOdataType, "allowNonBlockingAppInstallation": true, "selectedMobileAppIds": []any{}}, true, false},
		{"non-blocking with apps", map[string]any{"@odata.type": kEspOdataType, "allowNonBlockingAppInstallation": true, "selectedMobileAppIds": []any{"a1"}}, false, true},
		{"blocking with empty apps", map[string]any{"@odata.type": kEspOdataType, "allowNonBlockingAppInstallation": false, "selectedMobileAppIds": []any{}}, false, true},
		{"other type", map[string]any{"@odata.type": "#microsoft.graph.deviceEnrollmentLimitConfiguration", "allowNonBlockingAppInstallation": true, "selectedMobileAppIds": []any{}}, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var patches []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.Method != http.MethodPatch || req.URL.Path != "/beta/deviceManagement/deviceEnrollmentConfigurations/e1" {
					t.Errorf("unexpected request: %s %s", req.Method, req.URL)
				}
				b, _ := io.ReadAll(req.Body)
				patches = append(patches, string(b))
				w.WriteHeader(http.StatusNoContent)
			}))
			defer server.Close()
			client := msgraph.NewClient(msgraph.VersionBeta)
			client.Endpoint = server.URL
			client.RetryableClient.RetryWaitMin = time.Millisecond
			client.RetryableClient.RetryWaitMax = 2 * time.Millisecond
			r := &generic.GenericResource{AccessParams: generic.AccessParams{BaseUri: DeviceEnrollmentConfigurationResource.AccessParams.BaseUri}}
			r.AccessParams.InitializeGuarded(&client)

			a := &deviceEnrollmentConfigurationEspEmptyAppsWsa{}
			a.Initialize()
			diags := diag.Diagnostics{}
			wsaReq := &generic.WriteSubActionRequest{GenRes: r, Operation: generic.OperationUpdate, Id: "e1", RawVal: tc.rawVal, SubActionsData: map[string]any{}}
			a.ExecutePre(context.Background(), &diags, wsaReq)
			_, inBody := tc.rawVal["selectedMobileAppIds"]
			if inBody != tc.wantInRawBody {
				t.Errorf("selectedMobileAppIds in the entity write: got %v, want %v", inBody, tc.wantInRawBody)
			}
			a.ExecutePost(context.Background(), &diags, wsaReq)
			if diags.HasError() {
				t.Fatal(diags)
			}
			var want []string
			if tc.wantSeparate {
				want = []string{`{"@odata.type":"` + kEspOdataType + `","selectedMobileAppIds":[]}`}
			}
			if !reflect.DeepEqual(patches, want) {
				t.Errorf("separate PATCHes: got %v, want %v", patches, want)
			}
		})
	}
}
