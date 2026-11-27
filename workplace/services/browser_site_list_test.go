package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"
	"terraform-provider-microsoft365wp/workplace/generic"
	browsersitelistutils "terraform-provider-microsoft365wp/workplace/services/browser_site_list_utils"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type browserSiteListTestSite struct {
	AllowRedirect     types.Bool   `tfsdk:"allow_redirect"`
	Comment           types.String `tfsdk:"comment"`
	CompatibilityMode types.String `tfsdk:"compatibility_mode"`
	MergeType         types.String `tfsdk:"merge_type"`
	TargetEnvironment types.String `tfsdk:"target_environment"`
	WebUrl            types.String `tfsdk:"web_url"`
}

type browserSiteListTestCookie struct {
	Comment           types.String `tfsdk:"comment"`
	DisplayName       types.String `tfsdk:"display_name"`
	HostOnly          types.Bool   `tfsdk:"host_only"`
	HostOrDomain      types.String `tfsdk:"host_or_domain"`
	Path              types.String `tfsdk:"path"`
	SourceEnvironment types.String `tfsdk:"source_environment"`
}

func browserSiteListTestConfig(t *testing.T, sites []browserSiteListTestSite, cookies []browserSiteListTestCookie) tfsdk.Config {
	t.Helper()
	ctx := context.Background()
	state := tfsdk.State{Schema: browserSiteListResourceSchema, Raw: tftypes.NewValue(browserSiteListResourceSchema.Type().TerraformType(ctx), nil)}
	if diags := state.SetAttribute(ctx, path.Root("display_name"), "x"); diags.HasError() {
		t.Fatal(diags)
	}
	if diags := state.SetAttribute(ctx, path.Root("sites"), sites); diags.HasError() {
		t.Fatal(diags)
	}
	if diags := state.SetAttribute(ctx, path.Root("shared_cookies"), cookies); diags.HasError() {
		t.Fatal(diags)
	}
	return tfsdk.Config{Schema: state.Schema, Raw: state.Raw}
}

func TestBrowserSiteListConfigValidator(t *testing.T) {
	ctx := context.Background()
	site := func(url string) browserSiteListTestSite {
		return browserSiteListTestSite{WebUrl: types.StringValue(url)}
	}
	cookie := func(name, host string, hostOnly types.Bool) browserSiteListTestCookie {
		return browserSiteListTestCookie{DisplayName: types.StringValue(name), HostOrDomain: types.StringValue(host), Path: types.StringValue("/"), HostOnly: hostOnly}
	}
	for _, tc := range []struct {
		name      string
		sites     []browserSiteListTestSite
		cookies   []browserSiteListTestCookie
		wantError bool
	}{
		{"valid", []browserSiteListTestSite{site("https://a.test"), site("https://a.test/"), site("a.test")},
			[]browserSiteListTestCookie{cookie("c", "a.test", types.BoolValue(true)), cookie("c", ".a.test", types.BoolNull())}, false},
		{"duplicate site (case-insensitive)", []browserSiteListTestSite{site("https://a.test"), site("https://A.test")}, nil, true},
		{"unknown web_url is skipped", []browserSiteListTestSite{site("https://a.test"), {WebUrl: types.StringUnknown()}}, nil, false},
		{"duplicate cookie (case-insensitive)", nil,
			[]browserSiteListTestCookie{cookie("c", "a.test", types.BoolValue(true)), cookie("C", "A.test", types.BoolValue(true))}, true},
		{"domain cookie without dot", nil, []browserSiteListTestCookie{cookie("c", "a.test", types.BoolNull())}, true},
		{"host-only cookie with dot", nil, []browserSiteListTestCookie{cookie("c", ".a.test", types.BoolValue(true))}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := resource.ValidateConfigRequest{Config: browserSiteListTestConfig(t, tc.sites, tc.cookies)}
			resp := &resource.ValidateConfigResponse{}
			browsersitelistutils.ConfigValidator{}.ValidateResource(ctx, req, resp)
			if resp.Diagnostics.HasError() != tc.wantError {
				t.Errorf("got error %v, want %v (diags %v)", resp.Diagnostics.HasError(), tc.wantError, resp.Diagnostics)
			}
		})
	}
}

// TestBrowserSiteListSyncAndPublishWsa checks the requests of SyncAndPublishWsa against a fake MS Graph with pending
// entries from elsewhere.
func TestBrowserSiteListSyncAndPublishWsa(t *testing.T) {
	const listPath = "/beta/admin/edge/internetExplorerMode/siteLists/l1"

	site := func(id, url, status, comment string) map[string]any {
		return map[string]any{"id": id, "webUrl": url, "status": status, "comment": comment, "allowRedirect": false,
			"compatibilityMode": "default", "mergeType": "default", "targetEnvironment": "internetExplorerMode", "history": []any{}}
	}
	currentSites := []any{
		site("s1", "https://keep.test", "published", ""),             // unchanged: nothing
		site("s2", "https://edit.test", "published", "old"),          // changed comment: PATCH
		site("s3", "https://extra.test", "pendingAdd", ""),           // not planned: DELETE
		site("s4", "https://restore.test", "pendingDelete", ""),      // planned: PATCH (restores it)
		site("s5", "https://gone.test", "pendingDelete", ""),         // not planned: nothing (publish removes it)
		site("s6", "https://same.test", "pendingEdit", ""),           // equal to the plan, but pending: PATCH
		site("s7", "https://KEEP.test", "published", ""),             // duplicate key of s1: DELETE
		site("s8", "https://Case.test", "published", "case changed"), // case-only change of web_url: PATCH
	}
	currentCookies := []any{
		map[string]any{"id": "c1", "displayName": "old", "hostOrDomain": "a.test", "path": "/", "hostOnly": true, "status": "published",
			"comment": "", "sourceEnvironment": "microsoftEdge"},
	}

	var requests []string
	var publishBody string
	r := &generic.GenericResource{
		TypeNameSuffix: BrowserSiteListResource.TypeNameSuffix,
		SpecificSchema: browserSiteListResourceSchema,
		AccessParams:   generic.AccessParams{BaseUri: BrowserSiteListResource.AccessParams.BaseUri},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		p := strings.TrimPrefix(req.URL.Path, listPath)
		body, _ := io.ReadAll(req.Body)
		if req.Method != http.MethodPatch {
			w.Header().Set("Content-Type", "application/json")
		}
		switch {
		case req.Method == http.MethodGet && p == "/sites":
			_ = json.NewEncoder(w).Encode(map[string]any{"value": currentSites})
		case req.Method == http.MethodGet && p == "/sharedCookies":
			_ = json.NewEncoder(w).Encode(map[string]any{"value": currentCookies})
		case req.Method == http.MethodGet && p == "":
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "l1", "revision": "10.5"})
		case req.Method == http.MethodPost && p == "/publish":
			publishBody = string(body)
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "l1", "revision": "11", "status": "published"})
		case req.Method == http.MethodPost:
			requests = append(requests, fmt.Sprintf("POST %s %s", p, body))
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"id":"new"}`))
		case req.Method == http.MethodPatch:
			requests = append(requests, fmt.Sprintf("PATCH %s", p))
			w.WriteHeader(http.StatusNoContent)
		case req.Method == http.MethodDelete:
			requests = append(requests, fmt.Sprintf("DELETE %s", p))
			// MS Graph returns 200 with the entry (status pendingDelete) for published entries
			_, _ = w.Write([]byte(`{"status":"pendingDelete"}`))
		default:
			t.Errorf("unexpected request: %s %s", req.Method, req.URL)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	t.Cleanup(server.Close)
	client := msgraph.NewClient(msgraph.VersionBeta)
	client.Endpoint = server.URL
	client.RetryableClient.RetryWaitMin = time.Millisecond
	client.RetryableClient.RetryWaitMax = 2 * time.Millisecond
	r.AccessParams.InitializeGuarded(&client)

	planSite := func(url, comment string) any {
		return map[string]any{"webUrl": url, "comment": comment, "allowRedirect": false, "compatibilityMode": "default",
			"mergeType": "default", "targetEnvironment": "internetExplorerMode"}
	}
	rawVal := map[string]any{
		"displayName": "x",
		"sites": []any{
			planSite("https://keep.test", ""),
			planSite("https://edit.test", "new"),
			planSite("https://restore.test", ""),
			planSite("https://same.test", ""),
			planSite("https://case.test", "case changed"),
			planSite("https://new.test", ""),
		},
		"sharedCookies": []any{
			map[string]any{"displayName": "new", "hostOrDomain": "a.test", "path": "/", "hostOnly": true, "comment": "", "sourceEnvironment": "both"},
		},
	}

	wsa := &browsersitelistutils.SyncAndPublishWsa{}
	wsa.Initialize()
	diags := diag.Diagnostics{}
	wsaReq := &generic.WriteSubActionRequest{GenRes: r, Operation: generic.OperationUpdate, Id: "l1", RawVal: rawVal, SubActionsData: map[string]any{}}
	wsa.ExecutePre(context.Background(), &diags, wsaReq)
	wsa.ExecutePost(context.Background(), &diags, wsaReq)
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if _, ok := rawVal["sites"]; ok {
		t.Error("sites must not be part of the PATCH of the list itself")
	}

	want := []string{
		"DELETE /sites/s3",
		"DELETE /sites/s7",
		"PATCH /sites/s2",
		"PATCH /sites/s4",
		"PATCH /sites/s6",
		"PATCH /sites/s8",
		`POST /sites {"allowRedirect":false,"comment":"","compatibilityMode":"default","mergeType":"default","targetEnvironment":"internetExplorerMode","webUrl":"https://new.test"}`,
		"DELETE /sharedCookies/c1",
		`POST /sharedCookies {"comment":"","displayName":"new","hostOnly":true,"hostOrDomain":"a.test","path":"/","sourceEnvironment":"both"}`,
	}
	if !reflect.DeepEqual(requests, want) {
		t.Errorf("requests:\n%s\nwant:\n%s", strings.Join(requests, "\n"), strings.Join(want, "\n"))
	}
	if publishBody != `{"revision":"11"}` {
		t.Errorf("publish body = %s, want {\"revision\":\"11\"}", publishBody)
	}
}

// browserSiteListTestProvider is a minimal provider that only serves the browser site list resource, to run plans
// through the framework.
type browserSiteListTestProvider struct{}

func (browserSiteListTestProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "microsoft365wp"
}
func (browserSiteListTestProvider) Schema(context.Context, provider.SchemaRequest, *provider.SchemaResponse) {
}
func (browserSiteListTestProvider) Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse) {
}
func (browserSiteListTestProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{func() resource.Resource { return &BrowserSiteListResource }}
}
func (browserSiteListTestProvider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}

// TestBrowserSiteListPlanUnchangedConfig runs a plan through the framework with an unchanged config that omits
// defaults, and a proposed new state that (like Terraform/OpenTofu) fills null Optional+Computed attributes of the
// set elements from the prior state. The live test found configured values planned back to the defaults with default
// plan modifiers on the nested attributes (the framework pairs config and plan elements of sets by index).
func TestBrowserSiteListPlanUnchangedConfig(t *testing.T) {
	ctx := context.Background()
	server := providerserver.NewProtocol6(browserSiteListTestProvider{})()
	schemaResp, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil || len(schemaResp.Diagnostics) > 0 {
		t.Fatalf("GetProviderSchema: %v %v", err, schemaResp.Diagnostics)
	}
	objType := schemaResp.ResourceSchemas["microsoft365wp_browser_site_list"].ValueType()

	s := func(v string) types.String { return types.StringValue(v) }
	b := func(v bool) types.Bool { return types.BoolValue(v) }
	nullS, nullB := types.StringNull(), types.BoolNull()
	fullSite := func(url, env, compat, merge string, redirect bool, comment string) browserSiteListTestSite {
		return browserSiteListTestSite{WebUrl: s(url), TargetEnvironment: s(env), CompatibilityMode: s(compat), MergeType: s(merge), AllowRedirect: b(redirect), Comment: s(comment)}
	}

	priorSites := []browserSiteListTestSite{
		fullSite("erp.contoso.com", "internetExplorerMode", "default", "default", false, ""),
		fullSite("https://intranet.contoso.com/legacy", "internetExplorerMode", "internetExplorer8Enterprise", "noMerge", true, "needs IE8 document mode"),
		fullSite("https://intranet.contoso.com/new", "microsoftEdge", "default", "default", true, ""),
	}
	priorCookies := []browserSiteListTestCookie{
		{DisplayName: s("ASP.NET_SessionId"), HostOrDomain: s("erp.contoso.com"), Path: s("/"), HostOnly: b(true), SourceEnvironment: s("both"), Comment: s("")},
	}
	configSites := []browserSiteListTestSite{
		{WebUrl: s("erp.contoso.com"), TargetEnvironment: nullS, CompatibilityMode: nullS, MergeType: nullS, AllowRedirect: nullB, Comment: nullS},
		priorSites[1],
		// allow_redirect removed from the config: must be planned back to its default
		{WebUrl: s("https://intranet.contoso.com/new"), TargetEnvironment: s("microsoftEdge"), CompatibilityMode: nullS, MergeType: nullS, AllowRedirect: nullB, Comment: nullS},
	}
	configCookies := []browserSiteListTestCookie{
		{DisplayName: s("ASP.NET_SessionId"), HostOrDomain: s("erp.contoso.com"), Path: s("/"), HostOnly: b(true), SourceEnvironment: s("both"), Comment: nullS},
	}

	build := func(sites []browserSiteListTestSite, cookies []browserSiteListTestCookie, prior bool) tftypes.Value {
		state := tfsdk.State{Schema: browserSiteListResourceSchema, Raw: tftypes.NewValue(objType, nil)}
		set := func(name string, v any) {
			if diags := state.SetAttribute(ctx, path.Root(name), v); diags.HasError() {
				t.Fatalf("SetAttribute(%s): %v", name, diags)
			}
		}
		set("display_name", "tf-test-bsl-1")
		set("sites", sites)
		set("shared_cookies", cookies)
		if prior {
			set("id", "cd22a823-a55e-41fb-9f57-7d86e4caf841")
			set("description", "Legacy line-of-business apps")
			set("revision", "1")
			set("status", "published")
			set("published_date_time", "2026-11-27T12:00:00Z")
			set("last_modified_date_time", "2026-11-27T12:00:00Z")
		} else {
			set("description", "Legacy line-of-business apps")
		}
		return state.Raw
	}
	dyn := func(v tftypes.Value) *tfprotov6.DynamicValue {
		d, err := tfprotov6.NewDynamicValue(objType, v)
		if err != nil {
			t.Fatal(err)
		}
		return &d
	}

	prior := build(priorSites, priorCookies, true)
	config := build(configSites, configCookies, false)
	// Terraform/OpenTofu fill the null Optional+Computed attributes from the prior state (also the removed
	// allow_redirect), and order the set elements by their values, i.e. differently from the config
	proposed := build([]browserSiteListTestSite{priorSites[1], priorSites[2], priorSites[0]}, priorCookies, true)

	resp, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
		TypeName:         "microsoft365wp_browser_site_list",
		PriorState:       dyn(prior),
		ProposedNewState: dyn(proposed),
		Config:           dyn(config),
	})
	if err != nil || len(resp.Diagnostics) > 0 {
		for _, d := range resp.Diagnostics {
			t.Errorf("diagnostic: %s: %s", d.Summary, d.Detail)
		}
		t.Fatalf("PlanResourceChange: %v", err)
	}
	planned, err := resp.PlannedState.Unmarshal(objType)
	if err != nil {
		t.Fatal(err)
	}

	wantSites := []browserSiteListTestSite{priorSites[0], priorSites[1],
		fullSite("https://intranet.contoso.com/new", "microsoftEdge", "default", "default", false, "")}
	want := build(wantSites, priorCookies, true)
	plannedSites, _, _ := tftypes.WalkAttributePath(planned, tftypes.NewAttributePath().WithAttributeName("sites"))
	wantSitesTf, _, _ := tftypes.WalkAttributePath(want, tftypes.NewAttributePath().WithAttributeName("sites"))
	if !sameSetElements(plannedSites.(tftypes.Value), wantSitesTf.(tftypes.Value)) {
		t.Errorf("planned sites:\n%s\nwant:\n%s", plannedSites, wantSitesTf)
	}
	plannedCookies, _, _ := tftypes.WalkAttributePath(planned, tftypes.NewAttributePath().WithAttributeName("shared_cookies"))
	wantCookiesTf, _, _ := tftypes.WalkAttributePath(want, tftypes.NewAttributePath().WithAttributeName("shared_cookies"))
	if !sameSetElements(plannedCookies.(tftypes.Value), wantCookiesTf.(tftypes.Value)) {
		t.Errorf("planned shared_cookies:\n%s\nwant:\n%s", plannedCookies, wantCookiesTf)
	}

	// and with allow_redirect still in the config, the plan must equal the prior state (no diff at all)
	configSites[2].AllowRedirect = b(true)
	resp, err = server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
		TypeName:         "microsoft365wp_browser_site_list",
		PriorState:       dyn(prior),
		ProposedNewState: dyn(proposed),
		Config:           dyn(build(configSites, configCookies, false)),
	})
	if err != nil || len(resp.Diagnostics) > 0 {
		t.Fatalf("PlanResourceChange: %v %v", err, resp.Diagnostics)
	}
	planned, _ = resp.PlannedState.Unmarshal(objType)
	for _, name := range []string{"sites", "shared_cookies", "revision", "status", "published_date_time"} {
		got, _, _ := tftypes.WalkAttributePath(planned, tftypes.NewAttributePath().WithAttributeName(name))
		want, _, _ := tftypes.WalkAttributePath(prior, tftypes.NewAttributePath().WithAttributeName(name))
		if name == "sites" || name == "shared_cookies" {
			if !sameSetElements(got.(tftypes.Value), want.(tftypes.Value)) {
				t.Errorf("%s: expected no changes, got\n%s\nwant\n%s", name, got, want)
			}
		} else if !got.(tftypes.Value).Equal(want.(tftypes.Value)) {
			t.Errorf("%s: expected no changes, got %s, want %s", name, got, want)
		}
	}
}

// sameSetElements compares two set values independent of the order of their elements.
func sameSetElements(a, b tftypes.Value) bool {
	var as, bs []tftypes.Value
	if a.As(&as) != nil || b.As(&bs) != nil || len(as) != len(bs) {
		return false
	}
	for _, x := range as {
		found := false
		for _, y := range bs {
			if x.Equal(y) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
