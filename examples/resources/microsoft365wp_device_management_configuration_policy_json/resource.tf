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


resource "microsoft365wp_device_management_configuration_policy_json" "test1" {
  name = "TF Test JSON 1"

  settings_json = [for x in [
    {
      "@odata.type"       = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance",
      settingDefinitionId = "device_vendor_msft_bitlocker_requiredeviceencryption",
      choiceSettingValue = {
        value = "device_vendor_msft_bitlocker_requiredeviceencryption_0",
      }
    },
    {
      "@odata.type"       = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
      settingDefinitionId = "device_vendor_msft_policy_config_microsoft_edgev78diff~policy~microsoft_edge_adssettingforintrusiveadssites"
      choiceSettingValue = {
        value = "device_vendor_msft_policy_config_microsoft_edgev78diff~policy~microsoft_edge_adssettingforintrusiveadssites_1"
        children = [
          {
            "@odata.type"       = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
            settingDefinitionId = "device_vendor_msft_policy_config_microsoft_edgev78diff~policy~microsoft_edge_adssettingforintrusiveadssites_adssettingforintrusiveadssites"
            choiceSettingValue = {
              value = "device_vendor_msft_policy_config_microsoft_edgev78diff~policy~microsoft_edge_adssettingforintrusiveadssites_adssettingforintrusiveadssites_2"
            }
          }
        ]
      }
    },
  ] : jsonencode(x)]

  assignments = [
    for x in [
      "298fded6-b252-4166-a473-f405e935f58d",
      "62e39046-aad3-4423-98e0-b486e3538aff",
    ] :
    { target = { group = { group_id = x } } }
  ]
}

resource "microsoft365wp_device_management_configuration_policy_json" "test2" {
  name = "TF Test JSON 2"

  settings_json = [for x in [
    {
      "@odata.type"       = "#microsoft.graph.deviceManagementConfigurationSimpleSettingCollectionInstance"
      settingDefinitionId = "device_vendor_msft_policy_config_defender_excludedpaths"
      simpleSettingCollectionValue = [
        for path in [
          "C:\\Windows",
          "C:\\Temp",
        ] :
        {
          "@odata.type" = "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
          value         = path
        }
      ]
    }
  ] : jsonencode(x)]

  assignments = [
    { target = { all_devices = {} } },
    { target = { all_licensed_users = {} } },
    { target = { exclusion_group = { group_id = "62e39046-aad3-4423-98e0-b486e3538aff" } } },
  ]
}

resource "microsoft365wp_device_management_configuration_policy_json" "test_firewall" {
  name         = "TF Test JSON Firewall"
  technologies = "mdm,microsoftSense"
  template_reference = {
    template_id = "6078910e-d808-4a9f-a51d-1b8a7bacb7c0_1"
  }

  settings_json = [for x in [
    {
      "@odata.type" = "#microsoft.graph.deviceManagementConfigurationChoiceSettingInstance"
      settingInstanceTemplateReference = {
        settingInstanceTemplateId = "38329af6-2670-4a71-972d-482010ca97fc"
      }
      settingDefinitionId = "vendor_msft_firewall_mdmstore_global_disablestatefulftp"
      choiceSettingValue = {
        settingValueTemplateReference = {
          settingValueTemplateId = "559f6e01-53a9-4c10-9f10-d09d8fe7f903"
          useTemplateDefault     = false
        }
        value = "vendor_msft_firewall_mdmstore_global_disablestatefulftp_true"
      }
    }
  ] : jsonencode(x)]
}
