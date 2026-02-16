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


data "microsoft365wp_connector_group_member_connector" "one" {
  connector_group_id = "ebde4869-e643-44c9-8fda-f4d6ca5abcf3"
  id                 = "d0523f19-9ebf-4275-be47-3de67b67eb03"
}

output "microsoft365wp_connector_group_member_connector" {
  value = data.microsoft365wp_connector_group_member_connector.one
}
