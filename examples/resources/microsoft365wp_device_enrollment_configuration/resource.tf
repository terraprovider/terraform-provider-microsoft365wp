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


resource "microsoft365wp_device_enrollment_configuration" "windows10_esp" {
  display_name = "TF Test Windows 10 ESP"

  windows10_esp = {
    allow_device_reset_on_install_failure = true
  }

  assignments = [
    { target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } },
  ]
}


resource "microsoft365wp_device_enrollment_configuration" "device_enrollment_limit_default" {
  display_name = "###TfWorkplaceDefault###" # mandatory to target default config
  priority     = 0                          # mandatory to target default config

  device_enrollment_limit = {
    limit = 11
  }
}

resource "microsoft365wp_device_enrollment_configuration" "device_enrollment_limit_1" {
  display_name = "TF Test Enrollment Limit 1"

  device_enrollment_limit = {
    limit = 3
  }

  assignments = [
    { target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } },
  ]

  priority = 1
}

resource "microsoft365wp_device_enrollment_configuration" "device_enrollment_limit_2" {
  display_name = "TF Test Enrollment Limit 2"

  device_enrollment_limit = {
    limit = 4
  }

  priority = 2
}

resource "microsoft365wp_device_enrollment_configuration" "device_enrollment_limit_3" {
  display_name = "TF Test Enrollment Limit 3"

  device_enrollment_limit = {
    limit = 5
  }

  priority = 3
}


resource "microsoft365wp_device_enrollment_configuration" "platform_restrictions_default" {
  display_name = "###TfWorkplaceDefault###" # mandatory to target default config
  priority     = 0                          # mandatory to target default config

  platform_restrictions = {
    android_for_work_restriction = {
      platform_blocked = true
    }
  }
}


resource "microsoft365wp_device_enrollment_configuration" "single_platform_restriction" {
  display_name = "TF Test Platform Restriction"

  single_platform_restriction = {
    platform_type = "windows"
    platform_restriction = {
      personal_device_enrollment_blocked = true
      os_minimum_version                 = "10.0.19044"
    }
  }
}


resource "microsoft365wp_device_enrollment_configuration" "windows_hello_for_business_default" {
  display_name = "###TfWorkplaceDefault###" # mandatory to target default config
  priority     = 0                          # mandatory to target default config

  windows_hello_for_business = {
    pin_maximum_length = 126
  }
}
