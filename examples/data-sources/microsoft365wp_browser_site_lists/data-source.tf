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


data "microsoft365wp_browser_site_lists" "all" {
}

output "microsoft365wp_browser_site_lists" {
  value = { for x in data.microsoft365wp_browser_site_lists.all.browser_site_lists : x.id => x }
}
