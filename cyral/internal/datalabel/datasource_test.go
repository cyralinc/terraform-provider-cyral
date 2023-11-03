package datalabel_test

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/datalabel"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	datalabelDataSourceName = "data-datalabel"
)

func datalabelDataSourceTestDataLabels() []*datalabel.DataLabel {
	return []*datalabel.DataLabel{
		{
			Name:        utils.AccTestName(datalabelDataSourceName, "1"),
			Type:        datalabel.Custom,
			Description: "description-1",
			Tags:        []string{"tag-1", "tag-2"},
		},
		{
			Name:        utils.AccTestName(datalabelDataSourceName, "2"),
			Type:        datalabel.Custom,
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
		"type_filter_predefined", dataLabels, "", string(datalabel.Predefined))
	testConfigTypeFilterCustom, testFuncTypeFilterCustom := testDatalabelDataSource(t,
		"type_filter_custom", dataLabels, "", string(datalabel.Custom))

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"cyral": func() (*schema.Provider, error) {
				return provider.Provider(), nil
			},
		},
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
	dataLabels []*datalabel.DataLabel,
	nameFilter, typeFilter string,
) (
	string, resource.TestCheckFunc,
) {
	return testDatalabelDataSourceConfig(dsourceName, dataLabels, nameFilter, typeFilter),
		testDatalabelDataSourceChecks(t, dsourceName, dataLabels, nameFilter, typeFilter)
}

func testDatalabelDataSourceConfig(
	dsourceName string,
	dataLabels []*datalabel.DataLabel,
	nameFilter, typeFilter string,
) string {
	var config string
	var dependsOn []string
	for i, dataLabel := range dataLabels {
		ruleType, ruleCode, ruleStatus := "", "", ""
		if dataLabel.ClassificationRule != nil {
			ruleType = dataLabel.ClassificationRule.RuleType
			ruleCode = dataLabel.ClassificationRule.RuleCode
			ruleStatus = dataLabel.ClassificationRule.RuleStatus
		}
		resName := fmt.Sprintf("test_datalabel_%d", i)
		config += utils.FormatDataLabelIntoConfig(resName, dataLabel.Name, dataLabel.Description,
			ruleType, ruleCode, ruleStatus, dataLabel.Tags)
		dependsOn = append(dependsOn, utils.DatalabelConfigResourceFullName(resName))
	}
	config += datalabelDataSourceConfig(dsourceName, nameFilter, typeFilter, dependsOn)

	return config
}

func testDatalabelDataSourceChecks(
	t *testing.T,
	dsourceName string,
	dataLabels []*datalabel.DataLabel,
	nameFilter, typeFilter string,
) resource.TestCheckFunc {
	dataSourceFullName := fmt.Sprintf("data.cyral_datalabel.%s", dsourceName)

	if nameFilter == "" {
		return resource.ComposeTestCheckFunc(
			resource.TestMatchResourceAttr(dataSourceFullName,
				"datalabel_list.#",
				utils.NotZeroRegex(),
			),
			utils.DSourceCheckTypeFilter(
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

func filterDataLabels(dataLabels []*datalabel.DataLabel, nameFilter, typeFilter string) []*datalabel.DataLabel {
	var filteredDataLabels []*datalabel.DataLabel
	for _, dataLabel := range dataLabels {
		if (nameFilter == "" || dataLabel.Name == nameFilter) &&
			(typeFilter == "" || string(dataLabel.Type) == typeFilter) {
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
	}`, dsourceName, nameFilter, typeFilter, utils.ListToStr(dependsOn))
}
