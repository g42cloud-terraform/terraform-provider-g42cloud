package g42cloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVpcV1DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcV1Config(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceVpcV1Check("data.g42cloud_vpc.by_id", rName),
					testAccDataSourceVpcV1Check("data.g42cloud_vpc.by_cidr", rName),
					testAccDataSourceVpcV1Check("data.g42cloud_vpc.by_name", rName),
					resource.TestCheckResourceAttr(
						"data.g42cloud_vpc.by_id", "shared", "false"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_vpc.by_id", "status", "OK"),
				),
			},
		},
	})
}

func testAccDataSourceVpcV1Check(n, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", n)
		}

		vpcRs, ok := s.RootModule().Resources["g42cloud_vpc.test"]
		if !ok {
			return fmt.Errorf("can't find g42cloud_vpc.test in state")
		}

		attr := rs.Primary.Attributes

		if attr["id"] != vpcRs.Primary.Attributes["id"] {
			return fmt.Errorf(
				"id is %s; want %s",
				attr["id"],
				vpcRs.Primary.Attributes["id"],
			)
		}

		if attr["name"] != rName {
			return fmt.Errorf("bad vpc name %s", attr["name"])
		}

		return nil
	}
}

func testAccDataSourceVpcV1Config(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_vpc" "test" {
  name = "%s"
  cidr = "172.16.10.0/24"
}

data "g42cloud_vpc" "by_id" {
  id = g42cloud_vpc.test.id
}

data "g42cloud_vpc" "by_cidr" {
  cidr = g42cloud_vpc.test.cidr
}

data "g42cloud_vpc" "by_name" {
  name = g42cloud_vpc.test.name
}
`, rName)
}
