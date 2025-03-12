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

	CloudPcProvisioningPolicyPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&CloudPcProvisioningPolicyResource, "")
)

var cloudPcProvisioningPolicyReadOptions = generic.ReadOptions{
	ODataExpand: "assignments",
	DataSource: generic.DataSourceOptions{
		NoFilterSupport: true,
		Plural: generic.PluralOptions{
			ExtraAttributes: []string{"image_id", "image_type", "managed_by", "provisioning_type"},
		},
	},
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
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The unique identifier associated with the provisioning policy. This ID is auto populated during the creation of a new provisioning policy. Read-only. Supports `$filter`, `$select`, and `$orderBy`.",
		},
		"alternate_resource_url": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The URL of the alternate resource that links to this provisioning policy. Read-only.",
		},
		"autopatch": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // cloudPcProvisioningPolicyAutopatch
				"autopatch_group_id": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "The unique identifier (ID) of a Windows Autopatch group. An Autopatch group is a logical container or unit that groups several Microsoft Entra groups and software update policies. Devices with the same Autopatch group ID share unified software update management. The default value is `null` that indicates that no Autopatch group is associated with the provisioning policy.",
				},
			},
			MarkdownDescription: "The specific settings for Windows Autopatch that enable its customers to experience it on Cloud PC. The settings take effect when the tenant enrolls in Windows Autopatch and the **managedType** of the **microsoftManagedDesktop** property is set as `starterManaged`. Supports `$select`. / Represents specific settings for Windows Autopatch that enable its customers to experience it on Cloud PC. / https://learn.microsoft.com/en-us/graph/api/resources/cloudpcprovisioningpolicyautopatch?view=graph-rest-beta",
		},
		"autopilot_configuration": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // cloudPcAutopilotConfiguration
				"application_timeout_in_minutes": schema.Int64Attribute{
					Required:            true,
					MarkdownDescription: "Indicates the number of minutes allowed for the Autopilot application to apply the device preparation profile (DPP) configurations to the device. If the Autopilot application doesn't finish within the specified time (**applicationTimeoutInMinutes**), the application error is added to the **statusDetail** property of the [cloudPC](../resources/cloudpc.md) object. The supported value is an integer between 10 and 360. Required.",
				},
				"device_preparation_profile_id": schema.StringAttribute{
					Required:            true,
					MarkdownDescription: "The unique identifier (ID) of the Autopilot device preparation profile (DPP) that links a Windows Autopilot device preparation policy to ensure that devices are ready for users after provisioning. Required.",
				},
				"on_failure_device_access_denied": schema.BoolAttribute{
					Required:            true,
					MarkdownDescription: "Indicates whether the access to the device is allowed when the application of Autopilot device preparation profile (DPP) configurations fails or times out. If `true`, the **status** of the device is `failed` and the device is unable to access; otherwise, the **status** of the device is `provisionedWithWarnings` and the device is allowed to access. The default value is `false`. Required.",
				},
			},
			MarkdownDescription: "The specific settings for Windows Autopilot that enable Windows 365 customers to experience it on Cloud PC. Supports `$select`. / Represents specific settings for Windows Autopilot that enable Windows 365 customers to experience it on Cloud PC. / https://learn.microsoft.com/en-us/graph/api/resources/cloudpcautopilotconfiguration?view=graph-rest-beta",
		},
		"cloud_pc_group_display_name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The display name of the Cloud PC group that the Cloud PCs reside in. Read-only.",
		},
		"cloud_pc_naming_template": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The template used to name Cloud PCs provisioned using this policy. The template can contain custom text and replacement tokens, including `%USERNAME:x%` and `%RAND:x%`, which represent the user's name and a randomly generated number, respectively. For example, `CPC-%USERNAME:4%-%RAND:5%` means that the name of the Cloud PC starts with `CPC-`, followed by a four-character username, a `-` character, and then five random characters. The total length of the text generated by the template can't exceed 15 characters. Supports `$filter`, `$select`, and `$orderby`.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "The provisioning policy description. Supports `$filter`, `$select`, and `$orderBy`. The _provider_ default value is `\"\"`.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The display name for the provisioning policy.",
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
						MarkdownDescription: "Specifies the method by which the provisioned Cloud PC joins Microsoft Entra ID. If you choose the `hybridAzureADJoin` type, only provide a value for the **onPremisesConnectionId** property and leave the **regionName** property empty. If you choose the `azureADJoin` type, provide a value for either the **onPremisesConnectionId** or the **regionName** property. / Possible values are: `azureADJoin`, `hybridAzureADJoin`, `unknownFutureValue`. The _provider_ default value is `\"azureADJoin\"`.",
					},
					"on_premises_connection_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The Azure network connection ID that matches the virtual network IT admins want the provisioning policy to use when they create Cloud PCs. You can use this property in both domain join types: _Azure AD joined_ or _Hybrid Microsoft Entra joined_. If you enter an **onPremisesConnectionId**, leave the **regionName** property empty.",
					},
					"region_group": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("default", "australia", "canada", "usCentral", "usEast", "usWest", "france", "germany", "europeUnion", "unitedKingdom", "japan", "asia", "india", "southAmerica", "euap", "usGovernment", "usGovernmentDOD", "unknownFutureValue", "norway", "switzerland", "southKorea", "middleEast", "mexico", "australasia", "europe"),
						},
						MarkdownDescription: "The logical geographic group this region belongs to. Multiple regions can belong to one region group. A customer can select a **regionGroup** when they provision a Cloud PC, and the Cloud PC is put in one of the regions in the group based on resource status. For example, the Europe region group contains the Northern Europe and Western Europe regions.and `southKorea`. Read-only. / Possible values are: `default`, `australia`, `canada`, `usCentral`, `usEast`, `usWest`, `france`, `germany`, `europeUnion`, `unitedKingdom`, `japan`, `asia`, `india`, `southAmerica`, `euap`, `usGovernment`, `usGovernmentDOD`, `unknownFutureValue`, `norway`, `switzerland`, `southKorea`, `middleEast`, `mexico`, `australasia`, `europe`",
					},
					"region_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The supported Azure region where the IT admin wants the provisioning policy to create Cloud PCs. The underlying virtual network is created and managed by the Windows 365 service. This can only be entered if the IT admin chooses Microsoft Entra joined as the domain join type. If you enter a **regionName**, leave the **onPremisesConnectionId** property empty.",
					},
					"type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("azureADJoin", "hybridAzureADJoin", "unknownFutureValue"),
						},
						PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("azureADJoin")},
						Computed:            true,
						MarkdownDescription: "Possible values are: `azureADJoin`, `hybridAzureADJoin`, `unknownFutureValue`. The _provider_ default value is `\"azureADJoin\"`.",
					},
				},
			},
			MarkdownDescription: "Specifies a list ordered by priority on how Cloud PCs join Microsoft Entra ID (Azure AD). Supports `$select`. / Represents a defined configuration of how a provisioned Cloud PC device joins to Microsoft Entra ID. / https://learn.microsoft.com/en-us/graph/api/resources/cloudpcdomainjoinconfiguration?view=graph-rest-beta",
		},
		"enable_single_sign_on": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "`True` if single sign-on can access the provisioned Cloud PC. `False` indicates that the provisioned Cloud PC doesn't support this feature. The default value is `false`. Windows 365 users can use single sign-on to authenticate to Microsoft Entra ID with passwordless options (for example, FIDO keys) to access their Cloud PC. Optional. The _provider_ default value is `false`.",
		},
		"grace_period_in_hours": schema.Int64Attribute{
			Optional:            true,
			MarkdownDescription: "The number of hours to wait before reprovisioning/deprovisioning happens. Read-only.",
		},
		"image_display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The display name of the operating system image that is used for provisioning. For example, `Windows 11 Preview + Microsoft 365 Apps 23H2 23H2`. Supports `$filter`, `$select`, and `$orderBy`.",
		},
		"image_id": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The unique identifier that represents an operating system image that is used for provisioning new Cloud PCs. The format for a gallery type image is: {publisherName_offerName_skuName}. Supported values for each of the parameters are:<ul><li>publisher: `Microsoftwindowsdesktop`</li> <li>offer: `windows-ent-cpc`</li> <li>sku: `21h1-ent-cpc-m365`, `21h1-ent-cpc-os`, `20h2-ent-cpc-m365`, `20h2-ent-cpc-os`, `20h1-ent-cpc-m365`, `20h1-ent-cpc-os`, `19h2-ent-cpc-m365`, and `19h2-ent-cpc-os`</li></ul> Supports `$filter`, `$select`, and `$orderBy`.",
		},
		"image_type": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("gallery", "custom", "unknownFutureValue"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("gallery")},
			Computed:            true,
			MarkdownDescription: "The type of operating system image (custom or gallery) that is used for provisioning on Cloud PCs. The default value is `gallery`. Supports $filter, $select, and $orderBy. / Possible values are: `gallery`, `custom`, `unknownFutureValue`. The _provider_ default value is `\"gallery\"`.",
		},
		"local_admin_enabled": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "When `true`, the local admin is enabled for Cloud PCs; `false` indicates that the local admin isn't enabled for Cloud PCs. The default value is `false`. Supports `$filter`, `$select`, and `$orderBy`.",
		},
		"managed_by": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				wpvalidator.FlagEnumValues("windows365", "devBox", "unknownFutureValue", "rpaBox"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("windows365")},
			Computed:            true,
			MarkdownDescription: "Indicates the service that manages the provisioning policy. The default value is `windows365`. Use the `Prefer: include-unknown-enum-members` request header to get the following value in this [evolvable enum](/graph/best-practices-concept#handling-future-members-in-evolvable-enumerations): `rpaBox`. Supports `$filter`, `$select`, and `$orderBy`. / Possible values are: `windows365`, `devBox`, `unknownFutureValue`, `rpaBox`. The _provider_ default value is `\"windows365\"`.",
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
					MarkdownDescription: "Indicates the provisioning policy associated with Microsoft Managed Desktop settings. The default value is `notManaged`. / Possible values are: `notManaged`, `premiumManaged`, `standardManaged`, `starterManaged`, `unknownFutureValue`. The _provider_ default value is `\"notManaged\"`.",
				},
				"profile": schema.StringAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
					Computed:            true,
					MarkdownDescription: "The name of the Microsoft Managed Desktop profile that the Windows 365 Cloud PC is associated with. The _provider_ default value is `\"\"`.",
				},
				"type": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf("notManaged", "premiumManaged", "standardManaged", "starterManaged", "unknownFutureValue"),
					},
					PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notManaged")},
					Computed:            true,
					MarkdownDescription: "Possible values are: `notManaged`, `premiumManaged`, `standardManaged`, `starterManaged`, `unknownFutureValue`. The _provider_ default value is `\"notManaged\"`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "The specific settings to **microsoftManagedDesktop** that enables Microsoft Managed Desktop customers to get device managed experience for Cloud PC. To enable **microsoftManagedDesktop** to provide more value, an admin needs to specify certain settings in it. Supports `$filter`, `$select`, and `$orderBy`. / Represents specific settings for the Microsoft Managed Desktop that enables customers to get a managed device experience for a Cloud PC. / https://learn.microsoft.com/en-us/graph/api/resources/microsoftmanageddesktop?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
		"provisioning_type": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("dedicated", "shared", "unknownFutureValue", "sharedByUser", "sharedByEntraGroup"),
			},
			MarkdownDescription: "Specifies the type of licenses to be used when provisioning Cloud PCs using this policy. The possible values are `dedicated`, `shared`, `unknownFutureValue`, `sharedByUser`, `sharedByEntraGroup`. Use the `Prefer: include-unknown-enum-members` request header to get the following values from this [evolvable enum](/graph/best-practices-concept#handling-future-members-in-evolvable-enumerations): `sharedByUser`, `sharedByEntraGroup`. The `shared` member is deprecated and will stop returning on April 30, 2027; going forward, use the `sharedByUser` member. For example, a `dedicated` service plan can be assigned to only one user and provision only one Cloud PC. The `shared` and `sharedByUser` plans require customers to purchase a shared service plan. Each shared license purchased can enable up to three Cloud PCs, with only one user signed in at a time. The `sharedByEntraGroup` plan also requires the purchase of a shared service plan. Each shared license under this plan can enable one Cloud PC, which is shared for the group according to the assignments of this policy. By default, the license type is `dedicated` if the **provisioningType** isn't specified when you create the **cloudPcProvisioningPolicy**. You can't change this property after the **cloudPcProvisioningPolicy** is created. / Possible values are: `dedicated`, `shared`, `unknownFutureValue`, `sharedByUser`, `sharedByEntraGroup`",
		},
		"scope_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: ". The _provider_ default value is `[\"0\"]`.",
		},
		"windows_setting": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // cloudPcWindowsSetting
				"locale": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "The Windows language or region tag to use for language pack configuration and localization of the Cloud PC. The default value is `en-US`, which corresponds to English (United States).",
				},
			},
			MarkdownDescription: "Indicates a specific Windows setting to configure during the creation of Cloud PCs for this provisioning policy. Supports `$select`. / Represents a specific Windows setting to configure during the creation of Cloud PCs for a provisioning policy. / https://learn.microsoft.com/en-us/graph/api/resources/cloudpcwindowssetting?view=graph-rest-beta",
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
											Required:            true,
											MarkdownDescription: "The ID of the target group for the assignment.",
										},
										"service_plan_id": schema.StringAttribute{
											Optional:            true,
											MarkdownDescription: "The unique identifier for the service plan that indicates which size of the Cloud PC to provision for the user. Use a `null` value, when the **provisioningType** is `dedicated`.",
										},
									},
									Validators: []validator.Object{
										cloudPcProvisioningPolicyCloudPcManagementAssignmentTargetValidator,
									},
									MarkdownDescription: "Complex type that represents the assignment target group. / https://learn.microsoft.com/en-us/graph/api/resources/cloudpcmanagementgroupassignmenttarget?view=graph-rest-beta",
								},
							},
						},
						MarkdownDescription: "The assignment target for the provisioning policy. Currently, the only target supported for this policy is a user group. For details, see [cloudPcManagementGroupAssignmentTarget](cloudpcmanagementgroupassignmenttarget.md). / Represents an abstract base type for assignment targets.\n\nBase type of [cloudPcManagementGroupAssignmentTarget](cloudpcmanagementgroupassignmenttarget.md). / https://learn.microsoft.com/en-us/graph/api/resources/cloudpcmanagementassignmenttarget?view=graph-rest-beta",
					},
				},
			},
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "A defined collection of provisioning policy assignments. Represents the set of Microsoft 365 groups and security groups in Microsoft Entra ID that have provisioning policy assigned. Returned only on `$expand`. For an example about how to get the assignments relationship, see [Get cloudPcProvisioningPolicy](../api/cloudpcprovisioningpolicy-get.md). / Represents a defined collection of provisioning policy assignments. / https://learn.microsoft.com/en-us/graph/api/resources/cloudpcprovisioningpolicyassignment?view=graph-rest-beta. The _provider_ default value is `[]`.",
		},
	},
	MarkdownDescription: "Represents a Cloud PC provisioning policy. / https://learn.microsoft.com/en-us/graph/api/resources/cloudpcprovisioningpolicy?view=graph-rest-beta ||| MS Graph: Cloud PC",
}

var cloudPcProvisioningPolicyCloudPcManagementAssignmentTargetValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("group"),
)
