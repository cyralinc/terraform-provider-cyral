package main

import (
	"github.com/cyralinc/terraform-provider-cyral/cyral"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cyral.Provider,
	})
}
