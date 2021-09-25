package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialSidecarConfig SidecarData = SidecarData{
	Name:   "sidecar-test",
	Labels: []string{"test1"},
	SidecarProperty: SidecarProperty{
		DeploymentMethod: "cloudFormation",
	},
}

var updatedSidecarConfigDocker SidecarData = SidecarData{
	Name:   "sidecar-updated-test",
	Labels: []string{"test2"},
	SidecarProperty: SidecarProperty{
		DeploymentMethod: "docker",
	},
}

var updatedSidecarConfigHelm SidecarData = SidecarData{
	Name:   "sidecar-updated-test",
	Labels: []string{"test3"},
	SidecarProperty: SidecarProperty{
		DeploymentMethod: "helm",
	},
}

var updatedSidecarConfigTF SidecarData = SidecarData{
	Name:   "sidecar-updated-test",
	Labels: []string{"test4"},
	SidecarProperty: SidecarProperty{
		DeploymentMethod: "terraform",
	},
}

func TestAccSidecarResource(t *testing.T) {
	testConfig, testFunc := setupSidecarTest(initialSidecarConfig)
	testUpdateConfigDocker, testUpdateFuncDocker := setupSidecarTest(updatedSidecarConfigDocker)
	testUpdateConfigHelm, testUpdateFuncHelm := setupSidecarTest(updatedSidecarConfigHelm)
	testUpdateConfigTF, testUpdateFuncTF := setupSidecarTest(updatedSidecarConfigTF)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdateConfigDocker,
				Check:  testUpdateFuncDocker,
			},
			{
				Config: testUpdateConfigHelm,
				Check:  testUpdateFuncHelm,
			},
			{
				Config: testUpdateConfigTF,
				Check:  testUpdateFuncTF,
			},
		},
	})
}

func setupSidecarTest(integrationData SidecarData) (string, resource.TestCheckFunc) {
	configuration := formatSidecarDataIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_sidecar.test_sidecar", "name", integrationData.Name),
		resource.TestCheckResourceAttr("cyral_sidecar.test_sidecar", "deployment_method", integrationData.SidecarProperty.DeploymentMethod),
	)

	return configuration, testFunction
}

func formatSidecarDataIntoConfig(data SidecarData) string {
	return fmt.Sprintf(`
      resource "cyral_sidecar" "test_sidecar" {
      	name = "%s"
      	deployment_method = "%s"
		labels = ["%s"]
      }`, data.Name, data.SidecarProperty.DeploymentMethod, data.Labels[0])
}
