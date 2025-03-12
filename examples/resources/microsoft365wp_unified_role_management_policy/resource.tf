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


data "microsoft365wp_unified_role_definition" "main" {
  display_name = "Application Administrator"
}

data "microsoft365wp_unified_role_management_policy_assignment" "main" {
  odata_filter       = "scopeType eq 'DirectoryRole' and scopeId eq '/'"
  role_definition_id = data.microsoft365wp_unified_role_definition.main.id
}

locals {
  urm_policy = data.microsoft365wp_unified_role_management_policy_assignment.main.policy
  urm_policy_rules = [for r in local.urm_policy.rules :
    r.id == "Enablement_Admin_Assignment" ? merge(r, { enablement = { enabled_rules = ["Justification"] } }) : r
  ]
}

resource "microsoft365wp_unified_role_management_policy" "test" {
  # This will never try to create but always just update an existing policy determined by the id attribute.
  id = local.urm_policy.id

  is_organization_default = local.urm_policy.is_organization_default
  rules                   = local.urm_policy_rules
}
