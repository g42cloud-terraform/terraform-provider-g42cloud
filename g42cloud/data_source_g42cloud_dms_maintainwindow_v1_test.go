package g42cloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDmsMaintainWindowV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsMaintainWindowV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsMaintainWindowV1DataSourceID("data.g42cloud_dms_maintainwindow.maintainwindow1"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_dms_maintainwindow.maintainwindow1", "seq", "1"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_dms_maintainwindow.maintainwindow1", "begin", "22"),
				),
			},
		},
	})
}

func testAccCheckDmsMaintainWindowV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find Dms maintainwindow data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Dms maintainwindow data source ID not set")
		}

		return nil
	}
}

var testAccDmsMaintainWindowV1DataSource_basic = fmt.Sprintf(`
data "g42cloud_dms_maintainwindow" "maintainwindow1" {
  seq = 1
}
`)
