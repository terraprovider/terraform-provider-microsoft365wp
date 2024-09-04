package wpvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
)

var _ validator.String = flagEnumValidator{}

// flagEnumValidator validates that the value matches one of expected values.
type flagEnumValidator struct {
	allowedFlagValues []string
}

func (v flagEnumValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v flagEnumValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("all individual values must be one of: %q", v.allowedFlagValues)
}

func (v flagEnumValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {

	configValue := request.ConfigValue
	if configValue.IsNull() || configValue.IsUnknown() || configValue.ValueString() == "" {
		return
	}

	for _, configFlag := range strings.Split(configValue.ValueString(), ",") {
		flagIsValid := false
		for _, allowedFlag := range v.allowedFlagValues {
			if configFlag == allowedFlag {
				flagIsValid = true
				break
			}
		}
		if !flagIsValid {
			response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
				request.Path,
				v.Description(ctx),
				configValue.String(),
			))
		}
	}
}

// FlagEnumValues checks that the String held in the attribute
// is one of the given `values`.
func FlagEnumValues(values ...string) validator.String {
	return flagEnumValidator{
		allowedFlagValues: values,
	}
}
