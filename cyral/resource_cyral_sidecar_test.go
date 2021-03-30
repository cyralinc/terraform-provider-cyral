package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialSidecarConfig SidecarData = SidecarData{
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

var updatedSidecarConfig SidecarData = SidecarData{
	Name: "sidecar-updated-test",
	SidecarProperty: SidecarProperty{
		DeploymentMethod:     "terraform",
		AWSRegion:            "us-west-1",
		KeyName:              "myEC2Key-updated",
		VPC:                  "vpc-123456789",
		Subnets:              "subnet-123456789,subnet-789101112",
		PubliclyAccessible:   "false",
		MetricsIntegrationID: "",
		LogIntegrationID:     "",
	},
}

func TestAccSidecarResource(t *testing.T) {
	testConfig, testFunc := setupSidecarTest(initialSidecarConfig)
	testUpdateConfig, testUpdateFunc := setupSidecarTest(updatedSidecarConfig)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdateConfig,
				Check:  testUpdateFunc,
			},
		},
	})
}

func setupSidecarTest(integrationData SidecarData) (string, resource.TestCheckFunc) {
	configuration := formatSidecarDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_sidecar.test_repo_binding_sidecar", "name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_sidecar.test_repo_binding_sidecar", "deployment_method", integrationData.SidecarProperty.DeploymentMethod),
	)

	return configuration, testFunction
}

func formatSidecarDataIntoConfig(data SidecarData) string {
	return fmt.Sprintf(`
	resource "cyral_sidecar" "test_repo_binding_sidecar" {
		name = "%s"
		deployment_method = "%s"
		aws_configuration {
			publicly_accessible = %s
			aws_region = "%s"
			key_name = "%s"
			vpc = "%s"
			subnets = "%s"
		}
	}`, data.Name, data.SidecarProperty.DeploymentMethod,
		data.SidecarProperty.PubliclyAccessible, data.SidecarProperty.AWSRegion,
		data.SidecarProperty.KeyName, data.SidecarProperty.VPC,
		data.SidecarProperty.Subnets)
}
