package mobileapputils

import (
	"context"
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func getNestedPathFromOdataType(ctx context.Context, odataType string) path.Path {

	const kLogPref = "mobileapputils.getNestedPathFromOdataType: "

	tflog.Debug(ctx, kLogPref+"odataType: "+odataType)

	nestedPath := path.Empty()
	switch odataType {
	case "#microsoft.graph.macOSDmgApp":
		nestedPath = path.Root("macos_dmg")
	case "#microsoft.graph.macOSPkgApp":
		nestedPath = path.Root("macos_pkg")
	case "#microsoft.graph.win32LobApp":
		nestedPath = path.Root("win32_lob")
	case "#microsoft.graph.windowsMobileMSI":
		nestedPath = path.Root("windows_msi")
	case "#microsoft.graph.windowsUniversalAppX":
		nestedPath = path.Root("windows_universal_appx")
	default:
		// nothing to do
	}

	tflog.Trace(ctx, kLogPref+"nestedPath: "+nestedPath.String())
	return nestedPath
}

func getStringAttribute(ctx context.Context, diags *diag.Diagnostics, source generic.GetAttributer, path path.Path) string {

	const kLogPref = "mobileapputils.getStringAttribute: "

	tflog.Debug(ctx, kLogPref+"path: "+path.String())

	var resultTypes types.String // might be null
	diags.Append(source.GetAttribute(ctx, path, &resultTypes)...)
	if diags.HasError() {
		return ""
	}
	result := resultTypes.ValueString() // nil will become ""

	tflog.Trace(ctx, kLogPref+"result: "+result)
	return result
}
