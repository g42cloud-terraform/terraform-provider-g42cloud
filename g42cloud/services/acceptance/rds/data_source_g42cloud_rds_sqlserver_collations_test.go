package rds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccDatasourceSQLServerCollations_basic(t *testing.T) {
	rName := "data.g42cloud_rds_sqlserver_collations.collations"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceSQLServerCollations_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "char_sets.#"),
				),
			},
		},
	})
}

func testAccDatasourceSQLServerCollations_basic() string {
	return `
data "g42cloud_rds_sqlserver_collations" "collations" {}
`
}
