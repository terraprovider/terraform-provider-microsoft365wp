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


data "microsoft365wp_unified_role_management_policy_assignment" "_1" {
  id = "DirectoryRole_810017af-c1d4-4a3b-9e18-d693c12b0df1_130dc0c9-bcb8-4d55-a076-5be815a6a1ed_9c6df0f2-1e7c-4dc3-b195-66dfbd24aa8f"
}

output "microsoft365wp_unified_role_management_policy_assignment_1" {
  value = data.microsoft365wp_unified_role_management_policy_assignment._1
}

data "microsoft365wp_unified_role_management_policy_assignment" "_2" {
  odata_filter       = "scopeType eq 'DirectoryRole' and scopeId eq '/'"
  role_definition_id = "9c6df0f2-1e7c-4dc3-b195-66dfbd24aa8f"
}

output "microsoft365wp_unified_role_management_policy_assignment_2" {
  value = data.microsoft365wp_unified_role_management_policy_assignment._2
}

data "microsoft365wp_unified_role_management_policy_assignment" "_3" {
  scope_type = "DirectoryRole"
  scope_id   = "/"
  policy_id  = "DirectoryRole_810017af-c1d4-4a3b-9e18-d693c12b0df1_c34bb862-3e07-4496-92ff-9f31214740fb"
}

output "microsoft365wp_unified_role_management_policy_assignment_3" {
  value = data.microsoft365wp_unified_role_management_policy_assignment._3
}

data "microsoft365wp_unified_role_management_policy_assignment" "_4" {
  scope_type         = "Group"
  scope_id           = "298fded6-b252-4166-a473-f405e935f58d"
  role_definition_id = "owner"
}

output "microsoft365wp_unified_role_management_policy_assignment_4" {
  value = data.microsoft365wp_unified_role_management_policy_assignment._4
}
