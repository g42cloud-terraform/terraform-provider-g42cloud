package g42cloud

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dcs/v1/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccDcsInstancesV1_basic(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resourceName := "g42cloud_dcs_instance.instance_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "ip"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
				),
			},
			{
				Config: testAccDcsV1Instance_updated(instanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("g42cloud_dcs_instance.instance_1", "backup_policy.0.begin_at", "01:00-02:00"),
					resource.TestCheckResourceAttr("g42cloud_dcs_instance.instance_1", "backup_policy.0.save_days", "2"),
					resource.TestCheckResourceAttr("g42cloud_dcs_instance.instance_1", "backup_policy.0.backup_at.#", "3"),
				),
			},
		},
	})
}

func TestAccDcsInstancesV1_single(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dcs_instance_%s", acctest.RandString(5))
	resourceName := "g42cloud_dcs_instance.instance_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsV1Instance_single(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "Redis"),
					resource.TestCheckResourceAttr(resourceName, "engine_version", "5.0"),
					resource.TestCheckResourceAttr(resourceName, "capacity", "2"),
				),
			},
		},
	})
}

func testAccCheckDcsV1InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	dcsClient, err := config.DcsV1Client(G42_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating G42Cloud instance client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "g42cloud_dcs_instance" {
			continue
		}

		_, err := instances.Get(dcsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("the DCS instance still exists")
		}
	}
	return nil
}

func testAccCheckDcsV1InstanceExists(n string, instance instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		dcsClient, err := config.DcsV1Client(G42_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating G42Cloud instance client: %s", err)
		}

		v, err := instances.Get(dcsClient, rs.Primary.ID).Extract()
		if err != nil {
			return fmt.Errorf("Error getting G42Cloud instance: %s, err: %s", rs.Primary.ID, err)
		}

		if v.InstanceID != rs.Primary.ID {
			return fmt.Errorf("the DCS instance not found")
		}
		instance = *v
		return nil
	}
}

func testAccDcsV1Instance_basic(instanceName string) string {
	return fmt.Sprintf(`
	data "g42cloud_availability_zones" "test" {}

	data "g42cloud_vpc" "test" {
	  name = "vpc-default"
	}

	data "g42cloud_vpc_subnet" "test" {
	  name = "subnet-default"
	}

	data "g42cloud_dcs_az" "az_1" {
	  code = data.g42cloud_availability_zones.test.names[0]
	}

	resource "g42cloud_dcs_instance" "instance_1" {
	  name              = "%s"
	  engine_version    = "5.0"
	  password          = "G42_test"
	  engine            = "Redis"
	  capacity          = "2"
	  vpc_id            = data.g42cloud_vpc.test.id
	  subnet_id         = data.g42cloud_vpc_subnet.test.id
	  available_zones   = [data.g42cloud_dcs_az.az_1.id]
	  product_id        = "redis.ha.xu1.large.r2.2-h"
      backup_policy {
        backup_type = "manual"
        begin_at    = "00:00-01:00"
        period_type = "weekly"
        backup_at = [4]
        save_days = 1
      }

	  tags = {
	    key = "value"
	    owner = "terraform"
	  }
	}
	`, instanceName)
}

func testAccDcsV1Instance_updated(instanceName string) string {
	return fmt.Sprintf(`
	data "g42cloud_availability_zones" "test" {}

	data "g42cloud_vpc" "test" {
	  name = "vpc-default"
	}

	data "g42cloud_vpc_subnet" "test" {
	  name = "subnet-default"
	}

	data "g42cloud_dcs_az" "az_1" {
	  code = data.g42cloud_availability_zones.test.names[0]
	}

	resource "g42cloud_dcs_instance" "instance_1" {
	  name              = "%s"
	  engine_version    = "5.0"
	  password          = "G42_test"
	  engine            = "Redis"
	  capacity          = "2"
	  vpc_id            = data.g42cloud_vpc.test.id
	  subnet_id         = data.g42cloud_vpc_subnet.test.id
	  available_zones   = [data.g42cloud_dcs_az.az_1.id]
	  product_id        = "redis.ha.xu1.large.r2.2-h"
      backup_policy {
        backup_type = "manual"
        begin_at    = "01:00-02:00"
        period_type = "weekly"
        backup_at = [1, 2, 4]
        save_days = 2
      }

	  tags = {
	    key = "value"
	    owner = "terraform"
	  }
	}
	`, instanceName)
}

func testAccDcsV1Instance_single(instanceName string) string {
	return fmt.Sprintf(`
	data "g42cloud_availability_zones" "test" {}

	data "g42cloud_vpc" "test" {
	  name = "vpc-default"
	}

	data "g42cloud_vpc_subnet" "test" {
	  name = "subnet-default"
	}

	data "g42cloud_dcs_az" "az_1" {
	  code = data.g42cloud_availability_zones.test.names[0]
	}

	resource "g42cloud_dcs_instance" "instance_1" {
	  name              = "%s"
	  engine_version    = "5.0"
	  password          = "G42_test"
	  engine            = "Redis"
	  capacity          = 2
	  vpc_id            = data.g42cloud_vpc.test.id
	  subnet_id         = data.g42cloud_vpc_subnet.test.id
	  available_zones   = [data.g42cloud_dcs_az.az_1.id]
	  product_id        = "redis.single.xu1.large.2-h"
	}
	`, instanceName)
}
