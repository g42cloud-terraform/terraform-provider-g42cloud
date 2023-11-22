package lb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccDataLBPools_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.g42cloud_lb_pools.test"

	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataLBPools_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "pools.0.name", name),
					resource.TestCheckResourceAttrPair(rName, "pools.0.id",
						"g42cloud_lb_pool.pool_1", "id"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.description",
						"g42cloud_lb_pool.pool_1", "description"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.protocol",
						"g42cloud_lb_pool.pool_1", "protocol"),
					resource.TestCheckResourceAttrPair(rName, "pools.0.lb_method",
						"g42cloud_lb_pool.pool_1", "lb_method"),
				),
			},
		},
	})
}

func testAccDataLBPools_basic(name string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_lb_pools" "test" {
  name = "%s"

  depends_on = [
    g42cloud_lb_pool.pool_1
  ]
}
`, testAccLBV2PoolConfig_basic(name), name)
}

func testAccLBV2PoolConfig_basic(rName string) string {
	return fmt.Sprintf(`
data "g42cloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "g42cloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "%s"
  vip_subnet_id = data.g42cloud_vpc_subnet.test.ipv4_subnet_id
}

resource "g42cloud_lb_listener" "listener_1" {
  name            = "%s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = g42cloud_lb_loadbalancer.loadbalancer_1.id
}

resource "g42cloud_lb_pool" "pool_1" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = g42cloud_lb_listener.listener_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, rName, rName, rName)
}
