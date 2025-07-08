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


data "microsoft365wp_device_compliance_script" "one" {
  id = "3971eaf9-afba-4ac0-b524-817f4c9c15fa"
}

output "microsoft365wp_device_compliance_script" {
  value = data.microsoft365wp_device_compliance_script.one
}
