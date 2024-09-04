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


data "microsoft365wp_notification_message_template" "one" {
  id = "23a1cb09-1b64-42e7-8155-5be7a3d55fc8"
}

output "microsoft365wp_notification_message_template" {
  value = data.microsoft365wp_notification_message_template.one
}
