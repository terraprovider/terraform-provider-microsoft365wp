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


variable "dep_onboarding_setting_id" {
  description = "Id of the enrollment program token (deviceManagement/depOnboardingSettings)"
  type        = string
}


resource "microsoft365wp_dep_enrollment_profile" "ios" {
  dep_onboarding_setting_id = var.dep_onboarding_setting_id

  display_name                                               = "TF Test iOS/iPadOS"
  description                                                = "Corporate iOS/iPadOS devices"
  requires_user_authentication                               = true
  require_company_portal_on_setup_assistant_enrolled_devices = true

  ios = {
    supervised_mode_enabled  = true
    is_mandatory             = true
    profile_removal_disabled = true
    itunes_pairing_mode      = "allow"
    support_department       = "IT"
    support_phone_number     = "+49 89 1234567"

    # Setup Assistant screens to skip: set the *_disabled attributes and enabled_skip_keys consistently
    apple_id_disabled             = true
    privacy_pane_disabled         = true
    siri_disabled                 = true
    terms_and_conditions_disabled = true
    enabled_skip_keys             = ["AppleID", "Privacy", "Siri", "TOS", "Intelligence"]
  }
}

resource "microsoft365wp_dep_enrollment_profile" "macos" {
  dep_onboarding_setting_id = var.dep_onboarding_setting_id

  display_name                 = "TF Test macOS"
  requires_user_authentication = true

  macos = {
    supervised_mode_enabled  = true
    is_mandatory             = true
    profile_removal_disabled = true

    file_vault_disabled = true
    siri_disabled       = true
    enabled_skip_keys   = ["FileVault", "Siri"]

    admin_account_user_name = "ADM-LAPS"
    admin_account_full_name = "Intune LAPS Administrator"
    hide_admin_account      = true
    dep_profile_admin_account_password_rotation_setting = {
      auto_rotation_period_in_days = 30
    }
  }
}


# Migration from msgraph_resource (provider microsoft/msgraph) without recreating the profile (Terraform >= 1.8,
# OpenTofu >= 1.10): move the state of the old resource to the new one (the id and the id of the enrollment program
# token are taken from the old state, everything else is read from MS Graph). After the next apply, the moved block
# can be deleted.

moved {
  from = msgraph_resource.dep_enrollment_profile
  to   = microsoft365wp_dep_enrollment_profile.migrated
}

resource "microsoft365wp_dep_enrollment_profile" "migrated" {
  dep_onboarding_setting_id = var.dep_onboarding_setting_id
  display_name              = "iOS/iPadOS - Default"
  ios = {
    supervised_mode_enabled = true
  }
}
