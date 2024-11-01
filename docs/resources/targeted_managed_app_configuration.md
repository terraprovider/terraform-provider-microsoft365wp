---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "microsoft365wp_targeted_managed_app_configuration Resource - microsoft365wp"
subcategory: ""
description: |-
  Configuration used to deliver a set of custom settings as-is to all users in the targeted security group / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-targetedManagedAppConfiguration?view=graph-rest-beta
---

# microsoft365wp_targeted_managed_app_configuration (Resource)

Configuration used to deliver a set of custom settings as-is to all users in the targeted security group / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-targetedManagedAppConfiguration?view=graph-rest-beta

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


resource "microsoft365wp_targeted_managed_app_configuration" "test" {
  display_name = "TF Test"

  app_group_type = "selectedPublicApps"
  apps = [
    { app_id = { android = { package_id = "com.microsoft.office.outlook" } } },
    { app_id = { ios = { bundle_id = "com.microsoft.office.outlook" } } },
  ]

  custom_settings = [
    {
      name  = "General Test"
      value = "General Value"
    },
    {
      name  = "com.microsoft.outlook.Mail.FocusedInbox"
      value = "false"
    },
  ]

  assignments = [
    { target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } },
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `app_group_type` (String) Public Apps selection: group or individual / Indicates a collection of apps to target which can be one of several pre-defined lists of apps or a manually selected list of apps; possible values are: `selectedPublicApps` (Target the collection of apps manually selected by the admin.), `allCoreMicrosoftApps` (Target the core set of Microsoft apps (Office, Edge, etc).), `allMicrosoftApps` (Target all apps with Microsoft as publisher.), `allApps` (Target all apps with an available assignment.)
- `apps` (Attributes Set) List of apps to which the policy is deployed. / The identifier for the deployment an app. (see [below for nested schema](#nestedatt--apps))
- `custom_settings` (Attributes Set) A set of string key and string value pairs to be sent to apps for users to whom the configuration is scoped, unalterned by this service / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyValuePair?view=graph-rest-beta (see [below for nested schema](#nestedatt--custom_settings))
- `display_name` (String)

### Optional

- `assignments` (Attributes Set) The list of assignments. (see [below for nested schema](#nestedatt--assignments))
- `description` (String)
- `role_scope_tag_ids` (Set of String)
- `targeted_app_management_levels` (String) The intended app management levels for this policy / Management levels for apps; possible values are: `unspecified` (Unspecified), `unmanaged` (Unmanaged), `mdm` (MDM), `androidEnterprise` (Android Enterprise), `androidEnterpriseDedicatedDevicesWithAzureAdSharedMode` (Android Enterprise dedicated devices with Azure AD Shared mode), `androidOpenSourceProjectUserAssociated` (Android Open Source Project (AOSP) devices), `androidOpenSourceProjectUserless` (Android Open Source Project (AOSP) userless devices), `unknownFutureValue` (Place holder for evolvable enum)

### Read-Only

- `created_date_time` (String)
- `deployed_app_count` (Number) Count of apps to which the current policy is deployed.
- `id` (String) The ID of this resource.
- `is_assigned` (Boolean) Indicates if the policy is deployed to any inclusion groups or not.
- `last_modified_date_time` (String)
- `version` (String)

<a id="nestedatt--apps"></a>
### Nested Schema for `apps`

Required:

- `app_id` (Attributes) The identifier for an app with it's operating system type. / The identifier for a mobile app. (see [below for nested schema](#nestedatt--apps--app_id))

<a id="nestedatt--apps--app_id"></a>
### Nested Schema for `apps.app_id`

Optional:

- `android` (Attributes) The identifier for an Android app. (see [below for nested schema](#nestedatt--apps--app_id--android))
- `ios` (Attributes) The identifier for an iOS app. (see [below for nested schema](#nestedatt--apps--app_id--ios))
- `mac` (Attributes) The identifier for a Mac app. (see [below for nested schema](#nestedatt--apps--app_id--mac))
- `windows` (Attributes) The identifier for a Windows app. (see [below for nested schema](#nestedatt--apps--app_id--windows))

<a id="nestedatt--apps--app_id--android"></a>
### Nested Schema for `apps.app_id.android`

Required:

- `package_id` (String) The identifier for an app, as specified in the play store.


<a id="nestedatt--apps--app_id--ios"></a>
### Nested Schema for `apps.app_id.ios`

Required:

- `bundle_id` (String) The identifier for an app, as specified in the app store.


<a id="nestedatt--apps--app_id--mac"></a>
### Nested Schema for `apps.app_id.mac`

Required:

- `bundle_id` (String) The identifier for an app, as specified in the app store.


<a id="nestedatt--apps--app_id--windows"></a>
### Nested Schema for `apps.app_id.windows`

Required:

- `app_id` (String) The identifier for an app, as specified in the app store.




<a id="nestedatt--custom_settings"></a>
### Nested Schema for `custom_settings`

Required:

- `name` (String)

Optional:

- `value` (String)


<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `target` (Attributes) Base type for assignment targets. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-deviceAndAppManagementAssignmentTarget?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--target))

<a id="nestedatt--assignments--target"></a>
### Nested Schema for `assignments.target`

Optional:

- `all_devices` (Attributes) Represents an assignment to all managed devices in the tenant. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-allDevicesAssignmentTarget?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--target--all_devices))
- `all_licensed_users` (Attributes) Represents an assignment to all licensed users in the tenant. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-allLicensedUsersAssignmentTarget?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--target--all_licensed_users))
- `exclusion_group` (Attributes) Represents a group that should be excluded from an assignment. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-exclusionGroupAssignmentTarget?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--target--exclusion_group))
- `filter_id` (String) The Id of the filter for the target assignment.
- `filter_type` (String) The type of filter of the target assignment i.e. Exclude or Include. / Represents type of the assignment filter; possible values are: `none` (Default value. Do not use.), `include` (Indicates in-filter, rule matching will offer the payload to devices.), `exclude` (Indicates out-filter, rule matching will not offer the payload to devices.)
- `group` (Attributes) Represents an assignment to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-groupAssignmentTarget?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--target--group))

<a id="nestedatt--assignments--target--all_devices"></a>
### Nested Schema for `assignments.target.all_devices`


<a id="nestedatt--assignments--target--all_licensed_users"></a>
### Nested Schema for `assignments.target.all_licensed_users`


<a id="nestedatt--assignments--target--exclusion_group"></a>
### Nested Schema for `assignments.target.exclusion_group`

Required:

- `group_id` (String) AAD Group Id.


<a id="nestedatt--assignments--target--group"></a>
### Nested Schema for `assignments.target.group`

Required:

- `group_id` (String) AAD Group Id.
