---
page_title: "microsoft365wp_device_configuration_custom Resource - microsoft365wp"
subcategory: "MS Graph: Device configuration"
---

# microsoft365wp_device_configuration_custom (Resource)

Device Configuration. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-deviceconfiguration?view=graph-rest-beta

Provider Note: When using values of type `base64`, `string` or `string_xml` the write permission `DeviceManagementConfiguration.ReadWrite.All` is required also for plan operations. This is due to the fact that values of these types are saved encrypted in MS Graph and need to be retrieved using the special MS Graph action `getOmaSettingPlainTextValue` (which requires write permissions) when reading. Also see https://learn.microsoft.com/en-us/graph/api/intune-deviceconfig-deviceconfiguration-getomasettingplaintextvalue

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


resource "microsoft365wp_device_configuration_custom" "win10" {
  display_name = "TF Test Windows 10 Custom"
  windows10 = {
    oma_settings = [
      {
        display_name = "OMA 1 Base64"
        oma_uri      = "./User/Vendor/C4A8/Base64"
        base64       = { value_base64 = base64encode("asdfa") }
      },
      {
        display_name = "OMA 2 Boolean"
        oma_uri      = "./User/Vendor/C4A8/Boolean"
        boolean      = { value = false }
      },
      {
        display_name = "OMA 3 DateTime"
        oma_uri      = "./User/Vendor/C4A8/DateTime"
        date_time    = { value = "2023-01-01T00:00:01Z" }
      },
      {
        display_name   = "OMA 4 FloatingPoint"
        oma_uri        = "./User/Vendor/C4A8/FloatingPoint"
        floating_point = { value = 42.1 }
      },
      {
        display_name = "OMA 5 Integer"
        oma_uri      = "./User/Vendor/C4A8/Integer"
        integer      = { value = 42 }
      },
      {
        display_name = "OMA 6 String"
        oma_uri      = "./User/Vendor/C4A8/String"
        string       = { value = "asdf" }
      },
      {
        display_name = "OMA 7 StringXml"
        oma_uri      = "./User/Vendor/C4A8/StringXml"
        string_xml   = { value = "<?xml version=\"1.0\" encoding=\"utf-8\"?><root>a</root>" }
      },
    ]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) Admin provided name of the device configuration.

### Optional

- `assignments` (Attributes Set) The list of assignments. (see [below for nested schema](#nestedatt--assignments))
- `description` (String) Admin provided description of the Device Configuration.
- `device_management_applicability_rule_device_mode` (Attributes) The device mode applicability rule for this Policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-devicemanagementapplicabilityruledevicemode?view=graph-rest-beta (see [below for nested schema](#nestedatt--device_management_applicability_rule_device_mode))
- `device_management_applicability_rule_os_edition` (Attributes) The OS edition applicability for this Policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-devicemanagementapplicabilityruleosedition?view=graph-rest-beta (see [below for nested schema](#nestedatt--device_management_applicability_rule_os_edition))
- `device_management_applicability_rule_os_version` (Attributes) The OS version applicability rule for this Policy. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-devicemanagementapplicabilityruleosversion?view=graph-rest-beta (see [below for nested schema](#nestedatt--device_management_applicability_rule_os_version))
- `role_scope_tag_ids` (Set of String) List of Scope Tags for this Entity instance. The _provider_ default value is `["0"]`.
- `windows10` (Attributes) This topic provides descriptions of the declared methods, properties and relationships exposed by the windows10CustomConfiguration resource. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-windows10customconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows10))

### Read-Only

- `created_date_time` (String) DateTime the object was created.
- `id` (String) Key of the entity.
- `last_modified_date_time` (String) DateTime the object was last modified.
- `supports_scope_tags` (Boolean) Indicates whether or not the underlying Device Configuration supports the assignment of scope tags. Assigning to the ScopeTags property is not allowed when this value is false and entities will not be visible to scoped users. This occurs for Legacy policies created in Silverlight and can be resolved by deleting and recreating the policy in the Azure Portal. This property is read-only.
- `version` (Number) Version of the device configuration.

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




<a id="nestedatt--device_management_applicability_rule_device_mode"></a>
### Nested Schema for `device_management_applicability_rule_device_mode`

Required:

- `device_mode` (String) Applicability rule for device mode. / Windows 10 Device Mode type; possible values are: `standardConfiguration` (Standard Configuration), `sModeConfiguration` (S Mode Configuration)
- `rule_type` (String) Applicability Rule type. / Supported Applicability rule types for Device Configuration; possible values are: `include` (Include), `exclude` (Exclude)

Optional:

- `name` (String) Name for object.


<a id="nestedatt--device_management_applicability_rule_os_edition"></a>
### Nested Schema for `device_management_applicability_rule_os_edition`

Required:

- `os_edition_types` (Set of String) Applicability rule OS edition type. / Windows 10 Edition type; possible values are: `windows10Enterprise` (Windows 10 Enterprise), `windows10EnterpriseN` (Windows 10 EnterpriseN), `windows10Education` (Windows 10 Education), `windows10EducationN` (Windows 10 EducationN), `windows10MobileEnterprise` (Windows 10 Mobile Enterprise), `windows10HolographicEnterprise` (Windows 10 Holographic Enterprise), `windows10Professional` (Windows 10 Professional), `windows10ProfessionalN` (Windows 10 ProfessionalN), `windows10ProfessionalEducation` (Windows 10 Professional Education), `windows10ProfessionalEducationN` (Windows 10 Professional EducationN), `windows10ProfessionalWorkstation` (Windows 10 Professional for Workstations), `windows10ProfessionalWorkstationN` (Windows 10 Professional for Workstations N), `notConfigured` (NotConfigured), `windows10Home` (Windows 10 Home), `windows10HomeChina` (Windows 10 Home China), `windows10HomeN` (Windows 10 Home N), `windows10HomeSingleLanguage` (Windows 10 Home Single Language), `windows10Mobile` (Windows 10 Mobile), `windows10IoTCore` (Windows 10 IoT Core), `windows10IoTCoreCommercial` (Windows 10 IoT Core Commercial)
- `rule_type` (String) Applicability Rule type. / Supported Applicability rule types for Device Configuration; possible values are: `include` (Include), `exclude` (Exclude)

Optional:

- `name` (String) Name for object.


<a id="nestedatt--device_management_applicability_rule_os_version"></a>
### Nested Schema for `device_management_applicability_rule_os_version`

Required:

- `rule_type` (String) Applicability Rule type. / Supported Applicability rule types for Device Configuration; possible values are: `include` (Include), `exclude` (Exclude)

Optional:

- `max_os_version` (String) Max OS version for Applicability Rule.
- `min_os_version` (String) Min OS version for Applicability Rule.
- `name` (String) Name for object.


<a id="nestedatt--windows10"></a>
### Nested Schema for `windows10`

Optional:

- `oma_settings` (Attributes Set) OMA settings. This collection can contain a maximum of 1000 elements. / OMA Settings definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omasetting?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--windows10--oma_settings))

<a id="nestedatt--windows10--oma_settings"></a>
### Nested Schema for `windows10.oma_settings`

Required:

- `display_name` (String) Display Name.
- `oma_uri` (String) OMA.

Optional:

- `base64` (Attributes) OMA Settings Base64 definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omasettingbase64?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows10--oma_settings--base64))
- `boolean` (Attributes) OMA Settings Boolean definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omasettingboolean?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows10--oma_settings--boolean))
- `date_time` (Attributes) OMA Settings DateTime definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omasettingdatetime?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows10--oma_settings--date_time))
- `description` (String) Description.
- `floating_point` (Attributes) OMA Settings Floating Point definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omasettingfloatingpoint?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows10--oma_settings--floating_point))
- `integer` (Attributes) OMA Settings Integer definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omasettinginteger?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows10--oma_settings--integer))
- `string` (Attributes) OMA Settings String definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omasettingstring?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows10--oma_settings--string))
- `string_xml` (Attributes) OMA Settings StringXML definition. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-omasettingstringxml?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows10--oma_settings--string_xml))

<a id="nestedatt--windows10--oma_settings--base64"></a>
### Nested Schema for `windows10.oma_settings.base64`

Required:

- `value_base64` (String) Value. (Base64 encoded string)

Optional:

- `file_name` (String) File name associated with the Value property (*.cer


<a id="nestedatt--windows10--oma_settings--boolean"></a>
### Nested Schema for `windows10.oma_settings.boolean`

Required:

- `value` (Boolean) Value.


<a id="nestedatt--windows10--oma_settings--date_time"></a>
### Nested Schema for `windows10.oma_settings.date_time`

Required:

- `value` (String) Value.


<a id="nestedatt--windows10--oma_settings--floating_point"></a>
### Nested Schema for `windows10.oma_settings.floating_point`

Required:

- `value` (Number) Value.


<a id="nestedatt--windows10--oma_settings--integer"></a>
### Nested Schema for `windows10.oma_settings.integer`

Required:

- `value` (Number) Value.


<a id="nestedatt--windows10--oma_settings--string"></a>
### Nested Schema for `windows10.oma_settings.string`

Required:

- `value` (String) Value.


<a id="nestedatt--windows10--oma_settings--string_xml"></a>
### Nested Schema for `windows10.oma_settings.string_xml`

Required:

- `value` (String) Value. (UTF8 encoded byte array)

Optional:

- `file_name` (String) File name associated with the Value property (*.xml).
