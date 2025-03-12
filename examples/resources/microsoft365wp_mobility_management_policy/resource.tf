terraform {
  required_providers {
    microsoft365wp = {
      source = "terraprovider/microsoft365wp"
    }
  }
}


data "azuread_application_published_app_ids" "well_known" {}

locals {
  device = {
    applies_to = "all" # "all", "selected" or "none"
    groups = [
      "298fded6-b252-4166-a473-f405e935f58d",
      "62e39046-aad3-4423-98e0-b486e3538aff",
    ]
  }
  app = {
    applies_to = "all" # "all", "selected" or "none"
    groups = [
      "298fded6-b252-4166-a473-f405e935f58d",
    ]
  }
}

resource "microsoft365wp_mobility_management_policy" "device" {
  policy_type = "device"
  id          = data.azuread_application_published_app_ids.well_known.result.MicrosoftIntune

  terms_of_use_url = "https://portal.manage.microsoft.com/TermsofUse.aspx"
  discovery_url    = "https://enrollment.manage.microsoft.com/enrollmentserver/discovery.svc"
  compliance_url   = "https://portal.manage.microsoft.com/?portalAction=Compliance"

  applies_to      = local.device.applies_to
  included_groups = local.device.applies_to == "selected" ? [for i in local.device.groups : { id = i }] : []
}

resource "microsoft365wp_mobility_management_policy" "app" {
  policy_type = "app"
  id          = data.azuread_application_published_app_ids.well_known.result.MicrosoftIntune

  discovery_url = "https://wip.mam.manage.microsoft.com/Enroll"

  applies_to      = local.app.applies_to
  included_groups = local.app.applies_to == "selected" ? [for i in local.app.groups : { id = i }] : []
}
