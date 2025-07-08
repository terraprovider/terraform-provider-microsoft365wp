package generic

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"golang.org/x/exp/slices"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &GenericDataSourcePlural{}
	_ datasource.DataSourceWithConfigure = &GenericDataSourcePlural{}
)

// GenericDataSourcePlural is the data source implementation.
type GenericDataSourcePlural struct {
	TypeNameSuffix      string
	SpecificSchema      dsschema.Schema
	AccessParams        AccessParams
	ResultAttributeName string

	odataFilterAttrsWithGraphNames map[string]string
	odataSelectGraphNames          []string
}

func CreateGenericDataSourcePluralFromResource(genericResource *GenericResource, typeNameSuffix string) GenericDataSourcePlural {

	// Dev Note: Try to keep this func similar to CreateGenericDataSourceSingularFromResource

	// IMPORTANT for the whole function: ConvertResourceAttributesToDataSourceSingular does not convert simple type attributes
	// to data source attributes but leaves them as reasource attributes (as their required interface is identical).
	// Therefore additional attributes will also be added from the resource schema.

	// Preparations
	accessParams := &genericResource.AccessParams
	accessParams.InitializeGuarded(nil)
	rsSchema := &genericResource.SpecificSchema
	if accessParams.ReadOptions.DataSource.Plural.NoDataSource {
		panic(fmt.Sprintf("Cannot create plural data source for %s", genericResource.TypeNameSuffix))
	}

	if typeNameSuffix == "" {
		typeNameSuffix = genericResource.TypeNameSuffix + "s"
		if strings.HasSuffix(typeNameSuffix, "ys") {
			typeNameSuffix = strings.TrimSuffix(typeNameSuffix, "ys") + "ies"
		}
	}

	dsResult := GenericDataSourcePlural{
		TypeNameSuffix: typeNameSuffix,
		SpecificSchema: dsschema.Schema{
			Attributes:          map[string]dsschema.Attribute{},
			Description:         rsSchema.Description,
			MarkdownDescription: rsSchema.MarkdownDescription,
		},
		AccessParams:                   *accessParams, // copylocks is fine here (but VSCode would only allow to disable it globally)
		ResultAttributeName:            typeNameSuffix,
		odataFilterAttrsWithGraphNames: map[string]string{},
		odataSelectGraphNames:          []string{},
	}
	dsAttributesRoot := dsResult.SpecificSchema.Attributes

	//
	// Select and convert attributes (for schema root as well as for nested list)
	// Also collect list of Graph prop names to be passed as OData $select on queries later

	dsAttributesNestedList := make(map[string]dsschema.Attribute)
	dsAttributesRoot[typeNameSuffix] = dsschema.ListNestedAttribute{
		Computed: true,
		NestedObject: dsschema.NestedAttributeObject{
			Attributes: dsAttributesNestedList,
		},
	}

	nestedListCandidatesAttrNames := []string{
		accessParams.EntityId.AttrNameTf,
		"display_name",
		"name",
		"role_scope_tag_ids",
		"role_scope_tags",
		"version",
		"created_date_time",
		"last_modified_date_time",
		"modified_date_time",
		"deleted_date_time",
	}
	nestedListCandidatesAttrNames = append(nestedListCandidatesAttrNames, accessParams.ReadOptions.DataSource.Plural.ExtraAttributes...)

	if !accessParams.ReadOptions.DataSource.Plural.NoSelectSupport {
		dsResult.odataSelectGraphNames = append(dsResult.odataSelectGraphNames, accessParams.ReadOptions.ODataSelect...)
	}

	for attrName, rsAttribute := range rsSchema.Attributes {

		// Schema root
		identifiesParent := accessParams.ParentEntities.ContainsFieldName(attrName)
		isExtraFilterAttr := !accessParams.ReadOptions.DataSource.NoFilterSupport && slices.Contains(accessParams.ReadOptions.DataSource.ExtraFilterAttributes, attrName)
		if identifiesParent || isExtraFilterAttr {
			// For ExtraFilterAttributes, .Required will be adjusted inside ODataPrepareFilterAttributes
			dsAttributesRoot[attrName] = convertResourceAttr2DataSourceAttr(rsAttribute, true, true, attrName)
		}

		// Nested List
		isForNestedList := slices.Contains(nestedListCandidatesAttrNames, attrName)
		// Also include virtual attributes that signal a derived OData type. They will not have any child attributes or content
		// at all, but it might still be helpful to be able to check for a derived type by comparing with null in Terraform.
		_, isDerivedTypeNested := rsAttribute.(OdataDerivedTypeNestedAttributeRs)
		if isForNestedList || isDerivedTypeNested {
			dsAttributesNestedList[attrName] = convertResourceAttr2DataSourceAttr(rsAttribute, false, true, attrName)
			// only "real" attributes existing in MS Graph need to be added to OData $select
			if isForNestedList && !accessParams.ReadOptions.DataSource.Plural.NoSelectSupport {
				attrNameGraph := (*ToFromGraphTranslator).GraphAttributeNameFromTerraformNameImpl(nil, attrName, rsAttribute)
				// attrNameGraph might have already been included in ReadOptions.ODataSelect
				if !slices.Contains(dsResult.odataSelectGraphNames, attrNameGraph) {
					dsResult.odataSelectGraphNames = append(dsResult.odataSelectGraphNames, attrNameGraph)
				}
			}
		}
	}

	dsPrepareFilterAttributes(accessParams, dsAttributesRoot, dsResult.odataFilterAttrsWithGraphNames, nil, false, true)

	return dsResult // copylocks is fine here (but VSCode would only allow to disable it globally)
}

// Metadata returns the data source type name.
func (d *GenericDataSourcePlural) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, d.TypeNameSuffix)
}

// Configure adds the provider configured MSGraph client to the data source.
func (d *GenericDataSourcePlural) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	d.AccessParams.InitializeGuarded(req.ProviderData)
}

// Schema defines the schema for the data source.
func (d *GenericDataSourcePlural) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = d.SpecificSchema
}

// Read refreshes the Terraform state with the latest data.
func (d *GenericDataSourcePlural) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	uri := d.AccessParams.GetBaseUri(ctx, &resp.Diagnostics, "", &req.Config)
	if resp.Diagnostics.HasError() {
		return
	}

	odataFilter, odataOrderby, odataTop := dsGetFilterFromConfig(ctx, &resp.Diagnostics, &req.Config, &d.AccessParams, d.odataFilterAttrsWithGraphNames)
	if resp.Diagnostics.HasError() {
		return
	}

	// do not consider ODataExpand here as it would automatically include the attribute(s) in the result
	rawVal := d.AccessParams.ReadRaw3(ctx, &resp.Diagnostics, uri, "", odataFilter, d.odataSelectGraphNames, odataOrderby, odataTop, false)
	if resp.Diagnostics.HasError() {
		return
	}

	val := ConvertOdataRawToTerraform(ctx, &resp.Diagnostics, req.Config.Schema.(dsschema.Schema), rawVal, d.ResultAttributeName,
		d.AccessParams.GraphToTerraformMiddleware, "", d.AccessParams.GraphToTerraformMiddlewareTargetSetRunOnRawVal)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State = tfsdk.State{
		Schema: req.Config.Schema,
		Raw:    val,
	}
	d.AccessParams.PopulateStateParentIdsFromRequest(ctx, &resp.Diagnostics, &resp.State, req.Config)
}
