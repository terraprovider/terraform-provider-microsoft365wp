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


resource "microsoft365wp_device_and_app_management_assignment_filter" "devices_win10" {
  display_name = "TF Test Devices Win10"
  platform     = "windows10AndLater"
  rule         = "(device.manufacturer -eq \"some non-exsiting dummy value\")"
}

resource "microsoft365wp_device_configuration" "devices_win10" {
  display_name = "TF Test Devices Win10 Filter"
  windows_update_for_business = {
    user_pause_access     = "disabled"
    installation_schedule = { active_hours = {} }
  }

  assignments = [
    { target = {
      all_devices = {}
      filter_type = "include"
      filter_id   = microsoft365wp_device_and_app_management_assignment_filter.devices_win10.id
    } },
  ]
}


resource "microsoft365wp_device_and_app_management_assignment_filter" "apps_ios" {
  display_name = "TF Test Apps iOS"
  platform     = "iOSMobileApplicationManagement"
  rule         = "(app.appVersion -eq \"1.2.3.4\")"
}
