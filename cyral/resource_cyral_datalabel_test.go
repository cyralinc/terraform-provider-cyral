package cyral

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func initialDataLabelConfig() *DataLabel {
	return &DataLabel{
		Name:        "test-tf-label1",
		Description: "label1-description",
		Tags:        []string{"tag1", "tag2"},
	}
}

func updatedDataLabelConfig() *DataLabel {
	return &DataLabel{
		Name:        "test-tf-label2",
		Description: "label2-description",
		Tags:        []string{"tag1", "tag2"},
	}
}

func TestAccDatalabelResource(t *testing.T) {
	testInitialConfig, testInitialFunc := setupDatalabelTest(t,
		"datalabel_initial", initialDataLabelConfig())
	testUpdatedConfig, testUpdatedFunc := setupDatalabelTest(t,
		"datalabel_updated", updatedDataLabelConfig())

	resource.Test(t, resource.TestCase{
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
		tags = [%s]
	}`, resName, dataLabel.Name, dataLabel.Description,
		formatAttributes(dataLabel.Tags))
}
