package apig

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/acls"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getAclPolicyFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ApigV2Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	return acls.Get(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

// All elements' length are same.
// generateRandomStringArray is a method Used to generate the domain names and the domain IDs, and the name cannot start with a digit.
func generateRandomStringArray(count, strLen int) []string {
	if count < 1 || strLen < 1 {
		return nil
	}
	result := make([]string, count)
	for i := 0; i < count; i++ {
		result[i] = acctest.RandStringFromCharSet(strLen, "abcdef")
	}
	return result
}

func TestAccAclPolicy_basic(t *testing.T) {
	var (
		policy acls.Policy

		rName1 = "g42cloud_apig_acl_policy.ip_rule"
		rName2 = "g42cloud_apig_acl_policy.domain_rule"
		name   = acceptance.RandomAccResourceName() // The length is 13.

		basicDomainNames  = strings.Join(generateRandomStringArray(2, 4), ",")
		updateDomainNames = strings.Join(generateRandomStringArray(2, 4), ",")
		basicDomainIds    = strings.Join(generateRandomStringArray(2, 32), ",")
		updateDomainIds   = strings.Join(generateRandomStringArray(2, 32), ",")

		rc1 = acceptance.InitResourceCheck(rName1, &policy, getAclPolicyFunc)
		rc2 = acceptance.InitResourceCheck(rName2, &policy, getAclPolicyFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc1.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccApigAclPolicy_basic_step1(name, basicDomainNames, basicDomainIds),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "name", name+"_rule_ip"),
					resource.TestCheckResourceAttr(rName1, "type", "PERMIT"),
					resource.TestCheckResourceAttr(rName1, "entity_type", "IP"),
					resource.TestCheckResourceAttr(rName1, "value", "10.201.33.4,10.30.2.15"),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "name", name+"_rule_domain"),
					resource.TestCheckResourceAttr(rName2, "type", "PERMIT"),
					resource.TestCheckResourceAttr(rName2, "entity_type", "DOMAIN"),
					resource.TestCheckResourceAttr(rName2, "value", basicDomainNames),
				),
			},
			{
				Config: testAccApigAclPolicy_basic_step2(name, updateDomainNames, updateDomainIds),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "name", name+"_rule_ip_update"),
					resource.TestCheckResourceAttr(rName1, "type", "DENY"),
					resource.TestCheckResourceAttr(rName1, "entity_type", "IP"),
					resource.TestCheckResourceAttr(rName1, "value", "10.201.33.8,10.30.2.23"),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "name", name+"_rule_domain_update"),
					resource.TestCheckResourceAttr(rName2, "type", "DENY"),
					resource.TestCheckResourceAttr(rName2, "entity_type", "DOMAIN"),
					resource.TestCheckResourceAttr(rName2, "value", updateDomainNames),
				),
			},
			{
				ResourceName:      rName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAclPolicyImportStateFunc(rName1),
			},
			{
				ResourceName:      rName2,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAclPolicyImportStateFunc(rName2),
			},
		},
	})
}

func testAccAclPolicyImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}

func testAccApigAclPolicy_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "g42cloud_availability_zones" "test" {}

resource "g42cloud_apig_instance" "test" {
  name                  = "%[2]s"
  edition               = "BASIC"
  vpc_id                = g42cloud_vpc.test.id
  subnet_id             = g42cloud_vpc_subnet.test.id
  security_group_id     = g42cloud_networking_secgroup.test.id
  enterprise_project_id = "0"

  availability_zones = [
    data.g42cloud_availability_zones.test.names[0],
  ]
}
`, common.TestBaseNetwork(name), name)
}

func testAccApigAclPolicy_basic_step1(name, domainNames, domainIds string) string {
	return fmt.Sprintf(`
%[1]s

resource "g42cloud_apig_acl_policy" "ip_rule" {
  instance_id = g42cloud_apig_instance.test.id
  name        = "%[2]s_rule_ip"
  type        = "PERMIT"
  entity_type = "IP"
  value       = "10.201.33.4,10.30.2.15"
}

resource "g42cloud_apig_acl_policy" "domain_rule" {
  instance_id = g42cloud_apig_instance.test.id
  name        = "%[2]s_rule_domain"
  type        = "PERMIT"
  entity_type = "DOMAIN"
  value       = "%[3]s"
}
`, testAccApigAclPolicy_base(name), name, domainNames, domainIds)
}

func testAccApigAclPolicy_basic_step2(name, domainNames, domainIds string) string {
	return fmt.Sprintf(`
%[1]s

resource "g42cloud_apig_acl_policy" "ip_rule" {
  instance_id = g42cloud_apig_instance.test.id
  name        = "%[2]s_rule_ip_update"
  type        = "DENY"
  entity_type = "IP"
  value       = "10.201.33.8,10.30.2.23"
}

resource "g42cloud_apig_acl_policy" "domain_rule" {
  instance_id = g42cloud_apig_instance.test.id
  name        = "%[2]s_rule_domain_update"
  type        = "DENY"
  entity_type = "DOMAIN"
  value       = "%[3]s"
}
`, testAccApigAclPolicy_base(name), name, domainNames, domainIds)
}
