---
page_title: "Manage existing resources"
---

description: How to manage existing Entra- and Intune objects using this provider

# Manage existing resources

## Overview

If you already have a pre-existing environment of EntraID groups, Intune policies etc., you will probably want to create a Terraform configuration for this existing setup.

We provide a Quickstart Template that helps you get a headstart with managing a Entra/Intune tenant. This can be used for a fresh start and helps by offering a best-practises design regarding variable and resource handling.&#x20;

See [https://github.com/terraprovider/terraform-provider-microsoft365wp-quickstart](https://github.com/terraprovider/terraform-provider-microsoft365wp-quickstart)

## Exporting existing Resources

The template includes a PowerShell tool to export existing resources directly into the structure offered by the Quickstart Template.

You can find a [step-by-step](https://github.com/terraprovider/terraform-provider-microsoft365wp-quickstart/blob/main/README.MD#importing--managing-an-existing-tenant) guide as part of the repository.

This allows you to export e.g. EntraID groups, Intune Policies and Filters for

* Backup purposes
* To deploy to another tenant
* Manage existing resources

The tool includes logic to map ("import") existing resources to the generated configuration.
