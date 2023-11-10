package acceptance

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/cce/v1/namespaces"
)

func getNamespaceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CceV1Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating G42Cloud CCE v1 client: %s", err)
	}
	resp, err := namespaces.Get(c, state.Primary.Attributes["cluster_id"],
		state.Primary.Attributes["name"]).Extract()
	if resp == nil && err == nil {
		return resp, fmt.Errorf("Unable to find the namespace (%s)", state.Primary.ID)
	}
	return resp, err
}

func TestAccCCENamespaceV1_basic(t *testing.T) {
	var namespace namespaces.Namespace
	resourceName := "g42cloud_cce_namespace.test"
	randName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&namespace,
		getNamespaceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCCENamespaceV1_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "cluster_id",
						"${g42cloud_cce_cluster.test.id}"),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCCENamespaceImportStateIdFunc(randName),
			},
		},
	})
}

func TestAccCCENamespaceV1_generateName(t *testing.T) {
	var namespace namespaces.Namespace
	resourceName := "g42cloud_cce_namespace.test"
	randName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	rc := acceptance.InitResourceCheck(
		resourceName,
		&namespace,
		getNamespaceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCCENamespaceV1_generateName(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "cluster_id",
						"${g42cloud_cce_cluster.test.id}"),
					resource.TestCheckResourceAttr(resourceName, "prefix", randName),
					resource.TestCheckResourceAttr(resourceName, "status", "Active"),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(fmt.Sprintf(`^%s[a-z0-9-]*`, randName))),
				),
			},
		},
	})
}

func testAccCCENamespaceImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var clusterID string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "g42cloud_cce_cluster" {
				clusterID = rs.Primary.ID
			}
		}
		if clusterID == "" || name == "" {
			return "", fmtp.Errorf("resource not found: %s/%s", clusterID, name)
		}
		return fmt.Sprintf("%s/%s", clusterID, name), nil

	}
}

func testAccCCEClusterV3_base(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "g42cloud_vpc_subnet" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"

  //dns is required for cce node installing
  primary_dns   = "100.125.3.250"
  secondary_dns = "100.125.3.92"
  vpc_id        = g42cloud_vpc.test.id
}

resource "g42cloud_cce_cluster" "test" {
  name                   = "%[1]s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = g42cloud_vpc.test.id
  subnet_id              = g42cloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
}

`, rName)
}

func testAccCCENodeV3_base(rName string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_availability_zones" "test" {}

resource "g42cloud_compute_keypair" "test" {
  name       = "%s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "g42cloud_cce_node" "test" {
  cluster_id        = g42cloud_cce_cluster.test.id
  name              = "%s"
  flavor_id         = "c7n.large.2"
  availability_zone = data.g42cloud_availability_zones.test.names[0]
  key_pair          = g42cloud_compute_keypair.test.name
  os                = "CentOS 7.6"

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCCEClusterV3_base(rName), rName, rName)
}

func testAccCCENamespaceV1_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_cce_namespace" "test" {
  cluster_id = g42cloud_cce_cluster.test.id
  name       = "%s"
}
`, testAccCCEClusterV3_base(rName), rName)
}

func testAccCCENamespaceV1_generateName(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_cce_namespace" "test" {
  cluster_id = g42cloud_cce_cluster.test.id
  prefix     = "%s"
}
`, testAccCCEClusterV3_base(rName), rName)
}
