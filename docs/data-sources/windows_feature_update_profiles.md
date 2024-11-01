---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "microsoft365wp_windows_feature_update_profiles Data Source - microsoft365wp"
subcategory: ""
description: |-
  
---

# microsoft365wp_windows_feature_update_profiles (Data Source)



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


data "microsoft365wp_windows_feature_update_profiles" "all" {
}

output "microsoft365wp_windows_feature_update_profiles" {
  value = { for x in data.microsoft365wp_windows_feature_update_profiles.all.windows_feature_update_profiles : x.id => x }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `windows_feature_update_profiles` (Attributes Set) (see [below for nested schema](#nestedatt--windows_feature_update_profiles))

<a id="nestedatt--windows_feature_update_profiles"></a>
### Nested Schema for `windows_feature_update_profiles`

Required:

- `id` (String)

Read-Only:

- `created_date_time` (String) The date time that the profile was created.
- `display_name` (String) The display name of the profile.
- `last_modified_date_time` (String) The date time that the profile was last modified.
- `role_scope_tag_ids` (Set of String) List of Scope Tags for this Feature Update entity.
