// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// Code has mostly been lifted/copied from https://github.com/hashicorp/terraform-provider-azuread/blob/v3.7.0/internal/provider/provider.go
// To make updates easier, I tried to leave its structure as is as much as possible. Therefore it looks far from pretty ;-)

package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"

	"terraform-provider-microsoft365wp/workplace/external/msgraph"
	"terraform-provider-microsoft365wp/workplace/services"
	mobileappfuncs "terraform-provider-microsoft365wp/workplace/services/mobile_app_funcs"
	"terraform-provider-microsoft365wp/workplace/util/retryablehttputil"

	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/claims"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.ProviderWithFunctions = &workplaceProvider{}
)

// Helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &workplaceProvider{}
}

// Provider implementation.
type workplaceProvider struct{}

// Returns the provider type name.
func (p *workplaceProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "microsoft365wp"
}

// Defines the provider-level schema for configuration data.
func (p *workplaceProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {

	// --- Lifted/copied from internal/provider/provider.go func AzureADProvider() ---

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: pluginsdk.EnvDefaultFunc("ARM_CLIENT_ID", ""),
				Description: "The Client ID which should be used for service principal authentication",
			},
			"client_id_file_path": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: pluginsdk.EnvDefaultFunc("ARM_CLIENT_ID_FILE_PATH", ""),
				Description: "The path to a file containing the Client ID which should be used for service principal authentication",
			},
			"tenant_id": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: pluginsdk.EnvDefaultFunc("ARM_TENANT_ID", ""),
				Description: "The Tenant ID which should be used. Works with all authentication methods except Managed Identity",
			},
			"environment": schema.StringAttribute{
				Optional: true, // required is _not_ needed here since the default will work fine
				// DefaultFunc: pluginsdk.EnvDefaultFunc("ARM_ENVIRONMENT", "global"),
				Description: "The cloud environment which should be used. Possible values are: `global` (also `public`), `usgovernmentl4` (also `usgovernment`), `usgovernmentl5` (also `dod`), and `china`. Defaults to `global`. Not used and should not be specified when `metadata_host` is specified.",
			},
			"metadata_host": schema.StringAttribute{
				Optional: true, // required is _not_ needed here since the default will work fine
				// DefaultFunc: pluginsdk.EnvDefaultFunc("ARM_METADATA_HOSTNAME", ""),
				Description: "The Hostname which should be used for the Azure Metadata Service.",
			},

			// Client Certificate specific fields
			"client_certificate": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: pluginsdk.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE", ""),
				Description: "Base64 encoded PKCS#12 certificate bundle to use when authenticating as a Service Principal using a Client Certificate",
				Sensitive:   true,
			},
			"client_certificate_password": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: pluginsdk.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PASSWORD", ""),
				Description: "The password to decrypt the Client Certificate. For use when authenticating as a Service Principal using a Client Certificate",
				Sensitive:   true,
			},
			"client_certificate_path": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: pluginsdk.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PATH", ""),
				Description: "The path to the Client Certificate associated with the Service Principal for use when authenticating as a Service Principal using a Client Certificate",
			},

			// Client Secret specific fields
			"client_secret": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: pluginsdk.EnvDefaultFunc("ARM_CLIENT_SECRET", ""),
				Description: "The application password to use when authenticating as a Service Principal using a Client Secret",
				Sensitive:   true,
			},
			"client_secret_file_path": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: pluginsdk.EnvDefaultFunc("ARM_CLIENT_SECRET_FILE_PATH", ""),
				Description: "The path to a file containing the application password to use when authenticating as a Service Principal using a Client Secret",
			},

			// OIDC specific fields
			"use_oidc": schema.BoolAttribute{
				Optional: true,
				// DefaultFunc: pluginsdk.EnvDefaultFunc("ARM_USE_OIDC", false),
				Description: "Allow OpenID Connect to be used for authentication",
			},
			"oidc_token": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: pluginsdk.EnvDefaultFunc("ARM_OIDC_TOKEN", ""),
				Description: "The ID token for use when authenticating as a Service Principal using OpenID Connect.",
				Sensitive:   true,
			},
			"oidc_token_file_path": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: pluginsdk.EnvDefaultFunc("ARM_OIDC_TOKEN_FILE_PATH", ""),
				Description: "The path to a file containing an ID token for use when authenticating as a Service Principal using OpenID Connect.",
			},
			"ado_pipeline_service_connection_id": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ARM_ADO_PIPELINE_SERVICE_CONNECTION_ID", "ARM_OIDC_AZURE_SERVICE_CONNECTION_ID"}, nil),
				Description: "The Azure DevOps Pipeline Service Connection ID.",
			},
			"oidc_request_token": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: pluginsdk.MultiEnvDefaultFunc([]string{"ARM_OIDC_REQUEST_TOKEN", "ACTIONS_ID_TOKEN_REQUEST_TOKEN", "SYSTEM_ACCESSTOKEN"}, ""),
				Description: "The bearer token for the request to the OIDC provider. For use when authenticating as a Service Principal using OpenID Connect.",
				Sensitive:   true,
			},
			"oidc_request_url": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: pluginsdk.MultiEnvDefaultFunc([]string{"ARM_OIDC_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_URL", "SYSTEM_OIDCREQUESTURI"}, ""),
				Description: "The URL for the OIDC provider from which to request an ID token. For use when authenticating as a Service Principal using OpenID Connect.",
			},

			// Azure AKS Workload Identity fields
			"use_aks_workload_identity": schema.BoolAttribute{
				Optional: true,
				// DefaultFunc: pluginsdk.EnvDefaultFunc("ARM_USE_AKS_WORKLOAD_IDENTITY", false),
				Description: "Allow Azure AKS Workload Identity to be used for Authentication.",
			},

			// CLI authentication specific fields
			"use_cli": schema.BoolAttribute{
				Optional: true,
				// DefaultFunc: pluginsdk.EnvDefaultFunc("ARM_USE_CLI", true),
				Description: "Allow Azure CLI to be used for Authentication",
			},

			// Managed Identity specific fields
			"use_msi": schema.BoolAttribute{
				Optional: true,
				// DefaultFunc: pluginsdk.EnvDefaultFunc("ARM_USE_MSI", false),
				Description: "Allow Managed Identity to be used for Authentication",
			},
			"msi_endpoint": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: pluginsdk.EnvDefaultFunc("ARM_MSI_ENDPOINT", ""),
				Description: "The path to a custom endpoint for Managed Identity - in most circumstances this should be detected automatically",
			},

			"use_wgt": schema.BoolAttribute{
				Optional:    true,
				Description: "Allow tools/wpGetToken to be used for authentication",
			},
		},
		Description: "Terraform Provider for Microsoft 365",
	}
}

// Prepares an API client for data sources and resources.
func (p *workplaceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	if !req.Config.Raw.IsFullyKnown() {
		resp.Diagnostics.AddError("Unknown Value", "An attribute value is not yet known")
		return
	}

	pch := providerConfigHelper{
		ctx:    ctx,
		diags:  &resp.Diagnostics,
		config: &req.Config,
	}

	// --- Lifted/copied from internal/provider/provider.go func providerConfigure(...) ---

	var certData []byte
	if encodedCert := dGet(&pch, "client_certificate", "ARM_CLIENT_CERTIFICATE", "").(string); encodedCert != "" {
		var err error
		certData, err = decodeCertificate(encodedCert)
		if err != nil {
			diagsAddError(&pch, err)
			return
		}
	}

	idToken, err := getOidcToken(&pch)
	if err != nil {
		diagsAddError(&pch, err)
		return
	}

	clientSecret, err := getClientSecret(&pch)
	if err != nil {
		diagsAddError(&pch, err)
		return
	}

	clientId, err := getClientId(&pch)
	if err != nil {
		diagsAddError(&pch, err)
		return
	}

	tenantId, err := getTenantId(&pch)
	if err != nil {
		diagsAddError(&pch, err)
		return
	}

	var (
		env *environments.Environment

		envName      = dGet(&pch, "environment", "ARM_ENVIRONMENT", "global").(string)
		metadataHost = dGet(&pch, "metadata_host", "ARM_METADATA_HOSTNAME", "").(string)
	)

	if metadataHost != "" {
		if env, err = environments.FromEndpoint(ctx, fmt.Sprintf("https://%s", metadataHost)); err != nil {
			diagsAddError(&pch, err)
			return
		}
	} else {
		if env, err = environments.FromName(envName); err != nil {
			diagsAddError(&pch, err)
			return
		}
	}

	if env.MicrosoftGraph == nil {
		diagsAddError(&pch, errors.New("Microsoft Graph was not configured for the specified environment")) //lint:ignore ST1005 Company name
		return
	} else if endpoint, ok := env.MicrosoftGraph.Endpoint(); !ok || *endpoint == "" {
		diagsAddError(&pch, errors.New("Microsoft Graph endpoint could not be determined for the specified environment")) //lint:ignore ST1005 Company name
		return
	}

	var (
		enableAzureCli        = dGet(&pch, "use_cli", "ARM_USE_CLI", true).(bool)
		enableManagedIdentity = dGet(&pch, "use_msi", "ARM_USE_MSI", false).(bool)
		enableOidc            = dGet(&pch, "use_oidc", "ARM_USE_OIDC", false).(bool) || dGet(&pch, "use_aks_workload_identity", "ARM_USE_AKS_WORKLOAD_IDENTITY", false).(bool)
	)

	authConfig := &auth.Credentials{
		Environment:                    *env,
		ClientID:                       *clientId,
		TenantID:                       *tenantId,
		ClientCertificateData:          certData,
		ClientCertificatePassword:      dGet(&pch, "client_certificate_password", "ARM_CLIENT_CERTIFICATE_PASSWORD", "").(string),
		ClientCertificatePath:          dGet(&pch, "client_certificate_path", "ARM_CLIENT_CERTIFICATE_PATH", "").(string),
		ClientSecret:                   *clientSecret,
		OIDCAssertionToken:             *idToken,
		OIDCTokenRequestURL:            dGet(&pch, "oidc_request_url", "ARM_OIDC_REQUEST_URL", os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")).(string),
		OIDCTokenRequestToken:          dGet(&pch, "oidc_request_token", "ARM_OIDC_REQUEST_TOKEN", os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")).(string),
		ADOPipelineServiceConnectionID: dGet(&pch, "ado_pipeline_service_connection_id", "ARM_ADO_PIPELINE_SERVICE_CONNECTION_ID", os.Getenv("ARM_OIDC_AZURE_SERVICE_CONNECTION_ID")).(string),

		CustomManagedIdentityEndpoint: dGet(&pch, "msi_endpoint", "ARM_MSI_ENDPOINT", "").(string),

		EnableAuthenticatingUsingAzureCLI:          enableAzureCli,
		EnableAuthenticatingUsingClientCertificate: true,
		EnableAuthenticatingUsingClientSecret:      true,
		EnableAuthenticatingUsingManagedIdentity:   enableManagedIdentity,
		EnableAuthenticationUsingGitHubOIDC:        enableOidc,
		EnableAuthenticationUsingADOPipelineOIDC:   enableOidc,
		EnableAuthenticationUsingOIDC:              enableOidc,
	}

	var authorizer auth.Authorizer
	if dGet(&pch, "use_wgt", "ARM_USE_WGT", false).(bool) {
		authorizer, err = NewWgtAuthorizer(ctx)
		if err != nil {
			resp.Diagnostics.AddError("Could not configure WGT Authorizer", err.Error())
			return
		}
	} else {
		// --- Lifted/copied from internal/clients/builder.go func (b *ClientBuilder) Build(...) ---
		authorizer, err = auth.NewAuthorizerFromCredentials(ctx, *authConfig, authConfig.Environment.MicrosoftGraph)
		if err != nil {
			resp.Diagnostics.AddError("Unable to build authorizer", err.Error())
			return
		}
	}

	// --- Lifted/copied from internal/clients/client.go func (client *Client) build(...) ---

	// Acquire an access token upfront, so we can decode the JWT and populate the claims
	// The token will be cached and reused.
	token, err := authorizer.Token(ctx, &http.Request{})
	if err != nil {
		resp.Diagnostics.AddError("Unable to obtain access token", err.Error())
		return
	}

	claims, err := claims.ParseClaims(token)
	if err != nil {
		resp.Diagnostics.AddError("Unable to parse claims in access token", err.Error())
		return
	}

	// Log the claims for debugging
	claimsJson, err := json.Marshal(claims)
	if err != nil {
		tflog.Warn(ctx, "Unable to marshal access token claims for log output to JSON")
	} else if claimsJson == nil {
		tflog.Warn(ctx, "Marshaled access token claims JSON was nil")
	} else {
		tflog.Debug(ctx, fmt.Sprintf("Provider access token claims: %s", claimsJson))
	}

	// --- End of lifted/copied code ---

	// Log HTTP requests and responses
	requestLogger := func(req *http.Request) (*http.Request, error) {
		if req != nil {
			if dump, err := httputil.DumpRequestOut(req, true); err == nil {
				tflog.Info(ctx, fmt.Sprintf("%s\n", dump))
			}
		}
		return req, nil
	}
	responseLogger := func(req *http.Request, resp *http.Response) (*http.Response, error) {
		if resp != nil {
			if dump, err := httputil.DumpResponse(resp, true); err == nil {
				tflog.Info(ctx, fmt.Sprintf("%s\n", dump))
			}
		}
		return resp, nil
	}

	graphClient := msgraph.NewClient(msgraph.VersionBeta)
	graphClient.Authorizer = authorizer
	graphClient.RequestMiddlewares = &[]msgraph.RequestMiddleware{requestLogger}
	graphClient.ResponseMiddlewares = &[]msgraph.ResponseMiddleware{responseLogger}
	retryablehttputil.ConfigureClientRetryLimitsAndBackoff(graphClient.RetryableClient)

	// Make the graphClient available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = &graphClient
	resp.ResourceData = &graphClient
}

// Defines the data sources implemented in the provider.
func (p *workplaceProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return &services.AdministrativeUnitSingularDataSource },
		func() datasource.DataSource { return &services.AdministrativeUnitPluralDataSource },
		func() datasource.DataSource { return &services.AndroidManagedAppProtectionSingularDataSource },
		func() datasource.DataSource { return &services.AndroidManagedAppProtectionPluralDataSource },
		func() datasource.DataSource { return &services.ApplicationSingularDataSource },
		func() datasource.DataSource { return &services.ApplicationPluralDataSource },
		func() datasource.DataSource { return &services.AttributeSetSingularDataSource },
		func() datasource.DataSource { return &services.AttributeSetPluralDataSource },
		func() datasource.DataSource {
			return &services.AuthenticationCombinationConfigurationSingularDataSource
		},
		func() datasource.DataSource { return &services.AuthenticationCombinationConfigurationPluralDataSource },
		func() datasource.DataSource { return &services.AuthenticationContextClassReferenceSingularDataSource },
		func() datasource.DataSource { return &services.AuthenticationContextClassReferencePluralDataSource },
		func() datasource.DataSource { return &services.AuthenticationFlowsPolicySingularDataSource },
		func() datasource.DataSource { return &services.AuthenticationMethodsPolicySingularDataSource },
		func() datasource.DataSource { return &services.AuthenticationStrengthPolicySingularDataSource },
		func() datasource.DataSource { return &services.AuthenticationStrengthPolicyPluralDataSource },
		func() datasource.DataSource { return &services.AuthorizationPolicySingularDataSource },
		func() datasource.DataSource {
			return &services.AzureAdWindowsAutopilotDeploymentProfileSingularDataSource
		},
		func() datasource.DataSource {
			return &services.AzureAdWindowsAutopilotDeploymentProfilePluralDataSource
		},
		func() datasource.DataSource {
			return &services.AzureAdWindowsAutopilotDeploymentProfileAssignmentSingularDataSource
		},
		func() datasource.DataSource {
			return &services.AzureAdWindowsAutopilotDeploymentProfileAssignmentPluralDataSource
		},
		func() datasource.DataSource { return &services.CloudPcDeviceImageSingularDataSource },
		func() datasource.DataSource { return &services.CloudPcDeviceImagePluralDataSource },
		func() datasource.DataSource { return &services.CloudPcGalleryImageSingularDataSource },
		func() datasource.DataSource { return &services.CloudPcGalleryImagePluralDataSource },
		func() datasource.DataSource { return &services.CloudPcProvisioningPolicySingularDataSource },
		func() datasource.DataSource { return &services.CloudPcProvisioningPolicyPluralDataSource },
		func() datasource.DataSource { return &services.CloudPcUserSettingSingularDataSource },
		func() datasource.DataSource { return &services.CloudPcUserSettingPluralDataSource },
		func() datasource.DataSource { return &services.ConditionalAccessPolicySingularDataSource },
		func() datasource.DataSource { return &services.ConditionalAccessPolicyPluralDataSource },
		func() datasource.DataSource { return &services.ConnectedOrganizationSingularDataSource },
		func() datasource.DataSource { return &services.ConnectedOrganizationPluralDataSource },
		func() datasource.DataSource { return &services.ConnectorSingularDataSource },
		func() datasource.DataSource { return &services.ConnectorPluralDataSource },
		func() datasource.DataSource { return &services.ConnectorGroupSingularDataSource },
		func() datasource.DataSource { return &services.ConnectorGroupPluralDataSource },
		func() datasource.DataSource { return &services.ConnectorGroupApplicationSingularDataSource },
		func() datasource.DataSource { return &services.ConnectorGroupApplicationPluralDataSource },
		func() datasource.DataSource { return &services.ConnectorGroupMemberConnectorSingularDataSource },
		func() datasource.DataSource { return &services.ConnectorGroupMemberConnectorPluralDataSource },
		func() datasource.DataSource { return &services.CrossTenantAccessPolicySingularDataSource },
		func() datasource.DataSource {
			return &services.CrossTenantAccessPolicyConfigurationDefaultSingularDataSource
		},
		func() datasource.DataSource {
			return &services.CrossTenantAccessPolicyConfigurationPartnerSingularDataSource
		},
		func() datasource.DataSource {
			return &services.CrossTenantAccessPolicyConfigurationPartnerPluralDataSource
		},
		func() datasource.DataSource { return &services.CrossTenantIdentitySyncPolicyPartnerSingularDataSource },
		func() datasource.DataSource { return &services.CustomSecurityAttributeDefinitionSingularDataSource },
		func() datasource.DataSource { return &services.CustomSecurityAttributeDefinitionPluralDataSource },
		func() datasource.DataSource {
			return &services.DeviceAndAppManagementAssignmentFilterSingularDataSource
		},
		func() datasource.DataSource { return &services.DeviceAndAppManagementAssignmentFilterPluralDataSource },
		func() datasource.DataSource { return &services.DeviceCompliancePolicySingularDataSource },
		func() datasource.DataSource { return &services.DeviceCompliancePolicyPluralDataSource },
		func() datasource.DataSource { return &services.DeviceComplianceScriptSingularDataSource },
		func() datasource.DataSource { return &services.DeviceComplianceScriptPluralDataSource },
		func() datasource.DataSource { return &services.DeviceConfigurationSingularDataSource },
		func() datasource.DataSource { return &services.DeviceConfigurationPluralDataSource },
		func() datasource.DataSource { return &services.DeviceConfigurationCustomSingularDataSource },
		func() datasource.DataSource { return &services.DeviceConfigurationCustomPluralDataSource },
		func() datasource.DataSource { return &services.DeviceCustomAttributeShellScriptSingularDataSource },
		func() datasource.DataSource { return &services.DeviceCustomAttributeShellScriptPluralDataSource },
		func() datasource.DataSource { return &services.DeviceEnrollmentConfigurationSingularDataSource },
		func() datasource.DataSource { return &services.DeviceEnrollmentConfigurationPluralDataSource },
		func() datasource.DataSource { return &services.DeviceManagementConfigurationPolicySingularDataSource },
		func() datasource.DataSource { return &services.DeviceManagementConfigurationPolicyPluralDataSource },
		func() datasource.DataSource {
			return &services.DeviceManagementConfigurationPolicyJsonSingularDataSource
		},
		func() datasource.DataSource { return &services.DeviceManagementConfigurationPolicyJsonPluralDataSource },
		func() datasource.DataSource { return &services.DeviceManagementIntentSingularDataSource },
		func() datasource.DataSource { return &services.DeviceManagementIntentPluralDataSource },
		func() datasource.DataSource { return &services.DeviceManagementScriptSingularDataSource },
		func() datasource.DataSource { return &services.DeviceManagementScriptPluralDataSource },
		func() datasource.DataSource { return &services.DeviceRegistrationPolicySingularDataSource },
		func() datasource.DataSource { return &services.DeviceShellScriptSingularDataSource },
		func() datasource.DataSource { return &services.DeviceShellScriptPluralDataSource },
		func() datasource.DataSource { return &services.ExternalIdentitiesPolicySingularDataSource },
		func() datasource.DataSource { return &services.GroupSingularDataSource },
		func() datasource.DataSource { return &services.GroupPluralDataSource },
		func() datasource.DataSource { return &services.GroupAssignedLicenseSingularDataSource },
		func() datasource.DataSource { return &services.GroupAssignedLicensePluralDataSource },
		func() datasource.DataSource { return &services.IdentityGovernanceCustomTaskExtensionSingularDataSource },
		func() datasource.DataSource { return &services.IdentityGovernanceCustomTaskExtensionPluralDataSource },
		func() datasource.DataSource {
			return &services.IdentityGovernanceLifecycleManagementSettingsSingularDataSource
		},
		func() datasource.DataSource { return &services.IdentityGovernanceTaskDefinitionSingularDataSource },
		func() datasource.DataSource { return &services.IdentityGovernanceTaskDefinitionPluralDataSource },
		func() datasource.DataSource { return &services.IdentityGovernanceWorkflowSingularDataSource },
		func() datasource.DataSource { return &services.IdentityGovernanceWorkflowPluralDataSource },
		func() datasource.DataSource { return &services.IdentityGovernanceWorkflowVersionSingularDataSource },
		func() datasource.DataSource { return &services.IdentityGovernanceWorkflowVersionPluralDataSource },
		func() datasource.DataSource {
			return &services.IdentitySecurityDefaultsEnforcementPolicySingularDataSource
		},
		func() datasource.DataSource { return &services.IntuneBrandingProfileSingularDataSource },
		func() datasource.DataSource { return &services.IntuneBrandingProfilePluralDataSource },
		func() datasource.DataSource { return &services.IosManagedAppProtectionSingularDataSource },
		func() datasource.DataSource { return &services.IosManagedAppProtectionPluralDataSource },
		func() datasource.DataSource { return &services.ManagedDeviceMobileAppConfigurationSingularDataSource },
		func() datasource.DataSource { return &services.ManagedDeviceMobileAppConfigurationPluralDataSource },
		func() datasource.DataSource { return &services.MobileAppSingularDataSource },
		func() datasource.DataSource { return &services.MobileAppPluralDataSource },
		func() datasource.DataSource { return &services.MobileAppCategorySingularDataSource },
		func() datasource.DataSource { return &services.MobileAppCategoryPluralDataSource },
		func() datasource.DataSource { return &services.MobilityManagementPolicySingularDataSource },
		func() datasource.DataSource { return &services.MobilityManagementPolicyPluralDataSource },
		func() datasource.DataSource { return &services.NetworkaccessTenantStatusSingularDataSource },
		func() datasource.DataSource { return &services.NotificationMessageTemplateSingularDataSource },
		func() datasource.DataSource { return &services.NotificationMessageTemplatePluralDataSource },
		func() datasource.DataSource { return &services.ServicePrincipalSingularDataSource },
		func() datasource.DataSource { return &services.ServicePrincipalPluralDataSource },
		func() datasource.DataSource { return &services.SharepointSettingsSingularDataSource },
		func() datasource.DataSource { return &services.SubscribedSkuSingularDataSource },
		func() datasource.DataSource { return &services.SubscribedSkuPluralDataSource },
		func() datasource.DataSource { return &services.SynchronizationSchemaJsonSingularDataSource },
		func() datasource.DataSource { return &services.TargetedManagedAppConfigurationSingularDataSource },
		func() datasource.DataSource { return &services.TargetedManagedAppConfigurationPluralDataSource },
		func() datasource.DataSource { return &services.UnifiedRoleDefinitionSingularDataSource },
		func() datasource.DataSource { return &services.UnifiedRoleDefinitionPluralDataSource },
		func() datasource.DataSource { return &services.UnifiedRoleManagementPolicySingularDataSource },
		func() datasource.DataSource { return &services.UnifiedRoleManagementPolicyPluralDataSource },
		func() datasource.DataSource { return &services.UnifiedRoleManagementPolicyAssignmentSingularDataSource },
		func() datasource.DataSource { return &services.UnifiedRoleManagementPolicyAssignmentPluralDataSource },
		func() datasource.DataSource { return &services.UserSingularDataSource },
		func() datasource.DataSource { return &services.UserPluralDataSource },
		func() datasource.DataSource { return &services.WindowsDriverUpdateProfileSingularDataSource },
		func() datasource.DataSource { return &services.WindowsDriverUpdateProfilePluralDataSource },
		func() datasource.DataSource { return &services.WindowsFeatureUpdateProfileSingularDataSource },
		func() datasource.DataSource { return &services.WindowsFeatureUpdateProfilePluralDataSource },
		func() datasource.DataSource { return &services.WindowsManagementAppSingularDataSource },
	}
}

// Defines the resources implemented in the provider.
func (p *workplaceProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &services.AdministrativeUnitResource },
		func() resource.Resource { return &services.AndroidManagedAppProtectionResource },
		func() resource.Resource { return &services.AttributeSetResource },
		func() resource.Resource { return &services.AuthenticationCombinationConfigurationResource },
		func() resource.Resource { return &services.AuthenticationContextClassReferenceResource },
		func() resource.Resource { return &services.AuthenticationFlowsPolicyResource },
		func() resource.Resource { return &services.AuthenticationMethodsPolicyResource },
		func() resource.Resource { return &services.AuthenticationStrengthPolicyResource },
		func() resource.Resource { return &services.AuthorizationPolicyResource },
		func() resource.Resource { return &services.AzureAdWindowsAutopilotDeploymentProfileResource },
		func() resource.Resource { return &services.AzureAdWindowsAutopilotDeploymentProfileAssignmentResource },
		func() resource.Resource { return &services.CloudPcProvisioningPolicyResource },
		func() resource.Resource { return &services.CloudPcUserSettingResource },
		func() resource.Resource { return &services.ConditionalAccessPolicyResource },
		func() resource.Resource { return &services.ConnectedOrganizationResource },
		func() resource.Resource { return &services.ConnectorGroupResource },
		func() resource.Resource { return &services.ConnectorGroupApplicationResource },
		func() resource.Resource { return &services.ConnectorGroupMemberConnectorResource },
		func() resource.Resource { return &services.CrossTenantAccessPolicyResource },
		func() resource.Resource { return &services.CrossTenantAccessPolicyConfigurationDefaultResource },
		func() resource.Resource { return &services.CrossTenantAccessPolicyConfigurationPartnerResource },
		func() resource.Resource { return &services.CrossTenantIdentitySyncPolicyPartnerResource },
		func() resource.Resource { return &services.CustomSecurityAttributeDefinitionResource },
		func() resource.Resource { return &services.DeviceAndAppManagementAssignmentFilterResource },
		func() resource.Resource { return &services.DeviceCompliancePolicyResource },
		func() resource.Resource { return &services.DeviceComplianceScriptResource },
		func() resource.Resource { return &services.DeviceConfigurationResource },
		func() resource.Resource { return &services.DeviceConfigurationCustomResource },
		func() resource.Resource { return &services.DeviceCustomAttributeShellScriptResource },
		func() resource.Resource { return &services.DeviceEnrollmentConfigurationResource },
		func() resource.Resource { return &services.DeviceManagementConfigurationPolicyResource },
		func() resource.Resource { return &services.DeviceManagementConfigurationPolicyJsonResource },
		func() resource.Resource { return &services.DeviceManagementIntentResource },
		func() resource.Resource { return &services.DeviceManagementScriptResource },
		func() resource.Resource { return &services.DeviceRegistrationPolicyResource },
		func() resource.Resource { return &services.DeviceShellScriptResource },
		func() resource.Resource { return &services.ExternalIdentitiesPolicyResource },
		func() resource.Resource { return &services.GroupAssignedLicenseResource },
		func() resource.Resource { return &services.IdentityGovernanceCustomTaskExtensionResource },
		func() resource.Resource { return &services.IdentityGovernanceLifecycleManagementSettingsResource },
		func() resource.Resource { return &services.IdentityGovernanceWorkflowResource },
		func() resource.Resource { return &services.IdentitySecurityDefaultsEnforcementPolicyResource },
		func() resource.Resource { return &services.IntuneBrandingProfileResource },
		func() resource.Resource { return &services.IosManagedAppProtectionResource },
		func() resource.Resource { return &services.ManagedDeviceMobileAppConfigurationResource },
		func() resource.Resource { return &services.MobileAppResource },
		func() resource.Resource { return &services.MobileAppCategoryResource },
		func() resource.Resource { return &services.MobilityManagementPolicyResource },
		func() resource.Resource { return &services.NetworkaccessTenantStatusResource },
		func() resource.Resource { return &services.NotificationMessageTemplateResource },
		func() resource.Resource { return &services.SharepointSettingsResource },
		func() resource.Resource { return &services.SynchronizationSchemaJsonResource },
		func() resource.Resource { return &services.TargetedManagedAppConfigurationResource },
		func() resource.Resource { return &services.UnifiedRoleDefinitionResource },
		func() resource.Resource { return &services.UnifiedRoleManagementPolicyResource },
		func() resource.Resource { return &services.WindowsDriverUpdateProfileResource },
		func() resource.Resource { return &services.WindowsFeatureUpdateProfileResource },
		func() resource.Resource { return &services.WindowsManagementAppResource },
	}
}

func (p *workplaceProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		func() function.Function { return &mobileappfuncs.ParseIntunewinMetadataFunction{} },
		func() function.Function { return &mobileappfuncs.ParseAppxMsixMetadataFunction{} },
	}
}
