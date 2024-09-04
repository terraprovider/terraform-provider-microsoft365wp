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
  id = "3bc67ec2-2d41-4c12-a155-dee78e3f0174"
}

output "microsoft365wp_cloud_pc_provisioning_policy" {
  value = data.microsoft365wp_cloud_pc_provisioning_policy.one
}
