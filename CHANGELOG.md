# `microsoft365wp` Terraform Provider Changelog


## 0.18.0 (2026-02-16)

### added

- [connector](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/data-sources/connector)
- [connector_group](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/connector_group)
- [connector_group_application](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/connector_group_application)
- [connector_group_member_connector](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/connector_group_member_connector)
- [networkaccess_tenant_status](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/networkaccess_tenant_status)
- implement TF plugin framework default values

### fixed

- [group_assigned_license](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/group_assigned_license): fix concurrency issues

### changed

- [cloud_pc_provisioning_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/cloud_pc_provisioning_policy): fix examples
- [microsoft365wp_mobile_app.android_managed_store](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/mobile_app#nestedatt--android_managed_store): update docs
- [windows_management_app](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/windows_management_app): update docs
- Update MS Graph metadata & docs
- Update MS Graph metadata & docs (again)
- Update dependencies
- Enhance/cleanup [docs](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/)
- Changelog: Reformat, add links to docs
- msgraphutils.SingletonSyntheticIdGraphToTerraformMiddleware

### breaking

- [authentication_methods_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/authentication_methods_policy): add verifiable_credentials
- [authorization_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/authorization_policy): align default values
- [authorization_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/authorization_policy): align default values (again)
- [mobile_app](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/mobile_app): align default values


## 0.17.4 (2025-12-03)

### changed

- [conditional_access_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/conditional_access_policy): Update


## 0.17.3 (2025-12-02)

### added

- [administrative_unit](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/administrative_unit)

### changed

- Update dependencies


## 0.17.2 (2025-11-20)

### fixed

- [unified_role_management_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/unified_role_management_policy): Fix unknown value after create

### changed

- Update dependencies


## 0.17.1 (2025-10-21)

### fixed

- [authentication_methods_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/authentication_methods_policy): Ignore VerifiableCredentials
- [unified_role_management_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/unified_role_management_policy): Fix invisible state drift

### changed

- Update dependencies


## 0.17.0 (2025-07-08)

### added

- [application](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/data-sources/application) (data sources only)
- [attribute_set](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/attribute_set)
- [cross_tenant_identity_sync_policy_partner](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/cross_tenant_identity_sync_policy_partner)
- [custom_security_attribute_definition](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/custom_security_attribute_definition)
- [device_compliance_script](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/device_compliance_script)
- [external_identities_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/external_identities_policy)
- [group](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/data-sources/group) (data sources only)
- [group_assigned_license](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/group_assigned_license)
- [intune_branding_profile](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/intune_branding_profile)
- [service_principal](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/data-sources/service_principal) (data sources only)
- [subscribed_sku](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/data-sources/subscribed_sku) (data sources only)
- [user](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/data-sources/user) (data sources only)

### fixed

- assignmentTarget: Fix filter_id
- [azure_ad_windows_autopilot_deployment_profile](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/azure_ad_windows_autopilot_deployment_profile): Fix locale="" on create

### changed

- [authentication_methods_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/authentication_methods_policy): Convert null from MS Graph to "migrationComplete"
- Update dependencies
- Update MS Graph metadata & docs

### breaking

- [authentication_methods_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/authentication_methods_policy): Add qr_code_pin


## 0.16.0 (2025-03-12)

### added

- [authentication_combination_configuration](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/authentication_combination_configuration)
- [authentication_flows_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/authentication_flows_policy)
- [authorization_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/authorization_policy)
- [cloud_pc_device_image](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/data-sources/cloud_pc_device_image) (data sources only)
- [cloud_pc_gallery_image](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/data-sources/cloud_pc_gallery_image) (data sources only)
- [connected_organization](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/connected_organization)
- [cross_tenant_access_policy_configuration_default](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/cross_tenant_access_policy_configuration_default)
- [cross_tenant_access_policy_configuration_partner](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/cross_tenant_access_policy_configuration_partner)
- [cross_tenant_access_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/cross_tenant_access_policy)
- [device_configuration.windows10_general](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/device_configuration#nestedatt--windows10_general)
- [device_management_script](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/device_management_script)
- [device_registration_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/device_registration_policy)
- [identity_governance_custom_task_extension](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/identity_governance_custom_task_extension)
- [identity_governance_lifecycle_management_settings](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/identity_governance_lifecycle_management_settings)
- [identity_governance_task_definition](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/data-sources/identity_governance_task_definition) (data sources only)
- [identity_governance_workflow](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/identity_governance_workflow)
- [identity_governance_workflow_version](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/data-sources/identity_governance_workflow_version) (data sources only)
- [identity_security_defaults_enforcement_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/identity_security_defaults_enforcement_policy)
- [mobility_management_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/mobility_management_policy)
- [sharepoint_settings](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/sharepoint_settings)
- [synchronization_schema_json](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/synchronization_schema_json)
- [unified_role_definition](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/unified_role_definition)
- [unified_role_management_policy_assignment](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/data-sources/unified_role_management_policy_assignment) (data sources only)
- [unified_role_management_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/unified_role_management_policy)
- [windows_management_app](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/windows_management_app)

### changed

- Allow import of non-importable resources
- Update dependencies
- Update MS Graph metadata & docs
- Enhance data source filters
- Singular data sources: Support entity selection by attribute filters
- Singular data sources: Support entity selection by OData filter
- Update dependencies
- Update MS Graph metadata & docs

### breaking

- [android_managed_app_protection](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/android_managed_app_protection): Fix/update provider defaults
- [authentication_methods_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/authentication_methods_policy): Fix/update provider defaults
- [cloud_pc_provisioning_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/cloud_pc_provisioning_policy): Add image_display_name
- [device_management_configuration_policy_json](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/device_management_configuration_policy_json): Adapt to match device_management_configuration_policy
- [ios_managed_app_protection](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/ios_managed_app_protection): Fix/update provider defaults


## 0.15.1 (2024-12-11)

### fixed

- Fix authentication_methods_policy

### changed

- Update dependencies


## 0.15.0 (2024-12-10)

### added

- [authentication_context_class_reference](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/authentication_context_class_reference)
- [device_configuration](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/device_configuration): Add edition_upgrade
- [mobile_app](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/mobile_app)
- Add provider function parse_intunewin_metadata
- Add provider function parse_appxmsix_metadata
- [windows_driver_update_profile](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/windows_driver_update_profile)

### fixed

- [cloud_pc_provisioning_policy](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/cloud_pc_provisioning_policy): Fix region_name handling
- Fix handling of read-only attributes in nested derived types

### changed

- Docs: Major overhaul and various enhancements
- Update metadata: No relevant changes
- cloud_pc_*: Adapt for current MS Graph behavior

### breaking

- Cleanup handling of various base64 encoded attributes


## 0.14.0 (2024-10-30)

### changed

- Update MS Graph metadata
- Update dependencies


## 0.13.0 (2024-09-11)

### added

- [device_configuration](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/device_configuration): Support android_work_profile_general_device
- [device_configuration.ios_device_features](https://registry.terraform.io/providers/terraprovider/microsoft365wp/latest/docs/resources/device_configuration#nestedatt--ios_device_features): Support ios_single_sign_on_extension

### fixed

- Avoid state corruption when deleting

### changed

- Convert warning to info if entity has vanished from graph
- Docs: Include description and allow values for enum types
- Update graph metadata
