package generic

import (
	"context"
	"slices"
	"sync"
	"terraform-provider-microsoft365wp/workplace/external/strcase"
	"time"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type AccessParams struct {
	BaseUri        string
	ParentEntities ParentEntities
	UriNoId        bool   // use plain base (or parent) URI without appending id (e.g. for singletons)
	UriSuffix      string // append this to URI after the id (e.g. for singular childs)
	EntityId       EntityIdOptions

	IsSingleton  bool // just a shortcut (will be read and applied to other fields inside InitializeGuarded)
	ReadOptions  ReadOptions
	WriteOptions WriteOptions

	GraphToTerraformMiddleware                     GraphToTerraformMiddlewareFunc
	GraphToTerraformMiddlewareTargetSetRunOnRawVal bool
	TerraformToGraphMiddleware                     TerraformToGraphMiddlewareFunc

	CreateModifyFunc func(context.Context, *diag.Diagnostics, *CreateModifyFuncParams) // pass pointer to allow modifications
	UpdateModifyFunc func(context.Context, *diag.Diagnostics, *UpdateModifyFuncParams)
	DeleteModifyFunc func(context.Context, *diag.Diagnostics, *DeleteModifyFuncParams)

	CreateReplaceFunc func(context.Context, *diag.Diagnostics, *CreateReplaceFuncParams)
	UpdateReplaceFunc func(context.Context, *diag.Diagnostics, *UpdateReplaceFuncParams)
	DeleteReplaceFunc func(context.Context, *diag.Diagnostics, *DeleteReplaceFuncParams)

	initializeOnce sync.Once

	graphClient *msgraph.Client
}

type ParentEntities []ParentEntity

type ParentEntity struct {
	ParentIdField    path.Path
	UriSuffix        string
	ParentIdValueMap map[string]string
}

func (pes ParentEntities) ContainsFieldName(fieldName string) bool {
	return slices.ContainsFunc(pes, func(pe ParentEntity) bool {
		return fieldName == pe.ParentIdField.String()
	})
}

type EntityIdOptions struct {
	AttrNameGraph                             string
	AttrNameTf                                string
	FallbackIdGetterFunc                      func(ctx context.Context, diags *diag.Diagnostics, idAttributer GetAttributer) string
	CreateUpdateReadIdFromGraph               bool // read id attribute value from MS Graph (e.g. for singletons) - otherwise it needs to be in the TF config (e.g. for singular child resources the id of the parent)
	CreateUpdateReadIdFromGraphUseODataSelect bool // use $select=id when reading id attribute value from MS Graph (will only work for a few resources)
}

type ReadOptions struct {
	ODataExpand              string
	ODataFilter              string // GenericDataSourcePlural and SingleItemUseODataFilter only, not allowed otherwise
	ODataSelect              []string
	SingleItemUseODataFilter bool
	ValidStatusCodesExtra    []int
	ExtraRequests            []ReadExtraRequest
	ExtraRequestsCustom      []ReadExtraRequestCustom
	DataSource               DataSourceOptions
}

type ReadExtraRequest struct {
	Attribute string
	UriSuffix string
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

type DataSourceOptions struct {
	NoFilterSupport       bool
	NoIdFilterSupport     bool
	ExtraFilterAttributes []string
	Singular              SingularOptions
	Plural                PluralOptions
}

type SingularOptions struct {
	NoIdRequired bool
}

type PluralOptions struct {
	NoDataSource    bool // panic on trying to create plural data source (e.g. for singletons, singular childs etc.)
	NoSelectSupport bool
	ExtraAttributes []string
}

type WriteOptions struct {
	UpdateInsteadOfCreate  bool // update existing entity instead of trying to create a new one (e.g. for singletons, singular childs etc.)
	UpdateIfExistsOnCreate bool // update existing entity only in case one already exists (only works if id gets provided in TF config or by FallbackIdGetterFunc)
	UsePutForCreate        bool // will most likely only work if id gets provided in TF config or by FallbackIdGetterFunc (since PUT requests usually do not return any content and hence no id gets provided by MS Graph)
	UsePutForUpdate        bool
	SkipDelete             bool // just remove from TF state but do not try to delete from graph (e.g. for singletons)
	SubActions             []WriteSubAction
	SerializeWrites        bool
	SerialWritesDelay      time.Duration
}

type GraphToTerraformMiddlewareFunc func(context.Context, *diag.Diagnostics, *GraphToTerraformMiddlewareParams) GraphToTerraformMiddlewareReturns
type GraphToTerraformMiddlewareParams struct {
	ExpectedId          string
	RawVal              map[string]any
	IsTargetSetOnRawVal bool
}
type GraphToTerraformMiddlewareReturns error

type TerraformToGraphMiddlewareFunc func(context.Context, *diag.Diagnostics, *TerraformToGraphMiddlewareParams) TerraformToGraphMiddlewareReturns
type TerraformToGraphMiddlewareParams struct {
	Config   tfsdk.Config
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

type CreateReplaceFuncParams struct {
	R            *GenericResource
	Client       *msgraph.Client
	BaseUri      string
	IdAttributer GetAttributer
	RawVal       map[string]any
	RawResult    map[string]any
	Id           string
}

type UpdateReplaceFuncParams struct {
	R            *GenericResource
	Client       *msgraph.Client
	BaseUri      string
	Id           string
	IdAttributer GetAttributer
	RawVal       map[string]any
	RawResult    map[string]any
}

type DeleteReplaceFuncParams struct {
	R            *GenericResource
	Client       *msgraph.Client
	BaseUri      string
	Id           string
	IdAttributer GetAttributer
}

// lifted from github.com/hashicorp/terraform-plugin-framework/internal/privatestate
type PrivateDataGetSetter interface {
	GetKey(context.Context, string) ([]byte, diag.Diagnostics)
	SetKey(context.Context, string, []byte) diag.Diagnostics
}

func (ap *AccessParams) InitializeGuarded(providerData any) {

	// providerData may be empty on early calls, so we always check this here
	if ap.graphClient == nil && providerData != nil {
		if graphClient, ok := providerData.(*msgraph.Client); ok {
			ap.graphClient = graphClient
		}
	}

	ap.initializeOnce.Do(func() {

		if len(ap.ParentEntities) > 0 {
			parentIdFieldsCamel := make([]string, 0)
			for _, pe := range ap.ParentEntities {
				parentIdFieldsCamel = append(parentIdFieldsCamel, strcase.ToLowerCamel(pe.ParentIdField.String()))
				if pe.ParentIdValueMap == nil {
					pe.ParentIdValueMap = make(map[string]string)
				}
			}
			oldMiddleware := ap.TerraformToGraphMiddleware
			ap.TerraformToGraphMiddleware = func(ctx context.Context, diags *diag.Diagnostics, params *TerraformToGraphMiddlewareParams) TerraformToGraphMiddlewareReturns {
				for _, f := range parentIdFieldsCamel {
					delete(params.RawVal, f)
				}
				if oldMiddleware != nil {
					return oldMiddleware(ctx, diags, params)
				} else {
					return nil
				}
			}
		}

		if ap.EntityId.AttrNameGraph == "" {
			ap.EntityId.AttrNameGraph = "id"
		}
		if ap.EntityId.AttrNameTf == "" {
			ap.EntityId.AttrNameTf = strcase.ToSnake(ap.EntityId.AttrNameGraph)
		}

		if ap.IsSingleton {
			ap.UriNoId = true
			ap.EntityId.CreateUpdateReadIdFromGraph = true
			ap.ReadOptions.DataSource.NoFilterSupport = true
			ap.ReadOptions.DataSource.Singular.NoIdRequired = true
			ap.ReadOptions.DataSource.Plural.NoDataSource = true
			ap.WriteOptions.UpdateInsteadOfCreate = true
			ap.WriteOptions.SkipDelete = true
		}

		for _, wsa := range ap.WriteOptions.SubActions {
			wsa.Initialize()
		}

	})
}
