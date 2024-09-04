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
)

var (
	CloudPcUserSettingResource = generic.GenericResource{
		TypeNameSuffix: "cloud_pc_user_setting",
		SpecificSchema: cloudPcUserSettingResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:         "/deviceManagement/virtualEndpoint/userSettings",
			ReadOptions:     cloudPcUserSettingReadOptions,
			WriteSubActions: cloudPcUserSettingWriteSubActions,
		},
	}

	CloudPcUserSettingSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&CloudPcUserSettingResource)

	CloudPcUserSettingPluralDataSource = generic.CreateGenericDataSourcePluralFromSingular(
		&CloudPcUserSettingSingularDataSource, "")
)

var cloudPcUserSettingReadOptions = generic.ReadOptions{
	ODataExpand:           "assignments",
	PluralNoFilterSupport: true,
}

var cloudPcUserSettingWriteSubActions = []generic.WriteSubAction{
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			Attributes: []string{"assignments"},
			UriSuffix:  "assign",
		},
	},
}

var cloudPcUserSettingResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // cloudPcUserSetting
		"id": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"created_date_time": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"cross_region_disaster_recovery_setting": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // cloudPcCrossRegionDisasterRecoverySetting
				"cross_region_disaster_recovery_enabled": schema.BoolAttribute{
					Required: true,
				},
				"disaster_recovery_network_setting": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // cloudPcDisasterRecoveryNetworkSetting
						"azure_connection": generic.OdataDerivedTypeNestedAttributeRs{
							DerivedType: "#microsoft.graph.cloudPcDisasterRecoveryAzureConnectionSetting",
							SingleNestedAttribute: schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{ // cloudPcDisasterRecoveryAzureConnectionSetting
									"on_premises_connection_id": schema.StringAttribute{
										Required: true,
									},
								},
								Validators: []validator.Object{
									cloudPcUserSettingCloudPcDisasterRecoveryNetworkSettingValidator,
								},
							},
						},
						"microsoft_hosted_network": generic.OdataDerivedTypeNestedAttributeRs{
							DerivedType: "#microsoft.graph.cloudPcDisasterRecoveryMicrosoftHostedNetworkSetting",
							SingleNestedAttribute: schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{ // cloudPcDisasterRecoveryMicrosoftHostedNetworkSetting
									"region_group": schema.StringAttribute{
										Required: true,
										Validators: []validator.String{
											stringvalidator.OneOf("default", "australia", "canada", "usCentral", "usEast", "usWest", "france", "germany", "europeUnion", "unitedKingdom", "japan", "asia", "india", "southAmerica", "euap", "usGovernment", "usGovernmentDOD", "unknownFutureValue", "norway", "switzerland", "southKorea"),
										},
									},
									"region_name": schema.StringAttribute{
										Required: true,
									},
								},
								Validators: []validator.Object{
									cloudPcUserSettingCloudPcDisasterRecoveryNetworkSettingValidator,
								},
							},
						},
					},
				},
				"maintain_cross_region_restore_point_enabled": schema.BoolAttribute{
					Optional: true,
				},
			},
		},
		"display_name": schema.StringAttribute{
			Required: true,
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"local_admin_enabled": schema.BoolAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:      true,
		},
		"reset_enabled": schema.BoolAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:      true,
		},
		"restore_point_setting": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // cloudPcRestorePointSetting
				"frequency_in_hours": schema.Int64Attribute{
					Computed:      true,
					PlanModifiers: []planmodifier.Int64{wpplanmodifier.Int64UseStateForUnknown()},
				},
				"frequency_type": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf("default", "fourHours", "sixHours", "twelveHours", "sixteenHours", "twentyFourHours", "unknownFutureValue"),
					},
					PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("twelveHours")},
					Computed:      true,
				},
				"user_restore_enabled": schema.BoolAttribute{
					Optional:      true,
					PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
					Computed:      true,
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/cloudPcRestorePointSetting?view=graph-rest-beta`,
		},
		"assignments": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // cloudPcUserSettingAssignment
					"target": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{ // cloudPcManagementAssignmentTarget
							"group": generic.OdataDerivedTypeNestedAttributeRs{
								DerivedType: "#microsoft.graph.cloudPcManagementGroupAssignmentTarget",
								SingleNestedAttribute: schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{ // cloudPcManagementGroupAssignmentTarget
										"allotment_display_name": schema.StringAttribute{
											Optional: true,
										},
										"allotment_licenses_count": schema.Int64Attribute{
											Optional: true,
										},
										"group_id": schema.StringAttribute{
											Required: true,
										},
										"service_plan_id": schema.StringAttribute{
											Optional: true,
										},
									},
									Validators: []validator.Object{
										cloudPcUserSettingCloudPcManagementAssignmentTargetValidator,
									},
									MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/cloudPcManagementGroupAssignmentTarget?view=graph-rest-beta`,
								},
							},
						},
						MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/cloudPcManagementAssignmentTarget?view=graph-rest-beta`,
					},
				},
			},
			MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/cloudPcUserSettingAssignment?view=graph-rest-beta`,
		},
	},
	MarkdownDescription: `https://learn.microsoft.com/en-us/graph/api/resources/cloudPcUserSetting?view=graph-rest-beta`,
}

var cloudPcUserSettingCloudPcDisasterRecoveryNetworkSettingValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("azure_connection"),
	path.MatchRelative().AtParent().AtName("microsoft_hosted_network"),
)

var cloudPcUserSettingCloudPcManagementAssignmentTargetValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("group"),
)
