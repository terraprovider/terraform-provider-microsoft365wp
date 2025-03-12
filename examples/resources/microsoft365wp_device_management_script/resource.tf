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


resource "microsoft365wp_device_management_script" "test" {
  display_name = "TF Test"

  file_name      = "helloWorld.ps1"
  script_content = <<EOT
Write-Host "Hello World!"
EOT

  assignments = [
    { target = { group = { group_id = "62e39046-aad3-4423-98e0-b486e3538aff" } } },
  ]
}
