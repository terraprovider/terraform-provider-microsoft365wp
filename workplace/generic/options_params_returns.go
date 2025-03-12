package generic

import (
	"context"
	"slices"
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
	BaseUri        string
	IdNameGraph    string
	IdNameTf       string
	ParentEntities ParentEntities
	SingularEntity SingularEntity
	IsSingleton    bool // just a shortcut (will be read and applied to SingularEntity inside InitializeGuarded)

	ReadOptions ReadOptions

	GraphToTerraformMiddleware func(context.Context, *diag.Diagnostics, *GraphToTerraformMiddlewareParams) GraphToTerraformMiddlewareReturns
	TerraformToGraphMiddleware func(context.Context, *diag.Diagnostics, *TerraformToGraphMiddlewareParams) TerraformToGraphMiddlewareReturns

	CreateModifyFunc func(context.Context, *diag.Diagnostics, *CreateModifyFuncParams) // pass pointer to allow modifications
	UpdateModifyFunc func(context.Context, *diag.Diagnostics, *UpdateModifyFuncParams)
	DeleteModifyFunc func(context.Context, *diag.Diagnostics, *DeleteModifyFuncParams)

	WriteSubActions []WriteSubAction

	UsePutForUpdate bool

	ReplaceUpdateFunc func(context.Context, *diag.Diagnostics, *ReplaceUpdateFuncParams)

	SerializeWrites   bool
	SerialWritesDelay time.Duration

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

type SingularEntity struct {
	UriNoId                        bool   // use plain BaseUri without appending id (e.g. for singletons)
	ReadIdAttrFromGraph            bool   // read id attribute value from MS Graph (e.g. for singletons) - otherwise it needs to be in the TF config (e.g. for singular child resources the id of the parent)
	ReadIdAttrFromGraphUseSelectId bool   // use $select=id when reading id attribute value from MS Graph (will only work for a few resources)
	UriSuffix                      string // append this to (parent) uri (e.g. for singular childs)
	UpdateInsteadOfCreate          bool   // update existing entity instead of trying to create a new one (e.g. for singletons, singular childs etc.)
	SkipDelete                     bool   // just remove from TF state but do not try to delete from graph (e.g. for singletons)
	NoPluralDataSource             bool   // panic on trying to create plural data source (e.g. for singletons, singular childs etc.)
}

type ReadOptions struct {
	ODataExpand              string
	ODataFilter              string // GenericDataSourcePlural and SingleItemUseODataFilter only, not allowed otherwise
	SingleItemUseODataFilter bool
	ValidStatusCodesExtra    []int
	ExtraRequests            []ReadExtraRequest
	ExtraRequestsCustom      []ReadExtraRequestCustom
	DataSource               DataSourceOptions
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

type DataSourceOptions struct {
	NoFilterSupport       bool
	ExtraFilterAttributes []string
	Plural                PluralOptions
}

type PluralOptions struct {
	NoSelectSupport bool
	ExtraAttributes []string
}

type GraphToTerraformMiddlewareParams struct {
	RawVal map[string]any
}
type GraphToTerraformMiddlewareReturns error

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

type ReplaceUpdateFuncParams struct {
	R            *GenericResource
	BaseUri      string
	Id           string
	IdAttributer GetAttributer
	RawVal       map[string]any
	RawResult    map[string]any
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

		if ap.IdNameGraph == "" {
			ap.IdNameGraph = "id"
		}
		if ap.IdNameTf == "" {
			ap.IdNameTf = strcase.ToSnake(ap.IdNameGraph)
		}

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

		if ap.IsSingleton {
			singularEntity := &ap.SingularEntity
			singularEntity.UriNoId = true
			singularEntity.ReadIdAttrFromGraph = true
			singularEntity.UpdateInsteadOfCreate = true
			singularEntity.SkipDelete = true
			singularEntity.NoPluralDataSource = true
			ap.ReadOptions.DataSource.NoFilterSupport = true
		}

		for _, wsa := range ap.WriteSubActions {
			wsa.Initialize()
		}

	})
}
