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


data "microsoft365wp_mobile_threat_defense_connector_microsoft_defender" "singleton" {
}

output "microsoft365wp_mobile_threat_defense_connector_microsoft_defender" {
  value = data.microsoft365wp_mobile_threat_defense_connector_microsoft_defender.singleton
}
