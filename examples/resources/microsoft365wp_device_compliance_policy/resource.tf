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


resource "microsoft365wp_device_compliance_policy" "android" {
  display_name               = "TF Test Android"
  scheduled_actions_for_rule = [{ scheduled_action_configurations = [{}] }]
  android = {
    password_required = true
  }
}

resource "microsoft365wp_device_compliance_policy" "android_device_owner" {
  display_name               = "TF Test Android Device Owner"
  scheduled_actions_for_rule = [{ scheduled_action_configurations = [{}] }]
  android_device_owner = {
    password_required = true
  }
}

resource "microsoft365wp_device_compliance_policy" "android_work_profile" {
  display_name               = "TF Test Android Work Profile"
  scheduled_actions_for_rule = [{ scheduled_action_configurations = [{}] }]
  android_work_profile = {
    password_required = true
  }
}

resource "microsoft365wp_device_compliance_policy" "aosp_device_owner" {
  display_name               = "TF Test AOSP"
  scheduled_actions_for_rule = [{ scheduled_action_configurations = [{}] }]
  aosp_device_owner = {
    password_required = true
  }
}

resource "microsoft365wp_device_compliance_policy" "ios" {
  display_name               = "TF Test iOS"
  scheduled_actions_for_rule = [{ scheduled_action_configurations = [{}] }]
  ios = {
    passcode_required = true
  }
}

resource "microsoft365wp_device_compliance_policy" "macos" {
  display_name               = "TF Test macOS"
  scheduled_actions_for_rule = [{ scheduled_action_configurations = [{}] }]
  macos = {
    password_required = true
  }
  assignments = [
    for x in [
      "298fded6-b252-4166-a473-f405e935f58d",
      "62e39046-aad3-4423-98e0-b486e3538aff",
    ] :
    { target = { group = { group_id = x } } }
  ]
}

resource "microsoft365wp_device_compliance_policy" "win10" {
  display_name               = "TF Test Windows 10"
  scheduled_actions_for_rule = [{ scheduled_action_configurations = [{}] }]
  windows10 = {
    bitlocker_enabled = true
  }
  assignments = [
    { target = { all_devices = {} } },
    { target = { all_licensed_users = {} } },
    { target = { exclusion_group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } },
  ]
}
