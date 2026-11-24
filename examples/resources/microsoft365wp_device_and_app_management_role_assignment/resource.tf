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


resource "microsoft365wp_device_and_app_management_role_assignment" "test" {
  display_name       = "TF Test"
  description        = "TF Test Intune role assignment"
  role_definition_id = microsoft365wp_device_and_app_management_role_definition.test.id

  members         = ["62e39046-aad3-4423-98e0-b486e3538aff"] # admin group(s)
  scope_type      = "resourceScope"
  scope_members   = ["62e39046-aad3-4423-98e0-b486e3538aff"] # group(s) the admins may manage

  role_scope_tag_ids = ["0"]
}

resource "microsoft365wp_device_and_app_management_role_definition" "test" {
  display_name = "TF Test"

  role_permissions = [{
    resource_actions = [{
      allowed_resource_actions = ["Microsoft.Intune_DeviceConfigurations_Read"]
    }]
  }]
}
