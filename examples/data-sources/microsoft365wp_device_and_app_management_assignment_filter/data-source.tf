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


data "microsoft365wp_device_and_app_management_assignment_filter" "one" {
  id = "8808bb33-aa9c-4da7-a3fb-b92a4be8e001"
}

output "microsoft365wp_device_and_app_management_assignment_filter" {
  value = data.microsoft365wp_device_and_app_management_assignment_filter.one
}
