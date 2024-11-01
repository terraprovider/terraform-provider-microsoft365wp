---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "microsoft365wp_device_configuration_custom Data Source - microsoft365wp"
subcategory: ""
description: |-
  
---

# microsoft365wp_device_configuration_custom (Data Source)



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


data "microsoft365wp_device_configuration_custom" "one" {
  id = "1c527e3e-c8b3-4b6c-8ce0-671feb22f171"
}

output "microsoft365wp_device_configuration_custom" {
  value = data.microsoft365wp_device_configuration_custom.one
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `assignments` (Attributes Set) The list of assignments. (see [below for nested schema](#nestedatt--assignments))
- `created_date_time` (String) DateTime the object was created.
- `description` (String) Admin provided description of the Device Configuration.
- `device_management_applicability_rule_device_mode` (Attributes) The device mode applicability rule for this Policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceManagementApplicabilityRuleDeviceMode?view=graph-rest-beta (see [below for nested schema](#nestedatt--device_management_applicability_rule_device_mode))
- `device_management_applicability_rule_os_edition` (Attributes) The OS edition applicability for this Policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceManagementApplicabilityRuleOsEdition?view=graph-rest-beta (see [below for nested schema](#nestedatt--device_management_applicability_rule_os_edition))
- `device_management_applicability_rule_os_version` (Attributes) The OS version applicability rule for this Policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceManagementApplicabilityRuleOsVersion?view=graph-rest-beta (see [below for nested schema](#nestedatt--device_management_applicability_rule_os_version))
- `display_name` (String) Admin provided name of the device configuration.
- `id` (String) The ID of this resource.
- `last_modified_date_time` (String) DateTime the object was last modified.
- `role_scope_tag_ids` (Set of String) List of Scope Tags for this Entity instance.
- `supports_scope_tags` (Boolean)
- `version` (Number) Version of the device configuration.
- `windows10` (Attributes) This topic provides descriptions of the declared methods, properties and relationships exposed by the windows10CustomConfiguration resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windows10CustomConfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows10))

<a id="nestedatt--assignments"></a>
### Nested Schema for `assignments`

Read-Only:

- `target` (Attributes) Base type for assignment targets. / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-deviceAndAppManagementAssignmentTarget?view=graph-rest-beta (see [below for nested schema](#nestedatt--assignments--target))

<a id="nestedatt--assignments--target"></a>
### Nested Schema for `assignments.target`

Read-Only:

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

Read-Only:

- `group_id` (String) AAD Group Id.


<a id="nestedatt--assignments--target--group"></a>
### Nested Schema for `assignments.target.group`

Read-Only:

- `group_id` (String) AAD Group Id.




<a id="nestedatt--device_management_applicability_rule_device_mode"></a>
### Nested Schema for `device_management_applicability_rule_device_mode`

Read-Only:

- `device_mode` (String) Applicability rule for device mode. / Windows 10 Device Mode type; possible values are: `standardConfiguration` (Standard Configuration), `sModeConfiguration` (S Mode Configuration)
- `name` (String) Name for object.
- `rule_type` (String) Applicability Rule type. / Supported Applicability rule types for Device Configuration; possible values are: `include` (Include), `exclude` (Exclude)


<a id="nestedatt--device_management_applicability_rule_os_edition"></a>
### Nested Schema for `device_management_applicability_rule_os_edition`

Read-Only:

- `name` (String) Name for object.
- `os_edition_types` (Set of String) Applicability rule OS edition type. / Windows 10 Edition type; possible values are: `windows10Enterprise` (Windows 10 Enterprise), `windows10EnterpriseN` (Windows 10 EnterpriseN), `windows10Education` (Windows 10 Education), `windows10EducationN` (Windows 10 EducationN), `windows10MobileEnterprise` (Windows 10 Mobile Enterprise), `windows10HolographicEnterprise` (Windows 10 Holographic Enterprise), `windows10Professional` (Windows 10 Professional), `windows10ProfessionalN` (Windows 10 ProfessionalN), `windows10ProfessionalEducation` (Windows 10 Professional Education), `windows10ProfessionalEducationN` (Windows 10 Professional EducationN), `windows10ProfessionalWorkstation` (Windows 10 Professional for Workstations), `windows10ProfessionalWorkstationN` (Windows 10 Professional for Workstations N), `notConfigured` (NotConfigured), `windows10Home` (Windows 10 Home), `windows10HomeChina` (Windows 10 Home China), `windows10HomeN` (Windows 10 Home N), `windows10HomeSingleLanguage` (Windows 10 Home Single Language), `windows10Mobile` (Windows 10 Mobile), `windows10IoTCore` (Windows 10 IoT Core), `windows10IoTCoreCommercial` (Windows 10 IoT Core Commercial)
- `rule_type` (String) Applicability Rule type. / Supported Applicability rule types for Device Configuration; possible values are: `include` (Include), `exclude` (Exclude)


<a id="nestedatt--device_management_applicability_rule_os_version"></a>
### Nested Schema for `device_management_applicability_rule_os_version`

Read-Only:

- `max_os_version` (String) Max OS version for Applicability Rule.
- `min_os_version` (String) Min OS version for Applicability Rule.
- `name` (String) Name for object.
- `rule_type` (String) Applicability Rule type. / Supported Applicability rule types for Device Configuration; possible values are: `include` (Include), `exclude` (Exclude)


<a id="nestedatt--windows10"></a>
### Nested Schema for `windows10`

Read-Only:

- `oma_settings` (Attributes Set) OMA settings. This collection can contain a maximum of 1000 elements. / OMA Settings definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omaSetting?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows10--oma_settings))

<a id="nestedatt--windows10--oma_settings"></a>
### Nested Schema for `windows10.oma_settings`

Read-Only:

- `base64` (Attributes) OMA Settings Base64 definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omaSettingBase64?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows10--oma_settings--base64))
- `boolean` (Attributes) OMA Settings Boolean definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omaSettingBoolean?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows10--oma_settings--boolean))
- `date_time` (Attributes) OMA Settings DateTime definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omaSettingDateTime?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows10--oma_settings--date_time))
- `description` (String) Description.
- `display_name` (String) Display Name.
- `floating_point` (Attributes) OMA Settings Floating Point definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omaSettingFloatingPoint?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows10--oma_settings--floating_point))
- `integer` (Attributes) OMA Settings Integer definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omaSettingInteger?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows10--oma_settings--integer))
- `oma_uri` (String) OMA.
- `string` (Attributes) OMA Settings String definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omaSettingString?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows10--oma_settings--string))
- `string_xml` (Attributes) OMA Settings StringXML definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omaSettingStringXml?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows10--oma_settings--string_xml))

<a id="nestedatt--windows10--oma_settings--base64"></a>
### Nested Schema for `windows10.oma_settings.base64`

Read-Only:

- `file_name` (String) File name associated with the Value property (*.cer | *.crt | *.p7b | *.bin).
- `value_base64` (String) Value. (Base64 encoded string)


<a id="nestedatt--windows10--oma_settings--boolean"></a>
### Nested Schema for `windows10.oma_settings.boolean`

Read-Only:

- `value` (Boolean) Value.


<a id="nestedatt--windows10--oma_settings--date_time"></a>
### Nested Schema for `windows10.oma_settings.date_time`

Read-Only:

- `value` (String) Value.


<a id="nestedatt--windows10--oma_settings--floating_point"></a>
### Nested Schema for `windows10.oma_settings.floating_point`

Read-Only:

- `value` (Number) Value.


<a id="nestedatt--windows10--oma_settings--integer"></a>
### Nested Schema for `windows10.oma_settings.integer`

Read-Only:

- `value` (Number) Value.


<a id="nestedatt--windows10--oma_settings--string"></a>
### Nested Schema for `windows10.oma_settings.string`

Read-Only:

- `value` (String) Value.


<a id="nestedatt--windows10--oma_settings--string_xml"></a>
### Nested Schema for `windows10.oma_settings.string_xml`

Read-Only:

- `file_name` (String) File name associated with the Value property (*.xml).
- `value` (String) Value. (UTF8 encoded byte array)
