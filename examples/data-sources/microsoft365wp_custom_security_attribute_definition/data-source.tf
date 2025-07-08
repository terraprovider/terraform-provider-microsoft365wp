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


data "microsoft365wp_custom_security_attribute_definition" "one" {
  id = "TestKeepForTfDatasource_TestKeepForTfDatasource1"
}

output "microsoft365wp_custom_security_attribute_definition" {
  value = data.microsoft365wp_custom_security_attribute_definition.one
}
