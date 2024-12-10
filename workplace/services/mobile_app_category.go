package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var (
	MobileAppCategoryResource = generic.GenericResource{
		TypeNameSuffix: "mobile_app_category",
		SpecificSchema: mobileAppCategoryResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/deviceAppManagement/mobileAppCategories",
		},
	}

	MobileAppCategorySingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&MobileAppCategoryResource)

	MobileAppCategoryPluralDataSource = generic.CreateGenericDataSourcePluralFromSingular(
		&MobileAppCategorySingularDataSource, "")
)

var mobileAppCategoryResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // mobileAppCategory
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The key of the entity. This property is read-only.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The name of the app category.",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The date and time the mobileAppCategory was last modified. This property is read-only.  \nProvider Note: Warning: This attribute seems to always return the _current_ time for mobile app categories that have been created for the tenant (i.e. that have not been predefined by Microsoft). Therefore it can be expected to change with every query.",
		},
	},
	MarkdownDescription: "Contains properties for a single Intune app category. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileappcategory?view=graph-rest-beta ||| MS Graph: App management",
}
