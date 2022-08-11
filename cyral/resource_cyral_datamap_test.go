package cyral

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	datamapResourceName = "datamap"
)

func datamapResourceTestRepoConfig_MySQL() string {
	return formatBasicRepositoryIntoConfig(
		"repo_1",
		datamapResourceTestRepoName_MySQL(),
		"mysql",
		"some-host.com",
		3306,
	)
}

func datamapResourceTestRepoName_MySQL() string {
	return accTestName(datamapResourceName, "mysql")
}

func datamapResourceTestRepoConfig_MariaDB() string {
	return formatBasicRepositoryIntoConfig(
		"repo_2",
		datamapResourceTestRepoName_MariaDB(),
		"mariadb",
		"some-host.com",
		1234,
	)
}

func datamapResourceTestRepoName_MariaDB() string {
	return accTestName(datamapResourceName, "mariadb")
}

func TestAccDatamapResource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
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
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"mapping.", "last_updated"},
				ResourceName:            "cyral_datamap.datamap_1",
			},
		},
	})
}

func testAccDatamapConfig_InitialConfig() string {
	var config string
	config += datamapResourceTestRepoConfig_MySQL()
	config += `
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
	return config
}

func testAccDatamapCheck_InitialConfig() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckTypeSetElemNestedAttrs("cyral_datamap.datamap_1",
			"mapping.*",
			map[string]string{
				"label":                        "CNN",
				"data_location.#":              "1",
				"data_location.0.repo":         datamapResourceTestRepoName_MySQL(),
				"data_location.0.attributes.#": "4",
			},
		),
		resource.TestCheckNoResourceAttr("cyral_datamap.datamap_1",
			"last_updated",
		),
	)
}

func testAccDatamapConfig_UpdatedConfig() string {
	var config string
	config += datamapResourceTestRepoConfig_MySQL()
	config += datamapResourceTestRepoConfig_MariaDB()
	config += `
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
	return config
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
				"data_location.0.repo":         datamapResourceTestRepoName_MariaDB(),
				"data_location.0.attributes.#": "3",
			},
		),
		resource.TestCheckResourceAttrSet("cyral_datamap.datamap_1",
			"last_updated",
		),
	)
}
