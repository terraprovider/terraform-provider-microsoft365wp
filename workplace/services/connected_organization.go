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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	ConnectedOrganizationResource = generic.GenericResource{
		TypeNameSuffix: "connected_organization",
		SpecificSchema: connectedOrganizationResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri: "/identityGovernance/entitlementManagement/connectedOrganizations",
			ReadOptions: generic.ReadOptions{
				ExtraRequests: []generic.ReadExtraRequest{
					{
						Attribute: "internalSponsors",
					},
					{
						Attribute: "externalSponsors",
					},
				},
			},
			WriteOptions: generic.WriteOptions{
				SubActions: []generic.WriteSubAction{
					&generic.WriteSubActionIndividual{
						WriteSubActionBase: generic.WriteSubActionBase{
							Attributes: []string{"internalSponsors"},
							UriSuffix:  "internalSponsors",
							UpdateOnly: true,
						},
						ComparisonKeyAttribute: "id",
						SetNestedPath:          tftypes.NewAttributePath().WithAttributeName("internal_sponsors"),
						IsOdataReference:       true,
						OdataRefMapTypeToUriPrefix: map[string]string{
							"": "https://graph.microsoft.com/beta/directoryObjects/", // this will work for users and groups
						},
					},
					&generic.WriteSubActionIndividual{
						WriteSubActionBase: generic.WriteSubActionBase{
							Attributes: []string{"externalSponsors"},
							UriSuffix:  "externalSponsors",
							UpdateOnly: true,
						},
						ComparisonKeyAttribute: "id",
						SetNestedPath:          tftypes.NewAttributePath().WithAttributeName("external_sponsors"),
						IsOdataReference:       true,
						OdataRefMapTypeToUriPrefix: map[string]string{
							"": "https://graph.microsoft.com/beta/directoryObjects/", // this will work for users and groups
						},
					},
				},
			},
		},
	}

	ConnectedOrganizationSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&ConnectedOrganizationResource)

	ConnectedOrganizationPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&ConnectedOrganizationResource, "")
)

var connectedOrganizationResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // connectedOrganization
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "Read-only.",
		},
		"created_by": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "UPN of the user who created this resource. Read-only.",
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Read-only.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The description of the connected organization.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "The display name of the connected organization. Supports `$filter` (`eq`).",
		},
		"identity_sources": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{ // identitySource
					"azure_active_directory_tenant": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.azureActiveDirectoryTenant",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // azureActiveDirectoryTenant
								"display_name": schema.StringAttribute{
									Required:            true,
									PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
									MarkdownDescription: "The name of the Microsoft Entra tenant. Read only.",
								},
								"tenant_id": schema.StringAttribute{
									Required:            true,
									PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
									MarkdownDescription: "The ID of the Microsoft Entra tenant. Read only.",
								},
							},
							Validators:          []validator.Object{connectedOrganizationIdentitySourceValidator},
							MarkdownDescription: "Used in the identity sources of an [connectedOrganization](connectedOrganization.md). The `@odata.type` value `#microsoft.graph.azureActiveDirectoryTenant` indicates that this type identifies another Microsoft Entra tenant as an identity source for a connected organization. / https://learn.microsoft.com/en-us/graph/api/resources/azureactivedirectorytenant?view=graph-rest-beta",
						},
					},
					"cross_cloud_azure_active_directory_tenant": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.crossCloudAzureActiveDirectoryTenant",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // crossCloudAzureActiveDirectoryTenant
								"cloud_instance": schema.StringAttribute{
									Required:            true,
									PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
									MarkdownDescription: "The ID of the cloud where the tenant is located, one of `microsoftonline.com`, `microsoftonline.us` or `partner.microsoftonline.cn`. Read only.",
								},
								"display_name": schema.StringAttribute{
									Required:            true,
									PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
									MarkdownDescription: "The name of the Microsoft Entra tenant. Read only.",
								},
								"tenant_id": schema.StringAttribute{
									Required:            true,
									PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
									MarkdownDescription: "The ID of the Microsoft Entra tenant. Read only.",
								},
							},
							Validators:          []validator.Object{connectedOrganizationIdentitySourceValidator},
							MarkdownDescription: "Used in the identity sources of a [connectedOrganization](connectedOrganization.md) object. The `@odata.type` value `#microsoft.graph.crossCloudAzureActiveDirectoryTenant` indicates that this type identifies another Microsoft Entra tenant in a different cloud as an identity source for a connected organization. / https://learn.microsoft.com/en-us/graph/api/resources/crosscloudazureactivedirectorytenant?view=graph-rest-beta",
						},
					},
					"domain_identity_source": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.domainIdentitySource",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // domainIdentitySource
								"display_name": schema.StringAttribute{
									Required:            true,
									PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
									MarkdownDescription: "The name of the identity source, typically also the domain name. Read only.",
								},
								"domain_name": schema.StringAttribute{
									Required:            true,
									PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
									MarkdownDescription: "The domain name. Read only.",
								},
							},
							Validators:          []validator.Object{connectedOrganizationIdentitySourceValidator},
							MarkdownDescription: "Used in the identity sources of an [connectedOrganization](connectedOrganization.md). The `@odata.type` value `#microsoft.graph.domainIdentitySource` indicates that this type identifies a domain as an identity source for a connected organization. / https://learn.microsoft.com/en-us/graph/api/resources/domainidentitysource?view=graph-rest-beta",
						},
					},
					"external_domain_federation": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.externalDomainFederation",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // externalDomainFederation
								"display_name": schema.StringAttribute{
									Required:            true,
									PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
									MarkdownDescription: "The name of the identity source, typically also the domain name. Read only.",
								},
								"domain_name": schema.StringAttribute{
									Required:            true,
									PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
									MarkdownDescription: "The domain name. Read only.",
								},
								"issuer_uri": schema.StringAttribute{
									Required:            true,
									PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
									MarkdownDescription: "The issuerURI of the incoming federation. Read only.",
								},
							},
							Validators:          []validator.Object{connectedOrganizationIdentitySourceValidator},
							MarkdownDescription: "Used in the identity sources of an [connectedOrganization](connectedOrganization.md). The `@odata.type` value `#microsoft.graph.externalDomainFederation` indicates that this type identifies a domain with a configured identity provider as an identity source for a connected organization. / https://learn.microsoft.com/en-us/graph/api/resources/externaldomainfederation?view=graph-rest-beta",
						},
					},
					"social_identity_source": generic.OdataDerivedTypeNestedAttributeRs{
						DerivedType: "#microsoft.graph.socialIdentitySource",
						SingleNestedAttribute: schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{ // socialIdentitySource
								"display_name": schema.StringAttribute{
									Required:            true,
									PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
									MarkdownDescription: "The name of the identity source. Typically the same value as the **socialIdentitySourceType**.",
								},
								"social_identity_source_type": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf("facebook", "unknownFutureValue"),
									},
									PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
									MarkdownDescription: "The / Possible values are: `facebook`, `unknownFutureValue`",
								},
							},
							Validators:          []validator.Object{connectedOrganizationIdentitySourceValidator},
							MarkdownDescription: "Used in the identity sources of an [connectedOrganization](connectedOrganization.md). The `@odata.type` value `#microsoft.graph.socialIdentitySource` identifies a social identity as an identity source for a connected organization. / https://learn.microsoft.com/en-us/graph/api/resources/socialidentitysource?view=graph-rest-beta",
						},
					},
				},
			},
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "The identity sources in this connected organization, one of [azureActiveDirectoryTenant](azureactivedirectorytenant.md), [crossCloudAzureActiveDirectoryTenant](crosscloudazureactivedirectorytenant.md), [domainIdentitySource](domainidentitysource.md), [externalDomainFederation](externaldomainfederation.md), or [socialIdentitySource](socialidentitysource.md). Read-only. Nullable. Supports `$select` and `$filter`(`eq`). To filter by the derived types, you must declare the resource using its full OData cast, for example, `$filter=identitySources/any(is:is/microsoft.graph.azureActiveDirectoryTenant/tenantId eq 'bcfdfff4-cbc3-43f2-9000-ba7b7515054f')`. / The subtypes of this type, [azureActiveDirectoryTenant](azureactivedirectorytenant.md), [crossCloudAzureActiveDirectoryTenant](crosscloudazureactivedirectorytenant.md), [domainIdentitySource](domainidentitysource.md), [externalDomainFederation](externaldomainfederation.md), and [socialIdentitySource](socialidentitysource.md) are used in the identity sources of a [connectedOrganization](connectedOrganization.md). / https://learn.microsoft.com/en-us/graph/api/resources/identitysource?view=graph-rest-beta. The _provider_ default value is `[]`.",
		},
		"modified_by": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "UPN of the user who last modified this resource. Read-only.",
		},
		"modified_date_time": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Read-only.",
		},
		"state": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.OneOf("configured", "proposed", "unknownFutureValue"),
			},
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("configured")},
			Computed:            true,
			MarkdownDescription: "The state of a connected organization defines whether assignment policies with requestor scope type `AllConfiguredConnectedOrganizationSubjects` are applicable or not. / Possible values are: `configured`, `proposed`, `unknownFutureValue`. The _provider_ default value is `\"configured\"`.",
		},
		"external_sponsors": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: connectedOrganizationDirectoryObjectAttributes,
			},
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Nullable. / Represents a Microsoft Entra object. The **directoryObject** type is the base type for the following directory entity types generally referred to as directory objects:\n\nThis resource supports: / https://learn.microsoft.com/en-us/graph/api/resources/directoryobject?view=graph-rest-beta. The _provider_ default value is `[]`.",
		},
		"internal_sponsors": schema.SetNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: connectedOrganizationDirectoryObjectAttributes,
			},
			PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "Nullable. / Represents a Microsoft Entra object. The **directoryObject** type is the base type for the following directory entity types generally referred to as directory objects:\n\nThis resource supports: / https://learn.microsoft.com/en-us/graph/api/resources/directoryobject?view=graph-rest-beta. The _provider_ default value is `[]`.",
		},
	},
	MarkdownDescription: "In [Microsoft Entra entitlement management](entitlementmanagement-overview.md), a connected organization is a reference to a directory or domain of another organization whose users can request access. / https://learn.microsoft.com/en-us/graph/api/resources/connectedorganization?view=graph-rest-beta ||| MS Graph: Entitlement management",
}

var connectedOrganizationIdentitySourceValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("azure_active_directory_tenant"),
	path.MatchRelative().AtParent().AtName("cross_cloud_azure_active_directory_tenant"),
	path.MatchRelative().AtParent().AtName("domain_identity_source"),
	path.MatchRelative().AtParent().AtName("external_domain_federation"),
	path.MatchRelative().AtParent().AtName("social_identity_source"),
)

var connectedOrganizationDirectoryObjectAttributes = map[string]schema.Attribute{ // directoryObject
	"id": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The unique identifier for the object. For example, 12345678-9abc-def0-1234-56789abcde. The value of the **id** property is often but not exclusively in the form of a GUID; treat it as an opaque identifier and do not rely on it being a GUID. Key. Not nullable. Read-only.",
	},
}
