package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/mappings"
	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/providers"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getProviderConversionFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IAMNoVersionClient(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client without version: %s", err)
	}
	providerID := state.Primary.Attributes["provider_id"]
	conversionID := "mapping_" + providerID
	return mappings.Get(client, conversionID)
}

func TestAccIdentityProviderConversion_basic(t *testing.T) {
	var provider providers.Provider
	var name = acceptance.RandomAccResourceName()
	resourceName := "g42cloud_identity_provider_conversion.conversion"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&provider,
		getProviderConversionFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderConversion_conf(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.0.local.0.username", "Tom"),
				),
			},
			{
				Config: testAccIdentityProviderConversion_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.0.local.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.0.local.0.username", "Tom"),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.1.remote.0.value.#", "2"),
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

func testAccIdentityProviderConversion_conf(name string) string {
	return fmt.Sprintf(`
resource "g42cloud_identity_provider" "provider_1" {
  name     = "%s"
  protocol = "oidc"
}

resource "g42cloud_identity_provider_conversion" "conversion" {
  provider_id = g42cloud_identity_provider.provider_1.id

  conversion_rules {
    local {
      username = "Tom"
    }
    remote {
      attribute = "Tom"
    }
  }
}
`, name)
}

func testAccIdentityProviderConversion_update(name string) string {
	return fmt.Sprintf(`
resource "g42cloud_identity_provider" "provider_1" {
  name     = "%s"
  protocol = "oidc"
}

resource "g42cloud_identity_provider_conversion" "conversion" {
  provider_id = g42cloud_identity_provider.provider_1.id

  conversion_rules {
    local {
      username = "Tom"
    }
    local {
      username = "federateduser"
    }
    remote {
      attribute = "Tom"
    }
    remote {
      attribute = "federatedgroup"
    }
  }

  conversion_rules {
    local {
      username = "Jams"
    }
    remote {
      attribute = "username"
      condition = "any_one_of"
      value     = ["Tom", "Jerry"]
    }
  }
}
`, name)
}
