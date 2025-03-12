terraform {
  required_providers {
    microsoft365wp = {
      source = "terraprovider/microsoft365wp"
    }
  }
}


data "microsoft365wp_mobility_management_policies" "device" {
  policy_type = "device"
}

output "microsoft365wp_mobility_management_policies_device" {
  value = { for x in data.microsoft365wp_mobility_management_policies.device.mobility_management_policies : x.id => x }
}

data "microsoft365wp_mobility_management_policies" "app" {
  policy_type = "app"
}

output "microsoft365wp_mobility_management_policies_app" {
  value = { for x in data.microsoft365wp_mobility_management_policies.app.mobility_management_policies : x.id => x }
}
