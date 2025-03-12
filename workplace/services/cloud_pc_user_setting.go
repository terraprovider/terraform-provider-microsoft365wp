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

	CloudPcUserSettingPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&CloudPcUserSettingResource, "")
)

var cloudPcUserSettingReadOptions = generic.ReadOptions{
	ODataExpand: "assignments",
	DataSource: generic.DataSourceOptions{
		NoFilterSupport: true,
	},
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
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Unique identifier for the Cloud PC user setting. Read-only.",
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The date and time the setting was created. The timestamp type represents the date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 looks like this: '2014-01-01T00:00:00Z'.",
		},
		"cross_region_disaster_recovery_setting": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // cloudPcCrossRegionDisasterRecoverySetting
				"cross_region_disaster_recovery_enabled": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
					Computed:            true,
					MarkdownDescription: ". The _provider_ default value is `false`.",
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
										Required:            true,
										MarkdownDescription: "Indicates the unique ID of the virtual network that the new Cloud PC joins.",
									},
								},
								Validators: []validator.Object{
									cloudPcUserSettingCloudPcDisasterRecoveryNetworkSettingValidator,
								},
								MarkdownDescription: "Represents the Azure network connection configuration of backup Cloud PCs provisioned for cross-region disaster recovery. / https://learn.microsoft.com/en-us/graph/api/resources/cloudpcdisasterrecoveryazureconnectionsetting?view=graph-rest-beta",
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
											stringvalidator.OneOf("default", "australia", "canada", "usCentral", "usEast", "usWest", "france", "germany", "europeUnion", "unitedKingdom", "japan", "asia", "india", "southAmerica", "euap", "usGovernment", "usGovernmentDOD", "unknownFutureValue", "norway", "switzerland", "southKorea", "middleEast", "mexico", "australasia", "europe"),
										},
										MarkdownDescription: "Indicates the logic geographic group this region belongs to. Multiple regions can belong to one region group. When a region group is configured for disaster recovery, the new Cloud PC is assigned to one of the regions within the group based on resource availability. For example, the `europeUnion` region group contains the North Europe and West Europe regions.，`southKorea`, `middleEast`, `mexico`. Use the `Prefer: include-unknown-enum-members` request header to get the following values in this [evolvable enum](/graph/best-practices-concept#handling-future-members-in-evolvable-enumerations): `norway`, `switzerland`，`southKorea`, `middleEast`, `mexico`. / Possible values are: `default`, `australia`, `canada`, `usCentral`, `usEast`, `usWest`, `france`, `germany`, `europeUnion`, `unitedKingdom`, `japan`, `asia`, `india`, `southAmerica`, `euap`, `usGovernment`, `usGovernmentDOD`, `unknownFutureValue`, `norway`, `switzerland`, `southKorea`, `middleEast`, `mexico`, `australasia`, `europe`",
									},
									"region_name": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "Indicates the Azure region that the new Cloud PC is assigned to. The Windows 365 service creates and manages the underlying virtual network.",
									},
								},
								Validators: []validator.Object{
									cloudPcUserSettingCloudPcDisasterRecoveryNetworkSettingValidator,
								},
								MarkdownDescription: "Represents the configuration of Microsoft-hosted network for backup Cloud PCs provisioned for cross-region disaster recovery. / https://learn.microsoft.com/en-us/graph/api/resources/cloudpcdisasterrecoverymicrosofthostednetworksetting?view=graph-rest-beta",
							},
						},
					},
					MarkdownDescription: "Indicates the network settings of the Cloud PC during a cross-region disaster recovery operation. / An abstract type that represents the network configuration of backup Cloud PCs provisioned for cross-region disaster recovery.\n\nBase type of [cloudPcDisasterRecoveryAzureConnectionSetting](../resources/cloudpcdisasterrecoveryazureconnectionsetting.md) and [cloudPcDisasterRecoveryMicrosoftHostedNetworkSetting](../resources/cloudpcdisasterrecoverymicrosofthostednetworksetting.md) / https://learn.microsoft.com/en-us/graph/api/resources/cloudpcdisasterrecoverynetworksetting?view=graph-rest-beta",
				},
				"disaster_recovery_type": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf("notConfigured", "crossRegion", "premium", "unknownFutureValue"),
					},
					PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notConfigured")},
					Computed:            true,
					MarkdownDescription: "Indicates the type of disaster recovery to perform when a disaster occurs on the user's Cloud PC. The The default value is `notConfigured`. / Possible values are: `notConfigured`, `crossRegion`, `premium`, `unknownFutureValue`. The _provider_ default value is `\"notConfigured\"`.",
				},
				"maintain_cross_region_restore_point_enabled": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
					Computed:            true,
					MarkdownDescription: "Indicates whether Windows 365 maintain the cross-region disaster recovery function generated restore points. If `true`, the Windows 365 stored restore points; `false` indicates that Windows 365 doesn't generate or keep the restore point from the original Cloud PC. If a disaster occurs, the new Cloud PC can only be provisioned using the initial image. This limitation can result in the loss of some user data on the original Cloud PC. The default value is `false`. The _provider_ default value is `true`.",
				},
				"user_initiated_disaster_recovery_allowed": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
					Computed:            true,
					MarkdownDescription: "Indicates whether the client allows the end user to initiate a disaster recovery activation. `True` indicates that the client includes the option for the end user to activate Backup Cloud PC. When `false`, the end user doesn't have the option to activate disaster recovery. The default value is `false`. Currently, only premium disaster recovery is supported. The _provider_ default value is `false`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Defines whether the user's Cloud PC enables cross-region disaster recovery and specifies the network for the disaster recovery. / Represents the settings for cross-region disaster recovery on a Cloud PC. / https://learn.microsoft.com/en-us/graph/api/resources/cloudpccrossregiondisasterrecoverysetting?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The setting name displayed in the user interface.",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The last date and time the setting was modified. The timestamp type represents the date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 looks like this: '2014-01-01T00:00:00Z'.",
		},
		"local_admin_enabled": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Indicates whether the local admin option is enabled. Default value is `false`. To enable the local admin option, change the setting to `true`. If the local admin option is enabled, the end user can be an admin of the Cloud PC device. The _provider_ default value is `false`.",
		},
		"notification_setting": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // cloudPcNotificationSetting
				"restart_prompts_disabled": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
					Computed:            true,
					MarkdownDescription: "If `true`, doesn't prompt the user to restart the Cloud PC. If `false`, prompts the user to restart Cloud PC. The default value is `false`. The _provider_ default value is `false`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Defines the setting of the Cloud PC notification prompts for the Cloud PC user. / Represents the settings of a point-in-time restore of a Cloud PC. / https://learn.microsoft.com/en-us/graph/api/resources/cloudpcnotificationsetting?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
		"reset_enabled": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Indicates whether an end user is allowed to reset their Cloud PC. When `true`, the user is allowed to reset their Cloud PC. When `false`, end-user initiated reset isn't allowed. The default value is `false`. The _provider_ default value is `false`.",
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
					PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("twelveHours")},
					Computed:            true,
					MarkdownDescription: "The time interval in hours to take snapshots (restore points) of a Cloud PC automatically. The default value is `default` that indicates that the time interval for automatic capturing of restore point snapshots is set to 12 hours. / Possible values are: `default`, `fourHours`, `sixHours`, `twelveHours`, `sixteenHours`, `twentyFourHours`, `unknownFutureValue`. The _provider_ default value is `\"twelveHours\"`.",
				},
				"user_restore_enabled": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
					Computed:            true,
					MarkdownDescription: "If `true`, the user has the ability to use snapshots to restore Cloud PCs. If `false`, non-admin users can't use snapshots to restore the Cloud PC. The _provider_ default value is `false`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Defines how frequently a restore point is created that is, a snapshot is taken) for users' provisioned Cloud PCs (default is 12 hours), and whether the user is allowed to restore their own Cloud PCs to a backup made at a specific point in time. / Represents the settings of a point-in-time restore of a Cloud PC. / https://learn.microsoft.com/en-us/graph/api/resources/cloudpcrestorepointsetting?view=graph-rest-beta. The _provider_ default value is `{}`.",
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
											Required:            true,
											MarkdownDescription: "The ID of the target group for the assignment.",
										},
										"service_plan_id": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "The unique identifier for the service plan that indicates which size of the Cloud PC to provision for the user. Use a `null` value, when the **provisioningType** is `dedicated`.",
										},
									},
									Validators: []validator.Object{
										cloudPcUserSettingCloudPcManagementAssignmentTargetValidator,
									},
									MarkdownDescription: "Complex type that represents the assignment target group. / https://learn.microsoft.com/en-us/graph/api/resources/cloudpcmanagementgroupassignmenttarget?view=graph-rest-beta",
								},
							},
						},
						MarkdownDescription: "The assignment target for the user setting. Currently, the only target supported for this user setting is a user group. For details, see [cloudPcManagementGroupAssignmentTarget](cloudpcmanagementgroupassignmenttarget.md). / Represents an abstract base type for assignment targets.\n\nBase type of [cloudPcManagementGroupAssignmentTarget](cloudpcmanagementgroupassignmenttarget.md). / https://learn.microsoft.com/en-us/graph/api/resources/cloudpcmanagementassignmenttarget?view=graph-rest-beta",
					},
				},
			},
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Represents the set of Microsoft 365 groups and security groups in Microsoft Entra ID that have **cloudPCUserSetting** assigned. Returned only on `$expand`. For an example, see [Get cloudPcUserSettingample](../api/cloudpcusersetting-get.md). / Represents a defined collection of user setting assignments. / https://learn.microsoft.com/en-us/graph/api/resources/cloudpcusersettingassignment?view=graph-rest-beta. The _provider_ default value is `[]`.",
		},
	},
	MarkdownDescription: "Represents a Cloud PC user setting. / https://learn.microsoft.com/en-us/graph/api/resources/cloudpcusersetting?view=graph-rest-beta ||| MS Graph: Cloud PC",
}

var cloudPcUserSettingCloudPcDisasterRecoveryNetworkSettingValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("azure_connection"),
	path.MatchRelative().AtParent().AtName("microsoft_hosted_network"),
)

var cloudPcUserSettingCloudPcManagementAssignmentTargetValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("group"),
)
