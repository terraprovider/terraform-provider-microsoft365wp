---
page_title: "microsoft365wp_ios_managed_app_protection Resource - microsoft365wp"
subcategory: "MS Graph: App management"
---

# microsoft365wp_ios_managed_app_protection (Resource)

Policy used to configure detailed management settings targeted to specific security groups and for a specified set of apps on an iOS device / https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-iosmanagedappprotection?view=graph-rest-beta

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


resource "microsoft365wp_ios_managed_app_protection" "test1" {
  display_name = "TF Test iOS 1"

  app_group_type = "allCoreMicrosoftApps"

  assignments = [
    { target = { group = { group_id = "298fded6-b252-4166-a473-f405e935f58d" } } },
  ]
}

resource "microsoft365wp_ios_managed_app_protection" "test2" {
  display_name = "TF Test iOS 2"

  app_group_type = "selectedPublicApps"
  apps = [
    {
      app_id = {
        ios = {
          bundle_id = "com.microsoft.msedge"
        }
      }
    }
  ]

  allowed_outbound_data_transfer_destinations = "managedApps"
  exempted_app_protocols = [
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
- `display_name` (String) Policy display name.

### Optional

- `allow_widget_content_sync` (Boolean) Indicates  if content sync for widgets is allowed for iOS on App Protection Policies. The _provider_ default value is `false`.
- `allowed_data_ingestion_locations` (Set of String) Data storage locations where a user may store managed data. / Locations which can be used to bring data into organization documents; possible values are: `oneDriveForBusiness` (OneDrive for business), `sharePoint` (SharePoint Online), `camera` (The device's camera), `photoLibrary` (The device's photo library). The _provider_ default value is `["oneDriveForBusiness","sharePoint","camera","photoLibrary"]`.
- `allowed_data_storage_locations` (Set of String) Data storage locations where a user may store managed data. / Storage locations where managed apps can potentially store their data; possible values are: `oneDriveForBusiness` (OneDrive for business), `sharePoint` (SharePoint), `box` (Box), `localStorage` (Local storage on the device), `photoLibrary` (The device's photo library). The _provider_ default value is `[]`.
- `allowed_inbound_data_transfer_sources` (String) Sources from which data is allowed to be transferred. / Data can be transferred from/to these classes of apps; possible values are: `allApps` (All apps.), `managedApps` (Managed apps.), `none` (No apps.). The _provider_ default value is `"allApps"`.
- `allowed_ios_device_models` (String) Semicolon seperated list of device models allowed, as a string, for the managed app to work.
- `allowed_outbound_clipboard_sharing_exception_length` (Number) Specify the number of characters that may be cut or copied from Org data and accounts to any application. This setting overrides the AllowedOutboundClipboardSharingLevel restriction. Default value of '0' means no exception is allowed. The _provider_ default value is `0`.
- `allowed_outbound_clipboard_sharing_level` (String) The level to which the clipboard may be shared between apps on the managed device. / Represents the level to which the device's clipboard may be shared between apps; possible values are: `allApps` (Sharing is allowed between all apps, managed or not), `managedAppsWithPasteIn` (Sharing is allowed between all managed apps with paste in enabled), `managedApps` (Sharing is allowed between all managed apps), `blocked` (Sharing between apps is disabled). The _provider_ default value is `"managedAppsWithPasteIn"`.
- `allowed_outbound_data_transfer_destinations` (String) Destinations to which data is allowed to be transferred. / Data can be transferred from/to these classes of apps; possible values are: `allApps` (All apps.), `managedApps` (Managed apps.), `none` (No apps.). The _provider_ default value is `"allApps"`.
- `app_action_if_account_is_clocked_out` (String) Defines a managed app behavior, either block or warn, if the user is clocked out (non-working time). / An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting). The _provider_ default value is `"block"`.
- `app_action_if_device_compliance_required` (String) Defines a managed app behavior, either block or wipe, when the device is either rooted or jailbroken, if DeviceComplianceRequired is set to true. / An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting). The _provider_ default value is `"block"`.
- `app_action_if_ios_device_model_not_allowed` (String) Defines a managed app behavior, either block or wipe, if the specified device model is not allowed. / An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting). The _provider_ default value is `"block"`.
- `app_action_if_maximum_pin_retries_exceeded` (String) Defines a managed app behavior, either block or wipe, based on maximum number of incorrect pin retry attempts. / An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting). The _provider_ default value is `"block"`.
- `app_action_if_unable_to_authenticate_user` (String) If set, it will specify what action to take in the case where the user is unable to checkin because their authentication token is invalid. This happens when the user is deleted or disabled in AAD. / An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting)
- `app_data_encryption_type` (String) Type of encryption which should be used for data in a managed app. / Represents the level to which app data is encrypted for managed apps; possible values are: `useDeviceSettings` (App data is encrypted based on the default settings on the device.), `afterDeviceRestart` (App data is encrypted when the device is restarted.), `whenDeviceLockedExceptOpenFiles` (App data associated with this policy is encrypted when the device is locked, except data in files that are open), `whenDeviceLocked` (App data associated with this policy is encrypted when the device is locked). The _provider_ default value is `"whenDeviceLocked"`.
- `apps` (Attributes Set) List of apps to which the policy is deployed. / The identifier for the deployment an app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-managedmobileapp?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--apps))
- `assignments` (Attributes Set) The list of assignments. (see [below for nested schema](#nestedatt--assignments))
- `block_data_ingestion_into_organization_documents` (Boolean) Indicates whether a user can bring data into org documents. The _provider_ default value is `false`.
- `contact_sync_blocked` (Boolean) Indicates whether contacts can be synced to the user's device. The _provider_ default value is `false`.
- `custom_browser_protocol` (String) A custom browser protocol to open weblink on iOS. When this property is configured, ManagedBrowserToOpenLinksRequired should be true.
- `custom_dialer_app_protocol` (String) Protocol of a custom dialer app to click-to-open a phone number on iOS, for example, skype:.
- `data_backup_blocked` (Boolean) Indicates whether the backup of a managed app's data is blocked. The _provider_ default value is `false`.
- `description` (String) The policy's description. The _provider_ default value is `""`.
- `device_compliance_required` (Boolean) Indicates whether device compliance is required. The _provider_ default value is `true`.
- `dialer_restriction_level` (String) The classes of dialer apps that are allowed to click-to-open a phone number. / The classes of apps that are allowed to click-to-open a phone number, for making phone calls or sending text messages; possible values are: `allApps` (Sharing is allowed to all apps.), `managedApps` (Sharing is allowed to all managed apps.), `customApp` (Sharing is allowed to a custom app.), `blocked` (Sharing between apps is blocked.). The _provider_ default value is `"allApps"`.
- `disable_app_pin_if_device_pin_is_set` (Boolean) Indicates whether use of the app pin is required if the device pin is set. The _provider_ default value is `false`.
- `disable_protection_of_managed_outbound_open_in_data` (Boolean) Disable protection of data transferred to other apps through IOS OpenIn option. This setting is only allowed to be True when AllowedOutboundDataTransferDestinations is set to ManagedApps. The _provider_ default value is `false`.
- `exempted_app_protocols` (Attributes Set) Apps in this list will be exempt from the policy and will be able to receive data from managed apps. / Key value pair for storing custom settings / https://learn.microsoft.com/en-us/graph/api/resources/intune-deviceconfig-keyvaluepair?view=graph-rest-beta. The _provider_ default value is `[]`. (see [below for nested schema](#nestedatt--exempted_app_protocols))
- `exempted_universal_links` (Set of String) A list of custom urls that are allowed to invocate an unmanaged app. The _provider_ default value is `[]`.
- `face_id_blocked` (Boolean) Indicates whether use of the FaceID is allowed in place of a pin if PinRequired is set to True. The _provider_ default value is `false`.
- `filter_open_in_to_only_managed_apps` (Boolean) Defines if open-in operation is supported from the managed app to the filesharing locations selected. This setting only applies when AllowedOutboundDataTransferDestinations is set to ManagedApps and DisableProtectionOfManagedOutboundOpenInData is set to False. The _provider_ default value is `false`.
- `fingerprint_blocked` (Boolean) Indicates whether use of the fingerprint reader is allowed in place of a pin if PinRequired is set to True. The _provider_ default value is `false`.
- `grace_period_to_block_apps_during_off_clock_hours` (String) A grace period before blocking app access during off clock hours.
- `managed_browser` (String) Indicates in which managed browser(s) that internet links should be opened. When this property is configured, ManagedBrowserToOpenLinksRequired should be true. / Type of managed browser; possible values are: `notConfigured` (Not configured), `microsoftEdge` (Microsoft Edge). The _provider_ default value is `"notConfigured"`.
- `managed_browser_to_open_links_required` (Boolean) Indicates whether internet links should be opened in the managed browser app, or any custom browser specified by CustomBrowserProtocol (for iOS) or CustomBrowserPackageId/CustomBrowserDisplayName (for Android). The _provider_ default value is `false`.
- `managed_universal_links` (Set of String) A list of custom urls that are allowed to invocate a managed app. The _provider_ default value is `[]`.
- `maximum_allowed_device_threat_level` (String) Maximum allowed device threat level, as reported by the MTD app / The maxium threat level allowed for an app to be compliant; possible values are: `notConfigured` (Value not configured), `secured` (Device needs to have no threat), `low` (Device needs to have a low threat.), `medium` (Device needs to have not more than medium threat.), `high` (Device needs to have not more than high threat). The _provider_ default value is `"notConfigured"`.
- `maximum_pin_retries` (Number) Maximum number of incorrect pin retry attempts before the managed app is either blocked or wiped. The _provider_ default value is `5`.
- `maximum_required_os_version` (String) Versions bigger than the specified version will block the managed app from accessing company data.
- `maximum_warning_os_version` (String) Versions bigger than the specified version will block the managed app from accessing company data.
- `maximum_wipe_os_version` (String) Versions bigger than the specified version will block the managed app from accessing company data.
- `messaging_redirect_app_url_scheme` (String) When a specific app redirection is enforced by protectedMessagingRedirectAppType in an App Protection Policy, this value defines the app url redirect schemes which are allowed to be used.
- `minimum_pin_length` (Number) Minimum pin length required for an app-level pin if PinRequired is set to True. The _provider_ default value is `4`.
- `minimum_required_app_version` (String) Versions less than the specified version will block the managed app from accessing company data.
- `minimum_required_os_version` (String) Versions less than the specified version will block the managed app from accessing company data.
- `minimum_required_sdk_version` (String) Versions less than the specified version will block the managed app from accessing company data.
- `minimum_warning_app_version` (String) Versions less than the specified version will result in warning message on the managed app.
- `minimum_warning_os_version` (String) Versions less than the specified version will result in warning message on the managed app from accessing company data.
- `minimum_warning_sdk_version` (String) Versions less than the specified version will result in warning message on the managed app from accessing company data.
- `minimum_wipe_app_version` (String) Versions less than or equal to the specified version will wipe the managed app and the associated company data.
- `minimum_wipe_os_version` (String) Versions less than or equal to the specified version will wipe the managed app and the associated company data.
- `minimum_wipe_sdk_version` (String) Versions less than the specified version will block the managed app from accessing company data.
- `mobile_threat_defense_partner_priority` (String) Indicates how to prioritize which Mobile Threat Defense (MTD) partner is enabled for a given platform, when more than one is enabled. An app can only be actively using a single Mobile Threat Defense partner. When NULL, Microsoft Defender will be given preference. Otherwise setting the value to defenderOverThirdPartyPartner or thirdPartyPartnerOverDefender will make explicit which partner to prioritize. Possible values are: null, defenderOverThirdPartyPartner, thirdPartyPartnerOverDefender and unknownFutureValue. Default value is null / Determines the conflict resolution strategy, when more than one Mobile Threat Defense provider is enabled; possible values are: `defenderOverThirdPartyPartner` (Indicates use of Microsoft Defender Endpoint over 3rd party MTD connectors), `thirdPartyPartnerOverDefender` (Indicates use of a 3rd party MTD connector over Microsoft Defender Endpoint), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.)
- `mobile_threat_defense_remediation_action` (String) Determines what action to take if the mobile threat defense threat threshold isn't met. Warn isn't a supported value for this property / An admin initiated action to be applied on a managed app; possible values are: `block` (Indicates the user will be blocked from accessing the app and corporate data), `wipe` (Indicates the corporate data will be removed from the app), `warn` (Indicates user will be warned the when accessing the app), `blockWhenSettingIsSupported` (Indicates user will be blocked from accessing the app and corporate data if devices supports this setting). The _provider_ default value is `"block"`.
- `notification_restriction` (String) Specify app notification restriction / Restrict managed app notification; possible values are: `allow` (Share all notifications.), `blockOrganizationalData` (Do not share Orgnizational data in notifications.), `block` (Do not share notifications.). The _provider_ default value is `"allow"`.
- `organizational_credentials_required` (Boolean) Indicates whether organizational credentials are required for app use. The _provider_ default value is `false`.
- `period_before_pin_reset` (String) TimePeriod before the all-level pin must be reset if PinRequired is set to True. The _provider_ default value is `"PT0S"`.
- `period_offline_before_access_check` (String) The period after which access is checked when the device is not connected to the internet. The _provider_ default value is `"PT12H"`.
- `period_offline_before_wipe_is_enforced` (String) The amount of time an app is allowed to remain disconnected from the internet before all managed data it is wiped. The _provider_ default value is `"P90D"`.
- `period_online_before_access_check` (String) The period after which access is checked when the device is connected to the internet. The _provider_ default value is `"PT30M"`.
- `pin_character_set` (String) Character set which may be used for an app-level pin if PinRequired is set to True. / Character set which is to be used for a user's app PIN; possible values are: `numeric` (Numeric characters), `alphanumericAndSymbol` (Alphanumeric and symbolic characters). The _provider_ default value is `"numeric"`.
- `pin_required` (Boolean) Indicates whether an app-level pin is required. The _provider_ default value is `true`.
- `pin_required_instead_of_biometric_timeout` (String) Timeout in minutes for an app pin instead of non biometrics passcode. The _provider_ default value is `"PT30M"`.
- `previous_pin_block_count` (Number) Requires a pin to be unique from the number specified in this property. The _provider_ default value is `0`.
- `print_blocked` (Boolean) Indicates whether printing is allowed from managed apps. The _provider_ default value is `false`.
- `protect_inbound_data_from_unknown_sources` (Boolean) Protect incoming data from unknown source. This setting is only allowed to be True when AllowedInboundDataTransferSources is set to AllApps. The _provider_ default value is `false`.
- `protected_messaging_redirect_app_type` (String) Defines how app messaging redirection is protected by an App Protection Policy. Default is anyApp. / Defines how app messaging redirection is protected by an App Protection Policy. Default is anyApp; possible values are: `anyApp` (App protection policy will allow messaging redirection to any app.), `anyManagedApp` (App protection policy will allow messaging redirection to any managed application.), `specificApps` (App protection policy will allow messaging redirection only to specified applications in related App protection policy settings. See related settings `messagingRedirectAppDisplayName`, `messagingRedirectAppPackageId` and `messagingRedirectAppUrlScheme`.), `blocked` (App protection policy will block messaging redirection to any app.). The _provider_ default value is `"anyApp"`.
- `role_scope_tag_ids` (Set of String) List of Scope Tags for this Entity instance. The _provider_ default value is `["0"]`.
- `save_as_blocked` (Boolean) Indicates whether users may use the "Save As" menu item to save a copy of protected files. The _provider_ default value is `false`.
- `simple_pin_blocked` (Boolean) Indicates whether simplePin is blocked. The _provider_ default value is `false`.
- `targeted_app_management_levels` (String) The intended app management levels for this policy / Management levels for apps; possible values are: `unspecified` (Unspecified), `unmanaged` (Unmanaged), `mdm` (MDM), `androidEnterprise` (Android Enterprise), `androidEnterpriseDedicatedDevicesWithAzureAdSharedMode` (Android Enterprise dedicated devices with Azure AD Shared mode), `androidOpenSourceProjectUserAssociated` (Android Open Source Project (AOSP) devices), `androidOpenSourceProjectUserless` (Android Open Source Project (AOSP) userless devices), `unknownFutureValue` (Place holder for evolvable enum). The _provider_ default value is `"unspecified"`.
- `third_party_keyboards_blocked` (Boolean) Defines if third party keyboards are allowed while accessing a managed app. The _provider_ default value is `false`.

### Read-Only

- `created_date_time` (String) The date and time the policy was created.
- `deployed_app_count` (Number) Count of apps to which the current policy is deployed.
- `id` (String) Key of the entity.
- `is_assigned` (Boolean) Indicates if the policy is deployed to any inclusion groups or not.
- `last_modified_date_time` (String) Last time the policy was modified.
- `version` (String) Version of the entity.

<a id="nestedatt--apps"></a>
### Nested Schema for `apps`

Required:

- `app_id` (Attributes) The identifier for an app with it's operating system type. / The identifier for a mobile app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-mobileappidentifier?view=graph-rest-beta (see [below for nested schema](#nestedatt--apps--app_id))

<a id="nestedatt--apps--app_id"></a>
### Nested Schema for `apps.app_id`

Optional:

- `ios` (Attributes) The identifier for an iOS app. / https://learn.microsoft.com/en-us/graph/api/resources/intune-mam-iosmobileappidentifier?view=graph-rest-beta (see [below for nested schema](#nestedatt--apps--app_id--ios))

<a id="nestedatt--apps--app_id--ios"></a>
### Nested Schema for `apps.app_id.ios`

Required:

- `bundle_id` (String) The identifier for an app, as specified in the app store.




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




<a id="nestedatt--exempted_app_protocols"></a>
### Nested Schema for `exempted_app_protocols`

Required:

- `name` (String) Name for this key-value pair

Optional:

- `value` (String) Value for this key-value pair
