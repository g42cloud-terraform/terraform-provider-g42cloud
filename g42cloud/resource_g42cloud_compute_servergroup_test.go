package g42cloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/compute/v2/extensions/servergroups"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccComputeServerGroup_basic(t *testing.T) {
	var sg servergroups.ServerGroup
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "g42cloud_compute_servergroup.sg_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeServerGroupExists(resourceName, &sg),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
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

func TestAccComputeServerGroup_scheduler(t *testing.T) {
	var instance cloudservers.CloudServer
	var sg servergroups.ServerGroup
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "g42cloud_compute_servergroup.sg_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroup_scheduler(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeServerGroupExists(resourceName, &sg),
					testAccCheckComputeInstanceExists("g42cloud_compute_instance.instance_1", &instance),
					testAccCheckComputeInstanceInServerGroup(&instance, &sg),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func TestAccComputeServerGroup_members(t *testing.T) {
	var instance cloudservers.CloudServer
	var sg servergroups.ServerGroup
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "g42cloud_compute_servergroup.sg_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroup_members(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeServerGroupExists(resourceName, &sg),
					testAccCheckComputeInstanceExists("g42cloud_compute_instance.instance_1", &instance),
					testAccCheckComputeInstanceInServerGroup(&instance, &sg),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func testAccCheckComputeServerGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	computeClient, err := config.ComputeV2Client(G42_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating G42Cloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "g42cloud_compute_servergroup" {
			continue
		}

		_, err := servergroups.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("ServerGroup still exists")
		}
	}

	return nil
}

func testAccCheckComputeServerGroupExists(n string, kp *servergroups.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		computeClient, err := config.ComputeV2Client(G42_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating G42Cloud compute client: %s", err)
		}

		found, err := servergroups.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("ServerGroup not found")
		}

		*kp = *found

		return nil
	}
}

func testAccCheckComputeInstanceInServerGroup(instance *cloudservers.CloudServer, sg *servergroups.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(sg.Members) > 0 {
			for _, m := range sg.Members {
				if m == instance.ID {
					return nil
				}
			}
		}

		return fmt.Errorf("Instance %s is not part of Server Group %s", instance.ID, sg.ID)
	}
}

func testAccComputeServerGroup_basic(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_compute_servergroup" "sg_1" {
  name     = "%s"
  policies = ["anti-affinity"]
}
`, rName)
}

func testAccComputeServerGroup_scheduler(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_compute_servergroup" "sg_1" {
  name     = "%s"
  policies = ["anti-affinity"]
}

resource "g42cloud_compute_instance" "instance_1" {
  name               = "%s"
  image_id           = data.g42cloud_images_image.test.id
  flavor_id          = data.g42cloud_compute_flavors.test.ids[0]
  security_group_ids = [data.g42cloud_networking_secgroup.test.id]
  availability_zone  = data.g42cloud_availability_zones.test.names[0]

  scheduler_hints {
    group = g42cloud_compute_servergroup.sg_1.id
  }
  network {
    uuid = data.g42cloud_vpc_subnet.test.id
  }
}
`, testAccCompute_data, rName, rName)
}

func testAccComputeServerGroup_members(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_compute_servergroup" "sg_1" {
  name     = "%s"
  policies = ["anti-affinity"]
  members  = [g42cloud_compute_instance.instance_1.id]
}

resource "g42cloud_compute_instance" "instance_1" {
  name               = "%s"
  image_id           = data.g42cloud_images_image.test.id
  flavor_id          = data.g42cloud_compute_flavors.test.ids[0]
  security_group_ids = [data.g42cloud_networking_secgroup.test.id]
  availability_zone  = data.g42cloud_availability_zones.test.names[0]

  network {
    uuid = data.g42cloud_vpc_subnet.test.id
  }
}
`, testAccCompute_data, rName, rName)
}
