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


data "microsoft365wp_managed_device_mobile_app_configuration" "one" {
  id = "aa4d5339-6c93-4a52-8b3e-25f79805bdc6"
}

output "microsoft365wp_managed_device_mobile_app_configuration" {
  value = data.microsoft365wp_managed_device_mobile_app_configuration.one
}
