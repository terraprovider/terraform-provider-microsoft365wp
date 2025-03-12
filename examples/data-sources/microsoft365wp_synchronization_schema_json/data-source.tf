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


data "microsoft365wp_synchronization_schema_json" "one" {
  service_principal_id = "41dc044f-de44-430a-99a3-a1b15f6c8dd5"
  id                   = "Azure2Azure.810017afc1d44a3b9e18d693c12b0df1.12472e6c-34d1-4db9-9db0-f6649b14f5ba"
}

output "microsoft365wp_synchronization_schema_json" {
  value = data.microsoft365wp_synchronization_schema_json.one
}
