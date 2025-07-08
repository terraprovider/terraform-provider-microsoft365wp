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


data "microsoft365wp_user" "one" {
  id = "91ae5a52-67f6-4265-bbe5-b62268944675"
}

output "microsoft365wp_user" {
  value = data.microsoft365wp_user.one
}
