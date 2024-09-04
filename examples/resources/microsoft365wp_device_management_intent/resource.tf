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


resource "microsoft365wp_device_management_intent" "test_macos_filevault" {
  display_name = "TF Test (macOS FileVault)"

  template_id = "a239407c-698d-4ef8-b314-e3ae409204b8"

  settings = [

    {
      definition_id = "deviceConfiguration--macOSEndpointProtectionConfiguration_fileVaultHidePersonalRecoveryKey"
      boolean       = { value = true }
    },
    {
      definition_id = "deviceConfiguration--macOSEndpointProtectionConfiguration_fileVaultNumberOfTimesUserCanIgnore"
      integer       = { value = 3 }
    },
    {
      definition_id = "deviceConfiguration--macOSEndpointProtectionConfiguration_fileVaultPersonalRecoveryKeyHelpMessage"
      string        = { value = "Microsoft Intune Company Portal: portal.manage.microsoft.com" }
    },
    {
      definition_id = "deviceConfiguration--macOSEndpointProtectionConfiguration_fileVaultDisablePromptAtSignOut"
      boolean       = { value = true }
    },
    {
      definition_id = "deviceConfiguration--macOSEndpointProtectionConfiguration_fileVaultAllowDeferralUntilSignOut"
      boolean       = { value = true }
    },
    {
      definition_id = "deviceConfiguration--macOSEndpointProtectionConfiguration_fileVaultPersonalRecoveryKeyRotationInMonths"
      integer       = { value = 6 }
    },
    {
      definition_id = "deviceConfiguration--macOSEndpointProtectionConfiguration_fileVaultEnabled"
      boolean       = { value = true }
    },
    {
      definition_id = "deviceConfiguration--macOSEndpointProtectionConfiguration_fileVaultSelectedRecoveryKeyTypes"
      collection = {
        value_json = jsonencode(["personalRecoveryKey"])
      }
    },
  ]

  assignments = [
    { target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } },
    { target = { group = { group_id = "62e39046-aad3-4423-98e0-b486e3538aff" } } },
  ]
}
