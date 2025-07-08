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


resource "microsoft365wp_device_compliance_script" "test" {
  display_name = "TF Test"

  detection_script_content = <<-EOT
    $hash = @{ ThisMachineIsCompliant = $true }
    return $hash | ConvertTo-Json -Compress
  EOT

  assignments = [
    { target = { group = { group_id = "62e39046-aad3-4423-98e0-b486e3538aff" } } },
  ]
}
