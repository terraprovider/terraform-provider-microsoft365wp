package generic

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/manicminer/hamilton/msgraph"
)

//
// GetAttributer
//

type GetAttributer interface {
	GetAttribute(ctx context.Context, path path.Path, target any) diag.Diagnostics
}

//
// OperationType
//

type OperationType int

const (
	OperationCreate OperationType = iota
	OperationUpdate
	OperationDelete
)

//
// Get id
//

func (aps *AccessParams) GetId(ctx context.Context, diags *diag.Diagnostics, id string, idAttributer GetAttributer) string {

	if id == "" {
		diags.Append(idAttributer.GetAttribute(ctx, path.Root("id"), &id)...)
		if diags.HasError() {
			return ""
		}
	}

	return id
}

//
// Get*Uri*
//

func (aps *AccessParams) GetBaseUri(ctx context.Context, diags *diag.Diagnostics, baseUri string, parentIdAttributer GetAttributer) msgraph.Uri {
	uri, _ := aps.getUriInternal(ctx, diags, baseUri, false, "", parentIdAttributer, false, "")
	return uri
}

func (aps *AccessParams) GetUriWithIdForR(ctx context.Context, diags *diag.Diagnostics, baseUri string,
	id string, idAttributer GetAttributer, useODataFilterToSelectById bool, baseODataFilter string) (msgraph.Uri, string) {
	return aps.getUriInternal(ctx, diags, baseUri, true, id, idAttributer, useODataFilterToSelectById, baseODataFilter)
}

func (aps *AccessParams) GetUriWithIdForUD(ctx context.Context, diags *diag.Diagnostics, baseUri string,
	id string, idAttributer GetAttributer) msgraph.Uri {
	uri, _ := aps.getUriInternal(ctx, diags, baseUri, true, id, idAttributer, false, "")
	return uri
}

func (aps *AccessParams) getUriInternal(ctx context.Context, diags *diag.Diagnostics, baseUri string,
	withIdAppended bool, id string, idAttributer GetAttributer,
	useODataFilterToSelectById bool, baseODataFilter string) (msgraph.Uri, string) {

	if baseUri == "" {
		baseUri = aps.BaseUri
		if !aps.HasParentItem.ParentIdField.Equal(path.Path{}) {
			var parentId string
			diags.Append(idAttributer.GetAttribute(ctx, aps.HasParentItem.ParentIdField, &parentId)...)
			if diags.HasError() {
				return msgraph.Uri{}, ""
			}
			baseUri += fmt.Sprintf("/%s", parentId)
			baseUri += aps.HasParentItem.UriSuffix
		}
	}

	if !withIdAppended || aps.IsSingleton {
		return msgraph.Uri{Entity: baseUri}, ""
	}

	id = aps.GetId(ctx, diags, id, idAttributer)
	if diags.HasError() {
		return msgraph.Uri{}, ""
	}

	// baseODataFilter will only be returned for SingleItemUseODataFilter, also see ReadOptions.ODataFilter
	if useODataFilterToSelectById {
		odataFilter := fmt.Sprintf("id eq '%s'", id)
		if baseODataFilter != "" {
			odataFilter = fmt.Sprintf("(%s) and (%s)", baseODataFilter, odataFilter)
		}
		return msgraph.Uri{Entity: baseUri}, odataFilter
	} else {
		uriEntity := fmt.Sprintf("%s/%s", baseUri, id)
		return msgraph.Uri{Entity: uriEntity}, ""
	}
}

//
// EmptyDescriber
//

type EmptyDescriber struct{}

func (apm EmptyDescriber) Description(ctx context.Context) string         { return "" }
func (apm EmptyDescriber) MarkdownDescription(ctx context.Context) string { return "" }
