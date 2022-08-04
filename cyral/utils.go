package cyral

import (
	"fmt"
	"sort"
)

func urlQuery(kv map[string]string) string {
	queryStr := "?"
	for k, v := range kv {
		queryStr += fmt.Sprintf("&%s=%s", k, v)
	}
	return queryStr
}

func elementsMatch(this, other []string) bool {
	if len(this) != len(other) {
		return false
	}
	copyThis := append([]string{}, this...)
	copyOther := append([]string{}, other...)
	sort.Strings(copyThis)
	sort.Strings(copyOther)
	for i, elemThis := range copyThis {
		if elemOther := copyOther[i]; elemThis != elemOther {
			return false
		}
	}
	return true
}
