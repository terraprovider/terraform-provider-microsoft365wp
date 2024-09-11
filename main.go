package main

import (
	"context"
	"flag"
	"log"
	"terraform-provider-microsoft365wp/workplace"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs validate
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name microsoft365wp

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(context.Background(), workplace.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/terraprovider/microsoft365wp",
		Debug:   debug,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}
