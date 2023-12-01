package apig

import (
	"fmt"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance/common"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/signs"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getSignatureAssociateFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.ApigV2Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG v2 client: %s", err)
	}
	opts := signs.ListBindOpts{
		InstanceId:  state.Primary.Attributes["instance_id"],
		SignatureId: state.Primary.Attributes["signature_id"],
		Limit:       500,
	}
	resp, err := signs.ListBind(c, opts)
	if err != nil {
		return nil, err
	}
	if len(resp) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	return resp, nil
}

func TestAccSignatureAssociate_basic(t *testing.T) {
	var (
		apiDetails []signs.SignBindApiInfo

		name   = acceptance.RandomAccResourceName()
		rName1 = "g42cloud_apig_signature_associate.basic_bind"
		rName2 = "g42cloud_apig_signature_associate.hmac_bind"
		rName3 = "g42cloud_apig_signature_associate.aes_bind"

		rc1 = acceptance.InitResourceCheck(rName1, &apiDetails, getSignatureAssociateFunc)
		rc2 = acceptance.InitResourceCheck(rName2, &apiDetails, getSignatureAssociateFunc)
		rc3 = acceptance.InitResourceCheck(rName3, &apiDetails, getSignatureAssociateFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc1.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSignatureAssociate_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName1, "instance_id",
						"g42cloud_apig_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName1, "signature_id",
						"g42cloud_apig_signature.basic", "id"),
					resource.TestCheckResourceAttr(rName1, "publish_ids.#", "2"),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName2, "instance_id",
						"g42cloud_apig_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName2, "signature_id",
						"g42cloud_apig_signature.hmac", "id"),
					resource.TestCheckResourceAttr(rName2, "publish_ids.#", "2"),
					rc3.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName3, "instance_id",
						"g42cloud_apig_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName3, "signature_id",
						"g42cloud_apig_signature.aes", "id"),
					resource.TestCheckResourceAttr(rName3, "publish_ids.#", "2"),
				),
			},
			{
				Config: testAccSignatureAssociate_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "publish_ids.#", "2"),
					rc2.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName2, "publish_ids.#", "2"),
					rc3.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName3, "publish_ids.#", "2"),
				),
			},
			{
				ResourceName:      rName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureAssociateImportStateFunc(rName1),
			},
			{
				ResourceName:      rName2,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureAssociateImportStateFunc(rName2),
			},
			{
				ResourceName:      rName3,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSignatureAssociateImportStateFunc(rName3),
			},
		},
	})
}

func testAccSignatureAssociateImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.Attributes["signature_id"] == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<signature_id>', but got '%s/%s'",
				rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["signature_id"])
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["signature_id"]), nil
	}
}

func testAccSignatureAssociate_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "g42cloud_apig_instance" "test" {
  name                  = "%[2]s"
  edition               = "BASIC"
  vpc_id                = g42cloud_vpc.test.id
  subnet_id             = g42cloud_vpc_subnet.test.id
  security_group_id     = g42cloud_networking_secgroup.test.id
  enterprise_project_id = "0"

  availability_zones = try(slice(data.g42cloud_availability_zones.test.names, 0, 1), null)
}

resource "g42cloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.g42cloud_images_image.test.id
  flavor_id          = "m6.large.8"
  security_group_ids = [g42cloud_networking_secgroup.test.id]
  availability_zone  = data.g42cloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = g42cloud_vpc_subnet.test.id
  }
}

resource "g42cloud_apig_group" "test" {
  name        = "%[2]s"
  instance_id = g42cloud_apig_instance.test.id
}

resource "g42cloud_apig_vpc_channel" "test" {
  name        = "%[2]s"
  instance_id = g42cloud_apig_instance.test.id
  port        = 80
  algorithm   = "WRR"
  protocol    = "HTTP"
  path        = "/"
  http_code   = "201"

  members {
    id = g42cloud_compute_instance.test.id
  }
}

resource "g42cloud_apig_api" "test" {
  instance_id             = g42cloud_apig_instance.test.id
  group_id                = g42cloud_apig_group.test.id
  name                    = "%[2]s"
  type                    = "Public"
  request_protocol        = "HTTP"
  request_method          = "GET"
  request_path            = "/user_info/{user_age}"
  security_authentication = "APP"
  matching                = "Exact"
  success_response        = "Success response"
  failure_response        = "Failed response"
  description             = "Created by script"

  request_params {
    name     = "user_age"
    type     = "NUMBER"
    location = "PATH"
    required = true
    maximum  = 200
    minimum  = 0
  }
  
  backend_params {
    type     = "REQUEST"
    name     = "userAge"
    location = "PATH"
    value    = "user_age"
  }

  web {
    path             = "/getUserAge/{userAge}"
    vpc_channel_id   = g42cloud_apig_vpc_channel.test.id
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 30000
  }

  web_policy {
    name             = "%[2]s_policy1"
    request_protocol = "HTTP"
    request_method   = "GET"
    effective_mode   = "ANY"
    path             = "/getUserAge/{userAge}"
    timeout          = 30000
    vpc_channel_id   = g42cloud_apig_vpc_channel.test.id

    backend_params {
      type     = "REQUEST"
      name     = "userAge"
      location = "PATH"
      value    = "user_age"
    }

    conditions {
      source     = "param"
      param_name = "user_age"
      type       = "Equal"
      value      = "28"
    }
  }
}

resource "g42cloud_apig_environment" "test" {
  count = 6

  name        = "%[2]s_${count.index}"
  instance_id = g42cloud_apig_instance.test.id
}

resource "g42cloud_apig_api_publishment" "test" {
  count = 6

  instance_id = g42cloud_apig_instance.test.id
  api_id      = g42cloud_apig_api.test.id
  env_id      = g42cloud_apig_environment.test[count.index].id
}

resource "g42cloud_apig_signature" "basic" {
  instance_id = g42cloud_apig_instance.test.id
  name        = "%[2]s_basic"
  type        = "basic"
}

resource "g42cloud_apig_signature" "hmac" {
  instance_id = g42cloud_apig_instance.test.id
  name        = "%[2]s_hmac"
  type        = "hmac"
}

resource "g42cloud_apig_signature" "aes" {
  instance_id = g42cloud_apig_instance.test.id
  name        = "%[2]s_aes"
  type        = "aes"
  algorithm   = "aes-128-cfb"
}
`, common.TestBaseComputeResources(name), name)
}

func testAccSignatureAssociate_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "g42cloud_apig_signature_associate" "basic_bind" {
  instance_id  = g42cloud_apig_instance.test.id
  signature_id = g42cloud_apig_signature.basic.id
  publish_ids  = slice(g42cloud_apig_api_publishment.test[*].publish_id, 0, 2)
}

resource "g42cloud_apig_signature_associate" "hmac_bind" {
  instance_id  = g42cloud_apig_instance.test.id
  signature_id = g42cloud_apig_signature.hmac.id
  publish_ids  = slice(g42cloud_apig_api_publishment.test[*].publish_id, 2, 4)
}

resource "g42cloud_apig_signature_associate" "aes_bind" {
  instance_id  = g42cloud_apig_instance.test.id
  signature_id = g42cloud_apig_signature.aes.id
  publish_ids  = slice(g42cloud_apig_api_publishment.test[*].publish_id, 4, 6)
}
`, testAccSignatureAssociate_base(name))
}

func testAccSignatureAssociate_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "g42cloud_apig_signature_associate" "basic_bind" {
  instance_id  = g42cloud_apig_instance.test.id
  signature_id = g42cloud_apig_signature.basic.id
  publish_ids  = slice(g42cloud_apig_api_publishment.test[*].publish_id, 1, 3)
}

resource "g42cloud_apig_signature_associate" "hmac_bind" {
  instance_id  = g42cloud_apig_instance.test.id
  signature_id = g42cloud_apig_signature.hmac.id
  publish_ids  = slice(g42cloud_apig_api_publishment.test[*].publish_id, 3, 5)
}

resource "g42cloud_apig_signature_associate" "aes_bind" {
  instance_id  = g42cloud_apig_instance.test.id
  signature_id = g42cloud_apig_signature.aes.id
  publish_ids  = setunion(slice(g42cloud_apig_api_publishment.test[*].publish_id, 0, 1),
    slice(g42cloud_apig_api_publishment.test[*].publish_id, 5, 6))
}
`, testAccSignatureAssociate_base(name))
}
