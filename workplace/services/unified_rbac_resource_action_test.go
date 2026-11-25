package services

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// TestUnifiedRbacResourceActionWriteAuthenticationContextId checks that only authenticationContextId is sent (MS Graph
// rejects all other properties), that "" clears it, and that the value is read back with $select (otherwise MS Graph
// does not return it).
func TestUnifiedRbacResourceActionWriteAuthenticationContextId(t *testing.T) {
	oldOpts := generic.DefaultPatchAndVerifyOptions
	generic.DefaultPatchAndVerifyOptions = generic.PatchAndVerifyOptions{PollInterval: time.Millisecond, ResendAfter: time.Second, ConsecutiveReads: 2, Timeout: time.Second}
	defer func() { generic.DefaultPatchAndVerifyOptions = oldOpts }()

	const path = "/beta/roleManagement/directory/resourceNamespaces/microsoft.directory/resourceActions/microsoft.directory-conditionalAccessPolicies-create-post"

	for _, tc := range []struct {
		name      string
		value     any
		wantBody  map[string]any
		wantValue any
	}{
		{"set", "c1", map[string]any{"authenticationContextId": "c1"}, "c1"},
		{"clear with nil", nil, map[string]any{"authenticationContextId": nil}, nil},
		{"clear with empty string", "", map[string]any{"authenticationContextId": nil}, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var mu sync.Mutex
			var value any = "c9"
			var bodies []map[string]any
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				if req.URL.Path != path {
					t.Errorf("unexpected request: %s %s", req.Method, req.URL)
				}
				switch req.Method {
				case http.MethodPatch:
					b, _ := io.ReadAll(req.Body)
					var body map[string]any
					_ = json.Unmarshal(b, &body)
					bodies = append(bodies, body)
					value = body["authenticationContextId"]
					w.WriteHeader(http.StatusNoContent)
				case http.MethodGet:
					if req.URL.Query().Get("$select") == "" {
						t.Errorf("GET without $select: %s", req.URL)
					}
					w.Header().Set("Content-Type", "application/json")
					b, _ := json.Marshal(map[string]any{"id": "x", "authenticationContextId": value})
					_, _ = w.Write(b)
				}
			}))
			defer server.Close()
			client := msgraph.NewClient(msgraph.VersionBeta)
			client.Endpoint = server.URL

			diags := diag.Diagnostics{}
			uri := msgraph.Uri{Entity: UnifiedRbacResourceActionResource.AccessParams.BaseUri + "/microsoft.directory-conditionalAccessPolicies-create-post"}
			unifiedRbacResourceActionWriteAuthenticationContextId(context.Background(), &diags, &client, uri, tc.value)
			if diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}

			mu.Lock()
			defer mu.Unlock()
			if len(bodies) != 1 || !reflect.DeepEqual(bodies[0], tc.wantBody) {
				t.Errorf("PATCH bodies %v, want exactly one %v", bodies, tc.wantBody)
			}
			if value != tc.wantValue {
				t.Errorf("value %v, want %v", value, tc.wantValue)
			}
		})
	}
}
