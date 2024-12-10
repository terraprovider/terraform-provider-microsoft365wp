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


data "microsoft365wp_cloud_pc_user_setting" "one" {
  id = "226856e0-b728-44d8-8144-74f4e790a5e8"
}

output "microsoft365wp_cloud_pc_user_setting" {
  value = data.microsoft365wp_cloud_pc_user_setting.one
}
