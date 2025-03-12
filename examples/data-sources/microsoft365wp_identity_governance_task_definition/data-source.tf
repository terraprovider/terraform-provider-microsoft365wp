terraform {
  required_providers {
    microsoft365wp = {
      source = "terraprovider/microsoft365wp"
    }
  }
}

/*
.env
export ARM_TENANT_ID='...'
export ARM_CLIENT_ID='...'
export ARM_CLIENT_SECRET='...'
*/


data "microsoft365wp_identity_governance_task_definition" "one" {
  id = "c1ec1e76-f374-4375-aaa6-0bb6bd4c60be"
}

output "microsoft365wp_identity_governance_task_definition_1" {
  value = data.microsoft365wp_identity_governance_task_definition.one
}

data "microsoft365wp_identity_governance_task_definition" "two" {
  display_name = "Send email to notify manager of user move"
}

output "microsoft365wp_identity_governance_task_definition_2" {
  value = data.microsoft365wp_identity_governance_task_definition.two
}
