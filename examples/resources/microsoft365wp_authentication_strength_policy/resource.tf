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


resource "microsoft365wp_authentication_combination_configuration" "test" {
  authentication_strength_policy_id = microsoft365wp_authentication_strength_policy.test.id

  applies_to_combinations = ["fido2"]
  fido2 = {
    allowed_aaguids = [
      "a25342c0-3cdc-4414-8e46-f4807fca511c",
      "d7781e5d-e353-46aa-afe2-3ca49f13332a",
    ]
  }
}
