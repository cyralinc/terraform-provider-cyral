package datalabel

import (
	"fmt"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/src/cyral"
	cs "github.com/cyralinc/terraform-provider-cyral/src/cyral/datalabel/classificationrule"
	"github.com/cyralinc/terraform-provider-cyral/src/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	datalabelResourceName = "datalabel"
)

func initialDataLabelConfig() *DataLabel {
	return &DataLabel{
		Name:        utils.AccTestName(datalabelResourceName, "label1"),
		Description: "label1-description",
		Tags:        []string{"tag1", "tag2"},
		ClassificationRule: &cs.ClassificationRule{
			RuleType:   "UNKNOWN",
			RuleCode:   "",
			RuleStatus: "ENABLED",
		},
	}
}

func updatedDataLabelConfig() *DataLabel {
	return &DataLabel{
		Name:        utils.AccTestName(datalabelResourceName, "label2"),
		Description: "label2-description",
		Tags:        []string{"tag1", "tag2"},
		ClassificationRule: &cs.ClassificationRule{
			RuleType:   "REGO",
			RuleCode:   "int main() {cout << 'Hello World' << endl; return 0;}",
			RuleStatus: "DISABLED",
		},
	}
}

func TestAccDatalabelResource(t *testing.T) {
	testInitialConfig, testInitialFunc := setupDatalabelTest(t,
		"main_test", initialDataLabelConfig())
	testUpdatedConfig, testUpdatedFunc := setupDatalabelTest(t,
		"main_test", updatedDataLabelConfig())
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"cyral": func() (*schema.Provider, error) {
				return cyral.Provider(), nil
			},
		},
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
		resource.TestCheckResourceAttr(resourceFullName, "name", dataLabel.Name),
		resource.TestCheckResourceAttr(resourceFullName, "description", dataLabel.Description),
		resource.TestCheckResourceAttr(resourceFullName, "tags.#", "2"),
		resource.TestCheckResourceAttr(
			resourceFullName,
			"classification_rule.0.rule_type",
			dataLabel.ClassificationRule.RuleType,
		),
		resource.TestCheckResourceAttr(
			resourceFullName,
			"classification_rule.0.rule_code",
			dataLabel.ClassificationRule.RuleCode,
		),
		resource.TestCheckResourceAttr(
			resourceFullName,
			"classification_rule.0.rule_status",
			dataLabel.ClassificationRule.RuleStatus,
		),
	)

	return configuration, testFunction
}

func datalabelConfigResourceFullName(resName string) string {
	return fmt.Sprintf("cyral_datalabel.%s", resName)
}

func formatDataLabelIntoConfig(resName string, dataLabel *DataLabel) string {
	var classificationRuleConfig string
	if dataLabel.ClassificationRule != nil {
		classificationRuleConfig = fmt.Sprintf(`
 		classification_rule {
 			rule_type = "%s"
 			rule_code = "%s"
 			rule_status = "%s"
 		}`,
			dataLabel.ClassificationRule.RuleType,
			dataLabel.ClassificationRule.RuleCode,
			dataLabel.ClassificationRule.RuleStatus,
		)
	}
	return fmt.Sprintf(`
	resource "cyral_datalabel" "%s" {
		name  = "%s"
		description = "%s"
		tags = %s
		%s
	}`,
		resName,
		dataLabel.Name,
		dataLabel.Description,
		utils.ListToStr(dataLabel.Tags),
		classificationRuleConfig,
	)
}
