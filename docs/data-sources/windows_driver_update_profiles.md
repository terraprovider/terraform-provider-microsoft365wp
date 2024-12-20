---
page_title: "microsoft365wp_windows_driver_update_profiles Data Source - microsoft365wp"
subcategory: "MS Graph: Software updates"
---

# microsoft365wp_windows_driver_update_profiles (Data Source)

Windows Driver Update Profile / https://learn.microsoft.com/en-us/graph/api/resources/intune-softwareupdate-windowsdriverupdateprofile?view=graph-rest-beta

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


data "microsoft365wp_windows_driver_update_profiles" "all" {
}

output "microsoft365wp_windows_driver_update_profiles" {
  value = { for x in data.microsoft365wp_windows_driver_update_profiles.all.windows_driver_update_profiles : x.id => x }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `windows_driver_update_profiles` (Attributes Set) (see [below for nested schema](#nestedatt--windows_driver_update_profiles))

<a id="nestedatt--windows_driver_update_profiles"></a>
### Nested Schema for `windows_driver_update_profiles`

Required:

- `id` (String) The Intune policy id.

Read-Only:

- `created_date_time` (String) The date time that the profile was created.
- `display_name` (String) The display name for the profile.
- `last_modified_date_time` (String) The date time that the profile was last modified.
- `role_scope_tag_ids` (Set of String) List of Scope Tags for this Driver Update entity. The _provider_ default value is `["0"]`.
