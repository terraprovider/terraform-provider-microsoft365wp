package generic

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

//
// THIS IMPLEMENTATION WILL ONLY WORK FOR CHILD ELEMENTS THAT DO NOT CHANGE. They cannot even get a server-assigned id.
// Background is that Terraform will not be able to match elements with any small change to the respective plan element
// since element comparison within unordered sets takes place by comparing the whole value - hence it must not change.
// If you need something similar with server-assigned id (or similar) then implement them as separate objects (like GetAssignmentChildResource()).
//

type WriteSubActionIndividual struct {
	WriteSubActionBase
	ComparisonKeyAttribute     string
	SetNestedPath              *tftypes.AttributePath
	IdGetterFunc               func(context.Context, *diag.Diagnostics, map[string]any, string) string
	IsOdataReference           bool
	OdataRefMapTypeToUriPrefix map[string]string
	TerraformToGraphMiddleware func(context.Context, *diag.Diagnostics, *TerraformToGraphMiddlewareParams) TerraformToGraphMiddlewareReturns
}

var _ WriteSubAction = &WriteSubActionIndividual{}

func (a *WriteSubActionIndividual) Initialize() {
	a.WriteSubActionBase.Initialize()
}

func (a *WriteSubActionIndividual) CheckRunAction(wsaOperation OperationType) bool {
	return a.WriteSubActionBase.CheckRunAction(wsaOperation)
}

func (a *WriteSubActionIndividual) ExecutePre(ctx context.Context, diags *diag.Diagnostics, wsaReq *WriteSubActionRequest) {
	a.WriteSubActionBase.ExecutePre(ctx, diags, wsaReq, nil)
}

func (a *WriteSubActionIndividual) ExecutePost(ctx context.Context, diags *diag.Diagnostics, wsaReq *WriteSubActionRequest) {

	const ErrorSummary = "Error in WriteSubActionIndividual.ExecutePost"

	errorDetailPrefix := "Unable to retrieve state: "
	var stateTfSlice []tftypes.Value
	switch wsaReq.Operation {
	case OperationCreate:
		stateTfSlice = make([]tftypes.Value, 0)
	case OperationUpdate:
		stateTfRaw, _, _ := tftypes.WalkAttributePath(wsaReq.ReqState.Raw, a.SetNestedPath)
		stateTf, ok := stateTfRaw.(tftypes.Value)
		if !ok {
			diags.AddError(ErrorSummary, fmt.Sprintf(errorDetailPrefix+"stateTfRaw is not tftypes.Value but %T", stateTfRaw))
			return
		}
		if err := stateTf.As(&stateTfSlice); err != nil {
			diags.AddError(ErrorSummary, fmt.Sprintf(errorDetailPrefix+"%s", err.Error()))
			return
		}
	}

	errorDetailPrefix = "Unable to retrieve plan: "
	var planTfSlice []tftypes.Value
	planTfRaw, _, _ := tftypes.WalkAttributePath(wsaReq.ReqPlan.Raw, a.SetNestedPath)
	planTf, ok := planTfRaw.(tftypes.Value)
	if !ok {
		diags.AddError(ErrorSummary, fmt.Sprintf(errorDetailPrefix+"planTfRaw is not tftypes.Value but %T", planTfRaw))
		return
	}
	if err := planTf.As(&planTfSlice); err != nil {
		diags.AddError(ErrorSummary, fmt.Sprintf(errorDetailPrefix+"%s", err.Error()))
		return
	}

	planByKey := make(map[string]tftypes.Value)
	stateByKey := make(map[string]tftypes.Value)
	for _, v := range planTfSlice {
		errorDetailPrefix = fmt.Sprintf("Unable to retrieve plan element %s: ", v.String())
		var vMap map[string]tftypes.Value
		if err := v.As(&vMap); err != nil {
			diags.AddError(ErrorSummary, fmt.Sprintf(errorDetailPrefix+"%s", err.Error()))
			return
		}
		if keyRaw, ok := vMap[a.ComparisonKeyAttribute]; ok {
			keyString := keyRaw.String() // these are temporary keys for comparison only, so use String() to not require string attributes (but also e.g. numbers etc.)
			if keyString != "" {
				planByKey[keyString] = v
			}
		}
	}
	for _, v := range stateTfSlice {
		errorDetailPrefix = fmt.Sprintf("Unable to retrieve state element %s: ", v.String())
		var vMap map[string]tftypes.Value
		if err := v.As(&vMap); err != nil {
			diags.AddError(ErrorSummary, fmt.Sprintf(errorDetailPrefix+"%s", err.Error()))
			return
		}
		if keyRaw, ok := vMap[a.ComparisonKeyAttribute]; ok {
			keyString := keyRaw.String() // see above
			if keyString != "" {
				stateByKey[keyString] = v
			}
		}
	}

	elementsToAdd := make([]tftypes.Value, 0)
	elementsToUpdate := make([]tftypes.Value, 0)
	elementsToDelete := make([]tftypes.Value, 0)
	for k, vPlan := range planByKey {
		if vState, ok := stateByKey[k]; ok {
			if !vPlan.Equal(vState) {
				if !a.IsOdataReference {
					elementsToUpdate = append(elementsToUpdate, vPlan)
				} else {
					// OData references do not support updates and they should happen in very rare cases only anyway
					// (i.e. in case a practitioner had assigned the wrong type to a reference and changes/fixes that).
					diags.AddWarning("Warning from WriteSubActionIndividual.ExecutePost", "The element with id "+k+" has changed, but this is "+
						"not supported for OData references and should not happend under normal circumstances. Deleting and re-adding it instead.")
					elementsToDelete = append(elementsToDelete, vState)
					elementsToAdd = append(elementsToAdd, vPlan)
				}
			}
		} else {
			elementsToAdd = append(elementsToAdd, vPlan)
		}
	}
	for k, vState := range stateByKey {
		if _, ok := planByKey[k]; !ok {
			elementsToDelete = append(elementsToDelete, vState)
		}
	}

	if len(elementsToAdd) == 0 && len(elementsToUpdate) == 0 && len(elementsToDelete) == 0 {
		// nothing to do
		return
	}

	parentId := wsaReq.GenRes.AccessParams.GetId(ctx, diags, wsaReq.Id, wsaReq.IdAttributer)
	if diags.HasError() {
		return
	}

	parentUri := wsaReq.GenRes.AccessParams.GetUriWithIdForUD(ctx, diags, "", parentId, wsaReq.IdAttributer)
	if diags.HasError() {
		return
	}
	childBaseUri := fmt.Sprintf("%s/%s", parentUri.Entity, a.UriSuffix)

	setPlanSchemaRaw, _, _ := tftypes.WalkAttributePath(wsaReq.ReqPlan.Schema, a.SetNestedPath)
	setPlanSchema, ok := setPlanSchemaRaw.(schema.SetNestedAttribute)
	if !ok {
		diags.AddError(ErrorSummary, fmt.Sprintf("Unable to retrieve plan sub-schema: setPlanSchemaRaw is not schema.SetNestedAttribute but %T", setPlanSchemaRaw))
		return
	}
	elementTranslator := NewToFromGraphTranslator(setPlanSchema, false)

	// We do not check for errors after Graph operations but continue to try to finish other tasks as much as possible.
	// Our caller will check for errors and exit later.

	addElementsFunc := func() {

		createUri := childBaseUri
		if a.IsOdataReference {
			createUri += "/$ref"
		}

		for _, v := range elementsToAdd {

			var diags2 diag.Diagnostics
			vRaw, _ := a.getElementRawValAndId(ctx, &diags2, elementTranslator, v, parentId, false, false)
			diags.Append(diags2...)
			if diags2.HasError() {
				return
			}

			wsaReq.GenRes.AccessParams.CreateRaw(ctx, diags, createUri, nil, vRaw, true)
		}
	}

	updateElementsFunc := func() {

		// no need to check fo a.IsOdataReference here as updates are not supported for it

		for _, v := range elementsToUpdate {

			var diags2 diag.Diagnostics
			vRaw, id := a.getElementRawValAndId(ctx, &diags2, elementTranslator, v, parentId, true, false)
			diags.Append(diags2...)
			if diags2.HasError() {
				return
			}

			wsaReq.GenRes.AccessParams.UpdateRaw(ctx, diags, childBaseUri, id, nil, vRaw, false)
		}
	}

	deleteElementsFunc := func() {

		for _, v := range elementsToDelete {

			var diags2 diag.Diagnostics
			_, id := a.getElementRawValAndId(ctx, &diags2, elementTranslator, v, parentId, false, true)
			diags.Append(diags2...)
			if diags2.HasError() {
				return
			}

			if a.IsOdataReference {
				id += "/$ref"
			}

			wsaReq.GenRes.AccessParams.DeleteRaw(ctx, diags, childBaseUri, id, nil)
		}
	}

	if a.IsOdataReference {
		// delete first and then add to avoid problems duplicate elements on updates
		deleteElementsFunc()
		addElementsFunc()
	} else {
		// add, update and then delete to avoid problems with (missing) default elements etc.
		addElementsFunc()
		updateElementsFunc()
		deleteElementsFunc()
	}
}

func (a *WriteSubActionIndividual) getElementRawValAndId(ctx context.Context, diags *diag.Diagnostics,
	elementTranslator ToFromGraphTranslator, v tftypes.Value, parentId string, isUpdate bool, isDelete bool) (map[string]any, string) {

	const ErrorSummary = "Error in WriteSubActionIndividual.getElementRawValAndId"

	vRaw, err := elementTranslator.TerraformAsRaw(ctx, v)
	if err != nil {
		diags.AddError(ErrorSummary, fmt.Sprintf("Unable to create raw value from Terraform element %s: %s", v, err.Error()))
		return nil, ""
	}

	var id string
	if isUpdate || isDelete {
		id = a.getElementId(ctx, diags, vRaw, parentId)
		if diags.HasError() {
			return nil, ""
		}
	}

	if !isDelete {

		if a.IsOdataReference {
			// OData refs need the full OData URI of the target entity (instead of just a GUID id)

			odataType, ok := vRaw["@odata.type"].(string)
			if ok {
				delete(vRaw, "@odata.type")
			}
			uriPrefix, ok := a.OdataRefMapTypeToUriPrefix[odataType]
			if !ok {
				uriPrefix, ok = a.OdataRefMapTypeToUriPrefix[""]
				if !ok {
					diags.AddError(ErrorSummary, fmt.Sprintf("Unable to find OData URI prefix for type %s in OdataRefMapTypeToUriPrefix", odataType))
					return nil, ""
				}
			}

			id, ok := vRaw["id"].(string)
			if !ok {
				diags.AddError(ErrorSummary, fmt.Sprintf("Unable to find id in element %v", vRaw))
				return nil, ""
			}
			delete(vRaw, "id")

			vRaw["@odata.id"] = uriPrefix + id
		}

		if a.TerraformToGraphMiddleware != nil {
			params := TerraformToGraphMiddlewareParams{RawVal: vRaw, IsUpdate: isUpdate}
			if err := a.TerraformToGraphMiddleware(ctx, diags, &params); err != nil {
				diags.AddError(ErrorSummary, fmt.Sprintf("Error running middleware for element %s: %s", v, err.Error()))
			}
			if diags.HasError() {
				return nil, ""
			}
		}
	}

	return vRaw, id
}

func (a *WriteSubActionIndividual) getElementId(ctx context.Context, diags *diag.Diagnostics, vRaw map[string]any, parentId string) string {

	const ErrorSummary = "Error in WriteSubActionIndividual.getElementId"

	if a.IdGetterFunc != nil {

		return a.IdGetterFunc(ctx, diags, vRaw, parentId)

	} else {
		idRaw, ok := vRaw["id"]
		if !ok {
			diags.AddError(ErrorSummary, fmt.Sprintf("Unable to find attribute `id` in state element %s", vRaw))
			return ""
		}

		idString, ok := idRaw.(string)
		if !ok {
			diags.AddError(ErrorSummary, fmt.Sprintf("Value of attribute `id` in state element is not string but %T", idRaw))
			return ""
		}

		return idString
	}
}
