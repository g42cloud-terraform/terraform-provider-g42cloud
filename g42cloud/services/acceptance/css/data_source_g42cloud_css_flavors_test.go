package css

import (
	"testing"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCssFlavorsDataSource_basic(t *testing.T) {
	dataSourceName := "data.g42cloud_css_flavors.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCssFlavors_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.type", "ess"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.version", "7.9.3"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.region"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.memory"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.disk_range"),
				),
			},
		},
	})
}

const testAccDataSourceCssFlavors_basic = `
data "g42cloud_css_flavors" "test" {
  type    = "ess"
  version = "7.9.3"
}
`

func TestAccCssFlavorsDataSource_all(t *testing.T) {
	dataSourceName := "data.g42cloud_css_flavors.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCssFlavors_all,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.type", "ess"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.version", "7.9.3"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.id"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.region", "ae-ad-1"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.name", "ess.spec-ds.8xlarge.8"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.memory", "256"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.vcpus", "32"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.disk_range"),
				),
			},
		},
	})
}

const testAccDataSourceCssFlavors_all = `
data "g42cloud_css_flavors" "test" {
  type    = "ess"
  version = "7.9.3"
  vcpus   = 32
  memory  = 256
  region  = "ae-ad-1"
  name    = "ess.spec-ds.8xlarge.8"
}
`
