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


resource "microsoft365wp_connector_group" "test3" {
  name   = "TF Test 3"
  region = "eur"
}

data "microsoft365wp_application" "one" {
  app_id = "990d9e03-b1a0-4136-ad17-3cc3a9b829bf"
}

resource "microsoft365wp_connector_group_application" "test" {
  connector_group_id = microsoft365wp_connector_group.test3.id
  id                 = data.microsoft365wp_application.one.id
}

output "microsoft365wp_connector_group_application" {
  value = microsoft365wp_connector_group_application.test
}
