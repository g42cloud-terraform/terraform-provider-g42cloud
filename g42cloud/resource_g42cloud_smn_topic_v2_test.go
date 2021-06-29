package g42cloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/smn/v2/topics"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccSMNV2Topic_basic(t *testing.T) {
	var topic topics.TopicGet
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	displayName := fmt.Sprintf("The display name of %s", rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSMNTopicV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSMNV2TopicConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSMNV2TopicExists("g42cloud_smn_topic.topic_1", &topic),
					resource.TestCheckResourceAttr(
						"g42cloud_smn_topic.topic_1", "name", rName),
					resource.TestCheckResourceAttr(
						"g42cloud_smn_topic.topic_1", "display_name",
						displayName),
				),
			},
		},
	})
}

func testAccCheckSMNTopicV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	smnClient, err := config.SmnV2Client(G42_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating G42Cloud smn: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "g42cloud_smn_topic" {
			continue
		}

		_, err := topics.Get(smnClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Topic still exists")
		}
	}

	return nil
}

func testAccCheckSMNV2TopicExists(n string, topic *topics.TopicGet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		smnClient, err := config.SmnV2Client(G42_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating G42Cloud smn client: %s", err)
		}

		found, err := topics.Get(smnClient, rs.Primary.ID).ExtractGet()
		if err != nil {
			return err
		}

		if found.TopicUrn != rs.Primary.ID {
			return fmt.Errorf("Topic not found")
		}

		*topic = *found

		return nil
	}
}

func testAccSMNV2TopicConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_smn_topic" "topic_1" {
  name         = "%s"
  display_name = "The display name of %s"
}
`, rName, rName)
}
