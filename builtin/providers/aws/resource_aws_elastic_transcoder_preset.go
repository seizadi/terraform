package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsElasticTranscoderPreset() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsElasticTranscoderPresetCreate,
		Read:   resourceAwsElasticTranscoderPresetRead,
		Update: resourceAwsElasticTranscoderPresetUpdate,
		Delete: resourceAwsElasticTranscoderPresetDelete,

		Schema: map[string]*schema.Schema{
			"audio": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					// elastictranscoder.AudioParameters
					Schema: map[string]*schema.Schema{
						"audio_packing_mode": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"bitrate": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"channels": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"codec": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"codec_options": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"sample_rate": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"container": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"thumbnails": &schema.Schema{
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					// elastictranscoder.Thumbnails
					Schema: map[string]*schema.Schema{
						"aspect_ratio": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"format": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"interval": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_height": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_width": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"padding_policy": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"resolution:": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"sizing_policy": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"video": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					// elastictranscoder.VideoParameters
					Schema: map[string]*schema.Schema{
						"aspect_ratio": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"bitrate": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"codec": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"codec_options": &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,
						},
						"display_apect_ratio": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"fixed_gop": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"frame_rate": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"key_frame_max_dist": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_frame_rate": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_height": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_width": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"padding_policy": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"resolution": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"sizing_policy": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"watermarks": &schema.Schema{
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								// elastictranscoder.PresetWatermark
								Schema: map[string]*schema.Schema{
									"horizontal_align": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"horizaontal_offset": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"id": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"max_height": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"max_width": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"opacity": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"sizing_policy": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"target": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"vertical_align": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"vertical_offset": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAwsElasticTranscoderPresetCreate(d *schema.ResourceData, meta interface{}) error {
	elastictranscoderconn := meta.(*AWSClient).elastictranscoderconn

	req := &elastictranscoder.CreatePresetInput{}

	if audio, ok := d.GetOk("audio"); ok {
		req.Audio = expandETAudioParams(audio.(*schema.Set))
	}

	req.Container = aws.String(d.Get("container").(string))

	if desc, ok := d.GetOk("description"); ok {
		req.Description = aws.String(desc.(string))
	}

	req.Name = aws.String(d.Get("name").(string))

	if thumbs, ok := d.GetOk("thumbnails"); ok {
		req.Thumbnails = expandETThumbnails(thumbs.(*scheme.Set))
	}

	if video, ok := d.GetOk("video"); ok {
		req.Video = exapndETVideoParams(video.(*scheme.Set))
	}

	log.Printf("[DEBUG] Elastic Transcoder Preset create opts: %s", req)
	resp, err := elastictranscoder.CreatePreset(req)
	if err != nil {
		return fmt.Errorf("Error creating Elastic Transcoder Preset: %s", err)
	}

	return resourceAwsElasticTranscoderPipelineUpdate(d, meta)

}

func expandETThumbnails(s *schema.Set) *elastictranscoder.Thumbnails {
	if s == nil || s.Len() == 0 {
		return nil
	}

	thumbs := &elastictranscoder.Thumbnails{}

	t := s.List()[0].(map[string]interface{})

}

func expandETAudioParams(s *schema.Set) *elastictranscoder.AudioState {
	if s == nil || s.Len() == 0 {
		return nil
	}

	audioParams = &elastictranscoder.AudioParams{}

	audio := s.List()[0].(map[string]interface{})

	if a, ok := audio["audio_packing_mode"]; ok {
		audioParams.AudioPackingMode = aws.String(a.(string))
	}

	if b, ok := audio["bitrate"]; ok {
		audioParams.BitRate = aws.String(b.(string))
	}

	if c, ok := audio["channels"]; ok {
		audioParams.Channels = aws.String(c.(string))
	}

	if c, ok := audio["codec"]; ok {
		audioParams.Codec = aws.String(c.(string))
	}

	audioParams.CodecOptions = expandETAudioCodecOptions(audio["codec_options"].(*schema.Set))

	if s, ok := audio["sample_rate"]; ok {
		audioParams.SampleRate = aws.String(s.(string))
	}

	return audioParams
}

func expandETAudioCodecOptions(s *schema.Set) *slastictranscoder.AudioCodecOptions {
	if s == nil || s.Len() == 0 {
		return nil
	}

	codecOpts := &elastictranscoder.AudioCodecOptions{}
	codec := s.List()[0].(map[string]interface{})

	if b, ok := codec["bit_depth"]; ok {
		codeOpts.BitDepth = aws.String(b.(string))
	}

	if b, ok := codec["bit_order"]; ok {
		codecOpts.BitOrder = aws.String(b.(string))
	}

	if p, ok := codec["profile"]; ok {
		codecOpts.Profile = aws.String(p.(string))
	}

	if s, ok := codec["signed"]; ok {
		codecOpts.Signed = aws.String(s.(string))
	}

	return codecOpts
}

func resourceAwsElasticTranscoderPresetUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsElasticTranscoderPresetRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsElasticTranscoderPresetDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
