package services

//
// Pre-alpha stage and meant to be used for a data source for reading only!
//

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	MobileAppResource = generic.GenericResource{
		TypeNameSuffix: "mobile_app",
		SpecificSchema: mobileAppResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:         "/deviceAppManagement/mobileApps",
			ReadOptions:     mobileAppReadOptions,
			WriteSubActions: mobileAppWriteSubActions,
		},
	}

	MobileAppSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&MobileAppResource)

	MobileAppPluralDataSource = generic.CreateGenericDataSourcePluralFromSingular(
		&MobileAppSingularDataSource, "")
)

var mobileAppReadOptions = generic.ReadOptions{
	ODataExpand: "categories,assignments",
}

var mobileAppWriteSubActions = []generic.WriteSubAction{
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			Attributes: []string{"assignments"},
			UriSuffix:  "assign",
			UpdateOnly: true,
		},
	},
}

var mobileAppResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // mobileApp
		"id": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"created_date_time": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"dependent_app_count": schema.Int64Attribute{
			Computed:      true,
			PlanModifiers: []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
		},
		"description": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: `The description of the app.`,
		},
		"developer": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: `The developer of the app.`,
		},
		"display_name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: `The admin provided or imported title of the app.`,
		},
		"information_url": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: `The more information Url.`,
		},
		"is_assigned": schema.BoolAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
		},
		"is_featured": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: `The value indicating whether the app is marked as featured by the admin.`,
		},
		"large_icon": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // mimeContent
				"type": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: `Indicates the content mime type.`,
				},
				"value": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: `The byte array that contains the actual content.`,
				},
			},
			MarkdownDescription: `The large icon, to be displayed in the app details and used for upload of the icon. / Contains properties for a generic mime content. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mimeContent?view=graph-rest-beta`,
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"notes": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: `Notes for the app.`,
		},
		"owner": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: `The owner of the app.`,
		},
		"privacy_information_url": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: `The privacy statement Url.`,
		},
		"publisher": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: `The publisher of the app.`,
		},
		"publishing_state": schema.StringAttribute{
			Computed: true,
			Validators: []validator.String{
				stringvalidator.OneOf("notPublished", "processing", "published"),
			},
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: `List of scope tag ids for this mobile app.`,
		},
		"superseded_app_count": schema.Int64Attribute{
			Computed:      true,
			PlanModifiers: []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
		},
		"superseding_app_count": schema.Int64Attribute{
			Computed:      true,
			PlanModifiers: []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
		},
		"upload_state": schema.Int64Attribute{
			Computed:      true,
			PlanModifiers: []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
		},
		"assignments": deviceAndAppManagementAssignment,
		"categories": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // mobileAppCategory
					"id": schema.StringAttribute{
						Required: true,
					},
					"display_name": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: `The name of the app category.`,
					},
					"last_modified_date_time": schema.StringAttribute{
						Computed:      true,
						PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
					},
				},
			},
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: `The list of categories for this app. / Contains properties for a single Intune app category.`,
		},
		"android_managed_store": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.androidManagedStoreApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // androidManagedStoreApp
					"app_identifier": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The Identity Name.`,
					},
					"app_store_url": schema.StringAttribute{
						Optional: true,
					},
					"app_tracks": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // androidManagedStoreAppTrack
								"track_alias": schema.StringAttribute{
									Optional: true,
								},
								"track_id": schema.StringAttribute{
									Optional: true,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Contains track information for Android Managed Store apps.`,
					},
					"is_private": schema.BoolAttribute{
						Optional:      true,
						PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:      true,
					},
					"is_system_app": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Indicates whether the app is a preinstalled system app.`,
					},
					"package_id": schema.StringAttribute{
						Optional: true,
					},
					"supports_oem_config": schema.BoolAttribute{
						Optional:      true,
						PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:      true,
					},
					"total_license_count": schema.Int64Attribute{
						Optional:      true,
						PlanModifiers: []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
						Computed:      true,
					},
					"used_license_count": schema.Int64Attribute{
						Optional:      true,
						PlanModifiers: []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
						Computed:      true,
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: `Contains properties and inherited properties for Android Managed Store Apps.`,
			},
		},
		"ios_store": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.iosStoreApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // iosStoreApp
					"applicable_device_type": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // iosDeviceType
							"ipad": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `iPad`, // custom MS Graph attribute name
								MarkdownDescription: `Whether the app should run on iPads.`,
							},
							"iphone_and_ipod": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `iPhoneAndIPod`, // custom MS Graph attribute name
								MarkdownDescription: `Whether the app should run on iPhones and iPods.`,
							},
						},
						PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `The iOS architecture for which this app can run on. / Contains properties of the possible iOS device types the mobile app can run on.`,
					},
					"app_store_url": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The Apple App Store URL`,
					},
					"bundle_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The Identity Name.`,
					},
					"minimum_supported_operating_system": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // iosMinimumOperatingSystem
							"v10_0": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_0`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, only Version 10.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.`,
							},
							"v11_0": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v11_0`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, only Version 11.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.`,
							},
							"v12_0": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v12_0`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, only Version 12.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.`,
							},
							"v13_0": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v13_0`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, only Version 13.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.`,
							},
							"v14_0": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v14_0`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, only Version 14.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.`,
							},
							"v15_0": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v15_0`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, only Version 15.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.`,
							},
							"v16_0": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v16_0`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, only Version 16.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.`,
							},
							"v17_0": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v17_0`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, only Version 17.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.`,
							},
							"v8_0": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v8_0`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, only Version 8.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.`,
							},
							"v9_0": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v9_0`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, only Version 9.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.`,
							},
						},
						MarkdownDescription: `The value for the minimum applicable operating system. / Contains properties of the minimum operating system required for an iOS mobile app.`,
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: `Contains properties and inherited properties for iOS store apps.`,
			},
		},
		"ios_vpp": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.iosVppApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // iosVppApp
					"applicable_device_type": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // iosDeviceType
							"ipad": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `iPad`, // custom MS Graph attribute name
								MarkdownDescription: `Whether the app should run on iPads.`,
							},
							"iphone_and_ipod": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `iPhoneAndIPod`, // custom MS Graph attribute name
								MarkdownDescription: `Whether the app should run on iPhones and iPods.`,
							},
						},
						MarkdownDescription: `The applicable iOS Device Type. / Contains properties of the possible iOS device types the mobile app can run on.`,
					},
					"app_store_url": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The store URL.`,
					},
					"bundle_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The Identity Name.`,
					},
					"licensing_type": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // vppLicensingType
							"support_device_licensing": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: `Whether the program supports the device licensing type.`,
							},
							"supports_device_licensing": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: `Whether the program supports the device licensing type.`,
							},
							"supports_user_licensing": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: `Whether the program supports the user licensing type.`,
							},
							"support_user_licensing": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: `Whether the program supports the user licensing type.`,
							},
						},
						MarkdownDescription: `The supported License Type. / Contains properties for iOS Volume-Purchased Program (Vpp) Licensing Type.`,
					},
					"release_date_time": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The VPP application release date and time.`,
					},
					"revoke_license_action_results": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // iosVppAppRevokeLicensesActionResult
								"action_failure_reason": schema.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.OneOf("none", "appleFailure", "internalError", "expiredVppToken", "expiredApplePushNotificationCertificate"),
									},
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
									Computed:            true,
									MarkdownDescription: `The reason for the revoke licenses action failure.`,
								},
								"action_name": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `Action name`,
								},
								"action_state": schema.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.OneOf("none", "pending", "canceled", "active", "done", "failed", "notSupported"),
									},
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
									Computed:            true,
									MarkdownDescription: `State of the action`,
								},
								"failed_licenses_count": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: `A count of the number of licenses for which revoke failed.`,
								},
								"last_updated_date_time": schema.StringAttribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
									Computed:            true,
									MarkdownDescription: `Time the action state was last updated`,
								},
								"managed_device_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `DeviceId associated with the action.`,
								},
								"start_date_time": schema.StringAttribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
									Computed:            true,
									MarkdownDescription: `Time the action was initiated`,
								},
								"total_licenses_count": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: `A count of the number of licenses for which revoke was attempted.`,
								},
								"user_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `UserId associated with the action.`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Results of revoke license actions on this app. / Defines results for actions on iOS Vpp Apps, contains inherited properties for ActionResult.`,
					},
					"total_license_count": schema.Int64Attribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
						Computed:            true,
						MarkdownDescription: `The total number of VPP licenses.`,
					},
					"used_license_count": schema.Int64Attribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
						Computed:            true,
						MarkdownDescription: `The number of VPP licenses in use.`,
					},
					"vpp_token_account_type": schema.StringAttribute{
						Optional:            true,
						Validators:          []validator.String{stringvalidator.OneOf("business", "education")},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("business")},
						Computed:            true,
						MarkdownDescription: `The type of volume purchase program which the given Apple Volume Purchase Program Token is associated with. Possible values are: 'business', 'education'.`,
					},
					"vpp_token_apple_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The Apple Id associated with the given Apple Volume Purchase Program Token.`,
					},
					"vpp_token_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Identifier of the VPP token associated with this app.`,
					},
					"vpp_token_organization_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The organization associated with the Apple Volume Purchase Program Token`,
					},
					"assigned_licenses": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // iosVppAppAssignedLicense
								"id": schema.StringAttribute{
									Optional:      true,
									PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
									Computed:      true,
								},
								"user_email_address": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `The user email address.`,
								},
								"user_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `The user ID.`,
								},
								"user_name": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `The user name.`,
								},
								"user_principal_name": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `The user principal name.`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `The licenses assigned to this app. / iOS Volume Purchase Program license assignment. This class does not support Create, Delete, or Update.`,
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: `Contains properties and inherited properties for iOS Volume-Purchased Program (VPP) Apps.`,
			},
		},
		"macos_dmg": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.macOSDmgApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // macOSDmgApp
					"ignore_version_detection": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `When TRUE, indicates that the app's version will NOT be used to detect if the app is installed on a device. When FALSE, indicates that the app's version will be used to detect if the app is installed on a device. Set this to true for apps that use a self update feature. The default value is FALSE.`,
					},
					"included_apps": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // macOSIncludedApp
								"bundle_id": schema.StringAttribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
									Computed:            true,
									MarkdownDescription: `The bundleId of the app. This maps to the CFBundleIdentifier in the app's bundle configuration.`,
								},
								"bundle_version": schema.StringAttribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
									Computed:            true,
									MarkdownDescription: `The version of the app. This maps to the CFBundleShortVersion in the app's bundle configuration.`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `The list of .apps expected to be installed by the DMG (Apple Disk Image). This collection can contain a maximum of 500 elements. / Contains properties of an included .app in a MacOS app.`,
					},
					"minimum_supported_operating_system": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // macOSMinimumOperatingSystem
							"v10_10": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_10`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, indicates OS X 10.10 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.`,
							},
							"v10_11": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_11`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, indicates OS X 10.11 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.`,
							},
							"v10_12": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_12`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, indicates macOS 10.12 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.`,
							},
							"v10_13": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_13`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, indicates macOS 10.13 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.`,
							},
							"v10_14": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_14`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, indicates macOS 10.14 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.`,
							},
							"v10_15": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_15`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, indicates macOS 10.15 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.`,
							},
							"v10_7": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_7`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, indicates Mac OS X 10.7 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.`,
							},
							"v10_8": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_8`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, indicates OS X 10.8 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.`,
							},
							"v10_9": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_9`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, indicates OS X 10.9 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.`,
							},
							"v11_0": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v11_0`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, indicates macOS 11.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.`,
							},
							"v12_0": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v12_0`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, indicates macOS 12.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.`,
							},
							"v13_0": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v13_0`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, indicates macOS 13.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.`,
							},
							"v14_0": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v14_0`, // custom MS Graph attribute name
								MarkdownDescription: `When TRUE, indicates macOS 14.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.`,
							},
						},
						MarkdownDescription: `ComplexType macOSMinimumOperatingSystem that indicates the minimum operating system applicable for the application. / The minimum operating system required for a macOS app.`,
					},
					"primary_bundle_id": schema.StringAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
						Computed:            true,
						MarkdownDescription: `The bundleId of the primary .app in the DMG (Apple Disk Image). This maps to the CFBundleIdentifier in the app's bundle configuration.`,
					},
					"primary_bundle_version": schema.StringAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
						Computed:            true,
						MarkdownDescription: `The version of the primary .app in the DMG (Apple Disk Image). This maps to the CFBundleShortVersion in the app's bundle configuration.`,
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: `Contains properties and inherited properties for the MacOS DMG (Apple Disk Image) App.`,
			},
		},
		"macos_ms_defender": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.macOSMicrosoftDefenderApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional:   true,
				Attributes: map[string]schema.Attribute{ // macOSMicrosoftDefenderApp
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: `Contains properties and inherited properties for the macOS Microsoft Defender App.`,
			},
		},
		"macos_ms_edge": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.macOSMicrosoftEdgeApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // macOSMicrosoftEdgeApp
					"channel": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("dev", "beta", "stable", "unknownFutureValue"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("stable")},
						Computed:            true,
						MarkdownDescription: `The channel to install on target devices.`,
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: `Contains properties and inherited properties for the macOS Microsoft Edge App.`,
			},
		},
		"macos_office_suite": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.macOSOfficeSuiteApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional:   true,
				Attributes: map[string]schema.Attribute{ // macOSOfficeSuiteApp
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: `Contains properties and inherited properties for the MacOS Office Suite App.`,
			},
		},
		"macos_vpp": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.macOsVppApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // macOsVppApp
					"app_store_url": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The store URL.`,
					},
					"bundle_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The Identity Name.`,
					},
					"licensing_type": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // vppLicensingType
							"support_device_licensing": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: `Whether the program supports the device licensing type.`,
							},
							"supports_device_licensing": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: `Whether the program supports the device licensing type.`,
							},
							"supports_user_licensing": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: `Whether the program supports the user licensing type.`,
							},
							"support_user_licensing": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: `Whether the program supports the user licensing type.`,
							},
						},
						MarkdownDescription: `The supported License Type. / Contains properties for iOS Volume-Purchased Program (Vpp) Licensing Type.`,
					},
					"release_date_time": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The VPP application release date and time.`,
					},
					"revoke_license_action_results": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // macOsVppAppRevokeLicensesActionResult
								"action_failure_reason": schema.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.OneOf("none", "appleFailure", "internalError", "expiredVppToken", "expiredApplePushNotificationCertificate"),
									},
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
									Computed:            true,
									MarkdownDescription: `The reason for the revoke licenses action failure.`,
								},
								"action_name": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `Action name`,
								},
								"action_state": schema.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.OneOf("none", "pending", "canceled", "active", "done", "failed", "notSupported"),
									},
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
									Computed:            true,
									MarkdownDescription: `State of the action`,
								},
								"failed_licenses_count": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: `A count of the number of licenses for which revoke failed.`,
								},
								"last_updated_date_time": schema.StringAttribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
									Computed:            true,
									MarkdownDescription: `Time the action state was last updated`,
								},
								"managed_device_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `DeviceId associated with the action.`,
								},
								"start_date_time": schema.StringAttribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
									Computed:            true,
									MarkdownDescription: `Time the action was initiated`,
								},
								"total_licenses_count": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: `A count of the number of licenses for which revoke was attempted.`,
								},
								"user_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `UserId associated with the action.`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `Results of revoke license actions on this app. / Defines results for actions on MacOS Vpp Apps, contains inherited properties for ActionResult.`,
					},
					"total_license_count": schema.Int64Attribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
						Computed:            true,
						MarkdownDescription: `The total number of VPP licenses.`,
					},
					"used_license_count": schema.Int64Attribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
						Computed:            true,
						MarkdownDescription: `The number of VPP licenses in use.`,
					},
					"vpp_token_account_type": schema.StringAttribute{
						Optional:            true,
						Validators:          []validator.String{stringvalidator.OneOf("business", "education")},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("business")},
						Computed:            true,
						MarkdownDescription: `The type of volume purchase program which the given Apple Volume Purchase Program Token is associated with. Possible values are: 'business', 'education'.`,
					},
					"vpp_token_apple_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The Apple Id associated with the given Apple Volume Purchase Program Token.`,
					},
					"vpp_token_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Identifier of the VPP token associated with this app.`,
					},
					"vpp_token_organization_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The organization associated with the Apple Volume Purchase Program Token`,
					},
					"assigned_licenses": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // macOsVppAppAssignedLicense
								"id": schema.StringAttribute{
									Optional:      true,
									PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
									Computed:      true,
								},
								"user_email_address": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `The user email address.`,
								},
								"user_id": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `The user ID.`,
								},
								"user_name": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `The user name.`,
								},
								"user_principal_name": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `The user principal name.`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `The licenses assigned to this app. / MacOS Volume Purchase Program license assignment. This class does not support Create, Delete, or Update.`,
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: `Contains properties and inherited properties for MacOS Volume-Purchased Program (VPP) Apps.`,
			},
		},
		"windows_msi": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.windowsMobileMSI",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // windowsMobileMSI
					"command_line": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The command line.`,
					},
					"identity_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The identity version.`,
					},
					"ignore_version_detection": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `A boolean to control whether the app's version will be used to detect the app after it is installed on a device. Set this to true for Windows Mobile MSI Line of Business (LoB) apps that use a self update feature.`,
					},
					"product_code": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The product code.`,
					},
					"product_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The product version of Windows Mobile MSI Line of Business (LoB) app.`,
					},
					"use_device_context": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: `Indicates whether to install a dual-mode MSI in the device context. If true, app will be installed for all users. If false, app will be installed per-user. If null, service will use the MSI package's default install context. In case of dual-mode MSI, this default will be per-user.  Cannot be set for non-dual-mode apps.  Cannot be changed after initial creation of the application.`,
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: `Contains properties and inherited properties for Windows Mobile MSI Line Of Business apps.`,
			},
		},
		"windows_win32lob": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.win32LobApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // win32LobApp
					"allow_available_uninstall": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `When TRUE, indicates that uninstall is supported from the company portal for the Windows app (Win32) with an Available assignment. When FALSE, indicates that uninstall is not supported for the Windows app (Win32) with an Available assignment. Default value is FALSE.`,
					},
					"applicable_architectures": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							wpvalidator.FlagEnumValues("none", "x86", "x64", "arm", "neutral", "arm64"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
						Computed:            true,
						MarkdownDescription: `The Windows architecture(s) for which this app can run on.`,
					},
					"detection_rules": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // win32LobAppDetection
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `The detection rules to detect Win32 Line of Business (LoB) app. / Base class to detect a Win32 App`,
					},
					"display_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The version displayed in the UX for this app.`,
					},
					"install_command_line": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The command line to install this app`,
					},
					"install_experience": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // win32LobAppInstallExperience
							"device_restart_behavior": schema.StringAttribute{
								Optional: true,
								Validators: []validator.String{
									stringvalidator.OneOf("basedOnReturnCode", "allow", "suppress", "force"),
								},
								PlanModifiers: []planmodifier.String{
									wpdefaultvalue.StringDefaultValue("basedOnReturnCode"),
								},
								Computed:            true,
								MarkdownDescription: `Device restart behavior.`,
							},
							"max_run_time_in_minutes": schema.Int64Attribute{
								Optional:            true,
								MarkdownDescription: `The number of minutes the system will wait for install program to finish. Default value is 60 minutes.`,
							},
							"run_as_account": schema.StringAttribute{
								Optional:            true,
								Validators:          []validator.String{stringvalidator.OneOf("system", "user")},
								PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("system")},
								Computed:            true,
								MarkdownDescription: `Indicates the type of execution context the app runs in.`,
							},
						},
						MarkdownDescription: `The install experience for this app. / Contains installation experience properties for a Win32 App`,
					},
					"minimum_cpu_speed_in_m_hz": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `The value for the minimum CPU speed which is required to install this app.`,
					},
					"minimum_free_disk_space_in_m_b": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `The value for the minimum free disk space which is required to install this app.`,
					},
					"minimum_memory_in_m_b": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `The value for the minimum physical memory which is required to install this app.`,
					},
					"minimum_number_of_processors": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: `The value for the minimum number of processors which is required to install this app.`,
					},
					"minimum_supported_operating_system": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // windowsMinimumOperatingSystem
							"v10_0": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_0`, // custom MS Graph attribute name
								MarkdownDescription: `Windows version 10.0 or later.`,
							},
							"v10_1607": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_1607`, // custom MS Graph attribute name
								MarkdownDescription: `Windows 10 1607 or later.`,
							},
							"v10_1703": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_1703`, // custom MS Graph attribute name
								MarkdownDescription: `Windows 10 1703 or later.`,
							},
							"v10_1709": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_1709`, // custom MS Graph attribute name
								MarkdownDescription: `Windows 10 1709 or later.`,
							},
							"v10_1803": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_1803`, // custom MS Graph attribute name
								MarkdownDescription: `Windows 10 1803 or later.`,
							},
							"v10_1809": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_1809`, // custom MS Graph attribute name
								MarkdownDescription: `Windows 10 1809 or later.`,
							},
							"v10_1903": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_1903`, // custom MS Graph attribute name
								MarkdownDescription: `Windows 10 1903 or later.`,
							},
							"v10_1909": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_1909`, // custom MS Graph attribute name
								MarkdownDescription: `Windows 10 1909 or later.`,
							},
							"v10_2004": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_2004`, // custom MS Graph attribute name
								MarkdownDescription: `Windows 10 2004 or later.`,
							},
							"v10_21_h1": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_21H1`, // custom MS Graph attribute name
								MarkdownDescription: `Windows 10 21H1 or later.`,
							},
							"v10_2_h20": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v10_2H20`, // custom MS Graph attribute name
								MarkdownDescription: `Windows 10 2H20 or later.`,
							},
							"v8_0": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v8_0`, // custom MS Graph attribute name
								MarkdownDescription: `Windows version 8.0 or later.`,
							},
							"v8_1": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `v8_1`, // custom MS Graph attribute name
								MarkdownDescription: `Windows version 8.1 or later.`,
							},
						},
						MarkdownDescription: `The value for the minimum applicable operating system. / The minimum operating system required for a Windows mobile app.`,
					},
					"minimum_supported_windows_release": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The value for the minimum supported windows release.`,
					},
					"msi_information": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // win32LobAppMsiInformation
							"package_type": schema.StringAttribute{
								Optional: true,
								Validators: []validator.String{
									stringvalidator.OneOf("perMachine", "perUser", "dualPurpose"),
								},
								PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("perMachine")},
								Computed:            true,
								MarkdownDescription: `The MSI package type.`,
							},
							"product_code": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: `The MSI product code.`,
							},
							"product_name": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: `The MSI product name.`,
							},
							"product_version": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: `The MSI product version.`,
							},
							"publisher": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: `The MSI publisher.`,
							},
							"requires_reboot": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: `Whether the MSI app requires the machine to reboot to complete installation.`,
							},
							"upgrade_code": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: `The MSI upgrade code.`,
							},
						},
						MarkdownDescription: `The MSI details if this Win32 app is an MSI app. / Contains MSI app properties for a Win32 App.`,
					},
					"requirement_rules": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // win32LobAppRequirement
								"detection_value": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `The detection value`,
								},
								"operator": schema.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										wpvalidator.FlagEnumValues("notConfigured", "equal", "notEqual", "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual"),
									},
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
									Computed:            true,
									MarkdownDescription: `The operator for detection`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `The requirement rules to detect Win32 Line of Business (LoB) app. / Base class to detect a Win32 App`,
					},
					"return_codes": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // win32LobAppReturnCode
								"return_code": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: `Return code.`,
								},
								"type": schema.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.OneOf("failed", "success", "softReboot", "hardReboot", "retry"),
									},
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("failed")},
									Computed:            true,
									MarkdownDescription: `The type of return code.`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `The return codes for post installation behavior. / Contains return code properties for a Win32 App`,
					},
					"rules": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // win32LobAppRule
								"rule_type": schema.StringAttribute{
									Optional:            true,
									Validators:          []validator.String{stringvalidator.OneOf("detection", "requirement")},
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("detection")},
									Computed:            true,
									MarkdownDescription: `The rule type indicating the purpose of the rule.`,
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: `The detection and requirement rules for this app. / A base complex type to store the detection or requirement rule data for a Win32 LOB app.`,
					},
					"setup_file_path": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The relative path of the setup file in the encrypted Win32LobApp package.`,
					},
					"uninstall_command_line": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The command line to uninstall this app`,
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: `Contains properties and inherited properties for Win32 apps.`,
			},
		},
		"windows_winget": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.winGetApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // winGetApp
					"install_experience": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // winGetAppInstallExperience
							"run_as_account": schema.StringAttribute{
								Optional:            true,
								Validators:          []validator.String{stringvalidator.OneOf("system", "user")},
								PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("system")},
								Computed:            true,
								MarkdownDescription: `Indicates the type of execution context the app setup runs in on target devices. Options include values of the RunAsAccountType enum, which are System and User. Required at creation time, cannot be modified on existing objects.`,
							},
						},
						MarkdownDescription: `The install experience settings associated with this application, which are used to ensure the desired install experiences on the target device are taken into account. This includes the account type (System or User) that actions should be run as on target devices. Required at creation time. / Represents the install experience settings associated with WinGet apps. This is used to ensure the desired install experiences on the target device are taken into account. Required at creation time.`,
					},
					"manifest_hash": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Hash of package metadata properties used to validate that the application matches the metadata in the source repository.`,
					},
					"package_identifier": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `The PackageIdentifier from the WinGet source repository REST API. This also maps to the Id when using the WinGet client command line application. Required at creation time, cannot be modified on existing objects.`,
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: `A MobileApp that is based on a referenced application in a WinGet repository.`,
			},
		},
	},
	MarkdownDescription: `An abstract class containing the base properties for Intune mobile apps. Note: Listing mobile apps with '$expand=assignments' has been deprecated. Instead get the list of apps without the '$expand' query on 'assignments'. Then, perform the expansion on individual applications. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mobileApp?view=graph-rest-beta`,
}

var mobileAppMobileAppValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("android_managed_store"),
	path.MatchRelative().AtParent().AtName("ios_store"),
	path.MatchRelative().AtParent().AtName("ios_vpp"),
	path.MatchRelative().AtParent().AtName("macos_dmg"),
	path.MatchRelative().AtParent().AtName("macos_ms_defender"),
	path.MatchRelative().AtParent().AtName("macos_ms_edge"),
	path.MatchRelative().AtParent().AtName("macos_office_suite"),
	path.MatchRelative().AtParent().AtName("macos_vpp"),
	path.MatchRelative().AtParent().AtName("windows_msi"),
	path.MatchRelative().AtParent().AtName("windows_win32lob"),
	path.MatchRelative().AtParent().AtName("windows_winget"),
)
