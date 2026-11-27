package wpplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// SetAlignToConfigOrder returns a plan modifier for a SetNestedAttribute that orders the planned elements like the
// config elements. It never changes any value. It must run before the plan modifiers of the nested attributes (which
// the framework guarantees for set-level plan modifiers).
//
// Background (see issue #246): the framework runs the plan modifiers of the nested attributes of a set element with
// the config element at the same INDEX. Terraform/OpenTofu fill null Optional+Computed attributes of the proposed
// elements from the prior state, and cty orders set elements by value, so the proposed set can be ordered differently
// from the config. Default value plan modifiers (wpdefaultvaluemodifier) then see the null config value of another
// element and overwrite a configured value with the default (e.g. mobile_app assignment intent "uninstall" ->
// "required" for an unchanged config).
//
// Each config element is matched to the planned element it was derived from: all known, non-null config values must
// equal the planned values (null config values may have been filled from the prior state or by defaults). Exact
// matches are assigned first, so that a config element with many null values does not take the element of a more
// specific one. Planned elements without a match are kept at the end.
func SetAlignToConfigOrder() planmodifier.Set {
	return setAlignToConfigOrder{}
}

// IsSetAlignToConfigOrder reports whether m is a plan modifier returned by SetAlignToConfigOrder.
func IsSetAlignToConfigOrder(m planmodifier.Set) bool {
	_, ok := m.(setAlignToConfigOrder)
	return ok
}

var _ planmodifier.Set = (*setAlignToConfigOrder)(nil)

type setAlignToConfigOrder struct{}

func (pm setAlignToConfigOrder) Description(_ context.Context) string {
	return "Orders the planned set elements like the config elements, so that the plan modifiers of nested attributes get the matching config element."
}

func (pm setAlignToConfigOrder) MarkdownDescription(ctx context.Context) string {
	return pm.Description(ctx)
}

func (pm setAlignToConfigOrder) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || resp.PlanValue.IsNull() || resp.PlanValue.IsUnknown() {
		return
	}

	configElements := req.ConfigValue.Elements()
	planElements := resp.PlanValue.Elements()
	if len(planElements) <= 1 {
		return
	}

	configTf := make([]tftypes.Value, len(configElements))
	for i, e := range configElements {
		v, err := e.ToTerraformValue(ctx)
		if err != nil {
			return // leave the plan unchanged, the framework will report invalid values anyway
		}
		configTf[i] = v
	}
	planTf := make([]tftypes.Value, len(planElements))
	for i, e := range planElements {
		v, err := e.ToTerraformValue(ctx)
		if err != nil {
			return
		}
		planTf[i] = v
	}

	order := alignOrder(configTf, planTf)
	if order == nil {
		return
	}

	aligned := make([]attr.Value, 0, len(planElements))
	for _, i := range order {
		aligned = append(aligned, planElements[i])
	}

	planned, diags := types.SetValue(resp.PlanValue.ElementType(ctx), aligned)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.PlanValue = planned
}

// alignOrder returns the indexes of plan in the order of their matching config elements, followed by the indexes of
// the plan elements without a match (in their original order). It returns nil if the order does not change.
func alignOrder(config []tftypes.Value, plan []tftypes.Value) []int {
	planUsed := make([]bool, len(plan))
	planForConfig := make([]int, len(config))
	for i := range planForConfig {
		planForConfig[i] = -1
	}

	assign := func(match func(c, p tftypes.Value) bool) {
		for ci, c := range config {
			if planForConfig[ci] >= 0 {
				continue
			}
			for pi, p := range plan {
				if !planUsed[pi] && match(c, p) {
					planForConfig[ci] = pi
					planUsed[pi] = true
					break
				}
			}
		}
	}
	assign(func(c, p tftypes.Value) bool { return c.Equal(p) })
	assign(configMatchesPlan)

	order := make([]int, 0, len(plan))
	for _, pi := range planForConfig {
		if pi >= 0 {
			order = append(order, pi)
		}
	}
	for pi := range plan {
		if !planUsed[pi] {
			order = append(order, pi)
		}
	}

	for i, pi := range order {
		if i != pi {
			return order
		}
	}
	return nil
}

// configMatchesPlan reports whether a planned value can have been derived from a config value: unknown and null config
// values match anything (they get filled from the prior state, by defaults or as computed values), everything else
// must be equal (recursively for objects and collections).
func configMatchesPlan(c, p tftypes.Value) bool {
	if !c.IsKnown() || c.IsNull() {
		return true
	}
	if !p.IsKnown() || p.IsNull() {
		return false
	}

	typ := c.Type()
	switch {
	case typ.Is(tftypes.Object{}) || typ.Is(tftypes.Map{}):
		var cm, pm map[string]tftypes.Value
		if c.As(&cm) != nil || p.As(&pm) != nil {
			return false
		}
		if typ.Is(tftypes.Map{}) && len(cm) != len(pm) {
			return false
		}
		for k, cv := range cm {
			pv, ok := pm[k]
			if !ok || !configMatchesPlan(cv, pv) {
				return false
			}
		}
		return true

	case typ.Is(tftypes.List{}) || typ.Is(tftypes.Tuple{}):
		var cl, pl []tftypes.Value
		if c.As(&cl) != nil || p.As(&pl) != nil || len(cl) != len(pl) {
			return false
		}
		for i := range cl {
			if !configMatchesPlan(cl[i], pl[i]) {
				return false
			}
		}
		return true

	case typ.Is(tftypes.Set{}):
		var cs, ps []tftypes.Value
		if c.As(&cs) != nil || p.As(&ps) != nil || len(cs) != len(ps) {
			return false
		}
		used := make([]bool, len(ps))
		for _, cv := range cs {
			found := false
			for i, pv := range ps {
				if !used[i] && configMatchesPlan(cv, pv) {
					used[i] = true
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true

	default:
		return c.Equal(p)
	}
}
