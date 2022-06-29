package servicestage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/servicestage/v1/repositories"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccRepoPwdAuth_basic(t *testing.T) {
	var (
		auth         repositories.Authorization
		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "g42cloud_servicestage_repo_password_authorization.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&auth,
		getAuthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRepoPwdAuth(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRepoPwdAuth_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "type", "devcloud"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"user_name",
					"password",
				},
			},
		},
	})
}

func testAccRepoPwdAuth_basic(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_servicestage_repo_password_authorization" "test" {
  type      = "devcloud"
  name      = "%s"
  user_name = "%s/%s"
  password  = "%s"
}
`, rName, acceptance.G42_ACCOUNT_NAME, acceptance.G42_USERNAME, acceptance.G42_GITHUB_REPO_PWD)
}
