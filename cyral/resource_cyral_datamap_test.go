package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type DataMapConfig struct {
	Label string
}

var initialDataMapConfig DataMapConfig = DataMapConfig{
	Label: "CNN",
}
var updatedDataMapConfig DataMapConfig = DataMapConfig{
	Label: "CNN-updated",
}

func TestAccDatamapResource(t *testing.T) {
	testConfig, testFunc := setupDatamapTest(initialDataMapConfig)
	testUpdateConfig, testUpdateFunc := setupDatamapTest(updatedDataMapConfig)

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

func setupDatamapTest(integrationData DataMapConfig) (string, resource.TestCheckFunc) {
	configuration := formatDataMapIntoConfig(integrationData)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckTypeSetElemNestedAttrs("cyral_datamap.test_datamap", "mapping.*", map[string]string{
			"label": integrationData.Label,
		}))

	return configuration, testFunction
}

func formatDataMapIntoConfig(data DataMapConfig) string {
	return fmt.Sprintf(`
		resource "cyral_repository" "tf_test_repository" {
			type = "mysql"
			host = "http://mysql.local/"
			port = 3306
			name = "tf-test-mysql"
	  }
	  
	  resource "cyral_sidecar" "tf_test_sidecar" {
			name = "tf-test-sidecar"
			deployment_method = "cloudFormation"
	  }
	  
	  resource "cyral_repository_binding" "repo_binding" {
			enabled       = true
			repository_id = cyral_repository.tf_test_repository.id
			listener_port = 3307
			sidecar_id    = cyral_sidecar.tf_test_sidecar.id
	  }
	  
	  resource "cyral_datamap" "test_datamap" {
			mapping {
				label = "%s"
				data_location {
				repo       = cyral_repository.tf_test_repository.name
				attributes = ["database.table.column"]
				}
			}
	  }`, data.Label)
}

func formatAttibutes(list []string) string {
	currentResp := fmt.Sprintf("\"%s\"", list[0])
	if len(list) > 1 {
		for _, item := range list[1:] {
			currentResp = fmt.Sprintf("%s, \"%s\"", currentResp, item)
		}
	}
	return currentResp
}
