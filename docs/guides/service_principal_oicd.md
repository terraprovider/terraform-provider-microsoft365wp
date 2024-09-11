---
page_title: "Authenticating using a Service Principal and OpenID Connect"
subcategory: "Authentication"
---

# Authenticating using a Service Principal and OpenID Connect

* [Authenticating to Azure using a Service Principal and a Client Secret](service\_principal\_client\_secret.md)
* Authenticating to Azure using a Service Principal and OpenID Connect (covered in this guide)

***

Once you have configured a Service Principal as described in this guide, you should follow the [Configuring a Service Principal for managing EntraID](../service-principal-configuration.md) guide to grant the Service Principal necessary permissions to create and modify Microsoft EntraID objects.

***

## Setting up an Application and Service Principal

A Service Principal is a security principal within Microsoft EntraID which can be granted permissions to manage objects in Microsoft EntraID. To authenticate with a Service Principal, you will need to create an Application object within Microsoft EntraID, which you will use as a means of authentication, either using a [Client Secret](service\_principal\_client\_secret.md), Client Certificate, or OpenID Connect (which is documented in this guide). This can be done using the Azure Portal.

This guide will cover how to create an Application and linked Service Principal, and then how to assign federated identity credentials to the Application so that it can be used for authentication via OpenID Connect.

## Creating the Application and Service Principal

We're going to create the Application in the Azure Portal - to do this, navigate to the [**Microsoft EntraID** overview](https://portal.azure.com/#blade/Microsoft\_AAD\_IAM/ActiveDirectoryMenuBlade/Overview) within the [Azure Portal](https://portal.azure.com/) - then select the [**App Registrations** blade](https://portal.azure.com/#blade/Microsoft\_AAD\_IAM/ActiveDirectoryMenuBlade/RegisteredApps/RegisteredApps/Overview). Click the **New registration** button at the top to add a new Application within Microsoft EntraID. On this page, set the following values then press **Create**:

* **Name** - this is a friendly identifier and can be anything (e.g. "Terraform")
* **Supported Account Types** - this should be set to "Accounts in this organisational directory only (single tenant)"
* **Redirect URI** - you should choose "Web" in for the URI type. the actual value can be left blank

At this point the newly created Microsoft EntraID application should be visible on-screen - if it's not, navigate to the [**App Registration** blade](https://portal.azure.com/#blade/Microsoft\_AAD\_IAM/ActiveDirectoryMenuBlade/RegisteredApps/RegisteredApps/Overview) and select the Microsoft EntraID application.

At the top of this page, you'll need to take note of the "Application (client) ID" and the "Directory (tenant) ID", which you can use for the values of `client_id` and `tenant_id` respectively.

## Configure Microsoft EntraID Application to Trust a GitHub Repository

An application will need a federated credential specified for each GitHub Environment, Branch Name, Pull Request, or Tag based on your use case. For this example, we'll give permission to `main` branch workflow runs.

-> **Tip:** You can also configure the Application using the [azuread\_application](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/resources/application) and [azuread\_application\_federated\_identity\_credential](https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/resources/application\_federated\_identity\_credential) resources in the AzureAD Terraform Provider.

### Via the Portal

On the Microsoft EntraID application page, go to **Certificates and secrets**.

In the Federated credentials tab, select Add credential. The Add a credential blade opens. In the **Federated credential scenario** drop-down box select **GitHub actions deploying Azure resources**.

Specify the **Organization** and **Repository** for your GitHub Actions workflow. For **Entity type**, select **Environment**, **Branch**, **Pull request**, or **Tag** and specify the value. The values must exactly match the configuration in the GitHub workflow. For our example, let's select **Branch** and specify `main`.

Add a **Name** for the federated credential.

The **Issuer**, **Audiences**, and **Subject identifier** fields autopopulate based on the values you entered.

Click **Add** to configure the federated credential.

### Via the Azure API

```sh
az rest --method POST \
        --uri https://graph.microsoft.com/beta/applications/${APP_OBJ_ID}/federatedIdentityCredentials \
        --headers Content-Type='application/json' \
        --body @body.json
```

Where the body is:

```json
{
  "name":"${REPO_NAME}-pull-request",
  "issuer":"https://token.actions.githubusercontent.com",
  "subject":"repo:${REPO_OWNER}/${REPO_NAME}:refs:refs/heads/main",
  "description":"${REPO_OWNER} PR",
  "audiences":["api://AzureADTokenExchange"],
}
```

See the [official documentation](https://docs.microsoft.com/en-us/azure/active-directory/develop/workload-identity-federation-create-trust-github) for more details.

### Configure Microsoft EntraID Application to Trust a Generic Issuer

On the Microsoft EntraID application page, go to **Certificates and secrets**.

In the Federated credentials tab, select **Add credential**. The 'Add a credential' blade opens. Refer to the instructions from your OIDC provider for completing the form, before choosing a **Name** for the federated credential and clicking the **Add** button.

## Configuring Terraform to use OIDC

~> **Note:** If using the AzureRM Backend you may also need to configure OIDC there too, see [the documentation for the AzureRM Backend](https://www.terraform.io/language/settings/backends/azurerm) for more information.

As we've obtained the credentials for this Service Principal - it's possible to configure them in a few different ways.

When storing the credentials as Environment Variables, for example:

```bash
$ export ARM_CLIENT_ID="00000000-0000-0000-0000-000000000000"
$ export ARM_SUBSCRIPTION_ID="00000000-0000-0000-0000-000000000000"
$ export ARM_TENANT_ID="00000000-0000-0000-0000-000000000000"
$ export ARM_USE_OIDC=true
```

The provider will use the `ARM_OIDC_TOKEN` environment variable as an OIDC token. You can use this variable to specify the token provided by your OIDC provider. If your OIDC provider provides an ID token in a file, you can specify the path to this file with the `ARM_OIDC_TOKEN_FILE_PATH` environment variable.

When running in GitHub Actions, the provider will detect the `ACTIONS_ID_TOKEN_REQUEST_URL` and `ACTIONS_ID_TOKEN_REQUEST_TOKEN` environment variables set by the GitHub Actions runtime. You can also specify the `ARM_OIDC_REQUEST_TOKEN` and `ARM_OIDC_REQUEST_URL` environment variables.

For GitHub Actions workflows, you'll need to ensure the workflow has `write` permissions for the `id-token`.

```yaml
permissions:
  id-token: write
  contents: read
```

For more information about OIDC in GitHub Actions, see [official documentation](https://docs.github.com/en/actions/deployment/security-hardening-your-deployments/configuring-openid-connect-in-cloud-providers).

The following Terraform and Provider blocks can be specified - where `1.00.0` is the version of the Provider that you'd like to use:

```hcl
# We strongly recommend using the required_providers block to set the
# Azure Provider source and version being used
terraform {
  required_providers {
    microsoft365wp = {
      source  = "terraprovider/microsoft365wp"
      version = "=1.00.0"
    }
  }
}
# Configure the Microsoft Azure Provider
provider "microsoft365wp" {
  use_oidc = true # or use the environment variable "ARM_USE_OIDC=true"
  features {}
}
```



~> **Note:** If using the AzureRM Backend you may also need to configure OIDC there too, see [the documentation for the AzureRM Backend](https://www.terraform.io/language/settings/backends/azurerm) for more information.

More information on [the fields supported in the Provider block can be found here](../../index.md#argument-reference).

At this point running either `terraform plan` or `terraform apply` should allow Terraform to run using the Service Principal to authenticate.

Next you should follow the [Configuring a Service Principal for managing Microsoft EntraID](../service-principal-configuration.md) guide to grant the Service Principal necessary permissions to create and modify Microsoft EntraID objects such as users and groups.
