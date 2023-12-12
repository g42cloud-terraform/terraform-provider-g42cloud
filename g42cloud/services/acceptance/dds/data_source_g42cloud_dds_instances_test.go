package dds

import (
	"fmt"
	"testing"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccDatasourceDdsInstance_basic(t *testing.T) {
	rName := "data.g42cloud_dds_instances.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdsInstance_basic(name, 8800),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instances.0.name", name),
					resource.TestCheckResourceAttr(rName, "instances.0.mode", "Sharding"),
				),
			},
		},
	})
}

func testAccDatasourceDdsInstance_base(rName string, port int) string {
	return fmt.Sprintf(`
%s

data "g42cloud_availability_zones" "test" {}

resource "g42cloud_dds_instance" "test" {
  name              = "%s"
  availability_zone = data.g42cloud_availability_zones.test.names[0]
  vpc_id            = g42cloud_vpc.test.id
  subnet_id         = g42cloud_vpc_subnet.test.id
  security_group_id = g42cloud_networking_secgroup.test.id
  password          = "Terraform@123"
  mode              = "Sharding"
  port              = %d

  datastore {
    type           = "DDS-Community"
    version        = "3.4"
    storage_engine = "wiredTiger"
  }

  flavor {
    type      = "mongos"
    num       = 2
    spec_code = "dds.mongodb.c6.large.4.mongos"
  }

  flavor {
    type      = "shard"
    num       = 2
    storage   = "ULTRAHIGH"
    size      = 20
    spec_code = "dds.mongodb.c6.large.4.shard"
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
}`, common.TestBaseNetwork(rName), rName, port)
}

func testAccDatasourceDdsInstance_basic(name string, port int) string {
	return fmt.Sprintf(`
%s

data "g42cloud_dds_instances" "test" {
  name = g42cloud_dds_instance.test.name
}
`, testAccDatasourceDdsInstance_base(name, port))
}
