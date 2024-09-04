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


data "microsoft365wp_android_managed_app_protections" "all" {
  exclude_ids = ["01234567-89ab-cdef-0123-456789abcdee", "01234567-89ab-cdef-0123-456789abcdef"]
}

output "microsoft365wp_android_managed_app_protections" {
  value = { for x in data.microsoft365wp_android_managed_app_protections.all.android_managed_app_protections : x.id => x }
}
