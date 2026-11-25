package provider

import (
	"encoding/base64"
	"strings"
	"testing"

	"golang.org/x/oauth2"
)

func testToken(payload string) *oauth2.Token {
	enc := base64.RawURLEncoding.EncodeToString
	return &oauth2.Token{AccessToken: enc([]byte(`{"alg":"RS256","typ":"JWT"}`)) + "." + enc([]byte(payload)) + ".signature"}
}

func TestTokenClaimsJson(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		want    []string // substrings of the result
		notWant []string
	}{
		{
			name:    "claims challenge token",
			payload: `{"aud":"https://graph.microsoft.com","tid":"t","appid":"a","roles":["Policy.ReadWrite.ConditionalAccess"],"xms_cc":["cp1"],"xms_ssm":"1","acrs":["c20","pfdr"],"unknown":"x"}`,
			want:    []string{`"xms_cc":["cp1"]`, `"xms_ssm":"1"`, `"acrs":["c20","pfdr"]`, `"appid":"a"`, `"roles":["Policy.ReadWrite.ConditionalAccess"]`},
			notWant: []string{"unknown", "signature"},
		},
		{
			name:    "plain token",
			payload: `{"aud":"https://graph.microsoft.com","tid":"t"}`,
			want:    []string{`"aud":"https://graph.microsoft.com"`},
			notWant: []string{"xms_cc", "xms_ssm", "acrs"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tokenClaimsJson(testToken(tt.payload))
			if err != nil {
				t.Fatal(err)
			}
			for _, w := range tt.want {
				if !strings.Contains(got, w) {
					t.Errorf("missing %s in %s", w, got)
				}
			}
			for _, n := range tt.notWant {
				if strings.Contains(got, n) {
					t.Errorf("unexpected %s in %s", n, got)
				}
			}
		})
	}

	if _, err := tokenClaimsJson(&oauth2.Token{AccessToken: "not-a-jwt"}); err == nil {
		t.Error("want error for a token that is not a JWT")
	}
	if _, err := tokenClaimsJson(nil); err == nil {
		t.Error("want error for nil token")
	}
}
