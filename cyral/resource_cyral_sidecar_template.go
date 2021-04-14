package cyral

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cyralinc/terraform-provider-cyral/client"
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

func resourceDataSidecarTemplates() *schema.Resource {
	return &schema.Resource{
		Read: getSidecarTemplate,
		Schema: map[string]*schema.Schema{
			"sidecar_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sidecar_template": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func getSidecarTemplate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Init Get Sidecar Template")
	c := m.(*client.Client)

	url := formatUrl(d, c)

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return err
	}

	d.SetId(d.Get("sidecar_id").(string))
	d.Set("sidecar_template", string(body))

	log.Printf("[DEBUG] End Init Get Sidecar Template")

	return nil
}

func formatUrl(d *schema.ResourceData, c *client.Client) string {
	controlPlane := removePortFromURL(c.ControlPlane)
	return fmt.Sprintf("https://%s/deploy/cft/?SidecarId=%s&KeyName=%s&VPC=&SidecarName=%s&ControlPlane=%s&PublicSubnets=&ELKAddress=&publiclyAccessible=%t&logIntegrationType=&logIntegrationValue=&metricsIntegrationType=&metricsIntegrationValue=&",
		controlPlane, d.Get("sidecar_id").(string),
		"ec2_key",
		"name",
		controlPlane,
		true)
}

func removePortFromURL(url string) string {
	return strings.Split(url, ":")[0]
}
