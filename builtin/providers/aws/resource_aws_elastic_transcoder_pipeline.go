package aws

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elastictranscoder"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsElasticTranscoderPipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsElasticTranscoderPipelineCreate,
		Read:   resourceAwsElasticTranscoderPipelineRead,
		Update: resourceAwsElasticTranscoderPipelineUpdate,
		Delete: resourceAwsElasticTranscoderPipelineDelete,

		Schema: map[string]*schema.Schema{
			"arn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"aws_kms_key_arn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			// ContentConfig also requires ThumbnailConfig
			"content_config": pipelineOutputConfig(),

			"input_bucket": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if !regexp.MustCompile(`^[.0-9A-Za-z-_]+$`).MatchString(value) {
						errors = append(errors, fmt.Errorf(
							"only alphanumeric characters, hyphens, underscores, and periods allowed in %q", k))
					}
					if len(value) > 40 {
						errors = append(errors, fmt.Errorf("%q cannot be longer than 40 characters", k))
					}
					return
				},
			},

			"notifications": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"completed": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"error": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"progressing": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"warning": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			// One of output_bucket or content_config.bucket must be set
			"output_bucket": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"thumbnail_config": pipelineOutputConfig(),
		},
	}
}

func pipelineOutputConfig() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			// elastictranscoder.PipelineOutputConfing
			Schema: map[string]*schema.Schema{
				"bucket": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				"permissions": &schema.Schema{
					Type:     schema.TypeSet,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"access": &schema.Schema{
								Type:     schema.TypeList,
								Optional: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
							"grantee": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
							"grantee_type": &schema.Schema{
								Type:     schema.TypeString,
								Optional: true,
							},
						},
					},
				},
				"storage_class": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func resourceAwsElasticTranscoderPipelineCreate(d *schema.ResourceData, meta interface{}) error {
	elastictranscoderconn := meta.(*AWSClient).elastictranscoderconn

	req := &elastictranscoder.CreatePipelineInput{
		AwsKmsKeyArn:    getStringPtr(d, "aws_kms_key_arn"),
		ContentConfig:   expandETPiplineOutputConfig(d, "content_config"),
		InputBucket:     aws.String(d.Get("input_bucket").(string)),
		Name:            getStringPtr(d, "name"),
		Notifications:   expandETNotifications(d, "notifications"),
		OutputBucket:    getStringPtr(d, "output_bucket"),
		Role:            getStringPtr(d, "role"),
		ThumbnailConfig: expandETPiplineOutputConfig(d, "thumbnail_config"),
	}

	if (req.OutputBucket == nil && (req.ContentConfig == nil || req.ContentConfig.Bucket == nil)) ||
		(req.OutputBucket != nil && req.ContentConfig != nil && req.ContentConfig.Bucket != nil) {
		return fmt.Errorf("[ERROR] you must specifcy only one of output_bucket or content_config.bucket")
	}

	log.Printf("[DEBUG] Elastic Transcoder Pipeline create opts: %s", req)
	resp, err := elastictranscoderconn.CreatePipeline(req)
	if err != nil {
		return fmt.Errorf("Error creating Elastic Transcoder Pipeline: %s", err)
	}

	d.SetId(*resp.Pipeline.Id)

	for _, w := range resp.Warnings {
		log.Printf("[WARN] Elastic Transcoder Pipeline %s: %s", w.Code, w.Message)
	}

	return resourceAwsElasticTranscoderPipelineUpdate(d, meta)
}

func expandETNotifications(d *schema.ResourceData, key string) *elastictranscoder.Notifications {
	set, ok := d.GetOk(key)
	if !ok {
		return nil
	}

	s := set.(*schema.Set)
	if s == nil || s.Len() == 0 {
		return nil
	}

	m := s.List()[0].(map[string]interface{})

	return &elastictranscoder.Notifications{
		Completed:   getStringPtr(m, "completed"),
		Error:       getStringPtr(m, "error"),
		Progressing: getStringPtr(m, "progressing"),
		Warning:     getStringPtr(m, "warning"),
	}
}

func flattenETNotifications(n *elastictranscoder.Notifications) []map[string]interface{} {
	if n == nil {
		return nil
	}

	allEmpty := func(s ...*string) bool {
		for _, s := range s {
			if s != nil && *s != "" {
				return false
			}
		}
		return true
	}

	// the API always returns a Notifications value, even when all fields are nil
	if allEmpty(n.Completed, n.Error, n.Progressing, n.Warning) {
		return nil
	}

	m := make(map[string]interface{})

	if n.Completed != nil {
		m["completed"] = *n.Completed
	}

	if n.Error != nil {
		m["error"] = *n.Error
	}

	if n.Progressing != nil {
		m["progressing"] = *n.Progressing
	}

	if n.Warning != nil {
		m["warning"] = *n.Warning
	}

	return []map[string]interface{}{m}
}

func expandETPiplineOutputConfig(d *schema.ResourceData, key string) *elastictranscoder.PipelineOutputConfig {
	set, ok := d.GetOk(key)
	if !ok {
		return nil
	}

	s := set.(*schema.Set)
	if s == nil || s.Len() == 0 {
		return nil
	}

	cc := s.List()[0].(map[string]interface{})

	cfg := &elastictranscoder.PipelineOutputConfig{
		Bucket:       getStringPtr(cc, "bucket"),
		Permissions:  expandETPermList(cc["permissions"].(*schema.Set)),
		StorageClass: getStringPtr(cc, "storage_class"),
	}

	return cfg
}

func flattenETPipelineOutputConfig(cfg *elastictranscoder.PipelineOutputConfig) []map[string]interface{} {
	m := make(map[string]interface{})

	if cfg.Bucket != nil {
		m["bucket"] = *cfg.Bucket
	}

	if cfg.Permissions != nil {
		m["permissions"] = flattenETPermList(cfg.Permissions)
	}

	if cfg.StorageClass != nil {
		m["storage_class"] = *cfg.StorageClass
	}

	return []map[string]interface{}{m}
}

func expandETPermList(permissions *schema.Set) []*elastictranscoder.Permission {
	var perms []*elastictranscoder.Permission

	for _, p := range permissions.List() {
		m := p.(map[string]interface{})
		perm := &elastictranscoder.Permission{
			Access:      getStringPtrList(m, "access"),
			Grantee:     getStringPtr(m, "grantee"),
			GranteeType: getStringPtr(m, "grantee_type"),
		}
		perms = append(perms, perm)
	}
	return perms
}

func flattenETPermList(perms []*elastictranscoder.Permission) []map[string]interface{} {
	var set []map[string]interface{}

	for _, p := range perms {
		m := make(map[string]interface{})
		if p.Access != nil {
			m["access"] = flattenStringList(p.Access)
		}

		if p.Grantee != nil {
			m["grantee"] = *p.Grantee
		}

		if p.GranteeType != nil {
			m["grantee_type"] = *p.GranteeType
		}

		set = append(set, m)
	}
	return set
}

func resourceAwsElasticTranscoderPipelineUpdate(d *schema.ResourceData, meta interface{}) error {
	elastictranscoderconn := meta.(*AWSClient).elastictranscoderconn

	req := &elastictranscoder.UpdatePipelineInput{
		Id: aws.String(d.Id()),
	}

	if d.HasChange("aws_kms_key_arn") {
		req.AwsKmsKeyArn = getStringPtr(d, "aws_kms_key_arn")
	}

	if d.HasChange("content_config") {
		req.ContentConfig = expandETPiplineOutputConfig(d, "content_config")
	}

	if d.HasChange("input_bucket") {
		req.Role = getStringPtr(d, "input_bucket")
	}

	if d.HasChange("name") {
		req.Name = getStringPtr(d, "name")
	}

	if d.HasChange("notifications") {
		req.Notifications = expandETNotifications(d, "notifications")
	}

	if d.HasChange("role") {
		req.Role = getStringPtr(d, "role")
	}

	if d.HasChange("thumbnail_config") {
		req.ThumbnailConfig = expandETPiplineOutputConfig(d, "thumbnail_config")
	}

	log.Printf("[DEBUG] Updating Elastic Transcoder Pipeline: %#v", req)
	output, err := elastictranscoderconn.UpdatePipeline(req)
	if err != nil {
		return fmt.Errorf("Error updating Elastic Transcoder pipeline: %s", err)
	}

	for _, w := range output.Warnings {
		log.Printf("[WARN] Elastic Transcoder Pipeline %s: %s", w.Code, w.Message)
	}

	return resourceAwsElasticTranscoderPipelineRead(d, meta)
}

func resourceAwsElasticTranscoderPipelineRead(d *schema.ResourceData, meta interface{}) error {
	elastictranscoderconn := meta.(*AWSClient).elastictranscoderconn

	resp, err := elastictranscoderconn.ReadPipeline(&elastictranscoder.ReadPipelineInput{
		Id: aws.String(d.Id()),
	})

	if err != nil {
		if err, ok := err.(awserr.Error); ok && err.Code() == "ResourceNotFoundException" {
			d.SetId("")
			return nil
		}
		return err
	}

	log.Printf("[DEBUG] Elastic Transcoder Pipeline Read response: %#v", resp)

	pipeline := resp.Pipeline

	d.Set("arn", *pipeline.Arn)

	if arn := pipeline.AwsKmsKeyArn; arn != nil {
		d.Set("aws_kms_key_arn", *arn)
	}

	if pipeline.ContentConfig != nil {
		d.Set("content_config", flattenETPipelineOutputConfig(pipeline.ContentConfig))
	}

	d.Set("input_bucket", *pipeline.InputBucket)
	d.Set("name", *pipeline.Name)

	notifications := flattenETNotifications(pipeline.Notifications)
	if notifications != nil {
		d.Set("notifications", notifications)
	}

	if pipeline.OutputBucket != nil {
		d.Set("output_bucket", *pipeline.OutputBucket)
	}

	d.Set("role", pipeline.Role)

	if pipeline.ThumbnailConfig != nil {
		d.Set("thumbnail_config", flattenETPipelineOutputConfig(pipeline.ThumbnailConfig))
	}

	return nil
}

func resourceAwsElasticTranscoderPipelineDelete(d *schema.ResourceData, meta interface{}) error {
	elastictranscoderconn := meta.(*AWSClient).elastictranscoderconn

	log.Printf("[DEBUG] Elastic Transcoder Delete Pipeline: %s", d.Id())
	_, err := elastictranscoderconn.DeletePipeline(&elastictranscoder.DeletePipelineInput{
		Id: aws.String(d.Id()),
	})
	if err != nil {
		return err
	}
	return nil
}

// getNilString returns a *string version of the value taken from m, where m
// can be a map[string]interface{} or a *schema.ResourceData. If the key isn't
// present, getNilString returns nil.
func getStringPtr(m interface{}, key string) *string {
	switch m := m.(type) {
	case map[string]interface{}:
		if v, ok := m[key]; ok {
			s := v.(string)
			return &s
		}
	case *schema.ResourceData:
		if v, ok := m.GetOk(key); ok {
			s := v.(string)
			return &s
		}
	}
	return nil
}

// getNilStringList returns a []*string version of the map value. If the key
// isn't present, getNilStringList returns nil.
func getStringPtrList(m map[string]interface{}, key string) []*string {
	if v, ok := m[key]; ok {
		var stringList []*string
		for _, i := range v.([]interface{}) {
			s := i.(string)
			stringList = append(stringList, &s)
		}

		return stringList
	}
	return nil
}
