package modelarts

import (
	"fmt"
	"testing"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/modelarts/v1/notebook"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getNotebookResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.ModelArtsV1Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("error creating ModelArts v1 client, err=%s", err)
	}

	return notebook.Get(client, state.Primary.ID)
}

func TestAccResourceNotebook_basic(t *testing.T) {
	var instance notebook.CreateOpts
	resourceName := "g42cloud_modelarts_notebook.test"
	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getNotebookResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNotebook_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "modelarts.vm.cpu.2u"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "EFS"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.ownership", "MANAGED"),
					resource.TestCheckResourceAttr(resourceName, "auto_stop_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "image_name"),
					resource.TestCheckResourceAttrSet(resourceName, "image_swr_path"),
					resource.TestCheckResourceAttrSet(resourceName, "image_type"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
				),
			},
			{
				Config: testAccNotebook_basic(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "modelarts.vm.cpu.2u"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "EFS"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.ownership", "MANAGED"),
					resource.TestCheckResourceAttr(resourceName, "auto_stop_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "image_name"),
					resource.TestCheckResourceAttrSet(resourceName, "image_swr_path"),
					resource.TestCheckResourceAttrSet(resourceName, "image_type"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNotebook_basic(rName string) string {
	return fmt.Sprintf(`
data "g42cloud_modelarts_notebook_images" "test" {
  type     = "BUILD_IN"
  cpu_arch = "x86_64"
}

resource "g42cloud_modelarts_notebook" "test" {
  name      = "%s"
  flavor_id = "modelarts.vm.cpu.2u"
  image_id  = data.g42cloud_modelarts_notebook_images.test.images[0].id
  volume {
    type = "EFS"
  }
}
`, rName)
}

func TestAccResourceNotebook_all(t *testing.T) {
	var instance notebook.CreateOpts
	resourceName := "g42cloud_modelarts_notebook.test"
	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	ip := "10.1.1.2"
	updateIp := "10.1.1.3"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getNotebookResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNotebook_All(name, ip),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "modelarts.vm.cpu.2u"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "EFS"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.ownership", "MANAGED"),
					resource.TestCheckResourceAttr(resourceName, "auto_stop_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", name),
					resource.TestCheckResourceAttr(resourceName, "allowed_access_ips.0", ip),
					resource.TestCheckResourceAttrSet(resourceName, "image_name"),
					resource.TestCheckResourceAttrSet(resourceName, "image_swr_path"),
					resource.TestCheckResourceAttrSet(resourceName, "image_type"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
				),
			},
			{
				Config: testAccNotebook_All(updateName, updateIp),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "flavor_id", "modelarts.vm.cpu.2u"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.type", "EFS"),
					resource.TestCheckResourceAttr(resourceName, "volume.0.ownership", "MANAGED"),
					resource.TestCheckResourceAttr(resourceName, "auto_stop_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", updateName),
					resource.TestCheckResourceAttr(resourceName, "allowed_access_ips.0", updateIp),
					resource.TestCheckResourceAttrSet(resourceName, "image_name"),
					resource.TestCheckResourceAttrSet(resourceName, "image_swr_path"),
					resource.TestCheckResourceAttrSet(resourceName, "image_type"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNotebook_All(rName string, ip string) string {
	return fmt.Sprintf(`
data "g42cloud_modelarts_notebook_images" "test" {
  type     = "BUILD_IN"
  cpu_arch = "x86_64"
}

resource "g42cloud_compute_keypair" "test" {
  name       = "%s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD2LRDdahet4k5CjShfWF9GPWyRyIbIC4gkPrb1LTL1W8S8Hr03Q1k3yFKcJO18xyNWVj3Z9fgpDznvgwE6Uw9JXSlJiXwJ38PkA68CXs/bOqbGiEgYU+9J1KGwbegMLqQF8CM9Xx2r6jUQcu6L6SnbDSrF9Gf4eLoEyJG8ZYDDnpozaFz853Xxxww8Ldf3YZdYkvbv/pmjrMSLSuHINErxm3VyMcaAN1m5uMoppejiZesfD7Z4kWGnFtgUqIwPbuUiXYwGnIdI95k7hYeg5azev+87cgw6vh3J458YvMia3WFIZVkaqK0EuD4f6U6l+5qpRJRUtLRKx0Jn8gBEsR4v Generated-by-Nova\n"
}

resource "g42cloud_modelarts_notebook" "test" {
  name        = "%s"
  flavor_id   = "modelarts.vm.cpu.2u"
  image_id    = data.g42cloud_modelarts_notebook_images.test.images[0].id
  description = "%s"

  allowed_access_ips = ["%s"]
  key_pair           = g42cloud_compute_keypair.test.name

  volume {
    type = "EFS"
  }
}
`, rName, rName, rName, ip)
}
