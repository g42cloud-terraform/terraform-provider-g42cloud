package g42cloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/dms/v1/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccDmsInstancesV1_Rabbitmq(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dms_instance_%s", acctest.RandString(5))
	resourceName := "g42cloud_dms_instance.instance_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDmsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1Instance_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
		},
	})
}

func TestAccDmsInstancesV1_Kafka(t *testing.T) {
	var instance instances.Instance
	var instanceName = fmt.Sprintf("dms_instance_%s", acctest.RandString(5))
	resourceName := "g42cloud_dms_instance.instance_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDmsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1Instance_KafkaInstance(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsV1InstanceExists(resourceName, instance),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
		},
	})
}

func testAccCheckDmsV1InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	dmsClient, err := config.DmsV1Client(G42_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating G42Cloud instance client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "g42cloud_dms_instance" {
			continue
		}

		_, err := instances.Get(dmsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("The Dms instance still exists.")
		}
	}
	return nil
}

func testAccCheckDmsV1InstanceExists(n string, instance instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		dmsClient, err := config.DmsV1Client(G42_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating G42Cloud instance client: %s", err)
		}

		v, err := instances.Get(dmsClient, rs.Primary.ID).Extract()
		if err != nil {
			return fmt.Errorf("Error getting G42Cloud instance: %s, err: %s", rs.Primary.ID, err)
		}

		if v.InstanceID != rs.Primary.ID {
			return fmt.Errorf("The Dms instance not found.")
		}
		instance = *v
		return nil
	}
}

func testAccDmsV1Instance_base(name string) string {
	return fmt.Sprintf(`
data "g42cloud__dms_az" "test" {}

data "g42cloud_vpc" "test" {
  name = "vpc-default"
}

data "g42cloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "g42cloud_networking_secgroup" "test" {
  name = "%s"
}
`, name)
}

func testAccDmsV1Instance_basic(instanceName string) string {
	return fmt.Sprintf(`
%s
data "g42cloud_dms_product" "product_1" {
  engine        = "rabbitmq"
  instance_type = "single"
  version       = "3.7.17"
}

resource "g42cloud_dms_instance" "instance_1" {
  name              = "%s"
  engine            = "rabbitmq"
  access_user       = "user"
  password          = "Dmstest@123"
  vpc_id            = data.g42cloud_vpc.test.id
  subnet_id         = data.g42cloud_vpc_subnet.test.id
  security_group_id = g42cloud_networking_secgroup.test.id
  available_zones   = [data.g42cloud_dms_az.test.id]
  product_id        = data.g42cloud_dms_product.product_1.id
  engine_version    = data.g42cloud_dms_product.product_1.version
  storage_space     = data.g42cloud_dms_product.product_1.storage
  storage_spec_code = data.g42cloud_dms_product.product_1.storage_spec_code

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
	`, testAccDmsV1Instance_base(instanceName), instanceName)
}

func testAccDmsV1Instance_KafkaInstance(instanceName string) string {
	return fmt.Sprintf(`
%s
data "g42cloud_dms_product" "product_1" {
  engine        = "kafka"
  instance_type = "cluster"
  version       = "1.1.0"
}

resource "g42cloud_dms_instance" "instance_1" {
  name              = "%s"
  engine            = "kafka"
  vpc_id            = data.g42cloud_vpc.test.id
  subnet_id         = data.g42cloud_vpc_subnet.test.id
  security_group_id = g42cloud_networking_secgroup.test.id
  available_zones   = [data.g42cloud_dms_az.test.id]
  product_id        = data.g42cloud_dms_product.product_1.id
  engine_version    = data.g42cloud_dms_product.product_1.version
  specification     = data.g42cloud_dms_product.product_1.bandwidth
  partition_num     = data.g42cloud_dms_product.product_1.partition_num
  storage_space     = data.g42cloud_dms_product.product_1.storage
  storage_spec_code = data.g42cloud_dms_product.product_1.storage_spec_code

  tags = {
    key   = "value"
    owner = "terraform"
  }
}`, testAccDmsV1Instance_base(instanceName), instanceName)
}
