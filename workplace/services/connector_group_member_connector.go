package services

import (
	"context"
	"terraform-provider-microsoft365wp/workplace/external/msgraph"
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	ConnectorGroupMemberConnectorResource = generic.GenericResource{
		TypeNameSuffix: "connector_group_member_connector",
		SpecificSchema: connectorGroupMemberConnectorResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/onPremisesPublishingProfiles/applicationProxy/connectorGroups",
			ParentEntities: generic.ParentEntities{
				{
					ParentIdField: path.Root("connector_group_id"),
					UriSuffix:     "members",
				},
			},
			ReadOptions: generic.ReadOptions{
				ODataSelect:              []string{"id", "externalIp", "machineName"},
				SingleItemUseODataFilter: true,
				DataSource: generic.DataSourceOptions{
					ExtraFilterAttributes: []string{"external_ip", "machine_name"},
					Plural: generic.PluralOptions{
						ExtraAttributes: []string{"external_ip", "machine_name"},
					},
				},
			},
			CreateReplaceFunc: connectorGroupMemberConnectorCreateReplaceFunc,
			DeleteReplaceFunc: connectorGroupMemberConnectorDeleteReplaceFunc,
		},
	}

	ConnectorGroupMemberConnectorSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&ConnectorGroupMemberConnectorResource)

	ConnectorGroupMemberConnectorPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&ConnectorGroupMemberConnectorResource, "")
)

func connectorGroupMemberConnectorCreateReplaceFunc(ctx context.Context, diags *diag.Diagnostics, params *generic.CreateReplaceFuncParams) {
	connectorGroupMemberConnectorAddDelete(ctx, diags, &params.R.AccessParams, params.Client, params.BaseUri, params.Id, params.IdAttributer, false)
}

func connectorGroupMemberConnectorDeleteReplaceFunc(ctx context.Context, diags *diag.Diagnostics, params *generic.DeleteReplaceFuncParams) {
	connectorGroupMemberConnectorAddDelete(ctx, diags, &params.R.AccessParams, params.Client, params.BaseUri, params.Id, params.IdAttributer, true)
}

func connectorGroupMemberConnectorAddDelete(ctx context.Context, diags *diag.Diagnostics, aps *generic.AccessParams,
	client *msgraph.Client, baseUri string, id string, idAttributer generic.GetAttributer, isDelete bool) {

	parentIdAttributer := idAttributer
	// As connectors cannot be deleted from a connector group, we instead just assign (i.e. move) them to the default connector group.
	if isDelete {
		// Determine id of default connector group
		defaultConnectorGroupId := aps.ReadId(ctx, diags, aps.BaseUri, nil, true, "isDefault eq true", nil)
		if diags.HasError() {
			return
		}
		// Create a dummy state with default connector group id to use as parentIdAttributer.
		// From it, GetBaseUri will read the id of the connector group to which to assign the connector,
		// i.e. the connector will be assigned to the default connector group instead.
		parentIdAttributer = tfsdk.State{
			Schema: connectorGroupMemberConnectorResourceSchema,
			Raw: tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"connector_group_id": tftypes.String,
					},
				},
				map[string]tftypes.Value{
					"connector_group_id": tftypes.NewValue(tftypes.String, defaultConnectorGroupId),
				},
			),
		}
	}

	uri := aps.GetBaseUri(ctx, diags, baseUri, parentIdAttributer)
	if diags.HasError() {
		return
	}
	uri.Entity += "/$ref"

	id = aps.GetId(ctx, diags, id, idAttributer)
	if diags.HasError() {
		return
	}

	postRawVal := map[string]any{
		"@odata.id": "https://graph.microsoft.com/beta/onPremisesPublishingProfiles/applicationProxy/connectors/" + id,
	}
	generic.CreateRaw(ctx, diags, client, uri, postRawVal, nil, true, false)
}

var connectorGroupMemberConnectorResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // connector
		"connector_group_id": schema.StringAttribute{
			Required:            true,
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			MarkdownDescription: "_Provider_ Note: ID of the connector group to add the connector to. Required.",
		},
		"id": schema.StringAttribute{
			Required:            true,
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			MarkdownDescription: "The unique identifier of the connector. Read-only.  \n_Provider_ Note: In this special case, it is _not read-only_ but _required_ instead.",
		},
		"external_ip": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The external IP address as detected by the connector server. Read-only.",
		},
		"machine_name": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The name of the computer on which the connector is installed and runs on.",
		},
	},
	MarkdownDescription: "Use this resource to add/assign an existing connector to a connector group.\n\nNotes:  \n  - When deleting this resource, the connector will be assigned back to the default connector group.  \n  - To import this resource, an ID consisting of `connector_group_id` and `id` being joined by a forward slash (`/`) must be used. ||| MS Graph: On-premises publishing",
}
