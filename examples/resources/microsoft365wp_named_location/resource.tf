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
