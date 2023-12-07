package mrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccDatasourceMrsClusters_basic(t *testing.T) {
	rName := "data.g42cloud_mapreduce_clusters.clusters_1"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()
	pwd := acceptance.RandomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceMrsClusters_basic(name, pwd),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "clusters.0.name", "g42cloud_mapreduce_cluster.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "clusters.0.enterprise_project_id",
						"g42cloud_mapreduce_cluster.test", "enterprise_project_id"),
					resource.TestCheckResourceAttrPair(rName, "clusters.0.version",
						"g42cloud_mapreduce_cluster.test", "version"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.id"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.vpc_id"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.subnet_id"),
					resource.TestCheckResourceAttrSet(rName, "clusters.0.vnc"),
				),
			},
		},
	})
}

func testAccMrsCluster_base(name, pwd string) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

data "g42cloud_vpc" "test" {
  name = "vpc-default"
}

data "g42cloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "g42cloud_mapreduce_cluster" "test" {
  availability_zone  = data.g42cloud_availability_zones.test.names[0]
  name               = "%[1]s"
  type               = "ANALYSIS"
  version            = "MRS 3.2.0.1"
  manager_admin_pass = "%[2]s"
  node_admin_pass    = "%[2]s"
  subnet_id          = data.g42cloud_vpc_subnet.test.id
  vpc_id             = data.g42cloud_vpc.test.id
  component_list     = ["Hadoop", "ZooKeeper", "Ranger"]

  master_nodes {
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = 2
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
  }
  analysis_core_nodes {
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
  }
  analysis_task_nodes {
    flavor            = "c6.4xlarge.4.linux.bigdata"
    node_number       = 3
    root_volume_type  = "SAS"
    root_volume_size  = 480
    data_volume_type  = "SAS"
    data_volume_size  = 600
    data_volume_count = 1
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, pwd)
}

func testAccDatasourceMrsClusters_basic(name, pwd string) string {
	return fmt.Sprintf(`
%[1]s

data "g42cloud_mapreduce_clusters" "clusters_1" {
  status = "running"

  depends_on = [
    g42cloud_mapreduce_cluster.test
  ]
}
`, testAccMrsCluster_base(name, pwd), name)
}
