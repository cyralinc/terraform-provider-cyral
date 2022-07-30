package cyral

import (
	"fmt"
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

func importErrorf(id, fmtstr string, args ...interface{}) error {
	return fmt.Errorf("for resource ID %q:"+fmtstr, []interface{}{id, args}...)
}
