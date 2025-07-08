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


resource "microsoft365wp_azure_ad_windows_autopilot_deployment_profile" "test" {
  display_name = "TF Test 1"

  locale = "" // user select
}

// For this resource MS Graph requires for assignments to be individually created or deleted.
// Therefore we use a separate resource type for them.
resource "microsoft365wp_azure_ad_windows_autopilot_deployment_profile_assignment" "test" {
  azure_ad_windows_autopilot_deployment_profile_id = microsoft365wp_azure_ad_windows_autopilot_deployment_profile.test.id

  for_each = toset(["298fded6-b252-4166-a473-f405e935f58d", "62e39046-aad3-4423-98e0-b486e3538aff"])

  target = { group = { group_id = each.value } }
}
