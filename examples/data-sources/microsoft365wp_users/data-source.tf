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


data "microsoft365wp_users" "all" {
}

output "microsoft365wp_users_all" {
  value = { for x in data.microsoft365wp_users.all.users : x.id => x }
}

data "microsoft365wp_users" "with_e3" {
  # "assignedLicenses/any()" (without sku) does not seem to work for users, but it does work with a specific sku:
  odata_filter = "assignedLicenses/any(x:x/skuId eq 05e9a617-0261-4cee-bb44-138d3ef5d965)"
}

output "microsoft365wp_users_with_e3" {
  value = { for x in data.microsoft365wp_users.with_e3.users : x.id => x }
}
