package modelarts

import (
	"fmt"
	"testing"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDatasetVersions_basic(t *testing.T) {
	dataSourceName := "data.g42cloud_modelarts_dataset_versions.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	name := acceptance.RandomAccResourceName()
	obsName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDatasetVersions_basic(name, obsName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.id",
						"g42cloud_modelarts_dataset_version.test", "version_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.description",
						"g42cloud_modelarts_dataset_version.test", "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.name",
						"g42cloud_modelarts_dataset_version.test", "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.split_ratio",
						"g42cloud_modelarts_dataset_version.test", "split_ratio"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.status",
						"g42cloud_modelarts_dataset_version.test", "status"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.name",
						"g42cloud_modelarts_dataset_version.test", "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.files",
						"g42cloud_modelarts_dataset_version.test", "files"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.storage_path",
						"g42cloud_modelarts_dataset_version.test", "storage_path"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.is_current",
						"g42cloud_modelarts_dataset_version.test", "is_current"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.created_at",
						"g42cloud_modelarts_dataset_version.test", "created_at"),
					resource.TestCheckResourceAttrPair(dataSourceName, "versions.0.updated_at",
						"g42cloud_modelarts_dataset_version.test", "updated_at"),
				),
			},
			{
				Config: testAccDataSourceDatasetVersions_name(name, obsName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "versions.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceDatasetVersions_basic(rName, obsName string) string {
	datasetVersion := testAccDatasetVersion_basic(rName, obsName)
	return fmt.Sprintf(`
%s

data "g42cloud_modelarts_dataset_versions" "test" {
  dataset_id  = g42cloud_modelarts_dataset.test.id
  split_ratio = "0,2.9"

  depends_on = [
    g42cloud_modelarts_dataset_version.test
  ]
}
`, datasetVersion)
}

func testAccDataSourceDatasetVersions_name(rName, obsName string) string {
	datasetVersion := testAccDatasetVersion_basic(rName, obsName)
	return fmt.Sprintf(`
%s

data "g42cloud_modelarts_dataset_versions" "test" {
  dataset_id  = g42cloud_modelarts_dataset.test.id
  split_ratio = "0,2.9"
  name        = "wrong_name"

  depends_on = [
    g42cloud_modelarts_dataset_version.test
  ]
}
`, datasetVersion)
}
