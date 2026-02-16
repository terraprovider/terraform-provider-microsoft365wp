package services

import (
	"context"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvaluemodifier"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	MobilityManagementPolicyResource = generic.GenericResource{
		TypeNameSuffix: "mobility_management_policy",
		SpecificSchema: mobilityManagementPolicyResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/policies",
			ParentEntities: generic.ParentEntities{
				{
					ParentIdField: path.Root("policy_type"),
					ParentIdValueMap: map[string]string{
						"device": "mobileDeviceManagementPolicies",
						"app":    "mobileAppManagementPolicies",
					},
				},
			},
			ReadOptions: generic.ReadOptions{
				ODataExpand: "includedGroups",
			},
			WriteOptions: generic.WriteOptions{
				UpdateInsteadOfCreate: true,
				SkipDelete:            true,
				SubActions: []generic.WriteSubAction{
					&generic.WriteSubActionIndividual{
						ComparisonKeyAttribute: "id",
						SetNestedPath:          tftypes.NewAttributePath().WithAttributeName("included_groups"),
						IsOdataReference:       true,
						OdataRefMapTypeToUriPrefix: map[string]string{
							"": "https://graph.microsoft.com/odata/groups/",
						},
						WriteSubActionBase: generic.WriteSubActionBase{
							Attributes: []string{"includedGroups"},
							UriSuffix:  "includedGroups",
						},
					},
					&generic.WriteSubActionAllInOne{
						UsePatch: true,
						WriteSubActionBase: generic.WriteSubActionBase{
							Attributes: []string{"appliesTo"},
						},
					},
				},
			},
			TerraformToGraphMiddleware: mobilityManagementPolicyTerraformToGraphMiddleware,
		},
	}

	MobilityManagementPolicySingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&MobilityManagementPolicyResource)

	MobilityManagementPolicyPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&MobilityManagementPolicyResource, "")
)

func mobilityManagementPolicyTerraformToGraphMiddleware(ctx context.Context, diags *diag.Diagnostics, params *generic.TerraformToGraphMiddlewareParams) generic.TerraformToGraphMiddlewareReturns {
	// "selected" gets implied automatically if includedGroups contains any entries but it may not be written to Graph
	if params.IsUpdate && params.RawVal["appliesTo"] == "selected" {
		delete(params.RawVal, "appliesTo")
	}
	return nil
}

var mobilityManagementPolicyResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // mobilityManagementPolicy
		"policy_type": schema.StringAttribute{
			Required:            true,
			Validators:          []validator.String{stringvalidator.OneOf("device", "app")},
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			MarkdownDescription: "_Provider_ Note: Type of policy - must either be `device` or `app`",
		},
		"id": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Object Id of the mobility management application.",
		},
		"applies_to": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("none", "all", "selected", "unknownFutureValue"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvaluemodifier.StringDefaultValue("all")},
			Computed:            true,
			MarkdownDescription: "Indicates the user scope of the mobility management policy. <br/> _Provider_ allowed values are: `none`, `all`, `selected`, `unknownFutureValue`. The _provider_ default value is `\"all\"`.  \n_Provider_ Note: Please note that this attribute's value must be `selected` if `included_groups` contains any entries.",
		},
		"compliance_url": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Compliance URL of the mobility management application.",
		},
		"description": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Description of the mobility management application.",
		},
		"discovery_url": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Discovery URL of the mobility management application.",
		},
		"display_name": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Display name of the mobility management application.",
		},
		"is_valid": schema.BoolAttribute{
			Computed:            true,
			MarkdownDescription: "Whether policy is valid. Invalid policies may not be updated and should be deleted.",
		},
		"terms_of_use_url": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Terms of Use URL of the mobility management application.",
		},
		"included_groups": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // group
					"id": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The unique identifier for the group. <br/> Returned by default. Key. Not nullable. Read-only. <br/> Supports `$filter` (`eq`, `ne`, `not`, `in`).",
					},
					"display_name": schema.StringAttribute{
						Computed:            true,
						PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
						MarkdownDescription: "The display name for the group. Required. Maximum length is 256 characters. <br/> Returned by default. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values), `$search`, and `$orderby`.",
					},
				},
			},
			PlanModifiers:       []planmodifier.Set{wpdefaultvaluemodifier.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Microsoft Entra groups under the scope of the mobility management application if appliesTo is `selected` <br/> Represents a Microsoft Entra group, which can be a Microsoft 365 group, a team in Microsoft Teams, or a security group. <br/> For performance reasons, the [create](https://learn.microsoft.com/en-us/graph/api/group-post-groups?view=graph-rest-beta), [get](https://learn.microsoft.com/en-us/graph/api/group-get?view=graph-rest-beta), and [list](https://learn.microsoft.com/en-us/graph/api/group-list?view=graph-rest-beta) operations return only a subset of more commonly used properties by default. These _default_ properties are noted in the [Properties](#properties) section. To get any of the properties not returned by default, specify them in a `$select` OData query option. Also see [Microsoft docs for group](https://learn.microsoft.com/en-us/graph/api/resources/group?view=graph-rest-beta). <br/> The _provider_ default value is `[]`.  \n_Provider_ Note: Please note that this attribute's must not contain any entries if `applies_to` is `none` or `all`. <br> ",
		},
	},
	MarkdownDescription: "In Microsoft Entra ID, a mobility management policy represents an autoenrollment configuration for a mobility management (MDM or MAM) application. These policies are only applicable to devices based on Windows 10 OS and its derivatives such as Surface Hub and HoloLens. [Autoenrollment](https://learn.microsoft.com/en-us/windows/client-management/azure-ad-and-microsoft-intune-automatic-mdm-enrollment-in-the-new-portal) automatically enrolls Windows 10 devices into mobility management applications during [Microsoft Entra join](https://learn.microsoft.com/en-us/entra/identity/devices/concept-directory-join) or [Microsoft Entra register](https://learn.microsoft.com/en-us/entra/identity/devices/concept-device-registration) processes. <br/> Also see [Microsoft docs for mobilityManagementPolicy](https://learn.microsoft.com/en-us/graph/api/resources/mobilitymanagementpolicy?view=graph-rest-beta).\n\n_Provider_ Note: Please note that MS Graph does currently (as of 2025-02) not allow to update (or even read) this entity using\napplication permissions but only using work or school account delegated permissions (also see\nhttps://learn.microsoft.com/en-us/graph/api/mobiledevicemanagementpolicies-get?view=graph-rest-beta#permissions). ||| MS Graph: App management",
}
