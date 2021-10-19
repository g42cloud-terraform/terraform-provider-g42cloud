package g42cloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3/groups"
	"github.com/chnsz/golangsdk/openstack/identity/v3/projects"
	"github.com/chnsz/golangsdk/openstack/identity/v3/roles"
	"github.com/chnsz/golangsdk/pagination"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func extractRoleAssignmentID(roleAssignmentID string) (string, string, string, string) {
	split := strings.Split(roleAssignmentID, "/")
	return split[0], split[1], split[2], split[3]
}

func TestAccIdentityV3RoleAssignment_basic(t *testing.T) {
	var role roles.Role
	var group groups.Group
	var project projects.Project
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3RoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3RoleAssignment_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3RoleAssignmentExists("g42cloud_identity_role_assignment_v3.role_assignment_1", &role, &group, &project),
					resource.TestCheckResourceAttrPtr(
						"g42cloud_identity_role_assignment_v3.role_assignment_1", "project_id", &project.ID),
					resource.TestCheckResourceAttrPtr(
						"g42cloud_identity_role_assignment_v3.role_assignment_1", "group_id", &group.ID),
					resource.TestCheckResourceAttrPtr(
						"g42cloud_identity_role_assignment_v3.role_assignment_1", "role_id", &role.ID),
				),
			},
		},
	})
}

func testAccCheckIdentityV3RoleAssignmentDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	identityClient, err := config.IdentityV3Client(G42_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating G42Cloud identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "g42cloud_identity_role_assignment_v3" {
			continue
		}

		_, err := roles.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Role assignment still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3RoleAssignmentExists(n string, role *roles.Role, group *groups.Group, project *projects.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		identityClient, err := config.IdentityV3Client(G42_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating G42Cloud identity client: %s", err)
		}

		domainID, projectID, groupID, roleID := extractRoleAssignmentID(rs.Primary.ID)

		var opts roles.ListAssignmentsOpts
		opts = roles.ListAssignmentsOpts{
			GroupID:        groupID,
			ScopeDomainID:  domainID,
			ScopeProjectID: projectID,
		}

		pager := roles.ListAssignments(identityClient, opts)
		var assignment roles.RoleAssignment

		err = pager.EachPage(func(page pagination.Page) (bool, error) {
			assignmentList, err := roles.ExtractRoleAssignments(page)
			if err != nil {
				return false, err
			}

			for _, a := range assignmentList {
				if a.ID == roleID {
					assignment = a
					return false, nil
				}
			}

			return true, nil
		})
		if err != nil {
			return err
		}

		p, err := projects.Get(identityClient, projectID).Extract()
		if err != nil {
			return fmt.Errorf("Project not found")
		}
		*project = *p
		g, err := groups.Get(identityClient, groupID).Extract()
		if err != nil {
			return fmt.Errorf("Group not found")
		}
		*group = *g
		r, err := roles.Get(identityClient, assignment.ID).Extract()
		if err != nil {
			return fmt.Errorf("Role not found")
		}
		*role = *r

		return nil
	}
}

const testAccIdentityV3RoleAssignment_basic = `
resource "g42cloud_identity_group_v3" "group_1" {
  name = "group_1"
}

data "g42cloud_identity_role_v3" "role_1" {
  name = "ims_adm"
}

resource "g42cloud_identity_role_assignment_v3" "role_assignment_1" {
  group_id = "${g42cloud_identity_group_v3.group_1.id}"
  #project_id = "09338f629380276a2f95c0180c2bdb72"
  project_id = "0910fc31530026f82fd0c018a303517e"
  role_id = "${data.g42cloud_identity_role_v3.role_1.id}"
}
`
