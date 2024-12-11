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


resource "microsoft365wp_cloud_pc_provisioning_policy" "test_dedicated" {
  display_name = "TF Test - dedicated"

  image_id          = "microsoftwindowsdesktop_windows-ent-cpc_win11-23h2-ent-cpc-m365"
  provisioning_type = "dedicated"

  windows_setting = {
    locale = "en-US"
  }

  domain_join_configurations = [{
    region_group = "europeUnion"
    region_name  = "automatic"
  }]

  assignments = [
    for x in [
      "298fded6-b252-4166-a473-f405e935f58d",
      "62e39046-aad3-4423-98e0-b486e3538aff",
    ] :
    { target = { group = { group_id = x } } }
  ]
}


resource "microsoft365wp_cloud_pc_provisioning_policy" "test_shared" {
  display_name = "TF Test - shared"

  image_id          = "microsoftwindowsdesktop_windows-ent-cpc_win11-23h2-ent-cpc-m365"
  provisioning_type = "sharedByUser"

  windows_setting = {
    locale = "en-US"
  }

  domain_join_configurations = [{
    region_group = "europeUnion"
    region_name  = "automatic"
  }]

  assignments = [
    for x in [
      "298fded6-b252-4166-a473-f405e935f58d",
      "62e39046-aad3-4423-98e0-b486e3538aff",
    ] :
    { target = { group = {
      group_id        = x
      service_plan_id = "51855c77-4d2e-4736-be67-6dca605f2b57" // CPC_S_2C_4GB_128GB / Windows 365 Shared Use 2 vCPU, 4 GB, 128 GB
    } } }
  ]
}
