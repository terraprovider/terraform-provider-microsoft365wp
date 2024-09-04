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


data "microsoft365wp_device_enrollment_configurations" "all" {
}

output "microsoft365wp_device_enrollment_configurations" {
  value = { for x in data.microsoft365wp_device_enrollment_configurations.all.device_enrollment_configurations : x.id => x }
}
