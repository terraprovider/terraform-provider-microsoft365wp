package msgraph

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"golang.org/x/oauth2"
)

const challengeClaimsJson = `{"access_token":{"acrs":{"essential":true,"values":["c20"]}}}`

// fakeAuthorizer returns the token "base", or "claims:<claims>" for TokenWithClaims.
type fakeAuthorizer struct{}

func (fakeAuthorizer) Token(context.Context, *http.Request) (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: "base", TokenType: "Bearer"}, nil
}

func (fakeAuthorizer) AuxiliaryTokens(context.Context, *http.Request) ([]*oauth2.Token, error) {
	return nil, nil
}

type fakeClaimsAuthorizer struct{ fakeAuthorizer }

func (fakeClaimsAuthorizer) TokenWithClaims(_ context.Context, _ *http.Request, claims string) (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: "claims:" + claims, TokenType: "Bearer"}, nil
}

// challengeServer answers with a claims challenge (as sent by Graph for a protected action) unless the request carries
// the token for challengeClaimsJson (or always, if alwaysChallenge). It records the bodies it received.
func challengeServer(alwaysChallenge bool) (*httptest.Server, *atomic.Int32, *[]string) {
	var count atomic.Int32
	bodies := &[]string{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count.Add(1)
		b, _ := io.ReadAll(r.Body)
		*bodies = append(*bodies, string(b))
		w.Header().Set("Content-Type", "application/json")
		if alwaysChallenge || r.Header.Get("Authorization") != "Bearer claims:"+challengeClaimsJson {
			w.Header().Set("WWW-Authenticate", `Bearer realm="", authorization_uri="https://login.microsoftonline.com/common/oauth2/authorize", client_id="00000003-0000-0000-c000-000000000000", error="insufficient_claims", claims="`+
				base64.StdEncoding.EncodeToString([]byte(challengeClaimsJson))+`"`)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":{"code":"UnknownError","message":"insufficient_claims"}}`))
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"x"}`))
	}))
	return srv, &count, bodies
}

func postX(c Client) (int, error) {
	_, status, _, err := c.Post(context.Background(), PostHttpRequestInput{
		Uri:              Uri{Entity: "/identity/conditionalAccess/policies"},
		Body:             []byte(`{"displayName":"x"}`),
		ValidStatusCodes: []int{http.StatusCreated},
	})
	return status, err
}

func TestClaimsChallengeIsAnsweredOnce(t *testing.T) {
	srv, count, bodies := challengeServer(false)
	defer srv.Close()
	c := newTestClient(srv)
	c.Authorizer = fakeClaimsAuthorizer{}

	status, err := postX(c)
	if err != nil || status != http.StatusCreated {
		t.Fatalf("got status %d, err %v; want 201", status, err)
	}
	if n := count.Load(); n != 2 {
		t.Errorf("got %d requests, want 2 (challenge + retry)", n)
	}
	for i, b := range *bodies {
		if b != `{"displayName":"x"}` {
			t.Errorf("request %d had body %q, want the original body", i, b)
		}
	}
}

func TestClaimsChallengeIsNotAnsweredTwice(t *testing.T) {
	srv, count, _ := challengeServer(true)
	defer srv.Close()
	c := newTestClient(srv)
	c.Authorizer = fakeClaimsAuthorizer{}

	status, err := postX(c)
	if err == nil || status != http.StatusUnauthorized {
		t.Fatalf("got status %d, err %v; want 401 error", status, err)
	}
	if n := count.Load(); n != 2 {
		t.Errorf("got %d requests, want 2 (no second retry)", n)
	}
}

func TestClaimsChallengeWithoutClaimsAuthorizer(t *testing.T) {
	srv, count, _ := challengeServer(false)
	defer srv.Close()
	c := newTestClient(srv)
	c.Authorizer = fakeAuthorizer{}

	status, err := postX(c)
	if err == nil || status != http.StatusUnauthorized {
		t.Fatalf("got status %d, err %v; want 401 error", status, err)
	}
	if n := count.Load(); n != 1 {
		t.Errorf("got %d requests, want 1", n)
	}
}

func TestClaimsFromChallenge(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte(challengeClaimsJson))
	tests := []struct {
		name   string
		status int
		header []string
		want   string
	}{
		{"graph challenge", 401, []string{`Bearer realm="", error="insufficient_claims", claims="` + encoded + `"`}, challengeClaimsJson},
		{"unpadded", 401, []string{`Bearer error="insufficient_claims", claims="` + strings.TrimRight(encoded, "=") + `"`}, challengeClaimsJson},
		{"second header", 401, []string{`Basic realm="x"`, `Bearer claims="` + encoded + `", error="insufficient_claims"`}, challengeClaimsJson},
		{"other error", 401, []string{`Bearer error="invalid_token", claims="` + encoded + `"`}, ""},
		{"no claims", 401, []string{`Bearer error="insufficient_claims"`}, ""},
		{"not 401", 403, []string{`Bearer error="insufficient_claims", claims="` + encoded + `"`}, ""},
		{"no header", 401, nil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{StatusCode: tt.status, Header: http.Header{}}
			for _, h := range tt.header {
				resp.Header.Add("WWW-Authenticate", h)
			}
			got, ok := claimsFromChallenge(resp)
			if got != tt.want || ok != (tt.want != "") {
				t.Errorf("got %q, %v; want %q", got, ok, tt.want)
			}
		})
	}
}
