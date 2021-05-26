package cyral

import (
	"encoding/json"
	"errors"
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

func dataSourceSidecarTemplate() *schema.Resource {
	return &schema.Resource{
		Description: "Returns the deployment template for a given sidecar",
		Read:        getSidecarTemplate,
		Schema: map[string]*schema.Schema{
			"sidecar_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The sidecar id you want the template for",
			},
			"template": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The output variable that will contain the template",
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

	sidecarId := d.Get("sidecar_id").(string)

	properties, sidecarTypeErr := getSidecarData(c, d)
	if sidecarTypeErr != nil {
		return sidecarTypeErr
	}

	body, err := getTemplateForSidecarProperties(properties, c, d)
	if err != nil {
		return err
	}

	d.SetId(sidecarId)
	d.Set("template", string(body))

	log.Printf("[DEBUG] End Init Get Sidecar Template")

	return nil
}

func removePortFromURL(url string) string {
	return strings.Split(url, ":")[0]
}

func getSidecarData(c *client.Client, d *schema.ResourceData) (*SidecarData, error) {
	url := fmt.Sprintf("https://%s/v1/sidecars/%s", c.ControlPlane, d.Get("sidecar_id").(string))

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	response := SidecarData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func getTemplateForSidecarProperties(data *SidecarData, c *client.Client, d *schema.ResourceData) ([]byte, error) {
	controlPlane := removePortFromURL(c.ControlPlane)
	var url string
	switch data.SidecarProperty.DeploymentMethod {
	case "cloudFormation":
		url = fmt.Sprintf("https://%s/deploy/cft/?SidecarId=%s&KeyName=%s&VPC=&SidecarName=%s&ControlPlane=%s&PublicSubnets=&ELKAddress=&publiclyAccessible=%s&logIntegrationType=&logIntegrationValue=&metricsIntegrationType=&metricsIntegrationValue=&",
			controlPlane,
			d.Get("sidecar_id").(string),
			data.SidecarProperty.KeyName,
			data.Name,
			controlPlane,
			data.SidecarProperty.PubliclyAccessible)
	case "docker":
		url = fmt.Sprintf("https://%s/deploy/docker-compose?SidecarId=%s&SidecarName=%s&logIntegrationType=&logIntegrationValue=&metricsIntegrationType=&metricsIntegrationValue=&SplunkIndex=&SplunkHost=&SplunkPort=&SplunkTLS=&SplunkToken=&",
			controlPlane,
			d.Get("sidecar_id").(string),
			data.Name,
		)
	case "terraform":
		url = fmt.Sprintf("https://%s/deploy/terraform/?SidecarId=%s&AWSRegion=%s&KeyName=%s&VPC=%s&SidecarName=%s&ControlPlane=%s&PublicSubnets[]=%s&publiclyAccessible=%s&logIntegrationType=&logIntegrationValue=&metricsIntegrationType=&metricsIntegrationValue=&",
			controlPlane,
			d.Get("sidecar_id").(string),
			data.SidecarProperty.AWSRegion,
			data.SidecarProperty.KeyName,
			data.SidecarProperty.VPC,
			data.Name,
			controlPlane,
			data.SidecarProperty.Subnets,
			data.SidecarProperty.PubliclyAccessible,
		)
	case "helm":
		url = fmt.Sprintf("https://%s/deploy/helm/values.yaml?sidecarId=%s&logIntegrationType=&logIntegrationValue=&metricsIntegrationType=&metricsIntegrationValue=&SumologicHost=&SumologicUri=&",
			controlPlane,
			d.Get("sidecar_id").(string),
		)
	default:
		return nil, errors.New("invalid deployment method")
	}
	return c.DoRequest(url, http.MethodGet, nil)
}
