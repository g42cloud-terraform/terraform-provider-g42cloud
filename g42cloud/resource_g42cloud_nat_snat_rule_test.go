package g42cloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/hw_snatrules"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"
)

func TestAccNatSnatRule_basic(t *testing.T) {
	randSuffix := acctest.RandString(5)
	resourceName := "g42cloud_nat_snat_rule.snat_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNatV2SnatRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatV2SnatRule_basic(randSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNatV2GatewayExists("g42cloud_nat_gateway.nat_1"),
					testAccCheckNatV2SnatRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
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

func testAccCheckNatV2SnatRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*huaweicloud.Config)
	natClient, err := config.NatGatewayClient(G42_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating G42Cloud nat client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "g42cloud_nat_snat_rule" {
			continue
		}

		_, err := hw_snatrules.Get(natClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Snat rule still exists")
		}
	}

	return nil
}

func testAccCheckNatV2SnatRuleExists(n string) resource.TestCheckFunc {
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

		found, err := hw_snatrules.Get(natClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Snat rule not found")
		}

		return nil
	}
}

func testAccNatV2SnatRule_basic(suffix string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_vpc_eip" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "g42cloud_nat_gateway" "nat_1" {
  name                = "nat-gateway-basic-%s"
  description         = "test for terraform"
  spec                = "1"
  internal_network_id = g42cloud_vpc_subnet.subnet_1.id
  router_id           = g42cloud_vpc.vpc_1.id
}

resource "g42cloud_nat_snat_rule" "snat_1" {
  nat_gateway_id = g42cloud_nat_gateway.nat_1.id
  network_id     = g42cloud_vpc_subnet.subnet_1.id
  floating_ip_id = g42cloud_vpc_eip.eip_1.id
}
	`, testAccNatPreCondition(suffix), suffix)
}
