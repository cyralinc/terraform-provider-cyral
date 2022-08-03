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
	testCustomLabel    = "test-tf-custom-label"
)

func initialDataMapConfig() *DataMap {
	return &DataMap{
		Labels: map[string]*DataMapMapping{
			predefinedLabelCCN: &DataMapMapping{
				Attributes: []string{
					"schema1.table1.col1",
				},
			},
			predefinedLabelSSN: &DataMapMapping{
				Attributes: []string{
					// Important to have inverse order here,
					// to test that the resource diff is
					// consitent.
					"schema1.table1.col4",
					"schema1.table1.col3",
					"schema1.table1.col2",
				},
			},
		},
	}
}

func updatedDataMapConfigRemoveAttribute() *DataMap {
	return &DataMap{
		Labels: map[string]*DataMapMapping{
			predefinedLabelCCN: &DataMapMapping{
				Attributes: []string{
					"schema1.table1.col1",
				},
			},
			predefinedLabelSSN: &DataMapMapping{
				Attributes: []string{
					"schema1.table1.col3",
					"schema1.table1.col2",
				},
			},
		},
	}
}

func updatedDataMapConfigRemoveLabel() *DataMap {
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

func dataMapConfigWithDataLabel() (*DataMap, *DataLabel) {
	return &DataMap{
			Labels: map[string]*DataMapMapping{
				testCustomLabel: &DataMapMapping{
					Attributes: []string{
						"schema1.table1.col1",
					},
				},
			},
		}, &DataLabel{
			Name:        testCustomLabel,
			Description: "custom-label-description",
			Tags:        []string{"tag1"},
		}
}

func TestAccRepositoryDatamapResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			testRepositoryDatamapInitialConfig(t),
			testRepositoryDatamapUpdatedConfigRemoveAttribute(t),
			testRepositoryDatamapUpdatedConfigRemoveLabel(t),
			testRepositoryDatamapWithDataLabel(t),
		},
	})
}

func testRepositoryDatamapInitialConfig(t *testing.T) resource.TestStep {
	config := formatDataMapIntoConfig(t, initialDataMapConfig())
	check := setupRepositoryDatamapTestFunc(t, initialDataMapConfig())
	return resource.TestStep{Config: config, Check: check}
}

func testRepositoryDatamapUpdatedConfigRemoveAttribute(t *testing.T) resource.TestStep {
	config := formatDataMapIntoConfig(t, updatedDataMapConfigRemoveAttribute())
	check := setupRepositoryDatamapTestFunc(t, updatedDataMapConfigRemoveAttribute())
	return resource.TestStep{Config: config, Check: check}
}

func testRepositoryDatamapUpdatedConfigRemoveLabel(t *testing.T) resource.TestStep {
	config := formatDataMapIntoConfig(t, updatedDataMapConfigRemoveLabel())
	check := setupRepositoryDatamapTestFunc(t, updatedDataMapConfigRemoveLabel())
	return resource.TestStep{Config: config, Check: check}
}

func testRepositoryDatamapWithDataLabel(t *testing.T) resource.TestStep {
	configDM, configDL := dataMapConfigWithDataLabel()
	tfConfig := (formatDataMapIntoConfig(t, configDM) +
		formatDataLabelIntoConfig(t, configDL))
	check := setupRepositoryDatamapTestFunc(t, configDM)
	return resource.TestStep{Config: tfConfig, Check: check}
}

func setupRepositoryDatamapTestFunc(t *testing.T, dataMap *DataMap) resource.TestCheckFunc {
	testFunctions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(
			"cyral_repository_datamap.test_repository_datamap", "repository_id",
			"cyral_repository.test_repository", "id"),
	}

	require.NotNil(t, dataMap.Labels)
	idxMapping := 0
	sortedLabels := dataMapSortedLabels(dataMap)
	for _, label := range sortedLabels {
		mapping := dataMap.Labels[label]

		testFunctions = append(testFunctions, resource.TestCheckResourceAttr(
			"cyral_repository_datamap.test_repository_datamap",
			fmt.Sprintf("mapping.%d.label", idxMapping), label))
		testFunctions = append(testFunctions, resource.TestCheckResourceAttr(
			"cyral_repository_datamap.test_repository_datamap",
			fmt.Sprintf("mapping.%d.attributes.#", idxMapping),
			fmt.Sprintf("%d", len(mapping.Attributes))))

		idxMapping++
	}

	testFunction := resource.ComposeTestCheckFunc(testFunctions...)

	return testFunction
}

func formatDataMapIntoConfig(t *testing.T, dataMap *DataMap) string {
	dependsOnStr := ""
	mappingsStr := ""
	sortedLabels := dataMapSortedLabels(dataMap)
	for _, label := range sortedLabels {
		mapping := dataMap.Labels[label]

		require.NotNil(t, mapping)
		require.NotNil(t, mapping.Attributes)

		mappingsStr += fmt.Sprintf(`
		mapping {
			label = "%s"
			attributes = [%s]
		}`, label, formatAttributes(mapping.Attributes))

		if label == testCustomLabel {
			// If there is a custom label in the configuration, we
			// need to delete the data map first, otherwise the
			// label cannot be deleted. The depends_on Terraform
			// meta-argument forces the right deletion order.
			dependsOnStr = "depends_on = [cyral_datalabel.test_datalabel]"
		}
	}
	require.NotEmpty(t, mappingsStr)

	config := fmt.Sprintf(`
	resource "cyral_repository" "test_repository" {
		type  = "sqlserver"
		host  = "localhost"
		port  = 1433
		name  = "tf-test-sqlserver-1"
		labels = ["repo-label1", "repo-label2"]
	}

	resource "cyral_repository_datamap" "test_repository_datamap" {
		%s
		repository_id = cyral_repository.test_repository.id
		%s
	}`, dependsOnStr, mappingsStr)

	fmt.Printf("[DEBUG] Config: %s\n", config)

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
