package services

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func directorySettingTestSchema(t *testing.T) schema.Schema {
	t.Helper()
	resp := resource.SchemaResponse{}
	DirectorySettingResource.Schema(context.Background(), resource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("schema: %v", resp.Diagnostics)
	}
	return resp.Schema
}

func directorySettingTestValue(t *testing.T, s schema.Schema, id any, templateId string, values map[string]string) tftypes.Value {
	t.Helper()
	return directorySettingTestValueResetOnly(t, s, id, templateId, values, false)
}

func directorySettingTestValueResetOnly(t *testing.T, s schema.Schema, id any, templateId string, values map[string]string, resetOnly bool) tftypes.Value {
	t.Helper()
	objType := s.Type().TerraformType(context.Background()).(tftypes.Object)
	valuesTf := map[string]tftypes.Value{}
	for k, v := range values {
		valuesTf[k] = tftypes.NewValue(tftypes.String, v)
	}
	idVal := tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	if id != nil {
		idVal = tftypes.NewValue(tftypes.String, id)
	}
	return tftypes.NewValue(objType, map[string]tftypes.Value{
		"id":                                 idVal,
		"display_name":                       tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"template_id":                        tftypes.NewValue(tftypes.String, templateId),
		"values":                             tftypes.NewValue(objType.AttributeTypes["values"], valuesTf),
		"destroy_resets_managed_values_only": tftypes.NewValue(tftypes.Bool, resetOnly),
	})
}

func directorySettingTestValues(t *testing.T, val tftypes.Value) map[string]string {
	t.Helper()
	var attrs map[string]tftypes.Value
	if err := val.As(&attrs); err != nil {
		t.Fatalf("As object: %v", err)
	}
	var valuesTf map[string]tftypes.Value
	if err := attrs["values"].As(&valuesTf); err != nil {
		t.Fatalf("As map: %v", err)
	}
	result := map[string]string{}
	for k, v := range valuesTf {
		var s string
		if err := v.As(&s); err != nil {
			t.Fatalf("As string: %v", err)
		}
		result[k] = s
	}
	return result
}

func TestDirectorySettingTerraformToGraph(t *testing.T) {
	ctx := context.Background()
	s := directorySettingTestSchema(t)
	diags := diag.Diagnostics{}

	plan := directorySettingTestValueResetOnly(t, s, nil, "tmpl", map[string]string{"B": "2", "A": "true"}, true)
	raw := generic.ConvertTerraformToOdataRaw(ctx, &diags, s, false, plan, false, tfsdk.Config{}, directorySettingTerraformToGraphMiddleware, "Plan")
	if diags.HasError() {
		t.Fatalf("diags: %v", diags)
	}

	got, _ := json.Marshal(raw) // destroy_resets_managed_values_only is Terraform-only
	want := `{"templateId":"tmpl","values":[{"name":"A","value":"true"},{"name":"B","value":"2"}]}`
	if string(got) != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestDirectorySettingGraphToTerraform(t *testing.T) {
	ctx := context.Background()
	s := directorySettingTestSchema(t)

	graphRaw := func() map[string]any {
		return map[string]any{
			"id":          "s1",
			"displayName": "Group.Unified",
			"templateId":  "tmpl",
			"values": []any{
				map[string]any{"name": "EnableGroupCreation", "value": "false"},
				map[string]any{"name": "UsageGuidelinesUrl", "value": "https://example.invalid"},
				map[string]any{"name": "ClassificationList", "value": nil},
			},
		}
	}

	tests := []struct {
		name       string
		priorState map[string]string // nil: no prior state (import)
		want       map[string]string
	}{
		{"no prior state keeps all values", nil,
			map[string]string{"EnableGroupCreation": "false", "UsageGuidelinesUrl": "https://example.invalid", "ClassificationList": ""}},
		{"prior state filters values and keeps its spelling", map[string]string{"enablegroupcreation": "true", "NotInTemplate": "x"},
			map[string]string{"enablegroupcreation": "false"}},
		{"empty prior state map", map[string]string{}, map[string]string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := diag.Diagnostics{}
			raw := graphRaw()
			var reqState *tfsdk.State
			if tt.priorState != nil {
				reqState = &tfsdk.State{Schema: s, Raw: directorySettingTestValue(t, s, "s1", "tmpl", tt.priorState)}
			}
			directorySettingReadExtraRequestCustom(ctx, &diags, generic.ReadExtraRequestCustomParams{ReqState: reqState, RawVal: raw})
			val := generic.ConvertOdataRawToTerraform(ctx, &diags, s, raw, "", directorySettingGraphToTerraformMiddleware, "s1", false, reqState)
			if diags.HasError() {
				t.Fatalf("diags: %v", diags)
			}
			if got := directorySettingTestValues(t, val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func directorySettingTestClient(t *testing.T, handler http.HandlerFunc) *msgraph.Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client := msgraph.NewClient(msgraph.VersionBeta)
	client.Endpoint = server.URL
	client.RetryableClient.RetryWaitMin = 10 * time.Millisecond
	client.RetryableClient.RetryWaitMax = 20 * time.Millisecond
	return &client
}

func TestDirectorySettingImportModifyFunc(t *testing.T) {
	client := directorySettingTestClient(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet || req.URL.Path != "/beta/settings" {
			t.Errorf("unexpected request: %s %s", req.Method, req.URL)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"value":[
			{"id":"setting-1","templateId":"62375ab9-6b52-47ed-826b-58e47e0e304b","values":[]},
			{"id":"setting-2","templateId":"5cf42378-d67d-4f36-ba46-e8b86229381d","values":[]}]}`)
	})

	tests := []struct {
		importId string
		wantId   string // "" means error
	}{
		{"62375ab9-6b52-47ed-826b-58e47e0e304b", "setting-1"},
		{" 5CF42378-D67D-4F36-BA46-E8B86229381D ", "setting-2"},
		{"setting-2", "setting-2"},
		{"08d542b9-071f-4e16-94b0-74abb372e3d9", ""},
	}
	for _, tt := range tests {
		diags := diag.Diagnostics{}
		params := generic.ImportModifyFuncParams{Client: client, Id: tt.importId}
		directorySettingImportModifyFunc(context.Background(), &diags, &params)
		if tt.wantId == "" {
			if !diags.HasError() {
				t.Errorf("%q: expected an error, got id %q", tt.importId, params.Id)
			}
			continue
		}
		if diags.HasError() || params.Id != tt.wantId {
			t.Errorf("%q: got id %q (diags %v), want %q", tt.importId, params.Id, diags, tt.wantId)
		}
	}
}

func TestDirectorySettingWaitForValues(t *testing.T) {
	oldTimeout, oldInterval := directorySettingWaitTimeout, directorySettingWaitPollInterval
	directorySettingWaitTimeout, directorySettingWaitPollInterval = 300*time.Millisecond, time.Millisecond
	t.Cleanup(func() { directorySettingWaitTimeout, directorySettingWaitPollInterval = oldTimeout, oldInterval })

	const stale = `{"id":"s1","values":[{"name":"AllowToAddGuests","value":"true"},{"name":"GuestUsageGuidelinesUrl","value":"a"},{"name":"UsageGuidelinesUrl","value":null}]}`
	const fresh = `{"id":"s1","values":[{"name":"AllowToAddGuests","value":"True"},{"name":"GuestUsageGuidelinesUrl","value":"b"},{"name":"UsageGuidelinesUrl","value":null}]}`
	sent := []any{
		map[string]any{"name": "allowToAddGuests", "value": "True"},
		map[string]any{"name": "GuestUsageGuidelinesUrl", "value": "b"},
		map[string]any{"name": "UsageGuidelinesUrl", "value": ""},
	}

	for _, tc := range []struct {
		name        string
		responses   []string // per GET, the last one repeats; "404" for not found
		wantGets    int32    // 0: don't check
		wantWarning bool
	}{
		{"fresh immediately", []string{fresh}, 2, false},
		{"stale, not found, then fresh", []string{stale, "404", stale, fresh}, 5, false},
		{"fresh reads must be consecutive", []string{fresh, stale, fresh}, 4, false},
		{"stays stale", []string{stale}, 0, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var gets atomic.Int32
			client := directorySettingTestClient(t, func(w http.ResponseWriter, req *http.Request) {
				if req.Method != http.MethodGet || req.URL.Path != "/beta/settings/s1" {
					t.Errorf("unexpected request: %s %s", req.Method, req.URL)
				}
				i := int(gets.Add(1)) - 1
				body := tc.responses[min(i, len(tc.responses)-1)]
				w.Header().Set("Content-Type", "application/json")
				if body == "404" {
					w.WriteHeader(http.StatusNotFound)
					_, _ = io.WriteString(w, `{"error":{"code":"Request_ResourceNotFound","message":"not found"}}`)
					return
				}
				_, _ = io.WriteString(w, body)
			})
			client.DisableRetries = true // let the wait loop handle the 404 itself

			diags := diag.Diagnostics{}
			directorySettingWaitForValues(context.Background(), &diags, client, msgraph.Uri{Entity: "/settings/s1"}, sent)

			if diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}
			if gotWarning := diags.WarningsCount() > 0; gotWarning != tc.wantWarning {
				t.Errorf("warning: got %v, want %v (diags %v)", gotWarning, tc.wantWarning, diags)
			}
			if tc.wantGets > 0 && gets.Load() != tc.wantGets {
				t.Errorf("expected %d GETs, got %d", tc.wantGets, gets.Load())
			}
		})
	}
}

func TestDirectorySettingDeleteRetriesNotFound(t *testing.T) {
	for _, tc := range []struct {
		name        string
		notFounds   int32
		wantWarning bool
	}{
		{"404 resolves", 2, false},
		{"404 persists", 1000, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var deletes atomic.Int32
			client := directorySettingTestClient(t, func(w http.ResponseWriter, req *http.Request) {
				if req.Method != http.MethodDelete || req.URL.Path != "/beta/settings/s1" {
					t.Errorf("unexpected request: %s %s", req.Method, req.URL)
				}
				if deletes.Add(1) <= tc.notFounds {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusNotFound)
					_, _ = io.WriteString(w, `{"error":{"code":"Request_ResourceNotFound","message":"not found"}}`)
					return
				}
				w.WriteHeader(http.StatusNoContent)
			})
			if tc.wantWarning {
				client.RetryableClient.RetryMax = 3
			}

			r := &generic.GenericResource{AccessParams: generic.AccessParams{BaseUri: "/settings"}}
			diags := diag.Diagnostics{}
			params := generic.DeleteReplaceFuncParams{R: r, Client: client, Id: "s1"}
			directorySettingDeleteReplaceFunc(context.Background(), &diags, &params)

			if diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}
			if gotWarning := diags.WarningsCount() > 0; gotWarning != tc.wantWarning {
				t.Errorf("warning: got %v, want %v (diags %v)", gotWarning, tc.wantWarning, diags)
			}
			if !tc.wantWarning && deletes.Load() != 3 {
				t.Errorf("expected 3 DELETEs, got %d", deletes.Load())
			}
		})
	}
}

func TestDirectorySettingCreateRetriesTemplateConflict(t *testing.T) {
	const conflict = `{"error":{"code":"Request_BadRequest","message":"Another object with the same value for property templateId already exists.","details":[{"code":"ObjectConflict","target":"templateId"}]}}`

	for _, tc := range []struct {
		name      string
		conflicts int32
		wantError bool
	}{
		{"conflict resolves", 2, false},
		{"conflict persists", 1000, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var posts atomic.Int32
			client := directorySettingTestClient(t, func(w http.ResponseWriter, req *http.Request) {
				if req.Method != http.MethodPost || req.URL.Path != "/beta/settings" {
					t.Errorf("unexpected request: %s %s", req.Method, req.URL)
				}
				w.Header().Set("Content-Type", "application/json")
				if posts.Add(1) <= tc.conflicts {
					w.WriteHeader(http.StatusBadRequest)
					_, _ = io.WriteString(w, conflict)
					return
				}
				w.WriteHeader(http.StatusCreated)
				_, _ = io.WriteString(w, `{"id":"new-id","templateId":"tmpl","values":[]}`)
			})
			if tc.wantError {
				// stop retrying quickly: the deadline is directorySettingRetryDuration, so limit the attempts instead
				client.RetryableClient.RetryMax = 3
			}

			r := &generic.GenericResource{AccessParams: generic.AccessParams{BaseUri: "/settings"}}
			diags := diag.Diagnostics{}
			params := generic.CreateReplaceFuncParams{R: r, Client: client, RawVal: map[string]any{"templateId": "tmpl", "values": []any{}}}
			directorySettingCreateReplaceFunc(context.Background(), &diags, &params)

			if tc.wantError {
				if !diags.HasError() || !strings.Contains(diags.Errors()[0].Detail(), "import {") {
					t.Errorf("expected an error with import hint, got %v", diags)
				}
				return
			}
			if diags.HasError() || params.Id != "new-id" {
				t.Errorf("got id %q (diags %v), want new-id", params.Id, diags)
			}
			if got := posts.Load(); got != 3 {
				t.Errorf("expected 3 POSTs, got %d", got)
			}
		})
	}
}

func TestDirectorySettingReadSetsResetOnlyFromState(t *testing.T) {
	ctx := context.Background()
	s := directorySettingTestSchema(t)

	for _, tc := range []struct {
		name       string
		priorState *bool // nil: no prior state (import)
		want       bool
	}{
		{"no prior state", nil, false},
		{"prior state false", new(false), false},
		{"prior state true", new(true), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var reqState *tfsdk.State
			if tc.priorState != nil {
				reqState = &tfsdk.State{Schema: s, Raw: directorySettingTestValueResetOnly(t, s, "s1", "tmpl", map[string]string{}, *tc.priorState)}
			}
			raw := map[string]any{"id": "s1", "templateId": "tmpl", "values": []any{}}
			diags := diag.Diagnostics{}
			directorySettingReadExtraRequestCustom(ctx, &diags, generic.ReadExtraRequestCustomParams{ReqState: reqState, RawVal: raw})
			if diags.HasError() {
				t.Fatalf("diags: %v", diags)
			}
			if got := raw["destroyResetsManagedValuesOnly"]; got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDirectorySettingDeleteResetsManagedValuesOnly(t *testing.T) {
	oldInterval := directorySettingWaitPollInterval
	directorySettingWaitPollInterval = time.Millisecond
	t.Cleanup(func() { directorySettingWaitPollInterval = oldInterval })

	var patchBody atomic.Value
	client := directorySettingTestClient(t, func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
		}
		switch {
		case req.Method == http.MethodGet && req.URL.Path == "/beta/directorySettingTemplates/tmpl":
			_, _ = io.WriteString(w, `{"id":"tmpl","displayName":"Group.Unified","values":[
				{"name":"EnableGroupCreation","type":"System.Boolean","defaultValue":"true"},
				{"name":"AllowToAddGuests","type":"System.Boolean","defaultValue":"true"},
				{"name":"UsageGuidelinesUrl","type":"System.String","defaultValue":""}]}`)
		case req.Method == http.MethodPatch && req.URL.Path == "/beta/settings/s1":
			body, _ := io.ReadAll(req.Body)
			patchBody.Store(string(body))
			w.WriteHeader(http.StatusNoContent)
		case req.Method == http.MethodGet && req.URL.Path == "/beta/settings/s1":
			_, _ = io.WriteString(w, `{"id":"s1","templateId":"tmpl","values":[
				{"name":"EnableGroupCreation","value":"true"},
				{"name":"AllowToAddGuests","value":"false"},
				{"name":"UsageGuidelinesUrl","value":""}]}`)
		default:
			t.Errorf("unexpected request: %s %s", req.Method, req.URL)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	s := directorySettingTestSchema(t)
	reqState := &tfsdk.State{Schema: s, Raw: directorySettingTestValueResetOnly(t, s, "s1", "tmpl",
		map[string]string{"enablegroupcreation": "false", "UsageGuidelinesUrl": "https://example.invalid", "NotInTemplate": "x"}, true)}

	r := &generic.GenericResource{AccessParams: generic.AccessParams{BaseUri: "/settings"}}
	diags := diag.Diagnostics{}
	params := generic.DeleteReplaceFuncParams{R: r, Client: client, Id: "s1", ReqState: reqState}
	directorySettingDeleteReplaceFunc(context.Background(), &diags, &params)

	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	// AllowToAddGuests is not managed and must not be sent; names get the spelling of the template
	want := `{"values":[{"name":"UsageGuidelinesUrl","value":""},{"name":"EnableGroupCreation","value":"true"}]}`
	if got, _ := patchBody.Load().(string); got != want {
		t.Errorf("PATCH body: got %s, want %s", got, want)
	}
	if diags.WarningsCount() != 1 || !strings.Contains(diags.Warnings()[0].Detail(), "NotInTemplate") {
		t.Errorf("expected one warning about NotInTemplate, got %v", diags)
	}
}
