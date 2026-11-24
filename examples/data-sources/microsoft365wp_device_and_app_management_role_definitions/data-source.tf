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


data "microsoft365wp_device_and_app_management_role_definitions" "all" {
  exclude_ids = ["01234567-89ab-cdef-0123-456789abcdee", "01234567-89ab-cdef-0123-456789abcdef"]
}

output "microsoft365wp_device_and_app_management_role_definitions" {
  value = { for x in data.microsoft365wp_device_and_app_management_role_definitions.all.device_and_app_management_role_definitions : x.id => x }
}
