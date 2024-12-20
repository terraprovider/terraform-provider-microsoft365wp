---
page_title: "microsoft365wp_managed_device_mobile_app_configuration Resource - microsoft365wp"
subcategory: "MS Graph: App management"
---

# microsoft365wp_managed_device_mobile_app_configuration (Resource)

An abstract class for Mobile app configuration for enrolled devices. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-manageddevicemobileappconfiguration?view=graph-rest-beta

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


data "microsoft365wp_mobile_apps" "outlook_ios" {
  odata_filter = "contains(displayName, 'Outlook') and isof('microsoft.graph.iosStoreApp')"
}

resource "microsoft365wp_managed_device_mobile_app_configuration" "ios" {
  display_name = "TF Test iOS"

  targeted_mobile_apps = [
    one(data.microsoft365wp_mobile_apps.outlook_ios.mobile_apps[*].id),
  ]

  ios = {
    settings = [
      {
        app_config_key       = "Test Key"
        app_config_key_type  = "stringType"
        app_config_key_value = "Test Value"
      },
      {
        app_config_key       = "com.microsoft.outlook.Mail.FocusedInbox"
        app_config_key_type  = "booleanType"
        app_config_key_value = "true"
      },
    ]
  }
}


/*
##### Untested #####

resource "microsoft365wp_managed_device_mobile_app_configuration" "android_managed_store" {
  display_name = "TF Test Android Managed Store"

  targeted_mobile_apps = [
    "0ddf9bbb-abb6-4c9a-89a5-ff54d9640d40",
  ]

  android_managed_store = {
    profile_applicability = "default"
    package_id            = "app:com.microsoft.office.outlook"
    payload_json_base64 = base64encode(jsonencode({
      kind      = "androidenterprise#managedConfiguration"
      productId = "app:com.microsoft.office.outlook"
      managedProperty = [
        { key = "com.microsoft.outlook.EmailProfile.AccountType", valueString = "ModernAuth" },
        { key = "com.microsoft.outlook.EmailProfile.EmailUPN", valueString = "{{userprincipalname}}" },
        { key = "com.microsoft.outlook.EmailProfile.EmailAddress", valueString = "{{mail}}" },
        { key = "com.microsoft.outlook.Contacts.LocalSyncEnabled", valueBool = true },
        { key = "com.microsoft.outlook.Mail.DefaultSignatureEnabled", valueBool = false },
      ]
    }))
  }
}
*/
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) Admin provided name of the device configuration.
- `targeted_mobile_apps` (Set of String) the associated app.

### Optional

- `android_managed_store` (Attributes) Contains properties, inherited properties and actions for Android Enterprise mobile app configurations. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidmanagedstoreappconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--android_managed_store))
- `assignments` (Attributes Set) The list of assignments. (see [below for nested schema](#nestedatt--assignments))
- `description` (String) Admin provided description of the Device Configuration. The _provider_ default value is `""`.
- `ios` (Attributes) Contains properties, inherited properties and actions for iOS mobile app configurations. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosmobileappconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios))
- `role_scope_tag_ids` (Set of String) List of Scope Tags for this App configuration entity. The _provider_ default value is `["0"]`.

### Read-Only

- `created_date_time` (String) DateTime the object was created.
- `id` (String) Key of the entity.
- `last_modified_date_time` (String) DateTime the object was last modified.
- `version` (Number) Version of the device configuration.

<a id="nestedatt--android_managed_store"></a>
### Nested Schema for `android_managed_store`

Required:

- `package_id` (String) Android Enterprise app configuration package id.
- `profile_applicability` (String) Android Enterprise profile applicability (AndroidWorkProfile, DeviceOwner, or default (applies to both)). / Android profile applicability; possible values are: `default`, `androidWorkProfile`, `androidDeviceOwner`

Optional:

- `connected_apps_enabled` (Boolean) Setting to specify whether to allow ConnectedApps experience for this app. The _provider_ default value is `false`.
- `payload_json` (String) Android Enterprise app configuration JSON payload.
- `permission_actions` (Attributes Set) List of Android app permissions and corresponding permission actions. / Mapping between an Android app permission and the action Android should take when that permission is requested. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-androidpermissionaction?view=graph-rest-beta (see [below for nested schema](#nestedatt--android_managed_store--permission_actions))

Read-Only:

- `app_supports_oem_config` (Boolean) Whether or not this AppConfig is an OEMConfig policy. This property is read-only.

<a id="nestedatt--android_managed_store--permission_actions"></a>
### Nested Schema for `android_managed_store.permission_actions`

Required:

- `action` (String) Type of Android permission action. / Android action taken when an app requests a dangerous permission; possible values are: `prompt`, `autoGrant`, `autoDeny`

Optional:

- `permission` (String) Android permission string, defined in the official Android documentation.  Example 'android.permission.READ_CONTACTS'.



<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Required:

- `target` (Attributes) Base type for assignment targets. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-deviceandappmanagementassignmenttarget?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--target))

<a id="nestedatt--assignments--target"></a>
### Nested Schema for `assignments.target`

Optional:

- `all_devices` (Attributes) Represents an assignment to all managed devices in the tenant. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-alldevicesassignmenttarget?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--target--all_devices))
- `all_licensed_users` (Attributes) Represents an assignment to all licensed users in the tenant. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-alllicensedusersassignmenttarget?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--target--all_licensed_users))
- `exclusion_group` (Attributes) Represents a group that should be excluded from an assignment. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-exclusiongroupassignmenttarget?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--target--exclusion_group))
- `filter_id` (String) The ID of the filter for the target assignment.
- `filter_type` (String) The type of filter of the target assignment i.e. Exclude or Include. / Represents type of the assignment filter; possible values are: `none` (Default value. Do not use.), `include` (Indicates in-filter, rule matching will offer the payload to devices.), `exclude` (Indicates out-filter, rule matching will not offer the payload to devices.). The _provider_ default value is `"none"`.
- `group` (Attributes) Represents an assignment to a group. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-groupassignmenttarget?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--target--group))

<a id="nestedatt--assignments--target--all_devices"></a>
### Nested Schema for `assignments.target.all_devices`


<a id="nestedatt--assignments--target--all_licensed_users"></a>
### Nested Schema for `assignments.target.all_licensed_users`


<a id="nestedatt--assignments--target--exclusion_group"></a>
### Nested Schema for `assignments.target.exclusion_group`

Required:

- `group_id` (String) The group Id that is the target of the assignment.


<a id="nestedatt--assignments--target--group"></a>
### Nested Schema for `assignments.target.group`

Required:

- `group_id` (String) The group Id that is the target of the assignment.




<a id="nestedatt--ios"></a>
### Nested Schema for `ios`

Optional:

- `encoded_setting_xml` (String) mdm app configuration Base64 binary.
- `settings` (Attributes Set) app configuration setting items. / Contains properties for App configuration setting item. / https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-appconfigurationsettingitem?view=graph-rest-beta (see [below for nested schema](#nestedatt--ios--settings))

<a id="nestedatt--ios--settings"></a>
### Nested Schema for `ios.settings`

Required:

- `app_config_key` (String) app configuration key.
- `app_config_key_type` (String) app configuration key type. / App configuration key types; possible values are: `stringType`, `integerType`, `realType`, `booleanType`, `tokenType`
- `app_config_key_value` (String) app configuration key value.
