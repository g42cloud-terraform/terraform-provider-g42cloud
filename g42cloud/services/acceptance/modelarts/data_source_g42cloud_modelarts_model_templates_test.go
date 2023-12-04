package modelarts

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccDatasourceModelTemplates_basic(t *testing.T) {
	rName := "data.g42cloud_modelarts_model_templates.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceModelTemplates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "templates.0.id"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.name"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.description"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.arch"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.type"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.engine"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.environment"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.template_docs.#"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.template_inputs.#"),

					resource.TestCheckOutput("type_filter_is_useful", "true"),

					resource.TestCheckOutput("engine_filter_is_useful", "true"),

					resource.TestCheckOutput("environment_filter_is_useful", "true"),

					resource.TestCheckOutput("keyword_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceModelTemplates_basic() string {
	return `

data "g42cloud_modelarts_model_templates" "test" {
}

data "g42cloud_modelarts_model_templates" "type_filter" {
  type = "Common"
}
output "type_filter_is_useful" {
  value = length(data.g42cloud_modelarts_model_templates.type_filter.templates) > 0 && alltrue(
    [for v in data.g42cloud_modelarts_model_templates.type_filter.templates[*].type : v == "Common"]
  )
}

data "g42cloud_modelarts_model_templates" "engine_filter" {
  engine = "Caffe1.0 GPU"
}
output "engine_filter_is_useful" {
  value = length(data.g42cloud_modelarts_model_templates.engine_filter.templates) > 0 && alltrue(
    [for v in data.g42cloud_modelarts_model_templates.engine_filter.templates[*].engine : v == "Caffe"]
  )
}

data "g42cloud_modelarts_model_templates" "environment_filter" {
  environment = "python3.6
}
output "environment_filter_is_useful" {
  value = length(data.g42cloud_modelarts_model_templates.environment_filter.templates) > 0 && alltrue(
    [for v in data.g42cloud_modelarts_model_templates.environment_filter.templates[*].environment : strcontains(v, "python3.6")]
  )
}

data "g42cloud_modelarts_model_templates" "keyword_filter" {
  keyword = "CPU"
}
output "keyword_filter_is_useful" {
  value = length(data.g42cloud_modelarts_model_templates.keyword_filter.templates) > 0 && alltrue(
    [for v in data.g42cloud_modelarts_model_templates.keyword_filter.templates[*].description : strcontains(v, "CPU")]
  )
}
`
}
