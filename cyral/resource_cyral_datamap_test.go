package cyral

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatamapResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatamapConfig_InitialConfig(),
				Check:  testAccDatamapCheck_InitialConfig(),
			},
			{
				Config: testAccDatamapConfig_UpdatedConfig(),
				Check:  testAccDatamapCheck_UpdatedConfig(),
			},
		},
	})
}

func testAccDatamapConfig_InitialConfig() string {
	return `
	resource "cyral_repository" "repo_1" {
		type = "mysql"
		host = "some-host.com"
		port = 3306
		name = "tf-test-mysql"
	}

	resource "cyral_datamap" "datamap_1" {
		mapping {
			label = "CNN"
			data_location {
				repo       = cyral_repository.repo_1.name
				attributes = [
					"schema.table.attribute-a",
					"schema.table.attribute-b",
					"schema.table.attribute-c",
					"schema.table.attribute-d",
				]
			}
		}
	}`
}

func testAccDatamapCheck_InitialConfig() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckTypeSetElemNestedAttrs("cyral_datamap.datamap_1",
			"mapping.*",
			map[string]string{
				"label":                        "CNN",
				"data_location.#":              "1",
				"data_location.0.repo":         "tf-test-mysql",
				"data_location.0.attributes.#": "4",
			},
		),
		resource.TestCheckNoResourceAttr("cyral_datamap.datamap_1",
			"last_updated",
		),
	)
}

func testAccDatamapConfig_UpdatedConfig() string {
	return `
	resource "cyral_repository" "repo_1" {
		type = "mysql"
		host = "some-host.com"
		port = 3306
		name = "tf-test-mysql"
	}

	resource "cyral_repository" "repo_2" {
		type = "mariadb"
		host = "some-host.com"
		port = 1234
		name = "tf-test-mariadb"
	}

	resource "cyral_datamap" "datamap_1" {
		mapping {
			label = "CNN-UPDATED"
			data_location {
				repo       = cyral_repository.repo_1.name
				attributes = [
					"schema.table.attribute-a",
					"schema.table.attribute-b",
					"schema.table.attribute-c",
					"schema.table.attribute-d",
				]
			}
			data_location {
				repo       = cyral_repository.repo_2.name
				attributes = [
					"schema.table.attribute-1",
					"schema.table.attribute-2",
					"schema.table.attribute-3",
					"schema.table.attribute-4",
				]
			}
		}
		mapping {
			label = "EMAIL"
			data_location {
				repo       = cyral_repository.repo_2.name
				attributes = [
					"schema.table.attribute-5",
					"schema.table.attribute-6",
					"schema.table.attribute-7",
				]
			}
		}
	}`
}

func testAccDatamapCheck_UpdatedConfig() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckTypeSetElemNestedAttrs("cyral_datamap.datamap_1",
			"mapping.*",
			map[string]string{
				"label":           "CNN-UPDATED",
				"data_location.#": "2",
			},
		),
		resource.TestCheckTypeSetElemNestedAttrs("cyral_datamap.datamap_1",
			"mapping.*",
			map[string]string{
				"label":                        "EMAIL",
				"data_location.#":              "1",
				"data_location.0.repo":         "tf-test-mariadb",
				"data_location.0.attributes.#": "3",
			},
		),
		resource.TestCheckResourceAttrSet("cyral_datamap.datamap_1",
			"last_updated",
		),
	)
}
