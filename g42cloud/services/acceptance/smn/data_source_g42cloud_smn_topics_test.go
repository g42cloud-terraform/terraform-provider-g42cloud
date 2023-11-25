package smn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
)

func TestAccDataTopics_basic(t *testing.T) {
	dataSourceName := "data.g42cloud_smn_topics.test"
	resourcerName := "g42cloud_smn_topic.topic_1"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTopicsConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttrPair(dataSourceName, "topics.0.id", resourcerName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "topics.0.topic_urn", resourcerName, "topic_urn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "topics.0.display_name", resourcerName, "display_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "topics.0.tags.foo", resourcerName, "tags.foo"),
					resource.TestCheckResourceAttrPair(dataSourceName, "topics.0.tags.key", resourcerName, "tags.key"),
					resource.TestCheckResourceAttrPair(dataSourceName, "topics.0.enterprise_project_id",
						resourcerName, "enterprise_project_id"),
				),
			},
		},
	})
}

func testAccDataTopicsConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "g42cloud_smn_topics" "test" {
  name = "%s"

  depends_on = [
    g42cloud_smn_topic.topic_1
  ]
}
`, testAccSMNV2TopicConfig_basic(rName), rName)
}

func testAccSMNV2TopicConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_smn_topic" "topic_1" {
  name         = "%s"
  display_name = "The display name of %s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName, rName)
}
