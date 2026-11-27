package browsersitelistutils

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ConfigValidator = ConfigValidator{}

// ConfigValidator checks the sites and shared cookies of a site list at plan time (as far as their values are known):
//   - no duplicate keys (MS Graph rejects them with 409 for sites and 500 for shared cookies)
//   - host_or_domain of a shared cookie starts with a dot if and only if host_only is false (MS Graph adds the dot to a
//     domain cookie, which would cause a perpetual diff, and rejects it for a host-only cookie)
type ConfigValidator struct{}

func (v ConfigValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (ConfigValidator) MarkdownDescription(context.Context) string {
	return "Checks that sites and shared cookies are unique and that host_or_domain of a shared cookie starts with a dot if and only if host_only is false"
}

func (ConfigValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	for _, c := range childCollections {
		var set types.Set
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root(c.tfName), &set)...)
		if resp.Diagnostics.HasError() || set.IsNull() || set.IsUnknown() {
			continue
		}

		seen := make(map[string]bool)
		for _, element := range set.Elements() {
			obj, ok := element.(types.Object)
			if !ok || obj.IsNull() || obj.IsUnknown() {
				continue
			}
			raw, complete := objectToKeyRaw(obj.Attributes())

			if c.attribute == "sharedCookies" {
				validateHostOrDomain(resp, raw)
			}

			if complete {
				key := c.keyFunc(raw)
				if seen[key] {
					resp.Diagnostics.AddAttributeError(path.Root(c.tfName), "Duplicate entry in "+c.tfName,
						fmt.Sprintf("%s contains more than one entry with %s. MS Graph compares these values case-insensitively.", c.tfName, keyDescription(c, raw)))
				}
				seen[key] = true
			}
		}
	}
}

func validateHostOrDomain(resp *resource.ValidateConfigResponse, raw map[string]any) {
	hostOrDomain, ok := raw["hostOrDomain"].(string)
	if !ok {
		return
	}
	hostOnly, ok := raw["hostOnly"].(bool)
	if !ok {
		return
	}
	hasDot := strings.HasPrefix(hostOrDomain, ".")
	switch {
	case hostOnly && hasDot:
		resp.Diagnostics.AddAttributeError(path.Root("shared_cookies"), "Invalid host_or_domain of a shared cookie",
			fmt.Sprintf("host_or_domain %q of the host-only cookie %q must not start with a dot (MS Graph rejects it).", hostOrDomain, raw["displayName"]))
	case !hostOnly && !hasDot:
		resp.Diagnostics.AddAttributeError(path.Root("shared_cookies"), "Invalid host_or_domain of a shared cookie",
			fmt.Sprintf("host_or_domain %q of the domain cookie %q (host_only = false) must start with a dot, e.g. %q (MS Graph adds it otherwise), or set host_only = true.",
				hostOrDomain, raw["displayName"], "."+hostOrDomain))
	}
}

// objectToKeyRaw converts the attributes needed for the keys (and for validateHostOrDomain) to raw MS Graph values.
// complete is false if any of the key attributes is unknown. A null host_only counts as false (its default).
func objectToKeyRaw(attrs map[string]attr.Value) (raw map[string]any, complete bool) {
	raw = make(map[string]any)
	complete = true
	stringAttrs := map[string]string{"web_url": "webUrl", "display_name": "displayName", "host_or_domain": "hostOrDomain", "path": "path"}
	for tfName, graphName := range stringAttrs {
		v, ok := attrs[tfName].(types.String)
		if !ok {
			continue
		}
		if v.IsUnknown() {
			complete = false
			continue
		}
		raw[graphName] = v.ValueString()
	}
	if v, ok := attrs["host_only"].(types.Bool); ok {
		if v.IsUnknown() {
			complete = false
		} else {
			raw["hostOnly"] = v.ValueBool()
		}
	}
	return raw, complete
}
