// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Lifted from https://github.com/hashicorp/terraform-provider-awscc/blob/v0.48.0/internal/generic/default_value.go

package wpdefaultvalue

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type boolDefaultValueAttributePlanModifier struct {
	defaultValueAttributePlanModifier
	val types.Bool
}

func BoolDefaultValue(val bool) planmodifier.Bool {
	return boolDefaultValueAttributePlanModifier{
		val: types.BoolValue(val),
	}
}

func (attributePlanModifier boolDefaultValueAttributePlanModifier) PlanModifyBool(ctx context.Context, request planmodifier.BoolRequest, response *planmodifier.BoolResponse) {
	if request.ConfigValue.IsNull() {
		response.PlanValue = attributePlanModifier.val
	}
}

type float64DefaultValueAttributePlanModifier struct {
	defaultValueAttributePlanModifier
	val types.Float64
}

func Float64DefaultValue(val float64) planmodifier.Float64 {
	return float64DefaultValueAttributePlanModifier{
		val: types.Float64Value(val),
	}
}

func (attributePlanModifier float64DefaultValueAttributePlanModifier) PlanModifyFloat64(ctx context.Context, request planmodifier.Float64Request, response *planmodifier.Float64Response) {
	if request.ConfigValue.IsNull() {
		response.PlanValue = attributePlanModifier.val
	}
}

type int64DefaultValueAttributePlanModifier struct {
	defaultValueAttributePlanModifier
	val types.Int64
}

func Int64DefaultValue(val int64) planmodifier.Int64 {
	return int64DefaultValueAttributePlanModifier{
		val: types.Int64Value(val),
	}
}

func (attributePlanModifier int64DefaultValueAttributePlanModifier) PlanModifyInt64(ctx context.Context, request planmodifier.Int64Request, response *planmodifier.Int64Response) {
	if request.ConfigValue.IsNull() {
		response.PlanValue = attributePlanModifier.val
	}
}

type stringDefaultValueAttributePlanModifier struct {
	defaultValueAttributePlanModifier
	val types.String
}

func StringDefaultValue(val string) planmodifier.String {
	return stringDefaultValueAttributePlanModifier{
		val: types.StringValue(val),
	}
}

func (attributePlanModifier stringDefaultValueAttributePlanModifier) PlanModifyString(ctx context.Context, request planmodifier.StringRequest, response *planmodifier.StringResponse) {
	if request.ConfigValue.IsNull() {
		response.PlanValue = attributePlanModifier.val
	}
}
