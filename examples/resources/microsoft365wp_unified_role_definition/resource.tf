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


resource "microsoft365wp_unified_role_definition" "test" {
  display_name = "TF Test"
  description  = "A custom role created by TF"

  is_enabled = true

  role_permissions = [{
    allowed_resource_actions = [
      "microsoft.directory/applications/basic/update",
    ]
  }]
}
