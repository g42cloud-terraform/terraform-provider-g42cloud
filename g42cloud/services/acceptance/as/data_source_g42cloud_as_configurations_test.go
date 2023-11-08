package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccDataSourceASConfiguration_basic(t *testing.T) {
	dataSourceName := "data.g42cloud_as_configurations.configurations"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceASConfiguration_conf(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "configurations.0.scaling_configuration_name", name),
				),
			},
		},
	})
}

func testAccASV1Configuration_basic(rName string) string {
	return fmt.Sprintf(`
data "g42cloud_availability_zones" "test" {}

data "g42cloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "g42cloud_compute_flavors" "test" {
  availability_zone = data.g42cloud_availability_zones.test.names[0]
  performance_type  = "normal"
}

resource "g42cloud_compute_keypair" "hth_key" {
  name       = "%s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "g42cloud_as_configuration" "hth_as_config" {
  scaling_configuration_name = "%s"
  instance_config {
    image  = data.g42cloud_images_image.test.id
    flavor = data.g42cloud_compute_flavors.test.ids[0]
    disk {
      size        = 40
      volume_type = "SAS"
      disk_type   = "SYS"
    }
    key_name = "${g42cloud_compute_keypair.hth_key.id}"
  }
}
`, rName, rName)
}

func testAccDataSourceASConfiguration_conf(name string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_as_configurations" "configurations" {
  name     = g42cloud_as_configuration.hth_as_config.scaling_configuration_name
  image_id = g42cloud_as_configuration.hth_as_config.instance_config.0.image

  depends_on = [g42cloud_as_configuration.hth_as_config]
}
`, testAccASV1Configuration_basic(name))
}
