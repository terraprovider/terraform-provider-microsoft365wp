package wpvalidator

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ExactlyOneOfSiblings is an efficient replacement for objectvalidator.ExactlyOneOf
// when used with path.MatchRelative().AtParent().AtName("...") expressions.
//
// The standard ExactlyOneOf uses PathMatches which does a full tftypes.Walk of
// the entire config tree for EACH sibling check. For deeply nested schemas with
// many elements (e.g. 17 rules × 5 derived types), this causes O(N*M²) walks
// where N is the number of elements and M the number of derived types.
//
// This implementation avoids PathMatches and repeated sibling lookups by fetching
// the parent object once and accessing siblings directly from its attributes map,
// making validation O(M) per element.
func ExactlyOneOfSiblings(siblingNames ...string) validator.Object {
	return exactlyOneOfSiblingsValidator{siblingNames: siblingNames}
}

var _ validator.Object = exactlyOneOfSiblingsValidator{}

type exactlyOneOfSiblingsValidator struct {
	siblingNames []string
}

func (v exactlyOneOfSiblingsValidator) Description(_ context.Context) string {
	return fmt.Sprintf("Exactly one of these sibling attributes must be set: %s", formatSiblingNames(v.siblingNames))
}

func (v exactlyOneOfSiblingsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v exactlyOneOfSiblingsValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	// If current attribute is unknown, delay validation
	if req.ConfigValue.IsUnknown() {
		return
	}

	// Determine the parent path by removing the last step (this attribute's name)
	parentPath := req.Path.ParentPath()

	// Fetch the parent object once instead of calling GetAttribute for each sibling
	var parentObj types.Object
	diags := req.Config.GetAttribute(ctx, parentPath, &parentObj)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// If parent is null or unknown, delay validation
	if parentObj.IsNull() || parentObj.IsUnknown() {
		return
	}

	// Access parent's attributes map
	parentAttrs := parentObj.Attributes()

	// Count how many configured sibling attributes are non-null.
	count := 0

	// Inspect sibling attributes from the parent object's attributes map
	for _, name := range v.siblingNames {
		siblingVal, exists := parentAttrs[name]
		// Skip if sibling doesn't exist in parent (might be optional)
		if !exists {
			continue
		}

		// Delay validation if any sibling is unknown
		if siblingVal.IsUnknown() {
			return
		}

		if !siblingVal.IsNull() {
			count++
		}
	}

	if count == 0 {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			req.Path,
			"Invalid Attribute Combination",
			fmt.Sprintf("Exactly one of %s must be set.", formatSiblingNames(v.siblingNames)),
		))
	}

	if count > 1 {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			req.Path,
			"Invalid Attribute Combination",
			fmt.Sprintf("Exactly one of %s must be set; %d were provided.", formatSiblingNames(v.siblingNames), count),
		))
	}
}

func formatSiblingNames(names []string) string {
	if len(names) == 0 {
		return "[]"
	}

	formatted := make([]string, len(names))
	for i, name := range names {
		formatted[i] = fmt.Sprintf("`%s`", name)
	}

	return strings.Join(formatted, ", ")
}
