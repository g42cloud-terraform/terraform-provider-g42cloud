package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/swr/v2/namespaces"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getResourceOrganization(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	swrClient, err := conf.SwrV2Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating G42Cloud SWR client: %s", err)
	}

	return namespaces.Get(swrClient, state.Primary.ID).Extract()
}

func TestAccSWROrganization_basic(t *testing.T) {
	var org namespaces.Namespace
	rName := acceptance.RandomAccResourceName()
	resourceName := "g42cloud_swr_organization.test"
	loginServer := fmt.Sprintf("swr.%s.g42cloud.com", acceptance.G42_REGION_NAME)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&org,
		getResourceOrganization,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSWROrganization_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "permission", "Manage"),
					resource.TestCheckResourceAttr(resourceName, "login_server", loginServer),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSWROrganization_basic(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_swr_organization" "test" {
  name = "%s"
}
`, rName)
}
