package g42cloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVpcSubnetV1DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcSubnetV1Config(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceVpcSubnetV1Check("data.g42cloud_vpc_subnet.by_id", rName, "172.16.8.0/24", "172.16.8.1"),
					testAccDataSourceVpcSubnetV1Check("data.g42cloud_vpc_subnet.by_cidr", rName, "172.16.8.0/24", "172.16.8.1"),
					testAccDataSourceVpcSubnetV1Check("data.g42cloud_vpc_subnet.by_name", rName, "172.16.8.0/24", "172.16.8.1"),
					testAccDataSourceVpcSubnetV1Check("data.g42cloud_vpc_subnet.by_vpc_id", rName, "172.16.8.0/24", "172.16.8.1"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_vpc_subnet.by_id", "status", "ACTIVE"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_vpc_subnet.by_id", "dhcp_enable", "true"),
				),
			},
		},
	})
}

func testAccDataSourceVpcSubnetV1Check(n, name, cidr, gateway_ip string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", n)
		}

		subnetRs, ok := s.RootModule().Resources["g42cloud_vpc_subnet.test"]
		if !ok {
			return fmt.Errorf("can't find g42cloud_vpc_subnet.test in state")
		}

		attr := rs.Primary.Attributes

		if attr["id"] != subnetRs.Primary.Attributes["id"] {
			return fmt.Errorf(
				"id is %s; want %s",
				attr["id"],
				subnetRs.Primary.Attributes["id"],
			)
		}

		if attr["cidr"] != cidr {
			return fmt.Errorf("bad subnet cidr %s, expected: %s", attr["cidr"], cidr)
		}
		if attr["name"] != name {
			return fmt.Errorf("bad subnet name %s", attr["name"])
		}
		if attr["gateway_ip"] != gateway_ip {
			return fmt.Errorf("bad subnet gateway_ip %s", attr["gateway_ip"])
		}

		return nil
	}
}

func testAccDataSourceVpcSubnetV1Config(rName string) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

resource "g42cloud_vpc" "test" {
  name = "%s"
  cidr = "172.17.7.0/24"
}

resource "g42cloud_vpc_subnet" "test" {
  name              = "%s"
  cidr              = "172.17.7.0/24"
  gateway_ip        = "172.17.7.1"
  vpc_id            = g42cloud_vpc.test.id

  availability_zone = data.g42cloud_availability_zones.test.names[0]
}

data "g42cloud_vpc_subnet" "by_id" {
  id = g42cloud_vpc_subnet.test.id
}

data "g42cloud_vpc_subnet" "by_cidr" {
  cidr = g42cloud_vpc_subnet.test.cidr
}

data "g42cloud_vpc_subnet" "by_name" {
  name = g42cloud_vpc_subnet.test.name
}

data "g42cloud_vpc_subnet" "by_vpc_id" {
  vpc_id = g42cloud_vpc_subnet.test.vpc_id
}
`, rName, rName)
}
