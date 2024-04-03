package logging

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
)

const loggingApiUrl = "https://%s/v1/integrations/logging/%s"

func getLoggingConfig(resource *LoggingIntegration) (string, []interface{}, error) {
	var configType string
	var configScheme []interface{}
	switch {
	case resource.CloudWatch != nil:
		configType = CloudWatchKey
		configScheme = []interface{}{
			map[string]interface{}{
				"region": resource.CloudWatch.Region,
				"group":  resource.CloudWatch.Group,
				"stream": resource.CloudWatch.Stream,
			},
		}
	case resource.Datadog != nil:
		configType = DatadogKey
		configScheme = []interface{}{
			map[string]interface{}{
				"api_key": resource.Datadog.ApiKey,
			},
		}
	case resource.Elk != nil:
		configType = ElkKey
		elkConfig := map[string]interface{}{
			"es_url":     resource.Elk.EsURL,
			"kibana_url": resource.Elk.KibanaURL,
		}
		// Optional, so we need to verify separately
		if resource.Elk.EsCredentials != nil {
			elkConfig["es_credentials"] = []interface{}{
				map[string]interface{}{
					"username": resource.Elk.EsCredentials.Username,
					"password": resource.Elk.EsCredentials.Password,
				},
			}
		}
		configScheme = []interface{}{elkConfig}
	case resource.Splunk != nil:
		configType = SplunkKey
		configScheme = []interface{}{
			map[string]interface{}{
				"hostname":     resource.Splunk.Hostname,
				"hec_port":     resource.Splunk.HecPort,
				"access_token": resource.Splunk.AccessToken,
				"index":        resource.Splunk.Index,
				"use_tls":      resource.Splunk.UseTLS,
			},
		}
	case resource.SumoLogic != nil:
		configType = SumoLogicKey
		configScheme = []interface{}{
			map[string]interface{}{
				"address": resource.SumoLogic.Address,
			},
		}
	case resource.FluentBit != nil:
		configType = FluentbitKey
		configScheme = []interface{}{
			map[string]interface{}{
				"config":        resource.FluentBit.Config,
				"skip_validate": resource.FluentBit.SkipValidate,
			},
		}
	default:
		return configType, nil, fmt.Errorf("config scheme is required, log integration config is corrupt: %v", resource)
	}

	return configType, configScheme, nil
}

func (resource *LoggingIntegration) WriteToSchema(d *schema.ResourceData) error {
	if err := d.Set("name", resource.Name); err != nil {
		return fmt.Errorf("error setting 'name': %w", err)
	}
	if err := d.Set("receive_audit_logs", resource.ReceiveAuditLogs); err != nil {
		return fmt.Errorf("error setting 'receive_audit_logs': %w", err)
	}

	configType, configScheme, err := getLoggingConfig(resource)
	if err != nil {
		return err
	}

	if err := d.Set(configType, configScheme); err != nil {
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
	var configType string
	var config interface{}
	for _, integrationType := range allLogIntegrationConfigs {
		configAux, isConfigOk := d.GetOk(integrationType)
		if isConfigOk {
			config = configAux.(interface{})
			configType = integrationType
			break
		}
	}

	configDetails := config.(*schema.Set).List()

	m := configDetails[0].(map[string]interface{})

	switch configType {
	case CloudWatchKey:
		integrationLogConfig.CloudWatch = &CloudWatchConfig{
			Region: m["region"].(string),
			Group:  m["group"].(string),
			Stream: m["stream"].(string),
		}
	case DatadogKey:
		integrationLogConfig.Datadog = &DataDogConfig{
			ApiKey: m["api_key"].(string),
		}
	case ElkKey:
		integrationLogConfig.Elk = &ElkConfig{
			EsURL:     m["es_url"].(string),
			KibanaURL: m["kibana_url"].(string),
		}
		credentialsSet := m["es_credentials"].(*schema.Set).List()
		if len(credentialsSet) != 0 {
			credentialScheme := make(map[string]interface{})
			credentialScheme = credentialsSet[0].(map[string]interface{})
			integrationLogConfig.Elk.EsCredentials = &EsCredentials{
				Username: credentialScheme["username"].(string),
				Password: credentialScheme["password"].(string),
			}
		}
	case SplunkKey:
		integrationLogConfig.Splunk = &SplunkConfig{
			Hostname:    m["hostname"].(string),
			HecPort:     m["hec_port"].(string),
			AccessToken: m["access_token"].(string),
			Index:       m["index"].(string),
			UseTLS:      m["use_tls"].(bool),
		}
	case SumoLogicKey:
		integrationLogConfig.SumoLogic = &SumoLogicConfig{
			Address: m["address"].(string),
		}
	case FluentbitKey:
		integrationLogConfig.FluentBit = &FluentBitConfig{
			Config:       m["config"].(string),
			SkipValidate: m["skip_validate"].(bool),
		}
	default:
		return fmt.Errorf("unexpected config type [%s]", configType)
	}
	return nil
}

func CreateLoggingIntegration() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: "LoggingIntegrationCreate",
		Type:         operationtype.Create,
		HttpMethod:   http.MethodPost,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/logging", c.ControlPlane)
		},
		SchemaReaderFactory: func() core.SchemaReader { return &LoggingIntegration{} },
		SchemaWriterFactory: core.DefaultSchemaWriterFactory,
	}
}

var ReadLoggingIntegration = core.ResourceOperationConfig{
	ResourceName: "LoggingIntegrationRead",
	Type:         operationtype.Read,
	HttpMethod:   http.MethodGet,
	URLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf(loggingApiUrl, c.ControlPlane, d.Id())
	},
	SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &LoggingIntegration{} },
	RequestErrorHandler: &core.IgnoreHttpNotFound{ResName: "Integration logging"},
}

func UpdateLoggingIntegration() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: "LoggingIntegrationUpdate",
		Type:         operationtype.Update,
		HttpMethod:   http.MethodPut,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(loggingApiUrl, c.ControlPlane, d.Id())
		},
		SchemaReaderFactory: func() core.SchemaReader { return &LoggingIntegration{} },
	}
}

func DeleteLoggingIntegration() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: "LoggingIntegrationDelete",
		Type:         operationtype.Delete,
		HttpMethod:   http.MethodDelete,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(loggingApiUrl, c.ControlPlane, d.Id())
		},
	}
}

func ResourceIntegrationLogging() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a logging integration that can be used to push logs from Cyral to the corresponding logging system (E.g.: AWS CloudWatch, Splunk, SumoLogic, etc).",
		CreateContext: core.CreateResource(
			CreateLoggingIntegration(),
			ReadLoggingIntegration,
		),
		ReadContext:   core.ReadResource(ReadLoggingIntegration),
		UpdateContext: core.UpdateResource(UpdateLoggingIntegration(), ReadLoggingIntegration),
		DeleteContext: core.DeleteResource(DeleteLoggingIntegration()),
		Schema:        getIntegrationLogsSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
