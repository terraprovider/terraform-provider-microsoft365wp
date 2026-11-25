package generic

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

const unknownReferenceBody = `{"error":{"code":"BadRequest","message":"The server could not process the request because it is malformed or incorrect.","innerError":{"message":"1040: NamedLocation with id l does not exist in the directory."}}}`

func retryUnknownReference(resp *http.Response, o *odata.OData) bool {
	return resp.StatusCode == http.StatusBadRequest && o != nil && o.Error != nil && strings.Contains(o.Error.String(), "1040:")
}

// TestWriteRetryFunc checks that CreateRaw2 and UpdateRaw2 retry the responses WriteOptions.WriteRetryFunc accepts
// (for WriteOptions.WriteRetryFor), and only those.
func TestWriteRetryFunc(t *testing.T) {
	withWriteRetry := WriteOptions{WriteRetryFor: 300 * time.Millisecond, WriteRetryFunc: retryUnknownReference}

	for _, tc := range []struct {
		name      string
		options   WriteOptions
		method    string
		failures  int // number of 400 responses before the write succeeds
		wantCalls int // -2: at least 2
		wantError string
	}{
		{"create: retried", withWriteRetry, http.MethodPost, 2, 3, ""},
		{"create: without option", WriteOptions{}, http.MethodPost, 2, 1, "1040"},
		{"create: fails after the window", withWriteRetry, http.MethodPost, 1000, -2, "1040"},
		{"update: retried", withWriteRetry, http.MethodPatch, 2, 3, ""},
		{"update: without option", WriteOptions{}, http.MethodPatch, 2, 1, "1040"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var mu sync.Mutex
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				if req.Method != tc.method {
					t.Errorf("unexpected request: %s %s", req.Method, req.URL)
				}
				calls++
				w.Header().Set("Content-Type", "application/json")
				switch {
				case calls <= tc.failures:
					w.WriteHeader(http.StatusBadRequest)
					_, _ = io.WriteString(w, unknownReferenceBody)
				case req.Method == http.MethodPost:
					w.WriteHeader(http.StatusCreated)
					_, _ = io.WriteString(w, `{"id":"p"}`)
				default:
					w.Header().Del("Content-Type")
					w.WriteHeader(http.StatusNoContent)
				}
			}))
			defer server.Close()
			client := msgraph.NewClient(msgraph.VersionBeta)
			client.Endpoint = server.URL
			client.RetryableClient.RetryWaitMin = time.Millisecond
			client.RetryableClient.RetryWaitMax = 5 * time.Millisecond

			aps := &AccessParams{BaseUri: "/policies", WriteOptions: tc.options, graphClient: &client}
			aps.EntityId.AttrNameGraph = "id"
			diags := diag.Diagnostics{}
			if tc.method == http.MethodPost {
				aps.CreateRaw2(context.Background(), &diags, "", nil, map[string]any{"displayName": "x"}, false, false)
			} else {
				aps.UpdateRaw2(context.Background(), &diags, "", "p", nil, map[string]any{"displayName": "x"}, false)
			}

			if tc.wantError == "" && diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}
			if tc.wantError != "" && (!diags.HasError() || !strings.Contains(diagsString(diags), tc.wantError)) {
				t.Fatalf("want error containing %q, got %v", tc.wantError, diags)
			}
			mu.Lock()
			defer mu.Unlock()
			if tc.wantCalls == -2 && calls < 2 {
				t.Errorf("got %d calls, want at least 2", calls)
			}
			if tc.wantCalls >= 0 && calls != tc.wantCalls {
				t.Errorf("got %d calls, want %d", calls, tc.wantCalls)
			}
		})
	}
}
