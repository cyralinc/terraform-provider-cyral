package main

import (
	"github.com/cyralinc/terraform-provider-cyral/src/cyral"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cyral.Provider,
	})
}
