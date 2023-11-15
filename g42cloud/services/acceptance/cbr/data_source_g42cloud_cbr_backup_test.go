package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance/common"
)

func TestAccDataBackup_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceNameWithDash()
	dataSourceName := "data.g42cloud_cbr_backup.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataBackup_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccDataBackup_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "g42cloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.g42cloud_images_image.test.id
  flavor_id          = "s6.large.2"
  security_group_ids = [g42cloud_networking_secgroup.test.id]
  availability_zone  = data.g42cloud_availability_zones.test.names[0]

  network {
    uuid = g42cloud_vpc_subnet.test.id
  }

  data_disks {
    type = "SAS"
    size = "10"
  }
}

resource "g42cloud_cbr_vault" "test" {
  name             = "%[2]s"
  type             = "server"
  consistent_level = "app_consistent"
  protection_type  = "backup"
  size             = 200
}

resource "g42cloud_images_image" "test" {
  name        = "%[2]s"
  instance_id = g42cloud_compute_instance.test.id
  vault_id    = g42cloud_cbr_vault.test.id
}

data "g42cloud_cbr_backup" "test" {
  id = g42cloud_images_image.test.backup_id
}
`, common.TestBaseComputeResources(name), name)
}
