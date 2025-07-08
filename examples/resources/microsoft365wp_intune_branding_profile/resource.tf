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


resource "microsoft365wp_intune_branding_profile" "default" {
  profile_name = "###TfWorkplaceDefault###"
  display_name = "GKCorellia Inc."
  privacy_url  = "https://azure.microsoft.com/explore/trusted-cloud/privacy"

  landing_page_customized_image = {
    type         = "image/png"
    value_base64 = filebase64("terraprovider-logo.png")
  }

  lifecycle {
    ignore_changes = [profile_name]
  }
}

resource "microsoft365wp_intune_branding_profile" "test" {
  profile_name = "TF Test"
  display_name = "Acme Corporation"
  privacy_url  = "https://azure.microsoft.com/explore/trusted-cloud/privacy"

  ### There have been **nonspecific errors reported when trying to _assign_ with application permissions** (even
  ### though the MS docs state that it should be possible)! Empty `assignments` or using delegated permissions both
  ### seemed to work in our tests, though.
  assignments = [
    { target = { group = { group_id = "62e39046-aad3-4423-98e0-b486e3538aff" } } },
  ]
}
