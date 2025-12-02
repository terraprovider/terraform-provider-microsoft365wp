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


data "microsoft365wp_administrative_units" "all" {
  is_member_management_restricted = true
}

output "microsoft365wp_administrative_units" {
  value = { for x in data.microsoft365wp_administrative_units.all.administrative_units : x.id => x }
}
