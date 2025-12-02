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


data "microsoft365wp_administrative_unit" "one" {
  id = "9061a8e6-be66-4355-9d4d-3a039977f2a4"
}

output "microsoft365wp_administrative_unit" {
  value = data.microsoft365wp_administrative_unit.one
}
