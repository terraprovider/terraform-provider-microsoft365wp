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


data "microsoft365wp_unified_role_management_policy_assignments" "directory_root" {
  odata_filter = "scopeType eq 'DirectoryRole' and scopeId eq '/'"
}

output "microsoft365wp_unified_role_management_policy_assignments_directory_root" {
  value = { for x in data.microsoft365wp_unified_role_management_policy_assignments.directory_root.unified_role_management_policy_assignments : x.id => x }
}


data "microsoft365wp_unified_role_management_policy_assignments" "group" {
  odata_filter = "scopeType eq 'Group' and scopeId eq '298fded6-b252-4166-a473-f405e935f58d'"
}

output "microsoft365wp_unified_role_management_policy_assignments_group" {
  value = { for x in data.microsoft365wp_unified_role_management_policy_assignments.group.unified_role_management_policy_assignments : x.id => x }
}
