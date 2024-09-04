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


data "microsoft365wp_cloud_pc_user_settings" "all" {
}

output "microsoft365wp_cloud_pc_user_settings" {
  value = { for x in data.microsoft365wp_cloud_pc_user_settings.all.cloud_pc_user_settings : x.id => x }
}
