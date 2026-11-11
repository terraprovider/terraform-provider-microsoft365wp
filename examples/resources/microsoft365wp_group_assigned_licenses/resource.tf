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


resource "microsoft365wp_group_assigned_licenses" "test" {
  group_id = "298fded6-b252-4166-a473-f405e935f58d"

  # authoritative = true # remove any other licenses directly assigned to the group

  assigned_licenses = [
    {
      sku_id = "00ed1723-1992-4384-b7ce-1c3bf01eedc7"
    },
    {
      sku_id         = "05e9a617-0261-4cee-bb44-138d3ef5d965"
      disabled_plans = ["a23b959c-7ce8-4e57-9140-b90eb88a9e97"] # SWAY
    },
  ]
}
