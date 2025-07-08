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


data "microsoft365wp_attribute_sets" "all" {
}

output "microsoft365wp_attribute_sets" {
  value = { for x in data.microsoft365wp_attribute_sets.all.attribute_sets : x.id => x }
}
