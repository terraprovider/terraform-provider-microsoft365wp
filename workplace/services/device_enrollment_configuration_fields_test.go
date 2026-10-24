package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"terraform-provider-microsoft365wp/workplace/generic"

	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestEnrollmentESPFields(t *testing.T) {
	ctx := context.Background()
	r := &DeviceEnrollmentConfigurationResource
	esp := r.SpecificSchema.Attributes["windows10_esp"].(generic.OdataDerivedTypeNestedAttributeRs)
	espDS := DeviceEnrollmentConfigurationSingularDataSource.SpecificSchema.Attributes["windows10_esp"].(generic.OdataDerivedTypeNestedAttributeDs)
	for name, graphName := range map[string]string{
		"install_quality_updates":             "installQualityUpdates",
		"allow_non_blocking_app_installation": "allowNonBlockingAppInstallation",
	} {
		t.Run(name, func(t *testing.T) {
			attribute := esp.Attributes[name].(schema.BoolAttribute)
			if !attribute.Optional || !attribute.Computed || attribute.Default != nil || len(attribute.PlanModifiers) != 0 {
				t.Fatal("ESP field must be optional and computed without a forced default")
			}
			if attribute := espDS.Attributes[name].(dsschema.BoolAttribute); !attribute.Computed || attribute.Optional || attribute.Required {
				t.Fatal("ESP data source field must be computed")
			}
			settingPath := path.Root("windows10_esp").AtName(name)
			for surface, state := range map[string]tfsdk.State{
				"resource":    {Schema: r.SpecificSchema},
				"data_source": {Schema: DeviceEnrollmentConfigurationSingularDataSource.SpecificSchema},
			} {
				for _, value := range []types.Bool{types.BoolValue(true), types.BoolValue(false), types.BoolNull()} {
					t.Run(fmt.Sprintf("read/%s/%s", surface, value), func(t *testing.T) {
						raw := map[string]any{"@odata.type": esp.DerivedType, "priority": float64(1)}
						if !value.IsNull() {
							raw[graphName] = value.ValueBool()
						}
						var diags diag.Diagnostics
						state.Raw = generic.ConvertOdataRawToTerraform(ctx, &diags, state.Schema, raw, "", r.AccessParams.GraphToTerraformMiddleware, "", false, nil)
						if diags.HasError() {
							t.Fatal(diags)
						}
						var got types.Bool
						if diags := state.GetAttribute(ctx, settingPath, &got); diags.HasError() || !got.Equal(value) {
							t.Fatalf("read = %v, want %v; %v", got, value, diags)
						}
					})
				}
			}
			for _, update := range []bool{false, true} {
				for _, configured := range []types.Bool{types.BoolValue(true), types.BoolValue(false), types.BoolNull(), types.BoolUnknown()} {
					t.Run(fmt.Sprintf("update_%t/%s", update, configured), func(t *testing.T) {
						var diags diag.Diagnostics
						plan := tfsdk.State{Schema: r.SpecificSchema}
						plan.Raw = generic.ConvertOdataRawToTerraform(ctx, &diags, plan.Schema, map[string]any{
							"@odata.type": esp.DerivedType, "priority": float64(1), graphName: true,
						}, "", r.AccessParams.GraphToTerraformMiddleware, "", false, nil)
						if diags.HasError() {
							t.Fatal(diags)
						}
						config := plan
						if diags := config.SetAttribute(ctx, settingPath, configured); diags.HasError() {
							t.Fatal(diags)
						}
						if !configured.IsNull() && !configured.IsUnknown() {
							if diags := plan.SetAttribute(ctx, settingPath, configured); diags.HasError() {
								t.Fatal(diags)
							}
						}
						// An omitted setting may retain its prior value in Plan; only Config expresses write intent.
						raw := generic.ConvertTerraformToOdataRaw(ctx, &diags, plan.Schema, update, plan.Raw, update,
							tfsdk.Config{Schema: config.Schema, Raw: config.Raw}, r.AccessParams.TerraformToGraphMiddleware, "test")
						if diags.HasError() != configured.IsUnknown() {
							t.Fatalf("unexpected diagnostics: %v", diags)
						}
						if configured.IsUnknown() {
							if raw != nil {
								t.Fatal("unknown setting produced a write payload")
							}
							return
						}
						got, present := raw[graphName]
						if present != !configured.IsNull() || (present && got != configured.ValueBool()) {
							t.Fatalf("written value = %v, present = %t; configured = %v", got, present, configured)
						}
					})
				}
			}
		})
	}
}

func TestEnrollmentESPPluralKeepsTypeDiscriminator(t *testing.T) {
	ds := &DeviceEnrollmentConfigurationPluralDataSource
	list := ds.SpecificSchema.Attributes[ds.ResultAttributeName].(dsschema.ListNestedAttribute)
	esp := list.NestedObject.Attributes["windows10_esp"].(generic.OdataDerivedTypeNestedAttributeDs)
	if len(esp.Attributes) != 0 || esp.DerivedType != "#microsoft.graph.windows10EnrollmentCompletionPageConfiguration" {
		t.Fatal("plural data source must keep its existing ESP type discriminator")
	}
}

func TestEnrollmentESPFieldsSeparateWrite(t *testing.T) {
	for _, tc := range []struct {
		name                               string
		create, defaultProfile, omit       bool
		quality, nonblocking               bool
		omitNonblocking, requestPriority   bool
		preStatus, patchStatus, readStatus int
		malformed, mismatch                bool
	}{
		{name: "true_false", quality: true, preStatus: 200, patchStatus: 204, readStatus: 200},
		{name: "false_true", nonblocking: true, preStatus: 200, patchStatus: 204, readStatus: 200},
		{name: "mixed_explicit_and_omitted", quality: true, omitNonblocking: true, preStatus: 200, patchStatus: 204, readStatus: 200},
		{name: "omitted_preserves_service", omit: true, preStatus: 200, patchStatus: 204, readStatus: 200},
		{name: "create_explicit", create: true, quality: true, patchStatus: 204, readStatus: 200},
		{name: "create_omitted_uses_service_defaults", create: true, omit: true, readStatus: 200},
		{name: "default_create_preserves_service", create: true, defaultProfile: true, omit: true, preStatus: 200, patchStatus: 204, readStatus: 200},
		{name: "snapshot_denied_stops_write", quality: true, preStatus: 403},
		{name: "patch_denied_reads_actual_values", quality: true, preStatus: 200, patchStatus: 403, readStatus: 200},
		{name: "read_denied_retains_snapshot", quality: true, requestPriority: true, preStatus: 200, patchStatus: 204, readStatus: 403},
		{name: "malformed_read_retains_snapshot", quality: true, preStatus: 200, patchStatus: 204, readStatus: 200, malformed: true},
		{name: "mismatch_retains_observed", quality: true, preStatus: 200, patchStatus: 204, readStatus: 200, mismatch: true},
		{name: "failed_create_read_retains_id", create: true, quality: true, patchStatus: 204, readStatus: 403},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			id := "11111111-1111-1111-1111-111111111111_Windows10EnrollmentCompletionPageConfiguration"
			priority := float64(1)
			displayName := "Offline ESP fields"
			if tc.defaultProfile {
				id = "11111111-1111-1111-1111-111111111111_DefaultWindows10EnrollmentCompletionPageConfiguration"
				priority, displayName = 0, kDefConfigDisplayName
			}
			observed := map[string]any{"id": id, "installQualityUpdates": false, "allowNonBlockingAppInstallation": true}
			parentWrites, fieldWrites, snapshots, reads, assignments := 0, 0, 0, 0, 0
			r := priorityTestResource(t, func(req *http.Request) (*http.Response, error) {
				if req.Method == http.MethodGet {
					if req.URL.Query().Get("$filter") != "" {
						return priorityTestResponse(req, 200, map[string]any{"value": []any{map[string]any{"id": id}}})
					}
					status := tc.readStatus
					if parentWrites == 0 {
						snapshots++
						status = tc.preStatus
					} else {
						reads++
					}
					if status != 200 {
						return priorityTestResponse(req, status, map[string]any{"error": map[string]any{"code": "Forbidden", "message": "synthetic failure"}})
					}
					if parentWrites > 0 && tc.malformed {
						return priorityTestResponse(req, 200, map[string]any{"id": id, "installQualityUpdates": "invalid", "allowNonBlockingAppInstallation": false})
					}
					return priorityTestResponse(req, 200, observed)
				}
				if strings.HasSuffix(req.URL.Path, "/setPriority") {
					t.Fatal("priority changed after field failure")
				}
				var body map[string]any
				if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if strings.HasSuffix(req.URL.Path, "/assign") {
					assignments++
					if reads != 1 {
						t.Fatal("assignment ran before field readback")
					}
					return priorityTestResponse(req, 204, nil)
				}
				if _, isolated := body["allowNonBlockingAppInstallation"]; isolated {
					fieldWrites++
					if req.Method != http.MethodPatch || body["@odata.type"] != enrollmentESPType || len(body) != 3 || parentWrites != 1 {
						t.Fatalf("unexpected isolated field write: %s %v", req.Method, body)
					}
					if tc.patchStatus != 204 {
						return priorityTestResponse(req, tc.patchStatus, map[string]any{"error": map[string]any{"code": "Forbidden", "message": "synthetic field rejection"}})
					}
					if !tc.mismatch {
						observed["installQualityUpdates"] = body["installQualityUpdates"]
						observed["allowNonBlockingAppInstallation"] = body["allowNonBlockingAppInstallation"]
					}
					return priorityTestResponse(req, 204, nil)
				}
				parentWrites++
				for _, graphName := range enrollmentESPFields {
					if _, present := body[graphName]; present {
						t.Fatalf("field %s leaked into parent write", graphName)
					}
				}
				if tc.defaultProfile {
					for _, key := range []string{"displayName", "priority", "assignments"} {
						if _, present := body[key]; present {
							t.Fatalf("default parent write contains %s", key)
						}
					}
				}
				// Existing profiles lose this value during broad writes; a new profile keeps its service default.
				if !tc.create || tc.defaultProfile {
					observed["allowNonBlockingAppInstallation"] = false
				}
				if tc.create && !tc.defaultProfile {
					return priorityTestResponse(req, 201, observed)
				}
				return priorityTestResponse(req, 204, nil)
			})
			var diags diag.Diagnostics
			state := tfsdk.State{Schema: r.SpecificSchema}
			state.Raw = generic.ConvertOdataRawToTerraform(ctx, &diags, state.Schema, map[string]any{
				"@odata.type": enrollmentESPType, "id": id, "priority": priority, "displayName": displayName,
				"showInstallationProgress": true, "assignments": []any{}, "installQualityUpdates": false, "allowNonBlockingAppInstallation": true,
			}, "", r.AccessParams.GraphToTerraformMiddleware, id, false, nil)
			if diags.HasError() {
				t.Fatal(diags)
			}
			plan := state
			for name, value := range map[string]bool{"install_quality_updates": tc.quality, "allow_non_blocking_app_installation": tc.nonblocking} {
				if diags := plan.SetAttribute(ctx, path.Root("windows10_esp").AtName(name), types.BoolValue(value)); diags.HasError() {
					t.Fatal(diags)
				}
			}
			if tc.requestPriority {
				priorityTestSet(t, &plan, "priority", types.Int64Value(2))
			}
			config := plan
			if !tc.defaultProfile && !tc.requestPriority {
				priorityTestSet(t, &config, "priority", types.Int64Null())
			}
			for name := range enrollmentESPFields {
				if tc.omit || (tc.omitNonblocking && name == "allow_non_blocking_app_installation") {
					if diags := config.SetAttribute(ctx, path.Root("windows10_esp").AtName(name), types.BoolNull()); diags.HasError() {
						t.Fatal(diags)
					}
				}
			}
			var result tfsdk.State
			if tc.create {
				priorityTestSet(t, &plan, "id", types.StringUnknown())
				response := resource.CreateResponse{State: tfsdk.State{Schema: r.SpecificSchema}}
				r.Create(ctx, resource.CreateRequest{Config: tfsdk.Config{Schema: config.Schema, Raw: config.Raw}, Plan: tfsdk.Plan{Schema: plan.Schema, Raw: plan.Raw}}, &response)
				result, diags = response.State, response.Diagnostics
			} else {
				response := resource.UpdateResponse{State: state}
				r.Update(ctx, resource.UpdateRequest{State: state, Config: tfsdk.Config{Schema: config.Schema, Raw: config.Raw}, Plan: tfsdk.Plan{Schema: plan.Schema, Raw: plan.Raw}}, &response)
				result, diags = response.State, response.Diagnostics
			}
			wantError := tc.preStatus == 403 || tc.patchStatus == 403 || tc.readStatus == 403 || tc.malformed || tc.mismatch
			if diags.HasError() != wantError {
				t.Fatalf("diagnostics: %v", diags)
			}
			var gotID string
			if d := result.GetAttribute(ctx, path.Root("id"), &gotID); d.HasError() || gotID != id {
				t.Fatalf("id = %s; %v", gotID, d)
			}
			if tc.requestPriority {
				var got types.Int64
				if d := result.GetAttribute(ctx, path.Root("priority"), &got); d.HasError() || !got.Equal(types.Int64Value(1)) {
					t.Fatalf("unattempted priority = %v; %v", got, d)
				}
			}
			for name, graphName := range enrollmentESPFields {
				want := types.BoolValue(observed[graphName].(bool))
				if tc.readStatus == 403 || tc.malformed {
					want = types.BoolValue(name == "allow_non_blocking_app_installation")
					if tc.create && !tc.defaultProfile {
						want = types.BoolNull()
					}
				}
				var got types.Bool
				if d := result.GetAttribute(ctx, path.Root("windows10_esp").AtName(name), &got); d.HasError() || !got.Equal(want) {
					t.Fatalf("%s = %v, want %v; %v", name, got, want, d)
				}
			}
			if tc.preStatus == 403 {
				if parentWrites+fieldWrites+reads+assignments != 0 || snapshots != 1 {
					t.Fatal("snapshot failure did not stop all writes")
				}
			} else if parentWrites != 1 || reads != 1 || (!tc.create && snapshots != 1) || (wantError && assignments != 0) {
				t.Fatalf("unexpected requests: parent=%d field=%d snapshots=%d reads=%d assignments=%d", parentWrites, fieldWrites, snapshots, reads, assignments)
			}
			wantFieldWrites := 1
			if tc.preStatus == 403 || (tc.create && !tc.defaultProfile && tc.omit) {
				wantFieldWrites = 0
			}
			if fieldWrites != wantFieldWrites {
				t.Fatalf("field writes = %d, want %d", fieldWrites, wantFieldWrites)
			}
		})
	}
}
