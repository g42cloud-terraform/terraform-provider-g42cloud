package elb

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/elb/v3/listeners"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func getELBListenerResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.ElbV3Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("Error creating G42Cloud elb client: %s", err)
	}
	return listeners.Get(client, state.Primary.ID).Extract()
}

func TestAccElbV3Listener_basic(t *testing.T) {
	var listener listeners.Listener
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "g42cloud_elb_listener.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&listener,
		getELBListenerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3ListenerConfig_basic(testAccElbV3ListenerConfig_base(rName), rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "forward_eip", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccElbV3ListenerConfig_update(testAccElbV3ListenerConfig_base(rName), rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "forward_eip", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccElbV3ListenerConfig_base(rName string) string {
	return fmt.Sprintf(`
data "g42cloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "g42cloud_availability_zones" "test" {}

resource "g42cloud_elb_loadbalancer" "test" {
  name            = "%s"
  ipv4_subnet_id  = data.g42cloud_vpc_subnet.test.subnet_id
  ipv6_network_id = data.g42cloud_vpc_subnet.test.id

  availability_zone = [
    data.g42cloud_availability_zones.test.names[0]
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, rName)
}

func testAccElbV3ListenerConfig_basic(baseConfig, rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_elb_listener" "test" {
  name            = "%s"
  description     = "test description"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = g42cloud_elb_loadbalancer.test.id

  forward_eip = true

  idle_timeout = 62
  request_timeout = 63
  response_timeout = 64

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, baseConfig, rName)
}

func testAccElbV3ListenerConfig_update(baseConfig, rNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_elb_listener" "test" {
  name            = "%s"
  description     = "test description"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = g42cloud_elb_loadbalancer.test.id

  idle_timeout = 62
  request_timeout = 63
  response_timeout = 64

  tags = {
    key1  = "value1"
    owner = "terraform_update"
  }
}
`, baseConfig, rNameUpdate)
}
