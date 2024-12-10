---
page_title: "microsoft365wp_authentication_context_class_references Data Source - microsoft365wp"
subcategory: "MS Graph: Conditional access"
---

# microsoft365wp_authentication_context_class_references (Data Source)

Represents a Microsoft Entra authentication context class reference. Authentication context class references are custom values that define a Conditional Access authentication requirement. / https://learn.microsoft.com/en-us/graph/api/resources/authenticationcontextclassreference?view=graph-rest-beta

## Documentation Disclaimer

Please note that almost all information on this page has been sourced literally from the official Microsoft Graph API 
documentation and therefore is governed by Microsoft and not by the publishers of this provider.  
All supplements authored by the publishers of this provider have been explicitly marked as such.

## Example Usage

```terraform
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


data "microsoft365wp_authentication_context_class_references" "all" {
  # exclude_ids = ["c25"]
}

output "microsoft365wp_authentication_context_class_references" {
  value = { for x in data.microsoft365wp_authentication_context_class_references.all.authentication_context_class_references : x.id => x }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `authentication_context_class_references` (Attributes Set) (see [below for nested schema](#nestedatt--authentication_context_class_references))

<a id="nestedatt--authentication_context_class_references"></a>
### Nested Schema for `authentication_context_class_references`

Required:

- `id` (String) Identifier used to reference the authentication context class. The ID is used to trigger step-up authentication for the referenced authentication requirements and is the value that will be issued in the `acrs` claim of an access token. This value in the claim is used to verify that the required authentication context has been satisfied. The allowed values are `c1` through `c25`. <br/> Supports `$filter` (`eq`).

Read-Only:

- `display_name` (String) A friendly name that identifies the authenticationContextClassReference object when building user-facing admin experiences. For example, a selection UX.