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


resource "microsoft365wp_targeted_managed_app_configuration" "test" {
  display_name = "TF Test"

  app_group_type = "selectedPublicApps"
  apps = [
    { app_id = { android = { package_id = "com.microsoft.office.outlook" } } },
    { app_id = { ios = { bundle_id = "com.microsoft.office.outlook" } } },
  ]

  custom_settings = [
    {
      name  = "General Test"
      value = "General Value"
    },
    {
      name  = "com.microsoft.outlook.Mail.FocusedInbox"
      value = "false"
    },
  ]

  assignments = [
    { target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } },
  ]
}
