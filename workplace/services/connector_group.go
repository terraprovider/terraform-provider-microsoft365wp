package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	ConnectorGroupResource = generic.GenericResource{
		TypeNameSuffix: "connector_group",
		SpecificSchema: connectorGroupResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/onPremisesPublishingProfiles/applicationProxy/connectorGroups",
			ReadOptions: generic.ReadOptions{
				ODataExpand: "applications,members",
				DataSource: generic.DataSourceOptions{
					ExtraFilterAttributes: []string{"is_default"},
					Plural: generic.PluralOptions{
						ExtraAttributes: []string{"is_default"},
					},
				},
			},
		},
	}

	ConnectorGroupSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&ConnectorGroupResource)

	ConnectorGroupPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&ConnectorGroupResource, "")
)

var connectorGroupResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // connectorGroup
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Unique identifier for this connectorGroup. Read-only.",
		},
		"connector_group_type": schema.StringAttribute{
			Computed:            true,
			Validators:          []validator.String{stringvalidator.OneOf("applicationProxy")},
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Indicates the type of hybrid agent. This pre-set by the system. Read-only. <br/> _Provider_ allowed values are: `applicationProxy`.",
		},
		"is_default": schema.BoolAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.Bool{wpplanmodifier.BoolUseStateForUnknown()},
			MarkdownDescription: "Indicates if the connectorGroup is the default connectorGroup. Only a single connector group can be the default connectorGroup and this is pre-set by the system. Read-only.",
		},
		"name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The name associated with the connectorGroup.",
		},
		"region": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf("nam", "eur", "aus", "asia", "ind", "unknownFutureValue"),
			},
			PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			MarkdownDescription: "The region the connectorGroup is assigned to and will optimize traffic for. This region can only be set if **no connectors or applications** are assigned to the connectorGroup. <br/> _Provider_ allowed values are: `nam`, `eur`, `aus`, `asia`, `ind`, `unknownFutureValue`.",
		},
		"applications": schema.SetNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // application
					"id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Unique identifier for the application object. This property is referred to as **Object ID** in the Microsoft Entra admin center. Key. Not nullable. Read-only. Supports `$filter` (`eq`, `ne`, `not`, `in`).",
					},
					"app_id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The unique identifier for the application that is assigned by Microsoft Entra ID. Not nullable. Read-only. Alternate key. Supports `$filter` (`eq`).",
					},
					"display_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The display name for the application. Maximum length is 256 characters. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values), `$search`, and `$orderby`.",
					},
				},
			},
			MarkdownDescription: "Read-only. Nullable. / Represents an application. Any application that outsources authentication to Microsoft Entra ID must be registered in the Microsoft identity platform. Application registration involves telling Microsoft Entra ID about your application, including the URL where it's located, the URL to send replies after authentication, the URI to identify your application, and more. Also see [Microsoft docs for application](https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta). <br> ",
		},
		"members": schema.SetNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // connector
					"id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The unique identifier of the connector. Read-only.",
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
			},
			MarkdownDescription: "Read-only. Nullable. / Represents an Application Proxy connector. Connectors are lightweight agents that sit on-premises and facilitate the outbound connection to the [Microsoft Entra application proxy](https://learn.microsoft.com/en-us/entra/identity/app-proxy/overview-what-is-app-proxy) service. Each connector is part of a [connectorGroup](https://learn.microsoft.com/en-us/graph/api/resources/connectorgroup?view=graph-rest-beta). Also see [Microsoft docs for connector](https://learn.microsoft.com/en-us/graph/api/resources/connector?view=graph-rest-beta). <br> ",
		},
	},
	MarkdownDescription: "Each [Microsoft Entra application proxy](https://learn.microsoft.com/en-us/entra/identity/app-proxy/overview-what-is-app-proxy) connector is always part of a connector group. All the connectors that belong to the same connector group act as a separate unit for high-availability and load balancing. If you don't create connector groups, all your connectors will be part of the default group. When configuring an application with Application Proxy, you must also specify which connector group to assign the application to.\n\nAfter a connector group is created, you can add or move connectors to the connector group by using [Add connector](https://learn.microsoft.com/en-us/graph/api/connectorgroup-post-members?view=graph-rest-beta). You can also use [Add application](https://learn.microsoft.com/en-us/graph/api/connectorgroup-post-applications?view=graph-rest-beta) to assign an application to a connector group.\n\nAlso see [Microsoft docs for connectorGroup](https://learn.microsoft.com/en-us/graph/api/resources/connectorgroup?view=graph-rest-beta).\n\n_Provider_ Note: Use the separate resources `microsoft365wp_connector_group_member_connector` and `microsoft365wp_connector_group_application` to add/assign connectors or applications to a connector group. ||| MS Graph: On-premises publishing",
}
