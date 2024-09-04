package wpvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Dummy validator objects to signal that an attribute value being sent to or received from MS Graph might need to be
// translated to some other value for Terraform.

// We cannot use plan modifiers here since
// - this should also work for data sources and they do not have plan modifiers (but they do have validators)
// - they will only work under very specific conditions and otherwise throw "Provider produced invalid plan"

//
// AttributeValueTranslator
//

const (
	AttributeMising = attributeMissing(0)
)

type attributeMissing byte

type AttributeValueTranslator interface {
	TranslateTerraformToGraph(tfValue tftypes.Value) (any, bool, error)
	TranslateGraphToTerraform(graphValue any) (tftypes.Value, bool, error)
}

//
// translateValuesDefault
//

var _ validator.Bool = (*translateValuesDefault)(nil)
var _ validator.Float64 = (*translateValuesDefault)(nil)
var _ validator.String = (*translateValuesDefault)(nil)
var _ validator.List = (*translateValuesDefault)(nil)
var _ validator.Set = (*translateValuesDefault)(nil)
var _ validator.Map = (*translateValuesDefault)(nil)
var _ validator.Object = (*translateValuesDefault)(nil)

type translateValuesDefault struct {
}

func (pm translateValuesDefault) Description(ctx context.Context) string {
	return "TranslateValues"
}
func (pm translateValuesDefault) MarkdownDescription(ctx context.Context) string {
	return pm.Description(ctx)
}

// nothing to do here as it's just used as an internal flag on the attribute
func (translateValuesDefault) ValidateBool(ctx context.Context, req validator.BoolRequest, res *validator.BoolResponse) {
}
func (translateValuesDefault) ValidateFloat64(ctx context.Context, req validator.Float64Request, res *validator.Float64Response) {
}
func (translateValuesDefault) ValidateInt64(ctx context.Context, req validator.Int64Request, res *validator.Int64Response) {
}
func (translateValuesDefault) ValidateString(ctx context.Context, req validator.StringRequest, res *validator.StringResponse) {
}
func (translateValuesDefault) ValidateList(ctx context.Context, req validator.ListRequest, res *validator.ListResponse) {
}
func (translateValuesDefault) ValidateSet(ctx context.Context, req validator.SetRequest, res *validator.SetResponse) {
}
func (translateValuesDefault) ValidateMap(ctx context.Context, req validator.MapRequest, res *validator.MapResponse) {
}
func (translateValuesDefault) ValidateObject(ctx context.Context, req validator.ObjectRequest, res *validator.ObjectResponse) {
}
