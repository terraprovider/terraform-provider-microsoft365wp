package services

import (
	"context"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpvalidator"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	UnifiedRoleDefinitionResource = generic.GenericResource{
		TypeNameSuffix: "unified_role_definition",
		SpecificSchema: unifiedRoleDefinitionResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/roleManagement/directory/roleDefinitions",
			ReadOptions: generic.ReadOptions{
				DataSource: generic.DataSourceOptions{
					ExtraFilterAttributes: []string{"display_name"},
					Plural: generic.PluralOptions{
						ExtraAttributes: []string{"allowed_principal_types", "is_built_in", "is_enabled", "is_privileged"},
					},
				},
			},
			TerraformToGraphMiddleware: unifiedRoleDefinitionTerraformToGraphMiddleware,
		},
	}

	UnifiedRoleDefinitionSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&UnifiedRoleDefinitionResource)

	UnifiedRoleDefinitionPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&UnifiedRoleDefinitionResource, "")
)

func unifiedRoleDefinitionTerraformToGraphMiddleware(ctx context.Context, diags *diag.Diagnostics, params *generic.TerraformToGraphMiddlewareParams) generic.TerraformToGraphMiddlewareReturns {
	if params.IsUpdate {
		// do not send templateId as null on PATCH if not set
		if params.RawVal["templateId"] == nil {
			delete(params.RawVal, "templateId")
		}
	}
	return nil
}

var unifiedRoleDefinitionResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // unifiedRoleDefinition
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The unique identifier for the unifiedRoleDefinition. Key, not nullable, Read-only.  Supports `$filter` (`eq` operator only).",
		},
		"allowed_principal_types": schema.StringAttribute{
			Computed: true,
			Validators: []validator.String{
				wpvalidator.FlagEnumValues("user", "servicePrincipal", "group", "unknownFutureValue"),
			},
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Types of principals that can be assigned the role. Read-only. The This is a multi-valued enumeration that can contain up to three values as a comma-separated string. For example, `user, group`. Supports `$filter` (`eq`). / Possible values are: `user`, `servicePrincipal`, `group`, `unknownFutureValue`",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The description for the unifiedRoleDefinition. Read-only when **isBuiltIn** is `true`.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The display name for the unifiedRoleDefinition. Read-only when **isBuiltIn** is `true`. Required.  Supports `$filter` (`eq` and `startsWith`).",
		},
		"is_built_in": schema.BoolAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
			MarkdownDescription: "Flag indicating if the unifiedRoleDefinition is part of the default set included with the product or custom. Read-only.  Supports `$filter` (`eq`).",
		},
		"is_enabled": schema.BoolAttribute{
			Required:            true,
			MarkdownDescription: "Flag indicating if the role is enabled for assignment. If false the role is not available for assignment. Read-only when **isBuiltIn** is `true`.",
		},
		"is_privileged": schema.BoolAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.Bool{wpdefaultvalue.BoolDefaultValue(false)},
			Computed:            true,
			MarkdownDescription: "Flag indicating if the role is privileged. Microsoft Entra ID defines a role as privileged if it contains at least one sensitive resource action in the **rolePermissions** and **allowedResourceActions** objects. Applies only for actions in the `microsoft.directory` resource namespace. Read-only. Supports `$filter` (`eq`). The _provider_ default value is `false`.",
		},
		"role_permissions": schema.SetNestedAttribute{
			Required: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // unifiedRolePermission
					"allowed_resource_actions": schema.SetAttribute{
						ElementType:         types.StringType,
						Required:            true,
						MarkdownDescription: "Set of tasks that can be performed on a resource.",
					},
					"condition": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "Optional constraints that must be met for the permission to be effective. Not supported for custom roles.",
					},
				},
			},
			MarkdownDescription: "List of permissions included in the role. Read-only when **isBuiltIn** is `true`. Required. / Represents a collection of allowed resource actions and the conditions that must be met for the action to be effective. Resource actions are tasks that can be performed on a resource. For example, the application resource supports create, update, delete, and reset password resource actions. / https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolepermission?view=graph-rest-beta",
		},
		"template_id": schema.StringAttribute{
			Computed: true,
			Optional: true,
			PlanModifiers: []planmodifier.String{
				wpplanmodifier.StringUseStateForUnknown(),
				stringplanmodifier.RequiresReplace(),
			},
			MarkdownDescription: "Custom template identifier that can be set when isBuiltIn is `false`. This identifier is typically used if one needs an identifier to be the same across different directories. Read-only when **isBuiltIn** is `true`.",
		},
		"version": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Indicates the version of the unifiedRoleDefinition object. Read-only when **isBuiltIn** is `true`.",
		},
		"inherits_permissions_from": schema.SetNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // unifiedRoleDefinition
					"id": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The unique identifier for the unifiedRoleDefinition. Key, not nullable, Read-only.  Supports `$filter` (`eq` operator only).",
					},
				},
			},
			PlanModifiers:       []planmodifier.Set{wpplanmodifier.SetUseStateForUnknown()},
			MarkdownDescription: "Read-only collection of role definitions that the given role definition inherits from. Only Microsoft Entra built-in roles support this attribute.",
		},
	},
	MarkdownDescription: "Represents a collection of permissions listing the operations, such as read, write, and delete, that can be performed by an RBAC provider, as part of Microsoft 365 RBAC [role management](rolemanagement.md).\n\nThe following RBAC providers are currently supported: / https://learn.microsoft.com/en-us/graph/api/resources/unifiedroledefinition?view=graph-rest-beta ||| MS Graph: Role management",
}
