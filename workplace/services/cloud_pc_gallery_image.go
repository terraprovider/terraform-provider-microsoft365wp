package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	cloudPcGalleryImageResource = generic.GenericResource{
		TypeNameSuffix: "cloud_pc_gallery_image",
		SpecificSchema: cloudPcGalleryImageResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:     "/deviceManagement/virtualEndpoint/galleryImages",
			ReadOptions: cloudPcGalleryImageReadOptions,
		},
	}

	CloudPcGalleryImageSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&cloudPcGalleryImageResource)

	CloudPcGalleryImagePluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&cloudPcGalleryImageResource, "")
)

var cloudPcGalleryImageReadOptions = generic.ReadOptions{
	DataSource: generic.DataSourceOptions{
		ExtraFilterAttributes: []string{"offer_name", "publisher_name", "status", "sku_name"},
		Plural: generic.PluralOptions{
			ExtraAttributes: []string{"offer_name", "publisher_name", "status", "sku_name", "start_date", "end_date", "expiration_date"},
		},
	},
}

var cloudPcGalleryImageResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // cloudPcGalleryImage
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier (ID) of the gallery image resource on Cloud PC. The ID format is {publisherName_offerName_skuName}. For example, `MicrosoftWindowsDesktop_windows-ent-cpc_win11-22h2-ent-cpc-m365`. You can find the **publisherName**, **offerName**, and **skuName** in the Azure Marketplace. Read-only.",
		},
		"display_name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The display name of this gallery image. For example, `Windows 11 Enterprise + Microsoft 365 Apps 22H2`. Read-only.",
		},
		"end_date": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The date when the status of image becomes `supportedWithWarning`. Users can still provision new Cloud PCs if the current time is later than **endDate** and earlier than **expirationDate**. For example, assume the **endDate** of a gallery image is `2023-9-14` and **expirationDate** is `2024-3-14`, users are able to provision new Cloud PCs if today is 2023-10-01. Read-only.",
		},
		"expiration_date": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The date when the image is no longer available. Users are unable to provision new Cloud PCs if the current time is later than **expirationDate**. The value is usually **endDate** plus six months. For example, if the **startDate** is `2025-10-14`, the **expirationDate** is usually `2026-04-14`. Read-only.",
		},
		"offer_name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The offer name of this gallery image that is passed to ARM to retrieve the image resource. Read-only.",
		},
		"os_version_number": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The operating system version of this gallery image. For example, `10.0.22000.296`. Read-only.",
		},
		"publisher_name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The publisher name of this gallery image that is passed to ARM to retrieve the image resource. Read-only.",
		},
		"size_in_gb": schema.Int64Attribute{
			Optional:            true,
			Description:         `sizeInGB`, // custom MS Graph attribute name
			MarkdownDescription: "Indicates the size of this image in gigabytes. For example, `64`. Read-only.",
		},
		"sku_name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The SKU name of this image that is passed to ARM to retrieve the image resource. Read-only.",
		},
		"start_date": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The date when the Cloud PC image is available for provisioning new Cloud PCs. For example, `2022-09-20`. Read-only.",
		},
		"status": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("supported", "supportedWithWarning", "notSupported", "unknownFutureValue"),
			},
			MarkdownDescription: "The status of the gallery image on the Cloud PC. The default value is `supported`. Read-only. / Possible values are: `supported`, `supportedWithWarning`, `notSupported`, `unknownFutureValue`",
		},
	},
	MarkdownDescription: "Represents the gallery image resource of the current organization that can be used to provision a Cloud PC. / https://learn.microsoft.com/en-us/graph/api/resources/cloudpcgalleryimage?view=graph-rest-beta ||| MS Graph: Cloud PC",
}
