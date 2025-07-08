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


data "microsoft365wp_group_assigned_license" "one" {
  group_id = "298fded6-b252-4166-a473-f405e935f58d"

  sku_id = "00ed1723-1992-4384-b7ce-1c3bf01eedc7"
}

output "microsoft365wp_group_assigned_license" {
  value = data.microsoft365wp_group_assigned_license.one
}
