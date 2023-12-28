package modelarts

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccDatasourceResourceFlavors_basic(t *testing.T) {
	rName := "data.g42cloud_modelarts_resource_flavors.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceResourceFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.id"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.type"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.arch"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.cpu"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.memory"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.volume.#"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.billing_modes.#"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.az_status.%"),
					resource.TestCheckResourceAttrSet(rName, "flavors.0.tags.%"),

					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceResourceFlavors_basic() string {
	return `
data "g42cloud_modelarts_resource_flavors" "test" {
}

data "g42cloud_modelarts_resource_flavors" "type_filter" {
  type = "Dedicate"
}
output "type_filter_is_useful" {
  value = length(data.g42cloud_modelarts_resource_flavors.type_filter.flavors) > 0 && alltrue(
    [for v in data.g42cloud_modelarts_resource_flavors.type_filter.flavors[*].type : v == "Dedicate"]
  )
}
`
}
