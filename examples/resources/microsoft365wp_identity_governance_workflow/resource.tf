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


data "microsoft365wp_identity_governance_task_definition" "add_user_to_groups" {
  display_name = "Add user to groups"
}
data "microsoft365wp_identity_governance_task_definition" "add_user_to_teams" {
  display_name = "Add user to Teams"
}

resource "microsoft365wp_identity_governance_workflow" "test" {
  display_name = "TF Test"

  category = "joiner"

  execution_conditions = {
    trigger_and_scope_based_conditions = {
      trigger = {
        time_based_attribute = {
          time_based_attribute = "employeeHireDate"
          offset_in_days       = 7
        }
      }
      scope = {
        rule = { rule = "(department eq 'Marketing')" }
      }
    }
  }

  tasks = [
    for x in [
      merge(data.microsoft365wp_identity_governance_task_definition.add_user_to_groups, {
        is_enabled = true
        arguments = [{
          name  = "groupID"
          value = "298fded6-b252-4166-a473-f405e935f58d"
        }]
      }),
      merge(data.microsoft365wp_identity_governance_task_definition.add_user_to_teams, {
        is_enabled = false
        arguments = [{
          name  = "teamID"
          value = "5bdad799-d399-4241-b8ca-2a2ca837f443"
        }]
      }),
    ] :
    {
      task_definition_id = x.id
      category           = x.category
      display_name       = x.display_name
      arguments          = x.arguments
      is_enabled         = x.is_enabled
    }
  ]
}
