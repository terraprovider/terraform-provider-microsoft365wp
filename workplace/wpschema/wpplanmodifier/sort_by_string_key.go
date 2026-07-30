package wpplanmodifier

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ListSortByStringKey returns a plan modifier that sorts list elements by a
// nested string attribute. This ensures deterministic ordering so that both
// the plan (from config) and state (from API) use the same element order,
// preventing false diffs when using ListNestedAttribute instead of SetNestedAttribute.
func ListSortByStringKey(key string) planmodifier.List {
	return listSortByStringKey{key: key}
}

var _ planmodifier.List = (*listSortByStringKey)(nil)

type listSortByStringKey struct {
	key string
}

func (pm listSortByStringKey) Description(_ context.Context) string {
	return "Sorts list elements by a string key attribute to ensure deterministic ordering."
}

func (pm listSortByStringKey) MarkdownDescription(ctx context.Context) string {
	return pm.Description(ctx)
}

func (pm listSortByStringKey) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if resp.PlanValue.IsNull() || resp.PlanValue.IsUnknown() {
		return
	}

	elements := resp.PlanValue.Elements()
	if len(elements) <= 1 {
		return
	}

	genericSortElements(elements, pm.key)

	sorted, diags := types.ListValue(resp.PlanValue.ElementType(ctx), elements)
	resp.Diagnostics.Append(diags...)
	if !resp.Diagnostics.HasError() {
		resp.PlanValue = sorted
	}
}

// SetSortByStringKey returns a plan modifier that sorts set elements by a
// nested string attribute before converting to the plan value. This can be
// used to normalize set ordering for more predictable plans.
func SetSortByStringKey(key string) planmodifier.Set {
	return setSortByStringKey{key: key}
}

var _ planmodifier.Set = (*setSortByStringKey)(nil)

type setSortByStringKey struct {
	key string
}

func (pm setSortByStringKey) Description(_ context.Context) string {
	return "Sorts set elements by a string key attribute."
}

func (pm setSortByStringKey) MarkdownDescription(ctx context.Context) string {
	return pm.Description(ctx)
}

func (pm setSortByStringKey) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if resp.PlanValue.IsNull() || resp.PlanValue.IsUnknown() {
		return
	}

	elements := resp.PlanValue.Elements()
	if len(elements) <= 1 {
		return
	}

	genericSortElements(elements, pm.key)

	sorted, diags := types.SetValue(resp.PlanValue.ElementType(ctx), elements)
	resp.Diagnostics.Append(diags...)
	if !resp.Diagnostics.HasError() {
		resp.PlanValue = sorted
	}
}

// genericSortElements is a shared helper for sorting []attr.Value by a string key.
func genericSortElements(elements []attr.Value, key string) {
	sort.SliceStable(elements, func(i, j int) bool {
		iObj, iOk := elements[i].(types.Object)
		jObj, jOk := elements[j].(types.Object)
		if !iOk || !jOk {
			return false
		}
		iId, iOk := iObj.Attributes()[key].(types.String)
		jId, jOk := jObj.Attributes()[key].(types.String)
		if !iOk || !jOk || iId.IsNull() || iId.IsUnknown() || jId.IsNull() || jId.IsUnknown() {
			return false
		}
		return iId.ValueString() < jId.ValueString()
	})
}
