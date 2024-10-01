---
page_title: "Configuring a Service Principal for managing EntraID"
---

# Configuring a Service Principal for managing EntraID

Terraform supports a number of different methods for authenticating to Azure

## Creating a Service Principal

A Service Principal represents an application within Microsoft EntraID whose properties and authentication tokens can be used as the `tenant_id`, `client_id` and `client_secret` fields needed by Terraform.

* [Authenticating to Azure using a Service Principal and a Client Secret](service_principal_client_secret.md)
* [Authenticating to Azure using a Service Principal and OpenID Connect](service_principal_oicd.md)

***

## Microsoft EntraID permissions

Now that you have created and authenticated an Application / Service Principal pair, you will need to grant some permissions to administer Microsoft EntraID. You can choose either of the following methods to achieve similar results.

### API roles

Grant API roles to your Application and then grant consent for your Service Principal to access the APIs in its own capacity (i.e. not on behalf of a user).

Navigate to the [Microsoft Entra ID overview](https://portal.azure.com/#blade/Microsoft\_AAD\_IAM/ActiveDirectoryMenuBlade/Overview) and select the [App Registrations blade](https://portal.azure.com/#blade/Microsoft\_AAD\_IAM/ActiveDirectoryMenuBlade/RegisteredApps). Locate your registered Application and click on its display name to manage it.

Go to the API Permissions blade for the Application and click the "Add a permission" button. In the pane that opens, select `Microsoft Graph`.

Choose `Application Permissions` for the permission type and check the permissions you would like to assign. The permissions you need will depend on which directory objects you wish to manage with Terraform. The following table show the required permissions for some common resources:

| Resource(s)                                   | Description                                                       |
| --------------------------------------------- | ----------------------------------------------------------------- |
| `Application.ReadWrite.All`                   | Read and write all applications                                   |
| `DeviceManagementApps.ReadWrite.All`          | Read and write Microsoft Intune apps                              |
| `DeviceManagementConfiguration.ReadWrite.All` | Read and write Microsoft Intune device configuration and policies |
| `DeviceManagementServiceConfig.ReadWrite.All` | Read and write Microsoft Intune configuration                     |
| `Directory.ReadWrite.All`                     | Read and write directory data                                     |
| `Group.ReadWrite.All`                         | Read and write all groups                                         |
| `Policy.Read.All`                             | Read your organization's policies                                 |
| `Policy.ReadWrite.ConditionalAccess`          | Read and write your organization's conditional access policies    |
| `User.ReadWrite.All`                          | Read and write all users' full profiles                           |
| `CloudPC.ReadWrite.All`                       | Read and write Windows365 configurations                          |



!> **Caution** After assigning permissions, you will need to grant consent for the service principal to utilise them. The easiest way to do this is by clicking the Grant Admin Consent button in the same API Permissions pane. You will need to be signed in to the Portal as a Global Administrator.

-> **Permissions for other resources** If the resource you are using is not shown in the above table, consult the documentation page for the resource for a guide to the required permissions.

The Application now has the necessary permissions to administer your Microsoft EntraID tenant.
