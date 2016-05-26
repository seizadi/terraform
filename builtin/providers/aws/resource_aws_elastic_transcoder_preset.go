package aws

import "github.com/hashicorp/terraform/helper/schema"

func resourceAwsElasticYranscoderPreset() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsElasticTranscoderPresetCreate,
		Read:   resourceAwsElasticTranscoderPresetRead,
		Update: resourceAwsElasticTranscoderPresetUpdate,
		Delete: resourceAwsElasticTranscoderPresetDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

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
	return nil
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
