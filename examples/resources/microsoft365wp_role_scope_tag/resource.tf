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


resource "microsoft365wp_role_scope_tag" "test" {
  display_name = "TF Test"
  description  = "TF Test scope tag"

  # optional: omit to leave assignments untouched (e.g. when managed elsewhere), set to [] to clear them
  assignments = [
    { target = { group = { group_id = "62e39046-aad3-4423-98e0-b486e3538aff" } } },
  ]
}
