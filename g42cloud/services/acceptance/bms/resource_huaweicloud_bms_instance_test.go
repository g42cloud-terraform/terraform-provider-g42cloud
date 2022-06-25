package bms

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/bms/v1/baremetalservers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccBmsInstance_basic(t *testing.T) {
	var instance baremetalservers.CloudServer

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "g42cloud_bms_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckBms(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckBmsInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBmsInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBmsInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.G42_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckBmsInstanceDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	bmsClient, err := config.BmsV1Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating G42Cloud bms client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "g42cloud_bms_instance" {
			continue
		}

		server, err := baremetalservers.Get(bmsClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "DELETED" {
				return fmt.Errorf("Instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckBmsInstanceExists(n string, instance *baremetalservers.CloudServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		bmsClient, err := config.BmsV1Client(acceptance.G42_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating G42Cloud bms client: %s", err)
		}

		found, err := baremetalservers.Get(bmsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Instance not found")
		}

		*instance = *found

		return nil
	}
}

func testAccBmsInstance_basic(rName string) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

data "g42cloud_vpc" "test" {
  name = "vpc-default"
}

data "g42cloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "g42cloud_networking_secgroup" "test" {
  name = "default"
}

data "g42cloud_bms_flavors" "test" {
  availability_zone = data.g42cloud_availability_zones.test.names[0]
}

resource "g42cloud_compute_keypair" "test" {
  name = "%s"

  lifecycle {
    ignore_changes = [
      public_key,
    ]
  }
}

resource "g42cloud_vpc_eip" "myeip" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%s"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "g42cloud_bms_instance" "test" {
  name                  = "%s"
  user_id               = "%s"
  # CentOS 7.4 64bit for BareMetal
  image_id              = "519ea918-1fea-4ebc-911a-593739b1a3bc"
  flavor_id             = data.g42cloud_bms_flavors.test.flavors[0].id
  security_groups       = [data.g42cloud_networking_secgroup.test.id]
  availability_zone     = data.g42cloud_availability_zones.test.names[0]
  vpc_id                = data.g42cloud_vpc.test.id
  eip_id                = g42cloud_vpc_eip.myeip.id
  charging_mode         = "prePaid"
  period_unit           = "month"
  period                = "1"
  key_pair              = g42cloud_compute_keypair.test.name
  enterprise_project_id = "%s"

  nics {
    subnet_id = data.g42cloud_vpc_subnet.test.id
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName, rName, rName, acceptance.G42_USER_ID, acceptance.G42_ENTERPRISE_PROJECT_ID_TEST)
}
