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


//
// NOTES
// =====
//
// Default Enrollment Configurations
// ---------------------------------
// To target _default_ enrollment configurations instead of creating new ones, `display_name` MUST BE set to
// "###TfWorkplaceDefault###" and `priority` MUST BE set to 0. `description`, `roleScopeTagIds` and `assignments` should
// not be set (their values will be ignored).
// The provider will then automatically import the existing default config from MS Graph on resource creation instead of
// creating a new MS Graph entity. On destruction, the provider will only remove the resource from Terraform state but
// not try to delete the default config entity from MS Graph.
// To compare your own Terraform config to the actual values of a default configuration in MS Graph, you can still
// import the entity into Terraform state as can be done with any other resource (e.g.
// `terraform import microsoft365wp_device_enrollment_configuration.windows_hello_for_business_default 01234567-89ab-cdef-0123-456789abcdef_DefaultWindowsHelloForBusiness`
// - get the real `id` directly from MS Graph using e.g. some browser's development tools). Then run `terraform plan` as
// usual to see the required changes.
//
// Priorities
// ----------
// MS Graph seems to be quite picky about setting priorities. They have to be sequential (i.e. there cannot be gaps when
// creating) and when updating priorities, values can only be switched with another exiting config (e.g. if the config
// with priority 1 gets deleted (e.g. manually) then the list will start with priority 2 and the `setPriority` action
// will refuse to set any other config to priority 1 anymore but only to priority 2).
// Therefore setting of configuration priorities has been implemented as best-effort only: Any errors returned from
// MS Graph during the `setPriority` action will only be shown as warnings in Terraform and more than one apply action
// may be required to actually get every priority set correctly.
// And in case of a situation that MS Graph is not willing (or able) to fulfill, there will continously remain changes.
// If you are really, really sure that your configuration should actually be fulfillable by MS Graph, then it may help
// to destroy all configurations again and start over.
//
// Delays Between Writes
// ---------------------
// As there have been situations during testing that resulted in multiple enrollment configurations with duplicate
// priority values (which MS Graph was unable to correct or reorder ever again), all write operations will incur a small
// delay to ensure individual priority values.
//


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
  display_name = "###TfWorkplaceDefault###" // mandatory to target default config
  priority     = 0                          // mandatory to target default config

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
  display_name = "###TfWorkplaceDefault###" // mandatory to target default config
  priority     = 0                          // mandatory to target default config

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
  display_name = "###TfWorkplaceDefault###" // mandatory to target default config
  priority     = 0                          // mandatory to target default config

  windows_hello_for_business = {
    pin_maximum_length = 126
  }
}
