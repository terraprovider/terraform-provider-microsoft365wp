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


resource "microsoft365wp_browser_site_list" "test" {
  display_name = "Contoso IE mode sites"
  description  = "Legacy line-of-business apps"

  sites = [
    {
      web_url = "erp.contoso.com"
    },
    {
      web_url            = "https://intranet.contoso.com/legacy"
      target_environment = "internetExplorerMode"
      compatibility_mode = "internetExplorer8Enterprise"
      merge_type         = "noMerge"
      allow_redirect     = true
      comment            = "needs IE8 document mode"
    },
    {
      web_url            = "https://intranet.contoso.com/new"
      target_environment = "microsoftEdge"
    },
  ]

  shared_cookies = [
    {
      display_name       = "ASP.NET_SessionId"
      host_or_domain     = "erp.contoso.com"
      path               = "/"
      host_only          = true
      source_environment = "both"
    },
    {
      # domain cookie: host_or_domain must start with a dot
      display_name       = "LegacyAuth"
      host_or_domain     = ".contoso.com"
      path               = "/"
      source_environment = "internetExplorer11"
    },
  ]
}
