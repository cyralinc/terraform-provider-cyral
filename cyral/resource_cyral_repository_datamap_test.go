package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
)

const (
	testPredefinedLabel = "SSN"
)

func initialDataMapConfig() *DataMap {
	return &DataMap{
		Labels: map[string]*DataMapMapping{
			testPredefinedLabel: &DataMapMapping{
				Attributes: []string{
					"schema1.table1.col1",
					"schema1.table1.col2",
				},
			},
		},
	}
}

func updateDataMapConfig() *DataMap {
	return &DataMap{
		Labels: map[string]*DataMapMapping{
			testPredefinedLabel: &DataMapMapping{
				Attributes: []string{
					"schema1.table1.col3",
				},
			},
		},
	}
}

func TestAccRepositoryDatamapResource(t *testing.T) {
	testInitialConfig, testInitialFunc := setupRepositoryDatamapTest(t, initialDataMapConfig())
	testUpdateConfig, testUpdateFunc := setupRepositoryDatamapTest(t, updateDataMapConfig())

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testInitialConfig,
				Check:  testInitialFunc,
			},
			{
				Config: testUpdateConfig,
				Check:  testUpdateFunc,
			},
		},
	})
}

func setupRepositoryDatamapTest(t *testing.T, dataMap *DataMap) (string, resource.TestCheckFunc) {
	configuration := formatRepositoryDatamapIntoConfig(t, dataMap)

	testFunctions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(
			"cyral_repository_datamap.test_repository_datamap", "repo_id",
			"cyral_sidecar.test_sidecar", "id"),
	}

	require.NotNil(t, dataMap.Labels)
	idxMapping := 0
	for label, mapping := range dataMap.Labels {
		matchingMap := map[string]string{"label": label}
		for i, attribute := range mapping.Attributes {
			matchingMap[fmt.Sprintf("attributes.%d", i)] = attribute
		}
		testFunctions = append(testFunctions, resource.TestCheckTypeSetElemNestedAttrs(
			"cyral_repository_datamap.test_repository_datamap", fmt.Sprintf("mapping.%d", idxMapping),
			matchingMap))
		idxMapping++
	}

	testFunction := resource.ComposeTestCheckFunc(testFunctions...)

	return configuration, testFunction
}

func formatRepositoryDatamapIntoConfig(t *testing.T, dataMap *DataMap) string {
	mappingsStr := ""
	for label, mapping := range dataMap.Labels {
		require.NotNil(t, mapping)
		require.NotNil(t, mapping.Attributes)

		mappingsStr += fmt.Sprintf(`
		mapping {
			label = "%s"
			attributes = [%s]
		}`, label, formatAttributes(mapping.Attributes))
	}
	require.NotEmpty(t, mappingsStr)

	return fmt.Sprintf(`
	resource "cyral_repository" "repo1" {
		type  = "sqlserver"
		host  = "localhost"
		port  = 1433
		name  = "sqlserver-1"
		labels = ["repo-label1", "repo-label2"]
	}

	resource "cyral_repository_datamap" "test_repository_datamap" {
		repo_id = cyral_repository.repo1.id
		%s
	}`, mappingsStr)
}
