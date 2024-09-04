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


resource "microsoft365wp_android_managed_app_protection" "test1" {
  display_name = "TF Test Android 1"

  app_group_type = "allCoreMicrosoftApps"

  assignments = [
    { target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } },
  ]
}

resource "microsoft365wp_android_managed_app_protection" "test2" {
  display_name = "TF Test Android 2"

  app_group_type = "selectedPublicApps"
  apps = [
    {
      app_id = {
        android = {
          package_id = "com.adobe.reader"
        }
      }
    }
  ]

  allowed_outbound_data_transfer_destinations = "managedApps"
  exempted_app_packages = [
    {
      name  = "Name"
      value = "andValue"
    },
  ]
}
