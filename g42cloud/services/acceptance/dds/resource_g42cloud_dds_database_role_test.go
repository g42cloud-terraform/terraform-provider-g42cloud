package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dds/v3/roles"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDatabaseRoleFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DdsV3Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DDS v3 client: %s ", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	name := state.Primary.Attributes["name"]
	opts := roles.ListOpts{
		Name:   state.Primary.Attributes["name"],
		DbName: state.Primary.Attributes["db_name"],
	}
	resp, err := roles.List(client, instanceId, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting role (%s) from DDS instance (%s): %v", name, instanceId, err)
	}
	if len(resp) < 1 {
		return nil, fmt.Errorf("unable to find role (%s) from DDS instance (%s)", name, instanceId)
	}
	role := resp[0]
	return &role, nil
}

func TestAccDatabaseRole_basic(t *testing.T) {
	var role roles.Role
	rName := acceptance.RandomAccResourceName()
	resourceName := "g42cloud_dds_database_role.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&role,
		getDatabaseRoleFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseRole_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "roles.0.name",
						"g42cloud_dds_database_role.base", "name"),
					resource.TestCheckResourceAttrPair(resourceName, "inherited_privileges",
						"g42cloud_dds_database_role.base", "inherited_privileges"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccDatabaseRoleImportStateIdFunc(),
			},
		},
	})
}

func testAccDatabaseRoleImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, dbName, name string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "g42cloud_dds_database_role" {
				instanceId = rs.Primary.Attributes["instance_id"]
				dbName = rs.Primary.Attributes["db_name"]
				name = rs.Primary.Attributes["name"]
			}
		}
		if instanceId == "" || dbName == "" || name == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s", instanceId, dbName, name)
		}
		return fmt.Sprintf("%s/%s/%s", instanceId, dbName, name), nil
	}
}

func testAccDatabaseRole_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "g42cloud_dds_database_role" "base" {
  instance_id = g42cloud_dds_instance.test.id

  name    = "%[2]s-base"
  db_name = "admin"
}

resource "g42cloud_dds_database_role" "test" {
  instance_id = g42cloud_dds_instance.test.id

  name    = "%[2]s"
  db_name = "admin"

  roles {
    name    = g42cloud_dds_database_role.base.name
    db_name = "admin"
  }
}
`, testAccDatasourceDdsInstance_base(rName, 8800), rName)
}
