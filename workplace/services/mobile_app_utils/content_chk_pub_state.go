package mobileapputils

import (
	"context"
	"fmt"
	"terraform-provider-microsoft365wp/workplace/generic"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const kLogPrefChkPubStPost = "mobileapputils.CheckPublishingStateWsa Post: "

//
// CheckPublishingStateWsa
//

var _ generic.WriteSubAction = &CheckPublishingStateWsa{}

type CheckPublishingStateWsa struct {
	generic.WriteSubActionBase
	MessageSuffix     string
	AllowNotPublished bool
}

func (a *CheckPublishingStateWsa) Initialize() {
	a.WriteSubActionBase.Initialize()
}

func (a *CheckPublishingStateWsa) CheckRunAction(wsaOperation generic.OperationType) bool {
	return a.WriteSubActionBase.CheckRunAction(wsaOperation)
}

func (*CheckPublishingStateWsa) ExecutePre(ctx context.Context, diags *diag.Diagnostics, wsaReq *generic.WriteSubActionRequest) {
	// nothing to do here
}

func (wsa *CheckPublishingStateWsa) ExecutePost(ctx context.Context, diags *diag.Diagnostics, wsaReq *generic.WriteSubActionRequest) {

	errSummChkPubSt := fmt.Sprintf("Error checking mobile_app publish_state (%s)", wsa.MessageSuffix)

	tflog.Info(ctx, kLogPrefChkPubStPost+"Start", map[string]any{"wsaReq.Operation": wsaReq.Operation})

	tflog.Debug(ctx, kLogPrefCntWPost+"Get entityUri")
	entityUriUri, _ := wsaReq.GenRes.AccessParams.GetUriWithIdForR(ctx, diags, "", wsaReq.Id, wsaReq.IdAttributer, false, "")
	if diags.HasError() {
		return
	}
	tflog.Trace(ctx, kLogPrefCntWPost+"entityUri: "+entityUriUri.Entity)

	for {
		// include id for easier debugging
		body := wsaReq.GenRes.AccessParams.ReadRaw(ctx, diags, entityUriUri, "", "", []string{"id", "publishingState"}, false)
		if diags.HasError() {
			return
		}
		publishingState := body["publishingState"].(string)
		tflog.Trace(ctx, kLogPrefChkPubStPost+"publishingState: "+publishingState)

		if publishingState != "processing" {
			if !wsa.AllowNotPublished && publishingState != "published" {
				diags.AddError(errSummChkPubSt, "publishing_state is not 'published' but: "+publishingState)
			}
			return
		}

		time.Sleep(time.Millisecond * 200)
	}
}
