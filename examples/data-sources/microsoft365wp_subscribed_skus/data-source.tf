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


data "microsoft365wp_subscribed_skus" "all" {
}

output "microsoft365wp_subscribed_skus" {
  value = { for x in data.microsoft365wp_subscribed_skus.all.subscribed_skus : x.id => x }
}
