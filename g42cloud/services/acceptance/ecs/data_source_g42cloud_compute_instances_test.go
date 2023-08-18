package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccComputeInstancesDataSource_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	dataSourceName := "data.g42cloud_compute_instances.test"
	var instance cloudservers.CloudServer

	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rc := acceptance.InitResourceCheck(
		"g42cloud_compute_instance.test",
		&instance,
		getEcsInstanceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstancesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "instances.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.image_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.flavor_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.flavor_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.network.#"),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr("data.g42cloud_compute_instances.byID", "instances.#", "1"),
				),
			},
		},
	})
}

func testAccComputeInstancesDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.g42cloud_images_image.test.id
  flavor_id          = data.g42cloud_compute_flavors.test.ids[0]
  security_group_ids = [data.g42cloud_networking_secgroup.test.id]
  availability_zone  = data.g42cloud_availability_zones.test.names[0]

  network {
    uuid = data.g42cloud_vpc_subnet.test.id
  }

  tags = {
    foo = "bar"
  }
}

data "g42cloud_compute_instances" "test" {
  name = g42cloud_compute_instance.test.name
}

data "g42cloud_compute_instances" "byID" {
  instance_id = g42cloud_compute_instance.test.id
}
`, testAccCompute_data, rName)
}
