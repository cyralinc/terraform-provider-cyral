package cyral

import (
	"fmt"

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

type ElkConfig struct {
	EsURL     string `json:"esUrl"`
	KibanaURL string `json:"kibanaUrl"`
	EsUid     string `json:"esUid"`
	EsPwd     string `json:"esPwd"`
}

type SplunkConfig struct {
	Host        string `json:"host"`
	HecPort     string `json:"hecPort"`
	AccessToken string `json:"accessToken"`
	Index       string `json:"index"`
	UseTLS      bool   `json:"useTLS"`
}

type SumoLogicConfig struct {
	Address string `json:"address"`
}

type IntegrationLogConfig struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	EnableAuditLogs bool   `json:"enableAuditLogs"`
	IntegrationConfigScheme
}

type IntegrationConfigScheme struct {
	CloudWatch *CloudWatchConfig `json:"cloudWatch"`
	Datadog    *DataDogConfig    `json:"datadog"`
	Elk        *ElkConfig        `json:"elk"`
	Splunk     *SplunkConfig     `json:"splunk"`
	SumoLogic  *SumoLogicConfig  `json:"sumoLogic"`
}

var allLogIntegrationConfigs = []string{
	"cloud_watch",
	"datadog",
	"elk",
	"splunk",
	"sumo_logic",
}

func getIntegrationLogsSchema() map[string]*schema.Schema {
	configSchemeTypes := make([]string, 0, len(allLogIntegrationConfigs))
	for _, config := range allLogIntegrationConfigs {
		configSchemeTypes = append(configSchemeTypes,
			fmt.Sprintf("config_scheme.0.%s", config))
	}
	schema := map[string]*schema.Schema{
		"id": {
			Description: "The unique identifier of the logging integration.",
			Computed:    true,
			Type:        schema.TypeString,
		},
		"name": {
			Description: "The name of the logging integration config.",
			Required:    true,
			Type:        schema.TypeString,
			ForceNew:    true,
		},
		"enable_audit_logs": {
			Description: "Whether or not Cyral audit logs should be forwarded to this logging integration.",
			Optional:    true,
			Type:        schema.TypeBool,
			ForceNew:    true,
			Default:     true,
		},
		"config_scheme": {
			Description: "Config option for your integration. List of supported types: " +
				supportedTypesMarkdown(allLogIntegrationConfigs),
			Required: true,
			Type:     schema.TypeList,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cloud_watch": {
						Description:  "Represents the configuration data required for the `AWS` CloudWatch log management system.",
						Type:         schema.TypeSet,
						Optional:     true,
						ExactlyOneOf: configSchemeTypes,
						MaxItems:     1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"region": {
									Description: "The AWS region.",
									Required:    true,
									Type:        schema.TypeString,
									ForceNew:    true,
								},
								"group": {
									Description: "The CloudWatch log group.",
									Required:    true,
									Type:        schema.TypeString,
									ForceNew:    true,
								},
								"stream": {
									Description: "The CloudWatch log stream.",
									Required:    true,
									Type:        schema.TypeString,
									ForceNew:    true,
								},
							},
						},
					},

					"datadog": {
						Description:  "Represents the configuration data required for the Datadog's log management system.",
						Optional:     true,
						Type:         schema.TypeSet,
						ExactlyOneOf: configSchemeTypes,
						MaxItems:     1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"api_key": {
									Description: "The DataDog API key.",
									Required:    true,
									Type:        schema.TypeString,
									ForceNew:    true,
								},
							},
						},
					},

					"elk": {
						Description:  "Represents the configuration data required for the ELK stack log management system.",
						Optional:     true,
						Type:         schema.TypeSet,
						ExactlyOneOf: configSchemeTypes,
						MaxItems:     1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"es_url": {
									Description: "The Elasticsearch URL.",
									Required:    true,
									Type:        schema.TypeString,
									ForceNew:    true,
								},
								"kibana_url": {
									Description: "The Kibana URL.",
									Optional:    true,
									Type:        schema.TypeString,
									ForceNew:    true,
								},
								"es_uid": {
									Description: "The Elasticsearch username.",
									Optional:    true,
									Type:        schema.TypeString,
									ForceNew:    true,
								},
								"es_pwd": {
									Description: "The Elasticsearch password.",
									Optional:    true,
									Type:        schema.TypeString,
									ForceNew:    true,
								},
							},
						},
					},

					"splunk": {
						Description:  "Represents the configuration data required for the Splunk log management system.",
						Optional:     true,
						Type:         schema.TypeSet,
						ExactlyOneOf: configSchemeTypes,
						MaxItems:     1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"host": {
									Description: "The Splunk hostname.",
									Required:    true,
									Type:        schema.TypeString,
									ForceNew:    true,
								},
								"hec_port": {
									Description: "The Splunk HTTP Event Collector (HEC) port.",
									Required:    true,
									Type:        schema.TypeString,
									ForceNew:    true,
								},
								"access_token": {
									Description: "The Splunk access token.",
									Required:    true,
									Type:        schema.TypeString,
									ForceNew:    true,
								},
								"index": {
									Description: "The Slunk index.",
									Optional:    true,
									Type:        schema.TypeString,
									ForceNew:    true,
								},
								"use_tls": {
									Description: "Whether or not to use TLS.",
									Required:    true,
									Type:        schema.TypeBool,
									ForceNew:    true,
								},
							},
						},
					},

					"sumo_logic": {
						Description:  "Represents the configuration data required for the Sumo Logic log management system.",
						Optional:     true,
						Type:         schema.TypeSet,
						ExactlyOneOf: configSchemeTypes,
						MaxItems:     1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"address": {
									Description: "The Sumo Logic HTTP collector address. A full URL is expected",
									Required:    true,
									Type:        schema.TypeString,
									ForceNew:    true,
								},
							},
						},
					},
				},
			},
		},
	}

	return schema
}
