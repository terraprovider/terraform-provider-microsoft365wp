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
					MarkdownDescription: "Determines if Microsoft Entra join is allowed. / An abstract resource type indicating the scope that a device registration policy applies to. The derived types are [noDeviceRegistrationMembership](../resources/nodeviceregistrationmembership.md), [allDeviceRegistrationMembership](../resources/alldeviceregistrationmembership.md) and [enumeratedDeviceRegistrationMembership](../resources/enumerateddeviceregistrationmembership.md). / https://learn.microsoft.com/en-us/graph/api/resources/deviceregistrationmembership?view=graph-rest-beta. The _provider_ default value is `deviceRegistrationPolicyAll`.",
				},
				"is_admin_configurable": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
					Computed:            true,
					MarkdownDescription: "Determines if administrators can modify this policy. The _provider_ default value is `true`.",
				},
				"local_admins": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{ // localAdminSettings
						"enable_global_admins": schema.BoolAttribute{
							Optional:            true,
							PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
							Computed:            true,
							MarkdownDescription: "Indicates whether global administrators are local administrators on all Microsoft Entra-joined devices. This setting only applies to future registrations. Default is `true`. The _provider_ default value is `true`.",
						},
						"registering_users": schema.SingleNestedAttribute{
							Optional:            true,
							Attributes:          deviceRegistrationPolicyDeviceRegistrationMembershipAttributes,
							PlanModifiers:       []planmodifier.Object{deviceRegistrationPolicyAll},
							Computed:            true,
							MarkdownDescription: "Determines the users and groups that become local administrators on Microsoft Entra joined devices that they register. / An abstract resource type indicating the scope that a device registration policy applies to. The derived types are [noDeviceRegistrationMembership](../resources/nodeviceregistrationmembership.md), [allDeviceRegistrationMembership](../resources/alldeviceregistrationmembership.md) and [enumeratedDeviceRegistrationMembership](../resources/enumerateddeviceregistrationmembership.md). / https://learn.microsoft.com/en-us/graph/api/resources/deviceregistrationmembership?view=graph-rest-beta. The _provider_ default value is `deviceRegistrationPolicyAll`.",
						},
					},
					PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "Determines who becomes a local administrator on joined devices. / Controls local administrators on Microsoft Entra-joined devices. / https://learn.microsoft.com/en-us/graph/api/resources/localadminsettings?view=graph-rest-beta. The _provider_ default value is `{}`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			Description:         `azureADJoin`, // custom MS Graph attribute name
			MarkdownDescription: "Specifies the authorization policy for controlling registration of new devices using **Microsoft Entra join** within your organization. Required. For more information, see [What is a device identity?](/azure/active-directory/devices/overview). / Represents the policy scope of the Microsoft Entra tenant that controls the ability for users and groups to register device identities to your organization using Microsoft Entra join. / https://learn.microsoft.com/en-us/graph/api/resources/azureadjoinpolicy?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
		"azure_ad_registration": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // azureADRegistrationPolicy
				"allowed_to_register": schema.SingleNestedAttribute{
					Optional:            true,
					Attributes:          deviceRegistrationPolicyDeviceRegistrationMembershipAttributes,
					PlanModifiers:       []planmodifier.Object{deviceRegistrationPolicyAll},
					Computed:            true,
					MarkdownDescription: "Determines if Microsoft Entra registered is allowed. / An abstract resource type indicating the scope that a device registration policy applies to. The derived types are [noDeviceRegistrationMembership](../resources/nodeviceregistrationmembership.md), [allDeviceRegistrationMembership](../resources/alldeviceregistrationmembership.md) and [enumeratedDeviceRegistrationMembership](../resources/enumerateddeviceregistrationmembership.md). / https://learn.microsoft.com/en-us/graph/api/resources/deviceregistrationmembership?view=graph-rest-beta. The _provider_ default value is `deviceRegistrationPolicyAll`.",
				},
				"is_admin_configurable": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
					Computed:            true,
					MarkdownDescription: "Determines if administrators can modify this policy. The _provider_ default value is `false`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			Description:         `azureADRegistration`, // custom MS Graph attribute name
			MarkdownDescription: "Specifies the authorization policy for controlling registration of new devices using **Microsoft Entra registered** within your organization. Required. For more information, see [What is a device identity?](/azure/active-directory/devices/overview). / Represents the policy scope of the Microsoft Entra tenant that controls the ability for users and groups to register device identities to your organization using **Microsoft Entra registered**. For more information, see [What is a device identity?](/azure/active-directory/devices/overview). / https://learn.microsoft.com/en-us/graph/api/resources/azureadregistrationpolicy?view=graph-rest-beta. The _provider_ default value is `{}`.",
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
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
					Computed:            true,
					MarkdownDescription: "Specifies whether this policy scope is configurable by the admin. The default value is `false`. An admin can set it to true to enable Local Admin Password Solution (LAPS) within their organzation. The _provider_ default value is `false`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Specifies the setting for **Local Admin Password Solution (LAPS)** within your organization. / Represents the policy scope of the Microsoft Entra tenant that controls the Local Admin Password Solution (LAPS) setting. / https://learn.microsoft.com/en-us/graph/api/resources/localadminpasswordsettings?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
		"multi_factor_auth_configuration": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("notRequired", "required", "unknownFutureValue"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("notRequired")},
			Computed:            true,
			MarkdownDescription: "Specifies the authentication policy for a user to complete registration using **Microsoft Entra join** or **Microsoft Entra registered** within your organization. The The default value is `notRequired`. / Possible values are: `notRequired`, `required`, `unknownFutureValue`. The _provider_ default value is `\"notRequired\"`.",
		},
		"user_device_quota": schema.Int64Attribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(50)},
			Computed:            true,
			MarkdownDescription: "Specifies the maximum number of devices that a user can have within your organization before blocking new device registrations. The default value is set to 50. If this property isn't specified during the policy update operation, it's automatically reset to `0` to indicate that users aren't allowed to join any devices. The _provider_ default value is `50`.",
		},
	},
	MarkdownDescription: "Represents the policy scope that controls quota restrictions, additional authentication, and authorization policies to register device identities to your organization. / https://learn.microsoft.com/en-us/graph/api/resources/deviceregistrationpolicy?view=graph-rest-beta ||| MS Graph: Policies",
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
			MarkdownDescription: "Indicates that this device registration policy applies to all users and groups. / https://learn.microsoft.com/en-us/graph/api/resources/alldeviceregistrationmembership?view=graph-rest-beta",
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
					PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "List of groups that this policy applies to. The _provider_ default value is `[]`.",
				},
				"users": schema.SetAttribute{
					ElementType:         types.StringType,
					Optional:            true,
					PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
					Computed:            true,
					MarkdownDescription: "List of users that this policy applies to. The _provider_ default value is `[]`.",
				},
			},
			Validators: []validator.Object{
				deviceRegistrationPolicyDeviceRegistrationMembershipValidator,
			},
			MarkdownDescription: "Indicates that this device registration policy applies to the enumerated users and groups. / https://learn.microsoft.com/en-us/graph/api/resources/enumerateddeviceregistrationmembership?view=graph-rest-beta",
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
			MarkdownDescription: "Indicates that no users are allowed to join or register devices. / https://learn.microsoft.com/en-us/graph/api/resources/nodeviceregistrationmembership?view=graph-rest-beta",
		},
	},
}

var deviceRegistrationPolicyDeviceRegistrationMembershipValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("all"),
	path.MatchRelative().AtParent().AtName("enumerated"),
	path.MatchRelative().AtParent().AtName("no"),
)

var deviceRegistrationPolicyAll = wpdefaultvalue.ObjectDefaultValue(map[string]any{"all": map[string]any{}})
