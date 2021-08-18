package cyral

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var cftSidecarConfig SidecarData = SidecarData{
	Name: "sidecar-test",
	Tags: []string{"tag1", "tag2"},
	SidecarProperty: SidecarProperty{
		DeploymentMethod: "cloudFormation",
	},
}

func TestAccSidecarCftTemplateDataSource(t *testing.T) {
	cftConfig, cftFunc := setupSidecarCftTemplateTest(cftSidecarConfig)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: cftConfig,
				Check:  cftFunc,
			},
		},
	})
}

func setupSidecarCftTemplateTest(sidecarData SidecarData) (string, resource.TestCheckFunc) {
	deploymentTag := fmt.Sprintf("%s%s", DeploymentPrefix, sidecarData.SidecarProperty.DeploymentMethod)
	tags := append([]string{deploymentTag}, sidecarData.Tags...)

	configuration := formatSidecarDataIntoConfig(sidecarData.Name, tags) +
		formatSidecarCftTemplateDataIntoConfig()

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestMatchOutput("output_template", regexp.MustCompile(`\w+`)),
	)

	return configuration, testFunction
}

func formatSidecarCftTemplateDataIntoConfig() string {
	return `
	resource "cyral_integration_elk" "elk" {
		name = "my-elk-integration"
		kibana_url = "kibana.local"
		es_url = "es.local"
	}

	resource "cyral_integration_datadog" "datadog" {
		name = "my-datadog-integration"
		api_key = "datadog-api-key"
	}

	data "cyral_sidecar_cft_template" "test_template" {
		sidecar_id = cyral_sidecar.test_sidecar.id
		log_integration_id = cyral_integration_elk.elk.id
  	metrics_integration_id = cyral_integration_datadog.datadog.id
		aws_configuration {
			publicly_accessible = true
			key_name = "ec2-key-name"
		}
	}
	output "output_template" {
	    value = data.cyral_sidecar_cft_template.test_template.template
	}`
}
