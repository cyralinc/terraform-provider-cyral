package confanalysis

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO: v2 of this API should either return the repository ID
// or something else as an ID. Currently it accepts a `UserConfig`
// for the PUT payload, but returns a `RepositoryConfAnalysisData`.
// This makes the whole API utilization quite confusing.

type RepositoryConfAnalysisData struct {
	UserConfig UserConfig `json:"userConfig"`
}

type UserConfig struct {
	AlertOnViolation           bool     `json:"alertOnViolation"`
	BlockOnViolation           bool     `json:"blockOnViolation"`
	CommentAnnotationGroups    []string `json:"commentAnnotationGroups,omitempty"`
	DisableFilterAnalysis      bool     `json:"disableFilterAnalysis"`
	DisablePreConfiguredAlerts bool     `json:"disablePreConfiguredAlerts"`
	EnableDataMasking          bool     `json:"enableDataMasking"`
	MaskAllOccurrences         bool     `json:"maskAllOccurrences"`
	LogGroups                  []string `json:"logGroups,omitempty"`
	Redact                     string   `json:"redact"`
	EnableDatasetRewrites      bool     `json:"enableDatasetRewrites"`
}

func (r *RepositoryConfAnalysisData) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(d.Get("repository_id").(string))
	return r.UserConfig.WriteToSchema(d)
}

func (r *UserConfig) WriteToSchema(d *schema.ResourceData) error {
	logGroups := make([]interface{}, len(r.LogGroups))
	for index, logGroupItem := range r.LogGroups {
		logGroups[index] = logGroupItem
	}
	logGroupsSet := schema.NewSet(schema.HashString, logGroups)

	annotationGroups := make([]interface{}, len(r.CommentAnnotationGroups))
	for index, annotationGroupItem := range r.CommentAnnotationGroups {
		annotationGroups[index] = annotationGroupItem
	}
	annotationGroupsSet := schema.NewSet(schema.HashString, annotationGroups)

	d.Set("alert_on_violation", r.AlertOnViolation)
	d.Set("block_on_violation", r.BlockOnViolation)
	d.Set("comment_annotation_groups", annotationGroupsSet)
	d.Set("disable_filter_analysis", r.DisableFilterAnalysis)
	d.Set("disable_pre_configured_alerts", r.DisablePreConfiguredAlerts)
	d.Set("enable_data_masking", r.EnableDataMasking)
	d.Set("mask_all_occurrences", r.MaskAllOccurrences)
	d.Set("log_groups", logGroupsSet)
	d.Set("redact", r.Redact)
	d.Set("enable_dataset_rewrites", r.EnableDatasetRewrites)

	return nil
}

func (r *RepositoryConfAnalysisData) ReadFromSchema(d *schema.ResourceData) error {
	return r.UserConfig.ReadFromSchema(d)
}

func (r *UserConfig) ReadFromSchema(d *schema.ResourceData) error {
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

	r.AlertOnViolation = d.Get("alert_on_violation").(bool)
	r.BlockOnViolation = d.Get("block_on_violation").(bool)
	r.DisableFilterAnalysis = d.Get("disable_filter_analysis").(bool)
	r.DisablePreConfiguredAlerts = d.Get("disable_pre_configured_alerts").(bool)
	r.EnableDataMasking = d.Get("enable_data_masking").(bool)
	r.MaskAllOccurrences = d.Get("mask_all_occurrences").(bool)
	r.CommentAnnotationGroups = annotationGroups
	r.LogGroups = logGroups
	r.Redact = d.Get("redact").(string)
	r.EnableDatasetRewrites = d.Get("enable_dataset_rewrites").(bool)

	return nil
}
