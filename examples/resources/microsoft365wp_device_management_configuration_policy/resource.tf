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


resource "microsoft365wp_device_management_configuration_policy" "test1" {
  name = "TF Test Typed 1"

  settings = [
    { instance = {
      definition_id = "device_vendor_msft_policy_config_microsoft_edgev78diff~policy~microsoft_edge_adssettingforintrusiveadssites"
      choice = {
        value = {
          value = "device_vendor_msft_policy_config_microsoft_edgev78diff~policy~microsoft_edge_adssettingforintrusiveadssites_1"
          children = [
            {
              definition_id = "device_vendor_msft_policy_config_microsoft_edgev78diff~policy~microsoft_edge_adssettingforintrusiveadssites_adssettingforintrusiveadssites"
              choice        = { value = { value = "device_vendor_msft_policy_config_microsoft_edgev78diff~policy~microsoft_edge_adssettingforintrusiveadssites_adssettingforintrusiveadssites_2" } }
            },
          ]
        }
      }
    } },
    { instance = {
      definition_id = "device_vendor_msft_bitlocker_requiredeviceencryption"
      choice        = { value = { value = "device_vendor_msft_bitlocker_requiredeviceencryption_0" } }
    } },
    { instance = {
      definition_id = "device_vendor_msft_policy_config_internetexplorer_allowsitetozoneassignmentlist"
      choice = { value = {
        value = "device_vendor_msft_policy_config_internetexplorer_allowsitetozoneassignmentlist_1"
        children = [
          {
            definition_id = "device_vendor_msft_policy_config_internetexplorer_allowsitetozoneassignmentlist_iz_zonemapprompt"
            group_collection = { values = [
              { children = [
                {
                  definition_id = "device_vendor_msft_policy_config_internetexplorer_allowsitetozoneassignmentlist_iz_zonemapprompt_key"
                  simple        = { value = { string = { value = "contoso.sharepoint.com" } } }
                },
                {
                  definition_id = "device_vendor_msft_policy_config_internetexplorer_allowsitetozoneassignmentlist_iz_zonemapprompt_value"
                  simple        = { value = { string = { value = "2" } } }
                },
              ] },
              { children = [
                {
                  definition_id = "device_vendor_msft_policy_config_internetexplorer_allowsitetozoneassignmentlist_iz_zonemapprompt_key"
                  simple        = { value = { string = { value = "contoso-my.sharepoint.com" } } }
                },
                {
                  definition_id = "device_vendor_msft_policy_config_internetexplorer_allowsitetozoneassignmentlist_iz_zonemapprompt_value"
                  simple        = { value = { string = { value = "2" } } }
                },
              ] },
            ] }
          },
        ]
      } }
    } },
  ]

  assignments = [
    for x in [
      "298fded6-b252-4166-a473-f405e935f58d",
      "62e39046-aad3-4423-98e0-b486e3538aff",
    ] :
    { target = { group = { group_id = x } } }
  ]
}

resource "microsoft365wp_device_management_configuration_policy" "test2" {
  name = "TF Test Typed 2"

  settings = [
    { instance = {
      definition_id = "device_vendor_msft_policy_config_defender_excludedpaths"
      simple_collection = { values = [
        for path in [
          "C:\\Windows",
          "C:\\Temp",
        ] : { string = { value = path } }
      ] }
    } },
  ]

  assignments = [
    { target = { all_devices = {} } },
    { target = { all_licensed_users = {} } },
    { target = { exclusion_group = { group_id = "62e39046-aad3-4423-98e0-b486e3538aff" } } },
  ]
}

resource "microsoft365wp_device_management_configuration_policy" "test3" {
  name = "TF Test Typed 3"

  settings = [
    { instance = {
      definition_id = "device_vendor_msft_policy_config_chromeintunev1~policy~googlechrome~extensions_extensionallowedtypes"
      choice = { value = {
        value = "device_vendor_msft_policy_config_chromeintunev1~policy~googlechrome~extensions_extensionallowedtypes_1"
        children = [
          {
            definition_id = "device_vendor_msft_policy_config_chromeintunev1~policy~googlechrome~extensions_extensionallowedtypes_extensionallowedtypesdesc"
            simple_collection = { values = [
              { string = { value = "extension" } },
              { string = { value = "hosted_app" } },
              { string = { value = "platform_app" } },
              { string = { value = "theme" } },
            ] }
          },
        ]
      } }
    } },
    { instance = {
      definition_id = "device_vendor_msft_policy_config_chromeintunev1~policy~googlechrome~extensions_extensioninstallallowlist"
      choice = { value = {
        value = "device_vendor_msft_policy_config_chromeintunev1~policy~googlechrome~extensions_extensioninstallallowlist_1"
        children = [
          {
            definition_id = "device_vendor_msft_policy_config_chromeintunev1~policy~googlechrome~extensions_extensioninstallallowlist_extensioninstallallowlistdesc"
            simple_collection = { values = [
              { string = { value = "ajanlknhcmbhbdafadmkobjnfkhdiegm" } },
              { string = { value = "gojbdfnpnhogfdgjbigejoaolejmgdhk" } },
              { string = { value = "gpaiobkfhnonedkhhfjpmhdalgeoebfa" } },
              { string = { value = "ndjpnladcallmjemlbaebfadecfhkepb" } },
            ] }
          },
        ]
      } }
    } },
    { instance = {
      definition_id = "device_vendor_msft_policy_config_chromeintunev1~policy~googlechrome~extensions_extensioninstallblocklist"
      choice = { value = {
        value = "device_vendor_msft_policy_config_chromeintunev1~policy~googlechrome~extensions_extensioninstallblocklist_1"
        children = [
          {
            definition_id = "device_vendor_msft_policy_config_chromeintunev1~policy~googlechrome~extensions_extensioninstallblocklist_extensioninstallblocklistdesc"
            simple_collection = { values = [
              { string = { value = "*" } },
            ] }
          },
        ]
      } }
    } },
    { instance = {
      definition_id = "device_vendor_msft_policy_config_chromeintunev1~policy~googlechrome~extensions_extensioninstallforcelist"
      choice = { value = {
        value = "device_vendor_msft_policy_config_chromeintunev1~policy~googlechrome~extensions_extensioninstallforcelist_1"
        children = [
          {
            definition_id = "device_vendor_msft_policy_config_chromeintunev1~policy~googlechrome~extensions_extensioninstallforcelist_extensioninstallforcelistdesc"
            simple_collection = { values = [
              { string = { value = "ggjhpefgjjfobnfoldnjipclpcfbgbhl" } },
            ] }
          },
        ]
      } }
    } },
    { instance = {
      definition_id = "device_vendor_msft_policy_config_chromeintunev1~policy~googlechrome~extensions_blockexternalextensions"
      choice        = { value = { value = "device_vendor_msft_policy_config_chromeintunev1~policy~googlechrome~extensions_blockexternalextensions_1" } }
    } },
  ]
}

resource "microsoft365wp_device_management_configuration_policy" "test4" {
  name = "TF Test Typed 4 (macOS)"

  platforms    = "macOS"
  technologies = "mdm,appleRemoteManagement"

  settings = [
    { instance = {
      definition_id = "com.apple.tcc.configuration-profile-policy_com.apple.tcc.configuration-profile-policy"
      group_collection = { values = [
        {
          children = [
            {
              definition_id = "com.apple.tcc.configuration-profile-policy_services"
              group_collection = { values = [{
                children = [
                  {
                    definition_id = "com.apple.tcc.configuration-profile-policy_services_accessibility"
                    group_collection = { values = [{
                      children = [
                        {
                          definition_id = "com.apple.tcc.configuration-profile-policy_services_accessibility_item_allowed"
                          choice        = { value = { value = "com.apple.tcc.configuration-profile-policy_services_accessibility_item_allowed_true" } }
                        },
                        {
                          definition_id = "com.apple.tcc.configuration-profile-policy_services_accessibility_item_authorization"
                          choice        = { value = { value = "com.apple.tcc.configuration-profile-policy_services_accessibility_item_authorization_0" } }
                        },
                        {
                          definition_id = "com.apple.tcc.configuration-profile-policy_services_accessibility_item_coderequirement"
                          simple        = { value = { string = { value = "identifier \"com.microsoft.teams\" and anchor apple generic and certificate 1[field.1.2.840.113635.100.6.2.6] /* exists */ and certificate leaf[field.1.2.840.113635.100.6.1.13] /* exists */ and certificate leaf[subject.OU] = UBF8T346G9" } } }
                        },
                        {
                          definition_id = "com.apple.tcc.configuration-profile-policy_services_accessibility_item_identifier"
                          simple        = { value = { string = { value = "com.microsoft.teams" } } }
                        },
                        {
                          definition_id = "com.apple.tcc.configuration-profile-policy_services_accessibility_item_identifiertype"
                          choice        = { value = { value = "com.apple.tcc.configuration-profile-policy_services_accessibility_item_identifiertype_0" } }
                        },
                        {
                          definition_id = "com.apple.tcc.configuration-profile-policy_services_accessibility_item_staticcode"
                          choice        = { value = { value = "com.apple.tcc.configuration-profile-policy_services_accessibility_item_staticcode_false" } }
                        },
                      ]
                    }, ] }
                  },
                ] }, ]
              }
            },
          ]
        }, ]
      } }
    },
  ]
}

resource "microsoft365wp_device_management_configuration_policy" "test5" {
  name = "TF Test Typed 5"

  template_reference = { id = "4321b946-b76b-4450-8afd-769c08b16ffc_1" }

  settings = [
    { instance = {
      definition_id = "device_vendor_msft_policy_config_applicationcontrol_policies_{policyguid}_policiesoptions"
      choice = { value = {
        value              = "device_vendor_msft_policy_config_applicationcontrol_built_in_controls_selected"
        template_reference = { id = "b28c7dc4-c7b2-4ce2-8f51-6ebfd3ea69d3" }
        children = [
          {
            definition_id = "device_vendor_msft_policy_config_applicationcontrol_built_in_controls"
            group_collection = { values = [
              { children = [
                {
                  definition_id = "device_vendor_msft_policy_config_applicationcontrol_built_in_controls_enable_app_control"
                  choice = { value = {
                    value = "device_vendor_msft_policy_config_applicationcontrol_built_in_controls_enable_app_control_0"
                  } }
                },
                {
                  definition_id = "device_vendor_msft_policy_config_applicationcontrol_built_in_controls_trust_apps"
                  choice_collection = { values = [
                    { value = "device_vendor_msft_policy_config_applicationcontrol_built_in_controls_trust_apps_1" }
                  ] }
                },
              ] },
            ] }
          },
        ]
      } }
      template_reference = { id = "1de98212-6949-42dc-a89c-e0ff6e5da04b" }
    } },
  ]
}

resource "microsoft365wp_device_management_configuration_policy" "test_firewall" {
  name = "TF Test Typed Firewall"

  technologies       = "mdm,microsoftSense"
  template_reference = { id = "6078910e-d808-4a9f-a51d-1b8a7bacb7c0_1" }

  settings = [
    { instance = {
      definition_id = "vendor_msft_firewall_mdmstore_global_disablestatefulftp"
      choice = { value = {
        value              = "vendor_msft_firewall_mdmstore_global_disablestatefulftp_true"
        template_reference = { id = "559f6e01-53a9-4c10-9f10-d09d8fe7f903" }
      } }
      template_reference = { id = "38329af6-2670-4a71-972d-482010ca97fc" }
    } },
  ]
}

resource "microsoft365wp_device_management_configuration_policy" "test_atp" {
  name = "TF Test Typed ATP"

  technologies       = "mdm,microsoftSense"
  template_reference = { id = "0385b795-0f2f-44ac-8602-9f65bf6adede_1" }

  settings = [
    { instance = {
      definition_id = "device_vendor_msft_windowsadvancedthreatprotection_configurationtype"
      choice = { value = {
        value              = "device_vendor_msft_windowsadvancedthreatprotection_configurationtype_autofromconnector"
        template_reference = { id = "e5c7c98c-c854-4140-836e-bd22db59d651" }
        children = [
          {
            definition_id = "device_vendor_msft_windowsadvancedthreatprotection_onboarding_fromconnector"
            simple = { value = { secret = {
              value = "Microsoft ATP connector enabled"
              state = "notEncrypted"
            } } }
          },
        ]
      } }
      template_reference = { id = "23ab0ea3-1b12-429a-8ed0-7390cf699160" }
    } },
  ]

  # see note in docs about secret values
  lifecycle {
    ignore_changes = [
      settings[0].instance.choice.value.children[0].simple.value.secret
    ]
  }
}
