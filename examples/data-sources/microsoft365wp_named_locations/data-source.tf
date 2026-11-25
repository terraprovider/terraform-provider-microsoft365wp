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


data "microsoft365wp_named_locations" "all" {
}

output "microsoft365wp_named_locations" {
  value = { for x in data.microsoft365wp_named_locations.all.named_locations : x.display_name => x }
}
