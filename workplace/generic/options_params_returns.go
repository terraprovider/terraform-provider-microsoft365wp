package generic

import (
	"context"
	"sync"
	"terraform-provider-microsoft365wp/workplace/strcase"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/manicminer/hamilton/msgraph"
)

type AccessParams struct {
	BaseUri       string
	HasParentItem HasParentItem
	IsSingleton   bool

	ReadOptions                ReadOptions
	GraphToTerraformMiddleware func(GraphToTerraformMiddlewareParams) GraphToTerraformMiddlewareReturns

	TerraformToGraphMiddleware func(TerraformToGraphMiddlewareParams) TerraformToGraphMiddlewareReturns
	CreateModifyFunc           func(*CreateModifyFuncParams) // pass pointer to allow modifications
	UpdateModifyFunc           func(*UpdateModifyFuncParams)
	DeleteModifyFunc           func(*DeleteModifyFuncParams)
	WriteSubActions            []WriteSubAction
	UsePutForUpdate            bool
	SerializeWrites            bool
	SerialWritesDelay          time.Duration

	initializeOnce sync.Once

	graphClient *msgraph.Client
}

type HasParentItem struct {
	ParentIdField path.Path
	UriSuffix     string
}

type ReadOptions struct {
	ODataExpand              string
	ODataFilter              string // GenericDataSourcePlural and SingleItemUseODataFilter only, not allowed otherwise
	SingleItemUseODataFilter bool
	ExtraRequests            []ReadExtraRequest
	ExtraRequestsCustom      []ReadExtraRequestCustom
	PluralNoFilterSupport    bool
}

type ReadExtraRequest struct {
	ParentAttribute string
	UriSuffix       string
	ODataExpand     string
}

type ReadExtraRequestCustomFuncParams struct {
	Ctx              context.Context
	Diags            *diag.Diagnostics
	Client           *msgraph.Client
	Uri              msgraph.Uri
	TolerateNotFound bool
	RawVal           map[string]any
}
type ReadExtraRequestCustom func(ReadExtraRequestCustomFuncParams)

type GraphToTerraformMiddlewareParams struct {
	Ctx    context.Context
	RawVal map[string]any
}
type GraphToTerraformMiddlewareReturns error

type TerraformToGraphMiddlewareParams struct {
	Ctx      context.Context
	RawVal   map[string]any
	IsUpdate bool
}
type TerraformToGraphMiddlewareReturns error

type CreateModifyFuncParams struct {
	R              *GenericResource
	Ctx            context.Context
	Diags          *diag.Diagnostics
	Req            resource.CreateRequest
	Resp           *resource.CreateResponse
	BaseUri        string
	Id             string
	UpdateExisting bool
}

type UpdateModifyFuncParams struct {
	R       *GenericResource
	Ctx     context.Context
	Diags   *diag.Diagnostics
	Req     resource.UpdateRequest
	Resp    *resource.UpdateResponse
	BaseUri string
	Id      string
}

type DeleteModifyFuncParams struct {
	R          *GenericResource
	Ctx        context.Context
	Diags      *diag.Diagnostics
	Req        resource.DeleteRequest
	Resp       *resource.DeleteResponse
	BaseUri    string
	Id         string
	SkipDelete bool
}

func (ap *AccessParams) InitializeGuarded(providerData any) {

	// providerData may be empty on early calls, so we always check this here
	if ap.graphClient == nil {
		if graphClient, ok := providerData.(*msgraph.Client); ok {
			ap.graphClient = graphClient
		}
	}

	ap.initializeOnce.Do(func() {

		for _, wsa := range ap.WriteSubActions {
			wsa.Initialize()
		}

		if !ap.HasParentItem.ParentIdField.Equal(path.Path{}) && ap.TerraformToGraphMiddleware == nil {
			parentIdFieldCamel := strcase.ToLowerCamel(ap.HasParentItem.ParentIdField.String())
			ap.TerraformToGraphMiddleware = func(params TerraformToGraphMiddlewareParams) TerraformToGraphMiddlewareReturns {
				delete(params.RawVal, parentIdFieldCamel)
				return nil
			}
		}

	})
}
