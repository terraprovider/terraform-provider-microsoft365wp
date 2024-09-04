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


data "microsoft365wp_mobile_app" "one" {
  id = "255c9708-2cab-46b7-b668-dd83069cbad3"
}

output "microsoft365wp_mobile_app" {
  value = data.microsoft365wp_mobile_app.one
}
