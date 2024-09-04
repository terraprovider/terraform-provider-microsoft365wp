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


data "microsoft365wp_authentication_strength_policy" "one" {
  id = "ccde5395-3d28-45f2-88dd-8ea78d538016"
}

output "microsoft365wp_authentication_strength_policy" {
  value = data.microsoft365wp_authentication_strength_policy.one
}
