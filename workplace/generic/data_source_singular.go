package generic

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
	SpecificSchema schema.Schema
	AccessParams   AccessParams
}

func CreateGenericDataSourceSingularFromResource(genericResource *GenericResource) GenericDataSourceSingular {
	return GenericDataSourceSingular{
		TypeNameSuffix: genericResource.TypeNameSuffix,
		SpecificSchema: ConvertResourceSchemaToDataSourceSingular(&genericResource.SpecificSchema, &genericResource.AccessParams),
		AccessParams:   genericResource.AccessParams, // copylocks is fine here (but VSCode would only allow to disable it globally)
	}
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

	val := d.AccessParams.ReadSingleRaw(ctx, &resp.Diagnostics, req.Config.Schema,
		d.AccessParams.ReadOptions, "", req.Config, false, nil, nil)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update state from read item
	resp.State = tfsdk.State{
		Schema: req.Config.Schema,
		Raw:    val,
	}
}
