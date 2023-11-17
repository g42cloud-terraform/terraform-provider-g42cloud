package mpc

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	mpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/mpc/v1/model"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getTemplateResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.HcMpcV1Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating MPC client: %s", err)
	}

	id, err := strconv.ParseInt(state.Primary.ID, 10, 32)
	if err != nil {
		return nil, err
	}

	resp, err := c.ListTemplate(&mpc.ListTemplateRequest{TemplateId: &[]int32{int32(id)}})
	if err != nil {
		return nil, fmt.Errorf("error retrieving MPC transcoding template: %d", err)
	}

	templateList := *resp.TemplateArray
	template := templateList[0].Template
	if template == nil {
		return nil, fmt.Errorf("unable to retrieve MPC transcoding template: %d", id)
	}

	return template, nil
}

func TestAccTranscodingTemplate_basic(t *testing.T) {
	var template mpc.QueryTransTemplate
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := rName + "-update"
	resourceName := "g42cloud_mpc_transcoding_template.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&template,
		getTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTranscodingTemplate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "low_bitrate_hd", "true"),
					resource.TestCheckResourceAttr(resourceName, "output_format", "1"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.codec", "1"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.output_policy", "transcode"),
					resource.TestCheckResourceAttr(resourceName, "video.0.codec", "1"),
					resource.TestCheckResourceAttr(resourceName, "video.0.profile", "1"),
					resource.TestCheckResourceAttr(resourceName, "video.0.output_policy", "transcode"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testTranscodingTemplate_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "low_bitrate_hd", "true"),
					resource.TestCheckResourceAttr(resourceName, "output_format", "2"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.codec", "1"),
					resource.TestCheckResourceAttr(resourceName, "video.0.codec", "1"),
					resource.TestCheckResourceAttr(resourceName, "video.0.profile", "2"),
				),
			},
		},
	})
}

func TestAccTranscodingTemplate_basic2(t *testing.T) {
	var template mpc.QueryTransTemplate
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "g42cloud_mpc_transcoding_template.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&template,
		getTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTranscodingTemplate_basic2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "low_bitrate_hd", "true"),
					resource.TestCheckResourceAttr(resourceName, "output_format", "3"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.codec", "1"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.output_policy", "transcode"),
					resource.TestCheckResourceAttr(resourceName, "video.0.codec", "1"),
					resource.TestCheckResourceAttr(resourceName, "video.0.profile", "3"),
					resource.TestCheckResourceAttr(resourceName, "video.0.output_policy", "transcode"),
				),
			},
		},
	})
}

func TestAccTranscodingTemplate_basic3(t *testing.T) {
	var template mpc.QueryTransTemplate
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "g42cloud_mpc_transcoding_template.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&template,
		getTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTranscodingTemplate_basic3(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "low_bitrate_hd", "true"),
					resource.TestCheckResourceAttr(resourceName, "output_format", "4"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.codec", "2"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.output_policy", "transcode"),
					resource.TestCheckResourceAttr(resourceName, "video.0.codec", "2"),
					resource.TestCheckResourceAttr(resourceName, "video.0.profile", "4"),
					resource.TestCheckResourceAttr(resourceName, "video.0.output_policy", "transcode"),
				),
			},
		},
	})
}

func TestAccTranscodingTemplate_basic4(t *testing.T) {
	var template mpc.QueryTransTemplate
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "g42cloud_mpc_transcoding_template.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&template,
		getTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTranscodingTemplate_basic4(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "low_bitrate_hd", "true"),
					resource.TestCheckResourceAttr(resourceName, "output_format", "5"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.codec", "4"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.output_policy", "transcode"),
				),
			},
		},
	})
}

func TestAccTranscodingTemplate_basic5(t *testing.T) {
	var template mpc.QueryTransTemplate
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "g42cloud_mpc_transcoding_template.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&template,
		getTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTranscodingTemplate_basic5(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "low_bitrate_hd", "true"),
					resource.TestCheckResourceAttr(resourceName, "output_format", "6"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.codec", "2"),
					resource.TestCheckResourceAttr(resourceName, "audio.0.output_policy", "transcode"),
				),
			},
		},
	})
}

func testTranscodingTemplate_basic(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_mpc_transcoding_template" "test" {
  name                  = "%s"
  low_bitrate_hd        = true
  dash_segment_duration = 5
  hls_segment_duration  = 5
  output_format         = 1

  audio {
    codec         = 1
    sample_rate   = 1
    channels      = 1
    bitrate       = 0
    output_policy = "transcode"
  }

  video {
    codec                   = 1
    profile                 = 1
    level                   = 1
    quality                 = 1
    black_bar_removal       = 0
    max_consecutive_bframes = 7
    bitrate                 = 0
    fps                     = 0
    max_iframes_interval    = 5
    output_policy           = "transcode"
    height                  = 0
    width                   = 0
  }
}
`, rName)
}

func testTranscodingTemplate_update(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_mpc_transcoding_template" "test" {
  name                  = "%s"
  low_bitrate_hd        = true
  dash_segment_duration = 5
  hls_segment_duration  = 5
  output_format         = 2

  audio {
    codec         = 1
    sample_rate   = 2
    channels      = 1
    bitrate       = 0
    output_policy = "transcode"
  }

  video {
    codec                   = 1
    profile                 = 2
    level                   = 3
    quality                 = 1
    black_bar_removal       = 0
    max_consecutive_bframes = 7
    bitrate                 = 0
    fps                     = 0
    max_iframes_interval    = 5
    output_policy           = "transcode"
    height                  = 0
    width                   = 0
  }
}
`, rName)
}

func testTranscodingTemplate_basic2(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_mpc_transcoding_template" "test" {
  name                  = "%s"
  low_bitrate_hd        = true
  dash_segment_duration = 5
  hls_segment_duration  = 5
  output_format         = 3

  audio {
    codec         = 1
    sample_rate   = 3
    channels      = 2
    bitrate       = 0
    output_policy = "transcode"
  }

  video {
    codec                   = 1
    profile                 = 3
    level                   = 4
    quality                 = 3
    black_bar_removal       = 1
    max_consecutive_bframes = 7
    bitrate                 = 0
    fps                     = 0
    max_iframes_interval    = 5
    output_policy           = "transcode"
    height                  = 0
    width                   = 0
  }
}
`, rName)
}

func testTranscodingTemplate_basic3(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_mpc_transcoding_template" "test" {
  name                  = "%s"
  low_bitrate_hd        = true
  dash_segment_duration = 5
  hls_segment_duration  = 5
  output_format         = 4

  audio {
    codec         = 2
    sample_rate   = 4
    channels      = 2
    bitrate       = 0
    output_policy = "transcode"
  }

  video {
    codec                   = 2
    profile                 = 4
    level                   = 10
    quality                 = 2
    black_bar_removal       = 1
    max_consecutive_bframes = 7
    bitrate                 = 0
    fps                     = 0
    max_iframes_interval    = 5
    output_policy           = "transcode"
    height                  = 0
    width                   = 0
  }
}
`, rName)
}

func testTranscodingTemplate_basic4(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_mpc_transcoding_template" "test" {
  name                  = "%s"
  low_bitrate_hd        = true
  dash_segment_duration = 5
  hls_segment_duration  = 5
  output_format         = 5

  audio {
    codec         = 4
    sample_rate   = 5
    channels      = 2
    bitrate       = 0
    output_policy = "transcode"
  }
}
`, rName)
}

func testTranscodingTemplate_basic5(rName string) string {
	return fmt.Sprintf(`
resource "g42cloud_mpc_transcoding_template" "test" {
  name                  = "%s"
  low_bitrate_hd        = true
  dash_segment_duration = 5
  hls_segment_duration  = 5
  output_format         = 6

  audio {
    codec         = 2
    sample_rate   = 6
    channels      = 1
    bitrate       = 0
    output_policy = "transcode"
  }
}
`, rName)
}
