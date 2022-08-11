package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	datalabelResourceName = "datalabel"
)

func initialDataLabelConfig() *DataLabel {
	return &DataLabel{
		Name:        accTestName(datalabelResourceName, "label1"),
		Description: "label1-description",
		Tags:        []string{"tag1", "tag2"},
	}
}

func updatedDataLabelConfig() *DataLabel {
	return &DataLabel{
		Name:        accTestName(datalabelResourceName, "label2"),
		Description: "label2-description",
		Tags:        []string{"tag1", "tag2"},
	}
}

func TestAccDatalabelResource(t *testing.T) {
	testInitialConfig, testInitialFunc := setupDatalabelTest(t,
		"main_test", initialDataLabelConfig())
	testUpdatedConfig, testUpdatedFunc := setupDatalabelTest(t,
		"main_test", updatedDataLabelConfig())

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testInitialConfig,
				Check:  testInitialFunc,
			},
			{
				Config: testUpdatedConfig,
				Check:  testUpdatedFunc,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "cyral_datalabel.main_test",
			},
		},
	})
}

func setupDatalabelTest(t *testing.T, resName string, dataLabel *DataLabel) (string, resource.TestCheckFunc) {
	configuration := formatDataLabelIntoConfig(resName, dataLabel)

	resourceFullName := datalabelConfigResourceFullName(resName)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName,
			"name", dataLabel.Name),
		resource.TestCheckResourceAttr(resourceFullName,
			"description", dataLabel.Description),
		resource.TestCheckResourceAttr(resourceFullName,
			"tags.#", "2"),
	)

	return configuration, testFunction
}

func datalabelConfigResourceFullName(resName string) string {
	return fmt.Sprintf("cyral_datalabel.%s", resName)
}

func formatDataLabelIntoConfig(resName string, dataLabel *DataLabel) string {
	return fmt.Sprintf(`
	resource "cyral_datalabel" "%s" {
		name  = "%s"
		description = "%s"
		tags = %s
	}`, resName, dataLabel.Name, dataLabel.Description,
		listToStr(dataLabel.Tags))
}
