package cyral

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func datalabelDataSourceTestDataLabels() []*DataLabel {
	return []*DataLabel{
		{
			Name:        "tf-provider-test-datalabel-1",
			Type:        dataLabelTypeCustom,
			Description: "description-1",
			Tags:        []string{"tag-1", "tag-2"},
		},
		{
			Name:        "tf-provider-test-datalabel-2",
			Type:        dataLabelTypeCustom,
			Description: "description-2",
			Tags:        []string{"tag-3"},
		},
	}
}

func TestAccDatalabelDataSource(t *testing.T) {
	dataLabels := datalabelDataSourceTestDataLabels()

	testConfigNameFilter1, testFuncNameFilter1 := testDatalabelDataSource(t,
		"name_filter_1", dataLabels, dataLabels[0].Name, "")
	testConfigNameFilter2, testFuncNameFilter2 := testDatalabelDataSource(t,
		"name_filter_2", dataLabels, dataLabels[1].Name, "")
	testConfigTypeFilterPredefined, testFuncTypeFilterPredefined := testDatalabelDataSource(t,
		"type_filter_predefined", dataLabels, "", dataLabelTypePredefined)
	testConfigTypeFilterCustom, testFuncTypeFilterCustom := testDatalabelDataSource(t,
		"type_filter_custom", dataLabels, "", dataLabelTypeCustom)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigNameFilter1,
				Check:  testFuncNameFilter1,
			},
			{
				Config: testConfigNameFilter2,
				Check:  testFuncNameFilter2,
			},
			{
				Config: testConfigTypeFilterPredefined,
				Check:  testFuncTypeFilterPredefined,
			},
			{
				Config: testConfigTypeFilterCustom,
				Check:  testFuncTypeFilterCustom,
			},
		},
	})
}

func testDatalabelDataSource(
	t *testing.T,
	dsourceName string,
	dataLabels []*DataLabel,
	nameFilter, typeFilter string,
) (
	string, resource.TestCheckFunc,
) {
	return testDatalabelDataSourceConfig(dsourceName, dataLabels, nameFilter, typeFilter),
		testDatalabelDataSourceChecks(t, dsourceName, dataLabels, nameFilter, typeFilter)
}

func testDatalabelDataSourceConfig(
	dsourceName string,
	dataLabels []*DataLabel,
	nameFilter,
	typeFilter string,
) string {
	var config string
	var dependsOn []string
	for i, dataLabel := range dataLabels {
		resName := fmt.Sprintf("test_datalabel_%d", i)
		config += formatDataLabelIntoConfig(resName, dataLabel)
		dependsOn = append(dependsOn, datalabelConfigResourceFullName(resName))
	}
	config += datalabelDataSourceConfig(dsourceName, nameFilter, typeFilter, dependsOn)

	return config
}

func testDatalabelDataSourceChecks(
	t *testing.T,
	dsourceName string,
	dataLabels []*DataLabel,
	nameFilter, typeFilter string,
) resource.TestCheckFunc {
	dataSourceFullName := fmt.Sprintf("data.cyral_datalabel.%s", dsourceName)

	if nameFilter == "" {
		// In this case, we might encounter labels that we did not
		// create in the control plane, which can lead to
		// non-deterministic tests. We just check that the actual type
		// of the label matches the expected type.
		notZeroRegex := regexp.MustCompile("^[0-9]*[^0]$")
		checkFuncs := []resource.TestCheckFunc{
			resource.TestMatchResourceAttr(dataSourceFullName,
				"datalabel_list.#",
				notZeroRegex,
			),
			func(s *terraform.State) error {
				ds, ok := s.RootModule().Resources[dataSourceFullName]
				if !ok {
					return fmt.Errorf("Not found: %s", dataSourceFullName)
				}
				numDataLabels, err := strconv.Atoi(ds.Primary.Attributes["datalabel_list.#"])
				if err != nil {
					return err
				}
				for i := 0; i < numDataLabels; i++ {
					nameLocation := fmt.Sprintf("datalabel_list.%d.type", i)
					actualType := ds.Primary.Attributes[nameLocation]
					if actualType != typeFilter {
						return fmt.Errorf("Expected all labels to have "+
							"type equal to type filter %q, but got: "+
							"%s", typeFilter, actualType)
					}
				}
				return nil
			},
		}

		return resource.ComposeTestCheckFunc(checkFuncs...)

	}

	var checkFuncs []resource.TestCheckFunc
	filteredDataLabels := filterDataLabels(dataLabels, nameFilter, typeFilter)
	for i, label := range filteredDataLabels {
		checkFuncs = append(checkFuncs,
			resource.TestCheckResourceAttr(dataSourceFullName,
				fmt.Sprintf("datalabel_list.%d.name", i),
				label.Name,
			),
			resource.TestCheckResourceAttr(dataSourceFullName,
				fmt.Sprintf("datalabel_list.%d.description", i),
				label.Description,
			),
			resource.TestCheckResourceAttr(dataSourceFullName,
				fmt.Sprintf("datalabel_list.%d.tags.#", i),
				fmt.Sprintf("%d", len(label.Tags)),
			),
			resource.TestCheckResourceAttr(dataSourceFullName,
				fmt.Sprintf("datalabel_list.%d.implicit", i),
				"false",
			),
		)
	}

	return resource.ComposeTestCheckFunc(checkFuncs...)
}

func filterDataLabels(dataLabels []*DataLabel, nameFilter, typeFilter string) []*DataLabel {
	var filteredDataLabels []*DataLabel
	for _, dataLabel := range dataLabels {
		if (nameFilter == "" || dataLabel.Name == nameFilter) &&
			(typeFilter == "" || dataLabel.Type == typeFilter) {
			filteredDataLabels = append(filteredDataLabels, dataLabel)
		}
	}
	return filteredDataLabels
}

func datalabelDataSourceConfig(dsourceName, nameFilter, typeFilter string, dependsOn []string) string {
	return fmt.Sprintf(`
	data "cyral_datalabel" "%s" {
		name = "%s"
		type = "%s"
		depends_on = [%s]
	}`, dsourceName, nameFilter, typeFilter, formatAttributes(dependsOn))
}
