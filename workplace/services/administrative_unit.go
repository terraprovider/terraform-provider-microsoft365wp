package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var (
	AdministrativeUnitResource = generic.GenericResource{
		TypeNameSuffix: "administrative_unit",
		SpecificSchema: administrativeUnitResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/administrativeUnits",
			ReadOptions: generic.ReadOptions{
				DataSource: generic.DataSourceOptions{
					ExtraFilterAttributes: []string{"is_member_management_restricted", "membership_type"},
					Plural: generic.PluralOptions{
						ExtraAttributes: []string{"is_member_management_restricted", "membership_type"},
					},
				},
			},
		},
	}

	AdministrativeUnitSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&AdministrativeUnitResource)

	AdministrativeUnitPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&AdministrativeUnitResource, "")
)

var administrativeUnitResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // administrativeUnit
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Unique identifier for the administrative unit. Read-only. Supports `$filter` (`eq`).",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "An optional description for the administrative unit. Supports `$filter` (`eq`, `ne`, `in`, `startsWith`), `$search`.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Display name for the administrative unit. Maximum length is 256 characters. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values), `$search`, and `$orderby`.",
		},
		"is_member_management_restricted": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			MarkdownDescription: "`true` if members of this administrative unit should be treated as sensitive, which requires specific permissions to manage. If not set, the default value is `null` and the default behavior is false. Use this property to define administrative units with roles that don't inherit from tenant-level administrators, and where the management of individual member objects is limited to administrators scoped to a restricted management administrative unit. This property is immutable and can't be changed later. <br/><br/> For more information on how to work with restricted management administrative units, see [Restricted management administrative units in Microsoft Entra ID](/entra/identity/role-based-access-control/admin-units-restricted-management).",
		},
		"membership_rule": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The dynamic membership rule for the administrative unit. For more information about the rules you can use for dynamic administrative units and dynamic groups, see [Manage rules for dynamic membership groups in Microsoft Entra ID](/entra/identity/users/groups-dynamic-membership).",
		},
		"membership_rule_processing_state": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Controls whether the dynamic membership rule is actively processed. Set to `On` to activate the dynamic membership rule, or `Paused` to stop updating membership dynamically.",
		},
		"membership_type": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Indicates the membership type for the administrative unit. The If not set, the default value is `null` and the default behavior is assigned.",
		},
		"visibility": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Controls whether the administrative unit and its members are hidden or public. Can be set to `HiddenMembership` or `Public`. If not set, the default value is `null` and the default behavior is public. When set to `HiddenMembership`, only members of the administrative unit can list other members of the administrative unit.",
		},
	},
	MarkdownDescription: "An administrative unit provides a conceptual container for user, group, and device directory objects. With administrative units, a company administrator can now delegate administrative responsibilities to manage the users, groups, and devices contained within or scoped to an administrative unit to a regional or departmental administrator. For more information about administrative units, see [Administrative units in Microsoft Entra ID](/entra/identity/role-based-access-control/administrative-units).\n\nThis resource is an open type that allows other properties to be passed in.\n\nThis resource supports: / https://learn.microsoft.com/en-us/graph/api/resources/administrativeunit?view=graph-rest-beta\n\nProvider Note: This resource is only provided as an enhanced version of `azuread_administrative_unit` supporting the attribute `is_member_management_restricted`. Use `azuread_administrative_unit_member` and `azuread_administrative_unit_role_member` to modify membership as it is not planned to add further functionality here. ||| MS Graph: Directory management",
}
