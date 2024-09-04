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


data "microsoft365wp_windows_feature_update_profile" "one" {
  id = "99645f7a-f30a-457b-a78d-f4f5514fcd51"
}

output "microsoft365wp_windows_feature_update_profile" {
  value = data.microsoft365wp_windows_feature_update_profile.one
}
