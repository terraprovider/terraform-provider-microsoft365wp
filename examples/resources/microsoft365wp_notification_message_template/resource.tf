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


resource "microsoft365wp_notification_message_template" "test" {
  display_name = "TF Test 1"

  localized_notification_messages = [
    {
      locale           = "en-us"
      subject          = "Test en-us"
      message_template = "This is a test message."
      is_default       = true
    },
    {
      locale           = "de-de"
      subject          = "Test de-de"
      message_template = "Dies ist eine Testnachricht."
    },
    {
      locale           = "es-es"
      subject          = "Test es-es"
      message_template = "Este es un mensaje de prueba."
    }
  ]
}
