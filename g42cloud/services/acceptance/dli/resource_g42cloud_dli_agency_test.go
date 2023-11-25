package dli

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDliAgencyResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	region := acceptance.G42_REGION_NAME
	// getAgency: Query the Agency.
	var (
		getAgencyHttpUrl = "v2/{project_id}/agency"
		getAgencyProduct = "dli"
	)
	getAgencyClient, err := cfg.NewServiceClient(getAgencyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DLI Client: %s", err)
	}

	getAgencyPath := getAgencyClient.Endpoint + getAgencyHttpUrl
	getAgencyPath = strings.ReplaceAll(getAgencyPath, "{project_id}", getAgencyClient.ProjectID)

	getAgencyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getAgencyResp, err := getAgencyClient.Request("GET", getAgencyPath, &getAgencyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DLI Agency: %s", err)
	}

	getAgencyRespBody, err := utils.FlattenResponse(getAgencyResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DLI Agency: %s", err)
	}

	roles := utils.PathSearch("current_roles", getAgencyRespBody, nil)
	if v, ok := roles.([]interface{}); !ok || len(v) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}
	return roles, nil
}

func TestAccDliAgency_basic(t *testing.T) {
	var obj interface{}

	rName := "g42cloud_dli_agency.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDliAgencyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliAgency(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDliAgency_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "roles.#", "3"),
				),
			},
			{
				Config: testDliAgency_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "roles.#", "4"),
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

func testDliAgency_basic() string {
	return `
resource "g42cloud_dli_agency" "test" {
  roles = [
    "te_admin",
    "dis_adm",
    "vpc_netadm",
  ]
}
`
}

func testDliAgency_basic_update() string {
	return `
resource "g42cloud_dli_agency" "test" {
  roles = [
    "te_admin",
    "dis_adm",
    "vpc_netadm",
    "smn_adm",
  ]
}
`
}
