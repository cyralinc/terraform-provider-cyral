package cyral

import (
	"fmt"
	"strings"
)

func urlQuery(kv map[string]string) string {
	queryStr := "?"
	for k, v := range kv {
		queryStr += fmt.Sprintf("&%s=%s", k, v)
	}
	return queryStr
}

func marshalComposedID(id1, id2 string) string {
	return fmt.Sprintf("%s-%s", id1, id2)
}

func unmarshalComposedID(id string) (string, string) {
	ids := strings.Split(id, "-")
	return ids[0], ids[1]
}
