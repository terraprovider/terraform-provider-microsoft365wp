package generic

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &GenericDataSourceSingular{}
	_ datasource.DataSourceWithConfigure = &GenericDataSourceSingular{}
)

// GenericDataSourceSingular is the data source implementation.
type GenericDataSourceSingular struct {
	TypeNameSuffix string
	SpecificSchema dsschema.Schema
	AccessParams   AccessParams

	odataFilterAttrsWithGraphNames map[string]string
}

func CreateGenericDataSourceSingularFromResource(genericResource *GenericResource) GenericDataSourceSingular {

	// Dev Note: Try to keep contents similar to CreateGenericDataSourcePluralFromResource

	//
	// IMPORTANT for the whole function: ConvertResourceAttributesToDataSourceSingular does not convert simple type attributes
	// to data source attributes but leaves them as reasource attributes (as their required interface is identical).
	// Therefore additional attributes will also be added from the resource schema.
	//

	//
	// Preparations

	accessParams := &genericResource.AccessParams
	accessParams.InitializeGuarded(nil)

	result := GenericDataSourceSingular{
		TypeNameSuffix:                 genericResource.TypeNameSuffix,
		AccessParams:                   *accessParams, // copylocks is fine here (but VSCode would only allow to disable it globally)
		odataFilterAttrsWithGraphNames: map[string]string{},
	}

	//
	// Choose required and convert attributes

	isRequiredFunc := func(attrName string) bool { return accessParams.ParentEntities.ContainsFieldName(attrName) }
	dsAttributes := ConvertResourceAttributesToDataSourceSingular(genericResource.SpecificSchema.Attributes,
		isRequiredFunc)

	ODataPrepareFilterAttributes(accessParams, dsAttributes, result.odataFilterAttrsWithGraphNames, true, false)

	//
	// Create schema and return

	result.SpecificSchema = dsschema.Schema{
		Attributes:          dsAttributes,
		Description:         genericResource.SpecificSchema.Description,
		MarkdownDescription: genericResource.SpecificSchema.MarkdownDescription,
	}

	return result // copylocks is fine here (but VSCode would only allow to disable it globally)
}

// Metadata returns the data source type name.
func (d *GenericDataSourceSingular) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, d.TypeNameSuffix)
}

// Configure adds the provider configured MSGraph client to the data source.
func (d *GenericDataSourceSingular) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	d.AccessParams.InitializeGuarded(req.ProviderData)
}

// Schema defines the schema for the data source.
func (d *GenericDataSourceSingular) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = d.SpecificSchema
}

// Read refreshes the Terraform state with the latest data.
func (d *GenericDataSourceSingular) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	odataFilter, odataOrderby, odataTop := ODataGetFilterFromConfig(ctx, &resp.Diagnostics, &req.Config, &d.AccessParams, d.odataFilterAttrsWithGraphNames)
	if resp.Diagnostics.HasError() {
		return
	}

	val := d.AccessParams.ReadSingleCompleteTf3(ctx, &resp.Diagnostics, req.Config.Schema,
		d.AccessParams.ReadOptions, "", req.Config, odataFilter, odataOrderby, odataTop, false, nil, nil)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update state from read item
	resp.State = tfsdk.State{
		Schema: req.Config.Schema,
		Raw:    val,
	}
	d.AccessParams.PopulateStateParentIdsFromRequest(ctx, &resp.Diagnostics, &resp.State, req.Config)
}
