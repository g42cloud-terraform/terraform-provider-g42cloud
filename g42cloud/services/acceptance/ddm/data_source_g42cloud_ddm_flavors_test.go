package ddm

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccDatasourceDdmFlavors_basic(t *testing.T) {
	rName := "data.g42cloud_ddm_flavors.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdmFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "flavors.#", "1"),
					resource.TestCheckResourceAttr(rName, "flavors.0.cpu_arch", "X86"),
					resource.TestCheckResourceAttr(rName, "flavors.0.code", "ddm.c6.2xlarge.2"),
					resource.TestCheckResourceAttr(rName, "flavors.0.vcpus", "8"),
					resource.TestCheckResourceAttr(rName, "flavors.0.memory", "16"),
				),
			},
		},
	})
}

func TestAccDatasourceDdmFlavors_arm_basic(t *testing.T) {
	rName := "data.g42cloud_ddm_flavors.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDdmFlavors_arm_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "flavors.#", "1"),
					resource.TestCheckResourceAttr(rName, "flavors.0.cpu_arch", "ARM"),
					resource.TestCheckResourceAttr(rName, "flavors.0.code", "ddm.kc1.2xlarge.2"),
					resource.TestCheckResourceAttr(rName, "flavors.0.vcpus", "8"),
					resource.TestCheckResourceAttr(rName, "flavors.0.memory", "16"),
				),
			},
		},
	})
}

func testAccDatasourceDdmFlavors_basic() string {
	return `
data "g42cloud_ddm_engines" test {
}

data "g42cloud_ddm_flavors" "test" {
  engine_id = data.g42cloud_ddm_engines.test.engines[0].id
  cpu_arch  = "X86"
  code      = "ddm.c6.2xlarge.2"
  vcpus     = 8
  memory    = 16
}
`
}

func testAccDatasourceDdmFlavors_arm_basic() string {
	return `
data "g42cloud_ddm_engines" test {
}

data "g42cloud_ddm_flavors" "test" {
  engine_id = data.g42cloud_ddm_engines.test.engines[0].id
  cpu_arch  = "ARM"
  code      = "ddm.kc1.2xlarge.2"
  vcpus     = 8
  memory    = 16
}
`
}
