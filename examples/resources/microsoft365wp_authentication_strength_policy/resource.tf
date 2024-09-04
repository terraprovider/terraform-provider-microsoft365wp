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


resource "microsoft365wp_authentication_strength_policy" "test" {
  display_name = "TF Test AppFd2"
  allowed_combinations = [
    "fido2",
    "deviceBasedPush",
    "password,microsoftAuthenticatorPush",
  ]
}
