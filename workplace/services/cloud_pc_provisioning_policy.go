package services

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
	CloudPcProvisioningPolicyResource = generic.GenericResource{
		TypeNameSuffix: "cloud_pc_provisioning_policy",
		SpecificSchema: cloudPcProvisioningPolicyResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:         "/deviceManagement/virtualEndpoint/provisioningPolicies",
			ReadOptions:     cloudPcProvisioningPolicyReadOptions,
			WriteSubActions: cloudPcProvisioningPolicyWriteSubActions,
		},
	}

	CloudPcProvisioningPolicySingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&CloudPcProvisioningPolicyResource)

	CloudPcProvisioningPolicyPluralDataSource = generic.CreateGenericDataSourcePluralFromSingular(
		&CloudPcProvisioningPolicySingularDataSource, "")
)

var cloudPcProvisioningPolicyReadOptions = generic.ReadOptions{
	ODataExpand:           "assignments",
	PluralNoFilterSupport: true,
}

var cloudPcProvisioningPolicyWriteSubActions = []generic.WriteSubAction{
	&generic.WriteSubActionAllInOne{
		WriteSubActionBase: generic.WriteSubActionBase{
			Attributes:        []string{"assignments"},
			UriSuffix:         "assign",
			EmptyBeforeDelete: true,
		},
	},
}

var cloudPcProvisioningPolicyResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // cloudPcProvisioningPolicy
		"id": schema.StringAttribute{
			Computed:      true,
			PlanModifiers: []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
		},
		"alternate_resource_url": schema.StringAttribute{
			Optional: true,
		},
		"autopatch": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // cloudPcProvisioningPolicyAutopatch
				"autopatch_group_id": schema.StringAttribute{
					Optional: true,
				},
			},
			MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/cloudPcProvisioningPolicyAutopatch?view=graph-rest-beta",
		},
		"autopilot_configuration": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // cloudPcAutopilotConfiguration
				"application_timeout_in_minutes": schema.Int64Attribute{
					Required: true,
				},
				"device_preparation_profile_id": schema.StringAttribute{
					Required: true,
				},
				"on_failure_device_access_denied": schema.BoolAttribute{
					Required: true,
				},
			},
		},
		"cloud_pc_group_display_name": schema.StringAttribute{
			Optional: true,
		},
		"cloud_pc_naming_template": schema.StringAttribute{
			Optional: true,
		},
		"description": schema.StringAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:      true,
		},
		"display_name": schema.StringAttribute{
			Required: true,
		},
		"domain_join_configurations": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // cloudPcDomainJoinConfiguration
					"domain_join_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("azureADJoin", "hybridAzureADJoin", "unknownFutureValue"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("azureADJoin")},
						Computed:            true,
						MarkdownDescription: "; possible values are: `azureADJoin`, `hybridAzureADJoin`, `unknownFutureValue`",
					},
					"on_premises_connection_id": schema.StringAttribute{
						Optional: true,
					},
					"region_group": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("default", "australia", "canada", "usCentral", "usEast", "usWest", "france", "germany", "europeUnion", "unitedKingdom", "japan", "asia", "india", "southAmerica", "euap", "usGovernment", "usGovernmentDOD", "unknownFutureValue", "norway", "switzerland", "southKorea", "middleEast", "mexico"),
						},
						MarkdownDescription: "; possible values are: `default`, `australia`, `canada`, `usCentral`, `usEast`, `usWest`, `france`, `germany`, `europeUnion`, `unitedKingdom`, `japan`, `asia`, `india`, `southAmerica`, `euap`, `usGovernment`, `usGovernmentDOD`, `unknownFutureValue`, `norway`, `switzerland`, `southKorea`, `middleEast`, `mexico`",
					},
					"region_name": schema.StringAttribute{
						Optional:      true,
						PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("automatic")},
						Computed:      true,
					},
					"type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("azureADJoin", "hybridAzureADJoin", "unknownFutureValue"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("azureADJoin")},
						Computed:            true,
						MarkdownDescription: "; possible values are: `azureADJoin`, `hybridAzureADJoin`, `unknownFutureValue`",
					},
				},
			},
			MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/cloudPcDomainJoinConfiguration?view=graph-rest-beta",
		},
		"enable_single_sign_on": schema.BoolAttribute{
			Optional:      true,
			PlanModifiers: []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:      true,
		},
		"grace_period_in_hours": schema.Int64Attribute{
			Optional: true,
		},
		"image_id": schema.StringAttribute{
			Required: true,
		},
		"image_type": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("gallery", "custom", "unknownFutureValue"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("gallery")},
			Computed:            true,
			MarkdownDescription: "; possible values are: `gallery`, `custom`, `unknownFutureValue`",
		},
		"local_admin_enabled": schema.BoolAttribute{
			Optional: true,
		},
		"managed_by": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				wpvalidator.FlagEnumValues("windows365", "devBox", "unknownFutureValue", "rpaBox"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("windows365")},
			Computed:            true,
			MarkdownDescription: "; possible values are: `windows365`, `devBox`, `unknownFutureValue`, `rpaBox`",
		},
		"microsoft_managed_desktop": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // microsoftManagedDesktop
				"managed_type": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf("notManaged", "premiumManaged", "standardManaged", "starterManaged", "unknownFutureValue"),
					},
					PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notManaged")},
					Computed:            true,
					MarkdownDescription: "; possible values are: `notManaged`, `premiumManaged`, `standardManaged`, `starterManaged`, `unknownFutureValue`",
				},
				"profile": schema.StringAttribute{
					Optional:      true,
					PlanModifiers: []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
					Computed:      true,
				},
				"type": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf("notManaged", "premiumManaged", "standardManaged", "starterManaged", "unknownFutureValue"),
					},
					PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notManaged")},
					Computed:            true,
					MarkdownDescription: "; possible values are: `notManaged`, `premiumManaged`, `standardManaged`, `starterManaged`, `unknownFutureValue`",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/microsoftManagedDesktop?view=graph-rest-beta",
		},
		"provisioning_type": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("dedicated", "shared", "unknownFutureValue", "sharedByUser", "sharedByEntraGroup"),
			},
			MarkdownDescription: "; possible values are: `dedicated`, `shared`, `unknownFutureValue`, `sharedByUser`, `sharedByEntraGroup`",
		},
		"scope_ids": schema.SetAttribute{
			ElementType:   types.StringType,
			Optional:      true,
			PlanModifiers: []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:      true,
		},
		"windows_setting": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // cloudPcWindowsSetting
				"locale": schema.StringAttribute{
					Optional: true,
				},
			},
			MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/cloudPcWindowsSetting?view=graph-rest-beta",
		},
		"assignments": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // cloudPcProvisioningPolicyAssignment
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
										cloudPcProvisioningPolicyCloudPcManagementAssignmentTargetValidator,
									},
									MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/cloudPcManagementGroupAssignmentTarget?view=graph-rest-beta",
								},
							},
						},
						MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/cloudPcManagementAssignmentTarget?view=graph-rest-beta",
					},
				},
			},
			MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/cloudPcProvisioningPolicyAssignment?view=graph-rest-beta",
		},
	},
	MarkdownDescription: "https://learn.microsoft.com/en-us/graph/api/resources/cloudPcProvisioningPolicy?view=graph-rest-beta",
}

var cloudPcProvisioningPolicyCloudPcManagementAssignmentTargetValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("group"),
)
