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


data "microsoft365wp_cloud_pc_gallery_image" "one" {
  # Win 11 with M365 latest
  publisher_name = "microsoftwindowsdesktop"
  offer_name     = "windows-ent-cpc"
  sku_name       = "win11-*-m365"
  status         = "supported"
  odata_orderby  = "startDate desc"
  odata_top      = 1
}

output "microsoft365wp_cloud_pc_gallery_image" {
  value = data.microsoft365wp_cloud_pc_gallery_image.one
}
