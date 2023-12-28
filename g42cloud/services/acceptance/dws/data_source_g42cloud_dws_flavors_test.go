package dws

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccDwsFlavorsDataSource_basic(t *testing.T) {
	resourceName := "data.g42cloud_dws_flavors.test"
	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDwsFlavorsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.#"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.volumetype"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.availability_zones.#"),
					resource.TestCheckResourceAttr(resourceName, "flavors.0.vcpus", "4"),
				),
			},
		},
	})
}

func TestAccDwsFlavorsDataSource_memory(t *testing.T) {
	resourceName := "data.g42cloud_dws_flavors.test"
	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDwsFlavorsDataSource_memory,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.#"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.volumetype"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.availability_zones.#"),
					resource.TestCheckResourceAttr(resourceName, "flavors.0.memory", "16"),
				),
			},
		},
	})
}

func TestAccDwsFlavorsDataSource_all(t *testing.T) {
	resourceName := "data.g42cloud_dws_flavors.test"
	dc := acceptance.InitDataSourceCheck(resourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDwsFlavorsDataSource_all,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.#"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.volumetype"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.availability_zones.#"),
					resource.TestCheckResourceAttr(resourceName, "flavors.0.vcpus", "4"),
					resource.TestCheckResourceAttr(resourceName, "flavors.0.memory", "16"),
				),
			},
		},
	})
}

func TestAccDwsFlavorsDataSource_az(t *testing.T) {
	resourceName := "data.g42cloud_dws_flavors.test"
	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDwsFlavorsDataSource_az,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.#"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.volumetype"),
					resource.TestCheckResourceAttrSet(resourceName, "flavors.0.size"),
				),
			},
		},
	})
}

const testAccDwsFlavorsDataSource_basic = `
data "g42cloud_dws_flavors" "test" {
  vcpus = 4
}
`

const testAccDwsFlavorsDataSource_memory = `
data "g42cloud_dws_flavors" "test" {
  memory = 16
}
`

const testAccDwsFlavorsDataSource_all = `
data "g42cloud_availability_zones" "test" {}

data "g42cloud_dws_flavors" "test" {
  vcpus             = 4
  memory            = 16
  availability_zone = data.g42cloud_availability_zones.test.names[0]
}
`

const testAccDwsFlavorsDataSource_az = `
data "g42cloud_availability_zones" "test" {}

data "g42cloud_dws_flavors" "test" {
  availability_zone = data.g42cloud_availability_zones.test.names[2]
}
`
