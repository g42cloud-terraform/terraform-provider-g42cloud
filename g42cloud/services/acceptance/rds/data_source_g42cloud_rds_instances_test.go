package rds

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccRdsInstanceDataSource_basic(t *testing.T) {
	dataSourceName := "data.g42cloud_rds_instances.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	dc := acceptance.InitDataSourceCheck(dataSourceName)
	dbPwd := fmt.Sprintf("%s%s%d", acctest.RandString(5),
		acctest.RandStringFromCharSet(2, "!#%^*"), acctest.RandIntRange(10, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceDataSource_basic(rName, dbPwd),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "instances.#", regexp.MustCompile("\\d+")),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.name"),
				),
			},
		},
	})
}

func TestAccRdsInstanceDataSource_SQLServer_basic(t *testing.T) {
	dataSourceName := "data.g42cloud_rds_instances.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	dc := acceptance.InitDataSourceCheck(dataSourceName)
	dbPwd := fmt.Sprintf("%s%s%d", acctest.RandString(5),
		acctest.RandStringFromCharSet(2, "!#%^*"), acctest.RandIntRange(10, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceDataSource_SQLServer_basic(rName, dbPwd),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "instances.#", regexp.MustCompile("\\d+")),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.name"),
				),
			},
		},
	})
}

func TestAccRdsInstanceDataSource_PostgreSQL_basic(t *testing.T) {
	dataSourceName := "data.g42cloud_rds_instances.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	dc := acceptance.InitDataSourceCheck(dataSourceName)
	dbPwd := fmt.Sprintf("%s%s%d", acctest.RandString(5),
		acctest.RandStringFromCharSet(2, "!#%^*"), acctest.RandIntRange(10, 99))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceDataSource_PostgreSQL_basic(rName, dbPwd),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "instances.#", regexp.MustCompile("\\d+")),
					resource.TestCheckResourceAttrSet(dataSourceName, "instances.0.name"),
				),
			},
		},
	})
}

func testMySQL_base(rName, dbPwd string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_availability_zones" "test" {}

resource "g42cloud_networking_secgroup" "test" {
  name = "%[2]s"
}

resource "g42cloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.mysql.c6.large.4"
  availability_zone = [data.g42cloud_availability_zones.test.names[0]]
  security_group_id = g42cloud_networking_secgroup.test.id
  subnet_id         = g42cloud_vpc_subnet.test.id
  vpc_id            = g42cloud_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    password = "%[3]s"
    type     = "MySQL"
    version  = "8.0"
    port     = 8630
  }
  volume {
    type = "ULTRAHIGH"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  lifecycle {
    ignore_changes = [
      backup_strategy,
    ]
  }
}
`, common.TestVpc(rName), rName, dbPwd)
}

func testSQLServer_base(rName, dbPwd string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_availability_zones" "test" {}

resource "g42cloud_networking_secgroup" "test" {
  name = "%s"
}

resource "g42cloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.mssql.se.c6.large.4"
  availability_zone = [data.g42cloud_availability_zones.test.names[0]]
  security_group_id = g42cloud_networking_secgroup.test.id
  subnet_id         = g42cloud_vpc_subnet.test.id
  vpc_id            = g42cloud_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2019_SE"
    port     = 8631
  }
  volume {
    type = "ULTRAHIGH"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  lifecycle {
    ignore_changes = [
      backup_strategy,
    ]
  }
}
`, common.TestVpc(rName), rName, dbPwd)
}

func testPostgreSQL_base(rName, dbPwd string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_availability_zones" "test" {}

resource "g42cloud_networking_secgroup" "test" {
  name = "%s"
}

resource "g42cloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.pg.c6.large.4"
  availability_zone = [data.g42cloud_availability_zones.test.names[0]]
  security_group_id = g42cloud_networking_secgroup.test.id
  subnet_id         = g42cloud_vpc_subnet.test.id
  vpc_id            = g42cloud_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "14"
    port     = 8632
  }
  volume {
    type = "ULTRAHIGH"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  lifecycle {
    ignore_changes = [
      backup_strategy,
    ]
  }
}
`, common.TestVpc(rName), rName, dbPwd)
}

func testAccRdsInstanceDataSource_basic(name, dbPwd string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_rds_instances" "test" {
  depends_on = [
    g42cloud_rds_instance.test,
  ]
}
`, testMySQL_base(name, dbPwd))
}

func testAccRdsInstanceDataSource_SQLServer_basic(name, dbPwd string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_rds_instances" "test" {
  depends_on = [
    g42cloud_rds_instance.test,
  ]
}
`, testSQLServer_base(name, dbPwd))
}

func testAccRdsInstanceDataSource_PostgreSQL_basic(name, dbPwd string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_rds_instances" "test" {
  depends_on = [
    g42cloud_rds_instance.test,
  ]
}
`, testPostgreSQL_base(name, dbPwd))
}
