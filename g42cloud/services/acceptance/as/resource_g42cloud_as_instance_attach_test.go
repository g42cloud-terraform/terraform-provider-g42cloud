package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/instances"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getASInstanceAttachResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.G42_REGION_NAME
	client, err := cfg.AutoscalingV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating autoscaling client: %s", err)
	}

	groupID := state.Primary.Attributes["scaling_group_id"]
	instanceID := state.Primary.Attributes["instance_id"]
	page, err := instances.List(client, groupID, nil).AllPages()
	if err != nil {
		return nil, err
	}

	allInstances, err := page.(instances.InstancePage).Extract()
	if err != nil {
		return nil, fmt.Errorf("failed to fetching instances in AS group %s: %s", groupID, err)
	}

	for _, ins := range allInstances {
		if ins.ID == instanceID {
			return &ins, nil
		}
	}

	return nil, fmt.Errorf("can not find the instance %s in AS group %s", instanceID, groupID)
}

func TestAccASInstanceAttach_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "g42cloud_as_instance_attach.test0"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getASInstanceAttachResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testASInstanceAttach_conf(name, "false", "false"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "scaling_group_id", "g42cloud_as_group.hth_as_group", "id"),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "g42cloud_compute_instance.test.0", "id"),
					resource.TestCheckResourceAttr(rName, "protected", "false"),
					resource.TestCheckResourceAttr(rName, "standby", "false"),
					resource.TestCheckResourceAttr(rName, "status", "INSERVICE"),
				),
			},
			{
				Config: testASInstanceAttach_conf(name, "true", "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "protected", "true"),
					resource.TestCheckResourceAttr(rName, "standby", "false"),
					resource.TestCheckResourceAttr(rName, "status", "INSERVICE"),
				),
			},
			{
				Config: testASInstanceAttach_conf(name, "true", "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "protected", "true"),
					resource.TestCheckResourceAttr(rName, "standby", "true"),
					resource.TestCheckResourceAttr(rName, "status", "STANDBY"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"append_instance"},
			},
		},
	})
}

func testAsGroup_base(rName string) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

data "g42cloud_vpc" "test" {
  name = "vpc-default"
}

data "g42cloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "g42cloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "g42cloud_compute_flavors" "test" {
  availability_zone = data.g42cloud_availability_zones.test.names[0]
  performance_type  = "normal"
}

resource "g42cloud_networking_secgroup" "secgroup" {
  name        = "%[1]s"
  description = "This is a terraform test security group"
}

resource "g42cloud_compute_keypair" "hth_key" {
  name       = "%[1]s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "g42cloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "%[1]s"
  vip_subnet_id = data.g42cloud_vpc_subnet.test.subnet_id
}

resource "g42cloud_lb_listener" "listener_1" {
  name            = "%[1]s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = g42cloud_lb_loadbalancer.loadbalancer_1.id
}

resource "g42cloud_lb_pool" "pool_1" {
  name        = "%[1]s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = g42cloud_lb_listener.listener_1.id
}

resource "g42cloud_as_configuration" "hth_as_config"{
  scaling_configuration_name = "%[1]s"
  instance_config {
    image    = data.g42cloud_images_image.test.id
    flavor   = data.g42cloud_compute_flavors.test.ids[0]
    key_name = g42cloud_compute_keypair.hth_key.id

    disk {
      size        = 40
      volume_type = "SAS"
      disk_type   = "SYS"
    }
  }
}

resource "g42cloud_as_group" "hth_as_group"{
  scaling_group_name       = "%[1]s"
  scaling_configuration_id = g42cloud_as_configuration.hth_as_config.id
  vpc_id                   = data.g42cloud_vpc.test.id
  max_instance_number      = 3

  networks {
    id = data.g42cloud_vpc_subnet.test.id
  }
  security_groups {
    id = g42cloud_networking_secgroup.secgroup.id
  }
  lbaas_listeners {
    pool_id       = g42cloud_lb_pool.pool_1.id
    protocol_port = g42cloud_lb_listener.listener_1.protocol_port
  }
}
`, rName)
}

func testASInstanceAttach_conf(name, protection, standby string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_compute_instance" "test" {
  count              = 2
  name               = "%s-${count.index}"
  description        = "instance for AS attach"
  image_id           = data.g42cloud_images_image.test.id
  flavor_id          = data.g42cloud_compute_flavors.test.ids[0]
  security_group_ids = [g42cloud_networking_secgroup.secgroup.id]

  network {
    uuid = data.g42cloud_vpc_subnet.test.id
  }
}

resource "g42cloud_as_instance_attach" "test0" {
  scaling_group_id = g42cloud_as_group.hth_as_group.id
  instance_id      = g42cloud_compute_instance.test[0].id
  protected        = %[3]s
  standby          = %[4]s
}
`, testAsGroup_base(name), name, protection, standby)
}
