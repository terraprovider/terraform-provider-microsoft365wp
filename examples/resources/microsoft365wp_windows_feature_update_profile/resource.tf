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


locals {
  # take 25th of current month, add 10 days, and then take 1st of that month (i.e. take 1st of next month)
  start_date = formatdate("YYYY-MM-01'T'00:00:00Z", timeadd(formatdate("YYYY-MM-25'T'00:00:00Z", plantimestamp()), "240h"))
  end_date   = formatdate("YYYY-MM-15'T'00:00:00Z", local.start_date)
}

resource "microsoft365wp_windows_feature_update_profile" "test_immediate" {
  display_name = "TF Test Immediate"

  feature_update_version = "Windows 10, version 22H2"
  rollout_settings       = {}

  assignments = [
    { target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } },
  ]
}

resource "microsoft365wp_windows_feature_update_profile" "test_graduate" {
  display_name = "TF Test Graduate"

  feature_update_version = "Windows 10, version 22H2"
  rollout_settings = {
    offer_start_date_time_in_utc = local.start_date
    offer_end_date_time_in_utc   = local.end_date
    offer_interval_in_days       = 3
  }
}
