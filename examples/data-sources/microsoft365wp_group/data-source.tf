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


data "microsoft365wp_group" "one" {
  id = "298fded6-b252-4166-a473-f405e935f58d"
}

output "microsoft365wp_group" {
  value = data.microsoft365wp_group.one
}
