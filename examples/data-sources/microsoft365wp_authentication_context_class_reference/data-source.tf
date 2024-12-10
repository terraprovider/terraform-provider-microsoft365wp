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


data "microsoft365wp_authentication_context_class_reference" "one" {
  id = "c2"
}

output "microsoft365wp_authentication_context_class_reference" {
  value = data.microsoft365wp_authentication_context_class_reference.one
}
