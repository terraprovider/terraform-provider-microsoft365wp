package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	ManagedDeviceMobileAppConfigurationResource = generic.GenericResource{
		TypeNameSuffix: "managed_device_mobile_app_configuration",
		SpecificSchema: managedDeviceMobileAppConfigurationResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:         "/deviceAppManagement/mobileAppConfigurations",
			ReadOptions:     managedDeviceMobileAppConfigurationReadOptions,
			WriteSubActions: managedDeviceMobileAppConfigurationWriteSubActions,
		},
	}

	ManagedDeviceMobileAppConfigurationSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&ManagedDeviceMobileAppConfigurationResource)

	ManagedDeviceMobileAppConfigurationPluralDataSource = generic.CreateGenericDataSourcePluralFromSingular(
		&ManagedDeviceMobileAppConfigurationSingularDataSource, "")
)

var managedDeviceMobileAppConfigurationReadOptions = generic.ReadOptions{
	ODataExpand: "assignments",
}

var managedDeviceMobileAppConfigurationWriteSubActions = []generic.WriteSubAction{
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			Attributes: []string{"assignments"},
			UriSuffix:  "assign",
			UpdateOnly: true,
		},
	},
}

var managedDeviceMobileAppConfigurationResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // managedDeviceMobileAppConfiguration
		"id": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: `DateTime the object was created.`,
		},
		"description": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: `Admin provided description of the Device Configuration.`,
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: `Admin provided name of the device configuration.`,
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: `DateTime the object was last modified.`,
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: `List of Scope Tags for this App configuration entity.`,
		},
		"targeted_mobile_apps": schema.SetAttribute{
			ElementType:         types.StringType,
			Required:            true,
			MarkdownDescription: `the associated app.`,
		},
		"version": schema.Int64Attribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
			MarkdownDescription: `Version of the device configuration.`,
		},
		"assignments": deviceAndAppManagementAssignment,
		"android_managed_store": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.androidManagedStoreAppConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // androidManagedStoreAppConfiguration
					"app_supports_oem_config": schema.BoolAttribute{
						Optional:      true,
						PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:      true,
					},
					"connected_apps_enabled": schema.BoolAttribute{
						Optional:            true,
						PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
						Computed:            true,
						MarkdownDescription: `Setting to specify whether to allow ConnectedApps experience for this app.`,
					},
					"package_id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: `Android Enterprise app configuration package id.`,
					},
					"payload_json_base64": schema.StringAttribute{
						Optional:            true,
						Description:         `payloadJson`, // custom MS Graph attribute name
						MarkdownDescription: `Android Enterprise app configuration JSON payload.`,
					},
					"permission_actions": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // androidPermissionAction
								"action": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf("prompt", "autoGrant", "autoDeny"),
									},
									MarkdownDescription: `Type of Android permission action.`,
								},
								"permission": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: `Android permission string, defined in the official Android documentation.  Example 'android.permission.READ_CONTACTS'.`,
								},
							},
						},
						MarkdownDescription: `List of Android app permissions and corresponding permission actions. / Mapping between an Android app permission and the action Android should take when that permission is requested.`,
					},
					"profile_applicability": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("default", "androidWorkProfile", "androidDeviceOwner"),
						},
						MarkdownDescription: `Android Enterprise profile applicability (AndroidWorkProfile, DeviceOwner, or default (applies to both)).`,
					},
				},
				Validators: []validator.Object{
					managedDeviceMobileAppConfigurationManagedDeviceMobileAppConfigurationValidator,
				},
				MarkdownDescription: `Contains properties, inherited properties and actions for Android Enterprise mobile app configurations.`,
			},
		},
		"ios": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.iosMobileAppConfiguration",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // iosMobileAppConfiguration
					"encoded_setting_xml_base64": schema.StringAttribute{
						Optional:            true,
						Description:         `encodedSettingXml`, // custom MS Graph attribute name
						MarkdownDescription: `mdm app configuration Base64 binary.`,
					},
					"settings": schema.SetNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{ // appConfigurationSettingItem
								"app_config_key": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: `app configuration key.`,
								},
								"app_config_key_type": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf("stringType", "integerType", "realType", "booleanType", "tokenType"),
									},
									MarkdownDescription: `app configuration key type.`,
								},
								"app_config_key_value": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: `app configuration key value.`,
								},
							},
						},
						MarkdownDescription: `app configuration setting items. / Contains properties for App configuration setting item.`,
					},
				},
				Validators: []validator.Object{
					managedDeviceMobileAppConfigurationManagedDeviceMobileAppConfigurationValidator,
				},
				MarkdownDescription: `Contains properties, inherited properties and actions for iOS mobile app configurations.`,
			},
		},
	},
	MarkdownDescription: `An abstract class for Mobile app configuration for enrolled devices.`,
}

var managedDeviceMobileAppConfigurationManagedDeviceMobileAppConfigurationValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("android_managed_store"),
	path.MatchRelative().AtParent().AtName("ios"),
)
