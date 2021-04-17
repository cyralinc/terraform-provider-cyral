package cyral

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var cftSidecarConfig SidecarData = SidecarData{
	Name: "sidecar-test",
	SidecarProperty: SidecarProperty{
		DeploymentMethod:     "cloudFormation",
		AWSRegion:            "us-east-1",
		KeyName:              "myEC2Key",
		VPC:                  "vpc-123456",
		Subnets:              "subnet-123456,subnet-789101",
		PubliclyAccessible:   "true",
		MetricsIntegrationID: "default",
		LogIntegrationID:     "default",
	},
}

var dockerSidecarConfig SidecarData = SidecarData{
	Name: "sidecar-test",
	SidecarProperty: SidecarProperty{
		DeploymentMethod:     "docker",
		AWSRegion:            "",
		KeyName:              "",
		VPC:                  "",
		Subnets:              "",
		PubliclyAccessible:   "true",
		MetricsIntegrationID: "default",
		LogIntegrationID:     "default",
	},
}

var helmSidecarConfig SidecarData = SidecarData{
	Name: "sidecar-test",
	SidecarProperty: SidecarProperty{
		DeploymentMethod:     "helm",
		AWSRegion:            "",
		KeyName:              "",
		VPC:                  "",
		Subnets:              "",
		PubliclyAccessible:   "true",
		MetricsIntegrationID: "default",
		LogIntegrationID:     "default",
	},
}

var terraformSidecarConfig SidecarData = SidecarData{
	Name: "sidecar-test",
	SidecarProperty: SidecarProperty{
		DeploymentMethod:     "terraform",
		AWSRegion:            "us-east-1",
		KeyName:              "myEC2Key",
		VPC:                  "vpc-123456",
		Subnets:              "subnet-123456,subnet-789101",
		PubliclyAccessible:   "true",
		MetricsIntegrationID: "default",
		LogIntegrationID:     "default",
	},
}

func TestAccSidecarTemplateDataSource(t *testing.T) {
	cftConfig, cftFunc := setupSidecarTemplateTest(cftSidecarConfig)
	dockerConfig, dockerFunc := setupSidecarTemplateTest(dockerSidecarConfig)
	helmConfig, helmFunc := setupSidecarTemplateTest(helmSidecarConfig)
	tfConfig, tfFunc := setupSidecarTemplateTest(terraformSidecarConfig)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: cftConfig,
				Check:  cftFunc,
			},
			{
				Config: dockerConfig,
				Check:  dockerFunc,
			},
			{
				Config: helmConfig,
				Check:  helmFunc,
			},
			{
				Config: tfConfig,
				Check:  tfFunc,
			},
		},
	})
}

func setupSidecarTemplateTest(integrationData SidecarData) (string, resource.TestCheckFunc) {
	configuration := formatSidecarDataIntoConfig(integrationData) +
		formatSidecarTemplateDataIntoConfig()

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestMatchOutput("ret_template", regexp.MustCompile("\\w+")),
	)

	return configuration, testFunction
}

func formatSidecarTemplateDataIntoConfig() string {
	return `
	data "cyral_sidecar_template" "test_template" {
		sidecar_id = cyral_sidecar.test_repo_binding_sidecar.id
	}
	
	output "ret_template" {
	    value = data.cyral_sidecar_template.test_template.template
	}`
}
