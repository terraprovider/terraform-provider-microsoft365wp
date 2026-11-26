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
  display_name              = "TF Test iOS/iPadOS"
  ios = {
    supervised_mode_enabled = true
  }
}

# Make the classic profile the default profile of the enrollment program token for new iOS/iPadOS devices.
resource "microsoft365wp_dep_default_enrollment_profile" "ios" {
  dep_onboarding_setting_id = var.dep_onboarding_setting_id
  id                        = microsoft365wp_dep_enrollment_profile.ios.id
}


# For an ADE enrollment policy (settings catalog), use the id under which MS Graph lists it as an enrollment profile
# of the token.
variable "ade_enrollment_policy_id" {
  description = "Id of a macOS ADE enrollment policy (deviceManagement/configurationPolicies)"
  type        = string
}

resource "microsoft365wp_dep_default_enrollment_profile" "macos" {
  dep_onboarding_setting_id = var.dep_onboarding_setting_id
  id                        = "ECV2_${var.dep_onboarding_setting_id}_${var.ade_enrollment_policy_id}"
}
