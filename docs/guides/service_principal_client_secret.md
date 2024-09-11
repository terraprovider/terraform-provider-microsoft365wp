---
page_title: "Authenticating using a Service Principal with a Client Secret"
subcategory: "Authentication"
---

# Authenticating using a Service Principal with a Client Secret

* Authenticating to Azure using a Service Principal and a Client Secret (covered in this guide)
* [Authenticating to Azure using a Service Principal and OpenID Connect](service\_principal\_oicd.md)

***

Once you have configured a Service Principal as described in this guide, you should follow the [Configuring a Service Principal for managing EntraID](../service-principal-configuration.md) guide to grant the Service Principal necessary permissions to create and modify Microsoft EntraID objects.

***

## Setting up an Application and Service Principal

A Service Principal is a security principal within Microsoft EntraID which can be granted permissions to manage objects in Microsoft EntraID. To authenticate with a Service Principal, you will need to create an Application object within Microsoft EntraID, which you will use as a means of authentication.

This guide will cover how to create an Application and linked Service Principal, and then how to generate a Client Secret for the Application so that it can be used for authentication.

***

## Creating the Application and Service Principal

We're going to create the Application in the Azure Portal - to do this navigate to the [**Microsoft EntraID** overview](https://portal.azure.com/#blade/Microsoft\_AAD\_IAM/ActiveDirectoryMenuBlade/Overview) within the [Azure Portal](https://portal.azure.com/) - then select the \[**App Registrations** blade]\[azure-portal-applications-blade]. Click the **New registration** button at the top to add a new Application within Microsoft EntraID. On this page, set the following values then press **Create**:

* **Name** - this is a friendly identifier and can be anything (e.g. "Terraform")
* **Supported Account Types** - this should be set to "Accounts in this organisational directory only (single tenant)"
* **Redirect URI** - you should choose "Web" in for the URI type. the actual value can be left blank

At this point the newly created Microsoft EntraID application should be visible on-screen - if it's not, navigate to the \[**App Registration** blade]\[azure-portal-applications-blade] and select the Microsoft EntraID application.

At the top of this page, you'll need to take note of the "Application (client) ID" and the "Directory (tenant) ID", which you can use for the values of `client_id` and `tenant_id` respectively.

### Generating a Client Secret for the Microsoft EntraID Application

Now that the Microsoft EntraID Application exists we can create a Client Secret which can be used for authentication - to do this select **Certificates & secrets**. This screen displays the Certificates and Client Secrets (i.e. passwords) which are associated with this Microsoft EntraID Application.

Click the "New client secret" button, then enter a short description, choose an expiry period and click "Add". Once the Client Secret has been generated it will be displayed on screen - _the secret is only displayed once_ so **be sure to copy it now** (otherwise you will need to regenerate a new one). This is the `client_secret` you will need.

***

## Creating the Application and Service Principal

We're going to create the Application in the Azure Portal - to do this navigate to the [**Microsoft EntraID** overview](https://portal.azure.com/#blade/Microsoft\_AAD\_IAM/ActiveDirectoryMenuBlade/Overview) within the [Azure Portal](https://portal.azure.com/) - then select the \[**App Registrations** blade]\[azure-portal-applications-blade]. Click the **New registration** button at the top to add a new Application within Microsoft EntraID. On this page, set the following values then press **Create**:

* **Name** - this is a friendly identifier and can be anything (e.g. "Terraform")
* **Supported Account Types** - this should be set to "Accounts in this organisational directory only (single tenant)"
* **Redirect URI** - you should choose "Web" in for the URI type. the actual value can be left blank

At this point the newly created Microsoft EntraID application should be visible on-screen - if it's not, navigate to the \[**App Registration** blade]\[azure-portal-applications-blade] and select the Microsoft EntraID application.

At the top of this page, you'll need to take note of the "Application (client) ID" and the "Directory (tenant) ID", which you can use for the values of `client_id` and `tenant_id` respectively.

### Generating a Client Secret for the Microsoft EntraID Application

Now that the Microsoft EntraID Application exists we can create a Client Secret which can be used for authentication - to do this select **Certificates & secrets**. This screen displays the Certificates and Client Secrets (i.e. passwords) which are associated with this Microsoft EntraID Application.

Click the "New client secret" button, then enter a short description, choose an expiry period and click "Add". Once the Client Secret has been generated it will be displayed on screen - _the secret is only displayed once_ so **be sure to copy it now** (otherwise you will need to regenerate a new one). This is the `client_secret` you will need.

***

## Configuring Terraform to use the Client Secret

Now we have obtained the necessary credentials, it's possible to configure Terraform in a few different ways.

### 1.Environment Variables

```sh
export ARM_CLIENT_ID="00000000-0000-0000-0000-000000000000"
export ARM_CLIENT_SECRET="12345678-0000-0000-0000-000000000000"
export ARM_TENANT_ID="10000000-0000-0000-0000-000000000000"
export ARM_SUBSCRIPTION_ID="20000000-0000-0000-0000-000000000000"
```

At this point running either `terraform plan` or `terraform apply` should allow Terraform to authenticate using the Client Secret.

Next you should follow the [Configuring a Service Principal for managing Microsoft EntraID](../service-principal-configuration.md) guide to grant the Service Principal necessary permissions to create and modify Microsoft EntraID objects such as users and groups.

## Provider Block

It's also possible to configure these variables either directly, or from variables, in your provider block, like so:

~> **Caution** We recommend not defining these variables in-line since they could easily be checked into Source Control.

```hcl
# We strongly recommend using the required_providers block to set the
# Workplace Provider source and version being used
terraform {
  required_providers {
    microsoft365wp = {
      source  = "terraprovider/microsoft365wp"
      version = "1.0.0"
    }
  }
}

# Configure the Workplace Provider
provider "microsoft365wp" {
  client_id     = "00000000-0000-0000-0000-000000000000"
  client_secret = var.client_secret
  tenant_id     = "10000000-2000-3000-4000-500000000000"
}
```

More information on [the fields supported in the Provider block can be found here](../../index.md#argument-reference).

At this point running either `terraform plan` or `terraform apply` should allow Terraform to authenticate using the Client Secret.

Next you should follow the [Configuring a Service Principal for managing Microsoft EntraID](../service-principal-configuration.md) guide to grant the Service Principal necessary permissions to create and modify Microsoft EntraID objects such as users and groups.
