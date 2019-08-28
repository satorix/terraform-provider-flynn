package main

import (
	"github.com/hashicorp/terraform/plugin"
	"terraform-provider-flynn/flynn"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: flynn.Provider})
}
