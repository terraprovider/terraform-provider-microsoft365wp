package generic

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// TestWriteSubActionAllInOneSkipIfUnchangedInGraph checks that a setPriority-like action is skipped if MS Graph already
// has the value, and sent otherwise (also if the entity cannot be read).
func TestWriteSubActionAllInOneSkipIfUnchangedInGraph(t *testing.T) {
	for _, tc := range []struct {
		name       string
		skip       bool
		planned    any
		getStatus  int
		wantAction bool
		wantWarn   bool
	}{
		{"unchanged: skipped", true, float64(3), http.StatusOK, false, false},
		{"changed: sent", true, float64(4), http.StatusOK, true, false},
		{"read fails: sent with warning", true, float64(3), http.StatusInternalServerError, true, true},
		{"option off: always sent", false, float64(3), http.StatusOK, true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var mu sync.Mutex
			var requests []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				requests = append(requests, req.Method+" "+req.URL.Path)
				switch {
				case req.Method == http.MethodGet:
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(tc.getStatus)
					if tc.getStatus == http.StatusOK {
						_, _ = io.WriteString(w, `{"id":"c1","displayName":"x","priority":3}`)
					} else {
						_, _ = io.WriteString(w, `{"error":{"code":"InternalServerError","message":"x"}}`)
					}
				case req.Method == http.MethodPost && req.URL.Path == "/beta/deviceManagement/deviceEnrollmentConfigurations/c1/setPriority":
					w.WriteHeader(http.StatusNoContent)
				default:
					t.Errorf("unexpected request: %s %s", req.Method, req.URL)
					w.WriteHeader(http.StatusBadRequest)
				}
			}))
			defer server.Close()
			client := msgraph.NewClient(msgraph.VersionBeta)
			client.Endpoint = server.URL
			client.RetryableClient.RetryMax = 0
			client.RetryableClient.RetryWaitMin = time.Millisecond
			client.RetryableClient.RetryWaitMax = 2 * time.Millisecond

			r := &GenericResource{AccessParams: AccessParams{BaseUri: "/deviceManagement/deviceEnrollmentConfigurations", graphClient: &client}}
			r.AccessParams.EntityId.AttrNameGraph = "id"
			a := &WriteSubActionAllInOne{
				WriteSubActionBase:     WriteSubActionBase{Attributes: []string{"priority"}, UriSuffix: "setPriority"},
				GraphErrorsAsWarnings:  true,
				SkipIfUnchangedInGraph: tc.skip,
			}
			a.Initialize()
			diags := diag.Diagnostics{}
			wsaReq := &WriteSubActionRequest{GenRes: r, Operation: OperationUpdate, Id: "c1",
				SubActionsData: map[string]any{"priority": tc.planned}}
			a.ExecutePost(context.Background(), &diags, wsaReq)

			if diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}
			if got := diags.WarningsCount() > 0; got != tc.wantWarn {
				t.Errorf("warnings: got %v, want %v (%v)", got, tc.wantWarn, diags)
			}
			mu.Lock()
			defer mu.Unlock()
			sent := false
			for _, r := range requests {
				if r == "POST /beta/deviceManagement/deviceEnrollmentConfigurations/c1/setPriority" {
					sent = true
				}
			}
			if sent != tc.wantAction {
				t.Errorf("setPriority sent: got %v, want %v (requests %v)", sent, tc.wantAction, requests)
			}
		})
	}
}
