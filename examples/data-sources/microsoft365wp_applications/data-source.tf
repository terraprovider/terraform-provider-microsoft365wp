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


data "microsoft365wp_applications" "all" {
}

output "microsoft365wp_applications" {
  value = { for x in data.microsoft365wp_applications.all.applications : x.id => x }
}

data "microsoft365wp_applications" "api" {
  odata_filter = "identifierUris/any(x:startswith(x,'api://'))"
}

output "microsoft365wp_applications_api" {
  value = { for x in data.microsoft365wp_applications.api.applications : x.id => x }
}
