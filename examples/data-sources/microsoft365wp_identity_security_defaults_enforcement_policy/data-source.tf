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


data "microsoft365wp_identity_security_defaults_enforcement_policy" "singleton" {
}

output "microsoft365wp_identity_security_defaults_enforcement_policy" {
  value = data.microsoft365wp_identity_security_defaults_enforcement_policy.singleton
}
