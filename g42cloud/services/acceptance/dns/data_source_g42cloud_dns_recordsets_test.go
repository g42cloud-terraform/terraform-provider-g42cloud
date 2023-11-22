package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccDatasourceDNSRecordsets_basic(t *testing.T) {
	rName := "data.g42cloud_dns_recordsets.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := fmt.Sprintf("acpttest-recordset-%s.com.", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDNSRecordsets_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.id"),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.name"),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.zone_id"),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.zone_name"),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.type"),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.ttl"),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.records.#"),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.status"),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.default"),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.line_id"),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.weight"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("recordset_id_filter_is_useful", "true"),

					resource.TestCheckResourceAttr("data.g42cloud_dns_recordsets.tags_filter", "recordsets.#", "1"),
				),
			},
		},
	})
}

func TestAccDatasourceDNSRecordsets_private(t *testing.T) {
	rName := "data.g42cloud_dns_recordsets.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := fmt.Sprintf("acpttest-recordset-%s.com.", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDNSRecordsets_private(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.id"),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.name"),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.zone_id"),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.zone_name"),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.type"),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.ttl"),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.records.#"),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.status"),
					resource.TestCheckResourceAttrSet(rName, "recordsets.0.default"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("recordset_id_filter_is_useful", "true"),

					resource.TestCheckResourceAttr("data.g42cloud_dns_recordsets.tags_filter", "recordsets.#", "1"),
				),
			},
		},
	})
}

func testAccDatasourceDNSRecordsets_basic(name string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_dns_recordsets" "test" {
  zone_id = g42cloud_dns_recordset.test.zone_id
}

data "g42cloud_dns_recordsets" "line_id_filter" {
  zone_id = g42cloud_dns_recordset.test.zone_id
}
data "g42cloud_dns_recordsets" "status_filter" {
  zone_id = g42cloud_dns_recordset.test.zone_id
  status  = "ACTIVE"
}
data "g42cloud_dns_recordsets" "type_filter" {
  zone_id = g42cloud_dns_recordset.test.zone_id
  type    = g42cloud_dns_recordset.test.type
}
data "g42cloud_dns_recordsets" "name_filter" {
  zone_id = g42cloud_dns_recordset.test.zone_id
  name    = g42cloud_dns_recordset.test.name
}
data "g42cloud_dns_recordsets" "recordset_id_filter" {
  zone_id      = g42cloud_dns_recordset.test.zone_id
  recordset_id = split("/", g42cloud_dns_recordset.test.id).1
}

locals {
  status_filter_result = [for v in data.g42cloud_dns_recordsets.status_filter.recordsets[*].status : v == "ACTIVE"]
  type_filter_result = [for v in data.g42cloud_dns_recordsets.type_filter.recordsets[*].type : v == g42cloud_dns_recordset.test.type]
  name_filter_result = [for v in data.g42cloud_dns_recordsets.name_filter.recordsets[*].name : v == g42cloud_dns_recordset.test.name]
  recordset_id_filter_result = [for v in data.g42cloud_dns_recordsets.recordset_id_filter.recordsets[*].id :
v == split("/", g42cloud_dns_recordset.test.id).1]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

output "type_filter_is_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

output "recordset_id_filter_is_useful" {
  value = alltrue(local.recordset_id_filter_result) && length(local.recordset_id_filter_result) > 0
}

data "g42cloud_dns_recordsets" "tags_filter" {
  zone_id = g42cloud_dns_recordset.test.zone_id
  tags    = "key1,value1"
}
`, testDNSRecordset_basic(name))
}

func testAccDatasourceDNSRecordsets_private(name string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_dns_recordsets" "test" {
  zone_id = g42cloud_dns_recordset.test.zone_id
}

data "g42cloud_dns_recordsets" "status_filter" {
  zone_id = g42cloud_dns_recordset.test.zone_id
  status  = "ACTIVE"
}
data "g42cloud_dns_recordsets" "type_filter" {
  zone_id = g42cloud_dns_recordset.test.zone_id
  type    = g42cloud_dns_recordset.test.type
}
data "g42cloud_dns_recordsets" "name_filter" {
  zone_id = g42cloud_dns_recordset.test.zone_id
  name    = g42cloud_dns_recordset.test.name
}
data "g42cloud_dns_recordsets" "recordset_id_filter" {
  zone_id      = g42cloud_dns_recordset.test.zone_id
  recordset_id = split("/", g42cloud_dns_recordset.test.id).1
}

locals {
  status_filter_result = [for v in data.g42cloud_dns_recordsets.status_filter.recordsets[*].status : v == "ACTIVE"]
  type_filter_result = [for v in data.g42cloud_dns_recordsets.type_filter.recordsets[*].type : v == g42cloud_dns_recordset.test.type]
  name_filter_result = [for v in data.g42cloud_dns_recordsets.name_filter.recordsets[*].name : v == g42cloud_dns_recordset.test.name]
  recordset_id_filter_result = [for v in data.g42cloud_dns_recordsets.recordset_id_filter.recordsets[*].id :
v == split("/", g42cloud_dns_recordset.test.id).1]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

output "type_filter_is_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

output "recordset_id_filter_is_useful" {
  value = alltrue(local.recordset_id_filter_result) && length(local.recordset_id_filter_result) > 0
}

data "g42cloud_dns_recordsets" "tags_filter" {
  zone_id = g42cloud_dns_recordset.test.zone_id
  tags    = "foo,bar_private"
}
`, testDNSRecordset_privateZone(name))
}
