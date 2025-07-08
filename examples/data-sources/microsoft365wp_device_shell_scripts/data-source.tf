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


data "microsoft365wp_device_shell_scripts" "all" {
}

output "microsoft365wp_device_shell_scripts" {
  value = { for x in data.microsoft365wp_device_shell_scripts.all.device_shell_scripts : x.id => x }
}
