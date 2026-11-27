package services

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// TestAuthenticationStrengthPolicyCreateReplaceFunc checks that a create failing with "display names must be unique"
// names the existing policy (#247), and that other results pass through.
func TestAuthenticationStrengthPolicyCreateReplaceFunc(t *testing.T) {
	const duplicateBody = `{"error":{"code":"badRequest","message":"Authentication strength policy display names must be unique."}}`
	const listBody = `{"value":[{"id":"00000000-0000-0000-0000-000000000002","displayName":"Multifactor authentication"},` +
		`{"id":"leftover-1","displayName":"TF Test FIDO2"}]}`

	for _, tc := range []struct {
		name       string
		postStatus int
		postBody   string
		wantId     string
		wantError  string // "" = no error
	}{
		{"created", http.StatusCreated, `{"id":"new-1","displayName":"tf test fido2"}`, "new-1", ""},
		{"duplicate found (case-insensitive)", http.StatusBadRequest, duplicateBody, "", `id "leftover-1"`},
		{"other 400", http.StatusBadRequest, `{"error":{"code":"badRequest","message":"something else"}}`, "", "something else"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if req.URL.Path != "/beta/policies/authenticationStrengthPolicies" {
					t.Errorf("unexpected request: %s %s", req.Method, req.URL)
				}
				w.Header().Set("Content-Type", "application/json")
				switch req.Method {
				case http.MethodPost:
					w.WriteHeader(tc.postStatus)
					_, _ = io.WriteString(w, tc.postBody)
				case http.MethodGet:
					_, _ = io.WriteString(w, listBody)
				}
			}))
			defer server.Close()
			client := msgraph.NewClient(msgraph.VersionBeta)
			client.Endpoint = server.URL
			client.RetryableClient.RetryMax = 0
			client.RetryableClient.RetryWaitMin = time.Millisecond
			client.RetryableClient.RetryWaitMax = 2 * time.Millisecond
			r := &generic.GenericResource{AccessParams: generic.AccessParams{BaseUri: AuthenticationStrengthPolicyResource.AccessParams.BaseUri}}
			r.AccessParams.EntityId.AttrNameGraph = "id"
			r.AccessParams.InitializeGuarded(&client)

			diags := diag.Diagnostics{}
			params := &generic.CreateReplaceFuncParams{R: r, Client: &client, RawVal: map[string]any{"displayName": "tf test fido2"}}
			authenticationStrengthPolicyCreateReplaceFunc(context.Background(), &diags, params)

			if params.Id != tc.wantId {
				t.Errorf("id: got %q, want %q", params.Id, tc.wantId)
			}
			text := authenticationStrengthPolicyDiagsText(diags)
			if tc.wantError == "" && diags.HasError() {
				t.Fatalf("unexpected error: %s", text)
			}
			if tc.wantError != "" && (!diags.HasError() || !strings.Contains(text, tc.wantError)) {
				t.Fatalf("want error containing %q, got %q", tc.wantError, text)
			}
			if strings.Contains(tc.name, "duplicate") && !strings.Contains(text, "import {") {
				t.Errorf("the error should show an import block: %s", text)
			}
		})
	}
}
