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
  id = "29e75170-56c9-4653-8c6b-c9cc06df33d0"
}

output "microsoft365wp_device_management_intent" {
  value = data.microsoft365wp_device_management_intent.one
}
