package datalabel

import (
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
)

func DatalabelConfigResourceFullName(resName string) string {
	return fmt.Sprintf("cyral_datalabel.%s", resName)
}

func FormatDataLabelIntoConfig(resName string, dataLabel *DataLabel) string {
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
