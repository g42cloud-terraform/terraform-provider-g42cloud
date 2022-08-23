package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/ports"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getNetworkVipResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.NetworkingV1Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating G42Cloud VPC network v1 client: %s", err)
	}

	return ports.Get(client, state.Primary.ID)
}

func TestAccNetworkingVip_basic(t *testing.T) {
	var vip ports.Port
	resourceName := "g42cloud_networking_vip.test"
	rName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vip,
		getNetworkVipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2VIPConfig_ipv4(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "4"),
					resource.TestCheckResourceAttrSet(resourceName, "mac_address"),
				),
			},
			{
				Config: testAccNetworkingV2VIPConfig_ipv4(rName + "_update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_update"),
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

func TestAccNetworkingVip_ipv6(t *testing.T) {
	var vip ports.Port
	resourceName := "g42cloud_networking_vip.test"
	rName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vip,
		getNetworkVipResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2VIPConfig_ipv6(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "6"),
					resource.TestCheckResourceAttrSet(resourceName, "mac_address"),
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

func testAccNetworkingV2VIPConfig_ipv4(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "g42cloud_vpc_subnet" "test" {
  vpc_id     = g42cloud_vpc.test.id
  name       = "%[1]s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}

resource "g42cloud_networking_vip" "test" {
  name       = "%[1]s"
  network_id = g42cloud_vpc_subnet.test.id
}
`, rName)
}

func testAccNetworkingV2VIPConfig_ipv6(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "g42cloud_vpc_subnet" "test" {
  vpc_id      = g42cloud_vpc.test.id
  name        = "%[1]s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  ipv6_enable = true
}

resource "g42cloud_networking_vip" "test" {
  name       = "%[1]s"
  network_id = g42cloud_vpc_subnet.test.id
  ip_version = 6
}
`, rName)
}
