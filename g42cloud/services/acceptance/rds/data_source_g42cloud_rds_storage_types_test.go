package rds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccDatasourceStoragetype_basic(t *testing.T) {
	rName := "data.g42cloud_rds_storage_types.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceStoragetype_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.name"),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.az_status.%"),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.support_compute_group_type.#"),
				),
			},
		},
	})
}

func TestAccDatasourceStoragetype_PostgreSQL_basic(t *testing.T) {
	rName := "data.g42cloud_rds_storage_types.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceStoragetype_PostgreSQL_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.name"),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.az_status.%"),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.support_compute_group_type.#"),
				),
			},
		},
	})
}

func TestAccDatasourceStoragetype_SQLServer_basic(t *testing.T) {
	rName := "data.g42cloud_rds_storage_types.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceStoragetype_SQLServer_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.name"),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.az_status.%"),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.support_compute_group_type.#"),
				),
			},
		},
	})
}

func testAccDatasourceStoragetype_basic() string {
	return `
data "g42cloud_rds_storage_types" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "replica"
}`
}

func testAccDatasourceStoragetype_PostgreSQL_basic() string {
	return `
data "g42cloud_rds_storage_types" "test" {
  db_type       = "PostgreSQL"
  db_version    = "14"
  instance_mode = "ha"
}`
}

func testAccDatasourceStoragetype_SQLServer_basic() string {
	return `
data "g42cloud_rds_storage_types" "test" {
  db_type       = "SQLServer"
  db_version    = "2019_SE"
  instance_mode = "single"
}`
}
