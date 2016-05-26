package aws

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
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

	req := &elastictranscoder.CreatePipelineInput{}

	if key, ok := d.GetOk("aws_kms_key_arn"); ok {
		req.AwsKmsKeyArn = aws.String(key.(string))
	}

	if cc, ok := d.GetOk("content_config"); ok {
		req.ContentConfig = expandETPiplineOutputConfig(cc.(*schema.Set))
	}

	req.InputBucket = aws.String(d.Get("input_bucket").(string))

	req.Name = aws.String(d.Get("name").(string))

	if n, ok := d.GetOk("notifications"); ok {
		req.Notifications = expandETNotifications(n.(*schema.Set))
	}

	if ob, ok := d.GetOk("output_bucket"); ok {
		req.OutputBucket = aws.String(ob.(string))
	}

	req.Role = aws.String(d.Get("role").(string))

	if tc, ok := d.GetOk("thumbnail_config"); ok {
		req.ThumbnailConfig = expandETPiplineOutputConfig(tc.(*schema.Set))
	}

	if (req.OutputBucket == nil && (req.ContentConfig == nil || req.ContentConfig.Bucket == nil)) ||
		(req.OutputBucket != nil && req.ContentConfig != nil && req.ContentConfig.Bucket != nil) {
		return fmt.Errorf("[ERROR] you must specifcy only one of output_bucket or content_config.bucket")
	}

	log.Printf("[DEBUG] Elastic Transcoder Pipeline create opts: %s", req)
	resp, err := elastictranscoderconn.CreatePipeline(req)
	if err != nil {
		return fmt.Errorf("Error creating Elastic Transcoder pipeline: %s", err)
	}

	d.SetId(*resp.Pipeline.Id)

	for _, w := range resp.Warnings {
		log.Printf("[WARN] Elastic Transcoder Pipeline %s: %s", w.Code, w.Message)
	}

	return resourceAwsElasticTranscoderPipelineUpdate(d, meta)
}

func expandETNotifications(s *schema.Set) *elastictranscoder.Notifications {
	notifications := &elastictranscoder.Notifications{}

	if s.Len() == 0 {
		return nil
	}

	n := s.List()[0].(map[string]interface{})

	if c, ok := n["completed"]; ok {
		notifications.Completed = aws.String(c.(string))
	}

	if e, ok := n["error"]; ok {
		notifications.Error = aws.String(e.(string))
	}

	if p, ok := n["progressing"]; ok {
		notifications.Progressing = aws.String(p.(string))
	}

	if w, ok := n["warning"]; ok {
		notifications.Warning = aws.String(w.(string))
	}

	return notifications
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

func expandETPiplineOutputConfig(s *schema.Set) *elastictranscoder.PipelineOutputConfig {
	cfg := &elastictranscoder.PipelineOutputConfig{}

	if s.Len() == 0 {
		return nil
	}

	cc := s.List()[0].(map[string]interface{})

	if bucket, ok := cc["bucket"]; ok {
		cfg.Bucket = aws.String(bucket.(string))
	}

	if perms, ok := cc["permissions"]; ok {
		cfg.Permissions = expandETPermList(perms.(*schema.Set))
	}

	if sc, ok := cc["storage_class"]; ok {
		cfg.StorageClass = aws.String(sc.(string))
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
		permMap := p.(map[string]interface{})
		perm := &elastictranscoder.Permission{}

		if a, ok := permMap["access"]; ok {
			perm.Access = expandStringList(a.([]interface{}))
		}

		if g, ok := permMap["grantee"]; ok {
			perm.Grantee = aws.String(g.(string))
		}

		if gt, ok := permMap["grantee_type"]; ok {
			perm.GranteeType = aws.String(gt.(string))
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
		// not required, check for empty string
		key := d.Get("aws_kms_key_arn").(string)
		if key != "" {
			req.AwsKmsKeyArn = aws.String(key)
		}
	}

	if d.HasChange("content_config") {
		req.ContentConfig = expandETPiplineOutputConfig(d.Get("content_config").(*schema.Set))
	}

	if d.HasChange("input_bucket") {
		req.Role = aws.String(d.Get("input_bucket").(string))
	}

	if d.HasChange("name") {
		req.Name = aws.String(d.Get("name").(string))
	}

	if d.HasChange("notifications") {
		req.Notifications = expandETNotifications(d.Get("notifications").(*schema.Set))
	}

	if d.HasChange("role") {
		req.Role = aws.String(d.Get("role").(string))
	}

	if d.HasChange("thumbnail_config") {
		req.ThumbnailConfig = expandETPiplineOutputConfig(d.Get("thumbnail_config").(*schema.Set))
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
		// FIXME: check for NOT FOUND
		d.SetId("")
		return nil
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
		log.Println("[XXXX] NOTIFICATIONS", notifications)
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
