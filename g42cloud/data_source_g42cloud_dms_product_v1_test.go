package g42cloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDmsProductV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsProductV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsProductV1DataSourceID("data.g42cloud_dms_product.product1"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_dms_product.product1", "engine", "kafka"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_dms_product.product1", "bandwidth", "100MB"),
				),
			},
		},
	})
}

func TestAccDmsProductV1DataSource_rabbitmqSingle(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsProductV1DataSource_rabbitmqSingle,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsProductV1DataSourceID("data.g42cloud_dms_product.product1"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_dms_product.product1", "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_dms_product.product1", "node_num", "3"),
				),
			},
		},
	})
}

func TestAccDmsProductV1DataSource_rabbitmqCluster(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsProductV1DataSource_rabbitmqCluster,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsProductV1DataSourceID("data.g42cloud_dms_product.product1"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_dms_product.product1", "engine", "rabbitmq"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_dms_product.product1", "node_num", "5"),
				),
			},
		},
	})
}

func testAccCheckDmsProductV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find Dms product data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Dms product data source ID not set")
		}

		return nil
	}
}

var testAccDmsProductV1DataSource_basic = fmt.Sprintf(`
data "g42cloud_dms_product" "product1" {
  engine        = "kafka"
  version       = "1.1.0"
  instance_type = "cluster"
  bandwidth     = "100MB"
}
`)

var testAccDmsProductV1DataSource_rabbitmqSingle = fmt.Sprintf(`
data "g42cloud_dms_product" "product1" {
  engine        = "rabbitmq"
  version       = "3.7.17"
  instance_type = "single"
  node_num      = "3"
}
`)

var testAccDmsProductV1DataSource_rabbitmqCluster = fmt.Sprintf(`
data "g42cloud_dms_product" "product1" {
  engine        = "rabbitmq"
  version       = "3.7.17"
  instance_type = "cluster"
  node_num      = "5"
}
`)
