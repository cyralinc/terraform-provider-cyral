package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type DataMapConfig struct {
	Label        string
	DataLocation []RepoAttrs
}

var initialDataMapConfig DataMapConfig = DataMapConfig{
	Label: "CNN",
	DataLocation: []RepoAttrs{
		{
			Name:       "Repository",
			Attributes: []string{"applications.customers.credit_card_number"},
		},
	},
}

var updatedDataMapConfig DataMapConfig = DataMapConfig{
	Label: "CNN-Updated",
	DataLocation: []RepoAttrs{
		{
			Name:       "Repository-updated",
			Attributes: []string{"applications.customers.credit_card_number_updated"},
		},
	},
}

// This is loosely based on this example:
// https://github.com/hashicorp/terraform-provider-vault/blob/master/vault/resource_azure_secret_backend_role_test.go
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
		resource.TestCheckTypeSetElemNestedAttrs("cyral_datamap.test", "mapping.*", map[string]string{
			"label": integrationData.Label,
		}),
		resource.TestCheckTypeSetElemNestedAttrs("cyral_datamap.test", "mapping.0.data_location.*", map[string]string{
			"repo": integrationData.DataLocation[0].Name,
		}),
		resource.TestCheckTypeSetElemAttr("cyral_datamap.test", "mapping.0.data_location.0.attributes.*", integrationData.DataLocation[0].Attributes[0]),
		// resource.TestCheckResourceAttr("cyral_datamap.test", "mapping.#", "1"),
		// resource.TestCheckResourceAttr("cyral_datamap.test", "mapping.0.label", integrationData.Label),
	)

	return configuration, testFunction
}

func formatDataMapIntoConfig(data DataMapConfig) string {
	return fmt.Sprintf(`
	resource "cyral_datamap" "test" {
		mapping {
			label = "%s"
			data_location {
				repo = "%s"
				attributes = [%s]
			}
		}
	}`, data.Label, data.DataLocation[0].Name, formatAttibutes(data.DataLocation[0].Attributes))
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
