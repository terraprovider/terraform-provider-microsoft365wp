---
page_title: "microsoft365wp_device_configurations Data Source - microsoft365wp"
subcategory: "MS Graph: Device configuration"
---

# microsoft365wp_device_configurations (Data Source)

Device Configuration. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceconfiguration?view=graph-rest-beta

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


data "microsoft365wp_device_configurations" "all" {
  exclude_ids = ["01234567-89ab-cdef-0123-456789abcdee", "01234567-89ab-cdef-0123-456789abcdef"]
}

output "microsoft365wp_device_configurations" {
  value = { for x in data.microsoft365wp_device_configurations.all.device_configurations : x.id => x }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `exclude_ids` (Set of String) Filter query to exclude objects with these ids.
- `include_ids` (Set of String) Filter query to only return objects with these ids.
- `odata_filter` (String) Raw OData $filter string to pass to MS Graph.

### Read-Only

- `device_configurations` (Attributes Set) (see [below for nested schema](#nestedatt--device_configurations))

<a id="nestedatt--device_configurations"></a>
### Nested Schema for `device_configurations`

Required:

- `id` (String) Key of the entity.

Read-Only:

- `android_device_owner_general_device` (Attributes) Please note that this nested object does not have any attributes but only exists to be able to test if the parent object is of derived OData type `#microsoft.graph.androidDeviceOwnerGeneralDeviceConfiguration` (using e.g. `if x.android_device_owner_general_device != null`). (see [below for nested schema](#nestedatt--device_configurations--android_device_owner_general_device))
- `android_work_profile_general_device` (Attributes) Please note that this nested object does not have any attributes but only exists to be able to test if the parent object is of derived OData type `#microsoft.graph.androidWorkProfileGeneralDeviceConfiguration` (using e.g. `if x.android_work_profile_general_device != null`). (see [below for nested schema](#nestedatt--device_configurations--android_work_profile_general_device))
- `created_date_time` (String) DateTime the object was created.
- `display_name` (String) Admin provided name of the device configuration.
- `edition_upgrade` (Attributes) Please note that this nested object does not have any attributes but only exists to be able to test if the parent object is of derived OData type `#microsoft.graph.editionUpgradeConfiguration` (using e.g. `if x.edition_upgrade != null`). (see [below for nested schema](#nestedatt--device_configurations--edition_upgrade))
- `ios_custom` (Attributes) Please note that this nested object does not have any attributes but only exists to be able to test if the parent object is of derived OData type `#microsoft.graph.iosCustomConfiguration` (using e.g. `if x.ios_custom != null`). (see [below for nested schema](#nestedatt--device_configurations--ios_custom))
- `ios_device_features` (Attributes) Please note that this nested object does not have any attributes but only exists to be able to test if the parent object is of derived OData type `#microsoft.graph.iosDeviceFeaturesConfiguration` (using e.g. `if x.ios_device_features != null`). (see [below for nested schema](#nestedatt--device_configurations--ios_device_features))
- `ios_eas_email_profile` (Attributes) Please note that this nested object does not have any attributes but only exists to be able to test if the parent object is of derived OData type `#microsoft.graph.iosEasEmailProfileConfiguration` (using e.g. `if x.ios_eas_email_profile != null`). (see [below for nested schema](#nestedatt--device_configurations--ios_eas_email_profile))
- `ios_general_device` (Attributes) Please note that this nested object does not have any attributes but only exists to be able to test if the parent object is of derived OData type `#microsoft.graph.iosGeneralDeviceConfiguration` (using e.g. `if x.ios_general_device != null`). (see [below for nested schema](#nestedatt--device_configurations--ios_general_device))
- `ios_update` (Attributes) Please note that this nested object does not have any attributes but only exists to be able to test if the parent object is of derived OData type `#microsoft.graph.iosUpdateConfiguration` (using e.g. `if x.ios_update != null`). (see [below for nested schema](#nestedatt--device_configurations--ios_update))
- `ios_vpn` (Attributes) Please note that this nested object does not have any attributes but only exists to be able to test if the parent object is of derived OData type `#microsoft.graph.iosVpnConfiguration` (using e.g. `if x.ios_vpn != null`). (see [below for nested schema](#nestedatt--device_configurations--ios_vpn))
- `last_modified_date_time` (String) DateTime the object was last modified.
- `macos_custom` (Attributes) Please note that this nested object does not have any attributes but only exists to be able to test if the parent object is of derived OData type `#microsoft.graph.macOSCustomConfiguration` (using e.g. `if x.macos_custom != null`). (see [below for nested schema](#nestedatt--device_configurations--macos_custom))
- `macos_custom_app` (Attributes) Please note that this nested object does not have any attributes but only exists to be able to test if the parent object is of derived OData type `#microsoft.graph.macOSCustomAppConfiguration` (using e.g. `if x.macos_custom_app != null`). (see [below for nested schema](#nestedatt--device_configurations--macos_custom_app))
- `macos_device_features` (Attributes) Please note that this nested object does not have any attributes but only exists to be able to test if the parent object is of derived OData type `#microsoft.graph.macOSDeviceFeaturesConfiguration` (using e.g. `if x.macos_device_features != null`). (see [below for nested schema](#nestedatt--device_configurations--macos_device_features))
- `macos_extensions` (Attributes) Please note that this nested object does not have any attributes but only exists to be able to test if the parent object is of derived OData type `#microsoft.graph.macOSExtensionsConfiguration` (using e.g. `if x.macos_extensions != null`). (see [below for nested schema](#nestedatt--device_configurations--macos_extensions))
- `macos_software_update` (Attributes) Please note that this nested object does not have any attributes but only exists to be able to test if the parent object is of derived OData type `#microsoft.graph.macOSSoftwareUpdateConfiguration` (using e.g. `if x.macos_software_update != null`). (see [below for nested schema](#nestedatt--device_configurations--macos_software_update))
- `role_scope_tag_ids` (Set of String) List of Scope Tags for this Entity instance. The _provider_ default value is `["0"]`.
- `version` (Number) Version of the device configuration.
- `windows_health_monitoring` (Attributes) Please note that this nested object does not have any attributes but only exists to be able to test if the parent object is of derived OData type `#microsoft.graph.windowsHealthMonitoringConfiguration` (using e.g. `if x.windows_health_monitoring != null`). (see [below for nested schema](#nestedatt--device_configurations--windows_health_monitoring))
- `windows_update_for_business` (Attributes) Please note that this nested object does not have any attributes but only exists to be able to test if the parent object is of derived OData type `#microsoft.graph.windowsUpdateForBusinessConfiguration` (using e.g. `if x.windows_update_for_business != null`). (see [below for nested schema](#nestedatt--device_configurations--windows_update_for_business))

<a id="nestedatt--device_configurations--android_device_owner_general_device"></a>
### Nested Schema for `device_configurations.android_device_owner_general_device`


<a id="nestedatt--device_configurations--android_work_profile_general_device"></a>
### Nested Schema for `device_configurations.android_work_profile_general_device`


<a id="nestedatt--device_configurations--edition_upgrade"></a>
### Nested Schema for `device_configurations.edition_upgrade`


<a id="nestedatt--device_configurations--ios_custom"></a>
### Nested Schema for `device_configurations.ios_custom`


<a id="nestedatt--device_configurations--ios_device_features"></a>
### Nested Schema for `device_configurations.ios_device_features`


<a id="nestedatt--device_configurations--ios_eas_email_profile"></a>
### Nested Schema for `device_configurations.ios_eas_email_profile`


<a id="nestedatt--device_configurations--ios_general_device"></a>
### Nested Schema for `device_configurations.ios_general_device`


<a id="nestedatt--device_configurations--ios_update"></a>
### Nested Schema for `device_configurations.ios_update`


<a id="nestedatt--device_configurations--ios_vpn"></a>
### Nested Schema for `device_configurations.ios_vpn`


<a id="nestedatt--device_configurations--macos_custom"></a>
### Nested Schema for `device_configurations.macos_custom`


<a id="nestedatt--device_configurations--macos_custom_app"></a>
### Nested Schema for `device_configurations.macos_custom_app`


<a id="nestedatt--device_configurations--macos_device_features"></a>
### Nested Schema for `device_configurations.macos_device_features`


<a id="nestedatt--device_configurations--macos_extensions"></a>
### Nested Schema for `device_configurations.macos_extensions`


<a id="nestedatt--device_configurations--macos_software_update"></a>
### Nested Schema for `device_configurations.macos_software_update`


<a id="nestedatt--device_configurations--windows_health_monitoring"></a>
### Nested Schema for `device_configurations.windows_health_monitoring`


<a id="nestedatt--device_configurations--windows_update_for_business"></a>
### Nested Schema for `device_configurations.windows_update_for_business`
