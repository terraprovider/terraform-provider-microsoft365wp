package services

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"terraform-provider-microsoft365wp/workplace/external/msgraph"
	"terraform-provider-microsoft365wp/workplace/generic"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestDeviceConfigurationAssignmentIntent(t *testing.T) {
	intent := newDeviceConfigurationAssignmentResource().SpecificSchema.Attributes["intent"].(schema.StringAttribute)
	for _, value := range []string{"apply", "remove"} {
		resp := validator.StringResponse{}
		for _, check := range intent.Validators {
			check.ValidateString(context.Background(), validator.StringRequest{Path: path.Root("intent"), ConfigValue: types.StringValue(value)}, &resp)
		}
		if resp.Diagnostics.HasError() != (value != "apply") {
			t.Fatalf("unexpected intent validation for %q: %v", value, resp.Diagnostics)
		}
	}
}

func TestDeviceConfigurationAssignmentRead(t *testing.T) {
	const assignment = `{"id":"policy_group","intent":"apply","source":"direct","sourceId":"","target":{"@odata.type":"#microsoft.graph.groupAssignmentTarget","groupId":"group","deviceAndAppManagementAssignmentFilterType":"none"}}`
	for _, tc := range []struct {
		name         string
		status       int
		body         string
		wantAbsent   bool
		wantError    bool
		selectFields []string
		wantSelect   string
	}{
		{"select missing ID", 200, `{"value":[` + assignment + `]}`, false, false, []string{"intent"}, "intent,id"},
		{"select with ID", 200, `{"value":[` + assignment + `]}`, false, false, []string{"id", "intent"}, "id,intent"},
		{"present", 200, `{"value":[` + assignment + `]}`, false, false, nil, ""},
		{"present with metadata", 200, `{"value":[` + assignment + `],"@odata.context":"context","@odata.count":1,"metadata":"extra"}`, false, false, nil, ""},
		{"absent", 200, `{"value":[]}`, true, false, nil, ""},
		{"absent with metadata", 200, `{"value":[],"@odata.context":"context","@odata.count":0,"metadata":"extra"}`, true, false, nil, ""},
		{"missing parent", 404, `{"error":{"code":"ResourceNotFound","message":"missing parent"}}`, true, false, nil, ""},
		{"forbidden", 403, `{"error":{"code":"Forbidden","message":"access denied"}}`, false, true, nil, ""},
		{"invalid JSON", 200, `{`, false, true, nil, ""},
		{"null JSON", 200, `null`, false, true, nil, ""},
		{"missing array", 200, `{}`, false, true, nil, ""},
		{"null array", 200, `{"value":null}`, false, true, nil, ""},
		{"multiple matches", 200, `{"value":[` + assignment + `,` + assignment + `]}`, false, true, nil, ""},
		{"nonobject match", 200, `{"value":[null]}`, false, true, nil, ""},
		{"wrong ID", 200, `{"value":[{"id":"other"}]}`, false, true, nil, ""},
		{"missing ID", 200, `{"value":[{}]}`, false, true, nil, ""},
		{"continuation", 200, `{"value":[],"@odata.nextLink":"CONTINUATION"}`, false, true, nil, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			var server *httptest.Server
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				calls++
				if req.Method != "GET" || req.URL.Path != "/beta/deviceManagement/deviceConfigurations/policy/assignments" || req.URL.Query().Get("$filter") != "id eq 'policy_group'" || req.URL.Query().Get("$select") != tc.wantSelect {
					t.Errorf("unexpected request: %s %s", req.Method, req.URL)
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				io.WriteString(w, strings.ReplaceAll(tc.body, "CONTINUATION", server.URL+"/unexpected-next"))
			}))
			defer server.Close()
			client := msgraph.NewClient(msgraph.VersionBeta)
			client.Endpoint = server.URL
			client.DisableRetries = true
			r := newDeviceConfigurationAssignmentResource()
			r.AccessParams.ReadOptions.ODataSelect = slices.Clone(tc.selectFields)
			ctx := context.Background()
			schema := resource.SchemaResponse{}
			r.Schema(ctx, resource.SchemaRequest{}, &schema)
			if schema.Diagnostics.HasError() {
				t.Fatal(schema.Diagnostics)
			}
			r.Configure(ctx, resource.ConfigureRequest{ProviderData: &client}, &resource.ConfigureResponse{})
			translator := generic.NewToFromGraphTranslator(schema.Schema, false, nil)
			raw, err := translator.TerraformFromJsonString(ctx, assignment)
			if err != nil {
				t.Fatal(err)
			}
			state := tfsdk.State{Schema: schema.Schema, Raw: raw}
			if diags := state.SetAttribute(ctx, path.Root("device_configuration_id"), "policy"); diags.HasError() {
				t.Fatal(diags)
			}
			resp := resource.ReadResponse{State: state}
			r.Read(ctx, resource.ReadRequest{State: state}, &resp)
			if !slices.Equal(r.AccessParams.ReadOptions.ODataSelect, tc.selectFields) {
				t.Fatal("read changed the configured select fields")
			}
			if resp.Diagnostics.HasError() != tc.wantError {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
			if resp.State.Raw.IsNull() != tc.wantAbsent {
				t.Fatalf("absent = %v, want %v", resp.State.Raw.IsNull(), tc.wantAbsent)
			}
			if tc.wantError && !resp.State.Raw.Equal(state.Raw) {
				t.Fatal("failed read changed the existing state")
			}
			if !tc.wantAbsent && !tc.wantError {
				var parentID, intent, id, groupID string
				resp.State.GetAttribute(ctx, path.Root("device_configuration_id"), &parentID)
				resp.State.GetAttribute(ctx, path.Root("intent"), &intent)
				resp.State.GetAttribute(ctx, path.Root("id"), &id)
				resp.State.GetAttribute(ctx, path.Root("target").AtName("group").AtName("group_id"), &groupID)
				if parentID != "policy" || intent != "apply" || id != "policy_group" || groupID != "group" {
					t.Fatalf("unexpected state parent %q / intent %q / id %q / group %q", parentID, intent, id, groupID)
				}
			}
			if calls != 1 {
				t.Fatalf("expected one collection request, got %d", calls)
			}
		})
	}
}
