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
  id = "T_0f190b8d-9e67-43b6-8d19-e567fbfda584"
}

output "microsoft365wp_android_managed_app_protection" {
  value = data.microsoft365wp_android_managed_app_protection.one
}
