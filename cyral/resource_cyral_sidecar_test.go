package cyral

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var initialSidecarConfig SidecarData = SidecarData{
	Name: "sidecar-test",
	Tags: []string{"tag1", "tag2"},
	SidecarProperty: SidecarProperty{
		DeploymentMethod: "cloudFormation",
	},
}

var updatedSidecarConfigTags SidecarData = SidecarData{
	Name: "sidecar-test-updated",
	Tags: []string{"tag1", "tag2-modified", "tag3"},
	SidecarProperty: SidecarProperty{
		DeploymentMethod: "cloudFormation",
	},
}

var updatedSidecarConfigDocker SidecarData = SidecarData{
	Name: "sidecar-test-updated",
	Tags: []string{"tag1", "tag2"},
	SidecarProperty: SidecarProperty{
		DeploymentMethod: "docker",
	},
}

var updatedSidecarConfigHelm SidecarData = SidecarData{
	Name: "sidecar-test-updated",
	Tags: []string{"tag1", "tag2"},
	SidecarProperty: SidecarProperty{
		DeploymentMethod: "helm",
	},
}

var updatedSidecarConfigTF SidecarData = SidecarData{
	Name: "sidecar-test-updated",
	Tags: []string{"tag1", "tag2"},
	SidecarProperty: SidecarProperty{
		DeploymentMethod: "terraform",
	},
}

func TestAccSidecarResource(t *testing.T) {
	testConfig, testFunc := setupSidecarTest(initialSidecarConfig)
	testUpdateConfigTags, testUpdateFuncTags := setupSidecarTest(updatedSidecarConfigTags)
	testUpdateConfigDocker, testUpdateFuncDocker := setupSidecarTest(updatedSidecarConfigDocker)
	testUpdateConfigHelm, testUpdateFuncHelm := setupSidecarTest(updatedSidecarConfigHelm)
	testUpdateConfigTF, testUpdateFuncTF := setupSidecarTest(updatedSidecarConfigTF)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check:  testFunc,
			},
			{
				Config: testUpdateConfigTags,
				Check:  testUpdateFuncTags,
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

func setupSidecarTest(sidecarData SidecarData) (string, resource.TestCheckFunc) {
	deploymentTag := fmt.Sprintf("%s%s", DeploymentPrefix, sidecarData.SidecarProperty.DeploymentMethod)
	tags := append([]string{deploymentTag}, sidecarData.Tags...)

	configuration := formatSidecarDataIntoConfig(sidecarData.Name, tags)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_sidecar.test_sidecar", "name", sidecarData.Name),
		resource.TestCheckResourceAttr("cyral_sidecar.test_sidecar", "tags.0", deploymentTag),
		resource.TestCheckResourceAttr("cyral_sidecar.test_sidecar", "tags.#", fmt.Sprintf("%d", len(tags))),
	)

	return configuration, testFunction
}

func formatSidecarDataIntoConfig(name string, tags []string) string {
	return fmt.Sprintf(`
      resource "cyral_sidecar" "test_sidecar" {
      	name = "%s"
				tags = %s
      }`, name, formatSidecarTags(tags))
}

func formatSidecarTags(tags []string) string {
	return fmt.Sprintf("[\"%s\"]", strings.Join(tags, "\", \""))
}
