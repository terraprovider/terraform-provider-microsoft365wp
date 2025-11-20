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
    merge(r, r.id == "Approval_EndUser_Assignment" ? {
      approval = {
        setting = {
          is_approval_required                = true
          is_approval_required_for_extension  = false
          is_requestor_justification_required = true
          approval_mode                       = "SingleStage"
          approval_stages = [
            {
              approval_stage_time_out_in_days    = 1
              escalation_approvers               = []
              escalation_time_in_minutes         = 0
              is_approver_justification_required = true
              is_escalation_enabled              = false
              primary_approvers = [
                {
                  group_members = {
                    id = "298fded6-b252-4166-a473-f405e935f58d"
                  }
                },
                {
                  group_members = {
                    id = "62e39046-aad3-4423-98e0-b486e3538aff"
                  }
                }
              ]
            }
          ]
        }
      }
    } : {})
  ]
}

resource "microsoft365wp_unified_role_management_policy" "test" {
  # This will never try to create but always just update an existing policy determined by the id attribute.
  id = local.urm_policy.id

  is_organization_default = local.urm_policy.is_organization_default
  rules                   = local.urm_policy_rules
}
