---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "microsoft365wp_android_managed_app_protection Resource - microsoft365wp"
subcategory: ""
description: |-
  Policy used to configure detailed management settings targeted to specific security groups and for a specified set of apps on an Android device / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-androidManagedAppProtection?view=graph-rest-beta
---

# microsoft365wp_android_managed_app_protection (Resource)

Policy used to configure detailed management settings targeted to specific security groups and for a specified set of apps on an Android device / https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-androidManagedAppProtection?view=graph-rest-beta

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


resource "microsoft365wp_android_managed_app_protection" "test1" {
  display_name = "TF Test Android 1"

  app_group_type = "allCoreMicrosoftApps"

  assignments = [
    { target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } },
  ]
}

resource "microsoft365wp_android_managed_app_protection" "test2" {
  display_name = "TF Test Android 2"

  app_group_type = "selectedPublicApps"
  apps = [
    {
      app_id = {
        android = {
          package_id = "com.adobe.reader"
        }
      }
    }
  ]

  allowed_outbound_data_transfer_destinations = "managedApps"
  exempted_app_packages = [
    {
      name  = "Name"
      value = "andValue"
    },
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `app_group_type` (String) Public Apps selection: group or individual / Indicates a collection of apps to target which can be one of several pre-defined lists of apps or a manually selected list of apps; possible values are: `selectedPublicApps` (Target the collection of apps manually selected by the admin.), `allCoreMicrosoftApps` (Target the core set of Microsoft apps (Office, Edge, etc).), `allMicrosoftApps` (Target all apps with Microsoft as publisher.), `allApps` (Target all apps with an available assignment.)
- `display_name` (String)

### Optional

- `allowed_android_device_manufacturers` (String) Semicolon seperated list of device manufacturers allowed, as a string, for the managed app to work.
- `allowed_android_device_models` (Set of String) List of device models allowed, as a string, for the managed app to work.
- `allowed_data_ingestion_locations` (Set of String) Locations which can be used to bring data into organization documents; possible values are: `oneDriveForBusiness` (OneDrive for business), `sharePoint` (SharePoint Online), `camera` (The device's camera), `photoLibrary` (The device's photo library)
- `allowed_data_storage_locations` (Set of String) Storage locations where managed apps can potentially store their data; possible values are: `oneDriveForBusiness` (OneDrive for business), `sharePoint` (SharePoint), `box` (Box), `localStorage` (Local storage on the device), `photoLibrary` (The device's photo library)
- `allowed_inbound_data_transfer_sources` (String) Data can be transferred from/to these classes of apps; possible values are: `allApps` (All apps.), `managedApps` (Managed apps.), `none` (No apps.)
- `allowed_outbound_clipboard_sharing_exception_length` (Number)
- `allowed_outbound_clipboard_sharing_level` (String) Represents the level to which the device's clipboard may be shared between apps; possible values are: `allApps` (Sharing is allowed between all apps, managed or not), `managedAppsWithPasteIn` (Sharing is allowed between all managed apps with paste in enabled), `managedApps` (Sharing is allowed between all managed apps), `blocked` (Sharing between apps is disabled)
- `allowed_outbound_data_transfer_destinations` (String) Data can be transferred from/to these classes of apps; possible values are: `allApps` (All apps.), `managedApps` (Managed apps.), `none` (No apps.)
- `app_action_if_account_is_clocked_out` (String) Defines a managed app behavior, either block or warn, if the user is clocked out (non-working time). / An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting)
- `app_action_if_android_device_manufacturer_not_allowed` (String) Defines a managed app behavior, either block or wipe, if the specified device manufacturer is not allowed. / An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting)
- `app_action_if_android_device_model_not_allowed` (String) Defines a managed app behavior, either block or wipe, if the specified device model is not allowed. / An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting)
- `app_action_if_android_safety_net_apps_verification_failed` (String) Defines a managed app behavior, either warn or block, if the specified Android App Verification requirement fails. / An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting)
- `app_action_if_android_safety_net_device_attestation_failed` (String) Defines a managed app behavior, either warn or block, if the specified Android SafetyNet Attestation requirement fails. / An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting)
- `app_action_if_device_compliance_required` (String) An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting)
- `app_action_if_device_lock_not_set` (String) Defines a managed app behavior, either warn, block or wipe, if the screen lock is required on android device but is not set. / An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting)
- `app_action_if_device_passcode_complexity_less_than_high` (String) If the device does not have a passcode of high complexity or higher, trigger the stored action. / An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting)
- `app_action_if_device_passcode_complexity_less_than_low` (String) If the device does not have a passcode of low complexity or higher, trigger the stored action. / An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting)
- `app_action_if_device_passcode_complexity_less_than_medium` (String) If the device does not have a passcode of medium complexity or higher, trigger the stored action. / An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting)
- `app_action_if_maximum_pin_retries_exceeded` (String) An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting)
- `app_action_if_samsung_knox_attestation_required` (String) Defines the behavior of a managed app when Samsung Knox Attestation is required. Possible values are null, warn, block & wipe. If the admin does not set this action, the default is null, which indicates this setting is not configured. / An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting)
- `app_action_if_unable_to_authenticate_user` (String) An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting)
- `approved_keyboards` (Attributes Set) If Keyboard Restriction is enabled, only keyboards in this approved list will be allowed. A key should be Android package id for a keyboard and value should be a friendly name / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyValuePair?view=graph-rest-beta (see [below for nested schema](#nestedatt--approved_keyboards))
- `apps` (Attributes Set) List of apps to which the policy is deployed. / The identifier for the deployment an app. (see [below for nested schema](#nestedatt--apps))
- `assignments` (Attributes Set) The list of assignments. (see [below for nested schema](#nestedatt--assignments))
- `biometric_authentication_blocked` (Boolean) Indicates whether use of the biometric authentication is allowed in place of a pin if PinRequired is set to True.
- `block_after_company_portal_update_deferral_in_days` (Number) Maximum number of days Company Portal update can be deferred on the device or app access will be blocked.
- `block_data_ingestion_into_organization_documents` (Boolean)
- `connect_to_vpn_on_launch` (Boolean) Whether the app should connect to the configured VPN on launch.
- `contact_sync_blocked` (Boolean)
- `custom_browser_display_name` (String) Friendly name of the preferred custom browser to open weblink on Android. When this property is configured, ManagedBrowserToOpenLinksRequired should be true.
- `custom_browser_package_id` (String) Unique identifier of the preferred custom browser to open weblink on Android. When this property is configured, ManagedBrowserToOpenLinksRequired should be true.
- `custom_dialer_app_display_name` (String) Friendly name of a custom dialer app to click-to-open a phone number on Android.
- `custom_dialer_app_package_id` (String) PackageId of a custom dialer app to click-to-open a phone number on Android.
- `data_backup_blocked` (Boolean)
- `description` (String)
- `device_compliance_required` (Boolean)
- `device_lock_required` (Boolean) Defines if any kind of lock must be required on android device
- `dialer_restriction_level` (String) The classes of apps that are allowed to click-to-open a phone number, for making phone calls or sending text messages; possible values are: `allApps` (Sharing is allowed to all apps.), `managedApps` (Sharing is allowed to all managed apps.), `customApp` (Sharing is allowed to a custom app.), `blocked` (Sharing between apps is blocked.)
- `disable_app_encryption_if_device_encryption_is_enabled` (Boolean) When this setting is enabled, app level encryption is disabled if device level encryption is enabled
- `disable_app_pin_if_device_pin_is_set` (Boolean)
- `encrypt_app_data` (Boolean) Indicates whether application data for managed apps should be encrypted
- `exempted_app_packages` (Attributes Set) App packages in this list will be exempt from the policy and will be able to receive data from managed apps. / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyValuePair?view=graph-rest-beta (see [below for nested schema](#nestedatt--exempted_app_packages))
- `fingerprint_and_biometric_enabled` (Boolean) If null, this setting will be ignored. If false both fingerprints and biometrics will not be enabled. If true, both fingerprints and biometrics will be enabled.
- `fingerprint_blocked` (Boolean)
- `grace_period_to_block_apps_during_off_clock_hours` (String)
- `keyboards_restricted` (Boolean) Indicates if keyboard restriction is enabled. If enabled list of approved keyboards must be provided as well.
- `managed_browser` (String) Type of managed browser; possible values are: `notConfigured` (Not configured), `microsoftEdge` (Microsoft Edge)
- `managed_browser_to_open_links_required` (Boolean)
- `maximum_allowed_device_threat_level` (String) The maxium threat level allowed for an app to be compliant; possible values are: `notConfigured` (Value not configured), `secured` (Device needs to have no threat), `low` (Device needs to have a low threat.), `medium` (Device needs to have not more than medium threat.), `high` (Device needs to have not more than high threat)
- `maximum_pin_retries` (Number)
- `maximum_required_os_version` (String)
- `maximum_warning_os_version` (String)
- `maximum_wipe_os_version` (String)
- `messaging_redirect_app_display_name` (String) When a specific app redirection is enforced by protectedMessagingRedirectAppType in an App Protection Policy, this value defines the app name which is allowed to be used.
- `messaging_redirect_app_package_id` (String) When a specific app redirection is enforced by protectedMessagingRedirectAppType in an App Protection Policy, this value defines the app package id which is allowed to be used.
- `minimum_pin_length` (Number)
- `minimum_required_app_version` (String)
- `minimum_required_company_portal_version` (String) Minimum version of the Company portal that must be installed on the device or app access will be blocked
- `minimum_required_os_version` (String)
- `minimum_required_patch_version` (String) Define the oldest required Android security patch level a user can have to gain secure access to the app.
- `minimum_warning_app_version` (String)
- `minimum_warning_company_portal_version` (String) Minimum version of the Company portal that must be installed on the device or the user will receive a warning
- `minimum_warning_os_version` (String)
- `minimum_warning_patch_version` (String) Define the oldest recommended Android security patch level a user can have for secure access to the app.
- `minimum_wipe_app_version` (String)
- `minimum_wipe_company_portal_version` (String) Minimum version of the Company portal that must be installed on the device or the company data on the app will be wiped
- `minimum_wipe_os_version` (String)
- `minimum_wipe_patch_version` (String) Android security patch level  less than or equal to the specified value will wipe the managed app and the associated company data.
- `mobile_threat_defense_partner_priority` (String) Determines the conflict resolution strategy, when more than one Mobile Threat Defense provider is enabled; possible values are: `defenderOverThirdPartyPartner` (Indicates use of Microsoft Defender Endpoint over 3rd party MTD connectors), `thirdPartyPartnerOverDefender` (Indicates use of a 3rd party MTD connector over Microsoft Defender Endpoint), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.)
- `mobile_threat_defense_remediation_action` (String) An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting)
- `notification_restriction` (String) Restrict managed app notification; possible values are: `allow` (Share all notifications.), `blockOrganizationalData` (Do not share Orgnizational data in notifications.), `block` (Do not share notifications.)
- `organizational_credentials_required` (Boolean)
- `period_before_pin_reset` (String)
- `period_offline_before_access_check` (String)
- `period_offline_before_wipe_is_enforced` (String)
- `period_online_before_access_check` (String)
- `pin_character_set` (String) Character set which is to be used for a user's app PIN; possible values are: `numeric` (Numeric characters), `alphanumericAndSymbol` (Alphanumeric and symbolic characters)
- `pin_required` (Boolean)
- `pin_required_instead_of_biometric_timeout` (String)
- `previous_pin_block_count` (Number)
- `print_blocked` (Boolean)
- `protected_messaging_redirect_app_type` (String) Defines how app messaging redirection is protected by an App Protection Policy. Default is anyApp; possible values are: `anyApp` (App protection policy will allow messaging redirection to any app.), `anyManagedApp` (App protection policy will allow messaging redirection to any managed application.), `specificApps` (App protection policy will allow messaging redirection only to specified applications in related App protection policy settings. See related settings 'messagingRedirectAppDisplayName', 'messagingRedirectAppPackageId' and 'messagingRedirectAppUrlScheme'.), `blocked` (App protection policy will block messaging redirection to any app.)
- `require_class3_biometrics` (Boolean) Require user to apply Class 3 Biometrics on their Android device.
- `require_pin_after_biometric_change` (Boolean) A PIN prompt will override biometric prompts if class 3 biometrics are updated on the device.
- `required_android_safety_net_apps_verification_type` (String) Defines the Android SafetyNet Apps Verification requirement for a managed app to work. / An admin enforced Android SafetyNet Device Attestation requirement on a managed app; possible values are: `none` (no requirement set), `enabled` (require that Android device has SafetyNet Apps Verification enabled)
- `required_android_safety_net_device_attestation_type` (String) Defines the Android SafetyNet Device Attestation requirement for a managed app to work. / An admin enforced Android SafetyNet Device Attestation requirement on a managed app; possible values are: `none` (no requirement set), `basicIntegrity` (require that Android device passes SafetyNet Basic Integrity validation), `basicIntegrityAndDeviceCertification` (require that Android device passes SafetyNet Basic Integrity and Device Certification validations)
- `required_android_safety_net_evaluation_type` (String) Defines the Android SafetyNet evaluation type requirement for a managed app to work. / An admin enforced Android SafetyNet evaluation type requirement on a managed app; possible values are: `basic` (Require basic evaluation), `hardwareBacked` (Require hardware backed evaluation)
- `role_scope_tag_ids` (Set of String)
- `save_as_blocked` (Boolean)
- `screen_capture_blocked` (Boolean) Indicates whether a managed user can take screen captures of managed apps
- `simple_pin_blocked` (Boolean)
- `targeted_app_management_levels` (String) The intended app management levels for this policy / Management levels for apps; possible values are: `unspecified` (Unspecified), `unmanaged` (Unmanaged), `mdm` (MDM), `androidEnterprise` (Android Enterprise), `androidEnterpriseDedicatedDevicesWithAzureAdSharedMode` (Android Enterprise dedicated devices with Azure AD Shared mode), `androidOpenSourceProjectUserAssociated` (Android Open Source Project (AOSP) devices), `androidOpenSourceProjectUserless` (Android Open Source Project (AOSP) userless devices), `unknownFutureValue` (Place holder for evolvable enum)
- `warn_after_company_portal_update_deferral_in_days` (Number) Maximum number of days Company Portal update can be deferred on the device or the user will receive the warning
- `wipe_after_company_portal_update_deferral_in_days` (Number) Maximum number of days Company Portal update can be deferred on the device or the company data on the app will be wiped

### Read-Only

- `created_date_time` (String)
- `deployed_app_count` (Number) Count of apps to which the current policy is deployed.
- `id` (String) The ID of this resource.
- `is_assigned` (Boolean) Indicates if the policy is deployed to any inclusion groups or not.
- `last_modified_date_time` (String)
- `version` (String)

<a id="nestedatt--approved_keyboards"></a>
### Nested Schema for `approved_keyboards`

Required:

- `name` (String)

Optional:

- `value` (String)


<a id="nestedatt--apps"></a>
### Nested Schema for `apps`

Required:

- `app_id` (Attributes) The identifier for an app with it's operating system type. / The identifier for a mobile app. (see [below for nested schema](#nestedatt--apps--app_id))

<a id="nestedatt--apps--app_id"></a>
### Nested Schema for `apps.app_id`

Optional:

- `android` (Attributes) The identifier for an Android app. (see [below for nested schema](#nestedatt--apps--app_id--android))

<a id="nestedatt--apps--app_id--android"></a>
### Nested Schema for `apps.app_id.android`

Required:

- `package_id` (String) The identifier for an app, as specified in the play store.




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




<a id="nestedatt--exempted_app_packages"></a>
### Nested Schema for `exempted_app_packages`

Required:

- `name` (String)

Optional:

- `value` (String)
