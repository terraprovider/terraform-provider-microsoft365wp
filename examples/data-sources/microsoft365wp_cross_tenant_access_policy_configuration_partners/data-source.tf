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


data "microsoft365wp_cross_tenant_access_policy_configuration_partners" "all" {
}

output "microsoft365wp_cross_tenant_access_policy_configuration_partners" {
  value = { for x in data.microsoft365wp_cross_tenant_access_policy_configuration_partners.all.cross_tenant_access_policy_configuration_partners : x.tenant_id => x }
}
