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


data "microsoft365wp_azure_ad_windows_autopilot_deployment_profile" "one" {
  id = "23b5545c-f948-494c-9495-3785d779806a"
}

data "microsoft365wp_azure_ad_windows_autopilot_deployment_profile_assignments" "all" {
  azure_ad_windows_autopilot_deployment_profile_id = data.microsoft365wp_azure_ad_windows_autopilot_deployment_profile.one.id
}

data "microsoft365wp_azure_ad_windows_autopilot_deployment_profile_assignment" "one" {
  azure_ad_windows_autopilot_deployment_profile_id = data.microsoft365wp_azure_ad_windows_autopilot_deployment_profile.one.id
  id = tolist(
    data.microsoft365wp_azure_ad_windows_autopilot_deployment_profile_assignments.all.azure_ad_windows_autopilot_deployment_profile_assignments
  )[0].id
}

output "microsoft365wp_azure_ad_windows_autopilot_deployment_profile" {
  value = data.microsoft365wp_azure_ad_windows_autopilot_deployment_profile.one
}

output "microsoft365wp_azure_ad_windows_autopilot_deployment_profile_assignment" {
  value = data.microsoft365wp_azure_ad_windows_autopilot_deployment_profile_assignment.one
}
