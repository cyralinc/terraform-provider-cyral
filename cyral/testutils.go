package cyral

import (
	"fmt"
)

// TODO: don't assume the list is not empty --aholmquist 2022-06-27
func formatAttributes(list []string) string {
	currentResp := fmt.Sprintf("\"%s\"", list[0])
	if len(list) > 1 {
		for _, item := range list[1:] {
			currentResp = fmt.Sprintf("%s, \"%s\"", currentResp, item)
		}
	}
	return currentResp
}
