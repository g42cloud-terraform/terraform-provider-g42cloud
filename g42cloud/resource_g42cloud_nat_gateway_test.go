package g42cloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/natgateways"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"
)

func TestAccNatGateway_basic(t *testing.T) {
	randSuffix := acctest.RandString(5)
	resourceName := "g42cloud_nat_gateway.nat_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatV2GatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatV2Gateway_basic(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatV2GatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("nat-gateway-basic-%s", randSuffix)),
					resource.TestCheckResourceAttr(resourceName, "description", "test for terraform"),
					resource.TestCheckResourceAttr(resourceName, "spec", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNatV2Gateway_update(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("nat-gateway-updated-%s", randSuffix)),
					resource.TestCheckResourceAttr(resourceName, "description", "test for terraform updated"),
					resource.TestCheckResourceAttr(resourceName, "spec", "2"),
				),
			},
		},
	})
}

func testAccCheckNatV2GatewayDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*huaweicloud.Config)
	natClient, err := config.NatGatewayClient(G42_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating G42Cloud nat client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "g42cloud_nat_gateway" {
			continue
		}

		_, err := natgateways.Get(natClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Nat gateway still exists")
		}
	}

	return nil
}

func testAccCheckNatV2GatewayExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*huaweicloud.Config)
		natClient, err := config.NatGatewayClient(G42_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating G42Cloud nat client: %s", err)
		}

		found, err := natgateways.Get(natClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Nat gateway not found")
		}

		return nil
	}
}

func testAccNatPreCondition(suffix string) string {
	return fmt.Sprintf(`
resource "g42cloud_vpc" "vpc_1" {
  name = "nat-vpc-%s"
  cidr = "172.16.0.0/16"
}

resource "g42cloud_vpc_subnet" "subnet_1" {
  name       = "nat-sunnet-%s"
  cidr       = "172.16.10.0/24"
  gateway_ip = "172.16.10.1"
  vpc_id     = g42cloud_vpc.vpc_1.id
  dns_list   = ["100.125.1.250"]
}
	`, suffix, suffix)
}

func testAccNatV2Gateway_basic(suffix string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_nat_gateway" "nat_1" {
  name                = "nat-gateway-basic-%s"
  description         = "test for terraform"
  spec                = "1"
  internal_network_id = g42cloud_vpc_subnet.subnet_1.id
  router_id           = g42cloud_vpc.vpc_1.id
}
	`, testAccNatPreCondition(suffix), suffix)
}

func testAccNatV2Gateway_update(suffix string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_nat_gateway" "nat_1" {
  name                = "nat-gateway-updated-%s"
  description         = "test for terraform updated"
  spec                = "2"
  internal_network_id = g42cloud_vpc_subnet.subnet_1.id
  router_id           = g42cloud_vpc.vpc_1.id

}
	`, testAccNatPreCondition(suffix), suffix)
}
