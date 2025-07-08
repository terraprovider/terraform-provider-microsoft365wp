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


resource "microsoft365wp_authentication_methods_policy" "singleton" {

  authentication_method_configurations = [

    {
      id    = "Fido2"
      state = "enabled"
      fido2 = {}
    },

    {
      id                      = "MicrosoftAuthenticator"
      microsoft_authenticator = {}
      state                   = "enabled"
    },

    {
      id    = "Sms"
      sms   = {}
      state = "enabled"
    },

    {
      id          = "QRCodePin"
      qr_code_pin = {}
      state       = "disabled"
    },

    {
      id                    = "TemporaryAccessPass"
      temporary_access_pass = {}
      state                 = "enabled"
    },

    {
      id            = "HardwareOath"
      hardware_oath = {}
      state         = "disabled"
    },

    {
      id            = "SoftwareOath"
      software_oath = {}
      state         = "disabled"
    },

    {
      id    = "Voice"
      voice = {}
      state = "enabled"
    },

    {
      id    = "Email"
      email = {}
      state = "enabled"
    },

    {
      id               = "X509Certificate"
      x509_certificate = {}
      state            = "disabled"
    },

  ]

  policy_migration_state = "migrationComplete"
}
