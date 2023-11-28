package css

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/css/v1/cluster"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccCssCluster_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "g42cloud_css_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCssClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCssCluster_basic(rName, 1, 7, "bar"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "expect_node_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "elasticsearch"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccCssCluster_basic(rName, 2, 8, "bar_update"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "expect_node_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "8"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
		},
	})
}

func TestAccCssCluster_security(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "g42cloud_css_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCssClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCssCluster_security(rName, 1, "bar"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "expect_node_num", "1"),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "elasticsearch"),
					resource.TestCheckResourceAttr(resourceName, "security_mode", "true"),
				),
			},
		},
	})
}

func testAccVpc(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "g42cloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/24"
  vpc_id     = g42cloud_vpc.test.id
  gateway_ip = "192.168.0.1"
}

resource "g42cloud_networking_secgroup" "test" {
  name        = "%s"
  description = "terraform security group acceptance test"
}

data "g42cloud_availability_zones" "test" {}
`, rName, rName, rName)
}

func testAccCssCluster_basic(rName string, nodeNum int, keepDays int, tag string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_css_cluster" "test" {
  name            = "%s"
  engine_version  = "7.9.3"
  expect_node_num = %d

  node_config {
    flavor            = "ess.spec-8u16g"
    availability_zone = data.g42cloud_availability_zones.test.names[0]

    network_info {
      security_group_id = g42cloud_networking_secgroup.test.id
      subnet_id         = g42cloud_vpc_subnet.test.id
      vpc_id            = g42cloud_vpc.test.id
    }

    volume {
      volume_type = "HIGH"
      size        = 80
    }
  }

  backup_strategy {
    keep_days  = %d
    start_time = "00:00 GMT+08:00"
    prefix     = "snapshot"
  }

  tags = {
    foo = "%s"
    key = "value"
  }
}

`, testAccVpc(rName), rName, nodeNum, keepDays, tag)
}

func testAccCssCluster_security(rName string, nodeNum int, tag string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_css_cluster" "test" {
  name            = "%s"
  engine_version  = "7.9.3"
  expect_node_num = %d
  security_mode   = true
  password        = "Test@passw0rd"

  node_config {
    flavor            = "ess.spec-4u8g"
    availability_zone = data.g42cloud_availability_zones.test.names[0]

    network_info {
      security_group_id = g42cloud_networking_secgroup.test.id
      subnet_id         = g42cloud_vpc_subnet.test.id
      vpc_id            = g42cloud_vpc.test.id
    }

    volume {
      volume_type = "HIGH"
      size        = 40
    }  
  }

  tags = {
    foo = "%s"
    key = "value"
  }
}

`, testAccVpc(rName), rName, nodeNum, tag)
}

func testAccCheckCssClusterDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.CssV1Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("error creating CSS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "g42cloud_css_cluster" {
			continue
		}

		_, err := cluster.Get(client, rs.Primary.ID)
		if err == nil {
			return fmtp.Errorf("css cluster still exists, cluster_id:%s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckCssClusterExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.CssV1Client(acceptance.G42_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("error creating CSS client: %s", err)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmtp.Errorf("Error checking g42cloud_css_cluster exist, err=not found this resource")
		}

		_, errQueryDetail := cluster.Get(client, rs.Primary.ID)
		if errQueryDetail != nil {
			if _, ok := errQueryDetail.(golangsdk.ErrDefault404); ok {
				return fmtp.Errorf("g42cloud_css_cluster is not exist")
			}
			return fmtp.Errorf("checking g42cloud_css_cluster exist,err=send request failed:%s", errQueryDetail)
		}
		return nil
	}
}
