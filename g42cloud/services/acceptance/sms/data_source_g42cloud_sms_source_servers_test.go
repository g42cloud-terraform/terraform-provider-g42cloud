package sms

import (
	"fmt"
	"testing"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSourceServers_basic(t *testing.T) {
	basicDataSource := "data.g42cloud_sms_source_servers.all"
	byNameDataSource := "data.g42cloud_sms_source_servers.byName"
	nonExistentDataSource := "data.g42cloud_sms_source_servers.non-existent"
	basicDC := acceptance.InitDataSourceCheck(basicDataSource)
	name := acceptance.G42_SMS_SOURCE_SERVER

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSms(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceServers_basic(name),
				Check: resource.ComposeTestCheckFunc(
					basicDC.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(basicDataSource, "servers.#"),
					resource.TestCheckResourceAttrSet(basicDataSource, "servers.0.name"),

					resource.TestCheckResourceAttr(byNameDataSource, "servers.#", "1"),
					resource.TestCheckResourceAttr(byNameDataSource, "servers.0.name", name),
					resource.TestCheckResourceAttrSet(byNameDataSource, "servers.0.ip"),
					resource.TestCheckResourceAttrSet(byNameDataSource, "servers.0.state"),

					resource.TestCheckResourceAttr(nonExistentDataSource, "id", "0"),
					resource.TestCheckResourceAttr(nonExistentDataSource, "servers.#", "0"),
				),
			},
		},
	})
}

func testAccSourceServers_basic(name string) string {
	return fmt.Sprintf(`
data "g42cloud_sms_source_servers" "all" {
}

data "g42cloud_sms_source_servers" "byName" {
  name = "%s"
}

data "g42cloud_sms_source_servers" "non-existent" {
  name = "non-existent"
}
`, name)
}
