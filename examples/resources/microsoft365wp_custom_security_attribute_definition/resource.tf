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


resource "microsoft365wp_custom_security_attribute_definition" "test" {
  attribute_set = "TfProviderTest"
  description   = "A special attribute"
  name          = "Test10"
  type          = "String"

  use_pre_defined_values_only = true
  allowed_values = [
    { id = "one" },
    { id = "two" },
    { id = "three" },
    { id = "four" },
    { id = "four1", is_active = false },
  ]
}
