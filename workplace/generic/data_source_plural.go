package generic

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/manicminer/hamilton/msgraph"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &GenericDataSourcePlural{}
	_ datasource.DataSourceWithConfigure = &GenericDataSourcePlural{}
)

// GenericDataSourcePlural is the data source implementation.
type GenericDataSourcePlural struct {
	TypeNameSuffix      string
	SpecificSchema      schema.Schema
	AccessParams        AccessParams
	ResultAttributeName string

	odataSelect []string
}

func CreateGenericDataSourcePluralFromSingular(genericSingular *GenericDataSourceSingular, typeNameSuffix string) GenericDataSourcePlural {
	if genericSingular.AccessParams.IsSingleton {
		panic(fmt.Sprintf("Cannot create plural data source from Singleton (%s)", genericSingular.TypeNameSuffix))
	}

	if typeNameSuffix == "" {
		typeNameSuffix = genericSingular.TypeNameSuffix + "s"
		if strings.HasSuffix(typeNameSuffix, "ys") {
			typeNameSuffix = strings.TrimSuffix(typeNameSuffix, "ys") + "ies"
		}
	}
	return GenericDataSourcePlural{
		TypeNameSuffix:      typeNameSuffix,
		SpecificSchema:      ConvertDataSourceSingularSchemaToPlural(&genericSingular.SpecificSchema, typeNameSuffix, &genericSingular.AccessParams.HasParentItem),
		AccessParams:        genericSingular.AccessParams, // copylocks is fine here (but VSCode would only allow to disable it globally)
		ResultAttributeName: typeNameSuffix,
	}
}

// Metadata returns the data source type name.
func (d *GenericDataSourcePlural) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, d.TypeNameSuffix)
}

// Configure adds the provider configured MSGraph client to the data source.
func (d *GenericDataSourcePlural) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData != nil {
		d.AccessParams.graphClient = req.ProviderData.(*msgraph.Client)
	}
}

// Schema defines the schema for the data source.
func (d *GenericDataSourcePlural) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	if !d.AccessParams.ReadOptions.PluralNoFilterSupport {
		d.SpecificSchema.Attributes["odata_filter"] = schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: `Raw OData $filter string to pass to MS Graph.`,
		}
		d.SpecificSchema.Attributes["include_ids"] = schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: `Filter query to only return objects with these ids.`,
		}
		d.SpecificSchema.Attributes["exclude_ids"] = schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			MarkdownDescription: `Filter query to exclude objects with these ids.`,
		}
	}
	resp.Schema = d.SpecificSchema

	d.odataSelect = nil // clear in case of duplicate invocation
	resultSetNestedAttribute := d.SpecificSchema.Attributes[d.ResultAttributeName].(schema.SetNestedAttribute)
	translator := NewToFromGraphTranslator(d.SpecificSchema, false)
	for k, v := range resultSetNestedAttribute.NestedObject.Attributes {
		if _, ok := v.(DerivedTypeNestedAttribute); ok {
			// It's a virtual attribute only that signals a derived OData type. Do not add this to OData select (as MS
			// Graph itself does not know about it and would complain).
			continue
		}

		graphAttributeName, ok := translator.GraphAttributeNameFromTerraformName(ctx, resultSetNestedAttribute, k)
		if !ok {
			panic(fmt.Errorf("ToFromGraphTranslator.GraphAttributeNameFromTerraformName not ok for %s", k))
		}
		d.odataSelect = append(d.odataSelect, graphAttributeName)
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *GenericDataSourcePlural) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	uri := d.AccessParams.GetBaseUri(ctx, &resp.Diagnostics, "", &req.Config)
	if resp.Diagnostics.HasError() {
		return
	}

	odataFilter := d.AccessParams.ReadOptions.ODataFilter
	if !d.AccessParams.ReadOptions.PluralNoFilterSupport {

		var configOdataFilter *string
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("odata_filter"), &configOdataFilter)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if configOdataFilter != nil {
			appendODataFilter(&odataFilter, *configOdataFilter, false)
		}

		appendODataFilter(&odataFilter, buildIdFilterInner(ctx, &resp.Diagnostics, &req.Config, path.Root("include_ids")), false)
		appendODataFilter(&odataFilter, buildIdFilterInner(ctx, &resp.Diagnostics, &req.Config, path.Root("exclude_ids")), true)
	}

	// do not consider ODataExpand here as it would automatically include the attribute(s) in the result
	rawVal := d.AccessParams.ReadRaw(ctx, &resp.Diagnostics, uri, "", odataFilter, d.odataSelect, false)
	if resp.Diagnostics.HasError() {
		return
	}

	val := ConvertOdataRawToTerraform(ctx, &resp.Diagnostics, req.Config.Schema.(schema.Schema),
		rawVal, d.ResultAttributeName, d.AccessParams.GraphToTerraformMiddleware)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State = tfsdk.State{
		Schema: req.Config.Schema,
		Raw:    val,
	}
}

func appendODataFilter(baseFilter *string, appendFilter string, useNot bool) {
	if appendFilter == "" {
		return
	}

	if useNot {
		appendFilter = fmt.Sprintf("not(%s)", appendFilter)
	}
	if *baseFilter == "" {
		*baseFilter = appendFilter
	} else {
		*baseFilter = fmt.Sprintf("(%s) and (%s)", *baseFilter, appendFilter)
	}
}

func buildIdFilterInner(ctx context.Context, diags *diag.Diagnostics, filterAttributer GetAttributer, filterAttributPath path.Path) string {
	var ids []string
	diags.Append(filterAttributer.GetAttribute(ctx, filterAttributPath, &ids)...)
	if diags.HasError() || ids == nil {
		return ""
	}

	idFilter := ""
	for _, id := range ids {
		if idFilter != "" {
			idFilter += " or "
		}
		idFilter += fmt.Sprintf("id eq '%s'", id)
	}
	return idFilter
}
