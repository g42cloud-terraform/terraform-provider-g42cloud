package g42cloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccComputeV2Instance_basic(t *testing.T) {
	var instance cloudservers.CloudServer

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "g42cloud_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"stop_before_destroy",
					"force_delete",
				},
			},
		},
	})
}

func TestAccComputeV2Instance_tags(t *testing.T) {
	var instance cloudservers.CloudServer

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "g42cloud_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_tags(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					testAccCheckComputeInstanceTags(&instance, "foo", "bar"),
					testAccCheckComputeInstanceTags(&instance, "key", "value"),
				),
			},
			{
				Config: testAccComputeV2Instance_tags2(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					testAccCheckComputeInstanceTags(&instance, "foo2", "bar2"),
					testAccCheckComputeInstanceTags(&instance, "key", "value2"),
				),
			},
			{
				Config: testAccComputeV2Instance_notags(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					testAccCheckComputeInstanceNoTags(&instance),
				),
			},
			{
				Config: testAccComputeV2Instance_tags(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					testAccCheckComputeInstanceTags(&instance, "foo", "bar"),
					testAccCheckComputeInstanceTags(&instance, "key", "value"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_disks(t *testing.T) {
	var instance cloudservers.CloudServer

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "g42cloud_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_disks(rName, 50),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "system_disk_size", "50"),
					resource.TestCheckResourceAttrPair(
						resourceName, "system_disk_kms_key_id", "g42cloud_kms_key.test", "id"),
					resource.TestCheckResourceAttrPair(
						resourceName, "volume_attached.1.kms_key_id", "g42cloud_kms_key.test", "id"),
				),
			},
			{
				Config: testAccComputeV2Instance_disks(rName, 60),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "system_disk_size", "60"),
				),
			},
		},
	})
}

func testAccCheckComputeInstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	computeClient, err := config.ComputeV1Client(G42_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating G42Cloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_compute_instance" {
			continue
		}

		server, err := cloudservers.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "DELETED" {
				return fmt.Errorf("Instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckComputeInstanceExists(n string, instance *cloudservers.CloudServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		computeClient, err := config.ComputeV1Client(G42_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating G42cloud compute client: %s", err)
		}

		found, err := cloudservers.Get(computeClient, rs.Primary.ID).Extract()
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

func testAccCheckComputeInstanceTags(
	instance *cloudservers.CloudServer, k, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config.Config)
		client, err := config.ComputeV1Client(G42_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating G42Cloud compute v1 client: %s", err)
		}

		taglist, err := tags.Get(client, "cloudservers", instance.ID).Extract()
		for _, val := range taglist.Tags {
			if k != val.Key {
				continue
			}

			if v == val.Value {
				return nil
			}

			return fmt.Errorf("Bad value for %s: %s", k, val.Value)
		}

		return fmt.Errorf("Tag not found: %s", k)
	}
}

func testAccCheckComputeInstanceNoTags(
	instance *cloudservers.CloudServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*config.Config)
		client, err := config.ComputeV1Client(G42_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating G42Cloud compute v1 client: %s", err)
		}

		taglist, err := tags.Get(client, "cloudservers", instance.ID).Extract()

		if taglist.Tags == nil {
			return nil
		}
		if len(taglist.Tags) == 0 {
			return nil
		}

		return fmt.Errorf("Expected no tags, but found %v", taglist.Tags)
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

func testAccComputeV2Instance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_compute_instance" "test" {
  name              = "%s"
  image_id          = data.g42cloud_images_image.test.id
  flavor_id         = data.g42cloud_compute_flavors.test.ids[0]
  security_group_ids = [data.g42cloud_networking_secgroup.test.id]
  availability_zone = data.g42cloud_availability_zones.test.names[0]
  system_disk_type  = "SSD"

  network {
    uuid = data.g42cloud_vpc_subnet.test.id
  }
}
`, testAccCompute_data, rName)
}

func testAccComputeV2Instance_tags(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.g42cloud_images_image.test.id
  flavor_id          = data.g42cloud_compute_flavors.test.ids[0]
  security_group_ids = [data.g42cloud_networking_secgroup.test.id]
  availability_zone  = data.g42cloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = data.g42cloud_vpc_subnet.test.id
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCompute_data, rName)
}

func testAccComputeV2Instance_tags2(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.g42cloud_images_image.test.id
  flavor_id          = data.g42cloud_compute_flavors.test.ids[0]
  security_group_ids = [data.g42cloud_networking_secgroup.test.id]
  availability_zone  = data.g42cloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = data.g42cloud_vpc_subnet.test.id
  }

  tags = {
    foo2 = "bar2"
    key = "value2"
  }
}
`, testAccCompute_data, rName)
}

func testAccComputeV2Instance_notags(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.g42cloud_images_image.test.id
  flavor_id          = data.g42cloud_compute_flavors.test.ids[0]
  security_group_ids = [data.g42cloud_networking_secgroup.test.id]
  availability_zone  = data.g42cloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = data.g42cloud_vpc_subnet.test.id
  }
}
`, testAccCompute_data, rName)
}

func testAccComputeV2Instance_disks(rName string, systemDiskSize int) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_kms_key" "test" {
  key_alias    = "%s"
  pending_days = "7"
}

resource "g42cloud_compute_instance" "test" {
  name                        = "%s"
  image_id                    = data.g42cloud_images_image.test.id
  flavor_id                   = data.g42cloud_compute_flavors.test.ids[0]
  security_group_ids          = [data.g42cloud_networking_secgroup.test.id]
  availability_zone           = data.g42cloud_availability_zones.test.names[0]
  delete_disks_on_termination = true

  system_disk_type       = "SAS"
  system_disk_size       = %d
  system_disk_kms_key_id = g42cloud_kms_key.test.id

  data_disks {
    type       = "SAS"
    size       = "10"
	kms_key_id = g42cloud_kms_key.test.id
  }

  network {
    uuid = data.g42cloud_vpc_subnet.test.id
  }
}
`, testAccCompute_data, rName, rName, systemDiskSize)
}
