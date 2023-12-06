package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccDatasourceBackup_basic(t *testing.T) {
	rName := "data.g42cloud_rds_backups.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceBackup_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "backups.0.id", "g42cloud_rds_backup.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "backups.0.name", "g42cloud_rds_backup.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "backups.0.instance_id",
						"g42cloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "backups.0.type", "manual"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.size"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.status"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.begin_time"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.end_time"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.associated_with_ddm"),
					resource.TestCheckResourceAttr(rName, "backups.0.datastore.#", "1"),
					resource.TestCheckResourceAttr(rName, "backups.0.databases.#", "0"),
				),
			},
		},
	})
}

func testAccDatasourceBackup_basic() string {
	backupConfig := testBackup_mysql_basic(acceptance.RandomAccResourceName())
	return fmt.Sprintf(`
%s 

data "g42cloud_rds_backups" "test" {
  instance_id = g42cloud_rds_instance.test.id
  backup_type = "manual"

  depends_on = [
    g42cloud_rds_backup.test
  ]
}
`, backupConfig)
}
