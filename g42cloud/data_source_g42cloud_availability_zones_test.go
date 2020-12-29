package g42cloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAvailabilityZones_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAvailabilityZonesConfig_all,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.g42cloud_availability_zones.all", "names.#", regexp.MustCompile("[1-9]\\d*")),
				),
			},
		},
	})
}

const testAccAvailabilityZonesConfig_all = `
data "g42cloud_availability_zones" "all" {}
`
