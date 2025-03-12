package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	AuthorizationPolicyResource = generic.GenericResource{
		TypeNameSuffix: "authorization_policy",
		SpecificSchema: authorizationPolicyResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:     "/policies/authorizationPolicy/authorizationPolicy",
			IsSingleton: true,
			SingularEntity: generic.SingularEntity{
				ReadIdAttrFromGraphUseSelectId: true,
			},
		},
	}

	AuthorizationPolicySingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&AuthorizationPolicyResource)
)

var authorizationPolicyResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // authorizationPolicy
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "ID of the authorization policy. Required. Read-only.",
		},
		"description": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Description of this policy.",
		},
		"display_name": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Display name for this policy.",
		},
		"allowed_to_sign_up_email_based_subscriptions": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:            true,
			MarkdownDescription: "Indicates whether users can sign up for email based subscriptions. The _provider_ default value is `true`.",
		},
		"allowed_to_use_sspr": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
			Computed:            true,
			Description:         `allowedToUseSSPR`, // custom MS Graph attribute name
			MarkdownDescription: "Indicates whether administrators of the tenant can use the Self-Service Password Reset (SSPR). For more information, see [Self-service password reset for administrators](/entra/identity/authentication/concept-sspr-policy#administrator-reset-policy-differences). The _provider_ default value is `true`.",
		},
		"allow_email_verified_users_to_join_organization": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Indicates whether a user can join the tenant by email validation. The _provider_ default value is `false`.",
		},
		"allow_invites_from": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("none", "adminsAndGuestInviters", "adminsGuestInvitersAndAllMembers", "everyone", "unknownFutureValue"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("everyone")},
			Computed:            true,
			MarkdownDescription: "Indicates who can invite guests to the organization. `everyone` is the default setting for all cloud environments except US Government. For more information, see [allowInvitesFrom values](#allowinvitesfrom-values). / Possible values are: `none`, `adminsAndGuestInviters`, `adminsGuestInvitersAndAllMembers`, `everyone`, `unknownFutureValue`. The _provider_ default value is `\"everyone\"`.",
		},
		"allow_user_consent_for_risky_apps": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Indicates whether [user consent for risky apps](/azure/active-directory/manage-apps/configure-risk-based-step-up-consent) is allowed. Default value is `false`. We recommend that you keep the value set to `false`. The _provider_ default value is `false`.",
		},
		"block_msol_powershell": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			Description:         `blockMsolPowerShell`, // custom MS Graph attribute name
			MarkdownDescription: "To disable the use of the [MSOnline PowerShell module](/powershell/module/msonline) set this property to `true`. This also disables user-based access to the legacy service endpoint used by the MSOnline PowerShell module. This doesn't affect Microsoft Entra Connect or Microsoft Graph. The _provider_ default value is `false`.",
		},
		"default_user_role_permissions": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // defaultUserRolePermissions
				"allowed_to_create_apps": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
					Computed:            true,
					MarkdownDescription: "Indicates whether the default user role can create applications. This setting corresponds to the _Users can register applications_ setting in the [User settings menu in the Microsoft Entra admin center](/azure/active-directory/fundamentals/users-default-permissions?context=graph%2Fcontext#restrict-member-users-default-permissions). The _provider_ default value is `true`.",
				},
				"allowed_to_create_security_groups": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
					Computed:            true,
					MarkdownDescription: "Indicates whether the default user role can create security groups. This setting corresponds to the following menus in the Microsoft Entra admin center: <br/><li> _The Users can create security groups in Microsoft Entra admin centers, API or PowerShell_ setting in the [Group settings menu](/azure/active-directory/enterprise-users/groups-self-service-management). <li> _Users can create security groups_ setting in the [User settings menu](/azure/active-directory/fundamentals/users-default-permissions?context=graph%2Fcontext#restrict-member-users-default-permissions). The _provider_ default value is `true`.",
				},
				"allowed_to_create_tenants": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
					Computed:            true,
					MarkdownDescription: "Indicates whether the default user role can create tenants. This setting corresponds to the _Restrict non-admin users from creating tenants_ setting in the [User settings menu in the Microsoft Entra admin center](/azure/active-directory/fundamentals/users-default-permissions?context=graph%2Fcontext#restrict-member-users-default-permissions).<br/><br/> When this setting is `false`, users assigned the [Tenant Creator](/entra/identity/role-based-access-control/permissions-reference?context=graph%2Fcontext#tenant-creator) role can still create tenants. The _provider_ default value is `true`.",
				},
				"allowed_to_read_bitlocker_keys_for_owned_device": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
					Computed:            true,
					MarkdownDescription: "Indicates whether the registered owners of a device can read their own BitLocker recovery keys with default user role. The _provider_ default value is `true`.",
				},
				"allowed_to_read_other_users": schema.BoolAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(true)},
					Computed:            true,
					MarkdownDescription: "Indicates whether the default user role can read other users. **DO NOT SET THIS VALUE TO `false`**. The _provider_ default value is `true`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Specifies certain customizable permissions for default user role. / Contains certains customizable permissions of default user role in Microsoft Entra ID. / https://learn.microsoft.com/en-us/graph/api/resources/defaultuserrolepermissions?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
		"enabled_preview_features": schema.SetAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "List of features enabled for private preview on the tenant. The _provider_ default value is `[]`.",
		},
		"guest_user_role_id": schema.StringAttribute{
			Optional: true,
			PlanModifiers: []planmodifier.String{
				wpdefaultvalue.StringDefaultValue("10dae51f-b6af-4016-8d66-8c2a99b929b3"),
			},
			Computed:            true,
			MarkdownDescription: "Represents role templateId for the role that should be granted to guests. Refer to [List unifiedRoleDefinitions](../api/rbacapplication-list-roledefinitions.md) to find the list of available role templates. Currently following roles are supported:  *User* (`a0b1b346-4d3e-4e8b-98f8-753987be4970`), *Guest User* (`10dae51f-b6af-4016-8d66-8c2a99b929b3`), and *Restricted Guest User* (`2af84b1e-32c8-42b7-82bc-daa82404023b`). The _provider_ default value is `\"10dae51f-b6af-4016-8d66-8c2a99b929b3\"`.  \nProvider Note: The default value `10dae51f-b6af-4016-8d66-8c2a99b929b3` corresponds to \"Guest User\".",
		},
		"permission_grant_policy_ids_assigned_to_default_user_role": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			PlanModifiers: []planmodifier.Set{
				wpdefaultvalue.SetDefaultValue([]any{
					"ManagePermissionGrantsForSelf.microsoft-user-default-legacy",
					"ManagePermissionGrantsForOwnedResource.microsoft-dynamically-managed-permissions-for-team",
					"ManagePermissionGrantsForOwnedResource.microsoft-dynamically-managed-permissions-for-chat",
				}),
			},
			Computed:            true,
			MarkdownDescription: "Indicates if user consent to apps is allowed, and if it is, the [app consent policy](../resources/permissiongrantpolicy.md) that governs the permission for users to grant consent. Values should be in the format `managePermissionGrantsForSelf.{id}` for user consent policies or `managePermissionGrantsForOwnedResource.{id}` for resource-specific consent policies, where `{id}` is the **id** of a built-in or custom app consent policy. An empty list indicates user consent to apps is disabled. The _provider_ default value is `[\"ManagePermissionGrantsForSelf.microsoft-user-default-legacy\",\"ManagePermissionGrantsForOwnedResource.microsoft-dynamically-managed-permissions-for-team\",\"ManagePermissionGrantsForOwnedResource.microsoft-dynamically-managed-permissions-for-chat\"]`.",
		},
	},
	MarkdownDescription: "Represents a policy that can control Microsoft Entra authorization settings. It's a singleton that inherits from base policy type, and always exists for the tenant. / https://learn.microsoft.com/en-us/graph/api/resources/authorizationpolicy?view=graph-rest-beta ||| MS Graph: Policies",
}
