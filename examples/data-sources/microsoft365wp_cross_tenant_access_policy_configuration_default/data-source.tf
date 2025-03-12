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


data "microsoft365wp_cross_tenant_access_policy_configuration_default" "singleton" {
}

output "microsoft365wp_cross_tenant_access_policy_configuration_default" {
  value = data.microsoft365wp_cross_tenant_access_policy_configuration_default.singleton
}
