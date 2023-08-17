package cyral

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/src/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	sidecarCftTemplateDataSourceName = "data-sidecar-cft-template"
)

func TestAccSidecarCftTemplateDataSource(t *testing.T) {
	cftConfig, cftFunc := setupSidecarCftTemplateTest()

	resource.ParallelTest(
		t, resource.TestCase{
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config: cftConfig,
					Check:  cftFunc,
				},
			},
		},
	)
}

func setupSidecarCftTemplateTest() (string, resource.TestCheckFunc) {
	var configuration string
	configuration += utils.FormatBasicSidecarIntoConfig(
		BasicSidecarResName,
		utils.AccTestName(sidecarCftTemplateDataSourceName, "sidecar"),
		"cloudFormation", "",
	)
	configuration += formatELKIntegrationDataIntoConfig(
		ELKIntegration{
			Name:      utils.AccTestName(sidecarCftTemplateDataSourceName, "elk"),
			KibanaURL: "kibana.local",
			ESURL:     "es.local",
		},
	)
	configuration += formatDatadogIntegrationDataIntoConfig(
		DatadogIntegration{
			Name:   utils.AccTestName(sidecarCftTemplateDataSourceName, "datadog"),
			APIKey: "datadog-api-key",
		},
	)
	configuration += formatSidecarCftTemplateDataIntoConfig(
		utils.BasicSidecarID,
		"cyral_integration_elk.elk_integration.id",
		"cyral_integration_datadog.datadog_integration.id",
		true,
		"ec2-key-name",
	)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestMatchOutput("output_template", regexp.MustCompile(`\w+`)),
	)

	return configuration, testFunction
}

func formatSidecarCftTemplateDataIntoConfig(
	sidecarID, logIntegrationID, metricsIntegrationID string,
	publiclyAccessible bool,
	keyName string,
) string {
	return fmt.Sprintf(
		`
	data "cyral_sidecar_cft_template" "test_template" {
		sidecar_id             = %s
		log_integration_id     = %s
		metrics_integration_id = %s
		aws_configuration {
			publicly_accessible = %t
			key_name            = "%s"
		}
	}
	output "output_template" {
	    value = data.cyral_sidecar_cft_template.test_template.template
	}`, sidecarID, logIntegrationID, metricsIntegrationID, publiclyAccessible,
		keyName,
	)
}
