# TerraProvider - Terraform Provider for Microsoft 365

This provider instruments [Microsoft’s Graph API](https://learn.microsoft.com/en-us/graph/overview) to allow managing Intune and Entra ID configuration objects and policies via Terraform and OpenTofu. 

- see [Docs](https://docs.terraprovider.com/) for guides as well as data source and resource definitions
- see [TerraProvider QuickStart](https://github.com/terraprovider/terraform-provider-microsoft365wp-quickstart) for a simple reference implementation and samples of resources with this provider.

Visit [TerraProvider.com](https://terraprovider.com) for more information.

## Usage Example

```
# Configure Terraform
terraform {
  required_providers {
    azuread = {
      source  = "terraprovider/microsoft365wp"
    }
  }
}

# Create an Intune Compliance Policy
resource "microsoft365wp_device_compliance_policy" "all" {
  display_name = "Windows - Defender for Endpoint"
  assignments = [
    {
      target = { all_licensed_users = {} }
    }
  ]
  windows10 = {
    device_threat_protection_enabled = true
  }
  scheduled_actions_for_rule = [
    {
      scheduled_action_configurations = [
        {
          action_type        = "block"
          grace_period_hours = 6
        },
      ]
    },
  ]
}
```

## Authentication

You can authenticate using an Entra ID Service Principal, see either using a [Client Secret](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/guides/service_principal_client_secret) or [OpenID Connect](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/guides/service_principal_oicd).

In both cases the following ENV variables must be set:
- `ARM_TENANT_ID`
- `ARM_CLIENT_ID`

and if you use a Client Secret, also set `ARM_CLIENT_SECRET`
