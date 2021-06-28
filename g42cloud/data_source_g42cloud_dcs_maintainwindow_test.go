package g42cloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDcsMaintainWindowV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsMaintainWindowV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsMaintainWindowV1DataSourceID("data.g42cloud_dcs_maintainwindow.maintainwindow1"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_dcs_maintainwindow.maintainwindow1", "seq", "1"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_dcs_maintainwindow.maintainwindow1", "begin", "22"),
				),
			},
		},
	})
}

func testAccCheckDcsMaintainWindowV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Dcs maintainwindow data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Dcs maintainwindow data source ID not set")
		}

		return nil
	}
}

var testAccDcsMaintainWindowV1DataSource_basic = fmt.Sprintf(`
data "g42cloud_dcs_maintainwindow" "maintainwindow1" {
seq = 1
}
`)
