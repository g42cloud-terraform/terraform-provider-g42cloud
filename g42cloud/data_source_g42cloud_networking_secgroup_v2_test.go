package g42cloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccG42CloudNetworkingSecGroupV2DataSource_basic(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccG42CloudNetworkingSecGroupV2DataSource_group(rName),
			},
			{
				Config: testAccG42CloudNetworkingSecGroupV2DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.g42cloud_networking_secgroup.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_networking_secgroup.secgroup_1", "name", rName),
				),
			},
		},
	})
}

func TestAccG42CloudNetworkingSecGroupV2DataSource_secGroupID(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccG42CloudNetworkingSecGroupV2DataSource_group(rName),
			},
			{
				Config: testAccG42CloudNetworkingSecGroupV2DataSource_secGroupID(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.g42cloud_networking_secgroup.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_networking_secgroup.secgroup_1", "name", rName),
				),
			},
		},
	})
}

func testAccCheckNetworkingSecGroupV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find security group data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Security group data source ID not set")
		}

		return nil
	}
}

func testAccG42CloudNetworkingSecGroupV2DataSource_group(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_networking_secgroup" "secgroup_1" {
  name        = "%s"
  description = "My neutron security group"
}
`, rName)
}

func testAccG42CloudNetworkingSecGroupV2DataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_networking_secgroup" "secgroup_1" {
  name = "${g42cloud_networking_secgroup.secgroup_1.name}"
}
`, testAccG42CloudNetworkingSecGroupV2DataSource_group(rName))
}

func testAccG42CloudNetworkingSecGroupV2DataSource_secGroupID(rName string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_networking_secgroup" "secgroup_1" {
  name = "${g42cloud_networking_secgroup.secgroup_1.name}"
}
`, testAccG42CloudNetworkingSecGroupV2DataSource_group(rName))
}
