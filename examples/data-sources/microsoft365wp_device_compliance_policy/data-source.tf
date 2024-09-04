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


data "microsoft365wp_device_compliance_policy" "one" {
  id = "8b424a97-04d4-40c9-bdb3-b11894523c73"
}

output "microsoft365wp_device_compliance_policy" {
  value = data.microsoft365wp_device_compliance_policy.one
}
