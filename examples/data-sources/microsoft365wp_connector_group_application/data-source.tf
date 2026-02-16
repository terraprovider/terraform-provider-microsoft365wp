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


data "microsoft365wp_application" "one" {
  app_id = "3af43c8d-63c9-43bf-9fdb-1c3d2c2a80ce"
}

data "microsoft365wp_connector_group_application" "one" {
  connector_group_id = "ebde4869-e643-44c9-8fda-f4d6ca5abcf3"
  id                 = data.microsoft365wp_application.one.id
}

output "microsoft365wp_connector_group_application" {
  value = data.microsoft365wp_connector_group_application.one
}
