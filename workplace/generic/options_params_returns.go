package generic

import (
	"context"
	"sync"
	"terraform-provider-microsoft365wp/workplace/strcase"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/manicminer/hamilton/msgraph"
)

type AccessParams struct {
	BaseUri       string
	HasParentItem HasParentItem
	IsSingleton   bool

	ReadOptions                ReadOptions
	GraphToTerraformMiddleware func(context.Context, GraphToTerraformMiddlewareParams) GraphToTerraformMiddlewareReturns

	TerraformToGraphMiddleware func(context.Context, TerraformToGraphMiddlewareParams) TerraformToGraphMiddlewareReturns
	CreateModifyFunc           func(context.Context, *diag.Diagnostics, *CreateModifyFuncParams) // pass pointer to allow modifications
	UpdateModifyFunc           func(context.Context, *diag.Diagnostics, *UpdateModifyFuncParams)
	DeleteModifyFunc           func(context.Context, *diag.Diagnostics, *DeleteModifyFuncParams)
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
	ValidStatusCodesExtra    []int
	ExtraRequests            []ReadExtraRequest
	ExtraRequestsCustom      []ReadExtraRequestCustom
	PluralNoFilterSupport    bool
	PluralNoSelectSupport    bool
}

type ReadExtraRequest struct {
	ParentAttribute string
	UriSuffix       string
	ODataExpand     string
}

type ReadExtraRequestCustomParams struct {
	Client           *msgraph.Client
	Uri              msgraph.Uri
	TolerateNotFound bool
	ReqState         *tfsdk.State
	RespPrivate      PrivateDataGetSetter
	RawVal           map[string]any
}
type ReadExtraRequestCustom func(context.Context, *diag.Diagnostics, ReadExtraRequestCustomParams)

type GraphToTerraformMiddlewareParams struct {
	RawVal map[string]any
}
type GraphToTerraformMiddlewareReturns error

type TerraformToGraphMiddlewareParams struct {
	RawVal   map[string]any
	IsUpdate bool
}
type TerraformToGraphMiddlewareReturns error

type CreateModifyFuncParams struct {
	R              *GenericResource
	Req            resource.CreateRequest
	Resp           *resource.CreateResponse
	BaseUri        string
	Id             string
	UpdateExisting bool
}

type UpdateModifyFuncParams struct {
	R       *GenericResource
	Req     resource.UpdateRequest
	Resp    *resource.UpdateResponse
	BaseUri string
	Id      string
}

type DeleteModifyFuncParams struct {
	R          *GenericResource
	Req        resource.DeleteRequest
	Resp       *resource.DeleteResponse
	BaseUri    string
	Id         string
	SkipDelete bool
}

// lifted from github.com/hashicorp/terraform-plugin-framework/internal/privatestate
type PrivateDataGetSetter interface {
	GetKey(context.Context, string) ([]byte, diag.Diagnostics)
	SetKey(context.Context, string, []byte) diag.Diagnostics
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
			ap.TerraformToGraphMiddleware = func(_ context.Context, params TerraformToGraphMiddlewareParams) TerraformToGraphMiddlewareReturns {
				delete(params.RawVal, parentIdFieldCamel)
				return nil
			}
		}

	})
}
