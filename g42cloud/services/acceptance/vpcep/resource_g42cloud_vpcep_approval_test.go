package vpcep

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/vpcep/v1/endpoints"
	"github.com/chnsz/golangsdk/openstack/vpcep/v1/services"
	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVPCEndpointApproval_Basic(t *testing.T) {
	var service services.Service
	var endpoint endpoints.Endpoint

	rName := fmt.Sprintf("acc-test-%s", acctest.RandString(4))
	resourceName := "g42cloud_vpcep_approval.approval"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVPCEPServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEndpointApproval_Basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCEPServiceExists("g42cloud_vpcep_service.test", &service),
					testAccCheckVPCEndpointExists("g42cloud_vpcep_endpoint.test", &endpoint),
					resource.TestCheckResourceAttrPtr(resourceName, "id", &service.ID),
					resource.TestCheckResourceAttrPtr(resourceName, "connections.0.endpoint_id", &endpoint.ID),
					resource.TestCheckResourceAttr(resourceName, "connections.0.status", "accepted"),
				),
			},
			{
				Config: testAccVPCEndpointApproval_Update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceName, "connections.0.endpoint_id", &endpoint.ID),
					resource.TestCheckResourceAttr(resourceName, "connections.0.status", "rejected"),
				),
			},
		},
	})
}

func testAccVPCEndpointApproval_Basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_vpcep_service" "test" {
  name        = "%s"
  server_type = "VM"
  vpc_id      = data.g42cloud_vpc.myvpc.id
  port_id     = g42cloud_compute_instance.ecs.network[0].port
  approval    = true

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
  tags = {
    owner = "tf-acc"
  }
}

resource "g42cloud_vpcep_endpoint" "test" {
  service_id  = g42cloud_vpcep_service.test.id
  vpc_id      = data.g42cloud_vpc.myvpc.id
  network_id  = data.g42cloud_vpc_subnet.test.id
  enable_dns  = true

  tags = {
    owner = "tf-acc"
  }
  lifecycle {
    ignore_changes = [enable_dns]
  }
}

resource "g42cloud_vpcep_approval" "approval" {
  service_id = g42cloud_vpcep_service.test.id
  endpoints  = [g42cloud_vpcep_endpoint.test.id]
}
`, testAccVPCEndpoint_Precondition(rName), rName)
}

func testAccVPCEndpointApproval_Update(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_vpcep_service" "test" {
  name        = "%s"
  server_type = "VM"
  vpc_id      = data.g42cloud_vpc.myvpc.id
  port_id     = g42cloud_compute_instance.ecs.network[0].port
  approval    = true

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
  tags = {
    owner = "tf-acc"
  }
}

resource "g42cloud_vpcep_endpoint" "test" {
  service_id  = g42cloud_vpcep_service.test.id
  vpc_id      = data.g42cloud_vpc.myvpc.id
  network_id  = data.g42cloud_vpc_subnet.test.id
  enable_dns  = true

  tags = {
    owner = "tf-acc"
  }
  lifecycle {
    ignore_changes = [enable_dns]
  }
}

resource "g42cloud_vpcep_approval" "approval" {
  service_id = g42cloud_vpcep_service.test.id
  endpoints  = []
}
`, testAccVPCEndpoint_Precondition(rName), rName)
}
