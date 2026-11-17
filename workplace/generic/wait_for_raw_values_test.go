package generic

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func mustJson(t *testing.T, s string) any {
	t.Helper()
	var v any
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		t.Fatalf("json: %v", err)
	}
	return v
}

func TestRawContains(t *testing.T) {
	tests := []struct {
		name      string
		want, got string
		result    bool
	}{
		{"equal", `{"a":true,"b":"x"}`, `{"a":true,"b":"x"}`, true},
		{"extra keys in got", `{"a":true}`, `{"a":true,"b":1,"@odata.context":"c"}`, true},
		{"missing key in got", `{"a":true,"b":1}`, `{"a":true}`, false},
		{"stale value", `{"inboundTrust":{"isMfaAccepted":true}}`, `{"inboundTrust":{"isMfaAccepted":false}}`, false},
		{"nested subset", `{"inboundTrust":{"isMfaAccepted":true}}`, `{"inboundTrust":{"isMfaAccepted":true,"x":null}}`, true},
		{"null in want matches anything", `{"a":null}`, `{"a":{"b":1}}`, true},
		{"null in want matches missing key", `{"a":null}`, `{}`, true},
		{"empty list in want matches null", `{"a":[]}`, `{"a":null}`, true},
		{"empty object in want matches null", `{"a":{}}`, `{"a":null}`, true},
		{"string case-insensitive", `{"a":"AllUsers"}`, `{"a":"allusers"}`, true},
		{"list any order", `{"t":[{"target":"A"},{"target":"B"}]}`, `{"t":[{"target":"B","x":1},{"target":"A"}]}`, true},
		{"list different length", `{"t":[{"target":"A"}]}`, `{"t":[{"target":"A"},{"target":"B"}]}`, false},
		{"list duplicates need distinct matches", `{"t":["A","A"]}`, `{"t":["A","B"]}`, false},
		{"odata type ignored", `{"@odata.type":"#x","a":1}`, `{"a":1}`, true},
		{"type mismatch", `{"a":"1"}`, `{"a":1}`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RawContains(mustJson(t, tt.want), mustJson(t, tt.got)); got != tt.result {
				t.Errorf("RawContains(%s, %s) = %v, want %v", tt.want, tt.got, got, tt.result)
			}
		})
	}
}

func TestWaitForRawValues(t *testing.T) {
	const stale = `{"id":"p","inboundTrust":{"isCompliantNetworkAccepted":false}}`
	const fresh = `{"id":"p","inboundTrust":{"isCompliantNetworkAccepted":true}}`
	want := mustJson(t, `{"inboundTrust":{"isCompliantNetworkAccepted":true},"tenantRestrictions":null}`).(map[string]any)
	opts := WaitForRawValuesOptions{Timeout: 300 * time.Millisecond, PollInterval: time.Millisecond, ConsecutiveReads: 2}

	for _, tc := range []struct {
		name        string
		responses   []string // per GET, the last one repeats
		wantGets    int32    // 0: don't check
		wantWarning bool
	}{
		{"fresh immediately", []string{fresh}, 2, false},
		{"stale then fresh", []string{stale, stale, fresh}, 4, false},
		{"fresh reads must be consecutive", []string{fresh, stale, fresh}, 4, false},
		{"stays stale", []string{stale}, 0, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var gets atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.Method != http.MethodGet || req.URL.Path != "/beta/policies/p" {
					t.Errorf("unexpected request: %s %s", req.Method, req.URL)
				}
				i := int(gets.Add(1)) - 1
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, tc.responses[min(i, len(tc.responses)-1)])
			}))
			defer server.Close()
			client := msgraph.NewClient(msgraph.VersionBeta)
			client.Endpoint = server.URL

			diags := diag.Diagnostics{}
			WaitForRawValues(context.Background(), &diags, &client, msgraph.Uri{Entity: "/policies/p"}, want, opts)

			if diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}
			if gotWarning := diags.WarningsCount() > 0; gotWarning != tc.wantWarning {
				t.Errorf("warning: got %v, want %v (diags %v)", gotWarning, tc.wantWarning, diags)
			}
			if tc.wantGets > 0 && gets.Load() != tc.wantGets {
				t.Errorf("expected %d GETs, got %d", tc.wantGets, gets.Load())
			}
		})
	}
}
