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


resource "microsoft365wp_conditional_access_policy" "test_grant" {
  display_name = "TF Test Grant"
  conditions = {
    applications = {
      include_applications = ["None"]
    }
    users = {
      include_users = ["None"]
    }
  }
  grant_controls = {
    operator                = "AND"
    authentication_strength = { id = "00000000-0000-0000-0000-000000000004" }
    built_in_controls       = ["compliantDevice"]
  }
}

resource "microsoft365wp_conditional_access_policy" "test_session" {
  display_name = "TF Test Session"
  conditions = {
    applications = {
      include_applications = ["None"]
    }
    users = {
      include_users = ["298fded6-b252-4166-a473-f405e935f58d"]
      include_guests_or_external_users = {
        external_tenants = {
          membership_kind = "all"
          all             = {}
        }
        guest_or_external_user_types = "b2bCollaborationGuest,b2bCollaborationMember,b2bDirectConnectUser"
      }
      exclude_guests_or_external_users = {
        external_tenants = {
          membership_kind = "enumerated"
          enumerated = {
            members = ["62e39046-aad3-4423-98e0-b486e3538aff"]
          }
        }
        guest_or_external_user_types = "internalGuest,b2bCollaborationGuest,b2bCollaborationMember,b2bDirectConnectUser,otherExternalUser,serviceProvider"
      }
    }
  }
  session_controls = {
    sign_in_frequency = {
      authentication_type = "primaryAndSecondaryAuthentication"
      frequency_interval  = "timeBased"
      value               = 10
      type                = "hours"
      is_enabled          = true
    }
    disable_resilience_defaults = true
  }
}

resource "microsoft365wp_conditional_access_policy" "test_combined" {
  display_name = "TF Test Combined"
  conditions = {
    applications = {
      include_applications = ["All"]
    }
    users = {
      include_users = ["298fded6-b252-4166-a473-f405e935f58d"]
    }
    user_risk_levels = ["medium"]
  }
  grant_controls = {
    operator                = "AND"
    authentication_strength = { id = "00000000-0000-0000-0000-000000000003" }
    built_in_controls       = ["passwordChange"]
  }
  session_controls = {
    sign_in_frequency = {
      authentication_type = "primaryAndSecondaryAuthentication"
      frequency_interval  = "everyTime"
      is_enabled          = true
    }
  }
  state = "disabled"
}
