package generic

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var moveStateTestSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id":           schema.StringAttribute{Computed: true},
		"display_name": schema.StringAttribute{Optional: true},
	},
}

func moveStateTestIdFunc(rawState map[string]any) (string, error) {
	if rawState["kind"] != nil && rawState["kind"] != "thing" {
		return "", fmt.Errorf("kind %v is not a thing", rawState["kind"])
	}
	id, _ := rawState["id"].(string)
	return id, nil
}

// runStateMovers calls the StateMovers of r like the framework does (until one sets the target state or returns an
// error) and returns the last response and whether the state has been moved.
func runStateMovers(t *testing.T, r *GenericResource, s schema.Schema, req resource.MoveStateRequest) (resource.MoveStateResponse, bool) {
	t.Helper()
	ctx := context.Background()
	nullState := tftypes.NewValue(s.Type().TerraformType(ctx), nil)
	var resp resource.MoveStateResponse
	for _, sm := range r.MoveState(ctx) {
		resp = resource.MoveStateResponse{TargetState: tfsdk.State{Schema: s, Raw: nullState}}
		sm.StateMover(ctx, req, &resp)
		if resp.Diagnostics.HasError() {
			return resp, false
		}
		if !resp.TargetState.Raw.Equal(nullState) {
			return resp, true
		}
	}
	return resp, false
}

func TestMoveState(t *testing.T) {
	sources := []MoveStateSource{
		{ProviderType: "other", TypeName: "other_thing", IdFunc: moveStateTestIdFunc},
		{TypeName: "generic_resource", IdFunc: moveStateTestIdFunc},
	}
	tests := []struct {
		name            string
		providerAddress string
		typeName        string
		json            string // "-" for no raw state
		parentEntities  ParentEntities
		wantId          string // "" if the state must not be moved
		wantError       string
	}{
		{name: "match", providerAddress: "registry.opentofu.org/somebody/other", typeName: "other_thing", json: `{"id":"a1","display_name":"A"}`, wantId: "a1"},
		{name: "provider type is case-insensitive", providerAddress: "registry.terraform.io/x/Other", typeName: "other_thing", json: `{"id":"a2"}`, wantId: "a2"},
		{name: "empty ProviderType accepts any provider", providerAddress: "example.com/y/whatever", typeName: "generic_resource", json: `{"id":"b1","kind":"thing"}`, wantId: "b1"},
		{name: "other provider type", providerAddress: "registry.terraform.io/x/another", typeName: "other_thing", json: `{"id":"a1"}`},
		// OpenTofu passes this provider's own address for the further instances of a moved resource with several instances
		{name: "own provider address", providerAddress: "registry.opentofu.org/terraprovider/microsoft365wp", typeName: "other_thing", json: `{"id":"a3"}`, wantId: "a3"},
		{name: "own provider address, other resource type", providerAddress: "registry.opentofu.org/terraprovider/microsoft365wp", typeName: "other_stuff", json: `{"id":"a3"}`},
		{name: "other resource type", providerAddress: "registry.terraform.io/x/other", typeName: "other_stuff", json: `{"id":"a1"}`},
		{name: "IdFunc error", providerAddress: "registry.terraform.io/x/y", typeName: "generic_resource", json: `{"id":"b1","kind":"car"}`, wantError: "kind car is not a thing"},
		{name: "no id", providerAddress: "registry.terraform.io/x/other", typeName: "other_thing", json: `{"display_name":"A"}`, wantError: "the source state has no id"},
		{name: "invalid JSON", providerAddress: "registry.terraform.io/x/other", typeName: "other_thing", json: `{"id":`, wantError: "cannot decode the source state"},
		{name: "no raw state", providerAddress: "registry.terraform.io/x/other", typeName: "other_thing", json: "-", wantError: "no JSON data"},
		{name: "parent entities", providerAddress: "registry.terraform.io/x/other", typeName: "other_thing", json: `{"id":"a1"}`,
			parentEntities: ParentEntities{{ParentIdField: path.Root("parent_id")}}, wantError: "parent entities is not supported"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &GenericResource{
				SpecificSchema:   moveStateTestSchema,
				AccessParams:     AccessParams{ParentEntities: tt.parentEntities},
				MoveStateSources: sources,
			}
			req := resource.MoveStateRequest{SourceProviderAddress: tt.providerAddress, SourceTypeName: tt.typeName}
			if tt.json != "-" {
				req.SourceRawState = &tfprotov6.RawState{JSON: []byte(tt.json)}
			}
			resp, moved := runStateMovers(t, r, moveStateTestSchema, req)

			if tt.wantError != "" {
				if !resp.Diagnostics.HasError() || !strings.Contains(resp.Diagnostics.Errors()[0].Detail(), tt.wantError) {
					t.Fatalf("want error containing %q, got %v", tt.wantError, resp.Diagnostics)
				}
				return
			}
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %v", resp.Diagnostics)
			}
			if moved != (tt.wantId != "") {
				t.Fatalf("moved = %v, want %v", moved, tt.wantId != "")
			}
			if !moved {
				return
			}
			var id, displayName types.String
			resp.Diagnostics.Append(resp.TargetState.GetAttribute(context.Background(), path.Root("id"), &id)...)
			resp.Diagnostics.Append(resp.TargetState.GetAttribute(context.Background(), path.Root("display_name"), &displayName)...)
			if resp.Diagnostics.HasError() {
				t.Fatal(resp.Diagnostics)
			}
			if id.ValueString() != tt.wantId {
				t.Errorf("id = %q, want %q", id.ValueString(), tt.wantId)
			}
			// everything but the id gets read from MS Graph
			if !displayName.IsNull() {
				t.Errorf("display_name = %v, want null", displayName)
			}
		})
	}
}

var moveStateTestSchemaWithParent = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"parent_id":    schema.StringAttribute{Required: true},
		"id":           schema.StringAttribute{Computed: true},
		"display_name": schema.StringAttribute{Optional: true},
	},
}

func moveStateTestParentIdsFunc(rawState map[string]any) ([]string, error) {
	url, _ := rawState["url"].(string)
	parentId, ok := strings.CutPrefix(url, "parents/")
	if !ok {
		return nil, fmt.Errorf("unexpected url %q", url)
	}
	return strings.Split(parentId, ","), nil
}

func TestMoveStateParentEntities(t *testing.T) {
	sources := []MoveStateSource{
		{TypeName: "generic_resource", IdFunc: moveStateTestIdFunc, ParentIdsFunc: moveStateTestParentIdsFunc},
	}
	tests := []struct {
		name         string
		json         string
		wantId       string
		wantParentId string
		wantError    string
	}{
		{name: "match", json: `{"id":"c1","url":"parents/p1"}`, wantId: "c1", wantParentId: "p1"},
		{name: "ParentIdsFunc error", json: `{"id":"c1","url":"others/p1"}`, wantError: `unexpected url "others/p1"`},
		{name: "too many parent ids", json: `{"id":"c1","url":"parents/p1,p2"}`, wantError: "got 2 parent id(s) from the source state, but 1 are required"},
		{name: "empty parent id", json: `{"id":"c1","url":"parents/"}`, wantError: "empty parent id"},
		{name: "IdFunc error comes first", json: `{"id":"c1","kind":"car","url":"others/p1"}`, wantError: "kind car is not a thing"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &GenericResource{
				SpecificSchema:   moveStateTestSchemaWithParent,
				AccessParams:     AccessParams{ParentEntities: ParentEntities{{ParentIdField: path.Root("parent_id")}}},
				MoveStateSources: sources,
			}
			req := resource.MoveStateRequest{SourceProviderAddress: "registry.terraform.io/x/y", SourceTypeName: "generic_resource",
				SourceRawState: &tfprotov6.RawState{JSON: []byte(tt.json)}}
			resp, moved := runStateMovers(t, r, moveStateTestSchemaWithParent, req)

			if tt.wantError != "" {
				if !resp.Diagnostics.HasError() || !strings.Contains(resp.Diagnostics.Errors()[0].Detail(), tt.wantError) {
					t.Fatalf("want error containing %q, got %v", tt.wantError, resp.Diagnostics)
				}
				return
			}
			if resp.Diagnostics.HasError() || !moved {
				t.Fatalf("moved = %v, diagnostics: %v", moved, resp.Diagnostics)
			}
			var id, parentId, displayName types.String
			resp.Diagnostics.Append(resp.TargetState.GetAttribute(context.Background(), path.Root("id"), &id)...)
			resp.Diagnostics.Append(resp.TargetState.GetAttribute(context.Background(), path.Root("parent_id"), &parentId)...)
			resp.Diagnostics.Append(resp.TargetState.GetAttribute(context.Background(), path.Root("display_name"), &displayName)...)
			if resp.Diagnostics.HasError() {
				t.Fatal(resp.Diagnostics)
			}
			if id.ValueString() != tt.wantId || parentId.ValueString() != tt.wantParentId {
				t.Errorf("id, parent_id = %q, %q, want %q, %q", id.ValueString(), parentId.ValueString(), tt.wantId, tt.wantParentId)
			}
			if !displayName.IsNull() {
				t.Errorf("display_name = %v, want null", displayName)
			}
		})
	}
}
