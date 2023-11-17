package mpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	mpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/mpc/v1/model"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getTemplateGroupResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.HcMpcV1Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating MPC client: %s", err)
	}

	resp, err := c.ListTemplateGroup(&mpc.ListTemplateGroupRequest{GroupId: &[]string{state.Primary.ID}})
	if err != nil {
		return nil, fmt.Errorf("error retrieving MPC transcoding template group: %s", err)
	}

	templateGroupList := *resp.TemplateGroupList

	if len(templateGroupList) == 0 {
		return nil, fmt.Errorf("unable to retrieve MPC transcoding template group: %s", state.Primary.ID)
	}

	return templateGroupList[0], nil
}

func TestAccTranscodingTemplateGroup_basic(t *testing.T) {
	var templateGroup mpc.TemplateGroup
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := rName + "-update"
	resourceName := "g42cloud_mpc_transcoding_template_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&templateGroup,
		getTemplateGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTranscodingTemplateGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "low_bitrate_hd", "true"),
					resource.TestCheckResourceAttr(resourceName, "output_format", "1"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.codec", "1"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.output_policy", "transcode"),
					resource.TestCheckResourceAttr(resourceName, "video_common.0.codec", "1"),
					resource.TestCheckResourceAttr(resourceName, "video_common.0.profile", "1"),
					resource.TestCheckResourceAttr(resourceName, "video_common.0.output_policy", "transcode"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testTranscodingTemplateGroup_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "low_bitrate_hd", "false"),
					resource.TestCheckResourceAttr(resourceName, "output_format", "2"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.codec", "1"),
					resource.TestCheckResourceAttr(resourceName, "video_common.0.codec", "1"),
					resource.TestCheckResourceAttr(resourceName, "video_common.0.profile", "2"),
				),
			},
		},
	})
}

func TestAccTranscodingTemplateGroup_basic2(t *testing.T) {
	var templateGroup mpc.TemplateGroup
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "g42cloud_mpc_transcoding_template_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&templateGroup,
		getTemplateGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTranscodingTemplateGroup_basic2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "low_bitrate_hd", "true"),
					resource.TestCheckResourceAttr(resourceName, "output_format", "3"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.codec", "1"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.output_policy", "transcode"),
					resource.TestCheckResourceAttr(resourceName, "video_common.0.codec", "1"),
					resource.TestCheckResourceAttr(resourceName, "video_common.0.profile", "3"),
					resource.TestCheckResourceAttr(resourceName, "video_common.0.output_policy", "transcode"),
				),
			},
		},
	})
}

func TestAccTranscodingTemplateGroup_basic3(t *testing.T) {
	var templateGroup mpc.TemplateGroup
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "g42cloud_mpc_transcoding_template_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&templateGroup,
		getTemplateGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTranscodingTemplateGroup_basic3(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "low_bitrate_hd", "true"),
					resource.TestCheckResourceAttr(resourceName, "output_format", "4"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.codec", "2"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.output_policy", "transcode"),
					resource.TestCheckResourceAttr(resourceName, "video_common.0.codec", "2"),
					resource.TestCheckResourceAttr(resourceName, "video_common.0.profile", "4"),
					resource.TestCheckResourceAttr(resourceName, "video_common.0.output_policy", "transcode"),
				),
			},
		},
	})
}

func testTranscodingTemplateGroup_basic(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_mpc_transcoding_template_group" "test" {
  name                  = "%s"
  low_bitrate_hd        = true
  dash_segment_duration = 5
  hls_segment_duration  = 5
  output_format         = 1

  audio {
    codec         = 1
    sample_rate   = 1
    channels      = 6
    output_policy = "transcode"
    bitrate       = 0
  }

  video_common {
    codec                   = 1
    profile                 = 1
    level                   = 1
    quality                 = 1
    black_bar_removal       = 0
    max_consecutive_bframes = 7
    fps                     = 0
    max_iframes_interval    = 5
    output_policy           = "transcode"
  }

  videos {
    width   = 1920
    height  = 1080
    bitrate = 0
  }
}
`, rName)
}

func testTranscodingTemplateGroup_update(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_mpc_transcoding_template_group" "test" {
  name                  = "%s"
  low_bitrate_hd        = false
  dash_segment_duration = 5
  hls_segment_duration  = 5
  output_format         = 2

  audio {
    codec         = 1
    sample_rate   = 2
    channels      = 6
    output_policy = "transcode"
    bitrate       = 0
  }

  video_common {
    codec                   = 1
    profile                 = 2
    level                   = 4
    quality                 = 1
    black_bar_removal       = 0
    max_consecutive_bframes = 7
    fps                     = 0
    max_iframes_interval    = 5
    output_policy           = "transcode"
  }

  videos {
    width   = 3840
    height  = 2560
    bitrate = 0
  }
}
`, rName)
}

func testTranscodingTemplateGroup_basic2(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_mpc_transcoding_template_group" "test" {
  name                  = "%s"
  low_bitrate_hd        = true
  dash_segment_duration = 5
  hls_segment_duration  = 5
  output_format         = 3

  audio {
    codec         = 1
    sample_rate   = 3
    channels      = 2
    output_policy = "transcode"
    bitrate       = 0
  }

  video_common {
    codec                   = 1
    profile                 = 3
    level                   = 7
    quality                 = 2
    black_bar_removal       = 1
    max_consecutive_bframes = 7
    fps                     = 0
    max_iframes_interval    = 5
    output_policy           = "transcode"
  }

  videos {
    width   = 1920
    height  = 1080
    bitrate = 0
  }
}
`, rName)
}

func testTranscodingTemplateGroup_basic3(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_mpc_transcoding_template_group" "test" {
  name                  = "%s"
  low_bitrate_hd        = true
  dash_segment_duration = 5
  hls_segment_duration  = 5
  output_format         = 4

  audio {
    codec         = 2
    sample_rate   = 4
    channels      = 2
    output_policy = "transcode"
    bitrate       = 0
  }

  video_common {
    codec                   = 2
    profile                 = 4
    level                   = 10
    quality                 = 2
    black_bar_removal       = 1
    max_consecutive_bframes = 7
    fps                     = 0
    max_iframes_interval    = 5
    output_policy           = "transcode"
  }

  videos {
    width   = 1920
    height  = 1080
    bitrate = 0
  }
}
`, rName)
}
