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


data "microsoft365wp_authentication_context_class_references" "all" {
  # exclude_ids = ["c25"]
}

output "microsoft365wp_authentication_context_class_references" {
  value = { for x in data.microsoft365wp_authentication_context_class_references.all.authentication_context_class_references : x.id => x }
}
