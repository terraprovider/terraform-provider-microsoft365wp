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


resource "microsoft365wp_device_and_app_management_role_definition" "test" {
  display_name = "TF Test"
  description  = "TF Test Intune role"

  role_permissions = [{
    resource_actions = [{
      allowed_resource_actions = [
        "Microsoft.Intune_DeviceConfigurations_Read",
        "Microsoft.Intune_DeviceCompliancePolices_Read",
      ]
    }]
  }]
}
