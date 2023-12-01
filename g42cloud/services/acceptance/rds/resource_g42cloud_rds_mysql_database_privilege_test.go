package rds

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/pagination"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getMysqlDatabasePrivilegeResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.G42_REGION_NAME
	// getMysqlDatabasePrivilege: query RDS Mysql database privilege
	var (
		getMysqlDatabasePrivilegeHttpUrl = "v3/{project_id}/instances/{instance_id}/database/db_user"
		getMysqlDatabasePrivilegeProduct = "rds"
	)
	getMysqlDatabasePrivilegeClient, err := cfg.NewServiceClient(getMysqlDatabasePrivilegeProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and database from resource id
	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<db_name>")
	}
	instanceId := parts[0]
	dbName := parts[1]

	getMysqlDatabasePrivilegePath := getMysqlDatabasePrivilegeClient.Endpoint + getMysqlDatabasePrivilegeHttpUrl
	getMysqlDatabasePrivilegePath = strings.ReplaceAll(getMysqlDatabasePrivilegePath, "{project_id}",
		getMysqlDatabasePrivilegeClient.ProjectID)
	getMysqlDatabasePrivilegePath = strings.ReplaceAll(getMysqlDatabasePrivilegePath, "{instance_id}", instanceId)

	getMysqlDatabasePrivilegeQueryParams := buildGetMysqlDatabasePrivilegeQueryParams(dbName)
	getMysqlDatabasePrivilegePath += getMysqlDatabasePrivilegeQueryParams

	getMysqlDatabasePrivilegeResp, err := pagination.ListAllItems(
		getMysqlDatabasePrivilegeClient,
		"page",
		getMysqlDatabasePrivilegePath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving Mysql database privilege: %s", err)
	}

	getMysqlDatabasePrivilegeRespJson, err := json.Marshal(getMysqlDatabasePrivilegeResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Mysql database privilege: %s", err)
	}
	var getMysqlDatabasePrivilegeRespBody interface{}
	err = json.Unmarshal(getMysqlDatabasePrivilegeRespJson, &getMysqlDatabasePrivilegeRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Mysql database privilege: %s", err)
	}

	curJson := utils.PathSearch("users", getMysqlDatabasePrivilegeRespBody, make([]interface{}, 0))
	if len(curJson.([]interface{})) == 0 {
		return nil, fmt.Errorf("error get RDS Mysql database privilege")
	}

	return getMysqlDatabasePrivilegeRespBody, nil
}

func buildGetMysqlDatabasePrivilegeQueryParams(dbName string) string {
	return fmt.Sprintf("?db-name=%s&page=1&limit=100", dbName)
}

func TestAccRdsDatabasePrivilege_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "g42cloud_rds_mysql_database_privilege.test"
	dbPwd := fmt.Sprintf("%s%s%d", acctest.RandString(5), acctest.RandStringFromCharSet(2, "!#%^*"),
		acctest.RandIntRange(10, 99))

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getMysqlDatabasePrivilegeResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRdsDatabasePrivilege_basic(name, dbPwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"g42cloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "db_name",
						"g42cloud_rds_mysql_database.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "users.0.name",
						"g42cloud_rds_mysql_account.test_1", "name"),
					resource.TestCheckResourceAttr(rName, "users.0.readonly", "false"),
				),
			},
			{
				Config: testAccRdsDatabasePrivilege_basic_update(name, dbPwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"g42cloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "db_name",
						"g42cloud_rds_mysql_database.test", "name"),
					resource.TestCheckResourceAttrPair(rName, "users.0.name",
						"g42cloud_rds_mysql_account.test_2", "name"),
					resource.TestCheckResourceAttr(rName, "users.0.readonly", "true"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRdsDatabasePrivilege_basic(rName, dbPwd string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_rds_mysql_account" "test_1" {
  instance_id = g42cloud_rds_instance.test.id
  name        = "%s_1"
  password    = "Test@12345678"
}

resource "g42cloud_rds_mysql_database_privilege" "test" {
  instance_id = g42cloud_rds_instance.test.id
  db_name     = g42cloud_rds_mysql_database.test.name

  users {
    name = g42cloud_rds_mysql_account.test_1.name
  }
}
`, testMysqlDatabase_basic(rName, dbPwd, rName), rName)
}

func testAccRdsDatabasePrivilege_basic_update(rName, dbPwd string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_rds_mysql_account" "test_2" {
  instance_id = g42cloud_rds_instance.test.id
  name        = "%s_2"
  password    = "Test@12345678"
}

resource "g42cloud_rds_mysql_database_privilege" "test" {
  instance_id = g42cloud_rds_instance.test.id
  db_name     = g42cloud_rds_mysql_database.test.name

  users {
    name     = g42cloud_rds_mysql_account.test_2.name
    readonly = true
  }
}
`, testMysqlDatabase_basic(rName, dbPwd, rName), rName)
}
