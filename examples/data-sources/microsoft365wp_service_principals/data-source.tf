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


data "microsoft365wp_service_principals" "all" {
}

output "microsoft365wp_service_principals_all" {
  value = { for x in data.microsoft365wp_service_principals.all.service_principals : x.id => x }
}

data "microsoft365wp_service_principals" "enabled_and_managed_identities" {
  account_enabled        = true
  service_principal_type = "ManagedIdentity"
}

output "microsoft365wp_service_principals_enabled_and_managed_identities" {
  value = { for x in data.microsoft365wp_service_principals.enabled_and_managed_identities.service_principals : x.id => x }
}
