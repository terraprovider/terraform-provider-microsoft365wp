package generic

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithMoveState = &GenericResource{}

// MoveStateSource is a resource type (usually of another provider) that manages the same MS Graph entities as this
// resource, so that its state can be moved to this resource with a `moved` block (instead of an `import` block, which
// needs the id in the configuration). Only the id gets taken from the source state, the read that follows the move
// fills in everything else from MS Graph.
type MoveStateSource struct {
	// ProviderType is the type of the source provider (the last part of its address, e.g. "azuread"), so that any
	// registry and namespace (e.g. forks of a provider) are accepted. Empty accepts any provider.
	ProviderType string
	// TypeName is the source resource type, e.g. "azuread_named_location".
	TypeName string
	// IdFunc returns the MS Graph id of the entity from the source state (decoded from JSON). It returns an error if
	// the source state does not describe an entity of this resource (e.g. a generic resource with another URL).
	IdFunc func(rawState map[string]any) (string, error)
	// ParentIdsFunc returns the ids of the parent entities from the source state, in the order of
	// AccessParams.ParentEntities (e.g. taken from the URL of a generic resource). It is required for resources with
	// parent entities and must not be set for others.
	ParentIdsFunc func(rawState map[string]any) ([]string, error)
}

func (s *MoveStateSource) matches(req *resource.MoveStateRequest) bool {
	if req.SourceTypeName != s.TypeName {
		return false
	}
	if s.ProviderType == "" {
		return true
	}
	addressParts := strings.Split(req.SourceProviderAddress, "/")
	return strings.EqualFold(addressParts[len(addressParts)-1], s.ProviderType)
}

func (r *GenericResource) MoveState(_ context.Context) []resource.StateMover {
	stateMovers := make([]resource.StateMover, 0, len(r.MoveStateSources))
	for _, source := range r.MoveStateSources {
		stateMovers = append(stateMovers, resource.StateMover{
			StateMover: func(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
				r.moveState(ctx, &source, &req, resp)
			},
		})
	}
	return stateMovers
}

// moveState sets just the id (and the ids of the parent entities) in the target state if the request matches source
// (and leaves the target state null otherwise, so that the framework tries the next StateMover).
func (r *GenericResource) moveState(ctx context.Context, source *MoveStateSource, req *resource.MoveStateRequest, resp *resource.MoveStateResponse) {

	if !source.matches(req) {
		return
	}

	const errorSummary = "Unable to Move Resource State"
	errorDetail := func(format string, a ...any) string {
		return fmt.Sprintf("Cannot move the state of %s (provider %s) to this resource: %s", req.SourceTypeName,
			req.SourceProviderAddress, fmt.Sprintf(format, a...))
	}

	// MoveState may be called before Configure
	r.AccessParams.InitializeGuarded(nil)
	parentEntities := r.AccessParams.ParentEntities
	if len(parentEntities) > 0 && source.ParentIdsFunc == nil {
		resp.Diagnostics.AddError(errorSummary, errorDetail("moving the state of resources with parent entities is not supported for this source"))
		return
	}

	if req.SourceRawState == nil || len(req.SourceRawState.JSON) == 0 {
		resp.Diagnostics.AddError(errorSummary, errorDetail("the source state has no JSON data"))
		return
	}
	var rawState map[string]any
	if err := json.Unmarshal(req.SourceRawState.JSON, &rawState); err != nil {
		resp.Diagnostics.AddError(errorSummary, errorDetail("cannot decode the source state: %s", err))
		return
	}

	id, err := source.IdFunc(rawState)
	if err == nil && id == "" {
		err = fmt.Errorf("the source state has no id")
	}
	if err != nil {
		resp.Diagnostics.AddError(errorSummary, errorDetail("%s", err))
		return
	}

	var parentIds []string
	if len(parentEntities) > 0 {
		parentIds, err = source.ParentIdsFunc(rawState)
		if err == nil && len(parentIds) != len(parentEntities) {
			err = fmt.Errorf("got %d parent id(s) from the source state, but %d are required", len(parentIds), len(parentEntities))
		}
		if err == nil && slices.Contains(parentIds, "") {
			err = fmt.Errorf("the source state has an empty parent id")
		}
		if err != nil {
			resp.Diagnostics.AddError(errorSummary, errorDetail("%s", err))
			return
		}
	}

	resp.Diagnostics.Append(resp.TargetState.SetAttribute(ctx, path.Root(r.AccessParams.EntityId.AttrNameTf), id)...)
	for i, pe := range parentEntities {
		resp.Diagnostics.Append(resp.TargetState.SetAttribute(ctx, pe.ParentIdField, parentIds[i])...)
	}
}
