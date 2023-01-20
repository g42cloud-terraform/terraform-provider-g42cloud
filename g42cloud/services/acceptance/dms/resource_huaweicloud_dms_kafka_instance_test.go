package dms

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dms/v2/kafka/instances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getKafkaInstanceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DmsV2Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating G42cloud DMS client(V2): %s", err)
	}
	return instances.Get(client, state.Primary.ID).Extract()
}

func TestAccKafkaInstance_basic(t *testing.T) {
	var instance instances.Instance
	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := rName + "update"
	resourceName := "g42cloud_dms_kafka_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getKafkaInstanceFunc,
	)

	// DMS instances use the tenant-level shared lock, the instances cannot be created or modified in parallel.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "kafka"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccKafkaInstance_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "kafka test update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"manager_password",
					"used_storage_space",
					"cross_vpc_accesses",
				},
			},
		},
	})
}

func testAccKafkaInstance_base(rName string) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

resource "g42cloud_vpc" "test" {
  name        = "%s"
  cidr        = "192.168.11.0/24"
  description = "test for kafka"
}

resource "g42cloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.11.0/24"
  gateway_ip = "192.168.11.1"
  vpc_id     = g42cloud_vpc.test.id
}

resource "g42cloud_networking_secgroup" "test" {
  name        = "%s"
  description = "secgroup for kafka"
}
`, rName, rName, rName)
}

func testAccKafkaInstance_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_dms_kafka_flavors" "test" {
  type = "cluster"
}

resource "g42cloud_dms_kafka_instance" "test" {
  name               = "%s"
  description        = "kafka test"
  vpc_id             = g42cloud_vpc.test.id
  network_id         = g42cloud_vpc_subnet.test.id
  security_group_id  = g42cloud_networking_secgroup.test.id

  flavor_id          = data.g42cloud_dms_kafka_flavors.test.flavors[0].id
  storage_spec_code  = data.g42cloud_dms_kafka_flavors.test.flavors[0].ios[0].storage_spec_code
  availability_zones = [data.g42cloud_dms_kafka_flavors.test.flavors[0].ios[0].availability_zones[0]]
  engine_version     = data.g42cloud_dms_kafka_flavors.test.versions[0]
  storage_space      = data.g42cloud_dms_kafka_flavors.test.flavors[0].properties[0].min_broker * data.g42cloud_dms_kafka_flavors.test.flavors[0].properties[0].min_storage_per_node
  broker_num         = data.g42cloud_dms_kafka_flavors.test.flavors[0].properties[0].min_broker

  access_user      = "user"
  password         = "Kafkatest@123"
  manager_user     = "kafka-user"
  manager_password = "Kafkatest@123"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccKafkaInstance_base(rName), rName)
}

func testAccKafkaInstance_update(rName, updateName string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_dms_kafka_flavors" "test" {
  type = "cluster"
}

resource "g42cloud_dms_kafka_instance" "test" {
  name               = "%s"
  description        = "kafka test update"
  vpc_id             = g42cloud_vpc.test.id
  network_id         = g42cloud_vpc_subnet.test.id
  security_group_id  = g42cloud_networking_secgroup.test.id
  
  flavor_id          = data.g42cloud_dms_kafka_flavors.test.flavors[0].id
  storage_spec_code  = data.g42cloud_dms_kafka_flavors.test.flavors[0].ios[0].storage_spec_code
  availability_zones = [data.g42cloud_dms_kafka_flavors.test.flavors[0].ios[0].availability_zones[0]]
  engine_version     = data.g42cloud_dms_kafka_flavors.test.versions[0]
  storage_space      = data.g42cloud_dms_kafka_flavors.test.flavors[0].properties[0].min_broker * data.g42cloud_dms_kafka_flavors.test.flavors[0].properties[0].min_storage_per_node
  broker_num         = data.g42cloud_dms_kafka_flavors.test.flavors[0].properties[0].min_broker

  access_user      = "user"
  password         = "Kafkatest@123"
  manager_user     = "kafka-user"
  manager_password = "Kafkatest@123"

  tags = {
    key1  = "value"
    owner = "terraform_update"
  }
}
`, testAccKafkaInstance_base(rName), updateName)
}
