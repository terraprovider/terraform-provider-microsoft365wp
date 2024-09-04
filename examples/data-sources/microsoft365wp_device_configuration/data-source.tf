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


data "microsoft365wp_device_configuration" "one" {
  id = "8cefe928-c91b-40bb-b5d9-64dc5919dc37"
}

output "microsoft365wp_device_configuration" {
  value = data.microsoft365wp_device_configuration.one
}
