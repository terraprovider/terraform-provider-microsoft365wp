package browsersitelistutils

import (
	"context"
	"testing"

	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TestNextRevision(t *testing.T) {
	cases := []struct {
		current string
		want    string
		wantErr bool
	}{
		{"", "1", false},
		{"1", "2", false},
		{"1.0", "2", false},
		{"10.5", "11", false},
		{"10.10", "11", false},
		{"0.5", "1", false},
		{"abc", "", true},
		{"-1", "", true},
	}
	for _, c := range cases {
		got, err := nextRevision(c.current)
		if (err != nil) != c.wantErr {
			t.Errorf("nextRevision(%q): error = %v, wantErr %v", c.current, err, c.wantErr)
			continue
		}
		if got != c.want {
			t.Errorf("nextRevision(%q) = %q, want %q", c.current, got, c.want)
		}
	}
}

func TestKeys(t *testing.T) {
	if siteKey(map[string]any{"webUrl": "https://FULL.contoso.test/app"}) != siteKey(map[string]any{"webUrl": "https://full.contoso.test/app"}) {
		t.Error("siteKey must be case-insensitive")
	}
	if siteKey(map[string]any{"webUrl": "https://full.contoso.test/app"}) == siteKey(map[string]any{"webUrl": "https://full.contoso.test/app/"}) {
		t.Error("siteKey must distinguish a trailing slash")
	}

	cookie := func(name, host, path string, hostOnly bool) map[string]any {
		return map[string]any{"displayName": name, "hostOrDomain": host, "path": path, "hostOnly": hostOnly, "comment": "x"}
	}
	if sharedCookieKey(cookie("C2", "CONTOSO.test", "/", true)) != sharedCookieKey(cookie("c2", "contoso.test", "/", true)) {
		t.Error("sharedCookieKey must be case-insensitive")
	}
	if sharedCookieKey(cookie("c2", "contoso.test", "/", true)) == sharedCookieKey(cookie("c2", "contoso.test", "/", false)) {
		t.Error("sharedCookieKey must distinguish hostOnly")
	}
	if sharedCookieKey(cookie("c2", "contoso.test", "/", true)) == sharedCookieKey(cookie("c2", "contoso.test", "/a", true)) {
		t.Error("sharedCookieKey must distinguish path")
	}
}

func TestElementEquals(t *testing.T) {
	planned := map[string]any{"webUrl": "https://a.test", "allowRedirect": false, "comment": ""}
	current := map[string]any{"id": "1", "status": "published", "webUrl": "https://a.test", "allowRedirect": false, "comment": "", "history": []any{}}
	if !elementEquals(planned, current) {
		t.Error("expected equal")
	}
	current["comment"] = "changed"
	if elementEquals(planned, current) {
		t.Error("expected not equal after comment change")
	}
	current["comment"] = ""
	current["webUrl"] = "https://A.test"
	if elementEquals(planned, current) {
		t.Error("expected not equal after case-only webUrl change (must be written to MS Graph)")
	}
}

func TestExecutePreDuplicates(t *testing.T) {
	a := &SyncAndPublishWsa{}
	a.Initialize()
	wsaReq := &generic.WriteSubActionRequest{
		Operation: generic.OperationCreate,
		RawVal: map[string]any{
			"displayName": "x",
			"sites": []any{
				map[string]any{"webUrl": "https://a.test"},
				map[string]any{"webUrl": "https://A.test", "comment": "other"},
			},
			"sharedCookies": []any{},
		},
		SubActionsData: map[string]any{},
	}
	var diags diag.Diagnostics
	a.ExecutePre(context.Background(), &diags, wsaReq)
	if !diags.HasError() {
		t.Fatal("expected an error for duplicate sites")
	}
	if _, ok := wsaReq.RawVal["sites"]; ok {
		t.Error("sites must be removed from the list's own POST/PATCH body")
	}
}

func TestGraphToTerraformMiddleware(t *testing.T) {
	params := &generic.GraphToTerraformMiddlewareParams{RawVal: map[string]any{
		"sites": []any{
			map[string]any{"webUrl": "a", "status": "published"},
			map[string]any{"webUrl": "b", "status": "pendingDelete"},
			map[string]any{"webUrl": "c", "status": "pendingAdd"},
		},
		"sharedCookies": []any{
			map[string]any{"displayName": "x", "status": "pendingDelete"},
		},
	}}
	var diags diag.Diagnostics
	if err := GraphToTerraformMiddleware(context.Background(), &diags, params); err != nil {
		t.Fatal(err)
	}
	if got := len(params.RawVal["sites"].([]any)); got != 2 {
		t.Errorf("sites: got %d elements, want 2", got)
	}
	if got := len(params.RawVal["sharedCookies"].([]any)); got != 0 {
		t.Errorf("sharedCookies: got %d elements, want 0", got)
	}
}
