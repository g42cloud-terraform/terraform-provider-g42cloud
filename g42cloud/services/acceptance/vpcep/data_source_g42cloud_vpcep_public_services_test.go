package vpcep

import (
	"testing"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVPCEPPublicServicesDataSource_Basic(t *testing.T) {
	resourceName := "data.g42cloud_vpcep_public_services.services"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEPPublicServicesDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "services.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "services.0.service_name"),
					resource.TestCheckResourceAttrSet(resourceName, "services.0.service_type"),
				),
			},
		},
	})
}

var testAccVPCEPPublicServicesDataSourceBasic = `
data "g42cloud_vpcep_public_services" "services" {
  service_name = "dns"
}
`
