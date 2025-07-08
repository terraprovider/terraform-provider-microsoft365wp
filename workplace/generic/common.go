package generic

import (
	"context"
	"fmt"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//
// GetAttributer
//

type GetAttributer interface {
	GetAttribute(ctx context.Context, path path.Path, target any) diag.Diagnostics
}

//
// GetDescriptioner
//

type GetDescriptioner interface {
	GetDescription() string
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
		if idAttributer == nil {
			diags.AddError("GetId failed", "id is empty and idAttributer is nil")
			return ""
		}

		var idAttribute types.String // might be null
		diags.Append(idAttributer.GetAttribute(ctx, path.Root(aps.EntityId.AttrNameTf), &idAttribute)...)
		if diags.HasError() {
			return ""
		}
		id = idAttribute.ValueString() // nil will become ""

		if id == "" && aps.EntityId.FallbackIdGetterFunc != nil {
			id = aps.EntityId.FallbackIdGetterFunc(ctx, diags, idAttributer)
			if diags.HasError() {
				return ""
			}
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
		for _, pe := range aps.ParentEntities {
			var parentId string
			if idAttributer == nil {
				diags.AddError("GetUri... failed", "ParentEntities exist but idAttributer is nil")
				return msgraph.Uri{}, ""
			}
			diags.Append(idAttributer.GetAttribute(ctx, pe.ParentIdField, &parentId)...)
			if diags.HasError() {
				return msgraph.Uri{}, ""
			}
			if mappedParentIdValue, ok := pe.ParentIdValueMap[parentId]; ok {
				parentId = mappedParentIdValue
			}
			baseUri += fmt.Sprintf("/%s", parentId)
			if pe.UriSuffix != "" {
				baseUri += fmt.Sprintf("/%s", pe.UriSuffix)
			}
		}
	}

	if !withIdAppended || aps.UriNoId {
		return msgraph.Uri{Entity: baseUri}, baseODataFilter
	}

	id = aps.GetId(ctx, diags, id, idAttributer)
	if diags.HasError() {
		return msgraph.Uri{}, ""
	}
	if id == "" {
		if baseODataFilter != "" {
			// if selecting an entity by OData filter then id will/must be empty and we are done here
			return msgraph.Uri{Entity: baseUri}, baseODataFilter
		}
		diags.AddError("Unable to find id attribute value", "")
		return msgraph.Uri{}, ""
	}

	// baseODataFilter will only be returned for SingleItemUseODataFilter, also see ReadOptions.ODataFilter
	if useODataFilterToSelectById {
		odataFilter := fmt.Sprintf("%s eq '%s'", aps.EntityId.AttrNameGraph, id)
		if baseODataFilter != "" {
			odataFilter = fmt.Sprintf("(%s) and (%s)", baseODataFilter, odataFilter)
		}
		return msgraph.Uri{Entity: baseUri}, odataFilter
	} else {
		uriEntity := fmt.Sprintf("%s/%s", baseUri, id)
		if aps.UriSuffix != "" {
			uriEntity += fmt.Sprintf("/%s", aps.UriSuffix)
		}
		return msgraph.Uri{Entity: uriEntity}, ""
	}
}

//
// EmptyDescriber
//

type EmptyDescriber struct{}

func (apm EmptyDescriber) Description(ctx context.Context) string         { return "" }
func (apm EmptyDescriber) MarkdownDescription(ctx context.Context) string { return "" }
