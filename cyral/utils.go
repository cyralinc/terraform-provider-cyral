package cyral

import (
	"fmt"
	"sort"
	"strings"
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

func supportedTypesMarkdown(types []string) string {
	var s string
	for _, typ := range types {
		s += fmt.Sprintf("\n  - `%s`", typ)
	}
	return s
}

func marshalComposedID(ids []string, sep string) string {
	return strings.Join(ids, sep)
}

func unmarshalComposedID(id, sep string, numFields int) ([]string, error) {
	ids := strings.Split(id, sep)
	if len(ids) < numFields {
		return nil, fmt.Errorf("unexpected ID syntax. Correct import ID " +
			fmt.Sprintf("syntax uses separator %q and contains %d "+
				"fields", sep, numFields))
	}
	return ids, nil
}
