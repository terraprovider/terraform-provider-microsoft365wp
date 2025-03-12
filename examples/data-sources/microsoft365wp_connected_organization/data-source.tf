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


data "microsoft365wp_connected_organization" "one" {
  id = "99d8f723-a23e-447d-9a3b-8e59c00211ba"
}

output "microsoft365wp_connected_organization" {
  value = data.microsoft365wp_connected_organization.one
}
