package sidecar

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CreateSidecarResponse struct {
	ID string `json:"ID"`
}

type SidecarData struct {
	ID                       string                   `json:"id"`
	Name                     string                   `json:"name"`
	Labels                   []string                 `json:"labels"`
	SidecarProperties        *SidecarProperties       `json:"properties"`
	ServicesConfig           SidecarServicesConfig    `json:"services"`
	UserEndpoint             string                   `json:"userEndpoint"`
	CertificateBundleSecrets CertificateBundleSecrets `json:"certificateBundleSecrets,omitempty"`
}

func (sd *SidecarData) BypassMode() string {
	if sd.ServicesConfig != nil {
		if dispConfig, ok := sd.ServicesConfig["dispatcher"]; ok {
			if bypass_mode, ok := dispConfig["bypass"]; ok {
				return bypass_mode
			}
		}
	}
	return ""
}

func (r *SidecarData) WriteToSchema(d *schema.ResourceData) error {
	if err := d.Set("name", r.Name); err != nil {
		return fmt.Errorf("error setting 'name' field: %w", err)
	}
	if r.SidecarProperties != nil {
		if err := d.Set("deployment_method", r.SidecarProperties.DeploymentMethod); err != nil {
			return fmt.Errorf("error setting 'deployment_method' field: %w", err)
		}
		if err := d.Set("activity_log_integration_id", r.SidecarProperties.LogIntegrationID); err != nil {
			return fmt.Errorf("error setting 'activity_log_integration_id' field: %w", err)
		}
		if err := d.Set("diagnostic_log_integration_id", r.SidecarProperties.DiagnosticLogIntegrationID); err != nil {
			return fmt.Errorf("error setting 'diagnostic_log_integration_id' field: %w", err)
		}
	}
	if err := d.Set("labels", r.Labels); err != nil {
		return fmt.Errorf("error setting 'labels' field: %w", err)
	}
	if err := d.Set("user_endpoint", r.UserEndpoint); err != nil {
		return fmt.Errorf("error setting 'user_endpoint' field: %w", err)
	}
	if bypassMode := r.BypassMode(); bypassMode != "" {
		if err := d.Set("bypass_mode", bypassMode); err != nil {
			return fmt.Errorf("error setting 'bypass_mode' field: %w", err)
		}
	}

	if err := d.Set("certificate_bundle_secrets", flattenCertificateBundleSecrets(r.CertificateBundleSecrets)); err != nil {
		return fmt.Errorf("error setting 'certificate_bundle_secrets' field: %w", err)
	}
	return nil
}

func (r *SidecarData) ReadFromSchema(d *schema.ResourceData) error {
	activityLogIntegrationID := d.Get("activity_log_integration_id").(string)
	if activityLogIntegrationID == "" {
		activityLogIntegrationID = d.Get("log_integration_id").(string)
	}

	labels := d.Get("labels").([]interface{})
	sidecarDataLabels := []string{}
	for _, labelInterface := range labels {
		if label, ok := labelInterface.(string); ok {
			sidecarDataLabels = append(sidecarDataLabels, label)
		}
	}

	r.ID = d.Id()
	r.Name = d.Get("name").(string)
	r.Labels = sidecarDataLabels
	r.SidecarProperties = &SidecarProperties{
		DeploymentMethod:           d.Get("deployment_method").(string),
		LogIntegrationID:           activityLogIntegrationID,
		DiagnosticLogIntegrationID: d.Get("diagnostic_log_integration_id").(string),
	}
	r.ServicesConfig = SidecarServicesConfig{
		"dispatcher": map[string]string{
			"bypass": d.Get("bypass_mode").(string),
		},
	}
	r.UserEndpoint = d.Get("user_endpoint").(string)
	r.CertificateBundleSecrets = getCertificateBundleSecret(d)

	return nil
}

type SidecarProperties struct {
	DeploymentMethod           string `json:"deploymentMethod"`
	LogIntegrationID           string `json:"logIntegrationID,omitempty"`
	DiagnosticLogIntegrationID string `json:"diagnosticLogIntegrationID,omitempty"`
}

type SidecarServicesConfig map[string]map[string]string

type CertificateBundleSecrets map[string]*CertificateBundleSecret

type CertificateBundleSecret struct {
	Engine   string `json:"engine,omitempty"`
	SecretId string `json:"secretId,omitempty"`
	Type     string `json:"type,omitempty"`
}

func flattenCertificateBundleSecrets(cbs CertificateBundleSecrets) []interface{} {
	ctx := context.Background()
	tflog.Debug(ctx, "Init flattenCertificateBundleSecrets")
	var flatCBS []interface{}
	if cbs != nil {
		cb := make(map[string]interface{})

		for key, val := range cbs {
			// Ignore self-signed certificates
			if key != "sidecar-generated-selfsigned" {
				contentCB := make([]interface{}, 1)

				tflog.Debug(ctx, fmt.Sprintf("key: %v", key))
				tflog.Debug(ctx, fmt.Sprintf("val: %v", val))

				contentCBMap := make(map[string]interface{})
				contentCBMap["secret_id"] = val.SecretId
				contentCBMap["engine"] = val.Engine
				contentCBMap["type"] = val.Type

				contentCB[0] = contentCBMap
				cb[key] = contentCB
			}
		}

		if len(cb) > 0 {
			flatCBS = make([]interface{}, 1)
			flatCBS[0] = cb
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("end flattenCertificateBundleSecrets %v", flatCBS))
	return flatCBS
}

func getCertificateBundleSecret(d *schema.ResourceData) CertificateBundleSecrets {
	ctx := context.Background()
	tflog.Debug(ctx, "Init getCertificateBundleSecret")
	rdCBS := d.Get("certificate_bundle_secrets").(*schema.Set).List()
	ret := make(CertificateBundleSecrets)

	if len(rdCBS) > 0 {
		cbsMap := rdCBS[0].(map[string]interface{})
		for k, v := range cbsMap {
			vList := v.(*schema.Set).List()
			// 1. k = "sidecar" or other direct internal elements of certificate_bundle_secrets
			// 2. Also one element on this list due to MaxItems...
			// 3. Ignore self signed certificates
			if len(vList) > 0 && k != "sidecar-generated-selfsigned" {
				vMap := vList[0].(map[string]interface{})
				engine := ""
				if val, ok := vMap["engine"]; val != nil && ok {
					engine = val.(string)
				}
				cbsType := vMap["type"].(string)
				secretId := vMap["secret_id"].(string)
				cbs := CertificateBundleSecret{
					SecretId: secretId,
					Engine:   engine,
					Type:     cbsType,
				}
				ret[k] = &cbs
			}
		}
	}

	// If the occurrence of `sidecar` does not exist, set it to an empty certificate bundle
	// so that the API can remove the `sidecar` key from the persisted certificate bundle map.
	if _, ok := ret["sidecar"]; !ok {
		ret["sidecar"] = &CertificateBundleSecret{}
	}

	tflog.Debug(ctx, "end getCertificateBundleSecret")
	return ret
}
