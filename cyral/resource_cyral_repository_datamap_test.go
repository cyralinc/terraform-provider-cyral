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
			testRepositoryDatamapInitialConfigRemoveMapping(t),
			testRepositoryDatamapUpdatedConfigRemoveMapping(t),
			testRepositoryDatamapInitialConfigRemoveAttribute(t),
			testRepositoryDatamapUpdatedConfigRemoveAttribute(t),
			testRepositoryDatamapWithDataLabel(t),
		},
	})
}

func testRepositoryDatamapInitialConfigRemoveMapping(t *testing.T) resource.TestStep {
	resName := "test_remove_label"
	config := initialDataMapConfigRemoveMapping()
	terrConfig := formatDataMapIntoConfig(t, resName, config)
	check := setupRepositoryDatamapTestFunc(t, resName, config)
	return resource.TestStep{Config: terrConfig, Check: check}
}

func testRepositoryDatamapUpdatedConfigRemoveMapping(t *testing.T) resource.TestStep {
	resName := "test_remove_label"
	config := updatedDataMapConfigRemoveMapping()
	terrConfig := formatDataMapIntoConfig(t, resName, config)
	check := setupRepositoryDatamapTestFunc(t, resName, config)
	return resource.TestStep{Config: terrConfig, Check: check}
}

func testRepositoryDatamapInitialConfigRemoveAttribute(t *testing.T) resource.TestStep {
	resName := "test_remove_attribute"
	config := initialDataMapConfigRemoveAttribute()
	terrConfig := formatDataMapIntoConfig(t, resName, config)
	check := setupRepositoryDatamapTestFunc(t, resName, config)
	return resource.TestStep{Config: terrConfig, Check: check}
}

func testRepositoryDatamapUpdatedConfigRemoveAttribute(t *testing.T) resource.TestStep {
	resName := "test_remove_attribute"
	config := updatedDataMapConfigRemoveAttribute()
	terrConfig := formatDataMapIntoConfig(t, resName, config)
	check := setupRepositoryDatamapTestFunc(t, resName, config)
	return resource.TestStep{Config: terrConfig, Check: check}
}

func testRepositoryDatamapWithDataLabel(t *testing.T) resource.TestStep {
	resName := "test_with_datalabel"
	configDM, configDL := dataMapConfigWithDataLabel()
	tfConfig := (formatDataMapIntoConfig(t, resName, configDM) +
		formatDataLabelIntoConfig(t, configDL))
	check := setupRepositoryDatamapTestFunc(t, resName, configDM)
	return resource.TestStep{Config: tfConfig, Check: check}
}

func setupRepositoryDatamapTestFunc(t *testing.T, resName string, dataMap *DataMap) resource.TestCheckFunc {
	resFullName := fmt.Sprintf("cyral_repository_datamap.%s", resName)

	testFunctions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(
			resFullName, "repository_id",
			"cyral_repository.test_repository", "id"),
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

func formatDataMapIntoConfig(t *testing.T, resName string, dataMap *DataMap) string {
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

	resource "cyral_repository_datamap" "%s" {
		%s
		repository_id = cyral_repository.test_repository.id
		%s
	}`, resName, dependsOnStr, mappingsStr)

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
