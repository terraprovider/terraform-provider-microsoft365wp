---
page_title: "microsoft365wp_identity_governance_custom_task_extension Resource - microsoft365wp"
subcategory: "MS Graph: Lifecycle workflows"
---

# microsoft365wp_identity_governance_custom_task_extension (Resource)

Defines the attributes of a customTaskExtension that allows you to integrate Lifecycle Workflows with Azure Logic Apps. While Lifecycle Workflows provide multiple built-in tasks (known as taskDefinitions) to automate common scenarios during the user lifecycle, you may eventually reach the limits of these built-in tasks. You can create a customTaskExtension that contains information about an Azure Logic app, and trigger the Azure Logic app with the built-in task "Run a custom task extension" that references the corresponding customTaskExtension.

For more information about using custom task extensions, refer to the links in the [see also](#related-content) section. / https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-customtaskextension?view=graph-rest-beta

## Documentation Disclaimer

Please note that almost all information on this page has been sourced literally from the official Microsoft Graph API 
documentation and therefore is governed by Microsoft and not by the publishers of this provider.  
All supplements authored by the publishers of this provider have been explicitly marked as such.

## Example Usage

```terraform
terraform {
  required_providers {
    microsoft365wp = {
      source = "terraprovider/microsoft365wp"
    }
  }
}

# azurerm_logic_app_workflow etc. only
provider "azurerm" {
  resource_provider_registrations = "none" # This is only required when the User, Service Principal, or Identity running Terraform lacks the permissions to register Azure Resource Providers.
  features {}
}

/*
.env
export ARM_TENANT_ID='...'
export ARM_SUBSCRIPTION_ID='...'
export ARM_CLIENT_ID='...'
export ARM_CLIENT_SECRET='...'
*/


locals {
  resource_group_name     = "rg-tf-provider-helpers"
  logic_app_workflow_name = "logic-tf-provider-lifecycle-workflow-custom-extension-wait"
}

data "azurerm_client_config" "current" {}

data "azurerm_resource_group" "main" {
  name = local.resource_group_name
}

resource "azurerm_logic_app_workflow" "main" {
  resource_group_name = data.azurerm_resource_group.main.name
  name                = local.logic_app_workflow_name
  location            = data.azurerm_resource_group.main.location

  tags = { Purpose = "Azure AD Lifecycle Workflows" }

  identity { type = "SystemAssigned" }

  access_control {
    trigger {
      allowed_caller_ip_address_range = []
      open_authentication_policy {
        name = "AzureADLifecycleWorkflowsAuthPOPAuthPolicy"
        claim {
          name  = "appid"
          value = "ce79fdc4-cd1d-4ea5-8139-e74d7dbe0bb7" # hard coded
        }
        claim {
          name  = "iss"
          value = "https://sts.windows.net/${data.azurerm_client_config.current.tenant_id}/"
        }
        claim {
          name  = "m"
          value = "POST"
        }
        claim {
          name  = "p"
          value = "${data.azurerm_resource_group.main.id}/providers/Microsoft.Logic/workflows/${local.logic_app_workflow_name}"
        }
        claim {
          name  = "u"
          value = "management.azure.com"
        }
      }
    }
  }

}

resource "azurerm_logic_app_trigger_http_request" "main" {
  name         = "manual"
  logic_app_id = azurerm_logic_app_workflow.main.id

  schema = jsonencode(
    {
      properties = {
        data = {
          properties = {
            callbackUriPath = {
              description = "CallbackUriPath used for Resume Action"
              title       = "Data.CallbackUriPath"
              type        = "string"
            }
            subject = {
              properties = {
                displayName = {
                  description = "DisplayName of the Subject"
                  title       = "Subject.DisplayName"
                  type        = "string"
                }
                email = {
                  description = "Email of the Subject"
                  title       = "Subject.Email"
                  type        = "string"
                }
                id = {
                  description = "Id of the Subject"
                  title       = "Subject.Id"
                  type        = "string"
                }
                manager = {
                  properties = {
                    displayName = {
                      description = "DisplayName parameter for Manager"
                      title       = "Manager.DisplayName"
                      type        = "string"
                    }
                    email = {
                      description = "Mail parameter for Manager"
                      title       = "Manager.Mail"
                      type        = "string"
                    }
                    id = {
                      description = "Id parameter for Manager"
                      title       = "Manager.Id"
                      type        = "string"
                    }
                  }
                  type = "object"
                }
                userPrincipalName = {
                  description = "UserPrincipalName of the Subject"
                  title       = "Subject.UserPrincipalName"
                  type        = "string"
                }
              }
              type = "object"
            }
            task = {
              properties = {
                displayName = {
                  description = "DisplayName for Task Object"
                  title       = "Task.DisplayName"
                  type        = "string"
                }
                id = {
                  description = "Id for Task Object"
                  title       = "Task.Id"
                  type        = "string"
                }
              }
              type = "object"
            }
            taskProcessingResult = {
              properties = {
                createdDateTime = {
                  description = "CreatedDateTime for TaskProcessingResult Object"
                  title       = "TaskProcessingResult.CreatedDateTime"
                  type        = "string"
                }
                id = {
                  description = "Id for TaskProcessingResult Object"
                  title       = "TaskProcessingResult.Id"
                  type        = "string"
                }
              }
              type = "object"
            }
            workflow = {
              properties = {
                displayName = {
                  description = "DisplayName for Workflow Object"
                  title       = "Workflow.DisplayName"
                  type        = "string"
                }
                id = {
                  description = "Id for Workflow Object"
                  title       = "Workflow.Id"
                  type        = "string"
                }
                workflowVerson = {
                  description = "WorkflowVersion for Workflow Object"
                  title       = "Workflow.WorkflowVersion"
                  type        = "integer"
                }
              }
              type = "object"
            }
          }
          type = "object"
        }
        source = {
          description = "Context in which an event happened"
          title       = "Request.Source"
          type        = "string"
        }
        type = {
          description = "Value describing the type of event related to the originating occurrence."
          title       = "Request.Type"
          type        = "string"
        }
      }
      type = "object"
  })

}

resource "time_sleep" "spn_delay_after_create" {
  depends_on = [azurerm_logic_app_workflow.main]

  create_duration = "15s"
}

data "azuread_service_principal" "main" {
  object_id = one([for x in azurerm_logic_app_workflow.main.identity : x.principal_id if x.type == "SystemAssigned"])

  depends_on = [time_sleep.spn_delay_after_create]
}

resource "microsoft365wp_identity_governance_custom_task_extension" "test" {
  display_name = "TF Test 1"

  endpoint_configuration = {
    logic_app_trigger = {
      subscription_id         = data.azurerm_client_config.current.subscription_id
      resource_group_name     = azurerm_logic_app_workflow.main.resource_group_name
      logic_app_workflow_name = azurerm_logic_app_workflow.main.name
      # cut off url parameters, enforce (apparently mandatory) api version
      url = replace(azurerm_logic_app_trigger_http_request.main.callback_url, "/\\?.*$/", "?api-version=2016-10-01")
    }
  }

  callback_configuration = {
    task = {
      authorized_apps = [
        { "id" = data.azuread_service_principal.main.client_id },
      ]
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) A unique string that identifies the custom task extension. Required.<br><br>Supports `$filter`(`eq`, `ne`) and `$orderby`.
- `endpoint_configuration` (Attributes) Details for allowing the custom task extension to call the logic app. / Abstract base type that exposes the derived types used to configure the **endpointConfiguration** property of a custom extension. This abstract type is inherited by the following types: / https://learn.microsoft.com/en-us/graph/api/resources/customextensionendpointconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--endpoint_configuration))

### Optional

- `authentication_configuration` (Attributes) Configuration for securing the API call to the logic app. Required. / Abstract base type that exposes the configuration for the **authenticationConfiguration** property of the derived types that inherit from the [customCalloutExtension](customcalloutextension.md) abstract type.

This abstract type is inherited by the following resource types:

The type of token authentication used depends on the token security. If the token security value is normal, you use the [azureAdTokenAuthentication](../resources/azureadtokenauthentication.md) resource type. If the value is Proof of Possession, you use the [azureAdPopTokenAuthentication](../resources/azureAdPopTokenAuthentication.md) resource type. / https://learn.microsoft.com/en-us/graph/api/resources/customextensionauthenticationconfiguration?view=graph-rest-beta. The _provider_ default value is `{"azure_ad_pop_token":{}}`. (see [below for nested schema](#nestedatt--authentication_configuration))
- `callback_configuration` (Attributes) The callback configuration for a custom task extension. / Callback settings that define how long Microsoft Entra ID can wait for a resume signal for the callout that it made to the logic app. This is an abstract type that's inherited by [customTaskExtensionCallbackConfiguration](../resources/identitygovernance-customtaskextensioncallbackconfiguration.md). / https://learn.microsoft.com/en-us/graph/api/resources/customextensioncallbackconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--callback_configuration))
- `client_configuration` (Attributes) HTTP connection settings that define how long Microsoft Entra ID can wait for a connection to a logic app, how many times you can retry a timed-out connection and the exception scenarios when retries are allowed. / Connection settings that define how long Microsoft Entra ID can wait for a response from an external app before it shuts down the connection when trying to trigger the external app. / https://learn.microsoft.com/en-us/graph/api/resources/customextensionclientconfiguration?view=graph-rest-beta. The _provider_ default value is `{}`. (see [below for nested schema](#nestedatt--client_configuration))
- `description` (String) Describes the purpose of the custom task extension for administrative use. Optional. The _provider_ default value is `""`.

### Read-Only

- `created_by` (Attributes) The unique identifier of the Microsoft Entra user that created the custom task extension.<br><br>Supports `$filter`(`eq`, `ne`) and `$expand`. / Represents an Azure Active Directory user object. / https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-user?view=graph-rest-beta (see [below for nested schema](#nestedatt--created_by))
- `created_date_time` (String) When the custom task extension was created.<br><br>Supports `$filter`(`lt`, `le`, `gt`, `ge`, `eq`, `ne`) and `$orderby`.
- `id` (String) <br><br>Supports `$filter`(`eq`, `ne`) and `$orderby`.
- `last_modified_by` (Attributes) The unique identifier of the Microsoft Entra user that modified the custom task extension last.<br><br>Supports `$filter`(`eq`, `ne`) and `$expand`. / Represents an Azure Active Directory user object. / https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-user?view=graph-rest-beta (see [below for nested schema](#nestedatt--last_modified_by))
- `last_modified_date_time` (String) When the custom extension was last modified.<br><br>Supports `$filter`(`lt`, `le`, `gt`, `ge`, `eq`, `ne`) and `$orderby`.

<a id="nestedatt--endpoint_configuration"></a>
### Nested Schema for `endpoint_configuration`

Optional:

- `http_request` (Attributes) The HTTP endpoint that a custom extension calls. / https://learn.microsoft.com/en-us/graph/api/resources/httprequestendpoint?view=graph-rest-beta (see [below for nested schema](#nestedatt--endpoint_configuration--http_request))
- `logic_app_trigger` (Attributes) The configuration details for the logic app's endpoint that is associated with a custom access package workflow extension. Derived from the [customExtensionEndpointConfiguration](customextensionendpointconfiguration.md) abstract type. / https://learn.microsoft.com/en-us/graph/api/resources/logicapptriggerendpointconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--endpoint_configuration--logic_app_trigger))

<a id="nestedatt--endpoint_configuration--http_request"></a>
### Nested Schema for `endpoint_configuration.http_request`

Required:

- `target_url` (String) The HTTP endpoint that a custom extension calls.


<a id="nestedatt--endpoint_configuration--logic_app_trigger"></a>
### Nested Schema for `endpoint_configuration.logic_app_trigger`

Required:

- `logic_app_workflow_name` (String) The name of the logic app.
- `resource_group_name` (String) The Azure resource group name for the logic app.
- `subscription_id` (String) Identifier of the Azure subscription for the logic app.

Optional:

- `url` (String) The URL to the logic app endpoint that will be triggered. Only required for app-only token scenarios where app is creating a [customCalloutExtension](../resources/customcalloutextension.md) without a signed-in user.



<a id="nestedatt--authentication_configuration"></a>
### Nested Schema for `authentication_configuration`

Optional:

- `azure_ad_pop_token` (Attributes) Defines the Proof Of Possession (PoP) token authentication model to authenticate a logic app with a [accessPackageAssignmentRequestWorkflowExtensions](../resources/accessPackageAssignmentRequestWorkflowExtension.md) or a [accessPackageAssignmentWorkflowExtensions](../resources/accessPackageAssignmentWorkflowExtension.md) object. / https://learn.microsoft.com/en-us/graph/api/resources/azureadpoptokenauthentication?view=graph-rest-beta (see [below for nested schema](#nestedatt--authentication_configuration--azure_ad_pop_token))
- `azure_ad_token` (Attributes) Defines the Microsoft Entra application used to authenticate a logic app with a [custom access package workflow extension](../resources/customaccesspackageworkflowextension.md) or a [custom task extension](../resources/identitygovernance-customtaskextension.md). Only the app ID of the application is required. Derived from [customExtensionAuthenticationConfiguration](../resources/customextensionauthenticationconfiguration.md). / https://learn.microsoft.com/en-us/graph/api/resources/azureadtokenauthentication?view=graph-rest-beta (see [below for nested schema](#nestedatt--authentication_configuration--azure_ad_token))

<a id="nestedatt--authentication_configuration--azure_ad_pop_token"></a>
### Nested Schema for `authentication_configuration.azure_ad_pop_token`


<a id="nestedatt--authentication_configuration--azure_ad_token"></a>
### Nested Schema for `authentication_configuration.azure_ad_token`

Required:

- `resource_id` (String) The **appID** of the Microsoft Entra application to use to authenticate a logic app with a custom access package workflow extension.



<a id="nestedatt--callback_configuration"></a>
### Nested Schema for `callback_configuration`

Optional:

- `task` (Attributes) Defines if, and in, which time span a callback is expected from the Azure Logic App.

Inherits from  [customExtensionCallbackConfiguration](../resources/customextensioncallbackconfiguration.md). / https://learn.microsoft.com/en-us/graph/api/resources/identitygovernance-customtaskextensioncallbackconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--callback_configuration--task))
- `timeout_duration` (String) The maximum duration in ISO 8601 format that Microsoft Entra ID will wait for a resume action for the callout it sent to the logic app. The valid range for custom extensions in lifecycle workflows is five minutes to three hours. The valid range for custom extensions in entitlement management is between 5 minutes and 14 days. For example, `PT3H` refers to three hours, `P3D` refers to three days, `PT10M` refers to ten minutes. The _provider_ default value is `"PT30M"`.

<a id="nestedatt--callback_configuration--task"></a>
### Nested Schema for `callback_configuration.task`

Optional:

- `authorized_apps` (Attributes Set) A collection of unique identifiers or **appIds** of the applications that are allowed to [resume](../api/identitygovernance-taskprocessingresult-resume.md) a task processing result. / Represents an application. Any application that outsources authentication to Microsoft Entra ID must be registered in the Microsoft identity platform. Application registration involves telling Microsoft Entra ID about your application, including the URL where it's located, the URL to send replies after authentication, the URI to identify your application, and more.

This resource is an open type that allows other properties to be passed in.

This resource supports: / https://learn.microsoft.com/en-us/graph/api/resources/application?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--callback_configuration--task--authorized_apps))

<a id="nestedatt--callback_configuration--task--authorized_apps"></a>
### Nested Schema for `callback_configuration.task.authorized_apps`

Required:

- `id` (String) Unique identifier for the application object. This property is referred to as **Object ID** in the Microsoft Entra admin center. Key. Not nullable. Read-only. Supports `$filter` (`eq`, `ne`, `not`, `in`).




<a id="nestedatt--client_configuration"></a>
### Nested Schema for `client_configuration`

Optional:

- `maximum_retries` (Number) The max number of retries that Microsoft Entra ID makes to the external API. Values of 0 or 1 are supported. If `null`, the default for the service applies. The _provider_ default value is `1`.
- `timeout_in_milliseconds` (Number) The max duration in milliseconds that Microsoft Entra ID waits for a response from the external app before it shuts down the connection. The valid range is between `200` and `2000` milliseconds. If `null`, the default for the service applies. The _provider_ default value is `1000`.


<a id="nestedatt--created_by"></a>
### Nested Schema for `created_by`

Read-Only:

- `id` (String) The user identifier.


<a id="nestedatt--last_modified_by"></a>
### Nested Schema for `last_modified_by`

Read-Only:

- `id` (String) The user identifier.
