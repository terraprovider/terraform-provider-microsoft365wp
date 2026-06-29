// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// Code has mostly been lifted/copied from https://github.com/hashicorp/terraform-provider-azuread/blob/v3.7.0/internal/provider/helpers.go
// To make updates easier, I tried to leave its structure as is as much as possible. Therefore it looks far from pretty ;-)

package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type providerConfigHelper struct {
	ctx    context.Context
	diags  *diag.Diagnostics
	config *tfsdk.Config
}

func diagsAddError(pch *providerConfigHelper, err error) {
	pch.diags.AddError(
		"Error configuring microsoft365wp provider",
		fmt.Sprintf("Error configuring the microsoft365wp provider.\n%s\n", err),
	)
}

func dGet(pch *providerConfigHelper, attributeName string, envVarName string, defaultValue any) any {
	var tfTarget attr.Value
	pch.diags.Append(pch.config.GetAttribute(pch.ctx, path.Root(attributeName), &tfTarget)...)
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

func decodeCertificate(clientCertificate string) ([]byte, error) {
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

func getOidcToken(pch *providerConfigHelper) (*string, error) {
	idToken := dGet(pch, "oidc_token", "ARM_OIDC_TOKEN", "").(string)

	if path := dGet(pch, "oidc_token_file_path", "ARM_OIDC_TOKEN_FILE_PATH", "").(string); path != "" {
		fileTokenRaw, err := os.ReadFile(path)

		if err != nil {
			return nil, fmt.Errorf("reading OIDC Token from file %q: %v", path, err)
		}

		fileToken := strings.TrimSpace(string(fileTokenRaw))

		if idToken != "" && idToken != fileToken {
			return nil, fmt.Errorf("mismatch between supplied OIDC token and supplied OIDC token file contents - please either remove one or ensure they match")
		}

		idToken = fileToken
	}

	return &idToken, nil
}

func getClientId(pch *providerConfigHelper) (*string, error) {
	clientId := strings.TrimSpace(dGet(pch, "client_id", "ARM_CLIENT_ID", "").(string))

	if path := dGet(pch, "client_id_file_path", "ARM_CLIENT_ID_FILE_PATH", "").(string); path != "" {
		fileClientIdRaw, err := os.ReadFile(path)

		if err != nil {
			return nil, fmt.Errorf("reading Client ID from file %q: %v", path, err)
		}

		fileClientId := strings.TrimSpace(string(fileClientIdRaw))

		if clientId != "" && clientId != fileClientId {
			return nil, fmt.Errorf("mismatch between supplied Client ID and supplied Client ID file contents - please either remove one or ensure they match")
		}

		clientId = fileClientId
	}

	return &clientId, nil
}

func getClientSecret(pch *providerConfigHelper) (*string, error) {
	clientSecret := strings.TrimSpace(dGet(pch, "client_secret", "ARM_CLIENT_SECRET", "").(string))

	if path := dGet(pch, "client_secret_file_path", "ARM_CLIENT_SECRET_FILE_PATH", "").(string); path != "" {
		fileSecretRaw, err := os.ReadFile(path)

		if err != nil {
			return nil, fmt.Errorf("reading Client Secret from file %q: %v", path, err)
		}

		fileSecret := strings.TrimSpace(string(fileSecretRaw))

		if clientSecret != "" && clientSecret != fileSecret {
			return nil, fmt.Errorf("mismatch between supplied Client Secret and supplied Client Secret file contents - please either remove one or ensure they match")
		}

		clientSecret = fileSecret
	}

	return &clientSecret, nil
}

func getTenantId(pch *providerConfigHelper) (*string, error) {
	tenantId := strings.TrimSpace(dGet(pch, "tenant_id", "ARM_TENANT_ID", "").(string))

	if dGet(pch, "use_aks_workload_identity", "ARM_USE_AKS_WORKLOAD_IDENTITY", false).(bool) && os.Getenv("AZURE_TENANT_ID") != "" {
		aksTenantId := os.Getenv("AZURE_TENANT_ID")
		if tenantId != "" && tenantId != aksTenantId {
			return nil, fmt.Errorf("mismatch between supplied Tenant ID and that provided by AKS Workload Identity - please remove, ensure they match, or disable use_aks_workload_identity")
		}
		tenantId = aksTenantId
	}

	return &tenantId, nil
}
