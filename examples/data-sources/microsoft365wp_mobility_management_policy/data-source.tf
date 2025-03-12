terraform {
  required_providers {
    microsoft365wp = {
      source = "terraprovider/microsoft365wp"
    }
  }
}


data "azuread_application_published_app_ids" "well_known" {}

data "microsoft365wp_mobility_management_policy" "device" {
  policy_type = "device"
  id          = data.azuread_application_published_app_ids.well_known.result.MicrosoftIntune
}

output "microsoft365wp_mobility_management_policy_device" {
  value = data.microsoft365wp_mobility_management_policy.device
}

data "microsoft365wp_mobility_management_policy" "app" {
  policy_type = "app"
  id          = data.azuread_application_published_app_ids.well_known.result.MicrosoftIntune
}

output "microsoft365wp_mobility_management_policy_app" {
  value = data.microsoft365wp_mobility_management_policy.app
}
