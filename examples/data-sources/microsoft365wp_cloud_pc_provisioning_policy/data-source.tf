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


data "microsoft365wp_cloud_pc_provisioning_policy" "one" {
  id = "65253455-13c4-4d7f-80f9-ca35deb8439e"
}

output "microsoft365wp_cloud_pc_provisioning_policy" {
  value = data.microsoft365wp_cloud_pc_provisioning_policy.one
}
