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


data "microsoft365wp_ios_managed_app_protection" "one" {
  id = "T_b928a553-c00b-4aa3-bb9e-98305b7ba3b6"
}

output "microsoft365wp_ios_managed_app_protection" {
  value = data.microsoft365wp_ios_managed_app_protection.one
}
