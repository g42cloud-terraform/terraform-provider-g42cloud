package g42cloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCCEAddonTemplateV3DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddonTemplateV3DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.g42cloud_cce_addon_template.gpu_beta_test", "spec"),
					resource.TestCheckResourceAttrSet("data.g42cloud_cce_addon_template.autoscaler_test", "spec"),
				),
			},
		},
	})
}

func testAccCCEAddonTemplateV3DataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_cce_addon_template" "gpu_beta_test" {
  cluster_id = g42cloud_cce_cluster.test.id
  name       = "gpu-beta"
  version    = "1.1.11"
}

data "g42cloud_cce_addon_template" "autoscaler_test" {
  cluster_id = g42cloud_cce_cluster.test.id
  name       = "autoscaler"
  version    = "1.15.10"
}
`, testAccCCEClusterV3_basic(rName))
}
