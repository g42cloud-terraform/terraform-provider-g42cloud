package g42cloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/dds/v3/instances"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccDDSV3Instance_basic(t *testing.T) {
	var instance instances.Instance
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "g42cloud_dds_instance.instance"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDDSV3InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceV3Config_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSV3InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "ssl", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "08:00-09:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "8"),
				),
			},
			{
				Config: testAccDDSInstanceV3Config_updateBackupStrategy(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSV3InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.start_time", "00:00-01:00"),
					resource.TestCheckResourceAttr(resourceName, "backup_strategy.0.keep_days", "7"),
				),
			},
		},
	})
}

func TestAccDDSV3Instance_withEpsId(t *testing.T) {
	var instance instances.Instance
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "g42cloud_dds_instance.instance"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDDSV3InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSInstanceV3Config_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDDSV3InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", G42_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckDDSV3InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	client, err := config.DdsV3Client(G42_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating G42Cloud DDS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "g42cloud_dds_instance" {
			continue
		}

		opts := instances.ListInstanceOpts{
			Id: rs.Primary.ID,
		}
		allPages, err := instances.List(client, &opts).AllPages()
		if err != nil {
			return err
		}
		instances, err := instances.ExtractInstances(allPages)
		if err != nil {
			return err
		}

		if instances.TotalCount > 0 {
			return fmt.Errorf("Instance still exists. ")
		}
	}

	return nil
}

func testAccCheckDDSV3InstanceExists(n string, instance *instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s. ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set. ")
		}

		config := testAccProvider.Meta().(*config.Config)
		client, err := config.DdsV3Client(G42_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating G42Cloud DDS client: %s ", err)
		}

		opts := instances.ListInstanceOpts{
			Id: rs.Primary.ID,
		}
		allPages, err := instances.List(client, &opts).AllPages()
		if err != nil {
			return err
		}
		instances, err := instances.ExtractInstances(allPages)
		if err != nil {
			return err
		}
		if instances.TotalCount == 0 {
			return fmt.Errorf("dds instance not found.")
		}

		return nil
	}
}

func testAccDDSInstanceV3Config_basic(rName string) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

data "g42cloud_vpc" "test" {
  name = "vpc-default"
}

data "g42cloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "g42cloud_networking_secgroup" "secgroup_acc" {
  name = "%s"
}

resource "g42cloud_dds_instance" "instance" {
  name              = "%s"
  availability_zone = data.g42cloud_availability_zones.test.names[0]
  vpc_id            = data.g42cloud_vpc.test.id
  subnet_id         = data.g42cloud_vpc_subnet.test.id
  security_group_id = g42cloud_networking_secgroup.secgroup_acc.id
  password          = "Test@123"
  mode              = "Sharding"

  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.c6.large.2.mongos"
  }
  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 10
    spec_code = "dds.mongodb.c6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c6.large.2.config"
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = "8"
  }

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}`, rName, rName)
}

func testAccDDSInstanceV3Config_updateBackupStrategy(rName string) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

data "g42cloud_vpc" "test" {
  name = "vpc-default"
}

data "g42cloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "g42cloud_networking_secgroup" "secgroup_acc" {
  name = "%s"
}

resource "g42cloud_dds_instance" "instance" {
  name              = "%s"
  availability_zone = data.g42cloud_availability_zones.test.names[0]
  vpc_id            = data.g42cloud_vpc.test.id
  subnet_id         = data.g42cloud_vpc_subnet.test.id
  security_group_id = g42cloud_networking_secgroup.secgroup_acc.id
  password          = "Test@123"
  mode              = "Sharding"

  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.c6.large.2.mongos"
  }
  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 10
    spec_code = "dds.mongodb.c6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c6.large.2.config"
  }

  backup_strategy {
    start_time = "00:00-01:00"
    keep_days  = "7"
  }

  tags = {
	foo   = "bar"
    owner = "terraform"
  }
}`, rName, rName)
}

func testAccDDSInstanceV3Config_withEpsId(rName string) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

data "g42cloud_vpc" "test" {
  name = "vpc-default"
}

data "g42cloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "g42cloud_networking_secgroup" "secgroup_acc" {
  name = "%s"
}

resource "g42cloud_dds_instance" "instance" {
  name                  = "%s"
  availability_zone     = data.g42cloud_availability_zones.test.names[0]
  vpc_id                = data.g42cloud_vpc.test.id
  subnet_id             = data.g42cloud_vpc_subnet.test.id
  security_group_id     = g42cloud_networking_secgroup.secgroup_acc.id
  password              = "Test@123"
  mode                  = "Sharding"
  enterprise_project_id = "%s"

  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.c6.large.2.mongos"
  }
  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 10
    spec_code = "dds.mongodb.c6.large.2.shard"
  }
  flavor {
    type      = "config"
    num       = 1
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c6.large.2.config"
  }

  backup_strategy {
    start_time = "00:00-01:00"
    keep_days  = "7"
  }

  tags = {
	foo   = "bar"
    owner = "terraform"
  }
}`, rName, rName, G42_ENTERPRISE_PROJECT_ID_TEST)
}
