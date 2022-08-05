package cyral

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func formatAttributes(attributes []string) string {
	if len(attributes) == 0 {
		return ""
	}
	s := fmt.Sprintf(`"%s"`, attributes[0])
	if len(attributes) > 1 {
		for _, attribute := range attributes[1:] {
			s += fmt.Sprintf(`, "%s"`, attribute)
		}
	}
	return s
}

func importStateComposedIDFunc(
	resName string,
	idAtts []string,
	sep string,
) func(*terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		res, ok := s.RootModule().Resources[resName]
		if !ok {
			return "", fmt.Errorf("Resource not found: %s", resName)
		}
		var idParts []string
		for _, idAtt := range idAtts {
			idParts = append(idParts, res.Primary.Attributes[idAtt])
		}
		return marshalComposedID(idParts, sep), nil
	}
}
