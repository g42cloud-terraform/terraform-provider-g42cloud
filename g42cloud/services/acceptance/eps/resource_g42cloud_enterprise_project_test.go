package eps

import (
	"fmt"
	"testing"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/chnsz/golangsdk/openstack/eps/v1/enterpriseprojects"
)

func getResourceEnterpriseProject(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	epsClient, err := config.EnterpriseProjectClient(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Unable to create G42Cloud EPS client : %s", err)
	}

	return enterpriseprojects.Get(epsClient, state.Primary.ID).Extract()

}

func TestAccEnterpriseProject_basic(t *testing.T) {
	var project enterpriseprojects.Project
	rName := acceptance.RandomAccResourceName()
	updateName := rName + "update"
	resourceName := "g42cloud_enterprise_project.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&project,
		getResourceEnterpriseProject,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckEnterpriseProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnterpriseProject_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test"),
					resource.TestCheckResourceAttr(resourceName, "status", "1"),
				),
			},
			{
				Config: testAccEnterpriseProject_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform test update"),
					resource.TestCheckResourceAttr(resourceName, "status", "1"),
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

func testAccCheckEnterpriseProjectDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	epsClient, err := config.EnterpriseProjectClient(acceptance.G42_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Unable to create G42Cloud EPS client : %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "g42cloud_enterprise_project" {
			continue
		}

		project, err := enterpriseprojects.Get(epsClient, rs.Primary.ID).Extract()
		if err == nil {
			if project.Status != 2 {
				return fmt.Errorf("Project still active")
			}
		}
	}

	return nil
}

func testAccEnterpriseProject_basic(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_enterprise_project" "test" {
  name        = "%s"
  description = "terraform test"
}`, rName)
}

func testAccEnterpriseProject_update(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_enterprise_project" "test" {
  name        = "%s"
  description = "terraform test update"
}`, rName)
}
