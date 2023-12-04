package modelarts

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getModelartsModelResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.G42_REGION_NAME
	// getModel: Query the Modelarts model.
	var (
		getModelHttpUrl = "v1/{project_id}/models/{id}"
		getModelProduct = "modelarts"
	)
	getModelClient, err := cfg.NewServiceClient(getModelProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts Client: %s", err)
	}

	getModelPath := getModelClient.Endpoint + getModelHttpUrl
	getModelPath = strings.ReplaceAll(getModelPath, "{project_id}", getModelClient.ProjectID)
	getModelPath = strings.ReplaceAll(getModelPath, "{id}", state.Primary.ID)

	getModelOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getModelResp, err := getModelClient.Request("GET", getModelPath, &getModelOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Modelarts model: %s", err)
	}

	getModelRespBody, err := utils.FlattenResponse(getModelResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Modelarts model: %s", err)
	}

	return getModelRespBody, nil
}

func TestAccModelartsModel_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	bucketName := acceptance.RandomAccResourceNameWithDash()
	rName := "g42cloud_modelarts_model.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getModelartsModelResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testModelartsModel_basic(name, bucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "version", "0.0.2"),
					resource.TestCheckResourceAttr(rName, "description", "This is a demo"),
					resource.TestCheckResourceAttr(rName, "model_type", "PyTorch"),
					resource.TestCheckResourceAttr(rName, "runtime", "python3.6"),
					resource.TestCheckResourceAttr(rName, "version", "0.0.2"),
					resource.TestCheckResourceAttr(rName, "model_docs.0.doc_name", "guide"),
					resource.TestCheckResourceAttr(rName, "model_docs.0.doc_url", "https://doc.xxxx.yourdomain"),
					resource.TestCheckResourceAttrSet(rName, "model_size"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "model_source"),
					resource.TestCheckResourceAttrSet(rName, "tunable"),
					resource.TestCheckResourceAttrSet(rName, "market_flag"),
					resource.TestCheckResourceAttrSet(rName, "publishable_flag"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testModelartsModel_basic(name, bucketName string) string {
	return fmt.Sprintf(`
resource "g42cloud_obs_bucket" "test" {
  bucket        = "%[2]s"
  acl           = "private"
  force_destroy = true
}

resource "g42cloud_obs_bucket_object" "input" {
  bucket  = g42cloud_obs_bucket.test.bucket
  key     = "input/1.txt"
  content = "some_bucket_content"

}

resource "g42cloud_modelarts_model" "test" {
  name            = "%[1]s"
  version         = "0.0.2"
  description     = "This is a demo"
  source_location = "https://${g42cloud_obs_bucket.test.bucket_domain_name}/input"
  model_type      = "PyTorch"
  runtime         = "python3.6"

  model_docs {
    doc_name = "guide"
    doc_url  = "https://doc.xxxx.yourdomain"
  }

  dependencies {
    installer = "pip"
    packages {
      package_name = "new"
    }
  }

  initial_config = jsonencode(
    {
      "model_algorithm" : "predict_analysis",
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
`, name, bucketName)
}
