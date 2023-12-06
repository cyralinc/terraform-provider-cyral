package datalabel_test

import (
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/datalabel"
	cs "github.com/cyralinc/terraform-provider-cyral/cyral/internal/datalabel/classificationrule"
	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	datalabelResourceName = "datalabel"
)

func initialDataLabelConfig() *datalabel.DataLabel {
	return &datalabel.DataLabel{
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

func updatedDataLabelConfig() *datalabel.DataLabel {
	return &datalabel.DataLabel{
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
	// testUpdatedConfig, testUpdatedFunc := setupDatalabelTest(t,
	// 	"main_test", updatedDataLabelConfig())
	resource.Test(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testInitialConfig,
				Check:  testInitialFunc,
			},
			// {
			// 	Config: testUpdatedConfig,
			// 	Check:  testUpdatedFunc,
			// },
			// {
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	ResourceName:      "cyral_datalabel.main_test",
			// },
		},
	})
}

func setupDatalabelTest(t *testing.T, resName string, dataLabel *datalabel.DataLabel) (string, resource.TestCheckFunc) {
	ruleType, ruleCode, ruleStatus := "", "", ""
	if dataLabel.ClassificationRule != nil {
		ruleType = dataLabel.ClassificationRule.RuleType
		ruleCode = dataLabel.ClassificationRule.RuleCode
		ruleStatus = dataLabel.ClassificationRule.RuleStatus
	}
	config := utils.FormatDataLabelIntoConfig(resName, dataLabel.Name, dataLabel.Description,
		ruleType, ruleCode, ruleStatus, dataLabel.Tags)

	resourceFullName := utils.DatalabelConfigResourceFullName(resName)

	testFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "name", dataLabel.Name),
		// resource.TestCheckResourceAttr(resourceFullName, "description", dataLabel.Description),
		// resource.TestCheckResourceAttr(resourceFullName, "tags.#", "2"),
		// resource.TestCheckResourceAttr(
		// 	resourceFullName,
		// 	"classification_rule.0.rule_type",
		// 	dataLabel.ClassificationRule.RuleType,
		// ),
		// resource.TestCheckResourceAttr(
		// 	resourceFullName,
		// 	"classification_rule.0.rule_code",
		// 	dataLabel.ClassificationRule.RuleCode,
		// ),
		// resource.TestCheckResourceAttr(
		// 	resourceFullName,
		// 	"classification_rule.0.rule_status",
		// 	dataLabel.ClassificationRule.RuleStatus,
		// ),
	)

	return config, testFunction
}
