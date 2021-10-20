package g42cloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/rds/v3/instances"
)

func TestAccRdsReadReplicaInstance_basic(t *testing.T) {
	var replica instances.RdsInstanceResponse
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceType := "g42cloud_rds_read_replica_instance"
	resourceName := "g42cloud_rds_read_replica_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccReadRdsReplicaInstance_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &replica),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.c6.large.4.rr"),
					resource.TestCheckResourceAttr(resourceName, "type", "Replica"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccReadRdsReplicaInstance_update(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &replica),
					resource.TestCheckResourceAttr(resourceName, "flavor", "rds.pg.c6.xlarge.4.rr"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "ULTRAHIGH"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.size", "50"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"db",
				},
			},
		},
	})
}

func TestAccRdsReadReplicaInstance_withEpsId(t *testing.T) {
	var replica instances.RdsInstanceResponse
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceType := "g42cloud_rds_read_replica_instance"
	resourceName := "g42cloud_rds_read_replica_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRdsInstanceV3Destroy(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testAccReadRdsReplicaInstance_withEpsId(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsInstanceV3Exists(resourceName, &replica),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", G42_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccReadRdsReplicaInstanceV3_base(name string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_rds_instance" "test" {
  name              = "%s"
  flavor            = "rds.pg.c6.large.4"
  availability_zone = [data.g42cloud_availability_zones.test.names[0]]
  security_group_id = g42cloud_networking_secgroup.test.id
  vpc_id            = g42cloud_vpc.test.id
  subnet_id         = g42cloud_vpc_subnet.test.id

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "11"
    port     = 8635
  }
  volume {
    type = "ULTRAHIGH"
    size = 50
  }
}
`, testAccRdsInstanceV3_base(name), name)
}

func testAccReadRdsReplicaInstance_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_rds_read_replica_instance" "test" {
  name                = "%s"
  flavor              = "rds.pg.c6.large.4.rr"
  primary_instance_id = g42cloud_rds_instance.test.id
  availability_zone   = data.g42cloud_availability_zones.test.names[0]

  volume {
    type = "ULTRAHIGH"
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testAccReadRdsReplicaInstanceV3_base(name), name)
}

func testAccReadRdsReplicaInstance_update(name string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_rds_read_replica_instance" "test" {
  name                = "%s"
  flavor              = "rds.pg.c6.xlarge.4.rr"
  primary_instance_id = g42cloud_rds_instance.test.id
  availability_zone   = data.g42cloud_availability_zones.test.names[0]

  volume {
	type = "ULTRAHIGH"
  }

  tags = {
    key1 = "value"
    foo = "bar2"
  }
}
`, testAccReadRdsReplicaInstanceV3_base(name), name)
}

func testAccReadRdsReplicaInstance_withEpsId(name string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_rds_read_replica_instance" "test" {
  name                  = "%s"
  flavor                = "rds.pg.c6.large.4.rr"
  primary_instance_id   = g42cloud_rds_instance.test.id
  availability_zone     = data.g42cloud_availability_zones.test.names[0]
  enterprise_project_id = "%s"

  volume {
    type = "ULTRAHIGH"
  }
}
`, testAccReadRdsReplicaInstanceV3_base(name), name, G42_ENTERPRISE_PROJECT_ID_TEST)
}
