package generic

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/exp/slices"
)

func ODataPrepareFilterAttributes(accessParams *AccessParams, dsAttributes map[string]dsschema.Attribute,
	odataFilterAttrsWithGraphNames map[string]string, isSingular bool, isPlural bool) {

	//
	// Preparations

	// odataFilterAttrsWithGraphNames: attributes to be processed when running query later
	filterAttributes := []string{}    // to be adjusted to accept config values and to be added to odataFilterAttrsWithGraphNames
	validatorAttributes := []string{} // to be included in potential validator

	//
	// Add (additional) filter attributes

	if !accessParams.ReadOptions.DataSource.NoFilterSupport {
		// odata_filter
		dsAttributes["odata_filter"] = rsschema.StringAttribute{
			MarkdownDescription: "Literal OData `$filter` value to pass to MS Graph.",
		}
		filterAttributes = append(filterAttributes, "odata_filter")
		validatorAttributes = append(validatorAttributes, "odata_filter")

		// odata_orderby
		dsAttributes["odata_orderby"] = rsschema.StringAttribute{
			MarkdownDescription: "Literal OData `$orderby` value to pass to MS Graph.",
		}
		filterAttributes = append(filterAttributes, "odata_orderby")
		// not applicable for validator

		// odata_top (number!)
		// Cannot be added to filterAttributes and therefore must be completely defined here!
		dsAttributes["odata_top"] = rsschema.Int32Attribute{
			Optional:            true,
			MarkdownDescription: "Literal OData `$top` value to pass to MS Graph.",
		}
		odataFilterAttrsWithGraphNames["odata_top"] = "" // just needs to be part of map, value will be ignored
		// not applicable for validator

		if isPlural {

			// include_ids ([]string!)
			// Cannot be added to filterAttributes and therefore must be completely defined here!
			dsAttributes["include_ids"] = rsschema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Only return entities with these ids (using OData `$filter`).",
			}
			odataFilterAttrsWithGraphNames["include_ids"] = "" // just needs to be part of map, value will be ignored
			// not applicable for validator

			// exclude_ids ([]string!)
			// Cannot be added to filterAttributes and therefore must be completely defined here!
			dsAttributes["exclude_ids"] = rsschema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Exclude entities with these ids (using OData `$filter`).",
			}
			odataFilterAttrsWithGraphNames["exclude_ids"] = "" // just needs to be part of map, value will be ignored
			// not applicable for validator
		}
	}

	//
	// Collect lists of attributes to add to validator and to adjust to accept config values

	if !accessParams.SingularEntity.UriNoId {
		if isSingular {
			filterAttributes = append(filterAttributes, accessParams.IdNameTf)
			validatorAttributes = append(validatorAttributes, accessParams.IdNameTf)
		}
		if !accessParams.ReadOptions.DataSource.NoFilterSupport {
			filterAttributes = append(filterAttributes, accessParams.ReadOptions.DataSource.ExtraFilterAttributes...)
			validatorAttributes = append(validatorAttributes, accessParams.ReadOptions.DataSource.ExtraFilterAttributes...)
		}
	}

	//
	// For IsSingular, create AtLeastOneOf validator

	var filterAttrsValidator validator.String
	if isSingular {
		validatorPaths := []path.Expression{}
		for _, name := range validatorAttributes {
			validatorPaths = append(validatorPaths, path.MatchRoot(name))
		}
		filterAttrsValidator = stringvalidator.AtLeastOneOf(validatorPaths...)
	}

	//
	// Loop through filter attributes and adjust parameters
	// Also collect map of attribute names to MS Graph property names (odataFilterAttrsWithGraphNames) to use to construct OData $filter on queries later

	for attrName, attr := range dsAttributes {
		if slices.Contains(filterAttributes, attrName) {

			stringAttr := attr.(rsschema.StringAttribute)
			stringAttr.Required = false
			stringAttr.Optional = true
			if filterAttrsValidator != nil && slices.Contains(validatorAttributes, attrName) {
				stringAttr.Validators = append(stringAttr.Validators, filterAttrsValidator)
			}
			dsAttributes[attrName] = stringAttr

			// id later gets handled inside GetUri... and is therefore not required here
			if attrName != accessParams.IdNameTf {
				odataFilterAttrsWithGraphNames[attrName] =
					(*ToFromGraphTranslator).GraphAttributeNameFromTerraformNameImpl(nil, attrName, attr)
			}
		}
	}
}

func ODataGetFilterFromConfig(ctx context.Context, diags *diag.Diagnostics, config *tfsdk.Config, accessParams *AccessParams,
	odataFilterAttrsWithGraphNames map[string]string) (string, string, int) {

	odataFilter := accessParams.ReadOptions.ODataFilter
	odataOrderby := ""
	odataTop := 0

	if !accessParams.ReadOptions.DataSource.NoFilterSupport {

		for tfName, graphName := range odataFilterAttrsWithGraphNames {

			var intValue *int
			var stringValue *string
			attribPath := path.Root(tfName)

			if tfName == "odata_top" {
				// int
				diags.Append(config.GetAttribute(ctx, attribPath, &intValue)...)
				if intValue != nil {
					odataTop = *intValue
				}
			} else if tfName == "include_ids" {
				// []string
				odataFilterAppend(&odataFilter, odataFilterBuildIdFilterInner(ctx, diags, config, attribPath, accessParams.IdNameGraph), false)
			} else if tfName == "exclude_ids" {
				// []string
				odataFilterAppend(&odataFilter, odataFilterBuildIdFilterInner(ctx, diags, config, attribPath, accessParams.IdNameGraph), true)
			} else {
				// string
				diags.Append(config.GetAttribute(ctx, attribPath, &stringValue)...)
				if stringValue != nil && *stringValue != "" {

					if tfName == "odata_filter" {
						odataFilterAppend(&odataFilter, *stringValue, false)

					} else if tfName == "odata_orderby" {
						odataOrderby = *stringValue

					} else {
						odataFilterRegexp := regexp.MustCompile(`^(\*)?([^\*]+)(\*)?$|^([^\*]+)\*([^\*]+)$`)
						filterParts := odataFilterRegexp.FindStringSubmatch(*stringValue)
						if filterParts == nil {
							diags.AddAttributeError(attribPath, "Unable to parse attribute to create OData filter", "If using wildcards (`*`), they are only allowed at the start and/or the end of the attribute value or else exactly once inside")
							continue
						}

						if filterParts[2] != "" { // 1st alternative matched
							if filterParts[1] == "*" && filterParts[3] == "*" {
								odataFilterAppend(&odataFilter, fmt.Sprintf("contains(%s,'%s')", graphName, filterParts[2]), false)
							} else if filterParts[3] == "*" {
								odataFilterAppend(&odataFilter, fmt.Sprintf("startswith(%s,'%s')", graphName, filterParts[2]), false)
							} else if filterParts[1] == "*" {
								odataFilterAppend(&odataFilter, fmt.Sprintf("endswith(%s,'%s')", graphName, filterParts[2]), false)
							} else {
								odataFilterAppend(&odataFilter, fmt.Sprintf("%s eq '%s'", graphName, filterParts[2]), false)
							}
						} else { // 2nd alternative must have matched (otherwise the regex would not have matched at all)
							odataFilterAppend(&odataFilter, fmt.Sprintf("startswith(%s,'%s') and endswith(%s,'%s')", graphName, filterParts[4], graphName, filterParts[5]), false)
						}
					}
				}
			}
		}
	}

	return odataFilter, odataOrderby, odataTop
}

func odataFilterAppend(baseFilter *string, appendFilter string, useNot bool) {
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

func odataFilterBuildIdFilterInner(ctx context.Context, diags *diag.Diagnostics, filterAttributer GetAttributer, filterAttributPath path.Path, idNameGraph string) string {
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
		idFilter += fmt.Sprintf("%s eq '%s'", idNameGraph, id)
	}
	return idFilter
}
