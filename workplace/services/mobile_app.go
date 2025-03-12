package services

import (
	"context"
	"terraform-provider-microsoft365wp/workplace/generic"
	mobileapputils "terraform-provider-microsoft365wp/workplace/services/mobile_app_utils"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	MobileAppResource = generic.GenericResource{
		TypeNameSuffix: "mobile_app",
		SpecificSchema: mobileAppResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:                    "/deviceAppManagement/mobileApps",
			ReadOptions:                mobileAppReadOptions,
			WriteSubActions:            mobileAppWriteSubActions,
			TerraformToGraphMiddleware: mobileAppTerraformToGraphMiddleware,
		},
	}

	MobileAppSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&MobileAppResource)

	MobileAppPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&MobileAppResource, "")
)

var mobileAppReadOptions = generic.ReadOptions{
	ODataExpand: "categories,assignments",
	ExtraRequestsCustom: []generic.ReadExtraRequestCustom{
		mobileapputils.CheckContentStateInGraphRerc,
	},
}

var mobileAppWriteSubActions = []generic.WriteSubAction{
	&mobileapputils.CheckPublishingStateWsa{MessageSuffix: "after write", AllowNotPublished: true},
	&mobileapputils.WriteContentWsa{},
	&mobileapputils.CheckPublishingStateWsa{MessageSuffix: "after content"},
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			AttributesMap: map[string]string{"assignments": "mobileAppAssignments"},
			UriSuffix:     "assign",
		},
	},
	&generic.WriteSubActionIndividual{
		WriteSubActionBase: generic.WriteSubActionBase{
			Attributes: []string{"categories"},
			UriSuffix:  "categories",
		},
		ComparisonKeyAttribute: "id",
		SetNestedPath:          tftypes.NewAttributePath().WithAttributeName("categories"),
		IsOdataReference:       true,
		OdataRefMapTypeToUriPrefix: map[string]string{
			"": "https://graph.microsoft.com/beta/deviceAppManagement/mobileAppCategories/",
		},
	},
}

func mobileAppTerraformToGraphMiddleware(ctx context.Context, diags *diag.Diagnostics, params *generic.TerraformToGraphMiddlewareParams) generic.TerraformToGraphMiddlewareReturns {
	if params.IsUpdate {
		// remove unchangeable attributes from PATCH
		delete(params.RawVal, "appIdentifier")
		delete(params.RawVal, "appStoreUrl")
		delete(params.RawVal, "appUrl")
		delete(params.RawVal, "channel")
		delete(params.RawVal, "installExperience")
		delete(params.RawVal, "isSystemApp")
		delete(params.RawVal, "manifestHash")
		delete(params.RawVal, "packageIdentifier")
		if params.RawVal["@odata.type"].(string) == "#microsoft.graph.windowsMobileMSI" {
			delete(params.RawVal, "useDeviceContext")
		}
		if params.RawVal["@odata.type"].(string) == "#microsoft.graph.officeSuiteApp" {
			// "read-only" - they get set to "Microsoft" by the Intune Portal and shall probably just stay that way...
			delete(params.RawVal, "developer")
			delete(params.RawVal, "owner")
			delete(params.RawVal, "publisher")
			delete(params.RawVal, "officePlatformArchitecture")
			if params.RawVal["officeConfigurationXml"] == nil {
				delete(params.RawVal, "officeConfigurationXml")
			}
		}
		if params.RawVal["@odata.type"].(string) == "#microsoft.graph.windowsUniversalAppX" {
			delete(params.RawVal, "isBundle")
		}
	}
	return nil
}

var mobileAppResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // mobileApp
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Key of the entity.",
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The date and time the app was created.",
		},
		"dependent_app_count": schema.Int64Attribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
			MarkdownDescription: "The total number of dependencies the child app has.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "The description of the app. The _provider_ default value is `\"\"`.",
		},
		"developer": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The developer of the app.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The admin provided or imported title of the app.",
		},
		"information_url": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The more information Url.",
		},
		"is_assigned": schema.BoolAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
			MarkdownDescription: "The value indicating whether the app is assigned to at least one group.",
		},
		"is_featured": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "The value indicating whether the app is marked as featured by the admin. The _provider_ default value is `false`.",
		},
		"large_icon": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // mimeContent
				"type": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "Indicates the content mime type.",
				},
				"value_base64": schema.StringAttribute{
					Required:            true,
					Description:         `value`, // custom MS Graph attribute name
					MarkdownDescription: "The byte array that contains the actual content.",
				},
			},
			MarkdownDescription: "The large icon, to be displayed in the app details and used for upload of the icon. / Contains properties for a generic mime content. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mimecontent?view=graph-rest-beta",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The date and time the app was last modified.",
		},
		"notes": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Notes for the app.",
		},
		"owner": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The owner of the app.",
		},
		"privacy_information_url": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The privacy statement Url.",
		},
		"publisher": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The publisher of the app.",
		},
		"publishing_state": schema.StringAttribute{
			Computed: true,
			Validators: []validator.String{
				stringvalidator.OneOf("notPublished", "processing", "published"),
			},
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The publishing state for the app. The app cannot be assigned unless the app is published. / Indicates the publishing state of an app; possible values are: `notPublished` (The app is not yet published.), `processing` (The app is pending service-side processing.), `published` (The app is published.)",
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: "List of scope tag ids for this mobile app. The _provider_ default value is `[\"0\"]`.",
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
			Computed:            true,
			PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
			MarkdownDescription: "The upload state.",
		},
		"assignments": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: mobileAppMobileAppAssignmentAttributes,
			},
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "The list of group assignments for this mobile app. / A class containing the properties used for Group Assignment of a Mobile App. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileappassignment?view=graph-rest-beta. The _provider_ default value is `[]`.",
		},
		"categories": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // mobileAppCategory
					"id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The key of the entity. This property is read-only.",
					},
					"display_name": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The name of the app category.",
					},
				},
			},
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "The list of categories for this app. / Contains properties for a single Intune app category. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileappcategory?view=graph-rest-beta. The _provider_ default value is `[]`.",
		},
		"android_for_work": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.androidForWorkApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{ // androidForWorkApp
					"app_identifier": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The Identity Name. This property is read-only.",
					},
					"app_store_url": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The Play for Work Store app URL.",
					},
					"package_id": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The package identifier. This property is read-only.",
					},
					"total_license_count": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "The total number of VPP licenses.",
					},
					"used_license_count": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "The number of VPP licenses in use.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				PlanModifiers:       []planmodifier.Object{wpplanmodifier.ObjectUseStateForUnknown()},
				MarkdownDescription: "Contains properties and inherited properties for Android for Work (AFW) Apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidforworkapp?view=graph-rest-beta",
			},
		},
		"android_lob": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.androidLobApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{ // androidLobApp
					"committed_content_version": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The internal committed content version.",
					},
					"file_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The name of the main Lob application file.",
					},
					"size": schema.Int64Attribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
						MarkdownDescription: "The total size, including all uploaded files. This property is read-only.",
					},
					"minimum_supported_operating_system": schema.SingleNestedAttribute{
						Computed:            true,
						Attributes:          mobileAppAndroidMinimumOperatingSystemAttributes,
						MarkdownDescription: "The value for the minimum applicable operating system. / Contains properties for the minimum operating system required for an Android mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidminimumoperatingsystem?view=graph-rest-beta",
					},
					"package_id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The package identifier.",
					},
					"targeted_platforms": schema.StringAttribute{
						Computed: true,
						Validators: []validator.String{
							wpvalidator.FlagEnumValues("androidDeviceAdministrator", "androidOpenSourceProject", "unknownFutureValue"),
						},
						MarkdownDescription: "The platforms to which the application can be targeted. If not specified, will defauilt to Android Device Administrator. / Specifies which platform(s) can be targeted for a given Android LOB application or Managed Android LOB application; possible values are: `androidDeviceAdministrator` (Indicates the Android targeted platform is Android Device Administrator.), `androidOpenSourceProject` (Indicates the Android targeted platform is Android Open Source Project.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.)",
					},
					"version_code": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The version code of Android Line of Business (LoB) app.",
					},
					"version_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The version name of Android Line of Business (LoB) app.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				PlanModifiers:       []planmodifier.Object{wpplanmodifier.ObjectUseStateForUnknown()},
				MarkdownDescription: "Contains properties and inherited properties for Android Line Of Business apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidlobapp?view=graph-rest-beta",
			},
		},
		"android_managed_store": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.androidManagedStoreApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // androidManagedStoreApp
					"app_identifier": schema.StringAttribute{
						Required:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
						MarkdownDescription: "The Identity Name.",
					},
					"app_store_url": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The Play for Work Store app URL. This property is read-only.",
					},
					"app_tracks": schema.SetNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // androidManagedStoreAppTrack
								"track_alias": schema.StringAttribute{
									Computed:            true,
									PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
									MarkdownDescription: "Friendly name for track. This property is read-only.",
								},
								"track_id": schema.StringAttribute{
									Computed:            true,
									PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
									MarkdownDescription: "Unique track identifier. This property is read-only.",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpplanmodifier.SetUseStateForUnknown()},
						MarkdownDescription: "The tracks that are visible to this enterprise. This property is read-only. / Contains track information for Android Managed Store apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidmanagedstoreapptrack?view=graph-rest-beta",
					},
					"is_private": schema.BoolAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
						MarkdownDescription: "Indicates whether the app is only available to a given enterprise's users. This property is read-only.",
					},
					"is_system_app": schema.BoolAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.Bool{
							wpdefaultvalue.BoolDefaultValue(true),
							boolplanmodifier.RequiresReplace(),
						},
						Computed:            true,
						MarkdownDescription: "Indicates whether the app is a preinstalled system app. The _provider_ default value is `true`.",
					},
					"package_id": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The package identifier. This property is read-only.",
					},
					"supports_oem_config": schema.BoolAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
						MarkdownDescription: "Whether this app supports OEMConfig policy. This property is read-only.",
					},
					"total_license_count": schema.Int64Attribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
						MarkdownDescription: "The total number of VPP licenses. This property is read-only.",
					},
					"used_license_count": schema.Int64Attribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
						MarkdownDescription: "The number of VPP licenses in use. This property is read-only.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for Android Managed Store Apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidmanagedstoreapp?view=graph-rest-beta",
			},
		},
		"android_store": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.androidStoreApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // androidStoreApp
					"app_store_url": schema.StringAttribute{
						Required:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
						MarkdownDescription: "The Android app store URL.",
					},
					"minimum_supported_operating_system": schema.SingleNestedAttribute{
						Required:            true,
						Attributes:          mobileAppAndroidMinimumOperatingSystemAttributes,
						MarkdownDescription: "The value for the minimum applicable operating system. / Contains properties for the minimum operating system required for an Android mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidminimumoperatingsystem?view=graph-rest-beta",
					},
					"package_id": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The package identifier. This property is read-only.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for Android store apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidstoreapp?view=graph-rest-beta",
			},
		},
		"ios_lob": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.iosLobApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{ // iosLobApp
					"committed_content_version": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The internal committed content version.",
					},
					"file_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The name of the main Lob application file.",
					},
					"size": schema.Int64Attribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
						MarkdownDescription: "The total size, including all uploaded files. This property is read-only.",
					},
					"applicable_device_type": schema.SingleNestedAttribute{
						Computed:            true,
						Attributes:          mobileAppIosDeviceTypeAttributes,
						MarkdownDescription: "The iOS architecture for which this app can run on. / Contains properties of the possible iOS device types the mobile app can run on. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosdevicetype?view=graph-rest-beta",
					},
					"build_number": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The build number of iOS Line of Business (LoB) app.",
					},
					"bundle_id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The Identity Name.",
					},
					"expiration_date_time": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The expiration time.",
					},
					"minimum_supported_operating_system": schema.SingleNestedAttribute{
						Computed:            true,
						Attributes:          mobileAppIosMinimumOperatingSystemAttributes,
						MarkdownDescription: "The value for the minimum applicable operating system. / Contains properties of the minimum operating system required for an iOS mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosminimumoperatingsystem?view=graph-rest-beta",
					},
					"version_number": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The version number of iOS Line of Business (LoB) app.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				PlanModifiers:       []planmodifier.Object{wpplanmodifier.ObjectUseStateForUnknown()},
				MarkdownDescription: "Contains properties and inherited properties for iOS Line Of Business apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-ioslobapp?view=graph-rest-beta",
			},
		},
		"ios_store": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.iosStoreApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // iosStoreApp
					"applicable_device_type": schema.SingleNestedAttribute{
						Optional:            true,
						Attributes:          mobileAppIosDeviceTypeAttributes,
						PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "The iOS architecture for which this app can run on. / Contains properties of the possible iOS device types the mobile app can run on. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosdevicetype?view=graph-rest-beta. The _provider_ default value is `{}`.",
					},
					"app_store_url": schema.StringAttribute{
						Required:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
						MarkdownDescription: "The Apple App Store URL",
					},
					"bundle_id": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The Identity Name.",
					},
					"minimum_supported_operating_system": schema.SingleNestedAttribute{
						Required:            true,
						Attributes:          mobileAppIosMinimumOperatingSystemAttributes,
						MarkdownDescription: "The value for the minimum applicable operating system. / Contains properties of the minimum operating system required for an iOS mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosminimumoperatingsystem?view=graph-rest-beta",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for iOS store apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosstoreapp?view=graph-rest-beta",
			},
		},
		"ios_vpp": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.iosVppApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{ // iosVppApp
					"applicable_device_type": schema.SingleNestedAttribute{
						Computed:            true,
						Attributes:          mobileAppIosDeviceTypeAttributes,
						MarkdownDescription: "The applicable iOS Device Type. / Contains properties of the possible iOS device types the mobile app can run on. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosdevicetype?view=graph-rest-beta",
					},
					"app_store_url": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The store URL.",
					},
					"bundle_id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The Identity Name.",
					},
					"licensing_type": schema.SingleNestedAttribute{
						Computed:            true,
						Attributes:          mobileAppVppLicensingTypeAttributes,
						MarkdownDescription: "The supported License Type. / Contains properties for iOS Volume-Purchased Program (Vpp) Licensing Type. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-vpplicensingtype?view=graph-rest-beta",
					},
					"release_date_time": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The VPP application release date and time.",
					},
					"revoke_license_action_results": schema.SetNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // iosVppAppRevokeLicensesActionResult
								"action_failure_reason": schema.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.OneOf("none", "appleFailure", "internalError", "expiredVppToken", "expiredApplePushNotificationCertificate"),
									},
									MarkdownDescription: "The reason for the revoke licenses action failure. / Possible types of reasons for an Apple Volume Purchase Program token action failure; possible values are: `none` (None.), `appleFailure` (There was an error on Apple's service.), `internalError` (There was an internal error.), `expiredVppToken` (There was an error because the Apple Volume Purchase Program token was expired.), `expiredApplePushNotificationCertificate` (There was an error because the Apple Volume Purchase Program Push Notification certificate expired.)",
								},
								"action_name": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Action name",
								},
								"action_state": schema.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.OneOf("none", "pending", "canceled", "active", "done", "failed", "notSupported"),
									},
									MarkdownDescription: "State of the action. / State of the action on the device; possible values are: `none` (Not a valid action state), `pending` (Action is pending), `canceled` (Action has been cancelled.), `active` (Action is active.), `done` (Action completed without errors.), `failed` (Action failed), `notSupported` (Action is not supported.)",
								},
								"failed_licenses_count": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "A count of the number of licenses for which revoke failed.",
								},
								"last_updated_date_time": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Time the action state was last updated",
								},
								"managed_device_id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "DeviceId associated with the action.",
								},
								"start_date_time": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Time the action was initiated",
								},
								"total_licenses_count": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "A count of the number of licenses for which revoke was attempted.",
								},
								"user_id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "UserId associated with the action.",
								},
							},
						},
						MarkdownDescription: "Results of revoke license actions on this app. / Defines results for actions on iOS Vpp Apps, contains inherited properties for ActionResult. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosvppapprevokelicensesactionresult?view=graph-rest-beta",
					},
					"total_license_count": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "The total number of VPP licenses.",
					},
					"used_license_count": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "The number of VPP licenses in use.",
					},
					"vpp_token_account_type": schema.StringAttribute{
						Computed:            true,
						Validators:          []validator.String{stringvalidator.OneOf("business", "education")},
						MarkdownDescription: "The type of volume purchase program which the given Apple Volume Purchase Program Token is associated with. / Possible types of an Apple Volume Purchase Program token; possible values are: `business` (Apple Volume Purchase Program token associated with an business program.), `education` (Apple Volume Purchase Program token associated with an education program.)",
					},
					"vpp_token_apple_id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The Apple Id associated with the given Apple Volume Purchase Program Token.",
					},
					"vpp_token_display_name": schema.StringAttribute{
						Computed: true,
					},
					"vpp_token_id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Identifier of the VPP token associated with this app.",
					},
					"vpp_token_organization_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The organization associated with the Apple Volume Purchase Program Token",
					},
					"assigned_licenses": schema.SetNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // iosVppAppAssignedLicense
								"id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Key of the entity. This property is read-only.",
								},
								"user_email_address": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The user email address.",
								},
								"user_id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The user ID.",
								},
								"user_name": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The user name.",
								},
								"user_principal_name": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The user principal name.",
								},
							},
						},
						MarkdownDescription: "The licenses assigned to this app. / iOS Volume Purchase Program license assignment. This class does not support Create, Delete, or Update. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosvppappassignedlicense?view=graph-rest-beta",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				PlanModifiers:       []planmodifier.Object{wpplanmodifier.ObjectUseStateForUnknown()},
				MarkdownDescription: "Contains properties and inherited properties for iOS Volume-Purchased Program (VPP) Apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosvppapp?view=graph-rest-beta",
			},
		},
		"ios_webclip": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.iosiPadOSWebClip",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // iosiPadOSWebClip
					"app_url": schema.StringAttribute{
						Required:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
						MarkdownDescription: "Indicates iOS/iPadOS web clip app URL. Example: \"https://www.contoso.com\"",
					},
					"full_screen_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not to open the web clip as a full-screen web app. Defaults to false. If TRUE, opens the web clip as a full-screen web app. If FALSE, the web clip opens inside of another app, such as Safari or the app specified with targetApplicationBundleIdentifier.",
					},
					"ignore_manifest_scope": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not a full screen web clip can navigate to an external web site without showing the Safari UI. Defaults to false. If FALSE, the Safari UI appears when navigating away. If TRUE, the Safari UI will not be shown.",
					},
					"pre_composed_icon_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not the icon for the app is precomosed. Defaults to false. If TRUE, prevents SpringBoard from adding \"shine\" to the icon. If FALSE, SpringBoard can add \"shine\".",
					},
					"target_application_bundle_identifier": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Specifies the application bundle identifier which opens the URL. Available in iOS 14 and later.",
					},
					"use_managed_browser": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether or not to use managed browser. When TRUE, the app will be required to be opened in Microsoft Edge. When FALSE, the app will not be required to be opened in Microsoft Edge. By default, this property is set to FALSE. The _provider_ default value is `false`.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for iOS web apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosipadoswebclip?view=graph-rest-beta",
			},
		},
		"macos_dmg": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.macOSDmgApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // macOSDmgApp
					"committed_content_version": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The internal committed content version.",
					},
					"file_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The name of the main Lob application file.",
					},
					"size": schema.Int64Attribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
						MarkdownDescription: "The total size, including all uploaded files. This property is read-only.",
					},
					"ignore_version_detection": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "When TRUE, indicates that the app's version will NOT be used to detect if the app is installed on a device. When FALSE, indicates that the app's version will be used to detect if the app is installed on a device. Set this to true for apps that use a self update feature. The default value is FALSE. The _provider_ default value is `false`.",
					},
					"included_apps": schema.SetNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: mobileAppMacOSIncludedAppAttributes,
						},
						Validators:          []validator.Set{setvalidator.SizeAtLeast(1)},
						MarkdownDescription: "The list of .apps expected to be installed by the DMG (Apple Disk Image). This collection can contain a maximum of 500 elements. / Contains properties of an included .app in a MacOS app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosincludedapp?view=graph-rest-beta",
					},
					"minimum_supported_operating_system": schema.SingleNestedAttribute{
						Required:            true,
						Attributes:          mobileAppMacOSMinimumOperatingSystemAttributes,
						MarkdownDescription: "ComplexType macOSMinimumOperatingSystem that indicates the minimum operating system applicable for the application. / The minimum operating system required for a macOS app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosminimumoperatingsystem?view=graph-rest-beta",
					},
					"primary_bundle_id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The bundleId of the primary .app in the DMG (Apple Disk Image). This maps to the CFBundleIdentifier in the app's bundle configuration.",
					},
					"primary_bundle_version": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The version of the primary .app in the DMG (Apple Disk Image). This maps to the CFBundleShortVersion in the app's bundle configuration.",
					},
					"source_file": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Provider Note: The path to the DMG file to be uploaded.",
					},
					"source_sha256": schema.StringAttribute{
						Required:            true,
						Description:         `sourceSha256`, // custom MS Graph attribute name
						MarkdownDescription: "Provider Note: The SHA-256 sum of the source file (see example). Changing this attribute will trigger a new upload of the source file.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for the MacOS DMG (Apple Disk Image) App. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosdmgapp?view=graph-rest-beta",
			},
		},
		"macos_lob": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.macOSLobApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{ // macOSLobApp
					"committed_content_version": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The internal committed content version.",
					},
					"file_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The name of the main Lob application file.",
					},
					"size": schema.Int64Attribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
						MarkdownDescription: "The total size, including all uploaded files. This property is read-only.",
					},
					"build_number": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The build number of the package. This should match the package CFBundleShortVersionString of the .pkg file.",
					},
					"bundle_id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The primary bundleId of the package.",
					},
					"child_apps": schema.SetNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // macOSLobChildApp
								"build_number": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The build number of the app.",
								},
								"bundle_id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The bundleId of the app.",
								},
								"version_number": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The version number of the app.",
								},
							},
						},
						MarkdownDescription: "List of ComplexType macOSLobChildApp objects. Represents the apps expected to be installed by the package. / Contains properties of a macOS .app in the package / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macoslobchildapp?view=graph-rest-beta",
					},
					"ignore_version_detection": schema.BoolAttribute{
						Computed:            true,
						MarkdownDescription: "When TRUE, indicates that the app's version will NOT be used to detect if the app is installed on a device. When FALSE, indicates that the app's version will be used to detect if the app is installed on a device. Set this to true for apps that use a self update feature. The default value is FALSE.",
					},
					"install_as_managed": schema.BoolAttribute{
						Computed:            true,
						MarkdownDescription: "When TRUE, indicates that the app will be installed as managed (requires macOS 11.0 and other managed package restrictions). When FALSE, indicates that the app will be installed as unmanaged. The default value is FALSE.",
					},
					"md5_hash": schema.SetAttribute{
						ElementType:         types.StringType,
						Computed:            true,
						Description:         `md5Hash`, // custom MS Graph attribute name
						MarkdownDescription: "The MD5 hash codes. This is empty if the package was uploaded directly. If the Intune App Wrapping Tool is used to create a .intunemac, this value can be found inside the Detection.xml file.",
					},
					"md5_hash_chunk_size": schema.Int64Attribute{
						Computed:            true,
						Description:         `md5HashChunkSize`, // custom MS Graph attribute name
						MarkdownDescription: "The chunk size for MD5 hash. This is '0' or empty if the package was uploaded directly. If the Intune App Wrapping Tool is used to create a .intunemac, this value can be found inside the Detection.xml file.",
					},
					"minimum_supported_operating_system": schema.SingleNestedAttribute{
						Computed:            true,
						Attributes:          mobileAppMacOSMinimumOperatingSystemAttributes,
						MarkdownDescription: "ComplexType macOSMinimumOperatingSystem that indicates the minimum operating system applicable for the application. / The minimum operating system required for a macOS app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosminimumoperatingsystem?view=graph-rest-beta",
					},
					"version_number": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The version number of the package. This should match the package CFBundleVersion in the packageinfo file.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				PlanModifiers:       []planmodifier.Object{wpplanmodifier.ObjectUseStateForUnknown()},
				MarkdownDescription: "Contains properties and inherited properties for the macOS LOB App. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macoslobapp?view=graph-rest-beta",
			},
		},
		"macos_ms_defender": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.macOSMicrosoftDefenderApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional:   true,
				Attributes: map[string]schema.Attribute{ // macOSMicrosoftDefenderApp
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for the macOS Microsoft Defender App. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosmicrosoftdefenderapp?view=graph-rest-beta",
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
						PlanModifiers: []planmodifier.String{
							wpdefaultvalue.StringDefaultValue("stable"),
							stringplanmodifier.RequiresReplace(),
						},
						Computed:            true,
						MarkdownDescription: "The channel to install on target devices. / The enum to specify the channels for Microsoft Edge apps; possible values are: `dev` (The Dev Channel is intended to help you plan and develop with the latest capabilities of Microsoft Edge.), `beta` (The Beta Channel is intended for production deployment to a representative sample set of users. New features ship about every 4 weeks. Security and quality updates ship as needed.), `stable` (The Stable Channel is intended for broad deployment within organizations, and it's the channel that most users should be on. New features ship about every 4 weeks. Security and quality updates ship as needed.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.). The _provider_ default value is `\"stable\"`.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for the macOS Microsoft Edge App. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosmicrosoftedgeapp?view=graph-rest-beta",
			},
		},
		"macos_office_suite": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.macOSOfficeSuiteApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional:   true,
				Attributes: map[string]schema.Attribute{ // macOSOfficeSuiteApp
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for the MacOS Office Suite App. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosofficesuiteapp?view=graph-rest-beta",
			},
		},
		"macos_pkg": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.macOSPkgApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // macOSPkgApp
					"committed_content_version": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The internal committed content version.",
					},
					"file_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The name of the main Lob application file.",
					},
					"size": schema.Int64Attribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
						MarkdownDescription: "The total size, including all uploaded files. This property is read-only.",
					},
					"ignore_version_detection": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "When TRUE, indicates that the app's version will NOT be used to detect if the app is installed on a device. When FALSE, indicates that the app's version will be used to detect if the app is installed on a device. Set this to true for apps that use a self update feature. The default value is FALSE. The _provider_ default value is `false`.",
					},
					"included_apps": schema.SetNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: mobileAppMacOSIncludedAppAttributes,
						},
						Validators:          []validator.Set{setvalidator.SizeAtLeast(1)},
						MarkdownDescription: "The list of apps expected to be installed by the PKG. This collection can contain a maximum of 500 elements. / Contains properties of an included .app in a MacOS app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosincludedapp?view=graph-rest-beta",
					},
					"minimum_supported_operating_system": schema.SingleNestedAttribute{
						Required:            true,
						Attributes:          mobileAppMacOSMinimumOperatingSystemAttributes,
						MarkdownDescription: "ComplexType macOSMinimumOperatingSystem that indicates the minimum operating system applicable for the application. / The minimum operating system required for a macOS app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosminimumoperatingsystem?view=graph-rest-beta",
					},
					"post_install_script": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // macOSAppScript
							"script_content": schema.StringAttribute{
								Optional:            true,
								Validators:          []validator.String{wpvalidator.TranslateValueEncodeBase64()},
								PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
								Computed:            true,
								MarkdownDescription: "The base64 encoded shell script (.sh) that assists managing macOS apps. The _provider_ default value is `\"\"`.",
							},
						},
						MarkdownDescription: "ComplexType macOSAppScript the contains the post-install script for the app. This will execute on the macOS device after the app is installed. / Shell script used to assist installation of a macOS app. These scripts are used to perform additional tasks to help the app successfully be configured. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosappscript?view=graph-rest-beta",
					},
					"pre_install_script": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // macOSAppScript
							"script_content": schema.StringAttribute{
								Optional:            true,
								Validators:          []validator.String{wpvalidator.TranslateValueEncodeBase64()},
								PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
								Computed:            true,
								MarkdownDescription: "The base64 encoded shell script (.sh) that assists managing macOS apps. The _provider_ default value is `\"\"`.",
							},
						},
						MarkdownDescription: "ComplexType macOSAppScript the contains the post-install script for the app. This will execute on the macOS device after the app is installed. / Shell script used to assist installation of a macOS app. These scripts are used to perform additional tasks to help the app successfully be configured. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosappscript?view=graph-rest-beta",
					},
					"primary_bundle_id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The bundleId of the primary app in the PKG. This maps to the CFBundleIdentifier in the app's bundle configuration.",
					},
					"primary_bundle_version": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The version of the primary app in the PKG. This maps to the CFBundleShortVersion in the app's bundle configuration.",
					},
					"source_file": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Provider Note: The path to the PKG file to be uploaded.",
					},
					"source_sha256": schema.StringAttribute{
						Required:            true,
						Description:         `sourceSha256`, // custom MS Graph attribute name
						MarkdownDescription: "Provider Note: The SHA-256 sum of the source file (see example). Changing this attribute will trigger a new upload of the source file.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for the MacOSPkgApp. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macospkgapp?view=graph-rest-beta",
			},
		},
		"macos_vpp": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.macOsVppApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{ // macOsVppApp
					"app_store_url": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The store URL.",
					},
					"bundle_id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The Identity Name.",
					},
					"licensing_type": schema.SingleNestedAttribute{
						Computed:            true,
						Attributes:          mobileAppVppLicensingTypeAttributes,
						MarkdownDescription: "The supported License Type. / Contains properties for iOS Volume-Purchased Program (Vpp) Licensing Type. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-vpplicensingtype?view=graph-rest-beta",
					},
					"release_date_time": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The VPP application release date and time.",
					},
					"revoke_license_action_results": schema.SetNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // macOsVppAppRevokeLicensesActionResult
								"action_failure_reason": schema.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.OneOf("none", "appleFailure", "internalError", "expiredVppToken", "expiredApplePushNotificationCertificate"),
									},
									MarkdownDescription: "The reason for the revoke licenses action failure. / Possible types of reasons for an Apple Volume Purchase Program token action failure; possible values are: `none` (None.), `appleFailure` (There was an error on Apple's service.), `internalError` (There was an internal error.), `expiredVppToken` (There was an error because the Apple Volume Purchase Program token was expired.), `expiredApplePushNotificationCertificate` (There was an error because the Apple Volume Purchase Program Push Notification certificate expired.)",
								},
								"action_name": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Action name",
								},
								"action_state": schema.StringAttribute{
									Computed: true,
									Validators: []validator.String{
										stringvalidator.OneOf("none", "pending", "canceled", "active", "done", "failed", "notSupported"),
									},
									MarkdownDescription: "State of the action. / State of the action on the device; possible values are: `none` (Not a valid action state), `pending` (Action is pending), `canceled` (Action has been cancelled.), `active` (Action is active.), `done` (Action completed without errors.), `failed` (Action failed), `notSupported` (Action is not supported.)",
								},
								"failed_licenses_count": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "A count of the number of licenses for which revoke failed.",
								},
								"last_updated_date_time": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Time the action state was last updated",
								},
								"managed_device_id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "DeviceId associated with the action.",
								},
								"start_date_time": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Time the action was initiated",
								},
								"total_licenses_count": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "A count of the number of licenses for which revoke was attempted.",
								},
								"user_id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "UserId associated with the action.",
								},
							},
						},
						MarkdownDescription: "Results of revoke license actions on this app. / Defines results for actions on MacOS Vpp Apps, contains inherited properties for ActionResult. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosvppapprevokelicensesactionresult?view=graph-rest-beta",
					},
					"total_license_count": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "The total number of VPP licenses.",
					},
					"used_license_count": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "The number of VPP licenses in use.",
					},
					"vpp_token_account_type": schema.StringAttribute{
						Computed:            true,
						Validators:          []validator.String{stringvalidator.OneOf("business", "education")},
						MarkdownDescription: "The type of volume purchase program which the given Apple Volume Purchase Program Token is associated with. / Possible types of an Apple Volume Purchase Program token; possible values are: `business` (Apple Volume Purchase Program token associated with an business program.), `education` (Apple Volume Purchase Program token associated with an education program.)",
					},
					"vpp_token_apple_id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The Apple Id associated with the given Apple Volume Purchase Program Token.",
					},
					"vpp_token_display_name": schema.StringAttribute{
						Computed: true,
					},
					"vpp_token_id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Identifier of the VPP token associated with this app.",
					},
					"vpp_token_organization_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The organization associated with the Apple Volume Purchase Program Token",
					},
					"assigned_licenses": schema.SetNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // macOsVppAppAssignedLicense
								"id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Key of the entity. This property is read-only.",
								},
								"user_email_address": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The user email address.",
								},
								"user_id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The user ID.",
								},
								"user_name": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The user name.",
								},
								"user_principal_name": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The user principal name.",
								},
							},
						},
						MarkdownDescription: "The licenses assigned to this app. / MacOS Volume Purchase Program license assignment. This class does not support Create, Delete, or Update. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosvppappassignedlicense?view=graph-rest-beta",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				PlanModifiers:       []planmodifier.Object{wpplanmodifier.ObjectUseStateForUnknown()},
				MarkdownDescription: "Contains properties and inherited properties for MacOS Volume-Purchased Program (VPP) Apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosvppapp?view=graph-rest-beta",
			},
		},
		"macos_webclip": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.macOSWebClip",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // macOSWebClip
					"app_url": schema.StringAttribute{
						Required:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
						MarkdownDescription: "The web app URL starting with http:// or https://, such as https://learn.microsoft.com/mem/.",
					},
					"full_screen_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not to open the web clip as a full-screen web app. Defaults to false. If TRUE, opens the web clip as a full-screen web app. If FALSE, the web clip opens inside of another app.",
					},
					"pre_composed_icon_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether or not the icon for the app is precomosed. Defaults to false. If TRUE, prevents SpringBoard from adding \"shine\" to the icon. If FALSE, SpringBoard can add \"shine\".",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for macOS web apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macoswebclip?view=graph-rest-beta",
			},
		},
		"managed_android_lob": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.managedAndroidLobApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{ // managedAndroidLobApp
					"app_availability": schema.StringAttribute{
						Computed:            true,
						Validators:          []validator.String{stringvalidator.OneOf("global", "lineOfBusiness")},
						MarkdownDescription: "The Application's availability. This property is read-only. / A managed (MAM) application's availability; possible values are: `global` (A globally available app to all tenants.), `lineOfBusiness` (A line of business apps private to an organization.)",
					},
					"version": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The Application's version.",
					},
					"committed_content_version": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The internal committed content version.",
					},
					"file_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The name of the main Lob application file.",
					},
					"size": schema.Int64Attribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
						MarkdownDescription: "The total size, including all uploaded files. This property is read-only.",
					},
					"minimum_supported_operating_system": schema.SingleNestedAttribute{
						Computed:            true,
						Attributes:          mobileAppAndroidMinimumOperatingSystemAttributes,
						MarkdownDescription: "The value for the minimum applicable operating system. / Contains properties for the minimum operating system required for an Android mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidminimumoperatingsystem?view=graph-rest-beta",
					},
					"package_id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The package identifier.",
					},
					"targeted_platforms": schema.StringAttribute{
						Computed: true,
						Validators: []validator.String{
							wpvalidator.FlagEnumValues("androidDeviceAdministrator", "androidOpenSourceProject", "unknownFutureValue"),
						},
						MarkdownDescription: "The platforms to which the application can be targeted. If not specified, will defauilt to Android Device Administrator. / Specifies which platform(s) can be targeted for a given Android LOB application or Managed Android LOB application; possible values are: `androidDeviceAdministrator` (Indicates the Android targeted platform is Android Device Administrator.), `androidOpenSourceProject` (Indicates the Android targeted platform is Android Open Source Project.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.)",
					},
					"version_code": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The version code of managed Android Line of Business (LoB) app.",
					},
					"version_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The version name of managed Android Line of Business (LoB) app.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				PlanModifiers:       []planmodifier.Object{wpplanmodifier.ObjectUseStateForUnknown()},
				MarkdownDescription: "Contains properties and inherited properties for Managed Android Line Of Business apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-managedandroidlobapp?view=graph-rest-beta",
			},
		},
		"managed_android_store": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.managedAndroidStoreApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // managedAndroidStoreApp
					"app_availability": schema.StringAttribute{
						Computed:            true,
						Validators:          []validator.String{stringvalidator.OneOf("global", "lineOfBusiness")},
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The Application's availability. This property is read-only. / A managed (MAM) application's availability; possible values are: `global` (A globally available app to all tenants.), `lineOfBusiness` (A line of business apps private to an organization.)",
					},
					"version": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The Application's version.",
					},
					"app_store_url": schema.StringAttribute{
						Required:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
						MarkdownDescription: "The Android AppStoreUrl.",
					},
					"minimum_supported_operating_system": schema.SingleNestedAttribute{
						Required:            true,
						Attributes:          mobileAppAndroidMinimumOperatingSystemAttributes,
						MarkdownDescription: "The value for the minimum supported operating system. / Contains properties for the minimum operating system required for an Android mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidminimumoperatingsystem?view=graph-rest-beta",
					},
					"package_id": schema.StringAttribute{
						Required:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
						MarkdownDescription: "The app's package ID.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for Android store apps that you can manage with an Intune app protection policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-managedandroidstoreapp?view=graph-rest-beta",
			},
		},
		"managed_ios_lob": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.managedIOSLobApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{ // managedIOSLobApp
					"app_availability": schema.StringAttribute{
						Computed:            true,
						Validators:          []validator.String{stringvalidator.OneOf("global", "lineOfBusiness")},
						MarkdownDescription: "The Application's availability. This property is read-only. / A managed (MAM) application's availability; possible values are: `global` (A globally available app to all tenants.), `lineOfBusiness` (A line of business apps private to an organization.)",
					},
					"version": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The Application's version.",
					},
					"committed_content_version": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The internal committed content version.",
					},
					"file_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The name of the main Lob application file.",
					},
					"size": schema.Int64Attribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
						MarkdownDescription: "The total size, including all uploaded files. This property is read-only.",
					},
					"applicable_device_type": schema.SingleNestedAttribute{
						Computed:            true,
						Attributes:          mobileAppIosDeviceTypeAttributes,
						MarkdownDescription: "The iOS architecture for which this app can run on. / Contains properties of the possible iOS device types the mobile app can run on. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosdevicetype?view=graph-rest-beta",
					},
					"build_number": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The build number of managed iOS Line of Business (LoB) app.",
					},
					"bundle_id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The Identity Name.",
					},
					"expiration_date_time": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The expiration time.",
					},
					"identity_version": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The identity version.",
					},
					"minimum_supported_operating_system": schema.SingleNestedAttribute{
						Computed:            true,
						Attributes:          mobileAppIosMinimumOperatingSystemAttributes,
						MarkdownDescription: "The value for the minimum applicable operating system. / Contains properties of the minimum operating system required for an iOS mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosminimumoperatingsystem?view=graph-rest-beta",
					},
					"version_number": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The version number of managed iOS Line of Business (LoB) app.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				PlanModifiers:       []planmodifier.Object{wpplanmodifier.ObjectUseStateForUnknown()},
				MarkdownDescription: "Contains properties and inherited properties for Managed iOS Line Of Business apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-managedioslobapp?view=graph-rest-beta",
			},
		},
		"managed_ios_store": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.managedIOSStoreApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // managedIOSStoreApp
					"app_availability": schema.StringAttribute{
						Computed:            true,
						Validators:          []validator.String{stringvalidator.OneOf("global", "lineOfBusiness")},
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The Application's availability. This property is read-only. / A managed (MAM) application's availability; possible values are: `global` (A globally available app to all tenants.), `lineOfBusiness` (A line of business apps private to an organization.)",
					},
					"version": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The Application's version.",
					},
					"applicable_device_type": schema.SingleNestedAttribute{
						Optional:            true,
						Attributes:          mobileAppIosDeviceTypeAttributes,
						PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "The iOS architecture for which this app can run on. / Contains properties of the possible iOS device types the mobile app can run on. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosdevicetype?view=graph-rest-beta. The _provider_ default value is `{}`.",
					},
					"app_store_url": schema.StringAttribute{
						Required:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
						MarkdownDescription: "The Apple AppStoreUrl.",
					},
					"bundle_id": schema.StringAttribute{
						Required:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
						MarkdownDescription: "The app's Bundle ID.",
					},
					"minimum_supported_operating_system": schema.SingleNestedAttribute{
						Required:            true,
						Attributes:          mobileAppIosMinimumOperatingSystemAttributes,
						MarkdownDescription: "The value for the minimum supported operating system. / Contains properties of the minimum operating system required for an iOS mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosminimumoperatingsystem?view=graph-rest-beta",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for an iOS store app that you can manage with an Intune app protection policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-managediosstoreapp?view=graph-rest-beta",
			},
		},
		"web_link": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.webApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // webApp
					"app_url": schema.StringAttribute{
						Required:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
						MarkdownDescription: "The web app URL. This property cannot be PATCHed.",
					},
					"use_managed_browser": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "Whether or not to use managed browser. This property is only applicable for Android and IOS. The _provider_ default value is `false`.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for web apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-webapp?view=graph-rest-beta",
			},
		},
		"win32_lob": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.win32LobApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // win32LobApp
					"committed_content_version": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The internal committed content version.",
					},
					"file_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The name of the main Lob application file.",
					},
					"size": schema.Int64Attribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
						MarkdownDescription: "The total size, including all uploaded files. This property is read-only.",
					},
					"allow_available_uninstall": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
						Computed:            true,
						MarkdownDescription: "When TRUE, indicates that uninstall is supported from the company portal for the Windows app (Win32) with an Available assignment. When FALSE, indicates that uninstall is not supported for the Windows app (Win32) with an Available assignment. Default value is FALSE. The _provider_ default value is `true`.",
					},
					"applicable_architectures": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							wpvalidator.FlagEnumValues("none", "x86", "x64", "arm", "neutral", "arm64"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("x86,x64")},
						Computed:            true,
						MarkdownDescription: "The Windows architecture(s) for which this app can run on. / Contains properties for Windows architecture; possible values are: `none` (No flags set.), `x86` (Whether or not the X86 Windows architecture type is supported.), `x64` (Whether or not the X64 Windows architecture type is supported.), `arm` (Whether or not the Arm Windows architecture type is supported.), `neutral` (Whether or not the Neutral Windows architecture type is supported.), `arm64` (Whether or not the Arm64 Windows architecture type is supported.). The _provider_ default value is `\"x86,x64\"`.",
					},
					"detection_rules": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // win32LobAppDetection
								"filesystem": generic.OdataDerivedTypeNestedAttributeRs{
									DerivedType: "#microsoft.graph.win32LobAppFileSystemDetection",
									SingleNestedAttribute: schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{ // win32LobAppFileSystemDetection
											"check_32_bit_on_64_system": schema.BoolAttribute{
												Optional:            true,
												PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
												Computed:            true,
												Description:         `check32BitOn64System`, // custom MS Graph attribute name
												MarkdownDescription: "A value indicating whether this file or folder is for checking 32-bit app on 64-bit system. The _provider_ default value is `false`.",
											},
											"detection_type": schema.StringAttribute{
												Required: true,
												Validators: []validator.String{
													stringvalidator.OneOf("notConfigured", "exists", "modifiedDate", "createdDate", "version", "sizeInMB", "doesNotExist"),
												},
												MarkdownDescription: "The file system detection type. / Contains all supported file system detection type; possible values are: `notConfigured` (Not configured.), `exists` (Whether the specified file or folder exists.), `modifiedDate` (Last modified date.), `createdDate` (Created date.), `version` (Version value type.), `sizeInMB` (Size detection type.), `doesNotExist` (The specified file or folder does not exist.)",
											},
											"detection_value": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "The file or folder detection value",
											},
											"file_or_folder_name": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "The file or folder name to detect Win32 Line of Business (LoB) app",
											},
											"operator": schema.StringAttribute{
												Optional: true,
												Validators: []validator.String{
													wpvalidator.FlagEnumValues("notConfigured", "equal", "notEqual", "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual"),
												},
												PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
												Computed:            true,
												MarkdownDescription: "The operator for file or folder detection. / Contains properties for detection operator; possible values are: `notConfigured` (Not configured.), `equal` (Equal operator.), `notEqual` (Not equal operator.), `greaterThan` (Greater than operator.), `greaterThanOrEqual` (Greater than or equal operator.), `lessThan` (Less than operator.), `lessThanOrEqual` (Less than or equal operator.). The _provider_ default value is `\"notConfigured\"`.",
											},
											"path": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "The file or folder path to detect Win32 Line of Business (LoB) app",
											},
										},
										Validators:          []validator.Object{mobileAppWin32LobAppDetectionValidator},
										MarkdownDescription: "Contains file or folder path to detect a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobappfilesystemdetection?view=graph-rest-beta",
									},
								},
								"powershell_script": generic.OdataDerivedTypeNestedAttributeRs{
									DerivedType: "#microsoft.graph.win32LobAppPowerShellScriptDetection",
									SingleNestedAttribute: schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{ // win32LobAppPowerShellScriptDetection
											"enforce_signature_check": schema.BoolAttribute{
												Optional:            true,
												PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
												Computed:            true,
												MarkdownDescription: "A value indicating whether signature check is enforced. The _provider_ default value is `false`.",
											},
											"run_as_32_bit": schema.BoolAttribute{
												Optional:            true,
												PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
												Computed:            true,
												Description:         `runAs32Bit`, // custom MS Graph attribute name
												MarkdownDescription: "A value indicating whether this script should run as 32-bit. The _provider_ default value is `false`.",
											},
											"script_content": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "The base64 encoded script content to detect Win32 Line of Business (LoB) app",
											},
										},
										Validators:          []validator.Object{mobileAppWin32LobAppDetectionValidator},
										MarkdownDescription: "Contains PowerShell script properties to detect a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobapppowershellscriptdetection?view=graph-rest-beta",
									},
								},
								"product_code": generic.OdataDerivedTypeNestedAttributeRs{
									DerivedType: "#microsoft.graph.win32LobAppProductCodeDetection",
									SingleNestedAttribute: schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{ // win32LobAppProductCodeDetection
											"product_code": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "The product code of Win32 Line of Business (LoB) app.",
											},
											"product_version": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "The product version of Win32 Line of Business (LoB) app.",
											},
											"product_version_operator": schema.StringAttribute{
												Optional: true,
												Validators: []validator.String{
													wpvalidator.FlagEnumValues("notConfigured", "equal", "notEqual", "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual"),
												},
												PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
												Computed:            true,
												MarkdownDescription: "The operator to detect product version. / Contains properties for detection operator; possible values are: `notConfigured` (Not configured.), `equal` (Equal operator.), `notEqual` (Not equal operator.), `greaterThan` (Greater than operator.), `greaterThanOrEqual` (Greater than or equal operator.), `lessThan` (Less than operator.), `lessThanOrEqual` (Less than or equal operator.). The _provider_ default value is `\"notConfigured\"`.",
											},
										},
										Validators:          []validator.Object{mobileAppWin32LobAppDetectionValidator},
										MarkdownDescription: "Contains product code and version properties to detect a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobappproductcodedetection?view=graph-rest-beta",
									},
								},
								"registry": generic.OdataDerivedTypeNestedAttributeRs{
									DerivedType: "#microsoft.graph.win32LobAppRegistryDetection",
									SingleNestedAttribute: schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{ // win32LobAppRegistryDetection
											"check_32_bit_on_64_system": schema.BoolAttribute{
												Optional:            true,
												PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
												Computed:            true,
												Description:         `check32BitOn64System`, // custom MS Graph attribute name
												MarkdownDescription: "A value indicating whether this registry path is for checking 32-bit app on 64-bit system. The _provider_ default value is `false`.",
											},
											"detection_type": schema.StringAttribute{
												Required: true,
												Validators: []validator.String{
													stringvalidator.OneOf("notConfigured", "exists", "doesNotExist", "string", "integer", "version"),
												},
												MarkdownDescription: "The registry data detection type. / Contains all supported registry data detection type; possible values are: `notConfigured` (Not configured.), `exists` (The specified registry key or value exists.), `doesNotExist` (The specified registry key or value does not exist.), `string` (String value type.), `integer` (Integer value type.), `version` (Version value type.)",
											},
											"detection_value": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "The registry detection value",
											},
											"key_path": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "The registry key path to detect Win32 Line of Business (LoB) app",
											},
											"operator": schema.StringAttribute{
												Optional: true,
												Validators: []validator.String{
													wpvalidator.FlagEnumValues("notConfigured", "equal", "notEqual", "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual"),
												},
												PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
												Computed:            true,
												MarkdownDescription: "The operator for registry data detection. / Contains properties for detection operator; possible values are: `notConfigured` (Not configured.), `equal` (Equal operator.), `notEqual` (Not equal operator.), `greaterThan` (Greater than operator.), `greaterThanOrEqual` (Greater than or equal operator.), `lessThan` (Less than operator.), `lessThanOrEqual` (Less than or equal operator.). The _provider_ default value is `\"notConfigured\"`.",
											},
											"value_name": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "The registry value name",
											},
										},
										Validators:          []validator.Object{mobileAppWin32LobAppDetectionValidator},
										MarkdownDescription: "Contains registry properties to detect a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobappregistrydetection?view=graph-rest-beta",
									},
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "The detection rules to detect Win32 Line of Business (LoB) app. / Base class to detect a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobappdetection?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"display_version": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The version displayed in the UX for this app.",
					},
					"install_command_line": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The command line to install this app",
					},
					"install_experience": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // win32LobAppInstallExperience
							"device_restart_behavior": schema.StringAttribute{
								Optional: true,
								Validators: []validator.String{
									stringvalidator.OneOf("basedOnReturnCode", "allow", "suppress", "force"),
								},
								PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("allow")},
								Computed:            true,
								MarkdownDescription: "Device restart behavior. / Indicates the type of restart action; possible values are: `basedOnReturnCode` (Intune will restart the device after running the app installation if the operation returns a reboot code.), `allow` (Intune will not take any specific action on reboot codes resulting from app installations. Intune will not attempt to suppress restarts for MSI apps.), `suppress` (Intune will attempt to suppress restarts for MSI apps.), `force` (Intune will force the device to restart immediately after the app installation operation.). The _provider_ default value is `\"allow\"`.",
							},
							"max_run_time_in_minutes": schema.Int64Attribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(60)},
								Computed:            true,
								MarkdownDescription: "The number of minutes the system will wait for install program to finish. Default value is 60 minutes. The _provider_ default value is `60`.",
							},
							"run_as_account": schema.StringAttribute{
								Optional:            true,
								Validators:          []validator.String{stringvalidator.OneOf("system", "user")},
								PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("system")},
								Computed:            true,
								MarkdownDescription: "Indicates the type of execution context the app runs in. / Indicates the type of execution context the app runs in; possible values are: `system` (System context), `user` (User context). The _provider_ default value is `\"system\"`.",
							},
						},
						PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "The install experience for this app. / Contains installation experience properties for a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobappinstallexperience?view=graph-rest-beta. The _provider_ default value is `{}`.",
					},
					"minimum_cpu_speed_in_mhz": schema.Int64Attribute{
						Optional:            true,
						Description:         `minimumCpuSpeedInMHz`, // custom MS Graph attribute name
						MarkdownDescription: "The value for the minimum CPU speed which is required to install this app.",
					},
					"minimum_free_disk_space_in_mb": schema.Int64Attribute{
						Optional:            true,
						Description:         `minimumFreeDiskSpaceInMB`, // custom MS Graph attribute name
						MarkdownDescription: "The value for the minimum free disk space which is required to install this app.",
					},
					"minimum_memory_in_mb": schema.Int64Attribute{
						Optional:            true,
						Description:         `minimumMemoryInMB`, // custom MS Graph attribute name
						MarkdownDescription: "The value for the minimum physical memory which is required to install this app.",
					},
					"minimum_number_of_processors": schema.Int64Attribute{
						Optional:            true,
						MarkdownDescription: "The value for the minimum number of processors which is required to install this app.",
					},
					"minimum_supported_windows_release": schema.StringAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("1809")},
						Computed:            true,
						MarkdownDescription: "The value for the minimum supported windows release. The _provider_ default value is `\"1809\"`.",
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
								MarkdownDescription: "The MSI package type. / Indicates the package type of an MSI Win32LobApp; possible values are: `perMachine` (Indicates a per-machine app package.), `perUser` (Indicates a per-user app package.), `dualPurpose` (Indicates a dual-purpose app package.). The _provider_ default value is `\"perMachine\"`.",
							},
							"product_code": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "The MSI product code.",
							},
							"product_name": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "The MSI product name.",
							},
							"product_version": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "The MSI product version.",
							},
							"publisher": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "The MSI publisher.",
							},
							"requires_reboot": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: "Whether the MSI app requires the machine to reboot to complete installation. The _provider_ default value is `false`.",
							},
							"upgrade_code": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "The MSI upgrade code.",
							},
						},
						MarkdownDescription: "The MSI details if this Win32 app is an MSI app. / Contains MSI app properties for a Win32 App. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobappmsiinformation?view=graph-rest-beta",
					},
					"requirement_rules": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // win32LobAppRequirement
								"filesystem": generic.OdataDerivedTypeNestedAttributeRs{
									DerivedType: "#microsoft.graph.win32LobAppFileSystemRequirement",
									SingleNestedAttribute: schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{ // win32LobAppFileSystemRequirement
											"detection_value": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "The detection value",
											},
											"operator": schema.StringAttribute{
												Optional: true,
												Validators: []validator.String{
													wpvalidator.FlagEnumValues("notConfigured", "equal", "notEqual", "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual"),
												},
												PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
												Computed:            true,
												MarkdownDescription: "The operator for detection / Contains properties for detection operator; possible values are: `notConfigured` (Not configured.), `equal` (Equal operator.), `notEqual` (Not equal operator.), `greaterThan` (Greater than operator.), `greaterThanOrEqual` (Greater than or equal operator.), `lessThan` (Less than operator.), `lessThanOrEqual` (Less than or equal operator.). The _provider_ default value is `\"notConfigured\"`.",
											},
											"check_32_bit_on_64_system": schema.BoolAttribute{
												Optional:            true,
												PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
												Computed:            true,
												Description:         `check32BitOn64System`, // custom MS Graph attribute name
												MarkdownDescription: "A value indicating whether this file or folder is for checking 32-bit app on 64-bit system. The _provider_ default value is `false`.",
											},
											"detection_type": schema.StringAttribute{
												Required: true,
												Validators: []validator.String{
													stringvalidator.OneOf("notConfigured", "exists", "modifiedDate", "createdDate", "version", "sizeInMB", "doesNotExist"),
												},
												MarkdownDescription: "The file system detection type. / Contains all supported file system detection type; possible values are: `notConfigured` (Not configured.), `exists` (Whether the specified file or folder exists.), `modifiedDate` (Last modified date.), `createdDate` (Created date.), `version` (Version value type.), `sizeInMB` (Size detection type.), `doesNotExist` (The specified file or folder does not exist.)",
											},
											"file_or_folder_name": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "The file or folder name to detect Win32 Line of Business (LoB) app",
											},
											"path": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "The file or folder path to detect Win32 Line of Business (LoB) app",
											},
										},
										Validators:          []validator.Object{mobileAppWin32LobAppRequirementValidator},
										MarkdownDescription: "Contains file or folder path to detect a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobappfilesystemrequirement?view=graph-rest-beta",
									},
								},
								"powershell_script": generic.OdataDerivedTypeNestedAttributeRs{
									DerivedType: "#microsoft.graph.win32LobAppPowerShellScriptRequirement",
									SingleNestedAttribute: schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{ // win32LobAppPowerShellScriptRequirement
											"detection_value": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "The detection value",
											},
											"operator": schema.StringAttribute{
												Optional: true,
												Validators: []validator.String{
													wpvalidator.FlagEnumValues("notConfigured", "equal", "notEqual", "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual"),
												},
												PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
												Computed:            true,
												MarkdownDescription: "The operator for detection / Contains properties for detection operator; possible values are: `notConfigured` (Not configured.), `equal` (Equal operator.), `notEqual` (Not equal operator.), `greaterThan` (Greater than operator.), `greaterThanOrEqual` (Greater than or equal operator.), `lessThan` (Less than operator.), `lessThanOrEqual` (Less than or equal operator.). The _provider_ default value is `\"notConfigured\"`.",
											},
											"detection_type": schema.StringAttribute{
												Required: true,
												Validators: []validator.String{
													stringvalidator.OneOf("notConfigured", "string", "dateTime", "integer", "float", "version", "boolean"),
												},
												MarkdownDescription: "The detection type for script output. / Contains all supported Powershell Script output detection type; possible values are: `notConfigured` (Not configured.), `string` (Output data type is string.), `dateTime` (Output data type is date time.), `integer` (Output data type is integer.), `float` (Output data type is float.), `version` (Output data type is version.), `boolean` (Output data type is boolean.)",
											},
											"display_name": schema.StringAttribute{
												Optional:            true,
												PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
												Computed:            true,
												MarkdownDescription: "The unique display name for this rule. The _provider_ default value is `\"\"`.",
											},
											"enforce_signature_check": schema.BoolAttribute{
												Optional:            true,
												PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
												Computed:            true,
												MarkdownDescription: "A value indicating whether signature check is enforced. The _provider_ default value is `false`.",
											},
											"run_as_32_bit": schema.BoolAttribute{
												Optional:            true,
												PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
												Computed:            true,
												Description:         `runAs32Bit`, // custom MS Graph attribute name
												MarkdownDescription: "A value indicating whether this script should run as 32-bit. The _provider_ default value is `false`.",
											},
											"run_as_account": schema.StringAttribute{
												Optional:            true,
												Validators:          []validator.String{stringvalidator.OneOf("system", "user")},
												PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("system")},
												Computed:            true,
												MarkdownDescription: "Indicates the type of execution context the script runs in. / Indicates the type of execution context the app runs in; possible values are: `system` (System context), `user` (User context). The _provider_ default value is `\"system\"`.",
											},
											"script_content": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "The base64 encoded script content to detect Win32 Line of Business (LoB) app",
											},
										},
										Validators:          []validator.Object{mobileAppWin32LobAppRequirementValidator},
										MarkdownDescription: "Contains PowerShell script properties to detect a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobapppowershellscriptrequirement?view=graph-rest-beta",
									},
								},
								"registry": generic.OdataDerivedTypeNestedAttributeRs{
									DerivedType: "#microsoft.graph.win32LobAppRegistryRequirement",
									SingleNestedAttribute: schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{ // win32LobAppRegistryRequirement
											"detection_value": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "The detection value",
											},
											"operator": schema.StringAttribute{
												Optional: true,
												Validators: []validator.String{
													wpvalidator.FlagEnumValues("notConfigured", "equal", "notEqual", "greaterThan", "greaterThanOrEqual", "lessThan", "lessThanOrEqual"),
												},
												PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
												Computed:            true,
												MarkdownDescription: "The operator for detection / Contains properties for detection operator; possible values are: `notConfigured` (Not configured.), `equal` (Equal operator.), `notEqual` (Not equal operator.), `greaterThan` (Greater than operator.), `greaterThanOrEqual` (Greater than or equal operator.), `lessThan` (Less than operator.), `lessThanOrEqual` (Less than or equal operator.). The _provider_ default value is `\"notConfigured\"`.",
											},
											"check_32_bit_on_64_system": schema.BoolAttribute{
												Optional:            true,
												PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
												Computed:            true,
												Description:         `check32BitOn64System`, // custom MS Graph attribute name
												MarkdownDescription: "A value indicating whether this registry path is for checking 32-bit app on 64-bit system. The _provider_ default value is `false`.",
											},
											"detection_type": schema.StringAttribute{
												Required: true,
												Validators: []validator.String{
													stringvalidator.OneOf("notConfigured", "exists", "doesNotExist", "string", "integer", "version"),
												},
												MarkdownDescription: "The registry data detection type. / Contains all supported registry data detection type; possible values are: `notConfigured` (Not configured.), `exists` (The specified registry key or value exists.), `doesNotExist` (The specified registry key or value does not exist.), `string` (String value type.), `integer` (Integer value type.), `version` (Version value type.)",
											},
											"key_path": schema.StringAttribute{
												Required:            true,
												MarkdownDescription: "The registry key path to detect Win32 Line of Business (LoB) app",
											},
											"value_name": schema.StringAttribute{
												Optional:            true,
												MarkdownDescription: "The registry value name",
											},
										},
										Validators:          []validator.Object{mobileAppWin32LobAppRequirementValidator},
										MarkdownDescription: "Contains registry properties to detect a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobappregistryrequirement?view=graph-rest-beta",
									},
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "The requirement rules to detect Win32 Line of Business (LoB) app. / Base class to detect a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobapprequirement?view=graph-rest-beta. The _provider_ default value is `[]`.",
					},
					"return_codes": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // win32LobAppReturnCode
								"return_code": schema.Int64Attribute{
									Required:            true,
									MarkdownDescription: "Return code.",
								},
								"type": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf("failed", "success", "softReboot", "hardReboot", "retry"),
									},
									MarkdownDescription: "The type of return code. / Indicates the type of return code; possible values are: `failed` (Failed.), `success` (Success.), `softReboot` (Soft-reboot is required.), `hardReboot` (Hard-reboot is required.), `retry` (Retry.)",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{mobileAppWin32LobAppReturnCodesDefault},
						Computed:            true,
						MarkdownDescription: "The return codes for post installation behavior. / Contains return code properties for a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobappreturncode?view=graph-rest-beta. The _provider_ default value is `mobileAppWin32LobAppReturnCodesDefault`.",
					},
					"rules": schema.SetNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // win32LobAppRule
								"rule_type": schema.StringAttribute{
									Computed:            true,
									Validators:          []validator.String{stringvalidator.OneOf("detection", "requirement")},
									MarkdownDescription: "The rule type indicating the purpose of the rule. / Contains rule types for Win32 LOB apps; possible values are: `detection` (Detection rule.), `requirement` (Requirement rule.)",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpplanmodifier.SetUseStateForUnknown()},
						MarkdownDescription: "The detection and requirement rules for this app. / A base complex type to store the detection or requirement rule data for a Win32 LOB app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobapprule?view=graph-rest-beta",
					},
					"setup_file_path": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The relative path of the setup file in the encrypted Win32LobApp package.",
					},
					"uninstall_command_line": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The command line to uninstall this app",
					},
					"source_file": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Provider Note: The path to the IntuneWin file to be uploaded.",
					},
					"source_sha256": schema.StringAttribute{
						Required:            true,
						Description:         `sourceSha256`, // custom MS Graph attribute name
						MarkdownDescription: "Provider Note: The SHA-256 sum of the source file (see example). Changing this attribute will trigger a new upload of the source file.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for Win32 apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobapp?view=graph-rest-beta",
			},
		},
		"windows_ms_edge": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.windowsMicrosoftEdgeApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // windowsMicrosoftEdgeApp
					"channel": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("dev", "beta", "stable", "unknownFutureValue"),
						},
						PlanModifiers: []planmodifier.String{
							wpdefaultvalue.StringDefaultValue("stable"),
							stringplanmodifier.RequiresReplace(),
						},
						Computed:            true,
						MarkdownDescription: "The channel to install on target devices. The possible values are dev, beta, and stable. By default, this property is set to dev. / The enum to specify the channels for Microsoft Edge apps; possible values are: `dev` (The Dev Channel is intended to help you plan and develop with the latest capabilities of Microsoft Edge.), `beta` (The Beta Channel is intended for production deployment to a representative sample set of users. New features ship about every 4 weeks. Security and quality updates ship as needed.), `stable` (The Stable Channel is intended for broad deployment within organizations, and it's the channel that most users should be on. New features ship about every 4 weeks. Security and quality updates ship as needed.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.). The _provider_ default value is `\"stable\"`.",
					},
					"display_language_locale": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The language locale to use when the Edge app displays text to the user.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for the Microsoft Edge app on Windows. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-windowsmicrosoftedgeapp?view=graph-rest-beta",
			},
		},
		"windows_msi": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.windowsMobileMSI",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // windowsMobileMSI
					"committed_content_version": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The internal committed content version.",
					},
					"file_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The name of the main Lob application file.",
					},
					"size": schema.Int64Attribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
						MarkdownDescription: "The total size, including all uploaded files. This property is read-only.",
					},
					"command_line": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The command line.",
					},
					"identity_version": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The identity version.",
					},
					"ignore_version_detection": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "A boolean to control whether the app's version will be used to detect the app after it is installed on a device. Set this to true for Windows Mobile MSI Line of Business (LoB) apps that use a self update feature. The _provider_ default value is `false`.",
					},
					"product_code": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The product code.",
					},
					"product_version": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The product version of Windows Mobile MSI Line of Business (LoB) app.",
					},
					"use_device_context": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
						MarkdownDescription: "Indicates whether to install a dual-mode MSI in the device context. If true, app will be installed for all users. If false, app will be installed per-user. If null, service will use the MSI package's default install context. In case of dual-mode MSI, this default will be per-user.  Cannot be set for non-dual-mode apps.  Cannot be changed after initial creation of the application.",
					},
					"source_file": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Provider Note: The path to the MSI file to be uploaded.",
					},
					"source_sha256": schema.StringAttribute{
						Required:            true,
						Description:         `sourceSha256`, // custom MS Graph attribute name
						MarkdownDescription: "Provider Note: The SHA-256 sum of the source file (see example). Changing this attribute will trigger a new upload of the source file.",
					},
					"manifest_base64": schema.StringAttribute{
						Required:            true,
						Description:         `manifestBase64`, // custom MS Graph attribute name
						MarkdownDescription: "Provider Note: The manifest contains settings that are based on MSI properties (that cannot easily be parsed and read by the provider). Instead, either manually read the MSI file's properties (e.g. using MSI Orca on Windows) and based on their values select the appropriate values for the XML attributes (as shown in detail in the example above). Or open the Intune portal, enable network log recording in your browser's development tool and thereafter upload the MSI from there (creating a dummy app). Then search the network log for a POST request to https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/.../microsoft.graph.windowsMobileMSI/contentVersions/.../files and directly copy the base64 encoded manifest value from its payload.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for Windows Mobile MSI Line Of Business apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-windowsmobilemsi?view=graph-rest-beta",
			},
		},
		"windows_office_suite": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.officeSuiteApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // officeSuiteApp
					"auto_accept_eula": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
						Computed:            true,
						MarkdownDescription: "The value to accept the EULA automatically on the enduser's device. The _provider_ default value is `true`.",
					},
					"excluded_apps": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // excludedApps
							"access": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: "The value for if MS Office Access should be excluded or not. The _provider_ default value is `false`.",
							},
							"bing": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: "The value for if Microsoft Search as default should be excluded or not. The _provider_ default value is `false`.",
							},
							"excel": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: "The value for if MS Office Excel should be excluded or not. The _provider_ default value is `false`.",
							},
							"groove": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
								Computed:            true,
								MarkdownDescription: "The value for if MS Office OneDrive for Business - Groove should be excluded or not. The _provider_ default value is `true`.",
							},
							"infopath": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
								Computed:            true,
								Description:         `infoPath`, // custom MS Graph attribute name
								MarkdownDescription: "The value for if MS Office InfoPath should be excluded or not. The _provider_ default value is `true`.",
							},
							"lync": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
								Computed:            true,
								MarkdownDescription: "The value for if MS Office Skype for Business - Lync should be excluded or not. The _provider_ default value is `true`.",
							},
							"onedrive": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `oneDrive`, // custom MS Graph attribute name
								MarkdownDescription: "The value for if MS Office OneDrive should be excluded or not. The _provider_ default value is `false`.",
							},
							"onenote": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `oneNote`, // custom MS Graph attribute name
								MarkdownDescription: "The value for if MS Office OneNote should be excluded or not. The _provider_ default value is `false`.",
							},
							"outlook": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: "The value for if MS Office Outlook should be excluded or not. The _provider_ default value is `false`.",
							},
							"powerpoint": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								Description:         `powerPoint`, // custom MS Graph attribute name
								MarkdownDescription: "The value for if MS Office PowerPoint should be excluded or not. The _provider_ default value is `false`.",
							},
							"publisher": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: "The value for if MS Office Publisher should be excluded or not. The _provider_ default value is `false`.",
							},
							"sharepoint_designer": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
								Computed:            true,
								Description:         `sharePointDesigner`, // custom MS Graph attribute name
								MarkdownDescription: "The value for if MS Office SharePointDesigner should be excluded or not. The _provider_ default value is `true`.",
							},
							"teams": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: "The value for if MS Office Teams should be excluded or not. The _provider_ default value is `false`.",
							},
							"visio": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: "The value for if MS Office Visio should be excluded or not. The _provider_ default value is `false`.",
							},
							"word": schema.BoolAttribute{
								Optional:            true,
								PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
								Computed:            true,
								MarkdownDescription: "The value for if MS Office Word should be excluded or not. The _provider_ default value is `false`.",
							},
						},
						PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "The property to represent the apps which are excluded from the selected Office365 Product Id. / Contains properties for Excluded Office365 Apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-excludedapps?view=graph-rest-beta. The _provider_ default value is `{}`.",
					},
					"install_progress_display_level": schema.StringAttribute{
						Optional:            true,
						Validators:          []validator.String{stringvalidator.OneOf("none", "full")},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
						Computed:            true,
						MarkdownDescription: "To specify the level of display for the Installation Progress Setup UI on the Device. / The Enum to specify the level of display for the Installation Progress Setup UI on the Device; possible values are: `none`, `full`. The _provider_ default value is `\"none\"`.",
					},
					"locales_to_install": schema.SetAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "The property to represent the locales which are installed when the apps from Office365 is installed. It uses standard RFC 6033. Ref: https://technet.microsoft.com/library/cc179219(v=office.16).aspx. The _provider_ default value is `[]`.",
					},
					"office_configuration_xml": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The property to represent the XML configuration file that can be specified for Office ProPlus Apps. Takes precedence over all other properties. When present, the XML configuration file will be used to create the app.",
					},
					"office_platform_architecture": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							wpvalidator.FlagEnumValues("none", "x86", "x64", "arm", "neutral", "arm64"),
						},
						PlanModifiers: []planmodifier.String{
							wpdefaultvalue.StringDefaultValue("x64"),
							stringplanmodifier.RequiresReplace(),
						},
						Computed:            true,
						MarkdownDescription: "The property to represent the Office365 app suite version. / Contains properties for Windows architecture; possible values are: `none` (No flags set.), `x86` (Whether or not the X86 Windows architecture type is supported.), `x64` (Whether or not the X64 Windows architecture type is supported.), `arm` (Whether or not the Arm Windows architecture type is supported.), `neutral` (Whether or not the Neutral Windows architecture type is supported.), `arm64` (Whether or not the Arm64 Windows architecture type is supported.). The _provider_ default value is `\"x64\"`.",
					},
					"office_suite_app_default_file_format": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("notConfigured", "officeOpenXMLFormat", "officeOpenDocumentFormat", "unknownFutureValue"),
						},
						PlanModifiers: []planmodifier.String{
							wpdefaultvalue.StringDefaultValue("officeOpenDocumentFormat"),
						},
						Computed:            true,
						MarkdownDescription: "The property to represent the Office365 default file format type. / Describes the OfficeSuiteApp file format types that can be selected; possible values are: `notConfigured` (No file format selected), `officeOpenXMLFormat` (Office Open XML Format selected), `officeOpenDocumentFormat` (Office Open Document Format selected), `unknownFutureValue` (Placeholder for evolvable enum.). The _provider_ default value is `\"officeOpenDocumentFormat\"`.",
					},
					"product_ids": schema.SetAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.ValueStringsAre(
								stringvalidator.OneOf("o365ProPlusRetail", "o365BusinessRetail", "visioProRetail", "projectProRetail"),
							),
						},
						PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
						Computed:            true,
						MarkdownDescription: "The Product Ids that represent the Office365 Suite SKU. / The Enum to specify the Office365 ProductIds that represent the Office365 Suite SKUs; possible values are: `o365ProPlusRetail`, `o365BusinessRetail`, `visioProRetail`, `projectProRetail`. The _provider_ default value is `[]`.",
					},
					"should_uninstall_older_versions_of_office": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
						Computed:            true,
						MarkdownDescription: "The property to determine whether to uninstall existing Office MSI if an Office365 app suite is deployed to the device or not. The _provider_ default value is `true`.",
					},
					"target_version": schema.StringAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
						Computed:            true,
						MarkdownDescription: "The property to represent the specific target version for the Office365 app suite that should be remained deployed on the devices. The _provider_ default value is `\"\"`.",
					},
					"update_channel": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("none", "current", "deferred", "firstReleaseCurrent", "firstReleaseDeferred", "monthlyEnterprise"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("current")},
						Computed:            true,
						MarkdownDescription: "The property to represent the Office365 Update Channel. / The Enum to specify the Office365 Updates Channel; possible values are: `none`, `current`, `deferred`, `firstReleaseCurrent`, `firstReleaseDeferred`, `monthlyEnterprise`. The _provider_ default value is `\"current\"`.",
					},
					"update_version": schema.StringAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
						Computed:            true,
						MarkdownDescription: "The property to represent the update version in which the specific target version is available for the Office365 app suite. The _provider_ default value is `\"\"`.",
					},
					"use_shared_computer_activation": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: "The property to represent that whether the shared computer activation is used not for Office365 app suite. The _provider_ default value is `false`.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for the Office365 Suite App. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-officesuiteapp?view=graph-rest-beta",
			},
		},
		"windows_store": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.windowsStoreApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // windowsStoreApp
					"app_store_url": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The Windows app store URL.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for Windows Store apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-windowsstoreapp?view=graph-rest-beta",
			},
		},
		"windows_universal_appx": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.windowsUniversalAppX",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // windowsUniversalAppX
					"committed_content_version": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The internal committed content version.",
					},
					"file_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The name of the main Lob application file.",
					},
					"size": schema.Int64Attribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
						MarkdownDescription: "The total size, including all uploaded files. This property is read-only.",
					},
					"applicable_architectures": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							wpvalidator.FlagEnumValues("none", "x86", "x64", "arm", "neutral", "arm64"),
						},
						MarkdownDescription: "The Windows architecture(s) for which this app can run on.; default value is `none`. / Contains properties for Windows architecture; possible values are: `none` (No flags set.), `x86` (Whether or not the X86 Windows architecture type is supported.), `x64` (Whether or not the X64 Windows architecture type is supported.), `arm` (Whether or not the Arm Windows architecture type is supported.), `neutral` (Whether or not the Neutral Windows architecture type is supported.), `arm64` (Whether or not the Arm64 Windows architecture type is supported.)",
					},
					"applicable_device_types": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							wpvalidator.FlagEnumValues("none", "desktop", "mobile", "holographic", "team", "unknownFutureValue"),
						},
						MarkdownDescription: "The Windows device type(s) for which this app can run on.; default value is `none`. / Contains properties for Windows device type. Multiple values can be selected. Default value is `none`; possible values are: `none` (No device types supported. Default value.), `desktop` (Indicates support for Desktop Windows device type.), `mobile` (Indicates support for Mobile Windows device type.), `holographic` (Indicates support for Holographic Windows device type.), `team` (Indicates support for Team Windows device type.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.)",
					},
					"identity_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The Identity Name of the app, parsed from the appx file when it is uploaded through the Intune MEM console. For example: \"Contoso.DemoApp\".",
					},
					"identity_publisher_hash": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The Identity Publisher Hash of the app, parsed from the appx file when it is uploaded through the Intune MEM console. For example: \"AB82CD0XYZ\".",
					},
					"identity_resource_identifier": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The Identity Resource Identifier of the app, parsed from the appx file when it is uploaded through the Intune MEM console. For example: \"TestResourceId\".",
					},
					"identity_version": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The Identity Version of the app, parsed from the appx file when it is uploaded through the Intune MEM console.  For example: \"1.0.0.0\".",
					},
					"is_bundle": schema.BoolAttribute{
						Optional: true,
						PlanModifiers: []planmodifier.Bool{
							wpdefaultvalue.BoolDefaultValue(false),
							boolplanmodifier.RequiresReplace(),
						},
						Computed:            true,
						MarkdownDescription: "Whether or not the app is a bundle. If TRUE, app is a bundle; if FALSE, app is not a bundle. The _provider_ default value is `false`.",
					},
					"minimum_supported_operating_system": schema.SingleNestedAttribute{
						Required:            true,
						Attributes:          mobileAppWindowsMinimumOperatingSystemAttributes,
						MarkdownDescription: "The value for the minimum applicable Windows operating system. / The minimum operating system required for a Windows mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-windowsminimumoperatingsystem?view=graph-rest-beta  \nProvider Note: v10_0 seems to be required and the only minimum OS version to be accepted by MS Graph for Univeral APPX.",
					},
					"committed_contained_apps": schema.SetNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // mobileContainedApp
								"id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Key of the entity. This property is read-only.",
								},
							},
						},
						PlanModifiers:       []planmodifier.Set{wpplanmodifier.SetUseStateForUnknown()},
						MarkdownDescription: "The collection of contained apps in the committed mobileAppContent of a windowsUniversalAppX app. This property is read-only. / An abstract class that represents a contained app in a mobileApp acting as a package. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobilecontainedapp?view=graph-rest-beta",
					},
					"source_file": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Provider Note: The path to the APPX/MSIX file to be uploaded.",
					},
					"source_sha256": schema.StringAttribute{
						Required:            true,
						Description:         `sourceSha256`, // custom MS Graph attribute name
						MarkdownDescription: "Provider Note: The SHA-256 sum of the source file (see example). Changing this attribute will trigger a new upload of the source file.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for Windows Universal AppX Line Of Business apps. Inherits from `mobileLobApp`. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-windowsuniversalappx?view=graph-rest-beta",
			},
		},
		"windows_web_link": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.windowsWebApp",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // windowsWebApp
					"app_url": schema.StringAttribute{
						Required:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
						MarkdownDescription: "Indicates the Windows web app URL. Example: \"https://www.contoso.com\"",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "Contains properties and inherited properties for Windows web apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-windowswebapp?view=graph-rest-beta",
			},
		},
		"winget": generic.OdataDerivedTypeNestedAttributeRs{
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
								MarkdownDescription: "Indicates the type of execution context the app setup runs in on target devices. Options include values of the RunAsAccountType enum, which are System and User. Required at creation time, cannot be modified on existing objects. / Indicates the type of execution context the app runs in; possible values are: `system` (System context), `user` (User context). The _provider_ default value is `\"system\"`.",
							},
						},
						PlanModifiers: []planmodifier.Object{
							wpdefaultvalue.ObjectDefaultValueEmpty(),
							objectplanmodifier.RequiresReplace(),
						},
						Computed:            true,
						MarkdownDescription: "The install experience settings associated with this application, which are used to ensure the desired install experiences on the target device are taken into account. This includes the account type (System or User) that actions should be run as on target devices. Required at creation time. / Represents the install experience settings associated with WinGet apps. This is used to ensure the desired install experiences on the target device are taken into account. Required at creation time. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-wingetappinstallexperience?view=graph-rest-beta. The _provider_ default value is `{}`.",
					},
					"manifest_hash": schema.StringAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
						MarkdownDescription: "Hash of package metadata properties used to validate that the application matches the metadata in the source repository.",
					},
					"package_identifier": schema.StringAttribute{
						Required:            true,
						PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
						MarkdownDescription: "The PackageIdentifier from the WinGet source repository REST API. This also maps to the Id when using the WinGet client command line application. Required at creation time, cannot be modified on existing objects.",
					},
				},
				Validators:          []validator.Object{mobileAppMobileAppValidator},
				MarkdownDescription: "A MobileApp that is based on a referenced application in a WinGet repository. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-wingetapp?view=graph-rest-beta",
			},
		},
	},
	MarkdownDescription: "An abstract class containing the base properties for Intune mobile apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mobileapp?view=graph-rest-beta\n\nProvider Note: There have been reports about OpenTofu crashing (with a panic) when trying to generate configuration using `-generate-config-out=...`. Unfortunately this happens inside the code of OpenTofu itself and is not easily diagnoseable from outside or changeable from the provider. A workaround might be to instead use Terraform for generating the config (as it seems to work well for it) - the result should be fully compatible and usable also with OpenTofu. ||| MS Graph: App management",
}

var mobileAppMobileAppValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("android_for_work"),
	path.MatchRelative().AtParent().AtName("android_lob"),
	path.MatchRelative().AtParent().AtName("android_managed_store"),
	path.MatchRelative().AtParent().AtName("android_store"),
	path.MatchRelative().AtParent().AtName("ios_lob"),
	path.MatchRelative().AtParent().AtName("ios_store"),
	path.MatchRelative().AtParent().AtName("ios_vpp"),
	path.MatchRelative().AtParent().AtName("ios_webclip"),
	path.MatchRelative().AtParent().AtName("macos_dmg"),
	path.MatchRelative().AtParent().AtName("macos_lob"),
	path.MatchRelative().AtParent().AtName("macos_ms_defender"),
	path.MatchRelative().AtParent().AtName("macos_ms_edge"),
	path.MatchRelative().AtParent().AtName("macos_office_suite"),
	path.MatchRelative().AtParent().AtName("macos_pkg"),
	path.MatchRelative().AtParent().AtName("macos_vpp"),
	path.MatchRelative().AtParent().AtName("macos_webclip"),
	path.MatchRelative().AtParent().AtName("managed_android_lob"),
	path.MatchRelative().AtParent().AtName("managed_android_store"),
	path.MatchRelative().AtParent().AtName("managed_ios_lob"),
	path.MatchRelative().AtParent().AtName("managed_ios_store"),
	path.MatchRelative().AtParent().AtName("web_link"),
	path.MatchRelative().AtParent().AtName("win32_lob"),
	path.MatchRelative().AtParent().AtName("windows_ms_edge"),
	path.MatchRelative().AtParent().AtName("windows_msi"),
	path.MatchRelative().AtParent().AtName("windows_office_suite"),
	path.MatchRelative().AtParent().AtName("windows_store"),
	path.MatchRelative().AtParent().AtName("windows_universal_appx"),
	path.MatchRelative().AtParent().AtName("windows_web_link"),
	path.MatchRelative().AtParent().AtName("winget"),
)

var mobileAppMobileAppAssignmentAttributes = map[string]schema.Attribute{ // mobileAppAssignment
	"intent": schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			stringvalidator.OneOf("available", "required", "uninstall", "availableWithoutEnrollment"),
		},
		PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("required")},
		Computed:            true,
		MarkdownDescription: "The install intent defined by the admin. / Possible values for the install intent chosen by the admin; possible values are: `available` (Available install intent.), `required` (Required install intent.), `uninstall` (Uninstall install intent.), `availableWithoutEnrollment` (Available without enrollment install intent.). The _provider_ default value is `\"required\"`.",
	},
	"settings": schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{ // mobileAppAssignmentSettings
			"android_managed_store": generic.OdataDerivedTypeNestedAttributeRs{
				DerivedType: "#microsoft.graph.androidManagedStoreAppAssignmentSettings",
				SingleNestedAttribute: schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // androidManagedStoreAppAssignmentSettings
						"android_managed_store_app_track_ids": schema.SetAttribute{
							ElementType:         types.StringType,
							Optional:            true,
							PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
							Computed:            true,
							MarkdownDescription: "The track IDs to enable for this app assignment. The _provider_ default value is `[]`.",
						},
						"auto_update_mode": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("default", "postponed", "priority", "unknownFutureValue"),
							},
							PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("default")},
							Computed:            true,
							MarkdownDescription: "The prioritization of automatic updates for this app assignment. / Prioritization for automatic updates of Android Managed Store apps set on assignment; possible values are: `default` (Default update behavior (device must be connected to Wifi, charging and not actively used).), `postponed` (Updates are postponed for a maximum of 90 days after the app becomes out of date.), `priority` (The app is updated as soon as possible by the developer. If device is online, it will be updated within minutes.), `unknownFutureValue` (Unknown future mode (reserved, not used right now).). The _provider_ default value is `\"default\"`.",
						},
					},
					Validators:          []validator.Object{mobileAppMobileAppAssignmentSettingsValidator},
					MarkdownDescription: "Contains properties used to assign an Android Managed Store mobile app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-androidmanagedstoreappassignmentsettings?view=graph-rest-beta",
				},
			},
			"ios_lob": generic.OdataDerivedTypeNestedAttributeRs{
				DerivedType: "#microsoft.graph.iosLobAppAssignmentSettings",
				SingleNestedAttribute: schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // iosLobAppAssignmentSettings
						"is_removable": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "When TRUE, indicates that the app can be uninstalled by the user. When FALSE, indicates that the app cannot be uninstalled by the user. By default, this property is set to null which internally is treated as TRUE.",
						},
						"prevent_managed_app_backup": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to null which internally is treated as FALSE.",
						},
						"uninstall_on_device_removal": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "When TRUE, indicates that the app should be uninstalled when the device is removed from Intune. When FALSE, indicates that the app will not be uninstalled when the device is removed from Intune. By default, property is set to null which internally is treated as TRUE.",
						},
						"vpn_configuration_id": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "This is the unique identifier (Id) of the VPN Configuration to apply to the app.",
						},
					},
					Validators:          []validator.Object{mobileAppMobileAppAssignmentSettingsValidator},
					MarkdownDescription: "Contains properties used to assign an iOS LOB mobile app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-ioslobappassignmentsettings?view=graph-rest-beta",
				},
			},
			"ios_store": generic.OdataDerivedTypeNestedAttributeRs{
				DerivedType: "#microsoft.graph.iosStoreAppAssignmentSettings",
				SingleNestedAttribute: schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // iosStoreAppAssignmentSettings
						"is_removable": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "When TRUE, indicates that the app can be uninstalled by the user. When FALSE, indicates that the app cannot be uninstalled by the user. By default, this property is set to null which internally is treated as TRUE.",
						},
						"prevent_managed_app_backup": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to null which internally is treated as FALSE.",
						},
						"uninstall_on_device_removal": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "When TRUE, indicates that the app should be uninstalled when the device is removed from Intune. When FALSE, indicates that the app will not be uninstalled when the device is removed from Intune. By default, property is set to null which internally is treated as TRUE.",
						},
						"vpn_configuration_id": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "This is the unique identifier (Id) of the VPN Configuration to apply to the app.",
						},
					},
					Validators:          []validator.Object{mobileAppMobileAppAssignmentSettingsValidator},
					MarkdownDescription: "Contains properties used to assign an iOS Store mobile app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-iosstoreappassignmentsettings?view=graph-rest-beta",
				},
			},
			"ios_vpp": generic.OdataDerivedTypeNestedAttributeRs{
				DerivedType: "#microsoft.graph.iosVppAppAssignmentSettings",
				SingleNestedAttribute: schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // iosVppAppAssignmentSettings
						"is_removable": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Whether or not the app can be removed by the user.",
						},
						"prevent_auto_app_update": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "When TRUE, indicates that the app should not be automatically updated with the latest version from Apple app store. When FALSE, indicates that the app may be auto updated. By default, this property is set to null which internally is treated as FALSE.",
						},
						"prevent_managed_app_backup": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to null which internally is treated as FALSE.",
						},
						"uninstall_on_device_removal": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Whether or not to uninstall the app when device is removed from Intune.",
						},
						"use_device_licensing": schema.BoolAttribute{
							Optional:            true,
							PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
							Computed:            true,
							MarkdownDescription: "Whether or not to use device licensing. The _provider_ default value is `false`.",
						},
						"vpn_configuration_id": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "The VPN Configuration Id to apply for this app.",
						},
					},
					Validators:          []validator.Object{mobileAppMobileAppAssignmentSettingsValidator},
					MarkdownDescription: "Contains properties used to assign an iOS VPP mobile app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-iosvppappassignmentsettings?view=graph-rest-beta",
				},
			},
			"macos_lob": generic.OdataDerivedTypeNestedAttributeRs{
				DerivedType: "#microsoft.graph.macOsLobAppAssignmentSettings",
				SingleNestedAttribute: schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // macOsLobAppAssignmentSettings
						"uninstall_on_device_removal": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "When TRUE, indicates that the app should be uninstalled when the device is removed from Intune. When FALSE, indicates that the app will not be uninstalled when the device is removed from Intune.",
						},
					},
					Validators:          []validator.Object{mobileAppMobileAppAssignmentSettingsValidator},
					MarkdownDescription: "Contains properties used to assign a macOS LOB app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-macoslobappassignmentsettings?view=graph-rest-beta",
				},
			},
			"macos_vpp": generic.OdataDerivedTypeNestedAttributeRs{
				DerivedType: "#microsoft.graph.macOsVppAppAssignmentSettings",
				SingleNestedAttribute: schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // macOsVppAppAssignmentSettings
						"prevent_auto_app_update": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "When TRUE, indicates that the app should not be automatically updated with the latest version from Apple app store. When FALSE, indicates that the app may be auto updated. By default, this property is set to null which internally is treated as FALSE.",
						},
						"prevent_managed_app_backup": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to null which internally is treated as FALSE.",
						},
						"uninstall_on_device_removal": schema.BoolAttribute{
							Optional:            true,
							MarkdownDescription: "Whether or not to uninstall the app when device is removed from Intune.",
						},
						"use_device_licensing": schema.BoolAttribute{
							Optional:            true,
							PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
							Computed:            true,
							MarkdownDescription: "Whether or not to use device licensing. The _provider_ default value is `false`.",
						},
					},
					Validators:          []validator.Object{mobileAppMobileAppAssignmentSettingsValidator},
					MarkdownDescription: "Contains properties used to assign an Mac VPP mobile app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-macosvppappassignmentsettings?view=graph-rest-beta",
				},
			},
			"win32_lob": generic.OdataDerivedTypeNestedAttributeRs{
				DerivedType: "#microsoft.graph.win32LobAppAssignmentSettings",
				SingleNestedAttribute: schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // win32LobAppAssignmentSettings
						"auto_update_settings": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // win32LobAppAutoUpdateSettings
								"auto_update_superseded_apps_state": schema.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.OneOf("notConfigured", "enabled", "unknownFutureValue"),
									},
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
									Computed:            true,
									MarkdownDescription: "Contains value for auto-update superseded apps; possible values are: `notConfigured` (Indicates that the auto-update superseded apps state is not configured and the app will not auto-update the superseded apps.), `enabled` (Indicates that the auto-update superseded apps state is enabled and the app will auto-update the superseded apps if the superseded apps are installed on the device.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.). The _provider_ default value is `\"notConfigured\"`.",
								},
							},
							MarkdownDescription: "The auto-update settings to apply for this app assignment. / Contains properties used to perform the auto-update of an application. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-win32lobappautoupdatesettings?view=graph-rest-beta",
						},
						"delivery_optimization_priority": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("notConfigured", "foreground"),
							},
							PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
							Computed:            true,
							MarkdownDescription: "The delivery optimization priority for this app assignment. This setting is not supported in National Cloud environments. / Contains value for delivery optimization priority; possible values are: `notConfigured` (Not configured or background normal delivery optimization priority.), `foreground` (Foreground delivery optimization priority.). The _provider_ default value is `\"notConfigured\"`.",
						},
						"install_time_settings": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // mobileAppInstallTimeSettings
								"deadline_date_time": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The time at which the app should be installed.",
								},
								"start_date_time": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "The time at which the app should be available for installation.",
								},
								"use_local_time": schema.BoolAttribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
									Computed:            true,
									MarkdownDescription: "Whether the local device time or UTC time should be used when determining the available and deadline times. The _provider_ default value is `false`.",
								},
							},
							MarkdownDescription: "The install time settings to apply for this app assignment. / Contains properties used to determine when to offer an app to devices and when to install the app on devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mobileappinstalltimesettings?view=graph-rest-beta",
						},
						"notifications": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("showAll", "showReboot", "hideAll"),
							},
							PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("showAll")},
							Computed:            true,
							MarkdownDescription: "The notification status for this app assignment. / Contains value for notification status; possible values are: `showAll` (Show all notifications.), `showReboot` (Only show restart notification and suppress other notifications.), `hideAll` (Hide all notifications.). The _provider_ default value is `\"showAll\"`.",
						},
						"restart_settings": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // win32LobAppRestartSettings
								"countdown_display_before_restart_in_minutes": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: "The number of minutes before the restart time to display the countdown dialog for pending restarts. The _provider_ default value is `0`.",
								},
								"grace_period_in_minutes": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: "The number of minutes to wait before restarting the device after an app installation. The _provider_ default value is `0`.",
								},
								"restart_notification_snooze_duration_in_minutes": schema.Int64Attribute{
									Optional:            true,
									MarkdownDescription: "The number of minutes to snooze the restart notification dialog when the snooze button is selected.",
								},
							},
							MarkdownDescription: "The reboot settings to apply for this app assignment. / Contains properties describing restart coordination following an app installation. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-win32lobapprestartsettings?view=graph-rest-beta",
						},
					},
					Validators:          []validator.Object{mobileAppMobileAppAssignmentSettingsValidator},
					MarkdownDescription: "Contains properties used to assign an Win32 LOB mobile app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-win32lobappassignmentsettings?view=graph-rest-beta",
				},
			},
			"windows_universal_appx": generic.OdataDerivedTypeNestedAttributeRs{
				DerivedType: "#microsoft.graph.windowsUniversalAppXAppAssignmentSettings",
				SingleNestedAttribute: schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // windowsUniversalAppXAppAssignmentSettings
						"use_device_context": schema.BoolAttribute{
							Optional:            true,
							PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
							Computed:            true,
							MarkdownDescription: "If true, uses device execution context for Windows Universal AppX mobile app. Device-context install is not allowed when this type of app is targeted with Available intent. Defaults to false. The _provider_ default value is `false`.",
						},
					},
					Validators:          []validator.Object{mobileAppMobileAppAssignmentSettingsValidator},
					MarkdownDescription: "Contains properties used when assigning a Windows Universal AppX mobile app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-windowsuniversalappxappassignmentsettings?view=graph-rest-beta",
				},
			},
			"winget": generic.OdataDerivedTypeNestedAttributeRs{
				DerivedType: "#microsoft.graph.winGetAppAssignmentSettings",
				SingleNestedAttribute: schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // winGetAppAssignmentSettings
						"install_time_settings": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // winGetAppInstallTimeSettings
								"deadline_date_time": schema.StringAttribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
									Computed:            true,
									MarkdownDescription: "The time at which the app should be installed. The _provider_ default value is `\"\"`.",
								},
								"use_local_time": schema.BoolAttribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
									Computed:            true,
									MarkdownDescription: "Whether the local device time or UTC time should be used when determining the deadline times. The _provider_ default value is `false`.",
								},
							},
							MarkdownDescription: "The install time settings to apply for this app assignment. / Contains properties used to determine when to offer an app to devices and when to install the app on devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-wingetappinstalltimesettings?view=graph-rest-beta",
						},
						"notifications": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("showAll", "showReboot", "hideAll", "unknownFutureValue"),
							},
							PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("showAll")},
							Computed:            true,
							MarkdownDescription: "The notification status for this app assignment. / Contains value for notification status; possible values are: `showAll` (Show all notifications.), `showReboot` (Only show restart notification and suppress other notifications.), `hideAll` (Hide all notifications.), `unknownFutureValue` (Unknown future value, reserved for future usage as expandable enum.). The _provider_ default value is `\"showAll\"`.",
						},
						"restart_settings": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // winGetAppRestartSettings
								"countdown_display_before_restart_in_minutes": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: "The number of minutes before the restart time to display the countdown dialog for pending restarts. The _provider_ default value is `0`.",
								},
								"grace_period_in_minutes": schema.Int64Attribute{
									Optional:            true,
									PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(0)},
									Computed:            true,
									MarkdownDescription: "The number of minutes to wait before restarting the device after an app installation. The _provider_ default value is `0`.",
								},
								"restart_notification_snooze_duration_in_minutes": schema.Int64Attribute{
									Optional:            true,
									MarkdownDescription: "The number of minutes to snooze the restart notification dialog when the snooze button is selected.",
								},
							},
							MarkdownDescription: "The reboot settings to apply for this app assignment. / Contains properties describing restart coordination following an app installation. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-wingetapprestartsettings?view=graph-rest-beta",
						},
					},
					Validators:          []validator.Object{mobileAppMobileAppAssignmentSettingsValidator},
					MarkdownDescription: "Contains properties used to assign a WinGet app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-wingetappassignmentsettings?view=graph-rest-beta",
				},
			},
		},
		MarkdownDescription: "The settings for target assignment defined by the admin. / Abstract class to contain properties used to assign a mobile app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mobileappassignmentsettings?view=graph-rest-beta  \nProvider Note: For some mobileApp types it might be required to include appropriate (potentially empty) `settings` with each assignment to avoid recurring differences between plan and actual state, namely for the types: `android_managed_store`, `ios_store`, `managed_ios_store` (use `ios_store` as settings type), `winget`, `win32_lob`, `windows_universal_appx`. For details, see the example above.",
	},
	"target": deviceAndAppManagementAssignmentTarget,
}

var mobileAppMobileAppAssignmentSettingsValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("android_managed_store"),
	path.MatchRelative().AtParent().AtName("ios_lob"),
	path.MatchRelative().AtParent().AtName("ios_store"),
	path.MatchRelative().AtParent().AtName("ios_vpp"),
	path.MatchRelative().AtParent().AtName("macos_lob"),
	path.MatchRelative().AtParent().AtName("macos_vpp"),
	path.MatchRelative().AtParent().AtName("win32_lob"),
	path.MatchRelative().AtParent().AtName("windows_universal_appx"),
	path.MatchRelative().AtParent().AtName("winget"),
)

var mobileAppAndroidMinimumOperatingSystemAttributes = map[string]schema.Attribute{ // androidMinimumOperatingSystem
	"v10_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 10.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v11_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v11_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 11.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v12_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v12_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 12.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v13_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v13_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 13.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v14_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v14_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 14.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v15_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v15_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 15.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v4_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v4_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 4.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v4_0_3": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v4_0_3`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 4.0.3 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v4_1": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v4_1`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 4.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v4_2": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v4_2`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 4.2 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v4_3": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v4_3`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 4.3 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v4_4": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v4_4`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 4.4 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v5_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v5_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 5.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v5_1": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v5_1`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 5.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v6_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v6_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 6.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v7_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v7_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 7.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v7_1": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v7_1`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 7.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v8_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v8_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 8.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v8_1": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v8_1`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 8.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v9_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v9_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 9.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
}

var mobileAppIosDeviceTypeAttributes = map[string]schema.Attribute{ // iosDeviceType
	"ipad": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
		Computed:            true,
		Description:         `iPad`, // custom MS Graph attribute name
		MarkdownDescription: "Whether the app should run on iPads. The _provider_ default value is `true`.",
	},
	"iphone_and_ipod": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
		Computed:            true,
		Description:         `iPhoneAndIPod`, // custom MS Graph attribute name
		MarkdownDescription: "Whether the app should run on iPhones and iPods. The _provider_ default value is `true`.",
	},
}

var mobileAppIosMinimumOperatingSystemAttributes = map[string]schema.Attribute{ // iosMinimumOperatingSystem
	"v10_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 10.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v11_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v11_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 11.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v12_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v12_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 12.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v13_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v13_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 13.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v14_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v14_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 14.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v15_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v15_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 15.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v16_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v16_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 16.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v17_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v17_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 17.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v8_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v8_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 8.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
	"v9_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v9_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, only Version 9.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.",
	},
}

var mobileAppVppLicensingTypeAttributes = map[string]schema.Attribute{ // vppLicensingType
	"support_device_licensing": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		MarkdownDescription: "Whether the program supports the device licensing type. The _provider_ default value is `false`.",
	},
	"supports_device_licensing": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		MarkdownDescription: "Whether the program supports the device licensing type. The _provider_ default value is `false`.",
	},
	"supports_user_licensing": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		MarkdownDescription: "Whether the program supports the user licensing type. The _provider_ default value is `false`.",
	},
	"support_user_licensing": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		MarkdownDescription: "Whether the program supports the user licensing type. The _provider_ default value is `false`.",
	},
}

var mobileAppMacOSIncludedAppAttributes = map[string]schema.Attribute{ // macOSIncludedApp
	"bundle_id": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The bundleId of the app. This maps to the CFBundleIdentifier in the app's bundle configuration.",
	},
	"bundle_version": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The version of the app. This maps to the CFBundleShortVersion in the app's bundle configuration.",
	},
}

var mobileAppMacOSMinimumOperatingSystemAttributes = map[string]schema.Attribute{ // macOSMinimumOperatingSystem
	"v10_10": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_10`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, indicates OS X 10.10 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.",
	},
	"v10_11": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_11`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, indicates OS X 10.11 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.",
	},
	"v10_12": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_12`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, indicates macOS 10.12 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.",
	},
	"v10_13": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_13`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, indicates macOS 10.13 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.",
	},
	"v10_14": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_14`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, indicates macOS 10.14 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.",
	},
	"v10_15": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_15`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, indicates macOS 10.15 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.",
	},
	"v10_7": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_7`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, indicates Mac OS X 10.7 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.",
	},
	"v10_8": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_8`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, indicates OS X 10.8 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.",
	},
	"v10_9": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_9`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, indicates OS X 10.9 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.",
	},
	"v11_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v11_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, indicates macOS 11.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.",
	},
	"v12_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v12_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, indicates macOS 12.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.",
	},
	"v13_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v13_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, indicates macOS 13.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.",
	},
	"v14_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v14_0`, // custom MS Graph attribute name
		MarkdownDescription: "When TRUE, indicates macOS 14.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.",
	},
}

var mobileAppWin32LobAppDetectionValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("filesystem"),
	path.MatchRelative().AtParent().AtName("powershell_script"),
	path.MatchRelative().AtParent().AtName("product_code"),
	path.MatchRelative().AtParent().AtName("registry"),
)

var mobileAppWin32LobAppRequirementValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("filesystem"),
	path.MatchRelative().AtParent().AtName("powershell_script"),
	path.MatchRelative().AtParent().AtName("registry"),
)

var mobileAppWindowsMinimumOperatingSystemAttributes = map[string]schema.Attribute{ // windowsMinimumOperatingSystem
	"v10_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_0`, // custom MS Graph attribute name
		MarkdownDescription: "Windows version 10.0 or later. The _provider_ default value is `false`.",
	},
	"v10_1607": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_1607`, // custom MS Graph attribute name
		MarkdownDescription: "Windows 10 1607 or later. The _provider_ default value is `false`.",
	},
	"v10_1703": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_1703`, // custom MS Graph attribute name
		MarkdownDescription: "Windows 10 1703 or later. The _provider_ default value is `false`.",
	},
	"v10_1709": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_1709`, // custom MS Graph attribute name
		MarkdownDescription: "Windows 10 1709 or later. The _provider_ default value is `false`.",
	},
	"v10_1803": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_1803`, // custom MS Graph attribute name
		MarkdownDescription: "Windows 10 1803 or later. The _provider_ default value is `false`.",
	},
	"v10_1809": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_1809`, // custom MS Graph attribute name
		MarkdownDescription: "Windows 10 1809 or later. The _provider_ default value is `false`.",
	},
	"v10_1903": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_1903`, // custom MS Graph attribute name
		MarkdownDescription: "Windows 10 1903 or later. The _provider_ default value is `false`.",
	},
	"v10_1909": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_1909`, // custom MS Graph attribute name
		MarkdownDescription: "Windows 10 1909 or later. The _provider_ default value is `false`.",
	},
	"v10_2004": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_2004`, // custom MS Graph attribute name
		MarkdownDescription: "Windows 10 2004 or later. The _provider_ default value is `false`.",
	},
	"v10_21_h1": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_21H1`, // custom MS Graph attribute name
		MarkdownDescription: "Windows 10 21H1 or later. The _provider_ default value is `false`.",
	},
	"v10_2_h20": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v10_2H20`, // custom MS Graph attribute name
		MarkdownDescription: "Windows 10 2H20 or later. The _provider_ default value is `false`.",
	},
	"v8_0": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v8_0`, // custom MS Graph attribute name
		MarkdownDescription: "Windows version 8.0 or later. The _provider_ default value is `false`.",
	},
	"v8_1": schema.BoolAttribute{
		Optional:            true,
		PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
		Computed:            true,
		Description:         `v8_1`, // custom MS Graph attribute name
		MarkdownDescription: "Windows version 8.1 or later. The _provider_ default value is `false`.",
	},
}

var mobileAppWin32LobAppReturnCodesDefault = wpdefaultvalue.SetDefaultValue([]any{
	map[string]any{
		"return_code": 0,
		"type":        "success",
	},
	map[string]any{
		"return_code": 1707,
		"type":        "success",
	},
	map[string]any{
		"return_code": 3010,
		"type":        "softReboot",
	},
	map[string]any{
		"return_code": 1641,
		"type":        "hardReboot",
	},
	map[string]any{
		"return_code": 1618,
		"type":        "retry",
	},
})
