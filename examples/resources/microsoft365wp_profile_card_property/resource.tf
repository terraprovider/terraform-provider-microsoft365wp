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


resource "microsoft365wp_profile_card_property" "test" {
  directory_property_name = "customAttribute1"

  annotations = [
    {
      display_name = "Cost Center"
      localizations = [
        { language_tag = "de-DE", display_name = "Kostenstelle" },
      ]
    },
  ]
}
