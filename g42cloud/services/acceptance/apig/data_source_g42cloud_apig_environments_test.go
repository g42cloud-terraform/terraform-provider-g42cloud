package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataEnvironments_basic(t *testing.T) {
	var (
		dataSourceName = "data.g42cloud_apig_environments.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		rName          = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t) // The creation of APIG instance needs the enterprise project ID.
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnvironments_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "environments.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

func testAccDataEnvironments_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_apig_environment" "test" {
  name        = "%s"
  instance_id = g42cloud_apig_instance.test.id
  description = "Created by script"
}

data "g42cloud_apig_environments" "test" {
  instance_id = g42cloud_apig_instance.test.id
  name        = g42cloud_apig_environment.test.name
}
`, testAccApigApplication_base(rName), rName)
}
