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


data "microsoft365wp_role_scope_tags" "all" {
  exclude_ids = ["0"] # skip built-in "Default"
}

output "microsoft365wp_role_scope_tags" {
  value = { for x in data.microsoft365wp_role_scope_tags.all.role_scope_tags : x.id => x }
}
