package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dds/v3/users"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDatabaseUserFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DdsV3Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DDS v3 client: %s ", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	name := state.Primary.Attributes["name"]
	opts := users.ListOpts{
		Name:   state.Primary.Attributes["name"],
		DbName: state.Primary.Attributes["db_name"],
	}
	resp, err := users.List(client, instanceId, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting user (%s) from DDS instance (%s): %v", name, instanceId, err)
	}
	if len(resp) < 1 {
		return nil, fmt.Errorf("unable to find user (%s) from DDS instance (%s)", name, instanceId)
	}
	user := resp[0]
	return &user, nil
}

func TestAccDatabaseUser_basic(t *testing.T) {
	var user users.UserResp
	rName := acceptance.RandomAccResourceName()
	resourceName := "g42cloud_dds_database_user.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&user,
		getDatabaseUserFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseUser_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "roles.0.name",
						"g42cloud_dds_database_role.test", "name"),
					resource.TestCheckResourceAttrPair(resourceName, "inherited_privileges",
						"g42cloud_dds_database_role.test", "inherited_privileges"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDatabaseUserImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
		},
	})
}

func testAccDatabaseUserImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, dbName, userName string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "g42cloud_dds_database_user" {
				instanceId = rs.Primary.Attributes["instance_id"]
				dbName = rs.Primary.Attributes["db_name"]
				userName = rs.Primary.Attributes["name"]
			}
		}
		if instanceId == "" || dbName == "" || userName == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s", instanceId, dbName, userName)
		}
		return fmt.Sprintf("%s/%s/%s", instanceId, dbName, userName), nil
	}
}

func testAccDatabaseUser_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "g42cloud_dds_database_role" "test" {
  instance_id = g42cloud_dds_instance.test.id

  name    = "%[2]s"
  db_name = "admin"
}

resource "g42cloud_dds_database_user" "test" {
  instance_id = g42cloud_dds_instance.test.id

  name     = "%[2]s"
  password = "G42CloudTest@12345678"
  db_name  = "admin"

  roles {
    name    = g42cloud_dds_database_role.test.name
    db_name = "admin"
  }
}
`, testAccDatasourceDdsInstance_base(rName, 8860), rName)
}
