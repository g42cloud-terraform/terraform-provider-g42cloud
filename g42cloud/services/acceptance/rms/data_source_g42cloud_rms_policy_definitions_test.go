package rms

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccDataPolicyDefinitions_basic(t *testing.T) {
	var (
		dName = "data.g42cloud_rms_policy_definitions.test"
		dc    = acceptance.InitDataSourceCheck(dName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPolicyDefinitions_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dName, "definitions.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

const testAccDataPolicyDefinitions_basic = `
data "g42cloud_rms_policy_definitions" "test" {
  name = "allowed-ecs-flavors"
}
`

func TestAccDataPolicyDefinitions_keywords(t *testing.T) {
	var (
		dName = "data.g42cloud_rms_policy_definitions.test"
		dc    = acceptance.InitDataSourceCheck(dName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPolicyDefinitions_keywords,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dName, "definitions.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

const testAccDataPolicyDefinitions_keywords = `
data "g42cloud_rms_policy_definitions" "test" {
  keywords = ["ecs"]
}
`

func TestAccDataPolicyDefinitions_policyRuleType(t *testing.T) {
	var (
		dName = "data.g42cloud_rms_policy_definitions.test"
		dc    = acceptance.InitDataSourceCheck(dName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPolicyDefinitions_policyRuleType,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dName, "definitions.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

const testAccDataPolicyDefinitions_policyRuleType = `
data "g42cloud_rms_policy_definitions" "test" {
  policy_rule_type = "dsl"
}
`
