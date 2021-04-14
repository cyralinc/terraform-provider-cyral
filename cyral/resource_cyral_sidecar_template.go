package cyral

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SidecarTemplateData struct {
	SidecarId string
}

func (data SidecarTemplateData) WriteToSchema(d *schema.ResourceData) {
	data.SidecarId = d.Get("sidecar_id").(string)
	d.SetId(data.SidecarId)
}

func (data *SidecarTemplateData) ReadFromSchema(d *schema.ResourceData) {
	data.SidecarId = d.Id()
}

var getSidecarTemplate = ResourceOperationConfig{
	Name:       "SidecarTemplatesCreate",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		controlPlane := removePortFromURL(c.ControlPlane)
		return fmt.Sprintf("https://%s/deploy/cft/?SidecarId=%s&KeyName=%s&VPC=&SidecarName=%s&ControlPlane=%s&PublicSubnets=&ELKAddress=&publiclyAccessible=%t&logIntegrationType=&logIntegrationValue=&metricsIntegrationType=&metricsIntegrationValue=&",
			controlPlane, d.Get("sidecar_id").(string),
			"ec2_key",
			"name",
			controlPlane,
			true)
	},
	ResourceData: &SidecarTemplateData{},
	ResponseData: &SidecarTemplateResponse{},
}

func resourceDataSidecarTemplates() *schema.ResourceData {
	return &schema.ResourceData{
		CreateContext: getCyralSidecarTemplate(getSidecarTemplate),
		ReadContext: EmptyReadAction(
			ResourceOperationConfig{
				ResourceData: &SidecarTemplateData{},
			}),
		UpdateContext: updateCyralSidecarTemplate(getSidecarTemplate),
		DeleteContext: EmptyDeleteAction(
			ResourceOperationConfig{},
		),
		Schema: map[string]*schema.Schema{
			"sidecar_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func getCyralSidecarTemplate(config ResourceOperationConfig) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		log.Printf("[DEBUG] Init %s", config.Name)
		c := m.(*client.Client)

		config.ResourceData.ReadFromSchema(d)

		url := config.CreateURL(d, c)

		body, err := c.DoRequest(url, config.HttpMethod, config.ResourceData)
		if err != nil {
			return createError("Unable to create integration", fmt.Sprintf("%v", err))
		}

		log.Printf("[INFO]Sidecar Template:\n %v", body)

		config.ResponseData.WriteToSchema(d)

		log.Printf("[DEBUG] End %s", config.Name)

		return diag.Diagnostics{}
	}
}

func updateCyralSidecarTemplate(config ResourceOperationConfig) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		log.Printf("[DEBUG] Init %s", config.Name)
		c := m.(*client.Client)

		config.ResourceData.ReadFromSchema(d)

		url := config.CreateURL(d, c)

		body, err := c.DoRequest(url, config.HttpMethod, config.ResourceData)
		if err != nil {
			return createError("Unable to update integration", fmt.Sprintf("%v", err))
		}

		log.Printf("[INFO] Sidecar Template:\n %v", body)

		config.ResourceData.WriteToSchema(d)

		log.Printf("[DEBUG] End %s", config.Name)

		return diag.Diagnostics{}
	}
}

func EmptyReadAction(config ResourceOperationConfig) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		config.ResourceData.ReadFromSchema(d)
		return diag.Diagnostics{}
	}
}

func EmptyDeleteAction(config ResourceOperationConfig) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		return diag.Diagnostics{}
	}
}

func removePortFromURL(url string) string {
	return strings.Split(url, ":")[0]
}
