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


data "microsoft365wp_named_location" "one" {
  id = "01234567-89ab-cdef-0123-456789abcdef"
}

output "microsoft365wp_named_location" {
  value = data.microsoft365wp_named_location.one
}
