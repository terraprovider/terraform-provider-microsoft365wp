package msgraph

import (
	"context"
	"encoding/base64"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"golang.org/x/oauth2"
)

// ClaimsChallengeAuthorizer is an Authorizer that can also acquire tokens carrying additional claims, as requested by a
// claims challenge of Microsoft Graph (e.g. a Conditional Access authentication context required by an Entra protected
// action, or a Continuous Access Evaluation revocation).
//
// If the Authorizer of a Client implements this interface, the Client answers a 401 response with an
// `insufficient_claims` challenge by acquiring a token with the requested claims and sending the request once more.
type ClaimsChallengeAuthorizer interface {
	auth.Authorizer

	// TokenWithClaims returns a token whose token request included claims (a JSON object as decoded from the
	// challenge, e.g. `{"access_token":{"acrs":{"essential":true,"values":["c1"]}}}`).
	TokenWithClaims(ctx context.Context, req *http.Request, claims string) (*oauth2.Token, error)
}

var (
	reInsufficientClaims = regexp.MustCompile(`(?i)\berror\s*=\s*"insufficient_claims"`)
	reChallengeClaims    = regexp.MustCompile(`(?i)\bclaims\s*=\s*"([^"]*)"`)
)

// claimsFromChallenge returns the decoded claims of an `insufficient_claims` challenge in the WWW-Authenticate
// header(s) of resp, e.g.
//
//	WWW-Authenticate: Bearer realm="", authorization_uri="...", client_id="...", error="insufficient_claims", claims="eyJhY2Nlc3NfdG9rZW4iOnsi..."
func claimsFromChallenge(resp *http.Response) (string, bool) {
	if resp == nil || resp.StatusCode != http.StatusUnauthorized {
		return "", false
	}
	for _, h := range resp.Header.Values("WWW-Authenticate") {
		if !reInsufficientClaims.MatchString(h) {
			continue
		}
		m := reChallengeClaims.FindStringSubmatch(h)
		if m == nil {
			continue
		}
		if claims, ok := decodeChallengeClaims(m[1]); ok {
			return claims, true
		}
	}
	return "", false
}

// decodeChallengeClaims decodes the base64 claims value of a challenge. Graph sends standard base64 with padding, but
// be lenient and also accept the unpadded and URL-safe variants.
func decodeChallengeClaims(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	for _, enc := range []*base64.Encoding{base64.StdEncoding, base64.RawStdEncoding, base64.URLEncoding, base64.RawURLEncoding} {
		if b, err := enc.DecodeString(s); err == nil && len(b) > 0 {
			return string(b), true
		}
	}
	return "", false
}
