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


resource "microsoft365wp_cross_tenant_access_policy_configuration_default" "singleton" {
  b2b_direct_connect_outbound = {
    applications = {
      access_type = "allowed"
    }
    users_and_groups = {
      access_type = "allowed"
    }
  }
}
