package services

import (
	"context"
	"terraform-provider-microsoft365wp/workplace/generic"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpdefaultvalue"
	"terraform-provider-microsoft365wp/workplace/wpschema/wpplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	IdentityGovernanceCustomTaskExtensionResource = generic.GenericResource{
		TypeNameSuffix: "identity_governance_custom_task_extension",
		SpecificSchema: identityGovernanceCustomTaskExtensionResourceSchema,
		AccessParams: generic.AccessParams{
			BaseUri:                    "/identityGovernance/lifecycleWorkflows/customTaskExtensions",
			TerraformToGraphMiddleware: identityGovernanceCustomTaskExtensionTerraformToGraphMiddleware,
		},
	}

	IdentityGovernanceCustomTaskExtensionSingularDataSource = generic.CreateGenericDataSourceSingularFromResource(
		&IdentityGovernanceCustomTaskExtensionResource)

	IdentityGovernanceCustomTaskExtensionPluralDataSource = generic.CreateGenericDataSourcePluralFromResource(
		&IdentityGovernanceCustomTaskExtensionResource, "")
)

func identityGovernanceCustomTaskExtensionTerraformToGraphMiddleware(ctx context.Context, diags *diag.Diagnostics, params *generic.TerraformToGraphMiddlewareParams) generic.TerraformToGraphMiddlewareReturns {
	if params.IsUpdate {
		// do not send url on PATCH if not defined in TF config
		if epc, ok := params.RawVal["endpointConfiguration"].(map[string]any); ok && epc != nil {
			if epcType, ok := epc["@odata.type"].(string); ok && epcType == "#microsoft.graph.logicAppTriggerEndpointConfiguration" {
				var urlRef *string
				diags.Append(params.Config.GetAttribute(ctx, path.Root("endpoint_configuration").AtName("logic_app_trigger").AtName("url"), &urlRef)...)
				if diags.HasError() {
					return nil
				}
				if urlRef == nil {
					delete(epc, "url")
				}
			}
		}
	}
	return nil
}

var identityGovernanceCustomTaskExtensionResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{ // identityGovernance.customTaskExtension
		"id": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "<br><br>Supports `$filter`(`eq`, `ne`) and `$orderby`.",
		},
		"authentication_configuration": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // customExtensionAuthenticationConfiguration
				"azure_ad_pop_token": generic.OdataDerivedTypeNestedAttributeRs{
					DerivedType: "#microsoft.graph.azureAdPopTokenAuthentication",
					SingleNestedAttribute: schema.SingleNestedAttribute{
						Optional:   true,
						Attributes: map[string]schema.Attribute{ // azureAdPopTokenAuthentication
						},
						Validators: []validator.Object{
							identityGovernanceCustomTaskExtensionCustomExtensionAuthenticationConfigurationValidator,
						},
						MarkdownDescription: "Defines the Proof Of Possession (PoP) token authentication model to authenticate a logic app with a [accessPackageAssignmentRequestWorkflowExtensions](../resources/accessPackageAssignmentRequestWorkflowExtension.md) or a [accessPackageAssignmentWorkflowExtensions](../resources/accessPackageAssignmentWorkflowExtension.md) object. / https://learn.microsoft.com/en-us/graph/api/resources/azureadpoptokenauthentication?view=graph-rest-beta",
					},
				},
				"azure_ad_token": generic.OdataDerivedTypeNestedAttributeRs{
					DerivedType: "#microsoft.graph.azureAdTokenAuthentication",
					SingleNestedAttribute: schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // azureAdTokenAuthentication
							"resource_id": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "The **appID** of the Microsoft Entra application to use to authenticate a logic app with a custom access package workflow extension.",
							},
						},
						Validators: []validator.Object{
							identityGovernanceCustomTaskExtensionCustomExtensionAuthenticationConfigurationValidator,
						},
						MarkdownDescription: "Defines the Microsoft Entra application used to authenticate a logic app with a [custom access package workflow extension](../resources/customaccesspackageworkflowextension.md) or a [custom task extension](../resources/identitygovernance-customtaskextension.md). Only the app ID of the application is required. Derived from [customExtensionAuthenticationConfiguration](../resources/customextensionauthenticationconfiguration.md). / https://learn.microsoft.com/en-us/graph/api/resources/azureadtokenauthentication?view=graph-rest-beta",
					},
				},
			},
			PlanModifiers: []planmodifier.Object{
				wpdefaultvalue.ObjectDefaultValue(map[string]any{
					"azure_ad_pop_token": map[string]any{},
				}),
			},
			Computed:            true,
			MarkdownDescription: "Configuration for securing the API call to the logic app. Required. / Abstract base type that exposes the configuration for the **authenticationConfiguration** property of the derived types that inherit from the [customCalloutExtension](customcalloutextension.md) abstract type.\n\nThis abstract type is inherited by the following resource types:\n\nThe type of token authentication used depends on the token security. If the token security value is normal, you use the [azureAdTokenAuthentication](../resources/azureadtokenauthentication.md) resource type. If the value is Proof of Possession, you use the [azureAdPopTokenAuthentication](../resources/azureAdPopTokenAuthentication.md) resource type. / https://learn.microsoft.com/en-us/graph/api/resources/customextensionauthenticationconfiguration?view=graph-rest-beta. The _provider_ default value is `{\"azure_ad_pop_token\":{}}`.",
		},
		"client_configuration": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // customExtensionClientConfiguration
				"maximum_retries": schema.Int64Attribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(1)},
					Computed:            true,
					MarkdownDescription: "The max number of retries that Microsoft Entra ID makes to the external API. Values of 0 or 1 are supported. If `null`, the default for the service applies. The _provider_ default value is `1`.",
				},
				"timeout_in_milliseconds": schema.Int64Attribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.Int64{wpdefaultvalue.Int64DefaultValue(1000)},
					Computed:            true,
					MarkdownDescription: "The max duration in milliseconds that Microsoft Entra ID waits for a response from the external app before it shuts down the connection. The valid range is between `200` and `2000` milliseconds. If `null`, the default for the service applies. The _provider_ default value is `1000`.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpdefaultvalue.ObjectDefaultValueEmpty()},
			Computed:            true,
			MarkdownDescription: "HTTP connection settings that define how long Microsoft Entra ID can wait for a connection to a logic app, how many times you can retry a timed-out connection and the exception scenarios when retries are allowed. / Connection settings that define how long Microsoft Entra ID can wait for a response from an external app before it shuts down the connection when trying to trigger the external app. / https://learn.microsoft.com/en-us/graph/api/resources/customextensionclientconfiguration?view=graph-rest-beta. The _provider_ default value is `{}`.",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("")},
			Computed:            true,
			MarkdownDescription: "Describes the purpose of the custom task extension for administrative use. Optional. The _provider_ default value is `\"\"`.",
		},
		"display_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "A unique string that identifies the custom task extension. Required.<br><br>Supports `$filter`(`eq`, `ne`) and `$orderby`.",
		},
		"endpoint_configuration": schema.SingleNestedAttribute{
			Required: true,
			Attributes: map[string]schema.Attribute{ // customExtensionEndpointConfiguration
				"http_request": generic.OdataDerivedTypeNestedAttributeRs{
					DerivedType: "#microsoft.graph.httpRequestEndpoint",
					SingleNestedAttribute: schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // httpRequestEndpoint
							"target_url": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "The HTTP endpoint that a custom extension calls.",
							},
						},
						Validators: []validator.Object{
							identityGovernanceCustomTaskExtensionCustomExtensionEndpointConfigurationValidator,
						},
						MarkdownDescription: "The HTTP endpoint that a custom extension calls. / https://learn.microsoft.com/en-us/graph/api/resources/httprequestendpoint?view=graph-rest-beta",
					},
				},
				"logic_app_trigger": generic.OdataDerivedTypeNestedAttributeRs{
					DerivedType: "#microsoft.graph.logicAppTriggerEndpointConfiguration",
					SingleNestedAttribute: schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // logicAppTriggerEndpointConfiguration
							"logic_app_workflow_name": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "The name of the logic app.",
							},
							"resource_group_name": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "The Azure resource group name for the logic app.",
							},
							"subscription_id": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Identifier of the Azure subscription for the logic app.",
							},
							"url": schema.StringAttribute{
								Computed:            true,
								Optional:            true,
								PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
								MarkdownDescription: "The URL to the logic app endpoint that will be triggered. Only required for app-only token scenarios where app is creating a [customCalloutExtension](../resources/customcalloutextension.md) without a signed-in user.",
							},
						},
						Validators: []validator.Object{
							identityGovernanceCustomTaskExtensionCustomExtensionEndpointConfigurationValidator,
						},
						MarkdownDescription: "The configuration details for the logic app's endpoint that is associated with a custom access package workflow extension. Derived from the [customExtensionEndpointConfiguration](customextensionendpointconfiguration.md) abstract type. / https://learn.microsoft.com/en-us/graph/api/resources/logicapptriggerendpointconfiguration?view=graph-rest-beta",
					},
				},
			},
			MarkdownDescription: "Details for allowing the custom task extension to call the logic app. / Abstract base type that exposes the derived types used to configure the **endpointConfiguration** property of a custom extension. This abstract type is inherited by the following types: / https://learn.microsoft.com/en-us/graph/api/resources/customextensionendpointconfiguration?view=graph-rest-beta",
		},
		"callback_configuration": schema.SingleNestedAttribute{
			Optional: true,
			Attributes: map[string]schema.Attribute{ // customExtensionCallbackConfiguration
				"timeout_duration": schema.StringAttribute{
					Optional:            true,
					PlanModifiers:       []planmodifier.String{wpdefaultvalue.StringDefaultValue("PT30M")},
					Computed:            true,
					MarkdownDescription: "The maximum duration in ISO 8601 format that Microsoft Entra ID will wait for a resume action for the callout it sent to the logic app. The valid range for custom extensions in lifecycle workflows is five minutes to three hours. The valid range for custom extensions in entitlement management is between 5 minutes and 14 days. For example, `PT3H` refers to three hours, `P3D` refers to three days, `PT10M` refers to ten minutes. The _provider_ default value is `\"PT30M\"`.",
				},
				"task": generic.OdataDerivedTypeNestedAttributeRs{
					DerivedType: "#microsoft.graph.identityGovernance.customTaskExtensionCallbackConfiguration",
					SingleNestedAttribute: schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{ // identityGovernance.customTaskExtensionCallbackConfiguration
							"authorized_apps": schema.SetNestedAttribute{
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{ // application
										"id": schema.StringAttribute{
											Required:            true,
											MarkdownDescription: "Unique identifier for the application object. This property is referred to as **Object ID** in the Microsoft Entra admin center. Key. Not nullable. Read-only. Supports `$filter` (`eq`, `ne`, `not`, `in`).",
										},
									},
								},
								PlanModifiers:       []planmodifier.Set{wpdefaultvalue.SetDefaultValueEmpty()},
								Computed:            true,
								MarkdownDescription: "A collection of unique identifiers or **appIds** of the applications that are allowed to [resume](../api/identitygovernance-taskprocessingresult-resume.md) a task processing result. / Represents an application. Any application that outsources authentication to Microsoft Entra ID must be registered in the Microsoft identity platform. Application registration involves telling Microsoft Entra ID about your application, including the URL where it's located, the URL to send replies after authentication, the URI to identify your application, and more.\n\nThis resource is an open type that allows other properties to be passed in.\n\nThis resource supports: / https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta. The _provider_ default value is `[]`.",
							},
						},
						Validators: []validator.Object{
							identityGovernanceCustomTaskExtensionCustomExtensionCallbackConfigurationValidator,
						},
						MarkdownDescription: "Defines if, and in, which time span a callback is expected from the Azure Logic App.\n\nInherits from  [customExtensionCallbackConfiguration](../resources/customextensioncallbackconfiguration.md). / https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-customtaskextensioncallbackconfiguration?view=graph-rest-beta",
					},
				},
			},
			MarkdownDescription: "The callback configuration for a custom task extension. / Callback settings that define how long Microsoft Entra ID can wait for a resume signal for the callout that it made to the logic app. This is an abstract type that's inherited by [customTaskExtensionCallbackConfiguration](../resources/identitygovernance-customtaskextensioncallbackconfiguration.md). / https://learn.microsoft.com/en-us/graph/api/resources/customextensioncallbackconfiguration?view=graph-rest-beta",
		},
		"created_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "When the custom task extension was created.<br><br>Supports `$filter`(`lt`, `le`, `gt`, `ge`, `eq`, `ne`) and `$orderby`.",
		},
		"last_modified_date_time": schema.StringAttribute{
			Computed:            true,
			PlanModifiers:       []planmodifier.String{wpplanmodifier.StringUseStateForUnknown()},
			MarkdownDescription: "When the custom extension was last modified.<br><br>Supports `$filter`(`lt`, `le`, `gt`, `ge`, `eq`, `ne`) and `$orderby`.",
		},
		"created_by": schema.SingleNestedAttribute{
			Computed: true,
			Attributes: map[string]schema.Attribute{ // user
				"id": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The user identifier.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpplanmodifier.ObjectUseStateForUnknown()},
			MarkdownDescription: "The unique identifier of the Microsoft Entra user that created the custom task extension.<br><br>Supports `$filter`(`eq`, `ne`) and `$expand`. / Represents an Azure Active Directory user object. / https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-user?view=graph-rest-beta",
		},
		"last_modified_by": schema.SingleNestedAttribute{
			Computed: true,
			Attributes: map[string]schema.Attribute{ // user
				"id": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The user identifier.",
				},
			},
			PlanModifiers:       []planmodifier.Object{wpplanmodifier.ObjectUseStateForUnknown()},
			MarkdownDescription: "The unique identifier of the Microsoft Entra user that modified the custom task extension last.<br><br>Supports `$filter`(`eq`, `ne`) and `$expand`. / Represents an Azure Active Directory user object. / https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-user?view=graph-rest-beta",
		},
	},
	MarkdownDescription: "Defines the attributes of a customTaskExtension that allows you to integrate Lifecycle Workflows with Azure Logic Apps. While Lifecycle Workflows provide multiple built-in tasks (known as taskDefinitions) to automate common scenarios during the user lifecycle, you may eventually reach the limits of these built-in tasks. You can create a customTaskExtension that contains information about an Azure Logic app, and trigger the Azure Logic app with the built-in task \"Run a custom task extension\" that references the corresponding customTaskExtension.\n\nFor more information about using custom task extensions, refer to the links in the [see also](#related-content) section. / https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-customtaskextension?view=graph-rest-beta ||| MS Graph: Lifecycle workflows",
}

var identityGovernanceCustomTaskExtensionCustomExtensionAuthenticationConfigurationValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("azure_ad_pop_token"),
	path.MatchRelative().AtParent().AtName("azure_ad_token"),
)

var identityGovernanceCustomTaskExtensionCustomExtensionEndpointConfigurationValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("http_request"),
	path.MatchRelative().AtParent().AtName("logic_app_trigger"),
)

var identityGovernanceCustomTaskExtensionCustomExtensionCallbackConfigurationValidator = objectvalidator.ExactlyOneOf(
	path.MatchRelative().AtParent().AtName("task"),
)
