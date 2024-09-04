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


data "microsoft365wp_device_management_configuration_policy_json" "one" {
  id = "6c8633bd-7dbb-4b69-ba63-35b6cbe795e4"
}

output "microsoft365wp_device_management_configuration_policy_json" {
  value = data.microsoft365wp_device_management_configuration_policy_json.one
}
