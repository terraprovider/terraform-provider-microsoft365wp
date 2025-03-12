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


data "microsoft365wp_identity_governance_custom_task_extension" "one" {
  id = "57cd6b7a-58da-451e-9b2e-76ecf743e990"
}

output "microsoft365wp_identity_governance_custom_task_extension" {
  value = data.microsoft365wp_identity_governance_custom_task_extension.one
}
