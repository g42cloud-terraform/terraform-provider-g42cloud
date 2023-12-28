package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dli/v2/spark/resources"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
)

func getPackageResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.DliV2Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DLI v2 client: %s", err)
	}

	return dli.GetDliDependentPackageInfo(c, state.Primary.ID)
}

func TestAccDliPackage_basic(t *testing.T) {
	var pkg resources.Resource

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "g42cloud_dli_package.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pkg,
		getPackageResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliPackage_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "pyFile"),
					resource.TestCheckResourceAttr(resourceName, "object_path", fmt.Sprintf(
						"https://%s.obs.%s.g42cloud.com/dli/packages/simple_pyspark_test_DLF_refresh.py",
						rName, acceptance.G42_REGION_NAME)),
					resource.TestCheckResourceAttr(resourceName, "object_name", "simple_pyspark_test_DLF_refresh.py"),
					resource.TestCheckResourceAttr(resourceName, "status", "READY"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
		},
	})
}

func testAccDliPackage_basic(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_obs_bucket" "test" {
  bucket = "%s"
  acl    = "private"
}

resource "g42cloud_obs_bucket_object" "test" {
  bucket  = g42cloud_obs_bucket.test.bucket
  key     = "dli/packages/simple_pyspark_test_DLF_refresh.py"
  content = <<EOF
#!/usr/bin/env python
# _*_ coding: utf-8 _*_

import sys
import logging
from operator import add
import time

from pyspark.sql import SparkSession
from pyspark.sql import SQLContext

sparkSession = SparkSession.builder.appName("simple pyspark test DLF refresh").getOrCreate()
sc = SQLContext(sparkSession.sparkContext)

logging.basicConfig(format='%%(message)s', level=logging.INFO)
logger = logging.getLogger("Whatever")
logger.info("[DBmethods.py] HELLOOOOOOOOOOO")


sc._jsc.hadoopConfiguration().set("fs.obs.access.key", "%s")
sc._jsc.hadoopConfiguration().set("fs.obs.secret.key", "%s")
sc._jsc.hadoopConfiguration().set("fs.obs.endpoint", "obs.ae-ad-1.g42cloud.com")


# Read private bucket with encryption using AK/SK
private_encrypted_file = "obs://dedicated-for-terraform-acc-test/dli/spark/people.csv"

df = sparkSession.read.options(header='True', inferSchema='True', delimiter=',').csv(private_encrypted_file)
df.show()
df.printSchema()
print(df)
print(df.count())
print(time.time())


my_string_to_print = "{} - {}".format(int(time.time()), df.count()/2)
file_name = "my_file-{}-{}".format(int(time.time()), df.count()/2)


print(my_string_to_print)
print(file_name)

private_encrypted_output_folder = "obs://dedicated-for-terraform-acc-test/dli/result/"
# my_string_to_print.write.mode('overwrite').csv(private_encrypted_output_folder)

final_path = "{}{}".format(private_encrypted_output_folder, file_name)
print(final_path)


sparkSession.sparkContext.parallelize([my_string_to_print]).coalesce(1).saveAsTextFile(final_path)


EOF
  content_type = "text/py"
}

resource "g42cloud_dli_package" "test" {
  group_name  = "%s"
  type        = "pyFile"
  object_path = "https://${g42cloud_obs_bucket.test.bucket_domain_name}/dli/packages/simple_pyspark_test_DLF_refresh.py"
}
`, rName, acceptance.G42_ACCESS_KEY, acceptance.G42_SECRET_KEY, rName)
}
