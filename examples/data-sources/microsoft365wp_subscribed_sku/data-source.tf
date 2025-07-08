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


data "microsoft365wp_subscribed_sku" "one" {
  id = "810017af-c1d4-4a3b-9e18-d693c12b0df1_05e9a617-0261-4cee-bb44-138d3ef5d965"
}

output "microsoft365wp_subscribed_sku" {
  value = data.microsoft365wp_subscribed_sku.one
}
