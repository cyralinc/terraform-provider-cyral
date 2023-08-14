package cyral

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	datalabelDataSourceName = "data-datalabel"
)

func datalabelDataSourceTestDataLabels() []*model.DataLabel {
	return []*model.DataLabel{
		{
			Name:        accTestName(datalabelDataSourceName, "1"),
			Type:        model.DataLabelTypeCustom,
			Description: "description-1",
			Tags:        []string{"tag-1", "tag-2"},
		},
		{
			Name:        accTestName(datalabelDataSourceName, "2"),
			Type:        model.DataLabelTypeCustom,
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
		"type_filter_predefined", dataLabels, "", model.DataLabelTypePredefined)
	testConfigTypeFilterCustom, testFuncTypeFilterCustom := testDatalabelDataSource(t,
		"type_filter_custom", dataLabels, "", model.DataLabelTypeCustom)

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
	dataLabels []*model.DataLabel,
	nameFilter, typeFilter string,
) (
	string, resource.TestCheckFunc,
) {
	return testDatalabelDataSourceConfig(dsourceName, dataLabels, nameFilter, typeFilter),
		testDatalabelDataSourceChecks(t, dsourceName, dataLabels, nameFilter, typeFilter)
}

func testDatalabelDataSourceConfig(
	dsourceName string,
	dataLabels []*model.DataLabel,
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
		return resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(dataSourceFullName,
				"datalabel_list.#",
				notZeroRegex(),
			),
			dsourceCheckTypeFilter(
				dataSourceFullName,
				"datalabel_list.%d.type",
				typeFilter,
			),
		)
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

func filterDataLabels(dataLabels []*model.DataLabel, nameFilter, typeFilter string) []*DataLabel {
	var filteredDataLabels []*model.DataLabel
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
		depends_on = %s
	}`, dsourceName, nameFilter, typeFilter, listToStr(dependsOn))
}
