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


data "microsoft365wp_connector_group" "one" {
  id = "ebde4869-e643-44c9-8fda-f4d6ca5abcf3"
}

output "microsoft365wp_connector_group" {
  value = data.microsoft365wp_connector_group.one
}
