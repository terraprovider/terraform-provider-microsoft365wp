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


data "microsoft365wp_dep_enrollment_profile" "one" {
  dep_onboarding_setting_id = "437dc82e-9725-45e1-8126-ab8dc7014b1d"
  id                        = "437dc82e-9725-45e1-8126-ab8dc7014b1d_6de29b1b-9b37-40f4-99b8-9d8a1217ebc6"
}

output "microsoft365wp_dep_enrollment_profile" {
  value     = data.microsoft365wp_dep_enrollment_profile.one
  sensitive = true # macos.admin_account_password
}
