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


resource "microsoft365wp_connector_group" "test2" {
  name   = "TF Test 2"
  region = "eur"
}

data "microsoft365wp_connector" "gsaconn1" {
  machine_name = "gsaconn1"
}

resource "microsoft365wp_connector_group_member_connector" "test" {
  connector_group_id = microsoft365wp_connector_group.test2.id
  id                 = data.microsoft365wp_connector.gsaconn1.id
}

output "microsoft365wp_connector_group_member_connector" {
  value = microsoft365wp_connector_group_member_connector.test
}
