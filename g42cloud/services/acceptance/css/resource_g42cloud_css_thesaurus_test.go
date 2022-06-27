package css

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/css/v1/thesaurus"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccCssThesaurus_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "g42cloud_css_thesaurus.test"
	bucketName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCssThesaurusDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCssThesaurus_basic(rName, bucketName, "main.txt"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssThesaurusExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "bucket_name"),
					resource.TestCheckResourceAttr(resourceName, "main_object", "main.txt"),
				),
			},
			{
				Config: testAccCssThesaurus_basic(rName, bucketName, "main2.txt"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssThesaurusExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "bucket_name"),
					resource.TestCheckResourceAttr(resourceName, "main_object", "main2.txt"),
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

func testAccCssThesaurus_basic(rName string, bucketName string, obsObjectKey string) string {
	cssClusterBasic := testAccCssCluster_basic(rName, 1, 1, "value")

	return fmt.Sprintf(`
%s

resource "g42cloud_obs_bucket" "test" {
  bucket = "%s"
  acl    = "private"
}


resource "g42cloud_obs_bucket_object" "test" {
  bucket       = g42cloud_obs_bucket.test.bucket
  key          = "%s"
  content      = "123"
  content_type = "text/plain"
}

resource "g42cloud_css_thesaurus" "test" {
  cluster_id  = g42cloud_css_cluster.test.id
  bucket_name = g42cloud_obs_bucket.test.bucket
  main_object = g42cloud_obs_bucket_object.test.key
}

`, cssClusterBasic, bucketName, obsObjectKey)
}

func testAccCheckCssThesaurusDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.CssV1Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("error creating CSS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "g42cloud_css_thesaurus" {
			continue
		}

		resp, getErr := thesaurus.Get(client, rs.Primary.ID)
		if getErr != nil {
			if _, ok := getErr.(golangsdk.ErrDefault404); !ok {
				return fmtp.Errorf("Get CSS thesaurus failed.error=%s", getErr)
			}
		} else {
			if resp.Bucket != "" {
				return fmtp.Errorf("CSS thesaurus still exists, cluster_id:%s", rs.Primary.ID)
			}
		}

	}

	return nil
}

func testAccCheckCssThesaurusExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.CssV1Client(acceptance.G42_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("error creating CSS client: %s", err)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmtp.Errorf("Error checking g42cloud_css_thesaurus exist, err=not found this resource")
		}

		resp, errQueryDetail := thesaurus.Get(client, rs.Primary.ID)
		if errQueryDetail != nil {
			return fmtp.Errorf("error checking g42cloud_css_thesaurus exist,err=send request failed:%s", errQueryDetail)
		}

		if resp == nil || resp.Bucket == "" {
			return fmtp.Errorf("CSS thesaurus don't exists, cluster_id:%s", rs.Primary.ID)
		}

		return nil
	}
}
