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


locals {
  json_files_dir = "${path.module}/../../../../../_research/_EXPORT/SettingsCatalog"
  json_files     = fileset(local.json_files_dir, "*.json")
  policies       = { for f in local.json_files : f => jsondecode(file("${local.json_files_dir}/${f}")) }
}

resource "microsoft365wp_device_management_configuration_policy_json" "main" {
  for_each      = local.policies
  name          = each.value.name
  description   = each.value.description
  platforms     = each.value.platforms
  technologies  = each.value.technologies
  settings_json = [for x in each.value.settings : jsonencode(x.settingInstance)]
}
