package cyral

import (
	"fmt"
)

func urlQuery(kv map[string]string) string {
	queryStr := "?"
	for k, v := range kv {
		queryStr += fmt.Sprintf("&%s=%s", k, v)
	}
	return queryStr
}

func supportedTypesMarkdown(types []string) string {
	var s string
	for _, typ := range types {
		s += fmt.Sprintf("\n  - `%s`", typ)
	}
	return s
}
