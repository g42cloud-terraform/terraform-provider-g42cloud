package g42cloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccG42CloudImagesV2ImageDataSource_basic(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccG42CloudImagesV2ImageDataSource_ubuntu(rName),
			},
			{
				Config: testAccG42CloudImagesV2ImageDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.g42cloud_images_image.test"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_images_image.test", "name", rName),
					resource.TestCheckResourceAttr(
						"data.g42cloud_images_image.test", "protected", "false"),
					resource.TestCheckResourceAttr(
						"data.g42cloud_images_image.test", "visibility", "private"),
				),
			},
		},
	})
}

func TestAccG42CloudImagesV2ImageDataSource_testQueries(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccG42CloudImagesV2ImageDataSource_ubuntu(rName),
			},
			{
				Config: testAccG42CloudImagesV2ImageDataSource_querySizeMin(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.g42cloud_images_image.test"),
				),
			},
			{
				Config: testAccG42CloudImagesV2ImageDataSource_querySizeMax(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesV2DataSourceID("data.g42cloud_images_image.test"),
				),
			},
			{
				Config: testAccG42CloudImagesV2ImageDataSource_ubuntu(rName),
			},
		},
	})
}

func testAccCheckImagesV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find image data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Image data source ID not set")
		}

		return nil
	}
}

func testAccG42CloudImagesV2ImageDataSource_ubuntu(rName string) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

data "g42cloud_compute_flavors" "test" {
  availability_zone = data.g42cloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "g42cloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "g42cloud_compute_instance" "test" {
  name              = "%s"
  image_name        = "Ubuntu 18.04 server 64bit"
  flavor_id         = data.g42cloud_compute_flavors.test.ids[0]
  security_groups   = ["default"]
  availability_zone = data.g42cloud_availability_zones.test.names[0]
  system_disk_type  = "SSD"

  network {
    uuid = data.g42cloud_vpc_subnet.test.id
  }
}

resource "g42cloud_images_image" "test" {
  name        = "%s"
  instance_id = g42cloud_compute_instance.test.id
  description = "created by TerraformAccTest"
}

`, rName, rName)
}

func testAccG42CloudImagesV2ImageDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_images_image" "test" {
	most_recent = true
	name = g42cloud_images_image.test.name
}
`, testAccG42CloudImagesV2ImageDataSource_ubuntu(rName))
}

func testAccG42CloudImagesV2ImageDataSource_querySizeMin(rName string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_images_image" "test" {
	most_recent = true
	visibility = "private"
	size_min = "13000000"
}
`, testAccG42CloudImagesV2ImageDataSource_ubuntu(rName))
}

func testAccG42CloudImagesV2ImageDataSource_querySizeMax(rName string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_images_image" "test" {
	most_recent = true
	visibility = "private"
	size_max = "23000000"
}
`, testAccG42CloudImagesV2ImageDataSource_ubuntu(rName))
}
