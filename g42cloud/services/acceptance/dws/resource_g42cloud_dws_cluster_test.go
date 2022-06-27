package dws

import (
	"fmt"
	"testing"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/dws/v1/cluster"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDwsResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.DwsV1Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("error creating DWS v1 client, err=%s", err)
	}
	return cluster.Get(client, state.Primary.ID)
}

func TestAccResourceDWS_basic(t *testing.T) {
	var clusterInstance cluster.CreateOpts
	resourceName := "g42cloud_dws_cluster.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&clusterInstance,
		getDwsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDwsCluster_basic(name, 3, cluster.PublicBindTypeAuto, "cluster123@!"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "number_of_node", "3"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"user_pwd", "number_of_cn"},
			},
		},
	})
}

func testAccBaseResource(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "g42cloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/24"
  vpc_id     = g42cloud_vpc.test.id
  gateway_ip = "192.168.0.1"
}

resource "g42cloud_networking_secgroup" "test" {
  name        = "%s"
  description = "terraform security group acceptance test"
}

data "g42cloud_availability_zones" "test" {}
`, rName, rName, rName)
}

func testAccDwsCluster_basic(rName string, numberOfNode int, publicIpBindType string, password string) string {
	baseResource := testAccBaseResource(rName)

	return fmt.Sprintf(`
%s

resource "g42cloud_dws_cluster" "test" {
  name              = "%s"
  node_type         = "dws2.olap.4xlarge.i3"
  number_of_node    = %d
  vpc_id            = g42cloud_vpc.test.id
  network_id        = g42cloud_vpc_subnet.test.id
  security_group_id = g42cloud_networking_secgroup.test.id
  availability_zone = data.g42cloud_availability_zones.test.names[1]
  user_name         = "test_cluster_admin"
  user_pwd          = "%s"

  public_ip {
    public_bind_type = "%s"
  }
}
`, baseResource, rName, numberOfNode, password, publicIpBindType)
}
