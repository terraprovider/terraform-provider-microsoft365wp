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


data "microsoft365wp_mobile_app_category" "one" {
  id = "ed899483-3019-425e-a470-28e901b9790e" # "Productivity"
}

output "microsoft365wp_mobile_app_category" {
  value = data.microsoft365wp_mobile_app_category.one
}
