package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const priorityTestID = "11111111-1111-1111-1111-111111111111"

type priorityTestTransport func(*http.Request) (*http.Response, error)

func (f priorityTestTransport) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func priorityTestResource(t *testing.T, handler priorityTestTransport) *generic.GenericResource {
	t.Helper()
	source := &DeviceEnrollmentConfigurationResource
	r := &generic.GenericResource{
		TypeNameSuffix: source.TypeNameSuffix,
		SpecificSchema: source.SpecificSchema,
		AccessParams: generic.AccessParams{
			BaseUri: source.AccessParams.BaseUri, ReadOptions: source.AccessParams.ReadOptions,
			WriteOptions:               generic.WriteOptions{SubActions: source.AccessParams.WriteOptions.SubActions},
			TerraformToGraphMiddleware: source.AccessParams.TerraformToGraphMiddleware,
			GraphToTerraformMiddleware: source.AccessParams.GraphToTerraformMiddleware,
			CreateModifyFunc:           source.AccessParams.CreateModifyFunc,
			DeleteModifyFunc:           source.AccessParams.DeleteModifyFunc,
		},
	}
	client := msgraph.NewClient(msgraph.VersionBeta)
	client.HttpClient = &http.Client{Transport: handler}
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: &client}, &resource.ConfigureResponse{})
	var response resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &response)
	if response.Diagnostics.HasError() {
		t.Fatal(response.Diagnostics)
	}
	return r
}

func priorityTestResponse(req *http.Request, status int, body any) (*http.Response, error) {
	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(string(encoded))), Request: req}, nil
}

func priorityTestState(t *testing.T, r *generic.GenericResource) tfsdk.State {
	t.Helper()
	raw := map[string]any{
		"@odata.type": "#microsoft.graph.deviceEnrollmentLimitConfiguration",
		"id":          priorityTestID, "displayName": "Offline enrollment", "priority": float64(1),
		"limit": float64(2), "assignments": []any{},
	}
	var diags diag.Diagnostics
	value := generic.ConvertOdataRawToTerraform(context.Background(), &diags, r.SpecificSchema, raw, "", r.AccessParams.GraphToTerraformMiddleware, priorityTestID, false, nil)
	if diags.HasError() {
		t.Fatal(diags)
	}
	return tfsdk.State{Schema: r.SpecificSchema, Raw: value}
}

func priorityTestSet(t *testing.T, state *tfsdk.State, name string, value any) {
	t.Helper()
	if diags := state.SetAttribute(context.Background(), path.Root(name), value); diags.HasError() {
		t.Fatal(diags)
	}
}

func TestEnrollmentPriorityReadBeforeWrite(t *testing.T) {
	for _, create := range []bool{false, true} {
		for _, tc := range []struct {
			name                       string
			before, after              any
			readStatus, afterStatus    int
			actionStatus, assignStatus int
			want                       types.Int64
			wantError                  bool
			wantReads, wantActions     int
		}{
			{"already_matching", float64(2), nil, 200, 0, 0, 204, types.Int64Value(2), false, 1, 0},
			{"change_verified", float64(1), float64(2), 200, 200, 200, 204, types.Int64Value(2), false, 2, 1},
			{"action_rejected", float64(1), float64(1), 200, 200, 400, 204, types.Int64Value(1), true, 2, 1},
			{"action_rejected_after_change", float64(1), float64(2), 200, 200, 400, 204, types.Int64Value(2), true, 2, 1},
			{"change_not_applied", float64(1), float64(1), 200, 200, 204, 204, types.Int64Value(1), true, 2, 1},
			{"post_read_failed", float64(1), nil, 200, 403, 204, 204, types.Int64Value(1), true, 2, 1},
			{"pre_read_failed", nil, nil, 403, 0, 0, 204, types.Int64Null(), true, 1, 0},
			{"pre_read_invalid", float64(1.5), nil, 200, 0, 0, 204, types.Int64Null(), true, 1, 0},
			{"assignment_failed", nil, nil, 0, 0, 0, 400, types.Int64Null(), true, 0, 0},
		} {
			t.Run(fmt.Sprintf("create_%t/%s", create, tc.name), func(t *testing.T) {
				reads, actions := 0, 0
				var requests []string
				r := priorityTestResource(t, func(req *http.Request) (*http.Response, error) {
					requests = append(requests, req.Method+" "+req.URL.Path)
					switch {
					case req.Method == http.MethodGet:
						reads++
						status, observed := tc.readStatus, tc.before
						if reads == 2 {
							status, observed = tc.afterStatus, tc.after
						}
						if status != 200 {
							return priorityTestResponse(req, status, map[string]any{"error": map[string]any{"code": "Forbidden", "message": "synthetic failure"}})
						}
						return priorityTestResponse(req, status, map[string]any{"id": priorityTestID, "priority": observed})
					case strings.HasSuffix(req.URL.Path, "/setPriority"):
						actions++
						var body map[string]any
						if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
							t.Fatal(err)
						}
						if !reflect.DeepEqual(body, map[string]any{"priority": float64(2)}) {
							t.Fatal(body)
						}
						if tc.actionStatus != 204 && tc.actionStatus != 200 {
							return priorityTestResponse(req, tc.actionStatus, map[string]any{"error": map[string]any{"code": "BadRequest", "message": "synthetic rejection"}})
						}
						return priorityTestResponse(req, tc.actionStatus, nil)
					case strings.HasSuffix(req.URL.Path, "/assign"):
						if tc.assignStatus != 204 {
							return priorityTestResponse(req, tc.assignStatus, map[string]any{"error": map[string]any{"code": "BadRequest", "message": "synthetic assignment failure"}})
						}
						return priorityTestResponse(req, 204, nil)
					default:
						var body map[string]any
						if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
							t.Fatal(err)
						}
						if _, exists := body["priority"]; exists {
							t.Fatal("priority leaked into parent write")
						}
						if create {
							return priorityTestResponse(req, 201, map[string]any{"id": priorityTestID, "priority": tc.before})
						}
						return priorityTestResponse(req, 204, nil)
					}
				})
				state := priorityTestState(t, r)
				plan := state
				priorityTestSet(t, &plan, "priority", types.Int64Value(2))
				config := tfsdk.Config{Schema: plan.Schema, Raw: plan.Raw}
				var result tfsdk.State
				var diags diag.Diagnostics
				method := "PATCH"
				parentPath := "/beta/deviceManagement/deviceEnrollmentConfigurations/" + priorityTestID
				if create {
					method = "POST"
					parentPath = "/beta/deviceManagement/deviceEnrollmentConfigurations"
					priorityTestSet(t, &plan, "id", types.StringUnknown())
					response := resource.CreateResponse{State: tfsdk.State{Schema: r.SpecificSchema}}
					r.Create(context.Background(), resource.CreateRequest{Config: config, Plan: tfsdk.Plan{Schema: plan.Schema, Raw: plan.Raw}}, &response)
					result, diags = response.State, response.Diagnostics
				} else {
					response := resource.UpdateResponse{State: state}
					r.Update(context.Background(), resource.UpdateRequest{State: state, Config: config, Plan: tfsdk.Plan{Schema: plan.Schema, Raw: plan.Raw}}, &response)
					result, diags = response.State, response.Diagnostics
				}
				if diags.HasError() != tc.wantError || diags.WarningsCount() != 0 {
					t.Fatalf("diagnostics: %v", diags)
				}
				want := tc.want
				if !create && want.IsNull() {
					want = types.Int64Value(1)
				}
				var priority types.Int64
				var id types.String
				if d := result.GetAttribute(context.Background(), path.Root("priority"), &priority); d.HasError() || !priority.Equal(want) {
					t.Fatalf("priority=%v, want %v: %v", priority, want, d)
				}
				if d := result.GetAttribute(context.Background(), path.Root("id"), &id); d.HasError() || id.ValueString() != priorityTestID {
					t.Fatalf("id=%v: %v", id, d)
				}
				if reads != tc.wantReads || actions != tc.wantActions {
					t.Fatalf("reads=%d, actions=%d", reads, actions)
				}
				entity := "/beta/deviceManagement/deviceEnrollmentConfigurations/" + priorityTestID
				wantRequests := []string{method + " " + parentPath, "POST " + entity + "/assign"}
				if tc.wantReads > 0 {
					wantRequests = append(wantRequests, "GET "+entity)
				}
				if tc.wantActions > 0 {
					wantRequests = append(wantRequests, "POST "+entity+"/setPriority", "GET "+entity)
				}
				if !reflect.DeepEqual(requests, wantRequests) {
					t.Fatalf("requests=%v, want %v", requests, wantRequests)
				}
			})
		}
	}
}

func TestEnrollmentPriorityOmitted(t *testing.T) {
	for _, create := range []bool{false, true} {
		t.Run(fmt.Sprintf("create_%t", create), func(t *testing.T) {
			requests := 0
			r := priorityTestResource(t, func(req *http.Request) (*http.Response, error) {
				requests++
				if req.Method == http.MethodGet || strings.HasSuffix(req.URL.Path, "/setPriority") {
					t.Fatalf("omitted priority caused a request: %s %s", req.Method, req.URL.Path)
				}
				var body map[string]any
				if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if _, exists := body["priority"]; exists {
					t.Fatal("omitted priority leaked into write")
				}
				if create && !strings.HasSuffix(req.URL.Path, "/assign") {
					return priorityTestResponse(req, 201, map[string]any{"id": priorityTestID, "priority": float64(3)})
				}
				return priorityTestResponse(req, 204, nil)
			})
			state := priorityTestState(t, r)
			config := state
			priorityTestSet(t, &config, "priority", types.Int64Null())
			plan := state
			var result tfsdk.State
			var diags diag.Diagnostics
			want := types.Int64Value(1)
			if create {
				want = types.Int64Value(3)
				priorityTestSet(t, &plan, "id", types.StringUnknown())
				priorityTestSet(t, &plan, "priority", types.Int64Unknown())
				response := resource.CreateResponse{State: tfsdk.State{Schema: r.SpecificSchema}}
				r.Create(context.Background(), resource.CreateRequest{Config: tfsdk.Config{Schema: config.Schema, Raw: config.Raw}, Plan: tfsdk.Plan{Schema: plan.Schema, Raw: plan.Raw}}, &response)
				result, diags = response.State, response.Diagnostics
			} else {
				response := resource.UpdateResponse{State: state}
				r.Update(context.Background(), resource.UpdateRequest{State: state, Config: tfsdk.Config{Schema: config.Schema, Raw: config.Raw}, Plan: tfsdk.Plan{Schema: plan.Schema, Raw: plan.Raw}}, &response)
				result, diags = response.State, response.Diagnostics
			}
			if diags.HasError() || requests != 2 {
				t.Fatalf("diagnostics=%v, requests=%d", diags, requests)
			}
			var priority types.Int64
			if d := result.GetAttribute(context.Background(), path.Root("priority"), &priority); d.HasError() || !priority.Equal(want) {
				t.Fatalf("priority=%v, want %v: %v", priority, want, d)
			}
		})
	}
}

func TestEnrollmentPriorityUnknown(t *testing.T) {
	r := priorityTestResource(t, func(req *http.Request) (*http.Response, error) {
		t.Fatalf("unknown priority caused a request: %s %s", req.Method, req.URL.Path)
		return nil, nil
	})
	state := priorityTestState(t, r)
	config := state
	priorityTestSet(t, &config, "priority", types.Int64Unknown())
	response := resource.UpdateResponse{State: state}
	r.Update(context.Background(), resource.UpdateRequest{
		State: state, Config: tfsdk.Config{Schema: config.Schema, Raw: config.Raw},
		Plan: tfsdk.Plan{Schema: state.Schema, Raw: state.Raw},
	}, &response)
	if !response.Diagnostics.HasError() {
		t.Fatal("unknown priority was not rejected")
	}
}
