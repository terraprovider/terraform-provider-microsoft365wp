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


data "microsoft365wp_identity_governance_custom_task_extensions" "all" {
}

output "microsoft365wp_identity_governance_custom_task_extensions" {
  value = { for x in data.microsoft365wp_identity_governance_custom_task_extensions.all.identity_governance_custom_task_extensions : x.id => x }
}
