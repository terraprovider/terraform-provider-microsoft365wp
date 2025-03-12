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


data "microsoft365wp_unified_role_definition" "_1" {
  odata_filter = "displayName eq 'User Administrator'"
}

output "microsoft365wp_unified_role_definition_1" {
  value = data.microsoft365wp_unified_role_definition._1
}

data "microsoft365wp_unified_role_definition" "_2" {
  display_name = "Global Reader"
}

output "microsoft365wp_unified_role_definition_2" {
  value = data.microsoft365wp_unified_role_definition._2
}
