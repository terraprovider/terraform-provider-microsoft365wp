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


data "microsoft365wp_intune_branding_profile" "one" {
  id = "2d7aa309-0a8d-42a1-94bb-323c6d95040f"
}

output "microsoft365wp_intune_branding_profile_one" {
  value = data.microsoft365wp_intune_branding_profile.one
}

data "microsoft365wp_intune_branding_profile" "default" {
  is_default_profile = true
}

output "microsoft365wp_intune_branding_profile_default" {
  value = data.microsoft365wp_intune_branding_profile.default
}
