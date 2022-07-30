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

func updateDataLabelConfig() *DataLabel {
	return &DataLabel{
		Name:        "test-tf-label2",
		Description: "label2-description",
		Tags:        []string{"tag1", "tag2"},
	}
}

func TestAccDatalabelResource(t *testing.T) {
	testInitialConfig, testInitialFunc := setupDatalabelTest(t,
		initialDataLabelConfig())
	testUpdateConfig, testUpdateFunc := setupDatalabelTest(t,
		updateDataLabelConfig())

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
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "cyral_datalabel.test_datalabel",
			},
		},
	})
}

func setupDatalabelTest(t *testing.T, dataLabel *DataLabel) (string, resource.TestCheckFunc) {
	configuration := formatDataLabelIntoConfig(t, dataLabel)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("cyral_datalabel.test_datalabel",
			"name", dataLabel.Name),
		resource.TestCheckResourceAttr("cyral_datalabel.test_datalabel",
			"description", dataLabel.Description),
		resource.TestCheckResourceAttr("cyral_datalabel.test_datalabel",
			"tags.#", "2"),
	)

	return configuration, testFunction
}

func formatDataLabelIntoConfig(t *testing.T, dataLabel *DataLabel) string {
	return fmt.Sprintf(`
	resource "cyral_datalabel" "test_datalabel" {
		name  = "%s"
		description = "%s"
		tags = [%s]
	}`, dataLabel.Name, dataLabel.Description,
		formatAttributes(dataLabel.Tags))
}
