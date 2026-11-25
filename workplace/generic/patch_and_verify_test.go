package generic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// patchTestServer simulates an entity with one value (authenticationContextId) that loses the first loseWrites
// PATCHes (204 without effect).
type patchTestServer struct {
	mu           sync.Mutex
	value        any
	loseWrites   int
	patchStatus  int // 0: 204
	readNotFound bool
	patches      int
	gets         int
}

func (s *patchTestServer) handler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		if req.URL.Path != "/beta/actions/a" {
			t.Errorf("unexpected request: %s %s", req.Method, req.URL)
		}
		switch req.Method {
		case http.MethodPatch:
			s.patches++
			b, _ := io.ReadAll(req.Body)
			var body map[string]any
			if err := json.Unmarshal(b, &body); err != nil {
				t.Errorf("PATCH body: %v", err)
			}
			if s.patchStatus != 0 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(s.patchStatus)
				_, _ = io.WriteString(w, `{"error":{"code":"Directory_ObjectNotFound","message":"Unable to read the company information from the directory."}}`)
				return
			}
			if s.patches > s.loseWrites {
				s.value = body["authenticationContextId"]
			}
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			s.gets++
			if req.URL.Query().Get("$select") != "id,authenticationContextId" {
				t.Errorf("GET without the expected $select: %s", req.URL)
			}
			w.Header().Set("Content-Type", "application/json")
			if s.readNotFound {
				w.WriteHeader(http.StatusNotFound)
				_, _ = io.WriteString(w, `{"error":{"code":"Request_ResourceNotFound","message":"not found"}}`)
				return
			}
			b, _ := json.Marshal(map[string]any{"id": "a", "authenticationContextId": s.value})
			_, _ = w.Write(b)
		default:
			t.Errorf("unexpected method %s", req.Method)
		}
	}
}

func TestPatchAndVerifyRaw(t *testing.T) {
	opts := PatchAndVerifyOptions{PollInterval: time.Millisecond, ResendAfter: 20 * time.Millisecond, ConsecutiveReads: 2, Timeout: 200 * time.Millisecond}

	for _, tc := range []struct {
		name         string
		initial      any
		want         any
		loseWrites   int
		patchStatus  int
		readNotFound bool
		wantPatches  int // -1: don't check, -2: at least 2
		wantError    string
	}{
		{"set, visible at once", nil, "c1", 0, 0, false, 1, ""},
		{"clear, visible at once", "c1", nil, 0, 0, false, 1, ""},
		{"lost write is sent again", nil, "c1", 1, 0, false, 2, ""},
		{"lost clear is sent again", "c1", nil, 2, 0, false, 3, ""},
		{"404 for the PATCH, but already at target", nil, nil, 0, http.StatusNotFound, false, 1, ""},
		{"404 for the PATCH and not at target", "c1", nil, 0, http.StatusNotFound, false, -2, "did not apply"},
		{"entity does not exist", nil, "c1", 0, 0, true, 1, "does not exist"},
		{"other PATCH errors fail at once", nil, "c1", 0, http.StatusBadRequest, false, 1, "Unable to update"},
		{"never applied", nil, "c1", 1000, 0, false, -2, "did not apply"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := &patchTestServer{value: tc.initial, loseWrites: tc.loseWrites, patchStatus: tc.patchStatus, readNotFound: tc.readNotFound}
			server := httptest.NewServer(s.handler(t))
			defer server.Close()
			client := msgraph.NewClient(msgraph.VersionBeta)
			client.Endpoint = server.URL
			client.DisableRetries = true

			readUri := msgraph.Uri{Entity: "/actions/a", Params: map[string][]string{"$select": {"id,authenticationContextId"}}}
			matches := func(got map[string]any) bool { return got["authenticationContextId"] == tc.want }
			diags := diag.Diagnostics{}
			PatchAndVerifyRaw(context.Background(), &diags, &client, msgraph.Uri{Entity: "/actions/a"},
				map[string]any{"authenticationContextId": tc.want}, readUri, matches, opts)

			if tc.wantError == "" && diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}
			if tc.wantError != "" && (!diags.HasError() || !strings.Contains(patchTestDiagsText(diags), tc.wantError)) {
				t.Fatalf("want error containing %q, got %v", tc.wantError, diags)
			}
			s.mu.Lock()
			defer s.mu.Unlock()
			switch {
			case tc.wantPatches == -2 && s.patches < 2:
				t.Errorf("got %d PATCHes, want at least 2", s.patches)
			case tc.wantPatches >= 0 && s.patches != tc.wantPatches:
				t.Errorf("got %d PATCHes, want %d", s.patches, tc.wantPatches)
			}
		})
	}
}

func patchTestDiagsText(diags diag.Diagnostics) string {
	var sb strings.Builder
	for _, d := range diags {
		fmt.Fprintf(&sb, "%s: %s\n", d.Summary(), d.Detail())
	}
	return sb.String()
}
