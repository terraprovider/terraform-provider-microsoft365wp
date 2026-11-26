package services

import (
	"context"
	"encoding/json"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

const (
	depTestTokenId   = "437dc82e-9725-45e1-8126-ab8dc7014b1d"
	depTestProfileId = depTestTokenId + "_6de29b1b-9b37-40f4-99b8-9d8a1217ebc6"
)

// depTestRequest is a request that the test server of depTestResource received.
type depTestRequest struct {
	Method string
	Path   string
	Body   map[string]any
}

// depTestResource returns a resource like template that sends its requests to a test server, which answers GETs
// with getBody and records all requests.
func depTestResource(t *testing.T, template *generic.GenericResource, getBody string) (*generic.GenericResource, func() []depTestRequest) {
	t.Helper()
	var mutex sync.Mutex
	var requests []depTestRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		b, _ := io.ReadAll(req.Body)
		var body map[string]any
		if len(b) > 0 {
			if err := json.Unmarshal(b, &body); err != nil {
				t.Errorf("invalid JSON body %s: %v", b, err)
			}
		}
		mutex.Lock()
		requests = append(requests, depTestRequest{Method: req.Method, Path: req.URL.Path, Body: body})
		mutex.Unlock()
		if req.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, getBody)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(server.Close)
	client := msgraph.NewClient(msgraph.VersionBeta)
	client.Endpoint = server.URL
	client.RetryableClient.RetryWaitMin = time.Millisecond
	client.RetryableClient.RetryWaitMax = 2 * time.Millisecond
	r := &generic.GenericResource{
		TypeNameSuffix: template.TypeNameSuffix,
		SpecificSchema: template.SpecificSchema,
		AccessParams: generic.AccessParams{
			BaseUri:                    template.AccessParams.BaseUri,
			ParentEntities:             template.AccessParams.ParentEntities,
			ReadOptions:                template.AccessParams.ReadOptions,
			WriteOptions:               template.AccessParams.WriteOptions,
			TerraformToGraphMiddleware: template.AccessParams.TerraformToGraphMiddleware,
			GraphToTerraformMiddleware: template.AccessParams.GraphToTerraformMiddleware,
			UpdateReplaceFunc:          template.AccessParams.UpdateReplaceFunc,
			DeleteReplaceFunc:          template.AccessParams.DeleteReplaceFunc,
		},
	}
	r.AccessParams.InitializeGuarded(&client)
	return r, func() []depTestRequest {
		mutex.Lock()
		defer mutex.Unlock()
		return append([]depTestRequest{}, requests...)
	}
}

// depTestValue returns a value of s with the given attributes set (keys are paths with dots, values Go values).
func depTestValue(t *testing.T, s schema.Schema, attributes map[string]any) tftypes.Value {
	t.Helper()
	ctx := context.Background()
	state := tfsdk.State{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(ctx), nil)}
	for p, v := range attributes {
		steps := strings.Split(p, ".")
		attrPath := path.Root(steps[0])
		for _, step := range steps[1:] {
			attrPath = attrPath.AtName(step)
		}
		if diags := state.SetAttribute(ctx, attrPath, v); diags.HasError() {
			t.Fatalf("SetAttribute(%s): %v", p, diags)
		}
	}
	return state.Raw
}

// TestDepEnrollmentProfileUpdateSendsFullBody checks that an update that only changes the display name still PATCHes
// all Setup Assistant settings, as MS Graph resets them otherwise (Workplace Foundation issue #2840).
func TestDepEnrollmentProfileUpdateSendsFullBody(t *testing.T) {
	ctx := context.Background()
	s := depEnrollmentProfileResourceSchema
	stateAttributes := map[string]any{
		"dep_onboarding_setting_id":  depTestTokenId,
		"id":                         depTestProfileId,
		"configuration_endpoint_url": "https://example.com/enroll",
		"display_name":               "iOS/iPadOS - Default",
		"description":                "",
		"ios.siri_disabled":          true,
		"ios.apple_id_disabled":      true,
		"ios.enabled_skip_keys":      []string{"Siri", "AppleID", "Intelligence"},
		"ios.itunes_pairing_mode":    "allow",
		"ios.support_department":     "IT",
		"ios.is_default":             true,
	}
	planAttributes := maps.Clone(stateAttributes)
	planAttributes["display_name"] = "iOS/iPadOS - Default - deprecated"
	stateRaw := depTestValue(t, s, stateAttributes)
	planRaw := depTestValue(t, s, planAttributes)

	r, requests := depTestResource(t, &DepEnrollmentProfileResource, `{"@odata.type":"#microsoft.graph.depIOSEnrollmentProfile","id":"`+depTestProfileId+`"}`)
	req := resource.UpdateRequest{
		Config: tfsdk.Config{Schema: s, Raw: planRaw},
		Plan:   tfsdk.Plan{Schema: s, Raw: planRaw},
		State:  tfsdk.State{Schema: s, Raw: stateRaw},
	}
	resp := &resource.UpdateResponse{State: tfsdk.State{Schema: s, Raw: stateRaw}}
	r.Update(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %v", resp.Diagnostics)
	}

	var patch *depTestRequest
	for _, rq := range requests() {
		if rq.Method == http.MethodPatch {
			if patch != nil {
				t.Fatalf("more than one PATCH: %v", requests())
			}
			patch = &rq
		}
	}
	if patch == nil {
		t.Fatalf("no PATCH: %v", requests())
	}
	if want := "/beta/deviceManagement/depOnboardingSettings/" + depTestTokenId + "/enrollmentProfiles/" + depTestProfileId; patch.Path != want {
		t.Errorf("PATCH path %s, want %s", patch.Path, want)
	}
	for k, want := range map[string]any{
		"@odata.type":       "#microsoft.graph.depIOSEnrollmentProfile",
		"displayName":       "iOS/iPadOS - Default - deprecated",
		"siriDisabled":      true,
		"appleIdDisabled":   true,
		"enabledSkipKeys":   []any{"Siri", "AppleID", "Intelligence"},
		"iTunesPairingMode": "allow",
		"supportDepartment": "IT",
	} {
		got, ok := patch.Body[k]
		if !ok {
			t.Errorf("PATCH body lacks %s", k)
			continue
		}
		gotJson, _ := json.Marshal(got)
		wantJson, _ := json.Marshal(want)
		if k == "enabledSkipKeys" {
			// a set, so the order is not defined
			gotKeys, _ := got.([]any)
			if len(gotKeys) != 3 {
				t.Errorf("PATCH body %s = %s, want %s (any order)", k, gotJson, wantJson)
			}
			continue
		}
		if string(gotJson) != string(wantJson) {
			t.Errorf("PATCH body %s = %s, want %s", k, gotJson, wantJson)
		}
	}
	// all other writable attributes are sent as well (as null here, as the test plan has no defaults)
	for _, k := range []string{"description", "requiresUserAuthentication", "privacyPaneDisabled", "zoomDisabled", "managementCertificates"} {
		if _, ok := patch.Body[k]; !ok {
			t.Errorf("PATCH body lacks %s", k)
		}
	}
	for _, k := range []string{"id", "configurationEndpointUrl", "isDefault", "depOnboardingSettingId"} {
		if _, ok := patch.Body[k]; ok {
			t.Errorf("PATCH body contains read-only or parent attribute %s", k)
		}
	}
	for k := range patch.Body {
		if k == "accessibilityScreenDisabled" || k == "adminAccountPassword" || k == "fileVaultDisabled" {
			t.Errorf("PATCH body of an iOS profile contains macOS attribute %s", k)
		}
	}
}

func TestDepEnrollmentProfileImportModifyFunc(t *testing.T) {
	for _, tc := range []struct{ id, want string }{
		{depTestProfileId, depTestTokenId + "/" + depTestProfileId},
		{depTestTokenId + "/" + depTestProfileId, depTestTokenId + "/" + depTestProfileId},
		{"no-token-prefix", "no-token-prefix"}, // the generic import then fails with the expected format
	} {
		params := generic.ImportModifyFuncParams{Id: tc.id}
		diags := diag.Diagnostics{}
		depEnrollmentProfileImportModifyFunc(context.Background(), &diags, &params)
		if diags.HasError() || params.Id != tc.want {
			t.Errorf("%s: got %q (diags %v), want %q", tc.id, params.Id, diags, tc.want)
		}
	}
}

func TestDepDefaultEnrollmentProfileImportModifyFunc(t *testing.T) {
	const ecv2Id = "ECV2_" + depTestTokenId + "_0f9e8d7c-6b5a-4f3e-9d2c-1b0a9f8e7d6c"
	for _, tc := range []struct{ id, want string }{
		{depTestProfileId, depTestTokenId + "/" + depTestProfileId},
		{ecv2Id, depTestTokenId + "/" + ecv2Id},
		{depTestTokenId + "/" + ecv2Id, depTestTokenId + "/" + ecv2Id},
		{"a_b_c_d", "a_b_c_d"},
	} {
		params := generic.ImportModifyFuncParams{Id: tc.id}
		diags := diag.Diagnostics{}
		depDefaultEnrollmentProfileImportModifyFunc(context.Background(), &diags, &params)
		if diags.HasError() || params.Id != tc.want {
			t.Errorf("%s: got %q (diags %v), want %q", tc.id, params.Id, diags, tc.want)
		}
	}
}

func TestDepEnrollmentProfileMoveState(t *testing.T) {
	ctx := context.Background()
	s := depEnrollmentProfileResourceSchema
	for _, tc := range []struct {
		name      string
		json      string
		wantError string
	}{
		{"msgraph DEP profile", `{"id":"` + depTestProfileId + `","url":"deviceManagement/depOnboardingSettings/` + depTestTokenId + `/enrollmentProfiles","api_version":"beta"}`, ""},
		{"msgraph DEP profile with slashes", `{"id":"` + depTestProfileId + `","url":"/devicemanagement/depOnboardingSettings/` + depTestTokenId + `/enrollmentProfiles/"}`, ""},
		{"msgraph other url", `{"id":"x","url":"deviceManagement/androidDeviceOwnerEnrollmentProfiles"}`, `not one with url "deviceManagement/androidDeviceOwnerEnrollmentProfiles"`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := &generic.GenericResource{
				SpecificSchema:   s,
				AccessParams:     generic.AccessParams{ParentEntities: DepEnrollmentProfileResource.AccessParams.ParentEntities},
				MoveStateSources: DepEnrollmentProfileResource.MoveStateSources,
			}
			req := resource.MoveStateRequest{
				SourceProviderAddress: "registry.terraform.io/microsoft/msgraph",
				SourceTypeName:        "msgraph_resource",
				SourceRawState:        &tfprotov6.RawState{JSON: []byte(tc.json)},
			}
			nullState := tftypes.NewValue(s.Type().TerraformType(ctx), nil)
			resp := resource.MoveStateResponse{TargetState: tfsdk.State{Schema: s, Raw: nullState}}
			movers := r.MoveState(ctx)
			if len(movers) != 1 {
				t.Fatalf("got %d state movers, want 1", len(movers))
			}
			movers[0].StateMover(ctx, req, &resp)
			if tc.wantError != "" {
				if !resp.Diagnostics.HasError() || !strings.Contains(resp.Diagnostics.Errors()[0].Detail(), tc.wantError) {
					t.Fatalf("want error containing %q, got %v", tc.wantError, resp.Diagnostics)
				}
				return
			}
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %v", resp.Diagnostics)
			}
			var id, tokenId types.String
			resp.Diagnostics.Append(resp.TargetState.GetAttribute(ctx, path.Root("id"), &id)...)
			resp.Diagnostics.Append(resp.TargetState.GetAttribute(ctx, path.Root("dep_onboarding_setting_id"), &tokenId)...)
			if resp.Diagnostics.HasError() || id.ValueString() != depTestProfileId || tokenId.ValueString() != depTestTokenId {
				t.Errorf("got id %q, dep_onboarding_setting_id %q (diags %v)", id.ValueString(), tokenId.ValueString(), resp.Diagnostics)
			}
		})
	}
}

// TestDepEnrollmentProfileReadDuplicateSkipKeys checks that a read accepts the duplicate entries that MS Graph returns
// in enabledSkipKeys (e.g. after a full PATCH of a profile that has been reset, live test 2026-11-26).
func TestDepEnrollmentProfileReadDuplicateSkipKeys(t *testing.T) {
	ctx := context.Background()
	s := depEnrollmentProfileResourceSchema
	stateRaw := depTestValue(t, s, map[string]any{"dep_onboarding_setting_id": depTestTokenId, "id": depTestProfileId, "display_name": "x"})
	r, _ := depTestResource(t, &DepEnrollmentProfileResource, `{"@odata.type":"#microsoft.graph.depIOSEnrollmentProfile","id":"`+depTestProfileId+
		`","displayName":"x","siriDisabled":true,"enabledSkipKeys":["Siri","Privacy","Siri","Privacy","Intelligence"]}`)
	resp := &resource.ReadResponse{State: tfsdk.State{Schema: s, Raw: stateRaw}}
	r.Read(ctx, resource.ReadRequest{State: tfsdk.State{Schema: s, Raw: stateRaw}}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %v", resp.Diagnostics)
	}
	var skipKeys []string
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("ios").AtName("enabled_skip_keys"), &skipKeys)...)
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}
	slices.Sort(skipKeys)
	if want := []string{"Intelligence", "Privacy", "Siri"}; !slices.Equal(skipKeys, want) {
		t.Errorf("enabled_skip_keys = %v, want %v", skipKeys, want)
	}
}

// TestDepDefaultEnrollmentProfileCreate checks that making a profile the default only calls setDefaultProfile and
// never PATCHes the profile (which would reset its Setup Assistant settings).
func TestDepDefaultEnrollmentProfileCreate(t *testing.T) {
	ctx := context.Background()
	s := depDefaultEnrollmentProfileResourceSchema
	configRaw := depTestValue(t, s, map[string]any{"dep_onboarding_setting_id": depTestTokenId, "id": depTestProfileId})
	planState := tfsdk.State{Schema: s, Raw: configRaw}
	planState.SetAttribute(ctx, path.Root("is_default"), true)
	planState.SetAttribute(ctx, path.Root("display_name"), types.StringUnknown())

	r, requests := depTestResource(t, &DepDefaultEnrollmentProfileResource,
		`{"@odata.type":"#microsoft.graph.depIOSEnrollmentProfile","id":"`+depTestProfileId+`","displayName":"iOS/iPadOS - Default","isDefault":true}`)
	req := resource.CreateRequest{
		Config: tfsdk.Config{Schema: s, Raw: configRaw},
		Plan:   tfsdk.Plan{Schema: s, Raw: planState.Raw},
	}
	resp := &resource.CreateResponse{State: tfsdk.State{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(ctx), nil)}}
	r.Create(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %v", resp.Diagnostics)
	}

	profilePath := "/beta/deviceManagement/depOnboardingSettings/" + depTestTokenId + "/enrollmentProfiles/" + depTestProfileId
	sawAction := false
	for _, rq := range requests() {
		switch {
		case rq.Method == http.MethodPost && rq.Path == profilePath+"/setDefaultProfile":
			sawAction = true
		case rq.Method == http.MethodGet && rq.Path == profilePath:
		default:
			t.Errorf("unexpected request %s %s", rq.Method, rq.Path)
		}
	}
	if !sawAction {
		t.Errorf("no setDefaultProfile request: %v", requests())
	}
	var displayName types.String
	var isDefault types.Bool
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("display_name"), &displayName)...)
	resp.Diagnostics.Append(resp.State.GetAttribute(ctx, path.Root("is_default"), &isDefault)...)
	if resp.Diagnostics.HasError() || displayName.ValueString() != "iOS/iPadOS - Default" || !isDefault.ValueBool() {
		t.Errorf("state display_name %v, is_default %v (diags %v)", displayName, isDefault, resp.Diagnostics)
	}
}
