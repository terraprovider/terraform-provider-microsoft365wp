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


data "microsoft365wp_custom_security_attribute_definitions" "all" {
  status = "Available"
}

output "microsoft365wp_custom_security_attribute_definitions" {
  value = { for x in data.microsoft365wp_custom_security_attribute_definitions.all.custom_security_attribute_definitions : x.id => x }
}
