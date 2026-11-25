package services

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// TestNamedLocationRequiresReplaceOnTypeChange checks that a named location is only replaced if its type changes, not
// if the attributes of its type change (the attributes of the other types are unknown in the plan then).
func TestNamedLocationRequiresReplaceOnTypeChange(t *testing.T) {
	attrTypes := map[string]attr.Type{"is_trusted": types.BoolType}
	set := types.ObjectValueMust(attrTypes, map[string]attr.Value{"is_trusted": types.BoolValue(true)})
	changed := types.ObjectValueMust(attrTypes, map[string]attr.Value{"is_trusted": types.BoolValue(false)})
	null := types.ObjectNull(attrTypes)
	unknown := types.ObjectUnknown(attrTypes)

	// RequiresReplaceIf skips creates and destroys, so state and plan of the resource itself must not be null
	resourceRaw := tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{})

	for _, tc := range []struct {
		name        string
		state, plan types.Object
		want        bool
	}{
		{"same type, changed attributes", set, changed, false},
		{"other type, not configured (unknown in plan)", null, unknown, false},
		{"type is set now", null, set, true},
		{"type is removed (explicitly null)", set, null, true},
		{"type is removed (unknown as it is computed)", set, unknown, false}, // the newly set type triggers the replace
		{"unchanged", set, set, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := planmodifier.ObjectRequest{
				StateValue: tc.state,
				PlanValue:  tc.plan,
				State:      tfsdk.State{Raw: resourceRaw},
				Plan:       tfsdk.Plan{Raw: resourceRaw},
			}
			resp := &planmodifier.ObjectResponse{PlanValue: tc.plan}
			namedLocationRequiresReplaceOnTypeChange().PlanModifyObject(context.Background(), req, resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %v", resp.Diagnostics)
			}
			if resp.RequiresReplace != tc.want {
				t.Errorf("RequiresReplace = %v, want %v", resp.RequiresReplace, tc.want)
			}
		})
	}
}

func namedLocationTestResource(t *testing.T, handler http.HandlerFunc) *generic.GenericResource {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client := msgraph.NewClient(msgraph.VersionBeta)
	client.Endpoint = server.URL
	client.RetryableClient.RetryWaitMin = time.Millisecond
	client.RetryableClient.RetryWaitMax = 2 * time.Millisecond
	r := &generic.GenericResource{
		TypeNameSuffix: NamedLocationResource.TypeNameSuffix,
		SpecificSchema: namedLocationResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:      NamedLocationResource.AccessParams.BaseUri,
			WriteOptions: generic.WriteOptions{WriteRetryFor: time.Second, WriteRetryFunc: namedLocationWriteRetryFunc},
		},
	}
	r.AccessParams.InitializeGuarded(&client)
	return r
}

func TestNamedLocationWriteRetryFunc(t *testing.T) {
	const notFoundBody = `{"error":{"code":"ResourceNotFound","message":"NamedLocation with id x does not exist in the directory."}}`
	for _, tc := range []struct {
		status int
		body   string
		want   bool
		method string // default: DELETE
	}{
		{http.StatusNotFound, notFoundBody, true, http.MethodPatch}, // e.g. untrusting a named location that has just been created
		{http.StatusNotFound, notFoundBody, false, http.MethodDelete},
		{http.StatusNotFound, notFoundBody, false, http.MethodPost},
		{http.StatusForbidden, `{"error":{"code":"AccessDenied","message":"Access denied. CheckAccess returned 'NotAllowed' for action 'microsoft.directory/namedLocations/delete'."}}`, true, ""},
		{http.StatusForbidden, `{"error":{"code":"Authorization_RequestDenied","message":"Insufficient privileges to complete the operation."}}`, false, ""},
		{http.StatusBadRequest, `{"error":{"code":"BadRequest","message":"malformed","innerError":{"message":"1177: Named location: 'x' cannot be deleted because it is marked as a Trusted location."}}}`, true, ""},
		{http.StatusBadRequest, `{"error":{"code":"BadRequest","message":"malformed","innerError":{"message":"1178: Named location: 'x' cannot be deleted because it is referenced by one or more Conditional Access policies."}}}`, true, ""},
		{http.StatusBadRequest, `{"error":{"code":"BadRequest","message":"malformed","innerError":{"message":"1172: Deletion of 'Compliant Network' named location is not allowed."}}}`, false, ""},
		{http.StatusBadRequest, `{"error":{"code":"BadRequest","message":"malformed","innerError":{"message":"1042: Unexpected type: namedLocationType for NamedLocations."}}}`, false, ""},
	} {
		method := tc.method
		if method == "" {
			method = http.MethodDelete
		}
		req, _ := http.NewRequest(method, "https://graph.microsoft.com/beta/identity/conditionalAccess/namedLocations/x", nil)
		resp := &http.Response{StatusCode: tc.status, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(tc.body)), Request: req}
		o, err := odata.FromResponse(resp)
		if err != nil {
			t.Fatalf("odata.FromResponse: %v", err)
		}
		if got := namedLocationWriteRetryFunc(resp, o); got != tc.want {
			t.Errorf("%s %d %s: got %v, want %v", method, tc.status, tc.body, got, tc.want)
		}
	}
}

func TestNamedLocationFindCompliantNetwork(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
		want string
	}{
		{"exists", `{"value":[
			{"@odata.type":"#microsoft.graph.ipNamedLocation","id":"ip1","displayName":"x"},
			{"@odata.type":"#microsoft.graph.compliantNetworkNamedLocation","id":"cn1","displayName":"All Compliant Network locations"}]}`, "cn1"},
		{"does not exist", `{"value":[{"@odata.type":"#microsoft.graph.ipNamedLocation","id":"ip1","displayName":"x"}]}`, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := namedLocationTestResource(t, func(w http.ResponseWriter, req *http.Request) {
				if req.Method != http.MethodGet || req.URL.Path != "/beta/identity/conditionalAccess/namedLocations" {
					t.Errorf("unexpected request: %s %s", req.Method, req.URL)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = io.WriteString(w, tc.body)
			})
			diags := diag.Diagnostics{}
			if got := namedLocationFindCompliantNetwork(context.Background(), &diags, r); diags.HasError() || got != tc.want {
				t.Errorf("got %q (diags %v), want %q", got, diags, tc.want)
			}
		})
	}
}

// namedLocationTestState returns a state of the named location resource with the given attributes set.
func namedLocationTestState(t *testing.T, attributes map[string]any) tfsdk.State {
	t.Helper()
	ctx := context.Background()
	state := tfsdk.State{Schema: namedLocationResourceSchema, Raw: tftypes.NewValue(namedLocationResourceSchema.Type().TerraformType(ctx), nil)}
	for p, v := range attributes {
		steps := strings.Split(p, ".")
		attrPath := path.Root(steps[0])
		for _, s := range steps[1:] {
			attrPath = attrPath.AtName(s)
		}
		if diags := state.SetAttribute(ctx, attrPath, v); diags.HasError() {
			t.Fatalf("SetAttribute(%s): %v", p, diags)
		}
	}
	return state
}

func TestNamedLocationUntrustBeforeDeleteWsa(t *testing.T) {
	for _, tc := range []struct {
		name       string
		attributes map[string]any
		wantBody   string // "" means no request
	}{
		{"trusted ip location", map[string]any{"display_name": "x", "ip.is_trusted": true}, `{"@odata.type":"#microsoft.graph.ipNamedLocation","isTrusted":false}`},
		{"untrusted ip location", map[string]any{"display_name": "x", "ip.is_trusted": false}, ""},
		{"trusted compliant network location", map[string]any{"display_name": namedLocationCompliantNetworkDisplayName, "compliant_network.is_trusted": true},
			`{"@odata.type":"#microsoft.graph.compliantNetworkNamedLocation","isTrusted":false}`},
		{"country location", map[string]any{"display_name": "x", "country.include_unknown_countries_and_regions": true}, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var gotBody string
			r := namedLocationTestResource(t, func(w http.ResponseWriter, req *http.Request) {
				if req.Method != http.MethodPatch || req.URL.Path != "/beta/identity/conditionalAccess/namedLocations/l1" {
					t.Errorf("unexpected request: %s %s", req.Method, req.URL)
				}
				b, _ := io.ReadAll(req.Body)
				gotBody = string(b)
				w.WriteHeader(http.StatusNoContent)
			})
			state := namedLocationTestState(t, tc.attributes)
			diags := diag.Diagnostics{}
			wsaReq := &generic.WriteSubActionRequest{GenRes: r, Operation: generic.OperationDelete, Id: "l1", ReqState: &state}
			(&namedLocationUntrustBeforeDeleteWsa{}).ExecutePre(context.Background(), &diags, wsaReq)
			if diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}
			if tc.wantBody == "" && gotBody != "" {
				t.Errorf("unexpected PATCH with body %s", gotBody)
			}
			if tc.wantBody != "" {
				var got, want any
				_ = json.Unmarshal([]byte(gotBody), &got)
				_ = json.Unmarshal([]byte(tc.wantBody), &want)
				if !reflect.DeepEqual(got, want) {
					t.Errorf("got body %s, want %s", gotBody, tc.wantBody)
				}
			}
		})
	}
}

func TestNamedLocationCompliantNetworkDisplayNameValidator(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		name       string
		attributes map[string]any
		wantError  bool
	}{
		{"compliant network with the fixed name", map[string]any{"display_name": namedLocationCompliantNetworkDisplayName, "compliant_network.is_trusted": true}, false},
		{"compliant network with another name", map[string]any{"display_name": "All Compliant Network locations - TF", "compliant_network.is_trusted": true}, true},
		{"other type with any name", map[string]any{"display_name": "Teams SIP Gateways - TF", "ip.is_trusted": true}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			state := namedLocationTestState(t, tc.attributes)
			req := resource.ValidateConfigRequest{Config: tfsdk.Config{Schema: state.Schema, Raw: state.Raw}}
			resp := &resource.ValidateConfigResponse{}
			namedLocationCompliantNetworkDisplayNameValidator{}.ValidateResource(ctx, req, resp)
			if resp.Diagnostics.HasError() != tc.wantError {
				t.Errorf("got error %v, want %v (diags %v)", resp.Diagnostics.HasError(), tc.wantError, resp.Diagnostics)
			}
		})
	}
}

func TestNamedLocationMoveState(t *testing.T) {
	const id = "0a1b2c3d-4e5f-4a6b-8c7d-9e0f1a2b3c4d"
	tests := []struct {
		name            string
		providerAddress string
		typeName        string
		schemaVersion   int64
		json            string
		wantMoved       bool
		wantError       string
	}{
		{name: "azuread v3 fork", providerAddress: "registry.opentofu.org/glueckkanja/azuread", typeName: "azuread_named_location", schemaVersion: 1,
			json: `{"id":"/identity/conditionalAccess/namedLocations/` + id + `","display_name":"Teams SIP Gateways","ip":[{"ip_ranges":["192.0.2.0/24"],"trusted":false}],"country":[]}`, wantMoved: true},
		{name: "azuread v2", providerAddress: "registry.terraform.io/hashicorp/azuread", typeName: "azuread_named_location", schemaVersion: 0,
			json: `{"id":"` + id + `","display_name":"Trusted"}`, wantMoved: true},
		{name: "azuread unexpected id", providerAddress: "registry.terraform.io/hashicorp/azuread", typeName: "azuread_named_location", schemaVersion: 1,
			json: `{"id":"/policies/` + id + `"}`, wantError: "unexpected id"},
		{name: "azuread other resource type", providerAddress: "registry.terraform.io/hashicorp/azuread", typeName: "azuread_conditional_access_policy",
			json: `{"id":"` + id + `"}`},
		{name: "msgraph named location", providerAddress: "registry.terraform.io/microsoft/msgraph", typeName: "msgraph_resource",
			json: `{"id":"` + id + `","url":"identity/conditionalAccess/namedLocations","api_version":"beta","resource_url":"identity/conditionalAccess/namedLocations/` + id + `"}`, wantMoved: true},
		{name: "msgraph named location with slashes", providerAddress: "registry.opentofu.org/microsoft/msgraph", typeName: "msgraph_resource",
			json: `{"id":"` + id + `","url":"/identity/conditionalaccess/namedLocations/"}`, wantMoved: true},
		{name: "msgraph other url", providerAddress: "registry.terraform.io/microsoft/msgraph", typeName: "msgraph_resource",
			json: `{"id":"` + id + `","url":"identity/conditionalAccess/policies"}`, wantError: `not one with url "identity/conditionalAccess/policies"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			r := &generic.GenericResource{SpecificSchema: namedLocationResourceSchema, MoveStateSources: NamedLocationResource.MoveStateSources}
			req := resource.MoveStateRequest{
				SourceProviderAddress: tt.providerAddress,
				SourceTypeName:        tt.typeName,
				SourceSchemaVersion:   tt.schemaVersion,
				SourceRawState:        &tfprotov6.RawState{JSON: []byte(tt.json)},
			}
			nullState := tftypes.NewValue(namedLocationResourceSchema.Type().TerraformType(ctx), nil)
			var resp resource.MoveStateResponse
			moved := false
			for _, sm := range r.MoveState(ctx) {
				resp = resource.MoveStateResponse{TargetState: tfsdk.State{Schema: namedLocationResourceSchema, Raw: nullState}}
				sm.StateMover(ctx, req, &resp)
				if moved = !resp.TargetState.Raw.Equal(nullState); moved || resp.Diagnostics.HasError() {
					break
				}
			}

			if tt.wantError != "" {
				if !resp.Diagnostics.HasError() || !strings.Contains(resp.Diagnostics.Errors()[0].Detail(), tt.wantError) {
					t.Fatalf("want error containing %q, got %v", tt.wantError, resp.Diagnostics)
				}
				return
			}
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %v", resp.Diagnostics)
			}
			if moved != tt.wantMoved {
				t.Fatalf("moved = %v, want %v", moved, tt.wantMoved)
			}
			if !moved {
				return
			}
			want := namedLocationTestState(t, map[string]any{"id": id})
			if !resp.TargetState.Raw.Equal(want.Raw) {
				t.Errorf("target state = %v, want %v", resp.TargetState.Raw, want.Raw)
			}
		})
	}
}
