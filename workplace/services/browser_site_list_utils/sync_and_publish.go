// Package browsersitelistutils contains the custom logic of the microsoft365wp_browser_site_list resource (Microsoft
// Edge Internet Explorer mode site lists).
//
// MS Graph stores changes to the sites and shared cookies of a site list as pending changes (status pendingAdd,
// pendingEdit or pendingDelete) that only take effect when the list gets published. The resource does not model this
// draft stage: every create or update writes the sites and shared cookies and then publishes the whole list, so that
// the published list equals the configuration afterwards. Pending changes made elsewhere (e.g. in the Microsoft 365
// admin center) are treated as drift and get corrected before publishing, as the publish action always publishes all
// pending changes of the list.
package browsersitelistutils

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	kLogPref = "browsersitelistutils.SyncAndPublishWsa: "

	statusPendingDelete = "pendingDelete"
	statusPendingEdit   = "pendingEdit"
)

// childCollection describes a collection of child entities of a site list (sites or shared cookies).
type childCollection struct {
	attribute string // MS Graph attribute and URI suffix
	tfName    string // Terraform attribute (for messages)
	keyFunc   func(map[string]any) string
}

var childCollections = []childCollection{
	{attribute: "sites", tfName: "sites", keyFunc: siteKey},
	{attribute: "sharedCookies", tfName: "shared_cookies", keyFunc: sharedCookieKey},
}

// siteKey returns the key that MS Graph uses to detect duplicate sites: webUrl, compared case-insensitively (but e.g.
// the scheme or a trailing slash make a difference).
func siteKey(v map[string]any) string {
	webUrl, _ := v["webUrl"].(string)
	return strings.ToLower(webUrl)
}

// sharedCookieKey returns the key that MS Graph uses to detect duplicate shared cookies: displayName, hostOrDomain,
// path and hostOnly, the strings compared case-insensitively. (MS Graph rejects a duplicate with HTTP 500.)
func sharedCookieKey(v map[string]any) string {
	displayName, _ := v["displayName"].(string)
	hostOrDomain, _ := v["hostOrDomain"].(string)
	path, _ := v["path"].(string)
	hostOnly, _ := v["hostOnly"].(bool)
	return fmt.Sprintf("%s\x00%s\x00%s\x00%t", strings.ToLower(displayName), strings.ToLower(hostOrDomain), strings.ToLower(path), hostOnly)
}

// keyDescription returns a readable description of a child element's key for messages.
func keyDescription(c childCollection, v map[string]any) string {
	if c.attribute == "sites" {
		return fmt.Sprintf("web_url %q", v["webUrl"])
	}
	return fmt.Sprintf("display_name %q, host_or_domain %q, path %q, host_only %v", v["displayName"], v["hostOrDomain"], v["path"], v["hostOnly"])
}

//
// SyncAndPublishWsa
//

var _ generic.WriteSubAction = &SyncAndPublishWsa{}

// SyncAndPublishWsa writes the sites and shared cookies of a site list element by element (comparing the plan with the
// current entries in MS Graph, including pending ones) and then publishes the list.
type SyncAndPublishWsa struct {
	generic.WriteSubActionBase
}

func (a *SyncAndPublishWsa) Initialize() {
	a.Attributes = []string{"sites", "sharedCookies"}
	a.WriteSubActionBase.Initialize()
}

func (a *SyncAndPublishWsa) CheckRunAction(wsaOperation generic.OperationType) bool {
	return wsaOperation == generic.OperationCreate || wsaOperation == generic.OperationUpdate
}

func (a *SyncAndPublishWsa) ExecutePre(ctx context.Context, diags *diag.Diagnostics, wsaReq *generic.WriteSubActionRequest) {
	// remove sites and sharedCookies from the POST/PATCH of the list itself (they get written in ExecutePost)
	a.WriteSubActionBase.ExecutePre(ctx, diags, wsaReq, nil)
	if diags.HasError() {
		return
	}

	// check for duplicates before anything gets written
	for _, c := range childCollections {
		seen := make(map[string]bool)
		for _, v := range planElements(wsaReq, c) {
			key := c.keyFunc(v)
			if seen[key] {
				diags.AddError("Duplicate entry in "+c.tfName,
					fmt.Sprintf("%s contains more than one entry with %s. MS Graph compares these values case-insensitively.", c.tfName, keyDescription(c, v)))
			}
			seen[key] = true
		}
	}
}

func (a *SyncAndPublishWsa) ExecutePost(ctx context.Context, diags *diag.Diagnostics, wsaReq *generic.WriteSubActionRequest) {

	const errSumm = "Error writing browser site list"

	aps := &wsaReq.GenRes.AccessParams

	id := aps.GetId(ctx, diags, wsaReq.Id, wsaReq.IdAttributer)
	if diags.HasError() {
		return
	}
	listUri := aps.GetUriWithIdForUD(ctx, diags, "", id, wsaReq.IdAttributer)
	if diags.HasError() {
		return
	}

	for _, c := range childCollections {
		a.syncCollection(ctx, diags, aps, listUri.Entity, c, planElements(wsaReq, c))
		if diags.HasError() {
			return
		}
	}

	// publish the list (always, to also publish changes of the list itself and pending changes that already equal the
	// configuration)
	listRaw := aps.ReadRaw2(ctx, diags, listUri, "", "", []string{"id", "revision"}, false)
	if diags.HasError() {
		return
	}
	currentRevision, _ := listRaw["revision"].(string)
	revision, err := nextRevision(currentRevision)
	if err != nil {
		diags.AddError(errSumm, err.Error())
		return
	}
	tflog.Info(ctx, kLogPref+"publishing", map[string]any{"id": id, "currentRevision": currentRevision, "revision": revision})
	// only the revision, without sites or sharedCookies: this publishes all pending changes of the list
	aps.CreateRaw2(ctx, diags, listUri.Entity+"/publish", nil, map[string]any{"revision": revision}, true, false)
}

type childUpdate struct {
	id      string
	planned map[string]any
}

// syncCollection makes the (pending) entries of a collection in MS Graph equal the plan.
func (a *SyncAndPublishWsa) syncCollection(ctx context.Context, diags *diag.Diagnostics, aps *generic.AccessParams,
	listUri string, c childCollection, planned []map[string]any) {

	collectionUri := listUri + "/" + c.attribute

	currentRaw := aps.ReadRaw(ctx, diags, collectionUri, false)
	if diags.HasError() {
		return
	}
	currentItems, _ := currentRaw["value"].([]any)

	plannedByKey := make(map[string]map[string]any, len(planned))
	for _, v := range planned {
		plannedByKey[c.keyFunc(v)] = v
	}

	var toDelete []string
	var toUpdate []childUpdate
	matched := make(map[string]bool) // keys of planned elements that exist in MS Graph
	for _, item := range currentItems {
		current, ok := item.(map[string]any)
		if !ok {
			continue
		}
		currentId, _ := current["id"].(string)
		status, _ := current["status"].(string)
		key := c.keyFunc(current)
		plannedElement, isPlanned := plannedByKey[key]

		switch {
		case status == statusPendingDelete:
			if isPlanned && !matched[key] {
				// MS Graph rejects a POST with the same key (409), but a PATCH restores the entry (status pendingEdit)
				matched[key] = true
				toUpdate = append(toUpdate, childUpdate{id: currentId, planned: plannedElement})
			}
			// otherwise it will be removed by the publish
		case !isPlanned || matched[key]:
			toDelete = append(toDelete, currentId)
		default:
			matched[key] = true
			// A pending edit gets written again even if MS Graph returns the planned values, so that the publish
			// cannot publish values that have not been planned.
			if status == statusPendingEdit || !elementEquals(plannedElement, current) {
				toUpdate = append(toUpdate, childUpdate{id: currentId, planned: plannedElement})
			}
		}
	}

	var toCreate []map[string]any
	for _, v := range planned {
		if !matched[c.keyFunc(v)] {
			toCreate = append(toCreate, v)
		}
	}

	tflog.Info(ctx, kLogPref+"syncing "+c.attribute, map[string]any{"delete": len(toDelete), "update": len(toUpdate), "create": len(toCreate)})

	// Write sequentially. Delete first, so that keys get freed (e.g. for a case-only change of a webUrl that is not
	// matched because of a duplicate), then update and create.
	for _, childId := range toDelete {
		aps.DeleteRaw(ctx, diags, collectionUri, childId, nil)
		if diags.HasError() {
			return
		}
	}
	for _, u := range toUpdate {
		aps.UpdateRaw2(ctx, diags, collectionUri, u.id, nil, u.planned, false)
		if diags.HasError() {
			return
		}
	}
	for _, v := range toCreate {
		aps.CreateRaw2(ctx, diags, collectionUri, nil, v, false, false)
		if diags.HasError() {
			return
		}
	}
}

// planElements returns the planned elements of a collection as raw values (as extracted by ExecutePre).
func planElements(wsaReq *generic.WriteSubActionRequest, c childCollection) []map[string]any {
	items, _ := wsaReq.SubActionsData[c.attribute].([]any)
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if v, ok := item.(map[string]any); ok {
			result = append(result, v)
		}
	}
	return result
}

// elementEquals reports whether all planned attributes of an element have the same value in MS Graph.
func elementEquals(planned map[string]any, current map[string]any) bool {
	for k, v := range planned {
		if strings.HasPrefix(k, "@") {
			continue
		}
		if !reflect.DeepEqual(v, current[k]) {
			return false
		}
	}
	return true
}

// nextRevision returns the revision for the next publish. MS Graph only accepts a revision in the format "major" or
// "major.minor" that is greater than the last published revision (compared as a decimal number, so "10.10" is less
// than "10.5"). Therefore the next integer is used. A list that has never been published has the revision "".
func nextRevision(current string) (string, error) {
	if current == "" {
		return "1", nil
	}
	f, err := strconv.ParseFloat(current, 64)
	if err != nil || f < 0 || math.IsInf(f, 0) || math.IsNaN(f) {
		return "", fmt.Errorf("unable to determine the next revision: the current revision %q of the site list is not a number", current)
	}
	return strconv.FormatInt(int64(math.Floor(f))+1, 10), nil
}

// GraphToTerraformMiddleware removes the sites and shared cookies that are pending deletion, as they are no longer
// part of the list once it gets published (and would otherwise hide a planned element that has to be restored).
func GraphToTerraformMiddleware(ctx context.Context, diags *diag.Diagnostics, params *generic.GraphToTerraformMiddlewareParams) generic.GraphToTerraformMiddlewareReturns {
	for _, c := range childCollections {
		items, ok := params.RawVal[c.attribute].([]any)
		if !ok {
			continue
		}
		filtered := make([]any, 0, len(items))
		for _, item := range items {
			if v, ok := item.(map[string]any); ok && v["status"] == statusPendingDelete {
				continue
			}
			filtered = append(filtered, item)
		}
		params.RawVal[c.attribute] = filtered
	}
	return nil
}
