package vpc

import (
	"fmt"
	"log"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccNetworkingV2VIPAssociate_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkingV2VIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2VIPAssociateConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("g42cloud_networking_vip_associate.vip_associate_1",
						"port_ids.0", "g42cloud_compute_instance.test", "network.0.port"),
					resource.TestCheckResourceAttrPair("g42cloud_networking_vip_associate.vip_associate_1",
						"vip_id", "g42cloud_networking_vip.vip_1", "id"),
				),
			},
			{
				ResourceName:      "g42cloud_networking_vip_associate.vip_associate_1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNetworkingV2VIPAssociateImportStateIdFunc(),
			},
		},
	})
}

func testAccCheckNetworkingV2VIPAssociateDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV2Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "g42cloud_networking_vip_associate" {
			continue
		}

		vipID := rs.Primary.Attributes["vip_id"]
		_, err = ports.Get(networkingClient, vipID).Extract()
		if err != nil {
			// If the error is a 404, then the vip port does not exist,
			// and therefore the floating IP cannot be associated to it.
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}
	}

	log.Printf("[DEBUG] Destroy NetworkingVIPAssociated success!")
	return nil
}

func testAccNetworkingV2VIPAssociateImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		vip, ok := s.RootModule().Resources["g42cloud_networking_vip.vip_1"]
		if !ok {
			return "", fmt.Errorf("vip not found: %s", vip)
		}
		instance, ok := s.RootModule().Resources["g42cloud_compute_instance.test"]
		if !ok {
			return "", fmt.Errorf("port not found: %s", instance)
		}
		if vip.Primary.ID == "" || instance.Primary.Attributes["network.0.port"] == "" {
			return "", fmt.Errorf("resource not found: %s/%s", vip.Primary.ID,
				instance.Primary.Attributes["network.0.port"])
		}
		return fmt.Sprintf("%s/%s", vip.Primary.ID, instance.Primary.Attributes["network.0.port"]), nil
	}
}

const testAccCompute_data = `
data "g42cloud_availability_zones" "test" {}

data "g42cloud_compute_flavors" "test" {
  availability_zone = data.g42cloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "g42cloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "g42cloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "g42cloud_networking_secgroup" "test" {
  name = "default"
}
`

func testAccComputeInstance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_compute_instance" "test" {
  name                = "%s"
  image_id            = data.g42cloud_images_image.test.id
  flavor_id           = data.g42cloud_compute_flavors.test.ids[0]
  security_group_ids  = [data.g42cloud_networking_secgroup.test.id]
  stop_before_destroy = true

  network {
    uuid              = data.g42cloud_vpc_subnet.test.id
    source_dest_check = false
  }
}
`, testAccCompute_data, rName)
}

func testAccNetworkingV2VIPAssociateConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_networking_port" "port" {
  port_id = g42cloud_compute_instance.test.network[0].port
}

resource "g42cloud_networking_vip" "vip_1" {
  network_id = data.g42cloud_vpc_subnet.test.id
}

resource "g42cloud_networking_vip_associate" "vip_associate_1" {
  vip_id   = g42cloud_networking_vip.vip_1.id
  port_ids = [g42cloud_compute_instance.test.network[0].port]
}
`, testAccComputeInstance_basic(rName))
}
