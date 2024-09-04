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


data "microsoft365wp_azure_ad_windows_autopilot_deployment_profiles" "all" {
  exclude_ids = ["01234567-89ab-cdef-0123-456789abcdee", "01234567-89ab-cdef-0123-456789abcdef"]
}

locals {
  tf_profile = [
    for x in data.microsoft365wp_azure_ad_windows_autopilot_deployment_profiles.all.azure_ad_windows_autopilot_deployment_profiles :
    x if startswith(x.display_name, "TF Test")
  ][0]
}
data "microsoft365wp_azure_ad_windows_autopilot_deployment_profile_assignments" "all" {
  azure_ad_windows_autopilot_deployment_profile_id = local.tf_profile.id
}

output "microsoft365wp_azure_ad_windows_autopilot_deployment_profiles" {
  value = { for x in data.microsoft365wp_azure_ad_windows_autopilot_deployment_profiles.all.azure_ad_windows_autopilot_deployment_profiles : x.id => x }
}

output "microsoft365wp_azure_ad_windows_autopilot_deployment_profile_assignments" {
  value = { for x in data.microsoft365wp_azure_ad_windows_autopilot_deployment_profile_assignments.all.azure_ad_windows_autopilot_deployment_profile_assignments : x.id => x }
}
