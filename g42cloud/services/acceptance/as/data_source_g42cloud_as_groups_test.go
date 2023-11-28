package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccDataSourceASGroup_basic(t *testing.T) {
	dataSourceName := "data.g42cloud_as_groups.groups"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceASGroup_conf(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "groups.0.scaling_group_name", name),
				),
			},
		},
	})
}

func testAccDataSourceASGroup_conf(name string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_as_groups" "groups" {
  name                  = g42cloud_as_group.hth_as_group.scaling_group_name
  status                = "INSERVICE"
  enterprise_project_id = "0"

  depends_on = [g42cloud_as_group.hth_as_group]
}
`, testAsGroup_base(name))
}
