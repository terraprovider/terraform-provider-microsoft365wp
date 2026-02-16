package services

import (
	"context"
	"terraform-provider-microsoft365wp/workplace/external/msgraph"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

var (
	ConnectorGroupApplicationResource = generic.GenericResource{
		TypeNameSuffix: "connector_group_application",
		SpecificSchema: connectorGroupApplicationResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/onPremisesPublishingProfiles/applicationProxy/connectorGroups",
			ParentEntities: generic.ParentEntities{
				{
					ParentIdField: path.Root("connector_group_id"),
					UriSuffix:     "applications",
				},
			},
			ReadOptions: generic.ReadOptions{
				ODataSelect:              []string{"appId", "displayName"},
				SingleItemUseODataFilter: true,
				DataSource: generic.DataSourceOptions{
					Plural: generic.PluralOptions{
						ExtraAttributes: []string{"app_id", "display_name"},
					},
				},
			},
			CreateReplaceFunc: connectorGroupApplicationCreateReplaceFunc,
			DeleteReplaceFunc: connectorGroupApplicationDeleteReplaceFunc,
		},
	}

	ConnectorGroupApplicationSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&ConnectorGroupApplicationResource)

	ConnectorGroupApplicationPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&ConnectorGroupApplicationResource, "")
)

func connectorGroupApplicationCreateReplaceFunc(ctx context.Context, diags *diag.Diagnostics, params *generic.CreateReplaceFuncParams) {
	connectorGroupApplicationAddDelete(ctx, diags, &params.R.AccessParams, params.Client, params.BaseUri, params.Id, params.IdAttributer, false)
}

func connectorGroupApplicationDeleteReplaceFunc(ctx context.Context, diags *diag.Diagnostics, params *generic.DeleteReplaceFuncParams) {
	connectorGroupApplicationAddDelete(ctx, diags, &params.R.AccessParams, params.Client, params.BaseUri, params.Id, params.IdAttributer, true)
}

func connectorGroupApplicationAddDelete(ctx context.Context, diags *diag.Diagnostics, aps *generic.AccessParams,
	client *msgraph.Client, baseUri string, id string, idAttributer generic.GetAttributer, isDelete bool) {

	applicationId := aps.GetId(ctx, diags, id, idAttributer)
	if diags.HasError() {
		return
	}

	var connectorGroupId string
	if !isDelete {
		diags.Append(idAttributer.GetAttribute(ctx, aps.ParentEntities[0].ParentIdField, &connectorGroupId)...)
		if diags.HasError() {
			return
		}
	} else {
		// As applications cannot be deleted from a connector group, we instead just assign (i.e. move) them to the default connector group.
		connectorGroupId = aps.ReadId(ctx, diags, aps.BaseUri, nil, true, "isDefault eq true", nil)
		if diags.HasError() {
			return
		}
	}

	uri := msgraph.Uri{Entity: "/applications/" + applicationId + "/connectorGroup/$ref"}
	rawVal := map[string]any{
		"@odata.id": "https://graph.microsoft.com/beta/onPremisesPublishingProfiles/applicationproxy/connectorGroups/" + connectorGroupId,
	}
	generic.CreateRaw(ctx, diags, client, uri, rawVal, nil, true, true)
}

var connectorGroupApplicationResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // application
		"connector_group_id": schema.StringAttribute{
			Required:            true,
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			MarkdownDescription: "_Provider_ Note: ID of the connector group to add the application to. Required.",
		},
		"id": schema.StringAttribute{
			Required:            true,
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			MarkdownDescription: "Unique identifier for the application object. This property is referred to as **Object ID** in the Microsoft Entra admin center. Key. Not nullable. Read-only. Supports `$filter` (`eq`, `ne`, `not`, `in`).  \n_Provider_ Note: In this special case, it is _not read-only_ but _required_ instead. Set it to the **Object ID** of the **App registration** (which is different from both the Application/Client ID and the Object ID of the Enterprise application).",
		},
		"app_id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The unique identifier for the application that is assigned by Microsoft Entra ID. Not nullable. Read-only. Alternate key. Supports `$filter` (`eq`).",
		},
		"display_name": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The display name for the application. Maximum length is 256 characters. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values), `$search`, and `$orderby`.",
		},
	},
	MarkdownDescription: "Use this resource to add/assign an existing application to a connector group.\n\nNotes:  \n  - When deleting this resource, the application will be assigned back to the default connector group.  \n  - To import this resource, an ID consisting of `connector_group_id` and `id` being joined by a forward slash (`/`) must be used. ||| MS Graph: On-premises publishing",
}
