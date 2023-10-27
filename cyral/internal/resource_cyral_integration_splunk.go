package internal

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SplunkIntegration struct {
	Name                     string `json:"name"`
	AccessToken              string `json:"accessToken"`
	Port                     int    `json:"hecPort,string"`
	Host                     string `json:"host"`
	Index                    string `json:"index"`
	UseTLS                   bool   `json:"useTLS"`
	CyralActivityLogsEnabled bool   `json:"cyralActivityLogsEnabled"`
}

func (data SplunkIntegration) WriteToSchema(d *schema.ResourceData) error {
	d.Set("name", data.Name)
	d.Set("access_token", data.AccessToken)
	d.Set("port", data.Port)
	d.Set("host", data.Host)
	d.Set("index", data.Index)
	d.Set("use_tls", data.UseTLS)
	return nil
}

func (data *SplunkIntegration) ReadFromSchema(d *schema.ResourceData) error {
	data.Name = d.Get("name").(string)
	data.AccessToken = d.Get("access_token").(string)
	data.Port = d.Get("port").(int)
	data.Host = d.Get("host").(string)
	data.Index = d.Get("index").(string)
	data.UseTLS = d.Get("use_tls").(bool)
	return nil
}

var ReadSplunkConfig = core.ResourceOperationConfig{
	Name:       "SplunkResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/splunk/%s", c.ControlPlane, d.Id())
	},
	NewResponseData: func(_ *schema.ResourceData) core.ResponseData { return &SplunkIntegration{} },
}

func ResourceIntegrationSplunk() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Use resource `cyral_integration_logging` instead.",
		Description:        "Manages [integration with Splunk](https://cyral.com/docs/integrations/siem/splunk/#procedure).",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				Name:       "SplunkResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/splunk", c.ControlPlane)
				},
				NewResourceData: func() core.ResourceData { return &SplunkIntegration{} },
				NewResponseData: func(_ *schema.ResourceData) core.ResponseData { return &core.IDBasedResponse{} },
			}, ReadSplunkConfig,
		),
		ReadContext: core.ReadResource(ReadSplunkConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				Name:       "SplunkResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/splunk/%s", c.ControlPlane, d.Id())
				},
				NewResourceData: func() core.ResourceData { return &SplunkIntegration{} },
			}, ReadSplunkConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				Name:       "SplunkResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/integrations/splunk/%s", c.ControlPlane, d.Id())
				},
			},
		),

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Integration name that will be used internally in the control plane.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"access_token": {
				Description: "Splunk access token.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"port": {
				Description: "Splunk host port.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"host": {
				Description: "Splunk host.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"index": {
				Description: "Splunk data index name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"use_tls": {
				Description: "Should the communication with Splunk use TLS encryption?",
				Type:        schema.TypeBool,
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
