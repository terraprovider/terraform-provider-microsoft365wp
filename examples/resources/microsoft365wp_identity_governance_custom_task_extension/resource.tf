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
