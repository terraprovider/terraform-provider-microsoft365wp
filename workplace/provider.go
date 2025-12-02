// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package workplace

import (
	"context"
	"encoding/base64"
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
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
// Lifted from https://github.com/hashicorp/terraform-provider-azuread/blob/v2.36.0/internal/provider/provider.go
func (p *workplaceProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_ID", ""),
				Description: "The Client ID which should be used for service principal authentication",
			},
			"tenant_id": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("ARM_TENANT_ID", ""),
				Description: "The Tenant ID which should be used. Works with all authentication methods except Managed Identity",
			},
			"environment": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("ARM_ENVIRONMENT", "global"),
				Description: "The cloud environment which should be used. Possible values are: `global` (also `public`), `usgovernmentl4` (also `usgovernment`), `usgovernmentl5` (also `dod`), and `china`. Defaults to `global`",
			},
			"metadata_host": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("ARM_METADATA_HOSTNAME", ""),
				Description: "The Hostname which should be used for the Azure Metadata Service.",
			},

			// Client Certificate specific fields
			"client_certificate": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE", ""),
				Description: "Base64 encoded PKCS#12 certificate bundle to use when authenticating as a Service Principal using a Client Certificate",
				Sensitive:   true,
			},
			"client_certificate_password": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PASSWORD", ""),
				Description: "The password to decrypt the Client Certificate. For use when authenticating as a Service Principal using a Client Certificate",
				Sensitive:   true,
			},
			"client_certificate_path": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PATH", ""),
				Description: "The path to the Client Certificate associated with the Service Principal for use when authenticating as a Service Principal using a Client Certificate",
			},

			// Client Secret specific fields
			"client_secret": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET", ""),
				Description: "The application password to use when authenticating as a Service Principal using a Client Secret",
				Sensitive:   true,
			},

			// OIDC specific fields
			"use_oidc": schema.BoolAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("ARM_USE_OIDC", false),
				Description: "Allow OpenID Connect to be used for authentication",
			},
			"oidc_token": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("ARM_OIDC_TOKEN", ""),
				Description: "The ID token for use when authenticating as a Service Principal using OpenID Connect.",
				Sensitive:   true,
			},
			"oidc_token_file_path": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("ARM_OIDC_TOKEN_FILE_PATH", ""),
				Description: "The path to a file containing an ID token for use when authenticating as a Service Principal using OpenID Connect.",
			},
			"oidc_request_token": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ARM_OIDC_REQUEST_TOKEN", "ACTIONS_ID_TOKEN_REQUEST_TOKEN"}, ""),
				Description: "The bearer token for the request to the OIDC provider. For use when authenticating as a Service Principal using OpenID Connect.",
				Sensitive:   true,
			},
			"oidc_request_url": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: schema.MultiEnvDefaultFunc([]string{"ARM_OIDC_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_URL"}, ""),
				Description: "The URL for the OIDC provider from which to request an ID token. For use when authenticating as a Service Principal using OpenID Connect.",
			},

			// CLI authentication specific fields
			"use_cli": schema.BoolAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("ARM_USE_CLI", true),
				Description: "Allow Azure CLI to be used for Authentication",
			},

			// Managed Identity specific fields
			"use_msi": schema.BoolAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("ARM_USE_MSI", false),
				Description: "Allow Managed Identity to be used for Authentication",
			},
			"msi_endpoint": schema.StringAttribute{
				Optional: true,
				// DefaultFunc: schema.EnvDefaultFunc("ARM_MSI_ENDPOINT", ""),
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
// Lifted from https://github.com/hashicorp/terraform-provider-azuread/blob/v2.36.0/internal/provider/provider.go
// To make updates easier, I tried to leave its structure as is as much as possible. Therefore it looks far from pretty ;-)
func (p *workplaceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	if !req.Config.Raw.IsFullyKnown() {
		resp.Diagnostics.AddError("Unknown Value", "An attribute value is not yet known")
		return
	}

	// --- Helper functions ---

	dGet := func(attributeName string, envVarName string, defaultValue any) any {
		var tfTarget attr.Value
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root(attributeName), &tfTarget)...)
		switch typedTfTarget := tfTarget.(type) {

		case types.String:
			var result string
			if !typedTfTarget.IsNull() {
				result = typedTfTarget.ValueString()
			} else if envVarValue := os.Getenv(envVarName); envVarValue != "" {
				result = envVarValue
			} else {
				result = defaultValue.(string)
			}
			return result

		case types.Bool:
			var result bool
			if !typedTfTarget.IsNull() {
				result = typedTfTarget.ValueBool()
			} else if envVarValue := os.Getenv(envVarName); envVarValue != "" {
				result = envVarValue == "1" || envVarValue == "true"
			} else {
				result = defaultValue.(bool)
			}
			return result

		}

		panic(fmt.Sprintf("Don't know how to deal with config attribute of type %T", tfTarget))
	}

	addError := func(err error) {
		resp.Diagnostics.AddError(
			"Error configuring microsoft365wp provider",
			fmt.Sprintf("Error configuring the microsoft365wp provider.\n%s\n", err),
		)
	}

	// --- Copied from member functions ---

	decodeCertificate := func(clientCertificate string) ([]byte, error) {
		var pfx []byte
		if clientCertificate != "" {
			out := make([]byte, base64.StdEncoding.DecodedLen(len(clientCertificate)))
			n, err := base64.StdEncoding.Decode(out, []byte(clientCertificate))
			if err != nil {
				return pfx, fmt.Errorf("could not decode client certificate data: %v", err)
			}
			pfx = out[:n]
		}
		return pfx, nil
	}

	oidcToken := func() (string, error) {
		idToken := dGet("oidc_token", "ARM_OIDC_TOKEN", "").(string)

		if path := dGet("oidc_token_file_path", "ARM_OIDC_TOKEN_FILE_PATH", "").(string); path != "" {
			fileToken, err := os.ReadFile(path)

			if err != nil {
				return "", fmt.Errorf("reading OIDC Token from file %q: %v", path, err)
			}

			if idToken != "" && idToken != string(fileToken) {
				return "", fmt.Errorf("mismatch between supplied OIDC token and supplied OIDC token file contents - please either remove one or ensure they match")
			}

			idToken = string(fileToken)
		}

		return idToken, nil
	}

	// --- Copied from providerConfigure ---

	var certData []byte
	if encodedCert := dGet("client_certificate", "ARM_CLIENT_CERTIFICATE", "").(string); encodedCert != "" {
		var err error
		certData, err = decodeCertificate(encodedCert)
		if err != nil {
			addError(err)
			return
		}
	}

	var (
		env *environments.Environment
		err error

		envName      = dGet("environment", "ARM_ENVIRONMENT", "global").(string)
		metadataHost = dGet("metadata_host", "ARM_METADATA_HOSTNAME", "").(string)
	)

	if metadataHost != "" {
		if env, err = environments.FromEndpoint(ctx, fmt.Sprintf("https://%s", metadataHost), envName); err != nil {
			addError(err)
			return
		}
	} else if env, err = environments.FromName(envName); err != nil {
		addError(err)
		return
	}

	if env.MicrosoftGraph == nil {
		addError(errors.New("Microsoft Graph was not configured for the specified environment")) //lint:ignore ST1005 Company name
		return
	} else if endpoint, ok := env.MicrosoftGraph.Endpoint(); !ok || *endpoint == "" {
		addError(errors.New("Microsoft Graph endpoint could not be determined for the specified environment")) //lint:ignore ST1005 Company name
		return
	}

	idToken, err := oidcToken()
	if err != nil {
		addError(err)
		return
	}

	authConfig := auth.Credentials{
		Environment:                 *env,
		TenantID:                    dGet("tenant_id", "ARM_TENANT_ID", "").(string),
		ClientID:                    dGet("client_id", "ARM_CLIENT_ID", "").(string),
		ClientCertificateData:       certData,
		ClientCertificatePassword:   dGet("client_certificate_password", "ARM_CLIENT_CERTIFICATE_PASSWORD", "").(string),
		ClientCertificatePath:       dGet("client_certificate_path", "ARM_CLIENT_CERTIFICATE_PATH", "").(string),
		ClientSecret:                dGet("client_secret", "ARM_CLIENT_SECRET", "").(string),
		OIDCAssertionToken:          idToken,
		GitHubOIDCTokenRequestURL:   dGet("oidc_request_url", "ARM_OIDC_REQUEST_URL", os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")).(string),
		GitHubOIDCTokenRequestToken: dGet("oidc_request_token", "ARM_OIDC_REQUEST_TOKEN", os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")).(string),
		EnableAuthenticatingUsingClientCertificate: true,
		EnableAuthenticatingUsingClientSecret:      true,
		EnableAuthenticationUsingOIDC:              dGet("use_oidc", "ARM_USE_OIDC", false).(bool),
		EnableAuthenticationUsingGitHubOIDC:        dGet("use_oidc", "ARM_USE_OIDC", false).(bool),
		EnableAuthenticatingUsingAzureCLI:          dGet("use_cli", "ARM_USE_CLI", true).(bool),
		EnableAuthenticatingUsingManagedIdentity:   dGet("use_msi", "ARM_USE_MSI", false).(bool),
		CustomManagedIdentityEndpoint:              dGet("msi_endpoint", "ARM_MSI_ENDPOINT", "").(string),
	}

	api := authConfig.Environment.MicrosoftGraph

	var authorizer auth.Authorizer
	if dGet("use_wgt", "ARM_USE_WGT", false).(bool) {
		authorizer, err = NewWgtAuthorizer(ctx)
		if err != nil {
			resp.Diagnostics.AddError("Could not configure AzureCli Authorizer", err.Error())
		}
	} else {
		// --- Copied from internal/clients/ClientBuilder ---
		authorizer, err = auth.NewAuthorizerFromCredentials(ctx, authConfig, api)
		if err != nil {
			resp.Diagnostics.AddError("Unable to build authorizer", err.Error())
			return
		}
	}

	// --- Copied from internal/clients/Client ---
	// The token will be cached and reused.

	// Acquire an access token upfront, so we can decode the JWT and populate the claims
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
		tflog.Warn(ctx, "Unable to marshal access token claims for log to JSON")
	} else if claimsJson == nil {
		tflog.Warn(ctx, "Marshaled access token claims JSON was nil")
	} else {
		tflog.Debug(ctx, fmt.Sprintf("Provider access token claims: %s", claimsJson))
	}

	// --- End of copied code ---

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
