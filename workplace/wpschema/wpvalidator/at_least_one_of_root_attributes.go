package wpvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// DataSourceAtLeastOneOfRootAttributes is an efficient replacement for
// datasourcevalidator.AtLeastOneOf when only exact root attributes need to be
// checked.
func DataSourceAtLeastOneOfRootAttributes(attributeNames ...string) datasource.ConfigValidator {
	return atLeastOneOfRootAttributesValidator{attributeNames: attributeNames}
}

// ResourceAtLeastOneOfRootAttributes is an efficient replacement for
// resourcevalidator.AtLeastOneOf when only exact root attributes need to be
// checked.
func ResourceAtLeastOneOfRootAttributes(attributeNames ...string) resource.ConfigValidator {
	return atLeastOneOfRootAttributesValidator{attributeNames: attributeNames}
}

var _ datasource.ConfigValidator = atLeastOneOfRootAttributesValidator{}
var _ resource.ConfigValidator = atLeastOneOfRootAttributesValidator{}

type atLeastOneOfRootAttributesValidator struct {
	attributeNames []string
}

func (v atLeastOneOfRootAttributesValidator) Description(_ context.Context) string {
	return fmt.Sprintf("At least one of these root attributes must be configured: %s", formatSiblingNames(v.attributeNames))
}

func (v atLeastOneOfRootAttributesValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v atLeastOneOfRootAttributesValidator) ValidateDataSource(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	resp.Diagnostics.Append(v.validate(ctx, req.Config)...)
}

func (v atLeastOneOfRootAttributesValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	resp.Diagnostics.Append(v.validate(ctx, req.Config)...)
}

func (v atLeastOneOfRootAttributesValidator) validate(ctx context.Context, config tfsdk.Config) diag.Diagnostics {
	var diags diag.Diagnostics
	configuredCount := 0

	for _, name := range v.attributeNames {
		attributePath := path.Root(name)
		var value attr.Value

		getAttributeDiags := config.GetAttribute(ctx, attributePath, &value)
		diags.Append(getAttributeDiags...)
		if getAttributeDiags.HasError() {
			continue
		}

		if value.IsUnknown() {
			return diags
		}

		if value.IsNull() {
			continue
		}

		configuredCount++
	}

	if configuredCount == 0 && !diags.HasError() {
		diags.Append(diag.NewErrorDiagnostic(
			"Missing Attribute Configuration",
			v.Description(ctx),
		))
	}

	return diags
}
