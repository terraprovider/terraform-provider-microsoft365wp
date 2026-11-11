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


data "microsoft365wp_profile_card_properties" "all" {
}

output "microsoft365wp_profile_card_properties" {
  value = { for x in data.microsoft365wp_profile_card_properties.all.profile_card_properties : x.directory_property_name => x }
}
