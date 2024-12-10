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


data "microsoft365wp_windows_driver_update_profile" "one" {
  id = "f8ae8ad2-9612-4b7c-818b-5c4a15e1b59e"
}

output "microsoft365wp_windows_driver_update_profile" {
  value = data.microsoft365wp_windows_driver_update_profile.one
}
