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

resource "microsoft365wp_device_configuration" "windows10_general" {
  display_name = "TF Test Windows 10 General"
  windows10_general = {
    screen_capture_blocked = true
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
