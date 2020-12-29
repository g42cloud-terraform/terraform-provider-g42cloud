package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: g42cloud.Provider})
}
