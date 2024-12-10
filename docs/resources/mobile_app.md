---
page_title: "microsoft365wp_mobile_app Resource - microsoft365wp"
subcategory: "MS Graph: App management"
---

# microsoft365wp_mobile_app (Resource)

An abstract class containing the base properties for Intune mobile apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mobileapp?view=graph-rest-beta

Provider Note: There have been reports about OpenTofu crashing (with a panic) when trying to generate configuration using `-generate-config-out=...`. Unfortunately this happens inside the code of OpenTofu itself and is not easily diagnoseable from outside or changeable from the provider. A workaround might be to instead use Terraform for generating the config (as it seems to work well for it) - the result should be fully compatible and usable also with OpenTofu.

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


resource "microsoft365wp_mobile_app" "android_store" {
  display_name = "TF Test android_store"

  android_store = {
    app_store_url                      = "https://play.google.com/store/apps/details?id=com.microsoft.windowsintune.companyportal"
    minimum_supported_operating_system = { v7_0 = true }
  }

  assignments = [{ target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } }]

  categories = [
    { id = "2b73ae71-12c8-49be-b462-3dae769ccd9d" }, # Business
    { id = "ed899483-3019-425e-a470-28e901b9790e" }, # Productivity
  ]
}

resource "microsoft365wp_mobile_app" "managed_android_store" {
  display_name = "TF Test managed_android_store"

  managed_android_store = {
    app_store_url                      = "https://play.google.com/store/apps/details?id=com.microsoft.office.officelens"
    package_id                         = "com.microsoft.office.officelens"
    minimum_supported_operating_system = { v7_0 = true }
  }

  assignments = [{ target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } }]
}

resource "microsoft365wp_mobile_app" "android_managed_store" {
  display_name = "TF Test android_managed_store"

  publisher = "Microsoft"
  android_managed_store = {
    app_identifier = "com.microsoft.word"
  }

  assignments = [{
    target   = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } }
    settings = { android_managed_store = {} } # see docs
  }]
}

resource "microsoft365wp_mobile_app" "ios_store" {
  display_name = "TF Test ios_store"

  ios_store = {
    app_store_url                      = "https://apps.apple.com/us/app/intune-company-portal/id719171358?uo=4"
    minimum_supported_operating_system = { v14_0 = true }
  }

  assignments = [{
    target   = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } }
    settings = { ios_store = {} } # see docs
  }]
}

resource "microsoft365wp_mobile_app" "managed_ios_store" {
  display_name = "TF Test managed_ios_store"

  managed_ios_store = {
    app_store_url                      = "https://itunes.apple.com/us/app/microsoft-lens-pdf-scanner/id975925059?mt=8"
    bundle_id                          = "com.microsoft.officelens"
    minimum_supported_operating_system = { v14_0 = true }
  }

  assignments = [{
    target   = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } }
    settings = { ios_store = {} } # see docs
  }]
}

resource "microsoft365wp_mobile_app" "ios_webclip" {
  display_name = "TF Test ios_webclip"

  ios_webclip = {
    app_url = "https://apple.com"
  }

  assignments = [{ target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } }]
}

resource "microsoft365wp_mobile_app" "macos_ms_defender" {
  display_name = "TF Test macos_ms_defender"
  large_icon = {
    type         = "image/png"
    value_base64 = filebase64("MsDefenderLogo.png")
  }

  macos_ms_defender = {} # no properties

  assignments = [{ target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } }]
}

resource "microsoft365wp_mobile_app" "macos_ms_edge" {
  display_name = "TF Test macos_ms_edge"
  large_icon = {
    type         = "image/png"
    value_base64 = filebase64("MsEdgeLogo.png")
  }

  macos_ms_edge = {} # no required properties

  assignments = [{ target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } }]
}

resource "microsoft365wp_mobile_app" "macos_office_suite" {
  display_name = "TF Test macos_office_suite"
  large_icon = {
    type         = "image/png"
    value_base64 = filebase64("MsOfficeLogo.png")
  }

  macos_office_suite = {} # no properties

  assignments = [{ target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } }]
}

resource "microsoft365wp_mobile_app" "macos_webclip" {
  display_name = "TF Test macos_webclip"

  macos_webclip = {
    app_url = "https://apple.com"
  }

  assignments = [{ target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } }]
}

resource "microsoft365wp_mobile_app" "web_link" {
  display_name = "TF Test web_link"

  web_link = {
    app_url = "https://google.com"
  }

  assignments = [{ target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } }]
}

resource "microsoft365wp_mobile_app" "winget" {
  display_name = "TF Test winget"

  winget = {
    package_identifier = "9WZDNCRFJ3PZ" # Company Portal
  }

  assignments = [{
    target   = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } }
    settings = { winget = {} } # see docs
  }]
}

resource "microsoft365wp_mobile_app" "windows_ms_edge" {
  display_name = "TF Test windows_ms_edge"
  large_icon = {
    type         = "image/png"
    value_base64 = filebase64("MsEdgeLogo.png")
  }

  windows_ms_edge = {} # no required properties

  assignments = [{ target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } }]
}

resource "microsoft365wp_mobile_app" "windows_office_suite" {
  display_name = "TF Test windows_office_suite"
  large_icon = {
    type         = "image/png"
    value_base64 = filebase64("MsOfficeLogo.png")
  }

  windows_office_suite = {} # no required properties

  assignments = [{ target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } }]
}

resource "microsoft365wp_mobile_app" "windows_store" {
  display_name = "TF Test windows_store"

  windows_store = {
    app_store_url = "https://www.microsoft.com/en-us/p/company-portal/9wzdncrfj3pz"
  }

  assignments = [{
    intent = "available" # required (default) is not allowed here
    target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } }
  }]
}

resource "microsoft365wp_mobile_app" "windows_web_link" {
  display_name = "TF Test windows_web_link"

  windows_web_link = {
    app_url = "https://microsoft.com"
  }

  assignments = [{ target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } }]
}


locals {
  win32_lob_source_file = "${path.module}/test.intunewin"
  win32_lob_metadata    = provider::microsoft365wp::parse_intunewin_metadata(local.win32_lob_source_file)
}

resource "microsoft365wp_mobile_app" "win32_lob" {
  display_name = "TF Test win32_lob"

  win32_lob = {
    file_name       = local.win32_lob_metadata.file_name
    setup_file_path = local.win32_lob_metadata.setup_file

    install_command_line   = "cmd /c echo Hello"
    uninstall_command_line = "exit 1"

    detection_rules = [
      {
        filesystem = {
          detection_type      = "exists"
          path                = "C:\\"
          file_or_folder_name = "test.txt"
        }
      }
    ]

    source_file   = local.win32_lob_source_file
    source_sha256 = filesha256(local.win32_lob_source_file)
  }

  assignments = [
    for x in [
      "298fded6-b252-4166-a473-f405e935f58d",
      "62e39046-aad3-4423-98e0-b486e3538aff",
    ] :
    {
      target   = { group = { group_id = x } }
      settings = { win32_lob = {} } # see docs
    }
  ]
}


locals {
  macos_dmg_source_file = "${path.module}/test.dmg"
  macos_dmg_bundle_id   = "com.realvnc.vncviewer"
  macos_dmg_version     = "7.12.1"
}

resource "microsoft365wp_mobile_app" "macos_dmg" {
  display_name = "TF Test macos_dmg"

  macos_dmg = {
    file_name = basename(local.macos_dmg_source_file)

    included_apps = [{
      bundle_id      = local.macos_dmg_bundle_id
      bundle_version = local.macos_dmg_version
    }]
    primary_bundle_id      = local.macos_dmg_bundle_id
    primary_bundle_version = local.macos_dmg_version

    source_file   = local.macos_dmg_source_file
    source_sha256 = filesha256(local.macos_dmg_source_file)

    minimum_supported_operating_system = { v11_0 = true }
  }

  assignments = [
    for x in [
      "298fded6-b252-4166-a473-f405e935f58d",
      "62e39046-aad3-4423-98e0-b486e3538aff",
    ] :
    {
      target = { group = { group_id = x } }
    }
  ]
}


locals {
  macos_pkg_source_file = "${path.module}/test.pkg"
  macos_pkg_included_apps = [
    { bundle_id = "com.microsoft.autoupdate2", bundle_version = "4.73" },
    { bundle_id = "org.cocoapods.CocoaLumberjack", bundle_version = "3.8.2" },
    { bundle_id = "com.microsoft.remotehelp", bundle_version = "1.0" },
  ]
  macos_pkg_script = <<-EOT
    #!/bin/bash'
    echo "Hello World"
    EOT
}

resource "microsoft365wp_mobile_app" "macos_pkg" {
  display_name = "TF Test macos_pkg"

  macos_pkg = {
    file_name = basename(local.macos_pkg_source_file)

    included_apps          = local.macos_pkg_included_apps
    primary_bundle_id      = local.macos_pkg_included_apps[0].bundle_id
    primary_bundle_version = local.macos_pkg_included_apps[0].bundle_version

    source_file   = local.macos_pkg_source_file
    source_sha256 = filesha256(local.macos_pkg_source_file)

    minimum_supported_operating_system = { v11_0 = true }

    # base64 encoding is done by provider
    pre_install_script  = { script_content = local.macos_pkg_script }
    post_install_script = { script_content = local.macos_pkg_script }
  }

  assignments = [
    for x in [
      "298fded6-b252-4166-a473-f405e935f58d",
      "62e39046-aad3-4423-98e0-b486e3538aff",
    ] :
    {
      target = { group = { group_id = x } }
    }
  ]
}


locals {
  windows_msi_source_file = "${path.module}/test.msi"
  windows_msi_version     = "4.19.9"

  # manifest attributes are based on MSI property values (use MSI Orca on Windows or similar to read them)
  # conditions in comments below inferred from Intune Portal JavaScript source code
  windows_msi_manifest = {
    MsiExecutionContext           = "Any"                                    # ALLUSERS: "": "System", "1": "User", "2": "Any"
    MsiRequiresReboot             = false                                    # REBOOT == FORCE
    MsiUpgradeCode                = "{23170F69-40C1-2702-0000-000004000000}" # UpgradeCode
    MsiIsMachineInstall           = true                                     # MSIINSTALLPERUSER == "" && (ALLUSERS == "1" || ALLUSERS == "2")
    MsiIsUserInstall              = false                                    # !manifest.MsiIsMachineInstall
    MsiRequiresLogon              = false                                    # ALLUSERS == ""
    MsiIncludesServices           = false                                    # always false
    MsiContainsSystemRegistryKeys = false                                    # always false
    MsiContainsSystemFolders      = false                                    # always false
  }
  # poor man's xmlencode
  windows_msi_manifest_xml = <<-EOT
    <MobileMsiData
        MsiExecutionContext="${local.windows_msi_manifest.MsiExecutionContext}"
        MsiRequiresReboot="${local.windows_msi_manifest.MsiRequiresReboot}"
        MsiUpgradeCode="${local.windows_msi_manifest.MsiUpgradeCode}"
        MsiIsMachineInstall="${local.windows_msi_manifest.MsiIsMachineInstall}"
        MsiIsUserInstall="${local.windows_msi_manifest.MsiIsUserInstall}"
        MsiIncludesServices="${local.windows_msi_manifest.MsiIncludesServices}"
        MsiContainsSystemRegistryKeys="${local.windows_msi_manifest.MsiContainsSystemRegistryKeys}"
        MsiContainsSystemFolders="${local.windows_msi_manifest.MsiContainsSystemFolders}" >
    </MobileMsiData>
    EOT
  # also see docs for attribute manifest_base64
  windows_msi_manifest_base64 = base64encode(local.windows_msi_manifest_xml)
}

resource "microsoft365wp_mobile_app" "windows_msi" {
  display_name = "TF Test windows_msi"

  windows_msi = {
    file_name        = basename(local.windows_msi_source_file)
    identity_version = local.windows_msi_version
    product_version  = local.windows_msi_version

    source_file     = local.windows_msi_source_file
    source_sha256   = filesha256(local.windows_msi_source_file)
    manifest_base64 = local.windows_msi_manifest_base64
  }

  assignments = [{ target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } }]
}


locals {
  windows_universal_appx_source_file = "${path.module}/Microsoft.NET.Native.Runtime.2.2_2.2.28604.0_x64__8wekyb3d8bbwe.appx"
  windows_universal_appx_metadata    = provider::microsoft365wp::parse_appxmsix_metadata(local.windows_universal_appx_source_file)
}

resource "microsoft365wp_mobile_app" "windows_universal_appx" {
  display_name = "TF Test windows_universal_appx"
  description  = local.windows_universal_appx_metadata.properties.description
  publisher    = local.windows_universal_appx_metadata.properties.publisher_display_name
  large_icon   = local.windows_universal_appx_metadata.logo

  windows_universal_appx = {
    file_name               = basename(local.windows_universal_appx_source_file)
    identity_name           = local.windows_universal_appx_metadata.identity.name
    identity_publisher_hash = local.windows_universal_appx_metadata.identity.publisher_hash
    identity_version        = local.windows_universal_appx_metadata.identity.version

    applicable_device_types            = "desktop"
    applicable_architectures           = local.windows_universal_appx_metadata.identity.processor_architecture
    minimum_supported_operating_system = { v10_0 = true } # see docs

    is_msix       = endswith(lower(local.windows_universal_appx_source_file), ".msix")
    source_file   = local.windows_universal_appx_source_file
    source_sha256 = filesha256(local.windows_universal_appx_source_file)
  }

  assignments = [{
    target   = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } }
    settings = { windows_universal_appx = {} } # see docs
  }]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The admin provided or imported title of the app.

### Optional

- `android_managed_store` (Attributes) Contains properties and inherited properties for Android Managed Store Apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidmanagedstoreapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--android_managed_store))
- `android_store` (Attributes) Contains properties and inherited properties for Android store apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidstoreapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--android_store))
- `assignments` (Attributes Set) The list of group assignments for this mobile app. / A class containing the properties used for Group Assignment of a Mobile App. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileappassignment?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--assignments))
- `categories` (Attributes Set) The list of categories for this app. / Contains properties for a single Intune app category. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileappcategory?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--categories))
- `description` (String) The description of the app. The _provider_ default value is `""`.
- `developer` (String) The developer of the app.
- `information_url` (String) The more information Url.
- `ios_store` (Attributes) Contains properties and inherited properties for iOS store apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosstoreapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_store))
- `ios_webclip` (Attributes) Contains properties and inherited properties for iOS web apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosipadoswebclip?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_webclip))
- `is_featured` (Boolean) The value indicating whether the app is marked as featured by the admin. The _provider_ default value is `false`.
- `large_icon` (Attributes) The large icon, to be displayed in the app details and used for upload of the icon. / Contains properties for a generic mime content. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mimecontent?view=graph-rest-beta (see [below for nested schema](#nestedatt--large_icon))
- `macos_dmg` (Attributes) Contains properties and inherited properties for the MacOS DMG (Apple Disk Image) App. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosdmgapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_dmg))
- `macos_ms_defender` (Attributes) Contains properties and inherited properties for the macOS Microsoft Defender App. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosmicrosoftdefenderapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_ms_defender))
- `macos_ms_edge` (Attributes) Contains properties and inherited properties for the macOS Microsoft Edge App. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosmicrosoftedgeapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_ms_edge))
- `macos_office_suite` (Attributes) Contains properties and inherited properties for the MacOS Office Suite App. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosofficesuiteapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_office_suite))
- `macos_pkg` (Attributes) Contains properties and inherited properties for the MacOSPkgApp. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macospkgapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_pkg))
- `macos_webclip` (Attributes) Contains properties and inherited properties for macOS web apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macoswebclip?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_webclip))
- `managed_android_store` (Attributes) Contains properties and inherited properties for Android store apps that you can manage with an Intune app protection policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-managedandroidstoreapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--managed_android_store))
- `managed_ios_store` (Attributes) Contains properties and inherited properties for an iOS store app that you can manage with an Intune app protection policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-managediosstoreapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--managed_ios_store))
- `notes` (String) Notes for the app.
- `owner` (String) The owner of the app.
- `privacy_information_url` (String) The privacy statement Url.
- `publisher` (String) The publisher of the app.
- `role_scope_tag_ids` (Set of String) List of scope tag ids for this mobile app. The _provider_ default value is `["0"]`.
- `web_link` (Attributes) Contains properties and inherited properties for web apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-webapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--web_link))
- `win32_lob` (Attributes) Contains properties and inherited properties for Win32 apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--win32_lob))
- `windows_ms_edge` (Attributes) Contains properties and inherited properties for the Microsoft Edge app on Windows. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-windowsmicrosoftedgeapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows_ms_edge))
- `windows_msi` (Attributes) Contains properties and inherited properties for Windows Mobile MSI Line Of Business apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-windowsmobilemsi?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows_msi))
- `windows_office_suite` (Attributes) Contains properties and inherited properties for the Office365 Suite App. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-officesuiteapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows_office_suite))
- `windows_store` (Attributes) Contains properties and inherited properties for Windows Store apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-windowsstoreapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows_store))
- `windows_universal_appx` (Attributes) Contains properties and inherited properties for Windows Universal AppX Line Of Business apps. Inherits from `mobileLobApp`. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-windowsuniversalappx?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows_universal_appx))
- `windows_web_link` (Attributes) Contains properties and inherited properties for Windows web apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-windowswebapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows_web_link))
- `winget` (Attributes) A MobileApp that is based on a referenced application in a WinGet repository. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-wingetapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--winget))

### Read-Only

- `android_for_work` (Attributes) Contains properties and inherited properties for Android for Work (AFW) Apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidforworkapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--android_for_work))
- `android_lob` (Attributes) Contains properties and inherited properties for Android Line Of Business apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidlobapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--android_lob))
- `created_date_time` (String) The date and time the app was created.
- `dependent_app_count` (Number) The total number of dependencies the child app has.
- `id` (String) Key of the entity.
- `ios_lob` (Attributes) Contains properties and inherited properties for iOS Line Of Business apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-ioslobapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_lob))
- `ios_vpp` (Attributes) Contains properties and inherited properties for iOS Volume-Purchased Program (VPP) Apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosvppapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_vpp))
- `is_assigned` (Boolean) The value indicating whether the app is assigned to at least one group.
- `last_modified_date_time` (String) The date and time the app was last modified.
- `macos_lob` (Attributes) Contains properties and inherited properties for the macOS LOB App. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macoslobapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_lob))
- `macos_vpp` (Attributes) Contains properties and inherited properties for MacOS Volume-Purchased Program (VPP) Apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosvppapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_vpp))
- `managed_android_lob` (Attributes) Contains properties and inherited properties for Managed Android Line Of Business apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-managedandroidlobapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--managed_android_lob))
- `managed_ios_lob` (Attributes) Contains properties and inherited properties for Managed iOS Line Of Business apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-managedioslobapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--managed_ios_lob))
- `publishing_state` (String) The publishing state for the app. The app cannot be assigned unless the app is published. / Indicates the publishing state of an app; possible values are: `notPublished` (The app is not yet published.), `processing` (The app is pending service-side processing.), `published` (The app is published.)
- `superseded_app_count` (Number)
- `superseding_app_count` (Number)
- `upload_state` (Number) The upload state.

<a id="nestedatt--android_managed_store"></a>
### Nested Schema for `android_managed_store`

Required:

- `app_identifier` (String) The Identity Name.

Optional:

- `is_system_app` (Boolean) Indicates whether the app is a preinstalled system app. The _provider_ default value is `true`.

Read-Only:

- `app_store_url` (String) The Play for Work Store app URL. This property is read-only.
- `app_tracks` (Attributes Set) The tracks that are visible to this enterprise. This property is read-only. / Contains track information for Android Managed Store apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidmanagedstoreapptrack?view=graph-rest-beta (see [below for nested schema](#nestedatt--android_managed_store--app_tracks))
- `is_private` (Boolean) Indicates whether the app is only available to a given enterprise's users. This property is read-only.
- `package_id` (String) The package identifier. This property is read-only.
- `supports_oem_config` (Boolean) Whether this app supports OEMConfig policy. This property is read-only.
- `total_license_count` (Number) The total number of VPP licenses. This property is read-only.
- `used_license_count` (Number) The number of VPP licenses in use. This property is read-only.

<a id="nestedatt--android_managed_store--app_tracks"></a>
### Nested Schema for `android_managed_store.app_tracks`

Read-Only:

- `track_alias` (String) Friendly name for track. This property is read-only.
- `track_id` (String) Unique track identifier. This property is read-only.



<a id="nestedatt--android_store"></a>
### Nested Schema for `android_store`

Required:

- `app_store_url` (String) The Android app store URL.
- `minimum_supported_operating_system` (Attributes) The value for the minimum applicable operating system. / Contains properties for the minimum operating system required for an Android mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidminimumoperatingsystem?view=graph-rest-beta (see [below for nested schema](#nestedatt--android_store--minimum_supported_operating_system))

Read-Only:

- `package_id` (String) The package identifier. This property is read-only.

<a id="nestedatt--android_store--minimum_supported_operating_system"></a>
### Nested Schema for `android_store.minimum_supported_operating_system`

Optional:

- `v10_0` (Boolean) When TRUE, only Version 10.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v11_0` (Boolean) When TRUE, only Version 11.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v12_0` (Boolean) When TRUE, only Version 12.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v13_0` (Boolean) When TRUE, only Version 13.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v14_0` (Boolean) When TRUE, only Version 14.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v15_0` (Boolean) When TRUE, only Version 15.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_0` (Boolean) When TRUE, only Version 4.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_0_3` (Boolean) When TRUE, only Version 4.0.3 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_1` (Boolean) When TRUE, only Version 4.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_2` (Boolean) When TRUE, only Version 4.2 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_3` (Boolean) When TRUE, only Version 4.3 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_4` (Boolean) When TRUE, only Version 4.4 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v5_0` (Boolean) When TRUE, only Version 5.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v5_1` (Boolean) When TRUE, only Version 5.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v6_0` (Boolean) When TRUE, only Version 6.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v7_0` (Boolean) When TRUE, only Version 7.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v7_1` (Boolean) When TRUE, only Version 7.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v8_0` (Boolean) When TRUE, only Version 8.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v8_1` (Boolean) When TRUE, only Version 8.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v9_0` (Boolean) When TRUE, only Version 9.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.



<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `target` (Attributes) Base type for assignment targets. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-deviceandappmanagementassignmenttarget?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--target))

Optional:

- `intent` (String) The install intent defined by the admin. / Possible values for the install intent chosen by the admin; possible values are: `available` (Available install intent.), `required` (Required install intent.), `uninstall` (Uninstall install intent.), `availableWithoutEnrollment` (Available without enrollment install intent.). The _provider_ default value is `"required"`.
- `settings` (Attributes) The settings for target assignment defined by the admin. / Abstract class to contain properties used to assign a mobile app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mobileappassignmentsettings?view=graph-rest-beta  
Provider Note: For some mobileApp types it might be required to include appropriate (potentially empty) `settings` with each assignment to avoid recurring differences between plan and actual state, namely for the types: `android_managed_store`, `ios_store`, `managed_ios_store` (use `ios_store` as settings type), `winget`, `win32_lob`, `windows_universal_appx`. For details, see the example above. (see [below for nested schema](#nestedatt--assignments--settings))

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



<a id="nestedatt--assignments--settings"></a>
### Nested Schema for `assignments.settings`

Optional:

- `android_managed_store` (Attributes) Contains properties used to assign an Android Managed Store mobile app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-androidmanagedstoreappassignmentsettings?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--settings--android_managed_store))
- `ios_lob` (Attributes) Contains properties used to assign an iOS LOB mobile app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-ioslobappassignmentsettings?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--settings--ios_lob))
- `ios_store` (Attributes) Contains properties used to assign an iOS Store mobile app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-iosstoreappassignmentsettings?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--settings--ios_store))
- `ios_vpp` (Attributes) Contains properties used to assign an iOS VPP mobile app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-iosvppappassignmentsettings?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--settings--ios_vpp))
- `macos_lob` (Attributes) Contains properties used to assign a macOS LOB app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-macoslobappassignmentsettings?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--settings--macos_lob))
- `macos_vpp` (Attributes) Contains properties used to assign an Mac VPP mobile app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-macosvppappassignmentsettings?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--settings--macos_vpp))
- `win32_lob` (Attributes) Contains properties used to assign an Win32 LOB mobile app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-win32lobappassignmentsettings?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--settings--win32_lob))
- `windows_universal_appx` (Attributes) Contains properties used when assigning a Windows Universal AppX mobile app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-windowsuniversalappxappassignmentsettings?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--settings--windows_universal_appx))
- `winget` (Attributes) Contains properties used to assign a WinGet app to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-wingetappassignmentsettings?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--settings--winget))

<a id="nestedatt--assignments--settings--android_managed_store"></a>
### Nested Schema for `assignments.settings.android_managed_store`

Optional:

- `android_managed_store_app_track_ids` (Set of String) The track IDs to enable for this app assignment. The _provider_ default value is `[]`.
- `auto_update_mode` (String) The prioritization of automatic updates for this app assignment. / Prioritization for automatic updates of Android Managed Store apps set on assignment; possible values are: `default` (Default update behavior (device must be connected to Wifi, charging and not actively used).), `postponed` (Updates are postponed for a maximum of 90 days after the app becomes out of date.), `priority` (The app is updated as soon as possible by the developer. If device is online, it will be updated within minutes.), `unknownFutureValue` (Unknown future mode (reserved, not used right now).). The _provider_ default value is `"default"`.


<a id="nestedatt--assignments--settings--ios_lob"></a>
### Nested Schema for `assignments.settings.ios_lob`

Optional:

- `is_removable` (Boolean) When TRUE, indicates that the app can be uninstalled by the user. When FALSE, indicates that the app cannot be uninstalled by the user. By default, this property is set to null which internally is treated as TRUE.
- `prevent_managed_app_backup` (Boolean) When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to null which internally is treated as FALSE.
- `uninstall_on_device_removal` (Boolean) When TRUE, indicates that the app should be uninstalled when the device is removed from Intune. When FALSE, indicates that the app will not be uninstalled when the device is removed from Intune. By default, property is set to null which internally is treated as TRUE.
- `vpn_configuration_id` (String) This is the unique identifier (Id) of the VPN Configuration to apply to the app.


<a id="nestedatt--assignments--settings--ios_store"></a>
### Nested Schema for `assignments.settings.ios_store`

Optional:

- `is_removable` (Boolean) When TRUE, indicates that the app can be uninstalled by the user. When FALSE, indicates that the app cannot be uninstalled by the user. By default, this property is set to null which internally is treated as TRUE.
- `prevent_managed_app_backup` (Boolean) When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to null which internally is treated as FALSE.
- `uninstall_on_device_removal` (Boolean) When TRUE, indicates that the app should be uninstalled when the device is removed from Intune. When FALSE, indicates that the app will not be uninstalled when the device is removed from Intune. By default, property is set to null which internally is treated as TRUE.
- `vpn_configuration_id` (String) This is the unique identifier (Id) of the VPN Configuration to apply to the app.


<a id="nestedatt--assignments--settings--ios_vpp"></a>
### Nested Schema for `assignments.settings.ios_vpp`

Optional:

- `is_removable` (Boolean) Whether or not the app can be removed by the user.
- `prevent_auto_app_update` (Boolean) When TRUE, indicates that the app should not be automatically updated with the latest version from Apple app store. When FALSE, indicates that the app may be auto updated. By default, this property is set to null which internally is treated as FALSE.
- `prevent_managed_app_backup` (Boolean) When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to null which internally is treated as FALSE.
- `uninstall_on_device_removal` (Boolean) Whether or not to uninstall the app when device is removed from Intune.
- `use_device_licensing` (Boolean) Whether or not to use device licensing. The _provider_ default value is `false`.
- `vpn_configuration_id` (String) The VPN Configuration Id to apply for this app.


<a id="nestedatt--assignments--settings--macos_lob"></a>
### Nested Schema for `assignments.settings.macos_lob`

Optional:

- `uninstall_on_device_removal` (Boolean) When TRUE, indicates that the app should be uninstalled when the device is removed from Intune. When FALSE, indicates that the app will not be uninstalled when the device is removed from Intune.


<a id="nestedatt--assignments--settings--macos_vpp"></a>
### Nested Schema for `assignments.settings.macos_vpp`

Optional:

- `prevent_auto_app_update` (Boolean) When TRUE, indicates that the app should not be automatically updated with the latest version from Apple app store. When FALSE, indicates that the app may be auto updated. By default, this property is set to null which internally is treated as FALSE.
- `prevent_managed_app_backup` (Boolean) When TRUE, indicates that the app should not be backed up to iCloud. When FALSE, indicates that the app may be backed up to iCloud. By default, this property is set to null which internally is treated as FALSE.
- `uninstall_on_device_removal` (Boolean) Whether or not to uninstall the app when device is removed from Intune.
- `use_device_licensing` (Boolean) Whether or not to use device licensing. The _provider_ default value is `false`.


<a id="nestedatt--assignments--settings--win32_lob"></a>
### Nested Schema for `assignments.settings.win32_lob`

Optional:

- `auto_update_settings` (Attributes) The auto-update settings to apply for this app assignment. / Contains properties used to perform the auto-update of an application. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-win32lobappautoupdatesettings?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--settings--win32_lob--auto_update_settings))
- `delivery_optimization_priority` (String) The delivery optimization priority for this app assignment. This setting is not supported in National Cloud environments. / Contains value for delivery optimization priority; possible values are: `notConfigured` (Not configured or background normal delivery optimization priority.), `foreground` (Foreground delivery optimization priority.). The _provider_ default value is `"notConfigured"`.
- `install_time_settings` (Attributes) The install time settings to apply for this app assignment. / Contains properties used to determine when to offer an app to devices and when to install the app on devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-mobileappinstalltimesettings?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--settings--win32_lob--install_time_settings))
- `notifications` (String) The notification status for this app assignment. / Contains value for notification status; possible values are: `showAll` (Show all notifications.), `showReboot` (Only show restart notification and suppress other notifications.), `hideAll` (Hide all notifications.). The _provider_ default value is `"showAll"`.
- `restart_settings` (Attributes) The reboot settings to apply for this app assignment. / Contains properties describing restart coordination following an app installation. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-win32lobapprestartsettings?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--settings--win32_lob--restart_settings))

<a id="nestedatt--assignments--settings--win32_lob--auto_update_settings"></a>
### Nested Schema for `assignments.settings.win32_lob.auto_update_settings`

Optional:

- `auto_update_superseded_apps_state` (String) Contains value for auto-update superseded apps; possible values are: `notConfigured` (Indicates that the auto-update superseded apps state is not configured and the app will not auto-update the superseded apps.), `enabled` (Indicates that the auto-update superseded apps state is enabled and the app will auto-update the superseded apps if the superseded apps are installed on the device.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.). The _provider_ default value is `"notConfigured"`.


<a id="nestedatt--assignments--settings--win32_lob--install_time_settings"></a>
### Nested Schema for `assignments.settings.win32_lob.install_time_settings`

Optional:

- `deadline_date_time` (String) The time at which the app should be installed.
- `start_date_time` (String) The time at which the app should be available for installation.
- `use_local_time` (Boolean) Whether the local device time or UTC time should be used when determining the available and deadline times. The _provider_ default value is `false`.


<a id="nestedatt--assignments--settings--win32_lob--restart_settings"></a>
### Nested Schema for `assignments.settings.win32_lob.restart_settings`

Optional:

- `countdown_display_before_restart_in_minutes` (Number) The number of minutes before the restart time to display the countdown dialog for pending restarts. The _provider_ default value is `0`.
- `grace_period_in_minutes` (Number) The number of minutes to wait before restarting the device after an app installation. The _provider_ default value is `0`.
- `restart_notification_snooze_duration_in_minutes` (Number) The number of minutes to snooze the restart notification dialog when the snooze button is selected.



<a id="nestedatt--assignments--settings--windows_universal_appx"></a>
### Nested Schema for `assignments.settings.windows_universal_appx`

Optional:

- `use_device_context` (Boolean) If true, uses device execution context for Windows Universal AppX mobile app. Device-context install is not allowed when this type of app is targeted with Available intent. Defaults to false. The _provider_ default value is `false`.


<a id="nestedatt--assignments--settings--winget"></a>
### Nested Schema for `assignments.settings.winget`

Optional:

- `install_time_settings` (Attributes) The install time settings to apply for this app assignment. / Contains properties used to determine when to offer an app to devices and when to install the app on devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-wingetappinstalltimesettings?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--settings--winget--install_time_settings))
- `notifications` (String) The notification status for this app assignment. / Contains value for notification status; possible values are: `showAll` (Show all notifications.), `showReboot` (Only show restart notification and suppress other notifications.), `hideAll` (Hide all notifications.), `unknownFutureValue` (Unknown future value, reserved for future usage as expandable enum.). The _provider_ default value is `"showAll"`.
- `restart_settings` (Attributes) The reboot settings to apply for this app assignment. / Contains properties describing restart coordination following an app installation. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-wingetapprestartsettings?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--settings--winget--restart_settings))

<a id="nestedatt--assignments--settings--winget--install_time_settings"></a>
### Nested Schema for `assignments.settings.winget.install_time_settings`

Optional:

- `deadline_date_time` (String) The time at which the app should be installed. The _provider_ default value is `""`.
- `use_local_time` (Boolean) Whether the local device time or UTC time should be used when determining the deadline times. The _provider_ default value is `false`.


<a id="nestedatt--assignments--settings--winget--restart_settings"></a>
### Nested Schema for `assignments.settings.winget.restart_settings`

Optional:

- `countdown_display_before_restart_in_minutes` (Number) The number of minutes before the restart time to display the countdown dialog for pending restarts. The _provider_ default value is `0`.
- `grace_period_in_minutes` (Number) The number of minutes to wait before restarting the device after an app installation. The _provider_ default value is `0`.
- `restart_notification_snooze_duration_in_minutes` (Number) The number of minutes to snooze the restart notification dialog when the snooze button is selected.





<a id="nestedatt--categories"></a>
### Nested Schema for `categories`

Required:

- `id` (String) The key of the entity. This property is read-only.

Read-Only:

- `display_name` (String) The name of the app category.


<a id="nestedatt--ios_store"></a>
### Nested Schema for `ios_store`

Required:

- `app_store_url` (String) The Apple App Store URL
- `minimum_supported_operating_system` (Attributes) The value for the minimum applicable operating system. / Contains properties of the minimum operating system required for an iOS mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosminimumoperatingsystem?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_store--minimum_supported_operating_system))

Optional:

- `applicable_device_type` (Attributes) The iOS architecture for which this app can run on. / Contains properties of the possible iOS device types the mobile app can run on. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosdevicetype?view=graph-rest-beta. The _provider_ default value is `{}`. (see [below for nested schema](#nestedatt--ios_store--applicable_device_type))

Read-Only:

- `bundle_id` (String) The Identity Name.

<a id="nestedatt--ios_store--minimum_supported_operating_system"></a>
### Nested Schema for `ios_store.minimum_supported_operating_system`

Optional:

- `v10_0` (Boolean) When TRUE, only Version 10.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v11_0` (Boolean) When TRUE, only Version 11.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v12_0` (Boolean) When TRUE, only Version 12.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v13_0` (Boolean) When TRUE, only Version 13.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v14_0` (Boolean) When TRUE, only Version 14.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v15_0` (Boolean) When TRUE, only Version 15.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v16_0` (Boolean) When TRUE, only Version 16.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v17_0` (Boolean) When TRUE, only Version 17.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v8_0` (Boolean) When TRUE, only Version 8.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v9_0` (Boolean) When TRUE, only Version 9.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.


<a id="nestedatt--ios_store--applicable_device_type"></a>
### Nested Schema for `ios_store.applicable_device_type`

Optional:

- `ipad` (Boolean) Whether the app should run on iPads. The _provider_ default value is `true`.
- `iphone_and_ipod` (Boolean) Whether the app should run on iPhones and iPods. The _provider_ default value is `true`.



<a id="nestedatt--ios_webclip"></a>
### Nested Schema for `ios_webclip`

Required:

- `app_url` (String) Indicates iOS/iPadOS web clip app URL. Example: "https://www.contoso.com"

Optional:

- `full_screen_enabled` (Boolean) Whether or not to open the web clip as a full-screen web app. Defaults to false. If TRUE, opens the web clip as a full-screen web app. If FALSE, the web clip opens inside of another app, such as Safari or the app specified with targetApplicationBundleIdentifier.
- `ignore_manifest_scope` (Boolean) Whether or not a full screen web clip can navigate to an external web site without showing the Safari UI. Defaults to false. If FALSE, the Safari UI appears when navigating away. If TRUE, the Safari UI will not be shown.
- `pre_composed_icon_enabled` (Boolean) Whether or not the icon for the app is precomosed. Defaults to false. If TRUE, prevents SpringBoard from adding "shine" to the icon. If FALSE, SpringBoard can add "shine".
- `target_application_bundle_identifier` (String) Specifies the application bundle identifier which opens the URL. Available in iOS 14 and later.
- `use_managed_browser` (Boolean) Whether or not to use managed browser. When TRUE, the app will be required to be opened in Microsoft Edge. When FALSE, the app will not be required to be opened in Microsoft Edge. By default, this property is set to FALSE. The _provider_ default value is `false`.


<a id="nestedatt--large_icon"></a>
### Nested Schema for `large_icon`

Required:

- `type` (String) Indicates the content mime type.
- `value_base64` (String) The byte array that contains the actual content.


<a id="nestedatt--macos_dmg"></a>
### Nested Schema for `macos_dmg`

Required:

- `file_name` (String) The name of the main Lob application file.
- `included_apps` (Attributes Set) The list of .apps expected to be installed by the DMG (Apple Disk Image). This collection can contain a maximum of 500 elements. / Contains properties of an included .app in a MacOS app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosincludedapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_dmg--included_apps))
- `minimum_supported_operating_system` (Attributes) ComplexType macOSMinimumOperatingSystem that indicates the minimum operating system applicable for the application. / The minimum operating system required for a macOS app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosminimumoperatingsystem?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_dmg--minimum_supported_operating_system))
- `primary_bundle_id` (String) The bundleId of the primary .app in the DMG (Apple Disk Image). This maps to the CFBundleIdentifier in the app's bundle configuration.
- `primary_bundle_version` (String) The version of the primary .app in the DMG (Apple Disk Image). This maps to the CFBundleShortVersion in the app's bundle configuration.
- `source_file` (String) Provider Note: The path to the DMG file to be uploaded.
- `source_sha256` (String) Provider Note: The SHA-256 sum of the source file (see example). Changing this attribute will trigger a new upload of the source file.

Optional:

- `ignore_version_detection` (Boolean) When TRUE, indicates that the app's version will NOT be used to detect if the app is installed on a device. When FALSE, indicates that the app's version will be used to detect if the app is installed on a device. Set this to true for apps that use a self update feature. The default value is FALSE. The _provider_ default value is `false`.

Read-Only:

- `committed_content_version` (String) The internal committed content version.
- `size` (Number) The total size, including all uploaded files. This property is read-only.

<a id="nestedatt--macos_dmg--included_apps"></a>
### Nested Schema for `macos_dmg.included_apps`

Required:

- `bundle_id` (String) The bundleId of the app. This maps to the CFBundleIdentifier in the app's bundle configuration.
- `bundle_version` (String) The version of the app. This maps to the CFBundleShortVersion in the app's bundle configuration.


<a id="nestedatt--macos_dmg--minimum_supported_operating_system"></a>
### Nested Schema for `macos_dmg.minimum_supported_operating_system`

Optional:

- `v10_10` (Boolean) When TRUE, indicates OS X 10.10 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_11` (Boolean) When TRUE, indicates OS X 10.11 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_12` (Boolean) When TRUE, indicates macOS 10.12 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_13` (Boolean) When TRUE, indicates macOS 10.13 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_14` (Boolean) When TRUE, indicates macOS 10.14 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_15` (Boolean) When TRUE, indicates macOS 10.15 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_7` (Boolean) When TRUE, indicates Mac OS X 10.7 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_8` (Boolean) When TRUE, indicates OS X 10.8 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_9` (Boolean) When TRUE, indicates OS X 10.9 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v11_0` (Boolean) When TRUE, indicates macOS 11.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v12_0` (Boolean) When TRUE, indicates macOS 12.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v13_0` (Boolean) When TRUE, indicates macOS 13.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v14_0` (Boolean) When TRUE, indicates macOS 14.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.



<a id="nestedatt--macos_ms_defender"></a>
### Nested Schema for `macos_ms_defender`


<a id="nestedatt--macos_ms_edge"></a>
### Nested Schema for `macos_ms_edge`

Optional:

- `channel` (String) The channel to install on target devices. / The enum to specify the channels for Microsoft Edge apps; possible values are: `dev` (The Dev Channel is intended to help you plan and develop with the latest capabilities of Microsoft Edge.), `beta` (The Beta Channel is intended for production deployment to a representative sample set of users. New features ship about every 4 weeks. Security and quality updates ship as needed.), `stable` (The Stable Channel is intended for broad deployment within organizations, and it's the channel that most users should be on. New features ship about every 4 weeks. Security and quality updates ship as needed.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.). The _provider_ default value is `"stable"`.


<a id="nestedatt--macos_office_suite"></a>
### Nested Schema for `macos_office_suite`


<a id="nestedatt--macos_pkg"></a>
### Nested Schema for `macos_pkg`

Required:

- `file_name` (String) The name of the main Lob application file.
- `included_apps` (Attributes Set) The list of apps expected to be installed by the PKG. This collection can contain a maximum of 500 elements. / Contains properties of an included .app in a MacOS app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosincludedapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_pkg--included_apps))
- `minimum_supported_operating_system` (Attributes) ComplexType macOSMinimumOperatingSystem that indicates the minimum operating system applicable for the application. / The minimum operating system required for a macOS app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosminimumoperatingsystem?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_pkg--minimum_supported_operating_system))
- `primary_bundle_id` (String) The bundleId of the primary app in the PKG. This maps to the CFBundleIdentifier in the app's bundle configuration.
- `primary_bundle_version` (String) The version of the primary app in the PKG. This maps to the CFBundleShortVersion in the app's bundle configuration.
- `source_file` (String) Provider Note: The path to the PKG file to be uploaded.
- `source_sha256` (String) Provider Note: The SHA-256 sum of the source file (see example). Changing this attribute will trigger a new upload of the source file.

Optional:

- `ignore_version_detection` (Boolean) When TRUE, indicates that the app's version will NOT be used to detect if the app is installed on a device. When FALSE, indicates that the app's version will be used to detect if the app is installed on a device. Set this to true for apps that use a self update feature. The default value is FALSE. The _provider_ default value is `false`.
- `post_install_script` (Attributes) ComplexType macOSAppScript the contains the post-install script for the app. This will execute on the macOS device after the app is installed. / Shell script used to assist installation of a macOS app. These scripts are used to perform additional tasks to help the app successfully be configured. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosappscript?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_pkg--post_install_script))
- `pre_install_script` (Attributes) ComplexType macOSAppScript the contains the post-install script for the app. This will execute on the macOS device after the app is installed. / Shell script used to assist installation of a macOS app. These scripts are used to perform additional tasks to help the app successfully be configured. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosappscript?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_pkg--pre_install_script))

Read-Only:

- `committed_content_version` (String) The internal committed content version.
- `size` (Number) The total size, including all uploaded files. This property is read-only.

<a id="nestedatt--macos_pkg--included_apps"></a>
### Nested Schema for `macos_pkg.included_apps`

Required:

- `bundle_id` (String) The bundleId of the app. This maps to the CFBundleIdentifier in the app's bundle configuration.
- `bundle_version` (String) The version of the app. This maps to the CFBundleShortVersion in the app's bundle configuration.


<a id="nestedatt--macos_pkg--minimum_supported_operating_system"></a>
### Nested Schema for `macos_pkg.minimum_supported_operating_system`

Optional:

- `v10_10` (Boolean) When TRUE, indicates OS X 10.10 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_11` (Boolean) When TRUE, indicates OS X 10.11 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_12` (Boolean) When TRUE, indicates macOS 10.12 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_13` (Boolean) When TRUE, indicates macOS 10.13 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_14` (Boolean) When TRUE, indicates macOS 10.14 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_15` (Boolean) When TRUE, indicates macOS 10.15 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_7` (Boolean) When TRUE, indicates Mac OS X 10.7 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_8` (Boolean) When TRUE, indicates OS X 10.8 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_9` (Boolean) When TRUE, indicates OS X 10.9 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v11_0` (Boolean) When TRUE, indicates macOS 11.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v12_0` (Boolean) When TRUE, indicates macOS 12.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v13_0` (Boolean) When TRUE, indicates macOS 13.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v14_0` (Boolean) When TRUE, indicates macOS 14.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.


<a id="nestedatt--macos_pkg--post_install_script"></a>
### Nested Schema for `macos_pkg.post_install_script`

Optional:

- `script_content` (String) The base64 encoded shell script (.sh) that assists managing macOS apps. The _provider_ default value is `""`.


<a id="nestedatt--macos_pkg--pre_install_script"></a>
### Nested Schema for `macos_pkg.pre_install_script`

Optional:

- `script_content` (String) The base64 encoded shell script (.sh) that assists managing macOS apps. The _provider_ default value is `""`.



<a id="nestedatt--macos_webclip"></a>
### Nested Schema for `macos_webclip`

Required:

- `app_url` (String) The web app URL starting with http:// or https://, such as https://learn.microsoft.com/mem/.

Optional:

- `full_screen_enabled` (Boolean) Whether or not to open the web clip as a full-screen web app. Defaults to false. If TRUE, opens the web clip as a full-screen web app. If FALSE, the web clip opens inside of another app.
- `pre_composed_icon_enabled` (Boolean) Whether or not the icon for the app is precomosed. Defaults to false. If TRUE, prevents SpringBoard from adding "shine" to the icon. If FALSE, SpringBoard can add "shine".


<a id="nestedatt--managed_android_store"></a>
### Nested Schema for `managed_android_store`

Required:

- `app_store_url` (String) The Android AppStoreUrl.
- `minimum_supported_operating_system` (Attributes) The value for the minimum supported operating system. / Contains properties for the minimum operating system required for an Android mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidminimumoperatingsystem?view=graph-rest-beta (see [below for nested schema](#nestedatt--managed_android_store--minimum_supported_operating_system))
- `package_id` (String) The app's package ID.

Read-Only:

- `app_availability` (String) The Application's availability. This property is read-only. / A managed (MAM) application's availability; possible values are: `global` (A globally available app to all tenants.), `lineOfBusiness` (A line of business apps private to an organization.)
- `version` (String) The Application's version.

<a id="nestedatt--managed_android_store--minimum_supported_operating_system"></a>
### Nested Schema for `managed_android_store.minimum_supported_operating_system`

Optional:

- `v10_0` (Boolean) When TRUE, only Version 10.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v11_0` (Boolean) When TRUE, only Version 11.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v12_0` (Boolean) When TRUE, only Version 12.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v13_0` (Boolean) When TRUE, only Version 13.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v14_0` (Boolean) When TRUE, only Version 14.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v15_0` (Boolean) When TRUE, only Version 15.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_0` (Boolean) When TRUE, only Version 4.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_0_3` (Boolean) When TRUE, only Version 4.0.3 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_1` (Boolean) When TRUE, only Version 4.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_2` (Boolean) When TRUE, only Version 4.2 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_3` (Boolean) When TRUE, only Version 4.3 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_4` (Boolean) When TRUE, only Version 4.4 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v5_0` (Boolean) When TRUE, only Version 5.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v5_1` (Boolean) When TRUE, only Version 5.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v6_0` (Boolean) When TRUE, only Version 6.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v7_0` (Boolean) When TRUE, only Version 7.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v7_1` (Boolean) When TRUE, only Version 7.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v8_0` (Boolean) When TRUE, only Version 8.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v8_1` (Boolean) When TRUE, only Version 8.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v9_0` (Boolean) When TRUE, only Version 9.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.



<a id="nestedatt--managed_ios_store"></a>
### Nested Schema for `managed_ios_store`

Required:

- `app_store_url` (String) The Apple AppStoreUrl.
- `bundle_id` (String) The app's Bundle ID.
- `minimum_supported_operating_system` (Attributes) The value for the minimum supported operating system. / Contains properties of the minimum operating system required for an iOS mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosminimumoperatingsystem?view=graph-rest-beta (see [below for nested schema](#nestedatt--managed_ios_store--minimum_supported_operating_system))

Optional:

- `applicable_device_type` (Attributes) The iOS architecture for which this app can run on. / Contains properties of the possible iOS device types the mobile app can run on. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosdevicetype?view=graph-rest-beta. The _provider_ default value is `{}`. (see [below for nested schema](#nestedatt--managed_ios_store--applicable_device_type))

Read-Only:

- `app_availability` (String) The Application's availability. This property is read-only. / A managed (MAM) application's availability; possible values are: `global` (A globally available app to all tenants.), `lineOfBusiness` (A line of business apps private to an organization.)
- `version` (String) The Application's version.

<a id="nestedatt--managed_ios_store--minimum_supported_operating_system"></a>
### Nested Schema for `managed_ios_store.minimum_supported_operating_system`

Optional:

- `v10_0` (Boolean) When TRUE, only Version 10.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v11_0` (Boolean) When TRUE, only Version 11.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v12_0` (Boolean) When TRUE, only Version 12.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v13_0` (Boolean) When TRUE, only Version 13.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v14_0` (Boolean) When TRUE, only Version 14.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v15_0` (Boolean) When TRUE, only Version 15.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v16_0` (Boolean) When TRUE, only Version 16.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v17_0` (Boolean) When TRUE, only Version 17.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v8_0` (Boolean) When TRUE, only Version 8.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v9_0` (Boolean) When TRUE, only Version 9.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.


<a id="nestedatt--managed_ios_store--applicable_device_type"></a>
### Nested Schema for `managed_ios_store.applicable_device_type`

Optional:

- `ipad` (Boolean) Whether the app should run on iPads. The _provider_ default value is `true`.
- `iphone_and_ipod` (Boolean) Whether the app should run on iPhones and iPods. The _provider_ default value is `true`.



<a id="nestedatt--web_link"></a>
### Nested Schema for `web_link`

Required:

- `app_url` (String) The web app URL. This property cannot be PATCHed.

Optional:

- `use_managed_browser` (Boolean) Whether or not to use managed browser. This property is only applicable for Android and IOS. The _provider_ default value is `false`.


<a id="nestedatt--win32_lob"></a>
### Nested Schema for `win32_lob`

Required:

- `file_name` (String) The name of the main Lob application file.
- `install_command_line` (String) The command line to install this app
- `setup_file_path` (String) The relative path of the setup file in the encrypted Win32LobApp package.
- `source_file` (String) Provider Note: The path to the IntuneWin file to be uploaded.
- `source_sha256` (String) Provider Note: The SHA-256 sum of the source file (see example). Changing this attribute will trigger a new upload of the source file.
- `uninstall_command_line` (String) The command line to uninstall this app

Optional:

- `allow_available_uninstall` (Boolean) When TRUE, indicates that uninstall is supported from the company portal for the Windows app (Win32) with an Available assignment. When FALSE, indicates that uninstall is not supported for the Windows app (Win32) with an Available assignment. Default value is FALSE. The _provider_ default value is `true`.
- `applicable_architectures` (String) The Windows architecture(s) for which this app can run on. / Contains properties for Windows architecture; possible values are: `none` (No flags set.), `x86` (Whether or not the X86 Windows architecture type is supported.), `x64` (Whether or not the X64 Windows architecture type is supported.), `arm` (Whether or not the Arm Windows architecture type is supported.), `neutral` (Whether or not the Neutral Windows architecture type is supported.), `arm64` (Whether or not the Arm64 Windows architecture type is supported.). The _provider_ default value is `"x86,x64"`.
- `detection_rules` (Attributes Set) The detection rules to detect Win32 Line of Business (LoB) app. / Base class to detect a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobappdetection?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--win32_lob--detection_rules))
- `display_version` (String) The version displayed in the UX for this app.
- `install_experience` (Attributes) The install experience for this app. / Contains installation experience properties for a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobappinstallexperience?view=graph-rest-beta. The _provider_ default value is `{}`. (see [below for nested schema](#nestedatt--win32_lob--install_experience))
- `minimum_cpu_speed_in_mhz` (Number) The value for the minimum CPU speed which is required to install this app.
- `minimum_free_disk_space_in_mb` (Number) The value for the minimum free disk space which is required to install this app.
- `minimum_memory_in_mb` (Number) The value for the minimum physical memory which is required to install this app.
- `minimum_number_of_processors` (Number) The value for the minimum number of processors which is required to install this app.
- `minimum_supported_windows_release` (String) The value for the minimum supported windows release. The _provider_ default value is `"1809"`.
- `msi_information` (Attributes) The MSI details if this Win32 app is an MSI app. / Contains MSI app properties for a Win32 App. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobappmsiinformation?view=graph-rest-beta (see [below for nested schema](#nestedatt--win32_lob--msi_information))
- `requirement_rules` (Attributes Set) The requirement rules to detect Win32 Line of Business (LoB) app. / Base class to detect a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobapprequirement?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--win32_lob--requirement_rules))
- `return_codes` (Attributes Set) The return codes for post installation behavior. / Contains return code properties for a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobappreturncode?view=graph-rest-beta. The _provider_ default value is `mobileAppWin32LobAppReturnCodesDefault`. (see [below for nested schema](#nestedatt--win32_lob--return_codes))

Read-Only:

- `committed_content_version` (String) The internal committed content version.
- `rules` (Attributes Set) The detection and requirement rules for this app. / A base complex type to store the detection or requirement rule data for a Win32 LOB app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobapprule?view=graph-rest-beta (see [below for nested schema](#nestedatt--win32_lob--rules))
- `size` (Number) The total size, including all uploaded files. This property is read-only.

<a id="nestedatt--win32_lob--detection_rules"></a>
### Nested Schema for `win32_lob.detection_rules`

Optional:

- `filesystem` (Attributes) Contains file or folder path to detect a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobappfilesystemdetection?view=graph-rest-beta (see [below for nested schema](#nestedatt--win32_lob--detection_rules--filesystem))
- `powershell_script` (Attributes) Contains PowerShell script properties to detect a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobapppowershellscriptdetection?view=graph-rest-beta (see [below for nested schema](#nestedatt--win32_lob--detection_rules--powershell_script))
- `product_code` (Attributes) Contains product code and version properties to detect a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobappproductcodedetection?view=graph-rest-beta (see [below for nested schema](#nestedatt--win32_lob--detection_rules--product_code))
- `registry` (Attributes) Contains registry properties to detect a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobappregistrydetection?view=graph-rest-beta (see [below for nested schema](#nestedatt--win32_lob--detection_rules--registry))

<a id="nestedatt--win32_lob--detection_rules--filesystem"></a>
### Nested Schema for `win32_lob.detection_rules.filesystem`

Required:

- `detection_type` (String) The file system detection type. / Contains all supported file system detection type; possible values are: `notConfigured` (Not configured.), `exists` (Whether the specified file or folder exists.), `modifiedDate` (Last modified date.), `createdDate` (Created date.), `version` (Version value type.), `sizeInMB` (Size detection type.), `doesNotExist` (The specified file or folder does not exist.)
- `file_or_folder_name` (String) The file or folder name to detect Win32 Line of Business (LoB) app
- `path` (String) The file or folder path to detect Win32 Line of Business (LoB) app

Optional:

- `check_32_bit_on_64_system` (Boolean) A value indicating whether this file or folder is for checking 32-bit app on 64-bit system. The _provider_ default value is `false`.
- `detection_value` (String) The file or folder detection value
- `operator` (String) The operator for file or folder detection. / Contains properties for detection operator; possible values are: `notConfigured` (Not configured.), `equal` (Equal operator.), `notEqual` (Not equal operator.), `greaterThan` (Greater than operator.), `greaterThanOrEqual` (Greater than or equal operator.), `lessThan` (Less than operator.), `lessThanOrEqual` (Less than or equal operator.). The _provider_ default value is `"notConfigured"`.


<a id="nestedatt--win32_lob--detection_rules--powershell_script"></a>
### Nested Schema for `win32_lob.detection_rules.powershell_script`

Required:

- `script_content` (String) The base64 encoded script content to detect Win32 Line of Business (LoB) app

Optional:

- `enforce_signature_check` (Boolean) A value indicating whether signature check is enforced. The _provider_ default value is `false`.
- `run_as_32_bit` (Boolean) A value indicating whether this script should run as 32-bit. The _provider_ default value is `false`.


<a id="nestedatt--win32_lob--detection_rules--product_code"></a>
### Nested Schema for `win32_lob.detection_rules.product_code`

Required:

- `product_code` (String) The product code of Win32 Line of Business (LoB) app.

Optional:

- `product_version` (String) The product version of Win32 Line of Business (LoB) app.
- `product_version_operator` (String) The operator to detect product version. / Contains properties for detection operator; possible values are: `notConfigured` (Not configured.), `equal` (Equal operator.), `notEqual` (Not equal operator.), `greaterThan` (Greater than operator.), `greaterThanOrEqual` (Greater than or equal operator.), `lessThan` (Less than operator.), `lessThanOrEqual` (Less than or equal operator.). The _provider_ default value is `"notConfigured"`.


<a id="nestedatt--win32_lob--detection_rules--registry"></a>
### Nested Schema for `win32_lob.detection_rules.registry`

Required:

- `detection_type` (String) The registry data detection type. / Contains all supported registry data detection type; possible values are: `notConfigured` (Not configured.), `exists` (The specified registry key or value exists.), `doesNotExist` (The specified registry key or value does not exist.), `string` (String value type.), `integer` (Integer value type.), `version` (Version value type.)
- `key_path` (String) The registry key path to detect Win32 Line of Business (LoB) app

Optional:

- `check_32_bit_on_64_system` (Boolean) A value indicating whether this registry path is for checking 32-bit app on 64-bit system. The _provider_ default value is `false`.
- `detection_value` (String) The registry detection value
- `operator` (String) The operator for registry data detection. / Contains properties for detection operator; possible values are: `notConfigured` (Not configured.), `equal` (Equal operator.), `notEqual` (Not equal operator.), `greaterThan` (Greater than operator.), `greaterThanOrEqual` (Greater than or equal operator.), `lessThan` (Less than operator.), `lessThanOrEqual` (Less than or equal operator.). The _provider_ default value is `"notConfigured"`.
- `value_name` (String) The registry value name



<a id="nestedatt--win32_lob--install_experience"></a>
### Nested Schema for `win32_lob.install_experience`

Optional:

- `device_restart_behavior` (String) Device restart behavior. / Indicates the type of restart action; possible values are: `basedOnReturnCode` (Intune will restart the device after running the app installation if the operation returns a reboot code.), `allow` (Intune will not take any specific action on reboot codes resulting from app installations. Intune will not attempt to suppress restarts for MSI apps.), `suppress` (Intune will attempt to suppress restarts for MSI apps.), `force` (Intune will force the device to restart immediately after the app installation operation.). The _provider_ default value is `"allow"`.
- `max_run_time_in_minutes` (Number) The number of minutes the system will wait for install program to finish. Default value is 60 minutes. The _provider_ default value is `60`.
- `run_as_account` (String) Indicates the type of execution context the app runs in. / Indicates the type of execution context the app runs in; possible values are: `system` (System context), `user` (User context). The _provider_ default value is `"system"`.


<a id="nestedatt--win32_lob--msi_information"></a>
### Nested Schema for `win32_lob.msi_information`

Optional:

- `package_type` (String) The MSI package type. / Indicates the package type of an MSI Win32LobApp; possible values are: `perMachine` (Indicates a per-machine app package.), `perUser` (Indicates a per-user app package.), `dualPurpose` (Indicates a dual-purpose app package.). The _provider_ default value is `"perMachine"`.
- `product_code` (String) The MSI product code.
- `product_name` (String) The MSI product name.
- `product_version` (String) The MSI product version.
- `publisher` (String) The MSI publisher.
- `requires_reboot` (Boolean) Whether the MSI app requires the machine to reboot to complete installation. The _provider_ default value is `false`.
- `upgrade_code` (String) The MSI upgrade code.


<a id="nestedatt--win32_lob--requirement_rules"></a>
### Nested Schema for `win32_lob.requirement_rules`

Optional:

- `filesystem` (Attributes) Contains file or folder path to detect a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobappfilesystemrequirement?view=graph-rest-beta (see [below for nested schema](#nestedatt--win32_lob--requirement_rules--filesystem))
- `powershell_script` (Attributes) Contains PowerShell script properties to detect a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobapppowershellscriptrequirement?view=graph-rest-beta (see [below for nested schema](#nestedatt--win32_lob--requirement_rules--powershell_script))
- `registry` (Attributes) Contains registry properties to detect a Win32 App / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-win32lobappregistryrequirement?view=graph-rest-beta (see [below for nested schema](#nestedatt--win32_lob--requirement_rules--registry))

<a id="nestedatt--win32_lob--requirement_rules--filesystem"></a>
### Nested Schema for `win32_lob.requirement_rules.filesystem`

Required:

- `detection_type` (String) The file system detection type. / Contains all supported file system detection type; possible values are: `notConfigured` (Not configured.), `exists` (Whether the specified file or folder exists.), `modifiedDate` (Last modified date.), `createdDate` (Created date.), `version` (Version value type.), `sizeInMB` (Size detection type.), `doesNotExist` (The specified file or folder does not exist.)
- `file_or_folder_name` (String) The file or folder name to detect Win32 Line of Business (LoB) app
- `path` (String) The file or folder path to detect Win32 Line of Business (LoB) app

Optional:

- `check_32_bit_on_64_system` (Boolean) A value indicating whether this file or folder is for checking 32-bit app on 64-bit system. The _provider_ default value is `false`.
- `detection_value` (String) The detection value
- `operator` (String) The operator for detection / Contains properties for detection operator; possible values are: `notConfigured` (Not configured.), `equal` (Equal operator.), `notEqual` (Not equal operator.), `greaterThan` (Greater than operator.), `greaterThanOrEqual` (Greater than or equal operator.), `lessThan` (Less than operator.), `lessThanOrEqual` (Less than or equal operator.). The _provider_ default value is `"notConfigured"`.


<a id="nestedatt--win32_lob--requirement_rules--powershell_script"></a>
### Nested Schema for `win32_lob.requirement_rules.powershell_script`

Required:

- `detection_type` (String) The detection type for script output. / Contains all supported Powershell Script output detection type; possible values are: `notConfigured` (Not configured.), `string` (Output data type is string.), `dateTime` (Output data type is date time.), `integer` (Output data type is integer.), `float` (Output data type is float.), `version` (Output data type is version.), `boolean` (Output data type is boolean.)
- `script_content` (String) The base64 encoded script content to detect Win32 Line of Business (LoB) app

Optional:

- `detection_value` (String) The detection value
- `display_name` (String) The unique display name for this rule. The _provider_ default value is `""`.
- `enforce_signature_check` (Boolean) A value indicating whether signature check is enforced. The _provider_ default value is `false`.
- `operator` (String) The operator for detection / Contains properties for detection operator; possible values are: `notConfigured` (Not configured.), `equal` (Equal operator.), `notEqual` (Not equal operator.), `greaterThan` (Greater than operator.), `greaterThanOrEqual` (Greater than or equal operator.), `lessThan` (Less than operator.), `lessThanOrEqual` (Less than or equal operator.). The _provider_ default value is `"notConfigured"`.
- `run_as_32_bit` (Boolean) A value indicating whether this script should run as 32-bit. The _provider_ default value is `false`.
- `run_as_account` (String) Indicates the type of execution context the script runs in. / Indicates the type of execution context the app runs in; possible values are: `system` (System context), `user` (User context). The _provider_ default value is `"system"`.


<a id="nestedatt--win32_lob--requirement_rules--registry"></a>
### Nested Schema for `win32_lob.requirement_rules.registry`

Required:

- `detection_type` (String) The registry data detection type. / Contains all supported registry data detection type; possible values are: `notConfigured` (Not configured.), `exists` (The specified registry key or value exists.), `doesNotExist` (The specified registry key or value does not exist.), `string` (String value type.), `integer` (Integer value type.), `version` (Version value type.)
- `key_path` (String) The registry key path to detect Win32 Line of Business (LoB) app

Optional:

- `check_32_bit_on_64_system` (Boolean) A value indicating whether this registry path is for checking 32-bit app on 64-bit system. The _provider_ default value is `false`.
- `detection_value` (String) The detection value
- `operator` (String) The operator for detection / Contains properties for detection operator; possible values are: `notConfigured` (Not configured.), `equal` (Equal operator.), `notEqual` (Not equal operator.), `greaterThan` (Greater than operator.), `greaterThanOrEqual` (Greater than or equal operator.), `lessThan` (Less than operator.), `lessThanOrEqual` (Less than or equal operator.). The _provider_ default value is `"notConfigured"`.
- `value_name` (String) The registry value name



<a id="nestedatt--win32_lob--return_codes"></a>
### Nested Schema for `win32_lob.return_codes`

Required:

- `return_code` (Number) Return code.
- `type` (String) The type of return code. / Indicates the type of return code; possible values are: `failed` (Failed.), `success` (Success.), `softReboot` (Soft-reboot is required.), `hardReboot` (Hard-reboot is required.), `retry` (Retry.)


<a id="nestedatt--win32_lob--rules"></a>
### Nested Schema for `win32_lob.rules`

Read-Only:

- `rule_type` (String) The rule type indicating the purpose of the rule. / Contains rule types for Win32 LOB apps; possible values are: `detection` (Detection rule.), `requirement` (Requirement rule.)



<a id="nestedatt--windows_ms_edge"></a>
### Nested Schema for `windows_ms_edge`

Optional:

- `channel` (String) The channel to install on target devices. The possible values are dev, beta, and stable. By default, this property is set to dev. / The enum to specify the channels for Microsoft Edge apps; possible values are: `dev` (The Dev Channel is intended to help you plan and develop with the latest capabilities of Microsoft Edge.), `beta` (The Beta Channel is intended for production deployment to a representative sample set of users. New features ship about every 4 weeks. Security and quality updates ship as needed.), `stable` (The Stable Channel is intended for broad deployment within organizations, and it's the channel that most users should be on. New features ship about every 4 weeks. Security and quality updates ship as needed.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.). The _provider_ default value is `"stable"`.
- `display_language_locale` (String) The language locale to use when the Edge app displays text to the user.


<a id="nestedatt--windows_msi"></a>
### Nested Schema for `windows_msi`

Required:

- `file_name` (String) The name of the main Lob application file.
- `identity_version` (String) The identity version.
- `manifest_base64` (String) Provider Note: The manifest contains settings that are based on MSI properties (that cannot easily be parsed and read by the provider). Instead, either manually read the MSI file's properties (e.g. using MSI Orca on Windows) and based on their values select the appropriate values for the XML attributes (as shown in detail in the example above). Or open the Intune portal, enable network log recording in your browser's development tool and thereafter upload the MSI from there (creating a dummy app). Then search the network log for a POST request to https://graph.microsoft.com/beta/deviceAppManagement/mobileApps/.../microsoft.graph.windowsMobileMSI/contentVersions/.../files and directly copy the base64 encoded manifest value from its payload.
- `product_version` (String) The product version of Windows Mobile MSI Line of Business (LoB) app.
- `source_file` (String) Provider Note: The path to the MSI file to be uploaded.
- `source_sha256` (String) Provider Note: The SHA-256 sum of the source file (see example). Changing this attribute will trigger a new upload of the source file.

Optional:

- `command_line` (String) The command line.
- `ignore_version_detection` (Boolean) A boolean to control whether the app's version will be used to detect the app after it is installed on a device. Set this to true for Windows Mobile MSI Line of Business (LoB) apps that use a self update feature. The _provider_ default value is `false`.
- `product_code` (String) The product code.
- `use_device_context` (Boolean) Indicates whether to install a dual-mode MSI in the device context. If true, app will be installed for all users. If false, app will be installed per-user. If null, service will use the MSI package's default install context. In case of dual-mode MSI, this default will be per-user.  Cannot be set for non-dual-mode apps.  Cannot be changed after initial creation of the application.

Read-Only:

- `committed_content_version` (String) The internal committed content version.
- `size` (Number) The total size, including all uploaded files. This property is read-only.


<a id="nestedatt--windows_office_suite"></a>
### Nested Schema for `windows_office_suite`

Optional:

- `auto_accept_eula` (Boolean) The value to accept the EULA automatically on the enduser's device. The _provider_ default value is `true`.
- `excluded_apps` (Attributes) The property to represent the apps which are excluded from the selected Office365 Product Id. / Contains properties for Excluded Office365 Apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-excludedapps?view=graph-rest-beta. The _provider_ default value is `{}`. (see [below for nested schema](#nestedatt--windows_office_suite--excluded_apps))
- `install_progress_display_level` (String) To specify the level of display for the Installation Progress Setup UI on the Device. / The Enum to specify the level of display for the Installation Progress Setup UI on the Device; possible values are: `none`, `full`. The _provider_ default value is `"none"`.
- `locales_to_install` (Set of String) The property to represent the locales which are installed when the apps from Office365 is installed. It uses standard RFC 6033. Ref: https://technet.microsoft.com/library/cc179219(v=office.16).aspx. The _provider_ default value is `[]`.
- `office_configuration_xml` (String) The property to represent the XML configuration file that can be specified for Office ProPlus Apps. Takes precedence over all other properties. When present, the XML configuration file will be used to create the app.
- `office_platform_architecture` (String) The property to represent the Office365 app suite version. / Contains properties for Windows architecture; possible values are: `none` (No flags set.), `x86` (Whether or not the X86 Windows architecture type is supported.), `x64` (Whether or not the X64 Windows architecture type is supported.), `arm` (Whether or not the Arm Windows architecture type is supported.), `neutral` (Whether or not the Neutral Windows architecture type is supported.), `arm64` (Whether or not the Arm64 Windows architecture type is supported.). The _provider_ default value is `"x64"`.
- `office_suite_app_default_file_format` (String) The property to represent the Office365 default file format type. / Describes the OfficeSuiteApp file format types that can be selected; possible values are: `notConfigured` (No file format selected), `officeOpenXMLFormat` (Office Open XML Format selected), `officeOpenDocumentFormat` (Office Open Document Format selected), `unknownFutureValue` (Placeholder for evolvable enum.). The _provider_ default value is `"officeOpenDocumentFormat"`.
- `product_ids` (Set of String) The Product Ids that represent the Office365 Suite SKU. / The Enum to specify the Office365 ProductIds that represent the Office365 Suite SKUs; possible values are: `o365ProPlusRetail`, `o365BusinessRetail`, `visioProRetail`, `projectProRetail`. The _provider_ default value is `[]`.
- `should_uninstall_older_versions_of_office` (Boolean) The property to determine whether to uninstall existing Office MSI if an Office365 app suite is deployed to the device or not. The _provider_ default value is `true`.
- `target_version` (String) The property to represent the specific target version for the Office365 app suite that should be remained deployed on the devices. The _provider_ default value is `""`.
- `update_channel` (String) The property to represent the Office365 Update Channel. / The Enum to specify the Office365 Updates Channel; possible values are: `none`, `current`, `deferred`, `firstReleaseCurrent`, `firstReleaseDeferred`, `monthlyEnterprise`. The _provider_ default value is `"current"`.
- `update_version` (String) The property to represent the update version in which the specific target version is available for the Office365 app suite. The _provider_ default value is `""`.
- `use_shared_computer_activation` (Boolean) The property to represent that whether the shared computer activation is used not for Office365 app suite. The _provider_ default value is `false`.

<a id="nestedatt--windows_office_suite--excluded_apps"></a>
### Nested Schema for `windows_office_suite.excluded_apps`

Optional:

- `access` (Boolean) The value for if MS Office Access should be excluded or not. The _provider_ default value is `false`.
- `bing` (Boolean) The value for if Microsoft Search as default should be excluded or not. The _provider_ default value is `false`.
- `excel` (Boolean) The value for if MS Office Excel should be excluded or not. The _provider_ default value is `false`.
- `groove` (Boolean) The value for if MS Office OneDrive for Business - Groove should be excluded or not. The _provider_ default value is `true`.
- `infopath` (Boolean) The value for if MS Office InfoPath should be excluded or not. The _provider_ default value is `true`.
- `lync` (Boolean) The value for if MS Office Skype for Business - Lync should be excluded or not. The _provider_ default value is `true`.
- `onedrive` (Boolean) The value for if MS Office OneDrive should be excluded or not. The _provider_ default value is `false`.
- `onenote` (Boolean) The value for if MS Office OneNote should be excluded or not. The _provider_ default value is `false`.
- `outlook` (Boolean) The value for if MS Office Outlook should be excluded or not. The _provider_ default value is `false`.
- `powerpoint` (Boolean) The value for if MS Office PowerPoint should be excluded or not. The _provider_ default value is `false`.
- `publisher` (Boolean) The value for if MS Office Publisher should be excluded or not. The _provider_ default value is `false`.
- `sharepoint_designer` (Boolean) The value for if MS Office SharePointDesigner should be excluded or not. The _provider_ default value is `true`.
- `teams` (Boolean) The value for if MS Office Teams should be excluded or not. The _provider_ default value is `false`.
- `visio` (Boolean) The value for if MS Office Visio should be excluded or not. The _provider_ default value is `false`.
- `word` (Boolean) The value for if MS Office Word should be excluded or not. The _provider_ default value is `false`.



<a id="nestedatt--windows_store"></a>
### Nested Schema for `windows_store`

Required:

- `app_store_url` (String) The Windows app store URL.


<a id="nestedatt--windows_universal_appx"></a>
### Nested Schema for `windows_universal_appx`

Required:

- `applicable_architectures` (String) The Windows architecture(s) for which this app can run on.; default value is `none`. / Contains properties for Windows architecture; possible values are: `none` (No flags set.), `x86` (Whether or not the X86 Windows architecture type is supported.), `x64` (Whether or not the X64 Windows architecture type is supported.), `arm` (Whether or not the Arm Windows architecture type is supported.), `neutral` (Whether or not the Neutral Windows architecture type is supported.), `arm64` (Whether or not the Arm64 Windows architecture type is supported.)
- `applicable_device_types` (String) The Windows device type(s) for which this app can run on.; default value is `none`. / Contains properties for Windows device type. Multiple values can be selected. Default value is `none`; possible values are: `none` (No device types supported. Default value.), `desktop` (Indicates support for Desktop Windows device type.), `mobile` (Indicates support for Mobile Windows device type.), `holographic` (Indicates support for Holographic Windows device type.), `team` (Indicates support for Team Windows device type.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.)
- `file_name` (String) The name of the main Lob application file.
- `identity_publisher_hash` (String) The Identity Publisher Hash of the app, parsed from the appx file when it is uploaded through the Intune MEM console. For example: "AB82CD0XYZ".
- `identity_version` (String) The Identity Version of the app, parsed from the appx file when it is uploaded through the Intune MEM console.  For example: "1.0.0.0".
- `minimum_supported_operating_system` (Attributes) The value for the minimum applicable Windows operating system. / The minimum operating system required for a Windows mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-windowsminimumoperatingsystem?view=graph-rest-beta  
Provider Note: v10_0 seems to be required and the only minimum OS version to be accepted by MS Graph for Univeral APPX. (see [below for nested schema](#nestedatt--windows_universal_appx--minimum_supported_operating_system))
- `source_file` (String) Provider Note: The path to the APPX/MSIX file to be uploaded.
- `source_sha256` (String) Provider Note: The SHA-256 sum of the source file (see example). Changing this attribute will trigger a new upload of the source file.

Optional:

- `identity_name` (String) The Identity Name of the app, parsed from the appx file when it is uploaded through the Intune MEM console. For example: "Contoso.DemoApp".
- `identity_resource_identifier` (String) The Identity Resource Identifier of the app, parsed from the appx file when it is uploaded through the Intune MEM console. For example: "TestResourceId".
- `is_bundle` (Boolean) Whether or not the app is a bundle. If TRUE, app is a bundle; if FALSE, app is not a bundle. The _provider_ default value is `false`.

Read-Only:

- `committed_contained_apps` (Attributes Set) The collection of contained apps in the committed mobileAppContent of a windowsUniversalAppX app. This property is read-only. / An abstract class that represents a contained app in a mobileApp acting as a package. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobilecontainedapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows_universal_appx--committed_contained_apps))
- `committed_content_version` (String) The internal committed content version.
- `size` (Number) The total size, including all uploaded files. This property is read-only.

<a id="nestedatt--windows_universal_appx--minimum_supported_operating_system"></a>
### Nested Schema for `windows_universal_appx.minimum_supported_operating_system`

Optional:

- `v10_0` (Boolean) Windows version 10.0 or later. The _provider_ default value is `false`.
- `v10_1607` (Boolean) Windows 10 1607 or later. The _provider_ default value is `false`.
- `v10_1703` (Boolean) Windows 10 1703 or later. The _provider_ default value is `false`.
- `v10_1709` (Boolean) Windows 10 1709 or later. The _provider_ default value is `false`.
- `v10_1803` (Boolean) Windows 10 1803 or later. The _provider_ default value is `false`.
- `v10_1809` (Boolean) Windows 10 1809 or later. The _provider_ default value is `false`.
- `v10_1903` (Boolean) Windows 10 1903 or later. The _provider_ default value is `false`.
- `v10_1909` (Boolean) Windows 10 1909 or later. The _provider_ default value is `false`.
- `v10_2004` (Boolean) Windows 10 2004 or later. The _provider_ default value is `false`.
- `v10_21_h1` (Boolean) Windows 10 21H1 or later. The _provider_ default value is `false`.
- `v10_2_h20` (Boolean) Windows 10 2H20 or later. The _provider_ default value is `false`.
- `v8_0` (Boolean) Windows version 8.0 or later. The _provider_ default value is `false`.
- `v8_1` (Boolean) Windows version 8.1 or later. The _provider_ default value is `false`.


<a id="nestedatt--windows_universal_appx--committed_contained_apps"></a>
### Nested Schema for `windows_universal_appx.committed_contained_apps`

Read-Only:

- `id` (String) Key of the entity. This property is read-only.



<a id="nestedatt--windows_web_link"></a>
### Nested Schema for `windows_web_link`

Required:

- `app_url` (String) Indicates the Windows web app URL. Example: "https://www.contoso.com"


<a id="nestedatt--winget"></a>
### Nested Schema for `winget`

Required:

- `package_identifier` (String) The PackageIdentifier from the WinGet source repository REST API. This also maps to the Id when using the WinGet client command line application. Required at creation time, cannot be modified on existing objects.

Optional:

- `install_experience` (Attributes) The install experience settings associated with this application, which are used to ensure the desired install experiences on the target device are taken into account. This includes the account type (System or User) that actions should be run as on target devices. Required at creation time. / Represents the install experience settings associated with WinGet apps. This is used to ensure the desired install experiences on the target device are taken into account. Required at creation time. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-wingetappinstallexperience?view=graph-rest-beta. The _provider_ default value is `{}`. (see [below for nested schema](#nestedatt--winget--install_experience))
- `manifest_hash` (String) Hash of package metadata properties used to validate that the application matches the metadata in the source repository.

<a id="nestedatt--winget--install_experience"></a>
### Nested Schema for `winget.install_experience`

Optional:

- `run_as_account` (String) Indicates the type of execution context the app setup runs in on target devices. Options include values of the RunAsAccountType enum, which are System and User. Required at creation time, cannot be modified on existing objects. / Indicates the type of execution context the app runs in; possible values are: `system` (System context), `user` (User context). The _provider_ default value is `"system"`.



<a id="nestedatt--android_for_work"></a>
### Nested Schema for `android_for_work`

Read-Only:

- `app_identifier` (String) The Identity Name. This property is read-only.
- `app_store_url` (String) The Play for Work Store app URL.
- `package_id` (String) The package identifier. This property is read-only.
- `total_license_count` (Number) The total number of VPP licenses.
- `used_license_count` (Number) The number of VPP licenses in use.


<a id="nestedatt--android_lob"></a>
### Nested Schema for `android_lob`

Read-Only:

- `committed_content_version` (String) The internal committed content version.
- `file_name` (String) The name of the main Lob application file.
- `minimum_supported_operating_system` (Attributes) The value for the minimum applicable operating system. / Contains properties for the minimum operating system required for an Android mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidminimumoperatingsystem?view=graph-rest-beta (see [below for nested schema](#nestedatt--android_lob--minimum_supported_operating_system))
- `package_id` (String) The package identifier.
- `size` (Number) The total size, including all uploaded files. This property is read-only.
- `targeted_platforms` (String) The platforms to which the application can be targeted. If not specified, will defauilt to Android Device Administrator. / Specifies which platform(s) can be targeted for a given Android LOB application or Managed Android LOB application; possible values are: `androidDeviceAdministrator` (Indicates the Android targeted platform is Android Device Administrator.), `androidOpenSourceProject` (Indicates the Android targeted platform is Android Open Source Project.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.)
- `version_code` (String) The version code of Android Line of Business (LoB) app.
- `version_name` (String) The version name of Android Line of Business (LoB) app.

<a id="nestedatt--android_lob--minimum_supported_operating_system"></a>
### Nested Schema for `android_lob.minimum_supported_operating_system`

Optional:

- `v10_0` (Boolean) When TRUE, only Version 10.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v11_0` (Boolean) When TRUE, only Version 11.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v12_0` (Boolean) When TRUE, only Version 12.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v13_0` (Boolean) When TRUE, only Version 13.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v14_0` (Boolean) When TRUE, only Version 14.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v15_0` (Boolean) When TRUE, only Version 15.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_0` (Boolean) When TRUE, only Version 4.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_0_3` (Boolean) When TRUE, only Version 4.0.3 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_1` (Boolean) When TRUE, only Version 4.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_2` (Boolean) When TRUE, only Version 4.2 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_3` (Boolean) When TRUE, only Version 4.3 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_4` (Boolean) When TRUE, only Version 4.4 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v5_0` (Boolean) When TRUE, only Version 5.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v5_1` (Boolean) When TRUE, only Version 5.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v6_0` (Boolean) When TRUE, only Version 6.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v7_0` (Boolean) When TRUE, only Version 7.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v7_1` (Boolean) When TRUE, only Version 7.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v8_0` (Boolean) When TRUE, only Version 8.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v8_1` (Boolean) When TRUE, only Version 8.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v9_0` (Boolean) When TRUE, only Version 9.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.



<a id="nestedatt--ios_lob"></a>
### Nested Schema for `ios_lob`

Read-Only:

- `applicable_device_type` (Attributes) The iOS architecture for which this app can run on. / Contains properties of the possible iOS device types the mobile app can run on. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosdevicetype?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_lob--applicable_device_type))
- `build_number` (String) The build number of iOS Line of Business (LoB) app.
- `bundle_id` (String) The Identity Name.
- `committed_content_version` (String) The internal committed content version.
- `expiration_date_time` (String) The expiration time.
- `file_name` (String) The name of the main Lob application file.
- `minimum_supported_operating_system` (Attributes) The value for the minimum applicable operating system. / Contains properties of the minimum operating system required for an iOS mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosminimumoperatingsystem?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_lob--minimum_supported_operating_system))
- `size` (Number) The total size, including all uploaded files. This property is read-only.
- `version_number` (String) The version number of iOS Line of Business (LoB) app.

<a id="nestedatt--ios_lob--applicable_device_type"></a>
### Nested Schema for `ios_lob.applicable_device_type`

Optional:

- `ipad` (Boolean) Whether the app should run on iPads. The _provider_ default value is `true`.
- `iphone_and_ipod` (Boolean) Whether the app should run on iPhones and iPods. The _provider_ default value is `true`.


<a id="nestedatt--ios_lob--minimum_supported_operating_system"></a>
### Nested Schema for `ios_lob.minimum_supported_operating_system`

Optional:

- `v10_0` (Boolean) When TRUE, only Version 10.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v11_0` (Boolean) When TRUE, only Version 11.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v12_0` (Boolean) When TRUE, only Version 12.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v13_0` (Boolean) When TRUE, only Version 13.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v14_0` (Boolean) When TRUE, only Version 14.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v15_0` (Boolean) When TRUE, only Version 15.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v16_0` (Boolean) When TRUE, only Version 16.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v17_0` (Boolean) When TRUE, only Version 17.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v8_0` (Boolean) When TRUE, only Version 8.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v9_0` (Boolean) When TRUE, only Version 9.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.



<a id="nestedatt--ios_vpp"></a>
### Nested Schema for `ios_vpp`

Read-Only:

- `app_store_url` (String) The store URL.
- `applicable_device_type` (Attributes) The applicable iOS Device Type. / Contains properties of the possible iOS device types the mobile app can run on. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosdevicetype?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_vpp--applicable_device_type))
- `assigned_licenses` (Attributes Set) The licenses assigned to this app. / iOS Volume Purchase Program license assignment. This class does not support Create, Delete, or Update. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosvppappassignedlicense?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_vpp--assigned_licenses))
- `bundle_id` (String) The Identity Name.
- `licensing_type` (Attributes) The supported License Type. / Contains properties for iOS Volume-Purchased Program (Vpp) Licensing Type. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-vpplicensingtype?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_vpp--licensing_type))
- `release_date_time` (String) The VPP application release date and time.
- `revoke_license_action_results` (Attributes Set) Results of revoke license actions on this app. / Defines results for actions on iOS Vpp Apps, contains inherited properties for ActionResult. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosvppapprevokelicensesactionresult?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios_vpp--revoke_license_action_results))
- `total_license_count` (Number) The total number of VPP licenses.
- `used_license_count` (Number) The number of VPP licenses in use.
- `vpp_token_account_type` (String) The type of volume purchase program which the given Apple Volume Purchase Program Token is associated with. / Possible types of an Apple Volume Purchase Program token; possible values are: `business` (Apple Volume Purchase Program token associated with an business program.), `education` (Apple Volume Purchase Program token associated with an education program.)
- `vpp_token_apple_id` (String) The Apple Id associated with the given Apple Volume Purchase Program Token.
- `vpp_token_id` (String) Identifier of the VPP token associated with this app.
- `vpp_token_organization_name` (String) The organization associated with the Apple Volume Purchase Program Token

<a id="nestedatt--ios_vpp--applicable_device_type"></a>
### Nested Schema for `ios_vpp.applicable_device_type`

Optional:

- `ipad` (Boolean) Whether the app should run on iPads. The _provider_ default value is `true`.
- `iphone_and_ipod` (Boolean) Whether the app should run on iPhones and iPods. The _provider_ default value is `true`.


<a id="nestedatt--ios_vpp--assigned_licenses"></a>
### Nested Schema for `ios_vpp.assigned_licenses`

Read-Only:

- `id` (String) Key of the entity. This property is read-only.
- `user_email_address` (String) The user email address.
- `user_id` (String) The user ID.
- `user_name` (String) The user name.
- `user_principal_name` (String) The user principal name.


<a id="nestedatt--ios_vpp--licensing_type"></a>
### Nested Schema for `ios_vpp.licensing_type`

Optional:

- `support_device_licensing` (Boolean) Whether the program supports the device licensing type. The _provider_ default value is `false`.
- `support_user_licensing` (Boolean) Whether the program supports the user licensing type. The _provider_ default value is `false`.
- `supports_device_licensing` (Boolean) Whether the program supports the device licensing type. The _provider_ default value is `false`.
- `supports_user_licensing` (Boolean) Whether the program supports the user licensing type. The _provider_ default value is `false`.


<a id="nestedatt--ios_vpp--revoke_license_action_results"></a>
### Nested Schema for `ios_vpp.revoke_license_action_results`

Read-Only:

- `action_failure_reason` (String) The reason for the revoke licenses action failure. / Possible types of reasons for an Apple Volume Purchase Program token action failure; possible values are: `none` (None.), `appleFailure` (There was an error on Apple's service.), `internalError` (There was an internal error.), `expiredVppToken` (There was an error because the Apple Volume Purchase Program token was expired.), `expiredApplePushNotificationCertificate` (There was an error because the Apple Volume Purchase Program Push Notification certificate expired.)
- `action_name` (String) Action name
- `action_state` (String) State of the action. / State of the action on the device; possible values are: `none` (Not a valid action state), `pending` (Action is pending), `canceled` (Action has been cancelled.), `active` (Action is active.), `done` (Action completed without errors.), `failed` (Action failed), `notSupported` (Action is not supported.)
- `failed_licenses_count` (Number) A count of the number of licenses for which revoke failed.
- `last_updated_date_time` (String) Time the action state was last updated
- `managed_device_id` (String) DeviceId associated with the action.
- `start_date_time` (String) Time the action was initiated
- `total_licenses_count` (Number) A count of the number of licenses for which revoke was attempted.
- `user_id` (String) UserId associated with the action.



<a id="nestedatt--macos_lob"></a>
### Nested Schema for `macos_lob`

Read-Only:

- `build_number` (String) The build number of the package. This should match the package CFBundleShortVersionString of the .pkg file.
- `bundle_id` (String) The primary bundleId of the package.
- `child_apps` (Attributes Set) List of ComplexType macOSLobChildApp objects. Represents the apps expected to be installed by the package. / Contains properties of a macOS .app in the package / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macoslobchildapp?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_lob--child_apps))
- `committed_content_version` (String) The internal committed content version.
- `file_name` (String) The name of the main Lob application file.
- `ignore_version_detection` (Boolean) When TRUE, indicates that the app's version will NOT be used to detect if the app is installed on a device. When FALSE, indicates that the app's version will be used to detect if the app is installed on a device. Set this to true for apps that use a self update feature. The default value is FALSE.
- `install_as_managed` (Boolean) When TRUE, indicates that the app will be installed as managed (requires macOS 11.0 and other managed package restrictions). When FALSE, indicates that the app will be installed as unmanaged. The default value is FALSE.
- `md5_hash` (Set of String) The MD5 hash codes. This is empty if the package was uploaded directly. If the Intune App Wrapping Tool is used to create a .intunemac, this value can be found inside the Detection.xml file.
- `md5_hash_chunk_size` (Number) The chunk size for MD5 hash. This is '0' or empty if the package was uploaded directly. If the Intune App Wrapping Tool is used to create a .intunemac, this value can be found inside the Detection.xml file.
- `minimum_supported_operating_system` (Attributes) ComplexType macOSMinimumOperatingSystem that indicates the minimum operating system applicable for the application. / The minimum operating system required for a macOS app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosminimumoperatingsystem?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_lob--minimum_supported_operating_system))
- `size` (Number) The total size, including all uploaded files. This property is read-only.
- `version_number` (String) The version number of the package. This should match the package CFBundleVersion in the packageinfo file.

<a id="nestedatt--macos_lob--child_apps"></a>
### Nested Schema for `macos_lob.child_apps`

Read-Only:

- `build_number` (String) The build number of the app.
- `bundle_id` (String) The bundleId of the app.
- `version_number` (String) The version number of the app.


<a id="nestedatt--macos_lob--minimum_supported_operating_system"></a>
### Nested Schema for `macos_lob.minimum_supported_operating_system`

Optional:

- `v10_10` (Boolean) When TRUE, indicates OS X 10.10 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_11` (Boolean) When TRUE, indicates OS X 10.11 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_12` (Boolean) When TRUE, indicates macOS 10.12 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_13` (Boolean) When TRUE, indicates macOS 10.13 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_14` (Boolean) When TRUE, indicates macOS 10.14 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_15` (Boolean) When TRUE, indicates macOS 10.15 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_7` (Boolean) When TRUE, indicates Mac OS X 10.7 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_8` (Boolean) When TRUE, indicates OS X 10.8 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v10_9` (Boolean) When TRUE, indicates OS X 10.9 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v11_0` (Boolean) When TRUE, indicates macOS 11.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v12_0` (Boolean) When TRUE, indicates macOS 12.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v13_0` (Boolean) When TRUE, indicates macOS 13.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.
- `v14_0` (Boolean) When TRUE, indicates macOS 14.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE. The _provider_ default value is `false`.



<a id="nestedatt--macos_vpp"></a>
### Nested Schema for `macos_vpp`

Read-Only:

- `app_store_url` (String) The store URL.
- `assigned_licenses` (Attributes Set) The licenses assigned to this app. / MacOS Volume Purchase Program license assignment. This class does not support Create, Delete, or Update. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosvppappassignedlicense?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_vpp--assigned_licenses))
- `bundle_id` (String) The Identity Name.
- `licensing_type` (Attributes) The supported License Type. / Contains properties for iOS Volume-Purchased Program (Vpp) Licensing Type. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-vpplicensingtype?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_vpp--licensing_type))
- `release_date_time` (String) The VPP application release date and time.
- `revoke_license_action_results` (Attributes Set) Results of revoke license actions on this app. / Defines results for actions on MacOS Vpp Apps, contains inherited properties for ActionResult. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-macosvppapprevokelicensesactionresult?view=graph-rest-beta (see [below for nested schema](#nestedatt--macos_vpp--revoke_license_action_results))
- `total_license_count` (Number) The total number of VPP licenses.
- `used_license_count` (Number) The number of VPP licenses in use.
- `vpp_token_account_type` (String) The type of volume purchase program which the given Apple Volume Purchase Program Token is associated with. / Possible types of an Apple Volume Purchase Program token; possible values are: `business` (Apple Volume Purchase Program token associated with an business program.), `education` (Apple Volume Purchase Program token associated with an education program.)
- `vpp_token_apple_id` (String) The Apple Id associated with the given Apple Volume Purchase Program Token.
- `vpp_token_id` (String) Identifier of the VPP token associated with this app.
- `vpp_token_organization_name` (String) The organization associated with the Apple Volume Purchase Program Token

<a id="nestedatt--macos_vpp--assigned_licenses"></a>
### Nested Schema for `macos_vpp.assigned_licenses`

Read-Only:

- `id` (String) Key of the entity. This property is read-only.
- `user_email_address` (String) The user email address.
- `user_id` (String) The user ID.
- `user_name` (String) The user name.
- `user_principal_name` (String) The user principal name.


<a id="nestedatt--macos_vpp--licensing_type"></a>
### Nested Schema for `macos_vpp.licensing_type`

Optional:

- `support_device_licensing` (Boolean) Whether the program supports the device licensing type. The _provider_ default value is `false`.
- `support_user_licensing` (Boolean) Whether the program supports the user licensing type. The _provider_ default value is `false`.
- `supports_device_licensing` (Boolean) Whether the program supports the device licensing type. The _provider_ default value is `false`.
- `supports_user_licensing` (Boolean) Whether the program supports the user licensing type. The _provider_ default value is `false`.


<a id="nestedatt--macos_vpp--revoke_license_action_results"></a>
### Nested Schema for `macos_vpp.revoke_license_action_results`

Read-Only:

- `action_failure_reason` (String) The reason for the revoke licenses action failure. / Possible types of reasons for an Apple Volume Purchase Program token action failure; possible values are: `none` (None.), `appleFailure` (There was an error on Apple's service.), `internalError` (There was an internal error.), `expiredVppToken` (There was an error because the Apple Volume Purchase Program token was expired.), `expiredApplePushNotificationCertificate` (There was an error because the Apple Volume Purchase Program Push Notification certificate expired.)
- `action_name` (String) Action name
- `action_state` (String) State of the action. / State of the action on the device; possible values are: `none` (Not a valid action state), `pending` (Action is pending), `canceled` (Action has been cancelled.), `active` (Action is active.), `done` (Action completed without errors.), `failed` (Action failed), `notSupported` (Action is not supported.)
- `failed_licenses_count` (Number) A count of the number of licenses for which revoke failed.
- `last_updated_date_time` (String) Time the action state was last updated
- `managed_device_id` (String) DeviceId associated with the action.
- `start_date_time` (String) Time the action was initiated
- `total_licenses_count` (Number) A count of the number of licenses for which revoke was attempted.
- `user_id` (String) UserId associated with the action.



<a id="nestedatt--managed_android_lob"></a>
### Nested Schema for `managed_android_lob`

Read-Only:

- `app_availability` (String) The Application's availability. This property is read-only. / A managed (MAM) application's availability; possible values are: `global` (A globally available app to all tenants.), `lineOfBusiness` (A line of business apps private to an organization.)
- `committed_content_version` (String) The internal committed content version.
- `file_name` (String) The name of the main Lob application file.
- `minimum_supported_operating_system` (Attributes) The value for the minimum applicable operating system. / Contains properties for the minimum operating system required for an Android mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidminimumoperatingsystem?view=graph-rest-beta (see [below for nested schema](#nestedatt--managed_android_lob--minimum_supported_operating_system))
- `package_id` (String) The package identifier.
- `size` (Number) The total size, including all uploaded files. This property is read-only.
- `targeted_platforms` (String) The platforms to which the application can be targeted. If not specified, will defauilt to Android Device Administrator. / Specifies which platform(s) can be targeted for a given Android LOB application or Managed Android LOB application; possible values are: `androidDeviceAdministrator` (Indicates the Android targeted platform is Android Device Administrator.), `androidOpenSourceProject` (Indicates the Android targeted platform is Android Open Source Project.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.)
- `version` (String) The Application's version.
- `version_code` (String) The version code of managed Android Line of Business (LoB) app.
- `version_name` (String) The version name of managed Android Line of Business (LoB) app.

<a id="nestedatt--managed_android_lob--minimum_supported_operating_system"></a>
### Nested Schema for `managed_android_lob.minimum_supported_operating_system`

Optional:

- `v10_0` (Boolean) When TRUE, only Version 10.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v11_0` (Boolean) When TRUE, only Version 11.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v12_0` (Boolean) When TRUE, only Version 12.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v13_0` (Boolean) When TRUE, only Version 13.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v14_0` (Boolean) When TRUE, only Version 14.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v15_0` (Boolean) When TRUE, only Version 15.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_0` (Boolean) When TRUE, only Version 4.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_0_3` (Boolean) When TRUE, only Version 4.0.3 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_1` (Boolean) When TRUE, only Version 4.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_2` (Boolean) When TRUE, only Version 4.2 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_3` (Boolean) When TRUE, only Version 4.3 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v4_4` (Boolean) When TRUE, only Version 4.4 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v5_0` (Boolean) When TRUE, only Version 5.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v5_1` (Boolean) When TRUE, only Version 5.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v6_0` (Boolean) When TRUE, only Version 6.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v7_0` (Boolean) When TRUE, only Version 7.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v7_1` (Boolean) When TRUE, only Version 7.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v8_0` (Boolean) When TRUE, only Version 8.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v8_1` (Boolean) When TRUE, only Version 8.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v9_0` (Boolean) When TRUE, only Version 9.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.



<a id="nestedatt--managed_ios_lob"></a>
### Nested Schema for `managed_ios_lob`

Read-Only:

- `app_availability` (String) The Application's availability. This property is read-only. / A managed (MAM) application's availability; possible values are: `global` (A globally available app to all tenants.), `lineOfBusiness` (A line of business apps private to an organization.)
- `applicable_device_type` (Attributes) The iOS architecture for which this app can run on. / Contains properties of the possible iOS device types the mobile app can run on. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosdevicetype?view=graph-rest-beta (see [below for nested schema](#nestedatt--managed_ios_lob--applicable_device_type))
- `build_number` (String) The build number of managed iOS Line of Business (LoB) app.
- `bundle_id` (String) The Identity Name.
- `committed_content_version` (String) The internal committed content version.
- `expiration_date_time` (String) The expiration time.
- `file_name` (String) The name of the main Lob application file.
- `identity_version` (String) The identity version.
- `minimum_supported_operating_system` (Attributes) The value for the minimum applicable operating system. / Contains properties of the minimum operating system required for an iOS mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosminimumoperatingsystem?view=graph-rest-beta (see [below for nested schema](#nestedatt--managed_ios_lob--minimum_supported_operating_system))
- `size` (Number) The total size, including all uploaded files. This property is read-only.
- `version` (String) The Application's version.
- `version_number` (String) The version number of managed iOS Line of Business (LoB) app.

<a id="nestedatt--managed_ios_lob--applicable_device_type"></a>
### Nested Schema for `managed_ios_lob.applicable_device_type`

Optional:

- `ipad` (Boolean) Whether the app should run on iPads. The _provider_ default value is `true`.
- `iphone_and_ipod` (Boolean) Whether the app should run on iPhones and iPods. The _provider_ default value is `true`.


<a id="nestedatt--managed_ios_lob--minimum_supported_operating_system"></a>
### Nested Schema for `managed_ios_lob.minimum_supported_operating_system`

Optional:

- `v10_0` (Boolean) When TRUE, only Version 10.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v11_0` (Boolean) When TRUE, only Version 11.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v12_0` (Boolean) When TRUE, only Version 12.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v13_0` (Boolean) When TRUE, only Version 13.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v14_0` (Boolean) When TRUE, only Version 14.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v15_0` (Boolean) When TRUE, only Version 15.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v16_0` (Boolean) When TRUE, only Version 16.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v17_0` (Boolean) When TRUE, only Version 17.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v8_0` (Boolean) When TRUE, only Version 8.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.
- `v9_0` (Boolean) When TRUE, only Version 9.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE. The _provider_ default value is `false`.