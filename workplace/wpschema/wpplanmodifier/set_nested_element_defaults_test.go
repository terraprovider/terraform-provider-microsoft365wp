package wpplanmodifier

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSetNestedElementDefaults(t *testing.T) {
	ctx := context.Background()
	attrTypes := map[string]attr.Type{"web_url": types.StringType, "target_environment": types.StringType, "allow_redirect": types.BoolType}
	elemType := types.ObjectType{AttrTypes: attrTypes}
	obj := func(url string, env, redirect attr.Value) attr.Value {
		return types.ObjectValueMust(attrTypes, map[string]attr.Value{"web_url": types.StringValue(url), "target_environment": env, "allow_redirect": redirect})
	}
	pm := SetNestedElementDefaults(map[string]attr.Value{
		"target_environment": types.StringValue("internetExplorerMode"),
		"allow_redirect":     types.BoolValue(false),
	})

	// the scenario of the live test: config with defaults omitted, the proposed plan filled from the prior state
	config := types.SetValueMust(elemType, []attr.Value{
		obj("erp", types.StringNull(), types.BoolNull()),
		obj("legacy", types.StringValue("internetExplorerMode"), types.BoolValue(true)),
		obj("new", types.StringValue("microsoftEdge"), types.BoolNull()),
		obj("unknown", types.StringUnknown(), types.BoolNull()),
	})
	proposed := types.SetValueMust(elemType, []attr.Value{
		obj("erp", types.StringValue("internetExplorerMode"), types.BoolValue(false)),
		obj("legacy", types.StringValue("internetExplorerMode"), types.BoolValue(true)),
		obj("new", types.StringValue("microsoftEdge"), types.BoolValue(true)), // allow_redirect removed from config
		obj("unknown", types.StringUnknown(), types.BoolUnknown()),
	})
	want := types.SetValueMust(elemType, []attr.Value{
		obj("erp", types.StringValue("internetExplorerMode"), types.BoolValue(false)),
		obj("legacy", types.StringValue("internetExplorerMode"), types.BoolValue(true)),
		obj("new", types.StringValue("microsoftEdge"), types.BoolValue(false)),
		obj("unknown", types.StringUnknown(), types.BoolValue(false)),
	})

	resp := &planmodifier.SetResponse{PlanValue: proposed}
	pm.PlanModifySet(ctx, planmodifier.SetRequest{ConfigValue: config, PlanValue: proposed}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}
	if !resp.PlanValue.Equal(want) {
		t.Errorf("got %s, want %s", resp.PlanValue, want)
	}

	// null config (default of the set applies) and unknown config are left alone
	for _, cfg := range []types.Set{types.SetNull(elemType), types.SetUnknown(elemType)} {
		resp := &planmodifier.SetResponse{PlanValue: proposed}
		pm.PlanModifySet(ctx, planmodifier.SetRequest{ConfigValue: cfg, PlanValue: proposed}, resp)
		if !resp.PlanValue.Equal(proposed) {
			t.Errorf("config %s: plan changed to %s", cfg, resp.PlanValue)
		}
	}
}
