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


data "microsoft365wp_identity_governance_workflow_versions" "all" {
  workflow_id = "5fe4688c-3b34-481b-b5d2-a5f3033acade"
}

output "microsoft365wp_identity_governance_workflow_versions" {
  value = { for x in data.microsoft365wp_identity_governance_workflow_versions.all.identity_governance_workflow_versions : x.version_number => x }
}
