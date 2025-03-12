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


data "microsoft365wp_cross_tenant_access_policy_configuration_partner" "one" {
  tenant_id = "eb46eb45-4f35-4332-987d-deda508c6721"
}

output "microsoft365wp_cross_tenant_access_policy_configuration_partner" {
  value = data.microsoft365wp_cross_tenant_access_policy_configuration_partner.one
}
