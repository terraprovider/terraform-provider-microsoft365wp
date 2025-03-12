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


data "microsoft365wp_device_management_script" "one" {
  id = "345389ce-f037-4566-8950-45c6f603652d"
}

output "microsoft365wp_device_management_script" {
  value = data.microsoft365wp_device_management_script.one
}
