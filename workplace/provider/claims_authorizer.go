package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"golang.org/x/oauth2"
)

// Claims challenges (https://learn.microsoft.com/entra/identity-platform/claims-challenge)
//
// When an Entra protected action (e.g. creating a Conditional Access policy) requires a Conditional Access
// authentication context, Microsoft Graph rejects a token without the matching `acrs` claim. If the token declares the
// client capability `cp1`, Graph answers with a 401 and an `insufficient_claims` challenge naming the required claims
// instead of a plain 403. A client credentials token request may then ask for these claims via the `claims` parameter.
// Tokens with `cp1` are CAE tokens, so a CAE revocation arrives as the same kind of challenge.
//
// go-azure-sdk has no support for the `claims` parameter. Therefore:
//   - claimsChallengeAuthorizer keeps one (cached) SDK authorizer per claims set and passes the claims to its token
//     requests in the request context,
//   - tokenClaimsHTTPClient (installed as auth.Client, which the SDK uses for all token requests) adds the claims from
//     the context to client credentials token requests.
//
// This works for all authorizers that use the client credentials grant against the token endpoint (client secret,
// client certificate and the OIDC variants), but not for Azure CLI, managed identity or wpGetToken.
//
// A CAE revocation does not invalidate the cached base token: each later request then costs one challenge round trip
// until the base token is renewed. That is acceptable since a provider process only lives for one plan/apply.

// capabilityClaims declares the client capability `cp1` (the client can handle claims challenges).
const capabilityClaims = `{"access_token":{"xms_cc":{"values":["cp1"]}}}`

var _ msgraph.ClaimsChallengeAuthorizer = &claimsChallengeAuthorizer{}

type claimsChallengeAuthorizer struct {
	// newAuthorizer returns a new (cached) SDK authorizer, used for each additional claims set
	newAuthorizer func() (auth.Authorizer, error)

	mutex sync.Mutex
	// byClaims holds one authorizer per claims set (normalized JSON) so each one caches the tokens for its claims
	byClaims map[string]auth.Authorizer
}

// newClaimsChallengeAuthorizer returns a claimsChallengeAuthorizer, using base for tokens with the client capability
// only and newAuthorizer for tokens with additional (challenge) claims.
func newClaimsChallengeAuthorizer(base auth.Authorizer, newAuthorizer func() (auth.Authorizer, error)) (*claimsChallengeAuthorizer, error) {
	baseClaims, err := mergeClaims(capabilityClaims, "")
	if err != nil {
		return nil, err
	}
	installTokenClaimsHTTPClient()
	return &claimsChallengeAuthorizer{
		newAuthorizer: newAuthorizer,
		byClaims:      map[string]auth.Authorizer{baseClaims: base},
	}, nil
}

// supportsClaimsChallenge reports whether a (go-azure-sdk) authorizer requests its tokens with the client credentials
// grant via auth.Client, so that tokenClaimsHTTPClient can add claims to its token requests.
func supportsClaimsChallenge(a auth.Authorizer) bool {
	cached, ok := a.(*auth.CachedAuthorizer)
	if !ok {
		return false
	}
	switch cached.Source.(type) {
	case *auth.ClientSecretAuthorizer, *auth.ClientAssertionAuthorizer, *auth.GitHubOIDCAuthorizer, *auth.ADOPipelineOIDCAuthorizer:
		return true
	}
	return false
}

// Token returns a token with the client capability `cp1`.
func (a *claimsChallengeAuthorizer) Token(ctx context.Context, req *http.Request) (*oauth2.Token, error) {
	return a.TokenWithClaims(ctx, req, "")
}

// TokenWithClaims returns a token with the client capability `cp1` and the given (challenge) claims.
func (a *claimsChallengeAuthorizer) TokenWithClaims(ctx context.Context, req *http.Request, claims string) (*oauth2.Token, error) {
	claims, authorizer, err := a.authorizerFor(claims)
	if err != nil {
		return nil, err
	}
	return authorizer.Token(context.WithValue(ctx, tokenClaimsContextKey{}, claims), req)
}

// AuxiliaryTokens returns the auxiliary tokens of the base authorizer (client capability `cp1` only).
func (a *claimsChallengeAuthorizer) AuxiliaryTokens(ctx context.Context, req *http.Request) ([]*oauth2.Token, error) {
	claims, authorizer, err := a.authorizerFor("")
	if err != nil {
		return nil, err
	}
	return authorizer.AuxiliaryTokens(context.WithValue(ctx, tokenClaimsContextKey{}, claims), req)
}

// authorizerFor returns the claims to request (challenge claims merged with the client capability, normalized) and
// the authorizer to request them with.
func (a *claimsChallengeAuthorizer) authorizerFor(challengeClaims string) (string, auth.Authorizer, error) {
	claims, err := mergeClaims(capabilityClaims, challengeClaims)
	if err != nil {
		return "", nil, fmt.Errorf("invalid claims %q: %w", challengeClaims, err)
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()
	if authorizer, ok := a.byClaims[claims]; ok {
		return claims, authorizer, nil
	}
	authorizer, err := a.newAuthorizer()
	if err != nil {
		return "", nil, fmt.Errorf("configuring authorizer for claims %s: %w", claims, err)
	}
	a.byClaims[claims] = authorizer
	return claims, authorizer, nil
}

// mergeClaims deep-merges the JSON objects base and extra (extra wins on conflicting non-object values) and returns
// the result as normalized JSON (sorted keys), so it can be used as a cache key. An empty extra is ignored.
func mergeClaims(base string, extra string) (string, error) {
	var merged map[string]any
	if err := json.Unmarshal([]byte(base), &merged); err != nil {
		return "", err
	}
	if strings.TrimSpace(extra) != "" {
		var e map[string]any
		if err := json.Unmarshal([]byte(extra), &e); err != nil {
			return "", err
		}
		deepMerge(merged, e)
	}
	b, err := json.Marshal(merged)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func deepMerge(dst map[string]any, src map[string]any) {
	for k, sv := range src {
		if sm, ok := sv.(map[string]any); ok {
			if dm, ok := dst[k].(map[string]any); ok {
				deepMerge(dm, sm)
				continue
			}
		}
		dst[k] = sv
	}
}

// tokenClaimsContextKey is the context key for the claims (JSON) to add to token requests.
type tokenClaimsContextKey struct{}

var installTokenClaimsHTTPClientOnce sync.Once

// installTokenClaimsHTTPClient wraps auth.Client (process wide, once) with tokenClaimsHTTPClient. Token requests
// without claims in their context are not changed.
func installTokenClaimsHTTPClient() {
	installTokenClaimsHTTPClientOnce.Do(func() {
		auth.Client = tokenClaimsHTTPClient{inner: auth.Client}
	})
}

// tokenClaimsHTTPClient adds the claims from the request context as `claims` parameter to client credentials token
// requests (form posts with grant_type=client_credentials).
type tokenClaimsHTTPClient struct {
	inner auth.HTTPClient
}

func (c tokenClaimsHTTPClient) Do(req *http.Request) (*http.Response, error) {
	claims, _ := req.Context().Value(tokenClaimsContextKey{}).(string)
	if claims == "" || req.Method != http.MethodPost || req.Body == nil ||
		!strings.HasPrefix(req.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		return c.inner.Do(req)
	}

	body, err := io.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("reading token request body: %w", err)
	}
	if form, err := url.ParseQuery(string(body)); err == nil && form.Get("grant_type") == "client_credentials" && !form.Has("claims") {
		form.Set("claims", claims)
		body = []byte(form.Encode())
	}

	req = req.Clone(req.Context())
	req.Body = io.NopCloser(bytes.NewReader(body))
	req.ContentLength = int64(len(body))
	req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(body)), nil }
	return c.inner.Do(req)
}
