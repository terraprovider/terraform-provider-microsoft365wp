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


data "microsoft365wp_authentication_strength_policy" "one" {
  id = "ccde5395-3d28-45f2-88dd-8ea78d538016"
}

data "microsoft365wp_authentication_combination_configuration" "all_for_one" {
  authentication_strength_policy_id = data.microsoft365wp_authentication_strength_policy.one.id

  for_each = toset(data.microsoft365wp_authentication_strength_policy.one.combination_configurations[*].id)
  id       = each.key
}

output "microsoft365wp_authentication_strength_policy" {
  value = merge(
    data.microsoft365wp_authentication_strength_policy.one,
    { combination_configurations = [for k, v in data.microsoft365wp_authentication_combination_configuration.all_for_one : v] }
  )
}
