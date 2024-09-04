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


data "microsoft365wp_device_configuration_custom" "one" {
  id = "1c527e3e-c8b3-4b6c-8ce0-671feb22f171"
}

output "microsoft365wp_device_configuration_custom" {
  value = data.microsoft365wp_device_configuration_custom.one
}
