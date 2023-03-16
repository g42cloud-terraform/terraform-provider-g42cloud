package sms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/sms/v3/tasks"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getMigrationTaskResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.SmsV3Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SMS client: %s", err)
	}

	return tasks.Get(client, state.Primary.ID)
}

func TestAccMigrationTask_basic(t *testing.T) {
	var migration tasks.MigrateTask
	name := acceptance.RandomAccResourceName()
	resourceName := "g42cloud_sms_task.migration"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&migration,
		getMigrationTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSms(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMigrationTask_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "state", "READY"),
					resource.TestCheckResourceAttr(resourceName, "use_public_ip", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "target_server_disks.#"),
					resource.TestCheckResourceAttrSet(resourceName, "target_server_disks.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "target_server_disks.0.size"),
					resource.TestCheckResourceAttrPair(resourceName, "vm_template_id",
						"g42cloud_sms_server_template.test", "id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"use_public_ip", "syncing", "action"},
			},
		},
	})
}

func testAccMigrationTask_basic(name string) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

data "g42cloud_sms_source_servers" "source" {
  name = "%s"
}

resource "g42cloud_sms_server_template" "test" {
  name              = "%s"
  availability_zone = data.g42cloud_availability_zones.test.names[0]
}

resource "g42cloud_sms_task" "migration" {
  type             = "MIGRATE_FILE"
  os_type          = "LINUX"
  source_server_id = data.g42cloud_sms_source_servers.source.servers[0].id
  vm_template_id   = g42cloud_sms_server_template.test.id
}
`, acceptance.G42_SMS_SOURCE_SERVER, name)
}
