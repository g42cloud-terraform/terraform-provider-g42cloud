package vpc

import (
	"fmt"
	"testing"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVpcRouteTableDataSource_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dataSourceName := "data.g42cloud_vpc_route_table.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRouteTable_default(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "default", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.#", "1"),
				),
			},
		},
	})
}

func TestAccVpcRouteTableDataSource_custom(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dataSourceName := "data.g42cloud_vpc_route_table.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRouteTable_custom(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "default", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceRouteTable_base(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_vpc" "test" {
  name = "%s"
  cidr = "172.16.0.0/16"
}

resource "g42cloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "172.16.10.0/24"
  gateway_ip = "172.16.10.1"
  vpc_id     = g42cloud_vpc.test.id
}
`, rName, rName)
}

func testAccDataSourceRouteTable_default(rName string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_vpc_route_table" "test" {
  vpc_id = g42cloud_vpc.test.id

  depends_on = [
    g42cloud_vpc_subnet.test
  ]
}
`, testAccDataSourceRouteTable_base(rName))
}

func testAccDataSourceRouteTable_custom(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_vpc_route_table" "test" {
  name        = "%s"
  vpc_id      = g42cloud_vpc.test.id
  description = "created by terraform"
}

data "g42cloud_vpc_route_table" "test" {
  vpc_id = g42cloud_vpc.test.id
  name   = "%s"

  depends_on = [
    g42cloud_vpc_route_table.test
  ]
}
`, testAccDataSourceRouteTable_base(rName), rName, rName)
}
