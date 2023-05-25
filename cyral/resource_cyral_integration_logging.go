package cyral

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

const loggingApiUrl = "https://%s/v1/integrations/logging/%s"

func getLoggingConfig(resource *LoggingIntegration) ([]interface{}, error) {
	var configScheme []interface{}
	switch {
	case resource.CloudWatch != nil:
		configScheme = []interface{}{
			map[string]interface{}{
				"cloud_watch": []interface{}{
					map[string]interface{}{
						"region":             resource.CloudWatch.Region,
						"group":              resource.CloudWatch.Group,
						"stream":             resource.CloudWatch.Stream,
						"log_retention_days": resource.CloudWatch.LogRetentionDays,
					},
				},
			},
		}
	case resource.Datadog != nil:
		configScheme = []interface{}{
			map[string]interface{}{
				"datadog": []interface{}{
					map[string]interface{}{
						"api_key": resource.Datadog.ApiKey,
					},
				},
			},
		}
	case resource.Elk != nil:
		configScheme = []interface{}{
			map[string]interface{}{
				"elk": []interface{}{
					map[string]interface{}{
						"es_url":     resource.Elk.EsURL,
						"kibana_url": resource.Elk.KibanaURL,
						"es_credentials": []interface{}{
							map[string]interface{}{
								"username": resource.Elk.EsCredentials.Username,
								"password": resource.Elk.EsCredentials.Password,
							},
						},
					},
				},
			},
		}
	case resource.Splunk != nil:
		configScheme = []interface{}{
			map[string]interface{}{
				"splunk": []interface{}{
					map[string]interface{}{
						"hostname":     resource.Splunk.Hostname,
						"hec_port":     resource.Splunk.HecPort,
						"access_token": resource.Splunk.AccessToken,
						"index":        resource.Splunk.Index,
						"use_tls":      resource.Splunk.UseTLS,
					},
				},
			},
		}
	case resource.SumoLogic != nil:
		configScheme = []interface{}{
			map[string]interface{}{
				"sumo_logic": []interface{}{
					map[string]interface{}{
						"address": resource.SumoLogic.Address,
					},
				},
			},
		}
	case resource.FluentBit != nil:
		configScheme = []interface{}{
			map[string]interface{}{
				"fluentbit": []interface{}{
					map[string]interface{}{
						"config": resource.FluentBit.Config,
					},
				},
			},
		}
	default:
		return nil, fmt.Errorf("config scheme is required, log integration config is corrupt: %v", resource)
	}

	return configScheme, nil
}

func (resource *LoggingIntegration) WriteToSchema(d *schema.ResourceData) error {

	if err := d.Set("name", resource.Name); err != nil {
		return fmt.Errorf("error setting 'name': %w", err)
	}
	if err := d.Set("receive_audit_logs", resource.ReceiveAuditLogs); err != nil {
		return fmt.Errorf("error setting 'receive_audit_logs': %w", err)
	}

	configScheme, err := getLoggingConfig(resource)
	if err != nil {
		return err
	}

	if err := d.Set("config", configScheme); err != nil {
		return fmt.Errorf("error setting 'config': %w", err)
	}

	return nil
}

// ReadFromSchema is used to parse the resource schema into a logging integration structure that is expected by the API
func (integrationLogConfig *LoggingIntegration) ReadFromSchema(d *schema.ResourceData) error {
	integrationLogConfig.Id = d.Id() //Get("integration_id").(string)
	integrationLogConfig.Name = d.Get("name").(string)
	integrationLogConfig.ReceiveAuditLogs = d.Get("receive_audit_logs").(bool)

	// Handle Config Scheme (required field).
	configSchemeSet := d.Get("config").([]interface{})
	if len(configSchemeSet) != 1 {
		return fmt.Errorf(
			"exactly one config attribute is required",
		)
	}

	configSchemes := configSchemeSet[0].(map[string]interface{})

	for k, v := range configSchemes {
		configSchemeDetails := v.(*schema.Set).List()
		if len(configSchemeDetails) == 0 {
			continue
		}
		m := configSchemeDetails[0].(map[string]interface{})

		switch k {
		case "cloud_watch":
			integrationLogConfig.CloudWatch = &CloudWatchConfig{
				Region:           m["region"].(string),
				Group:            m["group"].(string),
				Stream:           m["stream"].(string),
				LogRetentionDays: m["log_retention_days"].(int),
			}
		case "datadog":
			integrationLogConfig.Datadog = &DataDogConfig{
				ApiKey: m["api_key"].(string),
			}
		case "elk":
			credentialsSet := m["es_credentials"].(*schema.Set).List()
			credentialScheme := make(map[string]interface{})
			if len(credentialsSet) != 0 {
				credentialScheme = credentialsSet[0].(map[string]interface{})
			}
			integrationLogConfig.Elk = &ElkConfig{
				EsURL:     m["es_url"].(string),
				KibanaURL: m["kibana_url"].(string),
				EsCredentials: EsCredentials{
					Username: credentialScheme["username"].(string),
					Password: credentialScheme["password"].(string),
				},
			}
		case "splunk":
			integrationLogConfig.Splunk = &SplunkConfig{
				Hostname:    m["hostname"].(string),
				HecPort:     m["hec_port"].(string),
				AccessToken: m["access_token"].(string),
				Index:       m["index"].(string),
				UseTLS:      m["use_tls"].(bool),
			}
		case "sumo_logic":
			integrationLogConfig.SumoLogic = &SumoLogicConfig{
				Address: m["address"].(string),
			}
		case "fluentbit":
			integrationLogConfig.FluentBit = &FluentBitConfig{
				Config: m["config"].(string),
			}
		default:
			return fmt.Errorf("unexpected config [%s]", k)
		}
	}
	return nil
}

func CreateLoggingIntegration() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "LoggingIntegrationCreate",
		HttpMethod: http.MethodPost,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/logging", c.ControlPlane)
		},
		NewResourceData: func() ResourceData { return &LoggingIntegration{} },
		NewResponseData: func(_ *schema.ResourceData) ResponseData { return &IDBasedResponse{} },
	}
}

var ReadLoggingIntegration = ResourceOperationConfig{
	Name:       "LoggingIntegrationRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf(loggingApiUrl, c.ControlPlane, d.Id())
	},
	NewResponseData: func(_ *schema.ResourceData) ResponseData { return &LoggingIntegration{} },
}

func UpdateLoggingIntegration() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "LoggingIntegrationUpdate",
		HttpMethod: http.MethodPut,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(loggingApiUrl, c.ControlPlane, d.Id())
		},
		NewResourceData: func() ResourceData { return &LoggingIntegration{} },
	}
}

func DeleteLoggingIntegration() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "LoggingIntegrationDelete",
		HttpMethod: http.MethodDelete,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(loggingApiUrl, c.ControlPlane, d.Id())
		},
	}
}

func resourceIntegrationLogging() *schema.Resource {

	return &schema.Resource{
		Description: "Manages a logging integration that can be used to push logs from Cyral to the corresponding logging system (E.g.: AWS CloudWatch, Splunk, SumoLogic, etc).",
		CreateContext: CreateResource(
			CreateLoggingIntegration(),
			ReadLoggingIntegration,
		),
		ReadContext:   ReadResource(ReadLoggingIntegration),
		UpdateContext: UpdateResource(UpdateLoggingIntegration(), ReadLoggingIntegration),
		DeleteContext: DeleteResource(DeleteLoggingIntegration()),
		Schema:        getIntegrationLogsSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
