package acceptance

import (
	"fmt"
	"testing"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCCEClustersDataSource_basic(t *testing.T) {
	dataSourceName := "data.g42cloud_cce_clusters.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEClustersDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "clusters.0.name", rName),
					resource.TestCheckResourceAttr(dataSourceName, "clusters.0.status", "Available"),
					resource.TestCheckResourceAttr(dataSourceName, "clusters.0.cluster_type", "VirtualMachine"),
				),
			},
		},
	})
}

func testAccCCEClustersDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_cce_clusters" "test" {
  name = g42cloud_cce_cluster.test.name

  depends_on = [g42cloud_cce_cluster.test]
}
`, testAccCCEClusterV3_base(rName))
}
