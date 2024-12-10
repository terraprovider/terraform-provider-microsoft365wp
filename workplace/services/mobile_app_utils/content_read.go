package mobileapputils

import (
	"context"
	"encoding/json"
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const kLogPrefCntR = "mobileapputils.CheckContentStateInGraphRerc: "
const kErrSummCntR = "Error reading mobile_app entity"

type contentPrivateData struct {
	ContentVersionIdCommited string
	Md5SumCommited           string
}

func (*contentPrivateData) GetPrivateKey() string {
	return "mobileApp.Content"
}

func CheckContentStateInGraphRerc(ctx context.Context, diags *diag.Diagnostics, params generic.ReadExtraRequestCustomParams) {

	tflog.Info(ctx, kLogPrefCntR+"Start")

	if params.RespPrivate == nil {
		return
	}

	tflog.Debug(ctx, kLogPrefCntR+"Get private data")
	privateBytes, diags2 := params.RespPrivate.GetKey(ctx, (*contentPrivateData).GetPrivateKey(nil))
	diags.Append(diags2...)
	if diags.HasError() || privateBytes == nil || len(privateBytes) == 0 {
		return
	}
	tflog.Trace(ctx, kLogPrefCntR+"Private data: "+string(privateBytes))
	private := contentPrivateData{}
	if err := json.Unmarshal(privateBytes, &private); err != nil {
		diags.AddWarning(kErrSummCntR, "Error decoding private JSON: "+err.Error())
		return
	}

	committedContentVersionRead, ok := params.RawVal["committedContentVersion"].(string)
	if !ok {
		diags.AddWarning(kErrSummCntR, "committedContentVersion not found or not string")
		// committedContentVersionRead should now be empty and trigger an update below
	}

	tflog.Debug(ctx, kLogPrefCntR+"Set sourceFile & sourceSha256 accordingly", map[string]any{
		"committedContentVersionRead": committedContentVersionRead,
		"ContentVersionIdCommited":    private.ContentVersionIdCommited,
	})
	if committedContentVersionRead == private.ContentVersionIdCommited {
		tflog.Debug(ctx, kLogPrefCntR+"Copy state values")
		odataType := params.RawVal["@odata.type"].(string)
		nestedPath := getNestedPathFromOdataType(ctx, odataType)
		params.RawVal["sourceFile"] = getStringAttribute(ctx, diags, params.ReqState, nestedPath.AtName("source_file"))
		params.RawVal["sourceSha256"] = getStringAttribute(ctx, diags, params.ReqState, nestedPath.AtName("source_sha256"))
		if odataType == "#microsoft.graph.windowsMobileMSI" {
			params.RawVal["manifestBase64"] = getStringAttribute(ctx, diags, params.ReqState, nestedPath.AtName("manifest_base64"))
		}
	} else {
		params.RawVal["sourceSha256"] = "(MS Graph differs from TF State)"
	}
}
