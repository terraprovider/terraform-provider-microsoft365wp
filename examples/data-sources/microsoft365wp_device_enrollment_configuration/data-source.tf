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


data "microsoft365wp_device_enrollment_configuration" "one" {
  id = "d186b27b-2a2f-44f0-920e-6fddd22ba192_Windows10EnrollmentCompletionPageConfiguration"
}

output "microsoft365wp_device_enrollment_configuration" {
  value = data.microsoft365wp_device_enrollment_configuration.one
}
