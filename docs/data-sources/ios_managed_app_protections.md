---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "microsoft365wp_ios_managed_app_protections Data Source - microsoft365wp"
subcategory: ""
description: |-
  
---

# microsoft365wp_ios_managed_app_protections (Data Source)



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


data "microsoft365wp_ios_managed_app_protections" "all" {
  exclude_ids = ["01234567-89ab-cdef-0123-456789abcdee", "01234567-89ab-cdef-0123-456789abcdef"]
}

output "microsoft365wp_ios_managed_app_protections" {
  value = { for x in data.microsoft365wp_ios_managed_app_protections.all.ios_managed_app_protections : x.id => x }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `exclude_ids` (Set of String) Filter query to exclude objects with these ids.
- `include_ids` (Set of String) Filter query to only return objects with these ids.
- `odata_filter` (String) Raw OData $filter string to pass to MS Graph.

### Read-Only

- `ios_managed_app_protections` (Attributes Set) (see [below for nested schema](#nestedatt--ios_managed_app_protections))

<a id="nestedatt--ios_managed_app_protections"></a>
### Nested Schema for `ios_managed_app_protections`

Required:

- `id` (String)

Read-Only:

- `created_date_time` (String)
- `display_name` (String)
- `last_modified_date_time` (String)
- `role_scope_tag_ids` (Set of String)
- `version` (String)
