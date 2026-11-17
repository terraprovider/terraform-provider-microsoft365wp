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


data "microsoft365wp_directory_setting" "one" {
  id = "95e0fc18-6b86-43f6-b53b-364ec1099f5f" # id of the setting (not of the template)
}

output "microsoft365wp_directory_setting" {
  value = data.microsoft365wp_directory_setting.one
}
