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


data "microsoft365wp_device_custom_attribute_shell_script" "one" {
  id = "fa3d00f2-2830-46f0-9bb4-4a0979ee7f46"
}

output "microsoft365wp_device_custom_attribute_shell_script" {
  value = data.microsoft365wp_device_custom_attribute_shell_script.one
}
