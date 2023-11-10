package acceptance

import (
	"fmt"
	"testing"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNodesDataSource_basic(t *testing.T) {
	dataSourceName := "data.g42cloud_cce_nodes.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNodesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "nodes.0.name", rName),
				),
			},
		},
	})
}

func testAccNodesDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_cce_nodes" "test" {
  cluster_id = g42cloud_cce_cluster.test.id
  name       = g42cloud_cce_node.test.name

  depends_on = [g42cloud_cce_node.test]
}
`, testAccCCENodeV3_base(rName))
}
