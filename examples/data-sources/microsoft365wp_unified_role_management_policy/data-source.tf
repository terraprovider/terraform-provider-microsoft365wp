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


data "microsoft365wp_unified_role_management_policy" "one" {
  id = "DirectoryRole_810017af-c1d4-4a3b-9e18-d693c12b0df1_7da17cbf-6f42-4103-b73d-e500dca5f92d"
}

output "microsoft365wp_unified_role_management_policy" {
  value = data.microsoft365wp_unified_role_management_policy.one
}
