package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/logtanks"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getELBLogTankResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.ElbV3Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB client: %s", err)
	}
	return logtanks.Get(client, state.Primary.ID).Extract()
}

func TestAccElbLogTank_basic(t *testing.T) {
	var logTanks logtanks.LogTank
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "g42cloud_elb_logtank.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&logTanks,
		getELBLogTankResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbLogTankConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_id",
						"g42cloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "log_topic_id",
						"g42cloud_lts_stream.test", "id"),
				),
			},
			{
				Config: testAccElbLogTankConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "log_group_id",
						"g42cloud_lts_group.test_update", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "log_topic_id",
						"g42cloud_lts_stream.test_update", "id"),
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

func testAccElbLogTankConfig_base(rName, updateName string) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

resource "g42cloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "g42cloud_vpc_subnet" "test" {
  name        = "%[1]s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  vpc_id      = g42cloud_vpc.test.id
  ipv6_enable = true
}

resource "g42cloud_elb_loadbalancer" "test" {
  name            = "%[1]s"
  ipv4_subnet_id  = g42cloud_vpc_subnet.test.ipv4_subnet_id
  ipv6_network_id = g42cloud_vpc_subnet.test.id

  availability_zone = [
    data.g42cloud_availability_zones.test.names[0]
  ]
}

resource "g42cloud_lts_group" "%[2]s" {
  group_name  = "%[2]s"
  ttl_in_days = 1
}

resource "g42cloud_lts_stream" "%[2]s" {
  group_id    = g42cloud_lts_group.%[2]s.id
  stream_name = "%[2]s"
}
`, rName, updateName)
}

func testAccElbLogTankConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_elb_logtank" "test" {
  loadbalancer_id = g42cloud_elb_loadbalancer.test.id
  log_group_id    = g42cloud_lts_group.test.id
  log_topic_id    = g42cloud_lts_stream.test.id
}
`, testAccElbLogTankConfig_base(rName, "test"))
}

func testAccElbLogTankConfig_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_elb_logtank" "test" {
  loadbalancer_id = g42cloud_elb_loadbalancer.test.id
  log_group_id    = g42cloud_lts_group.test_update.id
  log_topic_id    = g42cloud_lts_stream.test_update.id
}
`, testAccElbLogTankConfig_base(rName, "test_update"))
}
