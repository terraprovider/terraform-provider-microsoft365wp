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


resource "microsoft365wp_device_configuration_custom" "win10" {
  display_name = "TF Test Windows 10 Custom"
  windows10 = {
    oma_settings = [
      {
        display_name = "OMA 1 Base64"
        oma_uri      = "./User/Vendor/C4A8/Base64"
        base64       = { value_base64 = base64encode("asdfa") }
      },
      {
        display_name = "OMA 2 Boolean"
        oma_uri      = "./User/Vendor/C4A8/Boolean"
        boolean      = { value = false }
      },
      {
        display_name = "OMA 3 DateTime"
        oma_uri      = "./User/Vendor/C4A8/DateTime"
        date_time    = { value = "2023-01-01T00:00:01Z" }
      },
      {
        display_name   = "OMA 4 FloatingPoint"
        oma_uri        = "./User/Vendor/C4A8/FloatingPoint"
        floating_point = { value = 42.1 }
      },
      {
        display_name = "OMA 5 Integer"
        oma_uri      = "./User/Vendor/C4A8/Integer"
        integer      = { value = 42 }
      },
      {
        display_name = "OMA 6 String"
        oma_uri      = "./User/Vendor/C4A8/String"
        string       = { value = "asdf" }
      },
      {
        display_name = "OMA 7 StringXml"
        oma_uri      = "./User/Vendor/C4A8/StringXml"
        string_xml   = { value = "<?xml version=\"1.0\" encoding=\"utf-8\"?><root>a</root>" }
      },
    ]
  }
}
