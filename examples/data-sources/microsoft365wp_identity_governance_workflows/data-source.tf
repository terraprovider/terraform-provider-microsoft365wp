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


data "microsoft365wp_identity_governance_workflows" "all" {}

output "microsoft365wp_identity_governance_workflows" {
  value = { for x in data.microsoft365wp_identity_governance_workflows.all.identity_governance_workflows : x.id => x }
}
