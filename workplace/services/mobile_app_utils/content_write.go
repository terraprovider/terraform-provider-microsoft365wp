package mobileapputils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/util/azurestorage"
	"terraform-provider-microsoft365wp/workplace/util/mobileappcontent"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const kErrSummCntW = "Error uploading mobile_app content"

var kDebugContent = os.Getenv("TF_M365WP_MOBILE_APP_DEBUG_CONTENT") == "1"

//
// WriteContentWsa
//

var _ generic.WriteSubAction = &WriteContentWsa{}

type WriteContentWsa struct {
	generic.WriteSubActionBase
	writeMutex sync.Mutex
}

func (a *WriteContentWsa) Initialize() {
	a.WriteSubActionBase.Initialize()
}

func (a *WriteContentWsa) CheckRunAction(wsaOperation generic.OperationType) bool {
	return a.WriteSubActionBase.CheckRunAction(wsaOperation)
}

//
// ExecutePre
//

func (*WriteContentWsa) ExecutePre(ctx context.Context, diags *diag.Diagnostics, wsaReq *generic.WriteSubActionRequest) {

	const kLogPrefCntWPre = "mobileapputils.WriteContentWsa Pre: "

	tflog.Info(ctx, kLogPrefCntWPre+"Start", map[string]any{"wsaReq.Operation": wsaReq.Operation})

	tflog.Debug(ctx, kLogPrefCntWPre+"Get and check @odata.type")
	odataType := wsaReq.RawVal["@odata.type"].(string)
	nestedPath := getNestedPathFromOdataType(ctx, odataType)
	if nestedPath.Equal(path.Empty()) {
		// nothing to do
		return
	}

	// artifical attributes that may not be pushed to MS Graph
	delete(wsaReq.RawVal, "sourceFile")
	delete(wsaReq.RawVal, "sourceSha256")
	delete(wsaReq.RawVal, "manifestBase64")

	tflog.Debug(ctx, kLogPrefCntWPre+"Get plan values")
	sourceFile := getStringAttribute(ctx, diags, wsaReq.ReqPlan, nestedPath.AtName("source_file"))
	sourceSha256Plan := getStringAttribute(ctx, diags, wsaReq.ReqPlan, nestedPath.AtName("source_sha256"))
	var manifestPlan string
	if odataType == "#microsoft.graph.windowsMobileMSI" {
		manifestPlan = getStringAttribute(ctx, diags, wsaReq.ReqPlan, nestedPath.AtName("manifest_base64"))
	}
	if diags.HasError() {
		return
	}

	if wsaReq.Operation == generic.OperationUpdate {
		tflog.Debug(ctx, kLogPrefCntWPre+"Get state source_sha256 and compare to plan")
		sourceSha256State := getStringAttribute(ctx, diags, wsaReq.ReqState, nestedPath.AtName("source_sha256"))
		var manifestState string
		if odataType == "#microsoft.graph.windowsMobileMSI" {
			manifestState = getStringAttribute(ctx, diags, wsaReq.ReqState, nestedPath.AtName("manifest_base64"))
		}
		if sourceSha256State == sourceSha256Plan && manifestState == manifestPlan {
			tflog.Info(ctx, kLogPrefCntWPre+"source_sha256 (and manifest_base64 for windows_msi) are identical, no need to update mobileApp content")
			return
		}
	}

	tflog.Info(ctx, kLogPrefCntWPre+"Not an update or source_sha256 differ: Prepare to update mobileApp content")

	sad := writeContentWsaData{OdataType: odataType}
	var err error
	switch odataType {

	case "#microsoft.graph.win32LobApp":

		tflog.Debug(ctx, kLogPrefCntWPre+"Get IntuneWin encryption info and encrypted content", map[string]any{"sourceFile": sourceFile})
		var itwMetadata mobileappcontent.IntunewinMetadata
		itwMetadata, sad.FileEncryptionInfo, sad.FileNameEncrypted, err = mobileappcontent.ParseAndExtractIntunewinOuter(ctx, sourceFile)
		if err != nil {
			diags.AddError(kErrSummCntW, "Error reading IntuneWin file: "+err.Error())
			return
		}

		sad.ContentFileInfo.Name = itwMetadata.FileName
		sad.ContentFileInfo.Size = itwMetadata.UnencryptedContentSize
		sad.ContentFileInfo.SizeEncrypted = itwMetadata.EncryptedContentSize
		// Manifest is not required for IntuneWin

	default:

		tflog.Debug(ctx, kLogPrefCntWPre+"Stat source file")
		sourceFileInfo, err := os.Stat(sourceFile)
		if err != nil {
			diags.AddError(kErrSummCntW, "Error getting source file info: "+err.Error())
			return
		}
		sourceFileSize := sourceFileInfo.Size()
		tflog.Trace(ctx, kLogPrefCntWPre+fmt.Sprintf("sourceFileSize: %d", sourceFileSize))

		tflog.Debug(ctx, kLogPrefCntWPre+"Encrypt source file", map[string]any{"sourceFile": sourceFile})
		var fileSizeEnc int64
		sad.FileEncryptionInfo, sad.FileNameEncrypted, fileSizeEnc, err = mobileappcontent.EnryptFile(sourceFile)
		if err != nil {
			diags.AddError(kErrSummCntW, "Error encrypting source file: "+err.Error())
			return
		}

		sad.ContentFileInfo.Name = filepath.Base(sourceFile)
		sad.ContentFileInfo.Size = sourceFileSize
		sad.ContentFileInfo.SizeEncrypted = fileSizeEnc

		switch odataType {

		case "#microsoft.graph.windowsMobileMSI":
			sad.ContentFileInfo.Manifest = manifestPlan
			if diags.HasError() {
				return
			}

		case "#microsoft.graph.windowsUniversalAppX":
			sad.ContentFileInfo.Manifest = base64.StdEncoding.EncodeToString([]byte(sad.ContentFileInfo.Name))

		}

	}

	// logging only, use own scope
	{
		subActionsDataJson, _ := json.Marshal(sad)
		logMessage := kLogPrefCntWPre + "subActionsData: " + string(subActionsDataJson)
		if kDebugContent {
			tflog.Warn(ctx, logMessage)
		} else {
			tflog.Trace(ctx, logMessage)
		}
	}

	wsaReq.SubActionsData[(*writeContentWsaData).GetWsadKey(nil)] = sad
}

//
// ExecutePost
//

const kLogPrefCntWPost = "mobileapputils.WriteContentWsa Post: "

func (wsa *WriteContentWsa) ExecutePost(ctx context.Context, diags *diag.Diagnostics, wsaReq *generic.WriteSubActionRequest) {

	tflog.Info(ctx, kLogPrefCntWPost+"Start", map[string]any{"wsaReq.Operation": wsaReq.Operation})

	sadKey := (*writeContentWsaData).GetWsadKey(nil)
	tflog.Debug(ctx, kLogPrefCntWPost+fmt.Sprintf("Check existence of wsaReq.SubActionsData[\"%s\"]", sadKey))
	sadAny, ok := wsaReq.SubActionsData[sadKey]
	if !ok {
		tflog.Info(ctx, kLogPrefCntWPost+"No need to update mobileApp content")
		return // nothing to do for us
	}
	sad := sadAny.(writeContentWsaData)

	// serialize content updates for now due to problems with concurring content files creation
	if !wsa.writeMutex.TryLock() {
		tflog.Info(ctx, kLogPrefCntWPost+"Wait for serialization mutex")
		wsa.writeMutex.Lock()
		tflog.Info(ctx, kLogPrefCntWPost+"Lock of serialization mutex succeeded, delay a bit (just in case)")
		time.Sleep(time.Second * 3)
	}
	defer wsa.writeMutex.Unlock()

	tflog.Info(ctx, kLogPrefCntWPost+"Execute update of mobileApp content")

	if !kDebugContent {
		defer os.Remove(sad.FileNameEncrypted) // this has been created in Pre
	}

	tflog.Debug(ctx, kLogPrefCntWPost+"Get entityUri")
	entityUriUri := wsaReq.GenRes.AccessParams.GetUriWithIdForUD(ctx, diags, "", wsaReq.Id, wsaReq.IdAttributer)
	if diags.HasError() {
		return
	}
	entityUri := entityUriUri.Entity
	tflog.Trace(ctx, kLogPrefCntWPost+"entityUri: "+entityUri)

	// (hopefully) avoid Precondition Failed / ConditionNotMet errors
	tflog.Debug(ctx, kLogPrefCntWPost+"Delay before create contentVersion")
	time.Sleep(time.Second * 1)

	tflog.Debug(ctx, kLogPrefCntWPost+"Create contentVersion")
	contentVersionId, versionUri := wsa.createContentVersion(ctx, diags, wsaReq.GenRes, entityUri, sad.OdataType)
	if diags.HasError() {
		return
	}
	tflog.Trace(ctx, kLogPrefCntWPost+"contentVersionId: "+contentVersionId)
	tflog.Trace(ctx, kLogPrefCntWPost+"versionUri: "+versionUri)

	tflog.Debug(ctx, kLogPrefCntWPost+"Create contentFile")
	_, fileUri := wsa.createContentFile(ctx, diags, wsaReq.GenRes, versionUri, sad.ContentFileInfo)
	if diags.HasError() {
		return
	}
	tflog.Trace(ctx, kLogPrefCntWPost+"fileUri: "+fileUri)

	tflog.Debug(ctx, kLogPrefCntWPost+"Wait for azureStorageUriRequestSuccess")
	azureStorageUri := wsa.ensureContentFileUploadStateSuccess(ctx, diags, wsaReq.GenRes, fileUri, "azureStorageUriRequest", "azureStorageUri")
	if diags.HasError() {
		return
	}
	tflog.Trace(ctx, kLogPrefCntWPost+"azureStorageUri: "+azureStorageUri)

	tflog.Debug(ctx, kLogPrefCntWPost+"Upload IntuneWin to Azure Storage")
	fileMd5Sum, err := azurestorage.UploadBlobInBlocks(ctx, sad.FileNameEncrypted, azureStorageUri)
	if err != nil {
		diags.AddError(kErrSummCntW, "Error uploading IntuneWin file to Azure Storage: "+err.Error())
		return
	}
	tflog.Trace(ctx, kLogPrefCntWPost+"fileMd5Sum: "+fileMd5Sum)

	tflog.Debug(ctx, kLogPrefCntWPost+"Commit contentFile")
	wsa.commitContentFile(ctx, diags, wsaReq.GenRes, fileUri, sad.FileEncryptionInfo)
	if diags.HasError() {
		return
	}

	tflog.Debug(ctx, kLogPrefCntWPost+"Wait for commitFileSuccess")
	wsa.ensureContentFileUploadStateSuccess(ctx, diags, wsaReq.GenRes, fileUri, "commitFile", "")
	if diags.HasError() {
		return
	}

	tflog.Debug(ctx, kLogPrefCntWPost+"Patch mobileApp with new data")
	wsa.patchMobileApp(ctx, diags, wsaReq.GenRes, wsaReq.Id, wsaReq.IdAttributer, sad.OdataType, contentVersionId)
	if diags.HasError() {
		return
	}

	// TF will not accept any changes on attributes with UseStateForUnknown on _modify_ (only on read), so we cannot set committed_content_version here.

	tflog.Debug(ctx, kLogPrefCntWPost+"Set private data")
	privateData := contentPrivateData{
		ContentVersionIdCommited: contentVersionId,
		Md5SumCommited:           fileMd5Sum,
	}
	privateDataJson, err := json.Marshal(privateData)
	if err != nil {
		diags.AddError(kErrSummCntW, "Error encoding privateData to JSON: "+err.Error())
		return
	}
	tflog.Trace(ctx, kLogPrefCntWPost+"Private data: "+string(privateDataJson))
	diags.Append(wsaReq.RespPrivate.SetKey(ctx, privateData.GetPrivateKey(), privateDataJson)...)
	if diags.HasError() {
		return
	}
}

func (*WriteContentWsa) createContentVersion(ctx context.Context, diags *diag.Diagnostics, r *generic.GenericResource,
	mobileAppUri string, odatatype string) (string, string) {

	baseUri := fmt.Sprintf("%s/%s/contentVersions", mobileAppUri, strings.TrimPrefix(odatatype, "#"))
	versionId, _ := r.AccessParams.CreateRaw(ctx, diags, baseUri, nil, map[string]any{}, false)
	versionUri := fmt.Sprintf("%s/%s", baseUri, versionId)

	return versionId, versionUri
}

func (*WriteContentWsa) createContentFile(ctx context.Context, diags *diag.Diagnostics, r *generic.GenericResource,
	versionUri string, contentFileInfo mobileappcontent.ContentFileInfo) (string, string) {

	baseUri := fmt.Sprintf("%s/files", versionUri)
	body := map[string]any{
		"@odata.type":   "#microsoft.graph.mobileAppContentFile",
		"name":          contentFileInfo.Name,
		"size":          contentFileInfo.Size,
		"sizeEncrypted": contentFileInfo.SizeEncrypted,
		"manifest":      contentFileInfo.Manifest,
		"isDependency":  false,
	}
	fileId, _ := r.AccessParams.CreateRaw(ctx, diags, baseUri, nil, body, false)
	fileUri := fmt.Sprintf("%s/%s", baseUri, fileId)

	return fileId, fileUri
}

func (*WriteContentWsa) ensureContentFileUploadStateSuccess(ctx context.Context, diags *diag.Diagnostics, r *generic.GenericResource,
	fileUri string, statePrefix string, returnAttributeName string) (returnValue string) {

	for {
		fileBody := r.AccessParams.ReadRaw(ctx, diags, fileUri, false)
		if diags.HasError() {
			return
		}

		uploadState := fileBody["uploadState"].(string)
		tflog.Trace(ctx, kLogPrefCntWPost+"uploadState: "+uploadState)
		if uploadState == statePrefix+"Success" {
			if returnAttributeName != "" {
				returnValue = fileBody[returnAttributeName].(string)
			}
			return
		}
		if uploadState != statePrefix+"Pending" {
			diags.AddError(kErrSummCntW, "mobileAppContentFile.uploadState does not indicate success: "+uploadState)
			return
		}

		time.Sleep(time.Millisecond * 500)
	}
}

func (*WriteContentWsa) commitContentFile(ctx context.Context, diags *diag.Diagnostics, r *generic.GenericResource,
	fileUri string, fileEncryptionInfo mobileappcontent.FileEncryptionInfo) {

	baseUri := fmt.Sprintf("%s/commit", fileUri)

	// hacky, but simplest way I could find
	body := map[string]any{"fileEncryptionInfo": fileEncryptionInfo}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		diags.AddError(kErrSummCntW, "Error encoding fileEncryptionInfo to JSON: "+err.Error())
		return
	}
	var bodyRaw map[string]any
	if err := json.Unmarshal(bodyJson, &bodyRaw); err != nil {
		diags.AddError(kErrSummCntW, "Error decoding fileEncryptionInfo JSON: "+err.Error())
		return
	}

	r.AccessParams.CreateRaw(ctx, diags, baseUri, nil, bodyRaw, true)
}

func (*WriteContentWsa) patchMobileApp(ctx context.Context, diags *diag.Diagnostics, r *generic.GenericResource,
	mobileAppId string, mobileAppIdAttributer generic.GetAttributer, odatatype string, contentVersionId string) {

	body := map[string]any{
		"@odata.type":             odatatype,
		"committedContentVersion": contentVersionId,
	}
	r.AccessParams.UpdateRaw(ctx, diags, r.AccessParams.BaseUri, mobileAppId, mobileAppIdAttributer, body, false)
}

//
// writeContentWsaData
//

type writeContentWsaData struct {
	OdataType          string
	ContentFileInfo    mobileappcontent.ContentFileInfo
	FileEncryptionInfo mobileappcontent.FileEncryptionInfo
	FileNameEncrypted  string
}

func (*writeContentWsaData) GetWsadKey() string {
	return ".mobileApp.WriteContent"
}
