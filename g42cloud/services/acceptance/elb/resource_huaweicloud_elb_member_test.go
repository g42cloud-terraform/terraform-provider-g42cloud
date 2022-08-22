package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/pools"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccElbV3Member_basic(t *testing.T) {
	var member_1 pools.Member
	var member_2 pools.Member
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckElbV3MemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3MemberConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElbV3MemberExists("g42cloud_elb_member.member_1", &member_1),
					testAccCheckElbV3MemberExists("g42cloud_elb_member.member_2", &member_2),
				),
			},
			{
				Config: testAccElbV3MemberConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("g42cloud_elb_member.member_1", "weight", "10"),
					resource.TestCheckResourceAttr("g42cloud_elb_member.member_2", "weight", "15"),
				),
			},
			{
				ResourceName:      "g42cloud_elb_member.member_1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccELBMemberImportStateIdFunc(),
			},
		},
	})
}

func testAccCheckElbV3MemberDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	elbClient, err := config.ElbV3Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating G42Cloud elb client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "g42cloud_elb_member" {
			continue
		}

		poolId := rs.Primary.Attributes["pool_id"]
		_, err := pools.GetMember(elbClient, poolId, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Member still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckElbV3MemberExists(n string, member *pools.Member) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		elbClient, err := config.ElbV3Client(acceptance.G42_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating G42Cloud elb client: %s", err)
		}

		poolId := rs.Primary.Attributes["pool_id"]
		found, err := pools.GetMember(elbClient, poolId, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Member not found")
		}

		*member = *found

		return nil
	}
}

func testAccELBMemberImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		pool, ok := s.RootModule().Resources["g42cloud_elb_pool.test"]
		if !ok {
			return "", fmt.Errorf("pool not found: %s", pool)
		}
		member, ok := s.RootModule().Resources["g42cloud_elb_member.member_1"]
		if !ok {
			return "", fmt.Errorf("member not found: %s", member)
		}
		if pool.Primary.ID == "" || member.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", pool.Primary.ID, member.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", pool.Primary.ID, member.Primary.ID), nil
	}
}

func testAccElbV3MemberConfig_basic(rName string) string {
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
}

resource "g42cloud_elb_listener" "test" {
  name            = "%s"
  description     = "test description"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = g42cloud_elb_loadbalancer.test.id

  forward_eip = true

  idle_timeout = 60
  request_timeout = 60
  response_timeout = 60
}

resource "g42cloud_elb_pool" "test" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = g42cloud_elb_listener.test.id
}

resource "g42cloud_elb_member" "member_1" {
  address       = "192.168.0.10"
  protocol_port = 8080
  pool_id       = g42cloud_elb_pool.test.id
  subnet_id     = data.g42cloud_vpc_subnet.test.subnet_id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "g42cloud_elb_member" "member_2" {
  address       = "192.168.0.11"
  protocol_port = 8080
  pool_id       = g42cloud_elb_pool.test.id
  subnet_id     = data.g42cloud_vpc_subnet.test.subnet_id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, rName, rName, rName)
}

func testAccElbV3MemberConfig_update(rName string) string {
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
}

resource "g42cloud_elb_listener" "test" {
  name            = "%s"
  description     = "test description"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = g42cloud_elb_loadbalancer.test.id

  forward_eip = true

  idle_timeout = 60
  request_timeout = 60
  response_timeout = 60
}

resource "g42cloud_elb_pool" "test" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = g42cloud_elb_listener.test.id
}

resource "g42cloud_elb_member" "member_1" {
  address        = "192.168.0.10"
  protocol_port  = 8080
  weight         = 10
  pool_id        = g42cloud_elb_pool.test.id
  subnet_id      = data.g42cloud_vpc_subnet.test.subnet_id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "g42cloud_elb_member" "member_2" {
  address        = "192.168.0.11"
  protocol_port  = 8080
  weight         = 15
  pool_id        = g42cloud_elb_pool.test.id
  subnet_id      = data.g42cloud_vpc_subnet.test.subnet_id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, rName, rName, rName)
}
