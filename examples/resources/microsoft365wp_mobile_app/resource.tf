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
