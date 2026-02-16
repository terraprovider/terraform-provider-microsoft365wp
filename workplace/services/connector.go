package services

import (
	"terraform-provider-microsoft365wp/workplace/generic"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	connectorResource = generic.GenericResource{
		TypeNameSuffix: "connector",
		SpecificSchema: connectorResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/onPremisesPublishingProfiles/applicationProxy/connectors",
			ReadOptions: generic.ReadOptions{
				ODataExpand:              "memberOf",
				SingleItemUseODataFilter: true,
				DataSource: generic.DataSourceOptions{
					ExtraFilterAttributes: []string{"external_ip", "machine_name"},
					Plural: generic.PluralOptions{
						ExtraAttributes: []string{"external_ip", "machine_name"},
					},
				},
			},
		},
	}

	ConnectorSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&connectorResource)

	ConnectorPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&connectorResource, "")
)

var connectorResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // connector
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier of the connector. Read-only.",
		},
		"external_ip": schema.StringAttribute{
			MarkdownDescription: "The external IP address as detected by the connector server. Read-only.",
		},
		"machine_name": schema.StringAttribute{
			MarkdownDescription: "The name of the computer on which the connector is installed and runs on.",
		},
		"status": schema.StringAttribute{
			Validators:          []validator.String{stringvalidator.OneOf("active", "inactive")},
			MarkdownDescription: "Indicates the status of the connector. Read-only. <br/> _Provider_ allowed values are: `active`, `inactive`.",
		},
		"version": schema.StringAttribute{
			MarkdownDescription: "The version of the connector. Read-only.",
		},
		"member_of": schema.SetNestedAttribute{
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // connectorGroup
					"id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Unique identifier for this connectorGroup. Read-only.",
					},
					"is_default": schema.BoolAttribute{
						Computed:            true,
						MarkdownDescription: "Indicates if the connectorGroup is the default connectorGroup. Only a single connector group can be the default connectorGroup and this is pre-set by the system. Read-only.",
					},
					"name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "The name associated with the connectorGroup.",
					},
				},
			},
			MarkdownDescription: "The **connectorGroup** that the connector is a member of. Read-only. <br/> Each [Microsoft Entra application proxy](https://learn.microsoft.com/en-us/entra/identity/app-proxy/overview-what-is-app-proxy) connector is always part of a connector group. All the connectors that belong to the same connector group act as a separate unit for high-availability and load balancing. If you don't create connector groups, all your connectors will be part of the default group. When configuring an application with Application Proxy, you must also specify which connector group to assign the application to. <br/> After a connector group is created, you can add or move connectors to the connector group by using [Add connector](https://learn.microsoft.com/en-us/graph/api/connectorgroup-post-members?view=graph-rest-beta). You can also use [Add application](https://learn.microsoft.com/en-us/graph/api/connectorgroup-post-applications?view=graph-rest-beta) to assign an application to a connector group. Also see [Microsoft docs for connectorGroup](https://learn.microsoft.com/en-us/graph/api/resources/connectorgroup?view=graph-rest-beta). <br> ",
		},
	},
	MarkdownDescription: "Represents an Application Proxy connector. Connectors are lightweight agents that sit on-premises and facilitate the outbound connection to the [Microsoft Entra application proxy](https://learn.microsoft.com/en-us/entra/identity/app-proxy/overview-what-is-app-proxy) service. Each connector is part of a [connectorGroup](https://learn.microsoft.com/en-us/graph/api/resources/connectorgroup?view=graph-rest-beta). <br/> Also see [Microsoft docs for connector](https://learn.microsoft.com/en-us/graph/api/resources/connector?view=graph-rest-beta). ||| MS Graph: On-premises publishing",
}
