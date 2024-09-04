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


data "microsoft365wp_conditional_access_policy" "one" {
  id = "6aaba8b8-539b-46f9-bd49-da9c53248272"
}

output "microsoft365wp_conditional_access_policy" {
  value = data.microsoft365wp_conditional_access_policy.one
}
