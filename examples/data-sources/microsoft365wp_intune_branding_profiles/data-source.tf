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


data "microsoft365wp_intune_branding_profiles" "all" {
}

output "microsoft365wp_intune_branding_profiles_all" {
  value = { for x in data.microsoft365wp_intune_branding_profiles.all.intune_branding_profiles : x.id => x }
}

data "microsoft365wp_intune_branding_profiles" "not_default" {
  is_default_profile = false
}

output "microsoft365wp_intune_branding_profiles_not_default" {
  value = { for x in data.microsoft365wp_intune_branding_profiles.not_default.intune_branding_profiles : x.id => x }
}
