package confanalysis

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RepositoryConfAnalysisData struct {
	Config UserFacingConfig `json:"userConfig"`
}

type UserFacingConfig struct {
	AlertOnViolation           bool     `json:"alertOnViolation"`
	BlockOnViolation           bool     `json:"blockOnViolation"`
	CommentAnnotationGroups    []string `json:"commentAnnotationGroups,omitempty"`
	DisableFilterAnalysis      bool     `json:"disableFilterAnalysis"`
	DisablePreConfiguredAlerts bool     `json:"disablePreConfiguredAlerts"`
	EnableDataMasking          bool     `json:"enableDataMasking"`
	LogGroups                  []string `json:"logGroups,omitempty"`
	Redact                     string   `json:"redact"`
	EnableDatasetRewrites      bool     `json:"enableDatasetRewrites"`
}

func repositoryConfAnalysisResourceSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"repository_id": {
				Description: "The ID of an existing data repository resource that will be configured.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"redact": {
				Description:  "Valid values are: `all`, `none` and `watched`. If set to `all` it will enable the redact of all literal values, `none` will disable it, and `watched` will only redact values from tracked fields set in the Datamap.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "all",
				ValidateFunc: client.ValidateRepositoryConfAnalysisRedact(),
			},
			"alert_on_violation": {
				Description: "If set to `true` it will enable alert on policy violations.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"disable_pre_configured_alerts": {
				Description: "If set to `true` it will *disable* preconfigured alerts.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"enable_data_masking": {
				Description: "If set to `true` it will allow policies to force the masking " +
					" of specified data fields in the results of queries. " +
					"[Learn more](https://cyral.com/docs/using-cyral/masking/).",
				Type:     schema.TypeBool,
				Optional: true,
			},
			"block_on_violation": {
				Description: "If set to `true` it will enable query blocking in case of a " +
					"policy violation.",
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disable_filter_analysis": {
				Description: "If set to `true` it will *disable* filter analysis.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"enable_dataset_rewrites": {
				Description: "If set to `true` it will enable rewriting queries.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"comment_annotation_groups": {
				Description: "Valid values are: `identity`, `client`, `repo`, `sidecar`. The " +
					"default behavior is to set only the `identity` when this option is " +
					"enabled, but you can also opt to add the contents of `client`, `repo`, " +
					" `sidecar` logging blocks as query comments. " +
					" [Learn more](https://support.cyral.com/support/solutions/articles/44002218978).",
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: client.ValidateRepositoryConfAnalysisCommentAnnotationGroups(),
				},
			},
			"log_groups": {
				Description: "Responsible for configuring the Log Settings. Valid values are documented below. The `log_groups` list support the following values: " +
					"\n  - `everything` - Enables all the Log Settings." +
					"\n  - `dql` - Enables the `DQLs` setting for `all requests`." +
					"\n  - `dml` - Enables the `DMLs` setting for `all requests`." +
					"\n  - `ddl` - Enables the `DDLs` setting for `all requests`." +
					"\n  - `sensitive & dql` - Enables the `DQLs` setting for `logged fields`." +
					"\n  - `sensitive & dml` - Enables the `DMLs` setting for `logged fields`." +
					"\n  - `sensitive & ddl` - Enables the `DDLs` setting for `logged fields`." +
					"\n  - `privileged` - Enables the `Privileged commands` setting." +
					"\n  - `port-scan` - Enables the `Port scans` setting." +
					"\n  - `auth-failure` - Enables the `Authentication failures` setting." +
					"\n  - `full-table-scan` - Enables the `Full scans` setting." +
					"\n  - `violations` - Enables the `Policy violations` setting." +
					"\n  - `connections` - Enables the `Connection activity` setting." +
					"\n  - `sensitive` - Log all queries manipulating sensitive fields (watches)" +
					"\n  - `data-classification` - Log all queries whose response was automatically classified as sensitive (credit card numbers, emails and so on)." +
					"\n  - `audit` - Log `sensitive`, `DQLs`, `DDLs`, `DMLs` and `privileged`." +
					"\n  - `error` - Log analysis errors." +
					"\n  - `new-connections` - Log new connections." +
					"\n  - `closed-connections` - Log closed connections.",
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: client.ValidateRepositoryConfAnalysisLogSettings(),
				},
			},
		},
	}
}

// Previously, the ID for cyral_repository_conf_analysis had the format
// {repository_id}/ConfAnalysis. The goal of this state upgrade is to remove
// this suffix `ConfAnalysis`.
func UpgradeRepositoryConfAnalysisV0(
	_ context.Context,
	rawState map[string]interface{},
	_ interface{},
) (map[string]interface{}, error) {
	rawState["id"] = rawState["repository_id"]
	return rawState, nil
}

func ResourceRepositoryConfAnalysis() *schema.Resource {
	return &schema.Resource{
		Description: "Manages Repository Analysis Configuration. This resource allows configuring both " +
			"[Log Settings](https://cyral.com/docs/manage-repositories/repo-log-volume) " +
			"and [Advanced settings](https://cyral.com/docs/manage-repositories/repo-advanced-settings) " +
			"(Logs, Alerts, Analysis and Enforcement) configurations for Data Repositories.",
		CreateContext: resourceRepositoryConfAnalysisCreate,
		ReadContext:   resourceRepositoryConfAnalysisRead,
		UpdateContext: resourceRepositoryConfAnalysisUpdate,
		DeleteContext: resourceRepositoryConfAnalysisDelete,

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type: repositoryConfAnalysisResourceSchemaV0().
					CoreConfigSchema().ImpliedType(),
				Upgrade: UpgradeRepositoryConfAnalysisV0,
			},
		},

		Schema: repositoryConfAnalysisResourceSchemaV0().Schema,

		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m interface{},
			) ([]*schema.ResourceData, error) {
				d.Set("repository_id", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceRepositoryConfAnalysisCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Init resourceRepositoryConfAnalysisCreate")
	c := m.(*client.Client)

	resourceData, err := getConfAnalysisDataFromResource(d)
	if err != nil {
		return utils.CreateError("Unable to create conf analysis", fmt.Sprintf("%v", err))
	}

	url := fmt.Sprintf("https://%s/v1/repos/%s/conf/analysis", c.ControlPlane, d.Get("repository_id"))

	body, err := c.DoRequest(ctx, url, http.MethodPut, resourceData.Config)
	if err != nil {
		return utils.CreateError("Unable to create conf analysis", fmt.Sprintf("%v", err))
	}

	response := RepositoryConfAnalysisData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return utils.CreateError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}

	tflog.Debug(ctx, fmt.Sprintf("Response body (unmarshalled): %#v", response))

	d.SetId(d.Get("repository_id").(string))

	tflog.Debug(ctx, "End resourceRepositoryConfAnalysisCreate")

	return resourceRepositoryConfAnalysisRead(ctx, d, m)
}

func resourceRepositoryConfAnalysisRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Init resourceRepositoryConfAnalysisRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/repos/%s/conf/analysis", c.ControlPlane, d.Get("repository_id"))

	body, err := c.DoRequest(ctx, url, http.MethodGet, nil)
	if err != nil {
		return utils.CreateError(fmt.Sprintf("Unable to read conf analysis. Conf Analysis Id: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	response := RepositoryConfAnalysisData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return utils.CreateError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}

	tflog.Debug(ctx, fmt.Sprintf("Response body (unmarshalled): %#v", response))

	setConfAnalysisDataToResource(d, response)

	tflog.Debug(ctx, "End resourceRepositoryConfAnalysisRead")

	return diag.Diagnostics{}
}

func resourceRepositoryConfAnalysisUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Init resourceRepositoryConfAnalysisUpdate")
	c := m.(*client.Client)

	resourceData, err := getConfAnalysisDataFromResource(d)
	if err != nil {
		return utils.CreateError("Unable to update conf analysis", fmt.Sprintf("%v", err))
	}

	url := fmt.Sprintf("https://%s/v1/repos/%s/conf/analysis", c.ControlPlane, d.Get("repository_id"))

	if _, err := c.DoRequest(ctx, url, http.MethodPut, resourceData.Config); err != nil {
		return utils.CreateError("Unable to update conf analysis", fmt.Sprintf("%v", err))
	}

	tflog.Debug(ctx, "End resourceRepositoryConfAnalysisUpdate")

	return resourceRepositoryConfAnalysisRead(ctx, d, m)
}

func resourceRepositoryConfAnalysisDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Init resourceRepositoryConfAnalysisDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/repos/%s/conf/analysis", c.ControlPlane, d.Get("repository_id"))

	if _, err := c.DoRequest(ctx, url, http.MethodDelete, nil); err != nil {
		return utils.CreateError("Unable to delete conf analysis", fmt.Sprintf("%v", err))
	}

	tflog.Debug(ctx, "End resourceRepositoryConfAnalysisDelete")

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
			AlertOnViolation:           d.Get("alert_on_violation").(bool),
			BlockOnViolation:           d.Get("block_on_violation").(bool),
			DisableFilterAnalysis:      d.Get("disable_filter_analysis").(bool),
			DisablePreConfiguredAlerts: d.Get("disable_pre_configured_alerts").(bool),
			EnableDataMasking:          d.Get("enable_data_masking").(bool),
			CommentAnnotationGroups:    annotationGroups,
			LogGroups:                  logGroups,
			Redact:                     d.Get("redact").(string),
			EnableDatasetRewrites:      d.Get("enable_dataset_rewrites").(bool),
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

	d.Set("alert_on_violation", resourceData.Config.AlertOnViolation)
	d.Set("block_on_violation", resourceData.Config.BlockOnViolation)
	d.Set("comment_annotation_groups", annotationGroupsSet)
	d.Set("disable_filter_analysis", resourceData.Config.DisableFilterAnalysis)
	d.Set("disable_pre_configured_alerts", resourceData.Config.DisablePreConfiguredAlerts)
	d.Set("enable_data_masking", resourceData.Config.EnableDataMasking)
	d.Set("log_groups", logGroupsSet)
	d.Set("redact", resourceData.Config.Redact)
	d.Set("enable_dataset_rewrites", resourceData.Config.EnableDatasetRewrites)
}
