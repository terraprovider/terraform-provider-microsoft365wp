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

data "microsoft365wp_cloud_pc_gallery_image" "dedicated" {
  # Win11 with M365 latest
  publisher_name = "microsoftwindowsdesktop"
  offer_name     = "windows-ent-cpc"
  sku_name       = "win11-*-m365"
  status         = "supported"
  odata_orderby  = "startDate desc"
  odata_top      = 1
}

resource "microsoft365wp_cloud_pc_provisioning_policy" "test_dedicated" {
  display_name = "TF Test - dedicated"

  provisioning_type = "dedicated"

  image_id           = data.microsoft365wp_cloud_pc_gallery_image.dedicated.id
  image_display_name = data.microsoft365wp_cloud_pc_gallery_image.dedicated.display_name
  # for custom images, use microsoft365wp_cloud_pc_device_image to lookup details

  windows_setting = {
    locale = "en-US"
  }

  domain_join_configurations = [{
    region_group = "europeUnion"
    region_name  = "automatic"
  }]

  assignments = [
    for x in [
      "298fded6-b252-4166-a473-f405e935f58d",
      "62e39046-aad3-4423-98e0-b486e3538aff",
    ] :
    { target = { group = { group_id = x } } }
  ]
}


data "microsoft365wp_cloud_pc_gallery_image" "shared" {
  id = "microsoftwindowsdesktop_windows-ent-cpc_win11-23h2-ent-cpc-m365"
}

resource "microsoft365wp_cloud_pc_provisioning_policy" "test_shared" {
  display_name = "TF Test - shared"

  provisioning_type = "sharedByUser"

  image_id           = data.microsoft365wp_cloud_pc_gallery_image.shared.id
  image_display_name = data.microsoft365wp_cloud_pc_gallery_image.shared.display_name
  # for custom images, use microsoft365wp_cloud_pc_device_image to lookup details

  windows_setting = {
    locale = "en-US"
  }

  domain_join_configurations = [{
    region_group = "europeUnion"
    region_name  = "automatic"
  }]

  assignments = [
    for x in [
      "298fded6-b252-4166-a473-f405e935f58d",
      "62e39046-aad3-4423-98e0-b486e3538aff",
    ] :
    { target = { group = {
      group_id        = x
      service_plan_id = "51855c77-4d2e-4736-be67-6dca605f2b57" // CPC_S_2C_4GB_128GB / Windows 365 Shared Use 2 vCPU, 4 GB, 128 GB
    } } }
  ]
}
