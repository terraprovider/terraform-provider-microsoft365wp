package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvaluemodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	CloudPcProvisioningPolicyResource = generic.GenericResource{
		TypeNameSuffix: "cloud_pc_provisioning_policy",
		SpecificSchema: cloudPcProvisioningPolicyResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/deviceManagement/virtualEndpoint/provisioningPolicies",
			ReadOptions: generic.ReadOptions{
				ODataExpand: "assignments",
				DataSource: generic.DataSourceOptions{
					NoFilterSupport: true,
					Plural: generic.PluralOptions{
						ExtraAttributes: []string{"image_id", "image_type", "managed_by", "provisioning_type"},
					},
				},
			},
			WriteOptions: generic.WriteOptions{
				SubActions: []generic.WriteSubAction{
					&generic.WriteSubActionAllInOne{
						WriteSubActionBase: generic.WriteSubActionBase{
							Attributes:        []string{"assignments"},
							UriSuffix:         "assign",
							EmptyBeforeDelete: true,
						},
					},
				},
			},
		},
	}

	CloudPcProvisioningPolicySingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&CloudPcProvisioningPolicyResource)

	CloudPcProvisioningPolicyPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&CloudPcProvisioningPolicyResource, "")
)

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
			MarkdownDescription: "Indicates the Windows Autopatch settings for Cloud PCs using this provisioning policy. The settings take effect when the tenant enrolls in Autopatch and the **managedType** of the **microsoftManagedDesktop** property is set as `starterManaged`. Supports `$select`. / Indicates the Windows Autopatch settings for Cloud PCs using this provisioning policy. Also see [Microsoft docs for cloudPcProvisioningPolicyAutopatch](https://learn.microsoft.com/en-us/graph/api/resources/cloudpcprovisioningpolicyautopatch?view=graph-rest-beta). <br> ",
		},
		"autopilot_configuration": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // cloudPcAutopilotConfiguration
				"application_timeout_in_minutes": schema.Int64Attribute{
					Required:            true,
					MarkdownDescription: "Indicates the number of minutes allowed for the Autopilot application to apply the device preparation profile (DPP) configurations to the device. If the Autopilot application doesn't finish within the specified time (**applicationTimeoutInMinutes**), the application error is added to the **statusDetail** property of the [cloudPC](https://learn.microsoft.com/en-us/graph/api/resources/cloudpc?view=graph-rest-beta) object. The supported value is an integer between `30` and `360`. Required.",
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
			MarkdownDescription: "The specific settings for Windows Autopilot that enable Windows 365 customers to experience it on Cloud PC. Supports `$select`. / Represents specific settings for Windows Autopilot that enable Windows 365 customers to experience it on Cloud PC. Also see [Microsoft docs for cloudPcAutopilotConfiguration](https://learn.microsoft.com/en-us/graph/api/resources/cloudpcautopilotconfiguration?view=graph-rest-beta). <br> ",
		},
		"cloud_pc_group_display_name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The display name of the Cloud PC group that the Cloud PCs reside in. Read-only.",
		},
		"cloud_pc_naming_template": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The template used to name Cloud PCs provisioned using this policy. The template can contain custom text and replacement tokens, including `%USERNAME:x%` and `%RAND:x%`, which represent the user's name and a randomly generated number, respectively. For example, `CPC-%USERNAME:4%-%RAND:5%` means that the name of the Cloud PC starts with `CPC-`, followed by a four-character username, a `-` character, and then five random characters. The total length of the text generated by the template can't exceed 15 characters. Supports `$filter`, `$select`, and `$orderby`.",
		},
		"created_by": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The unique ID of the user who created this policy. For example, `5ccb8d35-dd04-473e-a287-69bb4473208b`. Read-only. Supports `$select`.",
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The timestamp when this provisioning policy was created. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Read-only. Supports `$select` and `$orderBy`.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvaluemodifier.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "The provisioning policy description. Supports `$filter`, `$select`, and `$orderBy`. <br/> The _provider_ default value is `\"\"`.",
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
						PlanModifiers: []planmodifier.String{
							wpdefaultvaluemodifier.StringDefaultValue("azureADJoin"),
						},
						Computed:            true,
						MarkdownDescription: "Specifies the method by which the provisioned Cloud PC joins Microsoft Entra ID. If you choose the `hybridAzureADJoin` type, only provide a value for the **onPremisesConnectionId** property and leave the **regionName** property empty. If you choose the `azureADJoin` type, provide a value for either the **onPremisesConnectionId** or the **regionName** property. <br/> _Provider_ allowed values are: `azureADJoin`, `hybridAzureADJoin`, `unknownFutureValue`. The _provider_ default value is `\"azureADJoin\"`.",
					},
					"geographic_location_type": schema.StringAttribute{
						Computed: true,
						Validators: []validator.String{
							stringvalidator.OneOf("default", "asia", "australasia", "canada", "europe", "india", "africa", "usCentral", "usEast", "usWest", "southAmerica", "middleEast", "centralAmerica", "usGovernment", "unknownFutureValue", "mexico"),
						},
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The geographic location where the region is located. Read-only. <br/> Represents the geographic location where a region is located for Microsoft-hosted network for backup Cloud PCs. <br/> This is an evolvable enumeration. <br/> _Provider_ allowed values are: `default`, `asia`, `australasia`, `canada`, `europe`, `india`, `africa`, `usCentral`, `usEast`, `usWest`, `southAmerica`, `middleEast`, `centralAmerica`, `usGovernment`, `unknownFutureValue`, `mexico`.",
					},
					"on_premises_connection_id": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The Azure network connection ID that matches the virtual network IT admins want the provisioning policy to use when they create Cloud PCs. You can use this property in both domain join types: _Azure AD joined_ or _Hybrid Microsoft Entra joined_. If you enter an **onPremisesConnectionId**, leave the **regionName** property empty.",
					},
					"region_group": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("default", "australia", "canada", "usCentral", "usEast", "usWest", "france", "germany", "europeUnion", "unitedKingdom", "japan", "asia", "india", "southAmerica", "euap", "usGovernment", "usGovernmentDOD", "unknownFutureValue", "norway", "switzerland", "southKorea", "middleEast", "mexico", "australasia", "europe", "singapore", "hongKong", "ireland", "sweden", "poland", "italy", "spain", "netherlands", "brazil", "israel", "automatic", "indonesia", "taiwan", "malaysia", "newZealand", "austria", "denmark", "belgium", "kenya"),
						},
						MarkdownDescription: "The logical geographic group this region belongs to. Multiple regions can belong to one region group. A customer can select a **regionGroup** when they provision a Cloud PC, and the Cloud PC is put in one of the regions in the group based on resource status. For example, the Europe region group contains the Northern Europe and Western Europe regions. Read-only. <br/> Represents the logical geographic group that a region belongs to for Microsoft-hosted network for backup Cloud PCs. Multiple regions can belong to one region group. <br/> This is an evolvable enumeration. <br/> _Provider_ allowed values are: `default`, `australia`, `canada`, `usCentral`, `usEast`, `usWest`, `france`, `germany`, `europeUnion`, `unitedKingdom`, `japan`, `asia`, `india`, `southAmerica`, `euap`, `usGovernment`, `usGovernmentDOD`, `unknownFutureValue`, `norway`, `switzerland`, `southKorea`, `middleEast`, `mexico`, `australasia`, `europe`, `singapore`, `hongKong`, `ireland`, `sweden`, `poland`, `italy`, `spain`, `netherlands`, `brazil`, `israel`, `automatic`, `indonesia`, `taiwan`, `malaysia`, `newZealand`, `austria`, `denmark`, `belgium`, `kenya`.",
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
						PlanModifiers: []planmodifier.String{
							wpdefaultvaluemodifier.StringDefaultValue("azureADJoin"),
						},
						Computed:            true,
						MarkdownDescription: "_Provider_ allowed values are: `azureADJoin`, `hybridAzureADJoin`, `unknownFutureValue`. The _provider_ default value is `\"azureADJoin\"`.",
					},
				},
			},
			MarkdownDescription: "Specifies a list ordered by priority on how Cloud PCs join Microsoft Entra ID (Azure AD). Supports `$select`. / Represents a defined configuration of how a provisioned Cloud PC device joins to Microsoft Entra ID. Also see [Microsoft docs for cloudPcDomainJoinConfiguration](https://learn.microsoft.com/en-us/graph/api/resources/cloudpcdomainjoinconfiguration?view=graph-rest-beta). <br> ",
		},
		"enable_single_sign_on": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvaluemodifier.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "`True` if single sign-on can access the provisioned Cloud PC. `False` indicates that the provisioned Cloud PC doesn't support this feature. The default value is `false`. Windows 365 users can use single sign-on to authenticate to Microsoft Entra ID with passwordless options (for example, FIDO keys) to access their Cloud PC. Optional. <br/> The _provider_ default value is `false`.",
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
			MarkdownDescription: "The unique identifier that represents an operating system image that is used for provisioning new Cloud PCs. The format for a gallery type image is: {publisherName_offerName_skuName}. Supported values for each of the parameters are: <br/> - publisher: `Microsoftwindowsdesktop` <br/> - offer: `windows-ent-cpc` <br/> - sku: `21h1-ent-cpc-m365`, `21h1-ent-cpc-os`, `20h2-ent-cpc-m365`, `20h2-ent-cpc-os`, `20h1-ent-cpc-m365`, `20h1-ent-cpc-os`, `19h2-ent-cpc-m365`, and `19h2-ent-cpc-os` <br/> Supports `$filter`, `$select`, and `$orderBy`.",
		},
		"image_type": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("gallery", "custom", "unknownFutureValue"),
			},
			PlanModifiers: []planmodifier.String{
				wpdefaultvaluemodifier.StringDefaultValue("gallery"),
			},
			Computed:            true,
			MarkdownDescription: "The type of operating system image (custom or gallery) that is used for provisioning on Cloud PCs. The default value is `gallery`. Supports $filter, $select, and $orderBy. <br/> _Provider_ allowed values are: `gallery`, `custom`, `unknownFutureValue`. The _provider_ default value is `\"gallery\"`.",
		},
		"last_modified_by": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The unique ID of the user who last updated this policy. For example, `5ccb8d35-dd04-473e-a287-69bb4473208b`. Read-only. Supports `$select`.",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The timestamp when this provisioning policy was last modified. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Read-only. Supports `$select` and `$orderBy`.",
		},
		"local_admin_enabled": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "When `true`, the local admin is enabled for Cloud PCs; `false` indicates that the local admin isn't enabled for Cloud PCs. The default value is `false`. Supports `$filter`, `$select`, and `$orderBy`.",
		},
		"managed_by": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				wpvalidator.FlagEnumValues("windows365", "devBox", "unknownFutureValue", "rpaBox", "microsoft365Opal", "microsoft365BizChat"),
			},
			PlanModifiers: []planmodifier.String{
				wpdefaultvaluemodifier.StringDefaultValue("windows365"),
			},
			Computed:            true,
			MarkdownDescription: "Indicates the service that manages the provisioning policy. The default value is `windows365`. Supports `$filter`, `$select`, and `$orderBy`. <br/> _Provider_ allowed values are: `windows365`, `devBox`, `unknownFutureValue`, `rpaBox`, `microsoft365Opal`, `microsoft365BizChat`. The _provider_ default value is `\"windows365\"`.",
		},
		"microsoft_managed_desktop": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // microsoftManagedDesktop
				"managed_type": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf("notManaged", "premiumManaged", "standardManaged", "starterManaged", "unknownFutureValue"),
					},
					PlanModifiers: []planmodifier.String{
						wpdefaultvaluemodifier.StringDefaultValue("notManaged"),
					},
					Computed:            true,
					MarkdownDescription: "Indicates the provisioning policy associated with Microsoft Managed Desktop settings. The default value is `notManaged`. <br/> _Provider_ allowed values are: `notManaged`, `premiumManaged`, `standardManaged`, `starterManaged`, `unknownFutureValue`. The _provider_ default value is `\"notManaged\"`.",
				},
				"profile": schema.StringAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.String{wpdefaultvaluemodifier.StringDefaultValue("")},
					Computed:            true,
					MarkdownDescription: "The name of the Microsoft Managed Desktop profile that the Windows 365 Cloud PC is associated with. <br/> The _provider_ default value is `\"\"`.",
				},
				"type": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf("notManaged", "premiumManaged", "standardManaged", "starterManaged", "unknownFutureValue"),
					},
					PlanModifiers: []planmodifier.String{
						wpdefaultvaluemodifier.StringDefaultValue("notManaged"),
					},
					Computed:            true,
					MarkdownDescription: "_Provider_ allowed values are: `notManaged`, `premiumManaged`, `standardManaged`, `starterManaged`, `unknownFutureValue`. The _provider_ default value is `\"notManaged\"`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvaluemodifier.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "The specific settings to **microsoftManagedDesktop** that enables Microsoft Managed Desktop customers to get device managed experience for Cloud PC. To enable **microsoftManagedDesktop** to provide more value, an admin needs to specify certain settings in it. Supports `$filter`, `$select`, and `$orderBy`. / Represents specific settings for the Microsoft Managed Desktop that enables customers to get a managed device experience for a Cloud PC. Also see [Microsoft docs for microsoftManagedDesktop](https://learn.microsoft.com/en-us/graph/api/resources/microsoftmanageddesktop?view=graph-rest-beta). <br/> The _provider_ default value is `{}`. <br> ",
		},
		"provisioning_type": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("dedicated", "shared", "unknownFutureValue", "sharedByUser", "sharedByEntraGroup", "reserve"),
			},
			MarkdownDescription: "Specifies the type of licenses to be used when provisioning Cloud PCs using this policy. The possible values are `dedicated`, `shared`, `unknownFutureValue`, `sharedByUser`, `sharedByEntraGroup`, `reserve`. The `shared` member is deprecated and will stop returning on April 30, 2027; going forward, use the `sharedByUser` member. For example, a `dedicated` service plan can be assigned to only one user and provision only one Cloud PC. The `shared` and `sharedByUser` plans require customers to purchase a shared service plan. Each shared license purchased can enable up to three Cloud PCs, with only one user signed in at a time. The `sharedByEntraGroup` plan also requires the purchase of a shared service plan. Each shared license under this plan can enable one Cloud PC, which is shared for the group according to the assignments of this policy. By default, the license type is `dedicated` if the **provisioningType** isn't specified when you create the **cloudPcProvisioningPolicy**. You can't change this property after the **cloudPcProvisioningPolicy** is created. <br/> _Provider_ allowed values are: `dedicated`, `shared`, `unknownFutureValue`, `sharedByUser`, `sharedByEntraGroup`, `reserve`.",
		},
		"scope_ids": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvaluemodifier.SetDefaultValue([]any{"0"})},
			Computed:            true,
			MarkdownDescription: "The _provider_ default value is `[\"0\"]`.",
		},
		"user_experience_type": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("cloudPc", "cloudApp", "unknownFutureValue"),
			},
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			MarkdownDescription: "Specifies the type of cloud object the end user can access. `cloudPc` indicates that the end user can access the entire desktop. `cloudApp` indicates that the end user can only access apps published under this provisioning policy. The type can't be changed once the provisioning policy is created. If not specified during creation, the default value is `cloudPc`. When `cloudApp` is selected, the **provisioningType** must be `sharedByEntraGroup`. Supports `$filter`, `$select`, `$orderBy`. <br/> _Provider_ allowed values are: `cloudPc`, `cloudApp`, `unknownFutureValue`.",
		},
		"user_settings_persistence_configuration": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // cloudPcUserSettingsPersistenceConfiguration
				"user_settings_persistence_enabled": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvaluemodifier.BoolDefaultValue(false)},
					Computed:            true,
					MarkdownDescription: "Indicates whether user application settings are persisted between Cloud PC sessions. The default value is `false`. When `true`, user settings persistence is enabled, and Windows 365 automatically saves any user-specific application data in a central cloud storage location. Anytime the user connects to a Cloud PC within this provisioning policy, Windows 365 reconnects the user to that persisted storage. When `false`, this feature isn't used. The persistent storage can only be accessed by Cloud PC; IT admins can't access it. <br/> The _provider_ default value is `false`.",
				},
				"user_settings_persistence_storage_size_category": schema.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf("fourGB", "eightGB", "sixteenGB", "thirtyTwoGB", "sixtyFourGB", "unknownFutureValue"),
					},
					MarkdownDescription: "Indicates the storage size for persisting user application settings. The default value is `fourGB`. <br/> _Provider_ allowed values are: `fourGB`, `eightGB`, `sixteenGB`, `thirtyTwoGB`, `sixtyFourGB`, `unknownFutureValue`.",
				},
			},
			MarkdownDescription: "Indicates specific settings that enable the persistence of user application settings between Cloud PC sessions. The default value is `null`. This feature is only available for Cloud PC provisioning policies of type `sharedByEntraGroup`. Supports `$select`. / Indicates the user settings persistence configuration when you create Cloud PCs for this [provisioning policy](https://learn.microsoft.com/en-us/graph/api/resources/cloudpcprovisioningpolicy?view=graph-rest-beta). Also see [Microsoft docs for cloudPcUserSettingsPersistenceConfiguration](https://learn.microsoft.com/en-us/graph/api/resources/cloudpcusersettingspersistenceconfiguration?view=graph-rest-beta). <br> ",
		},
		"windows_setting": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // cloudPcWindowsSetting
				"locale": schema.StringAttribute{
					Optional:            true,
					MarkdownDescription: "The Windows language or region tag to use for language pack configuration and localization of the Cloud PC. The default value is `en-US`, which corresponds to English (United States).",
				},
			},
			MarkdownDescription: "Indicates a specific Windows setting to configure during the creation of Cloud PCs for this provisioning policy. Supports `$select`. / Represents a specific Windows setting to configure during the creation of Cloud PCs for a provisioning policy. Also see [Microsoft docs for cloudPcWindowsSetting](https://learn.microsoft.com/en-us/graph/api/resources/cloudpcwindowssetting?view=graph-rest-beta). <br> ",
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
									MarkdownDescription: "Complex type that represents the assignment target group. Also see [Microsoft docs for cloudPcManagementGroupAssignmentTarget](https://learn.microsoft.com/en-us/graph/api/resources/cloudpcmanagementgroupassignmenttarget?view=graph-rest-beta). <br> ",
								},
							},
						},
						MarkdownDescription: "The assignment target for the provisioning policy. Currently, the only target supported for this policy is a user group. For details, see [cloudPcManagementGroupAssignmentTarget](cloudpcmanagementgroupassignmenttarget.md). / Represents an abstract base type for assignment targets. Also see [Microsoft docs for cloudPcManagementAssignmentTarget](https://learn.microsoft.com/en-us/graph/api/resources/cloudpcmanagementassignmenttarget?view=graph-rest-beta). <br> ",
					},
					"user_settings_persistence_detail": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{ // cloudPCUserSettingsPersistenceDetail
							"id": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Indicates the unique identifier for the Cloud PC user settings persistence configuration for a single policy collection. Required. Read-only.",
							},
							"grace_period_end_date_time": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Indicates the grace period end time when user settings persistence exceeds the available quota. If usage exceeds the available quota when the grace period expires, the system automatically deletes the profile with the oldest last attached timestamp. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Read-only.",
							},
						},
						PlanModifiers:       []planmodifier.Object{wpplanmodifier.ObjectUseStateForUnknown()},
						MarkdownDescription: "The assignment targeted user settings persistence for the provisioning policy. It allows user application data and Windows settings to be saved and applied between sessions. / Represents the configuration that indicates whether Cloud PC user settings persistence is enabled. When enabled, Windows 365 saves user-specific application data in a central cloud storage location and reconnects the user to that storage upon each connection. Also see [Microsoft docs for cloudPCUserSettingsPersistenceDetail](https://learn.microsoft.com/en-us/graph/api/resources/cloudpcusersettingspersistencedetail?view=graph-rest-beta). <br> ",
					},
				},
			},
			PlanModifiers:       []planmodifier.Set{wpdefaultvaluemodifier.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "A defined collection of provisioning policy assignments. Represents the set of Microsoft 365 groups and security groups in Microsoft Entra ID that have provisioning policy assigned. Returned only on `$expand`. For an example about how to get the assignments relationship, see [Get cloudPcProvisioningPolicy](https://learn.microsoft.com/en-us/graph/api/cloudpcprovisioningpolicy-get?view=graph-rest-beta). / Represents a defined collection of provisioning policy assignments. Also see [Microsoft docs for cloudPcProvisioningPolicyAssignment](https://learn.microsoft.com/en-us/graph/api/resources/cloudpcprovisioningpolicyassignment?view=graph-rest-beta). <br/> The _provider_ default value is `[]`. <br> ",
		},
	},
	MarkdownDescription: "Represents a Cloud PC provisioning policy. <br/> Also see [Microsoft docs for cloudPcProvisioningPolicy](https://learn.microsoft.com/en-us/graph/api/resources/cloudpcprovisioningpolicy?view=graph-rest-beta). ||| MS Graph: Cloud PC",
}

var cloudPcProvisioningPolicyCloudPcManagementAssignmentTargetValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("group"),
)
