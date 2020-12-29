package main

import (
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: g42cloud.Provider})
}
