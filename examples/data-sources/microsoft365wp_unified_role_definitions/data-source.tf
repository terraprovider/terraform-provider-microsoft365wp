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


data "microsoft365wp_unified_role_definitions" "all" {
}

output "microsoft365wp_unified_role_definitions" {
  value = { for x in data.microsoft365wp_unified_role_definitions.all.unified_role_definitions : x.id => x }
}
