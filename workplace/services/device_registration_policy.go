package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvaluemodifier"
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
	DeviceRegistrationPolicyResource = generic.GenericResource{
		TypeNameSuffix: "device_registration_policy",
		SpecificSchema: deviceRegistrationPolicyResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:     "/policies/deviceRegistrationPolicy",
			IsSingleton: true,
			WriteOptions: generic.WriteOptions{
				UsePutForUpdate: true,
			},
		},
	}

	DeviceRegistrationPolicySingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&DeviceRegistrationPolicyResource)
)

var deviceRegistrationPolicyResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // deviceRegistrationPolicy
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The identifier of the device registration policy. It's always set to `deviceRegistrationPolicy`. Read-only.",
		},
		"azure_ad_join": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // azureADJoinPolicy
				"allowed_to_join": schema.SingleNestedAttribute{
					Optional:            true,
					Attributes:          deviceRegistrationPolicyDeviceRegistrationMembershipAttributes,
					PlanModifiers:       []planmodifier.Object{deviceRegistrationPolicyAll},
					Computed:            true,
					MarkdownDescription: "Determines if Microsoft Entra join is allowed. / An abstract resource type indicating the scope that a device registration policy applies to. The derived types are [noDeviceRegistrationMembership](https://learn.microsoft.com/en-us/graph/api/resources/nodeviceregistrationmembership?view=graph-rest-beta), [allDeviceRegistrationMembership](https://learn.microsoft.com/en-us/graph/api/resources/alldeviceregistrationmembership?view=graph-rest-beta) and [enumeratedDeviceRegistrationMembership](https://learn.microsoft.com/en-us/graph/api/resources/enumerateddeviceregistrationmembership?view=graph-rest-beta). Also see [Microsoft docs for deviceRegistrationMembership](https://learn.microsoft.com/en-us/graph/api/resources/deviceregistrationmembership?view=graph-rest-beta). <br/> The _provider_ default value is `deviceRegistrationPolicyAll`. <br> ",
				},
				"is_admin_configurable": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvaluemodifier.BoolDefaultValue(true)},
					Computed:            true,
					MarkdownDescription: "Determines if administrators can modify this policy. <br/> The _provider_ default value is `true`.",
				},
				"local_admins": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // localAdminSettings
						"enable_global_admins": schema.BoolAttribute{
							Optional:            true,
							PlanModifiers:       []planmodifier.Bool{wpdefaultvaluemodifier.BoolDefaultValue(true)},
							Computed:            true,
							MarkdownDescription: "Indicates whether global administrators are local administrators on all Microsoft Entra-joined devices. This setting only applies to future registrations. Default is `true`. <br/> The _provider_ default value is `true`.",
						},
						"registering_users": schema.SingleNestedAttribute{
							Optional:            true,
							Attributes:          deviceRegistrationPolicyDeviceRegistrationMembershipAttributes,
							PlanModifiers:       []planmodifier.Object{deviceRegistrationPolicyAll},
							Computed:            true,
							MarkdownDescription: "Determines the users and groups that become local administrators on Microsoft Entra joined devices that they register. / An abstract resource type indicating the scope that a device registration policy applies to. The derived types are [noDeviceRegistrationMembership](https://learn.microsoft.com/en-us/graph/api/resources/nodeviceregistrationmembership?view=graph-rest-beta), [allDeviceRegistrationMembership](https://learn.microsoft.com/en-us/graph/api/resources/alldeviceregistrationmembership?view=graph-rest-beta) and [enumeratedDeviceRegistrationMembership](https://learn.microsoft.com/en-us/graph/api/resources/enumerateddeviceregistrationmembership?view=graph-rest-beta). Also see [Microsoft docs for deviceRegistrationMembership](https://learn.microsoft.com/en-us/graph/api/resources/deviceregistrationmembership?view=graph-rest-beta). <br/> The _provider_ default value is `deviceRegistrationPolicyAll`. <br> ",
						},
					},
					PlanModifiers:       []planmodifier.Object{wpdefaultvaluemodifier.ObjectDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "Determines who becomes a local administrator on joined devices. / Controls local administrators on Microsoft Entra-joined devices. Also see [Microsoft docs for localAdminSettings](https://learn.microsoft.com/en-us/graph/api/resources/localadminsettings?view=graph-rest-beta). <br/> The _provider_ default value is `{}`. <br> ",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvaluemodifier.ObjectDefaultValueEmpty()},
			Computed:            true,
			Description:         `azureADJoin`, // custom MS Graph attribute name
			MarkdownDescription: "Specifies the authorization policy for controlling registration of new devices using **Microsoft Entra join** within your organization. Required. For more information, see [What is a device identity?](https://learn.microsoft.com/en-us/azure/active-directory/devices/overview). / Represents the policy scope of the Microsoft Entra tenant that controls the ability for users and groups to register device identities to your organization using Microsoft Entra join. Also see [Microsoft docs for azureADJoinPolicy](https://learn.microsoft.com/en-us/graph/api/resources/azureadjoinpolicy?view=graph-rest-beta). <br/> The _provider_ default value is `{}`. <br> ",
		},
		"azure_ad_registration": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // azureADRegistrationPolicy
				"allowed_to_register": schema.SingleNestedAttribute{
					Optional:            true,
					Attributes:          deviceRegistrationPolicyDeviceRegistrationMembershipAttributes,
					PlanModifiers:       []planmodifier.Object{deviceRegistrationPolicyAll},
					Computed:            true,
					MarkdownDescription: "Determines if Microsoft Entra registered is allowed. / An abstract resource type indicating the scope that a device registration policy applies to. The derived types are [noDeviceRegistrationMembership](https://learn.microsoft.com/en-us/graph/api/resources/nodeviceregistrationmembership?view=graph-rest-beta), [allDeviceRegistrationMembership](https://learn.microsoft.com/en-us/graph/api/resources/alldeviceregistrationmembership?view=graph-rest-beta) and [enumeratedDeviceRegistrationMembership](https://learn.microsoft.com/en-us/graph/api/resources/enumerateddeviceregistrationmembership?view=graph-rest-beta). Also see [Microsoft docs for deviceRegistrationMembership](https://learn.microsoft.com/en-us/graph/api/resources/deviceregistrationmembership?view=graph-rest-beta). <br/> The _provider_ default value is `deviceRegistrationPolicyAll`. <br> ",
				},
				"is_admin_configurable": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvaluemodifier.BoolDefaultValue(false)},
					Computed:            true,
					MarkdownDescription: "Determines if administrators can modify this policy. <br/> The _provider_ default value is `false`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvaluemodifier.ObjectDefaultValueEmpty()},
			Computed:            true,
			Description:         `azureADRegistration`, // custom MS Graph attribute name
			MarkdownDescription: "Specifies the authorization policy for controlling registration of new devices using **Microsoft Entra registered** within your organization. Required. For more information, see [What is a device identity?](https://learn.microsoft.com/en-us/azure/active-directory/devices/overview). / Represents the policy scope of the Microsoft Entra tenant that controls the ability for users and groups to register device identities to your organization using **Microsoft Entra registered**. For more information, see [What is a device identity?](https://learn.microsoft.com/en-us/azure/active-directory/devices/overview). Also see [Microsoft docs for azureADRegistrationPolicy](https://learn.microsoft.com/en-us/graph/api/resources/azureadregistrationpolicy?view=graph-rest-beta). <br/> The _provider_ default value is `{}`. <br> ",
		},
		"description": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The description of the device registration policy. It's always set to `Tenant-wide policy that manages intial provisioning controls using quota restrictions, additional authentication and authorization checks`. Read-only.",
		},
		"display_name": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The name of the device registration policy. It's always set to `Device Registration Policy`. Read-only.",
		},
		"local_admin_password": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // localAdminPasswordSettings
				"is_enabled": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvaluemodifier.BoolDefaultValue(false)},
					Computed:            true,
					MarkdownDescription: "Specifies whether this policy scope is configurable by the admin. The default value is `false`. An admin can set it to true to enable Local Admin Password Solution (LAPS) within their organzation. <br/> The _provider_ default value is `false`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvaluemodifier.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Specifies the setting for **Local Admin Password Solution (LAPS)** within your organization. / Represents the policy scope of the Microsoft Entra tenant that controls the Local Admin Password Solution (LAPS) setting. Also see [Microsoft docs for localAdminPasswordSettings](https://learn.microsoft.com/en-us/graph/api/resources/localadminpasswordsettings?view=graph-rest-beta). <br/> The _provider_ default value is `{}`. <br> ",
		},
		"multi_factor_auth_configuration": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("notRequired", "required", "unknownFutureValue"),
			},
			PlanModifiers: []planmodifier.String{
				wpdefaultvaluemodifier.StringDefaultValue("notRequired"),
			},
			Computed:            true,
			MarkdownDescription: "Specifies the authentication policy for a user to complete registration using **Microsoft Entra join** or **Microsoft Entra registered** within your organization. The default value is `notRequired`. <br/> _Provider_ allowed values are: `notRequired`, `required`, `unknownFutureValue`. The _provider_ default value is `\"notRequired\"`.",
		},
		"user_device_quota": schema.Int64Attribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Int64{wpdefaultvaluemodifier.Int64DefaultValue(50)},
			Computed:            true,
			MarkdownDescription: "Specifies the maximum number of devices that a user can have within your organization before blocking new device registrations. The default value is set to 50. If this property isn't specified during the policy update operation, it's automatically reset to `0` to indicate that users aren't allowed to join any devices. <br/> The _provider_ default value is `50`.",
		},
	},
	MarkdownDescription: "Represents the policy scope that controls quota restrictions, additional authentication, and authorization policies to register device identities to your organization. <br/> Also see [Microsoft docs for deviceRegistrationPolicy](https://learn.microsoft.com/en-us/graph/api/resources/deviceregistrationpolicy?view=graph-rest-beta). ||| MS Graph: Policies",
}

var deviceRegistrationPolicyDeviceRegistrationMembershipAttributes = map[string]schema.Attribute{ // deviceRegistrationMembership
	"all": generic.OdataDerivedTypeNestedAttributeRs{
		DerivedType: "#microsoft.graph.allDeviceRegistrationMembership",
		SingleNestedAttribute: schema.SingleNestedAttribute{
			Optional:   true,
			Attributes: map[string]schema.Attribute{ // allDeviceRegistrationMembership
			},
			Validators: []validator.Object{
				deviceRegistrationPolicyDeviceRegistrationMembershipValidator,
			},
			MarkdownDescription: "Indicates that this device registration policy applies to all users and groups. Also see [Microsoft docs for allDeviceRegistrationMembership](https://learn.microsoft.com/en-us/graph/api/resources/alldeviceregistrationmembership?view=graph-rest-beta). <br> ",
		},
	},
	"enumerated": generic.OdataDerivedTypeNestedAttributeRs{
		DerivedType: "#microsoft.graph.enumeratedDeviceRegistrationMembership",
		SingleNestedAttribute: schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // enumeratedDeviceRegistrationMembership
				"groups": schema.SetAttribute{
					ElementType:         types.StringType,
					Optional:            true,
					PlanModifiers:       []planmodifier.Set{wpdefaultvaluemodifier.SetDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "List of groups that this policy applies to. <br/> The _provider_ default value is `[]`.",
				},
				"users": schema.SetAttribute{
					ElementType:         types.StringType,
					Optional:            true,
					PlanModifiers:       []planmodifier.Set{wpdefaultvaluemodifier.SetDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "List of users that this policy applies to. <br/> The _provider_ default value is `[]`.",
				},
			},
			Validators: []validator.Object{
				deviceRegistrationPolicyDeviceRegistrationMembershipValidator,
			},
			MarkdownDescription: "Indicates that this device registration policy applies to the enumerated users and groups. Also see [Microsoft docs for enumeratedDeviceRegistrationMembership](https://learn.microsoft.com/en-us/graph/api/resources/enumerateddeviceregistrationmembership?view=graph-rest-beta). <br> ",
		},
	},
	"no": generic.OdataDerivedTypeNestedAttributeRs{
		DerivedType: "#microsoft.graph.noDeviceRegistrationMembership",
		SingleNestedAttribute: schema.SingleNestedAttribute{
			Optional:   true,
			Attributes: map[string]schema.Attribute{ // noDeviceRegistrationMembership
			},
			Validators: []validator.Object{
				deviceRegistrationPolicyDeviceRegistrationMembershipValidator,
			},
			MarkdownDescription: "Indicates that no users are allowed to join or register devices. Also see [Microsoft docs for noDeviceRegistrationMembership](https://learn.microsoft.com/en-us/graph/api/resources/nodeviceregistrationmembership?view=graph-rest-beta). <br> ",
		},
	},
}

var deviceRegistrationPolicyDeviceRegistrationMembershipValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("all"),
	path.MatchRelative().AtParent().AtName("enumerated"),
	path.MatchRelative().AtParent().AtName("no"),
)

var deviceRegistrationPolicyAll = wpdefaultvaluemodifier.ObjectDefaultValue(map[string]any{"all": map[string]any{}})
