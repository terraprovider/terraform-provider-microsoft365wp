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


data "microsoft365wp_groups" "all" {
}

output "microsoft365wp_groups_all" {
  value = { for x in data.microsoft365wp_groups.all.groups : x.id => x }
}

data "microsoft365wp_groups" "with_licenses" {
  odata_filter = "assignedLicenses/any()"
}

output "microsoft365wp_groups_with_licenses" {
  value = { for x in data.microsoft365wp_groups.with_licenses.groups : x.id => x }
}
