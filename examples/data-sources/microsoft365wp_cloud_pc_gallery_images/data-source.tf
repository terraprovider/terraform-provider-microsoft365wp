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


data "microsoft365wp_cloud_pc_gallery_images" "all" {
}

output "microsoft365wp_cloud_pc_gallery_images_all" {
  value = { for x in data.microsoft365wp_cloud_pc_gallery_images.all.cloud_pc_gallery_images : x.id => x }
}

data "microsoft365wp_cloud_pc_gallery_images" "win11_m365_supported" {
  # publisher_name/publisherName must be compared case-insensitive but this specific query would produce the same result
  # if only filtering by offer_name, sku_name and status (and ommitting the odata_filter clause completely)
  odata_filter = "tolower(publisherName) eq 'microsoftwindowsdesktop'"
  offer_name   = "windows-ent-cpc"
  sku_name     = "win11-*-m365"
  status       = "supported"

  odata_orderby = "startDate desc"
}

output "microsoft365wp_cloud_pc_gallery_images_win11_m365_supported" {
  value = data.microsoft365wp_cloud_pc_gallery_images.win11_m365_supported.cloud_pc_gallery_images
}
