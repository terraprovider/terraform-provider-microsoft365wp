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


resource "microsoft365wp_windows_driver_update_profile" "automatic" {
  display_name = "TF Test Automatic"

  approval_type               = "automatic"
  deployment_deferral_in_days = 3

  assignments = [
    { target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } },
    { target = { group = { group_id = "62e39046-aad3-4423-98e0-b486e3538aff" } } },
  ]
}

resource "microsoft365wp_windows_driver_update_profile" "manual" {
  display_name = "TF Test Manual"

  approval_type = "manual"

  assignments = [
    { target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } },
  ]
}
