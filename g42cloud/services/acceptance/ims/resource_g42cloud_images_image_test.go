package ims

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/ims/v2/cloudimages"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ims"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccImsImage_basic(t *testing.T) {
	var image cloudimages.Image

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameUpdate := rName + "-update"
	resourceName := "g42cloud_images_image.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckImsImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImage_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageExists(resourceName, &image),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccImsImage_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageExists(resourceName, &image),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
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

func TestAccImsImage_withEpsId(t *testing.T) {
	var image cloudimages.Image

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "g42cloud_images_image.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckImsImageDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImsImage_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImsImageExists(resourceName, &image),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.G42_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckImsImageDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	imageClient, err := config.ImageV2Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating G42Cloud Image: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "g42cloud_images_image" {
			continue
		}

		_, err := ims.GetCloudImage(imageClient, rs.Primary.ID)
		if err == nil {
			return fmtp.Errorf("Image still exists")
		}
	}

	return nil
}

func testAccCheckImsImageExists(n string, image *cloudimages.Image) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("IMS Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		imageClient, err := config.ImageV2Client(acceptance.G42_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating G42Cloud Image: %s", err)
		}

		found, err := ims.GetCloudImage(imageClient, rs.Primary.ID)
		if err != nil {
			return err
		}

		*image = *found
		return nil
	}
}

func testAccImsImage_basic(rName string) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

data "g42cloud_compute_flavors" "test" {
  availability_zone = data.g42cloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 1
  memory_size       = 2
}

data "g42cloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "g42cloud_networking_secgroup" "test" {
  name = "default"
}

resource "g42cloud_compute_instance" "test" {
  name               = "%s"
  image_name         = "Ubuntu 18.04 server 64bit"
  flavor_id          = data.g42cloud_compute_flavors.test.ids[0]
  security_group_ids = [data.g42cloud_networking_secgroup.test.id]
  availability_zone  = data.g42cloud_availability_zones.test.names[0]

  network {
    uuid = data.g42cloud_vpc_subnet.test.id
  }
}

resource "g42cloud_images_image" "test" {
  name        = "%s"
  instance_id = g42cloud_compute_instance.test.id
  description = "created by Terraform AccTest"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName, rName)
}

func testAccImsImage_update(rName string) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

data "g42cloud_compute_flavors" "test" {
  availability_zone = data.g42cloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 1
  memory_size       = 2
}

data "g42cloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "g42cloud_networking_secgroup" "test" {
  name = "default"
}

resource "g42cloud_compute_instance" "test" {
  name               = "%s"
  image_name         = "Ubuntu 18.04 server 64bit"
  flavor_id          = data.g42cloud_compute_flavors.test.ids[0]
  security_group_ids = [data.g42cloud_networking_secgroup.test.id]
  availability_zone  = data.g42cloud_availability_zones.test.names[0]

  network {
    uuid = data.g42cloud_vpc_subnet.test.id
  }
}

resource "g42cloud_images_image" "test" {
  name        = "%s"
  instance_id = g42cloud_compute_instance.test.id
  description = "created by Terraform AccTest"

  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, rName, rName)
}

func testAccImsImage_withEpsId(rName string) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

data "g42cloud_compute_flavors" "test" {
  availability_zone = data.g42cloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 1
  memory_size       = 2
}

data "g42cloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "g42cloud_networking_secgroup" "test" {
  name = "default"
}

resource "g42cloud_compute_instance" "test" {
  name               = "%s"
  image_name         = "Ubuntu 18.04 server 64bit"
  flavor_id          = data.g42cloud_compute_flavors.test.ids[0]
  security_group_ids = [data.g42cloud_networking_secgroup.test.id]
  availability_zone  = data.g42cloud_availability_zones.test.names[0]

  network {
    uuid = data.g42cloud_vpc_subnet.test.id
  }
}

resource "g42cloud_images_image" "test" {
  name                  = "%s"
  instance_id           = g42cloud_compute_instance.test.id
  description           = "created by Terraform AccTest"
  enterprise_project_id = "%s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName, rName, acceptance.G42_ENTERPRISE_PROJECT_ID_TEST)
}
