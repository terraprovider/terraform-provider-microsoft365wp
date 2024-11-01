---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "microsoft365wp_device_and_app_management_assignment_filter Data Source - microsoft365wp"
subcategory: ""
description: |-
  
---

# microsoft365wp_device_and_app_management_assignment_filter (Data Source)



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


data "microsoft365wp_device_and_app_management_assignment_filter" "one" {
  id = "8808bb33-aa9c-4da7-a3fb-b92a4be8e001"
}

output "microsoft365wp_device_and_app_management_assignment_filter" {
  value = data.microsoft365wp_device_and_app_management_assignment_filter.one
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `assignment_filter_management_type` (String) Indicates filter is applied to either 'devices' or 'apps' management type. Possible values are devices, apps. Default filter will be applied to 'devices' / Supported filter management types whether its devices or apps; possible values are: `devices` (Indicates when filter is supported based on device properties. This is the default value when management type resolution fails.), `apps` (Indicates when filter is supported based on app properties.), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use.)
- `created_date_time` (String) The creation time of the assignment filter. The value cannot be modified and is automatically populated during new assignment filter process. The timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 would look like this: '2014-01-01T00:00:00Z'.
- `description` (String) Optional description of the Assignment Filter.
- `display_name` (String) The name of the Assignment Filter.
- `id` (String) The ID of this resource.
- `last_modified_date_time` (String) Last modified time of the Assignment Filter. The timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 would look like this: '2014-01-01T00:00:00Z'
- `platform` (String) Indicates filter is applied to which flatform. Possible values are android,androidForWork,iOS,macOS,windowsPhone81,windows81AndLater,windows10AndLater,androidWorkProfile, unknown, androidAOSP, androidMobileApplicationManagement, iOSMobileApplicationManagement, windowsMobileApplicationManagement. Default filter will be applied to 'unknown'. / Supported platform types; possible values are: `android` (Android.), `androidForWork` (AndroidForWork.), `iOS` (iOS.), `macOS` (MacOS.), `windowsPhone81` (WindowsPhone 8.1.), `windows81AndLater` (Windows 8.1 and later), `windows10AndLater` (Windows 10 and later.), `androidWorkProfile` (Android Work Profile.), `unknown` (Unknown.), `androidAOSP` (Android AOSP.), `androidMobileApplicationManagement` (Indicates Mobile Application Management (MAM) for android devices.), `iOSMobileApplicationManagement` (Indicates Mobile Application Management (MAM) for iOS devices), `unknownFutureValue` (Evolvable enumeration sentinel value. Do not use), `windowsMobileApplicationManagement` (Indicates Mobile Application Management (MAM) for Windows devices.)
- `role_scope_tags` (Set of String) Indicates role scope tags assigned for the assignment filter.
- `rule` (String) Rule definition of the assignment filter.
