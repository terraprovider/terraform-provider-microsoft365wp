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


resource "microsoft365wp_group_assigned_license" "test1" {
  group_id = "298fded6-b252-4166-a473-f405e935f58d"

  sku_id = "00ed1723-1992-4384-b7ce-1c3bf01eedc7"
}

resource "microsoft365wp_group_assigned_license" "test2" {
  group_id = "298fded6-b252-4166-a473-f405e935f58d"

  sku_id = "05e9a617-0261-4cee-bb44-138d3ef5d965"
}
