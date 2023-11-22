package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccDatasourceDNSZones_basic(t *testing.T) {
	rName := "data.g42cloud_dns_zones.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := fmt.Sprintf("acpttest-zone-%s.com.", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDNSZones_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "zones.0.id"),
					resource.TestCheckResourceAttrSet(rName, "zones.0.name"),
					resource.TestCheckResourceAttrSet(rName, "zones.0.description"),
					resource.TestCheckResourceAttrSet(rName, "zones.0.email"),
					resource.TestCheckResourceAttrSet(rName, "zones.0.zone_type"),
					resource.TestCheckResourceAttrSet(rName, "zones.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(rName, "zones.0.routers.0.router_id"),
					resource.TestCheckResourceAttrSet(rName, "zones.0.ttl"),
					resource.TestCheckResourceAttrSet(rName, "zones.0.status"),
					resource.TestCheckResourceAttrSet(rName, "zones.0.record_num"),
					resource.TestCheckResourceAttrSet(rName, "zones.0.masters.#"),

					resource.TestCheckOutput("tags_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func TestAccDatasourceDNSZones_public(t *testing.T) {
	rName := "data.g42cloud_dns_zones.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := fmt.Sprintf("acpttest-zone-%s.com.", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDNSZones_public(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "zones.0.id"),
					resource.TestCheckResourceAttrSet(rName, "zones.0.name"),
					resource.TestCheckResourceAttrSet(rName, "zones.0.description"),
					resource.TestCheckResourceAttrSet(rName, "zones.0.email"),
					resource.TestCheckResourceAttrSet(rName, "zones.0.zone_type"),
					resource.TestCheckResourceAttrSet(rName, "zones.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(rName, "zones.0.ttl"),
					resource.TestCheckResourceAttrSet(rName, "zones.0.status"),
					resource.TestCheckResourceAttrSet(rName, "zones.0.record_num"),
					resource.TestCheckResourceAttrSet(rName, "zones.0.masters.#"),

					resource.TestCheckOutput("tags_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceDNSZones_basic(name string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_dns_zones" "test" {
  zone_type = g42cloud_dns_zone.zone_1.zone_type
}

data "g42cloud_dns_zones" "tags_filter" {
  zone_type = g42cloud_dns_zone.zone_1.zone_type
  tags      = "zone_type,private"
}
data "g42cloud_dns_zones" "name_filter" {
  zone_type = g42cloud_dns_zone.zone_1.zone_type
  name      = g42cloud_dns_zone.zone_1.name
}
data "g42cloud_dns_zones" "status_filter" {
  zone_type = g42cloud_dns_zone.zone_1.zone_type
  status    = data.g42cloud_dns_zones.test.zones.0.status
}
data "g42cloud_dns_zones" "enterprise_project_id_filter" {
  zone_type             = g42cloud_dns_zone.zone_1.zone_type
  enterprise_project_id = g42cloud_dns_zone.zone_1.enterprise_project_id
}

locals {
  tags_filter_result = [for v in data.g42cloud_dns_zones.tags_filter.zones[*].tags : v.zone_type == "private"]
  name_filter_result = [for v in data.g42cloud_dns_zones.name_filter.zones[*].name : v == g42cloud_dns_zone.zone_1.name]
  status_filter_result = [for v in data.g42cloud_dns_zones.status_filter.zones[*].status : v == data.g42cloud_dns_zones.test.zones.0.status]
  enterprise_project_id_filter_result = [for v in data.g42cloud_dns_zones.enterprise_project_id_filter.zones[*].enterprise_project_id :
v == g42cloud_dns_zone.zone_1.enterprise_project_id]
}

output "tags_filter_is_useful" {
  value = alltrue(local.tags_filter_result) && length(local.tags_filter_result) > 0
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

output "enterprise_project_id_filter_is_useful" {
  value = alltrue(local.enterprise_project_id_filter_result) && length(local.enterprise_project_id_filter_result) > 0
}
`, testAccDNSZone_private(name))
}

func testAccDatasourceDNSZones_public(name string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_dns_zones" "test" {
  zone_type = g42cloud_dns_zone.zone_1.zone_type
}

data "g42cloud_dns_zones" "tags_filter" {
  zone_type = g42cloud_dns_zone.zone_1.zone_type
  tags      = "zone_type,public"
}
data "g42cloud_dns_zones" "name_filter" {
  zone_type = g42cloud_dns_zone.zone_1.zone_type
  name      = g42cloud_dns_zone.zone_1.name
}
data "g42cloud_dns_zones" "status_filter" {
  zone_type = g42cloud_dns_zone.zone_1.zone_type
  status    = data.g42cloud_dns_zones.test.zones.0.status
}
data "g42cloud_dns_zones" "enterprise_project_id_filter" {
  zone_type             = g42cloud_dns_zone.zone_1.zone_type
  enterprise_project_id = g42cloud_dns_zone.zone_1.enterprise_project_id
}

locals {
  tags_filter_result = [for v in data.g42cloud_dns_zones.tags_filter.zones[*].tags : v.zone_type == "public"]
  name_filter_result = [for v in data.g42cloud_dns_zones.name_filter.zones[*].name : v == g42cloud_dns_zone.zone_1.name]
  status_filter_result = [for v in data.g42cloud_dns_zones.status_filter.zones[*].status : v == data.g42cloud_dns_zones.test.zones.0.status]
  enterprise_project_id_filter_result = [for v in data.g42cloud_dns_zones.enterprise_project_id_filter.zones[*].enterprise_project_id :
v == g42cloud_dns_zone.zone_1.enterprise_project_id]
}

output "tags_filter_is_useful" {
  value = alltrue(local.tags_filter_result) && length(local.tags_filter_result) > 0
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

output "enterprise_project_id_filter_is_useful" {
  value = alltrue(local.enterprise_project_id_filter_result) && length(local.enterprise_project_id_filter_result) > 0
}
`, testAccDNSZone_basic(name))
}

func testDNSRecordset_privateZone(name string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_dns_recordset" "test" {
  zone_id     = g42cloud_dns_zone.zone_1.id
  name        = "%s"
  type        = "A"
  description = "a private record set"
  ttl         = 600
  records     = ["10.1.0.3"]

  tags = {
    foo = "bar_private"
    key = "value_private"
  }
}
`, testAccDNSZone_private(name), name)
}

func testDNSRecordset_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_dns_recordset" "test" {
  zone_id     = g42cloud_dns_zone.zone_1.id
  name        = "%s"
  type        = "A"
  description = "a recordset description"
  ttl         = 300
  records     = ["10.1.0.0"]

  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}
`, testAccDNSZone_basic(name), name)
}

func testAccDNSZone_basic(zoneName string) string {
	return fmt.Sprintf(`
resource "g42cloud_dns_zone" "zone_1" {
  name        = "%s"
  description = "a zone"
  ttl         = 300

  tags = {
    zone_type = "public"
    owner     = "terraform"
  }
}
`, zoneName)
}

func testAccDNSZone_private(zoneName string) string {
	return fmt.Sprintf(`
data "g42cloud_vpc" "default" {
  name = "vpc-default"
}

resource "g42cloud_dns_zone" "zone_1" {
  name        = "%s"
  email       = "email@example.com"
  description = "a private zone"
  zone_type   = "private"

  router {
    router_id = data.g42cloud_vpc.default.id
  }
  tags = {
    zone_type = "private"
    owner     = "terraform"
  }
}
`, zoneName)
}
