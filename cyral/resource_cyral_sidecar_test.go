package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialSidecarConfig SidecarData = SidecarData{
	Name: "sidecar-test",
	SidecarProperty: SidecarProperty{
		DeploymentMethod: "cloudFormation",
	},
}

var updatedSidecarConfig SidecarData = SidecarData{
	Name: "sidecar-updated-test",
	SidecarProperty: SidecarProperty{
		DeploymentMethod: "cloudFormation",
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
			publicly_accessible = true
		}
	}`, data.Name, data.SidecarProperty.DeploymentMethod)
}
