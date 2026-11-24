package provider

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDumpRequestRedacted(t *testing.T) {
	const token = "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.secret.signature"
	const body = `{"displayName":"x"}`

	req, err := http.NewRequest(http.MethodPost, "https://graph.microsoft.com/beta/identity/conditionalAccess/policies", bytes.NewBufferString(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	dump, err := dumpRequestRedacted(req)
	if err != nil {
		t.Fatal(err)
	}
	d := string(dump)
	if strings.Contains(d, "eyJ") || strings.Contains(d, "secret") {
		t.Errorf("dump contains the token:\n%s", d)
	}
	if !strings.Contains(d, "Authorization: <redacted>") {
		t.Errorf("dump does not show the redacted header:\n%s", d)
	}
	if !strings.Contains(d, body) || !strings.Contains(d, "Content-Type: application/json") {
		t.Errorf("dump lost the body or other headers:\n%s", d)
	}

	// The request itself must be unchanged and still sendable
	if got := req.Header.Get("Authorization"); got != token {
		t.Errorf("request Authorization header changed to %q", got)
	}
	b, err := io.ReadAll(req.Body)
	if err != nil || string(b) != body {
		t.Errorf("request body after dump: %q, %v; want %q", b, err, body)
	}
}

func TestDumpRequestRedactedWithoutAuthorization(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://graph.microsoft.com/beta/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	dump, err := dumpRequestRedacted(req)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(dump), "Authorization") {
		t.Errorf("dump shows an Authorization header that the request does not have:\n%s", dump)
	}
}
