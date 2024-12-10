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
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Key of the entity.",
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The date and time the policy was created.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "The policy's description. The _provider_ default value is `\"\"`.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Policy display name.",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Last time the policy was modified.",
		},
		"role_scope_tag_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: "List of Scope Tags for this Entity instance. The _provider_ default value is `[\"0\"]`.",
		},
		"version": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Version of the entity.",
		},
		"custom_settings": schema.SetNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // keyValuePair
					"name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Name for this key-value pair",
					},
					"value": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Value for this key-value pair",
					},
				},
			},
			MarkdownDescription: "A set of string key and string value pairs to be sent to apps for users to whom the configuration is scoped, unalterned by this service / Key value pair for storing custom settings / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyvaluepair?view=graph-rest-beta",
		},
		"app_group_type": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf("selectedPublicApps", "allCoreMicrosoftApps", "allMicrosoftApps", "allApps"),
			},
			MarkdownDescription: "Public Apps selection: group or individual. / Indicates a collection of apps to target which can be one of several pre-defined lists of apps or a manually selected list of apps; possible values are: `selectedPublicApps` (Target the collection of apps manually selected by the admin.), `allCoreMicrosoftApps` (Target the core set of Microsoft apps (Office, Edge, etc).), `allMicrosoftApps` (Target all apps with Microsoft as publisher.), `allApps` (Target all apps with an available assignment.)",
		},
		"deployed_app_count": schema.Int64Attribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
			MarkdownDescription: "Count of apps to which the current policy is deployed.",
		},
		"is_assigned": schema.BoolAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
			MarkdownDescription: "Indicates if the policy is deployed to any inclusion groups or not.",
		},
		"targeted_app_management_levels": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				wpvalidator.FlagEnumValues("unspecified", "unmanaged", "mdm", "androidEnterprise", "androidEnterpriseDedicatedDevicesWithAzureAdSharedMode", "androidOpenSourceProjectUserAssociated", "androidOpenSourceProjectUserless", "unknownFutureValue"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("unspecified")},
			Computed:            true,
			MarkdownDescription: "The intended app management levels for this policy. / Management levels for apps; possible values are: `unspecified` (Unspecified), `unmanaged` (Unmanaged), `mdm` (MDM), `androidEnterprise` (Android Enterprise), `androidEnterpriseDedicatedDevicesWithAzureAdSharedMode` (Android Enterprise dedicated devices with Azure AD Shared mode), `androidOpenSourceProjectUserAssociated` (Android Open Source Project (AOSP) devices), `androidOpenSourceProjectUserless` (Android Open Source Project (AOSP) userless devices), `unknownFutureValue` (Place holder for evolvable enum). The _provider_ default value is `\"unspecified\"`.",
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
											MarkdownDescription: "The identifier for an app, as specified in the play store.",
										},
									},
									Validators: []validator.Object{
										targetedManagedAppConfigurationMobileAppIdentifierValidator,
									},
									MarkdownDescription: "The identifier for an Android app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-androidmobileappidentifier?view=graph-rest-beta",
								},
							},
							"ios": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.iosMobileAppIdentifier",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // iosMobileAppIdentifier
										"bundle_id": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "The identifier for an app, as specified in the app store.",
										},
									},
									Validators: []validator.Object{
										targetedManagedAppConfigurationMobileAppIdentifierValidator,
									},
									MarkdownDescription: "The identifier for an iOS app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-iosmobileappidentifier?view=graph-rest-beta",
								},
							},
							"mac": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.macAppIdentifier",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // macAppIdentifier
										"bundle_id": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "The identifier for an app, as specified in the app store.",
										},
									},
									Validators: []validator.Object{
										targetedManagedAppConfigurationMobileAppIdentifierValidator,
									},
									MarkdownDescription: "The identifier for a Mac app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-macappidentifier?view=graph-rest-beta",
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
											MarkdownDescription: "The identifier for an app, as specified in the app store.",
										},
									},
									Validators: []validator.Object{
										targetedManagedAppConfigurationMobileAppIdentifierValidator,
									},
									MarkdownDescription: "The identifier for a Windows app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-windowsappidentifier?view=graph-rest-beta",
								},
							},
						},
						Description:         `mobileAppIdentifier`, // custom MS Graph attribute name
						MarkdownDescription: "The identifier for an app with it's operating system type. / The identifier for a mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-mobileappidentifier?view=graph-rest-beta",
					},
				},
			},
			PlanModifiers: []planmodifier.Set{
				&wpplanmodifier.IgnoreOnOtherAttributeValuePlanModifier{OtherAttributePath: path.Root("app_group_type"), ValuesRespect: []attr.Value{types.StringValue("selectedPublicApps")}},
			},
			MarkdownDescription: "List of apps to which the policy is deployed. / The identifier for the deployment an app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-managedmobileapp?view=graph-rest-beta",
		},
		"assignments": deviceAndAppManagementAssignment,
	},
	MarkdownDescription: "Configuration used to deliver a set of custom settings as-is to all users in the targeted security group / https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-targetedmanagedappconfiguration?view=graph-rest-beta ||| MS Graph: App management",
}

var targetedManagedAppConfigurationMobileAppIdentifierValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("android"),
	path.MatchRelative().AtParent().AtName("ios"),
	path.MatchRelative().AtParent().AtName("mac"),
	path.MatchRelative().AtParent().AtName("windows"),
)
