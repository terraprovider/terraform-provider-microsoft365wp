---
page_title: "microsoft365wp_named_location Data Source - microsoft365wp"
subcategory: "MS Graph: Conditional access"
---

# microsoft365wp_named_location (Data Source)

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
Conditional Access policies using it are not affected): replace the `azuread_named_location` with a
`microsoft365wp_named_location` and add a `moved` block from the old to the new resource address (requires
Terraform 1.8 or later, or OpenTofu 1.10 or later). The provider takes the `id` from the old state and reads
everything else from MS Graph, so the plan shows the move and at most an in-place update. A `msgraph_resource`
(provider `microsoft/msgraph`) with the `url` `identity/conditionalAccess/namedLocations` can be moved the
same way.

Alternatively, import the named location into a `microsoft365wp_named_location` with an `import` block (using
the `id` of the existing named location), and remove the `azuread_named_location` from state with a `removed`
block that has `lifecycle { destroy = false }`.

The `moved`, `import` and `removed` blocks can be deleted after the next apply. See the examples below.

Named locations whose actions are protected with Entra protected actions (see
`microsoft365wp_unified_rbac_resource_action`) can only be managed by app-only clients if the provider option
`enable_claims_challenge` is set.

## Documentation Disclaimer

Please note that almost all information on this page has been sourced literally from the official Microsoft Graph API 
documentation and therefore is governed by Microsoft and not by the publishers of this provider.  
All supplements authored by the publishers of this provider have been explicitly marked as such.

## Query Filters (if Supported)

If filtering by attribute values is supported (see schema below), then values set by the practitioner inside the config 
will be translated to a respective OData `$filter` clause. For string attributes (except enumerations!), simple 
wildcards (`*`) are supported at the start and/or the end of the attribute value or else exactly once inside and will be 
translated to corresponding OData predicates and functions (i.e. `eq`, `startswith`, `endswith` and `contains`). 
Multiple filter clauses will be combined using ` and `.  
If supported (see schema below), the attributes `odata_filter`, `odata_orderby` and `odata_top` can also be used to 
provide literal values for the respective OData options.

If this is a data source that returns a single element (singular data source), then the resulting OData query must 
result in exactly one returned entity! If supported, `odata_top = 1` and `odata_orderby` may be used to select a single 
entity from a list.

Please note that in the end all OData clauses/options will have to be interpreted by MS Graph, so MS Graph might impose 
further restrictions on what functionality may be used in practice.

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


data "microsoft365wp_named_location" "one" {
  id = "01234567-89ab-cdef-0123-456789abcdef"
}

output "microsoft365wp_named_location" {
  value = data.microsoft365wp_named_location.one
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (String) Identifier of a namedLocation object.
- `odata_filter` (String) Literal OData `$filter` value to pass to MS Graph.
- `odata_orderby` (String) Literal OData `$orderby` value to pass to MS Graph.
- `odata_top` (Number) Literal OData `$top` value to pass to MS Graph.

### Read-Only

- `compliant_network` (Attributes) Represents a Microsoft Entra ID named location defined by Global Secure Access. Automatically created with the name "All Compliant Network Locations" when you enable Global Secure Access signaling for Conditional Access. Named locations are custom rules that define network locations that can then be used in a Conditional Access policy. <br/> For more information, see [Enable compliant network check with Conditional Access](https://learn.microsoft.com/en-us/entra/global-secure-access/how-to-compliant-network). Also see [Microsoft docs for compliantNetworkNamedLocation](https://learn.microsoft.com/en-us/graph/api/resources/compliantnetworknamedlocation?view=graph-rest-beta). <br> (see [below for nested schema](#nestedatt--compliant_network))
- `country` (Attributes) Represents a Microsoft Entra ID named location defined by countries and regions. Named locations are custom rules that define network locations which can then be used in a Conditional Access policy. Also see [Microsoft docs for countryNamedLocation](https://learn.microsoft.com/en-us/graph/api/resources/countrynamedlocation?view=graph-rest-beta). <br> (see [below for nested schema](#nestedatt--country))
- `created_date_time` (String) The Timestamp type represents creation date and time of the location using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`.
- `display_name` (String) Human-readable name of the location.
- `ip` (Attributes) Represents a Microsoft Entra ID named location defined by IP ranges. Named locations are custom rules that define network locations that can then be used in a Conditional Access policy. Also see [Microsoft docs for ipNamedLocation](https://learn.microsoft.com/en-us/graph/api/resources/ipnamedlocation?view=graph-rest-beta). <br> (see [below for nested schema](#nestedatt--ip))
- `modified_date_time` (String) The Timestamp type represents last modified date and time of the location using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`.
- `private_link` (Attributes) (see [below for nested schema](#nestedatt--private_link))
- `service_tag` (Attributes) (see [below for nested schema](#nestedatt--service_tag))

<a id="nestedatt--compliant_network"></a>
### Nested Schema for `compliant_network`

Read-Only:

- `compliant_network_type` (String) Type of compliant network. Currently the only possible value is `allTenantCompliantNetworks`. <br/> _Provider_ allowed values are: `allTenantCompliantNetworks`, `unknownFutureValue`.
- `is_trusted` (Boolean) `true` if this location is explicitly trusted. Optional. <br/>


<a id="nestedatt--country"></a>
### Nested Schema for `country`

Read-Only:

- `countries_and_regions` (Set of String) List of countries and/or regions in two-letter format specified by ISO 3166-2.
- `country_lookup_method` (String) Determines what method is used to decide which country the user is located in. Possible values are `clientIpAddress` (default) and `authenticatorAppGps`. <br/> _Provider_ allowed values are: `clientIpAddress`, `authenticatorAppGps`, `unknownFutureValue`.
- `include_unknown_countries_and_regions` (Boolean) `true` if IP addresses that don't map to a country or region should be included in the named location. Optional. <br/>


<a id="nestedatt--ip"></a>
### Nested Schema for `ip`

Read-Only:

- `ip_ranges` (Attributes Set) List of IP address ranges in IPv4 CIDR format (for example, 1.2.3.4/32) or any allowable IPv6 format from IETF RFC5969. Required. / IP range base class for representing IPV4, IPV6 address ranges. Also see [Microsoft docs for ipRange](https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-iprange?view=graph-rest-beta). <br> (see [below for nested schema](#nestedatt--ip--ip_ranges))
- `is_trusted` (Boolean) `true` if this location is explicitly trusted. Optional. <br/>

<a id="nestedatt--ip--ip_ranges"></a>
### Nested Schema for `ip.ip_ranges`

Read-Only:

- `v4_cidr` (Attributes) Represents an IPv4 range using the Classless Inter-Domain Routing (CIDR) notation. Also see [Microsoft docs for iPv4CidrRange](https://learn.microsoft.com/en-us/graph/api/resources/ipv4cidrrange?view=graph-rest-beta). <br> (see [below for nested schema](#nestedatt--ip--ip_ranges--v4_cidr))
- `v6_cidr` (Attributes) Represents an IPv6 range using the Classless Inter-Domain Routing (CIDR) notation. Also see [Microsoft docs for iPv6CidrRange](https://learn.microsoft.com/en-us/graph/api/resources/ipv6cidrrange?view=graph-rest-beta). <br> (see [below for nested schema](#nestedatt--ip--ip_ranges--v6_cidr))

<a id="nestedatt--ip--ip_ranges--v4_cidr"></a>
### Nested Schema for `ip.ip_ranges.v4_cidr`

Read-Only:

- `cidr_address` (String) IPv4 address in CIDR notation. Not nullable.


<a id="nestedatt--ip--ip_ranges--v6_cidr"></a>
### Nested Schema for `ip.ip_ranges.v6_cidr`

Read-Only:

- `cidr_address` (String) IPv6 address in CIDR notation. Not nullable.




<a id="nestedatt--private_link"></a>
### Nested Schema for `private_link`

Read-Only:

- `is_trusted` (Boolean)
- `private_link_resource_policy_ids` (Set of String)


<a id="nestedatt--service_tag"></a>
### Nested Schema for `service_tag`

Read-Only:

- `is_trusted` (Boolean)
- `service_tags` (Set of String)
