package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
)

func TestMergeClaims(t *testing.T) {
	tests := []struct {
		name, extra, want string
	}{
		{"capability only", "", `{"access_token":{"xms_cc":{"values":["cp1"]}}}`},
		{"acrs challenge", `{"access_token":{"acrs":{"essential":true,"values":["c20"]}}}`,
			`{"access_token":{"acrs":{"essential":true,"values":["c20"]},"xms_cc":{"values":["cp1"]}}}`},
		{"challenge with cp1 already", `{"access_token":{"xms_cc":{"values":["cp1"]},"nbf":{"essential":true,"value":"1700000000"}}}`,
			`{"access_token":{"nbf":{"essential":true,"value":"1700000000"},"xms_cc":{"values":["cp1"]}}}`},
		{"id_token too", `{"id_token":{"auth_time":{"essential":true}}}`,
			`{"access_token":{"xms_cc":{"values":["cp1"]}},"id_token":{"auth_time":{"essential":true}}}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mergeClaims(capabilityClaims, tt.extra)
			if err != nil || got != tt.want {
				t.Errorf("got %s, %v; want %s", got, err, tt.want)
			}
		})
	}
	if _, err := mergeClaims(capabilityClaims, "not json"); err == nil {
		t.Error("want error for invalid claims")
	}
}

// tokenServer is a fake Entra token endpoint that returns the requested claims as access token.
func tokenServer(t *testing.T) (*httptest.Server, func() []url.Values) {
	var mu sync.Mutex
	var requests []url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		form, err := url.ParseQuery(string(b))
		if err != nil {
			t.Errorf("token request body: %v", err)
		}
		mu.Lock()
		requests = append(requests, form)
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"token_type":"Bearer","expires_in":3600,"access_token":%q}`, "claims="+form.Get("claims"))
	}))
	return srv, func() []url.Values {
		mu.Lock()
		defer mu.Unlock()
		return append([]url.Values{}, requests...)
	}
}

// TestClaimsChallengeAuthorizer uses a real go-azure-sdk client secret authorizer against a fake token endpoint.
func TestClaimsChallengeAuthorizer(t *testing.T) {
	srv, requests := tokenServer(t)
	defer srv.Close()

	env := environments.AzurePublic()
	env.Authorization.LoginEndpoint = srv.URL
	newAuthorizer := func() (auth.Authorizer, error) {
		return auth.NewAuthorizerFromCredentials(context.Background(), auth.Credentials{
			Environment:                           *env,
			TenantID:                              "tenant",
			ClientID:                              "client",
			ClientSecret:                          "secret",
			EnableAuthenticatingUsingClientSecret: true,
		}, env.MicrosoftGraph)
	}
	base, err := newAuthorizer()
	if err != nil {
		t.Fatal(err)
	}
	if !supportsClaimsChallenge(base) {
		t.Fatal("client secret authorizer should support claims challenges")
	}
	a, err := newClaimsChallengeAuthorizer(base, newAuthorizer)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	acrs := `{"access_token":{"acrs":{"essential":true,"values":["c20"]}}}`
	wantBase := `{"access_token":{"xms_cc":{"values":["cp1"]}}}`
	wantAcrs := `{"access_token":{"acrs":{"essential":true,"values":["c20"]},"xms_cc":{"values":["cp1"]}}}`

	for i, step := range []struct{ claims, want string }{
		{"", wantBase},
		{acrs, wantAcrs},
		{"", wantBase},   // cached
		{acrs, wantAcrs}, // cached
	} {
		token, err := a.TokenWithClaims(ctx, &http.Request{}, step.claims)
		if err != nil {
			t.Fatalf("step %d: %v", i, err)
		}
		if token.AccessToken != "claims="+step.want {
			t.Errorf("step %d: got token %q, want claims %s", i, token.AccessToken, step.want)
		}
	}

	reqs := requests()
	if len(reqs) != 2 {
		t.Fatalf("got %d token requests, want 2 (one per claims set, then cached)", len(reqs))
	}
	for _, r := range reqs {
		if r.Get("grant_type") != "client_credentials" || r.Get("client_secret") != "secret" {
			t.Errorf("token request lost its parameters: %v", r)
		}
	}

	// A token request without claims in its context (e.g. from a provider without enable_claims_challenge)
	// must not be changed.
	token, err := base.(*auth.CachedAuthorizer).Source.Token(ctx, &http.Request{})
	if err != nil {
		t.Fatal(err)
	}
	if token.AccessToken != "claims=" {
		t.Errorf("token request without claims got claims: %q", token.AccessToken)
	}
}

func TestSupportsClaimsChallenge(t *testing.T) {
	wgt, err := NewWgtAuthorizer(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if supportsClaimsChallenge(wgt) {
		t.Error("wpGetToken authorizer cannot request claims")
	}
}
