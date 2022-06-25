package mrs

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/mrs/v1/cluster"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

type GroupNodeNum struct {
	AnalysisCoreNum int
	StreamCoreNum   int
	AnalysisTaskNum int
	StreamTaskNum   int
}

func TestAccMrsMapReduceCluster_basic(t *testing.T) {
	var clusterGet cluster.Cluster
	resourceName := "g42cloud_mapreduce_cluster.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	password := fmt.Sprintf("TF%s%s%d", acctest.RandString(10), acctest.RandStringFromCharSet(1, "-_"),
		acctest.RandIntRange(0, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckMRSV2ClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMrsMapReduceClusterConfig_basic(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "STREAMING"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccMrsMapReduceClusterConfig_update(rName, password),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMRSV2ClusterExists(resourceName, &clusterGet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "STREAMING"),
					resource.TestCheckResourceAttr(resourceName, "safe_mode", "true"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo1", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "update_value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"manager_admin_pass",
					"node_admin_pass",
				},
			},
		},
	})
}

func testAccCheckMRSV2ClusterDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := config.MrsV1Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating g42cloud mrs: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "g42cloud_mapreduce_cluster" {
			continue
		}

		clusterGet, err := cluster.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return fmt.Errorf("MRS cluster (%s) is still exists", rs.Primary.ID)
		}
		if clusterGet.Clusterstate == "terminated" {
			return nil
		}
	}

	return nil
}

func testAccCheckMRSV2ClusterExists(n string, clusterGet *cluster.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource %s not found", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No MRS cluster ID")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		mrsClient, err := config.MrsV1Client(acceptance.G42_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating g42cloud MRS client: %s ", err)
		}

		found, err := cluster.Get(mrsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*clusterGet = *found
		return nil
	}
}

func testAccMrsMapReduceClusterConfig_base(rName string) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

resource "g42cloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "g42cloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/20"
  vpc_id     = g42cloud_vpc.test.id
  gateway_ip = "192.168.0.1"
}
`, rName, rName)
}

// The task node has not contain data disks.
func testAccMrsMapReduceClusterConfig_basic(rName, pwd string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_mapreduce_cluster" "test" {
  availability_zone  = data.g42cloud_availability_zones.test.names[0]
  name               = "%s"
  type               = "STREAMING"
  version            = "MRS 1.9.2"
  manager_admin_pass = "%s"
  node_admin_pass    = "%s"
  subnet_id          = g42cloud_vpc_subnet.test.id
  vpc_id             = g42cloud_vpc.test.id
  component_list     = ["Storm"]

  master_nodes {
    flavor            = "d3.2xlarge.8.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_core_nodes {
    flavor            = "d3.2xlarge.8.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_task_nodes {
    flavor            = "d3.2xlarge.8.linux.bigdata"
    node_number       = 1
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_count = 0
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, pwd)
}

func testAccMrsMapReduceClusterConfig_update(rName, pwd string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_mapreduce_cluster" "test" {
  availability_zone  = data.g42cloud_availability_zones.test.names[0]
  name               = "%s"
  type               = "STREAMING"
  version            = "MRS 1.9.2"
  manager_admin_pass = "%s"
  node_admin_pass    = "%s"
  subnet_id          = g42cloud_vpc_subnet.test.id
  vpc_id             = g42cloud_vpc.test.id
  component_list     = ["Storm"]

  master_nodes {
    flavor            = "d3.2xlarge.8.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_core_nodes {
    flavor            = "d3.2xlarge.8.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_type  = "SAS"
    data_volume_size  = 480
    data_volume_count = 1
  }
  streaming_task_nodes {
    flavor            = "d3.2xlarge.8.linux.bigdata"
    node_number       = 1
    root_volume_type  = "SAS"
    root_volume_size  = 300
    data_volume_count = 0
  }

  tags = {
    foo1 = "bar"
    key  = "update_value"
  }
}`, testAccMrsMapReduceClusterConfig_base(rName), rName, pwd, pwd)
}
