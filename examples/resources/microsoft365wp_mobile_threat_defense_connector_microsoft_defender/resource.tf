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


resource "microsoft365wp_mobile_threat_defense_connector_microsoft_defender" "singleton" {
  android_enabled = true
  ios_enabled     = true
  windows_enabled = true

  partner_unresponsiveness_threshold_in_days = 14
}
