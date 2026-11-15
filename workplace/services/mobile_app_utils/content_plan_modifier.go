package mobileapputils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Plan modifiers for computed attributes that change with a new content upload (e.g. committed_content_version and
// size). They have to run after UseStateForUnknown: if a new file gets uploaded, they reset the value to unknown, so
// that the actual value can be read from MS Graph after the upload (Terraform does not accept a value that differs
// from a known planned value).

// Sibling attributes of the nested derived type whose change triggers a new content upload (see WriteContentWsa).
var contentChangeAttributes = []string{"source_sha256", "manifest_base64"}

// contentChanges reports whether the content of the nested derived type that contains attrPath gets uploaded again.
func contentChanges(ctx context.Context, attrPath path.Path, plan tfsdk.Plan, state tfsdk.State) bool {

	const kLogPref = "mobileapputils.contentChanges: "

	if state.Raw.IsNull() || plan.Raw.IsNull() {
		// create (value is unknown anyway) or destroy
		return false
	}

	parentPath := attrPath.ParentPath()
	for _, name := range contentChangeAttributes {
		var planValue, stateValue types.String
		if plan.GetAttribute(ctx, parentPath.AtName(name), &planValue).HasError() {
			continue // attribute does not exist in this nested type (e.g. manifest_base64)
		}
		if state.GetAttribute(ctx, parentPath.AtName(name), &stateValue).HasError() {
			continue
		}
		if planValue.IsUnknown() || !planValue.Equal(stateValue) {
			tflog.Debug(ctx, kLogPref+"content changes", map[string]any{"path": attrPath.String(), "attribute": name})
			return true
		}
	}
	return false
}

type unknownOnContentChange struct{}

func (m unknownOnContentChange) Description(ctx context.Context) string {
	return "The value is unknown until after apply if a new file gets uploaded."
}

func (m unknownOnContentChange) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

// StringUnknownOnContentChange sets the planned value to unknown if a new file gets uploaded.
func StringUnknownOnContentChange() planmodifier.String {
	return unknownOnContentChange{}
}

func (m unknownOnContentChange) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if !resp.PlanValue.IsUnknown() && contentChanges(ctx, req.Path, req.Plan, req.State) {
		resp.PlanValue = types.StringUnknown()
	}
}

// Int64UnknownOnContentChange sets the planned value to unknown if a new file gets uploaded.
func Int64UnknownOnContentChange() planmodifier.Int64 {
	return unknownOnContentChange{}
}

func (m unknownOnContentChange) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if !resp.PlanValue.IsUnknown() && contentChanges(ctx, req.Path, req.Plan, req.State) {
		resp.PlanValue = types.Int64Unknown()
	}
}
