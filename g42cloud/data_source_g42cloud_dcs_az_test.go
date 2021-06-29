package g42cloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDcsAZV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsAZV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsAZV1DataSourceID("data.g42cloud_dcs_az.az1"),
					resource.TestCheckResourceAttr("data.g42cloud_dcs_az.az1", "code", "ae-ad-1a"),
					resource.TestCheckResourceAttr("data.g42cloud_dcs_az.az1", "port", "8002"),
				),
			},
		},
	})
}

func testAccCheckDcsAZV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Dcs az data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Dcs az data source ID not set")
		}

		return nil
	}
}

var testAccDcsAZV1DataSource_basic = fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

data "g42cloud_dcs_az" "az1" {
  code = data.g42cloud_availability_zones.test.names[0]
  port = "8002"
}
`)
