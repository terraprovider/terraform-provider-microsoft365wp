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


data "microsoft365wp_windows_feature_update_profiles" "all" {
}

output "microsoft365wp_windows_feature_update_profiles" {
  value = { for x in data.microsoft365wp_windows_feature_update_profiles.all.windows_feature_update_profiles : x.id => x }
}
