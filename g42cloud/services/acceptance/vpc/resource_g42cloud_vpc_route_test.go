package vpc

import (
	"fmt"
	"testing"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v2/routes"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccVpcRouteV2_basic(t *testing.T) {
	var route routes.Route

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "g42cloud_vpc_route.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRouteV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteV2_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteV2Exists(resourceName, &route),
					resource.TestCheckResourceAttr(resourceName, "destination", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "type", "peering"),
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

func testAccCheckRouteV2Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	routeClient, err := config.NetworkingV2Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud route client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "g42cloud_vpc_route" {
			continue
		}

		_, err := routes.Get(routeClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Route still exists")
		}
	}

	return nil
}

func testAccCheckRouteV2Exists(n string, route *routes.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		routeClient, err := config.NetworkingV2Client(acceptance.G42_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating huaweicloud route client: %s", err)
		}

		found, err := routes.Get(routeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.RouteID != rs.Primary.ID {
			return fmt.Errorf("route not found")
		}

		*route = *found

		return nil
	}
}

func testAccRouteV2_basic(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "g42cloud_vpc" "test2" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "g42cloud_vpc_peering_connection" "test" {
  name        = "%s"
  vpc_id      = g42cloud_vpc.test.id
  peer_vpc_id = g42cloud_vpc.test2.id
}

resource "g42cloud_vpc_route" "test" {
  type        = "peering"
  nexthop     = g42cloud_vpc_peering_connection.test.id
  destination = "192.168.0.0/16"
  vpc_id      = g42cloud_vpc.test.id

}
`, rName, rName+"2", rName)
}
