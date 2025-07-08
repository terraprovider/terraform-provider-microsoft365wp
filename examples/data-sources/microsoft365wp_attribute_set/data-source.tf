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


data "microsoft365wp_attribute_set" "one" {
  id = "TestKeepForTfDatasource"
}

output "microsoft365wp_attribute_set" {
  value = data.microsoft365wp_attribute_set.one
}
