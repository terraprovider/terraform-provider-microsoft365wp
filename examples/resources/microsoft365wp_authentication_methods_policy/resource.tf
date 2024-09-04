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
      state = "disabled"
      fido2 = {}
    },

    {
      id                      = "MicrosoftAuthenticator"
      microsoft_authenticator = {}
      state                   = "disabled"
    },

    {
      id    = "Sms"
      sms   = {}
      state = "disabled"
    },

    {
      id                    = "TemporaryAccessPass"
      temporary_access_pass = {}
      state                 = "disabled"
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
      state = "disabled"
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
}
