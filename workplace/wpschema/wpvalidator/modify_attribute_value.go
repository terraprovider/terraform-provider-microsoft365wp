package wpvalidator

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Dummy validator objects to signal that an attribute value being sent to or received from MS Graph might need to be
// translated to some other value for Terraform.

// We cannot use plan modifiers here since
// - this should also work for data sources and they do not have plan modifiers (but they do have validators)
// - they will only work under very specific conditions and otherwise throw "Provider produced invalid plan"

//
// AttributeValueModifier
//

type AttributeValueModifier interface {
	ModifyGraphToTerraform(graphValue any, priorState *tfsdk.State, path *tftypes.AttributePath) (any, error)
}

//
// modifyValueDefault
//

var _ validator.List = (*modifyValueDefault)(nil)

type modifyValueDefault struct {
}

func (pm modifyValueDefault) Description(ctx context.Context) string {
	return "ModifyValues"
}
func (pm modifyValueDefault) MarkdownDescription(ctx context.Context) string {
	return pm.Description(ctx)
}

// nothing to do here as it's just used as an internal flag on the attribute
func (modifyValueDefault) ValidateList(ctx context.Context, req validator.ListRequest, res *validator.ListResponse) {
}

//
// ModifyValueSortSet — reorders a raw JSON list to match the prior Terraform state order.
//

var _ AttributeValueModifier = (*ModifyValueSortSetImpl)(nil)

type ModifyValueSortSetImpl struct {
	modifyValueDefault
	key string // key field name used for matching elements (e.g. "id")
}

func ModifyValuesGraphToTerraformReorderByKey(key string) ModifyValueSortSetImpl {
	return ModifyValueSortSetImpl{key: key}
}

func (t ModifyValueSortSetImpl) ModifyGraphToTerraform(graphValue any, priorState *tfsdk.State, path *tftypes.AttributePath) (any, error) {
	rawList, ok := graphValue.([]any)
	if !ok || len(rawList) <= 1 {
		return graphValue, nil
	}

	// Derive TF attribute name from path
	tfAttr := ""
	if steps := path.Steps(); len(steps) > 0 {
		if an, ok := steps[len(steps)-1].(tftypes.AttributeName); ok {
			tfAttr = string(an)
		}
	}

	// Extract element order from prior state
	var stateOrder []string
	if tfAttr != "" && priorState != nil && !priorState.Raw.IsNull() && priorState.Raw.IsKnown() {
		var stateMap map[string]tftypes.Value
		if err := priorState.Raw.As(&stateMap); err == nil {
			if listVal, ok := stateMap[tfAttr]; ok && !listVal.IsNull() && listVal.IsKnown() {
				var elements []tftypes.Value
				if err := listVal.As(&elements); err == nil {
					for _, elem := range elements {
						var elemMap map[string]tftypes.Value
						if err := elem.As(&elemMap); err != nil {
							continue
						}
						if keyVal, ok := elemMap[t.key]; ok && !keyVal.IsNull() && keyVal.IsKnown() {
							var key string
							if err := keyVal.As(&key); err == nil {
								stateOrder = append(stateOrder, key)
							}
						}
					}
				}
			}
		}
	}

	if len(stateOrder) == 0 {
		// No prior state — fall back to alphabetical sort by key
		sort.SliceStable(rawList, func(i, j int) bool {
			iMap, iOk := rawList[i].(map[string]any)
			jMap, jOk := rawList[j].(map[string]any)
			if !iOk || !jOk {
				return false
			}
			iVal, _ := iMap[t.key].(string)
			jVal, _ := jMap[t.key].(string)
			return iVal < jVal
		})
		return rawList, nil
	}

	// Build index: key → raw element
	elemByKey := make(map[string]any, len(rawList))
	var unmatchedElems []any
	for _, elem := range rawList {
		if elemMap, ok := elem.(map[string]any); ok {
			if key, ok := elemMap[t.key].(string); ok {
				elemByKey[key] = elem
				continue
			}
		}
		unmatchedElems = append(unmatchedElems, elem)
	}

	// Reorder: first elements matching state order, then remaining sorted alphabetically
	result := make([]any, 0, len(rawList))

	for _, key := range stateOrder {
		if elem, ok := elemByKey[key]; ok {
			result = append(result, elem)
			delete(elemByKey, key)
		}
	}

	// Remaining elements not in prior state, sorted by key
	remaining := make([]string, 0, len(elemByKey))
	for key := range elemByKey {
		remaining = append(remaining, key)
	}
	sort.Strings(remaining)
	for _, key := range remaining {
		result = append(result, elemByKey[key])
	}

	result = append(result, unmatchedElems...)

	return result, nil
}
