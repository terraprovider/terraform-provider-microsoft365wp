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


data "microsoft365wp_service_principal" "one" {
  app_id = "12472e6c-34d1-4db9-9db0-f6649b14f5ba"
}

output "microsoft365wp_service_principal" {
  value = data.microsoft365wp_service_principal.one
}
