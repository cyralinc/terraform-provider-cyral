package deprecated

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/sidecar"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const CloudFormationDeploymentMethod = "cft-ec2"

func DataSourceSidecarCftTemplate() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This data source was deprecated. It will be removed in the next major version of " +
			"the provider.",
		Description: "Retrieves the CloudFormation deployment template for a given sidecar. This data source only " +
			"supports sidecars with `cft-ec2` deployment method. For Terraform template, use our " +
			"`terraform-cyral-sidecar-aws` module.",
		Read: getSidecarCftTemplate,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Same as `sidecar_id`.",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"sidecar_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the sidecar which the template will be generated.",
			},
			"log_integration_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "ID of the log integration that will be used by this template.",
			},
			"metrics_integration_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "ID of the metrics integration that will be used by this template.",
			},
			"aws_configuration": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"publicly_accessible": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Defines a public IP and an internet-facing LB if set to `true`.",
						},
						"key_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "Key-pair name that will be associated to the sidecar EC2 instances.",
						},
					},
				},
				Description: "AWS parameters for `cft-ec2` deployment method.",
			},
			"template": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Output variable with the template.",
			},
		},
	}
}

func getSidecarCftTemplate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Init Get Sidecar CFT Template")
	c := m.(*client.Client)

	sidecarId := d.Get("sidecar_id").(string)

	sidecarData, sidecarTypeErr := getSidecarData(c, d)
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

	body, err := getTemplateForSidecarProperties(sidecarData, logging, metrics, c, d)
	if err != nil {
		return err
	}

	d.SetId(sidecarId)
	d.Set("template", string(body))

	log.Printf("[DEBUG] End Get Sidecar CFT Template")

	return nil
}

func removePortFromURL(url string) string {
	return strings.Split(url, ":")[0]
}

func getSidecarData(c *client.Client, d *schema.ResourceData) (sidecar.SidecarData, error) {
	url := fmt.Sprintf("https://%s/v1/sidecars/%s", c.ControlPlane, d.Get("sidecar_id").(string))

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return sidecar.SidecarData{}, err
	}

	response := sidecar.SidecarData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return sidecar.SidecarData{}, err
	}

	return response, nil
}

func getLogIntegrations(c *client.Client, d *schema.ResourceData) ([]IntegrationsData, error) {
	url := fmt.Sprintf("https://%s/integrations/logging/", removePortFromURL(c.ControlPlane))

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	response := []IntegrationsData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func getMetricsIntegrations(c *client.Client, d *schema.ResourceData) ([]IntegrationsData, error) {
	url := fmt.Sprintf("https://%s/integrations/metrics", removePortFromURL(c.ControlPlane))

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	response := []IntegrationsData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func filterIntegrationData(integrations []IntegrationsData, id string) *IntegrationsData {
	for _, it := range integrations {
		if it.Id == id {
			return &it
		}
	}
	return NewDefaultIntegrationsData()
}

func getTemplateForSidecarProperties(
	sidecarData sidecar.SidecarData,
	logging []IntegrationsData,
	metrics []IntegrationsData,
	c *client.Client,
	d *schema.ResourceData,
) ([]byte, error) {
	controlPlane := removePortFromURL(c.ControlPlane)

	logIntegrationID := d.Get("log_integration_id").(string)
	logIntegration := filterIntegrationData(logging, logIntegrationID)

	metricsIntegrationID := d.Get("metrics_integration_id").(string)
	metricIntegration := filterIntegrationData(metrics, metricsIntegrationID)

	var keyName string
	var publiclyAccessible string

	awsConfig := d.Get("aws_configuration").(*schema.Set).List()
	for _, config := range awsConfig {
		config := config.(map[string]interface{})

		if v, ok := config["key_name"].(string); ok {
			keyName = v
		}
		if v, ok := config["publicly_accessible"].(bool); ok {
			publiclyAccessible = strconv.FormatBool(v)
		}
	}

	logIntegrationValue, err := logIntegration.GetValue()
	if err != nil {
		return nil, err
	}
	metricIntegrationValue, err := metricIntegration.GetValue()
	if err != nil {
		return nil, err
	}

	sidecarTemplatePropertiesKV := map[string]string{
		"SidecarId":               d.Get("sidecar_id").(string),
		"KeyName":                 keyName,
		"SidecarName":             sidecarData.Name,
		"ControlPlane":            controlPlane,
		"clientId":                "",
		"clientSecret":            "",
		"VPC":                     "",
		"PublicSubnets":           "",
		"ELKAddress":              "",
		"publiclyAccessible":      publiclyAccessible,
		"logIntegrationType":      logIntegration.Type,
		"logIntegrationValue":     logIntegrationValue,
		"metricsIntegrationType":  metricIntegration.Type,
		"metricsIntegrationValue": metricIntegrationValue,
	}

	var url string
	properties := sidecarData.SidecarProperties
	if properties != nil && properties.DeploymentMethod == CloudFormationDeploymentMethod {
		url = fmt.Sprintf("https://%s/deploy/cft/", controlPlane)
		url += utils.UrlQuery(sidecarTemplatePropertiesKV)
	} else {
		return nil, fmt.Errorf("invalid deployment method, only '%s' is supported",
			CloudFormationDeploymentMethod)
	}

	return c.DoRequest(url, http.MethodGet, nil)
}
