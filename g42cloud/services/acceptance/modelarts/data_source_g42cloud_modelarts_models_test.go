package modelarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccDatasourceModels_basic(t *testing.T) {
	rName := "data.g42cloud_modelarts_models.name_filter"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()
	name2 := acceptance.RandomAccResourceName()
	bucketName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceModels_basic(name, name2, bucketName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "models.0.id", "g42cloud_modelarts_model.test1", "id"),
					resource.TestCheckResourceAttr(rName, "models.0.name", name),
					resource.TestCheckResourceAttr(rName, "models.0.version", "0.0.1"),
					resource.TestCheckResourceAttr(rName, "models.0.model_type", "TensorFlow"),
					resource.TestCheckResourceAttr(rName, "models.0.description", "This is a demo"),
					resource.TestCheckResourceAttr(rName, "models.0.workspace_id", "0"),
					resource.TestCheckResourceAttr(rName, "models.0.status", "published"),
					resource.TestCheckResourceAttrSet(rName, "models.0.model_source"),
					resource.TestCheckResourceAttrSet(rName, "models.0.install_type.#"),
					resource.TestCheckResourceAttrSet(rName, "models.0.size"),
					resource.TestCheckResourceAttrSet(rName, "models.0.market_flag"),
					resource.TestCheckResourceAttrSet(rName, "models.0.tunable"),
					resource.TestCheckResourceAttrSet(rName, "models.0.publishable_flag"),

					resource.TestCheckOutput("default_filter_is_useful", "true"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),

					resource.TestCheckOutput("version_filter_is_useful", "true"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),

					resource.TestCheckOutput("description_filter_is_useful", "true"),

					resource.TestCheckOutput("workspace_id_filter_is_useful", "true"),

					resource.TestCheckOutput("model_type_filter_is_useful", "true"),

					resource.TestCheckOutput("not_model_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceModels_basic(name, name2, bucketName string) string {
	modelConfig := testAccDatasourceModels_config(name, name2, bucketName)

	return fmt.Sprintf(`
%[1]s

data "g42cloud_modelarts_models" "test" {
  depends_on = [g42cloud_modelarts_model.test1, g42cloud_modelarts_model.test2]
}

output "default_filter_is_useful" {
  value = length(data.g42cloud_modelarts_models.test.models) == 2
}

data "g42cloud_modelarts_models" "name_filter" {
  name        = g42cloud_modelarts_model.test1.name
  exact_match = "true"

  depends_on = [g42cloud_modelarts_model.test1, g42cloud_modelarts_model.test2]
}
output "name_filter_is_useful" {
  value = alltrue([for v in data.g42cloud_modelarts_models.name_filter.models[*].name : v == "%[2]s"])
}

data "g42cloud_modelarts_models" "version_filter" {
  version = g42cloud_modelarts_model.test1.version

  depends_on = [g42cloud_modelarts_model.test1, g42cloud_modelarts_model.test2]
}
output "version_filter_is_useful" {
  value = alltrue([for v in data.g42cloud_modelarts_models.version_filter.models[*].version : v == "0.0.1"])
}

data "g42cloud_modelarts_models" "status_filter" {
  status = "published"

  depends_on = [g42cloud_modelarts_model.test1, g42cloud_modelarts_model.test2]
}
output "status_filter_is_useful" {
  value = alltrue([for v in data.g42cloud_modelarts_models.status_filter.models[*].status : v == "published"])
}

data "g42cloud_modelarts_models" "description_filter" {
  description = g42cloud_modelarts_model.test2.description

  depends_on = [g42cloud_modelarts_model.test1, g42cloud_modelarts_model.test2]
}
output "description_filter_is_useful" {
  value = alltrue([for v in data.g42cloud_modelarts_models.description_filter.models[*].description : v == "%[3]s"])
}

data "g42cloud_modelarts_models" "workspace_id_filter" {
  workspace_id = g42cloud_modelarts_model.test1.workspace_id

  depends_on = [g42cloud_modelarts_model.test1, g42cloud_modelarts_model.test2]
}
output "workspace_id_filter_is_useful" {
  value = alltrue([for v in data.g42cloud_modelarts_models.workspace_id_filter.models[*].workspace_id : v == "0"])
}

data "g42cloud_modelarts_models" "model_type_filter" {
  model_type = g42cloud_modelarts_model.test1.model_type

  depends_on = [g42cloud_modelarts_model.test1, g42cloud_modelarts_model.test2]
}
output "model_type_filter_is_useful" {
  value = alltrue([for v in data.g42cloud_modelarts_models.model_type_filter.models[*].model_type : v == "TensorFlow"])
}

data "g42cloud_modelarts_models" "not_model_type_filter" {
  not_model_type = "TensorFlow"

  depends_on = [g42cloud_modelarts_model.test1, g42cloud_modelarts_model.test2]
}
output "not_model_type_filter_is_useful" {
  value = alltrue([for v in data.g42cloud_modelarts_models.not_model_type_filter.models[*].model_type : v != "TensorFlow"])
}
`, modelConfig, name, name2)
}

func testAccDatasourceModels_config(name, name2, bucketName string) string {
	return fmt.Sprintf(`
resource "g42cloud_obs_bucket" "test" {
  bucket        = "%[3]s"
  acl           = "private"
  force_destroy = true
}

resource "g42cloud_obs_bucket_object" "input" {
  bucket  = g42cloud_obs_bucket.test.bucket
  key     = "input/1.txt"
  content = "some_bucket_content"

}

resource "g42cloud_modelarts_model" "test1" {
  name            = "%[1]s"
  version         = "0.0.1"
  description     = "This is a demo"
  source_location = "https://${g42cloud_obs_bucket.test.bucket_domain_name}/input"
  model_type      = "TensorFlow"
  runtime         = "python3.6"

  model_docs {
    doc_name = "guide"
    doc_url  = "https://doc.xxxx.yourdomain"
  }

  initial_config = jsonencode(
    {
      "model_algorithm" : "object_detection",
      "metrics" : {},
      "apis" : [
        {
          "url" : "/",
          "method" : "post",
          "request" : {
            "Content-type" : "multipart/form-data",
            "data" : {
              "type" : "object",
              "properties" : {
                "images" : {
                  "type" : "file"
                }
              }
            }
          },
          "response" : {
            "Content-type" : "application/json",
            "data" : {
              "type" : "object",
              "properties" : {
                "mnist_result" : {
                  "type" : "array",
                  "item" : [
                    {
                      "type" : "string"
                    }
                  ]
                }
              }
            }
          }
        }
      ]
    }
  )
}

resource "g42cloud_modelarts_model" "test2" {
  name            = "%[2]s"
  version         = "0.0.2"
  description     = "%[2]s"
  source_location = "https://${g42cloud_obs_bucket.test.bucket_domain_name}/input"
  model_type      = "PyTorch"
  runtime         = "python3.6"

  model_docs {
    doc_name = "guide"
    doc_url  = "https://doc.xxxx.yourdomain"
  }

  initial_config = jsonencode(
    {
      "model_algorithm" : "object_detection",
      "metrics" : {},
      "apis" : [
        {
          "url" : "/",
          "method" : "post",
          "request" : {
            "Content-type" : "multipart/form-data",
            "data" : {
              "type" : "object",
              "properties" : {
                "images" : {
                  "type" : "file"
                }
              }
            }
          },
          "response" : {
            "Content-type" : "application/json",
            "data" : {
              "type" : "object",
              "properties" : {
                "mnist_result" : {
                  "type" : "array",
                  "item" : [
                    {
                      "type" : "string"
                    }
                  ]
                }
              }
            }
          }
        }
      ]
    }
  )
}
`, name, name2, bucketName)
}
