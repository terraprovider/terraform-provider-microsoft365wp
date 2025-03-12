package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	cloudPcDeviceImageResource = generic.GenericResource{
		TypeNameSuffix: "cloud_pc_device_image",
		SpecificSchema: cloudPcDeviceImageResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:     "/deviceManagement/virtualEndpoint/deviceImages",
			ReadOptions: cloudPcDeviceImageReadOptions,
		},
	}

	CloudPcDeviceImageSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&cloudPcDeviceImageResource)

	CloudPcDeviceImagePluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&cloudPcDeviceImageResource, "")
)

var cloudPcDeviceImageReadOptions = generic.ReadOptions{
	DataSource: generic.DataSourceOptions{
		ExtraFilterAttributes: []string{"operating_system", "os_build_number", "os_status", "status"},
		Plural: generic.PluralOptions{
			ExtraAttributes: []string{"expiration_date", "operating_system", "os_build_number", "os_status", "os_version_number", "status", "version"},
		},
	},
}

var cloudPcDeviceImageResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // cloudPcDeviceImage
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier (ID) of the image resource on the Cloud PC. Read-only.",
		},
		"display_name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The display name of the associated device image. The device image display name and the version are used to uniquely identify the Cloud PC device image. Read-only.",
		},
		"error_code": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("internalServerError", "sourceImageNotFound", "osVersionNotSupported", "sourceImageInvalid", "sourceImageNotGeneralized", "unknownFutureValue", "vmAlreadyAzureAdjoined", "paidSourceImageNotSupport", "sourceImageNotSupportCustomizeVMName", "sourceImageSizeExceedsLimitation", "sourceImageWithDataDiskNotSupported", "sourceImageWithDiskEncryptionSetNotSupported"),
			},
			MarkdownDescription: "The error code of the status of the image that indicates why the upload failed, if applicable. Use the `Prefer: include-unknown-enum-members` request header to get the following values from this [evolvable enum](/graph/best-practices-concept#handling-future-members-in-evolvable-enumerations): `vmAlreadyAzureAdJoined`, `paidSourceImageNotSupport`, `sourceImageNotSupportCustomizeVMName`, `sourceImageSizeExceedsLimitation`, `sourceImageWithDataDiskNotSupported`, `sourceImageWithDiskEncryptionSetNotSupported`. Read-only. / Possible values are: `internalServerError`, `sourceImageNotFound`, `osVersionNotSupported`, `sourceImageInvalid`, `sourceImageNotGeneralized`, `unknownFutureValue`, `vmAlreadyAzureAdjoined`, `paidSourceImageNotSupport`, `sourceImageNotSupportCustomizeVMName`, `sourceImageSizeExceedsLimitation`, `sourceImageWithDataDiskNotSupported`, `sourceImageWithDiskEncryptionSetNotSupported`",
		},
		"expiration_date": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The date when the image became unavailable. Read-only.",
		},
		"last_modified_date_time": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The data and time when the image was last modified. The timestamp represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Read-only.",
		},
		"operating_system": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The operating system of the image. For example, `Windows 10 Enterprise`. Read-only.",
		},
		"os_build_number": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The OS build version of the image. For example, `1909`. Read-only.",
		},
		"os_status": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("supported", "supportedWithWarning", "unknown", "unknownFutureValue"),
			},
			MarkdownDescription: "The OS status of this image. The default value is `unknown`. Read-only. / Possible values are: `supported`, `supportedWithWarning`, `unknown`, `unknownFutureValue`",
		},
		"os_version_number": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The operating system version of this image. For example, `10.0.22000.296`. Read-only.",
		},
		"scope_ids": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
		},
		"source_image_resource_id": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The unique identifier (ID) of the source image resource on Azure. The required ID format is: \"/subscriptions/{subscription-id}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/images/{imageName}\". Read-only.",
		},
		"status": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("pending", "ready", "failed", "unknownFutureValue", "warning"),
			},
			MarkdownDescription: "The status of the image on the Cloud PC. Read-only. / Possible values are: `pending`, `ready`, `failed`, `unknownFutureValue`, `warning`",
		},
		"status_details": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("internalServerError", "sourceImageNotFound", "osVersionNotSupported", "sourceImageInvalid", "sourceImageNotGeneralized", "unknownFutureValue", "vmAlreadyAzureAdjoined", "paidSourceImageNotSupport", "sourceImageNotSupportCustomizeVMName", "sourceImageSizeExceedsLimitation"),
			},
			MarkdownDescription: "Possible values are: `internalServerError`, `sourceImageNotFound`, `osVersionNotSupported`, `sourceImageInvalid`, `sourceImageNotGeneralized`, `unknownFutureValue`, `vmAlreadyAzureAdjoined`, `paidSourceImageNotSupport`, `sourceImageNotSupportCustomizeVMName`, `sourceImageSizeExceedsLimitation`",
		},
		"version": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The image version. For example, `0.0.1` and `1.5.13`. Read-only.",
		},
	},
	MarkdownDescription: "Represents the image resource on a Cloud PC. / https://learn.microsoft.com/en-us/graph/api/resources/cloudpcdeviceimage?view=graph-rest-beta ||| MS Graph: Cloud PC",
}
