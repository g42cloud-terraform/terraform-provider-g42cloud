package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getEcsInstanceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.ComputeV1Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating compute v1 client: %s", err)
	}

	resourceID := state.Primary.ID
	found, err := cloudservers.Get(client, resourceID).Extract()
	if err == nil && found.Status == "DELETED" {
		return nil, fmt.Errorf("the resource %s has been deleted", resourceID)
	}

	return found, err
}

func TestAccComputeInstanceDataSource_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	dataSourceName := "data.g42cloud_compute_instance.this"
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
				Config: testAccComputeInstanceDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "system_disk_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "security_groups.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "volume_attached.#"),
					resource.TestCheckResourceAttrSet("data.g42cloud_compute_instance.byID", "status"),
				),
			},
		},
	})
}

func testAccComputeInstanceDataSource_basic(rName string) string {
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
}

data "g42cloud_compute_instance" "this" {
  name = g42cloud_compute_instance.test.name
}

data "g42cloud_compute_instance" "byID" {
  instance_id = g42cloud_compute_instance.test.id
}
`, testAccCompute_data, rName)
}
