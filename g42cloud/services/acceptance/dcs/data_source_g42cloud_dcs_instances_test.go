package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccDatasourceDcsInstance_basic(t *testing.T) {
	rName := "data.g42cloud_dcs_instances.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDcsInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instances.0.name", name),
					resource.TestCheckResourceAttr(rName, "instances.0.port", "6388"),
					resource.TestCheckResourceAttr(rName, "instances.0.flavor", "redis.ha.xu1.tiny.r2.128"),
					resource.TestCheckResourceAttrPair(rName, "instances.0.flavor",
						"data.g42cloud_dcs_flavors.test", "flavors.0.name"),
				),
			},
		},
	})
}

func testAccDatasourceDcsInstance_basic(name string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_dcs_instances" "test" {
  name   = g42cloud_dcs_instance.instance_1.name
  status = "RUNNING"
}
`, testAccDcsV1Instance_basic(name))
}
