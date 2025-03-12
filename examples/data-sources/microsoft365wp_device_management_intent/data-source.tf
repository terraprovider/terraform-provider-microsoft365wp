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


data "microsoft365wp_device_management_intent" "one" {
  id = "9dacf1e4-5297-442d-8a7d-a11ab847becf"
}

output "microsoft365wp_device_management_intent" {
  value = data.microsoft365wp_device_management_intent.one
}
