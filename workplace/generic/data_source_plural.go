package generic

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

	// Dev Note: Try to keep contents similar to CreateGenericDataSourceSingularFromResource

	//
	// IMPORTANT for the whole function: ConvertResourceAttributesToDataSourceSingular does not convert simple type attributes
	// to data source attributes but leaves them as reasource attributes (as their required interface is identical).
	// Therefore additional attributes will also be added from the resource schema.
	//

	//
	// Preparations

	accessParams := &genericResource.AccessParams
	accessParams.InitializeGuarded(nil)

	if accessParams.SingularEntity.NoPluralDataSource {
		panic(fmt.Sprintf("Cannot create plural data source for %s", genericResource.TypeNameSuffix))
	}

	if typeNameSuffix == "" {
		typeNameSuffix = genericResource.TypeNameSuffix + "s"
		if strings.HasSuffix(typeNameSuffix, "ys") {
			typeNameSuffix = strings.TrimSuffix(typeNameSuffix, "ys") + "ies"
		}
	}

	result := GenericDataSourcePlural{
		TypeNameSuffix:                 typeNameSuffix,
		AccessParams:                   *accessParams, // copylocks is fine here (but VSCode would only allow to disable it globally)
		ResultAttributeName:            typeNameSuffix,
		odataFilterAttrsWithGraphNames: map[string]string{},
	}

	//
	// Select and convert attributes (for schema root as well as for nested list)
	// Also collect list of Graph prop names to be passed as OData $select on queries later

	odataSelectGraphNames := make([]string, 0)

	isRootAndRequiredAttrFunc := func(attrName string) bool {
		// For ExtraFilterAttributes, .Required will be reset in ODataPrepareFilterAttributes
		return accessParams.ParentEntities.ContainsFieldName(attrName) ||
			(!accessParams.ReadOptions.DataSource.NoFilterSupport && slices.Contains(accessParams.ReadOptions.DataSource.ExtraFilterAttributes, attrName))
	}

	nestedListCandidatesAttrNames := []string{
		accessParams.IdNameTf,
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

	isNestedListAttrFunc := func(attrName string, rsAttribute rsschema.Attribute) (bool, dsschema.Attribute) {
		if slices.Contains(nestedListCandidatesAttrNames, attrName) {
			graphPropName := (*ToFromGraphTranslator).GraphAttributeNameFromTerraformNameImpl(nil, attrName, rsAttribute)
			odataSelectGraphNames = append(odataSelectGraphNames, graphPropName)
			return true, nil
		} else if nestedDerivedOld, ok := rsAttribute.(OdataDerivedTypeNestedAttributeRs); ok {
			// Also include virtual attributes that signal a derived OData type. They will not have any child attributes or content
			// at all, but it might still be helpful to be able to check for a derived type by comparing with null in Terraform.
			return true, OdataDerivedTypeNestedAttributeDs{
				SingleNestedAttribute: dsschema.SingleNestedAttribute{
					CustomType:  nestedDerivedOld.CustomType,
					Required:    false,
					Optional:    false,
					Computed:    true,
					Sensitive:   nestedDerivedOld.Sensitive,
					Description: nestedDerivedOld.Description,
					MarkdownDescription: fmt.Sprintf("Please note that this nested object does not have any attributes "+
						"but only exists to be able to test if the parent object is of derived OData type `%s` "+
						"(using e.g. `if x.%s != null`).", nestedDerivedOld.DerivedType, attrName),
					DeprecationMessage: nestedDerivedOld.DeprecationMessage,
				},
				DerivedType: nestedDerivedOld.DerivedType,
			}
		}
		return false, nil
	}

	dsAttributes := ConvertResourceAttributesToDataSourcePlural(genericResource.SpecificSchema.Attributes, typeNameSuffix,
		isRootAndRequiredAttrFunc, isNestedListAttrFunc)
	if !accessParams.ReadOptions.DataSource.Plural.NoSelectSupport {
		result.odataSelectGraphNames = odataSelectGraphNames
	}

	ODataPrepareFilterAttributes(accessParams, dsAttributes, result.odataFilterAttrsWithGraphNames, false, true)

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

	odataFilter, odataOrderby, odataTop := ODataGetFilterFromConfig(ctx, &resp.Diagnostics, &req.Config, &d.AccessParams, d.odataFilterAttrsWithGraphNames)
	if resp.Diagnostics.HasError() {
		return
	}

	// do not consider ODataExpand here as it would automatically include the attribute(s) in the result
	rawVal := d.AccessParams.ReadRaw3(ctx, &resp.Diagnostics, uri, "", odataFilter, d.odataSelectGraphNames, odataOrderby, odataTop, false)
	if resp.Diagnostics.HasError() {
		return
	}

	val := ConvertOdataRawToTerraform(ctx, &resp.Diagnostics, req.Config.Schema.(dsschema.Schema),
		rawVal, d.ResultAttributeName, d.AccessParams.GraphToTerraformMiddleware)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State = tfsdk.State{
		Schema: req.Config.Schema,
		Raw:    val,
	}
	d.AccessParams.PopulateStateParentIdsFromRequest(ctx, &resp.Diagnostics, &resp.State, req.Config)
}
