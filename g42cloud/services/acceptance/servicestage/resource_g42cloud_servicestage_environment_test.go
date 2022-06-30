package servicestage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/servicestage/v2/environments"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getEnvResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.ServiceStageV2Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ServiceStage v2 client: %s", err)
	}
	return environments.Get(c, state.Primary.ID)
}

func TestAccEnvironment_basic(t *testing.T) {
	var (
		env          environments.Environment
		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "g42cloud_servicestage_environment.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&env,
		getEnvResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironment_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform test"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "g42cloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "basic_resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "optional_resources.#", "1"),
				),
			},
			{
				Config: testAccEnvironment_update(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by terraform test"),
					resource.TestCheckResourceAttr(resourceName, "basic_resources.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "optional_resources.#", "2"),
				),
			},
			{
				Config: testAccEnvironment_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "basic_resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "optional_resources.#", "1"),
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

func testAccEnvironment_base(rName string) string {
	return fmt.Sprintf(`
variable "subnet_config" {
  type = list(object({
    cidr       = string
    gateway_ip = string
  }))

  default = [
    {cidr = "192.168.192.0/18", gateway_ip = "192.168.192.1"},
    {cidr = "192.168.128.0/18", gateway_ip = "192.168.128.1"},
  ]
}

data "g42cloud_availability_zones" "test" {}

data "g42cloud_compute_flavors" "test" {
  availability_zone = data.g42cloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 8
  memory_size       = 16
}

data "g42cloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "g42cloud_compute_keypair" "test" {
  name = "%[1]s"
}

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

resource "g42cloud_networking_secgroup" "test" {
  name = "%[1]s"
}

%s

`, rName, testAccEnvironment_baseRes(rName))
}

func testAccEnvironment_baseRes(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_vpc_eip" "test" {
  count = 2

  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    size        = 5
    name        = "%[1]s_${count.index}"
    charge_mode = "traffic"
  }
}

resource "g42cloud_as_configuration" "test" {
  scaling_configuration_name = "%[1]s"

  instance_config {
    flavor   = data.g42cloud_compute_flavors.test.ids[0]
    image    = data.g42cloud_images_image.test.id
    key_name = g42cloud_compute_keypair.test.name

    disk {
      disk_type   = "SYS"
      volume_type = "SSD"
      size        = 40
    }
  }
}

resource "g42cloud_as_group" "test" {
  count = 2

  scaling_group_name       = "%[1]s"
  scaling_configuration_id = g42cloud_as_configuration.test.id
  vpc_id                   = g42cloud_vpc.test.id

  max_instance_number    = 3
  min_instance_number    = 0
  desire_instance_number = 2

  delete_instances = "yes"
  delete_publicip  = false

  cool_down_time = 86400

  networks {
    id = g42cloud_vpc_subnet.test.id
  }

  security_groups {
    id = g42cloud_networking_secgroup.test.id
  }
}
`, rName)
}

func testAccEnvironment_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_servicestage_environment" "test" {
  name        = "%s"
  description = "Created by terraform test"
  vpc_id      = g42cloud_vpc.test.id

  basic_resources {
    type = "as"
    id   = g42cloud_as_group.test[0].id
  }

  optional_resources {
    type = "eip"
    id   = g42cloud_vpc_eip.test[0].id
  }
}
`, testAccEnvironment_base(rName), rName)
}

func testAccEnvironment_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_servicestage_environment" "test" {
  name        = "%s-update"
  description = "Updated by terraform test"
  vpc_id      = g42cloud_vpc.test.id

  dynamic "basic_resources" {
    for_each = g42cloud_as_group.test[*].id
    content {
      type = "as"
      id   = basic_resources.value
    }
  }

  dynamic "optional_resources" {
    for_each = g42cloud_vpc_eip.test[*].id
    content {
      type = "eip"
      id   = optional_resources.value
    }
  }
}
`, testAccEnvironment_base(rName), rName)
}
