// Process: If the config value is null (i.e. no value provided in config) then return the hard-coded default value
// provided by the schema instead.
//
// Why we continue to use our own implementation even though the Framework supports default values natively by now: The
// Framework default value algorithm strictly follows the config and therefore does not work recursively on object
// attributes when the object itself had been set as a default value. Example: If the default for some_attr is an empty
// map/object, then the framework would not visit some_attr.some_prop and therefore it would be impossible to also set a
// default value for a child property.
//
// 2023-05-03 Felix Storm: Previously, we had also checked if there was no current plan value yet (to honor potential
// previous planModifier's result). But since TF seems to always copy the current state to the proposed plan if the
// resource already exists (before applying config values) this rendered the whole idea pretty much useless (as all
// config values not provided will have already been replaced by the corresponding values from the state before any plan
// modifier gets called). I am pretty sure that I had tested out that removing a value from the config will force it to
// become the default value again but cannot really understand now how this could have worked at all before...
//

package wpdefaultvalue

import (
	"context"
	"fmt"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpschemautil"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

//
// defaultValueAttributePlanModifier
//

type defaultValueAttributePlanModifier struct{}

func (attributePlanModifier defaultValueAttributePlanModifier) Description(_ context.Context) string {
	return "If the value of the attribute is missing, then the value is semantically the same as if the value was present with the default value hard-coded in the provider."
}

func (attributePlanModifier defaultValueAttributePlanModifier) MarkdownDescription(ctx context.Context) string {
	return attributePlanModifier.Description(ctx)
}

//
// defaultValuePlanModify
//

func defaultValuePlanModify[T attr.Value](ctx context.Context, diags *diag.Diagnostics, schema rsschema.Schema, path path.Path, configValue T, defaultValueRaw any, responsePlanValue *T) {

	if configValue.IsNull() {

		valueFw := wpschemautil.GetFwValueFromAny(ctx, diags, &schema, path, defaultValueRaw)
		if diags.HasError() {
			return
		}

		valueFwTyped, ok := valueFw.(T)
		if !ok {
			diags.AddError("Error setting default value", fmt.Sprintf("defaultValuePlanModify: valueFw is not %T but %T", *new(T), valueFw))
		}

		*responsePlanValue = valueFwTyped
	}
}
