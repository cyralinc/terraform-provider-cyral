package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RepositoryConfAnalysisData struct {
	Config UserFacingConfig `json:"userConfig"`
}

type UserFacingConfig struct {
	Redact                     string   `json:"redact"`
	AlertOnViolation           bool     `json:"alertOnViolation"`
	DisablePreConfiguredAlerts bool     `json:"disablePreConfiguredAlerts"`
	BlockOnViolation           bool     `json:"blockOnViolation"`
	DisableFilterAnalysis      bool     `json:"disableFilterAnalysis"`
	RewriteOnViolation         bool     `json:"rewriteOnViolation"`
	CommentAnnotationGroups    []string `json:"commentAnnotationGroups,omitempty"`
	LogGroups                  []string `json:"logGroups,omitempty"`
}

func resourceRepositoryConfAnalysis() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRepositoryConfAnalysisCreate,
		ReadContext:   resourceRepositoryConfAnalysisRead,
		UpdateContext: resourceRepositoryConfAnalysisUpdate,
		DeleteContext: resourceRepositoryConfAnalysisDelete,

		Schema: map[string]*schema.Schema{
			"repository_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"redact": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "all",
				ValidateFunc: client.ValidateRepositoryConfAnalysisRedact(),
			},
			"alert_on_violation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"disable_pre_configured_alerts": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"block_on_violation": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disable_filter_analysis": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"rewrite_on_violation": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"comment_annotation_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: client.ValidateRepositoryConfAnalysisCommentAnnotationGroups(),
				},
			},
			"log_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: client.ValidateRepositoryConfAnalysisLogSettings(),
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceRepositoryConfAnalysisCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryConfAnalysisCreate")
	c := m.(*client.Client)

	resourceData, err := getConfAnalysisDataFromResource(d)
	if err != nil {
		return createError("Unable to create conf analysis", fmt.Sprintf("%v", err))
	}

	url := fmt.Sprintf("https://%s/v1/repos/%s/conf/analysis", c.ControlPlane, d.Get("repository_id"))

	body, err := c.DoRequest(url, http.MethodPut, resourceData.Config)
	if err != nil {
		return createError("Unable to create conf analysis", fmt.Sprintf("%v", err))
	}

	response := RepositoryConfAnalysisData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(fmt.Sprintf("%s/ConfAnalysis", d.Get("repository_id")))

	log.Printf("[DEBUG] End resourceRepositoryConfAnalysisCreate")

	return resourceRepositoryConfAnalysisRead(ctx, d, m)
}

func resourceRepositoryConfAnalysisRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryConfAnalysisRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/repos/%s/conf/analysis", c.ControlPlane, d.Get("repository_id"))

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError(fmt.Sprintf("Unable to read conf analysis. Conf Analysis Id: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	response := RepositoryConfAnalysisData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	setConfAnalysisDataToResource(d, response)

	log.Printf("[DEBUG] End resourceRepositoryConfAnalysisRead")

	return diag.Diagnostics{}
}

func resourceRepositoryConfAnalysisUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryConfAnalysisUpdate")
	c := m.(*client.Client)

	resourceData, err := getConfAnalysisDataFromResource(d)
	if err != nil {
		return createError("Unable to update conf analysis", fmt.Sprintf("%v", err))
	}

	url := fmt.Sprintf("https://%s/v1/repos/%s/conf/analysis", c.ControlPlane, d.Get("repository_id"))

	if _, err := c.DoRequest(url, http.MethodPut, resourceData.Config); err != nil {
		return createError("Unable to update conf analysis", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceRepositoryConfAnalysisUpdate")

	return resourceRepositoryConfAnalysisRead(ctx, d, m)
}

func resourceRepositoryConfAnalysisDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryConfAnalysisDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/repos/%s/conf/analysis", c.ControlPlane, d.Get("repository_id"))

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return createError("Unable to delete conf analysis", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceRepositoryConfAnalysisDelete")

	return diag.Diagnostics{}
}

func getConfAnalysisDataFromResource(d *schema.ResourceData) (RepositoryConfAnalysisData, error) {
	var logGroups []string
	if logGroupsSet, ok := d.GetOk("log_groups"); ok {
		for _, logGroupItem := range logGroupsSet.(*schema.Set).List() {
			logGroups = append(logGroups, logGroupItem.(string))
		}
	}

	var annotationGroups []string
	if annotationGroupsSet, ok := d.GetOk("comment_annotation_groups"); ok {
		for _, annotationGroupItem := range annotationGroupsSet.(*schema.Set).List() {
			annotationGroups = append(annotationGroups, annotationGroupItem.(string))
		}
	}

	return RepositoryConfAnalysisData{
		Config: UserFacingConfig{
			Redact:                     d.Get("redact").(string),
			AlertOnViolation:           d.Get("alert_on_violation").(bool),
			DisablePreConfiguredAlerts: d.Get("disable_pre_configured_alerts").(bool),
			BlockOnViolation:           d.Get("block_on_violation").(bool),
			DisableFilterAnalysis:      d.Get("disable_filter_analysis").(bool),
			RewriteOnViolation:         d.Get("rewrite_on_violation").(bool),
			CommentAnnotationGroups:    annotationGroups,
			LogGroups:                  logGroups,
		},
	}, nil
}

func setConfAnalysisDataToResource(d *schema.ResourceData, resourceData RepositoryConfAnalysisData) {
	logGroups := make([]interface{}, len(resourceData.Config.LogGroups))
	for index, logGroupItem := range resourceData.Config.LogGroups {
		logGroups[index] = logGroupItem
	}
	logGroupsSet := schema.NewSet(schema.HashString, logGroups)

	annotationGroups := make([]interface{}, len(resourceData.Config.CommentAnnotationGroups))
	for index, annotationGroupItem := range resourceData.Config.CommentAnnotationGroups {
		annotationGroups[index] = annotationGroupItem
	}
	annotationGroupsSet := schema.NewSet(schema.HashString, annotationGroups)

	d.Set("redact", resourceData.Config.Redact)
	d.Set("alert_on_violation", resourceData.Config.AlertOnViolation)
	d.Set("disable_pre_configured_alerts", resourceData.Config.DisablePreConfiguredAlerts)
	d.Set("block_on_violation", resourceData.Config.BlockOnViolation)
	d.Set("disable_filter_analysis", resourceData.Config.DisableFilterAnalysis)
	d.Set("rewrite_on_violation", resourceData.Config.RewriteOnViolation)
	d.Set("comment_annotation_groups", annotationGroupsSet)
	d.Set("log_groups", logGroupsSet)
}
