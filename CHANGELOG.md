# `microsoft365wp` Terraform Provider Changelog


## 0.17.4 (2025-12-03)

### changed

- conditional_access_policy: update


## 0.17.3 (2025-12-02)

### added

- administrative_unit

### changed

- Update dependencies


## 0.17.2 (2025-11-20)

### fixed

- unified_role_management_policy: fix unknown value after create

### changed

- Update dependencies


## 0.17.1 (2025-10-21)

### fixed

- authentication_methods_policy: ignore VerifiableCredentials
- unified_role_management_policy: fix invisible state drift

### changed

- Update dependencies


## 0.17.0 (2025-07-08)

### added

- application (data sources only)
- attribute_set
- cross_tenant_identity_sync_policy_partner
- custom_security_attribute_definition
- device_compliance_script
- external_identities_policy
- group (data sources only)
- group_assigned_license
- intune_branding_profile
- service_principal (data sources only)
- subscribed_sku (data sources only)
- user (data sources only)

### fixed

- assignmentTarget: fix filter_id
- azure_ad_windows_autopilot_deployment_profile: fix locale="" on create

### changed

- authentication_methods_policy: convert null from MS Graph to "migrationComplete"
- Update dependencies
- Update MS Graph metadata & docs

### breaking

- authentication_methods_policy: add qr_code_pin


## 0.16.0 (2025-03-12)

### added

- authentication_combination_configuration
- authentication_flows_policy
- authorization_policy
- cloud_pc_device_image (data sources only)
- cloud_pc_gallery_image (data sources only)
- connected_organization
- cross_tenant_access_policy_configuration_default
- cross_tenant_access_policy_configuration_partner
- cross_tenant_access_policy
- device_configuration.windows10_general
- device_management_script
- device_registration_policy
- identity_governance_custom_task_extension
- identity_governance_lifecycle_management_settings
- identity_governance_task_definition (data sources only)
- identity_governance_workflow
- identity_governance_workflow_version (data sources only)
- identity_security_defaults_enforcement_policy
- mobility_management_policy
- sharepoint_settings
- synchronization_schema_json
- unified_role_definition
- unified_role_management_policy_assignment (data sources only)
- unified_role_management_policy
- windows_management_app

### changed

- Allow import of non-importable resources
- Update dependencies
- Update MS Graph metadata & docs
- Enhance data source filters
- Singular data sources: support entity selection by attribute filters
- Singular data sources: support entity selection by OData filter
- Update dependencies
- Update MS Graph metadata & docs

### breaking

- android_managed_app_protection: fix/update provider defaults
- authentication_methods_policy: fix/update provider defaults
- cloud_pc_provisioning_policy: add image_display_name
- device_management_configuration_policy_json: adapt to match device_management_configuration_policy
- ios_managed_app_protection: fix/update provider defaults


## 0.15.1 (2024-12-11)

### fixed

- Fix authentication_methods_policy

### changed

- Update dependencies


## 0.15.0 (2024-12-10)

### added

- authentication_context_class_reference
- device_configuration: add edition_upgrade
- mobile_app
- Add provider function parse_intunewin_metadata
- Add provider function parse_appxmsix_metadata
- windows_driver_update_profile

### fixed

- cloud_pc_provisioning_policy: fix region_name handling
- Fix handling of read-only attributes in nested derived types

### changed

- Docs: major overhaul and various enhancements
- Update metadata: no relevant changes
- cloud_pc_*: adapt for current MS Graph behavior

### breaking

- Cleanup handling of various base64 encoded attributes


## 0.14.0 (2024-10-30)

### changed

- Update MS Graph metadata
- Update dependencies


## 0.13.0 (2024-09-11)

### added

- device_configuration: support android_work_profile_general_device
- device_configuration.ios_device_features: support ios_single_sign_on_extension

### fixed

- Avoid state corruption when deleting...

### changed

- Convert warning to info if entity has vanished from graph
- Docs: Include description and allow values for enum types
- Update graph metadata
