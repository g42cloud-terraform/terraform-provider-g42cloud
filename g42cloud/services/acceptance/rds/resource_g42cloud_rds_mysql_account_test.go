package rds

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/pagination"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getMysqlAccountResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.G42_REGION_NAME
	// getMysqlAccount: query RDS Mysql account
	var (
		getMysqlAccountHttpUrl = "v3/{project_id}/instances/{instance_id}/db_user/detail?page=1&limit=100"
		getMysqlAccountProduct = "rds"
	)
	getMysqlAccountClient, err := cfg.NewServiceClient(getMysqlAccountProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and user from resource id
	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<name>")
	}
	instanceId := parts[0]
	accountName := parts[1]

	getMysqlAccountPath := getMysqlAccountClient.Endpoint + getMysqlAccountHttpUrl
	getMysqlAccountPath = strings.ReplaceAll(getMysqlAccountPath, "{project_id}", getMysqlAccountClient.ProjectID)
	getMysqlAccountPath = strings.ReplaceAll(getMysqlAccountPath, "{instance_id}", instanceId)

	getMysqlAccountResp, err := pagination.ListAllItems(
		getMysqlAccountClient,
		"page",
		getMysqlAccountPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS Mysql account")
	}

	getMysqlAccountRespJson, err := json.Marshal(getMysqlAccountResp)
	if err != nil {
		return nil, err
	}
	var getMysqlAccountRespBody interface{}
	err = json.Unmarshal(getMysqlAccountRespJson, &getMysqlAccountRespBody)
	if err != nil {
		return nil, err
	}

	account := utils.PathSearch(fmt.Sprintf("users[?name=='%s']|[0]", accountName), getMysqlAccountRespBody, nil)

	if account != nil {
		return account, nil
	}

	return nil, fmt.Errorf("error get RDS Mysql account by instanceID %s and account %s", instanceId, accountName)
}

func TestAccMysqlAccount_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "g42cloud_rds_mysql_account.test"
	dbPwd := fmt.Sprintf("%s%s%d", acctest.RandString(5),
		acctest.RandStringFromCharSet(2, "!#%^*"), acctest.RandIntRange(10, 99))

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getMysqlAccountResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testMysqlAccount_basic(name, dbPwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"g42cloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test_description"),
					resource.TestCheckResourceAttr(rName, "hosts.0", "10.10.%"),
				),
			},
			{
				Config: testMysqlAccount_basic_update(name, dbPwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"g42cloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test_description_update"),
					resource.TestCheckResourceAttr(rName, "hosts.0", "10.10.%"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testMysql_base(name, dbPwd string) string {
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
`, common.TestVpc(name), name, dbPwd)
}

func testMysqlAccount_basic(name, dbPwd string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_rds_mysql_account" "test" {
  instance_id = g42cloud_rds_instance.test.id
  name        = "%s"
  password    = "Test@12345678"
  description = "test_description"

  hosts = [
    "10.10.%%"
  ]
}
`, testMysql_base(name, dbPwd), name)
}

func testMysqlAccount_basic_update(name, dbPwd string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_rds_mysql_account" "test" {
  instance_id = g42cloud_rds_instance.test.id
  name        = "%s"
  password    = "Test@123456789"
  description = "test_description_update"

  hosts = [
    "10.10.%%"
  ]
}
`, testMysql_base(name, dbPwd), name)
}
