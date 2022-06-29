package servicestage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/servicestage/v2/instances"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getComponentInstanceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.ServiceStageV2Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ServiceStage v2 client: %s", err)
	}
	return instances.Get(c, state.Primary.Attributes["application_id"], state.Primary.Attributes["component_id"],
		state.Primary.ID)
}

func TestAccComponentInstance_basic(t *testing.T) {
	var (
		instance     instances.Instance
		randName     = acceptance.RandomAccResourceNameWithDash()
		resourceName = "g42cloud_servicestage_component_instance.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getComponentInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRepoTokenAuth(t)
			acceptance.TestAccPreCheckComponent(t)
			acceptance.TestAccPreCheckComponentDeployment(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComponentInstance_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "application_id", "g42cloud_servicestage_application.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "component_id", "g42cloud_servicestage_component.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "environment_id", "g42cloud_servicestage_environment.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.0"),
					resource.TestCheckResourceAttr(resourceName, "replica", "1"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "CUSTOM-10G:250m-250m:0.5Gi-0.5Gi"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform test"),
					resource.TestCheckResourceAttr(resourceName, "artifact.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "artifact.0.name", "g42cloud_servicestage_component.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "artifact.0.type", "image"),
					resource.TestCheckResourceAttr(resourceName, "artifact.0.storage", "swr"),
					resource.TestCheckResourceAttr(resourceName, "artifact.0.url", acceptance.G42_BUILD_IMAGE_URL),
					resource.TestCheckResourceAttr(resourceName, "artifact.0.auth_type", "iam"),
					resource.TestCheckResourceAttr(resourceName, "refer_resource.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.env_variable.0.name", "TZ"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.env_variable.0.value", "Asia/Shanghai"),
					resource.TestCheckResourceAttr(resourceName, "status", "RUNNING"),
				),
			},
			{
				Config: testAccComponentInstance_update(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0.2"),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "CUSTOM-15G:500m-500m:1Gi-1Gi"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "status", "RUNNING"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccInstanceImportStateIdFunc(),
			},
		},
	})
}

func testAccInstanceImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var appId, componentId, instance_id string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "g42cloud_servicestage_component_instance" {
				appId = rs.Primary.Attributes["application_id"]
				componentId = rs.Primary.Attributes["component_id"]
				instance_id = rs.Primary.ID
			}
		}
		if appId == "" || componentId == "" || instance_id == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s", appId, componentId, instance_id)
		}
		return fmt.Sprintf("%s/%s/%s", appId, componentId, instance_id), nil
	}
}

func testAccComponentInstance_base(rName string) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

data "g42cloud_compute_flavors" "test" {
  availability_zone = data.g42cloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 8
  memory_size       = 16
}

data "g42cloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "g42cloud_compute_keypair" "test" {
  name = "%[1]s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "g42cloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "g42cloud_vpc_subnet" "test" {
  name        = "%[1]s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  vpc_id      = g42cloud_vpc.test.id
  ipv6_enable = true
}

resource "g42cloud_networking_secgroup" "test" {
  name = "%[1]s"
}

resource "g42cloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    size        = 5
    name        = "%[1]s"
    charge_mode = "traffic"
  }
}

resource "g42cloud_as_configuration" "test" {
  scaling_configuration_name = "%[1]s"

  instance_config {
    flavor   = data.g42cloud_compute_flavors.test.ids[0]
    image    = data.g42cloud_images_image.test.id
    key_name = g42cloud_compute_keypair.test.name

    disk {
      disk_type   = "SYS"
      volume_type = "SSD"
      size        = 40
    }
  }
}

resource "g42cloud_as_group" "test" {
  scaling_group_name       = "%[1]s"
  scaling_configuration_id = g42cloud_as_configuration.test.id
  vpc_id                   = g42cloud_vpc.test.id

  max_instance_number    = 3
  min_instance_number    = 0
  desire_instance_number = 1

  delete_instances = "yes"
  delete_publicip  = false

  cool_down_time = 86400

  networks {
    id = g42cloud_vpc_subnet.test.id
  }

  security_groups {
    id = g42cloud_networking_secgroup.test.id
  }
}

resource "g42cloud_servicestage_environment" "test" {
  name        = "%[1]s"
  description = "Created by terraform test"
  vpc_id      = g42cloud_vpc.test.id

  basic_resources {
    type = "as"
    id   = g42cloud_as_group.test.id
  }

  optional_resources {
    type = "eip"
    id   = g42cloud_vpc_eip.test.id
  }
}

resource "g42cloud_servicestage_application" "test" {
  name = "%[1]s"
}

resource "g42cloud_servicestage_repo_token_authorization" "test" {
  type  = "github"
  name  = "%[1]s"
  host  = "%[2]s"
  token = "%[3]s"
}

resource "g42cloud_servicestage_component" "test" {
  application_id = g42cloud_servicestage_application.test.id

  name      = "%[1]s"
  type      = "MicroService"
  runtime   = "Docker"
  framework = "Java Classis"
}
`, rName, acceptance.G42_GITHUB_REPO_HOST, acceptance.G42_GITHUB_PERSONAL_TOKEN, acceptance.G42_GITHUB_REPO_URL,
		acceptance.G42_ACCOUNT_NAME)
}

func testAccComponentInstance_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "g42cloud_servicestage_component_instance" "test" {
  application_id = g42cloud_servicestage_application.test.id
  component_id   = g42cloud_servicestage_component.test.id
  environment_id = g42cloud_servicestage_environment.test.id

  name        = "%[2]s"
  version     = "1.0.0"
  replica     = 1
  flavor_id   = "CUSTOM-10G:250m-250m:0.5Gi-0.5Gi"
  description = "Created by terraform test"

  artifact {
    name      = g42cloud_servicestage_component.test.name
    type      = "image"
    storage   = "swr"
    url       = "%[3]s"
    auth_type = "iam"
  }

  refer_resource {
    type = "ecs"
    id   = "default"
  }

  configuration {
    env_variable {
      name  = "TZ"
      value = "Asia/Shanghai"
    }
  }

  lifecycle {
    ignore_changes = [
      configuration[0].env_variable,
    ]
  }
}
`, testAccComponentInstance_base(rName), rName, acceptance.G42_BUILD_IMAGE_URL)
}

func testAccComponentInstance_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "g42cloud_servicestage_component_instance" "test" {
  application_id = g42cloud_servicestage_application.test.id
  component_id   = g42cloud_servicestage_component.test.id
  environment_id = g42cloud_servicestage_environment.test.id

  name        = "%[2]s"
  version     = "1.0.2"
  replica     = 1
  flavor_id   = "CUSTOM-15G:500m-500m:1Gi-1Gi"

  artifact {
    name      = g42cloud_servicestage_component.test.name
    type      = "image"
    storage   = "swr"
    url       = "%[3]s"
    auth_type = "iam"
  }

  refer_resource {
    type = "ecs"
    id   = "default"
  }

  configuration {
    env_variable {
      name  = "TZ"
      value = "Asia/Shanghai"
    }
  }

  lifecycle {
    ignore_changes = [
      configuration[0].env_variable,
    ]
  }
}
`, testAccComponentInstance_base(rName), rName, acceptance.G42_BUILD_IMAGE_URL)
}
