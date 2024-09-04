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


resource "microsoft365wp_cloud_pc_user_setting" "test" {
  display_name = "TF Test"

  restore_point_setting = {
    frequency_type = "sixteenHours"
  }

  assignments = [
    for x in [
      "298fded6-b252-4166-a473-f405e935f58d",
      "62e39046-aad3-4423-98e0-b486e3538aff",
    ] :
    { target = { group = { group_id = x } } }
  ]
}
