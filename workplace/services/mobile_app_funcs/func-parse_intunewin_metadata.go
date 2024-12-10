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
var _ function.Function = &ParseIntunewinMetadataFunction{}

type ParseIntunewinMetadataFunction struct{}

func (f *ParseIntunewinMetadataFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "parse_intunewin_metadata"
}

func (f *ParseIntunewinMetadataFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Parses IntuneWin file metadata",
		Description: "Reads the file Detection.xml from an IntuneWin file and returns its contents" +
			generic.GetSubcategorySuffixFromMarkdownDescription(services.MobileAppResource.SpecificSchema.MarkdownDescription),

		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "intunewin_file",
				Description: "IntuneWin file to parse",
			},
		},
		Return: function.ObjectReturn{
			AttributeTypes: map[string]attr.Type{
				"name":                     types.StringType,
				"unencrypted_content_size": types.Int64Type,
				"file_name":                types.StringType,
				"setup_file":               types.StringType,
				"encryption_info": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"encryption_key":        types.StringType,
						"mac_key":               types.StringType,
						"initialization_vector": types.StringType,
						"mac":                   types.StringType,
						"profile_identifier":    types.StringType,
						"file_digest":           types.StringType,
						"file_digest_algorithm": types.StringType,
					},
				},
				"msi_info": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"msi_product_code":                  types.StringType,
						"msi_product_version":               types.StringType,
						"msi_package_code":                  types.StringType,
						"msi_upgrade_code":                  types.StringType,
						"msi_execution_context":             types.StringType,
						"msi_requires_logon":                types.BoolType,
						"msi_requires_reboot":               types.BoolType,
						"msi_is_machine_install":            types.BoolType,
						"msi_is_user_install":               types.BoolType,
						"msi_includes_services":             types.BoolType,
						"msi_includes_odbc_data_source":     types.BoolType,
						"msi_contains_system_registry_keys": types.BoolType,
						"msi_contains_system_folders":       types.BoolType,
						"msi_publisher":                     types.StringType,
					},
				},
			},
		},
	}
}

func (f *ParseIntunewinMetadataFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {

	var intunewinFile string
	resp.Error = req.Arguments.Get(ctx, &intunewinFile)
	if resp.Error != nil {
		return
	}

	itwMetadata, err := mobileappcontent.ParseIntunewinMetadata(ctx, intunewinFile)
	if err != nil {
		resp.Error = function.NewFuncError(fmt.Sprintf("Error parsing IntuneWin file: %v", err))
		return
	}

	resp.Error = resp.Result.Set(ctx, &itwMetadata)
	if resp.Error != nil {
		return
	}
}
