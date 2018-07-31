package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"

	"github.com/cedexis/terraform-provider-citrixitm/citrixitm"
)

func main() {
	plugin.Serve(
		&plugin.ServeOpts{
			ProviderFunc: func() terraform.ResourceProvider {
				return citrixitm.Provider()
			},
		})
}
