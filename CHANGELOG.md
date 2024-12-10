# `microsoft365wp` Terraform Provider Changelog


## 0.15.0 (2024-12-10)

### added (9 changes)

- mobile_app: add resource
- Add provider function parse_intunewin_metadata
- Add provider function parse_appxmsix_metadata
- device_configuration: add edition_upgrade
- windows_driver_update_profile: add resource
- cloud_pc_*: adapt for current MS Graph behavior
- authentication_context_class_reference: add resource
- Update metadata: no relevant changes
- Docs: major overhaul and various enhancements

### fixed (3 changes)

- Fix handling of read-only attributes in nested derived types
- cloud_pc_provisioning_policy: fix region_name handling
- Fix handling of various base64 encoded attributes


## 0.14.0 (2024-10-30)

### added (2 changes)

- Update MS Graph metadata
- Update dependencies


## 0.13.0 (2024-09-11)

### added (4 changes)

- Docs: Include description and allow values for enum types
- Update graph metadata
- device_configuration: support android_work_profile_general_device
- device_configuration.ios_device_features: support ios_single_sign_on_extension

### changed (1 change)

- Convert warning to info if entity has vanished from graph

### fixed (1 change)

- Avoid state corruption when deleting...
