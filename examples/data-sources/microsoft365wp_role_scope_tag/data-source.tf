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


data "microsoft365wp_role_scope_tag" "one" {
  id = "0" # "Default" (built-in)
}

output "microsoft365wp_role_scope_tag" {
  value = data.microsoft365wp_role_scope_tag.one
}
