package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	TargetedManagedAppConfigurationResource = generic.GenericResource{
		TypeNameSuffix: "targeted_managed_app_configuration",
		SpecificSchema: targetedManagedAppConfigurationResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:         "/deviceAppManagement/targetedManagedAppConfigurations",
			ReadOptions:     targetedManagedAppConfigurationReadOptions,
			WriteSubActions: targetedManagedAppConfigurationWriteSubActions,
		},
	}

	TargetedManagedAppConfigurationSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&TargetedManagedAppConfigurationResource)

	TargetedManagedAppConfigurationPluralDataSource = generic.CreateGenericDataSourcePluralFromSingular(
		&TargetedManagedAppConfigurationSingularDataSource, "")
)

var targetedManagedAppConfigurationReadOptions = generic.ReadOptions{
	ODataExpand: "apps,assignments",
}

var targetedManagedAppConfigurationWriteSubActions = []generic.WriteSubAction{
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			Attributes: []string{"appGroupType", "apps"},
			UriSuffix:  "targetApps",
			UpdateOnly: true,
		},
	},
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			Attributes: []string{"assignments"},
			UriSuffix:  "assign",
			UpdateOnly: true,
		},
	},
}

var targetedManagedAppConfigurationResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // targetedManagedAppConfiguration
		"id": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"created_date_time": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"description": schema.StringAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:      true,
		},
		"display_name": schema.StringAttribute{
			Required: true,
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:   types.StringType,
			Optional:      true,
			PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:      true,
		},
		"version": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"custom_settings": schema.SetNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // keyValuePair
					"name": schema.StringAttribute{
						Required: true,
					},
					"value": schema.StringAttribute{
						Optional: true,
					},
				},
			},
			MarkdownDescription: `A set of string key and string value pairs to be sent to apps for users to whom the configuration is scoped, unalterned by this service / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyValuePair?view=graph-rest-beta`,
		},
		"app_group_type": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf("selectedPublicApps", "allCoreMicrosoftApps", "allMicrosoftApps", "allApps"),
			},
			MarkdownDescription: `Public Apps selection: group or individual`,
		},
		"deployed_app_count": schema.Int64Attribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
			MarkdownDescription: `Count of apps to which the current policy is deployed.`,
		},
		"is_assigned": schema.BoolAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
			MarkdownDescription: `Indicates if the policy is deployed to any inclusion groups or not.`,
		},
		"targeted_app_management_levels": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				wpvalidator.FlagEnumValues("unspecified", "unmanaged", "mdm", "androidEnterprise", "androidEnterpriseDedicatedDevicesWithAzureAdSharedMode", "androidOpenSourceProjectUserAssociated", "androidOpenSourceProjectUserless", "unknownFutureValue"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unspecified")},
			Computed:            true,
			MarkdownDescription: `The intended app management levels for this policy`,
		},
		"apps": schema.SetNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // managedMobileApp
					"app_id": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{ // mobileAppIdentifier
							"android": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.androidMobileAppIdentifier",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // androidMobileAppIdentifier
										"package_id": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: `The identifier for an app, as specified in the play store.`,
										},
									},
									Validators: []validator.Object{
										targetedManagedAppConfigurationMobileAppIdentifierValidator,
									},
									MarkdownDescription: `The identifier for an Android app.`,
								},
							},
							"ios": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.iosMobileAppIdentifier",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // iosMobileAppIdentifier
										"bundle_id": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: `The identifier for an app, as specified in the app store.`,
										},
									},
									Validators: []validator.Object{
										targetedManagedAppConfigurationMobileAppIdentifierValidator,
									},
									MarkdownDescription: `The identifier for an iOS app.`,
								},
							},
							"mac": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.macAppIdentifier",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // macAppIdentifier
										"bundle_id": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: `The identifier for an app, as specified in the app store.`,
										},
									},
									Validators: []validator.Object{
										targetedManagedAppConfigurationMobileAppIdentifierValidator,
									},
									MarkdownDescription: `The identifier for a Mac app.`,
								},
							},
							"windows": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.windowsAppIdentifier",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // windowsAppIdentifier
										"app_id": schema.StringAttribute{
											Required:            true,
											Description:         `windowsAppId`, // custom MS Graph attribute name
											MarkdownDescription: `The identifier for an app, as specified in the app store.`,
										},
									},
									Validators: []validator.Object{
										targetedManagedAppConfigurationMobileAppIdentifierValidator,
									},
									MarkdownDescription: `The identifier for a Windows app.`,
								},
							},
						},
						Description:         `mobileAppIdentifier`, // custom MS Graph attribute name
						MarkdownDescription: `The identifier for an app with it's operating system type. / The identifier for a mobile app.`,
					},
				},
			},
			PlanModifiers: []planmodifier.Set{
				&wpplanmodifier.IgnoreOnOtherAttributeValuePlanModifier{OtherAttributePath: path.Root("app_group_type"), ValuesRespect: []attr.Value{types.StringValue("selectedPublicApps")}},
			},
			MarkdownDescription: `List of apps to which the policy is deployed. / The identifier for the deployment an app.`,
		},
		"assignments": deviceAndAppManagementAssignment,
	},
	MarkdownDescription: `Configuration used to deliver a set of custom settings as-is to all users in the targeted security group / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-targetedManagedAppConfiguration?view=graph-rest-beta`,
}

var targetedManagedAppConfigurationMobileAppIdentifierValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("android"),
	path.MatchRelative().AtParent().AtName("ios"),
	path.MatchRelative().AtParent().AtName("mac"),
	path.MatchRelative().AtParent().AtName("windows"),
)
