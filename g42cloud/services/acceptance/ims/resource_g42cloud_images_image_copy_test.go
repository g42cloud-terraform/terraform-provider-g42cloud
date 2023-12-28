package ims

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ims"
)

func getImsImageCopyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.G42_REGION_NAME

	targetRegion := state.Primary.Attributes["target_region"]
	if targetRegion != "" {
		region = targetRegion
	}

	imsClient, err := cfg.ImageV2Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating IMS client: %s", err)
	}

	img, err := ims.GetCloudImage(imsClient, state.Primary.ID)
	if err != nil {
		return nil, fmt.Errorf("image %s not found: %s", state.Primary.ID, err)
	}
	return img, nil
}

func TestAccImsImageCopy_basic(t *testing.T) {
	var obj interface{}

	sourceImageName := acceptance.RandomAccResourceName()
	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	rName := "g42cloud_images_image_copy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getImsImageCopyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testImsImageCopy_basic(sourceImageName, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "min_ram", "2048"),
					resource.TestCheckResourceAttr(rName, "max_ram", "4096"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(rName, "tags.key2", "value2"),
				),
			},
			{
				Config: testImsImageCopy_update(sourceImageName, updateName, 4096, 8192),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "it's a test"),
					resource.TestCheckResourceAttr(rName, "min_ram", "4096"),
					resource.TestCheckResourceAttr(rName, "max_ram", "8192"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1_update"),
					resource.TestCheckResourceAttr(rName, "tags.key3", "value3"),
					resource.TestCheckResourceAttr(rName, "tags.key4", "value4"),
				),
			},
			{
				Config: testImsImageCopy_update(sourceImageName, updateName, 0, 0),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "it's a test"),
					resource.TestCheckResourceAttr(rName, "min_ram", "0"),
					resource.TestCheckResourceAttr(rName, "max_ram", "0"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1_update"),
					resource.TestCheckResourceAttr(rName, "tags.key3", "value3"),
					resource.TestCheckResourceAttr(rName, "tags.key4", "value4"),
				),
			},
		},
	})
}

func TestAccImsImageCopy_basic_cross_region(t *testing.T) {
	var obj interface{}

	sourceImageName := acceptance.RandomAccResourceName()
	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	rName := "g42cloud_images_image_copy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getImsImageCopyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testImsImageCopy_basic_cross_region(sourceImageName, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "it's a test"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(rName, "tags.key2", "value2"),
				),
			},
			{
				Config: testImsImageCopy_update_cross_region(sourceImageName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "it's a test"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1_update"),
					resource.TestCheckResourceAttr(rName, "tags.key3", "value3"),
					resource.TestCheckResourceAttr(rName, "tags.key4", "value4"),
				),
			},
		},
	})
}

func testImsImageCopy_basic(baseImageName, copyImageName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_images_image_copy" "test" {
  source_image_id = g42cloud_images_image.test.id
  name            = "%s"
  description     = "terraform copied image"
  min_ram         = 2048
  max_ram         = 4096

  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}
`, testAccImsImage_basic(baseImageName), copyImageName)
}

func testImsImageCopy_update(baseImageName, copyImageName string, minRAM, maxRAM int) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_images_image_copy" "test" {
  source_image_id = g42cloud_images_image.test.id
  name            = "%[2]s"
  description     = "it's a test"
  min_ram         = %[3]d
  max_ram         = %[4]d

  tags = {
    key1 = "value1_update"
    key3 = "value3"
    key4 = "value4"
  }
}
`, testAccImsImage_basic(baseImageName), copyImageName, minRAM, maxRAM)
}

func testImsImageCopy_basic_cross_region(baseImageName, copyImageName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_images_image_copy" "test" {
  source_image_id = g42cloud_images_image.test.id
  name            = "%s"
  description     = "it's a test"
  target_region   = "%s"
  agency_name     = "ims_admin_agency"

  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}`, testAccImsImage_basic(baseImageName), copyImageName, acceptance.G42_DEST_REGION)
}

func testImsImageCopy_update_cross_region(baseImageName, copyImageName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_images_image_copy" "test" {
  source_image_id = g42cloud_images_image.test.id
  name            = "%s"
  description     = "it's a test"
  target_region   = "%s"
  agency_name     = "ims_admin_agency"

  tags = {
    key1 = "value1_update"
    key3 = "value3"
    key4 = "value4"
  }
}
`, testAccImsImage_basic(baseImageName), copyImageName, acceptance.G42_DEST_REGION)
}
