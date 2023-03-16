package sms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/sms/v3/templates"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getServerTemplateResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.SmsV3Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SMS client: %s", err)
	}

	return templates.Get(client, state.Primary.ID)
}

func TestAccServerTemplate_basic(t *testing.T) {
	var temp templates.TemplateResponse
	name := acceptance.RandomAccResourceName()
	resourceName := "g42cloud_sms_server_template.test"
	azDataName := "data.g42cloud_availability_zones.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&temp,
		getServerTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccServerTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "target_server_name", name),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "vpc_name", "autoCreate"),
					resource.TestCheckResourceAttr(resourceName, "subnet_ids.0", "autoCreate"),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.0", "autoCreate"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone", azDataName, "names.0"),
				),
			},
			{
				Config: testAccServerTemplate_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", name)),
					resource.TestCheckResourceAttr(resourceName, "target_server_name", name),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone", azDataName, "names.1"),
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

func TestAccServerTemplate_existing(t *testing.T) {
	var temp templates.TemplateResponse
	name := acceptance.RandomAccResourceName()
	resourceName := "g42cloud_sms_server_template.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&temp,
		getServerTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccServerTemplate_existing(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "target_server_name", name),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "vpc_name", name),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.0", "autoCreate"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "g42cloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_ids.0", "g42cloud_vpc_subnet.test", "id"),
				),
			},
			{
				Config: testAccServerTemplate_existing_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "g42cloud_vpc.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_ids.0", "g42cloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_ids.0", "g42cloud_networking_secgroup.test", "id"),
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

func testAccServerTemplate_basic(name string) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

resource "g42cloud_sms_server_template" "test" {
  name              = "%s"
  availability_zone = data.g42cloud_availability_zones.test.names[0]
}
`, name)
}

func testAccServerTemplate_update(name string) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

resource "g42cloud_sms_server_template" "test" {
  name               = "%s-update"
  target_server_name = "%s"
  availability_zone  = data.g42cloud_availability_zones.test.names[1]
  volume_type        = "GPSSD"
}
`, name, name)
}

func testAccServerTemplate_existing(name string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_availability_zones" "test" {}

resource "g42cloud_sms_server_template" "test" {
  name              = "%s"
  availability_zone = data.g42cloud_availability_zones.test.names[0]
  vpc_id            = g42cloud_vpc.test.id
  subnet_ids        = [ g42cloud_vpc_subnet.test.id ]
}
`, common.TestBaseNetwork(name), name)
}

func testAccServerTemplate_existing_update(name string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_availability_zones" "test" {}

resource "g42cloud_sms_server_template" "test" {
  name               = "%s"
  availability_zone  = data.g42cloud_availability_zones.test.names[0]
  vpc_id             = g42cloud_vpc.test.id
  subnet_ids         = [ g42cloud_vpc_subnet.test.id ]
  security_group_ids = [ g42cloud_networking_secgroup.test.id ]
}
`, common.TestBaseNetwork(name), name)
}
