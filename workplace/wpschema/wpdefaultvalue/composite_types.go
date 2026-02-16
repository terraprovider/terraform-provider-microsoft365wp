package wpdefaultvalue

import (
	"context"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpschemautil"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//
// ObjectDefaultValue / objectDefaultValue
//

var _ DefaultValueWithInit = &objectDefaultValue{}

type objectDefaultValue struct {
	defaultValueBase
	defaultValFw types.Object
}

func ObjectDefaultValue(defaultValueRaw map[string]any) defaults.Object {
	return &objectDefaultValue{defaultValueBase: defaultValueBase{defaultValRaw: defaultValueRaw}}
}

func (d *objectDefaultValue) Init(ctx context.Context, diags *diag.Diagnostics, schema *rsschema.Schema, path path.Path) {
	d.defaultValFw = wpschemautil.GetFwTypedValueFromAny[types.Object](ctx, diags, schema, path, d.defaultValRaw)
}

func (d objectDefaultValue) DefaultObject(ctx context.Context, req defaults.ObjectRequest, resp *defaults.ObjectResponse) {
	resp.PlanValue = d.defaultValFw
}

//
// SetDefaultValue / setDefaultValue
//

var _ DefaultValueWithInit = &setDefaultValue{}

type setDefaultValue struct {
	defaultValueBase
	defaultValFw types.Set
}

func SetDefaultValue(defaultValueRaw []any) defaults.Set {
	return &setDefaultValue{defaultValueBase: defaultValueBase{defaultValRaw: defaultValueRaw}}
}

func (d *setDefaultValue) Init(ctx context.Context, diags *diag.Diagnostics, schema *rsschema.Schema, path path.Path) {
	d.defaultValFw = wpschemautil.GetFwTypedValueFromAny[types.Set](ctx, diags, schema, path, d.defaultValRaw)
}

func (d setDefaultValue) DefaultSet(ctx context.Context, req defaults.SetRequest, resp *defaults.SetResponse) {
	resp.PlanValue = d.defaultValFw
}

//
// ListDefaultValue / listDefaultValue
//

var _ DefaultValueWithInit = &listDefaultValue{}

type listDefaultValue struct {
	defaultValueBase
	defaultValFw types.List
}

func ListDefaultValue(defaultValueRaw []any) defaults.List {
	return &listDefaultValue{defaultValueBase: defaultValueBase{defaultValRaw: defaultValueRaw}}
}

func (d *listDefaultValue) Init(ctx context.Context, diags *diag.Diagnostics, schema *rsschema.Schema, path path.Path) {
	d.defaultValFw = wpschemautil.GetFwTypedValueFromAny[types.List](ctx, diags, schema, path, d.defaultValRaw)
}

func (d listDefaultValue) DefaultList(ctx context.Context, req defaults.ListRequest, resp *defaults.ListResponse) {
	resp.PlanValue = d.defaultValFw
}
