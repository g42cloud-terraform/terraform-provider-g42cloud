package dew

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/csms/v1/secrets"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func geCsmsSecretFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.KmsV1Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("error creating G42Cloud CSMS(KMS) client: %s", err)
	}
	name := state.Primary.Attributes["name"]
	return secrets.Get(client, name)
}

func TestAccDewCsmsSecret_basic(t *testing.T) {
	var secret secrets.Secret
	name := acceptance.RandomAccResourceName()
	resourceName := "g42cloud_csms_secret.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secret,
		geCsmsSecretFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDewCsmsSecret_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(
						resourceName,
						"secret_text",
						utils.HashAndHexEncode("this is a password"),
					),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccDewCsmsSecret_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(
						resourceName,
						"secret_text",
						utils.HashAndHexEncode(`{"password":"123456","username":"admin"}`),
					),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "new_bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.hello", "world"),
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

func testAccDewCsmsSecret_basic(name string) string {
	return fmt.Sprintf(`
resource "g42cloud_csms_secret" "test" {
  name        = "%s"
  description = "csms secret test"
  secret_text = "this is a password"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func testAccDewCsmsSecret_update(name string) string {
	return fmt.Sprintf(`
resource "g42cloud_csms_secret" "test" {
  name        = "%s"
  description = "csms secret test"

  secret_text = jsonencode({
    username = "admin"
    password = "123456"
  })

  tags = {
    foo   = "new_bar"
    hello = "world"
  }
}
`, name)
}
