---
page_title: "microsoft365wp_device_configuration Resource - microsoft365wp"
subcategory: "MS Graph: Device configuration"
---

# microsoft365wp_device_configuration (Resource)

Device Configuration. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceconfiguration?view=graph-rest-beta

## Documentation Disclaimer

Please note that almost all information on this page has been sourced literally from the official Microsoft Graph API 
documentation and therefore is governed by Microsoft and not by the publishers of this provider.  
All supplements authored by the publishers of this provider have been explicitly marked as such.

## Example Usage

```terraform
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


resource "microsoft365wp_device_configuration" "android_device_owner_general_device" {
  display_name = "TF Test Android Enterprise Corp Owned Device Restrictions"
  android_device_owner_general_device = {
    camera_blocked = true
  }
}

resource "microsoft365wp_device_configuration" "android_work_profile_general_device" {
  display_name = "TF Test Android Enterprise Personally Owned Device Restrictions"
  android_work_profile_general_device = {
    work_profile_block_screen_capture = true
  }
}

resource "microsoft365wp_device_configuration" "edition_upgrade" {
  display_name = "TF Test Edition Upgrade"
  edition_upgrade = {
    target_edition = "windows10Enterprise"
    license_type   = "productKey"
    product_key    = "12345-22345-32345-42345-52345"
  }
}

resource "microsoft365wp_device_configuration" "ios_custom" {
  display_name = "TF Test iOS Custom"
  ios_custom = {
    name      = "Defender - Zero touch Control Filter"
    file_name = "Microsoft_Defender_for_Endpoint_Control_Filter_Zerotouch.mobileconfig"
    payload   = <<EOT
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>ConsentText</key>
	<dict>
		<key>default</key>
		<string>This profile will be used to analyze network traffic to ensure a safe browsing experience</string>
	</dict>
	<key>PayloadContent</key>
	<array>
		<dict>
			<key>FilterBrowsers</key>
			<true/>
			<key>FilterSockets</key>
			<false/>
			<key>FilterType</key>
			<string>Plugin</string>
			<key>PayloadDescription</key>
			<string>Configures content filtering settings</string>
			<key>PayloadDisplayName</key>
			<string>Microsoft Defender for Endpoint Content Filter</string>
			<key>PayloadIdentifier</key>
			<string>com.apple.webcontent-filter.64810102-01A9-4140-BB9C-B19DCCB1C5F4</string>
			<key>PayloadType</key>
			<string>com.apple.webcontent-filter</string>
			<key>PayloadUUID</key>
			<string>64810102-01A9-4140-BB9C-B19DCCB1C5F4</string>
			<key>PayloadVersion</key>
			<integer>1</integer>
			<key>PluginBundleID</key>
			<string>com.microsoft.scmx</string>
			<key>UserDefinedName</key>
			<string>Microsoft Defender for Endpoint Content Filter</string>
			<key>VendorConfig</key>
			<dict>
				<key>SilentOnboard</key>
				<string>true</string>
			</dict>
		</dict>
	</array>
	<key>PayloadDescription</key>
	<string>This profile will be used to analyze network traffic to ensure a safe browsing experience</string>
	<key>PayloadDisplayName</key>
	<string>Microsoft Defender for Endpoint Control Filter</string>
	<key>PayloadIdentifier</key>
	<string>com.microsoft.scmx.46C8CAF0-2D73-4F21-AC04-00283D167D28</string>
	<key>PayloadOrganization</key>
	<string>Microsoft</string>
	<key>PayloadRemovalDisallowed</key>
	<false/>
	<key>PayloadType</key>
	<string>Configuration</string>
	<key>PayloadUUID</key>
	<string>FCEC9441-5E6B-4063-B3D5-9C7A16EF9D41</string>
	<key>PayloadVersion</key>
	<integer>1</integer>
</dict>
</plist>
EOT
  }
}

resource "microsoft365wp_device_configuration" "ios_eas_email_profile" {
  display_name = "TF Test iOS EMail"
  ios_eas_email_profile = {
    host_name             = "outlook.office365.com"
    account_name          = "Corporate Email"
    username_aad_source   = "userPrincipalName"
    email_address_source  = "primarySmtpAddress"
    authentication_method = "usernameAndPassword"
  }
}

resource "microsoft365wp_device_configuration" "ios_device_features" {
  display_name = "TF Test iOS Device Features"
  ios_device_features = {
    lock_screen_footnote = "This device belongs to somebody."
    ios_single_sign_on_extension = {
      azure_ad = {
        configurations = [
          {
            key     = "browser_sso_interaction_enabled"
            integer = { value = 1 }
          },
          {
            key     = "disable_explicit_app_prompt_and_autologin"
            integer = { value = 1 }
          },
          {
            key     = "browser_sso_disable_mfa"
            integer = { value = 0 }
          },
          {
            key    = "device_registration"
            string = { value = "{{DEVICEREGISTRATION}}" }
          },
        ]
      }
    }
  }
}

resource "microsoft365wp_device_configuration" "ios_general_device" {
  display_name = "TF Test iOS General Device"
  ios_general_device = {
    passcode_required = true
  }
}

resource "microsoft365wp_device_configuration" "ios_update" {
  display_name = "TF Test Updates iOS"
  ios_update = {
    desired_os_version = "16.4"
  }
}

resource "microsoft365wp_device_configuration" "ios_vpn" {
  display_name = "TF Test iOS VPN"
  ios_vpn = {
    connection_type       = "customVpn"
    connection_name       = "VPN Connection 1"
    identifier            = "com.acme.vpnclient.applevpn.plugin"
    authentication_method = "usernameAndPassword"
    server = {
      address = "127.0.0.1"
    }
  }
}

resource "microsoft365wp_device_configuration" "macos_custom" {
  display_name = "TF Test macOS Custom"
  macos_custom = {
    name               = "Microsoft 365 Apps"
    deployment_channel = "deviceChannel"
    file_name          = "M365_Apps.mobileconfig"
    payload            = <<EOT
<?xml version="1.0" encoding="utf-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
		<key>PayloadContent</key>
		<array>
			<dict>
				<key>DiagnosticDataTypePreference</key>
				<string>BasicDiagnosticData</string>
				<key>HasUserSeenEnterpriseFREDialog</key>
				<true />
				<key>PayloadDisplayName</key>
				<string>Microsoft Office</string>
				<key>PayloadIdentifier</key>
				<string>com.microsoft.office.49555bd0-be2c-4ca6-a12f-06c7813531e0</string>
				<key>PayloadType</key>
				<string>com.microsoft.office</string>
				<key>PayloadUUID</key>
				<string>9d34c736-5337-4347-bf7e-16fc0cf12e5e</string>
				<key>PayloadVersion</key>
				<integer>1</integer>
				<key>SendAllTelemetryEnabled</key>
				<false />
			</dict>
		</array>
		<key>PayloadDisplayName</key>
		<string>M365 Apps</string>
		<key>PayloadIdentifier</key>
		<string>C4A8CND1272TQQ.a488d3df-763c-49f8-a92e-04642dcde833</string>
		<key>PayloadType</key>
		<string>Configuration</string>
		<key>PayloadUUID</key>
		<string>9d3a7d4a-40c3-4cfc-9fde-065ec5fa738a</string>
		<key>PayloadVersion</key>
		<integer>1</integer>
	</dict>
</plist>
EOT
  }
}

resource "microsoft365wp_device_configuration" "macos_custom_app" {
  display_name = "TF Test macOS Preferences File"
  macos_custom_app = {
    bundle_id         = "com.microsoft.Edge"
    file_name         = "Favorites.xml"
    configuration_xml = <<EOT
<key>ManagedFavorites</key>
<array>
  <dict>
    <key>toplevel_name</key>
    <string>Contoso</string>
  </dict>
  <dict>
    <key>name</key>
    <string>Contoso</string>
    <key>url</key>
    <string>https://www.contoso.com</string>
  </dict>
</array>
EOT
  }
}

resource "microsoft365wp_device_configuration" "macos_device_features" {
  display_name = "TF Test macOS Device Features"
  macos_device_features = {
    macos_single_sign_on_extension = {
      azure_ad = {
        configurations = [
          {
            key     = "browser_sso_interaction_enabled"
            integer = { value = 1 }
          },
          {
            key     = "disable_explicit_app_prompt"
            integer = { value = 1 }
          },
          {
            key     = "browser_sso_disable_mfa"
            integer = { value = 0 }
          },
        ]
      }
    }
  }
}

resource "microsoft365wp_device_configuration" "macos_extensions" {
  display_name = "TF Test macOS Extensions"
  macos_extensions = {
    kernel_extensions_allowed = [
      {
        bundle_id       = "com.company_1.application_a"
        team_identifier = ""
      }
    ]
    system_extensions_block_override = true
  }
}

resource "microsoft365wp_device_configuration" "macos_software_update" {
  display_name = "TF Test macOS Updates"
  macos_software_update = {
    all_other_update_behavior = "installASAP"
  }
}

resource "microsoft365wp_device_configuration" "windows_health_monitoring" {
  display_name = "TF Test Windows Health Monitoring"
  windows_health_monitoring = {
    allow = "enabled"
    scope = "bootPerformance,windowsUpdates"
  }
}

resource "microsoft365wp_device_configuration" "windows_update_for_business" {
  display_name = "TF Test Updates Windows"
  windows_update_for_business = {
    user_pause_access     = "disabled"
    installation_schedule = { active_hours = {} }
  }

  assignments = [
    { target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } },
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) Admin provided name of the device configuration.

### Optional

- `android_device_owner_general_device` (Attributes) This topic provides descriptions of the declared methods, properties and relationships exposed by the androidDeviceOwnerGeneralDeviceConfiguration resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceownergeneraldeviceconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--android_device_owner_general_device))
- `android_work_profile_general_device` (Attributes) Android Work Profile general device configuration. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androidworkprofilegeneraldeviceconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--android_work_profile_general_device))
- `assignments` (Attributes Set) The list of assignments. (see [below for nested schema](#nestedatt--assignments))
- `description` (String) Admin provided description of the Device Configuration.
- `device_management_applicability_rule_device_mode` (Attributes) The device mode applicability rule for this Policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-devicemanagementapplicabilityruledevicemode?view=graph-rest-beta (see [below for nested schema](#nestedatt--device_management_applicability_rule_device_mode))
- `device_management_applicability_rule_os_edition` (Attributes) The OS edition applicability for this Policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-devicemanagementapplicabilityruleosedition?view=graph-rest-beta (see [below for nested schema](#nestedatt--device_management_applicability_rule_os_edition))
- `device_management_applicability_rule_os_version` (Attributes) The OS version applicability rule for this Policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-devicemanagementapplicabilityruleosversion?view=graph-rest-beta (see [below for nested schema](#nestedatt--device_management_applicability_rule_os_version))
- `edition_upgrade` (Attributes) Windows 10 Edition Upgrade configuration. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-editionupgradeconfiguration?view=graph-rest-beta  
Provider Note: There have been reports that creating this type with an Entra Id application would fail even though it did have the permission `DeviceManagementConfiguration.ReadWrite.All`. Since this would be a limitation in MS Graph there is nothing that the provider could do about it. (see [below for nested schema](#nestedatt--edition_upgrade))
- `ios_custom` (Attributes) This topic provides descriptions of the declared methods, properties and relationships exposed by the iosCustomConfiguration resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioscustomconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_custom))
- `ios_device_features` (Attributes) iOS Device Features Configuration Profile. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosdevicefeaturesconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features))
- `ios_eas_email_profile` (Attributes) By providing configurations in this profile you can instruct the native email client on iOS devices to communicate with an Exchange server and get email, contacts, calendar, reminders, and notes. Furthermore, you can also specify how much email to sync and how often the device should sync. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioseasemailprofileconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_eas_email_profile))
- `ios_general_device` (Attributes) This topic provides descriptions of the declared methods, properties and relationships exposed by the iosGeneralDeviceConfiguration resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosgeneraldeviceconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_general_device))
- `ios_update` (Attributes) IOS Update Configuration, allows you to configure time window within week to install iOS updates / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosupdateconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_update))
- `ios_vpn` (Attributes) By providing the configurations in this profile you can instruct the iOS device to connect to desired VPN endpoint. By specifying the authentication method and security types expected by VPN endpoint you can make the VPN connection seamless for end user. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosvpnconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_vpn))
- `macos_custom` (Attributes) This topic provides descriptions of the declared methods, properties and relationships exposed by the macOSCustomConfiguration resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macoscustomconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_custom))
- `macos_custom_app` (Attributes) This topic provides descriptions of the declared methods, properties and relationships exposed by the macOSCustomAppConfiguration resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macoscustomappconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_custom_app))
- `macos_device_features` (Attributes) MacOS device features configuration profile. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macosdevicefeaturesconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features))
- `macos_extensions` (Attributes) MacOS extensions configuration profile. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macosextensionsconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_extensions))
- `macos_software_update` (Attributes) MacOS Software Update Configuration / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macossoftwareupdateconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_software_update))
- `role_scope_tag_ids` (Set of String) List of Scope Tags for this Entity instance. The _provider_ default value is `["0"]`.
- `windows_health_monitoring` (Attributes) Windows device health monitoring configuration / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowshealthmonitoringconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows_health_monitoring))
- `windows_update_for_business` (Attributes) Windows Update for business configuration, allows you to specify how and when Windows as a Service updates your Windows 10/11 devices with feature and quality updates. Supports ODATA clauses that DeviceConfiguration entity supports: $filter by types of DeviceConfiguration, $top, $select only DeviceConfiguration base properties, $orderby only DeviceConfiguration base properties, and $skip. The query parameter '$search' is not supported. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowsupdateforbusinessconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows_update_for_business))

### Read-Only

- `created_date_time` (String) DateTime the object was created.
- `id` (String) Key of the entity.
- `last_modified_date_time` (String) DateTime the object was last modified.
- `supports_scope_tags` (Boolean) Indicates whether or not the underlying Device Configuration supports the assignment of scope tags. Assigning to the ScopeTags property is not allowed when this value is false and entities will not be visible to scoped users. This occurs for Legacy policies created in Silverlight and can be resolved by deleting and recreating the policy in the Azure Portal. This property is read-only.
- `version` (Number) Version of the device configuration.

<a id="nestedatt--android_device_owner_general_device"></a>
### Nested Schema for `android_device_owner_general_device`

Optional:

- `accounts_block_modification` (Boolean) Indicates whether or not adding or removing accounts is disabled.
- `android_device_owner_delegated_scope_app_settings` (Attributes Set) Specifies the list of managed apps with app details and its associated delegated scope(s). This collection can contain a maximum of 500 elements. / Represents one item in the list of managed apps with app details and its associated delegated scope(s). / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceownerdelegatedscopeappsetting?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--android_device_owner_general_device--android_device_owner_delegated_scope_app_settings))
- `apps_allow_install_from_unknown_sources` (Boolean) Indicates whether or not the user is allowed to enable to unknown sources setting.
- `apps_auto_update_policy` (String) Indicates the value of the app auto update policy. / Android Device Owner possible values for states of the device's app auto update policy; possible values are: `notConfigured` (Not configured; this value is ignored.), `userChoice` (The user can control auto-updates.), `never` (Apps are never auto-updated.), `wiFiOnly` (Apps are auto-updated over Wi-Fi only.), `always` (Apps are auto-updated at any time. Data charges may apply.)
- `apps_default_permission_policy` (String) Indicates the permission policy for requests for runtime permissions if one is not defined for the app specifically. / Android Device Owner default app permission policy type; possible values are: `deviceDefault` (Device default value, no intent.), `prompt` (Prompt.), `autoGrant` (Auto grant.), `autoDeny` (Auto deny.)
- `apps_recommend_skipping_first_use_hints` (Boolean) Whether or not to recommend all apps skip any first-time-use hints they may have added.
- `azure_ad_shared_device_data_clear_apps` (Attributes Set) A list of managed apps that will have their data cleared during a global sign-out in AAD shared device mode. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--android_device_owner_general_device--azure_ad_shared_device_data_clear_apps))
- `bluetooth_block_configuration` (Boolean) Indicates whether or not to block a user from configuring bluetooth.
- `bluetooth_block_contact_sharing` (Boolean) Indicates whether or not to block a user from sharing contacts via bluetooth.
- `camera_blocked` (Boolean) Indicates whether or not to disable the use of the camera.
- `cellular_block_wifi_tethering` (Boolean) Indicates whether or not to block Wi-Fi tethering.
- `certificate_credential_configuration_disabled` (Boolean) Indicates whether or not to block users from any certificate credential configuration.
- `cross_profile_policies_allow_copy_paste` (Boolean) Indicates whether or not text copied from one profile (personal or work) can be pasted in the other.
- `cross_profile_policies_allow_data_sharing` (String) Indicates whether data from one profile (personal or work) can be shared with apps in the other profile. / An enum representing possible values for cross profile data sharing; possible values are: `notConfigured` (Not configured; this value defaults to CROSS_PROFILE_DATA_SHARING_UNSPECIFIED.), `crossProfileDataSharingBlocked` (Data cannot be shared from both the personal profile to work profile and the work profile to the personal profile.), `dataSharingFromWorkToPersonalBlocked` (Prevents users from sharing data from the work profile to apps in the personal profile. Personal data can be shared with work apps.), `crossProfileDataSharingAllowed` (Data from either profile can be shared with the other profile.), `unkownFutureValue` (Unknown future value (reserved, not used right now)). The _provider_ default value is `"notConfigured"`.
- `cross_profile_policies_show_work_contacts_in_personal_profile` (Boolean) Indicates whether or not contacts stored in work profile are shown in personal profile contact searches/incoming calls.
- `data_roaming_blocked` (Boolean) Indicates whether or not to block a user from data roaming.
- `date_time_configuration_blocked` (Boolean) Indicates whether or not to block the user from manually changing the date or time on the device
- `detailed_help_text` (Attributes) Represents the customized detailed help text provided to users when they attempt to modify managed settings on their device. / Represents a user-facing message with locale information as well as a default message to be used if the user's locale doesn't match with any of the localized messages / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceowneruserfacingmessage?view=graph-rest-beta (see [below for nested schema](#nestedatt--android_device_owner_general_device--detailed_help_text))
- `device_location_mode` (String) Indicates the location setting configuration for fully managed devices (COBO) and corporate owned devices with a work profile (COPE). / Android Device Owner Location Mode Type; possible values are: `notConfigured` (No restrictions on the location setting and no specific behavior is set or enforced. This is the default), `disabled` (Location setting is disabled on the device), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use)
- `device_owner_lock_screen_message` (Attributes) Represents the customized lock screen message provided to users when they attempt to modify managed settings on their device. / Represents a user-facing message with locale information as well as a default message to be used if the user's locale doesn't match with any of the localized messages / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceowneruserfacingmessage?view=graph-rest-beta (see [below for nested schema](#nestedatt--android_device_owner_general_device--device_owner_lock_screen_message))
- `enrollment_profile` (String) Indicates which enrollment profile you want to configure. / Android Device Owner Enrollment Profile types; possible values are: `notConfigured` (Not configured; this value is ignored.), `dedicatedDevice` (Dedicated device.), `fullyManaged` (Fully managed.). The _provider_ default value is `"notConfigured"`.
- `factory_reset_blocked` (Boolean) Indicates whether or not the factory reset option in settings is disabled.
- `factory_reset_device_administrator_emails` (Set of String) List of Google account emails that will be required to authenticate after a device is factory reset before it can be set up. The _provider_ default value is `[]`.
- `global_proxy` (Attributes) Proxy is set up directly with host, port and excluded hosts. / Android Device Owner Global Proxy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceownerglobalproxy?view=graph-rest-beta (see [below for nested schema](#nestedatt--android_device_owner_general_device--global_proxy))
- `google_accounts_blocked` (Boolean) Indicates whether or not google accounts will be blocked.
- `kiosk_customization_device_settings_blocked` (Boolean) Indicates whether a user can access the device's Settings app while in Kiosk Mode.
- `kiosk_customization_power_button_actions_blocked` (Boolean) Whether the power menu is shown when a user long presses the Power button of a device in Kiosk Mode.
- `kiosk_customization_status_bar` (String) Indicates whether system info and notifications are disabled in Kiosk Mode. / An enum representing possible values for kiosk customization system navigation; possible values are: `notConfigured` (Not configured; this value defaults to STATUS_BAR_UNSPECIFIED.), `notificationsAndSystemInfoEnabled` (System info and notifications are shown on the status bar in kiosk mode.), `systemInfoOnly` (Only system info is shown on the status bar in kiosk mode.). The _provider_ default value is `"notConfigured"`.
- `kiosk_customization_system_error_warnings` (Boolean) Indicates whether system error dialogs for crashed or unresponsive apps are shown in Kiosk Mode.
- `kiosk_customization_system_navigation` (String) Indicates which navigation features are enabled in Kiosk Mode. / An enum representing possible values for kiosk customization system navigation; possible values are: `notConfigured` (Not configured; this value defaults to NAVIGATION_DISABLED.), `navigationEnabled` (Home and overview buttons are enabled.), `homeButtonOnly` (Only the home button is enabled.). The _provider_ default value is `"notConfigured"`.
- `kiosk_mode_app_order_enabled` (Boolean) Whether or not to enable app ordering in Kiosk Mode.
- `kiosk_mode_app_positions` (Attributes Set) The ordering of items on Kiosk Mode Managed Home Screen. This collection can contain a maximum of 500 elements. / An item in the list of app positions that sets the order of items on the Managed Home Screen / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceownerkioskmodeapppositionitem?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--android_device_owner_general_device--kiosk_mode_app_positions))
- `kiosk_mode_apps` (Attributes Set) A list of managed apps that will be shown when the device is in Kiosk Mode. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--android_device_owner_general_device--kiosk_mode_apps))
- `kiosk_mode_apps_in_folder_ordered_by_name` (Boolean) Whether or not to alphabetize applications within a folder in Kiosk Mode.
- `kiosk_mode_bluetooth_configuration_enabled` (Boolean) Whether or not to allow a user to configure Bluetooth settings in Kiosk Mode.
- `kiosk_mode_debug_menu_easy_access_enabled` (Boolean) Whether or not to allow a user to easy access to the debug menu in Kiosk Mode.
- `kiosk_mode_exit_code` (String) Exit code to allow a user to escape from Kiosk Mode when the device is in Kiosk Mode.
- `kiosk_mode_flashlight_configuration_enabled` (Boolean) Whether or not to allow a user to use the flashlight in Kiosk Mode.
- `kiosk_mode_folder_icon` (String) Folder icon configuration for managed home screen in Kiosk Mode. / Android Device Owner Kiosk Mode folder icon type; possible values are: `notConfigured` (Not configured; this value is ignored.), `darkSquare` (Folder icon appears as dark square.), `darkCircle` (Folder icon appears as dark circle.), `lightSquare` (Folder icon appears as light square.), `lightCircle` (Folder icon appears as light circle  .)
- `kiosk_mode_grid_height` (Number) Number of rows for Managed Home Screen grid with app ordering enabled in Kiosk Mode. Valid values 1 to 9999999
- `kiosk_mode_grid_width` (Number) Number of columns for Managed Home Screen grid with app ordering enabled in Kiosk Mode. Valid values 1 to 9999999
- `kiosk_mode_icon_size` (String) Icon size configuration for managed home screen in Kiosk Mode. / Android Device Owner Kiosk Mode managed home screen icon size; possible values are: `notConfigured` (Not configured; this value is ignored.), `smallest` (Smallest icon size.), `small` (Small icon size.), `regular` (Regular icon size.), `large` (Large icon size.), `largest` (Largest icon size.)
- `kiosk_mode_lock_home_screen` (Boolean) Whether or not to lock home screen to the end user in Kiosk Mode.
- `kiosk_mode_managed_folders` (Attributes Set) A list of managed folders for a device in Kiosk Mode. This collection can contain a maximum of 500 elements. / A folder containing pages of apps and weblinks on the Managed Home Screen / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceownerkioskmodemanagedfolder?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--android_device_owner_general_device--kiosk_mode_managed_folders))
- `kiosk_mode_managed_home_screen_auto_signout` (Boolean) Whether or not to automatically sign-out of MHS and Shared device mode applications after inactive for Managed Home Screen.
- `kiosk_mode_managed_home_screen_inactive_sign_out_delay_in_seconds` (Number) Number of seconds to give user notice before automatically signing them out for Managed Home Screen. Valid values 0 to 9999999
- `kiosk_mode_managed_home_screen_inactive_sign_out_notice_in_seconds` (Number) Number of seconds device is inactive before automatically signing user out for Managed Home Screen. Valid values 0 to 9999999
- `kiosk_mode_managed_home_screen_pin_complexity` (String) Complexity of PIN for sign-in session for Managed Home Screen. / Complexity of PIN for Managed Home Screen sign-in session; possible values are: `notConfigured` (Not configured.), `simple` (Numeric values only.), `complex` (Alphanumerical value.)
- `kiosk_mode_managed_home_screen_pin_required` (Boolean) Whether or not require user to set a PIN for sign-in session for Managed Home Screen.
- `kiosk_mode_managed_home_screen_pin_required_to_resume` (Boolean) Whether or not required user to enter session PIN if screensaver has appeared for Managed Home Screen.
- `kiosk_mode_managed_home_screen_sign_in_background` (String) Custom URL background for sign-in screen for Managed Home Screen.
- `kiosk_mode_managed_home_screen_sign_in_branding_logo` (String) Custom URL branding logo for sign-in screen and session pin page for Managed Home Screen.
- `kiosk_mode_managed_home_screen_sign_in_enabled` (Boolean) Whether or not show sign-in screen for Managed Home Screen.
- `kiosk_mode_managed_settings_entry_disabled` (Boolean) Whether or not to display the Managed Settings entry point on the managed home screen in Kiosk Mode.
- `kiosk_mode_media_volume_configuration_enabled` (Boolean) Whether or not to allow a user to change the media volume in Kiosk Mode.
- `kiosk_mode_screen_orientation` (String) Screen orientation configuration for managed home screen in Kiosk Mode. / Android Device Owner Kiosk Mode managed home screen orientation; possible values are: `notConfigured` (Not configured; this value is ignored.), `portrait` (Portrait orientation.), `landscape` (Landscape orientation.), `autoRotate` (Auto rotate between portrait and landscape orientations.)
- `kiosk_mode_screen_saver_configuration_enabled` (Boolean) Whether or not to enable screen saver mode or not in Kiosk Mode.
- `kiosk_mode_screen_saver_detect_media_disabled` (Boolean) Whether or not the device screen should show the screen saver if audio/video is playing in Kiosk Mode.
- `kiosk_mode_screen_saver_display_time_in_seconds` (Number) The number of seconds that the device will display the screen saver for in Kiosk Mode. Valid values 0 to 9999999
- `kiosk_mode_screen_saver_image_url` (String) URL for an image that will be the device's screen saver in Kiosk Mode.
- `kiosk_mode_screen_saver_start_delay_in_seconds` (Number) The number of seconds the device needs to be inactive for before the screen saver is shown in Kiosk Mode. Valid values 1 to 9999999
- `kiosk_mode_show_app_notification_badge` (Boolean) Whether or not to display application notification badges in Kiosk Mode.
- `kiosk_mode_show_device_info` (Boolean) Whether or not to allow a user to access basic device information.
- `kiosk_mode_use_managed_home_screen_app` (String) Whether or not to use single app kiosk mode or multi-app kiosk mode. / Possible values of Android Kiosk Mode; possible values are: `notConfigured` (Not configured), `singleAppMode` (Run in single-app mode), `multiAppMode` (Run in multi-app mode). The _provider_ default value is `"notConfigured"`.
- `kiosk_mode_virtual_home_button_enabled` (Boolean) Whether or not to display a virtual home button when the device is in Kiosk Mode.
- `kiosk_mode_virtual_home_button_type` (String) Indicates whether the virtual home button is a swipe up home button or a floating home button. / Android Device Owner Kiosk Mode managed home screen virtual home button type; possible values are: `notConfigured` (Not configured; this value is ignored.), `swipeUp` (Swipe-up for home button.), `floating` (Floating home button.)
- `kiosk_mode_wallpaper_url` (String) URL to a publicly accessible image to use for the wallpaper when the device is in Kiosk Mode.
- `kiosk_mode_wifi_allowed_ssids` (Set of String) The restricted set of WIFI SSIDs available for the user to configure in Kiosk Mode. This collection can contain a maximum of 500 elements. The _provider_ default value is `[]`.
- `kiosk_mode_wifi_configuration_enabled` (Boolean) Whether or not to allow a user to configure Wi-Fi settings in Kiosk Mode.
- `locate_device_lost_mode_enabled` (Boolean) Indicates whether or not LocateDevice for devices with lost mode (COBO, COPE) is enabled.
- `locate_device_userless_disabled` (Boolean) Indicates whether or not LocateDevice for userless (COSU) devices is disabled.
- `microphone_force_mute` (Boolean) Indicates whether or not to block unmuting the microphone on the device.
- `microsoft_launcher_configuration_enabled` (Boolean) Indicates whether or not to you want configure Microsoft Launcher.
- `microsoft_launcher_custom_wallpaper_allow_user_modification` (Boolean) Indicates whether or not the user can modify the wallpaper to personalize their device.
- `microsoft_launcher_custom_wallpaper_enabled` (Boolean) Indicates whether or not to configure the wallpaper on the targeted devices.
- `microsoft_launcher_custom_wallpaper_image_url` (String) Indicates the URL for the image file to use as the wallpaper on the targeted devices.
- `microsoft_launcher_dock_presence_allow_user_modification` (Boolean) Indicates whether or not the user can modify the device dock configuration on the device.
- `microsoft_launcher_dock_presence_configuration` (String) Indicates whether or not you want to configure the device dock. / Microsoft Launcher Dock Presence selection; possible values are: `notConfigured` (Not configured; this value is ignored.), `show` (Indicates the device's dock will be displayed on the device.), `hide` (Indicates the device's dock will be hidden on the device, but the user can access the dock by dragging the handler on the bottom of the screen.), `disabled` (Indicates the device's dock will be disabled on the device.)
- `microsoft_launcher_feed_allow_user_modification` (Boolean) Indicates whether or not the user can modify the launcher feed on the device.
- `microsoft_launcher_feed_enabled` (Boolean) Indicates whether or not you want to enable the launcher feed on the device.
- `microsoft_launcher_search_bar_placement_configuration` (String) Indicates the search bar placement configuration on the device. / Microsoft Launcher Search Bar Placement selection; possible values are: `notConfigured` (Not configured; this value is ignored.), `top` (Indicates that the search bar will be displayed on the top of the device.), `bottom` (Indicates that the search bar will be displayed on the bottom of the device.), `hide` (Indicates that the search bar will be hidden on the device.)
- `network_escape_hatch_allowed` (Boolean) Indicates whether or not the device will allow connecting to a temporary network connection at boot time.
- `nfc_block_outgoing_beam` (Boolean) Indicates whether or not to block NFC outgoing beam.
- `password_block_keyguard` (Boolean) Indicates whether or not the keyguard is disabled.
- `password_block_keyguard_features` (Set of String) List of device keyguard features to block. This collection can contain a maximum of 11 elements. / Android keyguard feature; possible values are: `notConfigured` (Not configured; this value is ignored.), `camera` (Camera usage when on secure keyguard screens.), `notifications` (Showing notifications when on secure keyguard screens.), `unredactedNotifications` (Showing unredacted notifications when on secure keyguard screens.), `trustAgents` (Trust agent state when on secure keyguard screens.), `fingerprint` (Fingerprint sensor usage when on secure keyguard screens.), `remoteInput` (Notification text entry when on secure keyguard screens.), `allFeatures` (All keyguard features when on secure keyguard screens.), `face` (Face authentication on secure keyguard screens.), `iris` (Iris authentication on secure keyguard screens.), `biometrics` (All biometric authentication on secure keyguard screens.). The _provider_ default value is `[]`.
- `password_expiration_days` (Number) Indicates the amount of time that a password can be set for before it expires and a new password will be required. Valid values 1 to 365
- `password_minimum_length` (Number) Indicates the minimum length of the password required on the device. Valid values 4 to 16
- `password_minimum_letter_characters` (Number) Indicates the minimum number of letter characters required for device password. Valid values 1 to 16
- `password_minimum_lower_case_characters` (Number) Indicates the minimum number of lower case characters required for device password. Valid values 1 to 16
- `password_minimum_non_letter_characters` (Number) Indicates the minimum number of non-letter characters required for device password. Valid values 1 to 16
- `password_minimum_numeric_characters` (Number) Indicates the minimum number of numeric characters required for device password. Valid values 1 to 16
- `password_minimum_symbol_characters` (Number) Indicates the minimum number of symbol characters required for device password. Valid values 1 to 16
- `password_minimum_upper_case_characters` (Number) Indicates the minimum number of upper case letter characters required for device password. Valid values 1 to 16
- `password_minutes_of_inactivity_before_screen_timeout` (Number) Minutes of inactivity before the screen times out.
- `password_previous_password_count_to_block` (Number) Indicates the length of password history, where the user will not be able to enter a new password that is the same as any password in the history. Valid values 0 to 24
- `password_require_unlock` (String) Indicates the timeout period after which a device must be unlocked using a form of strong authentication. / An enum representing possible values for required password unlock; possible values are: `deviceDefault` (Timeout period before strong authentication is required is set to the device's default.), `daily` (Timeout period before strong authentication is required is set to 24 hours.), `unkownFutureValue` (Unknown future value (reserved, not used right now)). The _provider_ default value is `"deviceDefault"`.
- `password_required_type` (String) Indicates the minimum password quality required on the device. / Android Device Owner policy required password type; possible values are: `deviceDefault` (Device default value, no intent.), `required` (There must be a password set, but there are no restrictions on type.), `numeric` (At least numeric.), `numericComplex` (At least numeric with no repeating or ordered sequences.), `alphabetic` (At least alphabetic password.), `alphanumeric` (At least alphanumeric password), `alphanumericWithSymbols` (At least alphanumeric with symbols.), `lowSecurityBiometric` (Low security biometrics based password required.), `customPassword` (Custom password set by the admin.). The _provider_ default value is `"deviceDefault"`.
- `password_sign_in_failure_count_before_factory_reset` (Number) Indicates the number of times a user can enter an incorrect password before the device is wiped. Valid values 4 to 11
- `personal_profile_apps_allow_install_from_unknown_sources` (Boolean) Indicates whether the user can install apps from unknown sources on the personal profile.
- `personal_profile_camera_blocked` (Boolean) Indicates whether to disable the use of the camera on the personal profile.
- `personal_profile_personal_applications` (Attributes Set) Policy applied to applications in the personal profile. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--android_device_owner_general_device--personal_profile_personal_applications))
- `personal_profile_play_store_mode` (String) Used together with PersonalProfilePersonalApplications to control how apps in the personal profile are allowed or blocked. / Used together with personalApplications to control how apps in the personal profile are allowed or blocked; possible values are: `notConfigured` (Not configured.), `blockedApps` (Blocked Apps.), `allowedApps` (Allowed Apps.). The _provider_ default value is `"notConfigured"`.
- `personal_profile_screen_capture_blocked` (Boolean) Indicates whether to disable the capability to take screenshots on the personal profile.
- `play_store_mode` (String) Indicates the Play Store mode of the device. / Android Device Owner Play Store mode type; possible values are: `notConfigured` (Not Configured), `allowList` (Only apps that are in the policy are available and any app not in the policy will be automatically uninstalled from the device.), `blockList` (All apps are available and any app that should not be on the device should be explicitly marked as 'BLOCKED' in the applications policy.)
- `screen_capture_blocked` (Boolean) Indicates whether or not to disable the capability to take screenshots.
- `security_common_criteria_mode_enabled` (Boolean) Represents the security common criteria mode enabled provided to users when they attempt to modify managed settings on their device.
- `security_developer_settings_enabled` (Boolean) Indicates whether or not the user is allowed to access developer settings like developer options and safe boot on the device.
- `security_require_verify_apps` (Boolean) Indicates whether or not verify apps is required. The _provider_ default value is `true`.
- `share_device_location_disabled` (Boolean) Indicates whether or not location sharing is disabled for fully managed devices (COBO), and corporate owned devices with a work profile (COPE)
- `short_help_text` (Attributes) Represents the customized short help text provided to users when they attempt to modify managed settings on their device. / Represents a user-facing message with locale information as well as a default message to be used if the user's locale doesn't match with any of the localized messages / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceowneruserfacingmessage?view=graph-rest-beta (see [below for nested schema](#nestedatt--android_device_owner_general_device--short_help_text))
- `status_bar_blocked` (Boolean) Indicates whether or the status bar is disabled, including notifications, quick settings and other screen overlays.
- `stay_on_modes` (Set of String) List of modes in which the device's display will stay powered-on. This collection can contain a maximum of 4 elements. / Android Device Owner possible values for states of the device's plugged-in power modes; possible values are: `notConfigured` (Not configured; this value is ignored.), `ac` (Power source is an AC charger.), `usb` (Power source is a USB port.), `wireless` (Power source is wireless.). The _provider_ default value is `[]`.
- `storage_allow_usb` (Boolean) Indicates whether or not to allow USB mass storage.
- `storage_block_external_media` (Boolean) Indicates whether or not to block external media.
- `storage_block_usb_file_transfer` (Boolean) Indicates whether or not to block USB file transfer.
- `system_update_freeze_periods` (Attributes Set) Indicates the annually repeating time periods during which system updates are postponed. This collection can contain a maximum of 500 elements. / Represents one item in the list of freeze periods for Android Device Owner system updates / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceownersystemupdatefreezeperiod?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--android_device_owner_general_device--system_update_freeze_periods))
- `system_update_install_type` (String) The type of system update configuration. / System Update Types for Android Device Owner; possible values are: `deviceDefault` (Device default behavior, which typically prompts the user to accept system updates.), `postpone` (Postpone automatic install of updates up to 30 days.), `windowed` (Install automatically inside a daily maintenance window.), `automatic` (Automatically install updates as soon as possible.)
- `system_update_window_end_minutes_after_midnight` (Number) Indicates the number of minutes after midnight that the system update window ends. Valid values 0 to 1440
- `system_update_window_start_minutes_after_midnight` (Number) Indicates the number of minutes after midnight that the system update window starts. Valid values 0 to 1440
- `system_windows_blocked` (Boolean) Whether or not to block Android system prompt windows, like toasts, phone activities, and system alerts.
- `users_block_add` (Boolean) Indicates whether or not adding users and profiles is disabled.
- `users_block_remove` (Boolean) Indicates whether or not to disable removing other users from the device.
- `volume_block_adjustment` (Boolean) Indicates whether or not adjusting the master volume is disabled.
- `vpn_always_on_lockdown_mode` (Boolean) If an always on VPN package name is specified, whether or not to lock network traffic when that VPN is disconnected. The _provider_ default value is `false`.
- `vpn_always_on_package_identifier` (String) Android app package name for app that will handle an always-on VPN connection. The _provider_ default value is `""`.
- `wifi_block_edit_configurations` (Boolean) Indicates whether or not to block the user from editing the wifi connection settings.
- `wifi_block_edit_policy_defined_configurations` (Boolean) Indicates whether or not to block the user from editing just the networks defined by the policy.
- `work_profile_password_expiration_days` (Number) Indicates the number of days that a work profile password can be set before it expires and a new password will be required. Valid values 1 to 365
- `work_profile_password_minimum_length` (Number) Indicates the minimum length of the work profile password. Valid values 4 to 16
- `work_profile_password_minimum_letter_characters` (Number) Indicates the minimum number of letter characters required for the work profile password. Valid values 1 to 16
- `work_profile_password_minimum_lower_case_characters` (Number) Indicates the minimum number of lower-case characters required for the work profile password. Valid values 1 to 16
- `work_profile_password_minimum_non_letter_characters` (Number) Indicates the minimum number of non-letter characters required for the work profile password. Valid values 1 to 16
- `work_profile_password_minimum_numeric_characters` (Number) Indicates the minimum number of numeric characters required for the work profile password. Valid values 1 to 16
- `work_profile_password_minimum_symbol_characters` (Number) Indicates the minimum number of symbol characters required for the work profile password. Valid values 1 to 16
- `work_profile_password_minimum_upper_case_characters` (Number) Indicates the minimum number of upper-case letter characters required for the work profile password. Valid values 1 to 16
- `work_profile_password_previous_password_count_to_block` (Number) Indicates the length of the work profile password history, where the user will not be able to enter a new password that is the same as any password in the history. Valid values 0 to 24
- `work_profile_password_require_unlock` (String) Indicates the timeout period after which a work profile must be unlocked using a form of strong authentication. / An enum representing possible values for required password unlock; possible values are: `deviceDefault` (Timeout period before strong authentication is required is set to the device's default.), `daily` (Timeout period before strong authentication is required is set to 24 hours.), `unkownFutureValue` (Unknown future value (reserved, not used right now)). The _provider_ default value is `"deviceDefault"`.
- `work_profile_password_required_type` (String) Indicates the minimum password quality required on the work profile password. / Android Device Owner policy required password type; possible values are: `deviceDefault` (Device default value, no intent.), `required` (There must be a password set, but there are no restrictions on type.), `numeric` (At least numeric.), `numericComplex` (At least numeric with no repeating or ordered sequences.), `alphabetic` (At least alphabetic password.), `alphanumeric` (At least alphanumeric password), `alphanumericWithSymbols` (At least alphanumeric with symbols.), `lowSecurityBiometric` (Low security biometrics based password required.), `customPassword` (Custom password set by the admin.). The _provider_ default value is `"deviceDefault"`.
- `work_profile_password_sign_in_failure_count_before_factory_reset` (Number) Indicates the number of times a user can enter an incorrect work profile password before the device is wiped. Valid values 4 to 11

<a id="nestedatt--android_device_owner_general_device--android_device_owner_delegated_scope_app_settings"></a>
### Nested Schema for `android_device_owner_general_device.android_device_owner_delegated_scope_app_settings`

Optional:

- `app_detail` (Attributes) Information about the app like Name, AppStoreUrl, Publisher and AppId / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `{}`. (see [below for nested schema](#nestedatt--android_device_owner_general_device--android_device_owner_delegated_scope_app_settings--app_detail))
- `app_scopes` (Set of String) List of scopes an app has been assigned. / An enum representing possible values for delegated app scope; possible values are: `unspecified` (Unspecified; this value defaults to DELEGATED_SCOPE_UNSPECIFIED.), `certificateInstall` (Specifies that the admin has given app permission to install and manage certificates on device.), `captureNetworkActivityLog` (Specifies that the admin has given app permission to capture network activity logs on device. More info on Network activity logs: https://developer.android.com/work/dpc/logging), `captureSecurityLog` (Specified that the admin has given permission to capture security logs on device. More info on Security logs: https://developer.android.com/work/dpc/security#log_enterprise_device_activity), `unknownFutureValue` (Unknown future value (reserved, not used right now)). The _provider_ default value is `[]`.

<a id="nestedatt--android_device_owner_general_device--android_device_owner_delegated_scope_app_settings--app_detail"></a>
### Nested Schema for `android_device_owner_general_device.android_device_owner_delegated_scope_app_settings.app_detail`

Required:

- `name` (String) The application name

Optional:

- `app_id` (String) The application or bundle identifier of the application
- `app_store_url` (String) The Store URL of the application
- `publisher` (String) The publisher of the application



<a id="nestedatt--android_device_owner_general_device--azure_ad_shared_device_data_clear_apps"></a>
### Nested Schema for `android_device_owner_general_device.azure_ad_shared_device_data_clear_apps`

Required:

- `name` (String) The application name

Optional:

- `app_id` (String) The application or bundle identifier of the application
- `app_store_url` (String) The Store URL of the application
- `publisher` (String) The publisher of the application


<a id="nestedatt--android_device_owner_general_device--detailed_help_text"></a>
### Nested Schema for `android_device_owner_general_device.detailed_help_text`

Required:

- `default_message` (String) The default message displayed if the user's locale doesn't match with any of the localized messages

Optional:

- `localized_messages` (Attributes Set) The list of <locale, message> pairs. This collection can contain a maximum of 500 elements. / Key value pair for storing custom settings / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--android_device_owner_general_device--detailed_help_text--localized_messages))

<a id="nestedatt--android_device_owner_general_device--detailed_help_text--localized_messages"></a>
### Nested Schema for `android_device_owner_general_device.detailed_help_text.localized_messages`

Required:

- `name` (String) Name for this key-value pair

Optional:

- `value` (String) Value for this key-value pair



<a id="nestedatt--android_device_owner_general_device--device_owner_lock_screen_message"></a>
### Nested Schema for `android_device_owner_general_device.device_owner_lock_screen_message`

Required:

- `default_message` (String) The default message displayed if the user's locale doesn't match with any of the localized messages

Optional:

- `localized_messages` (Attributes Set) The list of <locale, message> pairs. This collection can contain a maximum of 500 elements. / Key value pair for storing custom settings / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--android_device_owner_general_device--device_owner_lock_screen_message--localized_messages))

<a id="nestedatt--android_device_owner_general_device--device_owner_lock_screen_message--localized_messages"></a>
### Nested Schema for `android_device_owner_general_device.device_owner_lock_screen_message.localized_messages`

Required:

- `name` (String) Name for this key-value pair

Optional:

- `value` (String) Value for this key-value pair



<a id="nestedatt--android_device_owner_general_device--global_proxy"></a>
### Nested Schema for `android_device_owner_general_device.global_proxy`


<a id="nestedatt--android_device_owner_general_device--kiosk_mode_app_positions"></a>
### Nested Schema for `android_device_owner_general_device.kiosk_mode_app_positions`

Optional:

- `item` (Attributes) Item to be arranged / Represents an item on the Android Device Owner Managed Home Screen (application, weblink or folder / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceownerkioskmodehomescreenitem?view=graph-rest-beta. The _provider_ default value is `{}`. (see [below for nested schema](#nestedatt--android_device_owner_general_device--kiosk_mode_app_positions--item))
- `position` (Number) Position of the item on the grid. Valid values 0 to 9999999. The _provider_ default value is `0`.

<a id="nestedatt--android_device_owner_general_device--kiosk_mode_app_positions--item"></a>
### Nested Schema for `android_device_owner_general_device.kiosk_mode_app_positions.item`



<a id="nestedatt--android_device_owner_general_device--kiosk_mode_apps"></a>
### Nested Schema for `android_device_owner_general_device.kiosk_mode_apps`

Required:

- `name` (String) The application name

Optional:

- `app_id` (String) The application or bundle identifier of the application
- `app_store_url` (String) The Store URL of the application
- `publisher` (String) The publisher of the application


<a id="nestedatt--android_device_owner_general_device--kiosk_mode_managed_folders"></a>
### Nested Schema for `android_device_owner_general_device.kiosk_mode_managed_folders`

Required:

- `folder_identifier` (String) Unique identifier for the folder
- `folder_name` (String) Display name for the folder

Optional:

- `items` (Attributes Set) Items to be added to managed folder. This collection can contain a maximum of 500 elements. / Represents an item that can be added to Android Device Owner folder (application or weblink) / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-androiddeviceownerkioskmodefolderitem?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--android_device_owner_general_device--kiosk_mode_managed_folders--items))

<a id="nestedatt--android_device_owner_general_device--kiosk_mode_managed_folders--items"></a>
### Nested Schema for `android_device_owner_general_device.kiosk_mode_managed_folders.items`



<a id="nestedatt--android_device_owner_general_device--personal_profile_personal_applications"></a>
### Nested Schema for `android_device_owner_general_device.personal_profile_personal_applications`

Required:

- `name` (String) The application name

Optional:

- `app_id` (String) The application or bundle identifier of the application
- `app_store_url` (String) The Store URL of the application
- `publisher` (String) The publisher of the application


<a id="nestedatt--android_device_owner_general_device--short_help_text"></a>
### Nested Schema for `android_device_owner_general_device.short_help_text`

Required:

- `default_message` (String) The default message displayed if the user's locale doesn't match with any of the localized messages

Optional:

- `localized_messages` (Attributes Set) The list of <locale, message> pairs. This collection can contain a maximum of 500 elements. / Key value pair for storing custom settings / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--android_device_owner_general_device--short_help_text--localized_messages))

<a id="nestedatt--android_device_owner_general_device--short_help_text--localized_messages"></a>
### Nested Schema for `android_device_owner_general_device.short_help_text.localized_messages`

Required:

- `name` (String) Name for this key-value pair

Optional:

- `value` (String) Value for this key-value pair



<a id="nestedatt--android_device_owner_general_device--system_update_freeze_periods"></a>
### Nested Schema for `android_device_owner_general_device.system_update_freeze_periods`

Optional:

- `end_day` (Number) The day of the end date of the freeze period. Valid values 1 to 31. The _provider_ default value is `0`.
- `end_month` (Number) The month of the end date of the freeze period. Valid values 1 to 12. The _provider_ default value is `0`.
- `start_day` (Number) The day of the start date of the freeze period. Valid values 1 to 31. The _provider_ default value is `0`.
- `start_month` (Number) The month of the start date of the freeze period. Valid values 1 to 12. The _provider_ default value is `0`.



<a id="nestedatt--android_work_profile_general_device"></a>
### Nested Schema for `android_work_profile_general_device`

Optional:

- `allowed_google_account_domains` (Set of String) Determine domains allow-list for accounts that can be added to work profile. The _provider_ default value is `[]`.
- `block_unified_password_for_work_profile` (Boolean) Prevent using unified password for unlocking device and work profile. The _provider_ default value is `false`.
- `password_block_face_unlock` (Boolean) Indicates whether or not to block face unlock. The _provider_ default value is `false`.
- `password_block_fingerprint_unlock` (Boolean) Indicates whether or not to block fingerprint unlock. The _provider_ default value is `false`.
- `password_block_iris_unlock` (Boolean) Indicates whether or not to block iris unlock. The _provider_ default value is `false`.
- `password_block_trust_agents` (Boolean) Indicates whether or not to block Smart Lock and other trust agents. The _provider_ default value is `false`.
- `password_expiration_days` (Number) Number of days before the password expires. Valid values 1 to 365
- `password_minimum_length` (Number) Minimum length of passwords. Valid values 4 to 16
- `password_minutes_of_inactivity_before_screen_timeout` (Number) Minutes of inactivity before the screen times out.
- `password_previous_password_block_count` (Number) Number of previous passwords to block. Valid values 0 to 24
- `password_required_type` (String) Type of password that is required. / Android Work Profile required password type; possible values are: `deviceDefault` (Device default value, no intent.), `lowSecurityBiometric` (Low security biometrics based password required.), `required` (Required.), `atLeastNumeric` (At least numeric password required.), `numericComplex` (Numeric complex password required.), `atLeastAlphabetic` (At least alphabetic password required.), `atLeastAlphanumeric` (At least alphanumeric password required.), `alphanumericWithSymbols` (At least alphanumeric with symbols password required.). The _provider_ default value is `"atLeastNumeric"`.
- `password_sign_in_failure_count_before_factory_reset` (Number) Number of sign in failures allowed before factory reset. Valid values 1 to 16
- `required_password_complexity` (String) Indicates the required device password complexity on Android. One of: NONE, LOW, MEDIUM, HIGH. This is a new API targeted to Android 12+. / The password complexity types that can be set on Android. One of: NONE, LOW, MEDIUM, HIGH. This is an API targeted to Android 11+; possible values are: `none` (Device default value, no password.), `low` (The required password complexity on the device is of type low as defined by the Android documentation.), `medium` (The required password complexity on the device is of type medium as defined by the Android documentation.), `high` (The required password complexity on the device is of type high as defined by the Android documentation.). The _provider_ default value is `"none"`.
- `security_require_verify_apps` (Boolean) Require the Android Verify apps feature is turned on. The _provider_ default value is `false`.
- `vpn_always_on_package_identifier` (String) Enable lockdown mode for always-on VPN.
- `vpn_enable_always_on_lockdown_mode` (Boolean) Enable lockdown mode for always-on VPN. The _provider_ default value is `false`.
- `work_profile_account_use` (String) Control user's ability to add accounts in work profile including Google accounts. / An enum representing possible values for account use in work profile; possible values are: `allowAllExceptGoogleAccounts` (Allow additon of all accounts except Google accounts in Android Work Profile.), `blockAll` (Block any account from being added in Android Work Profile.), `allowAll` (Allow addition of all accounts (including Google accounts) in Android Work Profile.), `unknownFutureValue` (Unknown future value for evolvable enum patterns.). The _provider_ default value is `"allowAllExceptGoogleAccounts"`.
- `work_profile_allow_app_installs_from_unknown_sources` (Boolean) Indicates whether to allow installation of apps from unknown sources. The _provider_ default value is `false`.
- `work_profile_allow_widgets` (Boolean) Allow widgets from work profile apps. The _provider_ default value is `false`.
- `work_profile_block_adding_accounts` (Boolean) Block users from adding/removing accounts in work profile. The _provider_ default value is `false`.
- `work_profile_block_camera` (Boolean) Block work profile camera. The _provider_ default value is `false`.
- `work_profile_block_cross_profile_caller_id` (Boolean) Block display work profile caller ID in personal profile. The _provider_ default value is `false`.
- `work_profile_block_cross_profile_contacts_search` (Boolean) Block work profile contacts availability in personal profile. The _provider_ default value is `false`.
- `work_profile_block_cross_profile_copy_paste` (Boolean) Boolean that indicates if the setting disallow cross profile copy/paste is enabled. The _provider_ default value is `true`.
- `work_profile_block_notifications_while_device_locked` (Boolean) Indicates whether or not to block notifications while device locked. The _provider_ default value is `false`.
- `work_profile_block_personal_app_installs_from_unknown_sources` (Boolean) Prevent app installations from unknown sources in the personal profile. The _provider_ default value is `false`.
- `work_profile_block_screen_capture` (Boolean) Block screen capture in work profile. The _provider_ default value is `false`.
- `work_profile_bluetooth_enable_contact_sharing` (Boolean) Allow bluetooth devices to access enterprise contacts. The _provider_ default value is `false`.
- `work_profile_data_sharing_type` (String) Type of data sharing that is allowed. / Android Work Profile cross profile data sharing type; possible values are: `deviceDefault` (Device default value, no intent.), `preventAny` (Prevent any sharing.), `allowPersonalToWork` (Allow data sharing request from personal profile to work profile.), `noRestrictions` (No restrictions on sharing.). The _provider_ default value is `"deviceDefault"`.
- `work_profile_default_app_permission_policy` (String) Type of password that is required. / Android Work Profile default app permission policy type; possible values are: `deviceDefault` (Device default value, no intent.), `prompt` (Prompt.), `autoGrant` (Auto grant.), `autoDeny` (Auto deny.). The _provider_ default value is `"deviceDefault"`.
- `work_profile_password_block_face_unlock` (Boolean) Indicates whether or not to block face unlock for work profile. The _provider_ default value is `false`.
- `work_profile_password_block_fingerprint_unlock` (Boolean) Indicates whether or not to block fingerprint unlock for work profile. The _provider_ default value is `false`.
- `work_profile_password_block_iris_unlock` (Boolean) Indicates whether or not to block iris unlock for work profile. The _provider_ default value is `false`.
- `work_profile_password_block_trust_agents` (Boolean) Indicates whether or not to block Smart Lock and other trust agents for work profile. The _provider_ default value is `false`.
- `work_profile_password_expiration_days` (Number) Number of days before the work profile password expires. Valid values 1 to 365
- `work_profile_password_min_letter_characters` (Number) Minimum # of letter characters required in work profile password. Valid values 1 to 10
- `work_profile_password_min_lower_case_characters` (Number) Minimum # of lower-case characters required in work profile password. Valid values 1 to 10
- `work_profile_password_min_non_letter_characters` (Number) Minimum # of non-letter characters required in work profile password. Valid values 1 to 10
- `work_profile_password_min_numeric_characters` (Number) Minimum # of numeric characters required in work profile password. Valid values 1 to 10
- `work_profile_password_min_symbol_characters` (Number) Minimum # of symbols required in work profile password. Valid values 1 to 10
- `work_profile_password_min_upper_case_characters` (Number) Minimum # of upper-case characters required in work profile password. Valid values 1 to 10
- `work_profile_password_minimum_length` (Number) Minimum length of work profile password. Valid values 4 to 16
- `work_profile_password_minutes_of_inactivity_before_screen_timeout` (Number) Minutes of inactivity before the screen times out.
- `work_profile_password_previous_password_block_count` (Number) Number of previous work profile passwords to block. Valid values 0 to 24
- `work_profile_password_required_type` (String) Type of work profile password that is required. / Android Work Profile required password type; possible values are: `deviceDefault` (Device default value, no intent.), `lowSecurityBiometric` (Low security biometrics based password required.), `required` (Required.), `atLeastNumeric` (At least numeric password required.), `numericComplex` (Numeric complex password required.), `atLeastAlphabetic` (At least alphabetic password required.), `atLeastAlphanumeric` (At least alphanumeric password required.), `alphanumericWithSymbols` (At least alphanumeric with symbols password required.). The _provider_ default value is `"deviceDefault"`.
- `work_profile_password_sign_in_failure_count_before_factory_reset` (Number) Number of sign in failures allowed before work profile is removed and all corporate data deleted. Valid values 1 to 16
- `work_profile_require_password` (Boolean) Password is required or not for work profile. The _provider_ default value is `false`.
- `work_profile_required_password_complexity` (String) Indicates the required work profile password complexity on Android. One of: NONE, LOW, MEDIUM, HIGH. This is a new API targeted to Android 12+. / The password complexity types that can be set on Android. One of: NONE, LOW, MEDIUM, HIGH. This is an API targeted to Android 11+; possible values are: `none` (Device default value, no password.), `low` (The required password complexity on the device is of type low as defined by the Android documentation.), `medium` (The required password complexity on the device is of type medium as defined by the Android documentation.), `high` (The required password complexity on the device is of type high as defined by the Android documentation.). The _provider_ default value is `"none"`.


<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `target` (Attributes) Base type for assignment targets. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-deviceandappmanagementassignmenttarget?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--target))

<a id="nestedatt--assignments--target"></a>
### Nested Schema for `assignments.target`

Optional:

- `all_devices` (Attributes) Represents an assignment to all managed devices in the tenant. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-alldevicesassignmenttarget?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--target--all_devices))
- `all_licensed_users` (Attributes) Represents an assignment to all licensed users in the tenant. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-alllicensedusersassignmenttarget?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--target--all_licensed_users))
- `exclusion_group` (Attributes) Represents a group that should be excluded from an assignment. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-exclusiongroupassignmenttarget?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--target--exclusion_group))
- `filter_id` (String) The ID of the filter for the target assignment.
- `filter_type` (String) The type of filter of the target assignment i.e. Exclude or Include. / Represents type of the assignment filter; possible values are: `none` (Default value. Do not use.), `include` (Indicates in-filter, rule matching will offer the payload to devices.), `exclude` (Indicates out-filter, rule matching will not offer the payload to devices.). The _provider_ default value is `"none"`.
- `group` (Attributes) Represents an assignment to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-groupassignmenttarget?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--target--group))

<a id="nestedatt--assignments--target--all_devices"></a>
### Nested Schema for `assignments.target.all_devices`


<a id="nestedatt--assignments--target--all_licensed_users"></a>
### Nested Schema for `assignments.target.all_licensed_users`


<a id="nestedatt--assignments--target--exclusion_group"></a>
### Nested Schema for `assignments.target.exclusion_group`

Required:

- `group_id` (String) The group Id that is the target of the assignment.


<a id="nestedatt--assignments--target--group"></a>
### Nested Schema for `assignments.target.group`

Required:

- `group_id` (String) The group Id that is the target of the assignment.




<a id="nestedatt--device_management_applicability_rule_device_mode"></a>
### Nested Schema for `device_management_applicability_rule_device_mode`

Required:

- `device_mode` (String) Applicability rule for device mode. / Windows 10 Device Mode type; possible values are: `standardConfiguration` (Standard Configuration), `sModeConfiguration` (S Mode Configuration)
- `rule_type` (String) Applicability Rule type. / Supported Applicability rule types for Device Configuration; possible values are: `include` (Include), `exclude` (Exclude)

Optional:

- `name` (String) Name for object.


<a id="nestedatt--device_management_applicability_rule_os_edition"></a>
### Nested Schema for `device_management_applicability_rule_os_edition`

Required:

- `os_edition_types` (Set of String) Applicability rule OS edition type. / Windows 10 Edition type; possible values are: `windows10Enterprise` (Windows 10 Enterprise), `windows10EnterpriseN` (Windows 10 EnterpriseN), `windows10Education` (Windows 10 Education), `windows10EducationN` (Windows 10 EducationN), `windows10MobileEnterprise` (Windows 10 Mobile Enterprise), `windows10HolographicEnterprise` (Windows 10 Holographic Enterprise), `windows10Professional` (Windows 10 Professional), `windows10ProfessionalN` (Windows 10 ProfessionalN), `windows10ProfessionalEducation` (Windows 10 Professional Education), `windows10ProfessionalEducationN` (Windows 10 Professional EducationN), `windows10ProfessionalWorkstation` (Windows 10 Professional for Workstations), `windows10ProfessionalWorkstationN` (Windows 10 Professional for Workstations N), `notConfigured` (NotConfigured), `windows10Home` (Windows 10 Home), `windows10HomeChina` (Windows 10 Home China), `windows10HomeN` (Windows 10 Home N), `windows10HomeSingleLanguage` (Windows 10 Home Single Language), `windows10Mobile` (Windows 10 Mobile), `windows10IoTCore` (Windows 10 IoT Core), `windows10IoTCoreCommercial` (Windows 10 IoT Core Commercial)
- `rule_type` (String) Applicability Rule type. / Supported Applicability rule types for Device Configuration; possible values are: `include` (Include), `exclude` (Exclude)

Optional:

- `name` (String) Name for object.


<a id="nestedatt--device_management_applicability_rule_os_version"></a>
### Nested Schema for `device_management_applicability_rule_os_version`

Required:

- `rule_type` (String) Applicability Rule type. / Supported Applicability rule types for Device Configuration; possible values are: `include` (Include), `exclude` (Exclude)

Optional:

- `max_os_version` (String) Max OS version for Applicability Rule.
- `min_os_version` (String) Min OS version for Applicability Rule.
- `name` (String) Name for object.


<a id="nestedatt--edition_upgrade"></a>
### Nested Schema for `edition_upgrade`

Optional:

- `license` (String) Edition Upgrade License File Content. The _provider_ default value is `""`.
- `license_type` (String) Edition Upgrade License Type. / Edition Upgrade License type; possible values are: `productKey` (Product Key Type), `licenseFile` (License File Type), `notConfigured` (NotConfigured). The _provider_ default value is `"notConfigured"`.
- `product_key` (String) Edition Upgrade Product Key. The _provider_ default value is `""`.  
Provider Note: MS Graph returns this attribute with most of the digits replaced by `#`. Therefore only the remaining visible digits will be compared to ensure that state matches plan, i.e. one of the visible digits must change for the entity to get updated in MS Graph.
- `target_edition` (String) Edition Upgrade Target Edition. / Windows 10 Edition type; possible values are: `windows10Enterprise` (Windows 10 Enterprise), `windows10EnterpriseN` (Windows 10 EnterpriseN), `windows10Education` (Windows 10 Education), `windows10EducationN` (Windows 10 EducationN), `windows10MobileEnterprise` (Windows 10 Mobile Enterprise), `windows10HolographicEnterprise` (Windows 10 Holographic Enterprise), `windows10Professional` (Windows 10 Professional), `windows10ProfessionalN` (Windows 10 ProfessionalN), `windows10ProfessionalEducation` (Windows 10 Professional Education), `windows10ProfessionalEducationN` (Windows 10 Professional EducationN), `windows10ProfessionalWorkstation` (Windows 10 Professional for Workstations), `windows10ProfessionalWorkstationN` (Windows 10 Professional for Workstations N), `notConfigured` (NotConfigured), `windows10Home` (Windows 10 Home), `windows10HomeChina` (Windows 10 Home China), `windows10HomeN` (Windows 10 Home N), `windows10HomeSingleLanguage` (Windows 10 Home Single Language), `windows10Mobile` (Windows 10 Mobile), `windows10IoTCore` (Windows 10 IoT Core), `windows10IoTCoreCommercial` (Windows 10 IoT Core Commercial). The _provider_ default value is `"notConfigured"`.
- `windows_s_mode` (String) S mode configuration. / The possible options to configure S mode unlock; possible values are: `noRestriction` (This option will remove all restrictions to unlock S mode - default), `block` (This option will block the user to unlock the device from S mode), `unlock` (This option will unlock the device from S mode). The _provider_ default value is `"noRestriction"`.


<a id="nestedatt--ios_custom"></a>
### Nested Schema for `ios_custom`

Required:

- `file_name` (String) Payload file name (*.mobileconfig
- `name` (String) Name that is displayed to the user.
- `payload` (String) Payload. (UTF8 encoded byte array)


<a id="nestedatt--ios_device_features"></a>
### Nested Schema for `ios_device_features`

Optional:

- `airprint_destinations` (Attributes Set) An array of AirPrint printers that should always be shown. This collection can contain a maximum of 500 elements. / Represents an AirPrint destination. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-airprintdestination?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--ios_device_features--airprint_destinations))
- `asset_tag_template` (String) Asset tag information for the device, displayed on the login window and lock screen.
- `content_filter_settings` (Attributes) Gets or sets iOS Web Content Filter settings, supervised mode only / Represents an iOS Web Content Filter setting base type. An empty and abstract base. Caller should use one of derived types for configurations. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioswebcontentfilterbase?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--content_filter_settings))
- `home_screen_dock_icons` (Attributes Set) A list of app and folders to appear on the Home Screen Dock. This collection can contain a maximum of 500 elements. / Represents an item on the iOS Home Screen / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioshomescreenitem?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--ios_device_features--home_screen_dock_icons))
- `home_screen_grid_height` (Number) Gets or sets the number of rows to render when configuring iOS home screen layout settings. If this value is configured, homeScreenGridWidth must be configured as well.
- `home_screen_grid_width` (Number) Gets or sets the number of columns to render when configuring iOS home screen layout settings. If this value is configured, homeScreenGridHeight must be configured as well.
- `home_screen_pages` (Attributes Set) A list of pages on the Home Screen. This collection can contain a maximum of 500 elements. / A page containing apps, folders, and web clips on the Home Screen. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioshomescreenpage?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--ios_device_features--home_screen_pages))
- `ios_single_sign_on_extension` (Attributes) Gets or sets a single sign-on extension profile. / An abstract base class for all iOS-specific single sign-on extension types. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iossinglesignonextension?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension))
- `lock_screen_footnote` (String) A footnote displayed on the login window and lock screen. Available in iOS 9.3.1 and later.
- `notification_settings` (Attributes Set) Notification settings for each bundle id. Applicable to devices in supervised mode only (iOS 9.3 and later). This collection can contain a maximum of 500 elements. / An item describing notification setting. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosnotificationsettings?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--ios_device_features--notification_settings))
- `single_sign_on_settings` (Attributes) The Kerberos login settings that enable apps on receiving devices to authenticate smoothly. / iOS Kerberos authentication settings for single sign-on / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iossinglesignonsettings?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--single_sign_on_settings))
- `wallpaper_display_location` (String) A wallpaper display location specifier. / An enum type for wallpaper display location specifier; possible values are: `notConfigured` (No location specified for wallpaper display.), `lockScreen` (A configured wallpaper image is displayed on Lock screen.), `homeScreen` (A configured wallpaper image is displayed on Home (icon list) screen.), `lockAndHomeScreens` (A configured wallpaper image is displayed on Lock screen and Home screen.). The _provider_ default value is `"notConfigured"`.
- `wallpaper_image` (Attributes) A wallpaper image must be in either PNG or JPEG format. It requires a supervised device with iOS 8 or later version. / Contains properties for a generic mime content. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mimecontent?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--wallpaper_image))

<a id="nestedatt--ios_device_features--airprint_destinations"></a>
### Nested Schema for `ios_device_features.airprint_destinations`

Required:

- `force_tls` (Boolean) If true AirPrint connections are secured by Transport Layer Security (TLS). Default is false. Available in iOS 11.0 and later.
- `ip_address` (String) The IP Address of the AirPrint destination.
- `resource_path` (String) The Resource Path associated with the printer. This corresponds to the rp parameter of the _ipps.tcp Bonjour record. For example: printers/Canon_MG5300_series, printers/Xerox_Phaser_7600, ipp/print, Epson_IPP_Printer.

Optional:

- `port` (Number) The listening port of the AirPrint destination. If this key is not specified AirPrint will use the default port. Available in iOS 11.0 and later.


<a id="nestedatt--ios_device_features--content_filter_settings"></a>
### Nested Schema for `ios_device_features.content_filter_settings`

Optional:

- `auto_filter` (Attributes) Represents an iOS Web Content Filter setting type, which enables iOS automatic filter feature and allows for additional URL access control. When constructed with no property values, the iOS device will enable the automatic filter regardless. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioswebcontentfilterautofilter?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--content_filter_settings--auto_filter))
- `specific_websites_access` (Attributes) Represents an iOS Web Content Filter setting type, which installs URL bookmarks into iOS built-in browser. An example scenario is in the classroom where teachers would like the students to navigate websites through browser bookmarks configured on their iOS devices, and no access to other sites. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioswebcontentfilterspecificwebsitesaccess?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--content_filter_settings--specific_websites_access))

<a id="nestedatt--ios_device_features--content_filter_settings--auto_filter"></a>
### Nested Schema for `ios_device_features.content_filter_settings.auto_filter`

Optional:

- `allowed_urls` (Set of String) Additional URLs allowed for access. The _provider_ default value is `[]`.
- `blocked_urls` (Set of String) Additional URLs blocked for access. The _provider_ default value is `[]`.


<a id="nestedatt--ios_device_features--content_filter_settings--specific_websites_access"></a>
### Nested Schema for `ios_device_features.content_filter_settings.specific_websites_access`

Optional:

- `specific_websites_only` (Attributes Set) URL bookmarks which will be installed into built-in browser and user is only allowed to access websites through bookmarks. This collection can contain a maximum of 500 elements. / iOS URL bookmark / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosbookmark?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--ios_device_features--content_filter_settings--specific_websites_access--specific_websites_only))
- `website_list` (Attributes Set) URL bookmarks which will be installed into built-in browser and user is only allowed to access websites through bookmarks. This collection can contain a maximum of 500 elements. / iOS URL bookmark / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosbookmark?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--ios_device_features--content_filter_settings--specific_websites_access--website_list))

<a id="nestedatt--ios_device_features--content_filter_settings--specific_websites_access--specific_websites_only"></a>
### Nested Schema for `ios_device_features.content_filter_settings.specific_websites_access.specific_websites_only`

Optional:

- `bookmark_folder` (String) The folder into which the bookmark should be added in Safari
- `display_name` (String) The display name of the bookmark. The _provider_ default value is `""`.
- `url` (String) URL allowed to access. The _provider_ default value is `""`.


<a id="nestedatt--ios_device_features--content_filter_settings--specific_websites_access--website_list"></a>
### Nested Schema for `ios_device_features.content_filter_settings.specific_websites_access.website_list`

Optional:

- `bookmark_folder` (String) The folder into which the bookmark should be added in Safari
- `display_name` (String) The display name of the bookmark. The _provider_ default value is `""`.
- `url` (String) URL allowed to access. The _provider_ default value is `""`.




<a id="nestedatt--ios_device_features--home_screen_dock_icons"></a>
### Nested Schema for `ios_device_features.home_screen_dock_icons`

Optional:

- `display_name` (String) Name of the app


<a id="nestedatt--ios_device_features--home_screen_pages"></a>
### Nested Schema for `ios_device_features.home_screen_pages`

Optional:

- `display_name` (String) Name of the page
- `icons` (Attributes Set) A list of apps, folders, and web clips to appear on a page. This collection can contain a maximum of 500 elements. / Represents an item on the iOS Home Screen / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioshomescreenitem?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--ios_device_features--home_screen_pages--icons))

<a id="nestedatt--ios_device_features--home_screen_pages--icons"></a>
### Nested Schema for `ios_device_features.home_screen_pages.icons`

Optional:

- `display_name` (String) Name of the app



<a id="nestedatt--ios_device_features--ios_single_sign_on_extension"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension`

Optional:

- `azure_ad` (Attributes) Represents an Azure AD-type Single Sign-On extension profile for iOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosazureadsinglesignonextension?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension--azure_ad))
- `credential` (Attributes) Represents a Credential-type Single Sign-On extension profile for iOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioscredentialsinglesignonextension?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension--credential))
- `kerberos` (Attributes) Represents a Kerberos-type Single Sign-On extension profile for iOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-ioskerberossinglesignonextension?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension--kerberos))
- `redirect` (Attributes) Represents a Redirect-type Single Sign-On extension profile for iOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosredirectsinglesignonextension?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension--redirect))

<a id="nestedatt--ios_device_features--ios_single_sign_on_extension--azure_ad"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension.azure_ad`

Optional:

- `bundle_id_access_control_list` (Set of String) An optional list of additional bundle IDs allowed to use the AAD extension for single sign-on. The _provider_ default value is `[]`.
- `configurations` (Attributes Set) Gets or sets a list of typed key-value pairs used to configure Credential-type profiles. This collection can contain a maximum of 500 elements. / A key-value pair with a string key and a typed value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keytypedvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension--azure_ad--configurations))
- `enable_shared_device_mode` (Boolean) Enables or disables shared device mode. The _provider_ default value is `false`.

<a id="nestedatt--ios_device_features--ios_single_sign_on_extension--azure_ad--configurations"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension.azure_ad.configurations`

Required:

- `key` (String) The string key of the key-value pair.

Optional:

- `boolean` (Attributes) A key-value pair with a string key and a Boolean value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keybooleanvaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension--azure_ad--configurations--boolean))
- `integer` (Attributes) A key-value pair with a string key and an integer value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyintegervaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension--azure_ad--configurations--integer))
- `real` (Attributes) A key-value pair with a string key and a real (floating-point) value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyrealvaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension--azure_ad--configurations--real))
- `string` (Attributes) A key-value pair with a string key and a string value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keystringvaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension--azure_ad--configurations--string))

<a id="nestedatt--ios_device_features--ios_single_sign_on_extension--azure_ad--configurations--boolean"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension.azure_ad.configurations.boolean`

Required:

- `value` (Boolean) The Boolean value of the key-value pair.


<a id="nestedatt--ios_device_features--ios_single_sign_on_extension--azure_ad--configurations--integer"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension.azure_ad.configurations.integer`

Required:

- `value` (Number) The integer value of the key-value pair.


<a id="nestedatt--ios_device_features--ios_single_sign_on_extension--azure_ad--configurations--real"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension.azure_ad.configurations.real`

Required:

- `value` (Number) The real (floating-point) value of the key-value pair.


<a id="nestedatt--ios_device_features--ios_single_sign_on_extension--azure_ad--configurations--string"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension.azure_ad.configurations.string`

Required:

- `value` (String) The string value of the key-value pair.




<a id="nestedatt--ios_device_features--ios_single_sign_on_extension--credential"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension.credential`

Required:

- `extension_identifier` (String) Gets or sets the bundle ID of the app extension that performs SSO for the specified URLs.
- `realm` (String) Gets or sets the case-sensitive realm name for this profile.
- `team_identifier` (String) Gets or sets the team ID of the app extension that performs SSO for the specified URLs.

Optional:

- `configurations` (Attributes Set) Gets or sets a list of typed key-value pairs used to configure Credential-type profiles. This collection can contain a maximum of 500 elements. / A key-value pair with a string key and a typed value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keytypedvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension--credential--configurations))
- `domains` (Set of String) Gets or sets a list of hosts or domain names for which the app extension performs SSO. The _provider_ default value is `[]`.

<a id="nestedatt--ios_device_features--ios_single_sign_on_extension--credential--configurations"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension.credential.configurations`

Required:

- `key` (String) The string key of the key-value pair.

Optional:

- `boolean` (Attributes) A key-value pair with a string key and a Boolean value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keybooleanvaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension--credential--configurations--boolean))
- `integer` (Attributes) A key-value pair with a string key and an integer value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyintegervaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension--credential--configurations--integer))
- `real` (Attributes) A key-value pair with a string key and a real (floating-point) value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyrealvaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension--credential--configurations--real))
- `string` (Attributes) A key-value pair with a string key and a string value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keystringvaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension--credential--configurations--string))

<a id="nestedatt--ios_device_features--ios_single_sign_on_extension--credential--configurations--boolean"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension.credential.configurations.boolean`

Required:

- `value` (Boolean) The Boolean value of the key-value pair.


<a id="nestedatt--ios_device_features--ios_single_sign_on_extension--credential--configurations--integer"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension.credential.configurations.integer`

Required:

- `value` (Number) The integer value of the key-value pair.


<a id="nestedatt--ios_device_features--ios_single_sign_on_extension--credential--configurations--real"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension.credential.configurations.real`

Required:

- `value` (Number) The real (floating-point) value of the key-value pair.


<a id="nestedatt--ios_device_features--ios_single_sign_on_extension--credential--configurations--string"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension.credential.configurations.string`

Required:

- `value` (String) The string value of the key-value pair.




<a id="nestedatt--ios_device_features--ios_single_sign_on_extension--kerberos"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension.kerberos`

Required:

- `block_active_directory_site_auto_discovery` (Boolean) Enables or disables whether the Kerberos extension can automatically determine its site name.
- `block_automatic_login` (Boolean) Enables or disables Keychain usage.
- `is_default_realm` (Boolean) When true, this profile's realm will be selected as the default. Necessary if multiple Kerberos-type profiles are configured.
- `managed_apps_in_bundle_id_acl_included` (Boolean) When set to True, the Kerberos extension allows managed apps, and any apps entered with the app bundle ID to access the credential. When set to False, the Kerberos extension allows all apps to access the credential. Available for devices running iOS and iPadOS versions 14 and later.
- `password_block_modification` (Boolean) Enables or disables password changes.
- `password_enable_local_sync` (Boolean) Enables or disables password syncing. This won't affect users logged in with a mobile account on macOS.
- `password_require_active_directory_complexity` (Boolean) Enables or disables whether passwords must meet Active Directory's complexity requirements.
- `realm` (String) Gets or sets the case-sensitive realm name for this profile.
- `require_user_presence` (Boolean) Gets or sets whether to require authentication via Touch ID, Face ID, or a passcode to access the keychain entry.

Optional:

- `active_directory_site_code` (String) Gets or sets the Active Directory site.
- `cache_name` (String) Gets or sets the Generic Security Services name of the Kerberos cache to use for this profile.
- `credential_bundle_id_access_control_list` (Set of String) Gets or sets a list of app Bundle IDs allowed to access the Kerberos Ticket Granting Ticket.
- `domain_realms` (Set of String) Gets or sets a list of realms for custom domain-realm mapping. Realms are case sensitive.
- `domains` (Set of String) Gets or sets a list of hosts or domain names for which the app extension performs SSO.
- `password_change_url` (String) Gets or sets the URL that the user will be sent to when they initiate a password change.
- `password_expiration_days` (Number) Overrides the default password expiration in days. For most domains, this value is calculated automatically.
- `password_expiration_notification_days` (Number) Gets or sets the number of days until the user is notified that their password will expire (default is 15).
- `password_minimum_age_days` (Number) Gets or sets the minimum number of days until a user can change their password again.
- `password_minimum_length` (Number) Gets or sets the minimum length of a password.
- `password_previous_password_block_count` (Number) Gets or sets the number of previous passwords to block.
- `password_requirements_description` (String) Gets or sets a description of the password complexity requirements.
- `sign_in_help_text` (String) Text displayed to the user at the Kerberos sign in window. Available for devices running iOS and iPadOS versions 14 and later.
- `user_principal_name` (String) Gets or sets the principle user name to use for this profile. The realm name does not need to be included.


<a id="nestedatt--ios_device_features--ios_single_sign_on_extension--redirect"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension.redirect`

Required:

- `extension_identifier` (String) Gets or sets the bundle ID of the app extension that performs SSO for the specified URLs.
- `team_identifier` (String) Gets or sets the team ID of the app extension that performs SSO for the specified URLs.
- `url_prefixes` (Set of String) One or more URL prefixes of identity providers on whose behalf the app extension performs single sign-on. URLs must begin with http:// or https://. All URL prefixes must be unique for all profiles.

Optional:

- `configurations` (Attributes Set) Gets or sets a list of typed key-value pairs used to configure Credential-type profiles. This collection can contain a maximum of 500 elements. / A key-value pair with a string key and a typed value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keytypedvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension--redirect--configurations))

<a id="nestedatt--ios_device_features--ios_single_sign_on_extension--redirect--configurations"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension.redirect.configurations`

Required:

- `key` (String) The string key of the key-value pair.

Optional:

- `boolean` (Attributes) A key-value pair with a string key and a Boolean value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keybooleanvaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension--redirect--configurations--boolean))
- `integer` (Attributes) A key-value pair with a string key and an integer value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyintegervaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension--redirect--configurations--integer))
- `real` (Attributes) A key-value pair with a string key and a real (floating-point) value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyrealvaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension--redirect--configurations--real))
- `string` (Attributes) A key-value pair with a string key and a string value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keystringvaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_device_features--ios_single_sign_on_extension--redirect--configurations--string))

<a id="nestedatt--ios_device_features--ios_single_sign_on_extension--redirect--configurations--boolean"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension.redirect.configurations.boolean`

Required:

- `value` (Boolean) The Boolean value of the key-value pair.


<a id="nestedatt--ios_device_features--ios_single_sign_on_extension--redirect--configurations--integer"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension.redirect.configurations.integer`

Required:

- `value` (Number) The integer value of the key-value pair.


<a id="nestedatt--ios_device_features--ios_single_sign_on_extension--redirect--configurations--real"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension.redirect.configurations.real`

Required:

- `value` (Number) The real (floating-point) value of the key-value pair.


<a id="nestedatt--ios_device_features--ios_single_sign_on_extension--redirect--configurations--string"></a>
### Nested Schema for `ios_device_features.ios_single_sign_on_extension.redirect.configurations.string`

Required:

- `value` (String) The string value of the key-value pair.





<a id="nestedatt--ios_device_features--notification_settings"></a>
### Nested Schema for `ios_device_features.notification_settings`

Optional:

- `alert_type` (String) Indicates the type of alert for notifications for this app. / Notification Settings Alert Type; possible values are: `deviceDefault` (Device default value, no intent.), `banner` (Banner.), `modal` (Modal.), `none` (None.). The _provider_ default value is `"deviceDefault"`.
- `app_name` (String) Application name to be associated with the bundleID.
- `badges_enabled` (Boolean) Indicates whether badges are allowed for this app.
- `bundle_id` (String) Bundle id of app to which to apply these notification settings. The _provider_ default value is `""`.
- `enabled` (Boolean) Indicates whether notifications are allowed for this app.
- `preview_visibility` (String) Overrides the notification preview policy set by the user on an iOS device. / Determines when notification previews are visible on an iOS device. Previews can include things like text (from Messages and Mail) and invitation details (from Calendar). When configured, it will override the user's defined preview settings; possible values are: `notConfigured` (Notification preview settings will not be overwritten.), `alwaysShow` (Always show notification previews.), `hideWhenLocked` (Only show notification previews when the device is unlocked.), `neverShow` (Never show notification previews.). The _provider_ default value is `"notConfigured"`.
- `publisher` (String) Publisher to be associated with the bundleID.
- `show_in_notification_center` (Boolean) Indicates whether notifications can be shown in notification center.
- `show_on_lock_screen` (Boolean) Indicates whether notifications can be shown on the lock screen.
- `sounds_enabled` (Boolean) Indicates whether sounds are allowed for this app.


<a id="nestedatt--ios_device_features--single_sign_on_settings"></a>
### Nested Schema for `ios_device_features.single_sign_on_settings`

Optional:

- `allowed_apps_list` (Attributes Set) List of app identifiers that are allowed to use this login. If this field is omitted, the login applies to all applications on the device. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--ios_device_features--single_sign_on_settings--allowed_apps_list))
- `allowed_urls` (Set of String) List of HTTP URLs that must be matched in order to use this login. With iOS 9.0 or later, a wildcard characters may be used. The _provider_ default value is `[]`.
- `display_name` (String) The display name of login settings shown on the receiving device.
- `kerberos_principal_name` (String) A Kerberos principal name. If not provided, the user is prompted for one during profile installation.
- `kerberos_realm` (String) A Kerberos realm name. Case sensitive.

<a id="nestedatt--ios_device_features--single_sign_on_settings--allowed_apps_list"></a>
### Nested Schema for `ios_device_features.single_sign_on_settings.allowed_apps_list`

Required:

- `name` (String) The application name

Optional:

- `app_id` (String) The application or bundle identifier of the application
- `app_store_url` (String) The Store URL of the application
- `publisher` (String) The publisher of the application



<a id="nestedatt--ios_device_features--wallpaper_image"></a>
### Nested Schema for `ios_device_features.wallpaper_image`

Optional:

- `type` (String) Indicates the content mime type.
- `value` (String) The byte array that contains the actual content.



<a id="nestedatt--ios_eas_email_profile"></a>
### Nested Schema for `ios_eas_email_profile`

Required:

- `account_name` (String) Account name.
- `email_address_source` (String) Email attribute that is picked from AAD and injected into this profile before installing on the device. / Possible values for username source or email source; possible values are: `userPrincipalName` (User principal name.), `primarySmtpAddress` (Primary SMTP address.)
- `host_name` (String) Exchange location that (URL) that the native mail app connects to.

Optional:

- `authentication_method` (String) Authentication method for this Email profile. / Exchange Active Sync authentication method; possible values are: `usernameAndPassword` (Authenticate with a username and password.), `certificate` (Authenticate with a certificate.), `derivedCredential` (Authenticate with derived credential.)
- `block_moving_messages_to_other_email_accounts` (Boolean) Indicates whether or not to block moving messages to other email accounts.
- `block_sending_email_from_third_party_apps` (Boolean) Indicates whether or not to block sending email from third party apps.
- `block_syncing_recently_used_email_addresses` (Boolean) Indicates whether or not to block syncing recently used email addresses, for instance - when composing new email.
- `custom_domain_name` (String) Custom domain name value used while generating an email profile before installing on the device.
- `duration_of_email_to_sync` (String) Duration of time email should be synced back to. . / Possible values for email sync duration; possible values are: `userDefined` (User Defined, default value, no intent.), `oneDay` (Sync one day of email.), `threeDays` (Sync three days of email.), `oneWeek` (Sync one week of email.), `twoWeeks` (Sync two weeks of email.), `oneMonth` (Sync one month of email.), `unlimited` (Sync an unlimited duration of email.). The _provider_ default value is `"userDefined"`.
- `eas_services` (String) Exchange data to sync. / Exchange Active Sync services; possible values are: `none`, `calendars` (Enables synchronization of calendars.), `contacts` (Enables synchronization of contacts.), `email` (Enables synchronization of email.), `notes` (Enables synchronization of notes.), `reminders` (Enables synchronization of reminders.)
- `eas_services_user_override_enabled` (Boolean) Allow users to change sync settings.
- `encryption_certificate_type` (String) Encryption Certificate type for this Email profile. / Supported certificate sources for email signing and encryption; possible values are: `none` (Do not use a certificate as a source.), `certificate` (Use an certificate for certificate source.), `derivedCredential` (Use a derived credential for certificate source.)
- `per_app_vpn_profile_id` (String) Profile ID of the Per-App VPN policy to be used to access emails from the native Mail client
- `require_ssl` (Boolean) Indicates whether or not to use SSL. The _provider_ default value is `true`.
- `signing_certificate_type` (String) Signing Certificate type for this Email profile. / Supported certificate sources for email signing and encryption; possible values are: `none` (Do not use a certificate as a source.), `certificate` (Use an certificate for certificate source.), `derivedCredential` (Use a derived credential for certificate source.)
- `use_oauth` (Boolean) Specifies whether the connection should use OAuth for authentication.
- `user_domain_name_source` (String) UserDomainname attribute that is picked from AAD and injected into this profile before installing on the device. / Domainname source; possible values are: `fullDomainName` (Full domain name.), `netBiosDomainName` (net bios domain name.)
- `username_aad_source` (String) Name of the AAD field, that will be used to retrieve UserName for email profile. / Username source; possible values are: `userPrincipalName` (User principal name.), `primarySmtpAddress` (Primary SMTP address.), `samAccountName` (The user sam account name.)
- `username_source` (String) Username attribute that is picked from AAD and injected into this profile before installing on the device. / Possible values for username source or email source; possible values are: `userPrincipalName` (User principal name.), `primarySmtpAddress` (Primary SMTP address.). The _provider_ default value is `"userPrincipalName"`.


<a id="nestedatt--ios_general_device"></a>
### Nested Schema for `ios_general_device`

Optional:

- `account_block_modification` (Boolean) Indicates whether or not to allow account modification when the device is in supervised mode. The _provider_ default value is `false`.
- `activation_lock_allow_when_supervised` (Boolean) Indicates whether or not to allow activation lock when the device is in the supervised mode. The _provider_ default value is `false`.
- `airdrop_blocked` (Boolean) Indicates whether or not to allow AirDrop when the device is in supervised mode. The _provider_ default value is `false`.
- `airdrop_force_unmanaged_drop_target` (Boolean) Indicates whether or not to cause AirDrop to be considered an unmanaged drop target (iOS 9.0 and later). The _provider_ default value is `false`.
- `airplay_force_pairing_password_for_outgoing_requests` (Boolean) Indicates whether or not to enforce all devices receiving AirPlay requests from this device to use a pairing password. The _provider_ default value is `false`.
- `airprint_block_credentials_storage` (Boolean) Indicates whether or not keychain storage of username and password for Airprint is blocked (iOS 11.0 and later). The _provider_ default value is `false`.
- `airprint_blocked` (Boolean) Indicates whether or not AirPrint is blocked (iOS 11.0 and later). The _provider_ default value is `false`.
- `airprint_blocki_beacon_discovery` (Boolean) Indicates whether or not iBeacon discovery of AirPrint printers is blocked. This prevents spurious AirPrint Bluetooth beacons from phishing for network traffic (iOS 11.0 and later). The _provider_ default value is `false`.
- `airprint_force_trusted_tls` (Boolean) Indicates if trusted certificates are required for TLS printing communication (iOS 11.0 and later). The _provider_ default value is `false`.
- `app_clips_blocked` (Boolean) Prevents a user from adding any App Clips and removes any existing App Clips on the device. The _provider_ default value is `false`.
- `app_removal_blocked` (Boolean) Indicates if the removal of apps is allowed. The _provider_ default value is `false`.
- `app_store_block_automatic_downloads` (Boolean) Indicates whether or not to block the automatic downloading of apps purchased on other devices when the device is in supervised mode (iOS 9.0 and later). The _provider_ default value is `false`.
- `app_store_block_in_app_purchases` (Boolean) Indicates whether or not to block the user from making in app purchases. The _provider_ default value is `false`.
- `app_store_block_ui_app_installation` (Boolean) Indicates whether or not to block the App Store app, not restricting installation through Host apps. Applies to supervised mode only (iOS 9.0 and later). The _provider_ default value is `false`.
- `app_store_blocked` (Boolean) Indicates whether or not to block the user from using the App Store. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.
- `app_store_require_password` (Boolean) Indicates whether or not to require a password when using the app store. The _provider_ default value is `false`.
- `apple_news_blocked` (Boolean) Indicates whether or not to block the user from using News when the device is in supervised mode (iOS 9.0 and later). The _provider_ default value is `false`.
- `apple_personalized_ads_blocked` (Boolean) Limits Apple personalized advertising when true. Available in iOS 14 and later. The _provider_ default value is `false`.
- `apple_watch_block_pairing` (Boolean) Indicates whether or not to allow Apple Watch pairing when the device is in supervised mode (iOS 9.0 and later). The _provider_ default value is `false`.
- `apple_watch_force_wrist_detection` (Boolean) Indicates whether or not to force a paired Apple Watch to use Wrist Detection (iOS 8.2 and later). The _provider_ default value is `false`.
- `apps_single_app_mode_list` (Attributes Set) Gets or sets the list of iOS apps allowed to autonomously enter Single App Mode. Supervised only. iOS 7.0 and later. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--ios_general_device--apps_single_app_mode_list))
- `apps_visibility_list` (Attributes Set) List of apps in the visibility list (either visible/launchable apps list or hidden/unlaunchable apps list, controlled by AppsVisibilityListType) (iOS 9.3 and later). This collection can contain a maximum of 10000 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--ios_general_device--apps_visibility_list))
- `apps_visibility_list_type` (String) Type of list that is in the AppsVisibilityList. / Possible values of the compliance app list; possible values are: `none` (Default value, no intent.), `appsInListCompliant` (The list represents the apps that will be considered compliant (only apps on the list are compliant).), `appsNotInListCompliant` (The list represents the apps that will be considered non compliant (all apps are compliant except apps on the list).). The _provider_ default value is `"none"`.
- `auto_fill_force_authentication` (Boolean) Indicates whether or not to force user authentication before autofilling passwords and credit card information in Safari and other apps on supervised devices. The _provider_ default value is `false`.
- `auto_unlock_blocked` (Boolean) Blocks users from unlocking their device with Apple Watch. Available for devices running iOS and iPadOS versions 14.5 and later. The _provider_ default value is `false`.
- `block_system_app_removal` (Boolean) Indicates whether or not the removal of system apps from the device is blocked on a supervised device (iOS 11.0 and later). The _provider_ default value is `false`.
- `bluetooth_block_modification` (Boolean) Indicates whether or not to allow modification of Bluetooth settings when the device is in supervised mode (iOS 10.0 and later). The _provider_ default value is `false`.
- `camera_blocked` (Boolean) Indicates whether or not to block the user from accessing the camera of the device. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.
- `cellular_block_data_roaming` (Boolean) Indicates whether or not to block data roaming. The _provider_ default value is `false`.
- `cellular_block_global_background_fetch_while_roaming` (Boolean) Indicates whether or not to block global background fetch while roaming. The _provider_ default value is `false`.
- `cellular_block_per_app_data_modification` (Boolean) Indicates whether or not to allow changes to cellular app data usage settings when the device is in supervised mode. The _provider_ default value is `false`.
- `cellular_block_personal_hotspot` (Boolean) Indicates whether or not to block Personal Hotspot. The _provider_ default value is `false`.
- `cellular_block_personal_hotspot_modification` (Boolean) Indicates whether or not to block the user from modifying the personal hotspot setting (iOS 12.2 or later). The _provider_ default value is `false`.
- `cellular_block_plan_modification` (Boolean) Indicates whether or not to allow users to change the settings of the cellular plan on a supervised device. The _provider_ default value is `false`.
- `cellular_block_voice_roaming` (Boolean) Indicates whether or not to block voice roaming. The _provider_ default value is `false`.
- `certificates_block_untrusted_tls_certificates` (Boolean) Indicates whether or not to block untrusted TLS certificates. The _provider_ default value is `false`.
- `classroom_app_block_remote_screen_observation` (Boolean) Indicates whether or not to allow remote screen observation by Classroom app when the device is in supervised mode (iOS 9.3 and later). The _provider_ default value is `false`.
- `classroom_app_force_unprompted_screen_observation` (Boolean) Indicates whether or not to automatically give permission to the teacher of a managed course on the Classroom app to view a student's screen without prompting when the device is in supervised mode. The _provider_ default value is `false`.
- `classroom_force_automatically_join_classes` (Boolean) Indicates whether or not to automatically give permission to the teacher's requests, without prompting the student, when the device is in supervised mode. The _provider_ default value is `false`.
- `classroom_force_request_permission_to_leave_classes` (Boolean) Indicates whether a student enrolled in an unmanaged course via Classroom will request permission from the teacher when attempting to leave the course (iOS 11.3 and later). The _provider_ default value is `false`.
- `classroom_force_unprompted_app_and_device_lock` (Boolean) Indicates whether or not to allow the teacher to lock apps or the device without prompting the student. Supervised only. The _provider_ default value is `false`.
- `compliant_app_list_type` (String) List that is in the AppComplianceList. / Possible values of the compliance app list; possible values are: `none` (Default value, no intent.), `appsInListCompliant` (The list represents the apps that will be considered compliant (only apps on the list are compliant).), `appsNotInListCompliant` (The list represents the apps that will be considered non compliant (all apps are compliant except apps on the list).). The _provider_ default value is `"none"`.
- `compliant_apps_list` (Attributes Set) List of apps in the compliance (either allow list or block list, controlled by CompliantAppListType). This collection can contain a maximum of 10000 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--ios_general_device--compliant_apps_list))
- `configuration_profile_block_changes` (Boolean) Indicates whether or not to block the user from installing configuration profiles and certificates interactively when the device is in supervised mode. The _provider_ default value is `false`.
- `contacts_allow_managed_to_unmanaged_write` (Boolean) Indicates whether or not managed apps can write contacts to unmanaged contacts accounts (iOS 12.0 and later). The _provider_ default value is `false`.
- `contacts_allow_unmanaged_to_managed_read` (Boolean) Indicates whether or not unmanaged apps can read from managed contacts accounts (iOS 12.0 or later). The _provider_ default value is `false`.
- `continuous_path_keyboard_blocked` (Boolean) Indicates whether or not to block the continuous path keyboard when the device is supervised (iOS 13 or later). The _provider_ default value is `false`.
- `date_and_time_force_set_automatically` (Boolean) Indicates whether or not the Date and Time "Set Automatically" feature is enabled and cannot be turned off by the user (iOS 12.0 and later). The _provider_ default value is `false`.
- `definition_lookup_blocked` (Boolean) Indicates whether or not to block definition lookup when the device is in supervised mode (iOS 8.1.3 and later ). The _provider_ default value is `false`.
- `device_block_enable_restrictions` (Boolean) Indicates whether or not to allow the user to enables restrictions in the device settings when the device is in supervised mode. The _provider_ default value is `false`.
- `device_block_erase_content_and_settings` (Boolean) Indicates whether or not to allow the use of the 'Erase all content and settings' option on the device when the device is in supervised mode. The _provider_ default value is `false`.
- `device_block_name_modification` (Boolean) Indicates whether or not to allow device name modification when the device is in supervised mode (iOS 9.0 and later). The _provider_ default value is `false`.
- `diagnostic_data_block_submission` (Boolean) Indicates whether or not to block diagnostic data submission. The _provider_ default value is `false`.
- `diagnostic_data_block_submission_modification` (Boolean) Indicates whether or not to allow diagnostics submission settings modification when the device is in supervised mode (iOS 9.3.2 and later). The _provider_ default value is `false`.
- `documents_block_managed_documents_in_unmanaged_apps` (Boolean) Indicates whether or not to block the user from viewing managed documents in unmanaged apps. The _provider_ default value is `false`.
- `documents_block_unmanaged_documents_in_managed_apps` (Boolean) Indicates whether or not to block the user from viewing unmanaged documents in managed apps. The _provider_ default value is `false`.
- `email_in_domain_suffixes` (Set of String) An email address lacking a suffix that matches any of these strings will be considered out-of-domain. The _provider_ default value is `[]`.
- `enterprise_app_block_trust` (Boolean) Indicates whether or not to block the user from trusting an enterprise app. The _provider_ default value is `false`.
- `enterprise_book_block_backup` (Boolean) Indicates whether or not Enterprise book back up is blocked. The _provider_ default value is `false`.
- `enterprise_book_block_metadata_sync` (Boolean) Indicates whether or not Enterprise book notes and highlights sync is blocked. The _provider_ default value is `false`.
- `esim_block_modification` (Boolean) Indicates whether or not to allow the addition or removal of cellular plans on the eSIM of a supervised device. The _provider_ default value is `false`.
- `facetime_blocked` (Boolean) Indicates whether or not to block the user from using FaceTime. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.
- `files_network_drive_access_blocked` (Boolean) Indicates if devices can access files or other resources on a network server using the Server Message Block (SMB) protocol. Available for devices running iOS and iPadOS, versions 13.0 and later. The _provider_ default value is `false`.
- `files_usb_drive_access_blocked` (Boolean) Indicates if sevices with access can connect to and open files on a USB drive. Available for devices running iOS and iPadOS, versions 13.0 and later. The _provider_ default value is `false`.
- `find_my_device_in_find_my_app_blocked` (Boolean) Indicates whether or not to block Find My Device when the device is supervised (iOS 13 or later). The _provider_ default value is `false`.
- `find_my_friends_blocked` (Boolean) Indicates whether or not to block changes to Find My Friends when the device is in supervised mode. The _provider_ default value is `false`.
- `find_my_friends_in_find_my_app_blocked` (Boolean) Indicates whether or not to block Find My Friends when the device is supervised (iOS 13 or later). The _provider_ default value is `false`.
- `game_center_blocked` (Boolean) Indicates whether or not to block the user from using Game Center when the device is in supervised mode. The _provider_ default value is `false`.
- `gaming_block_game_center_friends` (Boolean) Indicates whether or not to block the user from having friends in Game Center. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.
- `gaming_block_multiplayer` (Boolean) Indicates whether or not to block the user from using multiplayer gaming. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.
- `host_pairing_blocked` (Boolean) indicates whether or not to allow host pairing to control the devices an iOS device can pair with when the iOS device is in supervised mode. The _provider_ default value is `false`.
- `ibooks_store_block_erotica` (Boolean) Indicates whether or not to block the user from downloading media from the iBookstore that has been tagged as erotica. The _provider_ default value is `false`.
- `ibooks_store_blocked` (Boolean) Indicates whether or not to block the user from using the iBooks Store when the device is in supervised mode. The _provider_ default value is `false`.
- `icloud_block_activity_continuation` (Boolean) Indicates whether or not to block the user from continuing work they started on iOS device to another iOS or macOS device. The _provider_ default value is `false`.
- `icloud_block_backup` (Boolean) Indicates whether or not to block iCloud backup. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.
- `icloud_block_document_sync` (Boolean) Indicates whether or not to block iCloud document sync. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.
- `icloud_block_managed_apps_sync` (Boolean) Indicates whether or not to block Managed Apps Cloud Sync. The _provider_ default value is `false`.
- `icloud_block_photo_library` (Boolean) Indicates whether or not to block iCloud Photo Library. The _provider_ default value is `false`.
- `icloud_block_photo_stream_sync` (Boolean) Indicates whether or not to block iCloud Photo Stream Sync. The _provider_ default value is `false`.
- `icloud_block_shared_photo_stream` (Boolean) Indicates whether or not to block Shared Photo Stream. The _provider_ default value is `false`.
- `icloud_private_relay_blocked` (Boolean) iCloud private relay is an iCloud+ service that prevents networks and servers from monitoring a person's activity across the internet. By blocking iCloud private relay, Apple will not encrypt the traffic leaving the device. Available for devices running iOS 15 and later. The _provider_ default value is `false`.
- `icloud_require_encrypted_backup` (Boolean) Indicates whether or not to require backups to iCloud be encrypted. The _provider_ default value is `false`.
- `itunes_block_explicit_content` (Boolean) Indicates whether or not to block the user from accessing explicit content in iTunes and the App Store. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.
- `itunes_block_music_service` (Boolean) Indicates whether or not to block Music service and revert Music app to classic mode when the device is in supervised mode (iOS 9.3 and later and macOS 10.12 and later). The _provider_ default value is `false`.
- `itunes_block_radio` (Boolean) Indicates whether or not to block the user from using iTunes Radio when the device is in supervised mode (iOS 9.3 and later). The _provider_ default value is `false`.
- `itunes_blocked` (Boolean) Indicates whether or not to block the iTunes app. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.
- `keyboard_block_auto_correct` (Boolean) Indicates whether or not to block keyboard auto-correction when the device is in supervised mode (iOS 8.1.3 and later). The _provider_ default value is `false`.
- `keyboard_block_dictation` (Boolean) Indicates whether or not to block the user from using dictation input when the device is in supervised mode. The _provider_ default value is `false`.
- `keyboard_block_predictive` (Boolean) Indicates whether or not to block predictive keyboards when device is in supervised mode (iOS 8.1.3 and later). The _provider_ default value is `false`.
- `keyboard_block_shortcuts` (Boolean) Indicates whether or not to block keyboard shortcuts when the device is in supervised mode (iOS 9.0 and later). The _provider_ default value is `false`.
- `keyboard_block_spell_check` (Boolean) Indicates whether or not to block keyboard spell-checking when the device is in supervised mode (iOS 8.1.3 and later). The _provider_ default value is `false`.
- `keychain_block_cloud_sync` (Boolean) Indicates whether or not iCloud keychain synchronization is blocked. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.
- `kiosk_mode_allow_assistive_speak` (Boolean) Indicates whether or not to allow assistive speak while in kiosk mode. The _provider_ default value is `false`.
- `kiosk_mode_allow_assistive_touch_settings` (Boolean) Indicates whether or not to allow access to the Assistive Touch Settings while in kiosk mode. The _provider_ default value is `false`.
- `kiosk_mode_allow_color_inversion_settings` (Boolean) Indicates whether or not to allow access to the Color Inversion Settings while in kiosk mode. The _provider_ default value is `false`.
- `kiosk_mode_allow_voice_control_modification` (Boolean) Indicates whether or not to allow the user to toggle voice control in kiosk mode. The _provider_ default value is `false`.
- `kiosk_mode_allow_voice_over_settings` (Boolean) Indicates whether or not to allow access to the voice over settings while in kiosk mode. The _provider_ default value is `false`.
- `kiosk_mode_allow_zoom_settings` (Boolean) Indicates whether or not to allow access to the zoom settings while in kiosk mode. The _provider_ default value is `false`.
- `kiosk_mode_app_store_url` (String) URL in the app store to the app to use for kiosk mode. Use if KioskModeManagedAppId is not known.
- `kiosk_mode_app_type` (String) Type of app to run in kiosk mode. / App source options for iOS kiosk mode; possible values are: `notConfigured` (Device default value, no intent.), `appStoreApp` (The app to be run comes from the app store.), `managedApp` (The app to be run is built into the device.), `builtInApp` (The app to be run is a managed app.). The _provider_ default value is `"notConfigured"`.
- `kiosk_mode_block_auto_lock` (Boolean) Indicates whether or not to block device auto lock while in kiosk mode. The _provider_ default value is `false`.
- `kiosk_mode_block_ringer_switch` (Boolean) Indicates whether or not to block use of the ringer switch while in kiosk mode. The _provider_ default value is `false`.
- `kiosk_mode_block_screen_rotation` (Boolean) Indicates whether or not to block screen rotation while in kiosk mode. The _provider_ default value is `false`.
- `kiosk_mode_block_sleep_button` (Boolean) Indicates whether or not to block use of the sleep button while in kiosk mode. The _provider_ default value is `false`.
- `kiosk_mode_block_touchscreen` (Boolean) Indicates whether or not to block use of the touchscreen while in kiosk mode. The _provider_ default value is `false`.
- `kiosk_mode_block_volume_buttons` (Boolean) Indicates whether or not to block the volume buttons while in Kiosk Mode. The _provider_ default value is `false`.
- `kiosk_mode_built_in_app_id` (String) ID for built-in apps to use for kiosk mode. Used when KioskModeManagedAppId and KioskModeAppStoreUrl are not set.
- `kiosk_mode_enable_voice_control` (Boolean) Indicates whether or not to enable voice control in kiosk mode. The _provider_ default value is `false`.
- `kiosk_mode_managed_app_id` (String) Managed app id of the app to use for kiosk mode. If KioskModeManagedAppId is specified then KioskModeAppStoreUrl will be ignored.
- `kiosk_mode_require_assistive_touch` (Boolean) Indicates whether or not to require assistive touch while in kiosk mode. The _provider_ default value is `false`.
- `kiosk_mode_require_color_inversion` (Boolean) Indicates whether or not to require color inversion while in kiosk mode. The _provider_ default value is `false`.
- `kiosk_mode_require_mono_audio` (Boolean) Indicates whether or not to require mono audio while in kiosk mode. The _provider_ default value is `false`.
- `kiosk_mode_require_voice_over` (Boolean) Indicates whether or not to require voice over while in kiosk mode. The _provider_ default value is `false`.
- `kiosk_mode_require_zoom` (Boolean) Indicates whether or not to require zoom while in kiosk mode. The _provider_ default value is `false`.
- `lock_screen_block_control_center` (Boolean) Indicates whether or not to block the user from using control center on the lock screen. The _provider_ default value is `false`.
- `lock_screen_block_notification_view` (Boolean) Indicates whether or not to block the user from using the notification view on the lock screen. The _provider_ default value is `false`.
- `lock_screen_block_passbook` (Boolean) Indicates whether or not to block the user from using passbook when the device is locked. The _provider_ default value is `false`.
- `lock_screen_block_today_view` (Boolean) Indicates whether or not to block the user from using the Today View on the lock screen. The _provider_ default value is `false`.
- `managed_pasteboard_required` (Boolean) Open-in management controls how people share data between unmanaged and managed apps. Setting this to true enforces copy/paste restrictions based on how you configured <b>Block viewing corporate documents in unmanaged apps </b> and <b> Block viewing non-corporate documents in corporate apps.</b>. The _provider_ default value is `false`.
- `media_content_rating_apps` (String) Media content rating settings for Apps. / Apps rating as in media content; possible values are: `allAllowed` (Default value, allow all apps content), `allBlocked` (Do not allow any apps content), `agesAbove4` (4+, age 4 and above), `agesAbove9` (9+, age 9 and above), `agesAbove12` (12+, age 12 and above), `agesAbove17` (17+, age 17 and above). The _provider_ default value is `"allAllowed"`.
- `media_content_rating_australia` (Attributes) Media content rating settings for Australia / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediacontentratingaustralia?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_general_device--media_content_rating_australia))
- `media_content_rating_canada` (Attributes) Media content rating settings for Canada / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediacontentratingcanada?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_general_device--media_content_rating_canada))
- `media_content_rating_france` (Attributes) Media content rating settings for France / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediacontentratingfrance?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_general_device--media_content_rating_france))
- `media_content_rating_germany` (Attributes) Media content rating settings for Germany / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediacontentratinggermany?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_general_device--media_content_rating_germany))
- `media_content_rating_ireland` (Attributes) Media content rating settings for Ireland / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediacontentratingireland?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_general_device--media_content_rating_ireland))
- `media_content_rating_japan` (Attributes) Media content rating settings for Japan / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediacontentratingjapan?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_general_device--media_content_rating_japan))
- `media_content_rating_new_zealand` (Attributes) Media content rating settings for New Zealand / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediacontentratingnewzealand?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_general_device--media_content_rating_new_zealand))
- `media_content_rating_united_kingdom` (Attributes) Media content rating settings for United Kingdom / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediacontentratingunitedkingdom?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_general_device--media_content_rating_united_kingdom))
- `media_content_rating_united_states` (Attributes) Media content rating settings for United States / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-mediacontentratingunitedstates?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_general_device--media_content_rating_united_states))
- `messages_blocked` (Boolean) Indicates whether or not to block the user from using the Messages app on the supervised device. The _provider_ default value is `false`.
- `network_usage_rules` (Attributes Set) List of managed apps and the network rules that applies to them. This collection can contain a maximum of 1000 elements. / Network Usage Rules allow enterprises to specify how managed apps use networks, such as cellular data networks. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-iosnetworkusagerule?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--ios_general_device--network_usage_rules))
- `nfc_blocked` (Boolean) Disable NFC to prevent devices from pairing with other NFC-enabled devices. Available for iOS/iPadOS devices running 14.2 and later. The _provider_ default value is `false`.
- `notifications_block_settings_modification` (Boolean) Indicates whether or not to allow notifications settings modification (iOS 9.3 and later). The _provider_ default value is `false`.
- `on_device_only_dictation_forced` (Boolean) Disables connections to Siri servers so that users can’t use Siri to dictate text. Available for devices running iOS and iPadOS versions 14.5 and later. The _provider_ default value is `false`.
- `on_device_only_translation_forced` (Boolean) When set to TRUE, the setting disables connections to Siri servers so that users can’t use Siri to translate text. When set to FALSE, the setting allows connections to to Siri servers to users can use Siri to translate text. Available for devices running iOS and iPadOS versions 15.0 and later. The _provider_ default value is `false`.
- `passcode_block_fingerprint_modification` (Boolean) Block modification of registered Touch ID fingerprints when in supervised mode. The _provider_ default value is `false`.
- `passcode_block_fingerprint_unlock` (Boolean) Indicates whether or not to block fingerprint unlock. The _provider_ default value is `false`.
- `passcode_block_modification` (Boolean) Indicates whether or not to allow passcode modification on the supervised device (iOS 9.0 and later). The _provider_ default value is `false`.
- `passcode_block_simple` (Boolean) Indicates whether or not to block simple passcodes. The _provider_ default value is `false`.
- `passcode_expiration_days` (Number) Number of days before the passcode expires. Valid values 1 to 65535
- `passcode_minimum_character_set_count` (Number) Number of character sets a passcode must contain. Valid values 0 to 4
- `passcode_minimum_length` (Number) Minimum length of passcode. Valid values 4 to 14
- `passcode_minutes_of_inactivity_before_lock` (Number) Minutes of inactivity before a passcode is required.
- `passcode_minutes_of_inactivity_before_screen_timeout` (Number) Minutes of inactivity before the screen times out.
- `passcode_previous_passcode_block_count` (Number) Number of previous passcodes to block. Valid values 1 to 24
- `passcode_required` (Boolean) Indicates whether or not to require a passcode. The _provider_ default value is `false`.
- `passcode_required_type` (String) Type of passcode that is required. / Possible values of required passwords; possible values are: `deviceDefault` (Device default value, no intent.), `alphanumeric` (Alphanumeric password required.), `numeric` (Numeric password required.). The _provider_ default value is `"deviceDefault"`.
- `passcode_sign_in_failure_count_before_wipe` (Number) Number of sign in failures allowed before wiping the device. Valid values 2 to 11
- `password_block_airdrop_sharing` (Boolean) Indicates whether or not to block sharing passwords with the AirDrop passwords feature iOS 12.0 and later). The _provider_ default value is `false`.
- `password_block_auto_fill` (Boolean) Indicates if the AutoFill passwords feature is allowed (iOS 12.0 and later). The _provider_ default value is `false`.
- `password_block_proximity_requests` (Boolean) Indicates whether or not to block requesting passwords from nearby devices (iOS 12.0 and later). The _provider_ default value is `false`.
- `pki_block_ota_updates` (Boolean) Indicates whether or not over-the-air PKI updates are blocked. Setting this restriction to false does not disable CRL and OCSP checks (iOS 7.0 and later). The _provider_ default value is `false`.
- `podcasts_blocked` (Boolean) Indicates whether or not to block the user from using podcasts on the supervised device (iOS 8.0 and later). The _provider_ default value is `false`.
- `privacy_force_limit_ad_tracking` (Boolean) Indicates if ad tracking is limited.(iOS 7.0 and later). The _provider_ default value is `false`.
- `proximity_block_setup_to_new_device` (Boolean) Indicates whether or not to enable the prompt to setup nearby devices with a supervised device. The _provider_ default value is `false`.
- `safari_block_autofill` (Boolean) Indicates whether or not to block the user from using Auto fill in Safari. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.
- `safari_block_javascript` (Boolean) Indicates whether or not to block JavaScript in Safari. The _provider_ default value is `false`.
- `safari_block_popups` (Boolean) Indicates whether or not to block popups in Safari. The _provider_ default value is `false`.
- `safari_blocked` (Boolean) Indicates whether or not to block the user from using Safari. Requires a supervised device for iOS 13 and later. The _provider_ default value is `false`.
- `safari_cookie_settings` (String) Cookie settings for Safari. / Web Browser Cookie Settings; possible values are: `browserDefault` (Browser default value, no intent.), `blockAlways` (Always block cookies.), `allowCurrentWebSite` (Allow cookies from current Web site.), `allowFromWebsitesVisited` (Allow Cookies from websites visited.), `allowAlways` (Always allow cookies.). The _provider_ default value is `"browserDefault"`.
- `safari_managed_domains` (Set of String) URLs matching the patterns listed here will be considered managed. The _provider_ default value is `[]`.
- `safari_password_auto_fill_domains` (Set of String) Users can save passwords in Safari only from URLs matching the patterns listed here. Applies to devices in supervised mode (iOS 9.3 and later). The _provider_ default value is `[]`.
- `safari_require_fraud_warning` (Boolean) Indicates whether or not to require fraud warning in Safari. The _provider_ default value is `false`.
- `screen_capture_blocked` (Boolean) Indicates whether or not to block the user from taking Screenshots. The _provider_ default value is `false`.
- `shared_device_block_temporary_sessions` (Boolean) Indicates whether or not to block temporary sessions on Shared iPads (iOS 13.4 or later). The _provider_ default value is `false`.
- `siri_block_user_generated_content` (Boolean) Indicates whether or not to block Siri from querying user-generated content when used on a supervised device. The _provider_ default value is `false`.
- `siri_blocked` (Boolean) Indicates whether or not to block the user from using Siri. The _provider_ default value is `false`.
- `siri_blocked_when_locked` (Boolean) Indicates whether or not to block the user from using Siri when locked. The _provider_ default value is `false`.
- `siri_require_profanity_filter` (Boolean) Indicates whether or not to prevent Siri from dictating, or speaking profane language on supervised device. The _provider_ default value is `false`.
- `software_updates_enforced_delay_in_days` (Number) Sets how many days a software update will be delyed for a supervised device. Valid values 0 to 90
- `software_updates_force_delayed` (Boolean) Indicates whether or not to delay user visibility of software updates when the device is in supervised mode. The _provider_ default value is `false`.
- `spotlight_block_internet_results` (Boolean) Indicates whether or not to block Spotlight search from returning internet results on supervised device. The _provider_ default value is `false`.
- `unpaired_external_boot_to_recovery_allowed` (Boolean) Allow users to boot devices into recovery mode with unpaired devices. Available for devices running iOS and iPadOS versions 14.5 and later. The _provider_ default value is `false`.
- `usb_restricted_mode_blocked` (Boolean) Indicates if connecting to USB accessories while the device is locked is allowed (iOS 11.4.1 and later). The _provider_ default value is `false`.
- `voice_dialing_blocked` (Boolean) Indicates whether or not to block voice dialing. The _provider_ default value is `false`.
- `vpn_block_creation` (Boolean) Indicates whether or not the creation of VPN configurations is blocked (iOS 11.0 and later). The _provider_ default value is `false`.
- `wallpaper_block_modification` (Boolean) Indicates whether or not to allow wallpaper modification on supervised device (iOS 9.0 and later) . The _provider_ default value is `false`.
- `wifi_connect_only_to_configured_networks` (Boolean) Indicates whether or not to force the device to use only Wi-Fi networks from configuration profiles when the device is in supervised mode. Available for devices running iOS and iPadOS versions 14.4 and earlier. Devices running 14.5+ should use the setting, “WiFiConnectToAllowedNetworksOnlyForced. The _provider_ default value is `false`.
- `wifi_connect_to_allowed_networks_only_forced` (Boolean) Require devices to use Wi-Fi networks set up via configuration profiles. Available for devices running iOS and iPadOS versions 14.5 and later. The _provider_ default value is `false`.
- `wifi_power_on_forced` (Boolean) Indicates whether or not Wi-Fi remains on, even when device is in airplane mode. Available for devices running iOS and iPadOS, versions 13.0 and later. The _provider_ default value is `false`.

<a id="nestedatt--ios_general_device--apps_single_app_mode_list"></a>
### Nested Schema for `ios_general_device.apps_single_app_mode_list`

Required:

- `name` (String) The application name

Optional:

- `app_id` (String) The application or bundle identifier of the application
- `app_store_url` (String) The Store URL of the application
- `publisher` (String) The publisher of the application


<a id="nestedatt--ios_general_device--apps_visibility_list"></a>
### Nested Schema for `ios_general_device.apps_visibility_list`

Required:

- `name` (String) The application name

Optional:

- `app_id` (String) The application or bundle identifier of the application
- `app_store_url` (String) The Store URL of the application
- `publisher` (String) The publisher of the application


<a id="nestedatt--ios_general_device--compliant_apps_list"></a>
### Nested Schema for `ios_general_device.compliant_apps_list`

Required:

- `name` (String) The application name

Optional:

- `app_id` (String) The application or bundle identifier of the application
- `app_store_url` (String) The Store URL of the application
- `publisher` (String) The publisher of the application


<a id="nestedatt--ios_general_device--media_content_rating_australia"></a>
### Nested Schema for `ios_general_device.media_content_rating_australia`

Required:

- `movie_rating` (String) Movies rating selected for Australia. / Movies rating labels in Australia; possible values are: `allAllowed` (Default value, allow all movies content), `allBlocked` (Do not allow any movies content), `general` (The G classification is suitable for everyone), `parentalGuidance` (The PG recommends viewers under 15 with guidance from parents or guardians), `mature` (The M classification is not recommended for viewers under 15), `agesAbove15` (The MA15+ classification is not suitable for viewers under 15), `agesAbove18` (The R18+ classification is not suitable for viewers under 18)
- `tv_rating` (String) TV rating selected for Australia. / TV content rating labels in Australia; possible values are: `allAllowed` (Default value, allow all TV shows content), `allBlocked` (Do not allow any TV shows content), `preschoolers` (The P classification is intended for preschoolers), `children` (The C classification is intended for children under 14), `general` (The G classification is suitable for all ages), `parentalGuidance` (The PG classification is recommended for young viewers), `mature` (The M classification is recommended for viewers over 15), `agesAbove15` (The MA15+ classification is not suitable for viewers under 15), `agesAbove15AdultViolence` (The AV15+ classification is not suitable for viewers under 15, adult violence-specific)


<a id="nestedatt--ios_general_device--media_content_rating_canada"></a>
### Nested Schema for `ios_general_device.media_content_rating_canada`

Required:

- `movie_rating` (String) Movies rating selected for Canada. / Movies rating labels in Canada; possible values are: `allAllowed` (Default value, allow all movies content), `allBlocked` (Do not allow any movies content), `general` (The G classification is suitable for all ages), `parentalGuidance` (The PG classification advises parental guidance), `agesAbove14` (The 14A classification is suitable for viewers above 14 or older), `agesAbove18` (The 18A classification is suitable for viewers above 18 or older), `restricted` (The R classification is restricted to 18 years and older)
- `tv_rating` (String) TV rating selected for Canada. / TV content rating labels in Canada; possible values are: `allAllowed` (Default value, allow all TV shows content), `allBlocked` (Do not allow any TV shows content), `children` (The C classification is suitable for children ages of 2 to 7 years), `childrenAbove8` (The C8 classification is suitable for children ages 8+), `general` (The G classification is suitable for general audience), `parentalGuidance` (PG, Parental Guidance), `agesAbove14` (The 14+ classification is intended for viewers ages 14 and older), `agesAbove18` (The 18+ classification is intended for viewers ages 18 and older)


<a id="nestedatt--ios_general_device--media_content_rating_france"></a>
### Nested Schema for `ios_general_device.media_content_rating_france`

Required:

- `movie_rating` (String) Movies rating selected for France. / Movies rating labels in France; possible values are: `allAllowed` (Default value, allow all movies content), `allBlocked` (Do not allow any movies content), `agesAbove10` (The 10 classification prohibits the screening of the film to minors under 10), `agesAbove12` (The 12 classification prohibits the screening of the film to minors under 12), `agesAbove16` (The 16 classification prohibits the screening of the film to minors under 16), `agesAbove18` (The 18 classification prohibits the screening to minors under 18)
- `tv_rating` (String) TV rating selected for France. / TV content rating labels in France; possible values are: `allAllowed` (Default value, allow all TV shows content), `allBlocked` (Do not allow any TV shows content), `agesAbove10` (The -10 classification is not recommended for children under 10), `agesAbove12` (The -12 classification is not recommended for children under 12), `agesAbove16` (The -16 classification is not recommended for children under 16), `agesAbove18` (The -18 classification is not recommended for persons under 18)


<a id="nestedatt--ios_general_device--media_content_rating_germany"></a>
### Nested Schema for `ios_general_device.media_content_rating_germany`

Required:

- `movie_rating` (String) Movies rating selected for Germany. / Movies rating labels in Germany; possible values are: `allAllowed` (Default value, allow all movies content), `allBlocked` (Do not allow any movies content), `general` (Ab 0 Jahren, no age restrictions), `agesAbove6` (Ab 6 Jahren, ages 6 and older), `agesAbove12` (Ab 12 Jahren, ages 12 and older), `agesAbove16` (Ab 16 Jahren, ages 16 and older), `adults` (Ab 18 Jahren, adults only)
- `tv_rating` (String) TV rating selected for Germany. / TV content rating labels in Germany; possible values are: `allAllowed` (Default value, allow all TV shows content), `allBlocked` (Do not allow any TV shows content), `general` (Ab 0 Jahren, no age restrictions), `agesAbove6` (Ab 6 Jahren, ages 6 and older), `agesAbove12` (Ab 12 Jahren, ages 12 and older), `agesAbove16` (Ab 16 Jahren, ages 16 and older), `adults` (Ab 18 Jahren, adults only)


<a id="nestedatt--ios_general_device--media_content_rating_ireland"></a>
### Nested Schema for `ios_general_device.media_content_rating_ireland`

Required:

- `movie_rating` (String) Movies rating selected for Ireland. / Movies rating labels in Ireland; possible values are: `allAllowed` (Default value, allow all movies content), `allBlocked` (Do not allow any movies content), `general` (Suitable for children of school going age), `parentalGuidance` (The PG classification advises parental guidance), `agesAbove12` (The 12A classification is suitable for viewers of 12 or older), `agesAbove15` (The 15A classification is suitable for viewers of 15 or older), `agesAbove16` (The 16 classification is suitable for viewers of 16 or older), `adults` (The 18 classification, suitable only for adults)
- `tv_rating` (String) TV rating selected for Ireland. / TV content rating labels in Ireland; possible values are: `allAllowed` (Default value, allow all TV shows content), `allBlocked` (Do not allow any TV shows content), `general` (The GA classification is suitable for all audiences), `children` (The CH classification is suitable for children), `youngAdults` (The YA classification is suitable for teenage audience), `parentalSupervision` (The PS classification invites parents and guardians to consider restriction children’s access), `mature` (The MA classification is suitable for adults)


<a id="nestedatt--ios_general_device--media_content_rating_japan"></a>
### Nested Schema for `ios_general_device.media_content_rating_japan`

Required:

- `movie_rating` (String) Movies rating selected for Japan. / Movies rating labels in Japan; possible values are: `allAllowed` (Default value, allow all movies content), `allBlocked` (Do not allow any movies content), `general` (Suitable for all ages), `parentalGuidance` (The PG-12 classification requests parental guidance for young people under 12), `agesAbove15` (The R15+ classification is suitable for viewers of 15 or older), `agesAbove18` (The R18+ classification is suitable for viewers of 18 or older)
- `tv_rating` (String) TV rating selected for Japan. / TV content rating labels in Japan; possible values are: `allAllowed` (Default value, allow all TV shows content), `allBlocked` (Do not allow any TV shows content), `explicitAllowed` (All TV content is explicitly allowed)


<a id="nestedatt--ios_general_device--media_content_rating_new_zealand"></a>
### Nested Schema for `ios_general_device.media_content_rating_new_zealand`

Required:

- `movie_rating` (String) Movies rating selected for New Zealand. / Movies rating labels in New Zealand; possible values are: `allAllowed` (Default value, allow all movies content), `allBlocked` (Do not allow any movies content), `general` (Suitable for general audience), `parentalGuidance` (The PG classification recommends parental guidance), `mature` (The M classification is suitable for mature audience), `agesAbove13` (The R13 classification is restricted to persons 13 years and over), `agesAbove15` (The R15 classification is restricted to persons 15 years and over), `agesAbove16` (The R16 classification is restricted to persons 16 years and over), `agesAbove18` (The R18 classification is restricted to persons 18 years and over), `restricted` (The R classification is restricted to a certain audience), `agesAbove16Restricted` (The RP16 classification requires viewers under 16 accompanied by a parent or an adult)
- `tv_rating` (String) TV rating selected for New Zealand. / TV content rating labels in New Zealand; possible values are: `allAllowed` (Default value, allow all TV shows content), `allBlocked` (Do not allow any TV shows content), `general` (The G classification excludes materials likely to harm children under 14), `parentalGuidance` (The PGR classification encourages parents and guardians to supervise younger viewers), `adults` (The AO classification is not suitable for children)


<a id="nestedatt--ios_general_device--media_content_rating_united_kingdom"></a>
### Nested Schema for `ios_general_device.media_content_rating_united_kingdom`

Required:

- `movie_rating` (String) Movies rating selected for United Kingdom. / Movies rating labels in United Kingdom; possible values are: `allAllowed` (Default value, allow all movies content), `allBlocked` (Do not allow any movies content), `general` (The U classification is suitable for all ages), `universalChildren` (The UC classification is suitable for pre-school children, an old rating label), `parentalGuidance` (The PG classification is suitable for mature), `agesAbove12Video` (12, video release suitable for 12 years and over), `agesAbove12Cinema` (12A, cinema release suitable for 12 years and over), `agesAbove15` (15, suitable only for 15 years and older), `adults` (Suitable only for adults)
- `tv_rating` (String) TV rating selected for United Kingdom. / TV content rating labels in United Kingdom; possible values are: `allAllowed` (Default value, allow all TV shows content), `allBlocked` (Do not allow any TV shows content), `caution` (Allowing TV contents with a warning message)


<a id="nestedatt--ios_general_device--media_content_rating_united_states"></a>
### Nested Schema for `ios_general_device.media_content_rating_united_states`

Required:

- `movie_rating` (String) Movies rating selected for United States. / Movies rating labels in United States; possible values are: `allAllowed` (Default value, allow all movies content), `allBlocked` (Do not allow any movies content), `general` (G, all ages admitted), `parentalGuidance` (PG, some material may not be suitable for children), `parentalGuidance13` (PG13, some material may be inappropriate for children under 13), `restricted` (R, viewers under 17 require accompanying parent or adult guardian), `adults` (NC17, adults only)
- `tv_rating` (String) TV rating selected for United States. / TV content rating labels in United States; possible values are: `allAllowed` (Default value, allow all TV shows content), `allBlocked` (Do not allow any TV shows content), `childrenAll` (TV-Y, all children), `childrenAbove7` (TV-Y7, children age 7 and above), `general` (TV-G, suitable for all ages), `parentalGuidance` (TV-PG, parental guidance), `childrenAbove14` (TV-14, children age 14 and above), `adults` (TV-MA, adults only)


<a id="nestedatt--ios_general_device--network_usage_rules"></a>
### Nested Schema for `ios_general_device.network_usage_rules`

Required:

- `cellular_data_block_when_roaming` (Boolean) If set to true, corresponding managed apps will not be allowed to use cellular data when roaming.
- `cellular_data_blocked` (Boolean) If set to true, corresponding managed apps will not be allowed to use cellular data at any time.

Optional:

- `managed_apps` (Attributes Set) Information about the managed apps that this rule is going to apply to. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_general_device--network_usage_rules--managed_apps))

<a id="nestedatt--ios_general_device--network_usage_rules--managed_apps"></a>
### Nested Schema for `ios_general_device.network_usage_rules.managed_apps`

Required:

- `name` (String) The application name

Optional:

- `app_id` (String) The application or bundle identifier of the application
- `app_store_url` (String) The Store URL of the application
- `publisher` (String) The publisher of the application




<a id="nestedatt--ios_update"></a>
### Nested Schema for `ios_update`

Optional:

- `active_hours_end` (String) Active Hours End (active hours mean the time window when updates install should not happen). The _provider_ default value is `"00:00:00.0000000"`.
- `active_hours_start` (String) Active Hours Start (active hours mean the time window when updates install should not happen). The _provider_ default value is `"00:00:00.0000000"`.
- `custom_update_time_windows` (Attributes Set) If update schedule type is set to use time window scheduling, custom time windows when updates will be scheduled. This collection can contain a maximum of 20 elements. / Custom update time window / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-customupdatetimewindow?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--ios_update--custom_update_time_windows))
- `desired_os_version` (String) If left unspecified, devices will update to the latest version of the OS.
- `enforced_software_update_delay_in_days` (Number) Days before software updates are visible to iOS devices ranging from 0 to 90 inclusive
- `scheduled_install_days` (Set of String) Days in week for which active hours are configured. This collection can contain a maximum of 7 elements. / Possible values for a weekday; possible values are: `sunday` (Sunday.), `monday` (Monday.), `tuesday` (Tuesday.), `wednesday` (Wednesday.), `thursday` (Thursday.), `friday` (Friday.), `saturday` (Saturday.). The _provider_ default value is `[]`.
- `update_schedule_type` (String) Update schedule type. / Update schedule type for iOS software updates; possible values are: `updateOutsideOfActiveHours` (Update outside of active hours.), `alwaysUpdate` (Always update.), `updateDuringTimeWindows` (Update during time windows.), `updateOutsideOfTimeWindows` (Update outside of time windows.). The _provider_ default value is `"alwaysUpdate"`.
- `utc_time_offset_in_minutes` (Number) UTC Time Offset indicated in minutes

Read-Only:

- `is_enabled` (Boolean) Is setting enabled in UI

<a id="nestedatt--ios_update--custom_update_time_windows"></a>
### Nested Schema for `ios_update.custom_update_time_windows`

Required:

- `end_day` (String) End day of the time window. / Possible values for a weekday; possible values are: `sunday` (Sunday.), `monday` (Monday.), `tuesday` (Tuesday.), `wednesday` (Wednesday.), `thursday` (Thursday.), `friday` (Friday.), `saturday` (Saturday.)
- `end_time` (String) End time of the time window
- `start_day` (String) Start day of the time window. / Possible values for a weekday; possible values are: `sunday` (Sunday.), `monday` (Monday.), `tuesday` (Tuesday.), `wednesday` (Wednesday.), `thursday` (Thursday.), `friday` (Friday.), `saturday` (Saturday.)
- `start_time` (String) Start time of the time window



<a id="nestedatt--ios_vpn"></a>
### Nested Schema for `ios_vpn`

Required:

- `authentication_method` (String) Authentication method for this VPN connection. / VPN Authentication Method; possible values are: `certificate` (Authenticate with a certificate.), `usernameAndPassword` (Use username and password for authentication.), `sharedSecret` (Use Shared Secret for Authentication.  Only valid for iOS IKEv2.), `derivedCredential` (Use Derived Credential for Authentication.), `azureAD` (Use Azure AD for authentication.)
- `connection_name` (String) Connection name displayed to the user.
- `connection_type` (String) Connection type. / Apple VPN connection type; possible values are: `ciscoAnyConnect` (Cisco AnyConnect.), `pulseSecure` (Pulse Secure.), `f5EdgeClient` (F5 Edge Client.), `dellSonicWallMobileConnect` (Dell SonicWALL Mobile Connection.), `checkPointCapsuleVpn` (Check Point Capsule VPN.), `customVpn` (Custom VPN.), `ciscoIPSec` (Cisco (IPSec).), `citrix` (Citrix.), `ciscoAnyConnectV2` (Cisco AnyConnect V2.), `paloAltoGlobalProtect` (Palo Alto Networks GlobalProtect.), `zscalerPrivateAccess` (Zscaler Private Access.), `f5Access2018` (F5 Access 2018.), `citrixSso` (Citrix Sso.), `paloAltoGlobalProtectV2` (Palo Alto Networks GlobalProtect V2.), `ikEv2` (IKEv2.), `alwaysOn` (AlwaysOn.), `microsoftTunnel` (Microsoft Tunnel.), `netMotionMobility` (NetMotion Mobility.), `microsoftProtect` (Microsoft Protect.)
- `server` (Attributes) VPN Server on the network. Make sure end users can access this network location. / VPN Server definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-vpnserver?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_vpn--server))

Optional:

- `associated_domains` (Set of String) Associated Domains. The _provider_ default value is `[]`.
- `cloud_name` (String) Zscaler only. Zscaler cloud which the user is assigned to.
- `custom_data` (Attributes Set) Custom data when connection type is set to Custom VPN. Use this field to enable functionality not supported by Intune, but available in your VPN solution. Contact your VPN vendor to learn how to add these key/value pairs. This collection can contain a maximum of 25 elements. / Key Value definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyvalue?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--ios_vpn--custom_data))
- `disable_on_demand_user_override` (Boolean) Toggle to prevent user from disabling automatic VPN in the Settings app
- `disconnect_on_idle` (Boolean) Whether to disconnect after on-demand connection idles
- `disconnect_on_idle_timer_in_seconds` (Number) The length of time in seconds to wait before disconnecting an on-demand connection. Valid values 0 to 65535
- `enable_per_app` (Boolean) Setting this to true creates Per-App VPN payload which can later be associated with Apps that can trigger this VPN conneciton on the end user's iOS device.
- `enable_split_tunneling` (Boolean) Send all network traffic through VPN. The _provider_ default value is `false`.
- `exclude_list` (Set of String) Zscaler only. List of network addresses which are not sent through the Zscaler cloud. The _provider_ default value is `[]`.
- `excluded_domains` (Set of String) Domains that are accessed through the public internet instead of through VPN, even when per-app VPN is activated. The _provider_ default value is `[]`.
- `identifier` (String) Identifier provided by VPN vendor when connection type is set to Custom VPN. For example: Cisco AnyConnect uses an identifier of the form com.cisco.anyconnect.applevpn.plugin
- `login_group_or_domain` (String) Login group or domain when connection type is set to Dell SonicWALL Mobile Connection.
- `microsoft_tunnel_site_id` (String) Microsoft Tunnel site ID.
- `on_demand_rules` (Attributes Set) On-Demand Rules. This collection can contain a maximum of 500 elements. / VPN On-Demand Rule definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-vpnondemandrule?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--ios_vpn--on_demand_rules))
- `opt_in_to_device_id_sharing` (Boolean) Opt-In to sharing the device's Id to third-party vpn clients for use during network access control validation.
- `provider_type` (String) Provider type for per-app VPN. / Provider type for per-app VPN; possible values are: `notConfigured` (Tunnel traffic is not explicitly configured.), `appProxy` (Tunnel traffic at the application layer.), `packetTunnel` (Tunnel traffic at the IP layer.)
- `proxy_server` (Attributes) Proxy Server. / VPN Proxy Server. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-vpnproxyserver?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_vpn--proxy_server))
- `realm` (String) Realm when connection type is set to Pulse Secure.
- `role` (String) Role when connection type is set to Pulse Secure.
- `safari_domains` (Set of String) Safari domains when this VPN per App setting is enabled. In addition to the apps associated with this VPN, Safari domains specified here will also be able to trigger this VPN connection. The _provider_ default value is `[]`.
- `strict_enforcement` (Boolean) Zscaler only. Blocks network traffic until the user signs into Zscaler app. "True" means traffic is blocked.
- `targeted_mobile_apps` (Attributes Set) Targeted mobile apps. This collection can contain a maximum of 500 elements. / Represents an app in the list of managed applications / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-applistitem?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--ios_vpn--targeted_mobile_apps))
- `user_domain` (String) Zscaler only. Enter a static domain to pre-populate the login field with in the Zscaler app. If this is left empty, the user's Azure Active Directory domain will be used instead.

<a id="nestedatt--ios_vpn--server"></a>
### Nested Schema for `ios_vpn.server`

Required:

- `address` (String) Address (IP address, FQDN or URL)

Optional:

- `description` (String) Description.
- `is_default_server` (Boolean) Default server. The _provider_ default value is `true`.


<a id="nestedatt--ios_vpn--custom_data"></a>
### Nested Schema for `ios_vpn.custom_data`

Optional:

- `key` (String) Key.
- `value` (String) Value.


<a id="nestedatt--ios_vpn--on_demand_rules"></a>
### Nested Schema for `ios_vpn.on_demand_rules`

Required:

- `action` (String) Action. / VPN On-Demand Rule Connection Action; possible values are: `connect` (Connect.), `evaluateConnection` (Evaluate Connection.), `ignore` (Ignore.), `disconnect` (Disconnect.)

Optional:

- `dns_search_domains` (Set of String) DNS Search Domains. The _provider_ default value is `[]`.
- `dns_server_address_match` (Set of String) DNS Search Server Address. The _provider_ default value is `[]`.
- `domain_action` (String) Domain Action (Only applicable when Action is evaluate connection). / VPN On-Demand Rule Connection Domain Action; possible values are: `connectIfNeeded` (Connect if needed.), `neverConnect` (Never connect.). The _provider_ default value is `"connectIfNeeded"`.
- `domains` (Set of String) Domains (Only applicable when Action is evaluate connection). The _provider_ default value is `[]`.
- `interface_type_match` (String) Network interface to trigger VPN. / VPN On-Demand Rule Connection network interface type; possible values are: `notConfigured` (NotConfigured), `ethernet` (Ethernet.), `wiFi` (WiFi.), `cellular` (Cellular.). The _provider_ default value is `"notConfigured"`.
- `probe_required_url` (String) Probe Required Url (Only applicable when Action is evaluate connection and DomainAction is connect if needed).
- `probe_url` (String) A URL to probe. If this URL is successfully fetched (returning a 200 HTTP status code) without redirection, this rule matches.
- `ssids` (Set of String) Network Service Set Identifiers (SSIDs). The _provider_ default value is `[]`.


<a id="nestedatt--ios_vpn--proxy_server"></a>
### Nested Schema for `ios_vpn.proxy_server`

Optional:

- `address` (String) Address.
- `automatic_configuration_script_url` (String) Proxy's automatic configuration script url.
- `port` (Number) Port. Valid values 0 to 65535


<a id="nestedatt--ios_vpn--targeted_mobile_apps"></a>
### Nested Schema for `ios_vpn.targeted_mobile_apps`

Required:

- `name` (String) The application name

Optional:

- `app_id` (String) The application or bundle identifier of the application
- `app_store_url` (String) The Store URL of the application
- `publisher` (String) The publisher of the application



<a id="nestedatt--macos_custom"></a>
### Nested Schema for `macos_custom`

Required:

- `deployment_channel` (String) Indicates the channel used to deploy the configuration profile. Available choices are DeviceChannel, UserChannel. / Indicates the channel used to deploy the configuration profile. Available choices are DeviceChannel, UserChannel; possible values are: `deviceChannel` (Send payload down over Device Channel.), `userChannel` (Send payload down over User Channel.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.)
- `file_name` (String) Payload file name (*.mobileconfig
- `name` (String) Name that is displayed to the user.
- `payload` (String) Payload. (UTF8 encoded byte array)


<a id="nestedatt--macos_custom_app"></a>
### Nested Schema for `macos_custom_app`

Required:

- `bundle_id` (String) Bundle id for targeting.
- `configuration_xml` (String) Configuration xml. (UTF8 encoded byte array)
- `file_name` (String) Configuration file name (*.plist


<a id="nestedatt--macos_device_features"></a>
### Nested Schema for `macos_device_features`

Optional:

- `admin_show_host_info` (Boolean) Whether to show admin host information on the login window. The _provider_ default value is `false`.
- `airprint_destinations` (Attributes Set) An array of AirPrint printers that should always be shown. This collection can contain a maximum of 500 elements. / Represents an AirPrint destination. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-airprintdestination?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--macos_device_features--airprint_destinations))
- `app_associated_domains` (Attributes Set) Gets or sets a list that maps apps to their associated domains. Application identifiers must be unique. This collection can contain a maximum of 500 elements. / A mapping of application identifiers to associated domains. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macosassociateddomainsitem?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--macos_device_features--app_associated_domains))
- `authorized_users_list_hidden` (Boolean) Whether to show the name and password dialog or a list of users on the login window. The _provider_ default value is `false`.
- `authorized_users_list_hide_admin_users` (Boolean) Whether to hide admin users in the authorized users list on the login window. The _provider_ default value is `false`.
- `authorized_users_list_hide_local_users` (Boolean) Whether to show only network and system users in the authorized users list on the login window. The _provider_ default value is `false`.
- `authorized_users_list_hide_mobile_accounts` (Boolean) Whether to hide mobile users in the authorized users list on the login window. The _provider_ default value is `false`.
- `authorized_users_list_include_network_users` (Boolean) Whether to show network users in the authorized users list on the login window. The _provider_ default value is `false`.
- `authorized_users_list_show_other_managed_users` (Boolean) Whether to show other users in the authorized users list on the login window. The _provider_ default value is `false`.
- `auto_launch_items` (Attributes Set) List of applications, files, folders, and other items to launch when the user logs in. This collection can contain a maximum of 500 elements. / Represents an app in the list of macOS launch items / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macoslaunchitem?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--macos_device_features--auto_launch_items))
- `console_access_disabled` (Boolean) Whether the Other user will disregard use of the `console` special user name. The _provider_ default value is `false`.
- `content_caching_block_deletion` (Boolean) Prevents content caches from purging content to free up disk space for other apps. The _provider_ default value is `false`.
- `content_caching_client_listen_ranges` (Attributes Set) A list of custom IP ranges content caches will use to listen for clients. This collection can contain a maximum of 500 elements. / IP range base class for representing IPV4, IPV6 address ranges / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-iprange?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--macos_device_features--content_caching_client_listen_ranges))
- `content_caching_client_policy` (String) Determines the method in which content caching servers will listen for clients. / Determines which clients a content cache will serve; possible values are: `notConfigured` (Defaults to clients in local network.), `clientsInLocalNetwork` (Content caches will provide content to devices only in their immediate local network.), `clientsWithSamePublicIpAddress` (Content caches will provide content to devices that share the same public IP address.), `clientsInCustomLocalNetworks` (Content caches will provide content to devices in contentCachingClientListenRanges.), `clientsInCustomLocalNetworksWithFallback` (Content caches will provide content to devices in contentCachingClientListenRanges, contentCachingPeerListenRanges, and contentCachingParents.). The _provider_ default value is `"notConfigured"`.
- `content_caching_data_path` (String) The path to the directory used to store cached content. The value must be (or end with) /Library/Application Support/Apple/AssetCache/Data
- `content_caching_disable_connection_sharing` (Boolean) Disables internet connection sharing. The _provider_ default value is `false`.
- `content_caching_enabled` (Boolean) Enables content caching and prevents it from being disabled by the user. The _provider_ default value is `false`.
- `content_caching_force_connection_sharing` (Boolean) Forces internet connection sharing. contentCachingDisableConnectionSharing overrides this setting. The _provider_ default value is `false`.
- `content_caching_keep_awake` (Boolean) Prevent the device from sleeping if content caching is enabled. The _provider_ default value is `false`.
- `content_caching_log_client_identities` (Boolean) Enables logging of IP addresses and ports of clients that request cached content. The _provider_ default value is `false`.
- `content_caching_max_size_bytes` (Number) The maximum number of bytes of disk space that will be used for the content cache. A value of 0 (default) indicates unlimited disk space.
- `content_caching_parent_selection_policy` (String) Determines the method in which content caching servers will select parents if multiple are present. / Determines how content caches select a parent cache; possible values are: `notConfigured` (Defaults to round-robin strategy.), `roundRobin` (Rotate through the parents in order. Use this policy for load balancing.), `firstAvailable` (Always use the first available parent in the Parents list. Use this policy to designate permanent primary, secondary, and subsequent parents.), `urlPathHash` (Hash the path part of the requested URL so that the same parent is always used for the same URL. This is useful for maximizing the size of the combined caches of the parents.), `random` (Choose a parent at random. Use this policy for load balancing.), `stickyAvailable` (Use the first available parent that is available in the Parents list until it becomes unavailable, then advance to the next one. Use this policy for designating floating primary, secondary, and subsequent parents.). The _provider_ default value is `"notConfigured"`.
- `content_caching_parents` (Set of String) A list of IP addresses representing parent content caches. The _provider_ default value is `[]`.
- `content_caching_peer_filter_ranges` (Attributes Set) A list of custom IP ranges content caches will use to query for content from peers caches. This collection can contain a maximum of 500 elements. / IP range base class for representing IPV4, IPV6 address ranges / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-iprange?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--macos_device_features--content_caching_peer_filter_ranges))
- `content_caching_peer_listen_ranges` (Attributes Set) A list of custom IP ranges content caches will use to listen for peer caches. This collection can contain a maximum of 500 elements. / IP range base class for representing IPV4, IPV6 address ranges / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-iprange?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--macos_device_features--content_caching_peer_listen_ranges))
- `content_caching_peer_policy` (String) Determines the method in which content caches peer with other caches. / Determines which content caches other content caches will peer with; possible values are: `notConfigured` (Defaults to peers in local network.), `peersInLocalNetwork` (Content caches will only peer with caches in their immediate local network.), `peersWithSamePublicIpAddress` (Content caches will only peer with caches that share the same public IP address.), `peersInCustomLocalNetworks` (Content caches will use contentCachingPeerFilterRanges and contentCachingPeerListenRanges to determine which caches to peer with.). The _provider_ default value is `"notConfigured"`.
- `content_caching_port` (Number) Sets the port used for content caching. If the value is 0, a random available port will be selected. Valid values 0 to 65535
- `content_caching_public_ranges` (Attributes Set) A list of custom IP ranges that Apple's content caching service should use to match clients to content caches. This collection can contain a maximum of 500 elements. / IP range base class for representing IPV4, IPV6 address ranges / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-iprange?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--macos_device_features--content_caching_public_ranges))
- `content_caching_show_alerts` (Boolean) Display content caching alerts as system notifications. The _provider_ default value is `false`.
- `content_caching_type` (String) Determines what type of content is allowed to be cached by Apple's content caching service. / Indicates the type of content allowed to be cached by Apple's content caching service; possible values are: `notConfigured` (Default. Both user iCloud data and non-iCloud data will be cached.), `userContentOnly` (Allow Apple's content caching service to cache user iCloud data.), `sharedContentOnly` (Allow Apple's content caching service to cache non-iCloud data (e.g. app and software updates).). The _provider_ default value is `"notConfigured"`.
- `login_window_text` (String) Custom text to be displayed on the login window. The _provider_ default value is `""`.
- `logout_disabled_while_logged_in` (Boolean) Whether the Log Out menu item on the login window will be disabled while the user is logged in. The _provider_ default value is `false`.
- `macos_single_sign_on_extension` (Attributes) Gets or sets a single sign-on extension profile. / An abstract base class for all macOS-specific single sign-on extension types. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macossinglesignonextension?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension))
- `poweroff_disabled_while_logged_in` (Boolean) Whether the Power Off menu item on the login window will be disabled while the user is logged in. The _provider_ default value is `false`.
- `restart_disabled` (Boolean) Whether to hide the Restart button item on the login window. The _provider_ default value is `false`.
- `restart_disabled_while_logged_in` (Boolean) Whether the Restart menu item on the login window will be disabled while the user is logged in. The _provider_ default value is `false`.
- `screen_lock_disable_immediate` (Boolean) Whether to disable the immediate screen lock functions. The _provider_ default value is `false`.
- `shutdown_disabled` (Boolean) Whether to hide the Shut Down button item on the login window. The _provider_ default value is `false`.
- `shutdown_disabled_while_logged_in` (Boolean) Whether the Shut Down menu item on the login window will be disabled while the user is logged in. The _provider_ default value is `false`.
- `sleep_disabled` (Boolean) Whether to hide the Sleep menu item on the login window. The _provider_ default value is `false`.

<a id="nestedatt--macos_device_features--airprint_destinations"></a>
### Nested Schema for `macos_device_features.airprint_destinations`

Required:

- `force_tls` (Boolean) If true AirPrint connections are secured by Transport Layer Security (TLS). Default is false. Available in iOS 11.0 and later.
- `ip_address` (String) The IP Address of the AirPrint destination.
- `resource_path` (String) The Resource Path associated with the printer. This corresponds to the rp parameter of the _ipps.tcp Bonjour record. For example: printers/Canon_MG5300_series, printers/Xerox_Phaser_7600, ipp/print, Epson_IPP_Printer.

Optional:

- `port` (Number) The listening port of the AirPrint destination. If this key is not specified AirPrint will use the default port. Available in iOS 11.0 and later.


<a id="nestedatt--macos_device_features--app_associated_domains"></a>
### Nested Schema for `macos_device_features.app_associated_domains`

Required:

- `application_identifier` (String) The application identifier of the app to associate domains with.
- `direct_downloads_enabled` (Boolean) Determines whether data should be downloaded directly or via a CDN.
- `domains` (Set of String) The list of domains to associate.


<a id="nestedatt--macos_device_features--auto_launch_items"></a>
### Nested Schema for `macos_device_features.auto_launch_items`

Required:

- `hide` (Boolean) Whether or not to hide the item from the Users and Groups List.
- `path` (String) Path to the launch item.


<a id="nestedatt--macos_device_features--content_caching_client_listen_ranges"></a>
### Nested Schema for `macos_device_features.content_caching_client_listen_ranges`

Optional:

- `v4` (Attributes) IPv4 Range definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-ipv4range?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--content_caching_client_listen_ranges--v4))
- `v4_cidr` (Attributes) Represents an IPv4 range using the Classless Inter-Domain Routing (CIDR) notation. / https://learn.microsoft.com/en-us/graph/api/resources/ipv4cidrrange?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--content_caching_client_listen_ranges--v4_cidr))
- `v6` (Attributes) IPv6 Range definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-ipv6range?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--content_caching_client_listen_ranges--v6))
- `v6_cidr` (Attributes) Represents an IPv6 range using the Classless Inter-Domain Routing (CIDR) notation. / https://learn.microsoft.com/en-us/graph/api/resources/ipv6cidrrange?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--content_caching_client_listen_ranges--v6_cidr))

<a id="nestedatt--macos_device_features--content_caching_client_listen_ranges--v4"></a>
### Nested Schema for `macos_device_features.content_caching_client_listen_ranges.v4`

Required:

- `lower_address` (String) Lower address.
- `upper_address` (String) Upper address.


<a id="nestedatt--macos_device_features--content_caching_client_listen_ranges--v4_cidr"></a>
### Nested Schema for `macos_device_features.content_caching_client_listen_ranges.v4_cidr`

Required:

- `cidr_address` (String) IPv4 address in CIDR notation. Not nullable.


<a id="nestedatt--macos_device_features--content_caching_client_listen_ranges--v6"></a>
### Nested Schema for `macos_device_features.content_caching_client_listen_ranges.v6`

Required:

- `lower_address` (String) Lower address.
- `upper_address` (String) Upper address.


<a id="nestedatt--macos_device_features--content_caching_client_listen_ranges--v6_cidr"></a>
### Nested Schema for `macos_device_features.content_caching_client_listen_ranges.v6_cidr`

Required:

- `cidr_address` (String) IPv6 address in CIDR notation. Not nullable.



<a id="nestedatt--macos_device_features--content_caching_peer_filter_ranges"></a>
### Nested Schema for `macos_device_features.content_caching_peer_filter_ranges`

Optional:

- `v4` (Attributes) IPv4 Range definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-ipv4range?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--content_caching_peer_filter_ranges--v4))
- `v4_cidr` (Attributes) Represents an IPv4 range using the Classless Inter-Domain Routing (CIDR) notation. / https://learn.microsoft.com/en-us/graph/api/resources/ipv4cidrrange?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--content_caching_peer_filter_ranges--v4_cidr))
- `v6` (Attributes) IPv6 Range definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-ipv6range?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--content_caching_peer_filter_ranges--v6))
- `v6_cidr` (Attributes) Represents an IPv6 range using the Classless Inter-Domain Routing (CIDR) notation. / https://learn.microsoft.com/en-us/graph/api/resources/ipv6cidrrange?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--content_caching_peer_filter_ranges--v6_cidr))

<a id="nestedatt--macos_device_features--content_caching_peer_filter_ranges--v4"></a>
### Nested Schema for `macos_device_features.content_caching_peer_filter_ranges.v4`

Required:

- `lower_address` (String) Lower address.
- `upper_address` (String) Upper address.


<a id="nestedatt--macos_device_features--content_caching_peer_filter_ranges--v4_cidr"></a>
### Nested Schema for `macos_device_features.content_caching_peer_filter_ranges.v4_cidr`

Required:

- `cidr_address` (String) IPv4 address in CIDR notation. Not nullable.


<a id="nestedatt--macos_device_features--content_caching_peer_filter_ranges--v6"></a>
### Nested Schema for `macos_device_features.content_caching_peer_filter_ranges.v6`

Required:

- `lower_address` (String) Lower address.
- `upper_address` (String) Upper address.


<a id="nestedatt--macos_device_features--content_caching_peer_filter_ranges--v6_cidr"></a>
### Nested Schema for `macos_device_features.content_caching_peer_filter_ranges.v6_cidr`

Required:

- `cidr_address` (String) IPv6 address in CIDR notation. Not nullable.



<a id="nestedatt--macos_device_features--content_caching_peer_listen_ranges"></a>
### Nested Schema for `macos_device_features.content_caching_peer_listen_ranges`

Optional:

- `v4` (Attributes) IPv4 Range definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-ipv4range?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--content_caching_peer_listen_ranges--v4))
- `v4_cidr` (Attributes) Represents an IPv4 range using the Classless Inter-Domain Routing (CIDR) notation. / https://learn.microsoft.com/en-us/graph/api/resources/ipv4cidrrange?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--content_caching_peer_listen_ranges--v4_cidr))
- `v6` (Attributes) IPv6 Range definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-ipv6range?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--content_caching_peer_listen_ranges--v6))
- `v6_cidr` (Attributes) Represents an IPv6 range using the Classless Inter-Domain Routing (CIDR) notation. / https://learn.microsoft.com/en-us/graph/api/resources/ipv6cidrrange?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--content_caching_peer_listen_ranges--v6_cidr))

<a id="nestedatt--macos_device_features--content_caching_peer_listen_ranges--v4"></a>
### Nested Schema for `macos_device_features.content_caching_peer_listen_ranges.v4`

Required:

- `lower_address` (String) Lower address.
- `upper_address` (String) Upper address.


<a id="nestedatt--macos_device_features--content_caching_peer_listen_ranges--v4_cidr"></a>
### Nested Schema for `macos_device_features.content_caching_peer_listen_ranges.v4_cidr`

Required:

- `cidr_address` (String) IPv4 address in CIDR notation. Not nullable.


<a id="nestedatt--macos_device_features--content_caching_peer_listen_ranges--v6"></a>
### Nested Schema for `macos_device_features.content_caching_peer_listen_ranges.v6`

Required:

- `lower_address` (String) Lower address.
- `upper_address` (String) Upper address.


<a id="nestedatt--macos_device_features--content_caching_peer_listen_ranges--v6_cidr"></a>
### Nested Schema for `macos_device_features.content_caching_peer_listen_ranges.v6_cidr`

Required:

- `cidr_address` (String) IPv6 address in CIDR notation. Not nullable.



<a id="nestedatt--macos_device_features--content_caching_public_ranges"></a>
### Nested Schema for `macos_device_features.content_caching_public_ranges`

Optional:

- `v4` (Attributes) IPv4 Range definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-ipv4range?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--content_caching_public_ranges--v4))
- `v4_cidr` (Attributes) Represents an IPv4 range using the Classless Inter-Domain Routing (CIDR) notation. / https://learn.microsoft.com/en-us/graph/api/resources/ipv4cidrrange?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--content_caching_public_ranges--v4_cidr))
- `v6` (Attributes) IPv6 Range definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-ipv6range?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--content_caching_public_ranges--v6))
- `v6_cidr` (Attributes) Represents an IPv6 range using the Classless Inter-Domain Routing (CIDR) notation. / https://learn.microsoft.com/en-us/graph/api/resources/ipv6cidrrange?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--content_caching_public_ranges--v6_cidr))

<a id="nestedatt--macos_device_features--content_caching_public_ranges--v4"></a>
### Nested Schema for `macos_device_features.content_caching_public_ranges.v4`

Required:

- `lower_address` (String) Lower address.
- `upper_address` (String) Upper address.


<a id="nestedatt--macos_device_features--content_caching_public_ranges--v4_cidr"></a>
### Nested Schema for `macos_device_features.content_caching_public_ranges.v4_cidr`

Required:

- `cidr_address` (String) IPv4 address in CIDR notation. Not nullable.


<a id="nestedatt--macos_device_features--content_caching_public_ranges--v6"></a>
### Nested Schema for `macos_device_features.content_caching_public_ranges.v6`

Required:

- `lower_address` (String) Lower address.
- `upper_address` (String) Upper address.


<a id="nestedatt--macos_device_features--content_caching_public_ranges--v6_cidr"></a>
### Nested Schema for `macos_device_features.content_caching_public_ranges.v6_cidr`

Required:

- `cidr_address` (String) IPv6 address in CIDR notation. Not nullable.



<a id="nestedatt--macos_device_features--macos_single_sign_on_extension"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension`

Optional:

- `azure_ad` (Attributes) Represents an Azure AD-type Single Sign-On extension profile for macOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macosazureadsinglesignonextension?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension--azure_ad))
- `credential` (Attributes) Represents a Credential-type Single Sign-On extension profile for macOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macoscredentialsinglesignonextension?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension--credential))
- `kerberos` (Attributes) Represents a Kerberos-type Single Sign-On extension profile for macOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macoskerberossinglesignonextension?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension--kerberos))
- `redirect` (Attributes) Represents a Redirect-type Single Sign-On extension profile for macOS devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macosredirectsinglesignonextension?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension--redirect))

<a id="nestedatt--macos_device_features--macos_single_sign_on_extension--azure_ad"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension.azure_ad`

Optional:

- `bundle_id_access_control_list` (Set of String) An optional list of additional bundle IDs allowed to use the AAD extension for single sign-on. The _provider_ default value is `[]`.
- `configurations` (Attributes Set) Gets or sets a list of typed key-value pairs used to configure Credential-type profiles. This collection can contain a maximum of 500 elements. / A key-value pair with a string key and a typed value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keytypedvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension--azure_ad--configurations))
- `enable_shared_device_mode` (Boolean) Enables or disables shared device mode. The _provider_ default value is `false`.

<a id="nestedatt--macos_device_features--macos_single_sign_on_extension--azure_ad--configurations"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension.azure_ad.configurations`

Required:

- `key` (String) The string key of the key-value pair.

Optional:

- `boolean` (Attributes) A key-value pair with a string key and a Boolean value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keybooleanvaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension--azure_ad--configurations--boolean))
- `integer` (Attributes) A key-value pair with a string key and an integer value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyintegervaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension--azure_ad--configurations--integer))
- `real` (Attributes) A key-value pair with a string key and a real (floating-point) value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyrealvaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension--azure_ad--configurations--real))
- `string` (Attributes) A key-value pair with a string key and a string value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keystringvaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension--azure_ad--configurations--string))

<a id="nestedatt--macos_device_features--macos_single_sign_on_extension--azure_ad--configurations--boolean"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension.azure_ad.configurations.boolean`

Required:

- `value` (Boolean) The Boolean value of the key-value pair.


<a id="nestedatt--macos_device_features--macos_single_sign_on_extension--azure_ad--configurations--integer"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension.azure_ad.configurations.integer`

Required:

- `value` (Number) The integer value of the key-value pair.


<a id="nestedatt--macos_device_features--macos_single_sign_on_extension--azure_ad--configurations--real"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension.azure_ad.configurations.real`

Required:

- `value` (Number) The real (floating-point) value of the key-value pair.


<a id="nestedatt--macos_device_features--macos_single_sign_on_extension--azure_ad--configurations--string"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension.azure_ad.configurations.string`

Required:

- `value` (String) The string value of the key-value pair.




<a id="nestedatt--macos_device_features--macos_single_sign_on_extension--credential"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension.credential`

Required:

- `extension_identifier` (String) Gets or sets the bundle ID of the app extension that performs SSO for the specified URLs.
- `realm` (String) Gets or sets the case-sensitive realm name for this profile.
- `team_identifier` (String) Gets or sets the team ID of the app extension that performs SSO for the specified URLs.

Optional:

- `configurations` (Attributes Set) Gets or sets a list of typed key-value pairs used to configure Credential-type profiles. This collection can contain a maximum of 500 elements. / A key-value pair with a string key and a typed value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keytypedvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension--credential--configurations))
- `domains` (Set of String) Gets or sets a list of hosts or domain names for which the app extension performs SSO. The _provider_ default value is `[]`.

<a id="nestedatt--macos_device_features--macos_single_sign_on_extension--credential--configurations"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension.credential.configurations`

Required:

- `key` (String) The string key of the key-value pair.

Optional:

- `boolean` (Attributes) A key-value pair with a string key and a Boolean value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keybooleanvaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension--credential--configurations--boolean))
- `integer` (Attributes) A key-value pair with a string key and an integer value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyintegervaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension--credential--configurations--integer))
- `real` (Attributes) A key-value pair with a string key and a real (floating-point) value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyrealvaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension--credential--configurations--real))
- `string` (Attributes) A key-value pair with a string key and a string value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keystringvaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension--credential--configurations--string))

<a id="nestedatt--macos_device_features--macos_single_sign_on_extension--credential--configurations--boolean"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension.credential.configurations.boolean`

Required:

- `value` (Boolean) The Boolean value of the key-value pair.


<a id="nestedatt--macos_device_features--macos_single_sign_on_extension--credential--configurations--integer"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension.credential.configurations.integer`

Required:

- `value` (Number) The integer value of the key-value pair.


<a id="nestedatt--macos_device_features--macos_single_sign_on_extension--credential--configurations--real"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension.credential.configurations.real`

Required:

- `value` (Number) The real (floating-point) value of the key-value pair.


<a id="nestedatt--macos_device_features--macos_single_sign_on_extension--credential--configurations--string"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension.credential.configurations.string`

Required:

- `value` (String) The string value of the key-value pair.




<a id="nestedatt--macos_device_features--macos_single_sign_on_extension--kerberos"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension.kerberos`

Required:

- `block_active_directory_site_auto_discovery` (Boolean) Enables or disables whether the Kerberos extension can automatically determine its site name.
- `block_automatic_login` (Boolean) Enables or disables Keychain usage.
- `credentials_cache_monitored` (Boolean) When set to True, the credential is requested on the next matching Kerberos challenge or network state change. When the credential is expired or missing, a new credential is created. Available for devices running macOS versions 12 and later.
- `is_default_realm` (Boolean) When true, this profile's realm will be selected as the default. Necessary if multiple Kerberos-type profiles are configured.
- `kerberos_apps_in_bundle_id_acl_included` (Boolean) When set to True, the Kerberos extension allows any apps entered with the app bundle ID, managed apps, and standard Kerberos utilities, such as TicketViewer and klist, to access and use the credential. Available for devices running macOS versions 12 and later.
- `managed_apps_in_bundle_id_acl_included` (Boolean) When set to True, the Kerberos extension allows managed apps, and any apps entered with the app bundle ID to access the credential. When set to False, the Kerberos extension allows all apps to access the credential. Available for devices running iOS and iPadOS versions 14 and later.
- `password_block_modification` (Boolean) Enables or disables password changes.
- `password_enable_local_sync` (Boolean) Enables or disables password syncing. This won't affect users logged in with a mobile account on macOS.
- `password_require_active_directory_complexity` (Boolean) Enables or disables whether passwords must meet Active Directory's complexity requirements.
- `realm` (String) Gets or sets the case-sensitive realm name for this profile.
- `require_user_presence` (Boolean) Gets or sets whether to require authentication via Touch ID, Face ID, or a passcode to access the keychain entry.
- `tls_for_ldap_required` (Boolean) When set to True, LDAP connections are required to use Transport Layer Security (TLS). Available for devices running macOS versions 11 and later.
- `user_setup_delayed` (Boolean) When set to True, the user isn’t prompted to set up the Kerberos extension until the extension is enabled by the admin, or a Kerberos challenge is received. Available for devices running macOS versions 11 and later.

Optional:

- `active_directory_site_code` (String) Gets or sets the Active Directory site.
- `cache_name` (String) Gets or sets the Generic Security Services name of the Kerberos cache to use for this profile.
- `credential_bundle_id_access_control_list` (Set of String) Gets or sets a list of app Bundle IDs allowed to access the Kerberos Ticket Granting Ticket.
- `domain_realms` (Set of String) Gets or sets a list of realms for custom domain-realm mapping. Realms are case sensitive.
- `domains` (Set of String) Gets or sets a list of hosts or domain names for which the app extension performs SSO.
- `mode_credential_used` (String) Select how other processes use the Kerberos Extension credential.
- `password_change_url` (String) Gets or sets the URL that the user will be sent to when they initiate a password change.
- `password_expiration_days` (Number) Overrides the default password expiration in days. For most domains, this value is calculated automatically.
- `password_expiration_notification_days` (Number) Gets or sets the number of days until the user is notified that their password will expire (default is 15).
- `password_minimum_age_days` (Number) Gets or sets the minimum number of days until a user can change their password again.
- `password_minimum_length` (Number) Gets or sets the minimum length of a password.
- `password_previous_password_block_count` (Number) Gets or sets the number of previous passwords to block.
- `password_requirements_description` (String) Gets or sets a description of the password complexity requirements.
- `preferred_kdcs` (Set of String) Add creates an ordered list of preferred Key Distribution Centers (KDCs) to use for Kerberos traffic. This list is used when the servers are not discoverable using DNS. When the servers are discoverable, the list is used for both connectivity checks, and used first for Kerberos traffic. If the servers don’t respond, then the device uses DNS discovery. Delete removes an existing list, and devices use DNS discovery. Available for devices running macOS versions 12 and later.
- `sign_in_help_text` (String) Text displayed to the user at the Kerberos sign in window. Available for devices running iOS and iPadOS versions 14 and later.
- `user_principal_name` (String) Gets or sets the principle user name to use for this profile. The realm name does not need to be included.
- `username_label_custom` (String) This label replaces the user name shown in the Kerberos extension. You can enter a name to match the name of your company or organization. Available for devices running macOS versions 11 and later.


<a id="nestedatt--macos_device_features--macos_single_sign_on_extension--redirect"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension.redirect`

Required:

- `extension_identifier` (String) Gets or sets the bundle ID of the app extension that performs SSO for the specified URLs.
- `team_identifier` (String) Gets or sets the team ID of the app extension that performs SSO for the specified URLs.
- `url_prefixes` (Set of String) One or more URL prefixes of identity providers on whose behalf the app extension performs single sign-on. URLs must begin with http:// or https://. All URL prefixes must be unique for all profiles.

Optional:

- `configurations` (Attributes Set) Gets or sets a list of typed key-value pairs used to configure Credential-type profiles. This collection can contain a maximum of 500 elements. / A key-value pair with a string key and a typed value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keytypedvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension--redirect--configurations))

<a id="nestedatt--macos_device_features--macos_single_sign_on_extension--redirect--configurations"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension.redirect.configurations`

Required:

- `key` (String) The string key of the key-value pair.

Optional:

- `boolean` (Attributes) A key-value pair with a string key and a Boolean value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keybooleanvaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension--redirect--configurations--boolean))
- `integer` (Attributes) A key-value pair with a string key and an integer value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyintegervaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension--redirect--configurations--integer))
- `real` (Attributes) A key-value pair with a string key and a real (floating-point) value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyrealvaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension--redirect--configurations--real))
- `string` (Attributes) A key-value pair with a string key and a string value. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keystringvaluepair?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_device_features--macos_single_sign_on_extension--redirect--configurations--string))

<a id="nestedatt--macos_device_features--macos_single_sign_on_extension--redirect--configurations--boolean"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension.redirect.configurations.boolean`

Required:

- `value` (Boolean) The Boolean value of the key-value pair.


<a id="nestedatt--macos_device_features--macos_single_sign_on_extension--redirect--configurations--integer"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension.redirect.configurations.integer`

Required:

- `value` (Number) The integer value of the key-value pair.


<a id="nestedatt--macos_device_features--macos_single_sign_on_extension--redirect--configurations--real"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension.redirect.configurations.real`

Required:

- `value` (Number) The real (floating-point) value of the key-value pair.


<a id="nestedatt--macos_device_features--macos_single_sign_on_extension--redirect--configurations--string"></a>
### Nested Schema for `macos_device_features.macos_single_sign_on_extension.redirect.configurations.string`

Required:

- `value` (String) The string value of the key-value pair.






<a id="nestedatt--macos_extensions"></a>
### Nested Schema for `macos_extensions`

Optional:

- `kernel_extension_allowed_team_identifiers` (Set of String) All kernel extensions validly signed by the team identifiers in this list will be allowed to load. The _provider_ default value is `[]`.
- `kernel_extension_overrides_allowed` (Boolean) If set to true, users can approve additional kernel extensions not explicitly allowed by configurations profiles. The _provider_ default value is `false`.
- `kernel_extensions_allowed` (Attributes Set) A list of kernel extensions that will be allowed to load. . This collection can contain a maximum of 500 elements. / Represents a specific macOS kernel extension. A macOS kernel extension can be described by its team identifier plus its bundle identifier. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macoskernelextension?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--macos_extensions--kernel_extensions_allowed))
- `system_extensions_allowed` (Attributes Set) Gets or sets a list of allowed macOS system extensions. This collection can contain a maximum of 500 elements. / Represents a specific macOS system extension. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macossystemextension?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--macos_extensions--system_extensions_allowed))
- `system_extensions_allowed_team_identifiers` (Set of String) Gets or sets a list of allowed team identifiers. Any system extension signed with any of the specified team identifiers will be approved. The _provider_ default value is `[]`.
- `system_extensions_allowed_types` (Attributes Set) Gets or sets a list of allowed macOS system extension types. This collection can contain a maximum of 500 elements. / Represents a mapping between team identifiers for macOS system extensions and system extension types. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-macossystemextensiontypemapping?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--macos_extensions--system_extensions_allowed_types))
- `system_extensions_block_override` (Boolean) Gets or sets whether to allow the user to approve additional system extensions not explicitly allowed by configuration profiles. The _provider_ default value is `false`.

<a id="nestedatt--macos_extensions--kernel_extensions_allowed"></a>
### Nested Schema for `macos_extensions.kernel_extensions_allowed`

Required:

- `bundle_id` (String) Bundle ID of the kernel extension.

Optional:

- `team_identifier` (String) The team identifier that was used to sign the kernel extension.


<a id="nestedatt--macos_extensions--system_extensions_allowed"></a>
### Nested Schema for `macos_extensions.system_extensions_allowed`

Required:

- `bundle_id` (String) Gets or sets the bundle identifier of the system extension.

Optional:

- `team_identifier` (String) Gets or sets the team identifier that was used to sign the system extension.


<a id="nestedatt--macos_extensions--system_extensions_allowed_types"></a>
### Nested Schema for `macos_extensions.system_extensions_allowed_types`

Required:

- `allowed_types` (String) Gets or sets the allowed macOS system extension types. / Flag enum representing the allowed macOS system extension types; possible values are: `driverExtensionsAllowed` (Enables driver extensions.), `networkExtensionsAllowed` (Enables network extensions.), `endpointSecurityExtensionsAllowed` (Enables endpoint security extensions.)
- `team_identifier` (String) Gets or sets the team identifier used to sign the system extension.



<a id="nestedatt--macos_software_update"></a>
### Nested Schema for `macos_software_update`

Optional:

- `all_other_update_behavior` (String) Update behavior for all other updates. / Update behavior options for macOS software updates; possible values are: `notConfigured` (Not configured.), `default` (Download and/or install the software update, depending on the current device state.), `downloadOnly` (Download the software update without installing it.), `installASAP` (Install an already downloaded software update.), `notifyOnly` (Download the software update and notify the user via the App Store.), `installLater` (Download the software update and install it at a later time.). The _provider_ default value is `"notConfigured"`.
- `config_data_update_behavior` (String) Update behavior for configuration data file updates. / Update behavior options for macOS software updates; possible values are: `notConfigured` (Not configured.), `default` (Download and/or install the software update, depending on the current device state.), `downloadOnly` (Download the software update without installing it.), `installASAP` (Install an already downloaded software update.), `notifyOnly` (Download the software update and notify the user via the App Store.), `installLater` (Download the software update and install it at a later time.). The _provider_ default value is `"notConfigured"`.
- `critical_update_behavior` (String) Update behavior for critical updates. / Update behavior options for macOS software updates; possible values are: `notConfigured` (Not configured.), `default` (Download and/or install the software update, depending on the current device state.), `downloadOnly` (Download the software update without installing it.), `installASAP` (Install an already downloaded software update.), `notifyOnly` (Download the software update and notify the user via the App Store.), `installLater` (Download the software update and install it at a later time.). The _provider_ default value is `"notConfigured"`.
- `custom_update_time_windows` (Attributes Set) Custom Time windows when updates will be allowed or blocked. This collection can contain a maximum of 20 elements. / Custom update time window / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-customupdatetimewindow?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--macos_software_update--custom_update_time_windows))
- `firmware_update_behavior` (String) Update behavior for firmware updates. / Update behavior options for macOS software updates; possible values are: `notConfigured` (Not configured.), `default` (Download and/or install the software update, depending on the current device state.), `downloadOnly` (Download the software update without installing it.), `installASAP` (Install an already downloaded software update.), `notifyOnly` (Download the software update and notify the user via the App Store.), `installLater` (Download the software update and install it at a later time.). The _provider_ default value is `"notConfigured"`.
- `max_user_deferrals_count` (Number) The maximum number of times the system allows the user to postpone an update before it’s installed. Supported values: 0 - 366. Valid values 0 to 365
- `priority` (String) The scheduling priority for downloading and preparing the requested update. Default: Low. Possible values: Null, Low, High. / The scheduling priority options for downloading and preparing the requested mac OS update; possible values are: `low` (Indicates low scheduling priority for downloading and preparing the requested update), `high` (Indicates high scheduling priority for downloading and preparing the requested update), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.)
- `update_schedule_type` (String) Update schedule type. / Update schedule type for macOS software updates; possible values are: `alwaysUpdate` (Always update.), `updateDuringTimeWindows` (Update during time windows.), `updateOutsideOfTimeWindows` (Update outside of time windows.). The _provider_ default value is `"alwaysUpdate"`.
- `update_time_window_utc_offset_in_minutes` (Number) Minutes indicating UTC offset for each update time window

<a id="nestedatt--macos_software_update--custom_update_time_windows"></a>
### Nested Schema for `macos_software_update.custom_update_time_windows`

Required:

- `end_day` (String) End day of the time window. / Possible values for a weekday; possible values are: `sunday` (Sunday.), `monday` (Monday.), `tuesday` (Tuesday.), `wednesday` (Wednesday.), `thursday` (Thursday.), `friday` (Friday.), `saturday` (Saturday.)
- `end_time` (String) End time of the time window
- `start_day` (String) Start day of the time window. / Possible values for a weekday; possible values are: `sunday` (Sunday.), `monday` (Monday.), `tuesday` (Tuesday.), `wednesday` (Wednesday.), `thursday` (Thursday.), `friday` (Friday.), `saturday` (Saturday.)
- `start_time` (String) Start time of the time window



<a id="nestedatt--windows_health_monitoring"></a>
### Nested Schema for `windows_health_monitoring`

Optional:

- `allow` (String) Enables device health monitoring on the device. / Possible values of a property; possible values are: `notConfigured` (Device default value, no intent.), `enabled` (Enables the setting on the device.), `disabled` (Disables the setting on the device.). The _provider_ default value is `"notConfigured"`.
- `custom_scope` (String) Specifies custom set of events collected from the device where health monitoring is enabled
- `scope` (String) Specifies set of events collected from the device where health monitoring is enabled. / Device health monitoring scope; possible values are: `undefined` (Undefined), `healthMonitoring` (Basic events for windows device health monitoring), `bootPerformance` (Boot performance events), `windowsUpdates` (Windows updates events), `privilegeManagement` (PrivilegeManagement). The _provider_ default value is `"undefined"`.


<a id="nestedatt--windows_update_for_business"></a>
### Nested Schema for `windows_update_for_business`

Optional:

- `allow_windows11_upgrade` (Boolean) When TRUE, allows eligible Windows 10 devices to upgrade to Windows 11. When FALSE, implies the device stays on the existing operating system. Returned by default. Query parameters are not supported. The _provider_ default value is `false`.
- `auto_restart_notification_dismissal` (String) Specify the method by which the auto-restart required notification is dismissed. Possible values are: NotConfigured, Automatic, User. Returned by default. Query parameters are not supported. / Auto restart required notification dismissal method; possible values are: `notConfigured` (Not configured), `automatic` (Auto dismissal Indicates that the notification is automatically dismissed without user intervention), `user` (User dismissal. Allows the user to dismiss the notification), `unknownFutureValue` (Evolvable enum member). The _provider_ default value is `"notConfigured"`.
- `automatic_update_mode` (String) The Automatic Update Mode. Possible values are: UserDefined, NotifyDownload, AutoInstallAtMaintenanceTime, AutoInstallAndRebootAtMaintenanceTime, AutoInstallAndRebootAtScheduledTime, AutoInstallAndRebootWithoutEndUserControl, WindowsDefault. UserDefined is the default value, no intent. Returned by default. Query parameters are not supported. / Possible values for automatic update mode; possible values are: `userDefined` (User Defined, default value, no intent.), `notifyDownload` (Notify on download.), `autoInstallAtMaintenanceTime` (Auto-install at maintenance time.), `autoInstallAndRebootAtMaintenanceTime` (Auto-install and reboot at maintenance time.), `autoInstallAndRebootAtScheduledTime` (Auto-install and reboot at scheduled time.), `autoInstallAndRebootWithoutEndUserControl` (Auto-install and restart without end-user control), `windowsDefault` (Reset to Windows default value.). The _provider_ default value is `"autoInstallAtMaintenanceTime"`.
- `business_ready_updates_only` (String) Determines which branch devices will receive their updates from. Possible values are: UserDefined, All, BusinessReadyOnly, WindowsInsiderBuildFast, WindowsInsiderBuildSlow, WindowsInsiderBuildRelease. Returned by default. Query parameters are not supported. / Which branch devices will receive their updates from; possible values are: `userDefined` (Allow the user to set.), `all` (Semi-annual Channel (Targeted). Device gets all applicable feature updates from Semi-annual Channel (Targeted).), `businessReadyOnly` (Semi-annual Channel. Device gets feature updates from Semi-annual Channel.), `windowsInsiderBuildFast` (Windows Insider build - Fast), `windowsInsiderBuildSlow` (Windows Insider build - Slow), `windowsInsiderBuildRelease` (Release Windows Insider build). The _provider_ default value is `"userDefined"`.
- `deadline_for_feature_updates_in_days` (Number) Number of days before feature updates are installed automatically with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
- `deadline_for_quality_updates_in_days` (Number) Number of days before quality updates are installed automatically with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
- `deadline_grace_period_in_days` (Number) Number of days after deadline until restarts occur automatically with valid range from 0 to 7 days. Returned by default. Query parameters are not supported.
- `delivery_optimization_mode` (String) The Delivery Optimization Mode. Possible values are: UserDefined, HttpOnly, HttpWithPeeringNat, HttpWithPeeringPrivateGroup, HttpWithInternetPeering, SimpleDownload, BypassMode. UserDefined allows the user to set. Returned by default. Query parameters are not supported. / Delivery optimization mode for peer distribution; possible values are: `userDefined` (Allow the user to set.), `httpOnly` (HTTP only, no peering), `httpWithPeeringNat` (OS default – Http blended with peering behind the same network address translator), `httpWithPeeringPrivateGroup` (HTTP blended with peering across a private group), `httpWithInternetPeering` (HTTP blended with Internet peering), `simpleDownload` (Simple download mode with no peering), `bypassMode` (Bypass mode. Do not use Delivery Optimization and use BITS instead). The _provider_ default value is `"userDefined"`.
- `drivers_excluded` (Boolean) When TRUE, excludes Windows update Drivers. When FALSE, does not exclude Windows update Drivers. Returned by default. Query parameters are not supported. The _provider_ default value is `false`.
- `engaged_restart_deadline_in_days` (Number) Deadline in days before automatically scheduling and executing a pending restart outside of active hours, with valid range from 2 to 30 days. Returned by default. Query parameters are not supported.
- `engaged_restart_snooze_schedule_in_days` (Number) Number of days a user can snooze Engaged Restart reminder notifications with valid range from 1 to 3 days. Returned by default. Query parameters are not supported.
- `engaged_restart_transition_schedule_in_days` (Number) Number of days before transitioning from Auto Restarts scheduled outside of active hours to Engaged Restart, which requires the user to schedule, with valid range from 0 to 30 days. Returned by default. Query parameters are not supported.
- `feature_updates_deferral_period_in_days` (Number) Defer Feature Updates by these many days with valid range from 0 to 30 days. Returned by default. Query parameters are not supported. The _provider_ default value is `0`.
- `feature_updates_paused` (Boolean) When TRUE, assigned devices are paused from receiving feature updates for up to 35 days from the time you pause the ring. When FALSE, does not pause Feature Updates. Returned by default. Query parameters are not supported.s. The _provider_ default value is `false`.
- `feature_updates_rollback_window_in_days` (Number) The number of days after a Feature Update for which a rollback is valid with valid range from 2 to 60 days. Returned by default. Query parameters are not supported. The _provider_ default value is `10`.
- `installation_schedule` (Attributes) The Installation Schedule. Possible values are: ActiveHoursStart, ActiveHoursEnd, ScheduledInstallDay, ScheduledInstallTime. Returned by default. Query parameters are not supported. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowsupdateinstallscheduletype?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows_update_for_business--installation_schedule))
- `microsoft_update_service_allowed` (Boolean) When TRUE, allows Microsoft Update Service. When FALSE, does not allow Microsoft Update Service. Returned by default. Query parameters are not supported. The _provider_ default value is `true`.
- `postpone_reboot_until_after_deadline` (Boolean) When TRUE the device should wait until deadline for rebooting outside of active hours. When FALSE the device should not wait until deadline for rebooting outside of active hours. Returned by default. Query parameters are not supported.
- `prerelease_features` (String) The Pre-Release Features. Possible values are: UserDefined, SettingsOnly, SettingsAndExperimentations, NotAllowed. UserDefined is the default value, no intent. Returned by default. Query parameters are not supported. / Possible values for pre-release features; possible values are: `userDefined` (User Defined, default value, no intent.), `settingsOnly` (Settings only pre-release features.), `settingsAndExperimentations` (Settings and experimentations pre-release features.), `notAllowed` (Pre-release features not allowed.). The _provider_ default value is `"userDefined"`.
- `quality_updates_deferral_period_in_days` (Number) Defer Quality Updates by these many days with valid range from 0 to 30 days. Returned by default. Query parameters are not supported. The _provider_ default value is `0`.
- `quality_updates_paused` (Boolean) When TRUE, assigned devices are paused from receiving quality updates for up to 35 days from the time you pause the ring. When FALSE, does not pause Quality Updates. Returned by default. Query parameters are not supported. The _provider_ default value is `false`.
- `schedule_imminent_restart_warning_in_minutes` (Number) Specify the period for auto-restart imminent warning notifications. Supported values: 15, 30 or 60 (minutes). Returned by default. Query parameters are not supported.
- `schedule_restart_warning_in_hours` (Number) Specify the period for auto-restart warning reminder notifications. Supported values: 2, 4, 8, 12 or 24 (hours). Returned by default. Query parameters are not supported.
- `skip_checks_before_restart` (Boolean) When TRUE, skips all checks before restart: Battery level = 40%, User presence, Display Needed, Presentation mode, Full screen mode, phone call state, game mode etc. When FALSE, does not skip all checks before restart. Returned by default. Query parameters are not supported. The _provider_ default value is `false`.
- `update_notification_level` (String) Specifies what Windows Update notifications users see. Possible values are: NotConfigured, DefaultNotifications, RestartWarningsOnly, DisableAllNotifications. Returned by default. Query parameters are not supported. / Windows Update Notification Display Options; possible values are: `notConfigured` (Not configured), `defaultNotifications` (Use the default Windows Update notifications.), `restartWarningsOnly` (Turn off all notifications, excluding restart warnings.), `disableAllNotifications` (Turn off all notifications, including restart warnings.), `unknownFutureValue` (Evolvable enum member). The _provider_ default value is `"defaultNotifications"`.
- `update_weeks` (String) Schedule the update installation on the weeks of the month. Possible values are: UserDefined, FirstWeek, SecondWeek, ThirdWeek, FourthWeek, EveryWeek. Returned by default. Query parameters are not supported. / Scheduled the update installation on the weeks of the month; possible values are: `userDefined` (Allow the user to set.), `firstWeek` (Scheduled the update installation on the first week of the month), `secondWeek` (Scheduled the update installation on the second week of the month), `thirdWeek` (Scheduled the update installation on the third week of the month), `fourthWeek` (Scheduled the update installation on the fourth week of the month), `everyWeek` (Scheduled the update installation on every week of the month), `unknownFutureValue` (Evolvable enum member)
- `user_pause_access` (String) Specifies whether to enable end user’s access to pause software updates. Possible values are: NotConfigured, Enabled, Disabled. Returned by default. Query parameters are not supported. / Possible values of a property; possible values are: `notConfigured` (Device default value, no intent.), `enabled` (Enables the setting on the device.), `disabled` (Disables the setting on the device.). The _provider_ default value is `"enabled"`.
- `user_windows_update_scan_access` (String) Specifies whether to disable user’s access to scan Windows Update. Possible values are: NotConfigured, Enabled, Disabled. Returned by default. Query parameters are not supported. / Possible values of a property; possible values are: `notConfigured` (Device default value, no intent.), `enabled` (Enables the setting on the device.), `disabled` (Disables the setting on the device.). The _provider_ default value is `"enabled"`.

<a id="nestedatt--windows_update_for_business--installation_schedule"></a>
### Nested Schema for `windows_update_for_business.installation_schedule`

Optional:

- `active_hours` (Attributes) https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowsupdateactivehoursinstall?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows_update_for_business--installation_schedule--active_hours))
- `scheduled` (Attributes) https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windowsupdatescheduledinstall?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows_update_for_business--installation_schedule--scheduled))

<a id="nestedatt--windows_update_for_business--installation_schedule--active_hours"></a>
### Nested Schema for `windows_update_for_business.installation_schedule.active_hours`

Optional:

- `end` (String) Active Hours End. The _provider_ default value is `"17:00:00.0000000"`.
- `start` (String) Active Hours Start. The _provider_ default value is `"08:00:00.0000000"`.


<a id="nestedatt--windows_update_for_business--installation_schedule--scheduled"></a>
### Nested Schema for `windows_update_for_business.installation_schedule.scheduled`

Optional:

- `day` (String) Scheduled Install Day in week. / Possible values for a weekly schedule; possible values are: `userDefined` (User Defined, default value, no intent.), `everyday` (Everyday.), `sunday` (Sunday.), `monday` (Monday.), `tuesday` (Tuesday.), `wednesday` (Wednesday.), `thursday` (Thursday.), `friday` (Friday.), `saturday` (Saturday.), `noScheduledScan` (No Scheduled Scan). The _provider_ default value is `"everyday"`.
- `time` (String) Scheduled Install Time during day. The _provider_ default value is `"03:00:00.0000000"`.
