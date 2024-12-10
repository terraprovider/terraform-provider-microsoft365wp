package mobileappfuncs

import (
	"context"
	"fmt"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/services"
	"terraform-provider-microsoft365wp/workplace/util/mobileappcontent"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the desired interfaces.
var _ function.Function = &ParseAppxMsixMetadataFunction{}

type ParseAppxMsixMetadataFunction struct{}

func (f *ParseAppxMsixMetadataFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "parse_appxmsix_metadata"
}

func (f *ParseAppxMsixMetadataFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Parses APPX/MSIX file metadata",
		Description: "Reads the file AppxManifest.xml (and the logo if available) from an APPX/MSIX file and returns its contents" +
			generic.GetSubcategorySuffixFromMarkdownDescription(services.MobileAppResource.SpecificSchema.MarkdownDescription),

		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "appxmsix_file",
				Description: "APPX/MSIX file to parse",
			},
		},
		Return: function.ObjectReturn{
			AttributeTypes: map[string]attr.Type{
				"identity": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":                   types.StringType,
						"processor_architecture": types.StringType,
						"publisher":              types.StringType,
						"publisher_hash":         types.StringType,
						"version":                types.StringType,
					},
				},
				"properties": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"framework":              types.BoolType,
						"display_name":           types.StringType,
						"publisher_display_name": types.StringType,
						"description":            types.StringType,
						"logo":                   types.StringType,
					},
				},
				"logo": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"type":         types.StringType,
						"value_base64": types.StringType,
					},
				},
			},
		},
	}
}

func (f *ParseAppxMsixMetadataFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {

	var appxMsixFile string
	resp.Error = req.Arguments.Get(ctx, &appxMsixFile)
	if resp.Error != nil {
		return
	}

	appxMsixMetadata, err := mobileappcontent.ParseAppxMsixMetadata(ctx, appxMsixFile)
	if err != nil {
		resp.Error = function.NewFuncError(fmt.Sprintf("Error parsing APPX/MSIX file: %v", err))
		return
	}

	resp.Error = resp.Result.Set(ctx, &appxMsixMetadata)
	if resp.Error != nil {
		return
	}
}
