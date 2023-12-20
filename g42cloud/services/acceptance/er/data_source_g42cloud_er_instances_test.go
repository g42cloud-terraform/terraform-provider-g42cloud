package er

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccInstancesDataSource_basic(t *testing.T) {
	var (
		dName    = "data.g42cloud_er_instances.filter_by_name"
		name     = acceptance.RandomAccResourceName()
		bgpAsNum = acctest.RandIntRange(64512, 65534)

		dc = acceptance.InitDataSourceCheck(dName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckER(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstancesDataSource_filterByName(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccInstancesDataSource_base(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

resource "g42cloud_er_instance" "test" {
  availability_zones    = slice(data.g42cloud_availability_zones.test.names, 0, 1)
  name                  = "%[1]s"
  asn                   = %[2]d
  description           = "Created by terraform test"
  enterprise_project_id = "0"

  tags = {
    foo   = "bar"
    key   = "value"
    owner = "terraform"
  }
}
`, name, bgpAsNum)
}

func testAccInstancesDataSource_filterByName(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

data "g42cloud_er_instances" "filter_by_name" {
  depends_on = [
    g42cloud_er_instance.test,
  ]

  name = g42cloud_er_instance.test.name
}

output "is_name_filter_useful" {
  value = alltrue([for v in data.g42cloud_er_instances.filter_by_name.instances[*].id : v == g42cloud_er_instance.test.id])
}
`, testAccInstancesDataSource_base(name, bgpAsNum))
}

func TestAccInstancesDataSource_filterById(t *testing.T) {
	var (
		dName    = "data.g42cloud_er_instances.filter_by_id"
		name     = acceptance.RandomAccResourceName()
		bgpAsNum = acctest.RandIntRange(64512, 65534)

		dc = acceptance.InitDataSourceCheck(dName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckER(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstancesDataSource_filterById(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccInstancesDataSource_filterById(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

data "g42cloud_er_instances" "filter_by_id" {
  instance_id = g42cloud_er_instance.test.id
}

output "is_id_filter_useful" {
  value = alltrue([for v in data.g42cloud_er_instances.filter_by_id.instances[*].id : v == g42cloud_er_instance.test.id])
}
`, testAccInstancesDataSource_base(name, bgpAsNum))
}

func TestAccInstancesDataSource_filterByStatus(t *testing.T) {
	var (
		dName    = "data.g42cloud_er_instances.filter_by_status"
		name     = acceptance.RandomAccResourceName()
		bgpAsNum = acctest.RandIntRange(64512, 65534)

		dc = acceptance.InitDataSourceCheck(dName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckER(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstancesDataSource_filterByStatus(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccInstancesDataSource_filterByStatus(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

data "g42cloud_er_instances" "filter_by_status" {
  status = g42cloud_er_instance.test.status
}

output "is_status_filter_useful" {
  value = alltrue([for v in data.g42cloud_er_instances.filter_by_status.instances[*].status : v == g42cloud_er_instance.test.status])
}
`, testAccInstancesDataSource_base(name, bgpAsNum))
}

func TestAccInstancesDataSource_filterByEpsId(t *testing.T) {
	var (
		dName    = "data.g42cloud_er_instances.filter_by_eps_id"
		name     = acceptance.RandomAccResourceName()
		bgpAsNum = acctest.RandIntRange(64512, 65534)

		dc = acceptance.InitDataSourceCheck(dName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckER(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstancesDataSource_filterByEpsId(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccInstancesDataSource_filterByEpsId(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

data "g42cloud_er_instances" "filter_by_eps_id" {
  depends_on = [
    g42cloud_er_instance.test,
  ]

  // Query all instances belonging to the default enterprise project.
  enterprise_project_id = "0"
}

output "is_eps_id_filter_useful" {
  value = alltrue([for v in data.g42cloud_er_instances.filter_by_eps_id.instances[*].enterprise_project_id : v == g42cloud_er_instance.test.enterprise_project_id])
}
`, testAccInstancesDataSource_base(name, bgpAsNum))
}

func TestAccInstancesDataSource_filterByTags(t *testing.T) {
	var (
		dName    = "data.g42cloud_er_instances.filter_by_tags"
		name     = acceptance.RandomAccResourceName()
		bgpAsNum = acctest.RandIntRange(64512, 65534)

		dc = acceptance.InitDataSourceCheck(dName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckER(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInstancesDataSource_filterByTags(name, bgpAsNum),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccInstancesDataSource_filterByTags(name string, bgpAsNum int) string {
	return fmt.Sprintf(`
%[1]s

data "g42cloud_er_instances" "filter_by_tags" {
  depends_on = [
    g42cloud_er_instance.test,
  ]

  tags = {
    foo = "bar"
    key = "value"
  }
}

output "is_tags_filter_is_useful" {
  value = alltrue([for v in data.g42cloud_er_instances.filter_by_tags.instances[*].id : v == g42cloud_er_instance.test.id])
}
`, testAccInstancesDataSource_base(name, bgpAsNum))
}
