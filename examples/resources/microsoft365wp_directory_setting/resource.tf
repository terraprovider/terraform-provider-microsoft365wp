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
  directory_settings = {
    group_unified = {
      template_id = "62375ab9-6b52-47ed-826b-58e47e0e304b" # Group.Unified
      values = {
        EnableGroupCreation       = "true"
        AllowGuestsToAccessGroups = "true"
        UsageGuidelinesUrl        = "https://contoso.example/groups"
      }
    }
    password_rule_settings = {
      template_id = "5cf42378-d67d-4f36-ba46-e8b86229381d" # Password Rule Settings
      values = {
        LockoutThreshold         = "10"
        LockoutDurationInSeconds = "60"
      }
    }
  }
}

# Settings that already exist in the tenant must be imported first (by template id). This imports only those that
# exist, all others will be created.
data "microsoft365wp_directory_settings" "all" {
}

import {
  for_each = {
    for k, v in local.directory_settings : k => v.template_id
    if contains([for s in data.microsoft365wp_directory_settings.all.directory_settings : s.template_id], v.template_id)
  }
  to = microsoft365wp_directory_setting.this[each.key]
  id = each.value
}

resource "microsoft365wp_directory_setting" "this" {
  for_each    = local.directory_settings
  template_id = each.value.template_id
  values      = each.value.values
}
