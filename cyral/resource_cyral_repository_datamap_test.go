package cyral

import (
	"fmt"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
)

const (
	predefinedLabelCCN = "CCN"
	predefinedLabelSSN = "SSN"
)

func repositoryDatamapSampleRepositoryConfig(resName string) string {
	return formatBasicRepositoryIntoConfig(
		basicRepositoryResName,
		accTestName("repository-datamap", resName),
		"sqlserver",
		"localhost",
		1433,
	)
}

func testRepositoryDatamapCustomLabel() string {
	return accTestName("repository-datamap", "custom-label")
}

func initialDataMapConfigRemoveMapping() *DataMap {
	return &DataMap{
		Labels: map[string]*DataMapMapping{
			predefinedLabelCCN: &DataMapMapping{
				Attributes: []string{
					"schema1.table1.col1",
				},
			},
			predefinedLabelSSN: &DataMapMapping{
				Attributes: []string{
					"schema1.table1.col2",
				},
			},
		},
	}
}

func updatedDataMapConfigRemoveMapping() *DataMap {
	return &DataMap{
		Labels: map[string]*DataMapMapping{
			predefinedLabelSSN: &DataMapMapping{
				Attributes: []string{
					"schema1.table1.col2",
				},
			},
		},
	}
}

func initialDataMapConfigRemoveAttribute() *DataMap {
	return &DataMap{
		Labels: map[string]*DataMapMapping{
			predefinedLabelSSN: &DataMapMapping{
				Attributes: []string{
					"a.b.c",
					"b.c.d",
				},
			},
		},
	}
}

func updatedDataMapConfigRemoveAttribute() *DataMap {
	return &DataMap{
		Labels: map[string]*DataMapMapping{
			predefinedLabelSSN: &DataMapMapping{
				Attributes: []string{
					"a.b.c",
				},
			},
		},
	}
}

func dataMapConfigWithDataLabel() (*DataMap, *DataLabel) {
	return &DataMap{
			Labels: map[string]*DataMapMapping{
				testRepositoryDatamapCustomLabel(): &DataMapMapping{
					Attributes: []string{
						"schema1.table1.col1",
					},
				},
			},
		}, &DataLabel{
			Name:        testRepositoryDatamapCustomLabel(),
			Description: "custom-label-description",
			Tags:        []string{"tag1"},
		}
}

func TestAccRepositoryDatamapResource(t *testing.T) {
	importStateResName := "cyral_repository_datamap.test_with_datalabel"

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			testRepositoryDatamapInitialConfigRemoveMapping(t),
			testRepositoryDatamapUpdatedConfigRemoveMapping(t),
			testRepositoryDatamapInitialConfigRemoveAttribute(t),
			testRepositoryDatamapUpdatedConfigRemoveAttribute(t),
			testRepositoryDatamapWithDataLabel(t),
			testRepositoryDatamapImport(importStateResName),
		},
	})
}

func testRepositoryDatamapInitialConfigRemoveMapping(t *testing.T) resource.TestStep {
	resName := "test_remove_mapping"
	config := initialDataMapConfigRemoveMapping()
	var tfConfig string
	tfConfig += repositoryDatamapSampleRepositoryConfig(resName)
	tfConfig += formatDataMapIntoConfig(resName, basicRepositoryID, config)
	check := setupRepositoryDatamapTestFunc(t, resName, config)
	return resource.TestStep{Config: tfConfig, Check: check}
}

func testRepositoryDatamapUpdatedConfigRemoveMapping(t *testing.T) resource.TestStep {
	resName := "test_remove_mapping"
	config := updatedDataMapConfigRemoveMapping()
	var tfConfig string
	tfConfig += repositoryDatamapSampleRepositoryConfig(resName)
	tfConfig += formatDataMapIntoConfig(resName, basicRepositoryID, config)
	check := setupRepositoryDatamapTestFunc(t, resName, config)
	return resource.TestStep{Config: tfConfig, Check: check}
}

func testRepositoryDatamapInitialConfigRemoveAttribute(t *testing.T) resource.TestStep {
	resName := "test_remove_attribute"
	config := initialDataMapConfigRemoveAttribute()
	var tfConfig string
	tfConfig += repositoryDatamapSampleRepositoryConfig(resName)
	tfConfig += formatDataMapIntoConfig(resName, basicRepositoryID, config)
	check := setupRepositoryDatamapTestFunc(t, resName, config)
	return resource.TestStep{Config: tfConfig, Check: check}
}

func testRepositoryDatamapUpdatedConfigRemoveAttribute(t *testing.T) resource.TestStep {
	resName := "test_remove_attribute"
	config := updatedDataMapConfigRemoveAttribute()
	var tfConfig string
	tfConfig += repositoryDatamapSampleRepositoryConfig(resName)
	tfConfig += formatDataMapIntoConfig(resName, basicRepositoryID, config)
	check := setupRepositoryDatamapTestFunc(t, resName, config)
	return resource.TestStep{Config: tfConfig, Check: check}
}

func testRepositoryDatamapWithDataLabel(t *testing.T) resource.TestStep {
	resName := "test_with_datalabel"
	configDM, configDL := dataMapConfigWithDataLabel()
	var tfConfig string
	tfConfig += repositoryDatamapSampleRepositoryConfig(resName)
	tfConfig += (formatDataMapIntoConfig(resName, basicRepositoryID, configDM) +
		formatDataLabelIntoConfig(configDL.Name, configDL))
	check := setupRepositoryDatamapTestFunc(t, resName, configDM)
	return resource.TestStep{Config: tfConfig, Check: check}
}

func testRepositoryDatamapImport(importStateResName string) resource.TestStep {
	return resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		// TODO: Properly verify mappings -aholmquist 2022-08-05
		ImportStateVerifyIgnore: []string{"mapping."},
		ResourceName:            importStateResName,
	}
}

func setupRepositoryDatamapTestFunc(
	t *testing.T,
	resName string,
	dataMap *DataMap,
) resource.TestCheckFunc {
	resFullName := fmt.Sprintf("cyral_repository_datamap.%s", resName)

	testFunctions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(
			resFullName, "repository_id",
			fmt.Sprintf("cyral_repository.%s", basicRepositoryResName), "id",
		),
	}

	require.NotNil(t, dataMap.Labels)
	idxMapping := 0
	sortedLabels := dataMapSortedLabels(dataMap)
	for _, label := range sortedLabels {
		mapping := dataMap.Labels[label]

		testFunctions = append(testFunctions, resource.TestCheckResourceAttr(
			resFullName,
			fmt.Sprintf("mapping.%d.label", idxMapping), label))
		testFunctions = append(testFunctions, resource.TestCheckResourceAttr(
			resFullName,
			fmt.Sprintf("mapping.%d.attributes.#", idxMapping),
			fmt.Sprintf("%d", len(mapping.Attributes))))

		idxMapping++
	}

	testFunction := resource.ComposeTestCheckFunc(testFunctions...)

	return testFunction
}

func formatDataMapIntoConfig(
	resName, repositoryID string,
	dataMap *DataMap,
) string {
	dependsOnStr := ""
	mappingsStr := ""
	sortedLabels := dataMapSortedLabels(dataMap)
	for _, label := range sortedLabels {
		mapping := dataMap.Labels[label]

		mappingsStr += fmt.Sprintf(`
		mapping {
			label = "%s"
			attributes = [%s]
		}`, label, listToStr(mapping.Attributes))

		if label == testRepositoryDatamapCustomLabel() {
			// If there is a custom label in the configuration, we
			// need to delete the data map first, otherwise the
			// label cannot be deleted. The depends_on Terraform
			// meta-argument forces the right deletion order.
			dependsOnStr = fmt.Sprintf("depends_on = [%s]",
				datalabelConfigResourceFullName(label))
		}
	}

	config := fmt.Sprintf(`
	resource "cyral_repository_datamap" "%s" {
		repository_id = %s
		%s
		%s
	}`, resName, repositoryID, mappingsStr, dependsOnStr)

	return config
}

// dataMapSortedLabels exists to allow construction and checking of terraform
// configurations with the data map following the same order of the mappings.
func dataMapSortedLabels(dataMap *DataMap) []string {
	var labels []string
	for label, _ := range dataMap.Labels {
		labels = append(labels, label)
	}
	sort.Strings(labels)
	return labels
}
