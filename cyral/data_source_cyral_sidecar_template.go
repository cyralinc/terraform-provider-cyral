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

type integrationsData struct {
	Id    string `json:"id"`
	Label string `json:"label"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

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
		Read: getSidecarTemplate,
		Schema: map[string]*schema.Schema{
			"sidecar_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template": {
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

	sidecarId := d.Get("sidecar_id").(string)

	properties, sidecarTypeErr := getSidecarData(c, d)
	if sidecarTypeErr != nil {
		return sidecarTypeErr
	}

	logging, err := getLogIntegrations(c, d)
	if err != nil {
		return err
	}

	metrics, err := getMetricsIntegrations(c, d)
	if err != nil {
		return err
	}

	body, err := getTemplateForSidecarProperties(properties, logging, metrics, c, d)
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

func getLogIntegrations(c *client.Client, d *schema.ResourceData) (*[]integrationsData, error) {
	url := fmt.Sprintf("https://%s/integrations/logging/", removePortFromURL(c.ControlPlane))

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	response := []integrationsData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func getMetricsIntegrations(c *client.Client, d *schema.ResourceData) (*[]integrationsData, error) {
	url := fmt.Sprintf("https://%s/integrations/metrics", removePortFromURL(c.ControlPlane))

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	response := []integrationsData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func filterIntegrationData(integrations *[]integrationsData, id string) *integrationsData {
	for _, it := range *integrations {
		if it.Id == id {
			return &it
		}
	}
	return &integrationsData{
		Id:    "id",
		Type:  "default",
		Value: "default",
		Name:  "default",
		Label: "default",
	}
}

func getTemplateForSidecarProperties(data *SidecarData, logging *[]integrationsData, metrics *[]integrationsData, c *client.Client, d *schema.ResourceData) ([]byte, error) {
	controlPlane := removePortFromURL(c.ControlPlane)

	metric := filterIntegrationData(metrics, data.SidecarProperty.MetricsIntegrationID)

	log := filterIntegrationData(logging, data.SidecarProperty.LogIntegrationID)

	var url string
	switch data.SidecarProperty.DeploymentMethod {
	case "cloudFormation":
		url = fmt.Sprintf("https://%s/deploy/cft/?SidecarId=%s&KeyName=%s&VPC=&SidecarName=%s&ControlPlane=%s&PublicSubnets=&ELKAddress=&publiclyAccessible=%s&logIntegrationType=%s&logIntegrationValue=%s&metricsIntegrationType=%s&metricsIntegrationValue=%s&",
			controlPlane,
			d.Get("sidecar_id").(string),
			data.SidecarProperty.KeyName,
			data.Name,
			controlPlane,
			data.SidecarProperty.PubliclyAccessible,
			log.Type,
			log.Value,
			metric.Type,
			metric.Value,
		)
	case "docker":
		url = fmt.Sprintf("https://%s/deploy/docker-compose?SidecarId=%s&SidecarName=%s&logIntegrationType=%s&logIntegrationValue=%s&metricsIntegrationType=%s&metricsIntegrationValue=%s&SplunkIndex=&SplunkHost=&SplunkPort=&SplunkTLS=&SplunkToken=&",
			controlPlane,
			d.Get("sidecar_id").(string),
			data.Name,
			log.Type,
			log.Value,
			metric.Type,
			metric.Value,
		)
	case "terraform":
		url = fmt.Sprintf("https://%s/deploy/terraform/?SidecarId=%s&AWSRegion=%s&KeyName=%s&VPC=%s&SidecarName=%s&ControlPlane=%s&PublicSubnets[]=%s&publiclyAccessible=%s&logIntegrationType=%s&logIntegrationValue=%s&metricsIntegrationType=%s&metricsIntegrationValue=%s&",
			controlPlane,
			d.Get("sidecar_id").(string),
			data.SidecarProperty.AWSRegion,
			data.SidecarProperty.KeyName,
			data.SidecarProperty.VPC,
			data.Name,
			controlPlane,
			data.SidecarProperty.Subnets,
			data.SidecarProperty.PubliclyAccessible,
			log.Type,
			log.Value,
			metric.Type,
			metric.Value,
		)
	case "helm", "helm3":
		url = fmt.Sprintf("https://%s/deploy/helm/values.yaml?sidecarId=%s&logIntegrationType=%s&logIntegrationValue=%s&metricsIntegrationType=%s&metricsIntegrationValue=%s&SumologicHost=&SumologicUri=&",
			controlPlane,
			d.Get("sidecar_id").(string),
			log.Type,
			log.Value,
			metric.Type,
			metric.Value,
		)
	default:
		return nil, errors.New("invalid deployment method")
	}
	return c.DoRequest(url, http.MethodGet, nil)
}
