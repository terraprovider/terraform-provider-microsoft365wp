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


data "microsoft365wp_device_and_app_management_assignment_filters" "all" {
}

output "microsoft365wp_device_and_app_management_assignment_filters" {
  value = { for x in data.microsoft365wp_device_and_app_management_assignment_filters.all.device_and_app_management_assignment_filters : x.id => x }
}
