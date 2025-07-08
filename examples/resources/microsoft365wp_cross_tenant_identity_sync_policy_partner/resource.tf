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


resource "microsoft365wp_cross_tenant_identity_sync_policy_partner" "test" {
  tenant_id = "eb46eb45-4f35-4332-987d-deda508c6721"

  display_name = "eb46eb45-4f35-4332-987d-deda508c6721:sync"
  user_sync_inbound = {
    is_sync_allowed = true
  }
}
