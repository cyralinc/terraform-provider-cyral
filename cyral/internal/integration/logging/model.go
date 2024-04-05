package logging

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CloudWatchConfig struct {
	Region string `json:"region"`
	Group  string `json:"group"`
	Stream string `json:"stream"`
}

type DataDogConfig struct {
	ApiKey string `json:"apiKey"`
}

type EsCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ElkConfig struct {
	EsURL         string         `json:"esUrl"`
	KibanaURL     string         `json:"kibanaUrl"`
	EsCredentials *EsCredentials `json:"esCredentials"`
}

type SplunkConfig struct {
	Hostname    string `json:"hostname"`
	HecPort     string `json:"hecPort"`
	AccessToken string `json:"accessToken"`
	Index       string `json:"index"`
	UseTLS      bool   `json:"useTLS"`
}

type SumoLogicConfig struct {
	Address string `json:"address"`
}

type FluentBitConfig struct {
	Config       string `json:"config"`
	SkipValidate bool   `json:"skipValidate"`
}

type LoggingIntegration struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	ReceiveAuditLogs bool   `json:"receiveAuditLogs"`
	LoggingIntegrationConfig
}

type LoggingIntegrationConfig struct {
	CloudWatch *CloudWatchConfig `json:"cloudWatch"`
	Datadog    *DataDogConfig    `json:"datadog"`
	Elk        *ElkConfig        `json:"elk"`
	Splunk     *SplunkConfig     `json:"splunk"`
	SumoLogic  *SumoLogicConfig  `json:"sumoLogic"`
	FluentBit  *FluentBitConfig  `json:"fluentBit"`
}

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

type ListIntegrationLogsResponse struct {
	Integrations []LoggingIntegration `json:"integrations"`
}

func (resp *ListIntegrationLogsResponse) WriteToSchema(d *schema.ResourceData) error {
	integrations := make([]interface{}, len(resp.Integrations))
	for i, integration := range resp.Integrations {
		// write in config scheme
		configType, config, err := getLoggingConfig(&integration)
		if err != nil {
			return err
		}
		integrations[i] = map[string]interface{}{
			"id":                 integration.Id,
			"name":               integration.Name,
			"receive_audit_logs": integration.ReceiveAuditLogs,
			configType:           config,
		}
	}
	if err := d.Set("integrations", integrations); err != nil {
		return err
	}

	d.SetId(uuid.New().String())

	return nil
}

func getIntegrationLogsSchema() map[string]*schema.Schema {
	schema := map[string]*schema.Schema{
		"id": {
			Description: "Unique identifier of the logging integration.",
			Computed:    true,
			Type:        schema.TypeString,
		},
		"name": {
			Description: "Name of the logging integration config.",
			Required:    true,
			Type:        schema.TypeString,
		},
		"receive_audit_logs": {
			Description:   "Whether or not Cyral audit logs should be forwarded to this logging integration. Declaration not supported in conjunction with `fluent_bit` block.",
			Optional:      true,
			Type:          schema.TypeBool,
			Default:       false,
			ConflictsWith: []string{FluentbitKey},
		},
		CloudWatchKey: {
			Description:   "Represents the configuration data required for the `AWS` CloudWatch log management system.",
			Type:          schema.TypeSet,
			Optional:      true,
			ConflictsWith: []string{DatadogKey, ElkKey, SplunkKey, SumoLogicKey, FluentbitKey},
			MaxItems:      1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"region": {
						Description: "AWS region.",
						Required:    true,
						Type:        schema.TypeString,
					},
					"group": {
						Description: "CloudWatch log group.",
						Required:    true,
						Type:        schema.TypeString,
					},
					"stream": {
						Description: "CloudWatch log stream. Defaults to `cyral-sidecar` if not set.",
						Optional:    true,
						Type:        schema.TypeString,
					},
				},
			},
		},

		DatadogKey: {
			Description:   "Represents the configuration data required for the Datadog's log management system.",
			Optional:      true,
			Type:          schema.TypeSet,
			ConflictsWith: []string{CloudWatchKey, ElkKey, SplunkKey, SumoLogicKey, FluentbitKey},
			MaxItems:      1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"api_key": {
						Description: "DataDog API key.",
						Required:    true,
						Sensitive:   true,
						Type:        schema.TypeString,
					},
				},
			},
		},

		ElkKey: {
			Description:   "Represents the configuration data required for the ELK stack log management system.",
			Optional:      true,
			Type:          schema.TypeSet,
			ConflictsWith: []string{CloudWatchKey, DatadogKey, SplunkKey, SumoLogicKey, FluentbitKey},
			MaxItems:      1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"es_url": {
						Description: "Elasticsearch URL.",
						Required:    true,
						Type:        schema.TypeString,
					},
					"kibana_url": {
						Description: "Kibana URL.",
						Optional:    true,
						Type:        schema.TypeString,
					},
					"es_credentials": {
						Description: "Credentials used to authenticate to Elastic Search." +
							"Can be omitted for unprotected instances.",
						Optional: true,
						Type:     schema.TypeSet,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"username": {
									Description: "Elasticsearch username.",
									Required:    true,
									Type:        schema.TypeString,
								},
								"password": {
									Description: "Elasticsearch password.",
									Required:    true,
									Sensitive:   true,
									Type:        schema.TypeString,
								},
							},
						},
					},
				},
			},
		},

		SplunkKey: {
			Description:   "Represents the configuration data required for the Splunk log management system.",
			Optional:      true,
			Type:          schema.TypeSet,
			ConflictsWith: []string{CloudWatchKey, DatadogKey, ElkKey, SumoLogicKey, FluentbitKey},
			MaxItems:      1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"hostname": {
						Description: "Splunk hostname.",
						Required:    true,
						Type:        schema.TypeString,
					},
					"hec_port": {
						Description: "Splunk HTTP Event Collector (HEC) port.",
						Required:    true,
						Type:        schema.TypeString,
					},
					"access_token": {
						Description: "Splunk access token.",
						Required:    true,
						Sensitive:   true,
						Type:        schema.TypeString,
					},
					"index": {
						Description: "Splunk index which logs should be indexed to.",
						Optional:    true,
						Type:        schema.TypeString,
					},
					"use_tls": {
						Description: "Whether or not to use TLS.",
						Optional:    true,
						Type:        schema.TypeBool,
					},
				},
			},
		},

		SumoLogicKey: {
			Description:   "Represents the configuration data required for the Sumo Logic log management system.",
			Optional:      true,
			Type:          schema.TypeSet,
			ConflictsWith: []string{CloudWatchKey, DatadogKey, ElkKey, SplunkKey, FluentbitKey},
			MaxItems:      1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"address": {
						Description: "Sumo Logic HTTP collector address. A full URL is expected",
						Required:    true,
						Type:        schema.TypeString,
					},
				},
			},
		},

		FluentbitKey: {
			Description:   "Represents a custom Fluent Bit configuration which will be utilized by the sidecar's log shipper.",
			Optional:      true,
			Type:          schema.TypeSet,
			ConflictsWith: []string{CloudWatchKey, DatadogKey, ElkKey, SplunkKey, SumoLogicKey},
			MaxItems:      1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"config": {
						Description: "Fluent Bit configuration, in 'classic mode' INI format. For more details, see: https://docs.fluentbit.io/manual/administration/configuring-fluent-bit/classic-mode/configuration-file",
						Required:    true,
						Type:        schema.TypeString,
					},
					"skip_validate": {
						Optional:    true,
						Type:        schema.TypeBool,
						Description: "Whether to validate the Fluent Bit config.",
					},
				},
			},
		},
	}

	return schema
}
