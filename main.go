package main

import (
	"github.com/DeviaVir/terraform-provider-tfe/tfe"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: tfe.Provider})
}
