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


data "microsoft365wp_targeted_managed_app_configuration" "one" {
  id = "A_b04cdf30-859a-44e2-b4de-c81ac33eab6a"
}

output "microsoft365wp_targeted_managed_app_configuration" {
  value = data.microsoft365wp_targeted_managed_app_configuration.one
}
