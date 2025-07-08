package generic

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/exp/slices"
)

func dsPrepareFilterAttributes(accessParams *AccessParams, attributes map[string]dsschema.Attribute,
	odataFilterAttrsWithGraphNames map[string]string, dsConfigValidators *[]datasource.ConfigValidator, isSingular bool, isPlural bool) {

	dsOptions := &accessParams.ReadOptions.DataSource
	if dsOptions.NoFilterSupport && dsOptions.Singular.NoIdRequired {
		// nothing to do for use here
		return
	}

	// Preparations
	// odataFilterAttrsWithGraphNames: attributes to be processed when running query later
	filterAttributes := []string{}       // to be adjusted to accept config values and to be added to odataFilterAttrsWithGraphNames
	atLeastOneOfAttributes := []string{} // to be included in potential AtLeastOneOf validator

	// Include (existing) id attribute for singular data sources (if entity is not singleton)
	if isSingular && !dsOptions.Singular.NoIdRequired {
		filterAttributes = append(filterAttributes, accessParams.EntityId.AttrNameTf)
		atLeastOneOfAttributes = append(atLeastOneOfAttributes, accessParams.EntityId.AttrNameTf)
	}

	// Add (additional) filter attributes
	if !dsOptions.NoFilterSupport {

		attributes["odata_filter"] = rsschema.StringAttribute{
			MarkdownDescription: "Literal OData `$filter` value to pass to MS Graph.",
		}
		filterAttributes = append(filterAttributes, "odata_filter")
		atLeastOneOfAttributes = append(atLeastOneOfAttributes, "odata_filter")

		attributes["odata_orderby"] = rsschema.StringAttribute{
			MarkdownDescription: "Literal OData `$orderby` value to pass to MS Graph.",
		}
		filterAttributes = append(filterAttributes, "odata_orderby") // not applicable for atLeastOneOfAttributes

		attributes["odata_top"] = rsschema.Int32Attribute{
			MarkdownDescription: "Literal OData `$top` value to pass to MS Graph.",
		}
		filterAttributes = append(filterAttributes, "odata_top") // not applicable for atLeastOneOfAttributes

		if isPlural && !dsOptions.NoIdFilterSupport {
			attributes["include_ids"] = rsschema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Only return entities with these ids (using OData `$filter`).",
			}
			filterAttributes = append(filterAttributes, "include_ids") // not applicable for atLeastOneOfAttributes

			attributes["exclude_ids"] = rsschema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Exclude entities with these ids (using OData `$filter`).",
			}
			filterAttributes = append(filterAttributes, "exclude_ids") // not applicable for atLeastOneOfAttributes
		}

		filterAttributes = append(filterAttributes, dsOptions.ExtraFilterAttributes...)
		atLeastOneOfAttributes = append(atLeastOneOfAttributes, dsOptions.ExtraFilterAttributes...)
	}

	// Create AtLeastOneOf validator (if applicable) for singular data sources to ensure that any criteria has been specified
	if isSingular && len(atLeastOneOfAttributes) > 0 {
		validatorPaths := []path.Expression{}
		for _, name := range atLeastOneOfAttributes {
			validatorPaths = append(validatorPaths, path.MatchRoot(name))
		}
		*dsConfigValidators = append(*dsConfigValidators, datasourcevalidator.AtLeastOneOf(validatorPaths...))
	}

	// Now loop through filter attributes and adjust all required details
	// We only added stub attribute definitions before since we anyway need to update the existing attribute(s) for id and ExtraFilterAttributes.
	for attrName, attr := range attributes {
		if slices.Contains(filterAttributes, attrName) {

			switch typed := attr.(type) {
			case rsschema.StringAttribute:
				typed.Required = false
				typed.Optional = true
				attributes[attrName] = typed
			case rsschema.Int32Attribute:
				typed.Required = false
				typed.Optional = true
				attributes[attrName] = typed
			case rsschema.BoolAttribute:
				typed.Required = false
				typed.Optional = true
				attributes[attrName] = typed
			case rsschema.SetAttribute:
				typed.Required = false
				typed.Optional = true
				attributes[attrName] = typed
			}

			// Collect map of attribute names to MS Graph property names to use to construct OData $filter on queries later
			// Id will not be used in $filter (but inside GetUri...) and is therefore not required here
			if attrName != accessParams.EntityId.AttrNameTf {
				odataFilterAttrsWithGraphNames[attrName] =
					(*ToFromGraphTranslator).GraphAttributeNameFromTerraformNameImpl(nil, attrName, attr)
			}
		}
	}
}

func dsGetFilterFromConfig(ctx context.Context, diags *diag.Diagnostics, config *tfsdk.Config, accessParams *AccessParams,
	odataFilterAttrsWithGraphNames map[string]string) (odataFilter string, odataOrderby string, odataTop int) {

	odataFilter = accessParams.ReadOptions.ODataFilter
	odataOrderby = ""
	odataTop = 0

	if accessParams.ReadOptions.DataSource.NoFilterSupport {
		// nothing to do for use here
		return
	}

	for tfName, graphName := range odataFilterAttrsWithGraphNames {
		attribPath := path.Root(tfName)
		var attrValue attr.Value
		diags.Append(config.GetAttribute(ctx, attribPath, &attrValue)...)
		if !(attrValue.IsNull() || attrValue.IsUnknown()) {

			switch typedAttrValue := attrValue.(type) {

			case types.String:
				stringValue := typedAttrValue.ValueString()

				switch tfName {

				case "odata_filter":
					dsFilterAppend(&odataFilter, stringValue, false)
					continue

				case "odata_orderby":
					odataOrderby = stringValue
					continue

				default:
					// allow wildcards at some positions
					odataFilterRegexp := regexp.MustCompile(`^(\*)?([^\*]+)(\*)?$|^([^\*]+)\*([^\*]+)$`)
					filterParts := odataFilterRegexp.FindStringSubmatch(stringValue)
					if filterParts == nil {
						diags.AddAttributeError(attribPath, "Error creating OData filter",
							"If using wildcards (`*`), they are only allowed at the start and/or the end of the attribute value or else exactly once inside")
						continue
					}

					if filterParts[2] != "" { // 1st alternative matched
						if filterParts[1] == "*" && filterParts[3] == "*" {
							dsFilterAppend(&odataFilter, fmt.Sprintf("contains(%s,'%s')", graphName, filterParts[2]), false)
						} else if filterParts[3] == "*" {
							dsFilterAppend(&odataFilter, fmt.Sprintf("startswith(%s,'%s')", graphName, filterParts[2]), false)
						} else if filterParts[1] == "*" {
							dsFilterAppend(&odataFilter, fmt.Sprintf("endswith(%s,'%s')", graphName, filterParts[2]), false)
						} else {
							dsFilterAppend(&odataFilter, fmt.Sprintf("%s eq '%s'", graphName, filterParts[2]), false)
						}
					} else { // 2nd alternative must have matched (otherwise the regex would not have matched at all)
						dsFilterAppend(&odataFilter, fmt.Sprintf("startswith(%s,'%s') and endswith(%s,'%s')", graphName, filterParts[4], graphName, filterParts[5]), false)
					}
					continue
				}

			case types.Int32:
				int32Value := typedAttrValue.ValueInt32()
				switch tfName {

				case "odata_top":
					odataTop = int(int32Value)
					continue

				default:
					dsFilterAppend(&odataFilter, fmt.Sprintf("%s eq %d", graphName, int32Value), false)
					continue
				}

			case types.Bool:
				boolValue := typedAttrValue.ValueBool()
				dsFilterAppend(&odataFilter, fmt.Sprintf("%s eq %t", graphName, boolValue), false)
				continue

			case types.Set:
				switch tfName {

				case "include_ids", "exclude_ids":
					var ids []string
					diags.Append(typedAttrValue.ElementsAs(ctx, &ids, false)...)
					if diags.HasError() {
						return
					}
					idFilterInner := dsFilterBuildIdFilterInner(accessParams.EntityId.AttrNameGraph, ids)
					useNot := tfName == "exclude_ids"
					dsFilterAppend(&odataFilter, idFilterInner, useNot)
					continue
				}

			}

			diags.AddAttributeError(attribPath, "Error creating OData filter", fmt.Sprintf("Don't know how to handle attribute value of type %s", attrValue.Type(ctx).String()))
		}
	}

	return odataFilter, odataOrderby, odataTop
}

func dsFilterAppend(baseFilter *string, appendFilter string, useNot bool) {
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

func dsFilterBuildIdFilterInner(idNameGraph string, ids []string) string {
	idFilter := ""
	for _, id := range ids {
		if idFilter != "" {
			idFilter += " or "
		}
		idFilter += fmt.Sprintf("%s eq '%s'", idNameGraph, id)
	}
	return idFilter
}
