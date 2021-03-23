package g42cloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNatGatewayDataSource_basic(t *testing.T) {
	natgateway := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNatGatewayV2DataSource_basic(natgateway),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatGatewayV2DataSourceID("data.g42cloud_nat_gateway.nat_by_name"),
					testAccCheckNatGatewayV2DataSourceID("data.g42cloud_nat_gateway.nat_by_id"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_nat_gateway.nat_by_name", "name", natgateway),
					resource.TestCheckResourceAttr(
						"data.g42cloud_nat_gateway.nat_by_id", "name", natgateway),
				),
			},
		},
	})
}

func testAccCheckNatGatewayV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find natgateway data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("NatGateway data source ID not set")
		}

		return nil
	}
}

func testAccNatGatewayV2DataSource_basic(name string) string {
	return fmt.Sprintf(`
resource "g42cloud_vpc" "vpc_1" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "g42cloud_vpc_subnet" "subnet_1" {
  name       = "%s"
  cidr       = "192.168.199.0/24"
  gateway_ip = "192.168.199.1"
  vpc_id     = g42cloud_vpc.vpc_1.id
}

resource "g42cloud_nat_gateway" "nat_1" {
  name                  = "%s"
  description           = "test for terraform"
  spec                  = "1"
  internal_network_id   = g42cloud_vpc_subnet.subnet_1.id
  router_id             = g42cloud_vpc.vpc_1.id
}

data "g42cloud_nat_gateway" "nat_by_name" {
  name = g42cloud_nat_gateway.nat_1.name
}

data "g42cloud_nat_gateway" "nat_by_id" {
  id = g42cloud_nat_gateway.nat_1.id
}
`, name, name, name)
}
