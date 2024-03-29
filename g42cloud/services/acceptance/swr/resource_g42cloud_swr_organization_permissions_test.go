package swr

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/swr/v2/namespaces"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getResourcePermissions(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	swrClient, err := conf.SwrV2Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating G42Cloud SWR client: %s", err)
	}

	return namespaces.GetAccess(swrClient, state.Primary.ID).Extract()
}

func TestAccSwrOrganizationPermissions_basic(t *testing.T) {
	var permissions namespaces.Access
	organizationName := acceptance.RandomAccResourceName()
	userName1 := acceptance.RandomAccResourceName()
	userName2 := acceptance.RandomAccResourceName()
	resourceName := "g42cloud_swr_organization_permissions.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&permissions,
		getResourcePermissions,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccswrOrganizationPermissions_basic(organizationName, userName1, userName2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "organization",
						"${g42cloud_swr_organization.test.name}"),
					resource.TestCheckResourceAttr(resourceName, "users.0.user_name", userName1),
					resource.TestCheckResourceAttr(resourceName, "users.0.permission", "Read"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccswrOrganizationPermissions_update(organizationName, userName1, userName2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "organization",
						"${g42cloud_swr_organization.test.name}"),
					resource.TestCheckResourceAttr(resourceName, "users.0.user_name", userName1),
					resource.TestCheckResourceAttr(resourceName, "users.0.permission", "Write"),
					resource.TestCheckResourceAttr(resourceName, "users.1.user_name", userName2),
					resource.TestCheckResourceAttr(resourceName, "users.1.permission", "Read"),
				),
			},
		},
	})
}

func testAccswrOrganizationPermissions_basic(organizationName, userName1, userName2 string) string {
	return fmt.Sprintf(`
resource "g42cloud_swr_organization" "test" {
  name = "%s"
}

resource "g42cloud_identity_user" "user_1" {
  name     = "%s"
  enabled  = true
  password = "password12345!"
}

resource "g42cloud_swr_organization_permissions" "test" {
  organization = g42cloud_swr_organization.test.name

  users {
    user_name  = g42cloud_identity_user.user_1.name
    user_id    = g42cloud_identity_user.user_1.id
    permission = "Read"
  }
}
`, organizationName, userName1)
}

func testAccswrOrganizationPermissions_update(organizationName, userName1, userName2 string) string {
	return fmt.Sprintf(`
resource "g42cloud_swr_organization" "test" {
  name = "%s"
}

resource "g42cloud_identity_user" "user_1" {
  name     = "%s"
  enabled  = true
  password = "password12345!"
}

resource "g42cloud_identity_user" "user_2" {
  name     = "%s"
  enabled  = true
  password = "password12345!"
}

resource "g42cloud_swr_organization_permissions" "test" {
  organization = g42cloud_swr_organization.test.name

  users {
    user_name  = g42cloud_identity_user.user_1.name
    user_id    = g42cloud_identity_user.user_1.id
    permission = "Write"
  }

  users {
    user_name  = g42cloud_identity_user.user_2.name
    user_id    = g42cloud_identity_user.user_2.id
    permission = "Read"
  }
}
`, organizationName, userName1, userName2)
}
