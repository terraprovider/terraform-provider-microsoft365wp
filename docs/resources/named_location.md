---
page_title: "microsoft365wp_named_location Resource - microsoft365wp"
subcategory: "MS Graph: Conditional access"
---

# microsoft365wp_named_location (Resource)

This is the base class that represents a Microsoft Entra ID named location. Named locations are custom rules that define network locations which can then be used in a Conditional Access (CA) policy. <br/> Also see [Microsoft docs for namedLocation](https://learn.microsoft.com/en-us/graph/api/resources/namedlocation?view=graph-rest-beta).

_Provider_ Note: Named locations can be used in the `conditions.locations` of `microsoft365wp_conditional_access_policy`. Use
exactly one of `ip`, `country`, `compliant_network`, `service_tag` and `private_link` to choose the type of the
named location. Changing the type replaces the named location.

### Deletion
MS Graph does not allow to delete trusted named locations, so the provider sets `is_trusted` to `false` before
deleting them. Named locations that are still used by a Conditional Access policy cannot be deleted. If a
policy stops using a named location in the same apply that deletes it, the deletion is retried for up to a
minute until MS Graph has noticed the change of the policy. Note that Terraform deletes a named location
_before_ it updates the policies that use it, unless the named location has
`lifecycle { create_before_destroy = true }`. So set this for named locations that are used by policies, as
MS Graph will otherwise refuse to delete them (error 400, code 1178).

### Compliant Network Named Location
There can only be one named location of type `compliant_network` (it is also created by MS Graph when Global
Secure Access signaling for Conditional Access gets enabled). If it already exists, the provider will use (and
update) it instead of trying to create a new one. Its `display_name` must be `All Compliant Network locations`
(MS Graph does not allow to change it). MS Graph does not allow to delete it either, so on destroy the
provider only sets `is_trusted` to `false` and removes it from Terraform state.

### Migration from `azuread_named_location`
An existing named location can be moved to this resource without recreating it (it keeps its `id`, so
Conditional Access policies using it are not affected): import it into a `microsoft365wp_named_location` with
an `import` block (using the `id` of the existing named location), and remove the `azuread_named_location`
from state with a `removed` block that has `lifecycle { destroy = false }`. Both blocks can be deleted after
the next apply. See the example below.

Named locations whose actions are protected with Entra protected actions (see
`microsoft365wp_unified_rbac_resource_action`) can only be managed by app-only clients if the provider option
`enable_claims_challenge` is set.

## Documentation Disclaimer

Please note that almost all information on this page has been sourced literally from the official Microsoft Graph API 
documentation and therefore is governed by Microsoft and not by the publishers of this provider.  
All supplements authored by the publishers of this provider have been explicitly marked as such.

## Example Usage

```terraform
terraform {
  required_providers {
    microsoft365wp = {
      source = "terraprovider/microsoft365wp"
    }
  }
}

/*
.env
export ARM_TENANT_ID='...'
export ARM_CLIENT_ID='...'
export ARM_CLIENT_SECRET='...'
*/


resource "microsoft365wp_named_location" "ip" {
  display_name = "TF Test IP"
  ip = {
    is_trusted = true
    ip_ranges = [
      { v4_cidr = { cidr_address = "192.0.2.0/24" } },
      { v6_cidr = { cidr_address = "2001:db8::/48" } },
    ]
  }

  # for named locations used by Conditional Access policies: update the policies before deleting the named location
  lifecycle {
    create_before_destroy = true
  }
}

resource "microsoft365wp_named_location" "country" {
  display_name = "TF Test Country"
  country = {
    countries_and_regions = ["DE", "AT", "CH"]
  }
}

# There can only be one compliant network named location (Global Secure Access): if it exists already, it will be used.
# MS Graph does not allow to change its display name or to delete it.
resource "microsoft365wp_named_location" "compliant_network" {
  display_name = "All Compliant Network locations"
  compliant_network = {
    is_trusted = true
  }
}


# Migration from azuread_named_location without recreating the named location (OpenTofu / Terraform >= 1.7): import
# the existing named location into the new resource and remove the old resource from state without destroying it.
# After the next apply, both blocks can be deleted.

import {
  to = microsoft365wp_named_location.trusted_networks
  id = "01234567-89ab-cdef-0123-456789abcdef" # id of the existing named location
}

removed {
  from = azuread_named_location.trusted_networks
  lifecycle {
    destroy = false
  }
}

resource "microsoft365wp_named_location" "trusted_networks" {
  display_name = "Trusted Networks"
  ip = { # was: ip { ip_ranges = ["192.0.2.0/24"], trusted = true }
    is_trusted = true
    ip_ranges  = [{ v4_cidr = { cidr_address = "192.0.2.0/24" } }]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) Human-readable name of the location.

### Optional

- `compliant_network` (Attributes) Represents a Microsoft Entra ID named location defined by Global Secure Access. Automatically created with the name "All Compliant Network Locations" when you enable Global Secure Access signaling for Conditional Access. Named locations are custom rules that define network locations that can then be used in a Conditional Access policy. <br/> For more information, see [Enable compliant network check with Conditional Access](https://learn.microsoft.com/en-us/entra/global-secure-access/how-to-compliant-network). Also see [Microsoft docs for compliantNetworkNamedLocation](https://learn.microsoft.com/en-us/graph/api/resources/compliantnetworknamedlocation?view=graph-rest-beta). <br> (see [below for nested schema](#nestedatt--compliant_network))
- `country` (Attributes) Represents a Microsoft Entra ID named location defined by countries and regions. Named locations are custom rules that define network locations which can then be used in a Conditional Access policy. Also see [Microsoft docs for countryNamedLocation](https://learn.microsoft.com/en-us/graph/api/resources/countrynamedlocation?view=graph-rest-beta). <br> (see [below for nested schema](#nestedatt--country))
- `ip` (Attributes) Represents a Microsoft Entra ID named location defined by IP ranges. Named locations are custom rules that define network locations that can then be used in a Conditional Access policy. Also see [Microsoft docs for ipNamedLocation](https://learn.microsoft.com/en-us/graph/api/resources/ipnamedlocation?view=graph-rest-beta). <br> (see [below for nested schema](#nestedatt--ip))
- `private_link` (Attributes) (see [below for nested schema](#nestedatt--private_link))
- `service_tag` (Attributes) (see [below for nested schema](#nestedatt--service_tag))

### Read-Only

- `created_date_time` (String) The Timestamp type represents creation date and time of the location using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Read-only.
- `id` (String) Identifier of a namedLocation object. Read-only.
- `modified_date_time` (String) The Timestamp type represents last modified date and time of the location using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Read-only.

<a id="nestedatt--compliant_network"></a>
### Nested Schema for `compliant_network`

Optional:

- `compliant_network_type` (String) Type of compliant network. Currently the only possible value is `allTenantCompliantNetworks`. <br/> _Provider_ allowed values are: `allTenantCompliantNetworks`, `unknownFutureValue`. The _provider_ default value is `"allTenantCompliantNetworks"`.
- `is_trusted` (Boolean) `true` if this location is explicitly trusted. Optional. Default value is `false`. <br/> The _provider_ default value is `false`.


<a id="nestedatt--country"></a>
### Nested Schema for `country`

Required:

- `countries_and_regions` (Set of String) List of countries and/or regions in two-letter format specified by ISO 3166-2.

Optional:

- `country_lookup_method` (String) Determines what method is used to decide which country the user is located in. Possible values are `clientIpAddress` (default) and `authenticatorAppGps`. <br/> _Provider_ allowed values are: `clientIpAddress`, `authenticatorAppGps`, `unknownFutureValue`. The _provider_ default value is `"clientIpAddress"`.
- `include_unknown_countries_and_regions` (Boolean) `true` if IP addresses that don't map to a country or region should be included in the named location. Optional. Default value is `false`. <br/> The _provider_ default value is `false`.


<a id="nestedatt--ip"></a>
### Nested Schema for `ip`

Required:

- `ip_ranges` (Attributes Set) List of IP address ranges in IPv4 CIDR format (for example, 1.2.3.4/32) or any allowable IPv6 format from IETF RFC5969. Required. / IP range base class for representing IPV4, IPV6 address ranges. Also see [Microsoft docs for ipRange](https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-iprange?view=graph-rest-beta). <br> (see [below for nested schema](#nestedatt--ip--ip_ranges))

Optional:

- `is_trusted` (Boolean) `true` if this location is explicitly trusted. Optional. Default value is `false`. <br/> The _provider_ default value is `false`.

<a id="nestedatt--ip--ip_ranges"></a>
### Nested Schema for `ip.ip_ranges`

Optional:

- `v4_cidr` (Attributes) Represents an IPv4 range using the Classless Inter-Domain Routing (CIDR) notation. Also see [Microsoft docs for iPv4CidrRange](https://learn.microsoft.com/en-us/graph/api/resources/ipv4cidrrange?view=graph-rest-beta). <br> (see [below for nested schema](#nestedatt--ip--ip_ranges--v4_cidr))
- `v6_cidr` (Attributes) Represents an IPv6 range using the Classless Inter-Domain Routing (CIDR) notation. Also see [Microsoft docs for iPv6CidrRange](https://learn.microsoft.com/en-us/graph/api/resources/ipv6cidrrange?view=graph-rest-beta). <br> (see [below for nested schema](#nestedatt--ip--ip_ranges--v6_cidr))

<a id="nestedatt--ip--ip_ranges--v4_cidr"></a>
### Nested Schema for `ip.ip_ranges.v4_cidr`

Required:

- `cidr_address` (String) IPv4 address in CIDR notation. Not nullable.


<a id="nestedatt--ip--ip_ranges--v6_cidr"></a>
### Nested Schema for `ip.ip_ranges.v6_cidr`

Required:

- `cidr_address` (String) IPv6 address in CIDR notation. Not nullable.




<a id="nestedatt--private_link"></a>
### Nested Schema for `private_link`

Optional:

- `is_trusted` (Boolean) The _provider_ default value is `false`.
- `private_link_resource_policy_ids` (Set of String) The _provider_ default value is `[]`.


<a id="nestedatt--service_tag"></a>
### Nested Schema for `service_tag`

Optional:

- `is_trusted` (Boolean) The _provider_ default value is `false`.
- `service_tags` (Set of String) The _provider_ default value is `[]`.

## Import

Import is supported using an `import` block (or the `terraform import` command) with the value of `id`:

```terraform
import {
  to = microsoft365wp_named_location.example
  id = "<id>"
}
```
