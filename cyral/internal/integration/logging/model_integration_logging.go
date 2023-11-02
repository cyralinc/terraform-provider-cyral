package logging

import (
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

const (
	CloudWatchKey = "cloudwatch"
	DatadogKey    = "datadog"
	ElkKey        = "elk"
	SplunkKey     = "splunk"
	SumoLogicKey  = "sumo_logic"
	FluentbitKey  = "fluent_bit"
)

var allLogIntegrationConfigs = []string{
	CloudWatchKey,
	DatadogKey,
	ElkKey,
	SplunkKey,
	SumoLogicKey,
	FluentbitKey,
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
