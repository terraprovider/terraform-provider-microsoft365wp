package generic

import (
	"context"
	"fmt"
	"strings"
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
	_ resource.Resource                     = &GenericResource{}
	_ resource.ResourceWithConfigure        = &GenericResource{}
	_ resource.ResourceWithConfigValidators = &GenericResource{}
	_ resource.ResourceWithImportState      = &GenericResource{}
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

	tfVal := r.AccessParams.ReadSingleCompleteTf2(ctx, &resp.Diagnostics, req.State.Schema,
		r.AccessParams.ReadOptions, "", &req.State, "", true, &req.State, resp.Private)
	if resp.Diagnostics.HasError() {
		return
	}

	if !tfVal.IsNull() {
		newState := tfsdk.State{
			Schema: req.State.Schema,
			Raw:    tfVal,
		}
		r.AccessParams.PopulateStateParentIdsFromRequest(ctx, &resp.Diagnostics, &newState, req.State)
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
	skipDelete := r.AccessParams.WriteOptions.SkipDelete

	if r.AccessParams.DeleteModifyFunc != nil {
		id = r.AccessParams.GetId(ctx, diags, id, thisIdAttributer)
		if diags.HasError() {
			return
		}
		params := DeleteModifyFuncParams{
			R:          r,
			Req:        req,
			Resp:       resp,
			BaseUri:    baseUri,
			Id:         id,
			SkipDelete: false,
		}
		r.AccessParams.DeleteModifyFunc(ctx, diags, &params)
		if diags.HasError() {
			return
		}
		baseUri = params.BaseUri
		id = params.Id
		skipDelete = params.SkipDelete
	}

	subActionsData := r.executeWriteSubActionsPre(ctx, diags, OperationDelete, id, thisIdAttributer, nil, &req.State, nil, resp.Private)
	if diags.HasError() {
		return
	}

	if skipDelete {
		diags.AddWarning("Skipping deletion of entity from MS Graph, just removing resource from state", "Cannot delete entities that are singletons and/or have been created by MS Graph itself.")
	} else {
		if r.AccessParams.DeleteReplaceFunc == nil {
			r.AccessParams.DeleteRaw(ctx, diags, baseUri, id, thisIdAttributer)
		} else {
			params := DeleteReplaceFuncParams{
				R:            r,
				Client:       r.AccessParams.graphClient,
				BaseUri:      baseUri,
				Id:           id,
				IdAttributer: thisIdAttributer,
			}
			r.AccessParams.DeleteReplaceFunc(ctx, diags, &params)
		}
		if diags.HasError() {
			return
		}
	}

	r.executeWriteSubActionsPost(ctx, diags, OperationDelete, id, thisIdAttributer, subActionsData, &req.State, nil, &resp.State, resp.Private)
	if diags.HasError() {
		return
	}
}

func (r *GenericResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	// do not use resource.ImportStatePassthroughID here to be able to import entities with parents

	importPaths := []path.Path{}
	importAttributeNames := []string{}
	for _, v := range r.AccessParams.ParentEntities {
		importPaths = append(importPaths, v.ParentIdField)
		importAttributeNames = append(importAttributeNames, v.ParentIdField.String())
	}
	importPaths = append(importPaths, path.Root(r.AccessParams.EntityId.AttrNameTf))
	importAttributeNames = append(importAttributeNames, r.AccessParams.EntityId.AttrNameTf)

	idComponents := strings.Split(req.ID, "/")
	if len(idComponents) != len(importPaths) {
		resp.Diagnostics.AddError("Id for import does not have the correct format",
			fmt.Sprintf("To import this resource, an id with the following component(s) must be specified: '%s' (separated by forward slashes in case of multiple components)", strings.Join(importAttributeNames, "/")))
		return
	}

	for i, v := range idComponents {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, importPaths[i], strings.Trim(v, " "))...)
	}
}

func (r *GenericResource) createUpdate(ctx context.Context, operationType OperationType,
	createRequest *resource.CreateRequest, createResponse *resource.CreateResponse,
	updateRequest *resource.UpdateRequest, updateResponse *resource.UpdateResponse) {

	if r.checkAndLockMutex() {
		defer r.serializationMutex.Unlock()
	}

	// do this in here since we anyway need the original request and response for CreateModifyFuncParams and UpdateModifyFuncParams
	var diags *diag.Diagnostics
	var requestConfig *tfsdk.Config
	var requestPlan *tfsdk.Plan
	var requestState *tfsdk.State
	var responseState *tfsdk.State
	var responsePrivate PrivateDataGetSetter
	var parentIdAttributer GetAttributer = nil
	var thisIdAttributer GetAttributer = nil
	switch operationType {
	case OperationCreate:
		diags = &createResponse.Diagnostics
		requestConfig = &createRequest.Config
		requestPlan = &createRequest.Plan
		parentIdAttributer = createRequest.Config
		responseState = &createResponse.State
		responsePrivate = createResponse.Private
	case OperationUpdate:
		diags = &updateResponse.Diagnostics
		requestConfig = &updateRequest.Config
		requestPlan = &updateRequest.Plan
		requestState = &updateRequest.State
		thisIdAttributer = updateRequest.State
		responseState = &updateResponse.State
		responsePrivate = updateResponse.Private
	default:
		panic(fmt.Sprintf("Invalid operation type %d", operationType))
	}

	baseUri := ""
	id := ""
	updateExisting := operationType == OperationUpdate

	if operationType == OperationCreate {
		if r.AccessParams.WriteOptions.UpdateInsteadOfCreate {
			updateExisting = true
		} else if r.AccessParams.WriteOptions.UpdateIfExistsOnCreate {
			// id must exist in config or be provided by FallbackIdGetterFunc
			id = r.AccessParams.GetId(ctx, diags, id, createRequest.Config)
			if diags.HasError() {
				return
			}
			if r.AccessParams.CheckEntityExistence(ctx, diags, id, createRequest.Config) {
				updateExisting = true
			}
			if diags.HasError() {
				return
			}
		}
		if updateExisting {
			thisIdAttributer = createRequest.Config
		}
	}

	if updateExisting && id == "" {
		if operationType == OperationCreate && r.AccessParams.EntityId.CreateUpdateReadIdFromGraph {
			id = r.AccessParams.ReadId(ctx, diags, baseUri, parentIdAttributer, r.AccessParams.EntityId.CreateUpdateReadIdFromGraphUseODataSelect, "", nil)
		} else {
			id = r.AccessParams.GetId(ctx, diags, id, thisIdAttributer)
		}
		if diags.HasError() {
			return
		}
	}

	switch operationType {

	case OperationCreate:
		if r.AccessParams.CreateModifyFunc != nil {
			params := CreateModifyFuncParams{
				R:              r,
				Req:            *createRequest,
				Resp:           createResponse,
				BaseUri:        baseUri,
				Id:             id,
				UpdateExisting: updateExisting,
			}
			r.AccessParams.CreateModifyFunc(ctx, diags, &params)
			if diags.HasError() {
				return
			}
			baseUri = params.BaseUri
			id = params.Id
			updateExisting = params.UpdateExisting
		}

	case OperationUpdate:
		if r.AccessParams.UpdateModifyFunc != nil {
			params := UpdateModifyFuncParams{
				R:       r,
				Req:     *updateRequest,
				Resp:    updateResponse,
				BaseUri: baseUri,
				Id:      id,
			}
			r.AccessParams.UpdateModifyFunc(ctx, diags, &params)
			if diags.HasError() {
				return
			}
			baseUri = params.BaseUri
			id = params.Id
		}

	}

	if operationType == OperationCreate && updateExisting {
		diags.AddWarning("Using existing entity from MS Graph instead of creating new one",
			"The targeted entity is a singleton, has been created by MS Graph itself, is a descendant of another already existing entity and/or already exists itself and can only be modified.")
	}

	includeNullObjects := updateExisting && !r.AccessParams.WriteOptions.UsePutForUpdate
	rawVal := ConvertTerraformToOdataRaw(ctx, diags, requestPlan.Schema, includeNullObjects,
		requestPlan.Raw, updateExisting, *requestConfig, r.AccessParams.TerraformToGraphMiddleware, "Plan")
	if diags.HasError() {
		return
	}

	subActionsData := r.executeWriteSubActionsPre(ctx, diags, operationType, id, thisIdAttributer, rawVal, requestState, requestPlan, responsePrivate)
	if diags.HasError() {
		return
	}

	var createUpdateResultRaw map[string]any
	if updateExisting {
		if r.AccessParams.UpdateReplaceFunc == nil {
			r.AccessParams.UpdateRaw(ctx, diags, baseUri, id, thisIdAttributer, rawVal)
		} else {
			params := UpdateReplaceFuncParams{
				R:            r,
				Client:       r.AccessParams.graphClient,
				BaseUri:      baseUri,
				Id:           id,
				IdAttributer: thisIdAttributer,
				RawVal:       rawVal,
			}
			r.AccessParams.UpdateReplaceFunc(ctx, diags, &params)
			if params.RawResult != nil {
				createUpdateResultRaw = params.RawResult
			}
		}
	} else {
		if r.AccessParams.CreateReplaceFunc == nil {
			id, createUpdateResultRaw = r.AccessParams.CreateRaw(ctx, diags, baseUri, parentIdAttributer, rawVal)
		} else {
			params := CreateReplaceFuncParams{
				R:            r,
				Client:       r.AccessParams.graphClient,
				BaseUri:      baseUri,
				IdAttributer: parentIdAttributer,
				RawVal:       rawVal,
			}
			r.AccessParams.CreateReplaceFunc(ctx, diags, &params)
			id = params.Id
			createUpdateResultRaw = params.RawResult
		}
		if diags.HasError() {
			return
		}
		if id == "" {
			id = r.AccessParams.GetId(ctx, diags, id, parentIdAttributer)
		}
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
		diags.Append(responseStateTemp.SetAttribute(ctx, path.Root(r.AccessParams.EntityId.AttrNameTf), id)...)
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
	r.executeWriteSubActionsPost(ctx, diags, operationType, id, thisIdAttributer, subActionsData, requestState, requestPlan, responseState, responsePrivate)

	// Produce a wholly-known new state by determining the final values for any attributes left unknown in the planned state.
	diagsPopulateUnknowns := diag.Diagnostics{} // need separate diags here as diags might already contain error(s) from WriteSubActions
	populateSource := tftypes.Value{}
	if createUpdateResultRaw != nil {
		populateSource = ConvertOdataRawToTerraform(ctx, &diagsPopulateUnknowns, requestPlan.Schema, createUpdateResultRaw, "",
			r.AccessParams.GraphToTerraformMiddleware, id, r.AccessParams.GraphToTerraformMiddlewareTargetSetRunOnRawVal)
		if diagsPopulateUnknowns.HasError() {
			diags.Append(diagsPopulateUnknowns...)
			return
		}
	}
	diagsPopulateUnknowns.Append(r.populateUnknownValues(ctx, responseState, populateSource)...)
	if diagsPopulateUnknowns.HasError() { // still check here since it looks nicer and cleaner
		diags.Append(diagsPopulateUnknowns...)
		return
	}

}

func (r *GenericResource) checkAndLockMutex() bool {
	if r.AccessParams.WriteOptions.SerializeWrites {
		if !r.serializationMutex.TryLock() {
			r.serializationMutex.Lock()
			time.Sleep(r.AccessParams.WriteOptions.SerialWritesDelay)
		}
	}
	return r.AccessParams.WriteOptions.SerializeWrites
}

func (r *GenericResource) executeWriteSubActionsPre(ctx context.Context, diags *diag.Diagnostics, wsaOperation OperationType,
	id string, idAttributer GetAttributer, rawVal map[string]any, reqState *tfsdk.State, reqPlan *tfsdk.Plan,
	respPrivate PrivateDataGetSetter) map[string]any {

	subActionsData := make(map[string]any)
	for _, a := range r.AccessParams.WriteOptions.SubActions {
		if a.CheckRunAction(wsaOperation) {
			wsaRequest := WriteSubActionRequest{
				GenRes:         r,
				Operation:      wsaOperation,
				Id:             id,
				IdAttributer:   idAttributer,
				RawVal:         rawVal,
				SubActionsData: subActionsData,
				ReqState:       reqState,
				ReqPlan:        reqPlan,
				RespPrivate:    respPrivate,
			}
			a.ExecutePre(ctx, diags, &wsaRequest)
			if diags.HasError() {
				return subActionsData
			}
		}
	}

	return subActionsData
}

func (r *GenericResource) executeWriteSubActionsPost(ctx context.Context, diags *diag.Diagnostics, wsaOperation OperationType,
	id string, idAttributer GetAttributer, subActionsData map[string]any, reqState *tfsdk.State, reqPlan *tfsdk.Plan, respState *tfsdk.State,
	respPrivate PrivateDataGetSetter) {

	for _, a := range r.AccessParams.WriteOptions.SubActions {
		if a.CheckRunAction(wsaOperation) {
			wsaRequest := WriteSubActionRequest{
				GenRes:         r,
				Operation:      wsaOperation,
				Id:             id,
				IdAttributer:   idAttributer,
				SubActionsData: subActionsData,
				ReqState:       reqState,
				ReqPlan:        reqPlan,
				RespState:      respState,
				RespPrivate:    respPrivate,
			}
			a.ExecutePost(ctx, diags, &wsaRequest)
			if diags.HasError() {
				return
			}
		}
	}
}
