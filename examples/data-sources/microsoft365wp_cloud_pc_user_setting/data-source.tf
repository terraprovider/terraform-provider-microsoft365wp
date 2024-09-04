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
  id = "10a90ae9-4d43-45ca-9ba1-14047121178c"
}

output "microsoft365wp_cloud_pc_user_setting" {
  value = data.microsoft365wp_cloud_pc_user_setting.one
}
