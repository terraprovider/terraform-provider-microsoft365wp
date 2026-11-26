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


data "microsoft365wp_dep_enrollment_profiles" "all" {
  dep_onboarding_setting_id = "437dc82e-9725-45e1-8126-ab8dc7014b1d"
}

output "microsoft365wp_dep_enrollment_profiles" {
  value = { for x in data.microsoft365wp_dep_enrollment_profiles.all.dep_enrollment_profiles : x.id => x }
}
