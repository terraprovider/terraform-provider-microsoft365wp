package generic

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &GenericResource{}
	_ resource.ResourceWithConfigure   = &GenericResource{}
	_ resource.ResourceWithImportState = &GenericResource{}
)

// Resource implementation.
type GenericResource struct {
	TypeNameSuffix           string
	SpecificSchema           schema.Schema
	SpecificConfigValidators []resource.ConfigValidator
	AccessParams             AccessParams

	serializationMutex sync.Mutex
}

// Returns the resource type name.
func (r *GenericResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, r.TypeNameSuffix)
}

// Adds the provider configured MSGraph client to the resource.
func (r *GenericResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	r.AccessParams.InitializeGuarded(req.ProviderData)
}

// Defines the schema for the resource.
func (r *GenericResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = r.SpecificSchema
}

func (r *GenericResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return r.SpecificConfigValidators
}

func (r *GenericResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.createUpdate(ctx, OperationCreate, &req, resp, nil, nil)
}

func (r *GenericResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	tfVal := r.AccessParams.ReadSingleRaw(ctx, &resp.Diagnostics, req.State.Schema,
		r.AccessParams.ReadOptions, "", &req.State, true)
	if resp.Diagnostics.HasError() {
		return
	}

	if !tfVal.IsNull() {
		newState := tfsdk.State{
			Schema: req.State.Schema,
			Raw:    tfVal,
		}
		if !r.AccessParams.HasParentItem.ParentIdField.Equal(path.Path{}) {
			err := CopyValueAtPath(ctx, &newState, &req.State, r.AccessParams.HasParentItem.ParentIdField)
			if err != nil {
				resp.Diagnostics.AddError(fmt.Sprintf("CopyValueAtPath(), ParentIdField: %s", r.AccessParams.HasParentItem.ParentIdField), err.Error())
				return
			}
		}
		resp.State = newState
	} else {
		// item not found (anymore), remove from state
		resp.State.RemoveResource(ctx)
	}
}

func (r *GenericResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.createUpdate(ctx, OperationUpdate, nil, nil, &req, resp)
}

func (r *GenericResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	if r.checkAndLockMutex() {
		defer r.serializationMutex.Unlock()
	}

	thisIdAttributer := req.State
	diags := &resp.Diagnostics

	baseUri := ""
	id := ""
	skipDelete := false

	if r.AccessParams.IsSingleton {
		skipDelete = true
	}

	if r.AccessParams.DeleteModifyFunc != nil {
		id = r.AccessParams.GetId(ctx, diags, id, thisIdAttributer)
		if diags.HasError() {
			return
		}
		params := DeleteModifyFuncParams{
			R:          r,
			Ctx:        ctx,
			Diags:      diags,
			Req:        req,
			Resp:       resp,
			BaseUri:    baseUri,
			Id:         id,
			SkipDelete: false,
		}
		r.AccessParams.DeleteModifyFunc(&params)
		if diags.HasError() {
			return
		}
		baseUri = params.BaseUri
		id = params.Id
		skipDelete = params.SkipDelete
	}

	if skipDelete {
		diags.AddWarning("Skipping deletion of entity from MS Graph, just removing resource from state", "Cannot delete entities that are singletons and/or have been created by MS Graph itself.")
		return
	}

	subActionsData := r.executeWriteSubActionsPre(ctx, diags, OperationDelete, nil, id, thisIdAttributer)
	if diags.HasError() {
		return
	}

	r.AccessParams.DeleteRaw(ctx, diags, baseUri, id, thisIdAttributer)
	if diags.HasError() {
		return
	}

	r.executeWriteSubActionPost(ctx, diags, OperationDelete, id, thisIdAttributer, subActionsData, &req.State, nil)
	if diags.HasError() {
		return
	}
}

func (r *GenericResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *GenericResource) createUpdate(ctx context.Context, operationType OperationType,
	createRequest *resource.CreateRequest, createResponse *resource.CreateResponse,
	updateRequest *resource.UpdateRequest, updateResponse *resource.UpdateResponse) {

	if r.checkAndLockMutex() {
		defer r.serializationMutex.Unlock()
	}

	// do this in here since we anyway need the original request and response for CreateModifyFuncParams and UpdateModifyFuncParams
	var diags *diag.Diagnostics
	var requestPlan *tfsdk.Plan
	var requestState *tfsdk.State
	var responseState *tfsdk.State
	var parentIdAttributer GetAttributer = nil
	var thisIdAttributer GetAttributer = nil
	switch operationType {
	case OperationCreate:
		diags = &createResponse.Diagnostics
		requestPlan = &createRequest.Plan
		parentIdAttributer = createRequest.Config
		responseState = &createResponse.State
	case OperationUpdate:
		diags = &updateResponse.Diagnostics
		requestPlan = &updateRequest.Plan
		requestState = &updateRequest.State
		thisIdAttributer = updateRequest.State
		responseState = &updateResponse.State
	default:
		panic(fmt.Sprintf("Invalid operation type %d", operationType))
	}

	baseUri := ""
	id := ""
	updateExisting := operationType == OperationUpdate

	if operationType == OperationCreate && r.AccessParams.IsSingleton {
		updateExisting = true
		id = "singleton" // will only be used for TF resource but not for MS Graph for IsSingleton
	}

	switch operationType {

	case OperationCreate:
		if r.AccessParams.CreateModifyFunc != nil {
			params := CreateModifyFuncParams{
				R:              r,
				Ctx:            ctx,
				Diags:          diags,
				Req:            *createRequest,
				Resp:           createResponse,
				BaseUri:        baseUri,
				Id:             id,
				UpdateExisting: updateExisting,
			}
			r.AccessParams.CreateModifyFunc(&params)
			if diags.HasError() {
				return
			}
			baseUri = params.BaseUri
			id = params.Id
			updateExisting = params.UpdateExisting
		}

	case OperationUpdate:
		if r.AccessParams.UpdateModifyFunc != nil {
			id = r.AccessParams.GetId(ctx, diags, id, thisIdAttributer)
			if diags.HasError() {
				return
			}
			params := UpdateModifyFuncParams{
				R:       r,
				Ctx:     ctx,
				Diags:   diags,
				Req:     *updateRequest,
				Resp:    updateResponse,
				BaseUri: baseUri,
				Id:      id,
			}
			r.AccessParams.UpdateModifyFunc(&params)
			if diags.HasError() {
				return
			}
			baseUri = params.BaseUri
			id = params.Id
		}

	}

	if operationType == OperationCreate && updateExisting {
		diags.AddWarning("Using existing entity from MS Graph instead of creating new one", "The targeted entity is a singleton and/or has been created by MS Graph itself.")
	}

	includeNullObjects := updateExisting && !r.AccessParams.UsePutForUpdate
	rawVal := ConvertTerraformToOdataRaw(ctx, diags, requestPlan.Schema, includeNullObjects,
		requestPlan.Raw, updateExisting, r.AccessParams.TerraformToGraphMiddleware, "Plan")
	if diags.HasError() {
		return
	}

	subActionsData := r.executeWriteSubActionsPre(ctx, diags, operationType, rawVal, "", thisIdAttributer)
	if diags.HasError() {
		return
	}

	var createResultRaw map[string]any
	if updateExisting {
		r.AccessParams.UpdateRaw(ctx, diags, baseUri, id, thisIdAttributer, rawVal, r.AccessParams.UsePutForUpdate)
	} else {
		id, createResultRaw = r.AccessParams.CreateRaw(ctx, diags, baseUri, parentIdAttributer, rawVal)
	}
	if diags.HasError() {
		return
	}

	// TF will persist any state found in the response struct to real TF state even in case of terminating errors.
	// Since this can lead to nasty problems we try to abide by the follwing priorities for setting (new) state in the
	// response struct:
	// 1. We _never ever_ want to set any response state without a valid `id`: This will severely corrupt TF state (i.e.
	// the resource starts to exist in TF state but its `id` is empty and therefore the TF resource (and potentially the
	// created corresponding MS Graph entity) can only be removed manually from TF state (and MS Graph)).
	// 2. We never want to create any entity in MS Graph without informing TF about it (i.e. without setting response
	// state): This will lead to duplicate entities in MS Graph (since TF can only destroy and recreate tainted
	// resources that it got to know about).
	// 3. We would like to only set reponse state after creation or update has been successful: But this is not too
	// important as TF will detect missing entities in MS Graph or differences on next the state refresh anyway and
	// usually cleanup up automatically eventually.

	// For the creation or update to also succeed for TF, TF will only accept the exact same plan as a response state
	// (with the exception of previously unknown values that _must_ have been filled with the now known actual values)
	// and otherwise will throw a nasty error to the users (heavily accusing the plugin developers for being responsible
	// for it).
	// Therefore we only copy the plan to the response state (and later try to fill in unknown values) and also rely
	// here on any differences to be detected on next state refresh and to get sorted out eventually.

	switch operationType {

	case OperationCreate:
		if id == "" {
			diags.AddError("Error saving id to state",
				fmt.Sprintf("GenericResource.createUpdate: operationType is OperationCreate (updateExisting: %t), but id is (still) empty", updateExisting))
			return
		}
		// Create a copy to not modify request plan
		responseStateTemp := tfsdk.State{
			Raw:    requestPlan.Raw,
			Schema: responseState.Schema,
		}
		diags.Append(responseStateTemp.SetAttribute(ctx, path.Root("id"), id)...)
		if diags.HasError() {
			return
		}
		// Save plan with new id added
		responseState.Raw = responseStateTemp.Raw

	case OperationUpdate:
		// Just copy plan
		responseState.Raw = requestPlan.Raw

	}

	// Execute WriteSubActions first to ensure that populating unknown values can pick up potential changes.
	// But do not check for errors here yet as we want to ensure that populating runs even in case of errors here to get
	// the state as complete as possible. Terraform will still bring up the errors as they will still be there after
	// populating has run.
	r.executeWriteSubActionPost(ctx, diags, operationType, id, thisIdAttributer, subActionsData, requestState, requestPlan)

	// Produce a wholly-known new state by determining the final values for any attributes left unknown in the planned state.
	populateSource := tftypes.Value{}
	if createResultRaw != nil {
		populateSource = ConvertOdataRawToTerraform(ctx, diags, requestPlan.Schema, createResultRaw, "", r.AccessParams.GraphToTerraformMiddleware)
		if diags.HasError() {
			return
		}
	}
	diags.Append(r.populateUnknownValues(ctx, responseState, populateSource)...)
	if diags.HasError() { // still check here since it looks nicer and cleaner
		return
	}
}

func (r *GenericResource) checkAndLockMutex() bool {
	if r.AccessParams.SerializeWrites {
		if !r.serializationMutex.TryLock() {
			r.serializationMutex.Lock()
			time.Sleep(r.AccessParams.SerialWritesDelay)
		}
	}
	return r.AccessParams.SerializeWrites
}

func (r *GenericResource) executeWriteSubActionsPre(ctx context.Context, diags *diag.Diagnostics, wsaOperation OperationType,
	rawVal map[string]any, id string, idAttributer GetAttributer) map[string]any {

	result := make(map[string]any)
	for _, a := range r.AccessParams.WriteSubActions {
		if a.CheckRunAction(wsaOperation) {
			a.ExecutePre(ctx, diags, r, wsaOperation, rawVal, result, id, idAttributer)
			if diags.HasError() {
				return nil
			}
		}
	}

	return result
}

func (r *GenericResource) executeWriteSubActionPost(ctx context.Context, diags *diag.Diagnostics, wsaOperation OperationType,
	id string, idAttributer GetAttributer, subActionsData map[string]any, reqState *tfsdk.State, reqPlan *tfsdk.Plan) {

	for _, a := range r.AccessParams.WriteSubActions {
		if a.CheckRunAction(wsaOperation) {
			a.ExecutePost(ctx, diags, r, wsaOperation, id, idAttributer, subActionsData, reqState, reqPlan)
			if diags.HasError() {
				return
			}
		}
	}
}
