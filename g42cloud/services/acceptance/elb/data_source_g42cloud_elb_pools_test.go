package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccDatasourcePools_basic(t *testing.T) {
	rName := "data.g42cloud_elb_pools.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePools_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "pools.0.name", name),
					resource.TestCheckResourceAttrPair(rName, "pools.0.id",
						"g42cloud_elb_pool.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.description",
						"g42cloud_elb_pool.test", "description"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.protocol",
						"g42cloud_elb_pool.test", "protocol"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.lb_method",
						"g42cloud_elb_pool.test", "lb_method"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.type",
						"g42cloud_elb_pool.test", "type"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.protection_status",
						"g42cloud_elb_pool.test", "protection_status"),
					resource.TestCheckResourceAttr(rName, "pools.0.slow_start_enabled", "false"),
				),
			},
		},
	})
}

func testAccDatasourcePools_basic(name string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_elb_pools" "test" {
  name = "%s"

  depends_on = [
    g42cloud_elb_pool.test
  ]
}
`, testAccElbV3PoolConfig_basic(name), name)
}
