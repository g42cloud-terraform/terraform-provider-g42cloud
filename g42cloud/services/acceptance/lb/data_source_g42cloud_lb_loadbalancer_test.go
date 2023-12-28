package lb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/chnsz/golangsdk/openstack/elb/v2/loadbalancers"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func getLoadBalancerResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.LoadBalancerClient(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB v2 Client: %s", err)
	}
	resp, err := loadbalancers.Get(c, state.Primary.ID).Extract()
	if resp == nil && err == nil {
		return resp, fmt.Errorf("unable to find the LoadBalancer (%s)", state.Primary.ID)
	}
	return resp, err
}

func TestAccELBV2LoadbalancerDataSource_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	dataSourceName1 := "data.g42cloud_lb_loadbalancer.test_by_name"
	dc1 := acceptance.InitDataSourceCheck(dataSourceName1)
	dataSourceName2 := "data.g42cloud_lb_loadbalancer.test_by_description"
	dc2 := acceptance.InitDataSourceCheck(dataSourceName2)

	var lb loadbalancers.LoadBalancer
	resourceName := "g42cloud_lb_loadbalancer.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getLoadBalancerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccELBV2LoadbalancerDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName1, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName2, "name", rName),
				),
			},
		},
	})
}

func testAccELBV2LoadbalancerDataSource_basic(rName string) string {
	return fmt.Sprintf(`
data "g42cloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "g42cloud_lb_loadbalancer" "test" {
  name          = "%s"
  vip_subnet_id = data.g42cloud_vpc_subnet.test.ipv4_subnet_id
  description   = "test for load balancer data source"
}

data "g42cloud_lb_loadbalancer" "test_by_name" {
  name = g42cloud_lb_loadbalancer.test.name

  depends_on = [g42cloud_lb_loadbalancer.test]
}

data "g42cloud_lb_loadbalancer" "test_by_description" {
  description = g42cloud_lb_loadbalancer.test.description

  depends_on = [g42cloud_lb_loadbalancer.test]
}
`, rName)
}
