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


data "microsoft365wp_device_shell_script" "one" {
  id = "4f6034db-59e0-49a2-8e0a-fe5ac1364b92"
}

output "microsoft365wp_device_shell_script" {
  value = data.microsoft365wp_device_shell_script.one
}
