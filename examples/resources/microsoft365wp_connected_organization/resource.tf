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


resource "microsoft365wp_connected_organization" "test" {
  display_name = "TF Test 1"

  identity_sources = [
    {
      azure_active_directory_tenant = {
        display_name = "C4A8 Corellia Workplace Blueprint"
        tenant_id    = "eb46eb45-4f35-4332-987d-deda508c6721"
      }
    }
  ]

  internal_sponsors = [
    { id = "298fded6-b252-4166-a473-f405e935f58d" },
  ]
}
