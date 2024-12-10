---
page_title: "microsoft365wp_device_enrollment_configuration Resource - microsoft365wp"
subcategory: "MS Graph: Onboarding"
---

# microsoft365wp_device_enrollment_configuration (Resource)

The Base Class of Device Enrollment Configuration / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-deviceenrollmentconfiguration?view=graph-rest-beta

## Provider Notes

### Default Enrollment Configurations
To target _default_ enrollment configurations instead of creating new ones, `display_name` **must** be set to
`###TfWorkplaceDefault###` and `priority` **must** be set to `0`. `description`, `role_scope_tag_ids` and
`assignments` should not be set (their values will be ignored).  
The provider will then automatically import the existing default config from MS Graph on resource creation
instead of creating a new MS Graph entity. On destruction, the provider will only remove the resource from
Terraform state but not try to delete the default config entity from MS Graph.  
To compare your own Terraform config to the actual values of a default configuration in MS Graph, you can still
import the entity into Terraform state as can be done with any other resource (e.g.
`terraform import microsoft365wp_device_enrollment_configuration.windows_hello_for_business_default 01234567-89ab-cdef-0123-456789abcdef_DefaultWindowsHelloForBusiness` -
get the real `id` directly from MS Graph using e.g. some browser's development tools). Then run `terraform plan` as
usual to see the required changes.  

### Priorities
MS Graph seems to be quite picky about setting priorities. They have to be sequential (i.e. there cannot be gaps
when creating) and when updating priorities, values can only be switched with another exiting config (e.g. if
the config with priority 1 gets deleted (e.g. manually) then the list will start with priority 2 and the
`setPriority` action will refuse to set any other config to priority 1 anymore but only to priority 2).  
Therefore setting of configuration priorities has been implemented as best-effort only: Any errors returned from
MS Graph during the `setPriority` action will only be shown as warnings in Terraform and more than one apply
action may be required to actually get every priority set correctly.  
And in case of a situation that MS Graph is not willing (or able) to fulfill, there will continously remain
changes. If you are really, really sure that your configuration should actually be fulfillable by MS Graph, then
it may help to destroy all configurations again and start over.  

### Delays Between Writes
As there have been situations during testing that resulted in multiple enrollment configurations with duplicate
priority values (which MS Graph was unable to correct or reorder ever again), all write operations will incur a
small delay to ensure individual priority values.

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


resource "microsoft365wp_device_enrollment_configuration" "windows10_esp" {
  display_name = "TF Test Windows 10 ESP"

  windows10_esp = {
    allow_device_reset_on_install_failure = true
  }

  assignments = [
    { target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } },
  ]
}


resource "microsoft365wp_device_enrollment_configuration" "device_enrollment_limit_default" {
  display_name = "###TfWorkplaceDefault###" # mandatory to target default config
  priority     = 0                          # mandatory to target default config

  device_enrollment_limit = {
    limit = 11
  }
}

resource "microsoft365wp_device_enrollment_configuration" "device_enrollment_limit_1" {
  display_name = "TF Test Enrollment Limit 1"

  device_enrollment_limit = {
    limit = 3
  }

  assignments = [
    { target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } },
  ]

  priority = 1
}

resource "microsoft365wp_device_enrollment_configuration" "device_enrollment_limit_2" {
  display_name = "TF Test Enrollment Limit 2"

  device_enrollment_limit = {
    limit = 4
  }

  priority = 2
}

resource "microsoft365wp_device_enrollment_configuration" "device_enrollment_limit_3" {
  display_name = "TF Test Enrollment Limit 3"

  device_enrollment_limit = {
    limit = 5
  }

  priority = 3
}


resource "microsoft365wp_device_enrollment_configuration" "platform_restrictions_default" {
  display_name = "###TfWorkplaceDefault###" # mandatory to target default config
  priority     = 0                          # mandatory to target default config

  platform_restrictions = {
    android_for_work_restriction = {
      platform_blocked = true
    }
  }
}


resource "microsoft365wp_device_enrollment_configuration" "single_platform_restriction" {
  display_name = "TF Test Platform Restriction"

  single_platform_restriction = {
    platform_type = "windows"
    platform_restriction = {
      personal_device_enrollment_blocked = true
      os_minimum_version                 = "10.0.19044"
    }
  }
}


resource "microsoft365wp_device_enrollment_configuration" "windows_hello_for_business_default" {
  display_name = "###TfWorkplaceDefault###" # mandatory to target default config
  priority     = 0                          # mandatory to target default config

  windows_hello_for_business = {
    pin_maximum_length = 126
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `assignments` (Attributes Set) The list of assignments. (see [below for nested schema](#nestedatt--assignments))
- `description` (String) The description of the device enrollment configuration. The _provider_ default value is `""`.
- `device_enrollment_limit` (Attributes) Device Enrollment Configuration that restricts the number of devices a user can enroll / https://learn.microsoft.com/en-us/graph/api/resources/intune-onboarding-deviceenrollmentlimitconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--device_enrollment_limit))
- `display_name` (String) The display name of the device enrollment configuration
- `platform_restrictions` (Attributes) Device Enrollment Configuration that restricts the types of devices a user can enroll / https://learn.microsoft.com/en-us/graph/api/resources/intune-onboarding-deviceenrollmentplatformrestrictionsconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--platform_restrictions))
- `priority` (Number) Priority is used when a user exists in multiple groups that are assigned enrollment configuration. Users are subject only to the configuration with the lowest priority value.
- `role_scope_tag_ids` (Set of String) . The _provider_ default value is `["0"]`.
- `single_platform_restriction` (Attributes) Device Enrollment Configuration that restricts the types of devices a user can enroll for a single platform / https://learn.microsoft.com/en-us/graph/api/resources/intune-onboarding-deviceenrollmentplatformrestrictionconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--single_platform_restriction))
- `windows10_esp` (Attributes) Windows 10 Enrollment Status Page Configuration / https://learn.microsoft.com/en-us/graph/api/resources/intune-onboarding-windows10enrollmentcompletionpageconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows10_esp))
- `windows_hello_for_business` (Attributes) Windows Hello for Business settings lets users access their devices using a gesture, such as biometric authentication, or a PIN. Configure settings for enrolled Windows 10, Windows 10 Mobile and later. / https://learn.microsoft.com/en-us/graph/api/resources/intune-onboarding-deviceenrollmentwindowshelloforbusinessconfiguration?view=graph-rest-beta (see [below for nested schema](#nestedatt--windows_hello_for_business))

### Read-Only

- `created_date_time` (String) Created date time in UTC of the device enrollment configuration
- `id` (String) Unique Identifier for the account
- `last_modified_date_time` (String) Last modified date time in UTC of the device enrollment configuration
- `version` (Number) The version of the device enrollment configuration

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




<a id="nestedatt--device_enrollment_limit"></a>
### Nested Schema for `device_enrollment_limit`

Required:

- `limit` (Number) The maximum number of devices that a user can enroll


<a id="nestedatt--platform_restrictions"></a>
### Nested Schema for `platform_restrictions`

Optional:

- `android_for_work_restriction` (Attributes) Android for work restrictions based on platform, platform operating system version, and device ownership / Platform specific enrollment restrictions / https://learn.microsoft.com/en-us/graph/api/resources/intune-onboarding-deviceenrollmentplatformrestriction?view=graph-rest-beta. The _provider_ default value is `{}`. (see [below for nested schema](#nestedatt--platform_restrictions--android_for_work_restriction))
- `android_restriction` (Attributes) Android restrictions based on platform, platform operating system version, and device ownership / Platform specific enrollment restrictions / https://learn.microsoft.com/en-us/graph/api/resources/intune-onboarding-deviceenrollmentplatformrestriction?view=graph-rest-beta. The _provider_ default value is `{}`. (see [below for nested schema](#nestedatt--platform_restrictions--android_restriction))
- `ios_restriction` (Attributes) Ios restrictions based on platform, platform operating system version, and device ownership / Platform specific enrollment restrictions / https://learn.microsoft.com/en-us/graph/api/resources/intune-onboarding-deviceenrollmentplatformrestriction?view=graph-rest-beta. The _provider_ default value is `{}`. (see [below for nested schema](#nestedatt--platform_restrictions--ios_restriction))
- `macos_restriction` (Attributes) Mac restrictions based on platform, platform operating system version, and device ownership / Platform specific enrollment restrictions / https://learn.microsoft.com/en-us/graph/api/resources/intune-onboarding-deviceenrollmentplatformrestriction?view=graph-rest-beta. The _provider_ default value is `{}`. (see [below for nested schema](#nestedatt--platform_restrictions--macos_restriction))
- `windows_restriction` (Attributes) Windows restrictions based on platform, platform operating system version, and device ownership / Platform specific enrollment restrictions / https://learn.microsoft.com/en-us/graph/api/resources/intune-onboarding-deviceenrollmentplatformrestriction?view=graph-rest-beta. The _provider_ default value is `{}`. (see [below for nested schema](#nestedatt--platform_restrictions--windows_restriction))

<a id="nestedatt--platform_restrictions--android_for_work_restriction"></a>
### Nested Schema for `platform_restrictions.android_for_work_restriction`

Optional:

- `blocked_manufacturers` (Set of String) Collection of blocked Manufacturers. The _provider_ default value is `[]`.
- `blocked_skus` (Set of String) Collection of blocked Skus. The _provider_ default value is `[]`.
- `os_maximum_version` (String) Max OS version supported. The _provider_ default value is `""`.
- `os_minimum_version` (String) Min OS version supported. The _provider_ default value is `""`.
- `personal_device_enrollment_blocked` (Boolean) Block personally owned devices from enrolling. The _provider_ default value is `false`.
- `platform_blocked` (Boolean) Block the platform from enrolling. The _provider_ default value is `false`.


<a id="nestedatt--platform_restrictions--android_restriction"></a>
### Nested Schema for `platform_restrictions.android_restriction`

Optional:

- `blocked_manufacturers` (Set of String) Collection of blocked Manufacturers. The _provider_ default value is `[]`.
- `blocked_skus` (Set of String) Collection of blocked Skus. The _provider_ default value is `[]`.
- `os_maximum_version` (String) Max OS version supported. The _provider_ default value is `""`.
- `os_minimum_version` (String) Min OS version supported. The _provider_ default value is `""`.
- `personal_device_enrollment_blocked` (Boolean) Block personally owned devices from enrolling. The _provider_ default value is `false`.
- `platform_blocked` (Boolean) Block the platform from enrolling. The _provider_ default value is `false`.


<a id="nestedatt--platform_restrictions--ios_restriction"></a>
### Nested Schema for `platform_restrictions.ios_restriction`

Optional:

- `blocked_manufacturers` (Set of String) Collection of blocked Manufacturers. The _provider_ default value is `[]`.
- `blocked_skus` (Set of String) Collection of blocked Skus. The _provider_ default value is `[]`.
- `os_maximum_version` (String) Max OS version supported. The _provider_ default value is `""`.
- `os_minimum_version` (String) Min OS version supported. The _provider_ default value is `""`.
- `personal_device_enrollment_blocked` (Boolean) Block personally owned devices from enrolling. The _provider_ default value is `false`.
- `platform_blocked` (Boolean) Block the platform from enrolling. The _provider_ default value is `false`.


<a id="nestedatt--platform_restrictions--macos_restriction"></a>
### Nested Schema for `platform_restrictions.macos_restriction`

Optional:

- `blocked_manufacturers` (Set of String) Collection of blocked Manufacturers. The _provider_ default value is `[]`.
- `blocked_skus` (Set of String) Collection of blocked Skus. The _provider_ default value is `[]`.
- `os_maximum_version` (String) Max OS version supported. The _provider_ default value is `""`.
- `os_minimum_version` (String) Min OS version supported. The _provider_ default value is `""`.
- `personal_device_enrollment_blocked` (Boolean) Block personally owned devices from enrolling. The _provider_ default value is `false`.
- `platform_blocked` (Boolean) Block the platform from enrolling. The _provider_ default value is `false`.


<a id="nestedatt--platform_restrictions--windows_restriction"></a>
### Nested Schema for `platform_restrictions.windows_restriction`

Optional:

- `blocked_manufacturers` (Set of String) Collection of blocked Manufacturers. The _provider_ default value is `[]`.
- `blocked_skus` (Set of String) Collection of blocked Skus. The _provider_ default value is `[]`.
- `os_maximum_version` (String) Max OS version supported. The _provider_ default value is `""`.
- `os_minimum_version` (String) Min OS version supported. The _provider_ default value is `""`.
- `personal_device_enrollment_blocked` (Boolean) Block personally owned devices from enrolling. The _provider_ default value is `false`.
- `platform_blocked` (Boolean) Block the platform from enrolling. The _provider_ default value is `false`.



<a id="nestedatt--single_platform_restriction"></a>
### Nested Schema for `single_platform_restriction`

Required:

- `platform_restriction` (Attributes) Restrictions based on platform, platform operating system version, and device ownership / Platform specific enrollment restrictions / https://learn.microsoft.com/en-us/graph/api/resources/intune-onboarding-deviceenrollmentplatformrestriction?view=graph-rest-beta (see [below for nested schema](#nestedatt--single_platform_restriction--platform_restriction))
- `platform_type` (String) Type of platform for which this restriction applies. / This enum indicates the platform type for which the enrollment restriction applies; possible values are: `allPlatforms` (Indicates that the enrollment configuration applies to all platforms), `ios` (Indicates that the enrollment configuration applies only to iOS/iPadOS devices), `windows` (Indicates that the enrollment configuration applies only to Windows devices), `windowsPhone` (Indicates that the enrollment configuration applies only to Windows Phone devices), `android` (Indicates that the enrollment configuration applies only to Android devices), `androidForWork` (Indicates that the enrollment configuration applies only to Android for Work devices), `mac` (Indicates that the enrollment configuration applies only to macOS devices), `linux` (Indicates that the enrollment configuration applies only to Linux devices), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use)

<a id="nestedatt--single_platform_restriction--platform_restriction"></a>
### Nested Schema for `single_platform_restriction.platform_restriction`

Optional:

- `blocked_manufacturers` (Set of String) Collection of blocked Manufacturers. The _provider_ default value is `[]`.
- `blocked_skus` (Set of String) Collection of blocked Skus. The _provider_ default value is `[]`.
- `os_maximum_version` (String) Max OS version supported. The _provider_ default value is `""`.
- `os_minimum_version` (String) Min OS version supported. The _provider_ default value is `""`.
- `personal_device_enrollment_blocked` (Boolean) Block personally owned devices from enrolling. The _provider_ default value is `false`.
- `platform_blocked` (Boolean) Block the platform from enrolling. The _provider_ default value is `false`.



<a id="nestedatt--windows10_esp"></a>
### Nested Schema for `windows10_esp`

Optional:

- `allow_device_reset_on_install_failure` (Boolean) When TRUE, allows device reset on installation failure. When false, reset is blocked. The default is false. The _provider_ default value is `false`.
- `allow_device_use_on_install_failure` (Boolean) When TRUE, allows the user to continue using the device on installation failure. When false, blocks the user on installation failure. The default is false. The _provider_ default value is `false`.
- `allow_log_collection_on_install_failure` (Boolean) When TRUE, allows log collection on installation failure. When false, log collection is not allowed. The default is false. The _provider_ default value is `true`.
- `block_device_setup_retry_by_user` (Boolean) When TRUE, blocks user from retrying the setup on installation failure. When false, user is allowed to retry. The default is false. The _provider_ default value is `false`.
- `custom_error_message` (String) The custom error message to show upon installation failure. Max length is 10000. example: "Setup could not be completed. Please try again or contact your support person for help.". The _provider_ default value is `"Setup could not be completed. Please try again or contact your support person for help."`.
- `disable_user_status_tracking_after_first_user` (Boolean) When TRUE, disables showing installation progress for first user post enrollment. When false, enables showing progress. The default is false. The _provider_ default value is `true`.
- `install_progress_timeout_in_minutes` (Number) The installation progress timeout in minutes. Default is 60 minutes. The _provider_ default value is `60`.
- `selected_mobile_app_ids` (Set of String) Selected applications to track the installation status. It is in the form of an array of GUIDs. The _provider_ default value is `[]`.
- `show_installation_progress` (Boolean) When TRUE, shows installation progress to user. When false, hides installation progress. The default is false. The _provider_ default value is `true`.
- `track_install_progress_for_autopilot_only` (Boolean) When TRUE, installation progress is tracked for only Autopilot enrollment scenarios. When false, other scenarios are tracked as well. The default is false. The _provider_ default value is `true`.


<a id="nestedatt--windows_hello_for_business"></a>
### Nested Schema for `windows_hello_for_business`

Optional:

- `enhanced_biometrics_state` (String) Controls the ability to use the anti-spoofing features for facial recognition on devices which support it. If set to disabled, anti-spoofing features are not allowed. If set to Not Configured, the user can choose whether they want to use anti-spoofing. / Possible values of a property; possible values are: `notConfigured` (Device default value, no intent.), `enabled` (Enables the setting on the device.), `disabled` (Disables the setting on the device.). The _provider_ default value is `"notConfigured"`.
- `enhanced_sign_in_security` (Number) Setting to configure Enhanced sign-in security. Default is Not Configured. The _provider_ default value is `0`.
- `pin_expiration_in_days` (Number) Controls the period of time (in days) that a PIN can be used before the system requires the user to change it. This must be set between 0 and 730, inclusive. If set to 0, the user's PIN will never expire. The _provider_ default value is `0`.
- `pin_lowercase_characters_usage` (String) Controls the ability to use lowercase letters in the Windows Hello for Business PIN.  Allowed permits the use of lowercase letter(s), whereas Required ensures they are present. If set to Not Allowed, lowercase letters will not be permitted. / Windows Hello for Business pin usage options; possible values are: `allowed` (Allowed the usage of certain pin rule), `required` (Enforce the usage of certain pin rule), `disallowed` (Forbit the usage of certain pin rule). The _provider_ default value is `"disallowed"`.
- `pin_maximum_length` (Number) Controls the maximum number of characters allowed for the Windows Hello for Business PIN. This value must be between 4 and 127, inclusive. This value must be greater than or equal to the value set for the minimum PIN. The _provider_ default value is `127`.
- `pin_minimum_length` (Number) Controls the minimum number of characters required for the Windows Hello for Business PIN.  This value must be between 4 and 127, inclusive, and less than or equal to the value set for the maximum PIN. The _provider_ default value is `4`.
- `pin_previous_block_count` (Number) Controls the ability to prevent users from using past PINs. This must be set between 0 and 50, inclusive, and the current PIN of the user is included in that count. If set to 0, previous PINs are not stored. PIN history is not preserved through a PIN reset. The _provider_ default value is `0`.
- `pin_special_characters_usage` (String) Controls the ability to use special characters in the Windows Hello for Business PIN.  Allowed permits the use of special character(s), whereas Required ensures they are present. If set to Not Allowed, special character(s) will not be permitted. / Windows Hello for Business pin usage options; possible values are: `allowed` (Allowed the usage of certain pin rule), `required` (Enforce the usage of certain pin rule), `disallowed` (Forbit the usage of certain pin rule). The _provider_ default value is `"disallowed"`.
- `pin_uppercase_characters_usage` (String) Controls the ability to use uppercase letters in the Windows Hello for Business PIN.  Allowed permits the use of uppercase letter(s), whereas Required ensures they are present. If set to Not Allowed, uppercase letters will not be permitted. / Windows Hello for Business pin usage options; possible values are: `allowed` (Allowed the usage of certain pin rule), `required` (Enforce the usage of certain pin rule), `disallowed` (Forbit the usage of certain pin rule). The _provider_ default value is `"disallowed"`.
- `remote_passport_enabled` (Boolean) Controls the use of Remote Windows Hello for Business. Remote Windows Hello for Business provides the ability for a portable, registered device to be usable as a companion for desktop authentication. The desktop must be Azure AD joined and the companion device must have a Windows Hello for Business PIN. The _provider_ default value is `true`.
- `security_device_required` (Boolean) Controls whether to require a Trusted Platform Module (TPM) for provisioning Windows Hello for Business. A TPM provides an additional security benefit in that data stored on it cannot be used on other devices. If set to False, all devices can provision Windows Hello for Business even if there is not a usable TPM. The _provider_ default value is `false`.
- `security_key_for_sign_in` (String) Security key for Sign In provides the capacity for remotely turning ON/OFF Windows Hello Sercurity Keyl Not configured will honor configurations done on the clinet. / Possible values of a property; possible values are: `notConfigured` (Device default value, no intent.), `enabled` (Enables the setting on the device.), `disabled` (Disables the setting on the device.). The _provider_ default value is `"notConfigured"`.
- `state` (String) Controls whether to allow the device to be configured for Windows Hello for Business. If set to disabled, the user cannot provision Windows Hello for Business except on Azure Active Directory joined mobile phones if otherwise required. If set to Not Configured, Intune will not override client defaults. / Possible values of a property; possible values are: `notConfigured` (Device default value, no intent.), `enabled` (Enables the setting on the device.), `disabled` (Disables the setting on the device.). The _provider_ default value is `"enabled"`.
- `unlock_with_biometrics_enabled` (Boolean) Controls the use of biometric gestures, such as face and fingerprint, as an alternative to the Windows Hello for Business PIN.  If set to False, biometric gestures are not allowed. Users must still configure a PIN as a backup in case of failures. The _provider_ default value is `true`.
