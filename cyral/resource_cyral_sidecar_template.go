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
	SidecarId          string
	Name               string
	Ec2Key             string
	PubliclyAccessible bool
}

func (data SidecarTemplateData) WriteToSchema(d *schema.ResourceData) {
	d.Set("name", data.Name)
	d.Set("sidecar_id", data.SidecarId)
	d.Set("ec2_key", data.Ec2Key)
	d.Set("publicly_accessible", data.PubliclyAccessible)
}

func (data *SidecarTemplateData) ReadFromSchema(d *schema.ResourceData) {
	data.Name = d.Get("name").(string)
	data.SidecarId = d.Get("sidecar_id").(string)
	data.Ec2Key = d.Get("ec2_key").(string)
	data.PubliclyAccessible = d.Get("publicly_accessible").(bool)
}

func resourceSidecarTemplates() *schema.Resource {
	return &schema.Resource{
		CreateContext: getCyralSidecarTemplate(
			ResourceOperationConfig{
				Name:       "SidecarTemplatesCreate",
				HttpMethod: http.MethodGet,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					controlPlane := removePortFromURL(c.ControlPlane)
					return fmt.Sprintf("https://%s/deploy/cft/?SidecarId=%s&KeyName=%s&VPC=&SidecarName=%s&ControlPlane=%s&PublicSubnets=&ELKAddress=&publiclyAccessible=%t&logIntegrationType=&logIntegrationValue=&metricsIntegrationType=&metricsIntegrationValue=&",
						controlPlane, d.Get("sidecar_id").(string), d.Get("ec2_key").(string), d.Get("name").(string), controlPlane, d.Get("publicly_accessible").(bool))
				},
				ResourceData: &SidecarTemplateData{},
				ResponseData: &IDBasedResponse{},
			}),
		ReadContext:   EmptyAction,
		UpdateContext: EmptyAction,
		DeleteContext: EmptyAction,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sidecar_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ec2_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"publicly_accessible": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func EmptyAction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

func removePortFromURL(url string) string {
	return strings.Split(url, ":")[0]
}

func getCyralSidecarTemplate(config ResourceOperationConfig) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		log.Printf("[DEBUG] Init %s", config.Name)
		c := m.(*client.Client)

		config.ResourceData.ReadFromSchema(d)

		url := config.CreateURL(d, c)

		_, err := c.DoRequest(url, config.HttpMethod, config.ResourceData)
		if err != nil {
			return createError("Unable to create integration", fmt.Sprintf("%v", err))
		}

		// log.Printf("Sidecar Template: %v", body)

		// config.ResourceData.WriteToSchema(d)

		log.Printf("[DEBUG] End %s", config.Name)

		return diag.Diagnostics{}
	}
}
