package provider

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/claims"
	"golang.org/x/oauth2"
)

// logTokenClaims are the claims of an access token that are logged for debugging: those known to go-azure-sdk plus
// the ones that show whether claims challenges are in effect (see claims_authorizer.go). The token itself is never
// logged.
type logTokenClaims struct {
	claims.Claims

	// XmsCc lists the client capabilities, e.g. ["cp1"] if the client can handle claims challenges
	XmsCc any `json:"xms_cc,omitempty"`
	// XmsSsm is "1" for Continuous Access Evaluation (CAE) tokens
	XmsSsm any `json:"xms_ssm,omitempty"`
	// Acrs lists the satisfied authentication contexts, e.g. ["c1"] after a claims challenge
	Acrs any `json:"acrs,omitempty"`
}

// tokenClaimsJson returns the claims of token to be logged (see logTokenClaims) as JSON.
func tokenClaimsJson(token *oauth2.Token) (string, error) {
	if token == nil {
		return "", errors.New("token is nil")
	}
	jwt := strings.Split(token.AccessToken, ".")
	if len(jwt) != 3 {
		return "", errors.New("unexpected token format: does not have 3 parts")
	}
	payload, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(jwt[1], "="))
	if err != nil {
		return "", err
	}
	var c logTokenClaims
	if err := json.Unmarshal(payload, &c); err != nil {
		return "", err
	}
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
