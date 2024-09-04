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


data "microsoft365wp_mobile_apps" "outlook_ios" {
  odata_filter = "contains(displayName, 'Outlook') and isof('microsoft.graph.iosStoreApp')"
}

resource "microsoft365wp_managed_device_mobile_app_configuration" "ios" {
  display_name = "TF Test iOS"

  targeted_mobile_apps = [
    one(data.microsoft365wp_mobile_apps.outlook_ios.mobile_apps[*].id),
  ]

  ios = {
    settings = [
      {
        app_config_key       = "Test Key"
        app_config_key_type  = "stringType"
        app_config_key_value = "Test Value"
      },
      {
        app_config_key       = "com.microsoft.outlook.Mail.FocusedInbox"
        app_config_key_type  = "booleanType"
        app_config_key_value = "true"
      },
    ]
  }
}


/*
##### Untested #####

resource "microsoft365wp_managed_device_mobile_app_configuration" "android_managed_store" {
  display_name = "TF Test Android Managed Store"

  targeted_mobile_apps = [
    "0ddf9bbb-abb6-4c9a-89a5-ff54d9640d40",
  ]

  android_managed_store = {
    profile_applicability = "default"
    package_id            = "app:com.microsoft.office.outlook"
    payload_json_base64 = base64encode(jsonencode({
      kind      = "androidenterprise#managedConfiguration"
      productId = "app:com.microsoft.office.outlook"
      managedProperty = [
        { key = "com.microsoft.outlook.EmailProfile.AccountType", valueString = "ModernAuth" },
        { key = "com.microsoft.outlook.EmailProfile.EmailUPN", valueString = "{{userprincipalname}}" },
        { key = "com.microsoft.outlook.EmailProfile.EmailAddress", valueString = "{{mail}}" },
        { key = "com.microsoft.outlook.Contacts.LocalSyncEnabled", valueBool = true },
        { key = "com.microsoft.outlook.Mail.DefaultSignatureEnabled", valueBool = false },
      ]
    }))
  }
}
*/
