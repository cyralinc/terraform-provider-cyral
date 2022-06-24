package cyral

import (
	"fmt"
)

func formatAttributes(list []string) string {
	currentResp := fmt.Sprintf("\"%s\"", list[0])
	if len(list) > 1 {
		for _, item := range list[1:] {
			currentResp = fmt.Sprintf("%s, \"%s\"", currentResp, item)
		}
	}
	return currentResp
}
