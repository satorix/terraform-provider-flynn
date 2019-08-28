package main

import (
	"github.com/hashicorp/terraform/plugin"
	"gitlab.iexposure.com/satorix/terraform/terraform-provider-flynn/flynn"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: flynn.Provider})
}
