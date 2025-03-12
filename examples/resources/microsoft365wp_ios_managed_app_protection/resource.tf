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


resource "microsoft365wp_ios_managed_app_protection" "test1" {
  display_name = "TF Test iOS 1"

  app_group_type = "allCoreMicrosoftApps"

  pin_required_instead_of_biometric_timeout = "PT30M"

  assignments = [
    { target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } },
  ]
}

resource "microsoft365wp_ios_managed_app_protection" "test2" {
  display_name = "TF Test iOS 2"

  app_group_type = "selectedPublicApps"
  apps = [
    {
      app_id = {
        ios = {
          bundle_id = "com.microsoft.msedge"
        }
      }
    }
  ]

  allowed_outbound_data_transfer_destinations = "managedApps"
  exempted_app_protocols = [
    {
      name  = "Name"
      value = "andValue"
    },
  ]
}
