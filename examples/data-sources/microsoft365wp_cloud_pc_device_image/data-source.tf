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


data "microsoft365wp_cloud_pc_device_image" "does_not_exist" {
  id = "does_not_exist"
}

output "microsoft365wp_cloud_pc_device_image" {
  value = data.microsoft365wp_cloud_pc_device_image.does_not_exist
}
