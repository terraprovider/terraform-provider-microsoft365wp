---
page_title: "microsoft365wp_mobile_app_categories Data Source - microsoft365wp"
subcategory: "MS Graph: App management"
---

# microsoft365wp_mobile_app_categories (Data Source)

Contains properties for a single Intune app category. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-mobileappcategory?view=graph-rest-beta

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


data "microsoft365wp_mobile_app_categories" "all" {
  exclude_ids = ["01234567-89ab-cdef-0123-456789abcdee", "01234567-89ab-cdef-0123-456789abcdef"]
}

output "microsoft365wp_mobile_app_categories" {
  value = { for x in data.microsoft365wp_mobile_app_categories.all.mobile_app_categories : x.id => x }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `exclude_ids` (Set of String) Filter query to exclude objects with these ids.
- `include_ids` (Set of String) Filter query to only return objects with these ids.
- `odata_filter` (String) Raw OData $filter string to pass to MS Graph.

### Read-Only

- `mobile_app_categories` (Attributes Set) (see [below for nested schema](#nestedatt--mobile_app_categories))

<a id="nestedatt--mobile_app_categories"></a>
### Nested Schema for `mobile_app_categories`

Required:

- `id` (String) The key of the entity. This property is read-only.

Read-Only:

- `display_name` (String) The name of the app category.
- `last_modified_date_time` (String) The date and time the mobileAppCategory was last modified. This property is read-only.  
Provider Note: Warning: This attribute seems to always return the _current_ time for mobile app categories that have been created for the tenant (i.e. that have not been predefined by Microsoft). Therefore it can be expected to change with every query.
