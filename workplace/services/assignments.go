package services

import (
	"fmt"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var deviceAndAppManagementAssignment = schema.SetNestedAttribute{
	Optional: true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{ // ...Assignment
			// We cannot include any computed attributes here as SetUnknownValuesFromResourceModel will fail if values
			// are unknown on multiple levels (e.g. here and also inside `target`).
			"target": deviceAndAppManagementAssignmentTarget,
		},
	},
	Computed: true,
	PlanModifiers: []planmodifier.Set{
		wpdefaultvalue.SetDefaultValueEmpty(),
	},
	MarkdownDescription: `The list of assignments.`,
}

var deviceAndAppManagementAssignmentTarget = schema.SingleNestedAttribute{
	Required: true,
	Attributes: map[string]schema.Attribute{ // deviceAndAppManagementAssignmentTarget
		"filter_id": schema.StringAttribute{
			Computed:            true,
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			Description:         `deviceAndAppManagementAssignmentFilterId`, // custom MS Graph attribute name
			MarkdownDescription: "The ID of the filter for the target assignment.",
		},
		"filter_type": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("none", "include", "exclude"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("none")},
			Computed:            true,
			Description:         `deviceAndAppManagementAssignmentFilterType`, // custom MS Graph attribute name
			MarkdownDescription: "The type of filter of the target assignment i.e. Exclude or Include. / Represents type of the assignment filter; possible values are: `none` (Default value. Do not use.), `include` (Indicates in-filter, rule matching will offer the payload to devices.), `exclude` (Indicates out-filter, rule matching will not offer the payload to devices.). The _provider_ default value is `\"none\"`.",
		},
		"all_devices": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.allDevicesAssignmentTarget",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional:   true,
				Attributes: map[string]schema.Attribute{ // allDevicesAssignmentTarget
				},
				Validators: []validator.Object{
					deviceAndAppManagementAssignmentTargetDeviceAndAppManagementAssignmentTargetValidator,
				},
				MarkdownDescription: "Represents an assignment to all managed devices in the tenant. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-alldevicesassignmenttarget?view=graph-rest-beta",
			},
		},
		"all_licensed_users": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.allLicensedUsersAssignmentTarget",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional:   true,
				Attributes: map[string]schema.Attribute{ // allLicensedUsersAssignmentTarget
				},
				Validators: []validator.Object{
					deviceAndAppManagementAssignmentTargetDeviceAndAppManagementAssignmentTargetValidator,
				},
				MarkdownDescription: "Represents an assignment to all licensed users in the tenant. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-alllicensedusersassignmenttarget?view=graph-rest-beta",
			},
		},
		"exclusion_group": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.exclusionGroupAssignmentTarget",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // exclusionGroupAssignmentTarget
					"group_id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The group Id that is the target of the assignment.",
					},
				},
				Validators: []validator.Object{
					deviceAndAppManagementAssignmentTargetDeviceAndAppManagementAssignmentTargetValidator,
				},
				MarkdownDescription: "Represents a group that should be excluded from an assignment. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-exclusiongroupassignmenttarget?view=graph-rest-beta",
			},
		},
		"group": generic.OdataDerivedTypeNestedAttributeRs{
			DerivedType: "#microsoft.graph.groupAssignmentTarget",
			SingleNestedAttribute: schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{ // groupAssignmentTarget
					"group_id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The group Id that is the target of the assignment.",
					},
				},
				Validators: []validator.Object{
					deviceAndAppManagementAssignmentTargetDeviceAndAppManagementAssignmentTargetValidator,
				},
				MarkdownDescription: "Represents an assignment to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-groupassignmenttarget?view=graph-rest-beta",
			},
		},
	},
	MarkdownDescription: "Base type for assignment targets. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-deviceandappmanagementassignmenttarget?view=graph-rest-beta",
}

var deviceAndAppManagementAssignmentTargetDeviceAndAppManagementAssignmentTargetValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("all_devices"),
	path.MatchRelative().AtParent().AtName("all_licensed_users"),
	path.MatchRelative().AtParent().AtName("exclusion_group"),
	path.MatchRelative().AtParent().AtName("group"),
)

func GetAssignmentChildResource(parentResource *generic.GenericResource, singleItemUseODataFilter bool) generic.GenericResource {

	parentIdFieldName := fmt.Sprintf("%s_id", parentResource.TypeNameSuffix)

	attributes := map[string]schema.Attribute{ // deviceConfigurationAssignment
		parentIdFieldName: schema.StringAttribute{
			Required: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The key of the assignment.",
		},
		"source": schema.StringAttribute{
			Computed:            true,
			Validators:          []validator.String{stringvalidator.OneOf("direct", "policySets")},
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The source of the assignment. This property is read-only. / Represents source of assignment; possible values are: `direct` (Direct indicates a direct assignment.), `policySets` (PolicySets indicates assignment was made via PolicySet assignment.)",
		},
		"source_id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The identifier of the source of the assignment. This property is read-only.",
		},
	}

	targetAttribute := deviceAndAppManagementAssignmentTarget // copy it
	targetAttribute.PlanModifiers = append(targetAttribute.PlanModifiers, objectplanmodifier.RequiresReplace())
	attributes["target"] = targetAttribute

	return generic.GenericResource{
		TypeNameSuffix: fmt.Sprintf("%s_assignment", parentResource.TypeNameSuffix),
		SpecificSchema: schema.Schema{
			Attributes: attributes,
			MarkdownDescription: "\n\n**Note**: Due to technical difficulties trying to import this resource from MS Graph will result in an error, i.e. importing of this resource is not possible." +
				generic.GetSubcategorySuffixFromMarkdownDescription(parentResource.SpecificSchema.MarkdownDescription),
		},
		AccessParams: generic.AccessParams{
			BaseUri: parentResource.AccessParams.BaseUri,
			HasParentItem: generic.HasParentItem{
				ParentIdField: path.Root(parentIdFieldName),
				UriSuffix:     "/assignments",
			},
			ReadOptions: generic.ReadOptions{
				SingleItemUseODataFilter: singleItemUseODataFilter,
			},
		},
	}
}
