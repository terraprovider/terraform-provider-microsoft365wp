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


data "microsoft365wp_android_managed_app_protection" "one" {
  id = "T_ba297831-e919-4f23-bd6a-31e613f08b52"
}

output "microsoft365wp_android_managed_app_protection" {
  value = data.microsoft365wp_android_managed_app_protection.one
}
