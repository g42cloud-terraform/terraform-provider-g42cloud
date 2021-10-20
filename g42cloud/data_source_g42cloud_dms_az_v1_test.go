package g42cloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDmsAZV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsAZV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsAZV1DataSourceID("data.g42cloud_dms_az.az1"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_dms_az.az1", "port", "8002"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_dms_az.az1", "code", "ae-ad-1a"),
				),
			},
		},
	})
}

func testAccCheckDmsAZV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find Dms az data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Dms az data source ID not set")
		}

		return nil
	}
}

var testAccDmsAZV1DataSource_basic = fmt.Sprintf(`

data "g42cloud_availability_zones" "test" {}

data "g42cloud_dms_az" "az1" {
  port = "8002"
  code = data.g42cloud_availability_zones.test.names[0]
}
`)
